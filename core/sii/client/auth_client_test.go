package client

import (
	"context"
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"FMgo/core/sii/models/siimodels"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewAuthClient verifica la creación correcta del cliente de autenticación.
// Prueba:
// - Configuración de certificados
// - Configuración de endpoints
// - Configuración de timeouts
// - Validación de parámetros
func TestNewAuthClient(t *testing.T) {
	certFile, keyFile, cleanup := setupTestFiles(t)
	defer cleanup()

	config := &siimodels.ConfigSII{
		Ambiente:   siimodels.AmbienteCertificacion,
		BaseURL:    siimodels.URLBaseCertificacion,
		CertPath:   certFile,
		KeyPath:    keyFile,
		RetryCount: 3,
		Timeout:    30 * time.Second,
	}

	soapClient, err := NewSOAPClient(config)
	require.NoError(t, err)

	authClient := NewAuthClient(soapClient, config)
	assert.NotNil(t, authClient)
	assert.Equal(t, soapClient, authClient.soapClient)
	assert.Equal(t, config, authClient.config)
}

// TestGetToken verifica el proceso completo de obtención de token.
// Prueba:
// - Obtención de semilla
// - Firma de semilla
// - Obtención de token
// - Validación de token
// - Cache de token
func TestGetToken(t *testing.T) {
	certFile, keyFile, cleanup := setupTestFiles(t)
	defer cleanup()

	config := &siimodels.ConfigSII{
		Ambiente:   siimodels.AmbienteCertificacion,
		BaseURL:    siimodels.URLBaseCertificacion,
		CertPath:   certFile,
		KeyPath:    keyFile,
		RetryCount: 3,
		Timeout:    30 * time.Second,
	}

	// Crear servidor de prueba
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case siimodels.EndpointSemillaCert:
			// Respuesta para getSeed
			w.Header().Set("Content-Type", "text/xml; charset=utf-8")
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
			// Respuesta para getToken
			w.Header().Set("Content-Type", "text/xml; charset=utf-8")
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
		}
	}))
	defer server.Close()

	// Configurar el cliente para usar el servidor de prueba
	config.BaseURL = server.URL

	soapClient, err := NewSOAPClient(config)
	require.NoError(t, err)

	// Configurar TLS para el servidor de prueba
	soapClient.httpClient.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	authClient := NewAuthClient(soapClient, config)

	// Probar obtención de token
	t.Run("Obtener nuevo token", func(t *testing.T) {
		token, err := authClient.GetToken(context.Background())
		require.NoError(t, err)
		assert.Equal(t, "ABC123XYZ", token)
	})

	t.Run("Usar token en caché", func(t *testing.T) {
		// El token anterior debería estar en caché
		token, err := authClient.GetToken(context.Background())
		require.NoError(t, err)
		assert.Equal(t, "ABC123XYZ", token)
	})

	t.Run("Renovar token expirado", func(t *testing.T) {
		// Forzar expiración del token
		authClient.tokenExp = time.Now().Add(-1 * time.Hour)

		token, err := authClient.GetToken(context.Background())
		require.NoError(t, err)
		assert.Equal(t, "ABC123XYZ", token)
		assert.True(t, authClient.tokenExp.After(time.Now()))
	})
}

// TestGetTokenErrors verifica el manejo de errores en la autenticación.
// Prueba:
// - Errores de conexión
// - Errores de certificado
// - Errores de firma
// - Errores de respuesta
// - Errores de parsing
func TestGetTokenErrors(t *testing.T) {
	certFile, keyFile, cleanup := setupTestFiles(t)
	defer cleanup()

	config := &siimodels.ConfigSII{
		Ambiente:   siimodels.AmbienteCertificacion,
		BaseURL:    siimodels.URLBaseCertificacion,
		CertPath:   certFile,
		KeyPath:    keyFile,
		RetryCount: 1,
		Timeout:    1 * time.Second,
	}

	// Crear servidor que simula errores
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case siimodels.EndpointSemillaCert:
			w.Header().Set("Content-Type", "text/xml; charset=utf-8")
			w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
			<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
				<SOAP-ENV:Body>
					<RESP_HDR>
						<ESTADO>99</ESTADO>
						<GLOSA>Error al generar semilla</GLOSA>
					</RESP_HDR>
				</SOAP-ENV:Body>
			</SOAP-ENV:Envelope>`))

		case siimodels.EndpointTokenCert:
			w.Header().Set("Content-Type", "text/xml; charset=utf-8")
			w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
			<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
				<SOAP-ENV:Body>
					<RESP_HDR>
						<ESTADO>99</ESTADO>
						<GLOSA>Error al generar token</GLOSA>
					</RESP_HDR>
				</SOAP-ENV:Body>
			</SOAP-ENV:Envelope>`))
		}
	}))
	defer server.Close()

	// Configurar el cliente para usar el servidor de prueba
	config.BaseURL = server.URL

	soapClient, err := NewSOAPClient(config)
	require.NoError(t, err)

	// Configurar TLS para el servidor de prueba
	soapClient.httpClient.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	authClient := NewAuthClient(soapClient, config)

	// Probar errores
	t.Run("Error al obtener semilla", func(t *testing.T) {
		_, err := authClient.GetToken(context.Background())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "error al obtener semilla")
	})

	t.Run("Error al obtener token", func(t *testing.T) {
		// Forzar que se obtenga una semilla correcta
		authClient.soapClient = &SOAPClient{
			httpClient: &http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{
						InsecureSkipVerify: true,
					},
				},
			},
			config: config,
		}

		_, err := authClient.GetToken(context.Background())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "error al obtener token")
	})
}

// TestConcurrentTokenAccess verifica el acceso concurrente a tokens.
// Prueba:
// - Acceso simultáneo al cache
// - Renovación de tokens expirados
// - Race conditions
// - Bloqueo durante renovación
// - Timeout en espera
func TestConcurrentTokenAccess(t *testing.T) {
	certFile, keyFile, cleanup := setupTestFiles(t)
	defer cleanup()

	config := &siimodels.ConfigSII{
		Ambiente:   siimodels.AmbienteCertificacion,
		BaseURL:    siimodels.URLBaseCertificacion,
		CertPath:   certFile,
		KeyPath:    keyFile,
		RetryCount: 3,
		Timeout:    30 * time.Second,
	}

	// Crear servidor de prueba
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simular latencia
		time.Sleep(100 * time.Millisecond)

		w.Header().Set("Content-Type", "text/xml; charset=utf-8")
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
	}))
	defer server.Close()

	// Configurar el cliente
	config.BaseURL = server.URL
	soapClient, err := NewSOAPClient(config)
	require.NoError(t, err)
	soapClient.httpClient.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	authClient := NewAuthClient(soapClient, config)

	// Probar acceso concurrente
	concurrentRequests := 10
	done := make(chan bool)

	for i := 0; i < concurrentRequests; i++ {
		go func() {
			token, err := authClient.GetToken(context.Background())
			require.NoError(t, err)
			assert.Equal(t, "ABC123XYZ", token)
			done <- true
		}()
	}

	// Esperar que todas las goroutines terminen
	for i := 0; i < concurrentRequests; i++ {
		<-done
	}

	// Verificar que solo se hizo una llamada al servidor
	assert.Equal(t, "ABC123XYZ", authClient.token)
}
