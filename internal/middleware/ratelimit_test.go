package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestRateLimiterAllowsRequestsWithinLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)

	limiter := NewRateLimiter(2, time.Minute)

	router := gin.New()
	router.Use(limiter.Middleware())
	router.GET("/", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	for i := 0; i < 2; i++ {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = "192.0.2.1:1234"

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Fatalf(
				"request %d: expected %d, got %d",
				i+1,
				http.StatusOK,
				recorder.Code,
			)
		}
	}
}

func TestRateLimiterRejectsRequestOverLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)

	limiter := NewRateLimiter(1, time.Minute)

	router := gin.New()
	router.Use(limiter.Middleware())
	router.GET("/", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	first := httptest.NewRequest(http.MethodGet, "/", nil)
	first.RemoteAddr = "192.0.2.1:1234"

	firstRecorder := httptest.NewRecorder()
	router.ServeHTTP(firstRecorder, first)

	if firstRecorder.Code != http.StatusOK {
		t.Fatalf(
			"expected first request %d, got %d",
			http.StatusOK,
			firstRecorder.Code,
		)
	}

	second := httptest.NewRequest(http.MethodGet, "/", nil)
	second.RemoteAddr = "192.0.2.1:1234"

	secondRecorder := httptest.NewRecorder()
	router.ServeHTTP(secondRecorder, second)

	if secondRecorder.Code != http.StatusTooManyRequests {
		t.Fatalf(
			"expected second request %d, got %d",
			http.StatusTooManyRequests,
			secondRecorder.Code,
		)
	}
}

func TestRateLimiterSeparatesClients(t *testing.T) {
	gin.SetMode(gin.TestMode)

	limiter := NewRateLimiter(1, time.Minute)

	router := gin.New()
	router.Use(limiter.Middleware())
	router.GET("/", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	for _, address := range []string{
		"192.0.2.1:1234",
		"192.0.2.2:1234",
	} {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = address

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Fatalf(
				"client %s: expected %d, got %d",
				address,
				http.StatusOK,
				recorder.Code,
			)
		}
	}
}
