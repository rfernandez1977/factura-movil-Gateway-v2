package services

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"FMgo/models"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (m *MockRedisClient) Pipeline() redis.Pipeliner {
	args := m.Called()
	return args.Get(0).(redis.Pipeliner)
}

func (m *MockRedisClient) Incr(ctx context.Context, key string) *redis.IntCmd {
	args := m.Called(ctx, key)
	return args.Get(0).(*redis.IntCmd)
}

// MockPipeliner es un mock del pipeline de Redis
type MockPipeliner struct {
	mock.Mock
}

func (m *MockPipeliner) Incr(ctx context.Context, key string) *redis.IntCmd {
	args := m.Called(ctx, key)
	return args.Get(0).(*redis.IntCmd)
}

func (m *MockPipeliner) Expire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd {
	args := m.Called(ctx, key, expiration)
	return args.Get(0).(*redis.BoolCmd)
}

func (m *MockPipeliner) Exec(ctx context.Context) ([]redis.Cmder, error) {
	args := m.Called(ctx)
	return args.Get(0).([]redis.Cmder), args.Error(1)
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

func (m *MockMongoCollection) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(*mongo.Cursor), args.Error(1)
}

func (m *MockMongoCollection) ReplaceOne(ctx context.Context, filter interface{}, replacement interface{}, opts ...*options.ReplaceOptions) (*mongo.UpdateResult, error) {
	args := m.Called(ctx, filter, replacement)
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

func TestRetryService_AgregarReintento(t *testing.T) {
	// Preparar mocks
	db := new(MockMongoDatabase)
	cache := new(MockRedisClient)
	collection := new(MockMongoCollection)

	// Crear servicio
	service := NewRetryService(cache, db)

	// Preparar datos de prueba
	ctx := context.Background()
	reintento := &models.ColaReintentos{
		ID:                 primitive.NewObjectID(),
		FlujoID:            primitive.NewObjectID(),
		PasoID:             primitive.NewObjectID(),
		Estado:             models.EstadoReintentoPendiente,
		FechaCreacion:      time.Now(),
		FechaActualizacion: time.Now(),
	}

	// Configurar mocks
	db.On("Collection", "cola_reintentos").Return(collection)
	collection.On("InsertOne", ctx, mock.Anything).Return(&mongo.InsertOneResult{}, nil)
	statusCmd := redis.NewStatusCmd(ctx)
	cache.On("Set", ctx, mock.Anything, mock.Anything, TTLReintento).Return(statusCmd)

	// Ejecutar prueba
	err := service.AgregarReintento(ctx, reintento)

	// Verificar resultado
	assert.NoError(t, err)
	db.AssertExpectations(t)
	cache.AssertExpectations(t)
}

func TestRetryService_ProcesarReintentos(t *testing.T) {
	// Preparar mocks
	db := new(MockMongoDatabase)
	cache := new(MockRedisClient)
	collection := new(MockMongoCollection)

	// Crear servicio
	service := NewRetryService(cache, db)

	// Preparar datos de prueba
	ctx := context.Background()
	reintentos := []models.ColaReintentos{
		{
			ID:                 primitive.NewObjectID(),
			FlujoID:            primitive.NewObjectID(),
			PasoID:             primitive.NewObjectID(),
			Estado:             models.EstadoReintentoPendiente,
			FechaCreacion:      time.Now(),
			FechaActualizacion: time.Now(),
		},
	}

	// Configurar mock para Find
	cursor := mongo.NewCursorFromDocuments(reintentos, nil, nil)
	db.On("Collection", "cola_reintentos").Return(collection)
	collection.On("Find", ctx, mock.Anything).Return(cursor, nil)

	// Configurar mocks para obtener flujo y paso
	flujo := &models.FlujoIntegracion{
		ID:     reintentos[0].FlujoID,
		Nombre: "Flujo de prueba",
	}
	paso := &models.PasoFlujo{
		ID:            reintentos[0].PasoID,
		MaxReintentos: 3,
	}

	// Configurar mocks para caché de flujo
	flujoData, _ := json.Marshal(flujo)
	stringCmd := redis.NewStringCmd(ctx)
	stringCmd.SetVal(string(flujoData))
	cache.On("Get", ctx, mock.Anything).Return(stringCmd).Once()

	// Configurar mocks para caché de paso
	pasoData, _ := json.Marshal(paso)
	stringCmd2 := redis.NewStringCmd(ctx)
	stringCmd2.SetVal(string(pasoData))
	cache.On("Get", ctx, mock.Anything).Return(stringCmd2).Once()

	// Configurar mock para actualización
	collection.On("ReplaceOne", ctx, mock.Anything, mock.Anything).Return(&mongo.UpdateResult{}, nil)
	statusCmd := redis.NewStatusCmd(ctx)
	cache.On("Set", ctx, mock.Anything, mock.Anything, TTLReintento).Return(statusCmd)

	// Ejecutar prueba
	err := service.ProcesarReintentos(ctx)

	// Verificar resultado
	assert.NoError(t, err)
	db.AssertExpectations(t)
	cache.AssertExpectations(t)
}

func TestRetryService_Retry(t *testing.T) {
	// Preparar mocks
	db := new(MockMongoDatabase)
	cache := new(MockRedisClient)
	pipeliner := new(MockPipeliner)

	// Crear servicio
	service := NewRetryService(cache, db)

	// Preparar datos de prueba
	ctx := context.Background()
	operationID := "test-operation"
	config := DefaultRetryConfig()

	t.Run("Operación exitosa sin reintentos", func(t *testing.T) {
		// Configurar mock para obtener intentos
		stringCmd := redis.NewStringCmd(ctx)
		stringCmd.SetErr(redis.Nil)
		cache.On("Get", ctx, mock.Anything).Return(stringCmd).Once()

		// Configurar mock para limpiar intentos
		intCmd := redis.NewIntCmd(ctx)
		cache.On("Del", ctx, mock.Anything).Return(intCmd).Once()

		// Ejecutar prueba
		operation := func() error { return nil }
		err := service.Retry(ctx, operationID, config, operation)

		// Verificar resultado
		assert.NoError(t, err)
		cache.AssertExpectations(t)
	})

	t.Run("Operación con reintento exitoso", func(t *testing.T) {
		// Configurar mock para obtener intentos
		stringCmd := redis.NewStringCmd(ctx)
		stringCmd.SetVal("1")
		cache.On("Get", ctx, mock.Anything).Return(stringCmd).Once()

		// Configurar mock para pipeline
		cache.On("Pipeline").Return(pipeliner)
		intCmd := redis.NewIntCmd(ctx)
		pipeliner.On("Incr", ctx, mock.Anything).Return(intCmd)
		boolCmd := redis.NewBoolCmd(ctx)
		pipeliner.On("Expire", ctx, mock.Anything, TTLReintento).Return(boolCmd)
		pipeliner.On("Exec", ctx).Return([]redis.Cmder{intCmd, boolCmd}, nil)

		// Configurar mock para limpiar intentos
		intCmd2 := redis.NewIntCmd(ctx)
		cache.On("Del", ctx, mock.Anything).Return(intCmd2).Once()

		// Ejecutar prueba
		attempts := 0
		operation := func() error {
			if attempts == 0 {
				attempts++
				return fmt.Errorf("timeout")
			}
			return nil
		}
		err := service.Retry(ctx, operationID, config, operation)

		// Verificar resultado
		assert.NoError(t, err)
		cache.AssertExpectations(t)
		pipeliner.AssertExpectations(t)
	})
}

func TestRetryService_LimpiarCache(t *testing.T) {
	// Preparar mocks
	db := new(MockMongoDatabase)
	cache := new(MockRedisClient)

	// Crear servicio
	service := NewRetryService(cache, db)

	// Preparar datos de prueba
	ctx := context.Background()
	keys := []string{
		"reintento:123",
		"flujo:456",
		"paso:789",
		"intentos:test",
	}

	// Configurar mocks para Scan
	scanCmd1 := redis.NewScanCmd(ctx)
	scanCmd1.SetVal([]string{keys[0]}, 0)
	cache.On("Scan", ctx, uint64(0), PrefijoReintento+"*", int64(100)).Return(scanCmd1)

	scanCmd2 := redis.NewScanCmd(ctx)
	scanCmd2.SetVal([]string{keys[1]}, 0)
	cache.On("Scan", ctx, uint64(0), PrefijoFlujo+"*", int64(100)).Return(scanCmd2)

	scanCmd3 := redis.NewScanCmd(ctx)
	scanCmd3.SetVal([]string{keys[2]}, 0)
	cache.On("Scan", ctx, uint64(0), PrefijoPaso+"*", int64(100)).Return(scanCmd3)

	scanCmd4 := redis.NewScanCmd(ctx)
	scanCmd4.SetVal([]string{keys[3]}, 0)
	cache.On("Scan", ctx, uint64(0), PrefijoIntentos+"*", int64(100)).Return(scanCmd4)

	// Configurar mock para Del
	intCmd := redis.NewIntCmd(ctx)
	cache.On("Del", ctx, keys).Return(intCmd)

	// Ejecutar prueba
	err := service.LimpiarCache(ctx)

	// Verificar resultado
	assert.NoError(t, err)
	cache.AssertExpectations(t)
}
