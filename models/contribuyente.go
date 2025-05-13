package models

import (
	"time"
)

// InformacionContribuyente representa la información de un contribuyente en el SII
type InformacionContribuyente struct {
	RUT              string    `json:"rut"`
	RazonSocial      string    `json:"razon_social"`
	NombreFantasia   string    `json:"nombre_fantasia,omitempty"`
	Direccion        string    `json:"direccion"`
	Comuna           string    `json:"comuna"`
	Ciudad           string    `json:"ciudad"`
	Region           string    `json:"region"`
	Giro             string    `json:"giro"`
	FechaInicio      time.Time `json:"fecha_inicio,omitempty"`
	EstadoActividad  string    `json:"estado_actividad"`
	CategoriaEmpresa string    `json:"categoria_empresa,omitempty"`
	Acteco           []string  `json:"acteco,omitempty"`
	EmailDTE         string    `json:"email_dte,omitempty"`
	FechaResolucion  time.Time `json:"fecha_resolucion,omitempty"`
	NumeroResolucion int       `json:"numero_resolucion,omitempty"`
	FechaConsulta    time.Time `json:"fecha_consulta"`
}

// EstadoContribuyente representa el estado de un contribuyente en el SII
type EstadoContribuyente struct {
	RUT                       string    `json:"rut"`
	Estado                    string    `json:"estado"`
	Glosa                     string    `json:"glosa"`
	FechaConsulta             time.Time `json:"fecha_consulta"`
	AutorizadoEmitirDTE       bool      `json:"autorizado_emitir_dte"`
	AutorizadoRecibirDTE      bool      `json:"autorizado_recibir_dte"`
	TiposDocumentoAutorizados []string  `json:"tipos_documento_autorizados,omitempty"`
	FechaAutorizacion         time.Time `json:"fecha_autorizacion,omitempty"`
	NumeroResolucion          int       `json:"numero_resolucion,omitempty"`
}

// ResumenContribuyente representa un resumen de la información del contribuyente
type ResumenContribuyente struct {
	RUT                 string    `json:"rut"`
	RazonSocial         string    `json:"razon_social"`
	Giro                string    `json:"giro"`
	EstadoActividad     string    `json:"estado_actividad"`
	AutorizadoDTE       bool      `json:"autorizado_dte"`
	FechaConsulta       time.Time `json:"fecha_consulta"`
	UltimaActualizacion time.Time `json:"ultima_actualizacion,omitempty"`
}

// SIIResponseBasic representa una respuesta genérica del SII
type SIIResponseBasic struct {
	Codigo           int                    `json:"codigo"`
	Mensaje          string                 `json:"mensaje"`
	Detalle          string                 `json:"detalle,omitempty"`
	Timestamp        time.Time              `json:"timestamp"`
	DatosAdicionales map[string]interface{} `json:"datos_adicionales,omitempty"`
}
