package middlewares

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus"
)

// MetricsMiddleware — обёртка для сбора HTTP-метрик
func MetricsMiddleware(requests *prometheus.CounterVec, duration *prometheus.HistogramVec) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			next.ServeHTTP(wrapped, r)

			routePattern := chi.RouteContext(r.Context()).RoutePattern()
			if routePattern == "" {
				routePattern = "unknown"
			}

			duration.WithLabelValues(r.Method, routePattern).Observe(time.Since(start).Seconds())
			requests.WithLabelValues(r.Method, routePattern, strconv.Itoa(wrapped.statusCode)).Inc()
		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
