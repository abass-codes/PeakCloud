package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TestRequestIDGeneratesID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(RequestID())

	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	request := httptest.NewRequest(
		http.MethodGet,
		"/test",
		nil,
	)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	requestID := response.Header().Get(RequestIDHeader)

	if requestID == "" {
		t.Fatal("expected request ID header")
	}

	if _, err := uuid.Parse(requestID); err != nil {
		t.Fatalf(
			"expected generated UUID, got %q: %v",
			requestID,
			err,
		)
	}
}

func TestRequestIDPreservesIncomingID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(RequestID())

	router.GET("/test", func(c *gin.Context) {
		value, exists := c.Get("request_id")
		if !exists {
			t.Fatal("request_id missing from context")
		}

		requestID, ok := value.(string)
		if !ok {
			t.Fatal("request_id is not a string")
		}

		c.String(http.StatusOK, requestID)
	})

	request := httptest.NewRequest(
		http.MethodGet,
		"/test",
		nil,
	)

	request.Header.Set(
		RequestIDHeader,
		"test-request-123",
	)

	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if got := response.Header().Get(RequestIDHeader); got != "test-request-123" {
		t.Fatalf(
			"response request ID = %q; want %q",
			got,
			"test-request-123",
		)
	}

	if got := response.Body.String(); got != "test-request-123" {
		t.Fatalf(
			"context request ID = %q; want %q",
			got,
			"test-request-123",
		)
	}
}
