package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/supabase-community/postgrest-go"
)

func TestNewConfig(t *testing.T) {
	// Crear un cliente mock
	mockClient := &postgrest.Client{}

	config := NewConfig(mockClient)
	assert.NotNil(t, config)
	assert.Equal(t, "development", config.Env)
	assert.NotNil(t, config.Staging)
	assert.NotNil(t, config.Prod)
	assert.NotNil(t, config.Dev)
	assert.NotNil(t, config.Test)
	assert.NotNil(t, config.Cert)
	assert.NotNil(t, config.Integration)
}

func TestEnvironmentConfig(t *testing.T) {
	// Crear un cliente mock
	mockClient := &postgrest.Client{}

	config := NewConfig(mockClient)
	env := config.Staging

	// Test IsDebug
	assert.True(t, env.IsDebug())
	env.SetDebug(false)
	assert.False(t, env.IsDebug())

	// Test IsTrace
	assert.True(t, env.IsTrace())
	env.SetTrace(false)
	assert.False(t, env.IsTrace())

	// Test IsProfile
	assert.True(t, env.IsProfile())
	env.SetProfile(false)
	assert.False(t, env.IsProfile())

	// Test IsHotReload
	assert.True(t, env.IsHotReload())
	env.SetHotReload(false)
	assert.False(t, env.IsHotReload())
}

func TestLoadConfig(t *testing.T) {
	config, err := Load()
	assert.NoError(t, err)
	assert.NotNil(t, config)

	// Validar campos requeridos
	assert.NotEmpty(t, config.Supabase.URL)
	assert.NotEmpty(t, config.Supabase.AnonKey)
	assert.NotEmpty(t, config.Supabase.ServiceKey)
	assert.NotEmpty(t, config.Supabase.JWTSecret)
}

func TestConfigValidation(t *testing.T) {
	// Crear un cliente mock
	mockClient := &postgrest.Client{}

	config := NewConfig(mockClient)

	// Test validación de timeout
	config.Supabase.Timeout = -1
	err := config.Supabase.Validate()
	assert.Error(t, err)

	// Test validación de MaxRetries
	config.Supabase.MaxRetries = -1
	err = config.Supabase.Validate()
	assert.Error(t, err)

	// Test validación de RetryInterval
	config.Supabase.RetryInterval = -1
	err = config.Supabase.Validate()
	assert.Error(t, err)
}

func TestConfigOptions(t *testing.T) {
	// Crear un cliente mock
	mockClient := &postgrest.Client{}

	config := NewConfig(mockClient)

	// Test WithTimeout
	timeout := 60 * time.Second
	WithTimeout(timeout)(&config.Supabase)
	assert.Equal(t, timeout, config.Supabase.Timeout)

	// Test WithRetries
	maxRetries := 5
	retryInterval := 2 * time.Second
	WithRetries(maxRetries, retryInterval)(&config.Supabase)
	assert.Equal(t, maxRetries, config.Supabase.MaxRetries)
	assert.Equal(t, retryInterval, config.Supabase.RetryInterval)

	// Test WithSchema
	schema := "custom"
	WithSchema(schema)(&config.Supabase)
	assert.Equal(t, schema, config.Supabase.Schema)

	// Test WithHeaders
	headers := map[string]string{"custom": "header"}
	WithHeaders(headers)(&config.Supabase)
	assert.Equal(t, headers, config.Supabase.Headers)
}

func TestStagingConfig(t *testing.T) {
	mockClient := &postgrest.Client{}
	config := NewConfig(mockClient)
	assert.NotNil(t, config.Staging)
	assert.True(t, config.Staging.IsDebug())
	assert.True(t, config.Staging.IsTrace())
	assert.True(t, config.Staging.IsProfile())
	assert.True(t, config.Staging.IsHotReload())

	// Probar cambios de configuración
	config.Staging.SetDebug(false)
	assert.False(t, config.Staging.IsDebug())

	config.Staging.SetTrace(false)
	assert.False(t, config.Staging.IsTrace())

	config.Staging.SetProfile(false)
	assert.False(t, config.Staging.IsProfile())

	config.Staging.SetHotReload(false)
	assert.False(t, config.Staging.IsHotReload())
}

