package cache

import (
	"context"
	"time"
)

// TokenInfo contiene la información del token y su estado
type TokenInfo struct {
	Token      string    `json:"token"`
	CreatedAt  time.Time `json:"created_at"`
	ExpiresAt  time.Time `json:"expires_at"`
	RutEmisor  string    `json:"rut_emisor"`
	Ambiente   string    `json:"ambiente"`
	LastUsedAt time.Time `json:"last_used_at"`
}

// TokenCache define la interfaz para el caché de tokens
type TokenCache interface {
	// GetToken obtiene un token del caché
	GetToken(ctx context.Context, key string) (string, error)

	// SetToken almacena un token en el caché con una duración de expiración
	SetToken(ctx context.Context, key string, token string, expiration time.Duration) error

	// DeleteToken elimina un token del caché
	DeleteToken(ctx context.Context, key string) error

	// SaveToken guarda un token en el caché
	SaveToken(ctx context.Context, info *TokenInfo) error

	// InvalidateToken invalida un token específico
	InvalidateToken(ctx context.Context, rutEmisor string) error

	// CleanExpired limpia los tokens expirados
	CleanExpired(ctx context.Context) error

	// UpdateLastUsed actualiza la marca de último uso de un token
	UpdateLastUsed(ctx context.Context, rutEmisor string) error
}

// TokenCacheConfig contiene la configuración para el caché de tokens
type TokenCacheConfig struct {
	// Duración por defecto de un token
	TokenDuration time.Duration

	// Tiempo antes de la expiración para renovar
	RenewalWindow time.Duration

	// Máximo de intentos de renovación
	MaxRenewalAttempts int

	// Intervalo de limpieza de tokens expirados
	CleanupInterval time.Duration
}

// DefaultTokenCacheConfig retorna una configuración por defecto
func DefaultTokenCacheConfig() TokenCacheConfig {
	return TokenCacheConfig{
		TokenDuration:      4 * time.Hour,
		RenewalWindow:      30 * time.Minute,
		MaxRenewalAttempts: 3,
		CleanupInterval:    1 * time.Hour,
	}
}
