package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// Métricas de documentos
	DocumentCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "fmgo_document_total",
		Help: "Número total de documentos procesados",
	})

	DocumentErrors = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "fmgo_document_errors_total",
		Help: "Número total de errores al procesar documentos",
	})

	DocumentProcessingTime = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "fmgo_document_processing_seconds",
		Help:    "Tiempo de procesamiento de documentos en segundos",
		Buckets: prometheus.DefBuckets,
	})

	// Métricas de caché Redis
	CacheHits = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "fmgo_cache_hits_total",
		Help: "Número total de aciertos en caché",
	})

	CacheMisses = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "fmgo_cache_misses_total",
		Help: "Número total de fallos en caché",
	})

	CacheLatency = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "fmgo_cache_operation_seconds",
		Help:    "Tiempo de operaciones de caché en segundos",
		Buckets: []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1},
	})

	// Métricas de API
	APIRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "fmgo_api_request_duration_seconds",
			Help:    "Duración de las peticiones API en segundos",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint", "status"},
	)

	APIRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "fmgo_api_requests_total",
			Help: "Número total de peticiones API",
		},
		[]string{"method", "endpoint", "status"},
	)

	// Métricas de sistema
	GoroutinesGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "fmgo_goroutines",
		Help: "Número actual de goroutines",
	})

	MemoryUsage = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "fmgo_memory_bytes",
		Help: "Uso de memoria en bytes",
	})

	// Métricas de SII
	SIIRequestDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "fmgo_sii_request_duration_seconds",
		Help:    "Duración de las peticiones al SII en segundos",
		Buckets: prometheus.DefBuckets,
	})

	SIIRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "fmgo_sii_requests_total",
			Help: "Número total de peticiones al SII",
		},
		[]string{"operation", "status"},
	)
)

func init() {
	// Registro de métricas de documentos
	prometheus.MustRegister(DocumentCounter)
	prometheus.MustRegister(DocumentErrors)
	prometheus.MustRegister(DocumentProcessingTime)

	// Registro de métricas de caché
	prometheus.MustRegister(CacheHits)
	prometheus.MustRegister(CacheMisses)
	prometheus.MustRegister(CacheLatency)

	// Registro de métricas de API
	prometheus.MustRegister(APIRequestDuration)
	prometheus.MustRegister(APIRequestsTotal)

	// Registro de métricas de sistema
	prometheus.MustRegister(GoroutinesGauge)
	prometheus.MustRegister(MemoryUsage)

	// Registro de métricas de SII
	prometheus.MustRegister(SIIRequestDuration)
	prometheus.MustRegister(SIIRequestsTotal)
}

// MetricsHandler retorna el manejador HTTP para las métricas de Prometheus
func MetricsHandler() http.Handler {
	return promhttp.Handler()
}
