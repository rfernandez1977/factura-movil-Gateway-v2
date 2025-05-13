package models

import "time"

// BoletaElectronica representa una boleta electrónica
type BoletaElectronica struct {
	DocumentoTributario
	CAF               *CAF   `json:"caf,omitempty" bson:"caf,omitempty"`
	TimbreElectronico string `json:"timbre_electronico,omitempty" bson:"timbre_electronico,omitempty"`
	FirmaElectronica  string `json:"firma_electronica,omitempty" bson:"firma_electronica,omitempty"`
	// Campos adicionales requeridos por otros paquetes
	MontoNeto            float64             `json:"monto_neto" bson:"monto_neto"`
	MontoIVA             float64             `json:"monto_iva" bson:"monto_iva"`
	MontoExento          float64             `json:"monto_exento" bson:"monto_exento"`
	Items                []Item              `json:"items" bson:"items"`
	Referencias          []Referencia        `json:"referencias" bson:"referencias"`
	ImpuestosAdicionales []ImpuestoAdicional `json:"impuestos_adicionales,omitempty" bson:"impuestos_adicionales,omitempty"`
	Vendedor             string              `json:"vendedor" bson:"vendedor"`
}

// SolicitudBoleta representa una solicitud de boleta
type SolicitudBoleta struct {
	TipoDTE           TipoDTE   `json:"tipo_dte"`
	Folio             int       `json:"folio"`
	FechaEmision      time.Time `json:"fecha_emision"`
	RutEmisor         string    `json:"rut_emisor"`
	RazonSocialEmisor string    `json:"razon_social_emisor"`
	RutReceptor       string    `json:"rut_receptor"`
	RazonSocial       string    `json:"razon_social"`
	Direccion         string    `json:"direccion"`
	Comuna            string    `json:"comuna"`
	Ciudad            string    `json:"ciudad"`
	Giro              string    `json:"giro"`
	Items             []Item    `json:"items"`
	MontoNeto         float64   `json:"monto_neto"`
	MontoIVA          float64   `json:"monto_iva"`
	MontoTotal        float64   `json:"monto_total"`
}

// Detalle representa un ítem en una boleta
type Detalle struct {
	Descripcion    string  `json:"descripcion" binding:"required"`
	Cantidad       int     `json:"cantidad" binding:"required"`
	PrecioUnitario float64 `json:"precio_unitario" binding:"required"`
	MontoItem      float64 `json:"monto_item" binding:"required"`
}

// BoletaResponse representa la respuesta de una boleta
type BoletaResponse struct {
	TrackID string `json:"track_id"`
	Estado  string `json:"estado"`
	Glosa   string `json:"glosa"`
}
