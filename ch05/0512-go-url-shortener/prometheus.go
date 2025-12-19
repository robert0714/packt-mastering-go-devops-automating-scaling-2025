package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
)

// Prometheus metrics
var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"method", "route", "status_code"},
	)
	httpDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_duration_seconds",
			Help:    "Histogram of HTTP request durations.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "route"},
	)
	httpContentLength = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_content_length_bytes",
			Help:    "Histogram of HTTP response content lengths.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "route"},
	)
)

func init() {
	// Register metrics with Prometheus
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpDuration)
	prometheus.MustRegister(httpContentLength)
}

// Function to log request metrics
func logRequestMetrics(
	r *http.Request,
	statusCode int,
	duration time.Duration,
) {
	route := mux.CurrentRoute(r).GetName()
	httpRequestsTotal.WithLabelValues(
		r.Method, route, fmt.Sprintf("%d", statusCode)).Inc()
	httpDuration.WithLabelValues(r.Method, route).Observe(
		duration.Seconds())

	// Here, we'd normally calculate the actual content length
	contentLength := 700 // Stubbed content length for this example
	httpContentLength.WithLabelValues(r.Method, route).Observe(
		float64(contentLength),
	)
}
