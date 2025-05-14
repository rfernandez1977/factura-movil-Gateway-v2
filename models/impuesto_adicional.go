package models

// ImpuestoAdicional representa un impuesto adicional aplicado a un ítem
type ImpuestoAdicional struct {
	Codigo        string  `json:"codigo" bson:"codigo"`
	Nombre        string  `json:"nombre" bson:"nombre"`
	Porcentaje    float64 `json:"porcentaje" bson:"porcentaje"`
	MontoImpuesto float64 `json:"monto_impuesto" bson:"monto_impuesto"`
	BaseImponible float64 `json:"base_imponible" bson:"base_imponible"`
	Descripcion   string  `json:"descripcion,omitempty" bson:"descripcion,omitempty"`
}
