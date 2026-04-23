package scanner

import (
	"context"
	"sync"

	"github.com/hzj0523/cleanMyComputer/internal/models"
)

type ParallelScanner struct {
	workers int
	scanner *Scanner
}

func NewParallelScanner(workers int) *ParallelScanner {
	if workers < 1 {
		workers = 1
	}
	return &ParallelScanner{
		workers: workers,
		scanner: NewScanner(workers),
	}
}

func (p *ParallelScanner) Workers() int {
	return p.workers
}

// ScanRules 并行扫描多条规则
func (p *ParallelScanner) ScanRules(ctx context.Context, rules []*models.CleanRule) ([]*models.ScanItem, error) {
	type ruleResult struct {
		items []*models.ScanItem
		err   error
	}

	resultChan := make(chan ruleResult, len(rules))

	// 使用 semaphore 控制并发数
	sem := make(chan struct{}, p.workers)

	var wg sync.WaitGroup
	for _, rule := range rules {
		wg.Add(1)
		go func(r *models.CleanRule) {
			defer wg.Done()

			select {
			case <-ctx.Done():
				resultChan <- ruleResult{err: ctx.Err()}
				return
			case sem <- struct{}{}:
				defer func() { <-sem }()
			}

			items, err := p.scanner.ScanRule(ctx, r)
			resultChan <- ruleResult{items: items, err: err}
		}(rule)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	var allItems []*models.ScanItem
	var firstErr error
	for result := range resultChan {
		if result.err != nil && firstErr == nil {
			firstErr = result.err
		}
		allItems = append(allItems, result.items...)
	}

	return allItems, firstErr
}
