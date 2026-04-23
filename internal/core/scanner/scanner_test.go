package scanner

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/hzj0523/cleanMyComputer/internal/models"
)

func TestScanner_ScanTargets(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	os.WriteFile(testFile, []byte("test"), 0644)

	scanner := NewScanner(2)
	targets := []models.Target{
		{
			Type:      "folder",
			Path:      tmpDir,
			Pattern:   "*.txt",
			Recursive: false,
		},
	}

	ctx := context.Background()
	results, err := scanner.ScanTargets(ctx, targets)
	if err != nil {
		t.Fatalf("ScanTargets() error = %v", err)
	}

	if len(results) == 0 {
		t.Error("Expected scan results")
	}
}
