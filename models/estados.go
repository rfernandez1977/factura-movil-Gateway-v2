package models

import "time"

// EstadoRecepcion representa el estado de recepción de un documento
type EstadoRecepcion string

const (
	EstadoRecepcionPendiente  EstadoRecepcion = "PENDIENTE"
	EstadoRecepcionRecibido   EstadoRecepcion = "RECIBIDO"
	EstadoRecepcionRechazado  EstadoRecepcion = "RECHAZADO"
	EstadoRecepcionAceptado   EstadoRecepcion = "ACEPTADO"
	EstadoRecepcionConReparos EstadoRecepcion = "CON_REPAROS"
	EstadoRecepcionNoRecibido EstadoRecepcion = "NO_RECIBIDO"
	EstadoRecepcionEnProceso  EstadoRecepcion = "EN_PROCESO"
	EstadoRecepcionError      EstadoRecepcion = "ERROR"
)

// EstadoSII representa el estado de un documento en el SII
type EstadoSII struct {
	Codigo      int       `json:"codigo"`
	Descripcion string    `json:"descripcion"`
	Timestamp   time.Time `json:"timestamp"`
	// Campos adicionales que son utilizados por otros paquetes
	Estado  string     `json:"estado,omitempty"`
	Glosa   string     `json:"glosa,omitempty"`
	TrackID string     `json:"track_id,omitempty"`
	Fecha   time.Time  `json:"fecha,omitempty"`
	Detalle string     `json:"detalle,omitempty"`
	Errores []ErrorSII `json:"errores,omitempty"`
}

// EstadoDocumento representa el estado de un documento en el sistema
type EstadoDocumento string

const (
	EstadoDocumentoPendiente    EstadoDocumento = "PENDIENTE"
	EstadoDocumentoProcesando   EstadoDocumento = "PROCESANDO"
	EstadoDocumentoCompletado   EstadoDocumento = "COMPLETADO"
	EstadoDocumentoError        EstadoDocumento = "ERROR"
	EstadoDocumentoRechazado    EstadoDocumento = "RECHAZADO"
	EstadoDocumentoAnulado      EstadoDocumento = "ANULADO"
	EstadoDocumentoEnviado      EstadoDocumento = "ENVIADO"
	EstadoDocumentoAceptado     EstadoDocumento = "ACEPTADO"
	EstadoDocumentoRechazadoSII EstadoDocumento = "RECHAZADO_SII"
)

// EstadoFlujo representa el estado de un flujo de trabajo
type EstadoFlujo string

const (
	EstadoFlujoPendiente  EstadoFlujo = "PENDIENTE"
	EstadoFlujoProcesando EstadoFlujo = "PROCESANDO"
	EstadoFlujoCompletado EstadoFlujo = "COMPLETADO"
	EstadoFlujoError      EstadoFlujo = "ERROR"
	EstadoFlujoCancelado  EstadoFlujo = "CANCELADO"
	EstadoFlujoPausado    EstadoFlujo = "PAUSADO"
	EstadoFlujoReanudado  EstadoFlujo = "REANUDADO"
)

// EstadoPaso representa el estado de un paso en un flujo de trabajo
type EstadoPaso string

const (
	EstadoPasoPendiente  EstadoPaso = "PENDIENTE"
	EstadoPasoProcesando EstadoPaso = "PROCESANDO"
	EstadoPasoCompletado EstadoPaso = "COMPLETADO"
	EstadoPasoError      EstadoPaso = "ERROR"
	EstadoPasoCancelado  EstadoPaso = "CANCELADO"
	EstadoPasoPausado    EstadoPaso = "PAUSADO"
	EstadoPasoReanudado  EstadoPaso = "REANUDADO"
)

// EstadoNotificacion representa el estado de una notificación
type EstadoNotificacion string

const (
	EstadoNotificacionPendiente EstadoNotificacion = "PENDIENTE"
	EstadoNotificacionEnviada   EstadoNotificacion = "ENVIADA"
	EstadoNotificacionEntregada EstadoNotificacion = "ENTREGADA"
	EstadoNotificacionLeida     EstadoNotificacion = "LEIDA"
	EstadoNotificacionError     EstadoNotificacion = "ERROR"
	EstadoNotificacionCancelada EstadoNotificacion = "CANCELADA"
)

