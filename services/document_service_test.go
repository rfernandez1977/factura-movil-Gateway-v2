package services

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"FMgo/domain"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockDocumentRepository es un mock del repositorio de documentos
type MockDocumentRepository struct {
	mock.Mock
}

func (m *MockDocumentRepository) SaveDocumentoTributario(doc domain.DocumentoTributario) error {
	args := m.Called(doc)
	return args.Error(0)
}

func (m *MockDocumentRepository) GetDocumentoTributario(tipo string, folio int64) (*domain.DocumentoTributario, error) {
	args := m.Called(tipo, folio)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.DocumentoTributario), args.Error(1)
}

func (m *MockDocumentRepository) GetDocumentoTributarioByID(id primitive.ObjectID) (*domain.DocumentoTributario, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.DocumentoTributario), args.Error(1)
}

func (m *MockDocumentRepository) UpdateDocumentoTributario(doc domain.DocumentoTributario) error {
	args := m.Called(doc)
	return args.Error(0)
}

func (m *MockDocumentRepository) SaveEstadoDocumento(estado domain.EstadoDocumento) error {
	args := m.Called(estado)
	return args.Error(0)
}

func (m *MockDocumentRepository) GetEstadoDocumento(docID primitive.ObjectID) (*domain.EstadoDocumento, error) {
	args := m.Called(docID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.EstadoDocumento), args.Error(1)
}

func (m *MockDocumentRepository) SaveReferenciaDocumento(ref domain.ReferenciaDocumento) error {
	args := m.Called(ref)
	return args.Error(0)
}

func (m *MockDocumentRepository) GetReferenciasPorDocumento(tipoOrigen string, folioOrigen int64) ([]domain.ReferenciaDocumento, error) {
	args := m.Called(tipoOrigen, folioOrigen)
	return args.Get(0).([]domain.ReferenciaDocumento), args.Error(1)
}

// MockValidationService es un mock del servicio de validación
type MockValidationService struct {
	mock.Mock
}

func (m *MockValidationService) ValidarDocumento(doc *domain.DocumentoTributario) error {
	args := m.Called(doc)
	return args.Error(0)
}

// MockCAFService es un mock del servicio de CAF
type MockCAFService struct {
	mock.Mock
}

func (m *MockCAFService) ObtenerCAF(ctx context.Context, tipo string) (*domain.CAF, error) {
	args := m.Called(ctx, tipo)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.CAF), args.Error(1)
}

func (m *MockCAFService) ValidarCAF(ctx context.Context, caf *domain.CAF) error {
	args := m.Called(ctx, caf)
	return args.Error(0)
}

func (m *MockCAFService) ActualizarFolioActual(ctx context.Context, caf *domain.CAF) error {
	args := m.Called(ctx, caf)
	return args.Error(0)
}

// MockAuditService es un mock del servicio de auditoría
type MockAuditService struct {
	mock.Mock
}

func (m *MockAuditService) RegistrarOperacion(ctx context.Context, operacion string, entidad string, entidadID primitive.ObjectID, usuario string) error {
	args := m.Called(ctx, operacion, entidad, entidadID, usuario)
	return args.Error(0)
}

// MockRedisClient es un mock del cliente Redis
type MockRedisClient struct {
	mock.Mock
}

func (m *MockRedisClient) Get(ctx context.Context, key string) *redis.StringCmd {
	args := m.Called(ctx, key)
	return args.Get(0).(*redis.StringCmd)
}

func (m *MockRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	args := m.Called(ctx, key, value, expiration)
	return args.Get(0).(*redis.StatusCmd)
}

func (m *MockRedisClient) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	args := m.Called(ctx, keys)
	return args.Get(0).(*redis.IntCmd)
}

func (m *MockRedisClient) Scan(ctx context.Context, cursor uint64, match string, count int64) *redis.ScanCmd {
	args := m.Called(ctx, cursor, match, count)
	return args.Get(0).(*redis.ScanCmd)
}

