package analyzer

import (
	"path/filepath"
	"strings"
	"time"

	"github.com/hzj0523/cleanMyComputer/internal/models"
)

type RiskAnalyzer struct {
	systemPaths   map[string]bool
	protectedExts map[string]bool
}

func NewRiskAnalyzer() *RiskAnalyzer {
	return &RiskAnalyzer{
		systemPaths: map[string]bool{
			"C:\\Windows\\System32": true,
			"C:\\Program Files":     true,
		},
		protectedExts: map[string]bool{
			".exe": true, ".dll": true, ".sys": true,
		},
	}
}

func (r *RiskAnalyzer) CalculateRisk(item *models.ScanItem) int {
	score := item.RiskScore
	if r.isSystemPath(item.Path) {
		score += 30
	}
	ext := strings.ToLower(filepath.Ext(item.Path))
	if r.protectedExts[ext] {
		score += 20
	}
	if item.Size > 100*1024*1024 {
		score += 15
	}
	if time.Since(item.ModTime) < 7*24*time.Hour {
		score += 10
	}
	if score > 100 {
		score = 100
	}
	return score
}

func (r *RiskAnalyzer) isSystemPath(path string) bool {
	for sysPath := range r.systemPaths {
		if strings.HasPrefix(path, sysPath) {
			return true
		}
	}
	return false
}
