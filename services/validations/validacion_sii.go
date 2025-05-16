package validations

import (
	"fmt"
	"time"

	"github.com/fmgo/models"
)

// SIIBasicValidator proporciona métodos para validar documentos según las reglas del SII
type SIIBasicValidator struct {
	ambiente string
}

// NewSIIBasicValidator crea una nueva instancia de SIIBasicValidator
func NewSIIBasicValidator(ambiente string) *SIIBasicValidator {
	return &SIIBasicValidator{
		ambiente: ambiente,
	}
}

// ValidarDocumento valida un documento según las reglas del SII
func (v *SIIBasicValidator) ValidarDocumento(doc *models.DocumentoTributario) error {
	// Validar campos obligatorios
	if err := v.validarCamposObligatorios(doc); err != nil {
		return err
	}

	// Validar montos
	if err := v.validarMontos(doc); err != nil {
		return err
	}

	// Validar fechas
	if err := v.validarFechas(doc); err != nil {
		return err
	}

	// Validar referencias
	if err := v.validarReferencias(doc); err != nil {
		return err
	}

	return nil
}

// validarCamposObligatorios valida que todos los campos obligatorios estén presentes
func (v *SIIBasicValidator) validarCamposObligatorios(doc *models.DocumentoTributario) error {
	if doc.RUTEmisor == "" {
		return fmt.Errorf("RUT del emisor es obligatorio")
	}
	if doc.RUTReceptor == "" {
		return fmt.Errorf("RUT del receptor es obligatorio")
	}
	if doc.RazonSocialEmisor == "" {
		return fmt.Errorf("razón social del emisor es obligatoria")
	}
	if doc.RazonSocialReceptor == "" {
		return fmt.Errorf("razón social del receptor es obligatoria")
	}
	if doc.Folio <= 0 {
		return fmt.Errorf("folio debe ser mayor a 0")
	}
	return nil
}

// validarMontos valida que los montos sean correctos
func (v *SIIBasicValidator) validarMontos(doc *models.DocumentoTributario) error {
	// Validar que el monto total sea mayor a 0
	if doc.MontoTotal <= 0 {
		return fmt.Errorf("monto total debe ser mayor a 0")
	}

	return nil
}

// validarFechas valida que las fechas sean correctas
func (v *SIIBasicValidator) validarFechas(doc *models.DocumentoTributario) error {
	// La fecha de emisión no puede ser futura
	if doc.FechaEmision.After(time.Now()) {
		return fmt.Errorf("fecha de emisión no puede ser futura")
	}

	return nil
}

// validarReferencias valida que las referencias sean correctas
func (v *SIIBasicValidator) validarReferencias(doc *models.DocumentoTributario) error {
	if len(doc.Referencias) == 0 {
		return nil
	}

	for _, ref := range doc.Referencias {
		// El folio de la referencia debe ser mayor a 0
		if ref.Folio <= 0 {
			return fmt.Errorf("folio de la referencia debe ser mayor a 0")
		}

		// El tipo de documento de la referencia debe ser válido
		if ref.TipoDocumento == "" {
			return fmt.Errorf("tipo de documento de la referencia es obligatorio")
		}

		// La razón de la referencia debe estar presente
		if ref.RazonReferencia == "" {
			return fmt.Errorf("razón de la referencia es obligatoria")
		}
	}

	return nil
}
