package models

import (
	"time"
)

// ReporteAuditoria representa un reporte de auditoría de documentos
type ReporteAuditoria struct {
	ID                 string              `json:"id" bson:"_id"`
	FechaInicio        time.Time           `json:"fecha_inicio" bson:"fecha_inicio"`
	FechaFin           time.Time           `json:"fecha_fin" bson:"fecha_fin"`
	RutEmisor          string              `json:"rut_emisor" bson:"rut_emisor,omitempty"`
	RutReceptor        string              `json:"rut_receptor" bson:"rut_receptor,omitempty"`
	CambiosPorTipo     map[TipoDTE]int     `json:"cambios_por_tipo" bson:"cambios_por_tipo"`
	CambiosPorEstado   map[string]int      `json:"cambios_por_estado" bson:"cambios_por_estado"`
	CambiosPorUsuario  map[string]int      `json:"cambios_por_usuario" bson:"cambios_por_usuario"`
	RegistrosAuditoria []RegistroAuditoria `json:"registros_auditoria" bson:"registros_auditoria"`
	FechaGeneracion    time.Time           `json:"fecha_generacion" bson:"fecha_generacion"`
}

// RegistroAuditoria representa un registro individual de auditoría
type RegistroAuditoria struct {
	ID             string          `json:"id" bson:"_id"`
	DocumentoID    string          `json:"documento_id" bson:"documento_id"`
	TipoDocumento  TipoDTE         `json:"tipo_documento" bson:"tipo_documento"`
	EstadoAnterior string          `json:"estado_anterior" bson:"estado_anterior"`
	EstadoNuevo    string          `json:"estado_nuevo" bson:"estado_nuevo"`
	Usuario        string          `json:"usuario" bson:"usuario"`
	FechaCambio    time.Time       `json:"fecha_cambio" bson:"fecha_cambio"`
	Cambios        []CambioDetalle `json:"cambios" bson:"cambios"`
	IP             string          `json:"ip" bson:"ip"`
	UserAgent      string          `json:"user_agent" bson:"user_agent"`
}

// CambioDetalle representa un cambio específico en un documento
type CambioDetalle struct {
	Campo         string      `json:"campo" bson:"campo"`
	ValorAnterior interface{} `json:"valor_anterior" bson:"valor_anterior"`
	ValorNuevo    interface{} `json:"valor_nuevo" bson:"valor_nuevo"`
	TipoCambio    string      `json:"tipo_cambio" bson:"tipo_cambio"` // CREACION, MODIFICACION, ELIMINACION
}

// ReporteCumplimiento representa un reporte de cumplimiento normativo
type ReporteCumplimiento struct {
	ID                 string               `json:"id" bson:"_id"`
	FechaInicio        time.Time            `json:"fecha_inicio" bson:"fecha_inicio"`
	FechaFin           time.Time            `json:"fecha_fin" bson:"fecha_fin"`
	RutEmisor          string               `json:"rut_emisor" bson:"rut_emisor,omitempty"`
	RutReceptor        string               `json:"rut_receptor" bson:"rut_receptor,omitempty"`
	DocumentosVencidos []DocumentoVencido   `json:"documentos_vencidos" bson:"documentos_vencidos"`
	Alertas            []AlertaCumplimiento `json:"alertas" bson:"alertas"`
	FechaGeneracion    time.Time            `json:"fecha_generacion" bson:"fecha_generacion"`
}

// DocumentoVencido representa un documento que ha excedido su plazo
type DocumentoVencido struct {
	DocumentoID      string    `json:"documento_id" bson:"documento_id"`
	TipoDocumento    TipoDTE   `json:"tipo_documento" bson:"tipo_documento"`
	FechaEmision     time.Time `json:"fecha_emision" bson:"fecha_emision"`
	FechaVencimiento time.Time `json:"fecha_vencimiento" bson:"fecha_vencimiento"`
	DiasVencido      int       `json:"dias_vencido" bson:"dias_vencido"`
	Estado           string    `json:"estado" bson:"estado"`
}

// AlertaCumplimiento representa una alerta de cumplimiento
type AlertaCumplimiento struct {
	Tipo        string    `json:"tipo" bson:"tipo"`
	Descripcion string    `json:"descripcion" bson:"descripcion"`
	Severidad   string    `json:"severidad" bson:"severidad"` // BAJA, MEDIA, ALTA, CRITICA
	FechaAlerta time.Time `json:"fecha_alerta" bson:"fecha_alerta"`
	DocumentoID string    `json:"documento_id" bson:"documento_id"`
	Estado      string    `json:"estado" bson:"estado"` // PENDIENTE, RESUELTA, DESCARTADA
}
