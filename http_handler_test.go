package logusgcloud

import (
	"cloud.google.com/go/logging"
	"net/http"
	"testing"
)

func TestHttpHandlerForAppEngine(t *testing.T) {

	loggersCreated := 0
	newLogger := func(logID string, opts ...logging.LoggerOption) *logging.Logger {
		loggersCreated++
		return &logging.Logger{}
	}

	handlerCalled := 0
	var handler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		handlerCalled++
	}

	wrappedHandler := HttpHandlerForAppEngine(handler, newLogger)
	if wrappedHandler == nil {
		t.Error("Expected handler to be created, but got nil")

	}
	if loggersCreated != 1 {
		t.Errorf("Expected 1 logger to be created, but got %v", loggersCreated)
	}
	if handlerCalled != 0 {
		t.Errorf("Expected handler to be called 0 times, but got %v", handlerCalled)
	}
}
