package models

type Config struct {
	QuarantineRetentionHours int    `json:"quarantine_retention_hours"`
	AutoCleanEnabled         bool   `json:"auto_clean_enabled"`
	AutoCleanSchedule        string `json:"auto_clean_schedule"`
	OldFileDays              int    `json:"old_file_days"`
	ScanWorkers              int    `json:"scan_workers"`
	Language                 string `json:"language"`
}
