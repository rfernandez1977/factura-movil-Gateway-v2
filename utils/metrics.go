package utils

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics define las métricas de la aplicación
var Metrics = struct {
	HTTPRequestsTotal   *prometheus.CounterVec
	HTTPRequestDuration *prometheus.HistogramVec
	HTTPRequestSize     *prometheus.HistogramVec
	HTTPResponseSize    *prometheus.HistogramVec
	BoletasCreadas      *prometheus.Counter
	BoletasAnuladas     *prometheus.Counter
	BoletasEnviadas     *prometheus.Counter
	BoletasError        *prometheus.Counter
}{
	HTTPRequestsTotal: promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total de peticiones HTTP",
		},
		[]string{"method", "path", "status"},
	),
	HTTPRequestDuration: promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duración de las peticiones HTTP",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	),
	HTTPRequestSize: promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_size_bytes",
			Help:    "Tamaño de las peticiones HTTP",
			Buckets: prometheus.ExponentialBuckets(100, 10, 8),
		},
		[]string{"method", "path"},
	),
	HTTPResponseSize: promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_size_bytes",
			Help:    "Tamaño de las respuestas HTTP",
			Buckets: prometheus.ExponentialBuckets(100, 10, 8),
		},
		[]string{"method", "path"},
	),
	BoletasCreadas: func() *prometheus.Counter {
		c := promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "boletas_creadas_total",
				Help: "Total de boletas creadas",
			},
		)
		return &c
	}(),
	BoletasAnuladas: func() *prometheus.Counter {
		c := promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "boletas_anuladas_total",
				Help: "Total de boletas anuladas",
			},
		)
		return &c
	}(),
	BoletasEnviadas: func() *prometheus.Counter {
		c := promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "boletas_enviadas_total",
				Help: "Total de boletas enviadas por email",
			},
		)
		return &c
	}(),
	BoletasError: func() *prometheus.Counter {
		c := promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "boletas_error_total",
				Help: "Total de errores en boletas",
			},
		)
		return &c
	}(),
}

// RecordHTTPRequest registra una petición HTTP
func RecordHTTPRequest(method, path string, status int, duration float64, requestSize, responseSize float64) {
	Metrics.HTTPRequestsTotal.WithLabelValues(method, path, fmt.Sprintf("%d", status)).Inc()
	Metrics.HTTPRequestDuration.WithLabelValues(method, path).Observe(duration)
	Metrics.HTTPRequestSize.WithLabelValues(method, path).Observe(requestSize)
	Metrics.HTTPResponseSize.WithLabelValues(method, path).Observe(responseSize)
}

// RecordBoletaCreada registra una boleta creada
func RecordBoletaCreada() {
	(*Metrics.BoletasCreadas).Inc()
}

// RecordBoletaAnulada registra una boleta anulada
func RecordBoletaAnulada() {
	(*Metrics.BoletasAnuladas).Inc()
}

// RecordBoletaEnviada registra una boleta enviada
func RecordBoletaEnviada() {
	(*Metrics.BoletasEnviadas).Inc()
}

// RecordBoletaError registra un error en una boleta
func RecordBoletaError() {
	(*Metrics.BoletasError).Inc()
}
