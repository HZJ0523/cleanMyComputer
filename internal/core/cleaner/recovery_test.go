package cleaner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRecovery_RestoreFile(t *testing.T) {
	tmpDir := t.TempDir()
	qDir := filepath.Join(tmpDir, "quarantine")
	originalPath := filepath.Join(tmpDir, "test.txt")

	qm, _ := NewQuarantineManager(qDir)
	recovery := NewRecovery(qm)

	os.WriteFile(originalPath, []byte("test"), 0644)
	qm.Quarantine(originalPath)

	files, _ := os.ReadDir(qDir)
	if len(files) == 0 {
		t.Fatal("Expected quarantined file")
	}

	quarantinedPath := filepath.Join(qDir, files[0].Name())
	err := recovery.RestoreFile(quarantinedPath, originalPath)
	if err != nil {
		t.Fatalf("RestoreFile() error = %v", err)
	}

	if _, err := os.Stat(originalPath); os.IsNotExist(err) {
		t.Error("Expected file to be restored")
	}
}
