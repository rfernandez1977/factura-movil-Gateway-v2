package cache

import (
	"context"
	"sync"
	"time"

	"github.com/fmgo/core/sii/logger"
)

// TokenInfo representa la información de un token
type TokenInfo struct {
	Token     string
	ExpiresAt time.Time
}

// TokenCache maneja el caché de tokens del SII
type TokenCache struct {
	tokens   map[string]*TokenInfo
	mu       sync.RWMutex
	logger   *logger.Logger
	lifetime time.Duration
}

// NewTokenCache crea una nueva instancia del caché de tokens
func NewTokenCache(logger *logger.Logger, lifetime time.Duration) *TokenCache {
	return &TokenCache{
		tokens:   make(map[string]*TokenInfo),
		logger:   logger,
		lifetime: lifetime,
	}
}

// Get obtiene un token del caché
func (c *TokenCache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	info, exists := c.tokens[key]
	if !exists {
		return "", false
	}

	// Verificar si el token ha expirado
	if time.Now().After(info.ExpiresAt) {
		c.logger.Debug("Token expirado para la clave: %s", key)
		return "", false
	}

	return info.Token, true
}

// Set almacena un token en el caché
func (c *TokenCache) Set(key string, token string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.tokens[key] = &TokenInfo{
		Token:     token,
		ExpiresAt: time.Now().Add(c.lifetime),
	}
	c.logger.Debug("Token almacenado para la clave: %s", key)
}

// Delete elimina un token del caché
func (c *TokenCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.tokens, key)
	c.logger.Debug("Token eliminado para la clave: %s", key)
}

// StartCleanup inicia la limpieza periódica de tokens expirados
func (c *TokenCache) StartCleanup(ctx context.Context, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				c.cleanup()
			}
		}
	}()
}

// cleanup elimina los tokens expirados
func (c *TokenCache) cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for key, info := range c.tokens {
		if now.After(info.ExpiresAt) {
			delete(c.tokens, key)
			c.logger.Debug("Token expirado eliminado para la clave: %s", key)
		}
	}
}
