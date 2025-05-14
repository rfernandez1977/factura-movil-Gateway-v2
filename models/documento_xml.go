package models

import "time"

// DetalleXML representa el detalle de un documento en formato XML
type DetalleXML struct {
	NroLinDet      int            `xml:"NroLinDet"`
	CdgItem        CdgItem        `xml:"CdgItem,omitempty"`
	NmbItem        string         `xml:"NmbItem"`
	DscItem        string         `xml:"DscItem,omitempty"`
	QtyItem        float64        `xml:"QtyItem"`
	UnmdItem       string         `xml:"UnmdItem,omitempty"`
	PrcItem        float64        `xml:"PrcItem"`
	DescuentoMonto MontoDescuento `xml:"DescuentoMonto,omitempty"`
	MontoItem      int            `xml:"MontoItem"`
}

// CdgItem representa un código de ítem en XML
type CdgItem struct {
	TpoCodigo string `xml:"TpoCodigo"`
	VlrCodigo string `xml:"VlrCodigo"`
}

// MontoDescuento representa un descuento en XML
type MontoDescuento struct {
	TipoDesc  string `xml:"TipoDesc"`
	ValorDesc int    `xml:"ValorDesc"`
}

// DTEXMLModel representa un DTE en formato XML
type DTEXMLModel struct {
	Version   string            `xml:"version,attr"`
	Documento DocumentoXMLModel `xml:"Documento"`
	Signature *FirmaXMLModel    `xml:"Signature,omitempty"`
}

// DocumentoXMLModel representa la estructura de un documento XML
type DocumentoXMLModel struct {
	ID         string             `xml:"ID,attr"`
	Encabezado EncabezadoXMLModel `xml:"Encabezado"`
	Detalle    []DetalleDTEXML    `xml:"Detalle,omitempty"`
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
	RUT         string `xml:"RUTEmisor"`
	RazonSocial string `xml:"RznSoc"`
	Giro        string `xml:"GiroEmis"`
	Direccion   string `xml:"DirOrigen"`
	Comuna      string `xml:"CmnaOrigen"`
	Ciudad      string `xml:"CiudadOrigen"`
}

// ReceptorXMLModel representa al receptor en formato XML
type ReceptorXMLModel struct {
	RUT         string `xml:"RUTRecep"`
	RazonSocial string `xml:"RznSocRecep"`
	Giro        string `xml:"GiroRecep,omitempty"`
	Direccion   string `xml:"DirRecep"`
	Comuna      string `xml:"CmnaRecep"`
	Ciudad      string `xml:"CiudadRecep"`
}

// TotalesXMLModel representa los totales en formato XML
type TotalesXMLModel struct {
	MntNeto  *int64   `xml:"MntNeto,omitempty"`
	TasaIVA  *float64 `xml:"TasaIVA,omitempty"`
	IVA      *int64   `xml:"IVA,omitempty"`
	MntTotal int64    `xml:"MntTotal"`
}

// DetalleDTEXML representa un detalle de DTE en XML
type DetalleDTEXML struct {
	NroLinDet   int      `xml:"NroLinDet"`
	Nombre      string   `xml:"NmbItem"`
	Descripcion *string  `xml:"DscItem,omitempty"`
	Cantidad    *float64 `xml:"QtyItem,omitempty"`
	Precio      *float64 `xml:"PrcItem,omitempty"`
	MontoItem   int64    `xml:"MontoItem"`
}

// FirmaXMLModel representa la firma electrónica en formato XML
type FirmaXMLModel struct {
	SignedInfo struct {
		CanonicalizationMethod struct {
			Algorithm string `xml:"Algorithm,attr"`
		} `xml:"CanonicalizationMethod"`
		SignatureMethod struct {
			Algorithm string `xml:"Algorithm,attr"`
		} `xml:"SignatureMethod"`
		Reference struct {
			URI          string `xml:"URI,attr"`
			DigestMethod struct {
				Algorithm string `xml:"Algorithm,attr"`
			} `xml:"DigestMethod"`
			DigestValue string `xml:"DigestValue"`
		} `xml:"Reference"`
	} `xml:"SignedInfo"`
	SignatureValue string `xml:"SignatureValue"`
	KeyInfo        struct {
		X509Data struct {
			X509Certificate string `xml:"X509Certificate"`
		} `xml:"X509Data"`
	} `xml:"KeyInfo"`
}

// ErrorSII representa un error del SII en la respuesta
type ErrorSII struct {
	Codigo      string `xml:"Codigo" json:"codigo"`
	Descripcion string `xml:"Descripcion" json:"descripcion"`
	Detalle     string `xml:"Detalle" json:"detalle"`
}

// RespuestaSII representa la respuesta del SII a una consulta o envío
type RespuestaSII struct {
	TrackID      string     `xml:"TRACKID" json:"track_id"`
	Estado       string     `xml:"ESTADO" json:"estado"`
	Glosa        string     `xml:"GLOSA" json:"glosa"`
	NumAtencion  string     `xml:"NUMATENCION,omitempty" json:"num_atencion,omitempty"`
	FechaProceso time.Time  `xml:"FECHA_PROCESO" json:"fecha_proceso"`
	Errores      []ErrorSII `xml:"ERRORES>ERROR,omitempty" json:"errores,omitempty"`
}
