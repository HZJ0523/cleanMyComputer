package storage

import (
	"testing"
	"time"

	"github.com/hzj0523/cleanMyComputer/internal/models"
)

func TestHistory_SaveAndGetAll(t *testing.T) {
	db, _ := NewDB(":memory:")
	defer db.Close()

	history := NewHistory(db)
	record := &models.CleanRecord{
		StartTime:  time.Now(),
		EndTime:    time.Now(),
		ScanLevel:  1,
		TotalFiles: 10,
		TotalSize:  1024,
		FreedSize:  512,
		Status:     "success",
	}

	id, err := history.Save(record)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}
	if id == 0 {
		t.Error("Expected non-zero ID")
	}

	records, err := history.GetAll()
	if err != nil {
		t.Fatalf("GetAll() error = %v", err)
	}
	if len(records) != 1 {
		t.Errorf("Expected 1 record, got %d", len(records))
	}
}
