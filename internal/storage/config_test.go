package storage

import (
	"testing"
)

func TestConfig_SetAndGet(t *testing.T) {
	db, _ := NewDB(":memory:")
	defer db.Close()

	config := NewConfig(db)
	config.Set("test_key", "test_value")

	value, err := config.Get("test_key")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if value != "test_value" {
		t.Errorf("Get() = %s, want test_value", value)
	}
}
