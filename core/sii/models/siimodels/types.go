package siimodels

import (
	"encoding/xml"
	"fmt"
	"time"
)

// AmbienteSII representa el ambiente de ejecución del SII
type AmbienteSII string

const (
	// AmbienteCertificacion representa el ambiente de pruebas
	AmbienteCertificacion AmbienteSII = "certificacion"
	// AmbienteProduccion representa el ambiente de producción
	AmbienteProduccion AmbienteSII = "produccion"
)

// URLs base para cada ambiente
const (
	URLBaseCertificacion = "https://maullin.sii.cl"
	URLBaseProduccion    = "https://palena.sii.cl"
)

// Endpoints específicos para certificación
const (
	EndpointSemillaCert  = "/DTEWS/CrSeed.jws"
	EndpointTokenCert    = "/DTEWS/GetTokenFromSeed.jws"
	EndpointEnvioCert    = "/cgi_dte/UPL/DTEUpload"
	EndpointConsultaCert = "/DTEWS/QueryEstDte.jws"
)

// ConfigSII contiene la configuración para el cliente SII
type ConfigSII struct {
	// Configuración de ambiente
	Ambiente AmbienteSII `json:"ambiente"`
	BaseURL  string      `json:"base_url"`

	// Configuración de certificados
	CertPath       string `json:"cert_path"`
	KeyPath        string `json:"key_path"`
	RutEmpresa     string `json:"rut_empresa"`     // RUT de la empresa
	RutCertificado string `json:"rut_certificado"` // RUT asociado al certificado digital

	// Configuración de conexión
	RetryCount int           `json:"retry_count"`
	RetryDelay time.Duration `json:"retry_delay"`
	Timeout    time.Duration `json:"timeout"`

	// Configuración de esquemas
	SchemaPath string `json:"schema_path"`
}

// NewConfigSII crea una nueva configuración con valores por defecto para certificación
func NewConfigSII() *ConfigSII {
	return &ConfigSII{
		Ambiente:   AmbienteCertificacion,
		Timeout:    30 * time.Second,
		RetryCount: 3,
		RetryDelay: 5 * time.Second,
		BaseURL:    URLBaseCertificacion,
	}
}

// Validate valida la configuración
func (c *ConfigSII) Validate() error {
	if c.CertPath == "" {
		return fmt.Errorf("la ruta del certificado es requerida")
	}
	if c.KeyPath == "" {
		return fmt.Errorf("la ruta de la llave privada es requerida")
	}
	if c.Timeout < 1*time.Second {
		return fmt.Errorf("el timeout debe ser al menos 1 segundo")
	}
	if c.RetryCount < 0 {
		return fmt.Errorf("el número de reintentos no puede ser negativo")
	}
	return nil
}

// SoapEnvelope es la estructura para mensajes SOAP
type SoapEnvelope struct {
	XMLName xml.Name `xml:"Envelope"`
	XMLNS   string   `xml:"xmlns,attr"`
	Body    SoapBody `xml:"Body"`
}

// SoapBody es el cuerpo del mensaje SOAP
type SoapBody struct {
	XMLName xml.Name    `xml:"Body"`
	Content interface{} `xml:",any"`
	Fault   *SoapFault  `xml:"Fault,omitempty"`
}

// SoapFault representa un error en el mensaje SOAP
type SoapFault struct {
	XMLName     xml.Name `xml:"Fault"`
	FaultCode   string   `xml:"faultcode"`
	FaultString string   `xml:"faultstring"`
}

// RespuestaSII es la estructura base para las respuestas del SII
type RespuestaSII struct {
	XMLName xml.Name        `xml:"RESPUESTA"`
	Header  RespuestaHeader `xml:"RESP_HDR"`
	Body    RespuestaBody   `xml:"RESP_BODY"`
}

// RespuestaHeader contiene la información de estado de la respuesta
type RespuestaHeader struct {
	Estado string `xml:"ESTADO"`
	Glosa  string `xml:"GLOSA"`
}

// RespuestaBody contiene el cuerpo de la respuesta
type RespuestaBody struct {
	Semilla string `xml:"SEMILLA,omitempty"`
	Token   string `xml:"TOKEN,omitempty"`
	TrackID string `xml:"TRACKID,omitempty"`
}

// EstadoSII representa el estado de un documento en el SII
type EstadoSII struct {
	Estado  string `xml:"ESTADO"`
	Glosa   string `xml:"GLOSA"`
	TrackID string `xml:"TRACKID,omitempty"`
}

// RespuestaEnvio representa la respuesta a un envío de DTE
type RespuestaEnvio struct {
	XMLName        xml.Name  `xml:"RespuestaEnvio"`
	Version        string    `xml:"version,attr"`
	TrackID        string    `xml:"TrackID"`
	Estado         string    `xml:"Estado"`
	GlosaEstado    string    `xml:"GlosaEstado"`
	NumeroAtencion string    `xml:"NumeroAtencion,omitempty"`
	FechaRecepcion time.Time `xml:"FechaRecepcion"`
}
