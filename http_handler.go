package logusgcloud

import (
	"cloud.google.com/go/logging"
	"log"
	"net/http"
	"os"
)

func HttpHandlerForAppEngine(handler http.Handler, newLogger func(logID string, opts ...logging.LoggerOption) *logging.Logger) http.Handler {

	logger := newLogger("request_log_entries") // fmt.Sprintf("projects/%s/logs/request_log_entries", projectID)

	gaeInstanceID := os.Getenv("GAE_INSTANCE")

	var wrapper http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
		ctx = withAppEngineContext(ctx, r, projectID, gaeInstanceID)
		ctx = withLogger(ctx, logger)
		r = r.WithContext(ctx)
		defer func() {
			go func() {
				if err := logger.Flush(); err != nil {
					log.Printf("ERROR: failed to flush log entries: %v", err)
				}
			}()
		}()
		handler.ServeHTTP(w, r)
	}

	return wrapper
}
