package models

import (
	"fmt"
)

// ErrorCode representa un código de error específico del SII
type ErrorCode string

const (
	// Errores de autenticación
	ErrAuthInvalid  ErrorCode = "AUTH_001"
	ErrTokenExpired ErrorCode = "AUTH_002"
	ErrCertInvalid  ErrorCode = "AUTH_003"

	// Errores de validación de documentos
	ErrDTEInvalid      ErrorCode = "DTE_001"
	ErrDTEDuplicado    ErrorCode = "DTE_002"
	ErrDTENoEncontrado ErrorCode = "DTE_003"
	ErrFolioInvalido   ErrorCode = "DTE_004"
	ErrRUTInvalido     ErrorCode = "DTE_005"

	// Errores de comunicación
	ErrTimeout  ErrorCode = "COM_001"
	ErrConexion ErrorCode = "COM_002"
	ErrServidor ErrorCode = "COM_003"

	// Errores de procesamiento
	ErrProcesamiento ErrorCode = "PROC_001"
	ErrSchema        ErrorCode = "PROC_002"
	ErrFirma         ErrorCode = "PROC_003"
)

// SIIError representa un error específico del SII
type SIIError struct {
	Code    ErrorCode
	Message string
	Cause   error
}

func (e *SIIError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// IsRetryable determina si un error puede ser reintentado
func (e *SIIError) IsRetryable() bool {
	switch e.Code {
	case ErrTimeout, ErrConexion, ErrServidor:
		return true
	default:
		return false
	}
}

// NewSIIError crea un nuevo error del SII
func NewSIIError(code ErrorCode, message string, cause error) *SIIError {
	return &SIIError{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}

// ErrorFromEstado crea un error basado en el estado del SII
func ErrorFromEstado(estado string, glosa string) *SIIError {
	switch estado {
	case "":
		return NewSIIError(ErrProcesamiento, "Estado no disponible", nil)
	case "RCH":
		return NewSIIError(ErrDTEInvalido, glosa, nil)
	case "SOK":
		return nil // No es un error
	case "EPR":
		return nil // No es un error
	case "RPR":
		return nil // No es un error
	default:
		if glosa != "" {
			return NewSIIError(ErrProcesamiento, glosa, nil)
		}
		return NewSIIError(ErrProcesamiento, "Error desconocido", nil)
	}
}
