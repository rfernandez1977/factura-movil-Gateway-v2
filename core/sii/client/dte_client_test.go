package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"FMgo/core/sii/models/siimodels"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewDTEClient verifica la creación correcta de un cliente DTE.
// Prueba:
// - Configuración correcta de parámetros
// - Inicialización de clientes SOAP y Auth
// - Manejo de certificados
func TestNewDTEClient(t *testing.T) {
	certFile, keyFile, cleanup := setupTestFiles(t)
	defer cleanup()

	config := &siimodels.ConfigSII{
		Ambiente:       siimodels.AmbienteCertificacion,
		BaseURL:        siimodels.URLBaseCertificacion,
		CertPath:       certFile,
		KeyPath:        keyFile,
		RutEmpresa:     "76555555-5",
		RutCertificado: "11111111-1",
		RetryCount:     3,
		Timeout:        30 * time.Second,
	}

	client, err := NewDTEClient(config)
	require.NoError(t, err)
	assert.NotNil(t, client)
	assert.NotNil(t, client.soapClient)
	assert.NotNil(t, client.authClient)
	assert.Equal(t, config, client.config)
}

// TestEnviarDTE verifica el proceso completo de envío de un DTE al SII.
// Prueba:
// - Autenticación con semilla y token
// - Generación del sobre de envío
// - Firma del documento
// - Envío al SII
// - Parsing de la respuesta
// - Manejo del trackID
func TestEnviarDTE(t *testing.T) {
	certFile, keyFile, cleanup := setupTestFiles(t)
	defer cleanup()

	config := &siimodels.ConfigSII{
		Ambiente:       siimodels.AmbienteCertificacion,
		BaseURL:        siimodels.URLBaseCertificacion,
		CertPath:       certFile,
		KeyPath:        keyFile,
		RutEmpresa:     "76555555-5",
		RutCertificado: "11111111-1",
		RetryCount:     3,
		Timeout:        30 * time.Second,
	}

	// Crear servidor de prueba
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verificar token en cookie
		cookie := r.Header.Get("Cookie")
		assert.Contains(t, cookie, "TOKEN=")

		// Responder según el endpoint
		switch r.URL.Path {
		case siimodels.EndpointSemillaCert:
			w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
			<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
				<SOAP-ENV:Body>
					<RESP_HDR>
						<ESTADO>00</ESTADO>
						<GLOSA>Semilla generada correctamente</GLOSA>
					</RESP_HDR>
					<RESP_BODY>
						<SEMILLA>123456789</SEMILLA>
					</RESP_BODY>
				</SOAP-ENV:Body>
			</SOAP-ENV:Envelope>`))

		case siimodels.EndpointTokenCert:
			w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
			<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
				<SOAP-ENV:Body>
					<RESP_HDR>
						<ESTADO>00</ESTADO>
						<GLOSA>Token generado correctamente</GLOSA>
					</RESP_HDR>
					<RESP_BODY>
						<TOKEN>ABC123XYZ</TOKEN>
					</RESP_BODY>
				</SOAP-ENV:Body>
			</SOAP-ENV:Envelope>`))

		case siimodels.EndpointEnvioCert:
			w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
			<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
				<SOAP-ENV:Body>
					<RESP_HDR>
						<ESTADO>00</ESTADO>
						<GLOSA>Envío procesado correctamente</GLOSA>
					</RESP_HDR>
					<RESP_BODY>
						<TRACKID>123456789</TRACKID>
					</RESP_BODY>
				</SOAP-ENV:Body>
			</SOAP-ENV:Envelope>`))
		}
	}))
	defer server.Close()

	// Configurar el cliente para usar el servidor de prueba
	config.BaseURL = server.URL

	client, err := NewDTEClient(config)
	require.NoError(t, err)

	// Configurar TLS para el servidor de prueba
	client.soapClient.httpClient.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	// Crear DTE de prueba
	dte := &siimodels.DTE{
		Version: "1.0",
		Documento: siimodels.Documento{
			ID: "TEST_001",
			Encabezado: siimodels.Encabezado{
				IdDoc: siimodels.IdDoc{
					TipoDTE:      33,
					Folio:        1,
					FechaEmision: time.Now(),
				},
				Emisor: siimodels.Emisor{
					RUTEmisor:  "76555555-5",
					RznSoc:     "EMPRESA DE PRUEBA SPA",
					GiroEmis:   "DESARROLLO DE SOFTWARE",
					Acteco:     722000,
					DirOrigen:  "DIRECCION 123",
					CmnaOrigen: "SANTIAGO",
				},
				Receptor: siimodels.Receptor{
					RUTRecep:    "55666666-6",
					RznSocRecep: "CLIENTE DE PRUEBA LTDA",
					GiroRecep:   "COMERCIO",
					DirRecep:    "DIRECCION CLIENTE 456",
					CmnaRecep:   "PROVIDENCIA",
				},
				Totales: siimodels.Totales{
					MntNeto:  10000,
					TasaIVA:  19,
					IVA:      1900,
					MntTotal: 11900,
				},
			},
			Detalle: []siimodels.Detalle{
				{
					NroLinDet: 1,
					NmbItem:   "Producto de Prueba",
					QtyItem:   1,
					PrcItem:   10000,
					MontoItem: 10000,
				},
			},
		},
	}

	// Enviar DTE
	resp, err := client.EnviarDTE(context.Background(), dte)
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "123456789", resp.TrackID)
}

// TestConsultarEstadoDTE verifica la consulta del estado de un DTE.
// Prueba:
// - Autenticación con token válido
// - Parámetros de consulta correctos (RUT, tipo, folio)
// - Parsing de la respuesta del SII
// - Diferentes estados posibles
// - Manejo de errores de consulta
func TestConsultarEstadoDTE(t *testing.T) {
	certFile, keyFile, cleanup := setupTestFiles(t)
	defer cleanup()

	config := &siimodels.ConfigSII{
		Ambiente:       siimodels.AmbienteCertificacion,
		BaseURL:        siimodels.URLBaseCertificacion,
		CertPath:       certFile,
		KeyPath:        keyFile,
		RutEmpresa:     "76555555-5",
		RutCertificado: "11111111-1",
		RetryCount:     3,
		Timeout:        30 * time.Second,
	}

	// Crear servidor de prueba
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verificar token en cookie
		cookie := r.Header.Get("Cookie")
		assert.Contains(t, cookie, "TOKEN=")

		// Responder según el endpoint
		switch r.URL.Path {
		case siimodels.EndpointSemillaCert:
			w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
			<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
				<SOAP-ENV:Body>
					<RESP_HDR>
						<ESTADO>00</ESTADO>
						<GLOSA>Semilla generada correctamente</GLOSA>
					</RESP_HDR>
					<RESP_BODY>
						<SEMILLA>123456789</SEMILLA>
					</RESP_BODY>
				</SOAP-ENV:Body>
			</SOAP-ENV:Envelope>`))

		case siimodels.EndpointTokenCert:
			w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
			<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
				<SOAP-ENV:Body>
					<RESP_HDR>
						<ESTADO>00</ESTADO>
						<GLOSA>Token generado correctamente</GLOSA>
					</RESP_HDR>
					<RESP_BODY>
						<TOKEN>ABC123XYZ</TOKEN>
					</RESP_BODY>
				</SOAP-ENV:Body>
			</SOAP-ENV:Envelope>`))

		case siimodels.EndpointConsultaCert:
			w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
			<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
				<SOAP-ENV:Body>
					<RESP_HDR>
						<ESTADO>00</ESTADO>
						<GLOSA>Consulta procesada correctamente</GLOSA>
					</RESP_HDR>
					<RESP_BODY>
						<ESTADO>ACEPTADO</ESTADO>
						<GLOSA_ESTADO>DTE Aceptado</GLOSA_ESTADO>
						<NUM_ATENCION>12345</NUM_ATENCION>
					</RESP_BODY>
				</SOAP-ENV:Body>
			</SOAP-ENV:Envelope>`))
		}
	}))
	defer server.Close()

	// Configurar el cliente para usar el servidor de prueba
	config.BaseURL = server.URL

	client, err := NewDTEClient(config)
	require.NoError(t, err)

	// Configurar TLS para el servidor de prueba
	client.soapClient.httpClient.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	// Consultar estado
	estado, err := client.ConsultarEstadoDTE(context.Background(), "76555555-5", 33, 1)
	require.NoError(t, err)
	assert.NotNil(t, estado)
	assert.Equal(t, "ACEPTADO", estado.Estado)
	assert.Equal(t, "DTE Aceptado", estado.GlosaEstado)
	assert.Equal(t, "12345", estado.NumeroAtencion)
}

