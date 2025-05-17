package models

import "encoding/xml"

// DTEXMLModel representa un DTE en formato XML
type DTEXMLModel struct {
	XMLName   xml.Name           `xml:"DTE"`
	Version   string             `xml:"version,attr"`
	Documento *DocumentoXMLModel `xml:"Documento"`
	Signature *FirmaXMLModel     `xml:"Signature,omitempty"`
}

// DocumentoXMLModel representa la estructura de un documento XML
type DocumentoXMLModel struct {
	ID         string             `xml:"ID,attr"`
	Encabezado EncabezadoXMLModel `xml:"Encabezado"`
	Detalle    []DetalleDTEXML    `xml:"Detalle"`
}

// EncabezadoXMLModel representa el encabezado de un DTE en formato XML
type EncabezadoXMLModel struct {
	IdDoc    IDDocumentoXMLModel `xml:"IdDoc"`
	Emisor   EmisorXMLModel      `xml:"Emisor"`
	Receptor ReceptorXMLModel    `xml:"Receptor"`
	Totales  TotalesXMLModel     `xml:"Totales"`
}

// IDDocumentoXMLModel representa la identificación del documento en XML
type IDDocumentoXMLModel struct {
	TipoDTE      string `xml:"TipoDTE"`
	Folio        int    `xml:"Folio"`
	FechaEmision string `xml:"FechaEmision"`
}

// EmisorXMLModel representa al emisor en formato XML
type EmisorXMLModel struct {
	RUT         string `xml:"RUT"`
	RazonSocial string `xml:"RazonSocial"`
	Giro        string `xml:"Giro"`
	Direccion   string `xml:"Direccion"`
	Comuna      string `xml:"Comuna"`
	Ciudad      string `xml:"Ciudad"`
}

// ReceptorXMLModel representa al receptor en formato XML
type ReceptorXMLModel struct {
	RUT         string `xml:"RUT"`
	RazonSocial string `xml:"RazonSocial"`
	Giro        string `xml:"Giro"`
	Direccion   string `xml:"Direccion"`
	Comuna      string `xml:"Comuna"`
	Ciudad      string `xml:"Ciudad"`
}

// TotalesXMLModel representa los totales en formato XML
type TotalesXMLModel struct {
	MntNeto  *int64 `xml:"MntNeto,omitempty"`
	TasaIVA  *int   `xml:"TasaIVA,omitempty"`
	IVA      *int64 `xml:"IVA,omitempty"`
	MntTotal int64  `xml:"MntTotal"`
}

// DetalleDTEXML representa un detalle de DTE en XML
type DetalleDTEXML struct {
	NroLinDet      int      `xml:"NroLinDet"`
	CdgItem        string   `xml:"CdgItem,omitempty"`
	IndExe         *int     `xml:"IndExe,omitempty"`
	Nombre         string   `xml:"NmbItem"`
	Cantidad       float64  `xml:"QtyItem"`
	Unidad         string   `xml:"UnmdItem,omitempty"`
	PrecioUnit     float64  `xml:"PrcItem"`
	DescuentoMonto *float64 `xml:"DescuentoMonto,omitempty"`
	MontoItem      float64  `xml:"MontoItem"`
}

// FirmaXMLModel representa la firma electrónica en formato XML
type FirmaXMLModel struct {
	XMLName        xml.Name `xml:"Signature"`
	SignatureValue string   `xml:"SignatureValue"`
	KeyInfo        struct {
		X509Data struct {
			X509Certificate string `xml:"X509Certificate"`
		} `xml:"X509Data"`
	} `xml:"KeyInfo"`
}
