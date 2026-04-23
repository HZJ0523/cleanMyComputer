package ui

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/hzj0523/cleanMyComputer/internal/core/analyzer"
	"github.com/hzj0523/cleanMyComputer/internal/core/rule"
	"github.com/hzj0523/cleanMyComputer/internal/core/scanner"
	"github.com/hzj0523/cleanMyComputer/internal/models"
	"github.com/hzj0523/cleanMyComputer/internal/storage"
)

// CleanResult 记录一次清理的结果，用于回调通知
type CleanResult struct {
	Cleaned   int
	Failed    int
	FreedSize int64
	Duration  time.Duration
}

type AppState struct {
	mu sync.Mutex

	Engine          *rule.Engine
	ParallelScanner *scanner.ParallelScanner
	Analyzer        *analyzer.RiskAnalyzer
	DB              *storage.DB
	History         *storage.History

	ScanItems  []*models.ScanItem
	Rules      []*models.CleanRule
	IsScanning bool

	OnScanProgress  func(current, total int)
	OnCleanComplete func(result CleanResult)
}

func NewAppState() *AppState {
	loader := rule.NewLoader()
	engine := rule.NewEngine(loader)
	ps := scanner.NewParallelScanner(4)

	return &AppState{
		Engine:          engine,
		ParallelScanner: ps,
		Analyzer:        analyzer.NewRiskAnalyzer(),
	}
}

func (s *AppState) InitDB(path string) error {
	db, err := storage.NewDB(path)
	if err != nil {
		return err
	}
	s.DB = db
	s.History = storage.NewHistory(db)
	return nil
}

func (s *AppState) RunScan(level int) error {
	s.mu.Lock()
	s.IsScanning = true
	s.ScanItems = nil
	s.mu.Unlock()

	ctx := context.Background()

	if err := s.Engine.LoadRules(ctx, level); err != nil {
		return err
	}
	s.Rules = s.Engine.GetEnabledRules(level)

	// Parallel scan
	allItems, err := s.ParallelScanner.ScanRules(ctx, s.Rules)
	if err != nil {
		return err
	}

	// Calculate risk and filter forbidden
	var validItems []*models.ScanItem
	for _, item := range allItems {
		item.RiskScore = s.Analyzer.CalculateRisk(item)
		if !s.Analyzer.IsForbidden(item.Path) {
			validItems = append(validItems, item)
		}
	}

	s.mu.Lock()
	s.ScanItems = validItems
	s.IsScanning = false
	s.mu.Unlock()

	return nil
}

func (s *AppState) SaveCleanHistory(result CleanResult) {
	if s.History == nil {
		return
	}
	now := time.Now()
	record := &models.CleanRecord{
		StartTime:  now.Add(-result.Duration),
		EndTime:    now,
		ScanLevel:  1,
		TotalFiles: result.Cleaned + result.Failed,
		TotalSize:  result.FreedSize,
		FreedSize:  result.FreedSize,
		Status:     "success",
	}
	if result.Failed > 0 {
		record.Status = "partial"
	}
	_, err := s.History.Save(record)
	if err != nil {
		log.Printf("Failed to save clean history: %v", err)
	}
}

func (s *AppState) GetHistory() ([]*models.CleanRecord, error) {
	if s.History == nil {
		return nil, nil
	}
	return s.History.GetAll()
}
