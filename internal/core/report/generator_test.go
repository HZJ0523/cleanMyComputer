package report

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/hzj0523/cleanMyComputer/internal/models"
)

func TestGenerator_GenerateText(t *testing.T) {
	g := NewGenerator()
	r := &Report{
		ScanLevel: 1,
		StartTime: time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC),
		EndTime:   time.Date(2026, 1, 1, 10, 5, 0, 0, time.UTC),
		FreedSize: 1024 * 1024,
		Failed:    1,
		Items: []*models.ScanItem{
			{Path: "C:\\Users\\TestUser\\temp\\cache.tmp", Size: 512, RiskScore: 10},
			{Path: "C:\\Windows\\temp\\log.txt", Size: 2048, RiskScore: 50},
		},
	}

	text := g.GenerateText(r)
	if text == "" {
		t.Fatal("expected non-empty report text")
	}
	if len(r.Items) < 1 {
		t.Error("expected items in report")
	}
}

func TestGenerator_ExportToFile(t *testing.T) {
	g := NewGenerator()
	tmpDir := t.TempDir()
	exportPath := filepath.Join(tmpDir, "report.txt")

	r := &Report{
		ScanLevel: 2,
		StartTime: time.Now(),
		EndTime:   time.Now(),
		FreedSize: 1024,
		Items:     []*models.ScanItem{},
	}

	if err := g.ExportToFile(r, exportPath); err != nil {
		t.Fatalf("ExportToFile: %v", err)
	}

	data, err := os.ReadFile(exportPath)
	if err != nil {
		t.Fatalf("failed to read exported file: %v", err)
	}
	if len(data) == 0 {
		t.Error("expected non-empty exported file")
	}
}

func TestSanitizePath(t *testing.T) {
	username := os.Getenv("USERNAME")
	if username == "" {
		t.Skip("USERNAME env var not set")
	}

	input := "C:\\Users\\" + username + "\\temp\\cache.tmp"
	result := sanitizePath(input)
	if result == input {
		t.Error("expected username to be sanitized")
	}
}

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		input    int64
		contains string
	}{
		{0, "0 B"},
		{512, "512 B"},
		{1024, "1.0 KB"},
		{1024 * 1024, "1.0 MB"},
		{1024 * 1024 * 1024, "1.0 GB"},
	}
	for _, tt := range tests {
		result := formatBytes(tt.input)
		if result != tt.contains {
			t.Errorf("formatBytes(%d) = %q, want %q", tt.input, result, tt.contains)
		}
	}
}
