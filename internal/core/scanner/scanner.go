package scanner

import (
	"context"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/hzj0523/cleanMyComputer/internal/models"
)

var winEnvRe = regexp.MustCompile(`%([^%]+)%`)

func expandPath(path string) string {
	expanded := os.ExpandEnv(path)
	expanded = winEnvRe.ReplaceAllStringFunc(expanded, func(match string) string {
		name := match[1 : len(match)-1]
		return os.Getenv(name)
	})
	return filepath.Clean(expanded)
}

func isOldEnough(modTime time.Time, daysOld int) bool {
	if daysOld <= 0 {
		return true
	}
	return time.Since(modTime) >= time.Duration(daysOld)*24*time.Hour
}

type Scanner struct {
	workers int
	filter  *Filter
}

func NewScanner(workers int) *Scanner {
	return &Scanner{workers: workers, filter: NewFilter()}
}

func (s *Scanner) ScanRule(ctx context.Context, rule *models.CleanRule) ([]*models.ScanItem, error) {
	var results []*models.ScanItem
	for _, target := range rule.Targets {
		select {
		case <-ctx.Done():
			return results, ctx.Err()
		default:
		}

		switch target.Type {
		case "folder":
			items, err := s.scanFolder(ctx, &target, rule.ID)
			if err != nil {
				continue
			}
			results = append(results, items...)
		case "command":
			if target.Path != "" {
				results = append(results, &models.ScanItem{
					Path:      target.Path,
					RuleID:    rule.ID,
					RiskScore: rule.RiskScore,
				})
			}
		}
	}
	return results, nil
}

func (s *Scanner) scanFolder(ctx context.Context, target *models.Target, ruleID string) ([]*models.ScanItem, error) {
	expandedPath := expandPath(target.Path)

	var results []*models.ScanItem

	if target.Recursive {
		maxDepth := target.MaxDepth
		if maxDepth <= 0 {
			maxDepth = 10
		}
		filepath.WalkDir(expandedPath, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return nil
			}
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			rel, err := filepath.Rel(expandedPath, path)
			if err != nil {
				return nil
			}
			depth := 0
			for _, c := range rel {
				if c == filepath.Separator {
					depth++
				}
			}
			if depth > maxDepth {
				if d.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}

			if d.IsDir() {
				return nil
			}

			if target.Pattern != "" && target.Pattern != "*" {
				matched, _ := filepath.Match(target.Pattern, d.Name())
				if !matched {
					return nil
				}
			}

			if !s.filter.ShouldInclude(path) {
				return nil
			}
			for _, exclude := range target.ExcludeList {
				matched, _ := filepath.Match(exclude, d.Name())
				if matched {
					return nil
				}
			}

			info, err := d.Info()
			if err != nil {
				return nil
			}

			if !isOldEnough(info.ModTime(), target.DaysOld) {
				return nil
			}

			results = append(results, &models.ScanItem{
				Path:    path,
				Size:    info.Size(),
				ModTime: info.ModTime(),
				RuleID:  ruleID,
			})
			return nil
		})
	} else {
		matches, err := filepath.Glob(filepath.Join(expandedPath, target.Pattern))
		if err != nil {
			return nil, err
		}
		for _, match := range matches {
			info, err := os.Stat(match)
			if err != nil || info.IsDir() {
				continue
			}
			if !s.filter.ShouldInclude(match) {
				continue
			}
			excluded := false
			for _, exclude := range target.ExcludeList {
				matched, _ := filepath.Match(exclude, filepath.Base(match))
				if matched {
					excluded = true
					break
				}
			}
			if excluded {
				continue
			}
			if !isOldEnough(info.ModTime(), target.DaysOld) {
				continue
			}
			results = append(results, &models.ScanItem{
				Path:    match,
				Size:    info.Size(),
				ModTime: info.ModTime(),
				RuleID:  ruleID,
			})
		}
	}

	return results, nil
}

func (s *Scanner) ScanTargets(ctx context.Context, targets []models.Target) ([]*models.ScanItem, error) {
	var results []*models.ScanItem
	for _, target := range targets {
		expandedPath := expandPath(target.Path)
		matches, err := filepath.Glob(filepath.Join(expandedPath, target.Pattern))
		if err != nil {
			continue
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
