package models

import (
	"time"
)

// Configuracion representa la configuración general de una empresa
type Configuracion struct {
	ID                 string                  `json:"id"`
	EmpresaID          string                  `json:"empresa_id"`
	RazonSocial        string                  `json:"razon_social"`
	RUT                string                  `json:"rut"`
	Direccion          string                  `json:"direccion"`
	Comuna             string                  `json:"comuna"`
	Ciudad             string                  `json:"ciudad"`
	Giro               string                  `json:"giro"`
	ActividadEconomica string                  `json:"actividad_economica"`
	Logo               string                  `json:"logo"`
	FormaPago          string                  `json:"forma_pago"`
	PlazoPago          int                     `json:"plazo_pago"`
	RetencionIVA       bool                    `json:"retencion_iva"`
	ConfigSII          ConfiguracionSIIEmpresa `json:"config_sii"`
	ConfigEmail        ConfiguracionEmail      `json:"config_email"`
	CreatedAt          time.Time               `json:"created_at"`
	UpdatedAt          time.Time               `json:"updated_at"`
}

// ConfiguracionEmail representa la configuración de correo electrónico
type ConfiguracionEmail struct {
	ID              string    `json:"id"`
	EmpresaID       string    `json:"empresa_id"`
	Servidor        string    `json:"servidor"`
	Puerto          int       `json:"puerto"`
	Usuario         string    `json:"usuario"`
	Password        string    `json:"password"`
	SSL             bool      `json:"ssl"`
	EmailRemitente  string    `json:"email_remitente"`
	NombreRemitente string    `json:"nombre_remitente"`
	TemplateFactura string    `json:"template_factura"`
	TemplateBoleta  string    `json:"template_boleta"`
	PlantillaHTML   string    `json:"plantilla_html"`
	EnvioAutomatico bool      `json:"envio_automatico"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// ConfiguracionSIIEmpresa representa la configuración específica para la integración con el SII
type ConfiguracionSIIEmpresa struct {
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

// Funciones auxiliares para crear nuevas instancias
func NewConfiguracion(empresaID, razonSocial, rut, direccion string) *Configuracion {
	return &Configuracion{
		ID:           GenerateID(),
		EmpresaID:    empresaID,
		RazonSocial:  razonSocial,
		RUT:          rut,
		Direccion:    direccion,
		PlazoPago:    30,
		RetencionIVA: false,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

func NewConfiguracionEmail(empresaID, servidor, usuario, password string) *ConfiguracionEmail {
	return &ConfiguracionEmail{
		ID:              GenerateID(),
		EmpresaID:       empresaID,
		Servidor:        servidor,
		Puerto:          587,
		Usuario:         usuario,
		Password:        password,
		SSL:             true,
		EmailRemitente:  usuario,
		EnvioAutomatico: true,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
}
