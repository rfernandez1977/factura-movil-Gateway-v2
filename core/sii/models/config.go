package models

// Config representa la configuración del cliente SII
type Config struct {
	Ambiente         Ambiente
	RutEmpresa       string
	RutCertificado   string
	ClaveCertificado string
	PathCertificado  string
}

// NewConfig crea una nueva instancia de Config
func NewConfig(ambiente Ambiente, rutEmpresa, rutCert, claveCert, pathCert string) *Config {
	return &Config{
		Ambiente:         ambiente,
		RutEmpresa:       rutEmpresa,
		RutCertificado:   rutCert,
		ClaveCertificado: claveCert,
		PathCertificado:  pathCert,
	}
}

// SIIConfig representa la configuración específica para el SII
type SIIConfig struct {
	BaseURL    string   `json:"base_url"`
	CertPath   string   `json:"cert_path"`
	KeyPath    string   `json:"key_path"`
	RetryCount int      `json:"retry_count"`
	Timeout    int      `json:"timeout"`
	Ambiente   Ambiente `json:"ambiente"`
}

// Ambiente representa el ambiente de ejecución (Producción o Certificación)
type Ambiente string

const (
	Produccion    Ambiente = "PRODUCCION"
	Certificacion Ambiente = "CERTIFICACION"
)
