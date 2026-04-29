package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/hzj0523/cleanMyComputer/internal/core/analyzer"
	"github.com/hzj0523/cleanMyComputer/internal/core/cleaner"
	"github.com/hzj0523/cleanMyComputer/internal/core/rule"
	"github.com/hzj0523/cleanMyComputer/internal/core/scanner"
	"github.com/hzj0523/cleanMyComputer/internal/models"
	"github.com/hzj0523/cleanMyComputer/internal/storage"
)

var ErrScanInProgress = errors.New("scan already in progress")

type CleanSummary struct {
	Cleaned   int
	Failed    int
	FreedSize int64
	Duration  time.Duration
}

type Orchestrator struct {
	mu sync.Mutex

	engine   *rule.Engine
	scanner  *scanner.ParallelScanner
	analyzer *analyzer.RiskAnalyzer
	db       *storage.DB
	history  *storage.History
	config   *storage.Config

	ScanItems  []*models.ScanItem
	IsScanning bool
	scanLevel  int
}

func NewOrchestrator() *Orchestrator {
	loader := rule.NewLoader()
	return &Orchestrator{
		engine:   rule.NewEngine(loader),
		scanner:  scanner.NewParallelScanner(4),
		analyzer: analyzer.NewRiskAnalyzer(),
	}
}

func (o *Orchestrator) InitDB(path string) error {
	db, err := storage.NewDB(path)
	if err != nil {
		return err
	}
	o.db = db
	o.history = storage.NewHistory(db)
	o.config = storage.NewConfig(db)
	return nil
}

func (o *Orchestrator) CloseDB() {
	if o.db != nil {
		o.db.Close()
	}
}

func (o *Orchestrator) GetConfig(key string) (string, error) {
	if o.config == nil {
		return "", fmt.Errorf("database not initialized")
	}
	return o.config.Get(key)
}

func (o *Orchestrator) SetConfig(key, value string) error {
	if o.config == nil {
		return fmt.Errorf("database not initialized")
	}
	return o.config.Set(key, value)
}

func (o *Orchestrator) GetRuleStatus() (map[string]bool, error) {
	if o.db == nil {
		return nil, fmt.Errorf("database not initialized")
	}
	rows, err := o.db.Conn().Query("SELECT rule_id, enabled FROM rule_status")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	status := make(map[string]bool)
	for rows.Next() {
		var ruleID string
		var enabled bool
		if err := rows.Scan(&ruleID, &enabled); err != nil {
			log.Printf("Warning: skipping corrupt rule_status row: %v", err)
			continue
		}
		status[ruleID] = enabled
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed iterating rule_status rows: %w", err)
	}
	return status, nil
}

func (o *Orchestrator) SetRuleEnabled(ruleID string, enabled bool) error {
	if o.db == nil {
		return fmt.Errorf("database not initialized")
	}
	_, err := o.db.Conn().Exec(
		"INSERT INTO rule_status (rule_id, enabled, updated_at) VALUES (?, ?, CURRENT_TIMESTAMP) ON CONFLICT(rule_id) DO UPDATE SET enabled = ?, updated_at = CURRENT_TIMESTAMP",
		ruleID, enabled, enabled)
	return err
}

func (o *Orchestrator) GetAllRules() []*models.CleanRule {
	ctx := context.Background()
	o.engine.LoadRules(ctx, 3)

	rules := o.engine.GetAllRules()

	if o.db != nil {
		status, err := o.GetRuleStatus()
		if err != nil {
			log.Printf("Warning: failed to load rule status for display: %v", err)
			return rules
		}
		for _, rule := range rules {
			if enabled, ok := status[rule.ID]; ok {
				rule.Enabled = enabled
			}
		}
	}

	return rules
}

func (o *Orchestrator) GetScanItemsSafe() []*models.ScanItem {
	o.mu.Lock()
	defer o.mu.Unlock()
	items := make([]*models.ScanItem, len(o.ScanItems))
	copy(items, o.ScanItems)
	return items
}

func (o *Orchestrator) ClearScanItems() {
	o.mu.Lock()
	o.ScanItems = nil
	o.mu.Unlock()
}

