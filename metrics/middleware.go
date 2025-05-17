package metrics

import (
	"net/http"
	"runtime"
	"time"
)

// MetricsMiddleware es un middleware que registra métricas de las peticiones HTTP
func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap ResponseWriter para capturar el código de estado
		wrapped := &responseWriter{ResponseWriter: w, status: http.StatusOK}

		// Ejecutar el siguiente handler
		next.ServeHTTP(wrapped, r)

		// Registrar duración
		duration := time.Since(start).Seconds()

		// Actualizar métricas
		APIRequestDuration.WithLabelValues(
			r.Method,
			r.URL.Path,
			http.StatusText(wrapped.status),
		).Observe(duration)

		APIRequestsTotal.WithLabelValues(
			r.Method,
			r.URL.Path,
			http.StatusText(wrapped.status),
		).Inc()

		// Actualizar métricas de sistema
		GoroutinesGauge.Set(float64(runtime.NumGoroutine()))

		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		MemoryUsage.Set(float64(m.Alloc))
	})
}

// responseWriter es un wrapper para http.ResponseWriter que captura el código de estado
type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

// CacheMetricsMiddleware es un middleware para registrar métricas de caché
func CacheMetricsMiddleware(operation func() error) error {
	start := time.Now()
	err := operation()
	duration := time.Since(start).Seconds()

	CacheLatency.Observe(duration)

	if err != nil {
		CacheMisses.Inc()
	} else {
		CacheHits.Inc()
	}

	return err
}

// SIIMetricsMiddleware es un middleware para registrar métricas de peticiones al SII
func SIIMetricsMiddleware(operation string, fn func() error) error {
	start := time.Now()
	err := fn()
	duration := time.Since(start).Seconds()

	SIIRequestDuration.Observe(duration)

	status := "success"
	if err != nil {
		status = "error"
	}

	SIIRequestsTotal.WithLabelValues(operation, status).Inc()

	return err
}
