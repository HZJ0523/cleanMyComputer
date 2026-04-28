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

	executor := NewExecutor()

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
	if _, err := os.Stat(testFile); !os.IsNotExist(err) {
		t.Error("expected file to be deleted")
	}
}
