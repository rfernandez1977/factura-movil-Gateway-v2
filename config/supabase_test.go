package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewSupabaseConfig(t *testing.T) {
	tests := []struct {
		name       string
		url        string
		anonKey    string
		serviceKey string
		jwtSecret  string
		ambiente   string
		baseURL    string
		token      string
		opts       []SupabaseOption
		wantErr    bool
	}{
		{
			name:       "Configuración válida",
			url:        "https://test.supabase.co",
			anonKey:    "test-anon-key",
			serviceKey: "test-service-key",
			jwtSecret:  "test-jwt-secret",
			ambiente:   "CERTIFICACION",
			baseURL:    "https://api.test.cl",
			token:      "test-token",
			wantErr:    false,
		},
		{
			name:       "URL inválida",
			url:        "",
			anonKey:    "test-anon-key",
			serviceKey: "test-service-key",
			jwtSecret:  "test-jwt-secret",
			ambiente:   "CERTIFICACION",
			baseURL:    "https://api.test.cl",
			token:      "test-token",
			wantErr:    true,
		},
		{
			name:       "AnonKey inválida",
			url:        "https://test.supabase.co",
			anonKey:    "",
			serviceKey: "test-service-key",
			jwtSecret:  "test-jwt-secret",
			ambiente:   "CERTIFICACION",
			baseURL:    "https://api.test.cl",
			token:      "test-token",
			wantErr:    true,
		},
		{
			name:       "ServiceKey inválida",
			url:        "https://test.supabase.co",
			anonKey:    "test-anon-key",
			serviceKey: "",
			jwtSecret:  "test-jwt-secret",
			ambiente:   "CERTIFICACION",
			baseURL:    "https://api.test.cl",
			token:      "test-token",
			wantErr:    true,
		},
		{
			name:       "JWTSecret inválido",
			url:        "https://test.supabase.co",
			anonKey:    "test-anon-key",
			serviceKey: "test-service-key",
			jwtSecret:  "",
			ambiente:   "CERTIFICACION",
			baseURL:    "https://api.test.cl",
			token:      "test-token",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := NewSupabaseConfig(tt.url, tt.anonKey, tt.serviceKey, tt.jwtSecret, tt.ambiente, tt.baseURL, tt.token, tt.opts...)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, config)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, config)
				assert.Equal(t, tt.url, config.URL)
				assert.Equal(t, tt.anonKey, config.AnonKey)
				assert.Equal(t, tt.serviceKey, config.ServiceKey)
				assert.Equal(t, tt.jwtSecret, config.JWTSecret)
				assert.Equal(t, tt.ambiente, config.Ambiente)
				assert.Equal(t, tt.baseURL, config.BaseURL)
				assert.Equal(t, tt.token, config.Token)
			}
		})
	}
}

func TestWithTimeout(t *testing.T) {
	config := &SupabaseConfig{
		URL:        "https://test.supabase.co",
		AnonKey:    "test-anon-key",
		ServiceKey: "test-service-key",
		JWTSecret:  "test-jwt-secret",
		Ambiente:   "CERTIFICACION",
		BaseURL:    "https://api.test.cl",
		Token:      "test-token",
	}

	tests := []struct {
		name      string
		timeout   time.Duration
		wantValue time.Duration
	}{
		{
			name:      "Timeout válido",
			timeout:   60 * time.Second,
			wantValue: 60 * time.Second,
		},
		{
			name:      "Timeout cero",
			timeout:   0,
			wantValue: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithTimeout(tt.timeout)
			opt(config)
			assert.Equal(t, tt.wantValue, config.Timeout)
		})
	}
}

func TestWithRetries(t *testing.T) {
	config := &SupabaseConfig{
		URL:        "https://test.supabase.co",
		AnonKey:    "test-anon-key",
		ServiceKey: "test-service-key",
		JWTSecret:  "test-jwt-secret",
		Ambiente:   "CERTIFICACION",
		BaseURL:    "https://api.test.cl",
		Token:      "test-token",
	}

	tests := []struct {
		name          string
		maxRetries    int
		retryInterval time.Duration
		wantRetries   int
		wantInterval  time.Duration
	}{
		{
			name:          "Valores válidos",
			maxRetries:    5,
			retryInterval: 2 * time.Second,
			wantRetries:   5,
			wantInterval:  2 * time.Second,
		},
		{
			name:          "Retries negativo",
			maxRetries:    -1,
			retryInterval: 2 * time.Second,
			wantRetries:   -1,
			wantInterval:  2 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithRetries(tt.maxRetries, tt.retryInterval)
			opt(config)
			assert.Equal(t, tt.wantRetries, config.MaxRetries)
			assert.Equal(t, tt.wantInterval, config.RetryInterval)
		})
	}
}

