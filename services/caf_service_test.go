package services

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"FMgo/config"
	"FMgo/models"
	"FMgo/tests/mocks"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRedisClient implementa la interfaz de redis.Client para testing
type MockRedisClient struct {
	mock.Mock
	redis.Client
}

func (m *MockRedisClient) Get(ctx context.Context, key string) *redis.StringCmd {
	args := m.Called(ctx, key)
	return args.Get(0).(*redis.StringCmd)
}

func (m *MockRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	args := m.Called(ctx, key, value, expiration)
	return args.Get(0).(*redis.StatusCmd)
}

// MockSIIService implementa la interfaz de sii.SIIService para testing
type MockSIIService struct {
	mock.Mock
}

// ConsultarEstado implementa el método de la interfaz SIIService
func (m *MockSIIService) ConsultarEstado(trackID string) (*models.EstadoSII, error) {
	args := m.Called(trackID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.EstadoSII), args.Error(1)
}

// EnviarDTE implementa el método de la interfaz SIIService
func (m *MockSIIService) EnviarDTE(dte *models.DTEXMLModel) (*models.RespuestaSII, error) {
	args := m.Called(dte)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.RespuestaSII), args.Error(1)
}

// ConsultarDTE implementa el método de la interfaz SIIService
func (m *MockSIIService) ConsultarDTE(trackID string) (*models.RespuestaSII, error) {
	args := m.Called(trackID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.RespuestaSII), args.Error(1)
}

// VerificarComunicacion implementa el método de la interfaz SIIService
func (m *MockSIIService) VerificarComunicacion() error {
	args := m.Called()
	return args.Error(0)
}

func TestNewCAFService(t *testing.T) {
	config := &config.SupabaseConfig{}
	cache := &MockRedisClient{}
	siiService := mocks.NewMockSIIService()

	service := NewCAFService(config, &cache.Client, siiService, "cert.pem", "key.pem", "certificacion", "https://api.test.cl")
	assert.NotNil(t, service)
	assert.Equal(t, config, service.config)
	assert.Equal(t, &cache.Client, service.cache)
	assert.Equal(t, siiService, service.siiService)
	assert.Equal(t, "certificacion", service.ambiente)
	assert.Equal(t, "https://api.test.cl", service.baseURL)
}

func TestGetCAFDisponible(t *testing.T) {
	config := &config.SupabaseConfig{}
	cache := &MockRedisClient{}
	siiService := &MockSIIService{}
	service := NewCAFService(config, &cache.Client, siiService, "cert.pem", "key.pem", "certificacion", "https://api.test.cl")

	// Preparar datos de prueba
	caf := &models.CAF{
		ID:               "123",
		TipoDocumento:    "33",
		RutEmisor:        "12345678-9",
		RangoInicial:     1,
		RangoFinal:       1000,
		FechaVencimiento: time.Now().AddDate(0, 6, 0),
		Activo:           true,
		Estado:           "ACTIVO",
		FolioActual:      1,
		UpdatedAt:        time.Now(),
	}

	// Configurar mocks
	cache.On("Get", mock.Anything, mock.Anything).Return(redis.NewStringCmd(context.Background()))
	cache.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(redis.NewStatusCmd(context.Background()))

	// Test caso exitoso
	result, err := service.GetCAFDisponible(context.Background(), "33", "12345678-9")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, caf.ID, result.ID)

	// Test caso error - CAF no encontrado
	cache.On("Get", mock.Anything, mock.Anything).Return(redis.NewStringCmd(context.Background()))
	_, err = service.GetCAFDisponible(context.Background(), "33", "99999999-9")
	assert.Error(t, err)
}