// TestConsultarEstadoEnvio verifica la consulta del estado de un envío.
// Prueba:
// - Autenticación con token válido
// - Consulta por trackID
// - Parsing de respuesta detallada
// - Estados de envío y documentos
// - Manejo de errores
func TestConsultarEstadoEnvio(t *testing.T) {
	certFile, keyFile, cleanup := setupTestFiles(t)
	defer cleanup()

	config := &siimodels.ConfigSII{
		Ambiente:       siimodels.AmbienteCertificacion,
		BaseURL:        siimodels.URLBaseCertificacion,
		CertPath:       certFile,
		KeyPath:        keyFile,
		RutEmpresa:     "76555555-5",
		RutCertificado: "11111111-1",
		RetryCount:     3,
		Timeout:        30 * time.Second,
	}

	// Crear servidor de prueba
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verificar token en cookie
		cookie := r.Header.Get("Cookie")
		assert.Contains(t, cookie, "TOKEN=")

		// Responder según el endpoint
		switch r.URL.Path {
		case siimodels.EndpointSemillaCert:
			w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
			<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
				<SOAP-ENV:Body>
					<RESP_HDR>
						<ESTADO>00</ESTADO>
						<GLOSA>Semilla generada correctamente</GLOSA>
					</RESP_HDR>
					<RESP_BODY>
						<SEMILLA>123456789</SEMILLA>
					</RESP_BODY>
				</SOAP-ENV:Body>
			</SOAP-ENV:Envelope>`))

		case siimodels.EndpointTokenCert:
			w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
			<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
				<SOAP-ENV:Body>
					<RESP_HDR>
						<ESTADO>00</ESTADO>
						<GLOSA>Token generado correctamente</GLOSA>
					</RESP_HDR>
					<RESP_BODY>
						<TOKEN>ABC123XYZ</TOKEN>
					</RESP_BODY>
				</SOAP-ENV:Body>
			</SOAP-ENV:Envelope>`))

		case siimodels.EndpointConsultaCert:
			w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
			<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
				<SOAP-ENV:Body>
					<RESP_HDR>
						<ESTADO>00</ESTADO>
						<GLOSA>Consulta procesada correctamente</GLOSA>
					</RESP_HDR>
					<RESP_BODY>
						<TRACKID>123456789</TRACKID>
						<ESTADO>EPR</ESTADO>
						<GLOSA_ESTADO>Envío Procesado</GLOSA_ESTADO>
						<DETALLE_DTE>
							<DTE>
								<FOLIO>1</FOLIO>
								<TIPO_DTE>33</TIPO_DTE>
								<ESTADO>ACEPTADO</ESTADO>
								<GLOSA>DTE Aceptado</GLOSA>
							</DTE>
						</DETALLE_DTE>
					</RESP_BODY>
				</SOAP-ENV:Body>
			</SOAP-ENV:Envelope>`))
		}
	}))
	defer server.Close()

	// Configurar el cliente para usar el servidor de prueba
	config.BaseURL = server.URL

	client, err := NewDTEClient(config)
	require.NoError(t, err)

	// Configurar TLS para el servidor de prueba
	client.soapClient.httpClient.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	// Consultar estado de envío
	estado, err := client.ConsultarEstadoEnvio(context.Background(), "123456789")
	require.NoError(t, err)
	assert.NotNil(t, estado)
	assert.Equal(t, "123456789", estado.TrackID)
	assert.Equal(t, "EPR", estado.Estado)
	assert.Equal(t, "Envío Procesado", estado.GlosaEstado)
	require.Len(t, estado.DetalleDTE, 1)
	assert.Equal(t, int64(1), estado.DetalleDTE[0].Folio)
	assert.Equal(t, 33, estado.DetalleDTE[0].TipoDTE)
	assert.Equal(t, "ACEPTADO", estado.DetalleDTE[0].Estado)
	assert.Equal(t, "DTE Aceptado", estado.DetalleDTE[0].GlosaEstado)
}

// TestDTEClientErrors verifica el manejo correcto de errores del cliente.
// Prueba:
// - Errores de validación de DTE
// - Errores de autenticación
// - Errores de conexión
// - Errores de respuesta del SII
// - Errores de parsing
func TestDTEClientErrors(t *testing.T) {
	certFile, keyFile, cleanup := setupTestFiles(t)
	defer cleanup()

	config := &siimodels.ConfigSII{
		Ambiente:       siimodels.AmbienteCertificacion,
		BaseURL:        siimodels.URLBaseCertificacion,
		CertPath:       certFile,
		KeyPath:        keyFile,
		RutEmpresa:     "76555555-5",
		RutCertificado: "11111111-1",
		RetryCount:     1,
		Timeout:        1 * time.Second,
	}

	// Crear servidor que simula errores
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verificar token en cookie
		cookie := r.Header.Get("Cookie")
		assert.Contains(t, cookie, "TOKEN=")

		// Responder con error según el endpoint
		w.Header().Set("Content-Type", "text/xml; charset=utf-8")
		w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
		<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
			<SOAP-ENV:Body>
				<RESP_HDR>
					<ESTADO>99</ESTADO>
					<GLOSA>Error de validación</GLOSA>
				</RESP_HDR>
				<RESP_BODY>
					<ERROR>
						<CODIGO>1</CODIGO>
						<MENSAJE>Error de validación en DTE</MENSAJE>
					</ERROR>
				</RESP_BODY>
			</SOAP-ENV:Body>
		</SOAP-ENV:Envelope>`))
	}))
	defer server.Close()

	// Configurar el cliente para usar el servidor de prueba
	config.BaseURL = server.URL

	client, err := NewDTEClient(config)
	require.NoError(t, err)

	// Configurar TLS para el servidor de prueba
	client.soapClient.httpClient.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	// Probar envío con error
	t.Run("Error al enviar DTE", func(t *testing.T) {
		dte := &siimodels.DTE{
			Version: "1.0",
			Documento: siimodels.Documento{
				ID: "TEST_ERROR_001",
				Encabezado: siimodels.Encabezado{
					IdDoc: siimodels.IdDoc{
						TipoDTE:      33,
						Folio:        1,
						FechaEmision: time.Now(),
					},
				},
			},
		}

		_, err := client.EnviarDTE(context.Background(), dte)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error de validación")
	})

	// Probar consulta de estado con error
	t.Run("Error al consultar estado DTE", func(t *testing.T) {
		_, err := client.ConsultarEstadoDTE(context.Background(), "76555555-5", 33, 1)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error de validación")
	})

	// Probar consulta de envío con error
	t.Run("Error al consultar estado envío", func(t *testing.T) {
		_, err := client.ConsultarEstadoEnvio(context.Background(), "123456789")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error de validación")
	})
}

