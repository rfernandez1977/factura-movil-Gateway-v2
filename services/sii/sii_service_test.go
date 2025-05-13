// Package sii implementa el servicio de integración con el SII (Servicio de Impuestos Internos)
// para el envío y consulta de documentos tributarios electrónicos (DTE).
//
// Este paquete proporciona una interfaz para:
// - Enviar DTE al SII
// - Consultar el estado de DTE enviados
// - Verificar la comunicación con el SII
// - Manejar diferentes tipos de DTE (Facturas, Boletas, Guías, etc.)
//
// Las pruebas incluyen:
// - Validación de parámetros y configuración
// - Manejo de diferentes tipos de DTE
// - Validación de RUT y formatos de folio
// - Manejo de errores de red y certificados
// - Pruebas de concurrencia
// - Validación de respuestas y errores del SII

package sii

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockHTTPClient implementa la interfaz de http.Client para testing
type MockHTTPClient struct {
	mock.Mock
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*http.Response), args.Error(1)
}

type mockTransport struct {
	mock *MockHTTPClient
}

func (t *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return t.mock.Do(req)
}

func setupTestFiles(t *testing.T) (string, string, func()) {
	// Crear directorio temporal
	tmpDir, err := os.MkdirTemp("", "sii_test")
	if err != nil {
		t.Fatalf("Error al crear directorio temporal: %v", err)
	}

	// Crear archivos de certificado y llave
	certFile := filepath.Join(tmpDir, "cert.pem")
	keyFile := filepath.Join(tmpDir, "key.pem")

	// Escribir contenido de prueba en los archivos
	err = os.WriteFile(certFile, []byte("-----BEGIN CERTIFICATE-----\nMIICWDCCAcGgAwIBAgIJAP8m9/rSSJRvMA0GCSqGSIb3DQEBCwUAMEUxCzAJBgNV\n-----END CERTIFICATE-----"), 0600)
	if err != nil {
		t.Fatalf("Error al escribir archivo de certificado: %v", err)
	}

	err = os.WriteFile(keyFile, []byte("-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQC9QFi8Rf0o5IIp\n-----END PRIVATE KEY-----"), 0600)
	if err != nil {
		t.Fatalf("Error al escribir archivo de llave: %v", err)
	}

	// Función de limpieza
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}

	return certFile, keyFile, cleanup
}

// TestNewSIIService prueba la creación de una nueva instancia del servicio SII.
// Verifica que los parámetros se configuren correctamente.
func TestNewSIIService(t *testing.T) {
	mockClient := &MockHTTPClient{}
	service := &SIIServiceImpl{
		baseURL:    "https://api.test.cl",
		token:      "test-token",
		ambiente:   "CERTIFICACION",
		httpClient: &http.Client{Transport: &mockTransport{mockClient}},
	}

	// Verificar que el servicio se creó correctamente
	assert.NotNil(t, service)
	assert.Equal(t, "https://api.test.cl", service.baseURL)
	assert.Equal(t, "test-token", service.token)
	assert.Equal(t, "CERTIFICACION", service.ambiente)
}

