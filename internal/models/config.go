package models

import "errors"

type Config struct {
	AutoCleanEnabled  bool   `json:"auto_clean_enabled"`
	AutoCleanSchedule string `json:"auto_clean_schedule"`
	OldFileDays       int    `json:"old_file_days"`
	ScanWorkers       int    `json:"scan_workers"`
	Language          string `json:"language"`
}

func (c *Config) Validate() error {
	if c.ScanWorkers < 1 {
		return errors.New("scan_workers must be >= 1")
	}
	if c.OldFileDays < 1 {
		return errors.New("old_file_days must be >= 1")
	}
	return nil
}
