package integration

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"software.sslmate.com/src/go-pkcs12"

	"FMgo/core/sii/services"
)

func TestFirmaIntegration(t *testing.T) {
	// Configuración de prueba
	certPath := filepath.Join("testdata", "cert_test.p12")
	password := "test123"
	rutEmpresa := "76.555.555-5"

	// Crear certificado de prueba
	err := setupTestCertificate(certPath, password)
	require.NoError(t, err)
	defer os.Remove(certPath)

	// Crear servicio de firma
	firmaService, err := services.NewFirmaService(certPath, "", password, rutEmpresa)
	require.NoError(t, err)

	// Cargar XML de prueba
	xmlData, err := ioutil.ReadFile(filepath.Join("testdata", "documento_test.xml"))
	require.NoError(t, err)

	// Probar firma de documento
	signedXML, err := firmaService.FirmarXML(xmlData)
	assert.NoError(t, err)
	assert.NotNil(t, signedXML)

	// Validar firma
	err = firmaService.ValidarFirma(signedXML)
	assert.NoError(t, err)
}

func setupTestCertificate(certPath, password string) error {
	// Generar clave privada RSA
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("error al generar clave privada: %w", err)
	}

	// Crear plantilla de certificado
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName:   "Test Certificate",
			Organization: []string{"Test Company"},
			Country:      []string{"CL"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(1, 0, 0), // Válido por 1 año
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	// Crear certificado autofirmado
	certBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return fmt.Errorf("error al crear certificado: %w", err)
	}

	// Convertir a PKCS12
	pfxData, err := pkcs12.Encode(rand.Reader, privateKey, template, nil, password)
	if err != nil {
		return fmt.Errorf("error al codificar PKCS12: %w", err)
	}

	// Crear directorio si no existe
	err = os.MkdirAll(filepath.Dir(certPath), 0755)
	if err != nil {
		return fmt.Errorf("error al crear directorio: %w", err)
	}

	// Guardar archivo P12
	err = ioutil.WriteFile(certPath, pfxData, 0644)
	if err != nil {
		return fmt.Errorf("error al guardar certificado: %w", err)
	}

	return nil
}
