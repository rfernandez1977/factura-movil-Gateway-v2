package utils

import (
	"fmt"
	"time"

	"FMgo/models"
)

// ReferenceValidator define la validación de referencias
type ReferenceValidator struct {
	Reference            models.Referencia
	TipoDocumentoDestino string
	FolioDestino         string
	FechaDestino         time.Time
	RUTEmisorDestino     string
}

// NewReferenceValidator crea una nueva instancia de ReferenceValidator
func NewReferenceValidator(reference models.Referencia) *ReferenceValidator {
	return &ReferenceValidator{
		Reference: reference,
	}
}

// Validate valida la referencia
func (v *ReferenceValidator) Validate() error {
	if v.Reference.TipoDocumento == "" {
		return fmt.Errorf("tipo de documento es requerido")
	}

	if v.Reference.Folio <= 0 {
		return fmt.Errorf("folio es requerido")
	}

	if v.Reference.FechaReferencia.IsZero() {
		return fmt.Errorf("fecha de referencia es requerida")
	}

	if v.Reference.RazonReferencia == "" {
		return fmt.Errorf("razón de referencia es requerida")
	}

	// Validaciones específicas por tipo de documento
	// ...

	return nil
}

// ValidateReferenceChain valida la cadena de referencias
func ValidateReferenceChain(validators []ReferenceValidator) error {
	// Validar cada referencia individualmente
	for i, validator := range validators {
		if err := validator.Validate(); err != nil {
			return fmt.Errorf("referencia %d: %v", i+1, err)
		}
	}

	// Validar relaciones entre referencias
	// ...

	return nil
}

// ValidateReferenceTypes valida los tipos de documentos referenciados
func ValidateReferenceTypes(tipoDocumento string, tipoDocumentoReferencia string) error {
	// Validar que el tipo de documento de referencia sea válido para el tipo de documento actual
	// Por ejemplo, una nota de crédito solo puede referenciar a una factura, factura exenta o boleta electrónica

	// ... implementación detallada ...

	return nil
}

// getTipoDocumento convierte un tipo de documento interno a su representación para validaciones
func getTipoDocumento(tipoDTE string) string {
	switch tipoDTE {
	case "33":
		return "FACTURA"
	case "34":
		return "FACTURA_EXENTA"
	case "56":
		return "NOTA_DEBITO"
	case "61":
		return "NOTA_CREDITO"
	case "52":
		return "GUIA_DESPACHO"
	case "39":
		return "BOLETA_ELECTRONICA"
	case "41":
		return "BOLETA_EXENTA"
	default:
		return "DESCONOCIDO"
	}
}