func TestValidarCAF(t *testing.T) {
	config := &config.SupabaseConfig{}
	cache := &MockRedisClient{}
	siiService := &MockSIIService{}
	service := NewCAFService(config, &cache.Client, siiService, "cert.pem", "key.pem", "certificacion", "https://api.test.cl")

	tests := []struct {
		name    string
		caf     *models.CAF
		wantErr bool
	}{
		{
			name: "CAF válido",
			caf: &models.CAF{
				ID:               "123",
				TipoDocumento:    "33",
				RutEmisor:        "12345678-9",
				RangoInicial:     1,
				RangoFinal:       1000,
				FechaVencimiento: time.Now().AddDate(0, 6, 0),
				Activo:           true,
				Estado:           "ACTIVO",
				FolioActual:      1,
				UpdatedAt:        time.Now(),
			},
			wantErr: false,
		},
		{
			name: "CAF no activo",
			caf: &models.CAF{
				ID:               "123",
				TipoDocumento:    "33",
				RutEmisor:        "12345678-9",
				RangoInicial:     1,
				RangoFinal:       1000,
				FechaVencimiento: time.Now().AddDate(0, 6, 0),
				Activo:           false,
				Estado:           "INACTIVO",
				FolioActual:      1,
				UpdatedAt:        time.Now(),
			},
			wantErr: true,
		},
		{
			name: "CAF vencido",
			caf: &models.CAF{
				ID:               "123",
				TipoDocumento:    "33",
				RutEmisor:        "12345678-9",
				RangoInicial:     1,
				RangoFinal:       1000,
				FechaVencimiento: time.Now().AddDate(0, -1, 0), // Vencido hace 1 mes
				Activo:           true,
				Estado:           "ACTIVO",
				FolioActual:      1,
				UpdatedAt:        time.Now(),
			},
			wantErr: true,
		},
		{
			name: "CAF agotado",
			caf: &models.CAF{
				ID:               "123",
				TipoDocumento:    "33",
				RutEmisor:        "12345678-9",
				RangoInicial:     1,
				RangoFinal:       1000,
				FechaVencimiento: time.Now().AddDate(0, 6, 0),
				Activo:           true,
				Estado:           "ACTIVO",
				FolioActual:      1001, // Mayor que RangoFinal
				UpdatedAt:        time.Now(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.ValidarCAF(tt.caf)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSolicitarCAF(t *testing.T) {
	config := &config.SupabaseConfig{}
	cache := &MockRedisClient{}
	siiService := &MockSIIService{}
	service := NewCAFService(config, &cache.Client, siiService, "cert.pem", "key.pem", "certificacion", "https://api.test.cl")

	tests := []struct {
		name    string
		request *CAFRequest
		wantErr bool
	}{
		{
			name: "Solicitud válida",
			request: &CAFRequest{
				RUTEmisor:      "12345678-9",
				TipoDTE:        "33",
				FolioInicial:   1,
				FolioFinal:     1000,
				FechaSolicitud: time.Now(),
			},
			wantErr: false,
		},
		{
			name: "RUT emisor inválido",
			request: &CAFRequest{
				RUTEmisor:      "",
				TipoDTE:        "33",
				FolioInicial:   1,
				FolioFinal:     1000,
				FechaSolicitud: time.Now(),
			},
			wantErr: true,
		},
		{
			name: "Tipo DTE inválido",
			request: &CAFRequest{
				RUTEmisor:      "12345678-9",
				TipoDTE:        "",
				FolioInicial:   1,
				FolioFinal:     1000,
				FechaSolicitud: time.Now(),
			},
			wantErr: true,
		},
		{
			name: "Rango de folios inválido",
			request: &CAFRequest{
				RUTEmisor:      "12345678-9",
				TipoDTE:        "33",
				FolioInicial:   1000,
				FolioFinal:     1,
				FechaSolicitud: time.Now(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Configurar mocks
			cache.On("Get", mock.Anything, mock.Anything).Return(redis.NewStringCmd(context.Background()))
			cache.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(redis.NewStatusCmd(context.Background()))

			result, err := service.SolicitarCAF(context.Background(), tt.request)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

func TestConsultarEstadoCAF(t *testing.T) {
	config := &config.SupabaseConfig{}
	cache := &MockRedisClient{}
	siiService := &MockSIIService{}
	service := NewCAFService(config, &cache.Client, siiService, "cert.pem", "key.pem", "certificacion", "https://api.test.cl")

	tests := []struct {
		name    string
		trackID string
		estado  *models.EstadoSII
		err     error
		wantErr bool
	}{
		{
			name:    "Estado aceptado",
			trackID: "123",
			estado: &models.EstadoSII{
				Estado:  "ACEPTADO",
				Glosa:   "Solicitud procesada correctamente",
				TrackID: "123",
				Fecha:   time.Now(),
			},
			err:     nil,
			wantErr: false,
		},
		{
			name:    "Estado rechazado",
			trackID: "456",
			estado: &models.EstadoSII{
				Estado:  "RECHAZADO",
				Glosa:   "Error en la solicitud",
				TrackID: "456",
				Fecha:   time.Now(),
			},
			err:     nil,
			wantErr: false,
		},
		{
			name:    "Error en consulta",
			trackID: "789",
			estado:  nil,
			err:     fmt.Errorf("error de conexión"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Configurar mocks
			siiService.On("ConsultarEstado", mock.Anything).Return(tt.estado, tt.err)
			cache.On("Get", mock.Anything, mock.Anything).Return(redis.NewStringCmd(context.Background()))
			cache.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(redis.NewStatusCmd(context.Background()))

			result, err := service.ConsultarEstadoCAF(context.Background(), tt.trackID)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.estado.Estado, result.Estado)
				assert.Equal(t, tt.estado.TrackID, result.TrackID)
			}
		})
	}
}

func TestProgramarSolicitudCAF(t *testing.T) {
	config := &config.SupabaseConfig{}
	cache := &MockRedisClient{}
	siiService := &MockSIIService{}
	service := NewCAFService(config, &cache.Client, siiService, "cert.pem", "key.pem", "certificacion", "https://api.test.cl")

	tests := []struct {
		name              string
		rutEmisor         string
		tipoDTE           string
		umbralFolios      int
		foliosDisponibles int
		wantErr           bool
	}{
		{
			name:              "Solicitud programada exitosamente",
			rutEmisor:         "12345678-9",
			tipoDTE:           "33",
			umbralFolios:      100,
			foliosDisponibles: 50,
			wantErr:           false,
		},
		{
			name:              "Folios suficientes",
			rutEmisor:         "12345678-9",
			tipoDTE:           "33",
			umbralFolios:      100,
			foliosDisponibles: 150,
			wantErr:           false,
		},
		{
			name:              "Error al verificar folios",
			rutEmisor:         "12345678-9",
			tipoDTE:           "33",
			umbralFolios:      100,
			foliosDisponibles: -1,
			wantErr:           true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Configurar mocks
			cache.On("Get", mock.Anything, mock.Anything).Return(redis.NewStringCmd(context.Background()))
			cache.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(redis.NewStatusCmd(context.Background()))

			err := service.ProgramarSolicitudCAF(context.Background(), tt.rutEmisor, tt.tipoDTE, tt.umbralFolios)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDescargarCAF(t *testing.T) {
	config := &config.SupabaseConfig{}
	cache := &MockRedisClient{}
	siiService := &MockSIIService{}
	service := NewCAFService(config, &cache.Client, siiService, "cert.pem", "key.pem", "certificacion", "https://api.test.cl")

	tests := []struct {
		name      string
		url       string
		savePath  string
		rutEmisor string
		estado    *models.EstadoSII
		err       error
		wantErr   bool
	}{
		{
			name:      "Descarga exitosa",
			url:       "https://api.test.cl/caf/123",
			savePath:  "/tmp/caf.xml",
			rutEmisor: "12345678-9",
			estado: &models.EstadoSII{
				Estado:  "ACEPTADO",
				Glosa:   "Solicitud procesada correctamente",
				TrackID: "123",
				Fecha:   time.Now(),
			},
			err:     nil,
			wantErr: false,
		},
		{
			name:      "Error en descarga",
			url:       "https://api.test.cl/caf/456",
			savePath:  "/tmp/caf.xml",
			rutEmisor: "12345678-9",
			estado:    nil,
			err:       fmt.Errorf("error de conexión"),
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Configurar mocks
			siiService.On("ConsultarEstado", mock.Anything).Return(tt.estado, tt.err)
			cache.On("Get", mock.Anything, mock.Anything).Return(redis.NewStringCmd(context.Background()))
			cache.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(redis.NewStatusCmd(context.Background()))

			result, err := service.DescargarCAF(context.Background(), tt.url, tt.savePath, tt.rutEmisor)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.estado.Estado, result.Estado)
			}
		})
	}
}

func TestMonitorearEstadoCAF(t *testing.T) {
	config := &config.SupabaseConfig{}
	cache := &MockRedisClient{}
	siiService := &MockSIIService{}
	service := NewCAFService(config, &cache.Client, siiService, "cert.pem", "key.pem", "certificacion", "https://api.test.cl")

	tests := []struct {
		name    string
		trackID string
		estados []*models.EstadoSII
		errs    []error
		wantErr bool
	}{
		{
			name:    "Monitoreo exitoso",
			trackID: "123",
			estados: []*models.EstadoSII{
				{
					Estado:  "EN_PROCESO",
					Glosa:   "Procesando solicitud",
					TrackID: "123",
					Fecha:   time.Now(),
				},
				{
					Estado:  "ACEPTADO",
					Glosa:   "Solicitud procesada correctamente",
					TrackID: "123",
					Fecha:   time.Now(),
				},
			},
			errs:    []error{nil, nil},
			wantErr: false,
		},
		{
			name:    "Solicitud rechazada",
			trackID: "456",
			estados: []*models.EstadoSII{
				{
					Estado:  "EN_PROCESO",
					Glosa:   "Procesando solicitud",
					TrackID: "456",
					Fecha:   time.Now(),
				},
				{
					Estado:  "RECHAZADO",
					Glosa:   "Error en la solicitud",
					TrackID: "456",
					Fecha:   time.Now(),
				},
			},
			errs:    []error{nil, nil},
			wantErr: true,
		},
		{
			name:    "Error en consulta",
			trackID: "789",
			estados: []*models.EstadoSII{nil},
			errs:    []error{fmt.Errorf("error de conexión")},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Configurar mocks
			for i, estado := range tt.estados {
				siiService.On("ConsultarEstado", mock.Anything).Return(estado, tt.errs[i]).Once()
			}

			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			result, err := service.MonitorearEstadoCAF(ctx, tt.trackID, 100*time.Millisecond)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, "ACEPTADO", result.Estado)
			}
		})
	}
}

func TestRedisErrors(t *testing.T) {
	config := &config.SupabaseConfig{}
	cache := &MockRedisClient{}
	siiService := &MockSIIService{}
	service := NewCAFService(config, &cache.Client, siiService, "cert.pem", "key.pem", "certificacion", "https://api.test.cl")

	tests := []struct {
		name    string
		redisOp string
		err     error
		wantErr bool
	}{
		{
			name:    "Error de conexión Redis",
			redisOp: "Get",
			err:     redis.ErrClosed,
			wantErr: true,
		},
		{
			name:    "Error de timeout Redis",
			redisOp: "Get",
			err:     context.DeadlineExceeded,
			wantErr: true,
		},
		{
			name:    "Error al guardar en Redis",
			redisOp: "Set",
			err:     redis.ErrClosed,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.redisOp {
			case "Get":
				cmd := redis.NewStringCmd(context.Background())
				cmd.SetErr(tt.err)
				cache.On("Get", mock.Anything, mock.Anything).Return(cmd)
			case "Set":
				cmd := redis.NewStatusCmd(context.Background())
				cmd.SetErr(tt.err)
				cache.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(cmd)
			}

			_, err := service.GetCAFDisponible(context.Background(), "33", "12345678-9")
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "error de Redis")
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestConcurrentCAFOperations(t *testing.T) {
	config := &config.SupabaseConfig{}
	cache := &MockRedisClient{}
	siiService := &MockSIIService{}
	service := NewCAFService(config, &cache.Client, siiService, "cert.pem", "key.pem", "certificacion", "https://api.test.cl")

	// Configurar mocks para operaciones exitosas
	cmd := redis.NewStringCmd(context.Background())
	cmd.SetVal(`{"id":"123","tipoDocumento":"33","rutEmisor":"12345678-9","rangoInicial":1,"rangoFinal":1000,"fechaVencimiento":"2024-12-31T23:59:59Z","activo":true,"estado":"ACTIVO","folioActual":1}`)
	cache.On("Get", mock.Anything, mock.Anything).Return(cmd)
	cache.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(redis.NewStatusCmd(context.Background()))

	// Realizar operaciones concurrentes
	numOperations := 10
	results := make(chan error, numOperations)
	var wg sync.WaitGroup

	for i := 0; i < numOperations; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := service.GetCAFDisponible(context.Background(), "33", "12345678-9")
			results <- err
		}()
	}

	// Esperar a que todas las operaciones terminen
	wg.Wait()
	close(results)

	// Verificar resultados
	for err := range results {
		assert.NoError(t, err)
	}

	// Verificar que el mock fue llamado el número correcto de veces
	cache.AssertNumberOfCalls(t, "Get", numOperations)
}

func TestCAFValidationEdgeCases(t *testing.T) {
	config := &config.SupabaseConfig{}
	cache := &MockRedisClient{}
	siiService := &MockSIIService{}
	service := NewCAFService(config, &cache.Client, siiService, "cert.pem", "key.pem", "certificacion", "https://api.test.cl")

	tests := []struct {
		name    string
		caf     *models.CAF
		wantErr bool
	}{
		{
			name: "CAF con rango inicial negativo",
			caf: &models.CAF{
				ID:               "123",
				TipoDocumento:    "33",
				RutEmisor:        "12345678-9",
				RangoInicial:     -1,
				RangoFinal:       1000,
				FechaVencimiento: time.Now().AddDate(0, 6, 0),
				Activo:           true,
				Estado:           "ACTIVO",
				FolioActual:      1,
				UpdatedAt:        time.Now(),
			},
			wantErr: true,
		},
		{
			name: "CAF con rango final menor que inicial",
			caf: &models.CAF{
				ID:               "123",
				TipoDocumento:    "33",
				RutEmisor:        "12345678-9",
				RangoInicial:     1000,
				RangoFinal:       1,
				FechaVencimiento: time.Now().AddDate(0, 6, 0),
				Activo:           true,
				Estado:           "ACTIVO",
				FolioActual:      1,
				UpdatedAt:        time.Now(),
			},
			wantErr: true,
		},
		{
			name: "CAF con folio actual fuera de rango",
			caf: &models.CAF{
				ID:               "123",
				TipoDocumento:    "33",
				RutEmisor:        "12345678-9",
				RangoInicial:     1,
				RangoFinal:       1000,
				FechaVencimiento: time.Now().AddDate(0, 6, 0),
				Activo:           true,
				Estado:           "ACTIVO",
				FolioActual:      0,
				UpdatedAt:        time.Now(),
			},
			wantErr: true,
		},
		{
			name: "CAF con fecha de vencimiento en el pasado",
			caf: &models.CAF{
				ID:               "123",
				TipoDocumento:    "33",
				RutEmisor:        "12345678-9",
				RangoInicial:     1,
				RangoFinal:       1000,
				FechaVencimiento: time.Now().AddDate(0, 0, -1),
				Activo:           true,
				Estado:           "ACTIVO",
				FolioActual:      1,
				UpdatedAt:        time.Now(),
			},
			wantErr: true,
		},
		{
			name: "CAF con estado inválido",
			caf: &models.CAF{
				ID:               "123",
				TipoDocumento:    "33",
				RutEmisor:        "12345678-9",
				RangoInicial:     1,
				RangoFinal:       1000,
				FechaVencimiento: time.Now().AddDate(0, 6, 0),
				Activo:           true,
				Estado:           "INVALIDO",
				FolioActual:      1,
				UpdatedAt:        time.Now(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.ValidarCAF(tt.caf)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
