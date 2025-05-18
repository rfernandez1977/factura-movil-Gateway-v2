package firma

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"software.sslmate.com/src/go-pkcs12"
)

func setupTestCertificates(t *testing.T) (string, string) {
	// Crear directorio temporal
	tmpDir, err := os.MkdirTemp("", "firma_test")
	require.NoError(t, err)
	t.Cleanup(func() { os.RemoveAll(tmpDir) })

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

	// Crear archivo PKCS12
	pfxData, err := pkcs12.Encode(rand.Reader, privateKey, cert, nil, "test123")
	require.NoError(t, err)

	pfxPath := filepath.Join(tmpDir, "test.p12")
	err = os.WriteFile(pfxPath, pfxData, 0600)
	require.NoError(t, err)

	return pfxPath, "test123"
}

func TestNewBaseFirmaService(t *testing.T) {
	certPath, password := setupTestCertificates(t)

	tests := []struct {
		name        string
		config      *ConfiguracionFirma
		expectError bool
	}{
		{
			name: "Configuración válida",
			config: &ConfiguracionFirma{
				RutaCertificado: certPath,
				Password:        password,
			},
			expectError: false,
		},
		{
			name: "Ruta de certificado inválida",
			config: &ConfiguracionFirma{
				RutaCertificado: "invalid/path.p12",
				Password:        password,
			},
			expectError: true,
		},
		{
			name: "Contraseña incorrecta",
			config: &ConfiguracionFirma{
				RutaCertificado: certPath,
				Password:        "wrong_password",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, err := NewBaseFirmaService(tt.config)
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, service)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, service)
				assert.NotNil(t, service.privateKey)
				assert.NotNil(t, service.cert)
			}
		})
	}
}

func TestFirmarDocumento(t *testing.T) {
	certPath, password := setupTestCertificates(t)
	config := &ConfiguracionFirma{
		RutaCertificado: certPath,
		Password:        password,
	}

	service, err := NewBaseFirmaService(config)
	require.NoError(t, err)

	tests := []struct {
		name        string
		documento   string
		expectError bool
	}{
		{
			name: "Documento XML válido",
			documento: `<?xml version="1.0" encoding="UTF-8"?>
<Documento>
    <Contenido>Test</Contenido>
</Documento>`,
			expectError: false,
		},
		{
			name:        "Documento vacío",
			documento:   "",
			expectError: true,
		},
		{
			name:        "XML inválido",
			documento:   "<InvalidXML>",
			expectError: true,
		},
		{
			name: "XML con caracteres especiales",
			documento: `<?xml version="1.0" encoding="UTF-8"?>
<Documento>
    <Contenido>Test & áéíóú</Contenido>
</Documento>`,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultado, err := service.FirmarDocumento(tt.documento)
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, resultado)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resultado)
				assert.NotEmpty(t, resultado.XMLFirmado)
				assert.NotEmpty(t, resultado.DigestValue)
				assert.NotEmpty(t, resultado.SignatureValue)
				assert.NotEmpty(t, resultado.Timestamp)

				// Validar la firma
				estado, err := service.ValidarFirma(resultado.XMLFirmado)
				assert.NoError(t, err)
				assert.True(t, estado.Valida)
			}
		})
	}
}

func TestValidarFirma(t *testing.T) {
	certPath, password := setupTestCertificates(t)
	config := &ConfiguracionFirma{
		RutaCertificado: certPath,
		Password:        password,
	}

	service, err := NewBaseFirmaService(config)
	require.NoError(t, err)

	// Firmar un documento de prueba
	docOriginal := `<?xml version="1.0" encoding="UTF-8"?>
<Documento>
    <Contenido>Test</Contenido>
</Documento>`

	resultado, err := service.FirmarDocumento(docOriginal)
	require.NoError(t, err)

	tests := []struct {
		name        string
		xmlFirmado  string
		expectError bool
		expectValid bool
	}{
		{
			name:        "Firma válida",
			xmlFirmado:  resultado.XMLFirmado,
			expectError: false,
			expectValid: true,
		},
		{
			name:        "XML sin firma",
			xmlFirmado:  docOriginal,
			expectError: false,
			expectValid: false,
		},
		{
			name:        "XML inválido",
			xmlFirmado:  "<InvalidXML>",
			expectError: true,
			expectValid: false,
		},
		{
			name:        "XML vacío",
			xmlFirmado:  "",
			expectError: true,
			expectValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			estado, err := service.ValidarFirma(tt.xmlFirmado)
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, estado)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, estado)
				assert.Equal(t, tt.expectValid, estado.Valida)
				if tt.expectValid {
					assert.NotEmpty(t, estado.FechaValidacion)
					assert.NotEmpty(t, estado.CertificadoID)
					assert.Empty(t, estado.Error)
				}
			}
		})
	}
}

func TestObtenerCertificado(t *testing.T) {
	certPath, password := setupTestCertificates(t)
	config := &ConfiguracionFirma{
		RutaCertificado: certPath,
		Password:        password,
	}

	service, err := NewBaseFirmaService(config)
	require.NoError(t, err)

	// Probar obtención de certificado
	cert, err := service.ObtenerCertificado()
	assert.NoError(t, err)
	assert.NotNil(t, cert)
	assert.Equal(t, service.cert, cert)

	// Probar con servicio sin certificado
	invalidService := &BaseFirmaService{}
	cert, err = invalidService.ObtenerCertificado()
	assert.Error(t, err)
	assert.Nil(t, cert)
}

func TestConcurrencia(t *testing.T) {
	certPath, password := setupTestCertificates(t)
	config := &ConfiguracionFirma{
		RutaCertificado: certPath,
		Password:        password,
	}

	service, err := NewBaseFirmaService(config)
	require.NoError(t, err)

	docOriginal := `<?xml version="1.0" encoding="UTF-8"?>
<Documento>
    <Contenido>Test</Contenido>
</Documento>`

	// Probar firmas concurrentes
	numGoroutines := 10
	resultados := make(chan error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			_, err := service.FirmarDocumento(docOriginal)
			resultados <- err
		}()
	}

	// Verificar resultados
	for i := 0; i < numGoroutines; i++ {
		err := <-resultados
		assert.NoError(t, err)
	}
}
