package models

// Config representa la configuración genérica para servicios
type Config struct {
	// Configuración básica
	Ambiente string `json:"ambiente" bson:"ambiente"`
	Debug    bool   `json:"debug" bson:"debug"`
	Timeout  int    `json:"timeout" bson:"timeout"`
	MaxRetry int    `json:"max_retry" bson:"max_retry"`
	LogLevel string `json:"log_level" bson:"log_level"`

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
	}
}