func TestWithSchema(t *testing.T) {
	config := &SupabaseConfig{
		URL:        "https://test.supabase.co",
		AnonKey:    "test-anon-key",
		ServiceKey: "test-service-key",
		JWTSecret:  "test-jwt-secret",
		Ambiente:   "CERTIFICACION",
		BaseURL:    "https://api.test.cl",
		Token:      "test-token",
	}

	tests := []struct {
		name      string
		schema    string
		wantValue string
	}{
		{
			name:      "Schema válido",
			schema:    "custom",
			wantValue: "custom",
		},
		{
			name:      "Schema vacío",
			schema:    "",
			wantValue: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithSchema(tt.schema)
			opt(config)
			assert.Equal(t, tt.wantValue, config.Schema)
		})
	}
}

func TestWithHeaders(t *testing.T) {
	config := &SupabaseConfig{
		URL:        "https://test.supabase.co",
		AnonKey:    "test-anon-key",
		ServiceKey: "test-service-key",
		JWTSecret:  "test-jwt-secret",
		Ambiente:   "CERTIFICACION",
		BaseURL:    "https://api.test.cl",
		Token:      "test-token",
	}

	tests := []struct {
		name    string
		headers map[string]string
		want    map[string]string
	}{
		{
			name: "Headers válidos",
			headers: map[string]string{
				"custom": "header",
			},
			want: map[string]string{
				"custom": "header",
			},
		},
		{
			name:    "Headers vacíos",
			headers: map[string]string{},
			want:    map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithHeaders(tt.headers)
			opt(config)
			assert.Equal(t, tt.want, config.Headers)
		})
	}
}

func TestWithLogging(t *testing.T) {
	config := &SupabaseConfig{
		URL:        "https://test.supabase.co",
		AnonKey:    "test-anon-key",
		ServiceKey: "test-service-key",
		JWTSecret:  "test-jwt-secret",
		Ambiente:   "CERTIFICACION",
		BaseURL:    "https://api.test.cl",
		Token:      "test-token",
	}

	tests := []struct {
		name       string
		enable     bool
		level      string
		wantEnable bool
		wantLevel  string
	}{
		{
			name:       "Logging habilitado",
			enable:     true,
			level:      "debug",
			wantEnable: true,
			wantLevel:  "debug",
		},
		{
			name:       "Logging deshabilitado",
			enable:     false,
			level:      "info",
			wantEnable: false,
			wantLevel:  "info",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithLogging(tt.enable, tt.level)
			opt(config)
			assert.Equal(t, tt.wantEnable, config.EnableLogging)
			assert.Equal(t, tt.wantLevel, config.LogLevel)
		})
	}
}

