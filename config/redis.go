package config

import (
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisConfig contiene la configuración para Redis
type RedisConfig struct {
	Host            string        `json:"host"`
	Port            int           `json:"port"`
	Password        string        `json:"password"`
	DB              int           `json:"db"`
	MaxRetries      int           `json:"max_retries"`
	MinRetryBackoff time.Duration `json:"min_retry_backoff"`
	MaxRetryBackoff time.Duration `json:"max_retry_backoff"`
	DialTimeout     time.Duration `json:"dial_timeout"`
	ReadTimeout     time.Duration `json:"read_timeout"`
	WriteTimeout    time.Duration `json:"write_timeout"`
	PoolSize        int           `json:"pool_size"`
	MinIdleConns    int           `json:"min_idle_conns"`
	MaxConnAge      time.Duration `json:"max_conn_age"`
	PoolTimeout     time.Duration `json:"pool_timeout"`
	IdleTimeout     time.Duration `json:"idle_timeout"`
}

// DefaultRedisConfig retorna la configuración por defecto para Redis
func DefaultRedisConfig() *RedisConfig {
	return &RedisConfig{
		Host:            "localhost",
		Port:            6379,
		Password:        "",
		DB:              0,
		MaxRetries:      3,
		MinRetryBackoff: time.Millisecond * 100,
		MaxRetryBackoff: time.Second * 2,
		DialTimeout:     time.Second * 5,
		ReadTimeout:     time.Second * 3,
		WriteTimeout:    time.Second * 3,
		PoolSize:        10,
		MinIdleConns:    2,
		MaxConnAge:      time.Hour,
		PoolTimeout:     time.Second * 4,
		IdleTimeout:     time.Minute * 5,
	}
}

// NewRedisClient crea un nuevo cliente Redis con la configuración proporcionada
func NewRedisClient(config *RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:            fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password:        config.Password,
		DB:              config.DB,
		MaxRetries:      config.MaxRetries,
		MinRetryBackoff: config.MinRetryBackoff,
		MaxRetryBackoff: config.MaxRetryBackoff,
		DialTimeout:     config.DialTimeout,
		ReadTimeout:     config.ReadTimeout,
		WriteTimeout:    config.WriteTimeout,
		PoolSize:        config.PoolSize,
		MinIdleConns:    config.MinIdleConns,
		MaxConnAge:      config.MaxConnAge,
		PoolTimeout:     config.PoolTimeout,
		IdleTimeout:     config.IdleTimeout,
	})

	// Verificar conexión
	if err := client.Ping(client.Context()).Err(); err != nil {
		return nil, fmt.Errorf("error conectando a Redis: %w", err)
	}

	return client, nil
}
