package utils

import (
	"fmt"
	"math"
	"strings"
)

// AmountValidator define la validación de montos
type AmountValidator struct {
	maxAmount float64
	minAmount float64
}

// NewAmountValidator crea una nueva instancia de AmountValidator
func NewAmountValidator() *AmountValidator {
	return &AmountValidator{
		maxAmount: 9999999999.99, // Un billón menos un centavo
		minAmount: 0.01,          // Un centavo
	}
}

// ValidateAmount valida un monto
func (v *AmountValidator) ValidateAmount(amount float64, field string) error {
	if amount < v.minAmount {
		return fmt.Errorf("el %s debe ser mayor o igual a %.2f", field, v.minAmount)
	}
	if amount > v.maxAmount {
		return fmt.Errorf("el %s debe ser menor o igual a %.2f", field, v.maxAmount)
	}
	return nil
}

// ValidateTotalAmount valida un monto total
func (v *AmountValidator) ValidateTotalAmount(amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("el monto total debe ser mayor que cero")
	}
	if amount > v.maxAmount {
		return fmt.Errorf("el monto total debe ser menor o igual a %.2f", v.maxAmount)
	}
	return nil
}

// ValidateQuantity valida una cantidad
func (v *AmountValidator) ValidateQuantity(quantity float64) error {
	if quantity <= 0 {
		return fmt.Errorf("la cantidad debe ser mayor que cero")
	}
	if quantity > 999999 {
		return fmt.Errorf("la cantidad debe ser menor o igual a 999,999")
	}
	return nil
}

// ValidateUnitPrice valida un precio unitario
func (v *AmountValidator) ValidateUnitPrice(price float64) error {
	if price <= 0 {
		return fmt.Errorf("el precio unitario debe ser mayor que cero")
	}
	if price > 9999999999.99 {
		return fmt.Errorf("el precio unitario debe ser menor o igual a 9,999,999,999.99")
	}
	return nil
}

// ValidatePercentage valida un porcentaje
func (v *AmountValidator) ValidatePercentage(percentage float64, field string) error {
	if percentage < 0 {
		return fmt.Errorf("el %s no puede ser negativo", field)
	}
	if percentage > 100 {
		return fmt.Errorf("el %s no puede ser mayor a 100%%", field)
	}
	return nil
}

// CalculateSubtotal calcula el subtotal de un ítem
func (v *AmountValidator) CalculateSubtotal(quantity, unitPrice, discountPercentage float64) float64 {
	subtotal := quantity * unitPrice
	if discountPercentage > 0 {
		subtotal = subtotal * (1 - discountPercentage/100)
	}
	return v.RoundAmount(subtotal)
}

// CalculateIVA calcula el IVA
func (v *AmountValidator) CalculateIVA(amount, ivaPercentage float64) float64 {
	iva := amount * (ivaPercentage / 100)
	return v.RoundAmount(iva)
}

// RoundAmount redondea un monto a 2 decimales
func (v *AmountValidator) RoundAmount(amount float64) float64 {
	return math.Round(amount*100) / 100
}

// CalculateTaxes calcula los impuestos
func (v *AmountValidator) CalculateTaxes(amount float64, taxesPercentage []float64) []float64 {
	var taxes []float64
	for _, percentage := range taxesPercentage {
		tax := amount * (percentage / 100)
		taxes = append(taxes, v.RoundAmount(tax))
	}
	return taxes
}

// ValidateQuantityWithUnit valida una cantidad con su unidad específica
func (v *AmountValidator) ValidateQuantityWithUnit(quantity float64, unit string) error {
	if err := v.ValidateQuantity(quantity); err != nil {
		return err
	}

	// Validaciones específicas por unidad
	switch unit {
	case "KG", "kg":
		if quantity > 999999.999 {
			return fmt.Errorf("para kilogramos, la cantidad debe ser menor o igual a 999,999.999")
		}
	case "LT", "lt":
		if quantity > 999999.999 {
			return fmt.Errorf("para litros, la cantidad debe ser menor o igual a 999,999.999")
		}
	case "UN", "un":
		if quantity != float64(int(quantity)) {
			return fmt.Errorf("para unidades, la cantidad debe ser un número entero")
		}
	}

	return nil
}

