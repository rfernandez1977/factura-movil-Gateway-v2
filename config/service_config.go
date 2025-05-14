package config

// PDFConfig contiene la configuración para el servicio de PDF
type PDFConfig struct {
	TemplatePath    string
	TempPath        string
	DefaultTemplate string
	UseCompression  bool
	Dpi             int
	Orientation     string // portrait, landscape
	PaperSize       string // A4, Letter, etc.
}

// EmailConfig contiene la configuración para el servicio de email
type EmailConfig struct {
	SMTPServer   string
	SMTPPort     int
	SMTPUser     string
	SMTPPassword string
	FromEmail    string
	FromName     string
	ReplyTo      string
	UseTLS       bool
	UseSSL       bool
}

// SupabaseConfig contiene la configuración para el cliente de Supabase
type SupabaseConfig struct {
	URL             string
	APIKey          string
	ServiceKey      string
	AnonKey         string
	Timeout         int
	MaxRetries      int
	SchemaName      string
	TablaDocumentos string
	TablaEmpresas   string
	TablaUsuarios   string
	SIIBaseURL      string
	Ambiente        string
	JWTSecret       string
	pdfConfig       *PDFConfig   // Configuración de PDF
	emailConfig     *EmailConfig // Configuración de Email
}

// GetPDFConfig obtiene la configuración de PDF
func (c *SupabaseConfig) GetPDFConfig() *PDFConfig {
	if c.pdfConfig == nil {
		// Configuración por defecto si no existe
		c.pdfConfig = &PDFConfig{
			TemplatePath:    "templates/pdf",
			TempPath:        "temp/pdf",
			DefaultTemplate: "default.html",
			UseCompression:  true,
			Dpi:             300,
			Orientation:     "portrait",
			PaperSize:       "A4",
		}
	}
	return c.pdfConfig
}

// GetEmailConfig obtiene la configuración de Email
func (c *SupabaseConfig) GetEmailConfig() *EmailConfig {
	if c.emailConfig == nil {
		// Configuración por defecto si no existe
		c.emailConfig = &EmailConfig{
			SMTPServer:   "smtp.gmail.com",
			SMTPPort:     587,
			SMTPUser:     "user@example.com",
			SMTPPassword: "",
			FromEmail:    "sistema@ejemplo.com",
			FromName:     "Sistema de Facturación",
			ReplyTo:      "noreply@ejemplo.com",
			UseTLS:       true,
			UseSSL:       false,
		}
	}
	return c.emailConfig
}

// SetPDFConfig establece la configuración de PDF
func (c *SupabaseConfig) SetPDFConfig(config *PDFConfig) {
	c.pdfConfig = config
}

// SetEmailConfig establece la configuración de Email
func (c *SupabaseConfig) SetEmailConfig(config *EmailConfig) {
	c.emailConfig = config
}
