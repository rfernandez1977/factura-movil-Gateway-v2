package models

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSIIError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *SIIError
		expected string
	}{
		{
			name: "Error sin causa",
			err: &SIIError{
				Code:    ErrAuthInvalid,
				Message: "Token inválido",
			},
			expected: "[AUTH_001] Token inválido",
		},
		{
			name: "Error con causa",
			err: &SIIError{
				Code:    ErrConexion,
				Message: "Error de conexión",
				Cause:   errors.New("timeout"),
			},
			expected: "[COM_002] Error de conexión: timeout",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.err.Error())
		})
	}
}

func TestSIIError_IsRetryable(t *testing.T) {
	tests := []struct {
		name     string
		err      *SIIError
		expected bool
	}{
		{
			name: "Error de timeout",
			err: &SIIError{
				Code: ErrTimeout,
			},
			expected: true,
		},
		{
			name: "Error de conexión",
			err: &SIIError{
				Code: ErrConexion,
			},
			expected: true,
		},
		{
			name: "Error de servidor",
			err: &SIIError{
				Code: ErrServidor,
			},
			expected: true,
		},
		{
			name: "Error de autenticación",
			err: &SIIError{
				Code: ErrAuthInvalid,
			},
			expected: false,
		},
		{
			name: "Error de DTE",
			err: &SIIError{
				Code: ErrDTEInvalido,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.err.IsRetryable())
		})
	}
}

func TestErrorFromEstado(t *testing.T) {
	tests := []struct {
		name      string
		estado    string
		glosa     string
		expectErr bool
		errCode   ErrorCode
	}{
		{
			name:      "Estado vacío",
			estado:    "",
			glosa:     "",
			expectErr: true,
			errCode:   ErrProcesamiento,
		},
		{
			name:      "Estado rechazado",
			estado:    "RCH",
			glosa:     "DTE Rechazado",
			expectErr: true,
			errCode:   ErrDTEInvalido,
		},
		{
			name:      "Estado OK",
			estado:    "SOK",
			glosa:     "DTE Aceptado",
			expectErr: false,
		},
		{
			name:      "Estado procesado",
			estado:    "EPR",
			glosa:     "DTE En Proceso",
			expectErr: false,
		},
		{
			name:      "Estado desconocido",
			estado:    "XXX",
			glosa:     "Estado Desconocido",
			expectErr: true,
			errCode:   ErrProcesamiento,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ErrorFromEstado(tt.estado, tt.glosa)
			if tt.expectErr {
				assert.NotNil(t, err)
				assert.Equal(t, tt.errCode, err.Code)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
