package middleware

import (
	"net/http"
	"time"

	"github.com/suryamp/receipt-processor/metrics"
)

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		handler := r.URL.Path

		metrics.RequestsTotal.WithLabelValues(handler).Inc()

		next.ServeHTTP(w, r)

		duration := time.Since(start).Seconds()
		metrics.RequestDuration.WithLabelValues(handler).Observe(duration)
	})
}
