package models

import "time"

// TipoNotificacion define los tipos de notificaciones
type TipoNotificacion string

const (
	TipoNotificacionAlertaDocumento  TipoNotificacion = "ALERTA_DOCUMENTO"
	TipoNotificacionEstadoDocumento  TipoNotificacion = "ESTADO_DOCUMENTO"
	TipoNotificacionVencimiento      TipoNotificacion = "VENCIMIENTO"
	TipoNotificacionErrorSistema     TipoNotificacion = "ERROR_SISTEMA"
	TipoNotificacionActualizacionSII TipoNotificacion = "ACTUALIZACION_SII"
)

// Estados de las notificaciones
const (
	EstadoPendiente = "PENDIENTE"
	EstadoEnviada   = "ENVIADA"
	EstadoLeida     = "LEIDA"
	EstadoError     = "ERROR"
)

// Notificacion representa una notificación en el sistema
type Notificacion struct {
	ID            string                 `json:"id" bson:"_id,omitempty"`
	UsuarioID     string                 `json:"usuario_id" bson:"usuario_id"`
	Tipo          TipoNotificacion       `json:"tipo" bson:"tipo"`
	Titulo        string                 `json:"titulo" bson:"titulo"`
	Mensaje       string                 `json:"mensaje" bson:"mensaje"`
	Detalles      string                 `json:"detalles,omitempty" bson:"detalles,omitempty"`
	Data          map[string]interface{} `json:"data,omitempty" bson:"data,omitempty"`
	URL           string                 `json:"url,omitempty" bson:"url,omitempty"`
	Estado        string                 `json:"estado" bson:"estado"`
	IntentosEnvio int                    `json:"intentos_envio" bson:"intentos_envio"`
	FechaCreacion time.Time              `json:"fecha_creacion" bson:"fecha_creacion"`
	FechaEnvio    time.Time              `json:"fecha_envio,omitempty" bson:"fecha_envio,omitempty"`
	FechaLeida    time.Time              `json:"fecha_leida,omitempty" bson:"fecha_leida,omitempty"`
	Error         string                 `json:"error,omitempty" bson:"error,omitempty"`
}

// NotificacionRequest representa la solicitud para crear una notificación
type NotificacionRequest struct {
	Tipo         string                 `json:"tipo" binding:"required"`
	Titulo       string                 `json:"titulo" binding:"required"`
	Mensaje      string                 `json:"mensaje" binding:"required"`
	Destinatario string                 `json:"destinatario" binding:"required"`
	Prioridad    string                 `json:"prioridad" binding:"required"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

// NotificacionResponse representa la respuesta de una notificación
type NotificacionResponse struct {
	ID            string                 `json:"id"`
	Tipo          string                 `json:"tipo"`
	Titulo        string                 `json:"titulo"`
	Mensaje       string                 `json:"mensaje"`
	Destinatario  string                 `json:"destinatario"`
	FechaCreacion time.Time              `json:"fecha_creacion"`
	Estado        string                 `json:"estado"`
	Prioridad     string                 `json:"prioridad"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
	Leida         bool                   `json:"leida"`
	FechaLectura  time.Time              `json:"fecha_lectura,omitempty"`
}

// HorarioNotificacion representa un rango horario para recibir notificaciones
type HorarioNotificacion struct {
	Inicio time.Time `json:"inicio" bson:"inicio"`
	Fin    time.Time `json:"fin" bson:"fin"`
}

// PreferenciasNotificacion representa las preferencias de notificación de un usuario
type PreferenciasNotificacion struct {
	UsuarioID           string              `json:"usuario_id" bson:"_id,omitempty"`
	Email               string              `json:"email" bson:"email"`
	Telefono            string              `json:"telefono" bson:"telefono"`
	TokensPush          []string            `json:"tokens_push" bson:"tokens_push"`
	TiposNotificacion   []TipoNotificacion  `json:"tipos_notificacion" bson:"tipos_notificacion"`
	RecibirEmail        bool                `json:"recibir_email" bson:"recibir_email"`
	RecibirSMS          bool                `json:"recibir_sms" bson:"recibir_sms"`
	RecibirPush         bool                `json:"recibir_push" bson:"recibir_push"`
	HorarioNotificacion HorarioNotificacion `json:"horario_notificacion" bson:"horario_notificacion"`
	DiasNotificacion    []int               `json:"dias_notificacion" bson:"dias_notificacion"`
	FechaActualizacion  time.Time           `json:"fecha_actualizacion" bson:"fecha_actualizacion"`
}

// Validate valida que todos los campos obligatorios estén presentes
func (n *Notificacion) Validate() error {
	if n.UsuarioID == "" {
		return &ValidationFieldError{Field: "usuario_id", Message: "El ID del usuario es obligatorio"}
	}
	if n.Tipo == "" {
		return &ValidationFieldError{Field: "tipo", Message: "El tipo de notificación es obligatorio"}
	}
	if n.Titulo == "" {
		return &ValidationFieldError{Field: "titulo", Message: "El título es obligatorio"}
	}
	if n.Mensaje == "" {
		return &ValidationFieldError{Field: "mensaje", Message: "El mensaje es obligatorio"}
	}
	return nil
}

// Validate valida que todos los campos obligatorios estén presentes
func (p *PreferenciasNotificacion) Validate() error {
	if p.UsuarioID == "" {
		return &ValidationFieldError{Field: "usuario_id", Message: "El ID del usuario es obligatorio"}
	}
	if !p.RecibirEmail && !p.RecibirSMS && !p.RecibirPush {
		return &ValidationFieldError{Field: "recibir_*", Message: "Al menos un método de notificación debe estar activo"}
	}
	if p.RecibirEmail && p.Email == "" {
		return &ValidationFieldError{Field: "email", Message: "El email es obligatorio si se reciben notificaciones por email"}
	}
	if p.RecibirSMS && p.Telefono == "" {
		return &ValidationFieldError{Field: "telefono", Message: "El teléfono es obligatorio si se reciben notificaciones por SMS"}
	}
	if p.RecibirPush && len(p.TokensPush) == 0 {
		return &ValidationFieldError{Field: "tokens_push", Message: "Al menos un token push es obligatorio si se reciben notificaciones push"}
	}
	return nil
}
