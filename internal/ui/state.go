package ui

import (
	"context"
	"sync"
	"time"

	"github.com/hzj0523/cleanMyComputer/internal/core/analyzer"
	"github.com/hzj0523/cleanMyComputer/internal/core/rule"
	"github.com/hzj0523/cleanMyComputer/internal/core/scanner"
	"github.com/hzj0523/cleanMyComputer/internal/models"
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

	// Core components
	Engine   *rule.Engine
	Scanner  *scanner.Scanner
	Analyzer *analyzer.RiskAnalyzer

	// State
	ScanItems  []*models.ScanItem
	Rules      []*models.CleanRule
	IsScanning bool

	// OnCleanComplete 在清理完成时被调用
	OnCleanComplete func(result CleanResult)
}

func NewAppState() *AppState {
	loader := rule.NewLoader()
	engine := rule.NewEngine(loader)

	return &AppState{
		Engine:   engine,
		Scanner:  scanner.NewScanner(4),
		Analyzer: analyzer.NewRiskAnalyzer(),
	}
}

func (s *AppState) RunScan(level int) error {
	s.mu.Lock()
	s.IsScanning = true
	s.ScanItems = nil
	s.mu.Unlock()

	ctx := context.Background()

	// Load rules
	if err := s.Engine.LoadRules(ctx, level); err != nil {
		return err
	}
	s.Rules = s.Engine.GetEnabledRules(level)

	// Scan each rule separately to track RuleID
	var allItems []*models.ScanItem
	for _, r := range s.Rules {
		items, err := s.Scanner.ScanRule(ctx, r)
		if err != nil {
			continue
		}
		allItems = append(allItems, items...)
	}

	// Calculate risk scores
	for _, item := range allItems {
		item.RiskScore = s.Analyzer.CalculateRisk(item)
	}

	s.mu.Lock()
	s.ScanItems = allItems
	s.IsScanning = false
	s.mu.Unlock()

	return nil
}
