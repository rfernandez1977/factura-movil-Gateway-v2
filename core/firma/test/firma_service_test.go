package test

import (
	"encoding/xml"
	"path/filepath"
	"strings"
	"testing"

	"FMgo/core/firma/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFirmaService(t *testing.T) {
	// Obtener ruta del archivo de configuración
	configPath := filepath.Join("..", "..", "..", "dev", "config", "certs", "firma", "config.json")

	// Crear servicio de firma
	firmaService, err := services.NewFirmaService(configPath)
	require.NoError(t, err, "Error al crear servicio de firma")
	require.NotNil(t, firmaService, "Servicio de firma es nil")

	t.Run("Firma de Semilla", func(t *testing.T) {
		// Preparar semilla de prueba
		semilla := "1234567890"

		// Firmar semilla
		resultado, err := firmaService.FirmarSemilla(semilla)
		require.NoError(t, err, "Error al firmar semilla")
		require.NotNil(t, resultado, "Resultado de firma es nil")

		// Verificar estructura del XML firmado
		assert.True(t, strings.Contains(resultado.XMLFirmado, "<ds:Signature"), "XML no contiene firma")
		assert.True(t, strings.Contains(resultado.XMLFirmado, semilla), "XML no contiene semilla")
		assert.NotEmpty(t, resultado.DigestValue, "DigestValue está vacío")
		assert.NotEmpty(t, resultado.SignatureValue, "SignatureValue está vacío")

		// Validar que el XML sea válido
		var xmlDoc struct{}
		err = xml.Unmarshal([]byte(resultado.XMLFirmado), &xmlDoc)
		assert.NoError(t, err, "XML firmado no es válido")
	})

	t.Run("Firma de Documento", func(t *testing.T) {
		// Preparar documento de prueba
		documento := `<?xml version="1.0" encoding="UTF-8"?>
<Documento>
    <Contenido>Test de firma digital</Contenido>
</Documento>`

		// Firmar documento
		resultado, err := firmaService.FirmarDocumento(documento)
		require.NoError(t, err, "Error al firmar documento")
		require.NotNil(t, resultado, "Resultado de firma es nil")

		// Verificar estructura del XML firmado
		assert.True(t, strings.Contains(resultado.XMLFirmado, "<ds:Signature"), "XML no contiene firma")
		assert.True(t, strings.Contains(resultado.XMLFirmado, "Test de firma digital"), "XML no contiene contenido original")
		assert.NotEmpty(t, resultado.DigestValue, "DigestValue está vacío")
		assert.NotEmpty(t, resultado.SignatureValue, "SignatureValue está vacío")

		// Validar que el XML sea válido
		var xmlDoc struct{}
		err = xml.Unmarshal([]byte(resultado.XMLFirmado), &xmlDoc)
		assert.NoError(t, err, "XML firmado no es válido")
	})

	t.Run("Validación de Firma", func(t *testing.T) {
		// Preparar documento firmado
		documento := `<?xml version="1.0" encoding="UTF-8"?>
<Documento>
    <Contenido>Test de validación de firma</Contenido>
</Documento>`

		// Firmar documento
		resultado, err := firmaService.FirmarDocumento(documento)
		require.NoError(t, err, "Error al firmar documento")

		// Validar firma
		estado, err := firmaService.ValidarFirma(resultado.XMLFirmado)
		require.NoError(t, err, "Error al validar firma")
		require.NotNil(t, estado, "Estado de validación es nil")

		assert.True(t, estado.Valida, "La firma no es válida")
		assert.NotEmpty(t, estado.CertificadoID, "ID del certificado está vacío")
		assert.NotZero(t, estado.FechaValidacion, "Fecha de validación no está establecida")
	})
}
