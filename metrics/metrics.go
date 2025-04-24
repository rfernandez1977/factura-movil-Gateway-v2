package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/promhttp"
)

var (
	// RequestCounter counts the number of requests per endpoint and status
	RequestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gateway_requests_total",
			Help: "Total number of requests to Gateway endpoints",
		},
		[]string{"endpoint", "status"},
	)
)

// InitMetrics registers Prometheus metrics
func InitMetrics() {
	prometheus.MustRegister(RequestCounter)
}

// MetricsHandler exposes the Prometheus metrics endpoint
func MetricsHandler() promhttp.Handler {
	return promhttp.Handler()
}