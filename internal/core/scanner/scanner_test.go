package scanner

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/hzj0523/cleanMyComputer/internal/models"
)

func TestScanner_ScanRule_DirectoryPattern(t *testing.T) {
	tmpDir := t.TempDir()
	nodeModules := filepath.Join(tmpDir, "project", "node_modules")
	os.MkdirAll(nodeModules, 0755)
	os.WriteFile(filepath.Join(nodeModules, "a.js"), []byte("code"), 0644)
	os.WriteFile(filepath.Join(nodeModules, "b.js"), []byte("code"), 0644)

	otherDir := filepath.Join(tmpDir, "project", "src")
	os.MkdirAll(otherDir, 0755)
	os.WriteFile(filepath.Join(otherDir, "main.js"), []byte("main"), 0644)

	rule := &models.CleanRule{
		ID: "test_dir_pattern",
		Targets: []models.Target{
			{
				Type:      "folder",
				Path:      tmpDir,
				Pattern:   "node_modules",
				Recursive: true,
				MaxDepth:  5,
			},
		},
	}

	scanner := NewScanner()
	ctx := context.Background()
	results, err := scanner.ScanRule(ctx, rule)
	if err != nil {
		t.Fatalf("ScanRule() error = %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("Expected 1 result (node_modules dir), got %d", len(results))
	}

	item := results[0]
	if item.Type != "directory" {
		t.Errorf("Expected Type='directory', got '%s'", item.Type)
	}
	if item.Size != 8 {
		t.Errorf("Expected Size=8 (a.js=4 + b.js=4), got %d", item.Size)
	}
	if filepath.Base(item.Path) != "node_modules" {
		t.Errorf("Expected path base 'node_modules', got '%s'", filepath.Base(item.Path))
	}
}

func TestScanner_ScanRule_DirectoryPatternExcluded(t *testing.T) {
	tmpDir := t.TempDir()
	buildDir := filepath.Join(tmpDir, "myapp", "build")
	os.MkdirAll(buildDir, 0755)
	os.WriteFile(filepath.Join(buildDir, "out.js"), []byte("out"), 0644)

	rule := &models.CleanRule{
		ID: "test_dir_exclude",
		Targets: []models.Target{
			{
				Type:        "folder",
				Path:        tmpDir,
				Pattern:     "build",
				Recursive:   true,
				MaxDepth:    5,
				ExcludeList: []string{"build"},
			},
		},
	}

	scanner := NewScanner()
	ctx := context.Background()
	results, err := scanner.ScanRule(ctx, rule)
	if err != nil {
		t.Fatalf("ScanRule() error = %v", err)
	}

	if len(results) != 0 {
		t.Errorf("Expected 0 results (excluded), got %d", len(results))
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

	scanner := NewScanner()
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

	scanner := NewScanner()
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

	scanner := NewScanner()
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
