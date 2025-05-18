package siimodels

// Config contiene la configuración para el cliente SII
type Config struct {
	// Configuración general
	RutEmpresa       string `json:"rut_empresa"`
	RutCertificado   string `json:"rut_certificado"`
	ClaveCertificado string `json:"clave_certificado"`
	PathCertificado  string `json:"path_certificado"`

	// Configuración del servicio SII
	BaseURL    string   `json:"base_url"`
	CertPath   string   `json:"cert_path"`
	KeyPath    string   `json:"key_path"`
	RetryCount int      `json:"retry_count"`
	Timeout    int      `json:"timeout"`
	SchemaPath string   `json:"schema_path"` // Ruta al archivo XSD para validación
	Ambiente   Ambiente `json:"ambiente"`    // Ambiente de ejecución
}

// NewConfig crea una nueva instancia de Config con valores por defecto
func NewConfig() *Config {
	return &Config{
		BaseURL:    "https://palena.sii.cl",
		RetryCount: 3,
		Timeout:    30,
		Ambiente:   Certificacion,
	}
}
