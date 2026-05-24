package server

import (
	"encoding/json"
	"net/http"
	"sync/atomic"
	"time"
)

// Metrics collects request-level statistics for the file server.
// All fields use atomic operations — no mutex needed.
type Metrics struct {
	totalRequests   atomic.Int64
	activeRequests  atomic.Int64
	cacheHits       atomic.Int64
	cacheMisses     atomic.Int64
	totalResponseNs atomic.Int64
}

// NewMetrics returns a zero-valued Metrics ready for use.
func NewMetrics() *Metrics {
	return &Metrics{}
}

// RecordStart increments the active request counter and returns a function
// that should be deferred to record the request completion.
func (m *Metrics) RecordStart() func(cacheHit bool) {
	m.activeRequests.Add(1)
	m.totalRequests.Add(1)
	start := time.Now()

	return func(cacheHit bool) {
		m.activeRequests.Add(-1)
		m.totalResponseNs.Add(time.Since(start).Nanoseconds())
		if cacheHit {
			m.cacheHits.Add(1)
		} else {
			m.cacheMisses.Add(1)
		}
	}
}

// metricsResponse is the JSON structure returned by GET /metrics.
type metricsResponse struct {
	TotalRequests     int64   `json:"total_requests"`
	ActiveRequests    int64   `json:"active_requests"`
	CacheHits         int64   `json:"cache_hits"`
	CacheMisses       int64   `json:"cache_misses"`
	CacheHitRate      float64 `json:"cache_hit_rate"`
	AvgResponseTimeMs float64 `json:"avg_response_time_ms"`
}

// ServeHTTP writes the current metrics as JSON.
func (m *Metrics) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	total := m.totalRequests.Load()
	hits := m.cacheHits.Load()
	misses := m.cacheMisses.Load()
	totalNs := m.totalResponseNs.Load()

	var hitRate float64
	if hits+misses > 0 {
		hitRate = float64(hits) / float64(hits+misses)
	}

	var avgMs float64
	if total > 0 {
		avgMs = float64(totalNs) / float64(total) / 1e6
	}

	resp := metricsResponse{
		TotalRequests:     total,
		ActiveRequests:    m.activeRequests.Load(),
		CacheHits:         hits,
		CacheMisses:       misses,
		CacheHitRate:      hitRate,
		AvgResponseTimeMs: avgMs,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
