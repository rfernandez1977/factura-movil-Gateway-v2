package models

import "fmt"

// ValidationFieldError representa un error de validación de campo
type ValidationFieldError struct {
	Field   string      `json:"field"`
	Message string      `json:"message"`
	Code    string      `json:"code,omitempty"`
	Value   interface{} `json:"value,omitempty"`
}

// Error implementa la interfaz error
func (e *ValidationFieldError) Error() string {
	if e.Code != "" {
		return fmt.Sprintf("[%s] %s: %s", e.Code, e.Field, e.Message)
	}
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// NewValidationFieldError crea un nuevo error de validación de campo
func NewValidationFieldError(field, code, message string, value interface{}) error {
	return &ValidationFieldError{
		Field:   field,
		Code:    code,
		Message: message,
		Value:   value,
	}
}
