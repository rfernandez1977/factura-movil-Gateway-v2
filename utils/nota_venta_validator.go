package utils

import (
	"fmt"
	"time"

	"github.com/cursor/FMgo/models"
)

// NotaVentaValidator valida una nota de venta
type NotaVentaValidator struct {
	RUTEmisor       string
	RUTReceptor     string
	Folio           string
	FechaEmision    time.Time
	MontoNeto       float64
	MontoIVA        float64
	MontoTotal      float64
	TipoNotaVenta   string
	ReferenciaDoc   string
	Items           []models.ItemNotaVenta
	amountValidator *AmountValidator
}

// NewNotaVentaValidator crea una nueva instancia de NotaVentaValidator
func NewNotaVentaValidator() *NotaVentaValidator {
	return &NotaVentaValidator{
		amountValidator: NewAmountValidator(),
	}
}

// Validate valida todos los aspectos de una nota de venta
func (v *NotaVentaValidator) Validate() error {
	// Validar RUTs
	if err := v.validateRUTs(); err != nil {
		return err
	}

	// Validar folio
	if err := v.validateFolio(); err != nil {
		return err
	}

	// Validar fecha
	if err := v.validateFecha(); err != nil {
		return err
	}

	// Validar montos
	if err := v.validateMontos(); err != nil {
		return err
	}

	// Validar tipo de nota de venta
	if err := v.validateTipoNotaVenta(); err != nil {
		return err
	}

	// Validar items
	if err := v.validateItems(); err != nil {
		return err
	}

	return nil
}

// validateRUTs valida los RUTs de emisor y receptor
func (v *NotaVentaValidator) validateRUTs() error {
	if err := ValidateRUT(v.RUTEmisor); err != nil {
		return fmt.Errorf("RUT emisor inválido: %v", err)
	}
	if err := ValidateRUT(v.RUTReceptor); err != nil {
		return fmt.Errorf("RUT receptor inválido: %v", err)
	}
	return nil
}

// validateFolio valida el folio de la nota de venta
func (v *NotaVentaValidator) validateFolio() error {
	if v.Folio == "" {
		return fmt.Errorf("folio es requerido")
	}

	return nil
}

// validateFecha valida la fecha de emisión
func (v *NotaVentaValidator) validateFecha() error {
	if v.FechaEmision.IsZero() {
		return fmt.Errorf("fecha de emisión es obligatoria")
	}

	if v.FechaEmision.After(time.Now()) {
		return fmt.Errorf("fecha de emisión no puede ser futura")
	}

	return nil
}

// validateMontos valida los montos de la nota de venta
func (v *NotaVentaValidator) validateMontos() error {
	if err := v.amountValidator.ValidateAmount(v.MontoNeto, "monto neto"); err != nil {
		return err
	}

	if err := v.amountValidator.ValidateAmount(v.MontoIVA, "monto IVA"); err != nil {
		return err
	}

	if err := v.amountValidator.ValidateAmount(v.MontoTotal, "monto total"); err != nil {
		return err
	}

	// Validar que el total sea la suma de neto + IVA
	totalCalculado := v.amountValidator.RoundAmount(v.MontoNeto + v.MontoIVA)
	if totalCalculado != v.amountValidator.RoundAmount(v.MontoTotal) {
		return fmt.Errorf("el monto total no coincide con la suma de neto e IVA")
	}

	return nil
}

// validateTipoNotaVenta valida el tipo de nota de venta
func (v *NotaVentaValidator) validateTipoNotaVenta() error {
	tiposValidos := map[string]bool{
		"CONTADO":     true,
		"CREDITO":     true,
		"EXPORTACION": true,
	}

	if !tiposValidos[v.TipoNotaVenta] {
		return fmt.Errorf("tipo de nota de venta inválido: %s", v.TipoNotaVenta)
	}

	return nil
}

// validateItems valida los items de la nota de venta
func (v *NotaVentaValidator) validateItems() error {
	if len(v.Items) == 0 {
		return fmt.Errorf("la nota de venta debe tener al menos un ítem")
	}

	var totalItems float64
	for i, item := range v.Items {
		// Validar campos obligatorios
		if item.Codigo == "" {
			return fmt.Errorf("ítem %d: código es obligatorio", i+1)
		}
		if item.Descripcion == "" {
			return fmt.Errorf("ítem %d: descripción es obligatoria", i+1)
		}

		if err := v.amountValidator.ValidateQuantity(item.Cantidad); err != nil {
			return fmt.Errorf("ítem %d: %v", i+1, err)
		}

		if err := v.amountValidator.ValidateUnitPrice(item.PrecioUnitario); err != nil {
			return fmt.Errorf("ítem %d: %v", i+1, err)
		}

		if err := v.amountValidator.ValidateDiscount(item.PrecioUnitario*item.Cantidad, item.Descuento, fmt.Sprintf("ítem %d: descuento", i+1)); err != nil {
			return err
		}

		// Calcular subtotal
		subtotal := v.amountValidator.CalculateSubtotal(item.Cantidad, item.PrecioUnitario, item.Descuento)
		if v.amountValidator.RoundAmount(subtotal) != v.amountValidator.RoundAmount(item.Subtotal) {
			return fmt.Errorf("ítem %d: subtotal calculado no coincide", i+1)
		}

		totalItems += subtotal
	}

	// Validar que el total de items coincida con el monto neto
	if v.amountValidator.RoundAmount(totalItems) != v.amountValidator.RoundAmount(v.MontoNeto) {
		return fmt.Errorf("el total de items (%.2f) no coincide con el monto neto (%.2f)", totalItems, v.MontoNeto)
	}

	return nil
}

// CalculateTotals calcula los totales de la nota de venta
func (v *NotaVentaValidator) CalculateTotals() error {
	var totalNeto float64

	for _, item := range v.Items {
		subtotal := v.amountValidator.CalculateSubtotal(item.Cantidad, item.PrecioUnitario, item.Descuento)
		totalNeto += subtotal
	}

	v.MontoNeto = v.amountValidator.RoundAmount(totalNeto)
	v.MontoIVA = v.amountValidator.RoundAmount(totalNeto * 0.19)
	v.MontoTotal = v.amountValidator.RoundAmount(v.MontoNeto + v.MontoIVA)

	return nil
}
