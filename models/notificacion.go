package models

import "time"

// Notificacion representa una notificación en el sistema
type Notificacion struct {
	ID            string                 `json:"id" bson:"_id,omitempty"`
	Tipo          string                 `json:"tipo" bson:"tipo"`
	Titulo        string                 `json:"titulo" bson:"titulo"`
	Mensaje       string                 `json:"mensaje" bson:"mensaje"`
	Destinatario  string                 `json:"destinatario" bson:"destinatario"`
	FechaCreacion time.Time              `json:"fecha_creacion" bson:"fecha_creacion"`
	FechaEnvio    time.Time              `json:"fecha_envio,omitempty" bson:"fecha_envio,omitempty"`
	Estado        EstadoNotificacion     `json:"estado" bson:"estado"`
	Prioridad     string                 `json:"prioridad" bson:"prioridad"`
	Metadata      map[string]interface{} `json:"metadata,omitempty" bson:"metadata,omitempty"`
	Leida         bool                   `json:"leida" bson:"leida"`
	FechaLectura  time.Time              `json:"fecha_lectura,omitempty" bson:"fecha_lectura,omitempty"`
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
