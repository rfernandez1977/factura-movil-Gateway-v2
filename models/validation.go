package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// ValidationError representa un error de validación
type ValidationError struct {
	ID        string      `json:"id" bson:"id"`
	Field     string      `json:"field" bson:"field"`
	Code      string      `json:"code" bson:"code"`
	Message   string      `json:"message" bson:"message"`
	Value     interface{} `json:"value" bson:"value"`
	Timestamp time.Time   `json:"timestamp" bson:"timestamp"`
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// Validator interface para implementar validaciones
type Validator interface {
	Validate() error
}

// BaseValidator estructura base para validaciones
type BaseValidator struct {
	errors []*ValidationError
}

func (v *BaseValidator) AddError(field, message, code string, value ...string) {
	err := &ValidationError{
		Field:   field,
		Message: message,
		Code:    code,
	}
	if len(value) > 0 {
		err.Value = value[0]
	}
	v.errors = append(v.errors, err)
}

func (v *BaseValidator) HasErrors() bool {
	return len(v.errors) > 0
}

func (v *BaseValidator) GetErrors() []*ValidationError {
	return v.errors
}

// ValidationRule representa una regla de validación
type ValidationRule struct {
	ID          string    `json:"id" bson:"_id,omitempty"`
	Nombre      string    `json:"nombre" bson:"nombre"`
	Descripcion string    `json:"descripcion" bson:"descripcion"`
	Tipo        string    `json:"tipo" bson:"tipo"`
	Expresion   string    `json:"expresion" bson:"expresion"`
	Mensaje     string    `json:"mensaje" bson:"mensaje"`
	Activo      bool      `json:"activo" bson:"activo"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" bson:"updated_at"`
}

// ValidationResult representa el resultado de una validación
type ValidationResult struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	ReglaID   string    `json:"regla_id" bson:"regla_id"`
	Exitoso   bool      `json:"exitoso" bson:"exitoso"`
	Detalles  string    `json:"detalles" bson:"detalles"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}

// ValidationStatus representa el estado de una validación
type ValidationStatus struct {
	Estado  string    `json:"estado" bson:"estado"`
	Mensaje string    `json:"mensaje" bson:"mensaje"`
	Fecha   time.Time `json:"fecha" bson:"fecha"`
}

// ValidationType representa un tipo de validación
type ValidationType struct {
	ID          string    `json:"id" bson:"_id,omitempty"`
	Nombre      string    `json:"nombre" bson:"nombre"`
	Descripcion string    `json:"descripcion" bson:"descripcion"`
	Reglas      []string  `json:"reglas" bson:"reglas"`
	Activo      bool      `json:"activo" bson:"activo"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" bson:"updated_at"`
}

// ValidationConfig representa la configuración de validación
type ValidationConfig struct {
	ID          string           `json:"id" bson:"_id,omitempty"`
	Tipo        string           `json:"tipo" bson:"tipo"`
	Reglas      []ValidationRule `json:"reglas" bson:"reglas"`
	MaxErrores  int              `json:"max_errores" bson:"max_errores"`
	StopOnError bool             `json:"stop_on_error" bson:"stop_on_error"`
	CreatedAt   time.Time        `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at" bson:"updated_at"`
}

// ValidationRequest representa una solicitud de validación
type ValidationRequest struct {
	DocumentoID string                 `json:"documento_id" bson:"documento_id"`
	Tipo        string                 `json:"tipo" bson:"tipo"`
	Config      ValidationConfig       `json:"config" bson:"config"`
	Metadata    map[string]interface{} `json:"metadata" bson:"metadata"`
}

