package models

import (
	"errors"
	"time"
)

// CleanRule 定义清理规则
type CleanRule struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Category      string    `json:"category"`
	Level         int       `json:"level"`
	Description   string    `json:"description"`
	Targets       []Target  `json:"targets"`
	RiskScore     int       `json:"risk_score"`
	RequiresAdmin bool      `json:"requires_admin"`
	Enabled       bool      `json:"enabled"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// Target 定义清理目标
type Target struct {
	Type        string   `json:"type"`
	Path        string   `json:"path"`
	Pattern     string   `json:"pattern"`
	Recursive   bool     `json:"recursive"`
	MaxDepth    int      `json:"max_depth"`
	ExcludeList []string `json:"exclude_list"`
}

// Validate 验证规则有效性
func (r *CleanRule) Validate() error {
	if r.ID == "" {
		return errors.New("rule ID is required")
	}
	if r.Name == "" {
		return errors.New("rule name is required")
	}
	if r.Level < 1 || r.Level > 3 {
		return errors.New("rule level must be 1, 2, or 3")
	}
	if r.RiskScore < 0 || r.RiskScore > 100 {
		return errors.New("risk score must be between 0 and 100")
	}
	return nil
}
