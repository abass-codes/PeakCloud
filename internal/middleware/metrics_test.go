package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestHTTPMetricsCountsRequests(t *testing.T) {
	gin.SetMode(gin.TestMode)

	metrics := NewHTTPMetrics()

	router := gin.New()
	router.Use(metrics.Middleware())

	router.GET("/ok", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	router.GET("/missing", func(c *gin.Context) {
		c.Status(http.StatusNotFound)
	})

	for _, path := range []string{"/ok", "/missing"} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)
	}

	if got := metrics.RequestsTotal(); got != 2 {
		t.Fatalf("expected 2 requests, got %d", got)
	}

	if got := metrics.StatusTotal(http.StatusOK); got != 1 {
		t.Fatalf("expected one 200 response, got %d", got)
	}

	if got := metrics.StatusTotal(http.StatusNotFound); got != 1 {
		t.Fatalf("expected one 404 response, got %d", got)
	}
}

func TestHTTPMetricsHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	metrics := NewHTTPMetrics()

	app := gin.New()
	app.Use(metrics.Middleware())

	app.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	recorder := httptest.NewRecorder()

	app.ServeHTTP(recorder, req)

	metricsRouter := gin.New()
	metricsRouter.GET("/metrics", metrics.Handler)

	metricsRequest := httptest.NewRequest(
		http.MethodGet,
		"/metrics",
		nil,
	)

	metricsRecorder := httptest.NewRecorder()

	metricsRouter.ServeHTTP(
		metricsRecorder,
		metricsRequest,
	)

	if metricsRecorder.Code != http.StatusOK {
		t.Fatalf(
			"expected %d, got %d",
			http.StatusOK,
			metricsRecorder.Code,
		)
	}

	body := metricsRecorder.Body.String()

	expected := []string{
		"peakcloud_http_requests_total 1",
		`peakcloud_http_responses_total{status="204"} 1`,
	}

	for _, value := range expected {
		if !strings.Contains(body, value) {
			t.Fatalf(
				"expected metrics output to contain %q; body=%q",
				value,
				body,
			)
		}
	}
}
