package models

// ImpuestoAdicional representa un impuesto adicional
type ImpuestoAdicional struct {
	Tipo        string  `json:"tipo" bson:"tipo"`
	Codigo      string  `json:"codigo" bson:"codigo"`
	Tasa        float64 `json:"tasa" bson:"tasa"`
	Monto       float64 `json:"monto" bson:"monto"`
	Descripcion string  `json:"descripcion,omitempty" bson:"descripcion,omitempty"`
}
