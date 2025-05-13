package mocks

import (
	"context"
	"time"

	"github.com/cursor/FMgo/sii"
	"github.com/stretchr/testify/mock"
)

// MockSIIService implementa la interfaz SIIService para pruebas
type MockSIIService struct {
	mock.Mock
}

// ConsultarEstado implementa el método de la interfaz SIIService
func (m *MockSIIService) ConsultarEstado(trackID string) (*sii.EstadoSolicitud, error) {
	args := m.Called(trackID)
	return args.Get(0).(*sii.EstadoSolicitud), args.Error(1)
}

// EnviarDocumento implementa el método de la interfaz SIIService
func (m *MockSIIService) EnviarDocumento(ctx context.Context, documento interface{}) (*sii.RespuestaEnvio, error) {
	args := m.Called(ctx, documento)
	return args.Get(0).(*sii.RespuestaEnvio), args.Error(1)
}

// ConsultarDocumento implementa el método de la interfaz SIIService
func (m *MockSIIService) ConsultarDocumento(ctx context.Context, tipo string, folio int64) (*sii.EstadoDocumento, error) {
	args := m.Called(ctx, tipo, folio)
	return args.Get(0).(*sii.EstadoDocumento), args.Error(1)
}

// NewMockSIIService crea una nueva instancia del mock con valores por defecto
func NewMockSIIService() *MockSIIService {
	mock := &MockSIIService{}

	// Configurar respuestas por defecto
	mock.On("ConsultarEstado", Anything).Return(&sii.EstadoSolicitud{
		Estado:         "PENDIENTE",
		Glosa:          "Solicitud en proceso",
		TrackID:        "123456",
		FechaRespuesta: time.Now(),
	}, nil)

	mock.On("EnviarDocumento", Anything, Anything).Return(&sii.RespuestaEnvio{
		Estado:         "RECIBIDO",
		Glosa:          "Documento recibido correctamente",
		TrackID:        "123456",
		FechaRespuesta: time.Now(),
	}, nil)

	mock.On("ConsultarDocumento", Anything, Anything, Anything).Return(&sii.EstadoDocumento{
		Estado:         "ACEPTADO",
		Glosa:          "Documento aceptado",
		FechaRespuesta: time.Now(),
	}, nil)

	return mock
}
