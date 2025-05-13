package models

import (
	"time"
)

// Boleta representa una boleta electrónica
type Boleta struct {
	ID                  string           `json:"id"`
	TrackID             string           `json:"track_id"`
	Folio               int              `json:"folio"`
	FechaEmision        time.Time        `json:"fecha_emision"`
	MontoNeto           float64          `json:"monto_neto"`
	MontoIVA            float64          `json:"monto_iva"`
	MontoTotal          float64          `json:"monto_total"`
	RutEmisor           string           `json:"rut_emisor"`
	RutReceptor         string           `json:"rut_receptor"`
	RazonSocialEmisor   string           `json:"razon_social_emisor"`
	RazonSocialReceptor string           `json:"razon_social_receptor"`
	DireccionEmisor     string           `json:"direccion_emisor"`
	DireccionReceptor   string           `json:"direccion_receptor"`
	Estado              string           `json:"estado"`
	EstadoSII           string           `json:"estado_sii"`
	Detalles            []*DetalleBoleta `json:"detalles,omitempty"`
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
	Detalles    []*DetalleRequest `json:"detalles" binding:"required,min=1"`
}

// DetalleRequest representa un detalle en la solicitud de boleta
type DetalleRequest struct {
	Descripcion string  `json:"descripcion" binding:"required"`
	Cantidad    int     `json:"cantidad" binding:"required,gt=0"`
	Precio      float64 `json:"precio" binding:"required,gte=0"`
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
