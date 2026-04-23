package cleaner

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestExecutor_Execute(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	os.WriteFile(testFile, []byte("test"), 0644)

	qm, _ := NewQuarantineManager(filepath.Join(tmpDir, "quarantine"))
	executor := NewExecutor(qm)

	task := &CleanTask{
		Files: []*FileItem{
			{Path: testFile, Size: 4, RiskScore: 10},
		},
	}

	ctx := context.Background()
	result, err := executor.Execute(ctx, task)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if result.FreedSize != 4 {
		t.Errorf("FreedSize = %d, want 4", result.FreedSize)
	}
}

func TestExecutor_Execute_HighRisk(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	os.WriteFile(testFile, []byte("test"), 0644)

	qm, _ := NewQuarantineManager(filepath.Join(tmpDir, "quarantine"))
	executor := NewExecutor(qm)

	task := &CleanTask{
		Files: []*FileItem{
			{Path: testFile, Size: 4, RiskScore: 80},
		},
	}

	ctx := context.Background()
	result, err := executor.Execute(ctx, task)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if len(result.Cleaned) != 1 {
		t.Errorf("Expected 1 cleaned file, got %d", len(result.Cleaned))
	}
}
