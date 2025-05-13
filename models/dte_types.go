package models

import (
	"time"
)

// ReporteIntegracion representa un reporte de integración
type ReporteIntegracion struct {
	ID              string   `json:"id" bson:"_id,omitempty"`
	FechaInicio     string   `json:"fecha_inicio" bson:"fecha_inicio"`
	FechaFin        string   `json:"fecha_fin" bson:"fecha_fin"`
	TotalDocumentos int      `json:"total_documentos" bson:"total_documentos"`
	DocumentosOK    int      `json:"documentos_ok" bson:"documentos_ok"`
	DocumentosError int      `json:"documentos_error" bson:"documentos_error"`
	Errores         []string `json:"errores" bson:"errores"`
	Estado          string   `json:"estado" bson:"estado"`
}

// EstadoResult representa el resultado de una consulta de estado
type EstadoResult struct {
	Estado    string `json:"estado" xml:"Estado"`
	Mensaje   string `json:"mensaje" xml:"Mensaje"`
	CodigoSII string `json:"codigo_sii,omitempty" xml:"CodigoSII,omitempty"`
	GlosaSII  string `json:"glosa_sii,omitempty" xml:"GlosaSII,omitempty"`
	Exitoso   bool   `json:"exitoso" xml:"Exitoso"`
}

// Sobre representa un sobre para envío al SII
type Sobre struct {
	Documento []byte `json:"documento"`
	Token     string `json:"token"`
}

// DocumentoDTE representa un documento DTE
type DocumentoDTE struct {
	ID             string    `json:"id" bson:"_id,omitempty"`
	EmpresaID      string    `json:"empresa_id" bson:"empresa_id"`
	TipoDocumento  TipoDTE   `json:"tipo_documento" bson:"tipo_documento"`
	Folio          int       `json:"folio" bson:"folio"`
	Estado         string    `json:"estado" bson:"estado"`
	FechaEmision   time.Time `json:"fecha_emision" bson:"fecha_emision"`
	XMLOriginalID  string    `json:"xml_original_id" bson:"xml_original_id"`
	XMLFirmadoID   string    `json:"xml_firmado_id" bson:"xml_firmado_id"`
	XMLRespuestaID string    `json:"xml_respuesta_id" bson:"xml_respuesta_id"`
	FechaEnvio     time.Time `json:"fecha_envio" bson:"fecha_envio"`
	CreatedAt      time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" bson:"updated_at"`
}

// TipoDocumento representa los tipos de documentos tributarios
type TipoDocumento string

const (
	TipoDocumentoFactura           TipoDocumento = "FACTURA"
	TipoDocumentoFacturaExenta     TipoDocumento = "FACTURA_EXENTA"
	TipoDocumentoNotaCredito       TipoDocumento = "NOTA_CREDITO"
	TipoDocumentoNotaDebito        TipoDocumento = "NOTA_DEBITO"
	TipoDocumentoGuiaDespacho      TipoDocumento = "GUIA_DESPACHO"
	TipoDocumentoBoletaElectronica TipoDocumento = "BOLETA_ELECTRONICA"
	TipoDocumentoBoletaExenta      TipoDocumento = "BOLETA_EXENTA"
)

// TipoGuiaDespacho representa los tipos de guías de despacho
type TipoGuiaDespacho string

const (
	TipoGuiaDespachoTraslado    TipoGuiaDespacho = "TRASLADO"
	TipoGuiaDespachoVenta       TipoGuiaDespacho = "VENTA"
	TipoGuiaDespachoDevolucion  TipoGuiaDespacho = "DEVOLUCION"
	TipoGuiaDespachoExportacion TipoGuiaDespacho = "EXPORTACION"
	TipoGuiaDespachoImportacion TipoGuiaDespacho = "IMPORTACION"
)

// SobreDTEModel representa un sobre de envío de documentos tributarios electrónicos
type SobreDTEModel struct {
	XMLName   struct{} `xml:"EnvioDTE"`
	Version   string   `xml:"version,attr"`
	SetDTE    SetDTE   `xml:"SetDTE"`
	Signature string   `xml:"Signature,omitempty"`
}

// SetDTE representa un conjunto de documentos tributarios electrónicos
type SetDTE struct {
	ID       string        `xml:"ID,attr"`
	Caratula Caratula      `xml:"Caratula"`
	DTEs     []DTEXMLModel `xml:"DTE"`
}

// Caratula representa la información de encabezado del sobre
type Caratula struct {
	Version          string      `xml:"version,attr"`
	RutEmisor        string      `xml:"RutEmisor"`
	RutEnvia         string      `xml:"RutEnvia"`
	RutReceptor      string      `xml:"RutReceptor"`
	FechaResolucion  string      `xml:"FchResol"`
	NumeroResolucion int         `xml:"NroResol"`
	TmstFirmaEnv     string      `xml:"TmstFirmaEnv"`
	SubTotDTE        []SubTotDTE `xml:"SubTotDTE"`
}

// SubTotDTE representa los subtotales por tipo de documento
type SubTotDTE struct {
	TipoDTE string `xml:"TpoDTE"`
	NroDTE  int    `xml:"NroDTE"`
}

// ImpuestoAdicional representa un impuesto adicional aplicado a un documento
type ImpuestoAdicional struct {
	ID           string  `json:"id"`
	TipoImpuesto string  `json:"tipo_impuesto"`
	Tasa         float64 `json:"tasa"`
	Monto        float64 `json:"monto"`
}

// ImpuestoAdicionalItem representa un impuesto adicional aplicado a un ítem
type ImpuestoAdicionalItem struct {
	ID            string  `json:"id"`
	Item          Item    `json:"item"`
	TipoImpuesto  string  `json:"tipo_impuesto"`
	Porcentaje    float64 `json:"porcentaje"`
	BaseImponible float64 `json:"base_imponible"`
	Tasa          float64 `json:"tasa"`
	Monto         float64 `json:"monto"`
}

// SobreDTE representa un sobre con documentos tributarios electrónicos
type SobreDTE struct {
	ID           string                `json:"id" bson:"_id,omitempty"`
	RutEmisor    string                `json:"rut_emisor" bson:"rut_emisor"`
	RutEnvia     string                `json:"rut_envia" bson:"rut_envia"`
	RutReceptor  string                `json:"rut_receptor" bson:"rut_receptor"`
	Fecha        time.Time             `json:"fecha" bson:"fecha"`
	Documentos   []DocumentoTributario `json:"documentos" bson:"documentos"`
	Estado       string                `json:"estado" bson:"estado"`
	TrackID      string                `json:"track_id,omitempty" bson:"track_id,omitempty"`
	XMLEnvio     []byte                `json:"xml_envio,omitempty" bson:"xml_envio,omitempty"`
	XMLRespuesta []byte                `json:"xml_respuesta,omitempty" bson:"xml_respuesta,omitempty"`
	CreatedAt    time.Time             `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time             `json:"updated_at" bson:"updated_at"`
}
