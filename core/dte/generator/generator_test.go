package generator

import (
	"testing"
	"time"

	"FMgo/core/dte/types"
)

func TestGenerateDTE(t *testing.T) {
	generator := NewDefaultGenerator()

	tests := []struct {
		name    string
		input   *GenerateInput
		wantErr bool
	}{
		{
			name: "Generar DTE válido",
			input: &GenerateInput{
				TipoDTE:      "33",
				Folio:        1,
				FechaEmision: time.Now(),
				Emisor: types.Emisor{
					RUT:         "76212889-6",
					RazonSocial: "Empresa Test",
					Giro:        "Servicios",
					Direccion:   "Calle Test 123",
					Comuna:      "Santiago",
					Ciudad:      "Santiago",
				},
				Receptor: types.Receptor{
					RUT:         "13195458-1",
					RazonSocial: "Cliente Test",
					Giro:        "Comercio",
					Direccion:   "Av Test 456",
					Comuna:      "Santiago",
					Ciudad:      "Santiago",
				},
				Detalles: []types.Detalle{
					{
						NumeroLinea: 1,
						Nombre:      "Producto Test",
						Cantidad:    1,
						Precio:      1000,
						MontoItem:   1000,
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dte, err := generator.GenerateDTE(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateDTE() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verificar que el DTE se generó correctamente
				if dte == nil {
					t.Error("GenerateDTE() returned nil DTE")
					return
				}

				// Verificar ID
				if dte.ID == "" {
					t.Error("GenerateDTE() returned DTE with empty ID")
				}

				// Verificar estado inicial
				if dte.Estado != "CREADO" {
					t.Errorf("GenerateDTE() returned DTE with wrong state, got %s, want CREADO", dte.Estado)
				}

				// Verificar que no está firmado
				if dte.Firmado {
					t.Error("GenerateDTE() returned DTE marked as signed")
				}

				// Verificar totales
				totales := dte.Documento.Encabezado.Totales
				expectedIVA := 190.0    // 1000 * 0.19
				expectedTotal := 1190.0 // 1000 + 190

				if totales.IVA != expectedIVA {
					t.Errorf("GenerateDTE() wrong IVA calculation, got %.2f, want %.2f", totales.IVA, expectedIVA)
				}

				if totales.MontoTotal != expectedTotal {
					t.Errorf("GenerateDTE() wrong total calculation, got %.2f, want %.2f", totales.MontoTotal, expectedTotal)
				}
			}
		})
	}
}

func TestCalcularTotales(t *testing.T) {
	tests := []struct {
		name           string
		detalles       []types.Detalle
		wantMontoNeto  float64
		wantIVA        float64
		wantMontoTotal float64
	}{
		{
			name: "Cálculo básico",
			detalles: []types.Detalle{
				{
					MontoItem: 1000,
					Exento:    false,
				},
			},
			wantMontoNeto:  1000,
			wantIVA:        190,
			wantMontoTotal: 1190,
		},
		{
			name: "Múltiples detalles",
			detalles: []types.Detalle{
				{
					MontoItem: 1000,
					Exento:    false,
				},
				{
					MontoItem: 2000,
					Exento:    false,
				},
			},
			wantMontoNeto:  3000,
			wantIVA:        570,
			wantMontoTotal: 3570,
		},
		{
			name: "Con items exentos",
			detalles: []types.Detalle{
				{
					MontoItem: 1000,
					Exento:    false,
				},
				{
					MontoItem: 2000,
					Exento:    true,
				},
			},
			wantMontoNeto:  1000,
			wantIVA:        190,
			wantMontoTotal: 1190,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			montoNeto, iva, montoTotal := calcularTotales(tt.detalles)

			if montoNeto != tt.wantMontoNeto {
				t.Errorf("calcularTotales() montoNeto = %.2f, want %.2f", montoNeto, tt.wantMontoNeto)
			}
			if iva != tt.wantIVA {
				t.Errorf("calcularTotales() iva = %.2f, want %.2f", iva, tt.wantIVA)
			}
			if montoTotal != tt.wantMontoTotal {
				t.Errorf("calcularTotales() montoTotal = %.2f, want %.2f", montoTotal, tt.wantMontoTotal)
			}
		})
	}
}
