package firma

import (
	"encoding/xml"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testSemilla = "1234567890"
	testToken   = "ABC123XYZ"
	testCAF     = `<?xml version="1.0" encoding="UTF-8"?>
<CAF version="1.0">
    <DA>
        <RE>76555555-5</RE>
        <RS>EMPRESA DE PRUEBA</RS>
        <TD>33</TD>
        <RNG><D>1</D><H>100</H></RNG>
        <FA>2024-03-20</FA>
        <RSAPK><M>123</M><E>456</E></RSAPK>
        <IDK>1</IDK>
    </DA>
    <FRMA>FIRMA_SII</FRMA>
</CAF>`

	testDTE = `<?xml version="1.0" encoding="UTF-8"?>
<DTE version="1.0">
    <Documento>
        <Encabezado>
            <IdDoc>
                <TipoDTE>33</TipoDTE>
                <Folio>1</Folio>
                <FechaEmision>2024-03-20</FechaEmision>
            </IdDoc>
            <Emisor>
                <RUTEmisor>76555555-5</RUTEmisor>
                <RznSoc>Empresa de Prueba SpA</RznSoc>
                <GiroEmis>Servicios Informáticos</GiroEmis>
            </Emisor>
            <Receptor>
                <RUTRecep>66666666-6</RUTRecep>
                <RznSocRecep>Cliente de Prueba</RznSocRecep>
            </Receptor>
        </Encabezado>
    </Documento>
</DTE>`
)

func TestNewSIIFirmaService(t *testing.T) {
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
				RutEmpresa:      "76555555-5",
			},
			expectError: false,
		},
		{
			name: "RUT empresa inválido",
			config: &ConfiguracionFirma{
				RutaCertificado: certPath,
				Password:        password,
				RutEmpresa:      "invalid-rut",
			},
			expectError: true,
		},
		{
			name: "Sin RUT empresa",
			config: &ConfiguracionFirma{
				RutaCertificado: certPath,
				Password:        password,
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, err := NewSIIFirmaService(tt.config)
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, service)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, service)
				assert.NotNil(t, service.BaseFirmaService)
				assert.NotNil(t, service.certCache)
			}
		})
	}
}

func TestFirmarSemilla(t *testing.T) {
	certPath, password := setupTestCertificates(t)
	config := &ConfiguracionFirma{
		RutaCertificado: certPath,
		Password:        password,
		RutEmpresa:      "76555555-5",
	}

	service, err := NewSIIFirmaService(config)
	require.NoError(t, err)

	tests := []struct {
		name        string
		semilla     string
		expectError bool
	}{
		{
			name:        "Semilla válida",
			semilla:     testSemilla,
			expectError: false,
		},
		{
			name:        "Semilla vacía",
			semilla:     "",
			expectError: true,
		},
		{
			name:        "Semilla muy larga",
			semilla:     string(make([]byte, 1000)),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultado, err := service.FirmarSemilla(tt.semilla)
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, resultado)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resultado)
				assert.Contains(t, resultado.XMLFirmado, "<SemillaXML>")
				assert.Contains(t, resultado.XMLFirmado, tt.semilla)
				assert.Contains(t, resultado.XMLFirmado, "<ds:Signature")
				assert.NotEmpty(t, resultado.DigestValue)
				assert.NotEmpty(t, resultado.SignatureValue)
			}
		})
	}
}

func TestFirmarToken(t *testing.T) {
	certPath, password := setupTestCertificates(t)
	config := &ConfiguracionFirma{
		RutaCertificado: certPath,
		Password:        password,
		RutEmpresa:      "76555555-5",
	}

	service, err := NewSIIFirmaService(config)
	require.NoError(t, err)

	tests := []struct {
		name        string
		token       string
		expectError bool
	}{
		{
			name:        "Token válido",
			token:       testToken,
			expectError: false,
		},
		{
			name:        "Token vacío",
			token:       "",
			expectError: true,
		},
		{
			name:        "Token muy largo",
			token:       string(make([]byte, 1000)),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultado, err := service.FirmarToken(tt.token)
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, resultado)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resultado)
				assert.Contains(t, resultado.XMLFirmado, "<TokenXML>")
				assert.Contains(t, resultado.XMLFirmado, tt.token)
				assert.Contains(t, resultado.XMLFirmado, "<ds:Signature")
				assert.NotEmpty(t, resultado.DigestValue)
				assert.NotEmpty(t, resultado.SignatureValue)
			}
		})
	}
}

