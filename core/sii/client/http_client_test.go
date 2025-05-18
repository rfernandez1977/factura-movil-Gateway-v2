package client

import (
	"context"
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"FMgo/core/sii/infrastructure/certificates"
	"FMgo/core/sii/logger"
	"FMgo/core/sii/models"
	"github.com/stretchr/testify/assert"
)

func setupTestServer(t *testing.T, handler http.HandlerFunc) (*httptest.Server, *HTTPClient) {
	server := httptest.NewTLSServer(handler)

	config := &models.Config{
		SII: models.SIIConfig{
			BaseURL:    server.URL,
			CertPath:   "testdata/cert.pem",
			KeyPath:    "testdata/key.pem",
			RetryCount: 3,
			Timeout:    30,
		},
	}

	logger := logger.NewLogger()
	client, err := NewHTTPClient(config, logger)
	assert.NoError(t, err)

	return server, client
}

func TestObtenerSemilla(t *testing.T) {
	// Preparar respuesta mock
	expectedSemilla := "SEMILLA123"
	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)

		response := `<?xml version="1.0" encoding="UTF-8"?>
		<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
			<soap:Body>
				<getSeedResponse>
					<seed>` + expectedSemilla + `</seed>
				</getSeedResponse>
			</soap:Body>
		</soap:Envelope>`

		w.Header().Set("Content-Type", "text/xml")
		w.Write([]byte(response))
	}

	server, client := setupTestServer(t, handler)
	defer server.Close()

	// Ejecutar prueba
	ctx := context.Background()
	semilla, err := client.ObtenerSemilla(ctx)

	// Verificar resultados
	assert.NoError(t, err)
	assert.Equal(t, expectedSemilla, semilla)
}

func TestObtenerToken(t *testing.T) {
	// Preparar respuesta mock
	expectedToken := "TOKEN123"
	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)

		response := `<?xml version="1.0" encoding="UTF-8"?>
		<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
			<soap:Body>
				<getTokenResponse>
					<token>` + expectedToken + `</token>
				</getTokenResponse>
			</soap:Body>
		</soap:Envelope>`

		w.Header().Set("Content-Type", "text/xml")
		w.Write([]byte(response))
	}

	server, client := setupTestServer(t, handler)
	defer server.Close()

	// Ejecutar prueba
	ctx := context.Background()
	token, err := client.ObtenerToken(ctx, "SEMILLA123")

	// Verificar resultados
	assert.NoError(t, err)
	assert.Equal(t, expectedToken, token)
}

func TestEnviarDTE(t *testing.T) {
	// Preparar respuesta mock
	expectedTrackID := "TRACK123"
	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Contains(t, r.Header.Get("Cookie"), "TOKEN=")

		response := &models.RespuestaSII{
			TrackID: expectedTrackID,
			Estado:  "OK",
			Glosa:   "Envío exitoso",
		}

		w.Header().Set("Content-Type", "application/xml")
		xml.NewEncoder(w).Encode(response)
	}

	server, client := setupTestServer(t, handler)
	defer server.Close()

	// Ejecutar prueba
	ctx := context.Background()
	sobre := []byte("<DTE>...</DTE>")
	resp, err := client.EnviarDTE(ctx, sobre, "TOKEN123")

	// Verificar resultados
	assert.NoError(t, err)
	assert.Equal(t, expectedTrackID, resp.TrackID)
	assert.Equal(t, "OK", resp.Estado)
}

func TestConsultarEstado(t *testing.T) {
	// Preparar respuesta mock
	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Contains(t, r.URL.Query().Get("trackID"), "TRACK123")

		response := &models.EstadoSII{
			Estado:  "EPR",
			Glosa:   "Envío Procesado",
			TrackID: "TRACK123",
		}

		w.Header().Set("Content-Type", "application/xml")
		xml.NewEncoder(w).Encode(response)
	}

	server, client := setupTestServer(t, handler)
	defer server.Close()

	// Ejecutar prueba
	ctx := context.Background()
	estado, err := client.ConsultarEstado(ctx, "TRACK123")

	// Verificar resultados
	assert.NoError(t, err)
	assert.Equal(t, "EPR", estado.Estado)
	assert.Equal(t, "Envío Procesado", estado.Glosa)
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		name          string
		statusCode    int
		expectedError string
		responseBody  string
		contentType   string
	}{
		{
			name:          "Error de autenticación",
			statusCode:    http.StatusUnauthorized,
			expectedError: "Credenciales inválidas",
			contentType:   "text/xml",
		},
		{
			name:          "Error de certificado",
			statusCode:    http.StatusForbidden,
			expectedError: "Certificado inválido",
			contentType:   "text/xml",
		},
		{
			name:          "Timeout",
			statusCode:    http.StatusRequestTimeout,
			expectedError: "Tiempo de espera agotado",
			contentType:   "text/xml",
		},
		{
			name:          "Servidor no disponible",
			statusCode:    http.StatusServiceUnavailable,
			expectedError: "Servicio no disponible",
			contentType:   "text/xml",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", tt.contentType)
				w.WriteHeader(tt.statusCode)
				if tt.responseBody != "" {
					w.Write([]byte(tt.responseBody))
				}
			}

			server, client := setupTestServer(t, handler)
			defer server.Close()

			ctx := context.Background()
			_, err := client.ObtenerSemilla(ctx)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectedError)
		})
	}
}

