package models

import "time"

// ValidationRule representa una regla de validación
type ValidationRule struct {
	ID            string `json:"id" bson:"_id,omitempty"`
	Tipo          string `json:"tipo" bson:"tipo"` // required, format, range, custom
	Expresion     string `json:"expresion" bson:"expresion"`
	Mensaje       string `json:"mensaje" bson:"mensaje"`
	CampoRelativo string `json:"campo_relativo,omitempty" bson:"campo_relativo,omitempty"`
	Valor         string `json:"valor,omitempty" bson:"valor,omitempty"`
	Activo        bool   `json:"activo" bson:"activo"`
}

// ValidationConfig representa una configuración de validación
type ValidationConfig struct {
	ID          string           `json:"id" bson:"_id,omitempty"`
	TipoDTE     string           `json:"tipo_dte" bson:"tipo_dte"`
	Descripcion string           `json:"descripcion" bson:"descripcion"`
	Reglas      []ValidationRule `json:"reglas" bson:"reglas"`
	MaxErrores  int              `json:"max_errores" bson:"max_errores"`
	StopOnError bool             `json:"stop_on_error" bson:"stop_on_error"`
	Activo      bool             `json:"activo" bson:"activo"`
	CreatedAt   time.Time        `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at" bson:"updated_at"`
}

// Validate valida la configuración de validación
func (c *ValidationConfig) Validate() error {
	if c.TipoDTE == "" {
		return NewValidationFieldError("tipo_dte", "REQUIRED", "El tipo de DTE es requerido", nil)
	}
	if len(c.Reglas) == 0 {
		return NewValidationFieldError("reglas", "REQUIRED", "Se requiere al menos una regla de validación", nil)
	}
	return nil
}

// BaseValidator implementa un validador base
type BaseValidator struct {
	errors []*ValidationFieldError
}

// AddError agrega un error de validación
func (v *BaseValidator) AddError(campo, mensaje, codigo string, detalles ...string) {
	v.errors = append(v.errors, &ValidationFieldError{
		Field:   campo,
		Message: mensaje,
		Code:    codigo,
		Value:   detalles,
	})
}

// HasErrors indica si hay errores de validación
func (v *BaseValidator) HasErrors() bool {
	return len(v.errors) > 0
}

// GetErrors obtiene los errores de validación
func (v *BaseValidator) GetErrors() []*ValidationFieldError {
	return v.errors
}

// Suggestion representa una sugerencia de corrección
type Suggestion struct {
	Campo          string      `json:"campo" bson:"campo"`
	ValorActual    interface{} `json:"valor_actual" bson:"valor_actual"`
	Valor          interface{} `json:"valor" bson:"valor"`
	Descripcion    string      `json:"descripcion" bson:"descripcion"`
	Confianza      float64     `json:"confianza" bson:"confianza"`
	FuenteID       string      `json:"fuente_id,omitempty" bson:"fuente_id,omitempty"`
	TipoSugerencia string      `json:"tipo_sugerencia" bson:"tipo_sugerencia"`
}
