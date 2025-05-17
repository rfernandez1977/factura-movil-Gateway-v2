package service

import (
	"context"
	"testing"
	"time"

	"github.com/fmgo/core/sii/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// SIIClient define la interfaz para el cliente del SII
type SIIClient interface {
	ObtenerSemilla(ctx context.Context) (string, error)
	ObtenerToken(ctx context.Context, semilla string) (string, error)
	EnviarDTE(ctx context.Context, sobre []byte, token string) (*models.RespuestaSII, error)
	ConsultarEstado(ctx context.Context, trackID string) (*models.EstadoConsulta, error)
	ConsultarDTE(ctx context.Context, tipoDTE models.TipoDocumentoSII, folio int64, rutEmisor string) (*models.EstadoConsulta, error)
	VerificarComunicacion(ctx context.Context) error
}

// MockSIIClient es un mock del cliente SII para pruebas
type MockSIIClient struct {
	mock.Mock
}

func (m *MockSIIClient) ObtenerSemilla(ctx context.Context) (string, error) {
	args := m.Called(ctx)
	return args.String(0), args.Error(1)
}

func (m *MockSIIClient) ObtenerToken(ctx context.Context, semilla string) (string, error) {
	args := m.Called(ctx, semilla)
	return args.String(0), args.Error(1)
}

func (m *MockSIIClient) EnviarDTE(ctx context.Context, sobre []byte, token string) (*models.RespuestaSII, error) {
	args := m.Called(ctx, sobre, token)
	return args.Get(0).(*models.RespuestaSII), args.Error(1)
}

func (m *MockSIIClient) ConsultarEstado(ctx context.Context, trackID string) (*models.EstadoConsulta, error) {
	args := m.Called(ctx, trackID)
	return args.Get(0).(*models.EstadoConsulta), args.Error(1)
}

func (m *MockSIIClient) ConsultarDTE(ctx context.Context, tipoDTE models.TipoDocumentoSII, folio int64, rutEmisor string) (*models.EstadoConsulta, error) {
	args := m.Called(ctx, tipoDTE, folio, rutEmisor)
	return args.Get(0).(*models.EstadoConsulta), args.Error(1)
}

func (m *MockSIIClient) VerificarComunicacion(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// testContext retorna un contexto de prueba con request_id
func testContext() context.Context {
	return context.WithValue(context.Background(), "request_id", "test-123")
}

func TestEnviarDTE(t *testing.T) {
	mockClient := new(MockSIIClient)
	service := NewSIIService(mockClient)
	ctx := testContext()

	// Preparar datos de prueba
	dte := []byte(`<?xml version="1.0" encoding="UTF-8"?><DTE>...</DTE>`)
	respuestaEsperada := &models.RespuestaSII{
		TrackID:      "123456789",
		Estado:       models.EstadoRecibido,
		Glosa:        "DTE Recibido",
		FechaProceso: time.Now(),
	}

	// Configurar comportamiento esperado del mock
	mockClient.On("ObtenerSemilla", ctx).Return("SEMILLA123", nil)
	mockClient.On("ObtenerToken", ctx, "SEMILLA123").Return("TOKEN123", nil)
	mockClient.On("EnviarDTE", ctx, dte, "TOKEN123").Return(respuestaEsperada, nil)

	// Ejecutar prueba
	respuesta, err := service.EnviarDTE(ctx, dte)

	// Verificar resultados
	assert.NoError(t, err)
	assert.NotNil(t, respuesta)
	assert.Equal(t, respuestaEsperada.TrackID, respuesta.TrackID)
	assert.Equal(t, respuestaEsperada.Estado, respuesta.Estado)
	mockClient.AssertExpectations(t)
}

func TestConsultarEstado(t *testing.T) {
	mockClient := new(MockSIIClient)
	service := NewSIIService(mockClient)
	ctx := testContext()

	// Preparar datos de prueba
	trackID := "123456789"
	estadoEsperado := &models.EstadoConsulta{
		TrackID:        trackID,
		Estado:         models.EstadoAceptado,
		Glosa:          "DTE Aceptado",
		FechaRecepcion: time.Now(),
		FechaProceso:   time.Now(),
	}

	// Configurar comportamiento esperado del mock
	mockClient.On("ConsultarEstado", ctx, trackID).Return(estadoEsperado, nil)

	// Ejecutar prueba
	estado, err := service.ConsultarEstado(ctx, trackID)

	// Verificar resultados
	assert.NoError(t, err)
	assert.NotNil(t, estado)
	assert.Equal(t, estadoEsperado.TrackID, estado.TrackID)
	assert.Equal(t, estadoEsperado.Estado, estado.Estado)
	mockClient.AssertExpectations(t)
}

func TestConsultarDTE(t *testing.T) {
	mockClient := new(MockSIIClient)
	service := NewSIIService(mockClient)
	ctx := testContext()

	// Preparar datos de prueba
	tipoDTE := models.DTEFactura
	folio := int64(1234)
	rutEmisor := "76.123.456-7"
	estadoEsperado := &models.EstadoConsulta{
		TrackID:        "123456789",
		Estado:         models.EstadoAceptado,
		Glosa:          "DTE Aceptado",
		FechaRecepcion: time.Now(),
		FechaProceso:   time.Now(),
	}

	// Configurar comportamiento esperado del mock
	mockClient.On("ConsultarDTE", ctx, tipoDTE, folio, rutEmisor).Return(estadoEsperado, nil)

	// Ejecutar prueba
	estado, err := service.ConsultarDTE(ctx, tipoDTE, folio, rutEmisor)

	// Verificar resultados
	assert.NoError(t, err)
	assert.NotNil(t, estado)
	assert.Equal(t, estadoEsperado.TrackID, estado.TrackID)
	assert.Equal(t, estadoEsperado.Estado, estado.Estado)
	mockClient.AssertExpectations(t)
}

func TestVerificarComunicacion(t *testing.T) {
	mockClient := new(MockSIIClient)
	service := NewSIIService(mockClient)
	ctx := testContext()

	// Configurar comportamiento esperado del mock
	mockClient.On("VerificarComunicacion", ctx).Return(nil)

	// Ejecutar prueba
	err := service.VerificarComunicacion(ctx)

	// Verificar resultados
	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}