func TestProdConfig(t *testing.T) {
	mockClient := &postgrest.Client{}
	config := NewConfig(mockClient)
	assert.NotNil(t, config.Prod)
	assert.False(t, config.Prod.IsDebug())
	assert.False(t, config.Prod.IsTrace())
	assert.False(t, config.Prod.IsProfile())
	assert.False(t, config.Prod.IsHotReload())

	// Probar cambios de configuración
	config.Prod.SetDebug(true)
	assert.True(t, config.Prod.IsDebug())

	config.Prod.SetTrace(true)
	assert.True(t, config.Prod.IsTrace())

	config.Prod.SetProfile(true)
	assert.True(t, config.Prod.IsProfile())

	config.Prod.SetHotReload(true)
	assert.True(t, config.Prod.IsHotReload())
}

func TestDevConfig(t *testing.T) {
	mockClient := &postgrest.Client{}
	config := NewConfig(mockClient)
	assert.NotNil(t, config.Dev)
	assert.True(t, config.Dev.IsDebug())
	assert.True(t, config.Dev.IsTrace())
	assert.True(t, config.Dev.IsProfile())
	assert.True(t, config.Dev.IsHotReload())

	// Probar cambios de configuración
	config.Dev.SetDebug(false)
	assert.False(t, config.Dev.IsDebug())

	config.Dev.SetTrace(false)
	assert.False(t, config.Dev.IsTrace())

	config.Dev.SetProfile(false)
	assert.False(t, config.Dev.IsProfile())

	config.Dev.SetHotReload(false)
	assert.False(t, config.Dev.IsHotReload())
}

func TestTestConfig(t *testing.T) {
	mockClient := &postgrest.Client{}
	config := NewConfig(mockClient)
	assert.NotNil(t, config.Test)
	assert.True(t, config.Test.IsDebug())
	assert.True(t, config.Test.IsTrace())
	assert.True(t, config.Test.IsProfile())
	assert.True(t, config.Test.IsHotReload())

	// Probar cambios de configuración
	config.Test.SetDebug(false)
	assert.False(t, config.Test.IsDebug())

	config.Test.SetTrace(false)
	assert.False(t, config.Test.IsTrace())

	config.Test.SetProfile(false)
	assert.False(t, config.Test.IsProfile())

	config.Test.SetHotReload(false)
	assert.False(t, config.Test.IsHotReload())
}

func TestCertConfig(t *testing.T) {
	mockClient := &postgrest.Client{}
	config := NewConfig(mockClient)
	assert.NotNil(t, config.Cert)
	assert.True(t, config.Cert.IsDebug())
	assert.True(t, config.Cert.IsTrace())
	assert.True(t, config.Cert.IsProfile())
	assert.True(t, config.Cert.IsHotReload())

	// Probar cambios de configuración
	config.Cert.SetDebug(false)
	assert.False(t, config.Cert.IsDebug())

	config.Cert.SetTrace(false)
	assert.False(t, config.Cert.IsTrace())

	config.Cert.SetProfile(false)
	assert.False(t, config.Cert.IsProfile())

	config.Cert.SetHotReload(false)
	assert.False(t, config.Cert.IsHotReload())
}

func TestIntegrationConfig(t *testing.T) {
	mockClient := &postgrest.Client{}
	config := NewConfig(mockClient)
	assert.NotNil(t, config.Integration)
	assert.True(t, config.Integration.IsDebug())
	assert.True(t, config.Integration.IsTrace())
	assert.True(t, config.Integration.IsProfile())
	assert.True(t, config.Integration.IsHotReload())

	// Probar cambios de configuración
	config.Integration.SetDebug(false)
	assert.False(t, config.Integration.IsDebug())

	config.Integration.SetTrace(false)
	assert.False(t, config.Integration.IsTrace())

	config.Integration.SetProfile(false)
	assert.False(t, config.Integration.IsProfile())

	config.Integration.SetHotReload(false)
	assert.False(t, config.Integration.IsHotReload())
}
