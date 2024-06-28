package logusgcloud

import (
	"cloud.google.com/go/logging"
	"context"
	"fmt"
	"github.com/strongo/logus"
	"net/http"
	"strings"
)

func withAppEngineContext(ctx context.Context, r *http.Request, projectID, gaeInstanceID string) context.Context {
	traceID, spanID := getTraceAndSpanIDs(r, projectID)
	if traceID != "" {
		ctx = logus.WithTraceID(ctx, traceID)
	}
	if spanID != "" {
		ctx = logus.WithSpanID(ctx, spanID)
	}
	logus.WithLabels(ctx, map[string]string{
		"clone_id": gaeInstanceID,
	})
	return ctx
}

func getTraceAndSpanIDs(r *http.Request, projectID string) (traceID string, spanID string) {
	traceID = r.Header.Get("X-Cloud-Trace-Context")
	if i := strings.IndexByte(traceID, ';'); i > 0 {
		traceID = traceID[:i]
	}
	if i := strings.IndexByte(traceID, '/'); i > 0 {
		spanID = traceID[i+1:]
		traceID = traceID[:i]
	}
	if !strings.HasPrefix(traceID, "projects/") {
		if projectID == "" {
			return
		}
		traceID = fmt.Sprintf("projects/%s/traces/%s", projectID, traceID)
	}
	return
}

var gcLoggerKey = "gcLogger"

func withLogger(ctx context.Context, logger *logging.Logger) context.Context {
	return context.WithValue(ctx, &gcLoggerKey, logger)
}

func getLogger(ctx context.Context) *logging.Logger {
	if logger, ok := ctx.Value(&gcLoggerKey).(*logging.Logger); ok {
		return logger
	}
	return nil
}
