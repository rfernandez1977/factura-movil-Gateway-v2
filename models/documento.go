package models

import (
	"time"
)

// DocumentoTributarioBasico representa un documento tributario electrónico básico
type DocumentoTributarioBasico struct {
	ID           string    `json:"id" bson:"_id"`
	TipoDTE      string    `json:"tipo_dte" bson:"tipo_dte"`
	Folio        int64     `json:"folio" bson:"folio"`
	FechaEmision time.Time `json:"fecha_emision" bson:"fecha_emision"`
	RutEmisor    string    `json:"rut_emisor" bson:"rut_emisor"`
	RutReceptor  string    `json:"rut_receptor" bson:"rut_receptor"`
	MontoTotal   float64   `json:"monto_total" bson:"monto_total"`
	MontoNeto    float64   `json:"monto_neto" bson:"monto_neto"`
	MontoIVA     float64   `json:"monto_iva" bson:"monto_iva"`
	Estado       string    `json:"estado" bson:"estado"`
	TrackID      string    `json:"track_id,omitempty" bson:"track_id,omitempty"`
	XML          string    `json:"xml,omitempty" bson:"xml,omitempty"`
	XMLFirmado   string    `json:"xml_firmado,omitempty" bson:"xml_firmado,omitempty"`
}

// DocumentoRequest representa una solicitud para crear un documento
type DocumentoRequest struct {
	TipoDTE     string  `json:"tipo_dte" binding:"required"`
	RutEmisor   string  `json:"rut_emisor" binding:"required"`
	RutReceptor string  `json:"rut_receptor" binding:"required"`
	MontoNeto   float64 `json:"monto_neto" binding:"required,gte=0"`
	MontoIVA    float64 `json:"monto_iva" binding:"required,gte=0"`
	MontoTotal  float64 `json:"monto_total" binding:"required,gt=0"`
	Items       []Item  `json:"items" binding:"required,min=1"`
}

// DocumentoResponse representa la respuesta de un documento
type DocumentoResponse struct {
	ID           string    `json:"id"`
	TipoDTE      string    `json:"tipo_dte"`
	Folio        int64     `json:"folio"`
	FechaEmision time.Time `json:"fecha_emision"`
	RutEmisor    string    `json:"rut_emisor"`
	RutReceptor  string    `json:"rut_receptor"`
	MontoTotal   float64   `json:"monto_total"`
	MontoNeto    float64   `json:"monto_neto"`
	MontoIVA     float64   `json:"monto_iva"`
	Estado       string    `json:"estado"`
	TrackID      string    `json:"track_id,omitempty"`
}