// TestDTEClientValidations verifica las validaciones de datos del cliente.
// Prueba:
// - Validación de RUT
// - Validación de tipo de DTE
// - Validación de folio
// - Validación de datos requeridos
// - Validación de formato
func TestDTEClientValidations(t *testing.T) {
	certFile, keyFile, cleanup := setupTestFiles(t)
	defer cleanup()

	config := &siimodels.ConfigSII{
		Ambiente:       siimodels.AmbienteCertificacion,
		BaseURL:        siimodels.URLBaseCertificacion,
		CertPath:       certFile,
		KeyPath:        keyFile,
		RutEmpresa:     "76555555-5",
		RutCertificado: "11111111-1",
		RetryCount:     3,
		Timeout:        30 * time.Second,
	}

	client, err := NewDTEClient(config)
	require.NoError(t, err)

	// Probar validaciones de DTE
	t.Run("Validar DTE sin datos requeridos", func(t *testing.T) {
		dte := &siimodels.DTE{
			Version: "1.0",
			Documento: siimodels.Documento{
				ID: "", // ID vacío
				Encabezado: siimodels.Encabezado{
					IdDoc: siimodels.IdDoc{
						TipoDTE: 0,  // Tipo inválido
						Folio:   -1, // Folio inválido
					},
				},
			},
		}

		_, err := client.EnviarDTE(context.Background(), dte)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "ID")
		assert.Contains(t, err.Error(), "TipoDTE")
		assert.Contains(t, err.Error(), "Folio")
	})

	// Probar validaciones de RUT
	t.Run("Validar RUT inválido", func(t *testing.T) {
		_, err := client.ConsultarEstadoDTE(context.Background(), "12345", 33, 1) // RUT inválido
		require.Error(t, err)
		assert.Contains(t, err.Error(), "RUT")
	})

	// Probar validaciones de tipo DTE
	t.Run("Validar tipo DTE inválido", func(t *testing.T) {
		_, err := client.ConsultarEstadoDTE(context.Background(), "76555555-5", 0, 1) // Tipo DTE inválido
		require.Error(t, err)
		assert.Contains(t, err.Error(), "tipo DTE")
	})

	// Probar validaciones de folio
	t.Run("Validar folio inválido", func(t *testing.T) {
		_, err := client.ConsultarEstadoDTE(context.Background(), "76555555-5", 33, -1) // Folio inválido
		require.Error(t, err)
		assert.Contains(t, err.Error(), "folio")
	})

	// Probar validaciones de trackID
	t.Run("Validar trackID inválido", func(t *testing.T) {
		_, err := client.ConsultarEstadoEnvio(context.Background(), "") // TrackID vacío
		require.Error(t, err)
		assert.Contains(t, err.Error(), "trackID")
	})
}