func TestHTTPClient_ObtenerSemilla(t *testing.T) {
	// Crear servidor de prueba con respuesta SOAP correcta
	server := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		// Verificar headers si es necesario
		if strings.HasSuffix(r.URL.Path, "DTEUpload") {
			assert.Equal(t, "application/xml", r.Header.Get("Content-Type"))
			assert.Equal(t, "Bearer TEST-TOKEN", r.Header.Get("Authorization"))
		}

		// Configurar la respuesta
		w.Header().Set("Content-Type", "application/xml; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
	<soap:Body>
		<getSeedResponse xmlns="http://DefaultNamespace">
			<seed>123456789</seed>
		</getSeedResponse>
	</soap:Body>
</soap:Envelope>`))
	})
	defer server.Close()

	// Crear cliente de prueba y reemplazar la URL base
	client := &HTTPClient{
		client:   server.Client(),
		ambiente: models.Certificacion,
	}
	models.URLSemillaCert = server.URL

	// Ejecutar prueba
	semilla, err := client.ObtenerSemilla(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, "123456789", semilla)
}

func TestHTTPClient_ObtenerToken(t *testing.T) {
	// Crear servidor de prueba con respuesta SOAP correcta
	server := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		// Verificar headers si es necesario
		if strings.HasSuffix(r.URL.Path, "DTEUpload") {
			assert.Equal(t, "application/xml", r.Header.Get("Content-Type"))
			assert.Equal(t, "Bearer TEST-TOKEN", r.Header.Get("Authorization"))
		}

		// Configurar la respuesta
		w.Header().Set("Content-Type", "application/xml; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
	<soap:Body>
		<getTokenResponse xmlns="http://DefaultNamespace">
			<token>ABC123</token>
		</getTokenResponse>
	</soap:Body>
</soap:Envelope>`))
	})
	defer server.Close()

	// Crear cliente de prueba y reemplazar la URL base
	client := &HTTPClient{
		client:   server.Client(),
		ambiente: models.Certificacion,
	}
	models.URLTokenCert = server.URL

	// Ejecutar prueba
	token, err := client.ObtenerToken(context.Background(), "123456789")
	assert.NoError(t, err)
	assert.Equal(t, "ABC123", token)
}

func TestHTTPClient_EnviarDTE(t *testing.T) {
	// Crear servidor de prueba con respuesta SOAP correcta
	server := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		// Verificar headers si es necesario
		if strings.HasSuffix(r.URL.Path, "DTEUpload") {
			assert.Equal(t, "application/xml", r.Header.Get("Content-Type"))
			assert.Equal(t, "Bearer TEST-TOKEN", r.Header.Get("Authorization"))
		}

		// Configurar la respuesta
		w.Header().Set("Content-Type", "application/xml; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
	<soap:Body>
		<sendDTEResponse xmlns="http://DefaultNamespace">
			<estado>0</estado>
			<glosa>DTE Recibido</glosa>
			<trackid>123456</trackid>
		</sendDTEResponse>
	</soap:Body>
</soap:Envelope>`))
	})
	defer server.Close()

	// Crear cliente de prueba y reemplazar la URL base
	client := &HTTPClient{
		client:   server.Client(),
		ambiente: models.Certificacion,
	}
	models.URLEnvioDTECert = server.URL

	// Ejecutar prueba
	resp, err := client.EnviarDTE(context.Background(), []byte("<DTE></DTE>"), "TEST-TOKEN")
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "0", resp.Estado)
	assert.Equal(t, "DTE Recibido", resp.Glosa)
	assert.Equal(t, "123456", resp.TrackID)
}

func TestHTTPClient_ConsultarEstado(t *testing.T) {
	// Crear servidor de prueba con respuesta SOAP correcta
	server := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		// Verificar headers si es necesario
		if strings.HasSuffix(r.URL.Path, "DTEUpload") {
			assert.Equal(t, "application/xml", r.Header.Get("Content-Type"))
			assert.Equal(t, "Bearer TEST-TOKEN", r.Header.Get("Authorization"))
		}

		// Configurar la respuesta
		w.Header().Set("Content-Type", "application/xml; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
	<soap:Body>
		<getStatusResponse xmlns="http://DefaultNamespace">
			<estado>EPR</estado>
			<glosa>Envío Procesado</glosa>
			<trackid>123456</trackid>
		</getStatusResponse>
	</soap:Body>
</soap:Envelope>`))
	})
	defer server.Close()

	// Crear cliente de prueba y reemplazar la URL base
	client := &HTTPClient{
		client:   server.Client(),
		ambiente: models.Certificacion,
	}
	models.URLEstadoDTECert = server.URL

	// Ejecutar prueba
	estado, err := client.ConsultarEstado(context.Background(), "123456")
	assert.NoError(t, err)
	assert.NotNil(t, estado)
	assert.Equal(t, "EPR", estado.Estado)
	assert.Equal(t, "Envío Procesado", estado.Glosa)
	assert.Equal(t, "123456", estado.TrackID)
}

