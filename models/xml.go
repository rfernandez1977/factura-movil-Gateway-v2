package models

import (
	"encoding/xml"
	"time"
)

// Tipos XML principales utilizados para documentos tributarios electrónicos

// DTEXMLModel representa un documento tributario electrónico en formato XML
type DTEXMLModel struct {
	XMLName   xml.Name          `xml:"DTE"`
	Version   string            `xml:"version,attr"`
	Documento DocumentoXMLModel `xml:"Documento"`
	Signature SignatureXMLModel `xml:"Signature,omitempty"`
}

// DocumentoXMLModel representa un documento en XML
type DocumentoXMLModel struct {
	XMLName     xml.Name             `xml:"Documento"`
	ID          string               `xml:"ID,attr,omitempty"`
	Encabezado  EncabezadoXMLModel   `xml:"Encabezado"`
	Detalles    []DetalleXML         `xml:"Detalle"`
	Referencias []ReferenciaXMLModel `xml:"Referencias>Referencia,omitempty"`
}

// EncabezadoXMLModel representa el encabezado de un documento
type EncabezadoXMLModel struct {
	XMLName     xml.Name       `xml:"Encabezado"`
	IDDocumento IDDocumentoXML `xml:"IdDoc"`
	Emisor      EmisorXML      `xml:"Emisor"`
	Receptor    ReceptorXML    `xml:"Receptor"`
	Totales     TotalesXML     `xml:"Totales"`
}

// IDDocumentoXML representa la identificación del documento
type IDDocumentoXML struct {
	XMLName           xml.Name `xml:"IdDoc"`
	TipoDTE           string   `xml:"TipoDTE"`
	Folio             int      `xml:"Folio"`
	FechaEmision      string   `xml:"FchEmis"`
	TipoDespacho      string   `xml:"TipoDespacho,omitempty"`
	IndicadorServicio int      `xml:"IndServicio,omitempty"`
}

// EmisorXML representa los datos del emisor
type EmisorXML struct {
	XMLName     xml.Name `xml:"Emisor"`
	RUT         string   `xml:"RUTEmisor"`
	RazonSocial string   `xml:"RznSoc"`
	Giro        string   `xml:"GiroEmis"`
	Telefono    string   `xml:"Telefono,omitempty"`
	Email       string   `xml:"CorreoEmisor,omitempty"`
	Direccion   string   `xml:"DirOrigen"`
	Comuna      string   `xml:"CmnaOrigen"`
	Ciudad      string   `xml:"CiudadOrigen"`
}

// ReceptorXML representa los datos del receptor
type ReceptorXML struct {
	XMLName     xml.Name `xml:"Receptor"`
	RUT         string   `xml:"RUTRecep"`
	RazonSocial string   `xml:"RznSocRecep"`
	Giro        string   `xml:"GiroRecep,omitempty"`
	Telefono    string   `xml:"Contacto,omitempty"`
	Email       string   `xml:"CorreoRecep,omitempty"`
	Direccion   string   `xml:"DirRecep"`
	Comuna      string   `xml:"CmnaRecep"`
	Ciudad      string   `xml:"CiudadRecep"`
}

// TotalesXML representa los totales del documento
type TotalesXML struct {
	XMLName     xml.Name `xml:"Totales"`
	MontoNeto   int      `xml:"MntNeto,omitempty"`
	MontoExento int      `xml:"MntExe,omitempty"`
	TasaIVA     float64  `xml:"TasaIVA,omitempty"`
	IVA         int      `xml:"IVA,omitempty"`
	MontoTotal  int      `xml:"MntTotal"`
}

// DetalleXML representa un detalle de producto o servicio
type DetalleXML struct {
	XMLName        xml.Name      `xml:"Detalle"`
	NumeroLinea    int           `xml:"NroLinDet"`
	TipoDocumento  string        `xml:"TpoDocLiq,omitempty"`
	Codigo         string        `xml:"CdgItem>TpoCodigo,omitempty"`
	ValorCodigo    string        `xml:"CdgItem>VlrCodigo,omitempty"`
	Nombre         string        `xml:"NmbItem"`
	Descripcion    string        `xml:"DscItem,omitempty"`
	Cantidad       float64       `xml:"QtyItem,omitempty"`
	UnidadMedida   string        `xml:"UnmdItem,omitempty"`
	PrecioUnitario float64       `xml:"PrcItem,omitempty"`
	Descuento      float64       `xml:"DescuentoMonto,omitempty"`
	PorcentajeDesc float64       `xml:"DescuentoPct,omitempty"`
	SubTotal       int           `xml:"MontoItem"`
	Impuestos      []ImpuestoXML `xml:"ImptoReten,omitempty"`
}

// ImpuestoXML representa un impuesto en un detalle
type ImpuestoXML struct {
	XMLName xml.Name `xml:"ImptoReten"`
	Tipo    string   `xml:"TipoImp"`
	Tasa    float64  `xml:"TasaImp"`
	Monto   int      `xml:"MontoImp"`
}