// EstadoIntegracionERP representa el estado de una integración con un ERP
type EstadoIntegracionERP string

const (
	EstadoIntegracionERPPendiente  EstadoIntegracionERP = "PENDIENTE"
	EstadoIntegracionERPProcesando EstadoIntegracionERP = "PROCESANDO"
	EstadoIntegracionERPCompletado EstadoIntegracionERP = "COMPLETADO"
	EstadoIntegracionERPError      EstadoIntegracionERP = "ERROR"
	EstadoIntegracionERPCancelado  EstadoIntegracionERP = "CANCELADO"
)

// EstadoSesion representa el estado de una sesión de usuario
type EstadoSesion string

const (
	EstadoSesionActiva    EstadoSesion = "ACTIVA"
	EstadoSesionExpirada  EstadoSesion = "EXPIRADA"
	EstadoSesionCerrada   EstadoSesion = "CERRADA"
	EstadoSesionBloqueada EstadoSesion = "BLOQUEADA"
)

// EstadoCAF representa el estado de un CAF (Código de Autorización de Folios)
type EstadoCAF string

const (
	EstadoCAFActivo    EstadoCAF = "ACTIVO"
	EstadoCAFAgotado   EstadoCAF = "AGOTADO"
	EstadoCAFExpirado  EstadoCAF = "EXPIRADO"
	EstadoCAFAnulado   EstadoCAF = "ANULADO"
	EstadoCAFPendiente EstadoCAF = "PENDIENTE"
)

// EstadoDTE representa el estado de un DTE en el sistema
type EstadoDTE string

const (
	EstadoDTEPendiente    EstadoDTE = "PENDIENTE"
	EstadoDTEProcesando   EstadoDTE = "PROCESANDO"
	EstadoDTECompletado   EstadoDTE = "COMPLETADO"
	EstadoDTEError        EstadoDTE = "ERROR"
	EstadoDTERechazado    EstadoDTE = "RECHAZADO"
	EstadoDTEAnulado      EstadoDTE = "ANULADO"
	EstadoDTEEnviado      EstadoDTE = "ENVIADO"
	EstadoDTEAceptado     EstadoDTE = "ACEPTADO"
	EstadoDTERechazadoSII EstadoDTE = "RECHAZADO_SII"
)

// EstadoSIIType representa el estado de un documento en el SII
type EstadoSIIType string

const (
	EstadoSIIAceptado   EstadoSIIType = "ACEPTADO"
	EstadoSIIRechazado  EstadoSIIType = "RECHAZADO"
	EstadoSIIPendiente  EstadoSIIType = "PENDIENTE"
	EstadoSIIEnProceso  EstadoSIIType = "EN_PROCESO"
	EstadoSIIRecibido   EstadoSIIType = "RECIBIDO"
	EstadoSIINoRecibido EstadoSIIType = "NO_RECIBIDO"
	EstadoSIIReparo     EstadoSIIType = "CON_REPAROS"
	EstadoSIIError      EstadoSIIType = "ERROR"
)

// Funciones auxiliares para crear estados
func NewEstadoSII(codigo int, descripcion string) *EstadoSII {
	return &EstadoSII{
		Codigo:      codigo,
		Descripcion: descripcion,
		Timestamp:   time.Now(),
	}
}

// RespuestaSII representa la respuesta recibida del SII
type RespuestaSII struct {
	Estado       string     `json:"estado" xml:"ESTADO"`
	Glosa        string     `json:"glosa" xml:"GLOSA"`
	TrackID      string     `json:"track_id" xml:"TRACKID"`
	FechaProceso time.Time  `json:"fecha_proceso" xml:"FECHA_PROCESO"`
	Errores      []ErrorSII `json:"errores,omitempty" xml:"ERRORES>ERROR,omitempty"`
}

// TieneErrores verifica si la respuesta contiene errores
func (r *RespuestaSII) TieneErrores() bool {
	return r.Estado == "ERROR" || r.Estado == "RECHAZADO" || len(r.Errores) > 0
}
