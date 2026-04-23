package scanner

import "testing"

func TestFilter_ShouldInclude(t *testing.T) {
	filter := NewFilter()
	filter.AddExclude("*.log")

	if filter.ShouldInclude("app.log") {
		t.Error("Expected app.log to be excluded")
	}
	if !filter.ShouldInclude("app.txt") {
		t.Error("Expected app.txt to be included")
	}
}
