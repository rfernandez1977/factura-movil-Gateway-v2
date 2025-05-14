package models

import (
	"fmt"
	"time"
)

// ReporteDocumentosEstado representa un reporte de documentos por estado
type ReporteDocumentosEstado struct {
	ID               string                  `json:"id" bson:"_id"`
	FechaInicio      time.Time               `json:"fecha_inicio" bson:"fecha_inicio"`
	FechaFin         time.Time               `json:"fecha_fin" bson:"fecha_fin"`
	RutEmisor        string                  `json:"rut_emisor" bson:"rut_emisor,omitempty"`
	RutReceptor      string                  `json:"rut_receptor" bson:"rut_receptor,omitempty"`
	TotalesPorEstado map[EstadoDocumento]int `json:"totales_por_estado" bson:"totales_por_estado"`
	TotalesPorTipo   map[TipoDTE]int         `json:"totales_por_tipo" bson:"totales_por_tipo"`
	Documentos       []DocumentoTributario   `json:"documentos" bson:"documentos"`
	FechaGeneracion  time.Time               `json:"fecha_generacion" bson:"fecha_generacion"`
}

// ReporteRechazos representa un reporte de análisis de rechazos
type ReporteRechazos struct {
	ID                   string               `json:"id" bson:"_id"`
	FechaInicio          time.Time            `json:"fecha_inicio" bson:"fecha_inicio"`
	FechaFin             time.Time            `json:"fecha_fin" bson:"fecha_fin"`
	RutEmisor            string               `json:"rut_emisor" bson:"rut_emisor,omitempty"`
	RutReceptor          string               `json:"rut_receptor" bson:"rut_receptor,omitempty"`
	TotalRechazos        int                  `json:"total_rechazos" bson:"total_rechazos"`
	ErroresComunes       []ErrorFrecuente     `json:"errores_comunes" bson:"errores_comunes"`
	DocumentosRechazados []DocumentoRechazado `json:"documentos_rechazados" bson:"documentos_rechazados"`
	FechaGeneracion      time.Time            `json:"fecha_generacion" bson:"fecha_generacion"`
}

// ErrorFrecuente representa un error común en los rechazos
type ErrorFrecuente struct {
	Codigo      string  `json:"codigo" bson:"codigo"`
	Descripcion string  `json:"descripcion" bson:"descripcion"`
	Frecuencia  int     `json:"frecuencia" bson:"frecuencia"`
	Porcentaje  float64 `json:"porcentaje" bson:"porcentaje"`
}

// DocumentoRechazado representa un documento que fue rechazado
type DocumentoRechazado struct {
	ID            string     `json:"id" bson:"_id"`
	TipoDocumento string     `json:"tipo_documento" bson:"tipo_documento"`
	Folio         int64      `json:"folio" bson:"folio"`
	FechaEmision  time.Time  `json:"fecha_emision" bson:"fecha_emision"`
	RutEmisor     string     `json:"rut_emisor" bson:"rut_emisor"`
	RutReceptor   string     `json:"rut_receptor" bson:"rut_receptor"`
	MontoTotal    float64    `json:"monto_total" bson:"monto_total"`
	FechaRechazo  time.Time  `json:"fecha_rechazo" bson:"fecha_rechazo"`
	MotivoRechazo string     `json:"motivo_rechazo" bson:"motivo_rechazo"`
	Errores       []ErrorSII `json:"errores" bson:"errores"`
	EstadoActual  string     `json:"estado_actual" bson:"estado_actual"`
}

// ReporteMetricasRendimiento representa un reporte de métricas de rendimiento
type ReporteMetricasRendimiento struct {
	ID              string              `json:"id" bson:"_id"`
	FechaInicio     time.Time           `json:"fecha_inicio" bson:"fecha_inicio"`
	FechaFin        time.Time           `json:"fecha_fin" bson:"fecha_fin"`
	RutEmisor       string              `json:"rut_emisor" bson:"rut_emisor,omitempty"`
	Metricas        MetricasRendimiento `json:"metricas" bson:"metricas"`
	FechaGeneracion time.Time           `json:"fecha_generacion" bson:"fecha_generacion"`
}

