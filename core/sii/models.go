package sii

import (
	"encoding/xml"
)

// Envelope estructura base para mensajes SOAP
type Envelope struct {
	XMLName xml.Name `xml:"Envelope"`
	XMLNS   string   `xml:"xmlns,attr"`
	Body    Body     `xml:"Body"`
}

// Body estructura del cuerpo SOAP
type Body struct {
	XMLName xml.Name    `xml:"Body"`
	Content interface{} `xml:",any"`
	Fault   *Fault      `xml:"Fault,omitempty"`
}

// Fault estructura para errores SOAP
type Fault struct {
	XMLName     xml.Name `xml:"Fault"`
	FaultCode   string   `xml:"faultcode"`
	FaultString string   `xml:"faultstring"`
}

// SemillaResponse estructura de respuesta para la semilla
type SemillaResponse struct {
	XMLName xml.Name `xml:"getSeedResponse"`
	Return  string   `xml:"return"`
}

// TokenResponse estructura de respuesta para el token
type TokenResponse struct {
	XMLName xml.Name `xml:"getTokenResponse"`
	Return  string   `xml:"return"`
}

// SemillaRequest estructura para la solicitud de semilla
type SemillaRequest struct {
	XMLName xml.Name `xml:"getSeed"`
}

// TokenRequest estructura para la solicitud de token
type TokenRequest struct {
	XMLName xml.Name `xml:"getToken"`
	Token   string   `xml:"token"`
}

// Error estructura para errores del SII
type Error struct {
	Code    string `xml:"codigo"`
	Message string `xml:"mensaje"`
}
