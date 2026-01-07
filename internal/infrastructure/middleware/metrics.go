package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/M1ralai/go-modular-monolith-template/internal/infrastructure/metrics"
	"github.com/gorilla/mux"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip metrics for WebSocket to avoid ResponseWriter wrapping
		if r.URL.Path == "/ws" {
			next.ServeHTTP(w, r)
			return
		}

		start := time.Now()
		rw := NewResponseWriter(w)

		next.ServeHTTP(rw, r)

		duration := time.Since(start).Seconds()

		route := mux.CurrentRoute(r)
		path, _ := route.GetPathTemplate()
		if path == "" {
			path = r.URL.Path
		}

		status := strconv.Itoa(rw.statusCode)

		metrics.HttpRequestsTotal.WithLabelValues(r.Method, path, status).Inc()
		metrics.HttpRequestDuration.WithLabelValues(r.Method, path).Observe(duration)
	})
}