// MetricasRendimiento contiene las métricas de rendimiento
type MetricasRendimiento struct {
	TotalDocumentos         int           `json:"total_documentos" bson:"total_documentos"`
	DocumentosAceptados     int           `json:"documentos_aceptados" bson:"documentos_aceptados"`
	DocumentosRechazados    int           `json:"documentos_rechazados" bson:"documentos_rechazados"`
	DocumentosPendientes    int           `json:"documentos_pendientes" bson:"documentos_pendientes"`
	TasaAceptacion          float64       `json:"tasa_aceptacion" bson:"tasa_aceptacion"`
	TiempoPromedioRespuesta time.Duration `json:"tiempo_promedio_respuesta" bson:"tiempo_promedio_respuesta"`
	DocumentosPorHora       float64       `json:"documentos_por_hora" bson:"documentos_por_hora"`
	ErroresPorDocumento     float64       `json:"errores_por_documento" bson:"errores_por_documento"`
}

// ReporteTributario representa un reporte tributario
type ReporteTributario struct {
	ID                 string                `json:"id" bson:"_id"`
	FechaInicio        time.Time             `json:"fecha_inicio" bson:"fecha_inicio"`
	FechaFin           time.Time             `json:"fecha_fin" bson:"fecha_fin"`
	RutEmisor          string                `json:"rut_emisor" bson:"rut_emisor,omitempty"`
	RutReceptor        string                `json:"rut_receptor" bson:"rut_receptor,omitempty"`
	TotalesTributarios TotalesTributarios    `json:"totales_tributarios" bson:"totales_tributarios"`
	Documentos         []DocumentoTributario `json:"documentos" bson:"documentos"`
	FechaGeneracion    time.Time             `json:"fecha_generacion" bson:"fecha_generacion"`
}

// TotalesTributarios contiene los totales para el reporte tributario
type TotalesTributarios struct {
	MontoNetoTotal      float64                 `json:"monto_neto_total" bson:"monto_neto_total"`
	MontoIVATotal       float64                 `json:"monto_iva_total" bson:"monto_iva_total"`
	MontoRetencionTotal float64                 `json:"monto_retencion_total" bson:"monto_retencion_total"`
	MontoTotal          float64                 `json:"monto_total" bson:"monto_total"`
	TotalesPorTipo      map[TipoDTE]TotalesTipo `json:"totales_por_tipo" bson:"totales_por_tipo"`
}

// TotalesTipo contiene los totales por tipo de documento
type TotalesTipo struct {
	Cantidad       int     `json:"cantidad" bson:"cantidad"`
	MontoNeto      float64 `json:"monto_neto" bson:"monto_neto"`
	MontoIVA       float64 `json:"monto_iva" bson:"monto_iva"`
	MontoRetencion float64 `json:"monto_retencion" bson:"monto_retencion"`
	MontoTotal     float64 `json:"monto_total" bson:"monto_total"`
}

type SyncRecord struct {
	ID        int64  `json:"id"`
	Type      string `json:"type"`
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
}

// ReporteGenerator es el generador de reportes
type ReporteGenerator struct {
	client     ClienteReporte
	documentos []DocumentoTributario
	pdfs       []string
	errores    []*ErrorValidacion
	estados    []*EstadoSII
	outputDir  string
}

// ClienteReporte define la interfaz para el cliente de reportes
type ClienteReporte interface {
	ConsultarEstadoSII(id string) (*EstadoSII, error)
	DescargarPDF(id string, rutaDestino string) error
}

// ReporteConfig es la configuración para generar reportes
type ReporteConfig struct {
	FechaInicio   time.Time
	FechaFin      time.Time
	TipoDocumento string
	RutEmisor     string
	RutReceptor   string
	OutputDir     string
}

// ErrorReporte representa un error en la generación de reportes
type ErrorReporte struct {
	Codigo    string `json:"codigo"`
	Mensaje   string `json:"mensaje"`
	Campo     string `json:"campo,omitempty"`
	Valor     string `json:"valor,omitempty"`
	Timestamp string `json:"timestamp"`
}

// ObtenerPDFs obtiene los PDFs generados para los documentos en el rango de fechas
func (s *ReporteGenerator) ObtenerPDFs(config *ReporteConfig) error {
	// Para cada documento, obtener su PDF
	for _, doc := range s.documentos {
		if err := s.client.DescargarPDF(doc.ID, fmt.Sprintf("%s/%s.pdf", s.outputDir, doc.ID)); err != nil {
			s.errores = append(s.errores, &ErrorValidacion{
				Codigo:    "REPORTE-002",
				Mensaje:   fmt.Sprintf("Error al descargar PDF: %v", err),
				Campo:     "documento_id",
				Valor:     doc.ID,
				Timestamp: time.Now().Format(time.RFC3339),
			})
			continue
		}
		s.pdfs = append(s.pdfs, fmt.Sprintf("%s/%s.pdf", s.outputDir, doc.ID))
	}
	return nil
}
