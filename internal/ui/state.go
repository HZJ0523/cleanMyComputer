package ui

import (
	"context"
	"sync"

	"github.com/hzj0523/cleanMyComputer/internal/core/analyzer"
	"github.com/hzj0523/cleanMyComputer/internal/core/rule"
	"github.com/hzj0523/cleanMyComputer/internal/core/scanner"
	"github.com/hzj0523/cleanMyComputer/internal/models"
)

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

	// Collect all targets from all rules
	var allTargets []models.Target
	for _, r := range s.Rules {
		allTargets = append(allTargets, r.Targets...)
	}

	// Scan
	items, err := s.Scanner.ScanTargets(ctx, allTargets)
	if err != nil {
		return err
	}

	// Calculate risk scores
	for _, item := range items {
		item.RiskScore = s.Analyzer.CalculateRisk(item)
	}

	s.mu.Lock()
	s.ScanItems = items
	s.IsScanning = false
	s.mu.Unlock()

	return nil
}
