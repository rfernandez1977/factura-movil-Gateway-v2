package models

// ImpuestoAdicional representa un impuesto adicional aplicado a un documento
type ImpuestoAdicional struct {
	Codigo        string  `json:"codigo" bson:"codigo"`
	Nombre        string  `json:"nombre" bson:"nombre"`
	Porcentaje    float64 `json:"porcentaje" bson:"porcentaje"`
	Monto         float64 `json:"monto" bson:"monto"`
	BaseImponible float64 `json:"base_imponible,omitempty" bson:"base_imponible,omitempty"`
}

// ImpuestoAdicionalItem representa un impuesto adicional aplicado a un ítem
type ImpuestoAdicionalItem struct {
	TipoImpuesto  string  `json:"tipo_impuesto" bson:"tipo_impuesto"`
	Codigo        string  `json:"codigo" bson:"codigo"`
	Nombre        string  `json:"nombre" bson:"nombre"`
	Porcentaje    float64 `json:"porcentaje" bson:"porcentaje"`
	Monto         float64 `json:"monto" bson:"monto"`
	BaseImponible float64 `json:"base_imponible,omitempty" bson:"base_imponible,omitempty"`
}
