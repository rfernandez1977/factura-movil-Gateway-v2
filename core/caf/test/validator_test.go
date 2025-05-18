package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"FMgo/core/caf/models"
	"FMgo/core/caf/services"
)

type mockCacheService struct {
	cafs map[string]*models.CAF
}

func (m *mockCacheService) Get(ctx context.Context, key string) (interface{}, error) {
	if caf, ok := m.cafs[key]; ok {
		return caf, nil
	}
	return nil, fmt.Errorf("CAF no encontrado")
}

func (m *mockCacheService) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	if caf, ok := value.(*models.CAF); ok {
		m.cafs[key] = caf
		return nil
	}
	return fmt.Errorf("valor no es un CAF")
}

type mockLogger struct{}

func (m *mockLogger) Info(msg string, args ...interface{})  {}
func (m *mockLogger) Error(msg string, args ...interface{}) {}
func (m *mockLogger) Warn(msg string, args ...interface{})  {}
func (m *mockLogger) Debug(msg string, args ...interface{}) {}

func TestValidarCAF(t *testing.T) {
	cache := &mockCacheService{
		cafs: make(map[string]*models.CAF),
	}
	logger := &mockLogger{}

	validador := services.NewValidadorCAF(cache, logger)

	tests := []struct {
		nombre    string
		xmlCAF    []byte
		quiereErr bool
	}{
		{
			nombre: "CAF válido",
			xmlCAF: []byte(`<?xml version="1.0"?>
<AUTORIZACION version="1.0">
  <CAF>
    <DA>
      <RE>76329692-K</RE>
      <RS>EMPRESA DE PRUEBA</RS>
      <TD>33</TD>
      <RNG><D>1</D><H>100</H></RNG>
      <FA>2024-03-22</FA>
      <RSAPK><M>123</M><E>456</E></RSAPK>
      <IDK>1</IDK>
    </DA>
    <FRMA algoritmo="SHA1withRSA">ABC123</FRMA>
  </CAF>
</AUTORIZACION>`),
			quiereErr: false,
		},
		{
			nombre:    "XML inválido",
			xmlCAF:    []byte(`<xml malo>`),
			quiereErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			resultado, err := validador.ValidarCAF(context.Background(), tt.xmlCAF)

			if tt.quiereErr && err == nil {
				t.Error("se esperaba un error pero no se obtuvo ninguno")
			}

			if !tt.quiereErr && err != nil {
				t.Errorf("no se esperaba error pero se obtuvo: %v", err)
			}

			if err == nil && !resultado.Valido && !tt.quiereErr {
				t.Error("el CAF debería ser válido")
			}
		})
	}
}

func TestValidarFolio(t *testing.T) {
	cache := &mockCacheService{
		cafs: make(map[string]*models.CAF),
	}
	logger := &mockLogger{}

	validador := services.NewValidadorCAF(cache, logger)

	// Agregar un CAF de prueba al caché
	cafPrueba := &models.CAF{
		TipoDTE:    33,
		FolioDesde: 1,
		FolioHasta: 100,
	}
	cache.Set(context.Background(), "CAF_33", cafPrueba, time.Hour)

	tests := []struct {
		nombre    string
		folio     int
		tipoDTE   int
		quiereErr bool
		valido    bool
	}{
		{
			nombre:    "Folio válido",
			folio:     50,
			tipoDTE:   33,
			quiereErr: false,
			valido:    true,
		},
		{
			nombre:    "Folio fuera de rango",
			folio:     101,
			tipoDTE:   33,
			quiereErr: false,
			valido:    false,
		},
		{
			nombre:    "Tipo DTE inválido",
			folio:     1,
			tipoDTE:   34,
			quiereErr: true,
			valido:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			valido, err := validador.ValidarFolio(context.Background(), tt.folio, tt.tipoDTE)

			if tt.quiereErr && err == nil {
				t.Error("se esperaba un error pero no se obtuvo ninguno")
			}

			if !tt.quiereErr && err != nil {
				t.Errorf("no se esperaba error pero se obtuvo: %v", err)
			}

			if err == nil && valido != tt.valido {
				t.Errorf("se esperaba valido=%v pero se obtuvo valido=%v", tt.valido, valido)
			}
		})
	}
}
