package client

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWithRetry_Success(t *testing.T) {
	config := RetryConfig{
		MaxRetries:      3,
		InitialInterval: 10 * time.Millisecond,
		MaxInterval:     100 * time.Millisecond,
		Multiplier:      2.0,
	}

	attempts := 0
	result, err := withRetry(context.Background(), config, func() (string, error) {
		attempts++
		return "success", nil
	})

	assert.NoError(t, err)
	assert.Equal(t, "success", result)
	assert.Equal(t, 1, attempts, "Should succeed on first attempt")
}

func TestWithRetry_EventualSuccess(t *testing.T) {
	config := RetryConfig{
		MaxRetries:      3,
		InitialInterval: 10 * time.Millisecond,
		MaxInterval:     100 * time.Millisecond,
		Multiplier:      2.0,
	}

	attempts := 0
	result, err := withRetry(context.Background(), config, func() (string, error) {
		attempts++
		if attempts < 3 {
			return "", errors.New("error temporal")
		}
		return "success", nil
	})

	assert.NoError(t, err)
	assert.Equal(t, "success", result)
	assert.Equal(t, 3, attempts, "Should succeed on third attempt")
}

func TestWithRetry_Failure(t *testing.T) {
	config := RetryConfig{
		MaxRetries:      2,
		InitialInterval: 10 * time.Millisecond,
		MaxInterval:     100 * time.Millisecond,
		Multiplier:      2.0,
	}

	attempts := 0
	_, err := withRetry(context.Background(), config, func() (string, error) {
		attempts++
		return "", errors.New("error persistente")
	})

	assert.Error(t, err)
	assert.Equal(t, 3, attempts, "Should try MaxRetries + 1 times")
	assert.Contains(t, err.Error(), "después de 2 intentos")
}

func TestWithRetry_ContextCancellation(t *testing.T) {
	config := RetryConfig{
		MaxRetries:      3,
		InitialInterval: 10 * time.Millisecond,
		MaxInterval:     100 * time.Millisecond,
		Multiplier:      2.0,
	}

	ctx, cancel := context.WithCancel(context.Background())
	attempts := 0

	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	_, err := withRetry(ctx, config, func() (string, error) {
		attempts++
		time.Sleep(20 * time.Millisecond)
		return "", errors.New("error temporal")
	})

	assert.Error(t, err)
	assert.True(t, errors.Is(err, context.Canceled))
}
