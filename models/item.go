package models

// Item representa un ítem o línea de un documento tributario
type Item struct {
	ID                   string                  `json:"id" bson:"id"`
	NumeroLinea          int                     `json:"numero_linea" bson:"numero_linea"`
	Codigo               string                  `json:"codigo" bson:"codigo"`
	Nombre               string                  `json:"nombre" bson:"nombre"`
	Descripcion          string                  `json:"descripcion,omitempty" bson:"descripcion,omitempty"`
	Cantidad             float64                 `json:"cantidad" bson:"cantidad"`
	Unidad               string                  `json:"unidad,omitempty" bson:"unidad,omitempty"`
	PrecioUnitario       float64                 `json:"precio_unitario" bson:"precio_unitario"`
	Descuento            float64                 `json:"descuento,omitempty" bson:"descuento,omitempty"` // Porcentaje
	Subtotal             float64                 `json:"subtotal" bson:"subtotal"`                       // Precio unitario * cantidad - descuento
	MontoItem            float64                 `json:"monto_item" bson:"monto_item"`                   // Alias para Subtotal (compatibilidad)
	Exento               bool                    `json:"exento" bson:"exento"`
	PorcentajeIVA        float64                 `json:"porcentaje_iva" bson:"porcentaje_iva"`
	MontoIVA             float64                 `json:"monto_iva" bson:"monto_iva"`
	ImpuestosAdicionales []ImpuestoAdicionalItem `json:"impuestos_adicionales,omitempty" bson:"impuestos_adicionales,omitempty"`
}

// ItemRequest representa una solicitud para crear un ítem
type ItemRequest struct {
	Descripcion string  `json:"descripcion" binding:"required"`
	Cantidad    float64 `json:"cantidad" binding:"required,gt=0"`
	PrecioUnit  float64 `json:"precio_unit" binding:"required,gt=0"`
}

// ItemResponse representa la respuesta de un ítem
type ItemResponse struct {
	ID          string  `json:"id"`
	Descripcion string  `json:"descripcion"`
	Cantidad    float64 `json:"cantidad"`
	PrecioUnit  float64 `json:"precio_unit"`
	MontoNeto   float64 `json:"monto_neto"`
	MontoIVA    float64 `json:"monto_iva"`
	MontoTotal  float64 `json:"monto_total"`
}
