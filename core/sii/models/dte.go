package models

import (
	"encoding/xml"
	"time"
)

// Signature representa la firma digital XML
type Signature struct {
	XMLName        xml.Name   `xml:"Signature"`
	SignedInfo     SignedInfo `xml:"SignedInfo"`
	SignatureValue string     `xml:"SignatureValue"`
	KeyInfo        KeyInfo    `xml:"KeyInfo"`
}

// SignedInfo representa la información firmada
type SignedInfo struct {
	XMLName                xml.Name               `xml:"SignedInfo"`
	CanonicalizationMethod CanonicalizationMethod `xml:"CanonicalizationMethod"`
	SignatureMethod        SignatureMethod        `xml:"SignatureMethod"`
	Reference              Reference              `xml:"Reference"`
}

// CanonicalizationMethod representa el método de canonicalización
type CanonicalizationMethod struct {
	XMLName   xml.Name `xml:"CanonicalizationMethod"`
	Algorithm string   `xml:"Algorithm,attr"`
}

// SignatureMethod representa el método de firma
type SignatureMethod struct {
	XMLName   xml.Name `xml:"SignatureMethod"`
	Algorithm string   `xml:"Algorithm,attr"`
}

// Reference representa una referencia para la firma
type Reference struct {
	XMLName      xml.Name     `xml:"Reference"`
	URI          string       `xml:"URI,attr"`
	Transforms   Transforms   `xml:"Transforms"`
	DigestMethod DigestMethod `xml:"DigestMethod"`
	DigestValue  string       `xml:"DigestValue"`
}

// Transforms representa las transformaciones aplicadas
type Transforms struct {
	XMLName   xml.Name    `xml:"Transforms"`
	Transform []Transform `xml:"Transform"`
}

// Transform representa una transformación específica
type Transform struct {
	XMLName   xml.Name `xml:"Transform"`
	Algorithm string   `xml:"Algorithm,attr"`
}

// DigestMethod representa el método de digest
type DigestMethod struct {
	XMLName   xml.Name `xml:"DigestMethod"`
	Algorithm string   `xml:"Algorithm,attr"`
}

// KeyInfo representa la información de la clave
type KeyInfo struct {
	XMLName  xml.Name `xml:"KeyInfo"`
	KeyValue KeyValue `xml:"KeyValue"`
	X509Data X509Data `xml:"X509Data"`
}

// KeyValue representa el valor de la clave
type KeyValue struct {
	XMLName     xml.Name    `xml:"KeyValue"`
	RSAKeyValue RSAKeyValue `xml:"RSAKeyValue"`
}

// RSAKeyValue representa el valor de la clave RSA
type RSAKeyValue struct {
	XMLName  xml.Name `xml:"RSAKeyValue"`
	Modulus  string   `xml:"Modulus"`
	Exponent string   `xml:"Exponent"`
}

// X509Data representa datos del certificado X509
type X509Data struct {
	XMLName         xml.Name `xml:"X509Data"`
	X509Certificate string   `xml:"X509Certificate"`
}

// DTE representa un Documento Tributario Electrónico
type DTE struct {
	XMLName   xml.Name   `xml:"DTE"`
	Documento Documento  `xml:"Documento"`
	Signature *Signature `xml:"Signature,omitempty"`
	TmstFirma time.Time  `xml:"TmstFirma"`
}

// Documento representa la estructura principal del DTE
type Documento struct {
	XMLName      xml.Name     `xml:"Documento"`
	Encabezado   Encabezado   `xml:"Encabezado"`
	Detalle      []Detalle    `xml:"Detalle"`
	SubTotInfo   []SubTotInfo `xml:"SubTotInfo,omitempty"`
	DscRcgGlobal []DscRcg     `xml:"DscRcgGlobal,omitempty"`
	Referencia   []Referencia `xml:"Referencia,omitempty"`
	TED          *TED         `xml:"TED,omitempty"`
	TmstFirma    time.Time    `xml:"TmstFirma"`
}

// Encabezado contiene la información principal del DTE
type Encabezado struct {
	IdDoc      IdDoc       `xml:"IdDoc"`
	Emisor     Emisor      `xml:"Emisor"`
	Receptor   Receptor    `xml:"Receptor"`
	Totales    Totales     `xml:"Totales"`
	OtraMoneda *OtraMoneda `xml:"OtraMoneda,omitempty"`
}

// IdDoc contiene la identificación del documento
type IdDoc struct {
	TipoDTE      string    `xml:"TipoDTE"`
	Folio        int64     `xml:"Folio"`
	FchEmis      time.Time `xml:"FchEmis"`
	IndServicio  int       `xml:"IndServicio,omitempty"`
	IndMntNeto   int       `xml:"IndMntNeto,omitempty"`
	PeriodoDesde time.Time `xml:"PeriodoDesde,omitempty"`
	PeriodoHasta time.Time `xml:"PeriodoHasta,omitempty"`
	FchVenc      time.Time `xml:"FchVenc,omitempty"`
}

// Emisor contiene la información del emisor del DTE
type Emisor struct {
	RUTEmisor    string `xml:"RUTEmisor"`
	RznSoc       string `xml:"RznSoc"`
	GiroEmis     string `xml:"GiroEmis"`
	Acteco       string `xml:"Acteco"`
	DirOrigen    string `xml:"DirOrigen"`
	CmnaOrigen   string `xml:"CmnaOrigen"`
	CiudadOrigen string `xml:"CiudadOrigen,omitempty"`
}

// Receptor contiene la información del receptor del DTE
type Receptor struct {
	RUTRecep    string `xml:"RUTRecep"`
	RznSocRecep string `xml:"RznSocRecep"`
	GiroRecep   string `xml:"GiroRecep"`
	DirRecep    string `xml:"DirRecep"`
	CmnaRecep   string `xml:"CmnaRecep"`
	CiudadRecep string `xml:"CiudadRecep,omitempty"`
}

