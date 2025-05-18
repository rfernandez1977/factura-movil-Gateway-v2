package unit

import (
	"encoding/json"
	"testing"

	"FMgo/pkg/dte"
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
					"rut":          "76123456-7",
					"razon_social": "EMPRESA SPA",
				},
				"receptor": map[string]string{
					"rut":          "77654321-8",
					"razon_social": "CLIENTE LTDA",
				},
				"detalles": []map[string]interface{}{
					{
						"cantidad":        1,
						"descripcion":     "Servicio Profesional",
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
					"rut":          "76123456-8", // RUT inválido
					"razon_social": "EMPRESA SPA",
				},
				"receptor": map[string]string{
					"rut":          "77654321-8",
					"razon_social": "CLIENTE LTDA",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid_total",
			input: map[string]interface{}{
				"tipo_dte": "33",
				"emisor": map[string]string{
					"rut":          "76123456-7",
					"razon_social": "EMPRESA SPA",
				},
				"receptor": map[string]string{
					"rut":          "77654321-8",
					"razon_social": "CLIENTE LTDA",
				},
				"detalles": []map[string]interface{}{
					{
						"cantidad":        1,
						"descripcion":     "Servicio Profesional",
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

			var d dte.DTE
			if err := json.Unmarshal(jsonData, &d); err != nil {
				t.Fatalf("error unmarshaling input: %v", err)
			}

			// Validar DTE
			err = d.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("DTE.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
