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

	// Step 1: Scan
	s := scanner.NewScanner(2)
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

	// Step 2: Clean
	qm := cleaner.NewQuarantineManager(filepath.Join(tmpDir, "quarantine"))
	executor := cleaner.NewExecutor(qm)

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
}

func TestHighRiskWorkflow(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "important.sys")
	os.WriteFile(testFile, []byte("important"), 0644)

	ctx := context.Background()

	// Clean with high risk score - should go to quarantine
	qm := cleaner.NewQuarantineManager(filepath.Join(tmpDir, "quarantine"))
	executor := cleaner.NewExecutor(qm)

	task := &cleaner.CleanTask{
		Files: []*cleaner.FileItem{
			{Path: testFile, Size: 9, RiskScore: 80},
		},
	}

	result, err := executor.Execute(ctx, task)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if len(result.Cleaned) != 1 {
		t.Errorf("Expected 1 cleaned, got %d", len(result.Cleaned))
	}

	// File should be in quarantine, not deleted
	if _, err := os.Stat(testFile); !os.IsNotExist(err) {
		t.Error("Expected original file to be moved")
	}
}