// TestConsultarEstado prueba la consulta del estado de un DTE.
// Incluye pruebas para:
// - Estado aceptado
// - Estado rechazado
// - Error en la consulta
func TestConsultarEstado(t *testing.T) {
	mockClient := &MockHTTPClient{}
	service := &SIIServiceImpl{
		baseURL:    "https://api.test.cl",
		token:      "test-token",
		ambiente:   "CERTIFICACION",
		httpClient: &http.Client{Transport: &mockTransport{mockClient}},
	}

	tests := []struct {
		name    string
		trackID string
		resp    *http.Response
		err     error
		want    *EstadoSII
		wantErr bool
	}{
		{
			name:    "Estado aceptado",
			trackID: "123",
			resp: &http.Response{
				StatusCode: http.StatusOK,
				Body: io.NopCloser(strings.NewReader(`{
					"estado": "ACEPTADO",
					"glosa": "Solicitud procesada correctamente",
					"track_id": "123",
					"fecha_proceso": "2024-03-20T10:00:00Z"
				}`)),
			},
			err: nil,
			want: &EstadoSII{
				Estado:  "ACEPTADO",
				Glosa:   "Solicitud procesada correctamente",
				TrackID: "123",
				Fecha:   time.Date(2024, 3, 20, 10, 0, 0, 0, time.UTC),
			},
			wantErr: false,
		},
		{
			name:    "Estado rechazado",
			trackID: "456",
			resp: &http.Response{
				StatusCode: http.StatusBadRequest,
				Body: io.NopCloser(strings.NewReader(`{
					"estado": "RECHAZADO",
					"glosa": "Solicitud rechazada",
					"track_id": "456",
					"fecha_proceso": "2024-03-20T10:00:00Z",
					"errores": [
						{
							"codigo": "001",
							"descripcion": "Error de validación",
							"detalle": "Detalle del error"
						}
					]
				}`)),
			},
			err: nil,
			want: &EstadoSII{
				Estado:  "RECHAZADO",
				Glosa:   "Solicitud rechazada",
				TrackID: "456",
				Fecha:   time.Date(2024, 3, 20, 10, 0, 0, 0, time.UTC),
			},
			wantErr: true,
		},
		{
			name:    "Error en consulta",
			trackID: "error",
			resp:    nil,
			err:     assert.AnError,
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Configurar mock
			mockClient.On("Do", mock.Anything).Return(tt.resp, tt.err)

			// Ejecutar prueba
			result, err := service.ConsultarEstado(tt.trackID)

			// Verificar resultados
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.want.Estado, result.Estado)
				assert.Equal(t, tt.want.Glosa, result.Glosa)
				assert.Equal(t, tt.want.TrackID, result.TrackID)
				assert.Equal(t, tt.want.Fecha, result.Fecha)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

// TestConsultarEstadoCasosLimite prueba casos límite en la consulta de estado.
// Incluye pruebas para:
// - TrackID vacío
// - Estado inválido en la respuesta
// - JSON malformado
// - Error 500 del servidor
func TestConsultarEstadoCasosLimite(t *testing.T) {
	mockClient := &MockHTTPClient{}
	service := &SIIServiceImpl{
		baseURL:    "https://api.test.cl",
		token:      "test-token",
		ambiente:   "CERTIFICACION",
		httpClient: &http.Client{Transport: &mockTransport{mockClient}},
	}

	tests := []struct {
		name    string
		trackID string
		resp    *http.Response
		err     error
		wantErr bool
	}{
		{
			name:    "TrackID vacío",
			trackID: "",
			resp:    nil,
			err:     nil,
			wantErr: true,
		},
		{
			name:    "Estado inválido",
			trackID: "123",
			resp: &http.Response{
				StatusCode: http.StatusOK,
				Body: io.NopCloser(strings.NewReader(`{
					"estado": "ESTADO_INVALIDO",
					"glosa": "Estado no reconocido",
					"track_id": "123",
					"fecha_proceso": "2024-03-20T10:00:00Z"
				}`)),
			},
			err:     nil,
			wantErr: true,
		},
		{
			name:    "JSON malformado",
			trackID: "123",
			resp: &http.Response{
				StatusCode: http.StatusOK,
				Body: io.NopCloser(strings.NewReader(`{
					"estado": "ACEPTADO",
					"glosa": "Solicitud procesada correctamente",
					"track_id": "123",
					"fecha_proceso": "fecha_invalida"
				}`)),
			},
			err:     nil,
			wantErr: true,
		},
		{
			name:    "Error 500",
			trackID: "123",
			resp: &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body: io.NopCloser(strings.NewReader(`{
					"error": "Error interno del servidor"
				}`)),
			},
			err:     nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.resp != nil {
				mockClient.On("Do", mock.Anything).Return(tt.resp, tt.err)
			}

			result, err := service.ConsultarEstado(tt.trackID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

// TestTimeoutHandling prueba el manejo de timeouts en las solicitudes.
// Verifica el comportamiento cuando:
// - La solicitud excede el tiempo máximo de espera
// - El servidor no responde en el tiempo esperado
func TestTimeoutHandling(t *testing.T) {
	mockClient := &MockHTTPClient{}
	service := &SIIServiceImpl{
		baseURL:  "https://api.test.cl",
		token:    "test-token",
		ambiente: "CERTIFICACION",
		httpClient: &http.Client{
			Transport: &mockTransport{mockClient},
			Timeout:   100 * time.Millisecond,
		},
	}

	tests := []struct {
		name    string
		trackID string
		delay   time.Duration
		wantErr bool
	}{
		{
			name:    "Timeout por demora del servidor",
			trackID: "123",
			delay:   200 * time.Millisecond,
			wantErr: true,
		},
		{
			name:    "Respuesta dentro del timeout",
			trackID: "456",
			delay:   50 * time.Millisecond,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Configurar mock para simular delay
			mockClient.On("Do", mock.Anything).Run(func(args mock.Arguments) {
				time.Sleep(tt.delay)
			}).Return(&http.Response{
				StatusCode: http.StatusOK,
				Body: io.NopCloser(strings.NewReader(`{
					"estado": "ACEPTADO",
					"glosa": "Solicitud procesada correctamente",
					"track_id": "` + tt.trackID + `",
					"fecha_proceso": "2024-03-20T10:00:00Z"
				}`)),
			}, nil)

			result, err := service.ConsultarEstado(tt.trackID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), "timeout")
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

// TestConcurrentRequests prueba el manejo de solicitudes concurrentes.
// Verifica que el servicio maneje correctamente múltiples solicitudes simultáneas.
func TestConcurrentRequests(t *testing.T) {
	mockClient := &MockHTTPClient{}
	service := &SIIServiceImpl{
		baseURL:    "https://api.test.cl",
		token:      "test-token",
		ambiente:   "CERTIFICACION",
		httpClient: &http.Client{Transport: &mockTransport{mockClient}},
	}

	// Número de solicitudes concurrentes
	numRequests := 10
	results := make(chan error, numRequests)
	var wg sync.WaitGroup

	// Configurar mock para responder a todas las solicitudes
	for i := 0; i < numRequests; i++ {
		trackID := fmt.Sprintf("track_%d", i)
		mockClient.On("Do", mock.Anything).Return(&http.Response{
			StatusCode: http.StatusOK,
			Body: io.NopCloser(strings.NewReader(fmt.Sprintf(`{
				"estado": "ACEPTADO",
				"glosa": "Solicitud procesada correctamente",
				"track_id": "%s",
				"fecha_proceso": "2024-03-20T10:00:00Z"
			}`, trackID))),
		}, nil)
	}

	// Lanzar solicitudes concurrentes
	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			trackID := fmt.Sprintf("track_%d", i)
			_, err := service.ConsultarEstado(trackID)
			results <- err
		}(i)
	}

	// Esperar a que todas las solicitudes terminen
	wg.Wait()
	close(results)

	// Verificar resultados
	for err := range results {
		assert.NoError(t, err)
	}

	// Verificar que el mock fue llamado el número correcto de veces
	mockClient.AssertNumberOfCalls(t, "Do", numRequests)
}

// TestErrorHandling prueba el manejo de diferentes tipos de errores.
// Incluye pruebas para:
// - Errores de red
// - Errores de timeout
// - Errores de certificado
// - Errores de validación
func TestErrorHandling(t *testing.T) {
	mockClient := &MockHTTPClient{}
	service := &SIIServiceImpl{
		baseURL:    "https://api.test.cl",
		token:      "test-token",
		ambiente:   "CERTIFICACION",
		httpClient: &http.Client{Transport: &mockTransport{mockClient}},
	}

	tests := []struct {
		name    string
		trackID string
		resp    *http.Response
		err     error
		wantErr bool
		errMsg  string
	}{
		{
			name:    "Error de red",
			trackID: "123",
			resp:    nil,
			err:     &url.Error{Op: "Get", URL: "https://api.test.cl", Err: fmt.Errorf("connection refused")},
			wantErr: true,
			errMsg:  "connection refused",
		},
		{
			name:    "Error de timeout",
			trackID: "456",
			resp:    nil,
			err:     &url.Error{Op: "Get", URL: "https://api.test.cl", Err: context.DeadlineExceeded},
			wantErr: true,
			errMsg:  "timeout",
		},
		{
			name:    "Error de certificado",
			trackID: "789",
			resp:    nil,
			err:     &url.Error{Op: "Get", URL: "https://api.test.cl", Err: fmt.Errorf("certificate is not valid")},
			wantErr: true,
			errMsg:  "certificate",
		},
		{
			name:    "Error de validación",
			trackID: "012",
			resp: &http.Response{
				StatusCode: http.StatusBadRequest,
				Body: io.NopCloser(strings.NewReader(`{
					"error": "Invalid track ID format"
				}`)),
			},
			err:     nil,
			wantErr: true,
			errMsg:  "Invalid track ID format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.resp != nil {
				mockClient.On("Do", mock.Anything).Return(tt.resp, tt.err)
			} else {
				mockClient.On("Do", mock.Anything).Return(nil, tt.err)
			}

			result, err := service.ConsultarEstado(tt.trackID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

// TestConfigurationValidation prueba la validación de la configuración del servicio.
// Verifica que:
// - Los parámetros requeridos estén presentes
// - Los valores sean válidos
// - Se manejen correctamente los valores por defecto
func TestConfigurationValidation(t *testing.T) {
	tests := []struct {
		name     string
		baseURL  string
		token    string
		ambiente string
		certFile string
		keyFile  string
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "Configuración válida",
			baseURL:  "https://api.test.cl",
			token:    "test-token",
			ambiente: "CERTIFICACION",
			certFile: "cert.pem",
			keyFile:  "key.pem",
			wantErr:  false,
		},
		{
			name:     "BaseURL vacía",
			baseURL:  "",
			token:    "test-token",
			ambiente: "CERTIFICACION",
			certFile: "cert.pem",
			keyFile:  "key.pem",
			wantErr:  true,
			errMsg:   "baseURL es requerido",
		},
		{
			name:     "Token vacío",
			baseURL:  "https://api.test.cl",
			token:    "",
			ambiente: "CERTIFICACION",
			certFile: "cert.pem",
			keyFile:  "key.pem",
			wantErr:  true,
			errMsg:   "token es requerido",
		},
		{
			name:     "Certificado no existe",
			baseURL:  "https://api.test.cl",
			token:    "test-token",
			ambiente: "CERTIFICACION",
			certFile: "no_existe.pem",
			keyFile:  "key.pem",
			wantErr:  true,
			errMsg:   "error cargando certificado",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, err := NewSIIService(tt.baseURL, tt.token, tt.ambiente, tt.certFile, tt.keyFile)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, service)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, service)
				assert.Equal(t, tt.baseURL, service.baseURL)
				assert.Equal(t, tt.token, service.token)
				assert.Equal(t, tt.ambiente, service.ambiente)
			}
		})
	}
}

// TestInputValidation prueba la validación de los parámetros de entrada.
// Verifica que:
// - Los parámetros requeridos estén presentes
// - Los valores sean válidos
// - Se manejen correctamente los valores por defecto
func TestInputValidation(t *testing.T) {
	mockClient := &MockHTTPClient{}
	service := &SIIServiceImpl{
		baseURL:    "https://api.test.cl",
		token:      "test-token",
		ambiente:   "CERTIFICACION",
		httpClient: &http.Client{Transport: &mockTransport{mockClient}},
	}

	tests := []struct {
		name    string
		trackID string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "TrackID válido",
			trackID: "123456789",
			wantErr: false,
		},
		{
			name:    "TrackID vacío",
			trackID: "",
			wantErr: true,
			errMsg:  "trackID es requerido",
		},
		{
			name:    "TrackID inválido",
			trackID: "abc",
			wantErr: true,
			errMsg:  "trackID inválido",
		},
		{
			name:    "TrackID muy largo",
			trackID: strings.Repeat("1", 100),
			wantErr: true,
			errMsg:  "trackID demasiado largo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.ConsultarEstado(tt.trackID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}