func (o *Orchestrator) SetScanItemsForClean(items []*models.ScanItem) {
	o.mu.Lock()
	o.ScanItems = items
	o.mu.Unlock()
}

func (o *Orchestrator) GetScanItemCount() int {
	o.mu.Lock()
	defer o.mu.Unlock()
	return len(o.ScanItems)
}

func (o *Orchestrator) RunScan(level int) error {
	o.mu.Lock()
	if o.IsScanning {
		o.mu.Unlock()
		return ErrScanInProgress
	}
	o.IsScanning = true
	o.ScanItems = nil
	o.scanLevel = level
	o.mu.Unlock()

	defer func() {
		o.mu.Lock()
		o.IsScanning = false
		o.mu.Unlock()
	}()

	ctx := context.Background()

	if err := o.engine.LoadRules(ctx, level); err != nil {
		return err
	}

	rules := o.engine.GetEnabledRules(level)
	if o.db != nil {
		status, err := o.GetRuleStatus()
		if err != nil {
			log.Printf("Warning: failed to load rule status: %v", err)
		}
		var filtered []*models.CleanRule
		for _, r := range rules {
			if enabled, ok := status[r.ID]; ok {
				if !enabled {
					continue
				}
			}
			filtered = append(filtered, r)
		}
		rules = filtered
	}

	allItems, err := o.scanner.ScanRules(ctx, rules)
	if err != nil {
		return err
	}

	var validItems []*models.ScanItem
	for _, item := range allItems {
		item.RiskScore = o.analyzer.CalculateRisk(item)
		if o.analyzer.IsPathSafe(item.Path) {
			validItems = append(validItems, item)
		}
	}

	o.mu.Lock()
	o.ScanItems = validItems
	o.mu.Unlock()

	return nil
}

func (o *Orchestrator) RunClean() (CleanSummary, error) {
	o.mu.Lock()
	items := o.ScanItems
	level := o.scanLevel
	o.mu.Unlock()

	if len(items) == 0 {
		return CleanSummary{}, nil
	}

	var files []*cleaner.FileItem
	for _, item := range items {
		files = append(files, &cleaner.FileItem{
			Path:      item.Path,
			Size:      item.Size,
			RiskScore: item.RiskScore,
			Type:      item.Type,
		})
	}

	startTime := time.Now()
	task := &cleaner.CleanTask{Files: files}
	executor := cleaner.NewExecutor()

	result, err := executor.Execute(context.Background(), task)
	if err != nil {
		return CleanSummary{}, err
	}

	summary := CleanSummary{
		Cleaned:   len(result.Cleaned),
		Failed:    len(result.Failed),
		FreedSize: result.FreedSize,
		Duration:  time.Since(startTime),
	}

	o.SaveCleanHistory(summary, level)
	o.ClearScanItems()
	return summary, nil
}

func (o *Orchestrator) SaveCleanHistory(result CleanSummary, scanLevel int) {
	if o.history == nil {
		return
	}
	now := time.Now()
	record := &models.CleanRecord{
		StartTime:   now.Add(-result.Duration),
		EndTime:     now,
		ScanLevel:   scanLevel,
		TotalFiles:  result.Cleaned + result.Failed,
		TotalSize:   result.FreedSize,
		FreedSize:   result.FreedSize,
		FailedCount: result.Failed,
		Status:      "success",
	}
	if result.Failed > 0 {
		if result.Cleaned == 0 {
			record.Status = "failed"
		} else {
			record.Status = "partial"
		}
	}
	if _, err := o.history.Save(record); err != nil {
		log.Printf("Failed to save clean history: %v", err)
	}
}

func (o *Orchestrator) GetHistory() ([]*models.CleanRecord, error) {
	if o.history == nil {
		return nil, fmt.Errorf("database not initialized")
	}
	return o.history.GetAll()
}

func (o *Orchestrator) FindDuplicateFiles(root string) ([]analyzer.DuplicateGroup, error) {
	finder := analyzer.NewDuplicateFinder(1024)
	return finder.FindDuplicates(root)
}

func (o *Orchestrator) FindLargeFiles(root string, threshold int64) ([]analyzer.LargeFile, error) {
	finder := analyzer.NewLargeFileFinder(threshold)
	return finder.FindLargeFiles(root)
}
