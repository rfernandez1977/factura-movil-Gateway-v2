package services

import (
	"crypto/x509"
	"sync"
	"time"
)

// CertCache implementa un caché thread-safe para certificados digitales
type CertCache struct {
	cache    map[string]*cachedCert
	mu       sync.RWMutex
	maxAge   time.Duration
	maxItems int
}

type cachedCert struct {
	cert      *x509.Certificate
	timestamp time.Time
}

// NewCertCache crea una nueva instancia de CertCache
func NewCertCache(maxAge time.Duration, maxItems int) *CertCache {
	return &CertCache{
		cache:    make(map[string]*cachedCert),
		maxAge:   maxAge,
		maxItems: maxItems,
	}
}

// Get obtiene un certificado del caché
func (c *CertCache) Get(key string) (*x509.Certificate, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if item, exists := c.cache[key]; exists {
		if time.Since(item.timestamp) < c.maxAge {
			return item.cert, true
		}
		// Certificado expirado, eliminar del caché
		delete(c.cache, key)
	}
	return nil, false
}

// Set almacena un certificado en el caché
func (c *CertCache) Set(key string, cert *x509.Certificate) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Limpiar caché si está lleno
	if len(c.cache) >= c.maxItems {
		c.cleanup()
	}

	c.cache[key] = &cachedCert{
		cert:      cert,
		timestamp: time.Now(),
	}
}

// cleanup elimina los certificados más antiguos cuando el caché está lleno
func (c *CertCache) cleanup() {
	// Encontrar el certificado más antiguo
	var oldestKey string
	oldestTime := time.Now()

	for key, item := range c.cache {
		if item.timestamp.Before(oldestTime) {
			oldestTime = item.timestamp
			oldestKey = key
		}
	}

	// Eliminar el certificado más antiguo
	if oldestKey != "" {
		delete(c.cache, oldestKey)
	}
}

// Clear limpia todo el caché
func (c *CertCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache = make(map[string]*cachedCert)
}
