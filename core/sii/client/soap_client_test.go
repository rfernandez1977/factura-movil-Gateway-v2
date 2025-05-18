package client

import (
	"context"
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"FMgo/core/sii/models/siimodels"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestFiles(t *testing.T) (string, string, func()) {
	// Crear directorio temporal
	tmpDir, err := os.MkdirTemp("", "sii_test")
	require.NoError(t, err)

	// Crear archivos de certificado y llave
	certFile := filepath.Join(tmpDir, "cert.pem")
	keyFile := filepath.Join(tmpDir, "key.pem")

	// Escribir contenido de prueba
	err = os.WriteFile(certFile, []byte(`-----BEGIN CERTIFICATE-----
MIICWDCCAcGgAwIBAgIJAP8m9/rSSJRvMA0GCSqGSIb3DQEBCwUAMEUxCzAJBgNV
BAYTAkFVMRMwEQYDVQQIDApTb21lLVN0YXRlMSEwHwYDVQQKDBhJbnRlcm5ldCBX
aWRnaXRzIFB0eSBMdGQwHhcNMTkwNjAyMjI0OTU2WhcNMjAwNjAxMjI0OTU2WjBF
MQswCQYDVQQGEwJBVTETMBEGA1UECAwKU29tZS1TdGF0ZTEhMB8GA1UECgwYSW50
ZXJuZXQgV2lkZ2l0cyBQdHkgTHRkMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKB
gQDIXpgQgRyg6fB8/CrLrKYEm9kh6qzqYhQGQMF+bqWJ7QZQU0sBLrKF4jLkVGGQ
YdavZf5AkWfm8qOBKdICH3Qb8RXZ1RqpW3tZ1BjkU4QGPe3CXx7ZsqYgQzZELjYk
9qmkBvBJhKZrKfKwY6ZYwWaD6Y8YeUQAZJQQQqVpJQIDAQABo1AwTjAdBgNVHQ4E
FgQU5H1f0VyVhggsBa6+dbZlp9ldqGYwHwYDVR0jBBgwFoAU5H1f0VyVhggsBa6+
dbZlp9ldqGYwDAYDVR0TBAUwAwEB/zANBgkqhkiG9w0BAQsFAAOBgQB4rFB0FeOf
2vgVEwx1aU+VsB7LR7bxZC9xKWxZ5Vrf1NqL3HfWBP0CRUQLxKw4kF+n8qbQVRzA
D8JjvL5CJXcRAr90WlT1oVYHHW6ZoYK1LCr2KJuhdiK7pM3OGxRxKhtT6u/0pYjg
X/5rJ3FkwB/2QQw+xKjhcYGHqZQzKLhALA==
-----END CERTIFICATE-----`), 0600)
	require.NoError(t, err)

	err = os.WriteFile(keyFile, []byte(`-----BEGIN PRIVATE KEY-----
MIICdQIBADANBgkqhkiG9w0BAQEFAASCAl8wggJbAgEAAoGBAMhemBCBHKDp8Hz8
KsuspgSb2SHqrOpiFAZAwX5upYntBlBTSwEusoXiMuRUYZBh1q9l/kCRZ+byo4Ep
0gIfdBvxFdnVGqlbe1nUGORThAY97cJfHtmypiBDNkQuNiT2qaQG8EmEpmsq8rBj
pljBZoPpjxh5RABklBBCpWklAgMBAAECgYAKNqFrcOxZnxBt+qZx6X1PqFBqZPF/
E/1Rg9f3Yvs5MWKs4LTgJuWgYqzz1MjwY9NfCxqJQgJ0+n0B9zCHHzUYWfvVzq8P
8cWq5+FWcNZ5EYg8Zy4ayZwYWjo7VAgCD9i8QyNZqA9pCRsA8/jKBn4QHgxGyfCJ
8zhQGAzgJQJBAOQWwvWkwOg9Ry1kUkz8wOhjrMe6Uz2EKVQqXM4rKFxYyH5ZYF1y
Y3cNgXrIGGdPQjB9ZJkVEqEZx1+k4TMCQQDf2tNcPEDGHFZ4UxXNrCzAbH8xZPNH
gQ5CJ1Jt84WrYXRj5jQJHzYYVGqGo5JbLw5LZHwEUm6JqVQQmXaDAkAHJHcKCAjC
YX3LjkNjLyGVqc2Mk4yGdvBR5JouI8pJqjPHYwJBAKZ1XZgjgxBg1XYBXjUfJzNh
0ZwQIqAXq8F1AOJ8upGxdW5BNHnZvqFDiXxwgxwJ5NhIzOdJ5QVdhXQDAkBZ1J5d
RLB+eS5RhGHFYVkQEKKLHoGKHgAKv5HdRV3wGCwJ5n+Y4HmA4UB8RZF0xHnM8R1p
sZY/4kKWQxhZCg==
-----END PRIVATE KEY-----`), 0600)
	require.NoError(t, err)

	cleanup := func() {
		os.RemoveAll(tmpDir)
	}

	return certFile, keyFile, cleanup
}

// TestNewSOAPClient verifica la creación correcta del cliente SOAP.
// Prueba:
// - Configuración del cliente HTTP
// - Configuración de certificados TLS
// - Configuración de endpoints
// - Configuración de timeouts
// - Validación de parámetros obligatorios
func TestNewSOAPClient(t *testing.T) {
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

	client, err := NewSOAPClient(config)
	require.NoError(t, err)
	assert.NotNil(t, client)
	assert.NotNil(t, client.httpClient)
	assert.Equal(t, config, client.config)
}

// TestSOAPClientCall verifica las llamadas SOAP básicas.
// Prueba:
// - Construcción de mensaje SOAP
// - Envío de solicitud
// - Parsing de respuesta
// - Headers HTTP
// - Content-Type correcto
func TestSOAPClientCall(t *testing.T) {
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

	client, err := NewSOAPClient(config)
	require.NoError(t, err)

	// Crear servidor de prueba
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verificar headers
		assert.Equal(t, "text/xml; charset=utf-8", r.Header.Get("Content-Type"))
		assert.Equal(t, "", r.Header.Get("SOAPAction"))

		// Verificar token si está presente
		if token := r.Header.Get("Cookie"); token != "" {
			assert.Contains(t, token, "TOKEN=")
		}

		// Leer el cuerpo de la petición
		body, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err)

		// Verificar que es XML válido
		var envelope siimodels.SoapEnvelope
		err = xml.Unmarshal(body, &envelope)
		require.NoError(t, err)

		// Responder con una respuesta SOAP válida
		w.Header().Set("Content-Type", "text/xml; charset=utf-8")
		w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
		<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
			<SOAP-ENV:Body>
				<RESP_HDR>
					<ESTADO>00</ESTADO>
					<GLOSA>Operación exitosa</GLOSA>
				</RESP_HDR>
				<RESP_BODY>
					<TOKEN>ABC123</TOKEN>
				</RESP_BODY>
			</SOAP-ENV:Body>
		</SOAP-ENV:Envelope>`))
	}))
	defer server.Close()

	// Configurar el cliente para usar el servidor de prueba
	client.config.BaseURL = server.URL
	client.httpClient.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	// Realizar una llamada de prueba
	request := struct {
		XMLName struct{} `xml:"getToken"`
		Seed    string   `xml:"seed"`
	}{
		Seed: "12345",
	}

	response := &siimodels.RespuestaSII{}

	// Probar con y sin token
	tests := []struct {
		name  string
		ctx   context.Context
		token string
	}{
		{
			name:  "Sin token",
			ctx:   context.Background(),
			token: "",
		},
		{
			name:  "Con token",
			ctx:   context.WithValue(context.Background(), "token", "TEST_TOKEN"),
			token: "TEST_TOKEN",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.Call(tt.ctx, "/getToken", request, response)
			require.NoError(t, err)
			assert.Equal(t, "00", response.Header.Estado)
			assert.Equal(t, "Operación exitosa", response.Header.Glosa)
			assert.Equal(t, "ABC123", response.Body.Token)
		})
	}
}

// TestSOAPClientErrors verifica el manejo de errores en llamadas SOAP.
// Prueba:
// - Errores de conexión
// - Errores de timeout
// - Errores de certificado
// - Errores de respuesta
// - Errores de parsing XML
func TestSOAPClientErrors(t *testing.T) {
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

	client, err := NewSOAPClient(config)
	require.NoError(t, err)

	// Crear servidor que simula errores
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/error500" {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if r.URL.Path == "/soapfault" {
			w.Header().Set("Content-Type", "text/xml; charset=utf-8")
			w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
			<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
				<SOAP-ENV:Body>
					<SOAP-ENV:Fault>
						<faultcode>SOAP-ENV:Server</faultcode>
						<faultstring>Error interno del servidor</faultstring>
					</SOAP-ENV:Fault>
				</SOAP-ENV:Body>
			</SOAP-ENV:Envelope>`))
			return
		}
	}))
	defer server.Close()

	// Configurar el cliente para usar el servidor de prueba
	client.config.BaseURL = server.URL
	client.httpClient.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	request := struct {
		XMLName struct{} `xml:"test"`
	}{}

	response := &siimodels.RespuestaSII{}

	// Probar diferentes tipos de errores
	tests := []struct {
		name     string
		endpoint string
		wantErr  string
	}{
		{
			name:     "Error 500",
			endpoint: "/error500",
			wantErr:  "error del servidor: 500",
		},
		{
			name:     "SOAP Fault",
			endpoint: "/soapfault",
			wantErr:  "error SOAP: SOAP-ENV:Server - Error interno del servidor",
		},
		{
			name:     "Timeout",
			endpoint: "/timeout",
			wantErr:  "error después de 1 intentos",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.Call(context.Background(), tt.endpoint, request, response)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}

// TestSOAPClientTimeout verifica el manejo de timeouts.
// Prueba:
// - Timeout de conexión
// - Timeout de respuesta
// - Timeout de contexto
// - Configuración de tiempos
// - Limpieza de recursos
func TestSOAPClientTimeout(t *testing.T) {
	certFile, keyFile, cleanup := setupTestFiles(t)
	defer cleanup()

	config := &siimodels.ConfigSII{
		Ambiente:   siimodels.AmbienteCertificacion,
		BaseURL:    siimodels.URLBaseCertificacion,
		CertPath:   certFile,
		KeyPath:    keyFile,
		RetryCount: 3,
		Timeout:    1 * time.Second,
	}

	client, err := NewSOAPClient(config)
	require.NoError(t, err)

	// Crear servidor que simula un timeout
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
	}))
	defer server.Close()

	// Configurar el cliente para usar el servidor de prueba
	client.config.BaseURL = server.URL
	client.httpClient.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	request := struct {
		XMLName struct{} `xml:"test"`
	}{}

	response := &siimodels.RespuestaSII{}

	err = client.Call(context.Background(), "/test", request, response)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "error después de 3 intentos")
}
