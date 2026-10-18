package logusgcloud

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"cloud.google.com/go/logging"
	"cloud.google.com/go/logging/apiv2/loggingpb"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type fakeLoggingServer struct {
	loggingpb.UnimplementedLoggingServiceV2Server
}

func (s *fakeLoggingServer) WriteLogEntries(ctx context.Context, req *loggingpb.WriteLogEntriesRequest) (*loggingpb.WriteLogEntriesResponse, error) {
	return &loggingpb.WriteLogEntriesResponse{}, nil
}

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

func testLogger(t *testing.T) *logging.Logger {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = l.Close() })
	srv := grpc.NewServer()
	loggingpb.RegisterLoggingServiceV2Server(srv, &fakeLoggingServer{})
	go func() { _ = srv.Serve(l) }()
	t.Cleanup(func() { srv.Stop() })

	client, err := logging.NewClient(context.Background(), "projects/test-proj",
		option.WithEndpoint(l.Addr().String()),
		option.WithoutAuthentication(),
		option.WithGRPCDialOption(grpc.WithTransportCredentials(insecure.NewCredentials())),
	)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = client.Close() })
	return client.Logger("test")
}

func TestHttpHandlerForAppEngine_ServeHTTP(t *testing.T) {
	logger := testLogger(t)
	handlerCalled := 0
	var handler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		handlerCalled++
		if req := GetRequest(r.Context()); req == nil || req.URL.Path != r.URL.Path {
			t.Errorf("expected request in context, got %v", req)
		}
		if l := getLogger(r.Context()); l != logger {
			t.Errorf("expected logger in context, got %v", l)
		}
	}

	wrapped := HttpHandlerForAppEngine(handler, logger)
	logger.Log(logging.Entry{Payload: "hi"})
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	wrapped.ServeHTTP(w, req)
	if handlerCalled != 1 {
		t.Errorf("expected handlerCalled == 1, got %d", handlerCalled)
	}
	time.Sleep(20 * time.Millisecond) // allow goroutine Flush to complete
}

func TestHttpHandlerForAppEngine_FlushError(t *testing.T) {
	logger := testLogger(t)
	var handler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		// Log an invalid entry with LogName set, causing client.error() so Flush() returns an error
		logger.Log(logging.Entry{LogName: "invalid-should-not-be-set"})
	}

	wrapped := HttpHandlerForAppEngine(handler, logger)
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	wrapped.ServeHTTP(w, req)
	time.Sleep(50 * time.Millisecond) // allow goroutine Flush to complete
}
