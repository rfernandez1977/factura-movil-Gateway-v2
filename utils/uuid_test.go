package utils

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestParseUUID(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:    "UUID válido",
			input:   "81f6f4f3-f202-4916-8e3e-eb9a5eba4a50",
			want:    "81f6f4f3-f202-4916-8e3e-eb9a5eba4a50",
			wantErr: false,
		},
		{
			name:    "UUID con espacios",
			input:   " 81f6f4f3-f202-4916-8e3e-eb9a5eba4a50 ",
			want:    "81f6f4f3-f202-4916-8e3e-eb9a5eba4a50",
			wantErr: false,
		},
		{
			name:    "UUID inválido - formato incorrecto",
			input:   "81f6f4f3-f202-4916-8e3e-eb9a5eba4a5",
			want:    "",
			wantErr: true,
		},
		{
			name:    "UUID inválido - caracteres no hexadecimales",
			input:   "81f6f4f3-f202-4916-8e3e-eb9a5eba4a5g",
			want:    "",
			wantErr: true,
		},
		{
			name:    "UUID inválido - cadena vacía",
			input:   "",
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseUUID(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, uuid.Nil, got)
			} else {
				assert.NoError(t, err)
				assert.NotEqual(t, uuid.Nil, got)
				assert.Equal(t, tt.want, got.String())
			}
		})
	}
}

func TestGenerateUUID(t *testing.T) {
	// Generar varios UUIDs y verificar que sean únicos
	uuids := make(map[string]bool)
	for i := 0; i < 1000; i++ {
		uuid := GenerateUUID()
		assert.False(t, uuids[uuid], "UUID duplicado generado")
		uuids[uuid] = true
		assert.True(t, IsValidUUID(uuid), "UUID generado no es válido")
	}
}

func TestIsValidUUID(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "UUID válido",
			input: "81f6f4f3-f202-4916-8e3e-eb9a5eba4a50",
			want:  true,
		},
		{
			name:  "UUID inválido - formato incorrecto",
			input: "81f6f4f3-f202-4916-8e3e-eb9a5eba4a5",
			want:  false,
		},
		{
			name:  "UUID inválido - caracteres no hexadecimales",
			input: "81f6f4f3-f202-4916-8e3e-eb9a5eba4a5g",
			want:  false,
		},
		{
			name:  "UUID inválido - cadena vacía",
			input: "",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsValidUUID(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFormatUUID(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:    "UUID válido",
			input:   "81f6f4f3-f202-4916-8e3e-eb9a5eba4a50",
			want:    "81f6f4f3-f202-4916-8e3e-eb9a5eba4a50",
			wantErr: false,
		},
		{
			name:    "UUID con espacios",
			input:   " 81f6f4f3-f202-4916-8e3e-eb9a5eba4a50 ",
			want:    "81f6f4f3-f202-4916-8e3e-eb9a5eba4a50",
			wantErr: false,
		},
		{
			name:    "UUID inválido",
			input:   "81f6f4f3-f202-4916-8e3e-eb9a5eba4a5",
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FormatUUID(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestGenerateDocumentUUID(t *testing.T) {
	tests := []struct {
		name      string
		docType   string
		date      time.Time
		docNumber string
	}{
		{
			name:      "Factura normal",
			docType:   "FACTURA",
			date:      time.Date(2024, 4, 11, 0, 0, 0, 0, time.UTC),
			docNumber: "12345",
		},
		{
			name:      "Boleta de venta",
			docType:   "BOLETA",
			date:      time.Date(2024, 4, 11, 0, 0, 0, 0, time.UTC),
			docNumber: "67890",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uuid1 := GenerateDocumentUUID(tt.docType, tt.date, tt.docNumber)
			uuid2 := GenerateDocumentUUID(tt.docType, tt.date, tt.docNumber)

			// Verificar que el UUID es válido
			assert.True(t, IsValidUUID(uuid1))

			// Verificar que el mismo input genera el mismo UUID
			assert.Equal(t, uuid1, uuid2)

			// Verificar que diferentes inputs generan diferentes UUIDs
			uuid3 := GenerateDocumentUUID(tt.docType, tt.date, "diferente")
			assert.NotEqual(t, uuid1, uuid3)
		})
	}
}

func TestGenerateTransactionUUID(t *testing.T) {
	tests := []struct {
		name            string
		transactionType string
	}{
		{
			name:            "Venta",
			transactionType: "VENTA",
		},
		{
			name:            "Compra",
			transactionType: "COMPRA",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uuid1 := GenerateTransactionUUID(tt.transactionType)
			uuid2 := GenerateTransactionUUID(tt.transactionType)

			// Verificar que el UUID es válido
			assert.True(t, IsValidUUID(uuid1))

			// Verificar que diferentes llamadas generan diferentes UUIDs
			assert.NotEqual(t, uuid1, uuid2)
		})
	}
}

func TestGenerateClientUUID(t *testing.T) {
	tests := []struct {
		testName   string
		rut        string
		clientName string
	}{
		{
			testName:   "Cliente 1",
			rut:        "76212889-6",
			clientName: "FACTURA MOVIL SPA",
		},
		{
			testName:   "Cliente 2",
			rut:        "12345678-9",
			clientName: "OTRA EMPRESA SPA",
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			uuid1 := GenerateClientUUID(tt.rut, tt.clientName)
			uuid2 := GenerateClientUUID(tt.rut, tt.clientName)

			// Verificar que el UUID es válido
			assert.True(t, IsValidUUID(uuid1))

			// Verificar que el mismo input genera el mismo UUID
			assert.Equal(t, uuid1, uuid2)

			// Verificar que diferentes inputs generan diferentes UUIDs
			uuid3 := GenerateClientUUID(tt.rut, "nombre diferente")
			assert.NotEqual(t, uuid1, uuid3)
		})
	}
}

func TestGenerateProductUUID(t *testing.T) {
	tests := []struct {
		testName    string
		code        string
		productName string
	}{
		{
			testName:    "Producto 1",
			code:        "01",
			productName: "Servicio Mensual Plan Copihue",
		},
		{
			testName:    "Producto 2",
			code:        "02",
			productName: "Otro Servicio",
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			uuid1 := GenerateProductUUID(tt.code, tt.productName)
			uuid2 := GenerateProductUUID(tt.code, tt.productName)

			// Verificar que el UUID es válido
			assert.True(t, IsValidUUID(uuid1))

			// Verificar que el mismo input genera el mismo UUID
			assert.Equal(t, uuid1, uuid2)

			// Verificar que diferentes inputs generan diferentes UUIDs
			uuid3 := GenerateProductUUID(tt.code, "nombre diferente")
			assert.NotEqual(t, uuid1, uuid3)
		})
	}
}
