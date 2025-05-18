package dte

import (
	"context"
	"testing"

	"FMgo/core/caf/models"
)

// mockValidadorCAF implementa la interfaz ValidadorCAF para pruebas
type mockValidadorCAF struct {
	foliosValidos map[int]bool
}

func newMockValidadorCAF() *mockValidadorCAF {
	return &mockValidadorCAF{
		foliosValidos: map[int]bool{
			1: true,
			2: true,
			3: true,
			4: true,
			5: true,
		},
	}
}

func (m *mockValidadorCAF) ValidarCAF(ctx context.Context, xmlCAF []byte) (*models.ResultadoValidacion, error) {
	return &models.ResultadoValidacion{
		Valido: true,
	}, nil
}

func (m *mockValidadorCAF) ValidarFolio(ctx context.Context, folio int, tipoDTE int) (bool, error) {
	return m.foliosValidos[folio], nil
}

func (m *mockValidadorCAF) ObtenerRangoFolios(ctx context.Context, tipoDTE int) (int, int, error) {
	return 1, 5, nil
}

func TestValidadorDTE(t *testing.T) {
	validadorCAF := newMockValidadorCAF()
	validador := NewValidadorDTE(validadorCAF)

	tests := []struct {
		nombre    string
		dte       *DTE
		quiereErr bool
	}{
		{
			nombre: "DTE válido",
			dte: &DTE{
				TipoDTE:      "33",
				Folio:        1,
				RUTEmisor:    "76329692-K",
				RUTReceptor:  "76329693-1",
				MontoTotal:   1000,
				FechaEmision: "2024-03-22",
			},
			quiereErr: false,
		},
		{
			nombre: "Folio inválido",
			dte: &DTE{
				TipoDTE:      "33",
				Folio:        10, // Fuera del rango válido
				RUTEmisor:    "76329692-K",
				RUTReceptor:  "76329693-1",
				MontoTotal:   1000,
				FechaEmision: "2024-03-22",
			},
			quiereErr: true,
		},
		{
			nombre: "Tipo DTE no soportado",
			dte: &DTE{
				TipoDTE:      "99",
				Folio:        1,
				RUTEmisor:    "76329692-K",
				RUTReceptor:  "76329693-1",
				MontoTotal:   1000,
				FechaEmision: "2024-03-22",
			},
			quiereErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			// Validar DTE
			err := validador.ValidarDTE(tt.dte)
			if err != nil {
				t.Errorf("Error validando DTE: %v", err)
			}

			// Validar CAF
			err = validador.ValidarCAF(tt.dte)
			if tt.quiereErr && err == nil {
				t.Error("se esperaba un error pero no se obtuvo ninguno")
			}
			if !tt.quiereErr && err != nil {
				t.Errorf("no se esperaba error pero se obtuvo: %v", err)
			}
		})
	}
}
