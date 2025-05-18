package mocks

import (
	"FMgo/models"
	"github.com/stretchr/testify/mock"
)

// MockRepository implementa la interfaz DocumentRepository para pruebas
type MockRepository struct {
	mock.Mock
}

// NewMockRepository crea una nueva instancia de MockRepository
func NewMockRepository() *MockRepository {
	return &MockRepository{}
}

func (m *MockRepository) SaveDocumentoTributario(doc models.DocumentoTributario) error {
	args := m.Called(doc)
	return args.Error(0)
}

func (m *MockRepository) GetDocumentoTributario(tipo string, folio int) (*models.DocumentoTributario, error) {
	args := m.Called(tipo, folio)
	return args.Get(0).(*models.DocumentoTributario), args.Error(1)
}

func (m *MockRepository) UpdateDocumentoTributario(doc models.DocumentoTributario) error {
	args := m.Called(doc)
	return args.Error(0)
}

func (m *MockRepository) GetDocumentosPorEstado(estado string) ([]models.DocumentoTributario, error) {
	args := m.Called(estado)
	return args.Get(0).([]models.DocumentoTributario), args.Error(1)
}

func (m *MockRepository) SaveEstadoDocumento(estado models.EstadoDocumento) error {
	args := m.Called(estado)
	return args.Error(0)
}

func (m *MockRepository) GetEstadoDocumento(docID string) (*models.EstadoDocumento, error) {
	args := m.Called(docID)
	return args.Get(0).(*models.EstadoDocumento), args.Error(1)
}

func (m *MockRepository) UpdateEstadoDocumento(estado models.EstadoDocumento) error {
	args := m.Called(estado)
	return args.Error(0)
}

func (m *MockRepository) GetControlFolio(tipoDocumento string) (*models.ControlFolio, error) {
	args := m.Called(tipoDocumento)
	return args.Get(0).(*models.ControlFolio), args.Error(1)
}

func (m *MockRepository) UpdateControlFolio(control models.ControlFolio) error {
	args := m.Called(control)
	return args.Error(0)
}

func (m *MockRepository) SaveAsignacionFolio(asignacion models.AsignacionFolio) error {
	args := m.Called(asignacion)
	return args.Error(0)
}

func (m *MockRepository) CheckFolioUtilizado(tipoDocumento string, folio int) (bool, error) {
	args := m.Called(tipoDocumento, folio)
	return args.Bool(0), args.Error(1)
}

func (m *MockRepository) SaveReferencia(ref models.ReferenciaDocumento) error {
	args := m.Called(ref)
	return args.Error(0)
}

func (m *MockRepository) GetReferenciasPorDocumento(tipoOrigen string, folioOrigen int) ([]models.ReferenciaDocumento, error) {
	args := m.Called(tipoOrigen, folioOrigen)
	return args.Get(0).([]models.ReferenciaDocumento), args.Error(1)
}
