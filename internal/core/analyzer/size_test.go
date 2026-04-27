package analyzer

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindLargeFiles_FindsLarge(t *testing.T) {
	tmpDir := t.TempDir()
	large := make([]byte, 2048)
	os.WriteFile(filepath.Join(tmpDir, "big.dat"), large, 0644)

	finder := NewLargeFileFinder(1024)
	files, err := finder.FindLargeFiles(tmpDir)
	if err != nil {
		t.Fatalf("FindLargeFiles: %v", err)
	}
	if len(files) == 0 {
		t.Fatal("expected at least one large file")
	}
	if files[0].Size != 2048 {
		t.Errorf("expected size 2048, got %d", files[0].Size)
	}
}

func TestFindLargeFiles_SkipsSmall(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "small.dat"), []byte("x"), 0644)

	finder := NewLargeFileFinder(1024)
	files, err := finder.FindLargeFiles(tmpDir)
	if err != nil {
		t.Fatalf("FindLargeFiles: %v", err)
	}
	if len(files) != 0 {
		t.Errorf("expected 0 files, got %d", len(files))
	}
}

func TestFindLargeFiles_SortedDescending(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "a.dat"), make([]byte, 3000), 0644)
	os.WriteFile(filepath.Join(tmpDir, "b.dat"), make([]byte, 2000), 0644)
	os.WriteFile(filepath.Join(tmpDir, "c.dat"), make([]byte, 4000), 0644)

	finder := NewLargeFileFinder(1)
	files, err := finder.FindLargeFiles(tmpDir)
	if err != nil {
		t.Fatalf("FindLargeFiles: %v", err)
	}
	if len(files) < 3 {
		t.Fatalf("expected 3 files, got %d", len(files))
	}
	if files[0].Size < files[1].Size {
		t.Errorf("expected descending order: %d >= %d", files[0].Size, files[1].Size)
	}
}

func TestFindLargeFiles_NonExistentDir(t *testing.T) {
	finder := NewLargeFileFinder(1024)
	_, err := finder.FindLargeFiles("/nonexistent/path/12345")
	if err == nil {
		t.Error("expected error for nonexistent directory")
	}
}

func TestFindLargeFiles_DefaultThreshold(t *testing.T) {
	finder := NewLargeFileFinder(0)
	if finder.threshold != 100*1024*1024 {
		t.Errorf("expected default 100MB, got %d", finder.threshold)
	}
}
