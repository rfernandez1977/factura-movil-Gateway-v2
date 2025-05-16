package services

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/fmgo/models"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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

// MockMongoCollection es un mock de la colección de MongoDB
type MockMongoCollection struct {
	mock.Mock
}

func (m *MockMongoCollection) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	args := m.Called(ctx, document)
	return args.Get(0).(*mongo.InsertOneResult), args.Error(1)
}

func (m *MockMongoCollection) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	args := m.Called(ctx, filter)
	return args.Get(0).(*mongo.SingleResult)
}

func (m *MockMongoCollection) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	args := m.Called(ctx, filter, update)
	return args.Get(0).(*mongo.UpdateResult), args.Error(1)
}

// MockMongoDatabase es un mock de la base de datos MongoDB
type MockMongoDatabase struct {
	mock.Mock
}

func (m *MockMongoDatabase) Collection(name string) *mongo.Collection {
	args := m.Called(name)
	return args.Get(0).(*mongo.Collection)
}

func TestIntegrationService_IniciarSincronizacion(t *testing.T) {
	// Preparar mocks
	db := new(MockMongoDatabase)
	cache := new(MockRedisClient)
	collection := new(MockMongoCollection)

	// Crear servicio
	service := NewIntegrationService(db, cache, nil)

	// Preparar datos de prueba
	ctx := context.Background()
	datos := map[string]interface{}{
		"campo1": "valor1",
		"campo2": "valor2",
	}

	// Configurar expectativas
	db.On("Collection", "registros_sincronizacion").Return(collection)
	collection.On("InsertOne", ctx, mock.Anything).Return(&mongo.InsertOneResult{}, nil)
	statusCmd := redis.NewStatusCmd(ctx)
	cache.On("Set", ctx, mock.Anything, mock.Anything, TTLRegistro).Return(statusCmd)

	// Ejecutar prueba
	registro, err := service.IniciarSincronizacion(ctx, "erp1", "entidad1", "entrada", datos)

	// Verificar resultado
	assert.NoError(t, err)
	assert.NotNil(t, registro)
	assert.Equal(t, "erp1", registro.ERPID)
	assert.Equal(t, "entidad1", registro.Entidad)
	assert.Equal(t, "entrada", registro.Direccion)
	assert.Equal(t, models.EstadoPendiente, registro.Estado)
	assert.Equal(t, datos, registro.DatosOriginales)
	db.AssertExpectations(t)
	cache.AssertExpectations(t)
}

func TestIntegrationService_ProcesarSincronizacion(t *testing.T) {
	// Preparar mocks
	db := new(MockMongoDatabase)
	cache := new(MockRedisClient)
	collection := new(MockMongoCollection)
	workflowCollection := new(MockMongoCollection)

	// Crear servicio
	service := NewIntegrationService(db, cache, nil)

	// Preparar datos de prueba
	ctx := context.Background()
	registroID := "123"
	registro := &models.RegistroSincronizacion{
		ID:                 registroID,
		ERPID:              "erp1",
		Entidad:            "entidad1",
		Estado:             models.EstadoPendiente,
		FechaCreacion:      time.Now(),
		FechaActualizacion: time.Now(),
	}

	workflow := &models.Workflow{
		ERPID:   "erp1",
		Entidad: "entidad1",
		Pasos:   []models.PasoWorkflow{},
	}

	t.Run("Procesar con registro en caché", func(t *testing.T) {
		// Configurar mocks para registro en caché
		data, _ := json.Marshal(registro)
		stringCmd := redis.NewStringCmd(ctx)
		stringCmd.SetVal(string(data))
		cache.On("Get", ctx, mock.Anything).Return(stringCmd).Once()

		// Configurar mocks para workflow en caché
		workflowData, _ := json.Marshal(workflow)
		stringCmd2 := redis.NewStringCmd(ctx)
		stringCmd2.SetVal(string(workflowData))
		cache.On("Get", ctx, mock.Anything).Return(stringCmd2).Once()

		// Configurar mocks para actualización
		db.On("Collection", "registros_sincronizacion").Return(collection)
		collection.On("UpdateOne", ctx, mock.Anything, mock.Anything).Return(&mongo.UpdateResult{}, nil).Times(2)
		statusCmd := redis.NewStatusCmd(ctx)
		cache.On("Set", ctx, mock.Anything, mock.Anything, TTLRegistro).Return(statusCmd).Times(2)

		// Ejecutar prueba
		err := service.ProcesarSincronizacion(ctx, registroID)

		// Verificar resultado
		assert.NoError(t, err)
		db.AssertExpectations(t)
		cache.AssertExpectations(t)
	})

	t.Run("Procesar sin registro en caché", func(t *testing.T) {
		// Configurar mocks para registro no en caché
		stringCmd := redis.NewStringCmd(ctx)
		stringCmd.SetErr(redis.Nil)
		cache.On("Get", ctx, mock.Anything).Return(stringCmd).Once()

		// Configurar mocks para obtener de base de datos
		db.On("Collection", "registros_sincronizacion").Return(collection)
		singleResult := mongo.NewSingleResultFromDocument(registro, nil, nil)
		collection.On("FindOne", ctx, mock.Anything).Return(singleResult).Once()

		// Configurar mocks para workflow no en caché
		stringCmd2 := redis.NewStringCmd(ctx)
		stringCmd2.SetErr(redis.Nil)
		cache.On("Get", ctx, mock.Anything).Return(stringCmd2).Once()

		// Configurar mocks para obtener workflow de base de datos
		db.On("Collection", "workflows").Return(workflowCollection)
		workflowResult := mongo.NewSingleResultFromDocument(workflow, nil, nil)
		workflowCollection.On("FindOne", ctx, mock.Anything).Return(workflowResult).Once()

		// Configurar mocks para actualizaciones
		collection.On("UpdateOne", ctx, mock.Anything, mock.Anything).Return(&mongo.UpdateResult{}, nil).Times(2)
		statusCmd := redis.NewStatusCmd(ctx)
		cache.On("Set", ctx, mock.Anything, mock.Anything, mock.Anything).Return(statusCmd).Times(3)

		// Ejecutar prueba
		err := service.ProcesarSincronizacion(ctx, registroID)

		// Verificar resultado
		assert.NoError(t, err)
		db.AssertExpectations(t)
		cache.AssertExpectations(t)
	})
}

