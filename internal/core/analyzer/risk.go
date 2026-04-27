package analyzer

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hzj0523/cleanMyComputer/internal/models"
)

type RiskAnalyzer struct {
	systemPaths    map[string]bool
	protectedExts  map[string]bool
	forbiddenPaths []string
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
		forbiddenPaths: []string{
			"C:\\Windows\\System32\\config",
			"C:\\Windows\\System32\\drivers\\etc",
			"C:\\Windows\\Boot",
			"C:\\Windows\\EFI",
			"C:\\$Recycle.Bin",
			"C:\\Windows\\SysWOW64",
			"C:\\Program Files\\WindowsApps",
			"C:\\ProgramData\\Microsoft\\Windows\\Start Menu",
		},
	}
}

func (r *RiskAnalyzer) IsForbidden(path string) bool {
	lowerPath := strings.ToLower(path)
	for _, fp := range r.forbiddenPaths {
		lowerFP := strings.ToLower(fp)
		if lowerPath == lowerFP || strings.HasPrefix(lowerPath, lowerFP+string(os.PathSeparator)) {
			return true
		}
	}
	return false
}

func (r *RiskAnalyzer) IsPathSafe(path string) bool {
	if strings.Contains(path, "..") {
		return false
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		return false
	}
	return !r.IsForbidden(absPath)
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
	lowerPath := strings.ToLower(path)
	for sysPath := range r.systemPaths {
		lowerSys := strings.ToLower(sysPath)
		if lowerPath == lowerSys || strings.HasPrefix(lowerPath, lowerSys+string(os.PathSeparator)) {
			return true
		}
	}
	return false
}
