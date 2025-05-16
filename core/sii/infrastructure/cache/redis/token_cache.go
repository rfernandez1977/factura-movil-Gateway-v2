package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/fmgo/core/sii/infrastructure/cache"
	"github.com/go-redis/redis/v8"
)

// TokenData representa la estructura de datos del token en caché
type TokenData struct {
	Token      string    `json:"token"`
	CreatedAt  time.Time `json:"created_at"`
	ExpiresAt  time.Time `json:"expires_at"`
	RutEmisor  string    `json:"rut_emisor"`
	Ambiente   string    `json:"ambiente"`
	LastUsedAt time.Time `json:"last_used_at"`
}

// TokenInfo representa la información del token almacenada en caché
type TokenInfo struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

// RedisTokenCache implementa la interfaz TokenCache usando Redis
type RedisTokenCache struct {
	client *redis.Client
	prefix string
	config cache.TokenCacheConfig
}

// NewRedisTokenCache crea una nueva instancia de RedisTokenCache
func NewRedisTokenCache(client *redis.Client, prefix string) *RedisTokenCache {
	return &RedisTokenCache{
		client: client,
		prefix: prefix,
		config: cache.DefaultTokenCacheConfig(),
	}
}

// GetToken obtiene un token del caché
func (c *RedisTokenCache) GetToken(ctx context.Context, key string) (string, error) {
	fullKey := c.getFullKey(key)
	data, err := c.client.Get(ctx, fullKey).Bytes()
	if err != nil {
		if err == redis.Nil {
			return "", nil
		}
		return "", fmt.Errorf("error al obtener token de redis: %w", err)
	}

	var tokenData TokenData
	if err := json.Unmarshal(data, &tokenData); err != nil {
		return "", fmt.Errorf("error al deserializar token: %w", err)
	}

	// Verificar si el token ha expirado
	if time.Now().After(tokenData.ExpiresAt) {
		_ = c.client.Del(ctx, fullKey)
		return "", nil
	}

	return tokenData.Token, nil
}

// SetToken almacena un token en el caché
func (c *RedisTokenCache) SetToken(ctx context.Context, key string, token string, expiration time.Duration) error {
	tokenData := TokenData{
		Token:     token,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(expiration),
	}

	data, err := json.Marshal(tokenData)
	if err != nil {
		return fmt.Errorf("error al serializar token: %w", err)
	}

	fullKey := c.getFullKey(key)
	if err := c.client.Set(ctx, fullKey, data, expiration).Err(); err != nil {
		return fmt.Errorf("error al almacenar token en redis: %w", err)
	}

	return nil
}

// SaveToken guarda un token en el caché
func (c *RedisTokenCache) SaveToken(ctx context.Context, info *cache.TokenInfo) error {
	tokenData := TokenData{
		Token:      info.Token,
		CreatedAt:  info.CreatedAt,
		ExpiresAt:  info.ExpiresAt,
		RutEmisor:  info.RutEmisor,
		Ambiente:   info.Ambiente,
		LastUsedAt: info.LastUsedAt,
	}

	data, err := json.Marshal(tokenData)
	if err != nil {
		return fmt.Errorf("error al serializar token: %w", err)
	}

	fullKey := c.getFullKey(info.RutEmisor)
	expiration := info.ExpiresAt.Sub(time.Now())
	if err := c.client.Set(ctx, fullKey, data, expiration).Err(); err != nil {
		return fmt.Errorf("error al almacenar token en redis: %w", err)
	}

	return nil
}

// DeleteToken elimina un token del caché
func (c *RedisTokenCache) DeleteToken(ctx context.Context, key string) error {
	fullKey := c.getFullKey(key)
	if err := c.client.Del(ctx, fullKey).Err(); err != nil {
		return fmt.Errorf("error al eliminar token de redis: %w", err)
	}
	return nil
}

// InvalidateToken invalida un token específico
func (c *RedisTokenCache) InvalidateToken(ctx context.Context, rutEmisor string) error {
	return c.DeleteToken(ctx, rutEmisor)
}

// CleanExpired limpia los tokens expirados
func (c *RedisTokenCache) CleanExpired(ctx context.Context) error {
	// Redis maneja automáticamente la expiración de las claves
	return nil
}

// UpdateLastUsed actualiza la marca de último uso de un token
func (c *RedisTokenCache) UpdateLastUsed(ctx context.Context, rutEmisor string) error {
	fullKey := c.getFullKey(rutEmisor)
	data, err := c.client.Get(ctx, fullKey).Bytes()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("token no encontrado")
		}
		return fmt.Errorf("error al obtener token de redis: %w", err)
	}

	var tokenData TokenData
	if err := json.Unmarshal(data, &tokenData); err != nil {
		return fmt.Errorf("error al deserializar token: %w", err)
	}

	tokenData.LastUsedAt = time.Now()

	updatedData, err := json.Marshal(tokenData)
	if err != nil {
		return fmt.Errorf("error al serializar token: %w", err)
	}

	expiration := tokenData.ExpiresAt.Sub(time.Now())
	if err := c.client.Set(ctx, fullKey, updatedData, expiration).Err(); err != nil {
		return fmt.Errorf("error al actualizar token en redis: %w", err)
	}

	return nil
}

// getFullKey genera la clave completa para Redis
func (c *RedisTokenCache) getFullKey(key string) string {
	return fmt.Sprintf("%s:token:%s", c.prefix, key)
}

// TokenCache implementa el almacenamiento de tokens en Redis
type TokenCache struct {
	client *redis.Client
	ttl    time.Duration
}

// NewTokenCache crea una nueva instancia de TokenCache
func NewTokenCache(client *redis.Client, ttl time.Duration) *TokenCache {
	return &TokenCache{
		client: client,
		ttl:    ttl,
	}
}

// Set almacena un token en el caché
func (c *TokenCache) Set(ctx context.Context, key string, token *TokenInfo) error {
	data, err := json.Marshal(token)
	if err != nil {
		return fmt.Errorf("error al serializar token: %v", err)
	}

	return c.client.Set(ctx, key, data, c.ttl).Err()
}

// Get obtiene un token del caché
func (c *TokenCache) Get(ctx context.Context, key string) (*TokenInfo, error) {
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, fmt.Errorf("error al obtener token: %v", err)
	}

	var token TokenInfo
	if err := json.Unmarshal(data, &token); err != nil {
		return nil, fmt.Errorf("error al deserializar token: %v", err)
	}

	return &token, nil
}

// Delete elimina un token del caché
func (c *TokenCache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}
