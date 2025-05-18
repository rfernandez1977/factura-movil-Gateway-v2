package test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

type CertConfig struct {
	Certificado struct {
		RutFirmante string `json:"rut_firmante"`
		RutaKey     string `json:"ruta_key"`
		RutaPfx     string `json:"ruta_pfx"`
		Password    string `json:"password"`
	} `json:"certificado"`
}

func TestCertificadoDigital(t *testing.T) {
	// Obtener directorio raíz del proyecto
	pwd, err := os.Getwd()
	assert.NoError(t, err, "Error al obtener directorio actual")

	rootDir := filepath.Join(filepath.Dir(filepath.Dir(filepath.Dir(pwd))))

	// Cargar configuración
	configPath := filepath.Join(rootDir, "dev/config/certs/firma/config.json")
	configData, err := os.ReadFile(configPath)
	assert.NoError(t, err, "Error al leer archivo de configuración")

	var config CertConfig
	err = json.Unmarshal(configData, &config)
	assert.NoError(t, err, "Error al parsear configuración")

	// Verificar existencia de archivos
	keyPath := filepath.Join(rootDir, config.Certificado.RutaKey)
	_, err = os.Stat(keyPath)
	assert.NoError(t, err, "Archivo .key no encontrado")

	pfxPath := filepath.Join(rootDir, config.Certificado.RutaPfx)
	_, err = os.Stat(pfxPath)
	assert.NoError(t, err, "Archivo .pfx no encontrado")

	// Verificar formato del RUT
	assert.Regexp(t, `^\d{1,8}-[\dkK]$`, config.Certificado.RutFirmante, "Formato de RUT inválido")
}
