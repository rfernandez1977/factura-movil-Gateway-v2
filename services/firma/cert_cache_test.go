package firma

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTestCertificate(t *testing.T) *x509.Certificate {
	// Generar par de llaves RSA
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	// Crear certificado de prueba
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName:   "Test Certificate",
			Organization: []string{"Test Org"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(24 * time.Hour),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
	}

	// Auto-firmar el certificado
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	require.NoError(t, err)

	cert, err := x509.ParseCertificate(certDER)
	require.NoError(t, err)

	return cert
}

func TestNewCertCache(t *testing.T) {
	ttl := 1 * time.Hour
	maxItems := 10

	cache := NewCertCache(ttl, maxItems)
	assert.NotNil(t, cache)
	assert.Equal(t, ttl, cache.ttl)
	assert.Equal(t, maxItems, cache.maxItems)
	assert.Empty(t, cache.cache)
}

func TestCertCacheOperations(t *testing.T) {
	certPath, password := setupTestCertificates(t)
	config := &ConfiguracionFirma{
		RutaCertificado: certPath,
		Password:        password,
	}

	service, err := NewBaseFirmaService(config)
	require.NoError(t, err)

	cert := service.cert
	cache := NewCertCache(1*time.Hour, 2)

	// Probar Set y Get
	cache.Set("cert1", cert)
	assert.Equal(t, cert, cache.Get("cert1"))

	// Probar caché vacío
	assert.Nil(t, cache.Get("nonexistent"))

	// Probar límite de elementos
	cert2 := &x509.Certificate{SerialNumber: big.NewInt(2)}
	cert3 := &x509.Certificate{SerialNumber: big.NewInt(3)}

	cache.Set("cert2", cert2)
	cache.Set("cert3", cert3)

	// Verificar que se eliminó el elemento más antiguo
	assert.Nil(t, cache.Get("cert1"))
	assert.Equal(t, cert2, cache.Get("cert2"))
	assert.Equal(t, cert3, cache.Get("cert3"))

	// Probar Delete
	cache.Delete("cert2")
	assert.Nil(t, cache.Get("cert2"))

	// Probar Clear
	cache.Clear()
	assert.Nil(t, cache.Get("cert3"))
}

func TestCertCacheExpiration(t *testing.T) {
	certPath, password := setupTestCertificates(t)
	config := &ConfiguracionFirma{
		RutaCertificado: certPath,
		Password:        password,
	}

	service, err := NewBaseFirmaService(config)
	require.NoError(t, err)

	cert := service.cert
	cache := NewCertCache(100*time.Millisecond, 10)

	// Almacenar certificado
	cache.Set("cert1", cert)
	assert.Equal(t, cert, cache.Get("cert1"))

	// Esperar a que expire
	time.Sleep(200 * time.Millisecond)
	assert.Nil(t, cache.Get("cert1"))
}

func TestCertCacheConcurrency(t *testing.T) {
	certPath, password := setupTestCertificates(t)
	config := &ConfiguracionFirma{
		RutaCertificado: certPath,
		Password:        password,
	}

	service, err := NewBaseFirmaService(config)
	require.NoError(t, err)

	cert := service.cert
	cache := NewCertCache(1*time.Hour, 100)

	// Probar operaciones concurrentes
	const numGoroutines = 10
	done := make(chan bool)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			key := fmt.Sprintf("cert%d", id)
			cache.Set(key, cert)
			assert.Equal(t, cert, cache.Get(key))
			cache.Delete(key)
			assert.Nil(t, cache.Get(key))
			done <- true
		}(i)
	}

	// Esperar a que todas las goroutines terminen
	for i := 0; i < numGoroutines; i++ {
		<-done
	}
}

func TestCertCacheCleanup(t *testing.T) {
	cache := NewCertCache(100*time.Millisecond, 3)

	// Llenar el caché
	for i := 0; i < 5; i++ {
		cert := createTestCertificate(t)
		cache.Set(string(i), cert)
		time.Sleep(50 * time.Millisecond)
	}

	// Verificar que solo se mantienen los últimos 3 elementos
	assert.Nil(t, cache.Get("0"))
	assert.Nil(t, cache.Get("1"))
	assert.NotNil(t, cache.Get("2"))
	assert.NotNil(t, cache.Get("3"))
	assert.NotNil(t, cache.Get("4"))

	// Esperar a que expiren todos
	time.Sleep(200 * time.Millisecond)

	// Verificar que todos han expirado
	for i := 0; i < 5; i++ {
		assert.Nil(t, cache.Get(string(i)))
	}
}

func TestCertCacheEdgeCases(t *testing.T) {
	cache := NewCertCache(1*time.Hour, 2)
	cert := createTestCertificate(t)

	// Prueba con clave vacía
	cache.Set("", cert)
	retrieved := cache.Get("")
	assert.NotNil(t, retrieved)
	assert.Equal(t, cert.SerialNumber, retrieved.SerialNumber)

	// Prueba con certificado nil
	cache.Set("nil", nil)
	retrieved = cache.Get("nil")
	assert.Nil(t, retrieved)

	// Prueba de Delete con clave no existente
	cache.Delete("nonexistent")
	// No debería causar error

	// Prueba de Clear múltiples veces
	cache.Clear()
	cache.Clear()
	// No debería causar error
}
