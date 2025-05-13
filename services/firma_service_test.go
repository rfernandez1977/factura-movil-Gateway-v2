package services

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFirmaService(t *testing.T) {
	// Crear certificado y llave de prueba
	certPath := "testdata/test.crt"
	keyPath := "testdata/test.key"
	password := "test123"
	rutFirmante := "11.111.111-1"

	// Crear directorio de test si no existe
	if err := os.MkdirAll("testdata", 0755); err != nil {
		t.Fatalf("Error al crear directorio de test: %v", err)
	}

	// Generar certificado y llave de prueba
	if err := generarCertificadoPrueba(certPath, keyPath, password); err != nil {
		t.Fatalf("Error al generar certificado de prueba: %v", err)
	}
	defer os.RemoveAll("testdata")

	// Test casos exitosos
	t.Run("Crear servicio con certificado válido", func(t *testing.T) {
		service, err := NewFirmaService(certPath, keyPath, password, rutFirmante)
		assert.NoError(t, err)
		assert.NotNil(t, service)
		assert.Equal(t, rutFirmante, service.rutFirmante)
		assert.Equal(t, "SHA1", service.algoritmo)
	})

	t.Run("Crear servicio sin contraseña", func(t *testing.T) {
		// Generar certificado sin contraseña
		certSinPass := "testdata/test_nopass.crt"
		keySinPass := "testdata/test_nopass.key"
		if err := generarCertificadoSinPassword(certSinPass, keySinPass); err != nil {
			t.Fatalf("Error al generar certificado sin contraseña: %v", err)
		}

		service, err := NewFirmaService(certSinPass, keySinPass, "", rutFirmante)
		assert.NoError(t, err)
		assert.NotNil(t, service)
	})

	// Test casos de error
	t.Run("Error con certificado inválido", func(t *testing.T) {
		_, err := NewFirmaService("noexiste.crt", keyPath, password, rutFirmante)
		assert.Error(t, err)
	})

	t.Run("Error con llave inválida", func(t *testing.T) {
		_, err := NewFirmaService(certPath, "noexiste.key", password, rutFirmante)
		assert.Error(t, err)
	})

	t.Run("Error con contraseña incorrecta", func(t *testing.T) {
		_, err := NewFirmaService(certPath, keyPath, "wrongpass", rutFirmante)
		assert.Error(t, err)
	})
}

func TestFirmarXML(t *testing.T) {
	// Setup
	service := setupTestService(t)
	defer cleanupTestFiles(t)

	// Test casos exitosos
	t.Run("Firmar XML simple", func(t *testing.T) {
		xmlData := []byte(`<Documento ID="TEST001">
			<Datos>Test</Datos>
		</Documento>`)

		firmado, err := service.FirmarXML(xmlData, "TEST001")
		assert.NoError(t, err)
		assert.NotNil(t, firmado)
		assert.Contains(t, string(firmado), "<Signature")
	})

	t.Run("Firmar XML con SHA256", func(t *testing.T) {
		err := service.SetAlgoritmo("SHA256")
		assert.NoError(t, err)

		xmlData := []byte(`<Documento ID="TEST002">
			<Datos>Test SHA256</Datos>
		</Documento>`)

		firmado, err := service.FirmarXML(xmlData, "TEST002")
		assert.NoError(t, err)
		assert.NotNil(t, firmado)
		assert.Contains(t, string(firmado), "rsa-sha256")
	})

	// Test casos de error
	t.Run("Error con XML inválido", func(t *testing.T) {
		xmlData := []byte(`<Documento ID="TEST003">
			<Datos>Test</Datos>
		</Documento`)

		_, err := service.FirmarXML(xmlData, "TEST003")
		assert.Error(t, err)
	})

	t.Run("Error con ID no encontrado", func(t *testing.T) {
		xmlData := []byte(`<Documento>
			<Datos>Test</Datos>
		</Documento>`)

		_, err := service.FirmarXML(xmlData, "NOEXISTE")
		assert.Error(t, err)
	})
}

