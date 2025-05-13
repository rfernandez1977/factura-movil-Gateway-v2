package mocks

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/mock"
)

// MockRedisClient implementa la interfaz de redis.Client para testing
type MockRedisClient struct {
	mock.Mock
}

// NewMockRedisClient crea una nueva instancia de MockRedisClient
func NewMockRedisClient() *MockRedisClient {
	return &MockRedisClient{}
}

// Get implementa el método Get de redis.Client
func (m *MockRedisClient) Get(ctx context.Context, key string) *redis.StringCmd {
	args := m.Called(ctx, key)
	return args.Get(0).(*redis.StringCmd)
}

// Set implementa el método Set de redis.Client
func (m *MockRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	args := m.Called(ctx, key, value, expiration)
	return args.Get(0).(*redis.StatusCmd)
}

// Del implementa el método Del de redis.Client
func (m *MockRedisClient) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	args := m.Called(ctx, keys)
	return args.Get(0).(*redis.IntCmd)
}

// Exists implementa el método Exists de redis.Client
func (m *MockRedisClient) Exists(ctx context.Context, keys ...string) *redis.IntCmd {
	args := m.Called(ctx, keys)
	return args.Get(0).(*redis.IntCmd)
}

// Close implementa el método Close de redis.Client
func (m *MockRedisClient) Close() error {
	args := m.Called()
	return args.Error(0)
}

// NewMockRedisClientWithDefaults crea una nueva instancia con respuestas por defecto
func NewMockRedisClientWithDefaults() *MockRedisClient {
	mock := &MockRedisClient{}

	// Configurar respuestas por defecto
	mock.On("Get", Anything, Anything).Return(redis.NewStringCmd(context.Background()))
	mock.On("Set", Anything, Anything, Anything, Anything).Return(redis.NewStatusCmd(context.Background()))
	mock.On("Del", Anything, Anything).Return(redis.NewIntCmd(context.Background()))
	mock.On("Exists", Anything, Anything).Return(redis.NewIntCmd(context.Background()))
	mock.On("Close").Return(nil)

	return mock
}
