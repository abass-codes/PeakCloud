package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestRequestLoggerWritesStructuredRequestLog(t *testing.T) {
	gin.SetMode(gin.TestMode)

	core, logs := observer.New(zap.InfoLevel)
	logger := zap.New(core)

	router := gin.New()
	router.Use(RequestID())
	router.Use(RequestLogger(logger))

	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	req := httptest.NewRequest(
		http.MethodGet,
		"/test",
		nil,
	)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if logs.Len() != 1 {
		t.Fatalf(
			"expected 1 log entry, got %d",
			logs.Len(),
		)
	}

	entry := logs.All()[0]

	if entry.Message != "http_request" {
		t.Fatalf(
			"expected http_request log, got %q",
			entry.Message,
		)
	}

	fields := entry.ContextMap()

	expected := []string{
		"method",
		"path",
		"status",
		"request_id",
		"latency_ms",
		"client_ip",
	}

	for _, field := range expected {
		if _, ok := fields[field]; !ok {
			t.Fatalf(
				"expected log field %q",
				field,
			)
		}
	}

	if fields["method"] != http.MethodGet {
		t.Fatalf(
			"expected GET method, got %v",
			fields["method"],
		)
	}

	if fields["path"] != "/test" {
		t.Fatalf(
			"expected /test path, got %v",
			fields["path"],
		)
	}
}
