package services

import (
	"testing"
	"time"

	"github.com/cursor/FMgo/models"
	"github.com/cursor/FMgo/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository es un mock del repositorio para pruebas
type MockRepository struct {
	mock.Mock
}

// Implementación de los métodos necesarios del repositorio
func (m *MockRepository) GetControlFolio(tipoDocumento string) (*models.ControlFolio, error) {
	args := m.Called(tipoDocumento)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ControlFolio), args.Error(1)
}

func (m *MockRepository) UpdateControlFolio(control models.ControlFolio) error {
	args := m.Called(control)
	return args.Error(0)
}

func (m *MockRepository) CheckFolioUtilizado(tipoDocumento string, folio int) (bool, error) {
	args := m.Called(tipoDocumento, folio)
	return args.Bool(0), args.Error(1)
}

func (m *MockRepository) SaveAsignacionFolio(asignacion models.AsignacionFolio) error {
	args := m.Called(asignacion)
	return args.Error(0)
}

// Implementaciones vacías para satisfacer la interfaz DocumentRepository
func (m *MockRepository) SaveDocumentoTributario(doc models.DocumentoTributario) error {
	return nil
}

func (m *MockRepository) GetDocumentoTributario(tipo string, folio int) (*models.DocumentoTributario, error) {
	return nil, nil
}

func (m *MockRepository) UpdateDocumentoTributario(doc models.DocumentoTributario) error {
	return nil
}

func (m *MockRepository) GetDocumentosPorEstado(estado string) ([]models.DocumentoTributario, error) {
	return nil, nil
}

func (m *MockRepository) SaveEstadoDocumento(estado models.EstadoDocumento) error {
	return nil
}

func (m *MockRepository) GetEstadoDocumento(docID string) (*models.EstadoDocumento, error) {
	return nil, nil
}

func (m *MockRepository) UpdateEstadoDocumento(estado models.EstadoDocumento) error {
	return nil
}

func (m *MockRepository) AddEventoDocumento(docID string, evento models.EventoDocumento) error {
	return nil
}

func (m *MockRepository) SaveReferencia(ref models.ReferenciaDocumento) error {
	return nil
}

func (m *MockRepository) GetReferenciasPorDocumento(tipoOrigen string, folioOrigen int) ([]models.ReferenciaDocumento, error) {
	return nil, nil
}

func (m *MockRepository) SaveOperationLog(log models.OperationLog) error {
	return nil
}

func (m *MockRepository) SaveSecurityLog(log models.RegistroSeguridad) error {
	return nil
}

// TestObtenerSiguienteFolio prueba la obtención del siguiente folio disponible
func TestObtenerSiguienteFolio(t *testing.T) {
	// Crear el mock del repositorio
	mockRepo := new(MockRepository)

	// Crear un servicio de folios con el mock
	servicio := &FolioService{
		Repo: mockRepo,
	}

	// Configurar el comportamiento esperado del mock
	tiempo := time.Now()
	controlFolio := &models.ControlFolio{
		TipoDocumento:     "33",
		RangoInicial:      1,
		RangoFinal:        100,
		FolioActual:       5,
		FoliosDisponibles: 95,
		UltimoUso:         tiempo,
		EstadoCAF:         "ACTIVO",
		AlertaGenerada:    false,
	}

	// Configurar expectativas
	mockRepo.On("GetControlFolio", "33").Return(controlFolio, nil)
	mockRepo.On("CheckFolioUtilizado", "33", 6).Return(false, nil)
	
	// Configurar la expectativa para UpdateControlFolio con cualquier objeto ControlFolio
	mockRepo.On("UpdateControlFolio", mock.AnythingOfType("models.ControlFolio")).Return(nil)
	
	// Configurar la expectativa para SaveAsignacionFolio con cualquier objeto AsignacionFolio
	mockRepo.On("SaveAsignacionFolio", mock.AnythingOfType("models.AsignacionFolio")).Return(nil)

	// Ejecutar la función a probar
	folio, err := servicio.ObtenerSiguienteFolio("33")

	// Verificar resultados
	assert.NoError(t, err)
	assert.Equal(t, 6, folio)
	mockRepo.AssertExpectations(t)
}

// TestVerificarDisponibilidadFolios prueba la verificación de disponibilidad de folios
func TestVerificarDisponibilidadFolios(t *testing.T) {
	// Crear el mock del repositorio
	mockRepo := new(MockRepository)

	// Crear un servicio de folios con el mock
	servicio := &FolioService{
		Repo: mockRepo,
	}

	// Configurar el comportamiento esperado del mock
	tiempo := time.Now()
	controlFolio := &models.ControlFolio{
		TipoDocumento:     "33",
		RangoInicial:      1,
		RangoFinal:        100,
		FolioActual:       95,
		FoliosDisponibles: 5,
		UltimoUso:         tiempo,
		EstadoCAF:         "ACTIVO",
		AlertaGenerada:    false,
	}

	// Configurar expectativas
	mockRepo.On("GetControlFolio", "33").Return(controlFolio, nil)

	// Ejecutar la función a probar
	disponible, foliosRestantes, err := servicio.VerificarDisponibilidadFolios("33", 10)

	// Verificar resultados
	assert.NoError(t, err)
	assert.False(t, disponible)
	assert.Equal(t, 5, foliosRestantes)
	mockRepo.AssertExpectations(t)
}