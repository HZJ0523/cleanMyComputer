package analyzer

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindDuplicates_IdenticalFiles(t *testing.T) {
	tmpDir := t.TempDir()
	content := make([]byte, 2048)
	for i := range content {
		content[i] = byte(i % 256)
	}
	os.WriteFile(filepath.Join(tmpDir, "a.txt"), content, 0644)
	os.WriteFile(filepath.Join(tmpDir, "b.txt"), content, 0644)

	finder := NewDuplicateFinder(1024)
	groups, err := finder.FindDuplicates(tmpDir)
	if err != nil {
		t.Fatalf("FindDuplicates: %v", err)
	}
	if len(groups) == 0 {
		t.Fatal("expected at least one duplicate group")
	}
	if len(groups[0].Paths) != 2 {
		t.Errorf("expected 2 paths in group, got %d", len(groups[0].Paths))
	}
	if len(groups[0].Hash) != 32 {
		t.Errorf("expected hex-encoded hash of length 32, got %d", len(groups[0].Hash))
	}
}

func TestFindDuplicates_NoDuplicates(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "a.txt"), []byte("content-a-unique"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "b.txt"), []byte("content-b-unique"), 0644)

	finder := NewDuplicateFinder(1)
	groups, err := finder.FindDuplicates(tmpDir)
	if err != nil {
		t.Fatalf("FindDuplicates: %v", err)
	}
	if len(groups) != 0 {
		t.Errorf("expected 0 groups for different files, got %d", len(groups))
	}
}

func TestFindDuplicates_BelowMinSize(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "a.txt"), []byte("x"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "b.txt"), []byte("x"), 0644)

	finder := NewDuplicateFinder(1024)
	groups, err := finder.FindDuplicates(tmpDir)
	if err != nil {
		t.Fatalf("FindDuplicates: %v", err)
	}
	if len(groups) != 0 {
		t.Error("expected 0 groups for files below min size")
	}
}

func TestFindDuplicates_NonExistentDir(t *testing.T) {
	finder := NewDuplicateFinder(1024)
	_, err := finder.FindDuplicates("/nonexistent/path/12345")
	if err == nil {
		t.Error("expected error for nonexistent directory")
	}
}