func TestValidarCAF(t *testing.T) {
	certPath, password := setupTestCertificates(t)
	config := &ConfiguracionFirma{
		RutaCertificado: certPath,
		Password:        password,
		RutEmpresa:      "76555555-5",
	}

	service, err := NewSIIFirmaService(config)
	require.NoError(t, err)

	tests := []struct {
		name        string
		caf         string
		expectError bool
	}{
		{
			name:        "CAF válido",
			caf:         testCAF,
			expectError: false,
		},
		{
			name:        "CAF vacío",
			caf:         "",
			expectError: true,
		},
		{
			name:        "CAF inválido",
			caf:         "<InvalidCAF>",
			expectError: true,
		},
		{
			name: "CAF con RUT diferente",
			caf: `<?xml version="1.0" encoding="UTF-8"?>
<CAF version="1.0">
    <DA>
        <RE>77777777-7</RE>
        <RS>OTRA EMPRESA</RS>
        <TD>33</TD>
        <RNG><D>1</D><H>100</H></RNG>
        <FA>2024-03-20</FA>
        <RSAPK><M>123</M><E>456</E></RSAPK>
        <IDK>1</IDK>
    </DA>
    <FRMA>FIRMA_SII</FRMA>
</CAF>`,
			expectError: true,
		},
		{
			name: "CAF expirado",
			caf: `<?xml version="1.0" encoding="UTF-8"?>
<CAF version="1.0">
    <DA>
        <RE>76555555-5</RE>
        <RS>EMPRESA DE PRUEBA</RS>
        <TD>33</TD>
        <RNG><D>1</D><H>100</H></RNG>
        <FA>2020-01-01</FA>
        <RSAPK><M>123</M><E>456</E></RSAPK>
        <IDK>1</IDK>
    </DA>
    <FRMA>FIRMA_SII</FRMA>
</CAF>`,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.ValidarCAF([]byte(tt.caf))
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestFirmarDTE(t *testing.T) {
	certPath, password := setupTestCertificates(t)
	config := &ConfiguracionFirma{
		RutaCertificado: certPath,
		Password:        password,
		RutEmpresa:      "76555555-5",
	}

	service, err := NewSIIFirmaService(config)
	require.NoError(t, err)

	tests := []struct {
		name        string
		xml         string
		expectError bool
	}{
		{
			name:        "DTE válido",
			xml:         testDTE,
			expectError: false,
		},
		{
			name:        "XML vacío",
			xml:         "",
			expectError: true,
		},
		{
			name:        "XML inválido",
			xml:         "<InvalidXML>",
			expectError: true,
		},
		{
			name:        "DTE sin RUT emisor",
			xml:         strings.Replace(testDTE, "<RUTEmisor>76555555-5</RUTEmisor>", "", 1),
			expectError: true,
		},
		{
			name:        "DTE con RUT emisor diferente",
			xml:         strings.Replace(testDTE, "76555555-5", "77777777-7", 1),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultado, err := service.FirmarDTE(tt.xml)
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, resultado)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resultado)

				// Verificar estructura del XML firmado
				assert.Contains(t, resultado.XMLFirmado, "<DTE")
				assert.Contains(t, resultado.XMLFirmado, "<ds:Signature")
				assert.Contains(t, resultado.XMLFirmado, "<ds:SignedInfo")
				assert.Contains(t, resultado.XMLFirmado, "<ds:SignatureValue")
				assert.Contains(t, resultado.XMLFirmado, "<ds:KeyInfo")

				// Verificar que la firma es válida
				estado, err := service.ValidarFirma(resultado.XMLFirmado)
				assert.NoError(t, err)
				assert.True(t, estado.Valida)

				// Verificar que el RUT del emisor coincide con el certificado
				var dte struct {
					Documento struct {
						Encabezado struct {
							Emisor struct {
								RUTEmisor string
							}
						}
					}
				}
				err = xml.Unmarshal([]byte(resultado.XMLFirmado), &dte)
				assert.NoError(t, err)
				assert.Equal(t, config.RutEmpresa, dte.Documento.Encabezado.Emisor.RUTEmisor)
			}
		})
	}
}

func TestCertCache(t *testing.T) {
	certPath, password := setupTestCertificates(t)
	config := &ConfiguracionFirma{
		RutaCertificado: certPath,
		Password:        password,
		RutEmpresa:      "76555555-5",
	}

	service, err := NewSIIFirmaService(config)
	require.NoError(t, err)

	// Obtener certificado (debería guardarse en caché)
	cert1, err := service.ObtenerCertificado()
	require.NoError(t, err)
	require.NotNil(t, cert1)

	// Obtener nuevamente (debería venir del caché)
	cert2, err := service.ObtenerCertificado()
	require.NoError(t, err)
	require.NotNil(t, cert2)

	// Verificar que son el mismo certificado
	assert.Equal(t, cert1.SerialNumber, cert2.SerialNumber)

	// Verificar que el certificado está en caché
	cachedCert := service.certCache.Get(config.RutaCertificado)
	assert.NotNil(t, cachedCert)
	assert.Equal(t, cert1.SerialNumber, cachedCert.SerialNumber)
}

func TestConcurrencia(t *testing.T) {
	certPath, password := setupTestCertificates(t)
	config := &ConfiguracionFirma{
		RutaCertificado: certPath,
		Password:        password,
		RutEmpresa:      "76555555-5",
	}

	service, err := NewSIIFirmaService(config)
	require.NoError(t, err)

	// Probar firmas concurrentes de diferentes tipos de documentos
	numGoroutines := 10
	resultados := make(chan error, numGoroutines*3) // Para semilla, token y DTE

	for i := 0; i < numGoroutines; i++ {
		go func() {
			_, err := service.FirmarSemilla(testSemilla)
			resultados <- err
		}()

		go func() {
			_, err := service.FirmarToken(testToken)
			resultados <- err
		}()

		go func() {
			_, err := service.FirmarDocumento(testCAF)
			resultados <- err
		}()
	}

	// Verificar resultados
	for i := 0; i < numGoroutines*3; i++ {
		err := <-resultados
		assert.NoError(t, err)
	}
}
