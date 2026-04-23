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

func TestScanner_ScanRule_Recursive(t *testing.T) {
	tmpDir := t.TempDir()
	subDir := filepath.Join(tmpDir, "sub")
	os.MkdirAll(subDir, 0755)

	os.WriteFile(filepath.Join(tmpDir, "root.txt"), []byte("root"), 0644)
	os.WriteFile(filepath.Join(subDir, "deep.txt"), []byte("deep"), 0644)

	rule := &models.CleanRule{
		ID: "test_recursive",
		Targets: []models.Target{
			{
				Type:      "folder",
				Path:      tmpDir,
				Pattern:   "*.txt",
				Recursive: true,
				MaxDepth:  3,
			},
		},
	}

	scanner := NewScanner(2)
	ctx := context.Background()
	results, err := scanner.ScanRule(ctx, rule)
	if err != nil {
		t.Fatalf("ScanRule() error = %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 results (root + sub), got %d", len(results))
	}

	// Verify RuleID is set
	for _, item := range results {
		if item.RuleID != "test_recursive" {
			t.Errorf("Expected RuleID 'test_recursive', got '%s'", item.RuleID)
		}
	}
}

func TestScanner_ScanRule_EnvExpand(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "test.txt"), []byte("test"), 0644)

	// Set a custom env var for testing
	os.Setenv("TEST_SCAN_DIR", tmpDir)
	defer os.Unsetenv("TEST_SCAN_DIR")

	rule := &models.CleanRule{
		ID: "test_env",
		Targets: []models.Target{
			{
				Type:      "folder",
				Path:      "%TEST_SCAN_DIR%",
				Pattern:   "*.txt",
				Recursive: false,
			},
		},
	}

	scanner := NewScanner(2)
	ctx := context.Background()
	results, err := scanner.ScanRule(ctx, rule)
	if err != nil {
		t.Fatalf("ScanRule() error = %v", err)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 result, got %d", len(results))
	}
}

func TestScanner_ScanRule_ExcludeFilter(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "keep.log"), []byte("log"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "clean.txt"), []byte("txt"), 0644)

	rule := &models.CleanRule{
		ID: "test_exclude",
		Targets: []models.Target{
			{
				Type:        "folder",
				Path:        tmpDir,
				Pattern:     "*",
				Recursive:   false,
				ExcludeList: []string{"*.log"},
			},
		},
	}

	scanner := NewScanner(2)
	ctx := context.Background()
	results, err := scanner.ScanRule(ctx, rule)
	if err != nil {
		t.Fatalf("ScanRule() error = %v", err)
	}

	// Should only find clean.txt, not keep.log
	for _, item := range results {
		if filepath.Base(item.Path) == "keep.log" {
			t.Error("Expected keep.log to be excluded")
		}
	}
}
