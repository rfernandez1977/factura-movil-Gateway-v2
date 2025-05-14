package models

import "encoding/xml"

// SobreDTEModel representa un sobre de DTE para envío al SII
type SobreDTEModel struct {
	XMLName   xml.Name       `xml:"EnvioDTE"`
	SetDTE    *SetDTE        `xml:"SetDTE"`
	Signature *FirmaXMLModel `xml:"Signature,omitempty"`
}

// SetDTE representa el conjunto de DTEs dentro del sobre
type SetDTE struct {
	ID       string        `xml:"ID,attr"`
	Caratula *Caratula     `xml:"Caratula"`
	DTEs     []DTEXMLModel `xml:"DTE"`
}

// Caratula representa la información de caratula del envío
type Caratula struct {
	Version          string `xml:"version,attr"`
	RutEmisor        string `xml:"RutEmisor"`
	RutEnvia         string `xml:"RutEnvia"`
	RutReceptor      string `xml:"RutReceptor"`
	FechaResolucion  string `xml:"FechaResolucion"`
	NumeroResolucion int    `xml:"NumeroResolucion"`
	TmstFirmaEnv     string `xml:"TmstFirmaEnv"`
}