func TestHTTPClient_ConsultarDTE(t *testing.T) {
	// Crear servidor de prueba con respuesta SOAP correcta
	server := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		// Verificar headers si es necesario
		if strings.HasSuffix(r.URL.Path, "DTEUpload") {
			assert.Equal(t, "application/xml", r.Header.Get("Content-Type"))
			assert.Equal(t, "Bearer TEST-TOKEN", r.Header.Get("Authorization"))
		}

		// Configurar la respuesta
		w.Header().Set("Content-Type", "application/xml; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
	<soap:Body>
		<getEstDteResponse xmlns="http://DefaultNamespace">
			<estado>DTE_RECIBIDO</estado>
			<glosa>DTE Recibido</glosa>
			<trackid>123456</trackid>
		</getEstDteResponse>
	</soap:Body>
</soap:Envelope>`))
	})
	defer server.Close()

	// Crear cliente de prueba y reemplazar la URL base
	client := &HTTPClient{
		client:   server.Client(),
		ambiente: models.Certificacion,
	}
	models.URLEstadoDTECert = server.URL

	// Ejecutar prueba
	estado, err := client.ConsultarDTE(context.Background(), "33", 1234, "76212889-6")
	assert.NoError(t, err)
	assert.NotNil(t, estado)
	assert.Equal(t, "DTE_RECIBIDO", estado.Estado)
	assert.Equal(t, "DTE Recibido", estado.Glosa)
	assert.Equal(t, "123456", estado.TrackID)
}

func TestHTTPClient_VerificarComunicacion(t *testing.T) {
	// Crear servidor de prueba con respuesta SOAP correcta
	server := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		// Verificar headers si es necesario
		if strings.HasSuffix(r.URL.Path, "DTEUpload") {
			assert.Equal(t, "application/xml", r.Header.Get("Content-Type"))
			assert.Equal(t, "Bearer TEST-TOKEN", r.Header.Get("Authorization"))
		}

		// Configurar la respuesta
		w.Header().Set("Content-Type", "application/xml; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
	<soap:Body>
		<getSeedResponse xmlns="http://DefaultNamespace">
			<seed>123456789</seed>
		</getSeedResponse>
	</soap:Body>
</soap:Envelope>`))
	})
	defer server.Close()

	// Crear cliente de prueba y reemplazar la URL base
	client := &HTTPClient{
		client:   server.Client(),
		ambiente: models.Certificacion,
	}
	models.URLSemillaCert = server.URL

	// Ejecutar prueba
	err := client.VerificarComunicacion(context.Background())
	assert.NoError(t, err)
}

func TestNewHTTPClient(t *testing.T) {
	t.Run("error al crear con certificado inválido", func(t *testing.T) {
		client, err := NewHTTPClient("ruta/invalida.p12", "password", models.Produccion, DefaultRetryConfig())
		assert.Error(t, err)
		assert.Nil(t, client)
		assert.Contains(t, err.Error(), "error al crear gestor de certificados")
	})
}

func TestObtenerSemilla(t *testing.T) {
	t.Run("error por certificado inválido", func(t *testing.T) {
		mockCertManager := &mockCertManager{
			validateErr: fmt.Errorf("certificado expirado"),
		}
		client := &HTTPClient{
			certManager: mockCertManager,
			ambiente:    models.Produccion,
			retry:       DefaultRetryConfig(),
		}

		_, err := client.ObtenerSemilla(context.Background())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error de certificado")
	})
}

func TestObtenerToken(t *testing.T) {
	t.Run("error por certificado inválido", func(t *testing.T) {
		mockCertManager := &mockCertManager{
			validateErr: fmt.Errorf("certificado expirado"),
		}
		client := &HTTPClient{
			certManager: mockCertManager,
			ambiente:    models.Produccion,
			retry:       DefaultRetryConfig(),
		}

		_, err := client.ObtenerToken(context.Background(), "semilla-test")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error de certificado")
	})
}

// mockCertManager es un mock del gestor de certificados para pruebas
type mockCertManager struct {
	validateErr error
	certInfo    *certificates.CertificateInfo
}

func (m *mockCertManager) ValidateCertificate() error {
	return m.validateErr
}

func (m *mockCertManager) GetCertificateInfo() *certificates.CertificateInfo {
	return m.certInfo
}

func (m *mockCertManager) IsExpiringSoon(days int) bool {
	return false
}

func (m *mockCertManager) GetTLSConfig() *tls.Config {
	return &tls.Config{}
}
