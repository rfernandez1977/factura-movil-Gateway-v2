package models

import (
	"encoding/xml"
	"time"
)

// FirmaXML representa una firma digital XML
type FirmaXML struct {
	XMLName        xml.Name `xml:"Signature"`
	SignedInfo     SignedInfo
	SignatureValue string    `xml:"SignatureValue"`
	KeyInfo        KeyInfo   `xml:"KeyInfo"`
	Timestamp      time.Time `xml:"Timestamp,omitempty"`
}

// SignedInfo contiene la información firmada
type SignedInfo struct {
	XMLName                xml.Name  `xml:"SignedInfo"`
	CanonicalizationMethod Method    `xml:"CanonicalizationMethod"`
	SignatureMethod        Method    `xml:"SignatureMethod"`
	Reference              Reference `xml:"Reference"`
}

// Method representa un método algorítmico
type Method struct {
	Algorithm string `xml:"Algorithm,attr"`
}

// Reference contiene la referencia al documento firmado
type Reference struct {
	XMLName      xml.Name    `xml:"Reference"`
	URI          string      `xml:"URI,attr"`
	Transforms   []Transform `xml:"Transforms>Transform"`
	DigestMethod Method      `xml:"DigestMethod"`
	DigestValue  string      `xml:"DigestValue"`
}

// Transform representa una transformación XML
type Transform struct {
	Algorithm string `xml:"Algorithm,attr"`
}

// KeyInfo contiene la información de la clave
type KeyInfo struct {
	XMLName  xml.Name `xml:"KeyInfo"`
	X509Data X509Data `xml:"X509Data"`
}

// X509Data contiene los datos del certificado X509
type X509Data struct {
	XMLName         xml.Name `xml:"X509Data"`
	X509Certificate string   `xml:"X509Certificate"`
}

// EstadoFirma representa el estado de una firma
type EstadoFirma struct {
	Valida          bool      `json:"valida"`
	FechaValidacion time.Time `json:"fecha_validacion"`
	Error           string    `json:"error,omitempty"`
	CertificadoID   string    `json:"certificado_id"`
}
