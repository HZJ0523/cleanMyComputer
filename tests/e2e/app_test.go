package e2e

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/hzj0523/cleanMyComputer/internal/app"
	"github.com/hzj0523/cleanMyComputer/internal/models"
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

func TestScanItemsSafe(t *testing.T) {
	o := app.NewOrchestrator()

	if count := o.GetScanItemCount(); count != 0 {
		t.Errorf("expected 0 items, got %d", count)
	}

	items := o.GetScanItemsSafe()
	if len(items) != 0 {
		t.Errorf("expected empty slice, got %d", len(items))
	}

	o.ScanItems = []*models.ScanItem{
		{Path: "test.txt", Size: 100},
		{Path: "test2.txt", Size: 200},
	}

	if count := o.GetScanItemCount(); count != 2 {
		t.Errorf("expected 2 items, got %d", count)
	}

	safeCopy := o.GetScanItemsSafe()
	if len(safeCopy) != 2 {
		t.Fatalf("expected 2 safe items, got %d", len(safeCopy))
	}
	if safeCopy[0].Path != "test.txt" {
		t.Errorf("unexpected path: %s", safeCopy[0].Path)
	}

	o.ClearScanItems()
	if count := o.GetScanItemCount(); count != 0 {
		t.Errorf("expected 0 after clear, got %d", count)
	}
}

func TestRunCleanEmpty(t *testing.T) {
	o := app.NewOrchestrator()

	summary, err := o.RunClean()
	if err != nil {
		t.Fatalf("RunClean with no items should not error: %v", err)
	}
	if summary.Cleaned != 0 {
		t.Errorf("expected 0 cleaned, got %d", summary.Cleaned)
	}
}
