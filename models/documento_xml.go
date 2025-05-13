package models

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
