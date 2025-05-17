package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	// DefaultExpiration es el tiempo de expiración por defecto para las entradas en caché
	DefaultExpiration = 24 * time.Hour

	// KeyPrefix es el prefijo para todas las claves en Redis
	KeyPrefix = "fmgo:"
)

// RedisCache implementa el caché usando Redis
type RedisCache struct {
	client *redis.Client
	ctx    context.Context
}

// Config contiene la configuración para conectar a Redis
type Config struct {
	Host     string
	Port     int
	Password string
	DB       int
}

// NewRedisCache crea una nueva instancia de RedisCache
func NewRedisCache(cfg Config) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx := context.Background()

	// Verificar conexión
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("error conectando a Redis: %v", err)
	}

	return &RedisCache{
		client: client,
		ctx:    ctx,
	}, nil
}

// Set almacena un valor en el caché
func (c *RedisCache) Set(key string, value interface{}, expiration time.Duration) error {
	if expiration == 0 {
		expiration = DefaultExpiration
	}

	// Convertir valor a JSON
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("error codificando valor: %v", err)
	}

	// Almacenar en Redis
	key = KeyPrefix + key
	if err := c.client.Set(c.ctx, key, data, expiration).Err(); err != nil {
		return fmt.Errorf("error almacenando en Redis: %v", err)
	}

	return nil
}

// Get recupera un valor del caché
func (c *RedisCache) Get(key string, value interface{}) error {
	key = KeyPrefix + key
	data, err := c.client.Get(c.ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("clave no encontrada: %s", key)
		}
		return fmt.Errorf("error leyendo de Redis: %v", err)
	}

	// Decodificar valor desde JSON
	if err := json.Unmarshal(data, value); err != nil {
		return fmt.Errorf("error decodificando valor: %v", err)
	}

	return nil
}

// Delete elimina una clave del caché
func (c *RedisCache) Delete(key string) error {
	key = KeyPrefix + key
	if err := c.client.Del(c.ctx, key).Err(); err != nil {
		return fmt.Errorf("error eliminando de Redis: %v", err)
	}
	return nil
}

// Clear elimina todas las claves con el prefijo del proyecto
func (c *RedisCache) Clear() error {
	pattern := KeyPrefix + "*"
	iter := c.client.Scan(c.ctx, 0, pattern, 0).Iterator()

	for iter.Next(c.ctx) {
		if err := c.client.Del(c.ctx, iter.Val()).Err(); err != nil {
			return fmt.Errorf("error limpiando caché: %v", err)
		}
	}

	if err := iter.Err(); err != nil {
		return fmt.Errorf("error escaneando claves: %v", err)
	}

	return nil
}

// Close cierra la conexión con Redis
func (c *RedisCache) Close() error {
	return c.client.Close()
}
