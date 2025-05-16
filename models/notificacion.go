package models

import "time"

// EstadoNotificacion representa el estado de una notificación
type EstadoNotificacion string

// Estados de notificación
const (
	EstadoNotificacionPendiente EstadoNotificacion = "PENDIENTE"
	EstadoNotificacionEnviada   EstadoNotificacion = "ENVIADA"
	EstadoNotificacionEntregada EstadoNotificacion = "ENTREGADA"
	EstadoNotificacionLeida     EstadoNotificacion = "LEIDA"
	EstadoNotificacionError     EstadoNotificacion = "ERROR"
	EstadoNotificacionCancelada EstadoNotificacion = "CANCELADA"
)

// TipoNotificacion representa el tipo de notificación
type TipoNotificacion string

// Tipos de notificación
const (
	NotificacionEmail   TipoNotificacion = "EMAIL"
	NotificacionSMS     TipoNotificacion = "SMS"
	NotificacionPush    TipoNotificacion = "PUSH"
	NotificacionInApp   TipoNotificacion = "IN_APP"
	NotificacionWebhook TipoNotificacion = "WEBHOOK"
)

// Notificacion representa una notificación en el sistema
type Notificacion struct {
	ID             string                 `json:"id" bson:"_id,omitempty"`
	Tipo           TipoNotificacion       `json:"tipo" bson:"tipo"`
	Titulo         string                 `json:"titulo" bson:"titulo"`
	Mensaje        string                 `json:"mensaje" bson:"mensaje"`
	HTML           string                 `json:"html,omitempty" bson:"html,omitempty"`
	Destinatarios  []string               `json:"destinatarios" bson:"destinatarios"`
	Emisor         string                 `json:"emisor" bson:"emisor"`
	FechaCreacion  time.Time              `json:"fecha_creacion" bson:"fecha_creacion"`
	FechaEnvio     time.Time              `json:"fecha_envio,omitempty" bson:"fecha_envio,omitempty"`
	FechaEntrega   time.Time              `json:"fecha_entrega,omitempty" bson:"fecha_entrega,omitempty"`
	FechaLectura   time.Time              `json:"fecha_lectura,omitempty" bson:"fecha_lectura,omitempty"`
	Estado         EstadoNotificacion     `json:"estado" bson:"estado"`
	ErrorMensaje   string                 `json:"error_mensaje,omitempty" bson:"error_mensaje,omitempty"`
	Intentos       int                    `json:"intentos" bson:"intentos"`
	MaxIntentos    int                    `json:"max_intentos" bson:"max_intentos"`
	ProximoIntento time.Time              `json:"proximo_intento,omitempty" bson:"proximo_intento,omitempty"`
	Adjuntos       []AdjuntoNotificacion  `json:"adjuntos,omitempty" bson:"adjuntos,omitempty"`
	Datos          map[string]interface{} `json:"datos,omitempty" bson:"datos,omitempty"`
	EmpresaID      string                 `json:"empresa_id" bson:"empresa_id"`
	DocumentoID    string                 `json:"documento_id,omitempty" bson:"documento_id,omitempty"`
	CreatedAt      time.Time              `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at" bson:"updated_at"`
}

// AdjuntoNotificacion representa un archivo adjunto en una notificación
type AdjuntoNotificacion struct {
	ID            string    `json:"id" bson:"_id,omitempty"`
	Nombre        string    `json:"nombre" bson:"nombre"`
	NombreArchivo string    `json:"nombre_archivo" bson:"nombre_archivo"`
	ContentType   string    `json:"content_type" bson:"content_type"`
	Tamaño        int64     `json:"tamaño" bson:"tamaño"`
	Contenido     []byte    `json:"contenido,omitempty" bson:"contenido,omitempty"`
	URLPublica    string    `json:"url_publica,omitempty" bson:"url_publica,omitempty"`
	CreatedAt     time.Time `json:"created_at" bson:"created_at"`
}

// PlantillaNotificacion representa una plantilla para notificaciones
type PlantillaNotificacion struct {
	ID            string           `json:"id" bson:"_id,omitempty"`
	Nombre        string           `json:"nombre" bson:"nombre"`
	Descripcion   string           `json:"descripcion" bson:"descripcion"`
	Tipo          TipoNotificacion `json:"tipo" bson:"tipo"`
	Asunto        string           `json:"asunto" bson:"asunto"`
	Contenido     string           `json:"contenido" bson:"contenido"`
	ContenidoHTML string           `json:"contenido_html,omitempty" bson:"contenido_html,omitempty"`
	Variables     []string         `json:"variables" bson:"variables"`
	Activa        bool             `json:"activa" bson:"activa"`
	EmpresaID     string           `json:"empresa_id" bson:"empresa_id"`
	CreatedAt     time.Time        `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time        `json:"updated_at" bson:"updated_at"`
}

// NotificacionRequest representa una solicitud para enviar una notificación
type NotificacionRequest struct {
	Tipo          TipoNotificacion       `json:"tipo" binding:"required"`
	Titulo        string                 `json:"titulo" binding:"required"`
	Mensaje       string                 `json:"mensaje" binding:"required"`
	HTML          string                 `json:"html,omitempty"`
	Destinatarios []string               `json:"destinatarios" binding:"required,min=1"`
	Adjuntos      []AdjuntoRequest       `json:"adjuntos,omitempty"`
	Datos         map[string]interface{} `json:"datos,omitempty"`
	EmpresaID     string                 `json:"empresa_id" binding:"required"`
	DocumentoID   string                 `json:"documento_id,omitempty"`
	PlantillaID   string                 `json:"plantilla_id,omitempty"`
	Variables     map[string]string      `json:"variables,omitempty"`
}

// AdjuntoRequest representa una solicitud para adjuntar un archivo a una notificación
type AdjuntoRequest struct {
	Nombre        string `json:"nombre" binding:"required"`
	NombreArchivo string `json:"nombre_archivo" binding:"required"`
	ContentType   string `json:"content_type" binding:"required"`
	Contenido     []byte `json:"contenido" binding:"required"`
}

// NotificacionResponse representa la respuesta a una solicitud de notificación
type NotificacionResponse struct {
	ID            string             `json:"id"`
	Estado        EstadoNotificacion `json:"estado"`
	FechaCreacion time.Time          `json:"fecha_creacion"`
	Mensaje       string             `json:"mensaje"`
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
	if n.Emisor == "" {
		return &ValidationFieldError{Field: "emisor", Message: "El emisor es obligatorio"}
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
	if len(n.Destinatarios) == 0 {
		return &ValidationFieldError{Field: "destinatarios", Message: "Al menos un destinatario es obligatorio"}
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
