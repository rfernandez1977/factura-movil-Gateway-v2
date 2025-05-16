package validator

import (
	"fmt"

	"github.com/fmgo/core/dte/types"
	"github.com/fmgo/models"
	"github.com/fmgo/utils/validation"
)

// Constantes para los tipos de documentos
const (
	FacturaElectronica = "33"
	BoletaElectronica  = "39"
)

// ValidateDTE valida un DTE completo
func ValidateDTE(d *types.DTE) error {
	// Validaciones básicas
	if d.ID == "" {
		return models.NewValidationFieldError("ID", "REQUIRED_FIELD", "no puede estar vacío", nil)
	}

	if !d.Firmado && d.XMLFirmado != "" {
		return models.NewValidationFieldError("Firmado", "INCONSISTENT_STATE", "documento no firmado pero tiene XML firmado", nil)
	}

	return ValidateDocumento(&d.Documento)
}

// ValidateDocumento valida el documento
func ValidateDocumento(d *types.Documento) error {
	if err := ValidateEncabezado(&d.Encabezado); err != nil {
		return err
	}

	if len(d.Detalles) == 0 {
		return models.NewValidationFieldError("Detalles", "REQUIRED_FIELD", "debe contener al menos un ítem", nil)
	}

	for i, det := range d.Detalles {
		if err := ValidateDetalle(&det); err != nil {
			return models.NewValidationFieldError(
				fmt.Sprintf("Detalles[%d]", i),
				"VALIDATION_ERROR",
				err.Error(),
				det,
			)
		}
	}

	return nil
}

// ValidateEncabezado valida el encabezado
func ValidateEncabezado(e *types.Encabezado) error {
	if err := ValidateIDDocumento(&e.IDDocumento); err != nil {
		return err
	}
	if err := ValidateEmisor(&e.Emisor); err != nil {
		return err
	}
	if err := ValidateReceptor(&e.Receptor); err != nil {
		return err
	}
	if err := ValidateTotales(&e.Totales); err != nil {
		return err
	}
	return nil
}

// ValidateIDDocumento valida el ID del documento
func ValidateIDDocumento(id *types.IDDocumento) error {
	if id.TipoDTE != FacturaElectronica && id.TipoDTE != BoletaElectronica {
		return models.NewValidationFieldError("TipoDTE", "INVALID_VALUE", "tipo de documento inválido", id.TipoDTE)
	}
	if id.Folio <= 0 {
		return models.NewValidationFieldError("Folio", "INVALID_VALUE", "debe ser mayor que 0", id.Folio)
	}
	if id.FechaEmision.IsZero() {
		return models.NewValidationFieldError("FechaEmision", "REQUIRED_FIELD", "no puede estar vacía", nil)
	}
	return nil
}

// ValidateEmisor valida el emisor
func ValidateEmisor(e *types.Emisor) error {
	if err := validation.ValidateRUT(e.RUT); err != nil {
		return models.NewValidationFieldError("RUT", "VALIDATION_ERROR", err.Error(), e.RUT)
	}
	if e.RazonSocial == "" {
		return models.NewValidationFieldError("RazonSocial", "REQUIRED_FIELD", "no puede estar vacío", nil)
	}
	if e.Giro == "" {
		return models.NewValidationFieldError("Giro", "REQUIRED_FIELD", "no puede estar vacío", nil)
	}
	if e.Direccion == "" {
		return models.NewValidationFieldError("Direccion", "REQUIRED_FIELD", "no puede estar vacía", nil)
	}
	if e.Comuna == "" {
		return models.NewValidationFieldError("Comuna", "REQUIRED_FIELD", "no puede estar vacía", nil)
	}
	if e.Ciudad == "" {
		return models.NewValidationFieldError("Ciudad", "REQUIRED_FIELD", "no puede estar vacía", nil)
	}
	if e.Email != "" {
		if err := validation.ValidateEmail(e.Email); err != nil {
			return models.NewValidationFieldError("Email", "VALIDATION_ERROR", err.Error(), e.Email)
		}
	}
	return nil
}

