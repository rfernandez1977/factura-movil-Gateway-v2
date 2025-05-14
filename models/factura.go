package models

import (
	"time"

	"github.com/cursor/FMgo/domain"
)

// Factura representa una factura electrónica
type Factura struct {
	domain.DocumentoTributario
	ID                  string          `json:"id" bson:"_id"`
	TipoDocumento       TipoDTE         `json:"tipo_documento" bson:"tipo_documento"`
	Folio               int64           `json:"folio" bson:"folio"`
	FechaEmision        time.Time       `json:"fecha_emision" bson:"fecha_emision"`
	FechaVencimiento    time.Time       `json:"fecha_vencimiento" bson:"fecha_vencimiento"`
	RutEmisor           string          `json:"rut_emisor" bson:"rut_emisor"`
	RazonSocialEmisor   string          `json:"razon_social_emisor" bson:"razon_social_emisor"`
	RutReceptor         string          `json:"rut_receptor" bson:"rut_receptor"`
	RazonSocialReceptor string          `json:"razon_social_receptor" bson:"razon_social_receptor"`
	MontoTotal          float64         `json:"monto_total" bson:"monto_total"`
	MontoNeto           float64         `json:"monto_neto" bson:"monto_neto"`
	MontoExento         float64         `json:"monto_exento" bson:"monto_exento"`
	MontoIVA            float64         `json:"monto_iva" bson:"monto_iva"`
	FormaPago           string          `json:"forma_pago" bson:"forma_pago"`
	Vencimiento         int             `json:"vencimiento" bson:"vencimiento"`
	Estado              EstadoDocumento `json:"estado" bson:"estado"`
	Items               []domain.Item   `json:"items" bson:"items"`
	FechaCreacion       time.Time       `json:"fecha_creacion" bson:"fecha_creacion"`
	FechaActualizacion  time.Time       `json:"fecha_actualizacion" bson:"fecha_actualizacion"`
	CAF                 *domain.CAF     `json:"caf,omitempty" bson:"caf,omitempty"`
	TimbreElectronico   string          `json:"timbre_electronico,omitempty" bson:"timbre_electronico,omitempty"`
	FirmaElectronica    string          `json:"firma_electronica,omitempty" bson:"firma_electronica,omitempty"`
	Referencias         []Referencia    `json:"referencias,omitempty" bson:"referencias,omitempty"`
}

// FacturaRequest representa la solicitud para crear una factura
type FacturaRequest struct {
	RutEmisor        string        `json:"rut_emisor" binding:"required"`
	RutReceptor      string        `json:"rut_receptor" binding:"required"`
	FechaEmision     time.Time     `json:"fecha_emision" binding:"required"`
	FechaVencimiento time.Time     `json:"fecha_vencimiento"`
	FormaPago        string        `json:"forma_pago" binding:"required"`
	Vencimiento      int           `json:"vencimiento"`
	Items            []domain.Item `json:"items" binding:"required,min=1"`
}

// FacturaResponse representa la respuesta de una factura
type FacturaResponse struct {
	ID                  string        `json:"id"`
	TipoDocumento       string        `json:"tipo_documento"`
	Folio               int64         `json:"folio"`
	FechaEmision        time.Time     `json:"fecha_emision"`
	FechaVencimiento    time.Time     `json:"fecha_vencimiento"`
	RutEmisor           string        `json:"rut_emisor"`
	RazonSocialEmisor   string        `json:"razon_social_emisor"`
	RutReceptor         string        `json:"rut_receptor"`
	RazonSocialReceptor string        `json:"razon_social_receptor"`
	MontoTotal          float64       `json:"monto_total"`
	MontoNeto           float64       `json:"monto_neto"`
	MontoExento         float64       `json:"monto_exento"`
	MontoIVA            float64       `json:"monto_iva"`
	FormaPago           string        `json:"forma_pago"`
	Vencimiento         int           `json:"vencimiento"`
	Estado              string        `json:"estado"`
	Items               []domain.Item `json:"items"`
	TimbreElectronico   string        `json:"timbre_electronico,omitempty"`
	FirmaElectronica    string        `json:"firma_electronica,omitempty"`
}

// DetalleFactura representa un detalle de factura
type DetalleFactura struct {
	ID          string  `json:"id" bson:"_id"`
	FacturaID   string  `json:"factura_id" bson:"factura_id"`
	Descripcion string  `json:"descripcion" bson:"descripcion"`
	Cantidad    float64 `json:"cantidad" bson:"cantidad"`
	PrecioUnit  float64 `json:"precio_unit" bson:"precio_unit"`
	MontoTotal  float64 `json:"monto_total" bson:"monto_total"`
}
