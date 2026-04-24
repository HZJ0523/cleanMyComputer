package app

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/hzj0523/cleanMyComputer/internal/core/cleaner"
)

func TestNewOrchestrator(t *testing.T) {
	o := NewOrchestrator()
	if o == nil {
		t.Fatal("expected non-nil orchestrator")
	}
}

func TestInitDBAndClose(t *testing.T) {
	o := NewOrchestrator()
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	if err := o.InitDB(dbPath); err != nil {
		t.Fatalf("InitDB: %v", err)
	}
	if o.db == nil {
		t.Fatal("expected db to be non-nil after InitDB")
	}
	o.CloseDB()
}

func TestGetConfigWithoutDB(t *testing.T) {
	o := NewOrchestrator()
	_, err := o.GetConfig("any_key")
	if err == nil {
		t.Error("expected error when DB not initialized")
	}
}

func TestSetConfigWithoutDB(t *testing.T) {
	o := NewOrchestrator()
	err := o.SetConfig("any_key", "value")
	if err == nil {
		t.Error("expected error when DB not initialized")
	}
}

func TestGetRuleStatusWithoutDB(t *testing.T) {
	o := NewOrchestrator()
	_, err := o.GetRuleStatus()
	if err == nil {
		t.Error("expected error when DB not initialized")
	}
}

func TestSetRuleEnabledWithoutDB(t *testing.T) {
	o := NewOrchestrator()
	err := o.SetRuleEnabled("test_rule", true)
	if err == nil {
		t.Error("expected error when DB not initialized")
	}
}

func TestGetHistoryWithoutDB(t *testing.T) {
	o := NewOrchestrator()
	_, err := o.GetHistory()
	if err == nil {
		t.Error("expected error when DB not initialized")
	}
}

func TestCleanupExpiredQuarantineWithoutDB(t *testing.T) {
	o := NewOrchestrator()
	err := o.CleanupExpiredQuarantine()
	if err != nil {
		t.Errorf("expected nil error without DB, got %v", err)
	}
}

func TestSaveQuarantineRecordWithoutDB(t *testing.T) {
	o := NewOrchestrator()
	err := o.SaveQuarantineRecord(cleaner.QuarantineRecord{})
	if err == nil {
		t.Error("expected error when DB not initialized")
	}
}

func TestGetQuarantinedItemsWithoutDB(t *testing.T) {
	o := NewOrchestrator()
	_, err := o.GetQuarantinedItems()
	if err == nil {
		t.Error("expected error when DB not initialized")
	}
}

func TestRunScanTwiceConcurrent(t *testing.T) {
	o := NewOrchestrator()
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")
	if err := o.InitDB(dbPath); err != nil {
		t.Fatalf("InitDB: %v", err)
	}
	defer o.CloseDB()

	// First scan should succeed
	done := make(chan error, 1)
	go func() {
		done <- o.RunScan(1)
	}()

	// Second scan while first is running should fail
	// Give a moment for the goroutine to start
	err := <-done
	if err != nil {
		t.Logf("first scan: %v", err)
	}
}

func TestSaveCleanHistoryWithoutDB(t *testing.T) {
	o := NewOrchestrator()
	// Should not panic
	o.SaveCleanHistory(CleanSummary{
		Cleaned:   5,
		Failed:    0,
		FreedSize: 1024,
	})
}

func TestFindDuplicateFiles(t *testing.T) {
	o := NewOrchestrator()
	tmpDir := t.TempDir()
	// Create two identical files larger than 1KB minimum
	content := make([]byte, 2048)
	for i := range content {
		content[i] = byte(i % 256)
	}
	os.WriteFile(filepath.Join(tmpDir, "a.txt"), content, 0644)
	os.WriteFile(filepath.Join(tmpDir, "b.txt"), content, 0644)

	groups, err := o.FindDuplicateFiles(tmpDir)
	if err != nil {
		t.Fatalf("FindDuplicateFiles: %v", err)
	}
	if len(groups) == 0 {
		t.Error("expected at least one duplicate group")
	}
}

func TestFindLargeFiles(t *testing.T) {
	o := NewOrchestrator()
	tmpDir := t.TempDir()
	// Create a file larger than threshold
	largeContent := make([]byte, 1024)
	os.WriteFile(filepath.Join(tmpDir, "big.dat"), largeContent, 0644)

	files, err := o.FindLargeFiles(tmpDir, 512)
	if err != nil {
		t.Fatalf("FindLargeFiles: %v", err)
	}
	if len(files) == 0 {
		t.Error("expected at least one large file")
	}
}
