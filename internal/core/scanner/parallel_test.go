package scanner

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/hzj0523/cleanMyComputer/internal/models"
)

func TestParallelScanner_Workers(t *testing.T) {
	scanner := NewParallelScanner(4)
	if scanner.Workers() != 4 {
		t.Errorf("Workers() = %d, want 4", scanner.Workers())
	}
}

func TestParallelScanner_MinWorkers(t *testing.T) {
	scanner := NewParallelScanner(0)
	if scanner.Workers() != 1 {
		t.Errorf("Workers() = %d, want 1 (minimum)", scanner.Workers())
	}
}

func TestParallelScanner_ScanRules(t *testing.T) {
	tmpDir1 := t.TempDir()
	tmpDir2 := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir1, "test1.txt"), []byte("test1"), 0644)
	os.WriteFile(filepath.Join(tmpDir2, "test2.txt"), []byte("test2"), 0644)

	rules := []*models.CleanRule{
		{
			ID: "rule1",
			Targets: []models.Target{
				{Type: "folder", Path: tmpDir1, Pattern: "*.txt", Recursive: false},
			},
		},
		{
			ID: "rule2",
			Targets: []models.Target{
				{Type: "folder", Path: tmpDir2, Pattern: "*.txt", Recursive: false},
			},
		},
	}

	scanner := NewParallelScanner(2)
	ctx := context.Background()
	results, err := scanner.ScanRules(ctx, rules)
	if err != nil {
		t.Fatalf("ScanRules() error = %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(results))
	}
}
