package models

import "time"

// CAFDTEXML representa un CAF (Código de Autorización de Folios) en formato XML
type CAFDTEXML struct {
	XMLName struct{}     `xml:"AUTORIZACION"`
	Version string       `xml:"version,attr"`
	DA      DAXMLModel   `xml:"DA"`
	FRMA    FRMAXMLModel `xml:"FRMA"`
}

// DAXMLModel representa los datos de autorización en formato XML
type DAXMLModel struct {
	RUT         RutXMLModel   `xml:"RE"`
	RazonSocial string        `xml:"RS"`
	TipoDTE     string        `xml:"TD"`
	RangoDesde  int           `xml:"RNG>D"`
	RangoHasta  int           `xml:"RNG>H"`
	FechaAut    string        `xml:"FA"`
	RSAPK       RSAPKXMLModel `xml:"RSAPK"`
	IDK         int           `xml:"IDK"`
}

// RutXMLModel representa un RUT en formato XML
type RutXMLModel struct {
	Numero string `xml:",chardata"`
}

// RSAPKXMLModel representa una llave pública RSA en formato XML
type RSAPKXMLModel struct {
	Modulo    string `xml:"M"`
	Exponente string `xml:"E"`
}

// FRMAXMLModel representa la firma del CAF en formato XML
type FRMAXMLModel struct {
	Algoritmo string `xml:"algoritmo,attr"`
	Valor     string `xml:",chardata"`
}

// CAF representa un Código de Autorización de Folios
type CAF struct {
	ID               string    `json:"id" db:"id"`
	EmpresaID        string    `json:"empresa_id" db:"empresa_id"`
	TipoDocumento    string    `json:"tipo_documento" db:"tipo_documento"`
	Desde            int       `json:"desde" db:"desde"`
	Hasta            int       `json:"hasta" db:"hasta"`
	Archivo          []byte    `json:"archivo" db:"archivo"`
	FechaVencimiento string    `json:"fecha_vencimiento" db:"fecha_vencimiento"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

// NewCAF crea una nueva instancia de CAF
func NewCAF(empresaID, tipoDocumento string, desde, hasta int, archivo []byte, fechaVencimiento string) *CAF {
	return &CAF{
		EmpresaID:        empresaID,
		TipoDocumento:    tipoDocumento,
		Desde:            desde,
		Hasta:            hasta,
		Archivo:          archivo,
		FechaVencimiento: fechaVencimiento,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
}

// Validate valida que todos los campos obligatorios estén presentes
func (c *CAF) Validate() error {
	if c.EmpresaID == "" {
		return &ValidationError{Field: "empresa_id", Message: "El ID de la empresa es obligatorio"}
	}
	if c.TipoDocumento == "" {
		return &ValidationError{Field: "tipo_documento", Message: "El tipo de documento es obligatorio"}
	}
	if c.Desde <= 0 {
		return &ValidationError{Field: "desde", Message: "El rango inicial debe ser mayor a cero"}
	}
	if c.Hasta <= 0 {
		return &ValidationError{Field: "hasta", Message: "El rango final debe ser mayor a cero"}
	}
	if c.Hasta < c.Desde {
		return &ValidationError{Field: "hasta", Message: "El rango final debe ser mayor o igual al rango inicial"}
	}
	if c.Archivo == nil || len(c.Archivo) == 0 {
		return &ValidationError{Field: "archivo", Message: "El archivo del CAF es obligatorio"}
	}
	return nil
}

// CAFRequest representa una solicitud de CAF
type CAFRequest struct {
	TipoDocumento string `json:"tipo_documento"`
	RutEmisor     string `json:"rut_emisor"`
	Cantidad      int    `json:"cantidad"`
}
