package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cursor/FMgo/models"
	"github.com/supabase-community/postgrest-go"
)

// Config representa la configuración general del sistema, reexportada desde models
type Config = models.Config

// NewConfig creates a new config instance
func NewConfig() *Config {
	config := &Config{}

	// Inicializar el cliente por defecto
	if config.Supabase.URL != "" && config.Supabase.APIKey != "" {
		client := postgrest.NewClient(config.Supabase.URL, "", map[string]string{
			"apikey":        config.Supabase.APIKey,
			"Authorization": "Bearer " + config.Supabase.ServiceKey,
		})
		config.Client = client
	}

	return config
}

// WithTimeout configures timeout settings
func WithTimeout(config *Config, timeout int) *Config {
	config.Server.ReadTimeout = timeout
	config.Server.WriteTimeout = timeout
	config.SII.Timeout = timeout
	config.Supabase.Timeout = timeout
	return config
}

// WithRetries configures retry settings
func WithRetries(config *Config, retries int) *Config {
	config.SII.RetryCount = retries
	config.Supabase.MaxRetries = retries
	return config
}

// WithSchema configures schema settings
func WithSchema(config *Config, schema string) *Config {
	// Add schema configuration if needed
	return config
}

// WithHeaders configures header settings
func WithHeaders(config *Config, headers map[string]string) *Config {
	// Add headers configuration if needed
	return config
}

// Cargar configuración desde un archivo o valores por defecto
func Load(path string) (*Config, error) {
	if path == "" {
		path = "config.json"
	}

	// Verificar que el archivo existe
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo ruta absoluta: %v", err)
	}

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("archivo de configuración no encontrado: %s", absPath)
	}

	// Leer archivo
	data, err := os.ReadFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("error leyendo archivo: %v", err)
	}

	// Parsear JSON
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("error parseando JSON: %v", err)
	}

	return &config, nil
}

// SaveConfig guarda la configuración en un archivo
func SaveConfig(config *Config, path string) error {
	if path == "" {
		path = "config.json"
	}

	// Convertir a JSON
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("error serializando a JSON: %v", err)
	}

	// Guardar archivo
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("error guardando archivo: %v", err)
	}

	return nil
}

// GetDefaultConfig retorna una configuración por defecto
func GetDefaultConfig() *Config {
	config := &Config{}

	// Configuración del servidor
	config.Server.Port = 8080
	config.Server.Host = "localhost"
	config.Server.ReadTimeout = 30
	config.Server.WriteTimeout = 30

	// Configuración de la base de datos
	config.Database.Host = "localhost"
	config.Database.Port = 5432
	config.Database.User = "postgres"
	config.Database.Password = "postgres"
	config.Database.Name = "fmgo"
	config.Database.SSLMode = "disable"

	// Configuración del SII
	config.SII.BaseURL = "https://maullin.sii.cl/DTEWS/"
	config.SII.Timeout = 30
	config.SII.RetryCount = 3
	config.SII.RetryDelay = 5
	config.SII.CertPath = "cert/cert.p12"
	config.SII.KeyPath = "cert/key.pem"
	config.SII.CertPassword = ""

	// Configuración de logs
	config.Logging.Level = "info"
	config.Logging.FilePath = "logs/fmgo.log"
	config.Logging.MaxSize = 10
	config.Logging.MaxBackups = 5
	config.Logging.MaxAge = 30
	config.Logging.Compress = true

	// Configuración de seguridad
	config.Security.JWTSecret = "default-secret-key-change-in-production"
	config.Security.JWTExpiration = 24
	config.Security.CORSEnabled = true
	config.Security.CORSOrigins = []string{"*"}

	// Configuración de Supabase
	config.Supabase.URL = "https://example.supabase.co"
	config.Supabase.APIKey = "your-api-key"
	config.Supabase.ServiceKey = "your-service-key"
	config.Supabase.AnonKey = "your-anon-key"
	config.Supabase.Timeout = 30
	config.Supabase.MaxRetries = 3
	config.Supabase.TablaDocumentos = "documentos"
	config.Supabase.TablaEmpresas = "empresas"
	config.Supabase.TablaUsuarios = "usuarios"

	return config
}

// GetDSN returns the database connection string
func GetDSN(c *Config) string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.Name,
		c.Database.SSLMode)
}

// GetEnv returns the current environment
func GetEnv(c *Config) string {
	// Default to development environment
	return "development"
}
