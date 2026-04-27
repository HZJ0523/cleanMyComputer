package storage

import (
	"fmt"
	"log"

	"github.com/hzj0523/cleanMyComputer/internal/models"
)

type History struct {
	db *DB
}

func NewHistory(db *DB) *History {
	return &History{db: db}
}

func (h *History) Save(record *models.CleanRecord) (int64, error) {
	result, err := h.db.Conn().Exec(`
		INSERT INTO clean_history (start_time, end_time, scan_level, total_files, total_size, freed_size, failed_count, status)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, record.StartTime, record.EndTime, record.ScanLevel, record.TotalFiles, record.TotalSize, record.FreedSize, record.FailedCount, record.Status)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (h *History) GetAll() ([]*models.CleanRecord, error) {
	rows, err := h.db.Conn().Query(`SELECT id, start_time, end_time, scan_level, total_files, total_size, freed_size, failed_count, status FROM clean_history ORDER BY start_time DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []*models.CleanRecord
	for rows.Next() {
		var r models.CleanRecord
		if err := rows.Scan(&r.ID, &r.StartTime, &r.EndTime, &r.ScanLevel, &r.TotalFiles, &r.TotalSize, &r.FreedSize, &r.FailedCount, &r.Status); err != nil {
			log.Printf("Warning: skipping corrupt history record: %v", err)
			continue
		}
		records = append(records, &r)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed iterating history rows: %w", err)
	}
	return records, nil
}
