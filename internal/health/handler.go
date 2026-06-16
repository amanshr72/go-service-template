package health

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"runtime"
	"sync/atomic"
	"time"
)

var (
	startTime    = time.Now()
	requestCount int64
)

func IncrementRequest() {
	atomic.AddInt64(&requestCount, 1)
}

type response struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
}

func RegisterRoutes(mux *http.ServeMux, db *sql.DB) {
	mux.HandleFunc("GET /health", liveness())
	mux.HandleFunc("GET /health/ready", readiness(db))
	mux.HandleFunc("GET /health/metrics", metrics())
}

// Health godoc
// @Summary Health check
// @Description Liveness probe
// @Tags health
// @Produce json
// @Success 200 {object} response
// @Router /health [get]
func liveness() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, response{
			Status:    "ok",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		})
	}
}

// Readiness godoc
// @Summary Readiness check
// @Description Verifies service dependencies are available
// @Tags health
// @Produce json
// @Success 200 {object} response
// @Failure 503 {object} map[string]string
// @Router /health/ready [get]
func readiness(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if db != nil {
			if err := db.Ping(); err != nil {
				writeJSON(w, http.StatusServiceUnavailable, map[string]string{
					"status": "unavailable",
					"reason": "db ping failed: " + err.Error(),
				})
				return
			}
		}
		writeJSON(w, http.StatusOK, response{Status: "ready", Timestamp: time.Now().UTC().Format(time.RFC3339)})
	}
}

// Metrics godoc
// @Summary Runtime metrics
// @Description Returns uptime, goroutines, memory stats and request count
// @Tags health
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /metrics [get]
func metrics() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var mem runtime.MemStats
		runtime.ReadMemStats(&mem)

		writeJSON(w, http.StatusOK, map[string]interface{}{
			"uptime_seconds": time.Since(startTime).Seconds(),
			"goroutines":     runtime.NumGoroutine(),
			"request_count":  atomic.LoadInt64(&requestCount),
			"cpu_count":      runtime.NumCPU(),
			"go_version":     runtime.Version(),
			"memory": map[string]interface{}{
				"alloc_mb":       bToMb(mem.Alloc),
				"total_alloc_mb": bToMb(mem.TotalAlloc),
				"sys_mb":         bToMb(mem.Sys),
				"gc_cycles":      mem.NumGC,
			},
		})
	}
}

func bToMb(b uint64) float64 {
	return float64(b) / 1024 / 1024
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
