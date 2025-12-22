package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// wrapper for ResponseWriter class
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (r *responseWriter) WriteHeader(status int) {
	r.statusCode = status
	r.ResponseWriter.WriteHeader(status)
}

var (
	// The Prometheus metrics that will be exposed.
	httpHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "my_app_http_hit_total",
			Help: "Total number of HTTP hits.",
		},
	)

	httpStatusCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "response_status_count",
			Help: "Status of the HTTP response.",
		},
		[]string{"status"})

	httpMethodCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_method",
			Help: "HTTP method used.",
		},
		[]string{"method"})

	httpResponseTimeSeconds = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_response_time_bucket",
			Help: "Duration of the HTTP request.",
		},
		[]string{"endpoint"})

	// Add all metrics that will be collected
	metricsList = []prometheus.Collector{
		httpHits,
		httpStatusCounter,
		httpMethodCounter,
		httpResponseTimeSeconds,
	}

	// Prometheus Registry to register metrics.
	prometheusRegistry = prometheus.NewRegistry()
)

func init() {
	// Register metrics that will be exposed.
	prometheusRegistry.MustRegister(metricsList...)
}

func metricsHandler() http.Handler {
	return promhttp.HandlerFor(prometheusRegistry, promhttp.HandlerOpts{})
}

func count(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route := mux.CurrentRoute(r)
		path, _ := route.GetPathTemplate()

		rw := &responseWriter{w, http.StatusOK}

		timer := prometheus.NewTimer(httpResponseTimeSeconds.WithLabelValues(fmt.Sprintf("%s %s", r.Method, path)))
		defer timer.ObserveDuration()

		next.ServeHTTP(rw, r)

		httpHits.Inc()
		httpStatusCounter.WithLabelValues(strconv.Itoa(rw.statusCode)).Inc()
		httpMethodCounter.WithLabelValues(r.Method).Inc()
	})
}
