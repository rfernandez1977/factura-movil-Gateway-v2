package cache

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type testStruct struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func setupTestRedis(t *testing.T) (Service, func()) {
	mr, err := miniredis.Run()
	require.NoError(t, err)

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	logger, _ := zap.NewDevelopment()
	service := NewRedisService(client, logger)

	cleanup := func() {
		client.Close()
		mr.Close()
		logger.Sync()
	}

	return service, cleanup
}

func TestRedisService_SetAndGet(t *testing.T) {
	service, cleanup := setupTestRedis(t)
	defer cleanup()

	ctx := context.Background()
	key := "test:key"
	value := testStruct{
		ID:   "123",
		Name: "Test",
		Age:  30,
	}

	// Probar Set
	err := service.Set(ctx, key, value, time.Hour)
	assert.NoError(t, err)

	// Probar Get
	var result testStruct
	err = service.Get(ctx, key, &result)
	assert.NoError(t, err)
	assert.Equal(t, value, result)
}

func TestRedisService_GetNonExistent(t *testing.T) {
	service, cleanup := setupTestRedis(t)
	defer cleanup()

	ctx := context.Background()
	key := "non:existent"

	var result testStruct
	err := service.Get(ctx, key, &result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "clave no encontrada")
}

func TestRedisService_Delete(t *testing.T) {
	service, cleanup := setupTestRedis(t)
	defer cleanup()

	ctx := context.Background()
	key := "test:key"
	value := testStruct{
		ID:   "123",
		Name: "Test",
		Age:  30,
	}

	// Almacenar valor
	err := service.Set(ctx, key, value, time.Hour)
	assert.NoError(t, err)

	// Eliminar valor
	err = service.Delete(ctx, key)
	assert.NoError(t, err)

	// Verificar que fue eliminado
	var result testStruct
	err = service.Get(ctx, key, &result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "clave no encontrada")
}

func TestRedisService_Exists(t *testing.T) {
	service, cleanup := setupTestRedis(t)
	defer cleanup()

	ctx := context.Background()
	key := "test:key"
	value := testStruct{
		ID:   "123",
		Name: "Test",
		Age:  30,
	}

	// Verificar que no existe
	exists, err := service.Exists(ctx, key)
	assert.NoError(t, err)
	assert.False(t, exists)

	// Almacenar valor
	err = service.Set(ctx, key, value, time.Hour)
	assert.NoError(t, err)

	// Verificar que existe
	exists, err = service.Exists(ctx, key)
	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestRedisService_Clear(t *testing.T) {
	service, cleanup := setupTestRedis(t)
	defer cleanup()

	ctx := context.Background()
	value := testStruct{
		ID:   "123",
		Name: "Test",
		Age:  30,
	}

	// Almacenar múltiples valores
	for i := 0; i < 5; i++ {
		key := fmt.Sprintf("test:key:%d", i)
		err := service.Set(ctx, key, value, time.Hour)
		assert.NoError(t, err)
	}

	// Almacenar otros valores con diferente prefijo
	for i := 0; i < 3; i++ {
		key := fmt.Sprintf("other:key:%d", i)
		err := service.Set(ctx, key, value, time.Hour)
		assert.NoError(t, err)
	}

	// Limpiar solo las claves con prefijo "test:"
	err := service.Clear(ctx, "test:*")
	assert.NoError(t, err)

	// Verificar que las claves "test:" fueron eliminadas
	for i := 0; i < 5; i++ {
		key := fmt.Sprintf("test:key:%d", i)
		exists, err := service.Exists(ctx, key)
		assert.NoError(t, err)
		assert.False(t, exists)
	}

	// Verificar que las claves "other:" siguen existiendo
	for i := 0; i < 3; i++ {
		key := fmt.Sprintf("other:key:%d", i)
		exists, err := service.Exists(ctx, key)
		assert.NoError(t, err)
		assert.True(t, exists)
	}
}

func TestRedisService_Expiration(t *testing.T) {
	service, cleanup := setupTestRedis(t)
	defer cleanup()

	ctx := context.Background()
	key := "test:key"
	value := testStruct{
		ID:   "123",
		Name: "Test",
		Age:  30,
	}

	// Almacenar valor con expiración corta
	err := service.Set(ctx, key, value, time.Second)
	assert.NoError(t, err)

	// Verificar que existe
	exists, err := service.Exists(ctx, key)
	assert.NoError(t, err)
	assert.True(t, exists)

	// Esperar a que expire
	time.Sleep(2 * time.Second)

	// Verificar que ya no existe
	exists, err = service.Exists(ctx, key)
	assert.NoError(t, err)
	assert.False(t, exists)
}
