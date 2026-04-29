package windows

import (
	"os"
	"testing"
)

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
