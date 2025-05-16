package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/fmgo/core/firma/common"
	"github.com/fmgo/core/firma/models"
	"github.com/go-redis/redis/v8"
)

const (
	// TTLCaf tiempo de vida del CAF en caché
	TTLCaf = 24 * time.Hour
	// PrefijoCaf prefijo para las claves de CAF
	PrefijoCaf = "caf:"
	// PrefijoTipo prefijo para las claves de tipo de documento
	PrefijoTipo = "tipo:"
)

// CAFCache implementa el caché de CAF usando Redis
type CAFCache struct {
	client *redis.Client
	logger common.Logger
}

// NewCAFCache crea una nueva instancia del caché de CAF
func NewCAFCache(client *redis.Client, logger common.Logger) *CAFCache {
	return &CAFCache{
		client: client,
		logger: logger,
	}
}

// GuardarCAF guarda un CAF en el caché
func (c *CAFCache) GuardarCAF(ctx context.Context, caf *models.CAF) error {
	// Clave para el CAF individual
	key := PrefijoCaf + caf.ID

	// Serializar
	data, err := json.Marshal(caf)
	if err != nil {
		return fmt.Errorf("error serializando CAF: %w", err)
	}

	// Guardar CAF individual
	if err := c.client.Set(ctx, key, data, TTLCaf).Err(); err != nil {
		return fmt.Errorf("error guardando CAF en caché: %w", err)
	}

	// Clave para el conjunto de CAFs por tipo
	tipoKey := PrefijoTipo + caf.TipoDocumento

	// Agregar al conjunto de CAFs por tipo
	if err := c.client.SAdd(ctx, tipoKey, caf.ID).Err(); err != nil {
		return fmt.Errorf("error agregando CAF al conjunto por tipo: %w", err)
	}

	// Establecer TTL para el conjunto
	if err := c.client.Expire(ctx, tipoKey, TTLCaf).Err(); err != nil {
		c.logger.Warn("Error estableciendo TTL para conjunto de tipo",
			"tipo", caf.TipoDocumento,
			"error", err)
	}

	c.logger.Debug("CAF guardado en caché",
		"id", caf.ID,
		"tipo", caf.TipoDocumento,
		"ttl", TTLCaf)

	return nil
}

// ObtenerCAF obtiene un CAF del caché
func (c *CAFCache) ObtenerCAF(ctx context.Context, id string) (*models.CAF, error) {
	key := PrefijoCaf + id

	// Obtener del caché
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("CAF no encontrado en caché: %s", id)
		}
		return nil, fmt.Errorf("error obteniendo CAF del caché: %w", err)
	}

	// Deserializar
	var caf models.CAF
	if err := json.Unmarshal(data, &caf); err != nil {
		return nil, fmt.Errorf("error deserializando CAF: %w", err)
	}

	c.logger.Debug("CAF obtenido del caché",
		"id", id,
		"tipo", caf.TipoDocumento)

	return &caf, nil
}

// ListarCAFsPorTipo lista todos los CAFs de un tipo de documento desde el caché
func (c *CAFCache) ListarCAFsPorTipo(ctx context.Context, tipo string) ([]*models.CAF, error) {
	tipoKey := PrefijoTipo + tipo

	// Obtener IDs del conjunto
	ids, err := c.client.SMembers(ctx, tipoKey).Result()
	if err != nil {
		return nil, fmt.Errorf("error obteniendo IDs de CAF por tipo: %w", err)
	}

	if len(ids) == 0 {
		return nil, fmt.Errorf("no hay CAFs en caché para tipo: %s", tipo)
	}

	// Obtener cada CAF
	var cafs []*models.CAF
	for _, id := range ids {
		caf, err := c.ObtenerCAF(ctx, id)
		if err != nil {
			c.logger.Warn("Error obteniendo CAF del caché",
				"id", id,
				"error", err)
			continue
		}
		cafs = append(cafs, caf)
	}

	return cafs, nil
}

// EliminarCAF elimina un CAF del caché
func (c *CAFCache) EliminarCAF(ctx context.Context, id string) error {
	// Primero obtener el CAF para conocer su tipo
	caf, err := c.ObtenerCAF(ctx, id)
	if err != nil {
		return err
	}

	// Eliminar del conjunto por tipo
	tipoKey := PrefijoTipo + caf.TipoDocumento
	if err := c.client.SRem(ctx, tipoKey, id).Err(); err != nil {
		c.logger.Warn("Error eliminando CAF del conjunto por tipo",
			"id", id,
			"tipo", caf.TipoDocumento,
			"error", err)
	}

	// Eliminar CAF individual
	key := PrefijoCaf + id
	if err := c.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("error eliminando CAF del caché: %w", err)
	}

	c.logger.Debug("CAF eliminado del caché",
		"id", id,
		"tipo", caf.TipoDocumento)

	return nil
}

// LimpiarCache limpia todo el caché de CAF
func (c *CAFCache) LimpiarCache(ctx context.Context) error {
	// Obtener todas las claves con los prefijos
	var keys []string

	// Buscar claves de CAF individual
	iter := c.client.Scan(ctx, 0, PrefijoCaf+"*", 0).Iterator()
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}
	if err := iter.Err(); err != nil {
		return fmt.Errorf("error buscando claves de CAF: %w", err)
	}

	// Buscar claves de conjuntos por tipo
	iter = c.client.Scan(ctx, 0, PrefijoTipo+"*", 0).Iterator()
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}
	if err := iter.Err(); err != nil {
		return fmt.Errorf("error buscando claves de tipo: %w", err)
	}

	// Eliminar todas las claves
	if len(keys) > 0 {
		if err := c.client.Del(ctx, keys...).Err(); err != nil {
			return fmt.Errorf("error eliminando claves: %w", err)
		}
	}

	c.logger.Info("Caché de CAF limpiado")

	return nil
}
