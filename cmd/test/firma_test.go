package test

import (
	"encoding/json"
	"os"
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
	// Cargar configuración
	configData, err := os.ReadFile("../../dev/config/certs/firma/config.json")
	assert.NoError(t, err, "Error al leer archivo de configuración")

	var config CertConfig
	err = json.Unmarshal(configData, &config)
	assert.NoError(t, err, "Error al parsear configuración")

	// Verificar existencia de archivos
	_, err = os.Stat(config.Certificado.RutaKey)
	assert.NoError(t, err, "Archivo .key no encontrado")

	_, err = os.Stat(config.Certificado.RutaPfx)
	assert.NoError(t, err, "Archivo .pfx no encontrado")

	// Verificar formato del RUT
	assert.Regexp(t, `^\d{1,8}-[\dkK]$`, config.Certificado.RutFirmante, "Formato de RUT inválido")
}
