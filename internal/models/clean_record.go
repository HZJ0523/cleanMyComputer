package models

import "time"

type CleanRecord struct {
	ID          int64     `json:"id"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	ScanLevel   int       `json:"scan_level"`
	TotalFiles  int       `json:"total_files"`
	TotalSize   int64     `json:"total_size"`
	FreedSize   int64     `json:"freed_size"`
	FailedCount int       `json:"failed_count"`
	Status      string    `json:"status"`
}
