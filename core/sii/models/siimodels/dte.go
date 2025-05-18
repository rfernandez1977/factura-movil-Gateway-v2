package siimodels

import (
	"encoding/xml"
	"time"
)

// DTE representa un Documento Tributario Electrónico
type DTE struct {
	XMLName   xml.Name   `xml:"DTE"`
	Version   string     `xml:"version,attr"`
	Documento Documento  `xml:"Documento"`
	Signature *Signature `xml:"Signature,omitempty"`
}

// Documento representa el contenido del DTE
type Documento struct {
	XMLName     xml.Name     `xml:"Documento"`
	ID          string       `xml:"ID,attr"`
	Encabezado  Encabezado   `xml:"Encabezado"`
	Detalle     []Detalle    `xml:"Detalle"`
	Referencias []Referencia `xml:"Referencia,omitempty"`
	TED         *TED         `xml:"TED,omitempty"`
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
	TipoDTE      int       `xml:"TipoDTE"`
	Folio        int64     `xml:"Folio"`
	FechaEmision time.Time `xml:"FechaEmision"`
	TipoDespacho *int      `xml:"TipoDespacho,omitempty"`
	IndTraslado  *int      `xml:"IndTraslado,omitempty"`
}

// Emisor contiene la información del emisor del DTE
type Emisor struct {
	RUTEmisor    string `xml:"RUTEmisor"`
	RznSoc       string `xml:"RznSoc"`
	GiroEmis     string `xml:"GiroEmis"`
	Acteco       int    `xml:"Acteco"`
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
	MntNeto   int64 `xml:"MntNeto,omitempty"`
	TasaIVA   int   `xml:"TasaIVA,omitempty"`
	IVA       int64 `xml:"IVA,omitempty"`
	MntTotal  int64 `xml:"MntTotal"`
	MntExento int64 `xml:"MntExento,omitempty"`
}

// Detalle representa una línea de detalle del DTE
type Detalle struct {
	NroLinDet int      `xml:"NroLinDet"`
	CdgItem   *CdgItem `xml:"CdgItem,omitempty"`
	NmbItem   string   `xml:"NmbItem"`
	DscItem   string   `xml:"DscItem,omitempty"`
	QtyItem   float64  `xml:"QtyItem,omitempty"`
	UnmdItem  string   `xml:"UnmdItem,omitempty"`
	PrcItem   float64  `xml:"PrcItem,omitempty"`
	MontoItem int64    `xml:"MontoItem"`
	IndExe    *int     `xml:"IndExe,omitempty"`
}

// CdgItem representa el código del ítem
type CdgItem struct {
	TpoCodigo string `xml:"TpoCodigo"`
	VlrCodigo string `xml:"VlrCodigo"`
}

// Referencia representa una referencia a otro documento
type Referencia struct {
	NroLinRef int       `xml:"NroLinRef"`
	TpoDocRef string    `xml:"TpoDocRef"`
	FolioRef  string    `xml:"FolioRef"`
	FechaRef  time.Time `xml:"FechaRef"`
	CodRef    *int      `xml:"CodRef,omitempty"`
	RazonRef  string    `xml:"RazonRef,omitempty"`
}

// EnvioDTE representa el sobre de envío de DTE
type EnvioDTE struct {
	XMLName xml.Name `xml:"EnvioDTE"`
	SetDTE  SetDTE   `xml:"SetDTE"`
	Version string   `xml:"version,attr"`
}

// SetDTE representa un conjunto de DTE
type SetDTE struct {
	ID       string   `xml:"ID,attr"`
	Caratula Caratula `xml:"Caratula"`
	DTE      DTE      `xml:"DTE"`
}

// Caratula contiene la información de envío
type Caratula struct {
	RutEmisor  string    `xml:"RutEmisor"`
	RutEnvia   string    `xml:"RutEnvia"`
	FechaEnvio time.Time `xml:"FechaEnvio"`
	Version    string    `xml:"version"`
}

// ConsultaDTE representa una consulta de estado de DTE
type ConsultaDTE struct {
	RutEmisor     string    `xml:"RutEmisor"`
	TipoDTE       int       `xml:"TipoDTE"`
	Folio         int64     `xml:"Folio"`
	RutConsulta   string    `xml:"RutConsulta"`
	Token         string    `xml:"Token"`
	FechaConsulta time.Time `xml:"FechaConsulta"`
}

// ConsultaTrackID representa una consulta por TrackID
type ConsultaTrackID struct {
	RutEmpresa    string    `xml:"RutEmpresa"`
	TrackID       string    `xml:"TrackID"`
	Token         string    `xml:"Token"`
	FechaConsulta time.Time `xml:"FechaConsulta"`
}

// EstadoDTE representa el estado de un DTE
type EstadoDTE struct {
	Estado         string    `xml:"Estado"`
	GlosaEstado    string    `xml:"GlosaEstado"`
	NumeroAtencion string    `xml:"NumeroAtencion,omitempty"`
	FechaRecepcion time.Time `xml:"FechaRecepcion,omitempty"`
}

// EstadoEnvio representa el estado de un envío
type EstadoEnvio struct {
	TrackID     string       `xml:"TrackID"`
	Estado      string       `xml:"Estado"`
	GlosaEstado string       `xml:"GlosaEstado"`
	DetalleDTE  []DetalleDTE `xml:"DetalleDTE,omitempty"`
}

// DetalleDTE representa el detalle del estado de un DTE dentro de un envío
type DetalleDTE struct {
	TipoDTE     int    `xml:"TipoDTE"`
	Folio       int64  `xml:"Folio"`
	Estado      string `xml:"Estado"`
	GlosaEstado string `xml:"GlosaEstado"`
}
