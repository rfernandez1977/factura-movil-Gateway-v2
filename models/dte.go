package models

import (
	"encoding/xml"
	"time"
)

// DTEDocument representa un Documento Tributario Electrónico
type DTEDocument struct {
	XMLName   xml.Name `xml:"DTE"`
	Version   string   `xml:"version,attr"`
	Documento DTEDocumento
}

// DTEDocumento representa el contenido principal del DTE
type DTEDocumento struct {
	XMLName    xml.Name      `xml:"Documento"`
	ID         string        `xml:"ID,attr"`
	Encabezado DTEEncabezado `xml:"Encabezado"`
	Detalle    []DTEDetalle  `xml:"Detalle"`
}

// DTEEncabezado contiene la información principal del documento
type DTEEncabezado struct {
	IdDoc    DTEIdDoc    `xml:"IdDoc"`
	Emisor   DTEEmisor   `xml:"Emisor"`
	Receptor DTEReceptor `xml:"Receptor"`
	Totales  DTETotales  `xml:"Totales"`
}

// DTEIdDoc contiene la identificación del documento
type DTEIdDoc struct {
	TipoDTE string    `xml:"TipoDTE"`
	Folio   int       `xml:"Folio"`
	FchEmis time.Time `xml:"FchEmis"`
}

// DTEEmisor contiene la información del emisor del documento
type DTEEmisor struct {
	RUTEmisor  string `xml:"RUTEmisor"`
	RznSoc     string `xml:"RznSoc"`
	GiroEmis   string `xml:"GiroEmis"`
	Acteco     string `xml:"Acteco"`
	DirOrigen  string `xml:"DirOrigen"`
	CmnaOrigen string `xml:"CmnaOrigen"`
}

// DTEReceptor contiene la información del receptor del documento
type DTEReceptor struct {
	RUTRecep    string `xml:"RUTRecep"`
	RznSocRecep string `xml:"RznSocRecep"`
	GiroRecep   string `xml:"GiroRecep"`
	DirRecep    string `xml:"DirRecep"`
	CmnaRecep   string `xml:"CmnaRecep"`
}

// DTETotales contiene los totales del documento
type DTETotales struct {
	MntNeto  int     `xml:"MntNeto"`
	TasaIVA  float64 `xml:"TasaIVA"`
	IVA      int     `xml:"IVA"`
	MntTotal int     `xml:"MntTotal"`
}

// DTEDetalle representa una línea de detalle del documento
type DTEDetalle struct {
	NroLinDet int     `xml:"NroLinDet"`
	NmbItem   string  `xml:"NmbItem"`
	QtyItem   float64 `xml:"QtyItem"`
	PrcItem   float64 `xml:"PrcItem"`
	MontoItem int     `xml:"MontoItem"`
}

// EnvioDTEDocument representa un sobre de envío de DTEs
type EnvioDTEDocument struct {
	XMLName   xml.Name `xml:"EnvioDTE"`
	Version   string   `xml:"version,attr"`
	SetDTE    DTESet   `xml:"SetDTE"`
	Signature string   `xml:"Signature,omitempty"`
}

// DTESet representa un conjunto de DTEs
type DTESet struct {
	ID       string        `xml:"ID,attr"`
	Caratula DTECaratula   `xml:"Caratula"`
	DTE      []DTEDocument `xml:"DTE"`
}

// DTECaratula contiene la información de encabezado del envío
type DTECaratula struct {
	RutEmisor    string    `xml:"RutEmisor"`
	RutEnvia     string    `xml:"RutEnvia"`
	RutReceptor  string    `xml:"RutReceptor"`
	FchResol     time.Time `xml:"FchResol"`
	NroResol     int       `xml:"NroResol"`
	TmstFirmaEnv time.Time `xml:"TmstFirmaEnv"`
}
