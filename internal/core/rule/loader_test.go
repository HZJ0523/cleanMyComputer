package rule

import (
	"testing"
)

func TestLoader_LoadFromFile(t *testing.T) {
	loader := NewLoader()
	rules, err := loader.LoadFromFile("../../../configs/rules/level1_safe.json")
	if err != nil {
		t.Fatalf("LoadFromFile() error = %v", err)
	}

	if len(rules) == 0 {
		t.Error("Expected rules to be loaded")
	}

	for _, rule := range rules {
		if err := rule.Validate(); err != nil {
			t.Errorf("Invalid rule %s: %v", rule.ID, err)
		}
	}
}

func TestLoader_LoadFromFile_NotFound(t *testing.T) {
	loader := NewLoader()
	_, err := loader.LoadFromFile("nonexistent.json")
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}