// FormatAmount formatea un monto como string con 2 decimales
func (v *AmountValidator) FormatAmount(amount float64) string {
	// Validar que el monto sea válido
	if err := v.ValidateAmount(amount, "monto"); err != nil {
		return "0.00"
	}

	// Formatear el monto con 2 decimales
	return fmt.Sprintf("%.2f", v.RoundAmount(amount))
}

// isValidDecimal verifica si un número tiene el número correcto de decimales
func (v *AmountValidator) isValidDecimal(value float64) bool {
	// Validar que no sea infinito o NaN
	if math.IsInf(value, 0) || math.IsNaN(value) {
		return false
	}

	// Convertir a string para contar decimales
	str := fmt.Sprintf("%f", value)
	parts := strings.Split(str, ".")
	if len(parts) != 2 {
		return true // No tiene decimales
	}

	// Contar decimales significativos
	decimalPart := strings.TrimRight(parts[1], "0")
	return len(decimalPart) <= 2
}

// ValidateDecimal verifica si un número tiene la cantidad correcta de decimales
func (v *AmountValidator) ValidateDecimal(number float64, decimals int, fieldName string) error {
	// Validar que el número no sea infinito o NaN
	if math.IsInf(number, 0) || math.IsNaN(number) {
		return fmt.Errorf("%s no puede ser infinito o NaN", fieldName)
	}

	// Validar que el número de decimales sea válido
	if decimals < 0 || decimals > 2 {
		return fmt.Errorf("el número de decimales debe estar entre 0 y 2 para %s", fieldName)
	}

	// Validar que el número tenga el número correcto de decimales
	if !v.isValidDecimal(number) {
		return fmt.Errorf("%s no puede tener más de 2 decimales", fieldName)
	}

	return nil
}

// ValidateDiscount valida un descuento
func (v *AmountValidator) ValidateDiscount(amount float64, discount float64, field string) error {
	if discount < 0 {
		return fmt.Errorf("el %s no puede ser negativo", field)
	}

	if discount > 100 {
		return fmt.Errorf("el %s no puede ser mayor que 100%%", field)
	}

	// Validar que el monto final después del descuento sea positivo
	montoDescontado := amount * (1 - discount/100)
	if montoDescontado < 0 {
		return fmt.Errorf("el monto final después del descuento no puede ser negativo")
	}

	return nil
}

// CalculateDiscount calcula el monto de descuento
func (v *AmountValidator) CalculateDiscount(amount float64, discountPercentage float64) float64 {
	discount := amount * (discountPercentage / 100)
	return v.RoundAmount(discount)
}

// ValidateAmountsConsistency valida la consistencia de los montos
func (v *AmountValidator) ValidateAmountsConsistency(netAmount, exemptAmount, taxAmount, totalAmount float64, additionalTaxes ...float64) error {
	// Validar montos individuales
	if err := v.ValidateAmount(netAmount, "monto neto"); err != nil {
		return err
	}

	if exemptAmount > 0 {
		if err := v.ValidateAmount(exemptAmount, "monto exento"); err != nil {
			return err
		}
	}

	if taxAmount > 0 {
		if err := v.ValidateAmount(taxAmount, "monto impuesto"); err != nil {
			return err
		}
	}

	if err := v.ValidateTotalAmount(totalAmount); err != nil {
		return err
	}

	// Validar impuestos adicionales
	var totalAdditionalTaxes float64
	for _, tax := range additionalTaxes {
		if tax < 0 {
			return fmt.Errorf("impuesto adicional no puede ser negativo: %.2f", tax)
		}
		totalAdditionalTaxes += tax
	}

	// Validar consistencia de montos
	calculatedTotal := netAmount + exemptAmount + taxAmount + totalAdditionalTaxes
	calculatedTotal = v.RoundAmount(calculatedTotal)
	totalAmount = v.RoundAmount(totalAmount)

	// Permitir una pequeña diferencia por redondeo (máximo 1 centavo)
	if math.Abs(calculatedTotal-totalAmount) > 0.01 {
		return fmt.Errorf("inconsistencia en los montos: neto(%.2f) + exento(%.2f) + impuesto(%.2f) + adicionales(%.2f) = %.2f, pero el total es %.2f",
			netAmount, exemptAmount, taxAmount, totalAdditionalTaxes, calculatedTotal, totalAmount)
	}

	return nil
}
