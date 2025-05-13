package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

// CacheService maneja el almacenamiento en caché de documentos frecuentes
type CacheService struct {
	client *redis.Client
	logger *zap.Logger
}

type CacheConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

// NewCacheService crea una nueva instancia del servicio de caché
func NewCacheService(config CacheConfig, logger *zap.Logger) (*CacheService, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.Password,
		DB:       config.DB,
	})

	// Verificar conexión
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("error al conectar con Redis: %v", err)
	}

	return &CacheService{
		client: client,
		logger: logger,
	}, nil
}

// Get obtiene un documento de la caché
func (s *CacheService) Get(ctx context.Context, key string, result interface{}) error {
	val, err := s.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(val), result)
}

// Set almacena un documento en la caché
func (s *CacheService) Set(ctx context.Context, key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return s.client.Set(ctx, key, data, 0).Err()
}

// SetWithExpiration almacena un documento en la caché con tiempo de expiración
func (s *CacheService) SetWithExpiration(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return s.client.Set(ctx, key, data, expiration).Err()
}

// Delete elimina un documento de la caché
func (s *CacheService) Delete(ctx context.Context, key string) error {
	return s.client.Del(ctx, key).Err()
}

// GetOrSet obtiene un documento de la caché o lo carga desde la fuente
func (s *CacheService) GetOrSet(ctx context.Context, key string, result interface{}, loader func() (interface{}, error)) error {
	// Intentar obtener de la caché
	err := s.Get(ctx, key, result)
	if err == nil {
		return nil
	}

	// Cargar desde la fuente
	value, err := loader()
	if err != nil {
		return err
	}

	// Almacenar en caché
	err = s.Set(ctx, key, value)
	if err != nil {
		return err
	}

	// Asignar el valor al resultado
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, result)
}

// CacheKey genera una clave de caché para un documento
func (s *CacheService) CacheKey(collection string, id primitive.ObjectID) string {
	return collection + ":" + id.Hex()
}

// CacheKeys genera claves de caché para múltiples documentos
func (s *CacheService) CacheKeys(collection string, ids []primitive.ObjectID) []string {
	keys := make([]string, len(ids))
	for i, id := range ids {
		keys[i] = s.CacheKey(collection, id)
	}
	return keys
}

// BatchGet obtiene múltiples documentos de la caché
func (s *CacheService) BatchGet(ctx context.Context, keys []string, results interface{}) error {
	vals, err := s.client.MGet(ctx, keys...).Result()
	if err != nil {
		return err
	}

	// Filtrar valores nulos
	var data []string
	for _, val := range vals {
		if val != nil {
			data = append(data, val.(string))
		}
	}

	// Unmarshal de los resultados
	return json.Unmarshal([]byte("["+string(join(data, ","))+"]"), results)
}

// BatchSet almacena múltiples documentos en la caché
func (s *CacheService) BatchSet(ctx context.Context, items map[string]interface{}) error {
	pipe := s.client.Pipeline()
	for key, value := range items {
		data, err := json.Marshal(value)
		if err != nil {
			return err
		}
		pipe.Set(ctx, key, data, 0)
	}
	_, err := pipe.Exec(ctx)
	return err
}

// join une strings con un separador
func join(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	if len(strs) == 1 {
		return strs[0]
	}
	n := len(sep) * (len(strs) - 1)
	for i := 0; i < len(strs); i++ {
		n += len(strs[i])
	}

	b := make([]byte, n)
	bp := copy(b, strs[0])
	for _, s := range strs[1:] {
		bp += copy(b[bp:], sep)
		bp += copy(b[bp:], s)
	}
	return string(b)
}

// Increment incrementa un contador en el caché
func (s *CacheService) Increment(ctx context.Context, key string) (int64, error) {
	val, err := s.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("error al incrementar contador: %v", err)
	}
	return val, nil
}

// GetOrSetWithExpiration obtiene un valor del caché o lo establece si no existe con un tiempo de expiración
func (s *CacheService) GetOrSetWithExpiration(ctx context.Context, key string, value interface{}, expiration time.Duration, setter func() (interface{}, error)) error {
	// Intentar obtener del caché
	if err := s.Get(ctx, key, value); err != nil {
		return err
	}

	// Si el valor existe en caché, retornar
	if value != nil {
		return nil
	}

	// Obtener el valor usando la función setter
	newValue, err := setter()
	if err != nil {
		return err
	}

	// Guardar en caché
	if err := s.SetWithExpiration(ctx, key, newValue, expiration); err != nil {
		return err
	}

	// Asignar el valor
	value = newValue
	return nil
}

// ClearAll limpia todo el caché
func (s *CacheService) ClearAll(ctx context.Context) error {
	if err := s.client.FlushAll(ctx).Err(); err != nil {
		return fmt.Errorf("error al limpiar caché: %v", err)
	}
	return nil
}

// Close cierra la conexión con Redis
func (s *CacheService) Close() error {
	return s.client.Close()
}