func TestDocumentService_CrearDocumento(t *testing.T) {
	// Preparar mocks
	repo := new(MockDocumentRepository)
	validationSvc := new(MockValidationService)
	cafSvc := new(MockCAFService)
	auditSvc := new(MockAuditService)
	cache := new(MockRedisClient)

	// Crear servicio
	service := NewDocumentService(repo, validationSvc, cafSvc, auditSvc, cache)

	// Preparar datos de prueba
	ctx := context.Background()
	doc := &domain.DocumentoTributario{
		ID:            primitive.NewObjectID(),
		TipoDocumento: "33",
		Folio:         1,
		FechaEmision:  time.Now(),
		Estado:        "BORRADOR",
	}

	caf := &domain.CAF{
		ID:           "123",
		FolioActual:  1,
		FolioInicial: 1,
		FolioFinal:   100,
	}

	// Configurar expectativas
	validationSvc.On("ValidarDocumento", doc).Return(nil)
	cafSvc.On("ObtenerCAF", ctx, doc.TipoDocumento).Return(caf, nil)
	cafSvc.On("ValidarCAF", ctx, caf).Return(nil)
	repo.On("SaveDocumentoTributario", *doc).Return(nil)
	cafSvc.On("ActualizarFolioActual", ctx, caf).Return(nil)
	auditSvc.On("RegistrarOperacion", ctx, "CREAR_DOCUMENTO", "DocumentoTributario", doc.ID, "sistema").Return(nil)

	// Configurar mock de Redis
	data, _ := json.Marshal(doc)
	statusCmd := redis.NewStatusCmd(ctx)
	cache.On("Set", ctx, mock.Anything, data, TTLDocumento).Return(statusCmd)

	// Ejecutar prueba
	err := service.CrearDocumento(ctx, doc)

	// Verificar resultado
	assert.NoError(t, err)
	repo.AssertExpectations(t)
	validationSvc.AssertExpectations(t)
	cafSvc.AssertExpectations(t)
	auditSvc.AssertExpectations(t)
	cache.AssertExpectations(t)
}

func TestDocumentService_ObtenerDocumento(t *testing.T) {
	// Preparar mocks
	repo := new(MockDocumentRepository)
	validationSvc := new(MockValidationService)
	cafSvc := new(MockCAFService)
	auditSvc := new(MockAuditService)
	cache := new(MockRedisClient)

	// Crear servicio
	service := NewDocumentService(repo, validationSvc, cafSvc, auditSvc, cache)

	// Preparar datos de prueba
	ctx := context.Background()
	doc := &domain.DocumentoTributario{
		ID:            primitive.NewObjectID(),
		TipoDocumento: "33",
		Folio:         1,
		FechaEmision:  time.Now(),
		Estado:        "BORRADOR",
	}

	t.Run("Obtener desde caché", func(t *testing.T) {
		// Configurar mock de Redis
		data, _ := json.Marshal(doc)
		stringCmd := redis.NewStringCmd(ctx)
		stringCmd.SetVal(string(data))
		cache.On("Get", ctx, mock.Anything).Return(stringCmd).Once()

		// Ejecutar prueba
		result, err := service.ObtenerDocumento(ctx, doc.TipoDocumento, doc.Folio)

		// Verificar resultado
		assert.NoError(t, err)
		assert.Equal(t, doc.ID, result.ID)
		cache.AssertExpectations(t)
	})

	t.Run("Obtener desde base de datos", func(t *testing.T) {
		// Configurar mocks
		stringCmd := redis.NewStringCmd(ctx)
		stringCmd.SetErr(redis.Nil)
		cache.On("Get", ctx, mock.Anything).Return(stringCmd).Once()

		repo.On("GetDocumentoTributario", doc.TipoDocumento, doc.Folio).Return(doc, nil).Once()

		data, _ := json.Marshal(doc)
		statusCmd := redis.NewStatusCmd(ctx)
		cache.On("Set", ctx, mock.Anything, data, TTLDocumento).Return(statusCmd).Once()

		// Ejecutar prueba
		result, err := service.ObtenerDocumento(ctx, doc.TipoDocumento, doc.Folio)

		// Verificar resultado
		assert.NoError(t, err)
		assert.Equal(t, doc.ID, result.ID)
		repo.AssertExpectations(t)
		cache.AssertExpectations(t)
	})
}

func TestDocumentService_LimpiarCache(t *testing.T) {
	// Preparar mocks
	repo := new(MockDocumentRepository)
	validationSvc := new(MockValidationService)
	cafSvc := new(MockCAFService)
	auditSvc := new(MockAuditService)
	cache := new(MockRedisClient)

	// Crear servicio
	service := NewDocumentService(repo, validationSvc, cafSvc, auditSvc, cache)

	// Preparar datos de prueba
	ctx := context.Background()
	keys := []string{"doc:33:1", "doc:33:2", "ref:33:1"}

	// Configurar mocks
	scanCmd := redis.NewScanCmd(ctx)
	scanCmd.SetVal([]string{keys[0], keys[1]}, 0)
	cache.On("Scan", ctx, uint64(0), PrefijoDocumento+"*", int64(100)).Return(scanCmd)

	scanCmd2 := redis.NewScanCmd(ctx)
	scanCmd2.SetVal([]string{keys[2]}, 0)
	cache.On("Scan", ctx, uint64(0), PrefijoReferencias+"*", int64(100)).Return(scanCmd2)

	intCmd := redis.NewIntCmd(ctx)
	cache.On("Del", ctx, keys).Return(intCmd)

	// Ejecutar prueba
	err := service.LimpiarCache(ctx)

	// Verificar resultado
	assert.NoError(t, err)
	cache.AssertExpectations(t)
}
