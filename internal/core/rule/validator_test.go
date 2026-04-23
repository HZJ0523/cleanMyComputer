package rule

import (
	"testing"

	"github.com/hzj0523/cleanMyComputer/internal/models"
)

func TestValidator_ValidateRule(t *testing.T) {
	validator := NewValidator()

	tests := []struct {
		name    string
		rule    *models.CleanRule
		wantErr bool
	}{
		{
			name: "valid rule",
			rule: &models.CleanRule{
				ID: "test", Name: "Test", Level: 1, RiskScore: 10,
			},
			wantErr: false,
		},
		{
			name: "invalid risk score",
			rule: &models.CleanRule{
				ID: "test", Name: "Test", Level: 1, RiskScore: 150,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateRule(tt.rule)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateRule() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidator_ValidateRules(t *testing.T) {
	validator := NewValidator()

	rules := []*models.CleanRule{
		{ID: "r1", Name: "R1", Level: 1, RiskScore: 10},
		{ID: "r2", Name: "R2", Level: 99, RiskScore: 10},
	}

	errs := validator.ValidateRules(rules)
	if len(errs) != 1 {
		t.Errorf("Expected 1 error, got %d", len(errs))
	}
}
