package models

import (
	"time"
)

// Boleta representa una boleta electrónica
type Boleta struct {
	ID                  string           `json:"id" bson:"_id,omitempty"`
	TrackID             string           `json:"track_id,omitempty" bson:"track_id,omitempty"`
	Folio               int              `json:"folio" bson:"folio"`
	FechaEmision        time.Time        `json:"fecha_emision" bson:"fecha_emision"`
	TipoDocumento       TipoDTE          `json:"tipo_documento" bson:"tipo_documento"`
	RUTEmisor           string           `json:"rut_emisor" bson:"rut_emisor"`
	RazonSocialEmisor   string           `json:"razon_social_emisor" bson:"razon_social_emisor"`
	GiroEmisor          string           `json:"giro_emisor" bson:"giro_emisor"`
	DireccionEmisor     string           `json:"direccion_emisor" bson:"direccion_emisor"`
	ComunaEmisor        string           `json:"comuna_emisor" bson:"comuna_emisor"`
	RUTReceptor         string           `json:"rut_receptor,omitempty" bson:"rut_receptor,omitempty"`
	RazonSocialReceptor string           `json:"razon_social_receptor,omitempty" bson:"razon_social_receptor,omitempty"`
	DireccionReceptor   string           `json:"direccion_receptor,omitempty" bson:"direccion_receptor,omitempty"`
	MontoNeto           float64          `json:"monto_neto" bson:"monto_neto"`
	MontoExento         float64          `json:"monto_exento" bson:"monto_exento"`
	MontoIVA            float64          `json:"monto_iva" bson:"monto_iva"`
	TasaIVA             float64          `json:"tasa_iva" bson:"tasa_iva"`
	MontoTotal          float64          `json:"monto_total" bson:"monto_total"`
	Items               []*DetalleBoleta `json:"items" bson:"items"`
	Referencias         []Referencia     `json:"referencias,omitempty" bson:"referencias,omitempty"`
	Estado              string           `json:"estado" bson:"estado"`
	EstadoSII           string           `json:"estado_sii"`
	Detalles            []*DetalleBoleta `json:"detalles,omitempty"`
	CreatedAt           time.Time        `json:"created_at" bson:"created_at"`
	UpdatedAt           time.Time        `json:"updated_at" bson:"updated_at"`
}

// DetalleBoleta representa un detalle de boleta
type DetalleBoleta struct {
	ID          string  `json:"id"`
	BoletaID    string  `json:"boleta_id"`
	Descripcion string  `json:"descripcion"`
	Cantidad    int     `json:"cantidad"`
	Precio      float64 `json:"precio"`
	Total       float64 `json:"total"`
}

// BoletaRequest representa una solicitud de creación de boleta
type BoletaRequest struct {
	RutEmisor   string            `json:"rut_emisor" binding:"required"`
	RutReceptor string            `json:"rut_receptor" binding:"required"`
	MontoNeto   float64           `json:"monto_neto" binding:"required,gte=0"`
	MontoExento float64           `json:"monto_exento" binding:"gte=0"`
	Detalles    []*DetalleRequest `json:"detalles" binding:"required,min=1"`
}

// DetalleRequest representa un detalle en la solicitud de boleta
type DetalleRequest struct {
	Descripcion string  `json:"descripcion" binding:"required"`
	Cantidad    int     `json:"cantidad" binding:"required,gt=0"`
	Precio      float64 `json:"precio" binding:"required,gte=0"`
	Exento      bool    `json:"exento"`
}

// EstadoDocumentoSII representa el estado de un documento en el SII
type EstadoDocumentoSII struct {
	TrackID         string    `json:"track_id"`
	Estado          string    `json:"estado"`
	Glosa           string    `json:"glosa"`
	NumeroAtencion  string    `json:"numero_atencion,omitempty"`
	FechaRecepcion  time.Time `json:"fecha_recepcion"`
	FechaAceptacion time.Time `json:"fecha_aceptacion,omitempty"`
	FechaRechazo    time.Time `json:"fecha_rechazo,omitempty"`
	DetalleError    string    `json:"detalle_error,omitempty"`
}
