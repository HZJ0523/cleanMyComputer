package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/hzj0523/cleanMyComputer/internal/core/analyzer"
	"github.com/hzj0523/cleanMyComputer/internal/core/cleaner"
	"github.com/hzj0523/cleanMyComputer/internal/core/rule"
	"github.com/hzj0523/cleanMyComputer/internal/core/scanner"
	"github.com/hzj0523/cleanMyComputer/internal/models"
	"github.com/hzj0523/cleanMyComputer/internal/storage"
)

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
			continue
		}
		status[ruleID] = enabled
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
	return o.engine.GetEnabledRules(3)
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

func (o *Orchestrator) GetScanItemCount() int {
	o.mu.Lock()
	defer o.mu.Unlock()
	return len(o.ScanItems)
}

func (o *Orchestrator) RunScan(level int) error {
	o.mu.Lock()
	if o.IsScanning {
		o.mu.Unlock()
		return fmt.Errorf("扫描正在进行中，请稍候")
	}
	o.IsScanning = true
	o.ScanItems = nil
	o.scanLevel = level
	o.mu.Unlock()

	ctx := context.Background()

	if err := o.engine.LoadRules(ctx, level); err != nil {
		return err
	}

	// Apply user rule preferences
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
	o.IsScanning = false
	o.mu.Unlock()

	return nil
}

func (o *Orchestrator) RunClean() (CleanSummary, error) {
	o.mu.Lock()
	items := o.ScanItems
	o.mu.Unlock()

	if len(items) == 0 {
		return CleanSummary{}, nil
	}

	localAppData := os.Getenv("LOCALAPPDATA")
	if localAppData == "" {
		localAppData = os.TempDir()
	}
	qDir := filepath.Join(localAppData, "CleanMyComputer", "quarantine")
	qm, err := cleaner.NewQuarantineManager(qDir)
	if err != nil {
		return CleanSummary{}, fmt.Errorf("failed to create quarantine manager: %w", err)
	}

	var files []*cleaner.FileItem
	for _, item := range items {
		files = append(files, &cleaner.FileItem{
			Path:      item.Path,
			Size:      item.Size,
			RiskScore: item.RiskScore,
		})
	}

	startTime := time.Now()
	task := &cleaner.CleanTask{Files: files, TotalSize: 0}
	executor := cleaner.NewExecutor(qm)

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

	o.SaveCleanHistory(summary)
	o.ClearScanItems()
	return summary, nil
}

func (o *Orchestrator) SaveCleanHistory(result CleanSummary) {
	if o.history == nil {
		return
	}
	now := time.Now()
	record := &models.CleanRecord{
		StartTime:  now.Add(-result.Duration),
		EndTime:    now,
		ScanLevel:  o.scanLevel,
		TotalFiles: result.Cleaned + result.Failed,
		TotalSize:  result.FreedSize,
		FreedSize:  result.FreedSize,
		Status:     "success",
	}
	if result.Failed > 0 {
		record.Status = "partial"
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

func (o *Orchestrator) SaveQuarantineRecord(r cleaner.QuarantineRecord) error {
	if o.db == nil {
		return fmt.Errorf("database not initialized")
	}
	id := fmt.Sprintf("%d", time.Now().UnixNano())
	_, err := o.db.Conn().Exec(
		"INSERT INTO quarantine (id, original_path, quarantine_path, size_bytes, risk_score, created_at, expires_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
		id, r.OriginalPath, r.QuarantinePath, r.Size, 0, r.CreatedAt, r.ExpiresAt)
	return err
}

func (o *Orchestrator) GetQuarantinedItems() ([]cleaner.QuarantineRecord, error) {
	if o.db == nil {
		return nil, fmt.Errorf("database not initialized")
	}
	rows, err := o.db.Conn().Query("SELECT original_path, quarantine_path, size_bytes, created_at, expires_at FROM quarantine WHERE expires_at > datetime('now') ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []cleaner.QuarantineRecord
	for rows.Next() {
		var r cleaner.QuarantineRecord
		if err := rows.Scan(&r.OriginalPath, &r.QuarantinePath, &r.Size, &r.CreatedAt, &r.ExpiresAt); err != nil {
			continue
		}
		items = append(items, r)
	}
	return items, nil
}

func (o *Orchestrator) RestoreQuarantinedItem(quarantinePath, originalPath string) error {
	if err := cleaner.NewRecovery(nil).RestoreFile(quarantinePath, originalPath); err != nil {
		// Recovery uses os.Rename which works without QuarantineManager
		return err
	}
	_, err := o.db.Conn().Exec("DELETE FROM quarantine WHERE quarantine_path = ?", quarantinePath)
	return err
}

func (o *Orchestrator) DeleteQuarantinedItem(quarantinePath string) error {
	if err := os.Remove(quarantinePath); err != nil && !os.IsNotExist(err) {
		return err
	}
	_, err := o.db.Conn().Exec("DELETE FROM quarantine WHERE quarantine_path = ?", quarantinePath)
	return err
}

func (o *Orchestrator) CleanupExpiredQuarantine() error {
	if o.db == nil {
		return nil
	}
	rows, err := o.db.Conn().Query("SELECT quarantine_path FROM quarantine WHERE expires_at <= datetime('now')")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var path string
		if err := rows.Scan(&path); err != nil {
			continue
		}
		if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
				log.Printf("Warning: failed to remove expired quarantine file %s: %v", path, err)
			}
	}
	_, err = o.db.Conn().Exec("DELETE FROM quarantine WHERE expires_at <= datetime('now')")
	return err
}

func (o *Orchestrator) FindDuplicateFiles(root string) ([]analyzer.DuplicateGroup, error) {
	finder := analyzer.NewDuplicateFinder(1024)
	return finder.FindDuplicates(root)
}

func (o *Orchestrator) FindLargeFiles(root string, threshold int64) ([]analyzer.LargeFile, error) {
	finder := analyzer.NewLargeFileFinder(threshold)
	return finder.FindLargeFiles(root)
}
