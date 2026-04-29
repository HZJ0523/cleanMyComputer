package integration

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/hzj0523/cleanMyComputer/internal/core/cleaner"
	"github.com/hzj0523/cleanMyComputer/internal/core/scanner"
	"github.com/hzj0523/cleanMyComputer/internal/models"
)

func TestFullCleanWorkflow(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	os.WriteFile(testFile, []byte("test"), 0644)

	ctx := context.Background()

	s := scanner.NewScanner()
	targets := []models.Target{
		{Type: "folder", Path: tmpDir, Pattern: "*.txt", Recursive: false},
	}

	results, err := s.ScanTargets(ctx, targets)
	if err != nil {
		t.Fatalf("ScanTargets() error = %v", err)
	}

	if len(results) == 0 {
		t.Fatal("Expected scan results")
	}

	executor := cleaner.NewExecutor()

	task := &cleaner.CleanTask{
		Files: []*cleaner.FileItem{
			{Path: testFile, Size: 4, RiskScore: 10},
		},
	}

	cleanResult, err := executor.Execute(ctx, task)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if cleanResult.FreedSize != 4 {
		t.Errorf("FreedSize = %d, want 4", cleanResult.FreedSize)
	}
	if _, err := os.Stat(testFile); !os.IsNotExist(err) {
		t.Error("expected file to be deleted")
	}
}
