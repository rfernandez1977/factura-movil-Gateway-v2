package dte

import (
	"fmt"

	"github.com/cursor/FMgo/models"
	"github.com/cursor/FMgo/utils/validation"
)

// Constantes para los tipos de documentos
const (
	FacturaElectronica = "33"
	BoletaElectronica  = "39"
)

// Validate valida un DTE completo
func (d *DTE) Validate() error {
	// Validaciones básicas
	if d.ID == "" {
		return models.NewValidationError("ID", "REQUIRED_FIELD", "no puede estar vacío", nil)
	}

	if !d.Firmado && d.XMLFirmado != "" {
		return models.NewValidationError("Firmado", "INCONSISTENT_STATE", "documento no firmado pero tiene XML firmado", nil)
	}

	return d.Documento.Validate()
}

// Validate valida el documento
func (d *Documento) Validate() error {
	if err := d.Encabezado.Validate(); err != nil {
		return err
	}

	if len(d.Detalles) == 0 {
		return models.NewValidationError("Detalles", "REQUIRED_FIELD", "debe contener al menos un ítem", nil)
	}

	for i, det := range d.Detalles {
		if err := validateDetalle(&det); err != nil {
			return models.NewValidationError(
				fmt.Sprintf("Detalles[%d]", i),
				"VALIDATION_ERROR",
				err.Error(),
				det,
			)
		}
	}

	return nil
}

// Validate valida el encabezado
func (e *Encabezado) Validate() error {
	if err := e.IDDocumento.Validate(); err != nil {
		return err
	}
	if err := e.Emisor.Validate(); err != nil {
		return err
	}
	if err := e.Receptor.Validate(); err != nil {
		return err
	}
	if err := e.Totales.Validate(); err != nil {
		return err
	}
	return nil
}

// Validate valida el ID del documento
func (id *IDDocumento) Validate() error {
	if id.TipoDTE != FacturaElectronica && id.TipoDTE != BoletaElectronica {
		return models.NewValidationError("TipoDTE", "INVALID_VALUE", "tipo de documento inválido", id.TipoDTE)
	}
	if id.Folio <= 0 {
		return models.NewValidationError("Folio", "INVALID_VALUE", "debe ser mayor que 0", id.Folio)
	}
	if id.FechaEmision.IsZero() {
		return models.NewValidationError("FechaEmision", "REQUIRED_FIELD", "no puede estar vacía", nil)
	}
	return nil
}

// Validate valida el emisor
func (e *Emisor) Validate() error {
	if err := validation.ValidateRUT(e.RUT); err != nil {
		return models.NewValidationError("RUT", "VALIDATION_ERROR", err.Error(), e.RUT)
	}
	if e.RazonSocial == "" {
		return models.NewValidationError("RazonSocial", "REQUIRED_FIELD", "no puede estar vacío", nil)
	}
	if e.Giro == "" {
		return models.NewValidationError("Giro", "REQUIRED_FIELD", "no puede estar vacío", nil)
	}
	if e.Direccion == "" {
		return models.NewValidationError("Direccion", "REQUIRED_FIELD", "no puede estar vacía", nil)
	}
	if e.Comuna == "" {
		return models.NewValidationError("Comuna", "REQUIRED_FIELD", "no puede estar vacía", nil)
	}
	if e.Ciudad == "" {
		return models.NewValidationError("Ciudad", "REQUIRED_FIELD", "no puede estar vacía", nil)
	}
	if e.Email != "" {
		if err := validation.ValidateEmail(e.Email); err != nil {
			return models.NewValidationError("Email", "VALIDATION_ERROR", err.Error(), e.Email)
		}
	}
	return nil
}

// Validate valida el receptor
func (r *Receptor) Validate() error {
	if err := validation.ValidateRUT(r.RUT); err != nil {
		return models.NewValidationError("RUT", "VALIDATION_ERROR", err.Error(), r.RUT)
	}
	if r.RazonSocial == "" {
		return models.NewValidationError("RazonSocial", "REQUIRED_FIELD", "no puede estar vacío", nil)
	}
	if r.Giro == "" {
		return models.NewValidationError("Giro", "REQUIRED_FIELD", "no puede estar vacío", nil)
	}
	if r.Direccion == "" {
		return models.NewValidationError("Direccion", "REQUIRED_FIELD", "no puede estar vacía", nil)
	}
	if r.Comuna == "" {
		return models.NewValidationError("Comuna", "REQUIRED_FIELD", "no puede estar vacía", nil)
	}
	if r.Ciudad == "" {
		return models.NewValidationError("Ciudad", "REQUIRED_FIELD", "no puede estar vacía", nil)
	}
	return nil
}

// Validate valida los totales
func (t *Totales) Validate() error {
	if t.MontoNeto < 0 {
		return models.NewValidationError("MontoNeto", "INVALID_VALUE", "no puede ser negativo", t.MontoNeto)
	}
	if t.TasaIVA != 19 {
		return models.NewValidationError("TasaIVA", "INVALID_VALUE", "debe ser 19", t.TasaIVA)
	}
	if t.IVA < 0 {
		return models.NewValidationError("IVA", "INVALID_VALUE", "no puede ser negativo", t.IVA)
	}
	if t.MontoTotal < 0 {
		return models.NewValidationError("MontoTotal", "INVALID_VALUE", "no puede ser negativo", t.MontoTotal)
	}
	// Validar que el total sea igual a neto + IVA
	expectedTotal := t.MontoNeto + t.IVA
	if int(t.MontoTotal) != int(expectedTotal) {
		return models.NewValidationError(
			"MontoTotal",
			"INVALID_VALUE",
			fmt.Sprintf("debe ser igual a MontoNeto + IVA (%.2f)", expectedTotal),
			t.MontoTotal,
		)
	}
	return nil
}

// validateDetalle valida un detalle
func validateDetalle(d *Detalle) error {
	if d.NumeroLinea <= 0 {
		return models.NewValidationError("NumeroLinea", "INVALID_VALUE", "debe ser mayor que 0", d.NumeroLinea)
	}
	if d.Nombre == "" {
		return models.NewValidationError("Nombre", "REQUIRED_FIELD", "no puede estar vacío", nil)
	}
	if d.Cantidad <= 0 {
		return models.NewValidationError("Cantidad", "INVALID_VALUE", "debe ser mayor que 0", d.Cantidad)
	}
	if d.Precio <= 0 {
		return models.NewValidationError("Precio", "INVALID_VALUE", "debe ser mayor que 0", d.Precio)
	}
	if d.MontoItem <= 0 {
		return models.NewValidationError("MontoItem", "INVALID_VALUE", "debe ser mayor que 0", d.MontoItem)
	}
	// Validar que el monto sea igual a cantidad * precio
	expectedMonto := int(d.Cantidad * d.Precio)
	if int(d.MontoItem) != expectedMonto {
		return models.NewValidationError(
			"MontoItem",
			"INVALID_VALUE",
			fmt.Sprintf("debe ser igual a Cantidad * Precio (%.2f)", d.Cantidad*d.Precio),
			d.MontoItem,
		)
	}
	return nil
}
