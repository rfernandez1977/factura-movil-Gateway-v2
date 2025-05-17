package dte

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	// RUTs de prueba
	RutEmpresa  = "76212889-6"
	RutEnviador = "13195458-1"
	RutReceptor = "60803000-K" // SII
)

func TestDTEValidation(t *testing.T) {
	tests := []struct {
		name    string
		input   map[string]interface{}
		wantErr bool
	}{
		{
			name: "valid_dte",
			input: map[string]interface{}{
				"tipo_dte": "33",
				"emisor": map[string]string{
					"rut":          RutEmpresa,
					"razon_social": "EMPRESA DE PRUEBA SPA",
				},
				"receptor": map[string]string{
					"rut":          RutReceptor,
					"razon_social": "SERVICIO DE IMPUESTOS INTERNOS",
				},
				"detalles": []map[string]interface{}{
					{
						"cantidad":        1,
						"descripcion":     "Servicio de Prueba",
						"precio_unitario": 100000,
						"monto_total":     100000,
					},
				},
				"totales": map[string]int{
					"monto_neto": 100000,
					"tasa_iva":   19,
					"iva":        19000,
					"total":      119000,
				},
			},
			wantErr: false,
		},
		{
			name: "invalid_rut",
			input: map[string]interface{}{
				"tipo_dte": "33",
				"emisor": map[string]string{
					"rut":          "76212889-9", // RUT inválido
					"razon_social": "EMPRESA DE PRUEBA SPA",
				},
				"receptor": map[string]string{
					"rut":          RutReceptor,
					"razon_social": "SERVICIO DE IMPUESTOS INTERNOS",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid_total",
			input: map[string]interface{}{
				"tipo_dte": "33",
				"emisor": map[string]string{
					"rut":          RutEmpresa,
					"razon_social": "EMPRESA DE PRUEBA SPA",
				},
				"receptor": map[string]string{
					"rut":          RutReceptor,
					"razon_social": "SERVICIO DE IMPUESTOS INTERNOS",
				},
				"detalles": []map[string]interface{}{
					{
						"cantidad":        1,
						"descripcion":     "Servicio de Prueba",
						"precio_unitario": 100000,
						"monto_total":     100000,
					},
				},
				"totales": map[string]int{
					"monto_neto": 100000,
					"tasa_iva":   19,
					"iva":        19000,
					"total":      120000, // Total incorrecto
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Convertir input a DTE
			jsonData, err := json.Marshal(tt.input)
			if err != nil {
				t.Fatalf("error marshaling input: %v", err)
			}

			var d DTE
			if err := json.Unmarshal(jsonData, &d); err != nil {
				t.Fatalf("error unmarshaling input: %v", err)
			}

			// Validar DTE
			err = d.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
