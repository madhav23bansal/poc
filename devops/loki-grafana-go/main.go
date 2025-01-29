package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/time/rate"
)

// Prometheus metrics
var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests received",
		},
		[]string{"method", "path", "status"},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Histogram of response time (seconds)",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	requestSizeBytes = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_size_bytes",
			Help:    "Size of HTTP requests in bytes",
			Buckets: prometheus.ExponentialBuckets(100, 2, 10),
		},
		[]string{"method", "path"},
	)

	responseSizeBytes = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_size_bytes",
			Help:    "Size of HTTP responses in bytes",
			Buckets: prometheus.ExponentialBuckets(100, 2, 10),
		},
		[]string{"method", "path"},
	)
)

// Rate limiter for basic DDoS protection
var limiter = rate.NewLimiter(1, 5) // 1 request per second, burst of 5

func init() {
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)
	prometheus.MustRegister(requestSizeBytes)
	prometheus.MustRegister(responseSizeBytes)
}

// Middleware for Prometheus metrics and DDoS protection
func prometheusMiddleware(c *gin.Context) {
	start := time.Now()

	if !limiter.Allow() {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
		return
	}

	// Process request
	c.Next()

	// Calculate metrics
	duration := time.Since(start).Seconds()
	size := c.Writer.Size()
	method := c.Request.Method
	path := c.FullPath()
	status := c.Writer.Status()

	// Update Prometheus metrics
	httpRequestsTotal.WithLabelValues(method, path, http.StatusText(status)).Inc()
	httpRequestDuration.WithLabelValues(method, path).Observe(duration)
	requestSizeBytes.WithLabelValues(method, path).Observe(float64(c.Request.ContentLength))
	responseSizeBytes.WithLabelValues(method, path).Observe(float64(size))
}

func main() {
	r := gin.Default()

	// Attach Prometheus middleware
	r.Use(prometheusMiddleware)

	// Expose Prometheus metrics endpoint
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Sample routes
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	log.Println("Starting server on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
