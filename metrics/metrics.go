package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// DocumentCounter cuenta el número total de documentos procesados
	DocumentCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "document_total",
		Help: "Número total de documentos procesados",
	})

	// DocumentErrors cuenta el número de errores al procesar documentos
	DocumentErrors = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "document_errors_total",
		Help: "Número total de errores al procesar documentos",
	})

	// DocumentProcessingTime mide el tiempo de procesamiento de documentos
	DocumentProcessingTime = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "document_processing_seconds",
		Help:    "Tiempo de procesamiento de documentos en segundos",
		Buckets: prometheus.DefBuckets,
	})
)

func init() {
	prometheus.MustRegister(DocumentCounter)
	prometheus.MustRegister(DocumentErrors)
	prometheus.MustRegister(DocumentProcessingTime)
}

// MetricsHandler retorna el manejador HTTP para las métricas de Prometheus
func MetricsHandler() http.Handler {
	return promhttp.Handler()
}
