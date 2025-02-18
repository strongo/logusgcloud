package logusgcloud

import (
	"cloud.google.com/go/logging"
	"net/http"
	"testing"
)

func TestHttpHandlerForAppEngine(t *testing.T) {

	handlerCalled := 0
	var handler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		handlerCalled++
	}

	wrappedHandler := HttpHandlerForAppEngine(handler, &logging.Logger{})
	if wrappedHandler == nil {
		t.Error("Expected handler to be created, but got nil")

	}
	if handlerCalled != 0 {
		t.Errorf("Expected handler to be called 0 times, but got %v", handlerCalled)
	}
}
