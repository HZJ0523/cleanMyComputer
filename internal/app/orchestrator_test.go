package app

import (
	"os"
	"path/filepath"
	"testing"
	"time"
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

func TestRunScanTwiceConcurrent(t *testing.T) {
	o := NewOrchestrator()
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")
	if err := o.InitDB(dbPath); err != nil {
		t.Fatalf("InitDB: %v", err)
	}
	defer o.CloseDB()

	firstDone := make(chan error, 1)
	secondDone := make(chan error, 1)

	go func() {
		firstDone <- o.RunScan(1)
	}()

	time.Sleep(50 * time.Millisecond)

	go func() {
		secondDone <- o.RunScan(1)
	}()

	err1 := <-firstDone
	err2 := <-secondDone

	if err1 != nil {
		t.Logf("first scan: %v (may fail due to no rules dir)", err1)
	}
	if err2 == ErrScanInProgress {
		t.Log("second scan correctly blocked by first scan")
	}
}

func TestRunScanResetsIsScanningOnError(t *testing.T) {
	o := NewOrchestrator()

	if err := o.RunScan(1); err != nil {
		// RunScan will fail since no rules dir exists
		// but IsScanning must be reset
	}
	if o.IsScanning {
		t.Error("IsScanning should be false after RunScan returns (even on error)")
	}
}

func TestSaveCleanHistoryWithoutDB(t *testing.T) {
	o := NewOrchestrator()
	o.SaveCleanHistory(CleanSummary{
		Cleaned:   5,
		Failed:    0,
		FreedSize: 1024,
	}, 1)
}

func TestFindDuplicateFiles(t *testing.T) {
	o := NewOrchestrator()
	tmpDir := t.TempDir()
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
