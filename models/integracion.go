package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Constantes para estados no definidos previamente
const (
	EstadoFlujoEnProgreso = "EN_PROGRESO" // Este no está en estados.go
)

// Constantes para manejo de errores
const (
	ManejoErrorReintentar = "REINTENTAR"
	ManejoErrorDetener    = "DETENER"
	ManejoErrorIgnorar    = "IGNORAR"
)

// Constantes para tipos de paso
const (
	TipoPasoTransformacion = "TRANSFORMACION"
	TipoPasoValidacion     = "VALIDACION"
	TipoPasoIntegracion    = "INTEGRACION"
	TipoPasoNotificacion   = "NOTIFICACION"
)

// Constantes para estados de reintento
const (
	EstadoReintentoPendiente  = "PENDIENTE"
	EstadoReintentoCompletado = "COMPLETADO"
	EstadoReintentoError      = "ERROR"
)

// FlujoIntegracion representa un flujo de integración
type FlujoIntegracion struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Nombre      string             `json:"nombre" bson:"nombre"`
	Descripcion string             `json:"descripcion" bson:"descripcion"`
	Estado      string             `json:"estado" bson:"estado"`
	PasoActual  int                `json:"paso_actual" bson:"paso_actual"`
	Pasos       []PasoFlujo        `json:"pasos" bson:"pasos"`
	FechaInicio time.Time          `json:"fecha_inicio" bson:"fecha_inicio"`
	FechaFin    time.Time          `json:"fecha_fin,omitempty" bson:"fecha_fin,omitempty"`
	Error       string             `json:"error,omitempty" bson:"error,omitempty"`
	Metadata    map[string]string  `json:"metadata,omitempty" bson:"metadata,omitempty"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

// PasoFlujo representa un paso en un flujo de integración
type PasoFlujo struct {
	ID            primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Nombre        string             `json:"nombre" bson:"nombre"`
	Descripcion   string             `json:"descripcion" bson:"descripcion"`
	Tipo          string             `json:"tipo" bson:"tipo"`
	Estado        string             `json:"estado" bson:"estado"`
	Orden         int                `json:"orden" bson:"orden"`
	Configuracion map[string]string  `json:"configuracion,omitempty" bson:"configuracion,omitempty"`
	ManejoError   string             `json:"manejo_error" bson:"manejo_error"`
	MaxReintentos int                `json:"max_reintentos" bson:"max_reintentos"`
	Intentos      int                `json:"intentos" bson:"intentos"`
	FechaInicio   time.Time          `json:"fecha_inicio,omitempty" bson:"fecha_inicio,omitempty"`
	FechaFin      time.Time          `json:"fecha_fin,omitempty" bson:"fecha_fin,omitempty"`
	Error         string             `json:"error,omitempty" bson:"error,omitempty"`
}

// Validate valida que todos los campos obligatorios estén presentes en un FlujoIntegracion
func (f *FlujoIntegracion) Validate() error {
	if f.Nombre == "" {
		return &ValidationFieldError{Field: "nombre", Message: "El nombre del flujo es obligatorio"}
	}
	if len(f.Pasos) == 0 {
		return &ValidationFieldError{Field: "pasos", Message: "El flujo debe tener al menos un paso"}
	}
	return nil
}

// Validate valida que todos los campos obligatorios estén presentes en un PasoFlujo
func (p *PasoFlujo) Validate() error {
	if p.Nombre == "" {
		return &ValidationFieldError{Field: "nombre", Message: "El nombre del paso es obligatorio"}
	}
	if p.Tipo == "" {
		return &ValidationFieldError{Field: "tipo", Message: "El tipo del paso es obligatorio"}
	}
	if p.Orden < 0 {
		return &ValidationFieldError{Field: "orden", Message: "El orden debe ser mayor o igual a cero"}
	}
	if p.ManejoError == "" {
		return &ValidationFieldError{Field: "manejo_error", Message: "La estrategia de manejo de errores es obligatoria"}
	}
	return nil
}
