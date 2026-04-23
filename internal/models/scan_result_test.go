package models

import "testing"

func TestScanResult_TotalSize(t *testing.T) {
	result := &ScanResult{
		Items: []*ScanItem{
			{Path: "/tmp/file1.txt", Size: 1024},
			{Path: "/tmp/file2.txt", Size: 2048},
		},
	}

	expected := int64(3072)
	if got := result.TotalSize(); got != expected {
		t.Errorf("TotalSize() = %d, want %d", got, expected)
	}
}

func TestScanResult_TotalCount(t *testing.T) {
	result := &ScanResult{
		Items: []*ScanItem{
			{Path: "/tmp/file1.txt"},
			{Path: "/tmp/file2.txt"},
		},
	}

	if got := result.TotalCount(); got != 2 {
		t.Errorf("TotalCount() = %d, want 2", got)
	}
}
