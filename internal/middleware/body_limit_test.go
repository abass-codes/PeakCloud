package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestBodyLimitAllowsSmallRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(BodyLimit(10))
	router.POST("/", func(c *gin.Context) {
		_, err := c.GetRawData()
		if err != nil {
			c.Status(http.StatusRequestEntityTooLarge)
			return
		}

		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(
		http.MethodPost,
		"/",
		bytes.NewBufferString("small"),
	)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf(
			"expected %d, got %d",
			http.StatusOK,
			recorder.Code,
		)
	}
}

func TestBodyLimitRejectsLargeRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(BodyLimit(5))
	router.POST("/", func(c *gin.Context) {
		_, err := c.GetRawData()
		if err != nil {
			c.Status(http.StatusRequestEntityTooLarge)
			return
		}

		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(
		http.MethodPost,
		"/",
		bytes.NewBufferString("too-large"),
	)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf(
			"expected %d, got %d",
			http.StatusRequestEntityTooLarge,
			recorder.Code,
		)
	}
}

func TestBodyLimitRejectsKnownOversizedContentLength(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(BodyLimit(5))
	router.POST("/", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(
		http.MethodPost,
		"/",
		bytes.NewBufferString("too-large"),
	)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf(
			"expected %d, got %d",
			http.StatusRequestEntityTooLarge,
			recorder.Code,
		)
	}
}
