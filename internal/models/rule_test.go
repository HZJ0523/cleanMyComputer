package models

import (
	"testing"
)

func TestCleanRule_Validate(t *testing.T) {
	tests := []struct {
		name    string
		rule    CleanRule
		wantErr bool
	}{
		{
			name: "valid rule",
			rule: CleanRule{
				ID:          "test_rule",
				Name:        "Test Rule",
				Level:       1,
				Description: "Test description",
				RiskScore:   10,
				Enabled:     true,
			},
			wantErr: false,
		},
		{
			name: "missing ID",
			rule: CleanRule{
				Name:  "Test Rule",
				Level: 1,
			},
			wantErr: true,
		},
		{
			name: "invalid level",
			rule: CleanRule{
				ID:    "test_rule",
				Name:  "Test Rule",
				Level: 99,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.rule.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