// Totales contiene los montos totales del DTE
type Totales struct {
	MntNeto  int64   `xml:"MntNeto,omitempty"`
	MntExe   int64   `xml:"MntExe,omitempty"`
	TasaIVA  float64 `xml:"TasaIVA,omitempty"`
	IVA      int64   `xml:"IVA,omitempty"`
	MntTotal int64   `xml:"MntTotal"`
}

// Detalle representa cada ítem del DTE
type Detalle struct {
	NroLinDet      int     `xml:"NroLinDet"`
	CdgItem        CdgItem `xml:"CdgItem,omitempty"`
	IndExe         int     `xml:"IndExe,omitempty"`
	NmbItem        string  `xml:"NmbItem"`
	DscItem        string  `xml:"DscItem,omitempty"`
	QtyItem        float64 `xml:"QtyItem,omitempty"`
	UnmdItem       string  `xml:"UnmdItem,omitempty"`
	PrcItem        float64 `xml:"PrcItem,omitempty"`
	DescuentoPct   float64 `xml:"DescuentoPct,omitempty"`
	DescuentoMonto int64   `xml:"DescuentoMonto,omitempty"`
	RecargoPct     float64 `xml:"RecargoPct,omitempty"`
	RecargoMonto   int64   `xml:"RecargoMonto,omitempty"`
	MontoItem      int64   `xml:"MontoItem"`
}

// CdgItem representa los códigos de un ítem
type CdgItem struct {
	TpoCodigo string `xml:"TpoCodigo"`
	VlrCodigo string `xml:"VlrCodigo"`
}

// SubTotInfo representa información de subtotales
type SubTotInfo struct {
	NroSTI         int    `xml:"NroSTI"`
	GlosaSTI       string `xml:"GlosaSTI,omitempty"`
	SubTotMntNeto  int64  `xml:"SubTotMntNeto,omitempty"`
	SubTotMntIVA   int64  `xml:"SubTotMntIVA,omitempty"`
	SubTotMntExe   int64  `xml:"SubTotMntExe,omitempty"`
	SubTotMntTotal int64  `xml:"SubTotMntTotal"`
}

// DscRcg representa descuentos o recargos globales
type DscRcg struct {
	NroLinDR int     `xml:"NroLinDR"`
	TpoMov   string  `xml:"TpoMov"`
	GlosaDR  string  `xml:"GlosaDR,omitempty"`
	TpoValor string  `xml:"TpoValor"`
	ValorDR  float64 `xml:"ValorDR"`
	IndExeDR int     `xml:"IndExeDR,omitempty"`
}

// Referencia representa referencias a otros documentos
type Referencia struct {
	NroLinRef int       `xml:"NroLinRef"`
	TpoDocRef string    `xml:"TpoDocRef"`
	FolioRef  string    `xml:"FolioRef"`
	FchRef    time.Time `xml:"FchRef"`
	CodRef    string    `xml:"CodRef,omitempty"`
	RazonRef  string    `xml:"RazonRef,omitempty"`
}

// TED representa el Timbre Electrónico DTE
type TED struct {
	Version   string     `xml:"version,attr"`
	DD        DD         `xml:"DD"`
	FRMT      string     `xml:"FRMT"`
	Signature *Signature `xml:"Signature,omitempty"`
}

// DD representa los datos del DTE para el timbre
type DD struct {
	RE    string `xml:"RE"`    // RUT Emisor
	TD    string `xml:"TD"`    // Tipo DTE
	F     int64  `xml:"F"`     // Folio
	FE    string `xml:"FE"`    // Fecha Emisión
	RR    string `xml:"RR"`    // RUT Receptor
	RSR   string `xml:"RSR"`   // Razón Social Receptor
	MNT   int64  `xml:"MNT"`   // Monto Total
	IT1   string `xml:"IT1"`   // Item 1
	CAF   *CAF   `xml:"CAF"`   // Código Autorización de Folios
	TSTED string `xml:"TSTED"` // TimeStamp
}

// CAF representa el Código de Autorización de Folios
type CAF struct {
	Version string `xml:"version,attr"`
	DA      DA     `xml:"DA"`
	FRMA    string `xml:"FRMA"`
}

// DA representa los datos de autorización
type DA struct {
	RE    string    `xml:"RE"`    // RUT Emisor
	RS    string    `xml:"RS"`    // Razón Social
	TD    string    `xml:"TD"`    // Tipo DTE
	RNG   RNG       `xml:"RNG"`   // Rango
	FA    time.Time `xml:"FA"`    // Fecha Autorización
	RSAPK RSAPK     `xml:"RSAPK"` // Llave Pública
	IDK   int       `xml:"IDK"`   // ID Llave
}

// RNG representa el rango de folios autorizados
type RNG struct {
	D int64 `xml:"D"` // Desde
	H int64 `xml:"H"` // Hasta
}

// RSAPK representa la llave pública RSA
type RSAPK struct {
	M string `xml:"M"` // Módulo
	E string `xml:"E"` // Exponente
}

// OtraMoneda representa montos en otra moneda
type OtraMoneda struct {
	TpoMoneda        string  `xml:"TpoMoneda"`
	TpoCambio        float64 `xml:"TpoCambio"`
	MntNetoOtrMnda   int64   `xml:"MntNetoOtrMnda,omitempty"`
	MntExeOtrMnda    int64   `xml:"MntExeOtrMnda,omitempty"`
	MntFaecarOtrMnda int64   `xml:"MntFaecarOtrMnda,omitempty"`
	MntTotOtrMnda    int64   `xml:"MntTotOtrMnda"`
}
