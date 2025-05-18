package siimodels

import (
	"encoding/xml"
	"time"
)

// SolicitudSesion representa la solicitud de inicio de sesión electrónica con el SII
type SolicitudSesion struct {
	XMLName            xml.Name           `xml:"SolicitudSesion"`
	Version            string             `xml:"version,attr"`
	RutContribuyente   string             `xml:"RutContribuyente"`
	CertificadoDigital CertificadoDigital `xml:"CertificadoDigital"`
	Ambiente           AmbienteSII        `xml:"Ambiente"`
	Timestamp          time.Time          `xml:"Timestamp"`
}

// CertificadoDigital representa la información del certificado digital
type CertificadoDigital struct {
	SerialNumber string    `xml:"SerialNumber"`
	Issuer       string    `xml:"Issuer"`
	ValidFrom    time.Time `xml:"ValidFrom"`
	ValidTo      time.Time `xml:"ValidTo"`
}

// RespuestaSesion representa la respuesta a la solicitud de sesión
type RespuestaSesion struct {
	XMLName     xml.Name `xml:"RespuestaSesion"`
	Version     string   `xml:"version,attr"`
	Estado      string   `xml:"Estado"`
	Token       *Token   `xml:"Token,omitempty"`
	Mensaje     string   `xml:"Mensaje,omitempty"`
	CodigoError string   `xml:"CodigoError,omitempty"`
}

// Token representa el token de sesión
type Token struct {
	Valor  string    `xml:"Valor"`
	Expira time.Time `xml:"Expira"`
}

// RenovacionSesion representa la solicitud de renovación de sesión
type RenovacionSesion struct {
	XMLName   xml.Name  `xml:"RenovacionSesion"`
	Version   string    `xml:"version,attr"`
	Token     string    `xml:"Token"`
	Timestamp time.Time `xml:"Timestamp"`
}

// CierreSesion representa la solicitud de cierre de sesión
type CierreSesion struct {
	XMLName   xml.Name  `xml:"CierreSesion"`
	Version   string    `xml:"version,attr"`
	Token     string    `xml:"Token"`
	Timestamp time.Time `xml:"Timestamp"`
}

// Constantes para los estados de respuesta
const (
	EstadoOK    = "OK"
	EstadoError = "ERROR"
)

// NewSolicitudSesion crea una nueva solicitud de sesión
func NewSolicitudSesion(rutContribuyente string, certificado *CertificadoDigital, ambiente AmbienteSII) *SolicitudSesion {
	return &SolicitudSesion{
		Version:            "1.0",
		RutContribuyente:   rutContribuyente,
		CertificadoDigital: *certificado,
		Ambiente:           ambiente,
		Timestamp:          time.Now(),
	}
}
