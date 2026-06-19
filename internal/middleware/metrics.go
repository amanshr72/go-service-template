package middleware

import (
	"go-crud2/internal/metrics"
	"net/http"
	"strconv"
	"time"
)

func Metrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		metrics.ActiveRequests.Inc()
		defer metrics.ActiveRequests.Dec()

		start := time.Now()
		wrapped := &responseWriter{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(wrapped, r)

		duration := time.Since(start).Seconds()
		metrics.RequestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration)
		metrics.RequestsTotal.WithLabelValues(r.Method, r.URL.Path, strconv.Itoa(wrapped.status)).Inc()
	})
}