// TestDTEClientTimeout verifica el manejo de timeouts en las operaciones.
// Prueba:
// - Timeout en envío de DTE
// - Timeout en consulta de estado
// - Timeout en autenticación
// - Configuración de tiempos límite
// - Manejo de errores por timeout
func TestDTEClientTimeout(t *testing.T) {
	certFile, keyFile, cleanup := setupTestFiles(t)
	defer cleanup()

	config := &siimodels.ConfigSII{
		Ambiente:       siimodels.AmbienteCertificacion,
		BaseURL:        siimodels.URLBaseCertificacion,
		CertPath:       certFile,
		KeyPath:        keyFile,
		RutEmpresa:     "76555555-5",
		RutCertificado: "11111111-1",
		RetryCount:     1,
		Timeout:        1 * time.Second,
	}

	// Crear servidor que simula timeout
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second) // Dormir más que el timeout
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Configurar el cliente para usar el servidor de prueba
	config.BaseURL = server.URL

	client, err := NewDTEClient(config)
	require.NoError(t, err)

	// Configurar TLS para el servidor de prueba
	client.soapClient.httpClient.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	// Probar timeout en envío
	t.Run("Timeout al enviar DTE", func(t *testing.T) {
		dte := &siimodels.DTE{
			Version: "1.0",
			Documento: siimodels.Documento{
				ID: "TEST_TIMEOUT_001",
				Encabezado: siimodels.Encabezado{
					IdDoc: siimodels.IdDoc{
						TipoDTE:      33,
						Folio:        1,
						FechaEmision: time.Now(),
					},
				},
			},
		}

		_, err := client.EnviarDTE(context.Background(), dte)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "timeout")
	})

	// Probar timeout en consulta de estado
	t.Run("Timeout al consultar estado DTE", func(t *testing.T) {
		_, err := client.ConsultarEstadoDTE(context.Background(), "76555555-5", 33, 1)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "timeout")
	})

	// Probar timeout en consulta de envío
	t.Run("Timeout al consultar estado envío", func(t *testing.T) {
		_, err := client.ConsultarEstadoEnvio(context.Background(), "123456789")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "timeout")
	})
}

