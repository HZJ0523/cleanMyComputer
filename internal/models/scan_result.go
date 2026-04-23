package models

import "time"

type ScanResult struct {
	ID        int64       `json:"id"`
	Level     int         `json:"level"`
	StartTime time.Time   `json:"start_time"`
	EndTime   time.Time   `json:"end_time"`
	Items     []*ScanItem `json:"items"`
	Status    string      `json:"status"`
}

type ScanItem struct {
	Path      string    `json:"path"`
	Size      int64     `json:"size"`
	ModTime   time.Time `json:"mod_time"`
	RuleID    string    `json:"rule_id"`
	RiskScore int       `json:"risk_score"`
}

func (s *ScanResult) TotalSize() int64 {
	var total int64
	for _, item := range s.Items {
		total += item.Size
	}
	return total
}

func (s *ScanResult) TotalCount() int {
	return len(s.Items)
}
