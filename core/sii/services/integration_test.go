package services

import (
	"context"
	"testing"

	"FMgo/core/sii/logger"
	"FMgo/core/sii/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock del cliente SII
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

func (m *MockSIIClient) ConsultarEstado(ctx context.Context, trackID string) (*models.EstadoSII, error) {
	args := m.Called(ctx, trackID)
	return args.Get(0).(*models.EstadoSII), args.Error(1)
}

func (m *MockSIIClient) ConsultarDTE(ctx context.Context, tipoDTE string, folio int64, rutEmisor string) (*models.EstadoSII, error) {
	args := m.Called(ctx, tipoDTE, folio, rutEmisor)
	return args.Get(0).(*models.EstadoSII), args.Error(1)
}

// Mock del servicio de firma
type MockFirmaService struct {
	mock.Mock
}

func (m *MockFirmaService) FirmarDocumento(ctx context.Context, doc *models.Documento) (*models.Documento, error) {
	args := m.Called(ctx, doc)
	return args.Get(0).(*models.Documento), args.Error(1)
}

func TestEnviarDocumento(t *testing.T) {
	// Configurar mocks
	mockSII := new(MockSIIClient)
	mockFirma := new(MockFirmaService)
	logger := logger.NewLogger()

	service := NewIntegrationService(mockSII, mockFirma, logger)

	// Preparar documento de prueba
	doc := &models.Documento{
		ID:          "TEST001",
		TipoDTE:     "33",
		Folio:       1,
		RutEmisor:   "11111111-1",
		RutReceptor: "22222222-2",
		MontoTotal:  1000,
	}

	// Configurar comportamiento esperado
	docFirmado := &models.Documento{
		ID:  doc.ID,
		XML: []byte("<DTE>...</DTE>"),
	}
	mockFirma.On("FirmarDocumento", mock.Anything, doc).Return(docFirmado, nil)
	mockSII.On("ObtenerSemilla", mock.Anything).Return("SEMILLA123", nil)
	mockSII.On("ObtenerToken", mock.Anything, "SEMILLA123").Return("TOKEN123", nil)
	mockSII.On("EnviarDTE", mock.Anything, docFirmado.XML, "TOKEN123").Return(&models.RespuestaSII{
		TrackID: "TRACK123",
		Estado:  "OK",
		Glosa:   "Envío exitoso",
	}, nil)

	// Ejecutar prueba
	ctx := context.Background()
	resultado, err := service.EnviarDocumento(ctx, doc)

	// Verificar resultados
	assert.NoError(t, err)
	assert.NotNil(t, resultado)
	assert.Equal(t, "TRACK123", resultado.TrackID)
	assert.Equal(t, "OK", resultado.Estado)

	// Verificar que se llamaron todos los mocks como se esperaba
	mockFirma.AssertExpectations(t)
	mockSII.AssertExpectations(t)
}

func TestConsultarEstadoEnvio(t *testing.T) {
	// Configurar mocks
	mockSII := new(MockSIIClient)
	mockFirma := new(MockFirmaService)
	logger := logger.NewLogger()

	service := NewIntegrationService(mockSII, mockFirma, logger)

	// Configurar comportamiento esperado
	mockSII.On("ConsultarEstado", mock.Anything, "TRACK123").Return(&models.EstadoSII{
		Estado:  "EPR",
		Glosa:   "Envío Procesado",
		TrackID: "TRACK123",
	}, nil)

	// Ejecutar prueba
	ctx := context.Background()
	estado, err := service.ConsultarEstadoEnvio(ctx, "TRACK123")

	// Verificar resultados
	assert.NoError(t, err)
	assert.NotNil(t, estado)
	assert.Equal(t, "EPR", estado.Estado)
	assert.Equal(t, "Envío Procesado", estado.Glosa)

	// Verificar que se llamaron todos los mocks como se esperaba
	mockSII.AssertExpectations(t)
}

func TestValidarDocumento(t *testing.T) {
	// Configurar mocks
	mockSII := new(MockSIIClient)
	mockFirma := new(MockFirmaService)
	logger := logger.NewLogger()

	service := NewIntegrationService(mockSII, mockFirma, logger)

	// Preparar documento de prueba
	doc := &models.Documento{
		ID:          "TEST001",
		TipoDTE:     "33",
		Folio:       1,
		RutEmisor:   "11111111-1",
		RutReceptor: "22222222-2",
		MontoTotal:  1000,
	}

	// Configurar comportamiento esperado
	mockSII.On("ConsultarDTE", mock.Anything, doc.TipoDTE, doc.Folio, doc.RutEmisor).Return(&models.EstadoSII{
		Estado:  "DTE_RECIBIDO",
		Glosa:   "DTE Recibido",
		TrackID: "TRACK123",
	}, nil)

	// Ejecutar prueba
	ctx := context.Background()
	resultado, err := service.ValidarDocumento(ctx, doc)

	// Verificar resultados
	assert.NoError(t, err)
	assert.NotNil(t, resultado)
	assert.Equal(t, doc.Folio, resultado.Folio)
	assert.Equal(t, doc.TipoDTE, resultado.TipoDTE)
	assert.Equal(t, "DTE_RECIBIDO", resultado.Estado)

	// Verificar que se llamaron todos los mocks como se esperaba
	mockSII.AssertExpectations(t)
}
