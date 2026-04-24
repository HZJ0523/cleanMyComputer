package ui

import (
	"github.com/hzj0523/cleanMyComputer/internal/app"
	"github.com/hzj0523/cleanMyComputer/internal/models"
)

type CleanResult = app.CleanSummary

type AppState struct {
	*app.Orchestrator
	OnScanProgress func(current, total int)
}

func NewAppState() *AppState {
	return &AppState{
		Orchestrator: app.NewOrchestrator(),
	}
}

func (s *AppState) GetScanItems() []*models.ScanItem {
	return s.Orchestrator.ScanItems
}
