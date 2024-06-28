package logusgcloud

import (
	"cloud.google.com/go/logging"
	"context"
	"fmt"
	"github.com/strongo/logus"
	"net/http"
	"testing"
)

func Test_getLogger(t *testing.T) {
	ctx := context.Background()
	logger := &logging.Logger{}
	ctx = context.WithValue(ctx, &gcLoggerKey, logger)
	if got := getLogger(ctx); got != logger {
		t.Errorf("getLogger() = %v, want %v", got, logger)
	}
}

func Test_withLogger(t *testing.T) {
	ctx := context.Background()
	logger := &logging.Logger{}
	ctxWithLogger := withLogger(ctx, logger)
	ctxLogger := ctxWithLogger.Value(&gcLoggerKey).(*logging.Logger)
	if ctxLogger != logger {
		t.Errorf("withLogger() = %p, want %p", ctxLogger, logger)
	}
}

func Test_getTraceAndSpanIDs(t *testing.T) {
	type args struct {
		r         *http.Request
		projectID string
	}
	tests := []struct {
		name        string
		args        args
		wantTraceID string
		wantSpanID  string
	}{
		{
			name: "empty",
			args: args{
				r: &http.Request{},
			},
			wantTraceID: "",
			wantSpanID:  "",
		},
		{
			name: "empty",
			args: args{
				r: &http.Request{
					Header: map[string][]string{
						"X-Cloud-Trace-Context": {"trace1234567890/span1234567890;o=1"},
					},
				},
			},
			wantTraceID: "trace1234567890",
			wantSpanID:  "span1234567890",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTraceID, gotSpanID := getTraceAndSpanIDs(tt.args.r, tt.args.projectID)
			if gotTraceID != tt.wantTraceID {
				t.Errorf("getTraceAndSpanIDs() gotTraceID = %v, want %v", gotTraceID, tt.wantTraceID)
			}
			if gotSpanID != tt.wantSpanID {
				t.Errorf("getTraceAndSpanIDs() gotSpanID = %v, want %v", gotSpanID, tt.wantSpanID)
			}
		})
	}
}

func Test_withAppEngineContext(t *testing.T) {
	ctx := context.Background()
	const shortTraceID = "trace1234567890"
	const expectedTranceID = "projects/proj123/traces/" + shortTraceID
	const expectedSpanID = "span1234567890"
	r := &http.Request{
		Header: map[string][]string{
			"X-Cloud-Trace-Context": {fmt.Sprintf("%s/%s;o=1", shortTraceID, expectedSpanID)},
		},
	}
	const projectID = "proj123"
	const gaeInstanceID = "gaeInstance123"
	ctx = withAppEngineContext(ctx, r, projectID, gaeInstanceID)
	if ctx == nil {
		t.Errorf("withAppEngineContext() = nil")
	}
	if traceID := logus.GetTraceID(ctx); traceID != expectedTranceID {
		t.Errorf("withAppEngineContext() traceID = %s, want %s", traceID, expectedTranceID)
	}
	if spanID := logus.GetSpanID(ctx); spanID != expectedSpanID {
		t.Errorf("withAppEngineContext() spanID = %s, want %s", spanID, expectedSpanID)
	}
}
