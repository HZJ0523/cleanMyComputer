package ui

import (
	"time"

	"github.com/hzj0523/cleanMyComputer/internal/app"
)

type scheduler struct {
	*app.Scheduler
}

func newScheduler(orch *app.Orchestrator) *scheduler {
	return &scheduler{Scheduler: app.NewScheduler(orch)}
}

func (s *scheduler) SetInterval(hours int) {
	s.Scheduler.SetInterval(time.Duration(hours) * time.Hour)
}
