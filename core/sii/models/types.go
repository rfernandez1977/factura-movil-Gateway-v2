package models

// Ambiente representa el ambiente de ejecución del SII (certificación o producción)
type Ambiente string

const (
	// Certificacion representa el ambiente de pruebas del SII
	Certificacion Ambiente = "certificacion"
	// Produccion representa el ambiente de producción del SII
	Produccion Ambiente = "produccion"
)

// String devuelve la representación en string del ambiente
func (a Ambiente) String() string {
	switch a {
	case Certificacion:
		return "Ambiente de Certificación"
	case Produccion:
		return "Ambiente de Producción"
	default:
		return "Ambiente Desconocido"
	}
}

// IsValid verifica si el ambiente es válido
func (a Ambiente) IsValid() bool {
	return a == Certificacion || a == Produccion
}

// URLs base para cada ambiente del SII
const (
	URLBaseCertificacion = "https://maullin.sii.cl"
	URLBaseProduccion    = "https://palena.sii.cl"
)

// Endpoints específicos del SII
const (
	EndpointSemilla  = "/DTEWS/CrSeed.jws"           // Endpoint para obtener semilla
	EndpointToken    = "/DTEWS/GetTokenFromSeed.jws" // Endpoint para obtener token
	EndpointEnvio    = "/cgi_dte/UPL/DTEUpload"      // Endpoint para enviar DTE
	EndpointConsulta = "/DTEWS/QueryEstDte.jws"      // Endpoint para consultar estado
)

// TipoDocumentoSII representa los tipos de documentos tributarios electrónicos soportados por el SII
type TipoDocumentoSII string

const (
	// Tipos de documentos tributarios electrónicos
	TipoFactura       TipoDocumentoSII = "33" // Factura Electrónica
	TipoFacturaExenta TipoDocumentoSII = "34" // Factura Electrónica Exenta
	TipoBoleta        TipoDocumentoSII = "39" // Boleta Electrónica
	TipoBoletaExenta  TipoDocumentoSII = "41" // Boleta Electrónica Exenta
	TipoNotaCredito   TipoDocumentoSII = "61" // Nota de Crédito Electrónica
	TipoNotaDebito    TipoDocumentoSII = "56" // Nota de Débito Electrónica
	TipoGuiaDespacho  TipoDocumentoSII = "52" // Guía de Despacho Electrónica
)

// String devuelve la representación en string del tipo de documento
func (t TipoDocumentoSII) String() string {
	switch t {
	case TipoFactura:
		return "Factura Electrónica"
	case TipoFacturaExenta:
		return "Factura Electrónica Exenta"
	case TipoBoleta:
		return "Boleta Electrónica"
	case TipoBoletaExenta:
		return "Boleta Electrónica Exenta"
	case TipoNotaCredito:
		return "Nota de Crédito Electrónica"
	case TipoNotaDebito:
		return "Nota de Débito Electrónica"
	case TipoGuiaDespacho:
		return "Guía de Despacho Electrónica"
	default:
		return "Tipo de Documento Desconocido"
	}
}

// IsValid verifica si el tipo de documento es válido
func (t TipoDocumentoSII) IsValid() bool {
	switch t {
	case TipoFactura, TipoFacturaExenta, TipoBoleta, TipoBoletaExenta,
		TipoNotaCredito, TipoNotaDebito, TipoGuiaDespacho:
		return true
	default:
		return false
	}
}

// EstadoConsulta representa el estado de una consulta al SII
type EstadoConsulta struct {
	Estado  EstadoSII `xml:"ESTADO" json:"estado"`                        // Estado de la consulta
	Glosa   string    `xml:"GLOSA" json:"glosa"`                          // Descripción del estado
	TrackID string    `xml:"TRACKID,omitempty" json:"track_id,omitempty"` // ID de seguimiento
}
