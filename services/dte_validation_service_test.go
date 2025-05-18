package services

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"FMgo/models"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

func TestDTEValidationService_RegisterValidation(t *testing.T) {
	// Preparar mock
	cache := new(MockRedisClient)
	service := NewDTEValidationService(cache)

	// Preparar datos de prueba
	ctx := context.Background()
	config := &models.ValidationConfig{
		Reglas: []models.ValidationRule{
			{
				Expresion: "RUT",
				Tipo:      "required",
				Mensaje:   "RUT es requerido",
			},
		},
		StopOnError: true,
		MaxErrores:  1,
	}

	// Configurar mock de Redis
	statusCmd := redis.NewStatusCmd(ctx)
	cache.On("Set", ctx, mock.Anything, mock.Anything, TTLValidationConfig).Return(statusCmd)

	// Ejecutar prueba
	err := service.RegisterValidation(ctx, "33", config)

	// Verificar resultado
	assert.NoError(t, err)
	cache.AssertExpectations(t)
}

func TestDTEValidationService_ValidateDTE(t *testing.T) {
	// Preparar mock
	cache := new(MockRedisClient)
	service := NewDTEValidationService(cache)

	// Preparar datos de prueba
	ctx := context.Background()
	doc := &models.DocumentoTributario{
		TipoDTE: "33",
		Folio:   1,
	}

	config := &models.ValidationConfig{
		Reglas: []models.ValidationRule{
			{
				Expresion: "RUT",
				Tipo:      "required",
				Mensaje:   "RUT es requerido",
			},
		},
		StopOnError: true,
		MaxErrores:  1,
	}

	t.Run("Validar con resultado en caché", func(t *testing.T) {
		// Preparar resultado en caché
		errors := []*models.ValidationFieldError{
			{
				Campo:   "RUT",
				Mensaje: "RUT es requerido",
				Codigo:  "FIELD_REQUIRED",
			},
		}
		data, _ := json.Marshal(errors)

		// Configurar mock de Redis
		stringCmd := redis.NewStringCmd(ctx)
		stringCmd.SetVal(string(data))
		cache.On("Get", ctx, mock.Anything).Return(stringCmd).Once()

		// Ejecutar prueba
		result, err := service.ValidateDTE(ctx, doc)

		// Verificar resultado
		assert.NoError(t, err)
		assert.Equal(t, 1, len(result))
		assert.Equal(t, "RUT", result[0].Campo)
		cache.AssertExpectations(t)
	})

	t.Run("Validar sin resultado en caché", func(t *testing.T) {
		// Configurar mock de Redis para resultado no encontrado
		stringCmd := redis.NewStringCmd(ctx)
		stringCmd.SetErr(redis.Nil)
		cache.On("Get", ctx, mock.Anything).Return(stringCmd).Once()

		// Configurar mock de Redis para obtener configuración
		configData, _ := json.Marshal(config)
		stringCmd2 := redis.NewStringCmd(ctx)
		stringCmd2.SetVal(string(configData))
		cache.On("Get", ctx, mock.Anything).Return(stringCmd2).Once()

		// Configurar mock de Redis para guardar resultado
		statusCmd := redis.NewStatusCmd(ctx)
		cache.On("Set", ctx, mock.Anything, mock.Anything, TTLValidationResult).Return(statusCmd).Once()

		// Ejecutar prueba
		result, err := service.ValidateDTE(ctx, doc)

		// Verificar resultado
		assert.NoError(t, err)
		assert.Equal(t, 1, len(result))
		assert.Equal(t, "RUT", result[0].Campo)
		cache.AssertExpectations(t)
	})
}

func TestDTEValidationService_LimpiarCache(t *testing.T) {
	// Preparar mock
	cache := new(MockRedisClient)
	service := NewDTEValidationService(cache)

	// Preparar datos de prueba
	ctx := context.Background()
	keys := []string{"val_config:33", "val_result:33:1"}

	// Configurar mocks
	scanCmd := redis.NewScanCmd(ctx)
	scanCmd.SetVal([]string{keys[0]}, 0)
	cache.On("Scan", ctx, uint64(0), PrefijoValidationConfig+"*", int64(100)).Return(scanCmd)

	scanCmd2 := redis.NewScanCmd(ctx)
	scanCmd2.SetVal([]string{keys[1]}, 0)
	cache.On("Scan", ctx, uint64(0), PrefijoValidationResult+"*", int64(100)).Return(scanCmd2)

	intCmd := redis.NewIntCmd(ctx)
	cache.On("Del", ctx, keys).Return(intCmd)

	// Ejecutar prueba
	err := service.LimpiarCache(ctx)

	// Verificar resultado
	assert.NoError(t, err)
	cache.AssertExpectations(t)
}

func TestDTEValidationService_ApplySuggestions(t *testing.T) {
	// Preparar mock
	cache := new(MockRedisClient)
	service := NewDTEValidationService(cache)

	// Preparar datos de prueba
	ctx := context.Background()
	doc := &models.DocumentoTributario{
		TipoDTE: "33",
		Folio:   1,
	}
	suggestions := []*models.Suggestion{
		{
			Campo: "RUT",
			Valor: "12345678-9",
		},
	}

	// Configurar mock de Redis para invalidar caché
	intCmd := redis.NewIntCmd(ctx)
	cache.On("Del", ctx, mock.Anything).Return(intCmd)

	// Ejecutar prueba
	err := service.ApplySuggestions(ctx, doc, suggestions)

	// Verificar resultado
	assert.NoError(t, err)
	cache.AssertExpectations(t)
}
