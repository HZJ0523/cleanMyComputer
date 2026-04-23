package scanner

import (
	"context"
	"os"
	"path/filepath"

	"github.com/hzj0523/cleanMyComputer/internal/models"
)

type Scanner struct {
	workers int
	filter  *Filter
}

func NewScanner(workers int) *Scanner {
	return &Scanner{workers: workers, filter: NewFilter()}
}

func (s *Scanner) ScanTargets(ctx context.Context, targets []models.Target) ([]*models.ScanItem, error) {
	var results []*models.ScanItem
	for _, target := range targets {
		matches, err := filepath.Glob(filepath.Join(target.Path, target.Pattern))
		if err != nil {
			return nil, err
		}
		for _, match := range matches {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
			}
			info, err := os.Stat(match)
			if err != nil || info.IsDir() {
				continue
			}
			results = append(results, &models.ScanItem{
				Path:    match,
				Size:    info.Size(),
				ModTime: info.ModTime(),
			})
		}
	}
	return results, nil
}
