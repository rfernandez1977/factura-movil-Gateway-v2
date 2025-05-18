package firma

import (
	"crypto/x509"
	"sync"
	"time"
)

// CertCacheItem representa un elemento en el caché
type CertCacheItem struct {
	cert      *x509.Certificate
	timestamp time.Time
}

// CertCache implementa un caché de certificados con TTL
type CertCache struct {
	cache    map[string]CertCacheItem
	mu       sync.RWMutex
	ttl      time.Duration
	maxItems int
}

// NewCertCache crea una nueva instancia de CertCache
func NewCertCache(ttl time.Duration, maxItems int) *CertCache {
	return &CertCache{
		cache:    make(map[string]CertCacheItem),
		ttl:      ttl,
		maxItems: maxItems,
	}
}

// Get obtiene un certificado del caché
func (c *CertCache) Get(key string) *x509.Certificate {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if item, exists := c.cache[key]; exists {
		if time.Since(item.timestamp) < c.ttl {
			return item.cert
		}
		// Si expiró, lo eliminamos
		delete(c.cache, key)
	}
	return nil
}

// Set almacena un certificado en el caché
func (c *CertCache) Set(key string, cert *x509.Certificate) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Si alcanzamos el límite, eliminamos el elemento más antiguo
	if len(c.cache) >= c.maxItems {
		var oldestKey string
		var oldestTime time.Time
		first := true

		for k, v := range c.cache {
			if first || v.timestamp.Before(oldestTime) {
				oldestKey = k
				oldestTime = v.timestamp
				first = false
			}
		}
		delete(c.cache, oldestKey)
	}

	c.cache[key] = CertCacheItem{
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

// Clear elimina todos los certificados del caché
func (c *CertCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache = make(map[string]CertCacheItem)
}
