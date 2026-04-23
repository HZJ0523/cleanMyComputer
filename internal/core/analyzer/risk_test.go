package analyzer

import (
	"testing"
	"time"

	"github.com/hzj0523/cleanMyComputer/internal/models"
)

func TestRiskAnalyzer_CalculateRisk(t *testing.T) {
	analyzer := NewRiskAnalyzer()

	item := &models.ScanItem{
		Path:      "C:\\Windows\\System32\\test.dll",
		Size:      1024,
		ModTime:   time.Now(),
		RiskScore: 10,
	}

	score := analyzer.CalculateRisk(item)
	if score < 10 {
		t.Errorf("Expected risk score >= 10, got %d", score)
	}
}

func TestRiskAnalyzer_LowRiskFile(t *testing.T) {
	analyzer := NewRiskAnalyzer()

	item := &models.ScanItem{
		Path:      "C:\\Users\\test\\AppData\\Local\\Temp\\tmp.txt",
		Size:      100,
		ModTime:   time.Now().Add(-30 * 24 * time.Hour),
		RiskScore: 5,
	}

	score := analyzer.CalculateRisk(item)
	if score > 60 {
		t.Errorf("Expected low risk score, got %d", score)
	}
}