// ValidateReceptor valida el receptor
func ValidateReceptor(r *types.Receptor) error {
	if err := validation.ValidateRUT(r.RUT); err != nil {
		return models.NewValidationFieldError("RUT", "VALIDATION_ERROR", err.Error(), r.RUT)
	}
	if r.RazonSocial == "" {
		return models.NewValidationFieldError("RazonSocial", "REQUIRED_FIELD", "no puede estar vacío", nil)
	}
	if r.Giro == "" {
		return models.NewValidationFieldError("Giro", "REQUIRED_FIELD", "no puede estar vacío", nil)
	}
	if r.Direccion == "" {
		return models.NewValidationFieldError("Direccion", "REQUIRED_FIELD", "no puede estar vacía", nil)
	}
	if r.Comuna == "" {
		return models.NewValidationFieldError("Comuna", "REQUIRED_FIELD", "no puede estar vacía", nil)
	}
	if r.Ciudad == "" {
		return models.NewValidationFieldError("Ciudad", "REQUIRED_FIELD", "no puede estar vacía", nil)
	}
	return nil
}

// ValidateTotales valida los totales
func ValidateTotales(t *types.Totales) error {
	if t.MontoNeto < 0 {
		return models.NewValidationFieldError("MontoNeto", "INVALID_VALUE", "no puede ser negativo", t.MontoNeto)
	}
	if t.TasaIVA != 19 {
		return models.NewValidationFieldError("TasaIVA", "INVALID_VALUE", "debe ser 19", t.TasaIVA)
	}
	if t.IVA < 0 {
		return models.NewValidationFieldError("IVA", "INVALID_VALUE", "no puede ser negativo", t.IVA)
	}
	if t.MontoTotal < 0 {
		return models.NewValidationFieldError("MontoTotal", "INVALID_VALUE", "no puede ser negativo", t.MontoTotal)
	}
	// Validar que el total sea igual a neto + IVA
	expectedTotal := t.MontoNeto + t.IVA
	if int(t.MontoTotal) != int(expectedTotal) {
		return models.NewValidationFieldError(
			"MontoTotal",
			"INVALID_VALUE",
			fmt.Sprintf("debe ser igual a MontoNeto + IVA (%.2f)", expectedTotal),
			t.MontoTotal,
		)
	}
	return nil
}

// ValidateDetalle valida un detalle
func ValidateDetalle(d *types.Detalle) error {
	if d.NumeroLinea <= 0 {
		return models.NewValidationFieldError("NumeroLinea", "INVALID_VALUE", "debe ser mayor que 0", d.NumeroLinea)
	}
	if d.Nombre == "" {
		return models.NewValidationFieldError("Nombre", "REQUIRED_FIELD", "no puede estar vacío", nil)
	}
	if d.Cantidad <= 0 {
		return models.NewValidationFieldError("Cantidad", "INVALID_VALUE", "debe ser mayor que 0", d.Cantidad)
	}
	if d.Precio <= 0 {
		return models.NewValidationFieldError("Precio", "INVALID_VALUE", "debe ser mayor que 0", d.Precio)
	}
	if d.MontoItem <= 0 {
		return models.NewValidationFieldError("MontoItem", "INVALID_VALUE", "debe ser mayor que 0", d.MontoItem)
	}
	// Validar que el monto sea igual a cantidad * precio
	expectedMonto := int(d.Cantidad * d.Precio)
	if int(d.MontoItem) != expectedMonto {
		return models.NewValidationFieldError(
			"MontoItem",
			"INVALID_VALUE",
			fmt.Sprintf("debe ser igual a Cantidad * Precio (%.2f)", d.Cantidad*d.Precio),
			d.MontoItem,
		)
	}
	return nil
}
