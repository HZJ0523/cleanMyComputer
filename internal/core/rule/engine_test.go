package rule

import (
	"context"
	"testing"

	"github.com/hzj0523/cleanMyComputer/internal/models"
)

func TestEngine_LoadRules(t *testing.T) {
	loader := NewLoader()
	loader.rulesDir = "../../../configs/rules"
	engine := NewEngine(loader)

	ctx := context.Background()
	err := engine.LoadRules(ctx, 1)
	if err != nil {
		t.Fatalf("LoadRules() error = %v", err)
	}

	rules := engine.GetEnabledRules(1)
	if len(rules) == 0 {
		t.Error("Expected rules to be loaded")
	}
}

func TestEngine_GetEnabledRules(t *testing.T) {
	engine := NewEngine(nil)
	engine.rules = map[string]*models.CleanRule{
		"rule1": {ID: "rule1", Level: 1, Enabled: true},
		"rule2": {ID: "rule2", Level: 2, Enabled: true},
		"rule3": {ID: "rule3", Level: 1, Enabled: false},
	}

	rules := engine.GetEnabledRules(1)
	if len(rules) != 1 {
		t.Errorf("Expected 1 enabled rule, got %d", len(rules))
	}
}
