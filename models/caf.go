package models

import "time"

// CAFDTEXML representa un CAF (Código de Autorización de Folios) en formato XML
type CAFDTEXML struct {
	XMLName struct{}     `xml:"AUTORIZACION"`
	Version string       `xml:"version,attr"`
	DA      DAXMLModel   `xml:"DA"`
	FRMA    FRMAXMLModel `xml:"FRMA"`
}

// DAXMLModel representa los datos de autorización en formato XML
type DAXMLModel struct {
	RUT         RutXMLModel   `xml:"RE"`
	RazonSocial string        `xml:"RS"`
	TipoDTE     string        `xml:"TD"`
	RangoDesde  int           `xml:"RNG>D"`
	RangoHasta  int           `xml:"RNG>H"`
	FechaAut    string        `xml:"FA"`
	RSAPK       RSAPKXMLModel `xml:"RSAPK"`
	IDK         int           `xml:"IDK"`
}

// RutXMLModel representa un RUT en formato XML
type RutXMLModel struct {
	Numero string `xml:",chardata"`
}

// RSAPKXMLModel representa una llave pública RSA en formato XML
type RSAPKXMLModel struct {
	Modulo    string `xml:"M"`
	Exponente string `xml:"E"`
}

// FRMAXMLModel representa la firma del CAF en formato XML
type FRMAXMLModel struct {
	Algoritmo string `xml:"algoritmo,attr"`
	Valor     string `xml:",chardata"`
}

// CAF representa un certificado de autorización de folios
type CAF struct {
	ID                string    `json:"id" bson:"_id,omitempty"`
	RutEmisor         string    `json:"rut_emisor" bson:"rut_emisor"`
	TipoDTE           string    `json:"tipo_dte" bson:"tipo_dte"`
	RangoInicio       int       `json:"rango_inicio" bson:"rango_inicio"`
	RangoFin          int       `json:"rango_fin" bson:"rango_fin"`
	FolioActual       int       `json:"folio_actual" bson:"folio_actual"`
	Estado            string    `json:"estado" bson:"estado"`
	FechaAutorizacion time.Time `json:"fecha_autorizacion" bson:"fecha_autorizacion"`
	XML               []byte    `json:"xml" bson:"xml"`
	CreatedAt         time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" bson:"updated_at"`
}

// CAFRequest representa una solicitud de CAF
type CAFRequest struct {
	TipoDocumento string `json:"tipo_documento"`
	RutEmisor     string `json:"rut_emisor"`
	Cantidad      int    `json:"cantidad"`
}
