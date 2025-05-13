package models

import "time"

// NotaCredito representa una nota de crédito electrónica
type NotaCredito struct {
	DocumentoTributario
	DocumentoReferencia     *DocumentoTributario `json:"documento_referencia,omitempty" bson:"documento_referencia,omitempty"`
	DocumentoRef            string               `json:"documento_ref" bson:"documento_ref"`
	FolioReferencia         int64                `json:"folio_referencia" bson:"folio_referencia"`
	FechaReferencia         time.Time            `json:"fecha_referencia" bson:"fecha_referencia"`
	TipoReferencia          string               `json:"tipo_referencia" bson:"tipo_referencia"`
	Motivo                  string               `json:"motivo" bson:"motivo"`
	TipoDocumentoReferencia string               `json:"tipo_documento_referencia" bson:"tipo_documento_referencia"`
	RazonReferencia         string               `json:"razon_referencia" bson:"razon_referencia"`
	IndicadorServicio       string               `json:"indicador_servicio,omitempty" bson:"indicador_servicio,omitempty"`
	IndicadorVentas         string               `json:"indicador_ventas,omitempty" bson:"indicador_ventas,omitempty"`
	IndicadorExportacion    string               `json:"indicador_exportacion,omitempty" bson:"indicador_exportacion,omitempty"`
	FechaInicioServicio     time.Time            `json:"fecha_inicio_servicio,omitempty" bson:"fecha_inicio_servicio,omitempty"`
	FechaFinServicio        time.Time            `json:"fecha_fin_servicio,omitempty" bson:"fecha_fin_servicio,omitempty"`
	Periodicidad            string               `json:"periodicidad,omitempty" bson:"periodicidad,omitempty"`
	CodigoAduana            string               `json:"codigo_aduana,omitempty" bson:"codigo_aduana,omitempty"`
	CAF                     *CAF                 `json:"caf,omitempty" bson:"caf,omitempty"`
	TimbreElectronico       string               `json:"timbre_electronico,omitempty" bson:"timbre_electronico,omitempty"`
	FirmaElectronica        string               `json:"firma_electronica,omitempty" bson:"firma_electronica,omitempty"`
	MontoNeto               float64              `json:"monto_neto" bson:"monto_neto"`
	MontoExento             float64              `json:"monto_exento" bson:"monto_exento"`
	MontoIVA                float64              `json:"monto_iva" bson:"monto_iva"`
	Items                   []Item               `json:"items" bson:"items"`
}

// NotaCreditoRequest representa la solicitud para crear una nota de crédito
type NotaCreditoRequest struct {
	TipoDTE                 TipoDTE   `json:"tipo_dte"`
	Folio                   int       `json:"folio"`
	FechaEmision            time.Time `json:"fecha_emision"`
	RutEmisor               string    `json:"rut_emisor"`
	RazonSocialEmisor       string    `json:"razon_social_emisor"`
	RutReceptor             string    `json:"rut_receptor"`
	RazonSocialReceptor     string    `json:"razon_social_receptor"`
	TipoDocumentoReferencia string    `json:"tipo_documento_referencia"`
	FolioReferencia         int64     `json:"folio_referencia"`
	FechaReferencia         time.Time `json:"fecha_referencia"`
	RazonReferencia         string    `json:"razon_referencia"`
	IndicadorServicio       string    `json:"indicador_servicio,omitempty"`
	IndicadorVentas         string    `json:"indicador_ventas,omitempty"`
	IndicadorExportacion    string    `json:"indicador_exportacion,omitempty"`
	FechaInicioServicio     time.Time `json:"fecha_inicio_servicio,omitempty"`
	FechaFinServicio        time.Time `json:"fecha_fin_servicio,omitempty"`
	Periodicidad            string    `json:"periodicidad,omitempty"`
	CodigoAduana            string    `json:"codigo_aduana,omitempty"`
	Items                   []Item    `json:"items"`
	MontoTotal              float64   `json:"monto_total"`
	MontoNeto               float64   `json:"monto_neto"`
	MontoIVA                float64   `json:"monto_iva"`
}

// NotaCreditoResponse representa la respuesta de una nota de crédito
type NotaCreditoResponse struct {
	TrackID string `json:"track_id"`
	Estado  string `json:"estado"`
	Glosa   string `json:"glosa"`
}
