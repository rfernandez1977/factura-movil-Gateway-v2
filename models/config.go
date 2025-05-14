package models

import "github.com/supabase-community/postgrest-go"

// Config representa la configuración genérica para servicios
type Config struct {
	// Configuración básica
	Ambiente string `json:"ambiente" bson:"ambiente"`
	Debug    bool   `json:"debug" bson:"debug"`
	Timeout  int    `json:"timeout" bson:"timeout"`
	MaxRetry int    `json:"max_retry" bson:"max_retry"`
	LogLevel string `json:"log_level" bson:"log_level"`
	Env      string `json:"env" bson:"env"`

	// Configuración de almacenamiento
	StoragePath string `json:"storage_path,omitempty" bson:"storage_path,omitempty"`

	// Configuración de base de datos
	DatabaseURL      string `json:"database_url,omitempty" bson:"database_url,omitempty"`
	DatabaseName     string `json:"database_name,omitempty" bson:"database_name,omitempty"`
	DatabaseUsername string `json:"database_username,omitempty" bson:"database_username,omitempty"`
	DatabasePassword string `json:"database_password,omitempty" bson:"database_password,omitempty"`

	// Configuración de API
	APIPort     int    `json:"api_port,omitempty" bson:"api_port,omitempty"`
	APIHost     string `json:"api_host,omitempty" bson:"api_host,omitempty"`
	APIBasePath string `json:"api_base_path,omitempty" bson:"api_base_path,omitempty"`

	// Configuración de SII
	SIIEndpoint   string `json:"sii_endpoint,omitempty" bson:"sii_endpoint,omitempty"`
	SIICertPath   string `json:"sii_cert_path,omitempty" bson:"sii_cert_path,omitempty"`
	SIIKeyPath    string `json:"sii_key_path,omitempty" bson:"sii_key_path,omitempty"`
	SIIPassphrase string `json:"sii_passphrase,omitempty" bson:"sii_passphrase,omitempty"`

	// Configuración adicional como mapa
	AdditionalConfig map[string]interface{} `json:"additional_config,omitempty" bson:"additional_config,omitempty"`

	// Configuraciones específicas
	Server struct {
		Port         int    `json:"port"`
		Host         string `json:"host"`
		ReadTimeout  int    `json:"read_timeout"`
		WriteTimeout int    `json:"write_timeout"`
	} `json:"server"`

	Database struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
		Name     string `json:"name"`
		SSLMode  string `json:"ssl_mode"`
	} `json:"database"`

	SII struct {
		BaseURL      string `json:"base_url"`
		Timeout      int    `json:"timeout"`
		RetryCount   int    `json:"retry_count"`
		RetryDelay   int    `json:"retry_delay"`
		CertPath     string `json:"cert_path"`
		KeyPath      string `json:"key_path"`
		CertPassword string `json:"cert_password"`
	} `json:"sii"`

	Logging struct {
		Level      string `json:"level"`
		FilePath   string `json:"file_path"`
		MaxSize    int    `json:"max_size"`
		MaxBackups int    `json:"max_backups"`
		MaxAge     int    `json:"max_age"`
		Compress   bool   `json:"compress"`
	} `json:"logging"`

	Security struct {
		JWTSecret     string   `json:"jwt_secret"`
		JWTExpiration int      `json:"jwt_expiration"`
		CORSEnabled   bool     `json:"cors_enabled"`
		CORSOrigins   []string `json:"cors_origins"`
	} `json:"security"`

	Supabase struct {
		URL             string `json:"url"`
		APIKey          string `json:"api_key"`
		ServiceKey      string `json:"service_key"`
		AnonKey         string `json:"anon_key"`
		Timeout         int    `json:"timeout"`
		MaxRetries      int    `json:"max_retries"`
		TablaDocumentos string `json:"tabla_documentos"`
		TablaEmpresas   string `json:"tabla_empresas"`
		TablaUsuarios   string `json:"tabla_usuarios"`
	} `json:"supabase"`

	// Cliente externo
	Client *postgrest.Client `json:"-"`
}

// NewConfig crea una nueva configuración con valores predeterminados
func NewConfig() *Config {
	return &Config{
		Ambiente:         "desarrollo",
		Debug:            true,
		Timeout:          30,
		MaxRetry:         3,
		LogLevel:         "info",
		APIPort:          8080,
		APIHost:          "localhost",
		APIBasePath:      "/api/v1",
		AdditionalConfig: make(map[string]interface{}),
		Env:              "development",
	}
}
