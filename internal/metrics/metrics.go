package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// Counter: total HTTP requests, labeled by method/path/status for filtering in Grafana
	RequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	// Histogram: request duration — Grafana can show p50/p95/p99 latency from this
	RequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets, // default: .005 to 10 seconds
		},
		[]string{"method", "path"},
	)

	// Gauge: active in-flight requests
	ActiveRequests = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_requests_active",
			Help: "Number of in-flight HTTP requests",
		},
	)
)

func init() {
	prometheus.MustRegister(RequestsTotal, RequestDuration, ActiveRequests)
}

func Handler() http.Handler {
	return promhttp.Handler()
}
