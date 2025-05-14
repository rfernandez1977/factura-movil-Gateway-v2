package models

import "time"

// DocumentoTributario representa la estructura común para todos los documentos tributarios
type DocumentoTributario struct {
	ID                  string       `json:"id" bson:"_id,omitempty"`
	Folio               int          `json:"folio" bson:"folio"`
	FechaEmision        time.Time    `json:"fecha_emision" bson:"fecha_emision"`
	TipoDocumento       TipoDTE      `json:"tipo_documento" bson:"tipo_documento"`
	RUTEmisor           string       `json:"rut_emisor" bson:"rut_emisor"`
	RazonSocialEmisor   string       `json:"razon_social_emisor" bson:"razon_social_emisor"`
	GiroEmisor          string       `json:"giro_emisor" bson:"giro_emisor"`
	DireccionEmisor     string       `json:"direccion_emisor" bson:"direccion_emisor"`
	ComunaEmisor        string       `json:"comuna_emisor" bson:"comuna_emisor"`
	RUTReceptor         string       `json:"rut_receptor" bson:"rut_receptor"`
	RazonSocialReceptor string       `json:"razon_social_receptor" bson:"razon_social_receptor"`
	GiroReceptor        string       `json:"giro_receptor,omitempty" bson:"giro_receptor,omitempty"`
	DireccionReceptor   string       `json:"direccion_receptor" bson:"direccion_receptor"`
	ComunaReceptor      string       `json:"comuna_receptor,omitempty" bson:"comuna_receptor,omitempty"`
	MontoNeto           float64      `json:"monto_neto" bson:"monto_neto"`
	MontoExento         float64      `json:"monto_exento" bson:"monto_exento"`
	MontoIVA            float64      `json:"monto_iva" bson:"monto_iva"`
	TasaIVA             float64      `json:"tasa_iva" bson:"tasa_iva"`
	MontoTotal          float64      `json:"monto_total" bson:"monto_total"`
	Referencias         []Referencia `json:"referencias,omitempty" bson:"referencias,omitempty"`
	Estado              string       `json:"estado" bson:"estado"`
	TrackID             string       `json:"track_id,omitempty" bson:"track_id,omitempty"`
	PDF                 string       `json:"pdf,omitempty" bson:"pdf,omitempty"`
	XML                 string       `json:"xml,omitempty" bson:"xml,omitempty"`
	CreatedAt           time.Time    `json:"created_at" bson:"created_at"`
	UpdatedAt           time.Time    `json:"updated_at" bson:"updated_at"`
}