func TestVerificarFirma(t *testing.T) {
	// Setup
	service := setupTestService(t)
	defer cleanupTestFiles(t)

	// Test casos exitosos
	t.Run("Verificar firma válida", func(t *testing.T) {
		xmlData := []byte(`<Documento ID="TEST001">
			<Datos>Test</Datos>
		</Documento>`)

		firmado, err := service.FirmarXML(xmlData, "TEST001")
		assert.NoError(t, err)

		valido, err := service.VerificarFirma(firmado)
		assert.NoError(t, err)
		assert.True(t, valido)
	})

	t.Run("Verificar firma SHA256", func(t *testing.T) {
		err := service.SetAlgoritmo("SHA256")
		assert.NoError(t, err)

		xmlData := []byte(`<Documento ID="TEST002">
			<Datos>Test SHA256</Datos>
		</Documento>`)

		firmado, err := service.FirmarXML(xmlData, "TEST002")
		assert.NoError(t, err)

		valido, err := service.VerificarFirma(firmado)
		assert.NoError(t, err)
		assert.True(t, valido)
	})

	// Test casos de error
	t.Run("Error con XML sin firma", func(t *testing.T) {
		xmlData := []byte(`<Documento ID="TEST003">
			<Datos>Test</Datos>
		</Documento>`)

		valido, err := service.VerificarFirma(xmlData)
		assert.Error(t, err)
		assert.False(t, valido)
	})

	t.Run("Error con XML inválido", func(t *testing.T) {
		xmlData := []byte(`<Documento ID="TEST004">
			<Datos>Test</Datos>
		</Documento`)

		valido, err := service.VerificarFirma(xmlData)
		assert.Error(t, err)
		assert.False(t, valido)
	})
}

func TestValidarCertificado(t *testing.T) {
	// Setup
	service := setupTestService(t)
	defer cleanupTestFiles(t)

	// Test casos exitosos
	t.Run("Validar certificado válido", func(t *testing.T) {
		err := service.ValidarCertificado()
		assert.NoError(t, err)
	})

	// Test casos de error
	t.Run("Error con certificado expirado", func(t *testing.T) {
		// Crear certificado expirado
		certExpirado := "testdata/expired.crt"
		keyExpirado := "testdata/expired.key"
		if err := generarCertificadoExpirado(certExpirado, keyExpirado); err != nil {
			t.Fatalf("Error al generar certificado expirado: %v", err)
		}

		service, err := NewFirmaService(certExpirado, keyExpirado, "test123", "11.111.111-1")
		assert.NoError(t, err)

		err = service.ValidarCertificado()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "certificado fuera de fecha de validez")
	})
}

// Funciones auxiliares para los tests

func setupTestService(t *testing.T) *FirmaService {
	// Crear certificado y llave de prueba
	certPath := "testdata/test.crt"
	keyPath := "testdata/test.key"
	password := "test123"
	rutFirmante := "11.111.111-1"

	// Crear directorio de test si no existe
	if err := os.MkdirAll("testdata", 0755); err != nil {
		t.Fatalf("Error al crear directorio de test: %v", err)
	}

	// Generar certificado y llave de prueba
	if err := generarCertificadoPrueba(certPath, keyPath, password); err != nil {
		t.Fatalf("Error al generar certificado de prueba: %v", err)
	}

	// Crear servicio
	service, err := NewFirmaService(certPath, keyPath, password, rutFirmante)
	if err != nil {
		t.Fatalf("Error al crear servicio de prueba: %v", err)
	}

	return service
}

func cleanupTestFiles(t *testing.T) {
	if err := os.RemoveAll("testdata"); err != nil {
		t.Fatalf("Error al limpiar archivos de test: %v", err)
	}
}

func generarCertificadoPrueba(certPath, keyPath, password string) error {
	// TODO: Implementar generación de certificado de prueba
	// Por ahora, crear archivos vacíos para las pruebas
	if err := os.WriteFile(certPath, []byte("test certificate"), 0644); err != nil {
		return err
	}
	if err := os.WriteFile(keyPath, []byte("test key"), 0600); err != nil {
		return err
	}
	return nil
}

func generarCertificadoSinPassword(certPath, keyPath string) error {
	// TODO: Implementar generación de certificado sin contraseña
	// Por ahora, crear archivos vacíos para las pruebas
	if err := os.WriteFile(certPath, []byte("test certificate no pass"), 0644); err != nil {
		return err
	}
	if err := os.WriteFile(keyPath, []byte("test key no pass"), 0600); err != nil {
		return err
	}
	return nil
}

func generarCertificadoExpirado(certPath, keyPath string) error {
	// TODO: Implementar generación de certificado expirado
	// Por ahora, crear archivos vacíos para las pruebas
	if err := os.WriteFile(certPath, []byte("expired certificate"), 0644); err != nil {
		return err
	}
	if err := os.WriteFile(keyPath, []byte("expired key"), 0600); err != nil {
		return err
	}
	return nil
}
