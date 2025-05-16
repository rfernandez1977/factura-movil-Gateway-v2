package retry

import (
	"context"
	"time"
)

// Config representa la configuración para los reintentos
type Config struct {
	MaxRetries  int
	InitialWait time.Duration
	MaxWait     time.Duration
}

// DefaultConfig retorna una configuración por defecto
func DefaultConfig() *Config {
	return &Config{
		MaxRetries:  3,
		InitialWait: 1 * time.Second,
		MaxWait:     30 * time.Second,
	}
}

// Do ejecuta una función con reintentos según la configuración
func Do(ctx context.Context, cfg *Config, fn func() error) error {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	var err error
	wait := cfg.InitialWait

	for attempt := 0; attempt <= cfg.MaxRetries; attempt++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if err = fn(); err == nil {
				return nil
			}

			if attempt == cfg.MaxRetries {
				return err
			}

			time.Sleep(wait)
			wait = min(wait*2, cfg.MaxWait)
		}
	}

	return err
}

func min(a, b time.Duration) time.Duration {
	if a < b {
		return a
	}
	return b
}
