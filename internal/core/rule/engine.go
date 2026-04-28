package rule

import (
	"context"
	"log"
	"sort"
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

	var valid []*models.CleanRule
	for _, rule := range rules {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		if err := rule.Validate(); err == nil {
			valid = append(valid, rule)
		} else {
			log.Printf("Warning: skipping invalid rule %s: %v", rule.ID, err)
		}
	}

	e.mu.Lock()
	for _, rule := range valid {
		e.rules[rule.ID] = rule
	}
	e.mu.Unlock()
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
	sort.Slice(enabled, func(i, j int) bool {
		return enabled[i].ID < enabled[j].ID
	})
	return enabled
}

func (e *Engine) GetAllRules() []*models.CleanRule {
	e.mu.RLock()
	defer e.mu.RUnlock()
	var all []*models.CleanRule
	for _, rule := range e.rules {
		all = append(all, rule)
	}
	sort.Slice(all, func(i, j int) bool {
		return all[i].ID < all[j].ID
	})
	return all
}
