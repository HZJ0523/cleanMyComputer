package scanner

import "path/filepath"

type Filter struct {
	excludePatterns []string
}

func NewFilter() *Filter {
	return &Filter{excludePatterns: []string{}}
}

func (f *Filter) AddExclude(pattern string) {
	f.excludePatterns = append(f.excludePatterns, pattern)
}

func (f *Filter) ShouldInclude(path string) bool {
	for _, pattern := range f.excludePatterns {
		matched, _ := filepath.Match(pattern, filepath.Base(path))
		if matched {
			return false
		}
	}
	return true
}
