package parser

import (
	"encoding/xml"

	"FMgo/core/dte/types"
)

// Parser es la interfaz para parsear DTEs desde diferentes formatos
type Parser interface {
	ParseXML(xmlData []byte) (*types.DTE, error)
	GenerateXML(dte *types.DTE) ([]byte, error)
}

// XMLParser es la implementación por defecto del parser XML
type XMLParser struct{}

// NewXMLParser crea una nueva instancia del parser XML
func NewXMLParser() *XMLParser {
	return &XMLParser{}
}

// ParseXML parsea un DTE desde un documento XML
func (p *XMLParser) ParseXML(xmlData []byte) (*types.DTE, error) {
	var dte types.DTE
	if err := xml.Unmarshal(xmlData, &dte); err != nil {
		return nil, err
	}
	return &dte, nil
}

// GenerateXML genera el XML para un DTE
func (p *XMLParser) GenerateXML(dte *types.DTE) ([]byte, error) {
	return xml.MarshalIndent(dte, "", "  ")
}

// XMLDocument representa la estructura del documento XML
type XMLDocument struct {
	XMLName   xml.Name     `xml:"DTE"`
	Version   string       `xml:"version,attr"`
	Documento XMLDocumento `xml:"Documento"`
}

// XMLDocumento representa la estructura del documento en XML
type XMLDocumento struct {
	XMLName    xml.Name      `xml:"Documento"`
	Encabezado XMLEncabezado `xml:"Encabezado"`
	Detalles   []XMLDetalle  `xml:"Detalle"`
}

// XMLEncabezado representa el encabezado en XML
type XMLEncabezado struct {
	XMLName     xml.Name       `xml:"Encabezado"`
	IDDocumento XMLIDDocumento `xml:"IdDoc"`
	Emisor      XMLEmisor      `xml:"Emisor"`
	Receptor    XMLReceptor    `xml:"Receptor"`
	Totales     XMLTotales     `xml:"Totales"`
}

// XMLIDDocumento representa el ID del documento en XML
type XMLIDDocumento struct {
	XMLName          xml.Name `xml:"IdDoc"`
	TipoDTE          string   `xml:"TipoDTE"`
	Folio            int      `xml:"Folio"`
	FechaEmision     string   `xml:"FechaEmision"`
	FechaVencimiento string   `xml:"FechaVencimiento,omitempty"`
	FormaPago        string   `xml:"FormaPago,omitempty"`
}

// XMLEmisor representa el emisor en XML
type XMLEmisor struct {
	XMLName     xml.Name `xml:"Emisor"`
	RUT         string   `xml:"RUTEmisor"`
	RazonSocial string   `xml:"RznSoc"`
	Giro        string   `xml:"GiroEmis"`
	Direccion   string   `xml:"DirOrigen"`
	Comuna      string   `xml:"CmnaOrigen"`
	Ciudad      string   `xml:"CiudadOrigen"`
}

// XMLReceptor representa el receptor en XML
type XMLReceptor struct {
	XMLName     xml.Name `xml:"Receptor"`
	RUT         string   `xml:"RUTRecep"`
	RazonSocial string   `xml:"RznSocRecep"`
	Giro        string   `xml:"GiroRecep"`
	Direccion   string   `xml:"DirRecep"`
	Comuna      string   `xml:"CmnaRecep"`
	Ciudad      string   `xml:"CiudadRecep"`
}

// XMLTotales representa los totales en XML
type XMLTotales struct {
	XMLName    xml.Name `xml:"Totales"`
	MontoNeto  float64  `xml:"MntNeto"`
	TasaIVA    float64  `xml:"TasaIVA"`
	IVA        float64  `xml:"IVA"`
	MontoTotal float64  `xml:"MntTotal"`
}

// XMLDetalle representa un detalle en XML
type XMLDetalle struct {
	XMLName     xml.Name `xml:"Detalle"`
	NumeroLinea int      `xml:"NroLinDet"`
	Nombre      string   `xml:"NmbItem"`
	Descripcion string   `xml:"DscItem,omitempty"`
	Cantidad    float64  `xml:"QtyItem"`
	Precio      float64  `xml:"PrcItem"`
	MontoItem   float64  `xml:"MontoItem"`
}
