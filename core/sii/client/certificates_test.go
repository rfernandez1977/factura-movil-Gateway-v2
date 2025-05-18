package client

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func generateTestCertificate(t *testing.T, notBefore, notAfter time.Time) (*x509.Certificate, string) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.NoError(t, err)

	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName:   "Test Company SpA",
			Organization: []string{"RUT:76123456-7"},
		},
		EmailAddresses: []string{"test@example.com"},
		NotBefore:      notBefore,
		NotAfter:       notAfter,
		KeyUsage:       x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:    []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		IsCA:           false,
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, template, template, &privateKey.PublicKey, privateKey)
	assert.NoError(t, err)

	// Crear archivo temporal
	tmpfile, err := os.CreateTemp("", "test-cert-*.pem")
	assert.NoError(t, err)
	defer tmpfile.Close()

	// Escribir certificado en formato PEM
	err = pem.Encode(tmpfile, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})
	assert.NoError(t, err)

	return template, tmpfile.Name()
}

func TestDefaultCertManager_CargarCertificado(t *testing.T) {
	manager := NewCertificateManager()

	// Generar certificado de prueba válido
	now := time.Now()
	_, certPath := generateTestCertificate(t, now.Add(-1*time.Hour), now.Add(24*time.Hour))
	defer os.Remove(certPath)

	// Probar carga exitosa
	cert, err := manager.CargarCertificado(certPath, "")
	assert.NoError(t, err)
	assert.NotNil(t, cert)

	// Probar archivo no existente
	_, err = manager.CargarCertificado("no-existe.pem", "")
	assert.Error(t, err)
}

func TestDefaultCertManager_ValidarCertificado(t *testing.T) {
	manager := NewCertificateManager()
	now := time.Now()

	tests := []struct {
		name      string
		notBefore time.Time
		notAfter  time.Time
		wantErr   bool
	}{
		{
			name:      "Certificado válido",
			notBefore: now.Add(-1 * time.Hour),
			notAfter:  now.Add(24 * time.Hour),
			wantErr:   false,
		},
		{
			name:      "Certificado no vigente aún",
			notBefore: now.Add(1 * time.Hour),
			notAfter:  now.Add(24 * time.Hour),
			wantErr:   true,
		},
		{
			name:      "Certificado expirado",
			notBefore: now.Add(-24 * time.Hour),
			notAfter:  now.Add(-1 * time.Hour),
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cert, _ := generateTestCertificate(t, tt.notBefore, tt.notAfter)
			err := manager.ValidarCertificado(cert)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDefaultCertManager_ObtenerDatosFirmante(t *testing.T) {
	manager := NewCertificateManager()
	now := time.Now()

	// Generar certificado de prueba
	cert, _ := generateTestCertificate(t, now.Add(-1*time.Hour), now.Add(24*time.Hour))

	// Probar extracción de datos
	datos, err := manager.ObtenerDatosFirmante(cert)
	assert.NoError(t, err)
	assert.Equal(t, "76123456-7", datos.RUT)
	assert.Equal(t, "Test Company SpA", datos.RazonSocial)
	assert.Equal(t, "test@example.com", datos.Email)

	// Probar certificado nil
	_, err = manager.ObtenerDatosFirmante(nil)
	assert.Error(t, err)
}

func TestDefaultCertManager_RenovarCertificado(t *testing.T) {
	manager := NewCertificateManager()
	now := time.Now()

	// Generar certificado de prueba
	cert, _ := generateTestCertificate(t, now.Add(-1*time.Hour), now.Add(24*time.Hour))

	// Por ahora, la renovación no está implementada y debe retornar error
	err := manager.RenovarCertificado(cert)
	assert.Error(t, err)
}
