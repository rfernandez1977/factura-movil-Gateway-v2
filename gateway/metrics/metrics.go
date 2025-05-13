package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// RequestCounter cuenta las peticiones por endpoint y resultado
	RequestCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total de peticiones HTTP por endpoint y resultado",
		},
		[]string{"endpoint", "result"},
	)

	// RequestDuration mide la duración de las peticiones
	RequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duración de las peticiones HTTP en segundos",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"endpoint"},
	)
)
