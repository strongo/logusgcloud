package logusgcloud

import (
	"context"
	"testing"

	"github.com/strongo/logus"
)

func TestNewLogEntryHandler(t *testing.T) {
	if got := NewLogEntryHandler(); got == nil {
		t.Error("Expected log entry handler to be created, but got nil")
	}
}

func TestLogEntryHandler_Log(t *testing.T) {
	logger := testLogger(t)
	handler := NewLogEntryHandler()
	ctx := withLogger(context.Background(), logger)

	// 1. MessageFormat != "" with len(MessageArgs) == 0
	err := handler.Log(ctx, logus.LogEntry{
		Severity:      logus.SeverityInfo,
		MessageFormat: "hello with no args",
	})
	if err != nil {
		t.Fatalf("Log failed: %v", err)
	}

	// 2. MessageFormat != "" with len(MessageArgs) > 0
	err = handler.Log(ctx, logus.LogEntry{
		Severity:      logus.SeverityWarning,
		MessageFormat: "hello %s",
		MessageArgs:   []any{"world"},
	})
	if err != nil {
		t.Fatalf("Log failed: %v", err)
	}

	// 3. MessageFormat == "" with Payload != nil
	err = handler.Log(ctx, logus.LogEntry{
		Severity: logus.SeverityError,
		Payload:  map[string]string{"foo": "bar"},
	})
	if err != nil {
		t.Fatalf("Log failed: %v", err)
	}

	// 4. MessageFormat == "" and Payload == nil, and no logger in ctx
	err = handler.Log(context.Background(), logus.LogEntry{
		Severity: logus.SeverityDebug,
	})
	if err != nil {
		t.Fatalf("Log failed: %v", err)
	}
}