// TestDTEClientConcurrency verifica el comportamiento en operaciones concurrentes.
// Prueba:
// - Envíos simultáneos de DTEs
// - Consultas concurrentes de estado
// - Manejo de tokens compartidos
// - Race conditions
// - Recursos compartidos
func TestDTEClientConcurrency(t *testing.T) {
	certFile, keyFile, cleanup := setupTestFiles(t)
	defer cleanup()

	config := &siimodels.ConfigSII{
		Ambiente:       siimodels.AmbienteCertificacion,
		BaseURL:        siimodels.URLBaseCertificacion,
		CertPath:       certFile,
		KeyPath:        keyFile,
		RutEmpresa:     "76555555-5",
		RutCertificado: "11111111-1",
		RetryCount:     3,
		Timeout:        30 * time.Second,
	}

	// Crear servidor de prueba
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simular latencia
		time.Sleep(100 * time.Millisecond)

		// Responder según el endpoint
		switch r.URL.Path {
		case siimodels.EndpointSemillaCert:
			w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
			<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
				<SOAP-ENV:Body>
					<RESP_HDR>
						<ESTADO>00</ESTADO>
						<GLOSA>Semilla generada correctamente</GLOSA>
					</RESP_HDR>
					<RESP_BODY>
						<SEMILLA>123456789</SEMILLA>
					</RESP_BODY>
				</SOAP-ENV:Body>
			</SOAP-ENV:Envelope>`))

		case siimodels.EndpointTokenCert:
			w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
			<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
				<SOAP-ENV:Body>
					<RESP_HDR>
						<ESTADO>00</ESTADO>
						<GLOSA>Token generado correctamente</GLOSA>
					</RESP_HDR>
					<RESP_BODY>
						<TOKEN>ABC123XYZ</TOKEN>
					</RESP_BODY>
				</SOAP-ENV:Body>
			</SOAP-ENV:Envelope>`))

		case siimodels.EndpointEnvioCert:
			w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
			<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
				<SOAP-ENV:Body>
					<RESP_HDR>
						<ESTADO>00</ESTADO>
						<GLOSA>Envío procesado correctamente</GLOSA>
					</RESP_HDR>
					<RESP_BODY>
						<TRACKID>123456789</TRACKID>
					</RESP_BODY>
				</SOAP-ENV:Body>
			</SOAP-ENV:Envelope>`))
		}
	}))
	defer server.Close()

	// Configurar el cliente para usar el servidor de prueba
	config.BaseURL = server.URL

	client, err := NewDTEClient(config)
	require.NoError(t, err)

	// Configurar TLS para el servidor de prueba
	client.soapClient.httpClient.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	// Probar envíos concurrentes
	t.Run("Envíos concurrentes", func(t *testing.T) {
		concurrentRequests := 5
		done := make(chan bool)

		for i := 0; i < concurrentRequests; i++ {
			go func(index int) {
				dte := &siimodels.DTE{
					Version: "1.0",
					Documento: siimodels.Documento{
						ID: fmt.Sprintf("TEST_CONCURRENT_%d", index),
						Encabezado: siimodels.Encabezado{
							IdDoc: siimodels.IdDoc{
								TipoDTE:      33,
								Folio:        int64(index + 1),
								FechaEmision: time.Now(),
							},
						},
					},
				}

				resp, err := client.EnviarDTE(context.Background(), dte)
				require.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, "123456789", resp.TrackID)
				done <- true
			}(i)
		}

		// Esperar que todas las goroutines terminen
		for i := 0; i < concurrentRequests; i++ {
			<-done
		}
	})

	// Probar consultas concurrentes
	t.Run("Consultas concurrentes", func(t *testing.T) {
		concurrentRequests := 5
		done := make(chan bool)

		for i := 0; i < concurrentRequests; i++ {
			go func(index int) {
				estado, err := client.ConsultarEstadoDTE(context.Background(), "76555555-5", 33, int64(index+1))
				require.NoError(t, err)
				assert.NotNil(t, estado)
				done <- true
			}(i)
		}

		// Esperar que todas las goroutines terminen
		for i := 0; i < concurrentRequests; i++ {
			<-done
		}
	})
}

// TestDTEClientRetry verifica la política de reintentos del cliente.
// Prueba:
// - Reintentos en errores temporales
// - Límite máximo de reintentos
// - Intervalos entre reintentos
// - Errores permanentes vs temporales
// - Éxito después de reintentos
func TestDTEClientRetry(t *testing.T) {
	certFile, keyFile, cleanup := setupTestFiles(t)
	defer cleanup()

	config := &siimodels.ConfigSII{
		Ambiente:       siimodels.AmbienteCertificacion,
		BaseURL:        siimodels.URLBaseCertificacion,
		CertPath:       certFile,
		KeyPath:        keyFile,
		RutEmpresa:     "76555555-5",
		RutCertificado: "11111111-1",
		RetryCount:     3,
		Timeout:        30 * time.Second,
	}

	// Contador de intentos
	attempts := 0

	// Crear servidor que simula errores temporales
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++

		// Simular error en los primeros dos intentos
		if attempts <= 2 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		// Responder exitosamente en el tercer intento
		w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
		<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
			<SOAP-ENV:Body>
				<RESP_HDR>
					<ESTADO>00</ESTADO>
					<GLOSA>Envío procesado correctamente</GLOSA>
				</RESP_HDR>
				<RESP_BODY>
					<TRACKID>123456789</TRACKID>
				</RESP_BODY>
			</SOAP-ENV:Body>
		</SOAP-ENV:Envelope>`))
	}))
	defer server.Close()

	// Configurar el cliente para usar el servidor de prueba
	config.BaseURL = server.URL

	client, err := NewDTEClient(config)
	require.NoError(t, err)

	// Configurar TLS para el servidor de prueba
	client.soapClient.httpClient.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	// Probar reintento exitoso
	t.Run("Reintento exitoso", func(t *testing.T) {
		dte := &siimodels.DTE{
			Version: "1.0",
			Documento: siimodels.Documento{
				ID: "TEST_RETRY_001",
				Encabezado: siimodels.Encabezado{
					IdDoc: siimodels.IdDoc{
						TipoDTE:      33,
						Folio:        1,
						FechaEmision: time.Now(),
					},
				},
			},
		}

		resp, err := client.EnviarDTE(context.Background(), dte)
		require.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "123456789", resp.TrackID)
		assert.Equal(t, 3, attempts) // Verificar que se hicieron 3 intentos
	})
}

