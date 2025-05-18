package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

// Service define la interfaz para el servicio de caché
type Service interface {
	// Get obtiene un valor del caché
	Get(ctx context.Context, key string, result interface{}) error

	// Set almacena un valor en el caché
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error

	// Delete elimina un valor del caché
	Delete(ctx context.Context, key string) error

	// Exists verifica si una clave existe en el caché
	Exists(ctx context.Context, key string) (bool, error)

	// Clear limpia todas las claves que coincidan con el patrón
	Clear(ctx context.Context, pattern string) error

	// GetClient retorna el cliente Redis subyacente
	GetClient() *redis.Client
}

// RedisService implementa el servicio de caché usando Redis
type RedisService struct {
	client *redis.Client
	logger *zap.Logger
}

// NewRedisService crea una nueva instancia del servicio de caché Redis
func NewRedisService(client *redis.Client, logger *zap.Logger) Service {
	return &RedisService{
		client: client,
		logger: logger,
	}
}

// Get obtiene un valor del caché
func (s *RedisService) Get(ctx context.Context, key string, result interface{}) error {
	data, err := s.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("clave no encontrada: %s", key)
		}
		return fmt.Errorf("error obteniendo valor de Redis: %w", err)
	}

	if err := json.Unmarshal(data, result); err != nil {
		return fmt.Errorf("error deserializando valor: %w", err)
	}

	s.logger.Debug("Valor obtenido del caché",
		zap.String("key", key),
		zap.Int("size", len(data)))

	return nil
}

// Set almacena un valor en el caché
func (s *RedisService) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("error serializando valor: %w", err)
	}

	if err := s.client.Set(ctx, key, data, expiration).Err(); err != nil {
		return fmt.Errorf("error almacenando valor en Redis: %w", err)
	}

	s.logger.Debug("Valor almacenado en caché",
		zap.String("key", key),
		zap.Int("size", len(data)),
		zap.Duration("expiration", expiration))

	return nil
}

// Delete elimina un valor del caché
func (s *RedisService) Delete(ctx context.Context, key string) error {
	if err := s.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("error eliminando valor de Redis: %w", err)
	}

	s.logger.Debug("Valor eliminado del caché",
		zap.String("key", key))

	return nil
}

// Exists verifica si una clave existe en el caché
func (s *RedisService) Exists(ctx context.Context, key string) (bool, error) {
	result, err := s.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("error verificando existencia en Redis: %w", err)
	}

	exists := result > 0
	s.logger.Debug("Verificación de existencia en caché",
		zap.String("key", key),
		zap.Bool("exists", exists))

	return exists, nil
}

// Clear limpia todas las claves que coincidan con el patrón
func (s *RedisService) Clear(ctx context.Context, pattern string) error {
	iter := s.client.Scan(ctx, 0, pattern, 0).Iterator()

	var deletedKeys int
	for iter.Next(ctx) {
		if err := s.client.Del(ctx, iter.Val()).Err(); err != nil {
			return fmt.Errorf("error eliminando clave %s: %w", iter.Val(), err)
		}
		deletedKeys++
	}

	if err := iter.Err(); err != nil {
		return fmt.Errorf("error iterando claves: %w", err)
	}

	s.logger.Debug("Limpieza de caché completada",
		zap.String("pattern", pattern),
		zap.Int("deleted_keys", deletedKeys))

	return nil
}

// GetClient retorna el cliente Redis subyacente
func (s *RedisService) GetClient() *redis.Client {
	return s.client
}
