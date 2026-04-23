package scanner

import "testing"

func TestScanner_WorkerCount(t *testing.T) {
	scanner := NewParallelScanner(4)
	if scanner.Workers() != 4 {
		t.Errorf("Workers() = %d, want 4", scanner.Workers())
	}
}

func TestParallelScanner_MinWorkers(t *testing.T) {
	scanner := NewParallelScanner(0)
	if scanner.Workers() != 1 {
		t.Errorf("Workers() = %d, want 1 (minimum)", scanner.Workers())
	}
}
