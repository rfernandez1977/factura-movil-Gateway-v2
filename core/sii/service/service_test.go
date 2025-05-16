package service

import (
	"context"
	"testing"
	"time"

	"github.com/fmgo/core/sii/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

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
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.RespuestaSII), args.Error(1)
}

func (m *MockSIIClient) ConsultarEstado(ctx context.Context, trackID string) (*models.EstadoSII, error) {
	args := m.Called(ctx, trackID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.EstadoSII), args.Error(1)
}

func (m *MockSIIClient) ConsultarDTE(ctx context.Context, tipoDTE string, folio int64, rutEmisor string) (*models.EstadoSII, error) {
	args := m.Called(ctx, tipoDTE, folio, rutEmisor)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.EstadoSII), args.Error(1)
}

func (m *MockSIIClient) VerificarComunicacion(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func TestDefaultSIIService_EnviarDTE(t *testing.T) {
	// Crear mock del cliente
	mockClient := new(MockSIIClient)
	mockClient.On("ObtenerSemilla", mock.Anything).Return("SEMILLA-123", nil)
	mockClient.On("ObtenerToken", mock.Anything, "SEMILLA-123").Return("TOKEN-123", nil)
	mockClient.On("EnviarDTE", mock.Anything, []byte("<DTE></DTE>"), "TOKEN-123").Return(&models.RespuestaSII{
		EstadoSII: models.EstadoSII{
			Estado:  "0",
			Glosa:   "DTE Recibido",
			TrackID: "123456",
		},
	}, nil)

	// Crear servicio
	service := NewDefaultSIIService(mockClient)

	// Ejecutar prueba
	resp, err := service.EnviarDTE(context.Background(), []byte("<DTE></DTE>"))
	assert.NoError(t, err)
	assert.Equal(t, "0", resp.Estado)
	assert.Equal(t, "DTE Recibido", resp.Glosa)
	assert.Equal(t, "123456", resp.TrackID)

	// Verificar que se llamaron los métodos esperados
	mockClient.AssertExpectations(t)
}

func TestDefaultSIIService_ConsultarEstado(t *testing.T) {
	// Crear mock del cliente
	mockClient := new(MockSIIClient)
	mockClient.On("ObtenerSemilla", mock.Anything).Return("SEMILLA-123", nil)
	mockClient.On("ObtenerToken", mock.Anything, "SEMILLA-123").Return("TOKEN-123", nil)
	mockClient.On("ConsultarEstado", mock.Anything, "123456").Return(&models.EstadoSII{
		Estado: "EPR",
		Glosa:  "Envío Procesado",
	}, nil)

	// Crear servicio
	service := NewDefaultSIIService(mockClient)

	// Ejecutar prueba
	estado, err := service.ConsultarEstado(context.Background(), "123456")
	assert.NoError(t, err)
	assert.Equal(t, "EPR", estado.Estado)
	assert.Equal(t, "Envío Procesado", estado.Glosa)

	// Verificar que se llamaron los métodos esperados
	mockClient.AssertExpectations(t)
}

func TestDefaultSIIService_ConsultarDTE(t *testing.T) {
	// Crear mock del cliente
	mockClient := new(MockSIIClient)
	mockClient.On("ObtenerSemilla", mock.Anything).Return("SEMILLA-123", nil)
	mockClient.On("ObtenerToken", mock.Anything, "SEMILLA-123").Return("TOKEN-123", nil)
	mockClient.On("ConsultarDTE", mock.Anything, "33", int64(1234), "76212889-6").Return(&models.EstadoSII{
		Estado: "DTE_RECIBIDO",
		Glosa:  "DTE Recibido",
	}, nil)

	// Crear servicio
	service := NewDefaultSIIService(mockClient)

	// Ejecutar prueba
	estado, err := service.ConsultarDTE(context.Background(), "33", 1234, "76212889-6")
	assert.NoError(t, err)
	assert.Equal(t, "DTE_RECIBIDO", estado.Estado)
	assert.Equal(t, "DTE Recibido", estado.Glosa)

	// Verificar que se llamaron los métodos esperados
	mockClient.AssertExpectations(t)
}

func TestDefaultSIIService_VerificarComunicacion(t *testing.T) {
	// Crear mock del cliente
	mockClient := new(MockSIIClient)
	mockClient.On("VerificarComunicacion", mock.Anything).Return(nil)

	// Crear servicio
	service := NewDefaultSIIService(mockClient)

	// Ejecutar prueba
	err := service.VerificarComunicacion(context.Background())
	assert.NoError(t, err)

	// Verificar que se llamaron los métodos esperados
	mockClient.AssertExpectations(t)
}

func TestDefaultSIIService_TokenExpiration(t *testing.T) {
	// Crear mock del cliente
	mockClient := new(MockSIIClient)
	mockClient.On("ObtenerSemilla", mock.Anything).Return("SEMILLA-123", nil).Times(2)
	mockClient.On("ObtenerToken", mock.Anything, "SEMILLA-123").Return("TOKEN-123", nil).Times(2)
	mockClient.On("EnviarDTE", mock.Anything, []byte("<DTE></DTE>"), "TOKEN-123").Return(&models.RespuestaSII{
		EstadoSII: models.EstadoSII{
			Estado:  "0",
			Glosa:   "DTE Recibido",
			TrackID: "123456",
		},
	}, nil).Times(2)

	// Crear servicio
	service := NewDefaultSIIService(mockClient)

	// Primera llamada - obtiene nuevo token
	_, err := service.EnviarDTE(context.Background(), []byte("<DTE></DTE>"))
	assert.NoError(t, err)

	// Simular expiración del token
	service.tokenExp = time.Now().Add(-1 * time.Hour)

	// Segunda llamada - debe obtener nuevo token
	_, err = service.EnviarDTE(context.Background(), []byte("<DTE></DTE>"))
	assert.NoError(t, err)

	// Verificar que se llamaron los métodos esperados
	mockClient.AssertExpectations(t)
}
