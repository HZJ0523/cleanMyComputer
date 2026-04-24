package i18n

import "testing"

func TestTKeyFallback(t *testing.T) {
	result := T("nonexistent.key")
	if result != "nonexistent.key" {
		t.Errorf("expected key as fallback, got %q", result)
	}
}

func TestTEmptyKey(t *testing.T) {
	result := T("")
	if result != "" {
		t.Errorf("expected empty string, got %q", result)
	}
}
