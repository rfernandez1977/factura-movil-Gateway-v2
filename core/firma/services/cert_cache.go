package services

import (
	"crypto/x509"
	"sync"
	"time"
)

// CertCache implementa un caché de certificados con expiración
type CertCache struct {
	cache    map[string]*certEntry
	mu       sync.RWMutex
	ttl      time.Duration
	maxItems int
}

type certEntry struct {
	cert      *x509.Certificate
	timestamp time.Time
}

// NewCertCache crea una nueva instancia del caché de certificados
func NewCertCache(ttl time.Duration, maxItems int) *CertCache {
	cache := &CertCache{
		cache:    make(map[string]*certEntry),
		ttl:      ttl,
		maxItems: maxItems,
	}

	// Iniciar rutina de limpieza
	go cache.cleanupRoutine()

	return cache
}

// Get obtiene un certificado del caché
func (c *CertCache) Get(key string) *x509.Certificate {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if entry, exists := c.cache[key]; exists {
		if time.Since(entry.timestamp) < c.ttl {
			return entry.cert
		}
		// Certificado expirado, eliminar
		delete(c.cache, key)
	}
	return nil
}

// Set almacena un certificado en el caché
func (c *CertCache) Set(key string, cert *x509.Certificate) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Si el caché está lleno, eliminar la entrada más antigua
	if len(c.cache) >= c.maxItems {
		var oldestKey string
		var oldestTime time.Time
		for k, v := range c.cache {
			if oldestKey == "" || v.timestamp.Before(oldestTime) {
				oldestKey = k
				oldestTime = v.timestamp
			}
		}
		delete(c.cache, oldestKey)
	}

	c.cache[key] = &certEntry{
		cert:      cert,
		timestamp: time.Now(),
	}
}

// Delete elimina un certificado del caché
func (c *CertCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.cache, key)
}

// Clear limpia todo el caché
func (c *CertCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache = make(map[string]*certEntry)
}

// cleanupRoutine ejecuta la limpieza periódica del caché
func (c *CertCache) cleanupRoutine() {
	ticker := time.NewTicker(c.ttl / 2)
	for range ticker.C {
		c.cleanup()
	}
}

// cleanup elimina las entradas expiradas del caché
func (c *CertCache) cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for key, entry := range c.cache {
		if now.Sub(entry.timestamp) > c.ttl {
			delete(c.cache, key)
		}
	}
}
