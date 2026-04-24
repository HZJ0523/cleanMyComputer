package app

import (
	"log"
	"sync"
	"time"
)

type Scheduler struct {
	mu        sync.Mutex
	orch      *Orchestrator
	stopCh    chan struct{}
	running   bool
	interval  time.Duration
	scanLevel int
}

func NewScheduler(orch *Orchestrator) *Scheduler {
	return &Scheduler{
		orch:      orch,
		scanLevel: 1,
		interval:  24 * time.Hour,
	}
}

func (s *Scheduler) SetInterval(d time.Duration) {
	s.mu.Lock()
	s.interval = d
	s.mu.Unlock()
}

func (s *Scheduler) SetScanLevel(level int) {
	s.mu.Lock()
	s.scanLevel = level
	s.mu.Unlock()
}

func (s *Scheduler) Start() {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return
	}
	s.running = true
	s.stopCh = make(chan struct{})
	interval := s.interval
	s.mu.Unlock()

	go s.run(interval)
	log.Printf("[Scheduler] started with interval %v", interval)
}

func (s *Scheduler) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.running {
		return
	}
	close(s.stopCh)
	s.running = false
	log.Printf("[Scheduler] stopped")
}

func (s *Scheduler) Running() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.running
}

func (s *Scheduler) run(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-s.stopCh:
			return
		case <-ticker.C:
			s.mu.Lock()
			level := s.scanLevel
			s.mu.Unlock()

			log.Printf("[Scheduler] auto-clean triggered at %v", time.Now().Format(time.DateTime))
			if err := s.orch.RunScan(level); err != nil {
				log.Printf("[Scheduler] scan failed: %v", err)
				continue
			}

			count := s.orch.GetScanItemCount()
			if count == 0 {
				log.Printf("[Scheduler] no items to clean")
				continue
			}
			log.Printf("[Scheduler] scan found %d items, starting clean", count)

			summary, err := s.orch.RunClean()
			if err != nil {
				log.Printf("[Scheduler] clean failed: %v", err)
				continue
			}
			log.Printf("[Scheduler] auto-clean completed: cleaned=%d, failed=%d, freed=%d bytes",
				summary.Cleaned, summary.Failed, summary.FreedSize)
		}
	}
}