// ValidationResponse representa la respuesta de una validación
type ValidationResponse struct {
	ID          string             `json:"id" bson:"_id,omitempty"`
	DocumentoID string             `json:"documento_id" bson:"documento_id"`
	Exitoso     bool               `json:"exitoso" bson:"exitoso"`
	Resultados  []ValidationResult `json:"resultados" bson:"resultados"`
	Errores     []ValidationError  `json:"errores" bson:"errores"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
}

// ValidationMetadata representa metadatos de validación
type ValidationMetadata struct {
	ID          string                 `json:"id" bson:"_id,omitempty"`
	DocumentoID string                 `json:"documento_id" bson:"documento_id"`
	Version     string                 `json:"version" bson:"version"`
	Atributos   map[string]interface{} `json:"atributos" bson:"atributos"`
	CreatedAt   time.Time              `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at" bson:"updated_at"`
}

// Suggestion representa una sugerencia de corrección para un error de validación
type Suggestion struct {
	DocumentoID string      `json:"documento_id" bson:"documento_id"`
	ErrorID     string      `json:"error_id" bson:"error_id"`
	Campo       string      `json:"campo" bson:"campo"`
	Tipo        string      `json:"tipo" bson:"tipo"`
	Mensaje     string      `json:"mensaje" bson:"mensaje"`
	Valor       interface{} `json:"valor" bson:"valor"`
	Timestamp   time.Time   `json:"timestamp" bson:"timestamp"`
}

// NewValidationRule crea una nueva regla de validación
func NewValidationRule(nombre, descripcion, tipo, expresion, mensaje string) *ValidationRule {
	return &ValidationRule{
		ID:          GenerateID(),
		Nombre:      nombre,
		Descripcion: descripcion,
		Tipo:        tipo,
		Expresion:   expresion,
		Mensaje:     mensaje,
		Activo:      true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// NewValidationResult crea un nuevo resultado de validación
func NewValidationResult(reglaID string, exitoso bool, detalles string) *ValidationResult {
	return &ValidationResult{
		ID:        GenerateID(),
		ReglaID:   reglaID,
		Exitoso:   exitoso,
		Detalles:  detalles,
		CreatedAt: time.Now(),
	}
}

// NewValidationConfig crea una nueva configuración de validación
func NewValidationConfig(tipo string, reglas []ValidationRule, maxErrores int, stopOnError bool) *ValidationConfig {
	return &ValidationConfig{
		ID:          GenerateID(),
		Tipo:        tipo,
		Reglas:      reglas,
		MaxErrores:  maxErrores,
		StopOnError: stopOnError,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// NewValidationError crea un nuevo error de validación
func NewValidationError(field, code, message string, value interface{}) *ValidationError {
	return &ValidationError{
		ID:        uuid.New().String(),
		Field:     field,
		Code:      code,
		Message:   message,
		Value:     value,
		Timestamp: time.Now(),
	}
}

// Validate implementa la interfaz Validator para ValidationRule
func (r *ValidationRule) Validate() error {
	if r.Nombre == "" {
		return NewValidationError("nombre", "no puede estar vacío", "REQUIRED_FIELD", nil)
	}
	if r.Tipo == "" {
		return NewValidationError("tipo", "no puede estar vacío", "REQUIRED_FIELD", nil)
	}
	if r.Expresion == "" {
		return NewValidationError("expresion", "no puede estar vacía", "REQUIRED_FIELD", nil)
	}
	return nil
}

// Validate implementa la interfaz Validator para ValidationConfig
func (c *ValidationConfig) Validate() error {
	if c.Tipo == "" {
		return NewValidationError("tipo", "no puede estar vacío", "REQUIRED_FIELD", nil)
	}
	if len(c.Reglas) == 0 {
		return NewValidationError("reglas", "debe contener al menos una regla", "REQUIRED_FIELD", nil)
	}
	if c.MaxErrores <= 0 {
		return NewValidationError("max_errores", "debe ser mayor que 0", "INVALID_VALUE", nil)
	}
	return nil
}

// NewSuggestion crea una nueva sugerencia
func NewSuggestion(documentoID, errorID, campo, tipo, mensaje string, valor interface{}) *Suggestion {
	return &Suggestion{
		DocumentoID: documentoID,
		ErrorID:     errorID,
		Campo:       campo,
		Tipo:        tipo,
		Mensaje:     mensaje,
		Valor:       valor,
		Timestamp:   time.Now(),
	}
}
