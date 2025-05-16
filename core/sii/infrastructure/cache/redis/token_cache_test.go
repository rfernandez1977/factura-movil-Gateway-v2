package redis

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestRedis(t *testing.T) (*RedisTokenCache, func()) {
	mr, err := miniredis.Run()
	require.NoError(t, err)

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	cache := NewRedisTokenCache(client, "test")

	cleanup := func() {
		client.Close()
		mr.Close()
	}

	return cache, cleanup
}

func TestRedisTokenCache_SetAndGetToken(t *testing.T) {
	cache, cleanup := setupTestRedis(t)
	defer cleanup()

	ctx := context.Background()
	testKey := "test-key"
	testToken := "test-token"
	expiration := time.Hour

	// Probar almacenar token
	err := cache.SetToken(ctx, testKey, testToken, expiration)
	assert.NoError(t, err)

	// Probar obtener token
	token, err := cache.GetToken(ctx, testKey)
	assert.NoError(t, err)
	assert.Equal(t, testToken, token)
}

func TestRedisTokenCache_TokenExpiration(t *testing.T) {
	cache, cleanup := setupTestRedis(t)
	defer cleanup()

	ctx := context.Background()
	testKey := "test-key"
	testToken := "test-token"
	expiration := time.Second

	// Almacenar token con expiración corta
	err := cache.SetToken(ctx, testKey, testToken, expiration)
	assert.NoError(t, err)

	// Verificar que el token existe
	token, err := cache.GetToken(ctx, testKey)
	assert.NoError(t, err)
	assert.Equal(t, testToken, token)

	// Esperar a que expire
	time.Sleep(2 * time.Second)

	// Verificar que el token ha expirado
	token, err = cache.GetToken(ctx, testKey)
	assert.NoError(t, err)
	assert.Empty(t, token)
}

func TestRedisTokenCache_DeleteToken(t *testing.T) {
	cache, cleanup := setupTestRedis(t)
	defer cleanup()

	ctx := context.Background()
	testKey := "test-key"
	testToken := "test-token"
	expiration := time.Hour

	// Almacenar token
	err := cache.SetToken(ctx, testKey, testToken, expiration)
	assert.NoError(t, err)

	// Eliminar token
	err = cache.DeleteToken(ctx, testKey)
	assert.NoError(t, err)

	// Verificar que el token fue eliminado
	token, err := cache.GetToken(ctx, testKey)
	assert.NoError(t, err)
	assert.Empty(t, token)
}

func TestRedisTokenCache_NonExistentToken(t *testing.T) {
	cache, cleanup := setupTestRedis(t)
	defer cleanup()

	ctx := context.Background()
	testKey := "non-existent-key"

	// Intentar obtener un token que no existe
	token, err := cache.GetToken(ctx, testKey)
	assert.NoError(t, err)
	assert.Empty(t, token)
}

func TestRedisTokenCache_InvalidJSON(t *testing.T) {
	cache, cleanup := setupTestRedis(t)
	defer cleanup()

	ctx := context.Background()
	testKey := "test-key"

	// Almacenar datos inválidos directamente en Redis
	fullKey := cache.getFullKey(testKey)
	err := cache.client.Set(ctx, fullKey, "invalid-json", time.Hour).Err()
	assert.NoError(t, err)

	// Intentar obtener el token con datos inválidos
	token, err := cache.GetToken(ctx, testKey)
	assert.Error(t, err)
	assert.Empty(t, token)
}

func TestRedisTokenCache_SaveToken(t *testing.T) {
	cache, cleanup := setupTestRedis(t)
	defer cleanup()

	ctx := context.Background()
	now := time.Now()

	info := &cache.TokenInfo{
		Token:      "test-token",
		CreatedAt:  now,
		ExpiresAt:  now.Add(time.Hour),
		RutEmisor:  "12345678-9",
		Ambiente:   "CERTIFICACION",
		LastUsedAt: now,
	}

	// Guardar token
	err := cache.SaveToken(ctx, info)
	assert.NoError(t, err)

	// Verificar que se guardó correctamente
	token, err := cache.GetToken(ctx, info.RutEmisor)
	assert.NoError(t, err)
	assert.Equal(t, info.Token, token)
}

func TestRedisTokenCache_InvalidateToken(t *testing.T) {
	cache, cleanup := setupTestRedis(t)
	defer cleanup()

	ctx := context.Background()
	rutEmisor := "12345678-9"
	testToken := "test-token"
	expiration := time.Hour

	// Almacenar token
	err := cache.SetToken(ctx, rutEmisor, testToken, expiration)
	assert.NoError(t, err)

	// Invalidar token
	err = cache.InvalidateToken(ctx, rutEmisor)
	assert.NoError(t, err)

	// Verificar que el token fue invalidado
	token, err := cache.GetToken(ctx, rutEmisor)
	assert.NoError(t, err)
	assert.Empty(t, token)
}

func TestRedisTokenCache_UpdateLastUsed(t *testing.T) {
	cache, cleanup := setupTestRedis(t)
	defer cleanup()

	ctx := context.Background()
	now := time.Now()

	info := &cache.TokenInfo{
		Token:      "test-token",
		CreatedAt:  now,
		ExpiresAt:  now.Add(time.Hour),
		RutEmisor:  "12345678-9",
		Ambiente:   "CERTIFICACION",
		LastUsedAt: now,
	}

	// Guardar token inicial
	err := cache.SaveToken(ctx, info)
	assert.NoError(t, err)

	// Esperar un momento
	time.Sleep(time.Millisecond * 100)

	// Actualizar último uso
	err = cache.UpdateLastUsed(ctx, info.RutEmisor)
	assert.NoError(t, err)

	// Verificar que el token sigue siendo válido
	token, err := cache.GetToken(ctx, info.RutEmisor)
	assert.NoError(t, err)
	assert.Equal(t, info.Token, token)
}

func TestRedisTokenCache_UpdateLastUsed_NonExistent(t *testing.T) {
	cache, cleanup := setupTestRedis(t)
	defer cleanup()

	ctx := context.Background()
	rutEmisor := "non-existent"

	// Intentar actualizar último uso de un token que no existe
	err := cache.UpdateLastUsed(ctx, rutEmisor)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "token no encontrado")
}

func TestRedisTokenCache_CleanExpired(t *testing.T) {
	cache, cleanup := setupTestRedis(t)
	defer cleanup()

	ctx := context.Background()

	// CleanExpired no hace nada ya que Redis maneja la expiración automáticamente
	err := cache.CleanExpired(ctx)
	assert.NoError(t, err)
}
