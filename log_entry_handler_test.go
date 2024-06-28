package logusgcloud

import (
	"testing"
)

func TestNewLogEntryHandler(t *testing.T) {
	if got := NewLogEntryHandler(); got == nil {
		t.Error("Expected log entry handler to be created, but got nil")
	}
}
