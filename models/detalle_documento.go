package models

// DetalleDocumento representa un detalle en un documento tributario
type DetalleDocumento struct {
	Nombre         string  `json:"nombre" bson:"nombre"`
	Cantidad       int     `json:"cantidad" bson:"cantidad"`
	PrecioUnitario float64 `json:"precio_unitario" bson:"precio_unitario"`
	MontoItem      float64 `json:"monto_item" bson:"monto_item"`
	Exento         bool    `json:"exento" bson:"exento,omitempty"`
	Descuento      float64 `json:"descuento" bson:"descuento,omitempty"`
	Codigo         string  `json:"codigo" bson:"codigo,omitempty"`
	Unidad         string  `json:"unidad" bson:"unidad,omitempty"`
}
