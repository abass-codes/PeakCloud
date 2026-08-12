package middleware

import (
	"fmt"
	"net/http"
	"sort"
	"sync"
	"sync/atomic"

	"github.com/gin-gonic/gin"
)

type HTTPMetrics struct {
	requestsTotal atomic.Uint64

	mu       sync.RWMutex
	statuses map[int]uint64
}

func NewHTTPMetrics() *HTTPMetrics {
	return &HTTPMetrics{
		statuses: make(map[int]uint64),
	}
}

func (m *HTTPMetrics) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		status := c.Writer.Status()

		m.requestsTotal.Add(1)

		m.mu.Lock()
		m.statuses[status]++
		m.mu.Unlock()
	}
}

func (m *HTTPMetrics) RequestsTotal() uint64 {
	return m.requestsTotal.Load()
}

func (m *HTTPMetrics) StatusTotal(status int) uint64 {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.statuses[status]
}

func (m *HTTPMetrics) Handler(c *gin.Context) {
	m.mu.RLock()

	statuses := make(map[int]uint64, len(m.statuses))
	statusCodes := make([]int, 0, len(m.statuses))

	for status, count := range m.statuses {
		statuses[status] = count
		statusCodes = append(statusCodes, status)
	}

	m.mu.RUnlock()

	sort.Ints(statusCodes)

	c.Header(
		"Content-Type",
		"text/plain; version=0.0.4; charset=utf-8",
	)

	c.Writer.WriteHeader(http.StatusOK)

	_, _ = fmt.Fprintln(
		c.Writer,
		"# HELP peakcloud_http_requests_total Total HTTP requests processed.",
	)

	_, _ = fmt.Fprintln(
		c.Writer,
		"# TYPE peakcloud_http_requests_total counter",
	)

	_, _ = fmt.Fprintf(
		c.Writer,
		"peakcloud_http_requests_total %d\n",
		m.RequestsTotal(),
	)

	_, _ = fmt.Fprintln(
		c.Writer,
		"# HELP peakcloud_http_responses_total HTTP responses by status code.",
	)

	_, _ = fmt.Fprintln(
		c.Writer,
		"# TYPE peakcloud_http_responses_total counter",
	)

	for _, status := range statusCodes {
		_, _ = fmt.Fprintf(
			c.Writer,
			"peakcloud_http_responses_total{status=\"%d\"} %d\n",
			status,
			statuses[status],
		)
	}
}
