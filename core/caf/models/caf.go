package models

import (
	"encoding/xml"
	"time"
)

// CAF representa un Código de Autorización de Folios
type CAF struct {
	RUT              string
	RazonSocial      string
	TipoDTE          int
	FolioDesde       int
	FolioHasta       int
	FechaEmision     time.Time
	FechaVencimiento time.Time
	LlavePrivada     []byte
	LlavePublica     []byte
	XMLOriginal      []byte
}

// AutorizacionCAF representa la estructura XML del CAF
type AutorizacionCAF struct {
	XMLName xml.Name `xml:"AUTORIZACION"`
	Version string   `xml:"version,attr"`
	CAF     struct {
		DA struct {
			RE  string `xml:"RE"`
			RS  string `xml:"RS"`
			TD  int    `xml:"TD"`
			RNG struct {
				D int `xml:"D"`
				H int `xml:"H"`
			} `xml:"RNG"`
			FA    string `xml:"FA"`
			RSAPK struct {
				M string `xml:"M"`
				E string `xml:"E"`
			} `xml:"RSAPK"`
			IDK string `xml:"IDK"`
		} `xml:"DA"`
		FRMA struct {
			Algoritmo string `xml:"algoritmo,attr"`
			Valor     string `xml:",chardata"`
		} `xml:"FRMA"`
	} `xml:"CAF"`
}

// ResultadoValidacion representa el resultado de validar un CAF
type ResultadoValidacion struct {
	Valido    bool
	Error     error
	Detalles  string
	Timestamp time.Time
}
