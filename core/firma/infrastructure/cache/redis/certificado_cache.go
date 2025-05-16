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
	// TTLCertificado tiempo de vida del certificado en caché
	TTLCertificado = 24 * time.Hour
	// PrefijoCertificado prefijo para las claves de certificados
	PrefijoCertificado = "cert:"
)

// CertificadoCache implementa el caché de certificados usando Redis
type CertificadoCache struct {
	client *redis.Client
	logger common.Logger
}

// NewCertificadoCache crea una nueva instancia del caché de certificados
func NewCertificadoCache(client *redis.Client, logger common.Logger) *CertificadoCache {
	return &CertificadoCache{
		client: client,
		logger: logger,
	}
}

// ObtenerCertificado obtiene un certificado del caché
func (c *CertificadoCache) ObtenerCertificado(ctx context.Context, id string) (*models.Certificado, error) {
	key := PrefijoCertificado + id

	// Obtener del caché
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("certificado no encontrado en caché: %s", id)
		}
		return nil, fmt.Errorf("error obteniendo certificado del caché: %w", err)
	}

	// Deserializar
	var cert models.Certificado
	if err := json.Unmarshal(data, &cert); err != nil {
		return nil, fmt.Errorf("error deserializando certificado: %w", err)
	}

	c.logger.Debug("Certificado obtenido del caché",
		"id", id,
		"rut", cert.RutEmpresa)

	return &cert, nil
}

// GuardarCertificado guarda un certificado en el caché
func (c *CertificadoCache) GuardarCertificado(ctx context.Context, id string, cert *models.Certificado) error {
	key := PrefijoCertificado + id

	// Serializar
	data, err := json.Marshal(cert)
	if err != nil {
		return fmt.Errorf("error serializando certificado: %w", err)
	}

	// Guardar en caché con TTL
	if err := c.client.Set(ctx, key, data, TTLCertificado).Err(); err != nil {
		return fmt.Errorf("error guardando certificado en caché: %w", err)
	}

	c.logger.Debug("Certificado guardado en caché",
		"id", id,
		"rut", cert.RutEmpresa,
		"ttl", TTLCertificado)

	return nil
}

// EliminarCertificado elimina un certificado del caché
func (c *CertificadoCache) EliminarCertificado(ctx context.Context, id string) error {
	key := PrefijoCertificado + id

	// Eliminar del caché
	if err := c.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("error eliminando certificado del caché: %w", err)
	}

	c.logger.Debug("Certificado eliminado del caché", "id", id)

	return nil
}

// LimpiarCache limpia todo el caché de certificados
func (c *CertificadoCache) LimpiarCache(ctx context.Context) error {
	// Obtener todas las claves con el prefijo
	pattern := PrefijoCertificado + "*"
	iter := c.client.Scan(ctx, 0, pattern, 0).Iterator()

	// Eliminar cada clave
	for iter.Next(ctx) {
		if err := c.client.Del(ctx, iter.Val()).Err(); err != nil {
			return fmt.Errorf("error eliminando clave %s: %w", iter.Val(), err)
		}
	}

	if err := iter.Err(); err != nil {
		return fmt.Errorf("error iterando claves: %w", err)
	}

	c.logger.Info("Caché de certificados limpiado")

	return nil
}
