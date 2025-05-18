package models

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Config contiene la configuración para el cliente SII
type Config struct {
	// Ambiente especifica el ambiente de ejecución (certificación o producción)
	Ambiente Ambiente `json:"ambiente" validate:"required"`

	// CertPath es la ruta al archivo del certificado digital
	CertPath string `json:"cert_path" validate:"required,file"`

	// KeyPath es la ruta al archivo de la llave privada
	KeyPath string `json:"key_path" validate:"required,file"`

	// SchemaPath es la ruta al archivo de esquema XSD para validación
	SchemaPath string `json:"schema_path" validate:"omitempty,file"`

	// Timeout es el tiempo máximo de espera para las peticiones HTTP (en segundos)
	Timeout time.Duration `json:"timeout" validate:"required,min=1"`

	// RetryCount es el número máximo de reintentos para peticiones fallidas
	RetryCount int `json:"retry_count" validate:"min=0"`

	// RetryDelay es el tiempo de espera entre reintentos (en segundos)
	RetryDelay time.Duration `json:"retry_delay" validate:"required_with=RetryCount,min=1"`

	// BaseURL es la URL base del servicio (se configura automáticamente según el ambiente)
	BaseURL string `json:"-"`
}

// NewConfig crea una nueva configuración con valores por defecto
func NewConfig() *Config {
	return &Config{
		Ambiente:   Certificacion,
		Timeout:    30 * time.Second,
		RetryCount: 3,
		RetryDelay: 5 * time.Second,
		SchemaPath: "",
	}
}

// Validate valida la configuración y verifica la existencia de archivos
func (c *Config) Validate() error {
	if c == nil {
		return fmt.Errorf("la configuración no puede ser nil")
	}

	// Validar ambiente
	if !c.Ambiente.IsValid() {
		return fmt.Errorf("ambiente inválido: %s", c.Ambiente)
	}

	// Validar rutas de archivos
	if err := c.validateFilePath(c.CertPath, "certificado"); err != nil {
		return err
	}
	if err := c.validateFilePath(c.KeyPath, "llave privada"); err != nil {
		return err
	}
	if c.SchemaPath != "" {
		if err := c.validateFilePath(c.SchemaPath, "esquema XSD"); err != nil {
			return err
		}
	}

	// Validar timeouts y reintentos
	if c.Timeout < 1*time.Second {
		return fmt.Errorf("el timeout debe ser al menos 1 segundo")
	}
	if c.RetryCount < 0 {
		return fmt.Errorf("el número de reintentos no puede ser negativo")
	}
	if c.RetryCount > 0 && c.RetryDelay < 1*time.Second {
		return fmt.Errorf("el tiempo entre reintentos debe ser al menos 1 segundo")
	}

	return nil
}

// validateFilePath valida que una ruta de archivo exista y sea accesible
func (c *Config) validateFilePath(path, fileType string) error {
	if path == "" {
		return fmt.Errorf("la ruta del %s es requerida", fileType)
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("error obteniendo ruta absoluta del %s: %w", fileType, err)
	}

	info, err := os.Stat(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("el archivo %s no existe: %s", fileType, absPath)
		}
		return fmt.Errorf("error accediendo al %s: %w", fileType, err)
	}

	if info.IsDir() {
		return fmt.Errorf("la ruta del %s es un directorio: %s", fileType, absPath)
	}

	return nil
}

// SetupBaseURL configura la URL base según el ambiente
func (c *Config) SetupBaseURL() {
	if c == nil {
		return
	}

	switch c.Ambiente {
	case Certificacion:
		c.BaseURL = URLBaseCertificacion
	case Produccion:
		c.BaseURL = URLBaseProduccion
	default:
		c.BaseURL = URLBaseCertificacion
	}
}
