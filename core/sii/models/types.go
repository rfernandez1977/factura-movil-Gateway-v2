package models

import "time"

// TipoDocumentoSII representa los tipos de documentos soportados por el SII
type TipoDocumentoSII string

const (
	DTEFactura       TipoDocumentoSII = "33"
	DTEFacturaExenta TipoDocumentoSII = "34"
	DTEBoleta        TipoDocumentoSII = "39"
	DTEBoletaExenta  TipoDocumentoSII = "41"
	DTEFacturaCompra TipoDocumentoSII = "46"
	DTEGuiaDespacho  TipoDocumentoSII = "52"
	DTENotaDebito    TipoDocumentoSII = "56"
	DTENotaCredito   TipoDocumentoSII = "61"
)

// EstadoSII representa los posibles estados de un documento en el SII
type EstadoSII string

const (
	EstadoAceptado     EstadoSII = "ACEPTADO"
	EstadoRechazado    EstadoSII = "RECHAZADO"
	EstadoPendiente    EstadoSII = "PENDIENTE"
	EstadoEnProceso    EstadoSII = "EPR"
	EstadoRecibido     EstadoSII = "REC"
	EstadoReparo       EstadoSII = "REP"
	EstadoDTERecibido  EstadoSII = "DTE_RECIBIDO"
	EstadoDTERechazado EstadoSII = "DTE_RECHAZADO"
	EstadoDTEProcesado EstadoSII = "DTE_PROCESADO"
	EstadoDTEEnReparo  EstadoSII = "DTE_EN_REPARO"
)

// RespuestaSII representa la respuesta del SII al enviar un documento
type RespuestaSII struct {
	TrackID      string     `json:"track_id" xml:"trackid"`
	Estado       EstadoSII  `json:"estado" xml:"estado"`
	Glosa        string     `json:"glosa" xml:"glosa"`
	NumAtencion  string     `json:"num_atencion,omitempty" xml:"numatencion,omitempty"`
	FechaProceso time.Time  `json:"fecha_proceso" xml:"fecha_proceso"`
	Errores      []ErrorSII `json:"errores,omitempty" xml:"errores>error,omitempty"`
}

// ErrorSII representa un error devuelto por el SII
type ErrorSII struct {
	Codigo      string `json:"codigo" xml:"codigo"`
	Descripcion string `json:"descripcion" xml:"descripcion"`
	Detalle     string `json:"detalle,omitempty" xml:"detalle,omitempty"`
}

// EstadoConsulta representa el estado de un documento consultado en el SII
type EstadoConsulta struct {
	TrackID         string    `json:"track_id" xml:"trackid"`
	Estado          EstadoSII `json:"estado" xml:"estado"`
	Glosa           string    `json:"glosa" xml:"glosa"`
	FechaRecepcion  time.Time `json:"fecha_recepcion" xml:"fecha_recepcion"`
	FechaProceso    time.Time `json:"fecha_proceso" xml:"fecha_proceso"`
	FechaAceptacion time.Time `json:"fecha_aceptacion,omitempty" xml:"fecha_aceptacion,omitempty"`
	FechaRechazo    time.Time `json:"fecha_rechazo,omitempty" xml:"fecha_rechazo,omitempty"`
}

// ValidacionSII representa el resultado de una validación del SII
type ValidacionSII struct {
	CodigoValidacion string    `json:"codigo_validacion" bson:"codigo_validacion"`
	Resultado        bool      `json:"resultado" bson:"resultado"`
	MensajeError     string    `json:"mensaje_error,omitempty" bson:"mensaje_error,omitempty"`
	FechaValidacion  time.Time `json:"fecha_validacion" bson:"fecha_validacion"`
}

// ConfiguracionSII representa la configuración para conectarse al SII
type ConfiguracionSII struct {
	Ambiente         string    `json:"ambiente"`
	RutEmpresa       string    `json:"rut_empresa"`
	RutCertificado   string    `json:"rut_certificado"`
	PathCertificado  string    `json:"path_certificado"`
	ClavePrivada     string    `json:"clave_privada"`
	URLProduccion    string    `json:"url_produccion"`
	URLCertificacion string    `json:"url_certificacion"`
	NumResolucion    int       `json:"num_resolucion"`
	FechaResolucion  time.Time `json:"fecha_resolucion"`
}

// DetalleSII representa el detalle de una respuesta del SII
type DetalleSII struct {
	Tipo     string `json:"tipo" xml:"TIPO"`
	Folio    int64  `json:"folio" xml:"FOLIO"`
	Estado   string `json:"estado" xml:"ESTADO"`
	Glosa    string `json:"glosa" xml:"GLOSA"`
	NumError int    `json:"numError,omitempty" xml:"NUM_ERROR,omitempty"`
}
