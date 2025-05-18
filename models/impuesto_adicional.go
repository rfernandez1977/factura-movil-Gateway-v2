package models

// ImpuestoAdicional representa un impuesto adicional en un documento tributario
type ImpuestoAdicional struct {
	Codigo        string  `json:"codigo"`         // Código del impuesto según catálogo SII
	Porcentaje    float64 `json:"porcentaje"`     // Porcentaje del impuesto
	Monto         float64 `json:"monto"`          // Monto calculado del impuesto
	BaseImponible float64 `json:"base_imponible"` // Base sobre la cual se calcula el impuesto
}
