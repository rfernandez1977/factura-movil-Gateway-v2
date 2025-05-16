package utils

import (
	"fmt"

	"github.com/fmgo/models"
)

// ItemValidator valida los items de un documento
type ItemValidator struct {
	amountValidator *AmountValidator
}

// NewItemValidator crea una nueva instancia de ItemValidator
func NewItemValidator() *ItemValidator {
	return &ItemValidator{
		amountValidator: NewAmountValidator(),
	}
}

// ValidateItem valida un ítem individual
func (v *ItemValidator) ValidateItem(item *models.ItemNotaVenta) error {
	if item == nil {
		return fmt.Errorf("el ítem no puede ser nulo")
	}

	if item.Codigo == "" {
		return fmt.Errorf("el código del ítem es requerido")
	}

	if item.Descripcion == "" {
		return fmt.Errorf("la descripción del ítem es requerida")
	}

	if err := v.amountValidator.ValidateQuantity(item.Cantidad); err != nil {
		return fmt.Errorf("cantidad inválida: %v", err)
	}

	if err := v.amountValidator.ValidateUnitPrice(item.PrecioUnitario); err != nil {
		return fmt.Errorf("precio unitario inválido: %v", err)
	}

	if item.Descuento < 0 {
		return fmt.Errorf("el descuento no puede ser negativo")
	}

	if item.Descuento > item.PrecioUnitario*item.Cantidad {
		return fmt.Errorf("el descuento no puede ser mayor que el subtotal")
	}

	// Calcular el subtotal
	subtotal := item.PrecioUnitario * item.Cantidad
	subtotal = v.amountValidator.RoundAmount(subtotal - item.Descuento)

	if subtotal != item.Subtotal {
		return fmt.Errorf("el subtotal calculado no coincide con el proporcionado")
	}

	return nil
}

// ValidateItems valida una lista de items
func (v *ItemValidator) ValidateItems(items []*models.ItemNotaVenta) error {
	if len(items) == 0 {
		return fmt.Errorf("el documento debe tener al menos un item")
	}

	for i, item := range items {
		if err := v.ValidateItem(item); err != nil {
			return fmt.Errorf("error en item %d: %v", i+1, err)
		}
	}

	return nil
}