// ReferenciaXMLModel representa una referencia a otro documento
type ReferenciaXMLModel struct {
	XMLName    xml.Name `xml:"Referencia"`
	TipoDocRef string   `xml:"TpoDocRef"`
	FolioRef   string   `xml:"FolioRef"`
	FechaRef   string   `xml:"FchRef"`
	CodigoRef  string   `xml:"CodRef,omitempty"`
	RazonRef   string   `xml:"RazonRef,omitempty"`
}

// SignatureXMLModel representa la firma digital del documento
type SignatureXMLModel struct {
	XMLName        xml.Name      `xml:"Signature"`
	Xmlns          string        `xml:"xmlns,attr"`
	SignedInfo     SignedInfoXML `xml:"SignedInfo"`
	SignatureValue string        `xml:"SignatureValue"`
	KeyInfo        KeyInfoXML    `xml:"KeyInfo"`
}

// SignedInfoXML representa la información firmada
type SignedInfoXML struct {
	XMLName                xml.Name                  `xml:"SignedInfo"`
	CanonicalizationMethod CanonicalizationMethodXML `xml:"CanonicalizationMethod"`
	SignatureMethod        SignatureMethodXML        `xml:"SignatureMethod"`
	Reference              ReferenceXML              `xml:"Reference"`
}

// CanonicalizationMethodXML representa el método de canonicalización
type CanonicalizationMethodXML struct {
	XMLName   xml.Name `xml:"CanonicalizationMethod"`
	Algorithm string   `xml:"Algorithm,attr"`
}

// SignatureMethodXML representa el método de firma
type SignatureMethodXML struct {
	XMLName   xml.Name `xml:"SignatureMethod"`
	Algorithm string   `xml:"Algorithm,attr"`
}

// ReferenceXML representa una referencia para la firma
type ReferenceXML struct {
	XMLName      xml.Name        `xml:"Reference"`
	URI          string          `xml:"URI,attr"`
	Transforms   TransformsXML   `xml:"Transforms"`
	DigestMethod DigestMethodXML `xml:"DigestMethod"`
	DigestValue  string          `xml:"DigestValue"`
}

// TransformsXML representa transformaciones para la firma
type TransformsXML struct {
	XMLName   xml.Name       `xml:"Transforms"`
	Transform []TransformXML `xml:"Transform"`
}

// TransformXML representa una transformación para la firma
type TransformXML struct {
	XMLName   xml.Name `xml:"Transform"`
	Algorithm string   `xml:"Algorithm,attr"`
}

// DigestMethodXML representa el método de resumen
type DigestMethodXML struct {
	XMLName   xml.Name `xml:"DigestMethod"`
	Algorithm string   `xml:"Algorithm,attr"`
}

// KeyInfoXML representa la información de la clave
type KeyInfoXML struct {
	XMLName  xml.Name    `xml:"KeyInfo"`
	X509Data X509DataXML `xml:"X509Data"`
}

// X509DataXML representa datos de certificado X509
type X509DataXML struct {
	XMLName         xml.Name `xml:"X509Data"`
	X509Certificate string   `xml:"X509Certificate"`
}

// FacturaElectronica representa una factura electrónica
type FacturaElectronica struct {
	DTEBase `json:",inline"`
}

// BoletaElectronicaTipo representa una boleta electrónica
type BoletaElectronicaTipo struct {
	DTEBase `json:",inline"`
}

// GuiaDespachoElectronica representa una guía de despacho electrónica
type GuiaDespachoElectronica struct {
	DTEBase `json:",inline"`
}

// NotaDebitoElectronica representa una nota de débito electrónica
type NotaDebitoElectronica struct {
	DTEBase `json:",inline"`
}

// NotaCreditoElectronica representa una nota de crédito electrónica
type NotaCreditoElectronica struct {
	DTEBase `json:",inline"`
}

// DTEBase contiene los campos básicos compartidos por todos los tipos de DTE
type DTEBase struct {
	ID           string       `json:"id" bson:"_id,omitempty"`
	TipoDTE      TipoDTE      `json:"tipo_dte" bson:"tipo_dte"`
	Folio        int          `json:"folio" bson:"folio"`
	FechaEmision time.Time    `json:"fecha_emision" bson:"fecha_emision"`
	EmisorID     string       `json:"emisor_id" bson:"emisor_id"`
	ReceptorID   string       `json:"receptor_id" bson:"receptor_id"`
	MontoNeto    float64      `json:"monto_neto" bson:"monto_neto"`
	MontoExento  float64      `json:"monto_exento" bson:"monto_exento"`
	MontoIVA     float64      `json:"monto_iva" bson:"monto_iva"`
	MontoTotal   float64      `json:"monto_total" bson:"monto_total"`
	Estado       string       `json:"estado" bson:"estado"`
	EstadoSII    string       `json:"estado_sii" bson:"estado_sii"`
	TrackID      string       `json:"track_id,omitempty" bson:"track_id,omitempty"`
	XML          string       `json:"xml,omitempty" bson:"xml,omitempty"`
	XMLDTE       string       `json:"xml_dte,omitempty" bson:"xml_dte,omitempty"`
	URLPublica   string       `json:"url_publica,omitempty" bson:"url_publica,omitempty"`
	Detalles     []DetalleXML `json:"detalles,omitempty" bson:"detalles,omitempty"`
	CreatedAt    time.Time    `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at" bson:"updated_at"`
}