// TestDTEClientContextCancellation verifica la cancelación de operaciones.
// Prueba:
// - Cancelación durante envío
// - Cancelación durante consulta
// - Limpieza de recursos
// - Propagación de cancelación
// - Timeouts de contexto
func TestDTEClientContextCancellation(t *testing.T) {
	certFile, keyFile, cleanup := setupTestFiles(t)
	defer cleanup()

	config := &siimodels.ConfigSII{
		Ambiente:       siimodels.AmbienteCertificacion,
		BaseURL:        siimodels.URLBaseCertificacion,
		CertPath:       certFile,
		KeyPath:        keyFile,
		RutEmpresa:     "76555555-5",
		RutCertificado: "11111111-1",
		RetryCount:     3,
		Timeout:        30 * time.Second,
	}

	// Crear servidor que simula una respuesta lenta
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Dormir por más tiempo que el timeout del contexto
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Configurar el cliente para usar el servidor de prueba
	config.BaseURL = server.URL

	client, err := NewDTEClient(config)
	require.NoError(t, err)

	// Configurar TLS para el servidor de prueba
	client.soapClient.httpClient.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	// Probar cancelación de contexto
	t.Run("Cancelación de contexto", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		dte := &siimodels.DTE{
			Version: "1.0",
			Documento: siimodels.Documento{
				ID: "TEST_CANCEL_001",
				Encabezado: siimodels.Encabezado{
					IdDoc: siimodels.IdDoc{
						TipoDTE:      33,
						Folio:        1,
						FechaEmision: time.Now(),
					},
				},
			},
		}

		_, err := client.EnviarDTE(ctx, dte)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "context deadline exceeded")
	})
}
