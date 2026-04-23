package cleaner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestQuarantineManager_Quarantine(t *testing.T) {
	tmpDir := t.TempDir()
	qDir := filepath.Join(tmpDir, "quarantine")
	testFile := filepath.Join(tmpDir, "test.txt")
	os.WriteFile(testFile, []byte("test"), 0644)

	qm, _ := NewQuarantineManager(qDir)
	err := qm.Quarantine(testFile)
	if err != nil {
		t.Fatalf("Quarantine() error = %v", err)
	}

	if _, err := os.Stat(testFile); !os.IsNotExist(err) {
		t.Error("Expected file to be moved")
	}
}

func TestQuarantineManager_Restore(t *testing.T) {
	tmpDir := t.TempDir()
	qDir := filepath.Join(tmpDir, "quarantine")
	testFile := filepath.Join(tmpDir, "test.txt")
	os.WriteFile(testFile, []byte("test"), 0644)

	qm, _ := NewQuarantineManager(qDir)
	qm.Quarantine(testFile)

	files, _ := os.ReadDir(qDir)
	if len(files) == 0 {
		t.Fatal("Expected quarantined file")
	}

	quarantinedPath := filepath.Join(qDir, files[0].Name())
	err := qm.Restore(quarantinedPath, testFile)
	if err != nil {
		t.Fatalf("Restore() error = %v", err)
	}

	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Error("Expected file to be restored")
	}
}
