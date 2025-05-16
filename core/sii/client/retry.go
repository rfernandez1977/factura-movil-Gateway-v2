package client

import (
	"context"
	"fmt"
	"time"

	"github.com/fmgo/core/sii/models"
)

// RetryConfig define la configuración para reintentos
type RetryConfig struct {
	MaxRetries      int
	InitialInterval time.Duration
	MaxInterval     time.Duration
	Multiplier      float64
}

// DefaultRetryConfig retorna la configuración por defecto para reintentos
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxRetries:      3,
		InitialInterval: 1 * time.Second,
		MaxInterval:     30 * time.Second,
		Multiplier:      2.0,
	}
}

// isRetryableError determina si un error puede ser reintentado
func isRetryableError(err error) bool {
	if siiErr, ok := err.(*models.SIIError); ok {
		return siiErr.IsRetryable()
	}
	// Por defecto, reintentamos errores desconocidos
	return true
}

// withRetry ejecuta una operación con reintentos
func withRetry[T any](ctx context.Context, config RetryConfig, operation func() (T, error)) (T, error) {
	var lastErr error
	var result T

	interval := config.InitialInterval
	for attempt := 0; attempt <= config.MaxRetries; attempt++ {
		select {
		case <-ctx.Done():
			return result, ctx.Err()
		default:
		}

		if attempt > 0 {
			// Si el error no es reintentable, fallamos inmediatamente
			if !isRetryableError(lastErr) {
				return result, lastErr
			}

			time.Sleep(interval)
			// Incrementar el intervalo para el próximo intento
			interval = time.Duration(float64(interval) * config.Multiplier)
			if interval > config.MaxInterval {
				interval = config.MaxInterval
			}
		}

		result, err := operation()
		if err == nil {
			return result, nil
		}

		lastErr = err
	}

	if !isRetryableError(lastErr) {
		return result, lastErr
	}

	return result, fmt.Errorf("después de %d intentos: %v", config.MaxRetries, lastErr)
}
