package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogger_Basic(t *testing.T) {
	// Configurar archivo de log temporal
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "test.log")

	// Crear logger
	log, err := NewLogger(logPath, DEBUG)
	require.NoError(t, err)
	defer log.Close()

	// Probar diferentes niveles de log
	log.Debug("mensaje debug %s", "test")
	log.Info("mensaje info %s", "test")
	log.Warn("mensaje warn %s", "test")
	log.Error("mensaje error %s", "test")

	// Leer contenido del archivo
	content, err := os.ReadFile(logPath)
	require.NoError(t, err)

	// Verificar que todos los mensajes se escribieron
	logContent := string(content)
	assert.Contains(t, logContent, "DEBUG: ")
	assert.Contains(t, logContent, "INFO: ")
	assert.Contains(t, logContent, "WARN: ")
	assert.Contains(t, logContent, "ERROR: ")
	assert.Contains(t, logContent, "mensaje debug test")
	assert.Contains(t, logContent, "mensaje info test")
	assert.Contains(t, logContent, "mensaje warn test")
	assert.Contains(t, logContent, "mensaje error test")
}

func TestLogger_XMLOperation(t *testing.T) {
	// Configurar archivo de log temporal
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "test.log")

	// Crear logger
	log, err := NewLogger(logPath, DEBUG)
	require.NoError(t, err)
	defer log.Close()

	// XML de prueba
	xmlData := []byte(`<?xml version="1.0"?><test>contenido</test>`)

	// Probar log de operación exitosa
	log.LogXMLOperation("test", xmlData, nil)

	// Probar log de operación con error
	testError := fmt.Errorf("error de prueba")
	log.LogXMLOperation("test", xmlData, testError)

	// Leer contenido del archivo
	content, err := os.ReadFile(logPath)
	require.NoError(t, err)

	// Verificar contenido
	logContent := string(content)
	assert.Contains(t, logContent, "Operación XML 'test' exitosa")
	assert.Contains(t, logContent, "Operación XML 'test' falló")
	assert.Contains(t, logContent, "error de prueba")
	assert.Contains(t, logContent, "<test>contenido</test>")
}

func TestLogger_CertOperation(t *testing.T) {
	// Configurar archivo de log temporal
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "test.log")

	// Crear logger
	log, err := NewLogger(logPath, DEBUG)
	require.NoError(t, err)
	defer log.Close()

	// Probar log de operación exitosa
	log.LogCertOperation("test", "CN=Test Certificate", nil)

	// Probar log de operación con error
	testError := fmt.Errorf("error de certificado")
	log.LogCertOperation("test", "CN=Test Certificate", testError)

	// Leer contenido del archivo
	content, err := os.ReadFile(logPath)
	require.NoError(t, err)

	// Verificar contenido
	logContent := string(content)
	assert.Contains(t, logContent, "Operación de certificado 'test' exitosa")
	assert.Contains(t, logContent, "Operación de certificado 'test' falló")
	assert.Contains(t, logContent, "error de certificado")
	assert.Contains(t, logContent, "CN=Test Certificate")
}

func TestLogger_LevelFiltering(t *testing.T) {
	// Configurar archivo de log temporal
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "test.log")

	// Crear logger con nivel INFO
	log, err := NewLogger(logPath, INFO)
	require.NoError(t, err)
	defer log.Close()

	// Enviar mensajes de todos los niveles
	log.Debug("mensaje debug")
	log.Info("mensaje info")
	log.Warn("mensaje warn")
	log.Error("mensaje error")

	// Leer contenido del archivo
	content, err := os.ReadFile(logPath)
	require.NoError(t, err)

	// Verificar que solo se escribieron los mensajes de nivel INFO y superior
	logContent := string(content)
	assert.NotContains(t, logContent, "mensaje debug")
	assert.Contains(t, logContent, "mensaje info")
	assert.Contains(t, logContent, "mensaje warn")
	assert.Contains(t, logContent, "mensaje error")
}