func TestIntegrationService_RegistrarMetrica(t *testing.T) {
	// Preparar mocks
	db := new(MockMongoDatabase)
	cache := new(MockRedisClient)
	collection := new(MockMongoCollection)

	// Crear servicio
	service := NewIntegrationService(db, cache, nil)

	// Preparar datos de prueba
	ctx := context.Background()
	metrica := &models.MetricaIntegracion{
		ERPID:     "erp1",
		Tipo:      "RENDIMIENTO",
		Valor:     123.45,
		Timestamp: time.Now(),
	}

	// Configurar mocks
	db.On("Collection", "metricas_integracion").Return(collection)
	collection.On("InsertOne", ctx, mock.Anything).Return(&mongo.InsertOneResult{}, nil)

	// Configurar mock de Redis para obtener métricas existentes
	stringCmd := redis.NewStringCmd(ctx)
	stringCmd.SetErr(redis.Nil)
	cache.On("Get", ctx, mock.Anything).Return(stringCmd)

	// Configurar mock de Redis para guardar métricas
	statusCmd := redis.NewStatusCmd(ctx)
	cache.On("Set", ctx, mock.Anything, mock.Anything, TTLMetricas).Return(statusCmd)

	// Ejecutar prueba
	err := service.RegistrarMetrica(ctx, metrica)

	// Verificar resultado
	assert.NoError(t, err)
	db.AssertExpectations(t)
	cache.AssertExpectations(t)
}

func TestIntegrationService_LimpiarCache(t *testing.T) {
	// Preparar mocks
	db := new(MockMongoDatabase)
	cache := new(MockRedisClient)

	// Crear servicio
	service := NewIntegrationService(db, cache, nil)

	// Preparar datos de prueba
	ctx := context.Background()
	keys := []string{
		"workflow:erp1:entidad1",
		"registro:123",
		"metricas:erp1",
	}

	// Configurar mocks para Scan
	scanCmd1 := redis.NewScanCmd(ctx)
	scanCmd1.SetVal([]string{keys[0]}, 0)
	cache.On("Scan", ctx, uint64(0), PrefijoWorkflow+"*", int64(100)).Return(scanCmd1)

	scanCmd2 := redis.NewScanCmd(ctx)
	scanCmd2.SetVal([]string{keys[1]}, 0)
	cache.On("Scan", ctx, uint64(0), PrefijoRegistro+"*", int64(100)).Return(scanCmd2)

	scanCmd3 := redis.NewScanCmd(ctx)
	scanCmd3.SetVal([]string{keys[2]}, 0)
	cache.On("Scan", ctx, uint64(0), PrefijoMetricas+"*", int64(100)).Return(scanCmd3)

	// Configurar mock para Del
	intCmd := redis.NewIntCmd(ctx)
	cache.On("Del", ctx, keys).Return(intCmd)

	// Ejecutar prueba
	err := service.LimpiarCache(ctx)

	// Verificar resultado
	assert.NoError(t, err)
	cache.AssertExpectations(t)
}
