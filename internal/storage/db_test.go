package storage

import (
	"testing"
)

func TestDB_Init(t *testing.T) {
	db, err := NewDB(":memory:")
	if err != nil {
		t.Fatalf("NewDB() error = %v", err)
	}
	defer db.Close()

	if db.Conn() == nil {
		t.Error("Expected database connection")
	}
}
