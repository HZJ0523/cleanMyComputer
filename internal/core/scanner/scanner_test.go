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

func TestScanner_ScanRule_BrowserCacheWildcard(t *testing.T) {
	tmpDir := t.TempDir()
	userDataDir := filepath.Join(tmpDir, "User Data")
	profileDir := filepath.Join(userDataDir, "Default")
	os.MkdirAll(profileDir, 0755)

	cacheDirs := []string{"Cache", "Code Cache", "GPUCache", "ShaderCache"}
	for _, name := range cacheDirs {
		dir := filepath.Join(profileDir, name)
		os.MkdirAll(dir, 0755)
		os.WriteFile(filepath.Join(dir, "data.bin"), []byte("cached"), 0644)
	}

	rule := &models.CleanRule{
		ID: "browser_cache_test",
		Targets: []models.Target{
			{
				Type:      "folder",
				Path:      userDataDir,
				Pattern:   "*Cache",
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

	if len(results) != 4 {
		var names []string
		for _, r := range results {
			names = append(names, filepath.Base(r.Path))
		}
		t.Fatalf("Expected 4 cache dirs (Cache, Code Cache, GPUCache, ShaderCache), got %d: %v", len(results), names)
	}

	for _, item := range results {
		if item.Type != "directory" {
			t.Errorf("Expected directory type, got '%s' for %s", item.Type, item.Path)
		}
		if item.Size != 6 {
			t.Errorf("Expected size 6 for %s, got %d", item.Path, item.Size)
		}
	}
}

func TestScanner_ScanRule_MultipleTargets(t *testing.T) {
	tmpDir := t.TempDir()
	dir1 := filepath.Join(tmpDir, "WER-System")
	dir2 := filepath.Join(tmpDir, "WER-User")
	os.MkdirAll(dir1, 0755)
	os.MkdirAll(dir2, 0755)
	os.WriteFile(filepath.Join(dir1, "report.wer"), []byte("sys-report"), 0644)
	os.WriteFile(filepath.Join(dir2, "report.wer"), []byte("user-report"), 0644)

	os.Setenv("TEST_WER_USER", dir2)
	defer os.Unsetenv("TEST_WER_USER")

	rule := &models.CleanRule{
		ID: "error_reports_test",
		Targets: []models.Target{
			{Type: "folder", Path: dir1, Pattern: "*", Recursive: true, MaxDepth: 3},
			{Type: "folder", Path: "%TEST_WER_USER%", Pattern: "*", Recursive: true, MaxDepth: 3},
		},
	}

	scanner := NewScanner()
	ctx := context.Background()
	results, err := scanner.ScanRule(ctx, rule)
	if err != nil {
		t.Fatalf("ScanRule() error = %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("Expected 2 results (system + user WER), got %d", len(results))
	}
}

func TestScanner_ScanRule_LogPattern(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "CBS.log"), []byte("log-content"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "CBS.persist"), []byte("persist-data"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "other.dat"), []byte("other"), 0644)

	rule := &models.CleanRule{
		ID: "cbs_logs_test",
		Targets: []models.Target{
			{Type: "folder", Path: tmpDir, Pattern: "*.log", Recursive: false},
		},
	}

	scanner := NewScanner()
	ctx := context.Background()
	results, err := scanner.ScanRule(ctx, rule)
	if err != nil {
		t.Fatalf("ScanRule() error = %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("Expected 1 result (*.log only), got %d", len(results))
	}
	if filepath.Base(results[0].Path) != "CBS.log" {
		t.Errorf("Expected CBS.log, got %s", results[0].Path)
	}
}
