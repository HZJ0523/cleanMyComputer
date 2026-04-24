package e2e

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/hzj0523/cleanMyComputer/internal/app"
	"github.com/hzj0523/cleanMyComputer/internal/core/cleaner"
)

func TestOrchestratorScanAndCleanHistory(t *testing.T) {
	o := app.NewOrchestrator()

	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")
	if err := o.InitDB(dbPath); err != nil {
		t.Fatalf("InitDB: %v", err)
	}
	defer o.CloseDB()

	summary := app.CleanSummary{
		Cleaned:   10,
		Failed:    1,
		FreedSize: 1024 * 1024,
		Duration:  5 * time.Second,
	}
	o.SaveCleanHistory(summary)

	records, err := o.GetHistory()
	if err != nil {
		t.Fatalf("GetHistory: %v", err)
	}
	if len(records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(records))
	}
	r := records[0]
	if r.FreedSize != 1024*1024 {
		t.Errorf("FreedSize = %d, want %d", r.FreedSize, 1024*1024)
	}
	if r.Status != "partial" {
		t.Errorf("Status = %s, want partial", r.Status)
	}
}

func TestOrchestratorConfigRoundTrip(t *testing.T) {
	o := app.NewOrchestrator()

	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")
	if err := o.InitDB(dbPath); err != nil {
		t.Fatalf("InitDB: %v", err)
	}
	defer o.CloseDB()

	if err := o.SetConfig("test_key", "test_value"); err != nil {
		t.Fatalf("SetConfig: %v", err)
	}

	val, err := o.GetConfig("test_key")
	if err != nil {
		t.Fatalf("GetConfig: %v", err)
	}
	if val != "test_value" {
		t.Errorf("got %q, want %q", val, "test_value")
	}
}

func TestOrchestratorQuarantineLifecycle(t *testing.T) {
	o := app.NewOrchestrator()

	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")
	if err := o.InitDB(dbPath); err != nil {
		t.Fatalf("InitDB: %v", err)
	}
	defer o.CloseDB()

	qDir := filepath.Join(tmpDir, "quarantine")
	os.MkdirAll(qDir, 0755)
	qFile := filepath.Join(qDir, "test.dat")
	os.WriteFile(qFile, []byte("quarantined"), 0644)

	now := time.Now()
	record := cleaner.QuarantineRecord{
		OriginalPath:   "/fake/path/file.txt",
		QuarantinePath: qFile,
		Size:           12,
		CreatedAt:      now,
		ExpiresAt:      now.Add(24 * time.Hour),
	}

	if err := o.SaveQuarantineRecord(record); err != nil {
		t.Fatalf("SaveQuarantineRecord: %v", err)
	}

	items, err := o.GetQuarantinedItems()
	if err != nil {
		t.Fatalf("GetQuarantinedItems: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 quarantined item, got %d", len(items))
	}
	if items[0].OriginalPath != record.OriginalPath {
		t.Errorf("OriginalPath = %s, want %s", items[0].OriginalPath, record.OriginalPath)
	}

	if err := o.DeleteQuarantinedItem(qFile); err != nil {
		t.Fatalf("DeleteQuarantinedItem: %v", err)
	}

	items, err = o.GetQuarantinedItems()
	if err != nil {
		t.Fatalf("GetQuarantinedItems after delete: %v", err)
	}
	if len(items) != 0 {
		t.Errorf("expected 0 quarantined items after delete, got %d", len(items))
	}
}

func TestSchedulerStartStop(t *testing.T) {
	o := app.NewOrchestrator()
	s := app.NewScheduler(o)

	if s.Running() {
		t.Error("scheduler should not be running initially")
	}

	s.SetInterval(1 * time.Hour)
	s.Start()
	if !s.Running() {
		t.Error("scheduler should be running after Start()")
	}

	s.Stop()
	if s.Running() {
		t.Error("scheduler should not be running after Stop()")
	}
}
