package scanner

import (
	"path/filepath"
	"sync"
)

type Filter struct {
	mu              sync.RWMutex
	excludePatterns []string
}

func NewFilter() *Filter {
	return &Filter{excludePatterns: []string{}}
}

func (f *Filter) AddExclude(pattern string) {
	f.mu.Lock()
	f.excludePatterns = append(f.excludePatterns, pattern)
	f.mu.Unlock()
}

func (f *Filter) ShouldInclude(path string) bool {
	f.mu.RLock()
	defer f.mu.RUnlock()
	for _, pattern := range f.excludePatterns {
		matched, _ := filepath.Match(pattern, filepath.Base(path))
		if matched {
			return false
		}
	}
	return true
}
