package retry

import (
	"context"
	"fmt"
	"time"
)

// Config contiene la configuración para los reintentos
type Config struct {
	MaxRetries  int
	WaitTime    time.Duration
	MaxWaitTime time.Duration
}

// DefaultConfig retorna una configuración por defecto
func DefaultConfig() *Config {
	return &Config{
		MaxRetries:  3,
		WaitTime:    time.Second,
		MaxWaitTime: time.Second * 10,
	}
}

// Do ejecuta una función con reintentos según la configuración
func Do(ctx context.Context, config *Config, fn func() error) error {
	var lastErr error
	waitTime := config.WaitTime

	for i := 0; i <= config.MaxRetries; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if err := fn(); err == nil {
			return nil
		} else {
			lastErr = err
		}

		if i == config.MaxRetries {
			break
		}

		// Esperar antes del siguiente reintento
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(waitTime):
		}

		// Incrementar el tiempo de espera para el siguiente reintento
		waitTime *= 2
		if waitTime > config.MaxWaitTime {
			waitTime = config.MaxWaitTime
		}
	}

	return fmt.Errorf("máximo número de reintentos alcanzado: %w", lastErr)
}
