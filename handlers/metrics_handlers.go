package handlers

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/fmgo/api"
)

// MetricsHandlers maneja las métricas
type MetricsHandlers struct {
	client  *api.FacturaMovilClient
	metrics *DocumentMetrics
}

type DocumentMetrics struct {
	documentosEmitidos     *prometheus.CounterVec
	documentosRechazados   *prometheus.CounterVec
	tiempoProcesamientoSII *prometheus.HistogramVec
	foliosDisponibles      *prometheus.GaugeVec
	erroresProcesamiento   *prometheus.CounterVec
	estadosDocumentos      *prometheus.GaugeVec
}

// NewMetricsHandlers crea una nueva instancia de MetricsHandlers
func NewMetricsHandlers(client *api.FacturaMovilClient) *MetricsHandlers {
	metrics := &DocumentMetrics{
		documentosEmitidos: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "documentos_emitidos_total",
				Help: "Total de documentos tributarios emitidos",
			},
			[]string{"tipo_documento"},
		),
		documentosRechazados: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "documentos_rechazados_total",
				Help: "Total de documentos rechazados por el SII",
			},
			[]string{"tipo_documento", "motivo_rechazo"},
		),
		tiempoProcesamientoSII: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "tiempo_procesamiento_sii_segundos",
				Help:    "Tiempo de procesamiento de documentos en el SII",
				Buckets: prometheus.LinearBuckets(0, 30, 10),
			},
			[]string{"tipo_documento"},
		),
		foliosDisponibles: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "folios_disponibles",
				Help: "Cantidad de folios disponibles por tipo de documento",
			},
			[]string{"tipo_documento"},
		),
		erroresProcesamiento: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "errores_procesamiento_total",
				Help: "Total de errores en el procesamiento de documentos",
			},
			[]string{"tipo_error", "tipo_documento"},
		),
		estadosDocumentos: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "estados_documentos",
				Help: "Cantidad de documentos por estado",
			},
			[]string{"tipo_documento", "estado"},
		),
	}

	return &MetricsHandlers{
		client:  client,
		metrics: metrics,
	}
}

// GetMetricsHandler maneja la obtención de métricas
func (h *MetricsHandlers) GetMetricsHandler(c *gin.Context) {
	// Implementación pendiente
	c.JSON(200, gin.H{
		"message": "Metrics handler",
	})
}

func (h *MetricsHandlers) RegistrarEmisionHandler(c *gin.Context) {
	var emision struct {
		TipoDocumento string    `json:"tipoDocumento"`
		Estado        string    `json:"estado"`
		TiempoInicio  time.Time `json:"tiempoInicio"`
	}

	// Registrar métricas de emisión
	h.metrics.documentosEmitidos.WithLabelValues(emision.TipoDocumento).Inc()

	// Registrar tiempo de procesamiento
	tiempoProcesamiento := time.Since(emision.TiempoInicio).Seconds()
	h.metrics.tiempoProcesamientoSII.WithLabelValues(emision.TipoDocumento).Observe(tiempoProcesamiento)

	// Actualizar estado del documento
	h.metrics.estadosDocumentos.WithLabelValues(emision.TipoDocumento, emision.Estado).Inc()

	c.JSON(200, gin.H{
		"mensaje":             "Métricas registradas correctamente",
		"tiempoProcesamiento": tiempoProcesamiento,
	})
}

func (h *MetricsHandlers) ActualizarFoliosHandler(c *gin.Context) {
	var folios struct {
		TipoDocumento string `json:"tipoDocumento"`
		Cantidad      int    `json:"cantidad"`
	}

	// Actualizar métrica de folios disponibles
	h.metrics.foliosDisponibles.WithLabelValues(folios.TipoDocumento).Set(float64(folios.Cantidad))

	c.JSON(200, gin.H{
		"mensaje":           "Folios actualizados",
		"foliosDisponibles": folios.Cantidad,
	})
}

func (h *MetricsHandlers) RegistrarErrorHandler(c *gin.Context) {
	var error struct {
		TipoError     string `json:"tipoError"`
		TipoDocumento string `json:"tipoDocumento"`
		Detalle       string `json:"detalle"`
	}

	// Registrar error en métricas
	h.metrics.erroresProcesamiento.WithLabelValues(error.TipoError, error.TipoDocumento).Inc()

	c.JSON(200, gin.H{
		"mensaje": "Error registrado en métricas",
		"error":   error,
	})
}