func TestWithConnectionPool(t *testing.T) {
	config := &SupabaseConfig{
		URL:        "https://test.supabase.co",
		AnonKey:    "test-anon-key",
		ServiceKey: "test-service-key",
		JWTSecret:  "test-jwt-secret",
		Ambiente:   "CERTIFICACION",
		BaseURL:    "https://api.test.cl",
		Token:      "test-token",
	}

	tests := []struct {
		name          string
		poolSize      int
		maxIdleConns  int
		maxOpenConns  int
		maxLifetime   time.Duration
		wantPoolSize  int
		wantIdleConns int
		wantOpenConns int
		wantLifetime  time.Duration
	}{
		{
			name:          "Valores válidos",
			poolSize:      20,
			maxIdleConns:  10,
			maxOpenConns:  30,
			maxLifetime:   2 * time.Hour,
			wantPoolSize:  20,
			wantIdleConns: 10,
			wantOpenConns: 30,
			wantLifetime:  2 * time.Hour,
		},
		{
			name:          "Valores mínimos",
			poolSize:      1,
			maxIdleConns:  1,
			maxOpenConns:  1,
			maxLifetime:   1 * time.Second,
			wantPoolSize:  1,
			wantIdleConns: 1,
			wantOpenConns: 1,
			wantLifetime:  1 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithConnectionPool(tt.poolSize, tt.maxIdleConns, tt.maxOpenConns, tt.maxLifetime)
			opt(config)
			assert.Equal(t, tt.wantPoolSize, config.PoolSize)
			assert.Equal(t, tt.wantIdleConns, config.MaxIdleConns)
			assert.Equal(t, tt.wantOpenConns, config.MaxOpenConns)
			assert.Equal(t, tt.wantLifetime, config.ConnMaxLifetime)
		})
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  *SupabaseConfig
		wantErr bool
	}{
		{
			name: "Configuración válida",
			config: &SupabaseConfig{
				URL:             "https://test.supabase.co",
				AnonKey:         "test-anon-key",
				ServiceKey:      "test-service-key",
				JWTSecret:       "test-jwt-secret",
				Timeout:         30 * time.Second,
				MaxRetries:      3,
				RetryInterval:   1 * time.Second,
				PoolSize:        10,
				MaxIdleConns:    5,
				MaxOpenConns:    20,
				ConnMaxLifetime: 1 * time.Hour,
			},
			wantErr: false,
		},
		{
			name: "URL inválida",
			config: &SupabaseConfig{
				URL:             "",
				AnonKey:         "test-anon-key",
				ServiceKey:      "test-service-key",
				JWTSecret:       "test-jwt-secret",
				Timeout:         30 * time.Second,
				MaxRetries:      3,
				RetryInterval:   1 * time.Second,
				PoolSize:        10,
				MaxIdleConns:    5,
				MaxOpenConns:    20,
				ConnMaxLifetime: 1 * time.Hour,
			},
			wantErr: true,
		},
		{
			name: "Timeout inválido",
			config: &SupabaseConfig{
				URL:             "https://test.supabase.co",
				AnonKey:         "test-anon-key",
				ServiceKey:      "test-service-key",
				JWTSecret:       "test-jwt-secret",
				Timeout:         0,
				MaxRetries:      3,
				RetryInterval:   1 * time.Second,
				PoolSize:        10,
				MaxIdleConns:    5,
				MaxOpenConns:    20,
				ConnMaxLifetime: 1 * time.Hour,
			},
			wantErr: true,
		},
		{
			name: "MaxRetries inválido",
			config: &SupabaseConfig{
				URL:             "https://test.supabase.co",
				AnonKey:         "test-anon-key",
				ServiceKey:      "test-service-key",
				JWTSecret:       "test-jwt-secret",
				Timeout:         30 * time.Second,
				MaxRetries:      -1,
				RetryInterval:   1 * time.Second,
				PoolSize:        10,
				MaxIdleConns:    5,
				MaxOpenConns:    20,
				ConnMaxLifetime: 1 * time.Hour,
			},
			wantErr: true,
		},
		{
			name: "RetryInterval inválido",
			config: &SupabaseConfig{
				URL:             "https://test.supabase.co",
				AnonKey:         "test-anon-key",
				ServiceKey:      "test-service-key",
				JWTSecret:       "test-jwt-secret",
				Timeout:         30 * time.Second,
				MaxRetries:      3,
				RetryInterval:   0,
				PoolSize:        10,
				MaxIdleConns:    5,
				MaxOpenConns:    20,
				ConnMaxLifetime: 1 * time.Hour,
			},
			wantErr: true,
		},
		{
			name: "PoolSize inválido",
			config: &SupabaseConfig{
				URL:             "https://test.supabase.co",
				AnonKey:         "test-anon-key",
				ServiceKey:      "test-service-key",
				JWTSecret:       "test-jwt-secret",
				Timeout:         30 * time.Second,
				MaxRetries:      3,
				RetryInterval:   1 * time.Second,
				PoolSize:        0,
				MaxIdleConns:    5,
				MaxOpenConns:    20,
				ConnMaxLifetime: 1 * time.Hour,
			},
			wantErr: true,
		},
		{
			name: "MaxIdleConns inválido",
			config: &SupabaseConfig{
				URL:             "https://test.supabase.co",
				AnonKey:         "test-anon-key",
				ServiceKey:      "test-service-key",
				JWTSecret:       "test-jwt-secret",
				Timeout:         30 * time.Second,
				MaxRetries:      3,
				RetryInterval:   1 * time.Second,
				PoolSize:        10,
				MaxIdleConns:    0,
				MaxOpenConns:    20,
				ConnMaxLifetime: 1 * time.Hour,
			},
			wantErr: true,
		},
		{
			name: "MaxOpenConns inválido",
			config: &SupabaseConfig{
				URL:             "https://test.supabase.co",
				AnonKey:         "test-anon-key",
				ServiceKey:      "test-service-key",
				JWTSecret:       "test-jwt-secret",
				Timeout:         30 * time.Second,
				MaxRetries:      3,
				RetryInterval:   1 * time.Second,
				PoolSize:        10,
				MaxIdleConns:    5,
				MaxOpenConns:    0,
				ConnMaxLifetime: 1 * time.Hour,
			},
			wantErr: true,
		},
		{
			name: "ConnMaxLifetime inválido",
			config: &SupabaseConfig{
				URL:             "https://test.supabase.co",
				AnonKey:         "test-anon-key",
				ServiceKey:      "test-service-key",
				JWTSecret:       "test-jwt-secret",
				Timeout:         30 * time.Second,
				MaxRetries:      3,
				RetryInterval:   1 * time.Second,
				PoolSize:        10,
				MaxIdleConns:    5,
				MaxOpenConns:    20,
				ConnMaxLifetime: 0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
