package siimodels

import (
	"encoding/xml"
)

// Signature representa la firma XML del documento
type Signature struct {
	XMLName        xml.Name   `xml:"Signature"`
	SignedInfo     SignedInfo `xml:"SignedInfo"`
	SignatureValue string     `xml:"SignatureValue"`
	KeyInfo        KeyInfo    `xml:"KeyInfo"`
}

// SignedInfo contiene la información firmada
type SignedInfo struct {
	XMLName                xml.Name               `xml:"SignedInfo"`
	CanonicalizationMethod CanonicalizationMethod `xml:"CanonicalizationMethod"`
	SignatureMethod        SignatureMethod        `xml:"SignatureMethod"`
	Reference              Reference              `xml:"Reference"`
}

// CanonicalizationMethod representa el método de canonicalización
type CanonicalizationMethod struct {
	Algorithm string `xml:"Algorithm,attr"`
}

// SignatureMethod representa el método de firma
type SignatureMethod struct {
	Algorithm string `xml:"Algorithm,attr"`
}

// Reference representa la referencia a los datos firmados
type Reference struct {
	URI          string       `xml:"URI,attr"`
	Transforms   *Transforms  `xml:"Transforms,omitempty"`
	DigestMethod DigestMethod `xml:"DigestMethod"`
	DigestValue  string       `xml:"DigestValue"`
}

// Transforms representa las transformaciones aplicadas
type Transforms struct {
	Transform []Transform `xml:"Transform"`
}

// Transform representa una transformación específica
type Transform struct {
	Algorithm string `xml:"Algorithm,attr"`
}

// DigestMethod representa el método de digest
type DigestMethod struct {
	Algorithm string `xml:"Algorithm,attr"`
}

// KeyInfo contiene la información de la llave
type KeyInfo struct {
	X509Data X509Data `xml:"X509Data"`
}

// X509Data contiene los datos del certificado X509
type X509Data struct {
	X509Certificate string `xml:"X509Certificate"`
}

// TED representa el Timbre Electrónico del DTE
type TED struct {
	Version   string     `xml:"version,attr"`
	DD        DD         `xml:"DD"`
	FRMT      string     `xml:"FRMT"`
	Signature *Signature `xml:"Signature,omitempty"`
}

// DD contiene los datos del timbre
type DD struct {
	RE    string `xml:"RE"`    // RUT Emisor
	TD    int    `xml:"TD"`    // Tipo de DTE
	F     int64  `xml:"F"`     // Folio
	FE    string `xml:"FE"`    // Fecha Emisión
	RR    string `xml:"RR"`    // RUT Receptor
	RSR   string `xml:"RSR"`   // Razón Social Receptor
	MNT   int64  `xml:"MNT"`   // Monto Total
	IT1   string `xml:"IT1"`   // Primer Item
	CAF   *CAF   `xml:"CAF"`   // Código de Autorización de Folios
	TSTED string `xml:"TSTED"` // TimeStamp de TED
}

// CAF representa el Código de Autorización de Folios
type CAF struct {
	Version string `xml:"version,attr"`
	DA      DA     `xml:"DA"`
	FRMA    string `xml:"FRMA"`
}

// DA contiene los datos de autorización
type DA struct {
	RE    string `xml:"RE"`    // RUT Emisor
	RS    string `xml:"RS"`    // Razón Social
	TD    int    `xml:"TD"`    // Tipo de DTE
	RNG   RNG    `xml:"RNG"`   // Rango de Folios
	FA    string `xml:"FA"`    // Fecha de Autorización
	RSAPK RSAPK  `xml:"RSAPK"` // Llave Pública RSA
	RSASK string `xml:"RSASK"` // Llave Privada RSA (opcional)
	IDK   int    `xml:"IDK"`   // ID de Llave
}

// RNG representa el rango de folios
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
	TpoMoneda string  `xml:"TpoMoneda"`
	TpoCambio float64 `xml:"TpoCambio"`
	MntExeOM  int64   `xml:"MntExeOM,omitempty"`
	MntNetoOM int64   `xml:"MntNetoOM,omitempty"`
	IVAom     int64   `xml:"IVAom,omitempty"`
	MntTotOM  int64   `xml:"MntTotOM"`
}
