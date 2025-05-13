package models

import (
	"time"
)

// Config representa la configuración general del sistema
type Config struct {
	// Configuración del servidor
	Server struct {
		Port         int    `json:"port"`
		Host         string `json:"host"`
		ReadTimeout  int    `json:"readTimeout"`
		WriteTimeout int    `json:"writeTimeout"`
	} `json:"server"`

	// Configuración de la base de datos
	Database struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
		Name     string `json:"name"`
		SSLMode  string `json:"sslMode"`
	} `json:"database"`

	// Configuración del SII
	SII struct {
		BaseURL      string `json:"baseUrl"`
		Timeout      int    `json:"timeout"`
		RetryCount   int    `json:"retryCount"`
		RetryDelay   int    `json:"retryDelay"`
		CertPath     string `json:"certPath"`
		KeyPath      string `json:"keyPath"`
		CertPassword string `json:"certPassword"`
	} `json:"sii"`

	// Configuración de logs
	Logging struct {
		Level      string `json:"level"`
		FilePath   string `json:"filePath"`
		MaxSize    int    `json:"maxSize"`
		MaxBackups int    `json:"maxBackups"`
		MaxAge     int    `json:"maxAge"`
		Compress   bool   `json:"compress"`
	} `json:"logging"`

	// Configuración de caché
	Cache struct {
		Enabled  bool   `json:"enabled"`
		Type     string `json:"type"`
		Address  string `json:"address"`
		Password string `json:"password"`
		DB       int    `json:"db"`
	} `json:"cache"`

	// Configuración de seguridad
	Security struct {
		JWTSecret     string   `json:"jwtSecret"`
		JWTExpiration int      `json:"jwtExpiration"`
		CORSEnabled   bool     `json:"corsEnabled"`
		CORSOrigins   []string `json:"corsOrigins"`
	} `json:"security"`

	// Configuración de notificaciones
	Notifications struct {
		Email struct {
			Enabled  bool   `json:"enabled"`
			Host     string `json:"host"`
			Port     int    `json:"port"`
			Username string `json:"username"`
			Password string `json:"password"`
			From     string `json:"from"`
		} `json:"email"`
		SMS struct {
			Enabled  bool   `json:"enabled"`
			Provider string `json:"provider"`
			APIKey   string `json:"apiKey"`
			From     string `json:"from"`
		} `json:"sms"`
	} `json:"notifications"`

	// Configuración de Supabase
	Supabase struct {
		URL             string `json:"url"`
		APIKey          string `json:"apiKey"`
		ServiceKey      string `json:"serviceKey"`
		AnonKey         string `json:"anonKey"`
		Timeout         int    `json:"timeout"`
		MaxRetries      int    `json:"maxRetries"`
		TablaDocumentos string `json:"tablaDocumentos"`
		TablaEmpresas   string `json:"tablaEmpresas"`
		TablaUsuarios   string `json:"tablaUsuarios"`
	} `json:"supabase"`

	// Configuración de almacenamiento
	Storage struct {
		Type  string `json:"type"`
		Local struct {
			BasePath string `json:"basePath"`
		} `json:"local"`
		S3 struct {
			Bucket    string `json:"bucket"`
			Region    string `json:"region"`
			AccessKey string `json:"accessKey"`
			SecretKey string `json:"secretKey"`
		} `json:"s3"`
	} `json:"storage"`

	// Configuración de integración con ERP
	ERP struct {
		Enabled    bool   `json:"enabled"`
		Type       string `json:"type"`
		BaseURL    string `json:"baseUrl"`
		APIKey     string `json:"apiKey"`
		Timeout    int    `json:"timeout"`
		RetryCount int    `json:"retryCount"`
	} `json:"erp"`

	// Configuración de validación
	Validation struct {
		Enabled     bool     `json:"enabled"`
		Rules       []string `json:"rules"`
		MaxErrors   int      `json:"maxErrors"`
		StopOnError bool     `json:"stopOnError"`
	} `json:"validation"`

	// Configuración de monitoreo
	Monitoring struct {
		Enabled     bool   `json:"enabled"`
		Type        string `json:"type"`
		Endpoint    string `json:"endpoint"`
		APIKey      string `json:"apiKey"`
		Environment string `json:"environment"`
	} `json:"monitoring"`

	Client interface{} `json:"-" bson:"-"`
}

// ConfiguracionSII representa la configuración específica para la integración con el SII
type ConfiguracionSII struct {
	ID              string    `json:"id"`
	EmpresaID       string    `json:"empresa_id"`
	RUTEmisor       string    `json:"rutEmisor"`
	CertificadoPath string    `json:"certificadoPath"`
	ClavePrivada    string    `json:"clavePrivada"`
	Ambiente        string    `json:"ambiente"`
	Timeout         int       `json:"timeout"`
	RetryCount      int       `json:"retryCount"`
	RetryDelay      int       `json:"retryDelay"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

// ConfiguracionERP representa la configuración específica para la integración con el ERP
type ConfiguracionERP struct {
	ID         string    `json:"id"`
	TipoERP    TipoERP   `json:"tipo"`
	BaseURL    string    `json:"baseUrl"`
	APIKey     string    `json:"apiKey"`
	Timeout    int       `json:"timeout"`
	RetryCount int       `json:"retryCount"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

// ConfiguracionValidacion representa la configuración específica para la validación de documentos
type ConfiguracionValidacion struct {
	ID          string    `json:"id"`
	Tipo        string    `json:"tipo"`
	Reglas      []string  `json:"reglas"`
	MaxErrores  int       `json:"maxErrores"`
	StopOnError bool      `json:"stopOnError"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// ConfiguracionNotificacion representa la configuración específica para las notificaciones
type ConfiguracionNotificacion struct {
	ID           string    `json:"id"`
	Tipo         string    `json:"tipo"`
	Destinatario string    `json:"destinatario"`
	Template     string    `json:"template"`
	Activo       bool      `json:"activo"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// Funciones auxiliares para crear nuevas instancias de configuración
func NewConfig() *Config {
	return &Config{}
}

func NewConfiguracionSII(rutEmisor, certificadoPath, clavePrivada, ambiente string) *ConfiguracionSII {
	return &ConfiguracionSII{
		ID:              GenerateID(),
		EmpresaID:       "",
		RUTEmisor:       rutEmisor,
		CertificadoPath: certificadoPath,
		ClavePrivada:    clavePrivada,
		Ambiente:        ambiente,
		Timeout:         30,
		RetryCount:      3,
		RetryDelay:      5,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
}

func NewConfiguracionERP(tipo TipoERP, baseURL, apiKey string) *ConfiguracionERP {
	return &ConfiguracionERP{
		ID:         GenerateID(),
		TipoERP:    tipo,
		BaseURL:    baseURL,
		APIKey:     apiKey,
		Timeout:    30,
		RetryCount: 3,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}

func NewConfiguracionValidacion(tipo string, reglas []string) *ConfiguracionValidacion {
	return &ConfiguracionValidacion{
		ID:          GenerateID(),
		Tipo:        tipo,
		Reglas:      reglas,
		MaxErrores:  10,
		StopOnError: false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func NewConfiguracionNotificacion(tipo, destinatario, template string) *ConfiguracionNotificacion {
	return &ConfiguracionNotificacion{
		ID:           GenerateID(),
		Tipo:         tipo,
		Destinatario: destinatario,
		Template:     template,
		Activo:       true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}
