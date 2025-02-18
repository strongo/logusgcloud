package logusgcloud

import (
	"cloud.google.com/go/logging"
	"log"
	"net/http"
	"os"
)

func HttpHandlerForAppEngine(handler http.Handler, logger *logging.Logger) http.Handler {

	gaeInstanceID := os.Getenv("GAE_INSTANCE")

	var wrapper http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
		ctx = withAppEngineContext(ctx, r, projectID, gaeInstanceID)
		ctx = withLogger(ctx, logger)
		ctx = withRequest(ctx, r)
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
