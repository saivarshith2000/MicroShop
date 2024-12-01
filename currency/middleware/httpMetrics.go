package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

var HttpRequestsTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total Number of HTTP requests",
	},
	[]string{"path", "status_code"},
)

var RequestLatency = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "request_latency_microseconds",
		Help: "Request latency in microseconds",
	},
	[]string{"path"},
)

func HttpRequestCounterMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()

		path := c.Request.URL.Path

		// Process request
		c.Next()

		latency := time.Since(t)
		RequestLatency.WithLabelValues(path).Set(float64(latency.Microseconds()))

		status := c.Writer.Status()
		HttpRequestsTotal.WithLabelValues(path, fmt.Sprintf("%d", status)).Inc()
	}
}
