package scanner

import (
	"context"
	"os"
	"path/filepath"
	"regexp"

	"github.com/hzj0523/cleanMyComputer/internal/models"
)

// expandPath expands environment variables in both $VAR/${VAR} and Windows %VAR% forms.
var winEnvRe = regexp.MustCompile(`%([^%]+)%`)

func expandPath(path string) string {
	// First expand $VAR / ${VAR} (os.ExpandEnv handles these)
	expanded := os.ExpandEnv(path)
	// Then expand Windows %VAR% syntax
	expanded = winEnvRe.ReplaceAllStringFunc(expanded, func(match string) string {
		name := match[1 : len(match)-1]
		return os.Getenv(name)
	})
	return expanded
}

type Scanner struct {
	workers int
	filter  *Filter
}

func NewScanner(workers int) *Scanner {
	return &Scanner{workers: workers, filter: NewFilter()}
}

// ScanRule scans all targets of a single rule, returning ScanItems with RuleID set.
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
			// command targets are handled by the platform layer
			// TODO: delegate to platform adapter for command-type cleanup
		}
	}
	return results, nil
}

func (s *Scanner) scanFolder(ctx context.Context, target *models.Target, ruleID string) ([]*models.ScanItem, error) {
	// BUG-001 fix: expand environment variables including %VAR% syntax
	expandedPath := expandPath(target.Path)

	var results []*models.ScanItem

	if target.Recursive {
		// BUG-002 fix: recursive scanning with depth control
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

			// Compute depth by counting separators in relative path
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

			// Pattern matching
			if target.Pattern != "" && target.Pattern != "*" {
				matched, _ := filepath.Match(target.Pattern, d.Name())
				if !matched {
					return nil
				}
			}

			// Filter fix: check global exclude list
			if !s.filter.ShouldInclude(path) {
				return nil
			}
			// Check target-level exclude list
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

			results = append(results, &models.ScanItem{
				Path:    path,
				Size:    info.Size(),
				ModTime: info.ModTime(),
				RuleID:  ruleID,
			})
			return nil
		})
	} else {
		// Non-recursive: single-level scan
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
			// Check target-level exclude list
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

// ScanTargets retains backward compatibility (prefer ScanRule).
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
