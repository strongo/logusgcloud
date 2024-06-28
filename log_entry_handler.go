package logusgcloud

import (
	"cloud.google.com/go/logging"
	"context"
	"github.com/strongo/logus"
)

func NewLogEntryHandler() logus.LogEntryHandler {
	return logEntryHandler{}
}

type logEntryHandler struct {
}

func (l logEntryHandler) Log(ctx context.Context, entry logus.LogEntry) error {
	gcLogEntry := logging.Entry{
		Severity: getGCSeverity(entry.Severity),
		Trace:    logus.GetTraceID(ctx),
		SpanID:   logus.GetSpanID(ctx),
		Labels:   logus.GetLabels(ctx),
	}
	if entry.Message != "" {
		gcLogEntry.Payload = entry.Message
	} else if entry.Payload != nil {
		gcLogEntry.Payload = entry.Payload
	}
	if logger := getLogger(ctx); logger != nil {
		logger.Log(gcLogEntry)
	}
	return nil
}
