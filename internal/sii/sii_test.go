package sii

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func getProjectRoot() string {
	// Asumimos que estamos en internal/sii, así que subimos dos niveles
	pwd, _ := os.Getwd()
	return filepath.Join(filepath.Dir(filepath.Dir(pwd)))
}

const (
	// Rutas de certificados
	CertPath = "firma_test/1/firmaFM.pfx"
	KeyPath  = "firma_test/1/firma.key"

	// RUTs de prueba
	RutEmpresa  = "76212889-6"
	RutEnviador = "13195458-1"
)

func TestSIIAuthentication(t *testing.T) {
	tests := []struct {
		name     string
		certPath string
		wantErr  bool
		errType  error
	}{
		{
			name:     "valid_cert",
			certPath: "../../firma_test/mvp_firma/firmaFM.pfx",
			wantErr:  false,
		},
		{
			name:     "invalid_cert",
			certPath: "nonexistent.pfx",
			wantErr:  true,
			errType:  ErrCertificateNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.certPath)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {
					assert.ErrorIs(t, err, tt.errType)
				}
				return
			}

			assert.NoError(t, err)
			err = client.Authenticate()
			assert.NoError(t, err)
		})
	}
}

func TestSIIEnvio(t *testing.T) {
	client, err := NewClient("../../firma_test/mvp_firma/firmaFM.pfx")
	assert.NoError(t, err)

	err = client.Authenticate()
	assert.NoError(t, err)

	tests := []struct {
		name     string
		dteXML   string
		wantErr  bool
		errType  error
		wantCode string
	}{
		{
			name: "envio_valido",
			dteXML: `<?xml version="1.0" encoding="UTF-8"?>
<DTE version="1.0" xmlns="http://www.sii.cl/SiiDte">
    <Documento ID="DTE_76212889-6_33_1">
        <Encabezado>
            <IdDoc>
                <TipoDTE>33</TipoDTE>
                <Folio>1</Folio>
            </IdDoc>
            <Emisor>
                <RUTEmisor>76212889-6</RUTEmisor>
                <RznSoc>EMPRESA DE PRUEBA SPA</RznSoc>
            </Emisor>
            <Receptor>
                <RUTRecep>60803000-K</RUTRecep>
                <RznSocRecep>SERVICIO DE IMPUESTOS INTERNOS</RznSocRecep>
            </Receptor>
        </Encabezado>
    </Documento>
</DTE>`,
			wantErr:  false,
			wantCode: "123456789",
		},
		{
			name:    "xml_invalido",
			dteXML:  "contenido no xml",
			wantErr: true,
			errType: ErrInvalidXML,
		},
		{
			name:    "sin_autenticacion",
			dteXML:  "<DTE></DTE>",
			wantErr: true,
			errType: ErrNotAuthenticated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "sin_autenticacion" {
				client.token = "" // Simular cliente no autenticado
			}

			trackID, err := client.SendDTE(tt.dteXML)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {
					assert.ErrorIs(t, err, tt.errType)
				}
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.wantCode, trackID)
		})
	}
}

func TestRetryConfig(t *testing.T) {
	client, err := NewClient("../../firma_test/mvp_firma/firmaFM.pfx")
	assert.NoError(t, err)

	// Configurar reintentos personalizados
	customConfig := RetryConfig{
		MaxRetries:    2,
		RetryInterval: 100 * time.Millisecond,
	}
	client.SetRetryConfig(customConfig)

	assert.Equal(t, customConfig.MaxRetries, client.retryConfig.MaxRetries)
	assert.Equal(t, customConfig.RetryInterval, client.retryConfig.RetryInterval)
}

func TestResponseValidation(t *testing.T) {
	client, err := NewClient("../../firma_test/mvp_firma/firmaFM.pfx")
	assert.NoError(t, err)

	tests := []struct {
		name     string
		response []byte
		wantErr  bool
		errType  error
	}{
		{
			name:     "respuesta_valida",
			response: []byte(`<response><status>OK</status><message>Éxito</message></response>`),
			wantErr:  false,
		},
		{
			name:     "respuesta_error",
			response: []byte(`<response><status>ERROR</status><message>Error de validación</message></response>`),
			wantErr:  true,
			errType:  ErrInvalidResponse,
		},
		{
			name:     "respuesta_vacia",
			response: []byte{},
			wantErr:  true,
			errType:  ErrInvalidResponse,
		},
		{
			name:     "xml_invalido",
			response: []byte(`<invalid>`),
			wantErr:  true,
			errType:  ErrInvalidResponse,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.ValidateResponse(tt.response)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {
					assert.ErrorIs(t, err, tt.errType)
				}
				return
			}
			assert.NoError(t, err)
		})
	}
}
