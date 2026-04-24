package windows

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWindowsPlatform_ExpandPath(t *testing.T) {
	p := NewPlatform()

	tempDir := os.Getenv("TEMP")
	if tempDir == "" {
		t.Skip("TEMP env var not set")
	}

	result := p.ExpandPath("%TEMP%\\sub")
	expected := filepath.Join(tempDir, "sub")
	if result != expected {
		t.Errorf("ExpandPath() = %s, want %s", result, expected)
	}
}

func TestWindowsPlatform_GetCommonPaths(t *testing.T) {
	p := NewPlatform()
	paths := p.GetCommonPaths()

	if paths["TEMP"] == "" {
		t.Error("Expected TEMP path to be set")
	}
	if paths["HOME"] == "" {
		t.Error("Expected HOME path to be set")
	}
}

func TestWindowsPlatform_GetDiskUsage(t *testing.T) {
	p := NewPlatform()
	home, _ := os.UserHomeDir()
	usage, err := p.GetDiskUsage(home)
	if err != nil {
		t.Fatalf("GetDiskUsage() error = %v", err)
	}
	if usage.TotalGB <= 0 {
		t.Error("Expected positive total disk size")
	}
	if usage.FreeGB <= 0 {
		t.Error("Expected positive free disk space")
	}
}
