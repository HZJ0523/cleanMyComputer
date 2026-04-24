package ui

import (
	"github.com/hzj0523/cleanMyComputer/internal/app"
)

type CleanResult = app.CleanSummary

type AppState struct {
	*app.Orchestrator
}

func NewAppState() *AppState {
	return &AppState{
		Orchestrator: app.NewOrchestrator(),
	}
}
