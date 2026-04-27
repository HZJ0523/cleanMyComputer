package models

import "time"

type ScanItem struct {
	Path      string    `json:"path"`
	Size      int64     `json:"size"`
	ModTime   time.Time `json:"mod_time"`
	RuleID    string    `json:"rule_id"`
	RiskScore int       `json:"risk_score"`
	Type      string    `json:"type"` // "file" or "command"
}
