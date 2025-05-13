package models

// ItemNotaVenta representa un ítem de una nota de venta
type ItemNotaVenta struct {
	Codigo         string  `json:"codigo" bson:"codigo"`
	Descripcion    string  `json:"descripcion" bson:"descripcion"`
	Cantidad       float64 `json:"cantidad" bson:"cantidad"`
	PrecioUnitario float64 `json:"precio_unitario" bson:"precio_unitario"`
	Descuento      float64 `json:"descuento" bson:"descuento"`
	Subtotal       float64 `json:"subtotal" bson:"subtotal"`
}
