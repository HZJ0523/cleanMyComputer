package rule

import (
	"context"
	"sync"

	"github.com/hzj0523/cleanMyComputer/internal/models"
)

type Engine struct {
	rules  map[string]*models.CleanRule
	mu     sync.RWMutex
	loader *Loader
}

func NewEngine(loader *Loader) *Engine {
	return &Engine{
		rules:  make(map[string]*models.CleanRule),
		loader: loader,
	}
}

func (e *Engine) LoadRules(ctx context.Context, level int) error {
	rules, err := e.loader.LoadByLevel(level)
	if err != nil {
		return err
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	for _, rule := range rules {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		if err := rule.Validate(); err == nil {
			e.rules[rule.ID] = rule
		}
	}
	return nil
}

func (e *Engine) GetEnabledRules(level int) []*models.CleanRule {
	e.mu.RLock()
	defer e.mu.RUnlock()
	var enabled []*models.CleanRule
	for _, rule := range e.rules {
		if rule.Enabled && rule.Level <= level {
			enabled = append(enabled, rule)
		}
	}
	return enabled
}
