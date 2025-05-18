package validation

import (
	"time"
)

// ValidationConfig contiene la configuración para las validaciones
type ValidationConfig struct {
	// Modo estricto de validación
	StrictMode bool `json:"strict_mode"`

	// Validación de esquemas XSD
	ValidateSchema bool   `json:"validate_schema"`
	SchemaPath     string `json:"schema_path"`

	// Validación de CAF
	ValidateCAF bool   `json:"validate_caf"`
	CAFPath     string `json:"caf_path"`

	// Validación de firma
	ValidateFirma bool `json:"validate_firma"`

	// Timeouts
	ValidationTimeout time.Duration `json:"validation_timeout"`
}

// ValidationRules define las reglas de negocio para validación
type ValidationRules struct {
	// Reglas generales
	MaxMontoTotal     int64 `json:"max_monto_total"`
	MaxLineasDetalle  int   `json:"max_lineas_detalle"`
	MaxFoliosPerCAF   int   `json:"max_folios_per_caf"`
	MaxReferenciasDTE int   `json:"max_referencias_dte"`

	// Reglas por tipo de documento
	ReglasFactura      *ReglasDTE `json:"reglas_factura"`
	ReglasNotaCredito  *ReglasDTE `json:"reglas_nota_credito"`
	ReglasNotaDebito   *ReglasDTE `json:"reglas_nota_debito"`
	ReglasGuiaDespacho *ReglasDTE `json:"reglas_guia_despacho"`
}

// ReglasDTE define reglas específicas por tipo de documento
type ReglasDTE struct {
	MontoMinimo            int64    `json:"monto_minimo"`
	MontoMaximo            int64    `json:"monto_maximo"`
	RequiereReferencia     bool     `json:"requiere_referencia"`
	RequiereDetalle        bool     `json:"requiere_detalle"`
	TiposReferenciaValidos []string `json:"tipos_referencia_validos"`
}

// NewDefaultConfig crea una configuración por defecto
func NewDefaultConfig() *ValidationConfig {
	return &ValidationConfig{
		StrictMode:        true,
		ValidateSchema:    true,
		ValidateCAF:       true,
		ValidateFirma:     true,
		ValidationTimeout: 30 * time.Second,
	}
}

// NewDefaultRules crea reglas de validación por defecto
func NewDefaultRules() *ValidationRules {
	return &ValidationRules{
		MaxMontoTotal:     1000000000, // 1 billón
		MaxLineasDetalle:  60,
		MaxFoliosPerCAF:   100,
		MaxReferenciasDTE: 40,
		ReglasFactura: &ReglasDTE{
			MontoMinimo:            0,
			MontoMaximo:            1000000000,
			RequiereDetalle:        true,
			TiposReferenciaValidos: []string{"801", "802", "803"},
		},
		ReglasNotaCredito: &ReglasDTE{
			RequiereReferencia:     true,
			TiposReferenciaValidos: []string{"33", "34", "56", "61"},
		},
		ReglasNotaDebito: &ReglasDTE{
			RequiereReferencia:     true,
			TiposReferenciaValidos: []string{"33", "34"},
		},
		ReglasGuiaDespacho: &ReglasDTE{
			RequiereDetalle: true,
		},
	}
}
