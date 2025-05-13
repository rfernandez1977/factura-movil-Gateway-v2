package utils

import (
	"fmt"
	"time"
)

// PaymentValidator contiene utilidades para validar pagos
type PaymentValidator struct {
	amountValidator  *AmountValidator
	dateValidator    *DateValidator
	TipoNotaVenta    string
	MontoTotal       float64
	FechaEmision     time.Time
	FechaVencimiento time.Time
	Cuotas           []CuotaPago
	Interes          float64
	Descuento        float64
	Moneda           string
	TipoCambio       float64
}

// CuotaPago representa una cuota de pago
type CuotaPago struct {
	Numero           int
	Monto            float64
	FechaVencimiento time.Time
	Estado           string
	Interes          float64
	Descuento        float64
	Saldo            float64
}

// NewPaymentValidator crea una nueva instancia de PaymentValidator
func NewPaymentValidator() *PaymentValidator {
	return &PaymentValidator{
		amountValidator: NewAmountValidator(),
		dateValidator:   NewDateValidator(),
	}
}

// ValidatePayment valida un pago individual
func (v *PaymentValidator) ValidatePayment(monto float64, fechaVencimiento time.Time, tipo string) error {
	if monto <= 0 {
		return fmt.Errorf("el monto del %s debe ser mayor a cero", tipo)
	}

	if fechaVencimiento.IsZero() {
		return fmt.Errorf("la fecha de vencimiento del %s es requerida", tipo)
	}

	if fechaVencimiento.Before(time.Now()) {
		return fmt.Errorf("la fecha de vencimiento del %s no puede ser anterior a la fecha actual", tipo)
	}

	return nil
}

// ValidatePaymentSchedule valida un cronograma de pagos
func (v *PaymentValidator) ValidatePaymentSchedule(pagos []struct {
	Monto      float64
	FechaPago  time.Time
	MetodoPago string
}, montoTotal float64) error {
	if len(pagos) == 0 {
		return fmt.Errorf("el cronograma de pagos no puede estar vacío")
	}

	var totalPagos float64
	for _, pago := range pagos {
		if err := v.ValidatePayment(pago.Monto, pago.FechaPago, "pago"); err != nil {
			return err
		}
		totalPagos += pago.Monto
	}

	// Validar que la suma de los pagos coincida con el monto total
	if v.amountValidator.RoundAmount(totalPagos) != v.amountValidator.RoundAmount(montoTotal) {
		return fmt.Errorf("la suma de los pagos (%.2f) no coincide con el monto total (%.2f)", totalPagos, montoTotal)
	}

	return nil
}

// ValidatePaymentMethod valida un método de pago
func (v *PaymentValidator) ValidatePaymentMethod(metodoPago string) error {
	metodosValidos := map[string]bool{
		"EFECTIVO":      true,
		"TRANSFERENCIA": true,
		"CHEQUE":        true,
		"TARJETA":       true,
	}

	if !metodosValidos[metodoPago] {
		return fmt.Errorf("método de pago no válido: %s", metodoPago)
	}

	return nil
}

// ValidatePaymentDates valida las fechas de pago
func (v *PaymentValidator) ValidatePaymentDates(fechaEmision time.Time, fechaVencimiento time.Time) error {
	if fechaEmision.IsZero() {
		return fmt.Errorf("la fecha de emisión es requerida")
	}

	if fechaVencimiento.IsZero() {
		return fmt.Errorf("la fecha de vencimiento es requerida")
	}

	if fechaVencimiento.Before(fechaEmision) {
		return fmt.Errorf("la fecha de vencimiento no puede ser anterior a la fecha de emisión")
	}

	return nil
}

// Validate valida las condiciones de pago
func (v *PaymentValidator) Validate() error {
	// Validar tipo de nota de venta
	if err := v.validateTipoNotaVenta(); err != nil {
		return err
	}

	// Validar fechas
	if err := v.dateValidator.ValidateDocumentDates(v.FechaEmision, v.FechaVencimiento); err != nil {
		return err
	}

	// Validar cuotas
	if err := v.validateCuotas(); err != nil {
		return err
	}

	// Validar montos y tasas
	if err := v.validateMontos(); err != nil {
		return err
	}

	// Validar moneda y tipo de cambio
	if err := v.validateMoneda(); err != nil {
		return err
	}

	return nil
}

// validateTipoNotaVenta valida el tipo de nota de venta
func (v *PaymentValidator) validateTipoNotaVenta() error {
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

// validateCuotas valida las cuotas de pago
func (v *PaymentValidator) validateCuotas() error {
	// Para notas de venta al contado, no debe haber cuotas
	if v.TipoNotaVenta == "CONTADO" && len(v.Cuotas) > 0 {
		return fmt.Errorf("para notas de venta al contado, no debe haber cuotas")
	}

	// Para notas de venta a crédito, debe haber al menos una cuota
	if v.TipoNotaVenta == "CREDITO" && len(v.Cuotas) == 0 {
		return fmt.Errorf("para notas de venta a crédito, debe haber al menos una cuota")
	}

	// Validar cada cuota
	var totalCuotas float64
	for _, cuota := range v.Cuotas {
		if err := v.ValidatePayment(cuota.Monto, cuota.FechaVencimiento, "cuota"); err != nil {
			return err
		}
		totalCuotas += cuota.Monto
	}

	// Validar que la suma de las cuotas coincida con el monto total
	if v.amountValidator.RoundAmount(totalCuotas) != v.amountValidator.RoundAmount(v.MontoTotal) {
		return fmt.Errorf("la suma de las cuotas (%.2f) no coincide con el monto total (%.2f)", totalCuotas, v.MontoTotal)
	}

	return nil
}

// validateMontos valida los montos y tasas
func (v *PaymentValidator) validateMontos() error {
	// Validar monto total
	if err := v.amountValidator.ValidateTotalAmount(v.MontoTotal); err != nil {
		return err
	}

	// Validar interés
	if err := v.amountValidator.ValidatePercentage(v.Interes, "interés"); err != nil {
		return err
	}

	// Validar descuento
	if err := v.amountValidator.ValidateDiscount(v.MontoTotal, v.Descuento, "descuento"); err != nil {
		return err
	}

	return nil
}

// validateMoneda valida la moneda y tipo de cambio
func (v *PaymentValidator) validateMoneda() error {
	monedasValidas := map[string]bool{
		"CLP": true,
		"USD": true,
		"EUR": true,
	}

	if !monedasValidas[v.Moneda] {
		return fmt.Errorf("moneda no válida: %s", v.Moneda)
	}

	if v.Moneda != "CLP" && v.TipoCambio <= 0 {
		return fmt.Errorf("el tipo de cambio debe ser mayor que 0 para monedas extranjeras")
	}

	return nil
}

// CalculateCuotas calcula las cuotas de pago
func (v *PaymentValidator) CalculateCuotas(numeroCuotas int, fechaPrimeraCuota time.Time) error {
	if numeroCuotas <= 0 {
		return fmt.Errorf("el número de cuotas debe ser mayor que 0")
	}

	if err := v.dateValidator.ValidateDate(fechaPrimeraCuota, "fecha de primera cuota"); err != nil {
		return err
	}

	// Calcular monto de cada cuota
	montoCuota := v.MontoTotal / float64(numeroCuotas)
	montoCuota = v.amountValidator.RoundAmount(montoCuota)

	// Crear cuotas
	v.Cuotas = make([]CuotaPago, numeroCuotas)
	for i := 0; i < numeroCuotas; i++ {
		fechaVencimiento := fechaPrimeraCuota.AddDate(0, i, 0)
		v.Cuotas[i] = CuotaPago{
			Numero:           i + 1,
			Monto:            montoCuota,
			FechaVencimiento: fechaVencimiento,
			Estado:           "PENDIENTE",
			Interes:          v.Interes,
			Descuento:        v.Descuento,
			Saldo:            montoCuota,
		}
	}

	return nil
}

// GetPaymentSummary obtiene un resumen del pago
func (v *PaymentValidator) GetPaymentSummary() string {
	return fmt.Sprintf("Tipo: %s, Monto Total: %s %s, Cuotas: %d",
		v.TipoNotaVenta, v.amountValidator.FormatAmount(v.MontoTotal), v.Moneda, len(v.Cuotas))
}

// GetCuotaSummary obtiene un resumen de una cuota
func (v *PaymentValidator) GetCuotaSummary(numeroCuota int) (string, error) {
	if numeroCuota < 1 || numeroCuota > len(v.Cuotas) {
		return "", fmt.Errorf("número de cuota inválido")
	}

	cuota := v.Cuotas[numeroCuota-1]
	return fmt.Sprintf("Cuota %d: %s %s, Vencimiento: %s, Estado: %s",
		cuota.Numero, v.amountValidator.FormatAmount(cuota.Monto), v.Moneda,
		v.dateValidator.FormatDate(cuota.FechaVencimiento), cuota.Estado), nil
}

// UpdateCuotaEstado actualiza el estado de una cuota
func (v *PaymentValidator) UpdateCuotaEstado(numeroCuota int, nuevoEstado string) error {
	if numeroCuota < 1 || numeroCuota > len(v.Cuotas) {
		return fmt.Errorf("número de cuota inválido")
	}

	estadosValidos := map[string]bool{
		"PENDIENTE": true,
		"PAGADA":    true,
		"VENCIDA":   true,
	}

	if !estadosValidos[nuevoEstado] {
		return fmt.Errorf("estado de cuota inválido: %s", nuevoEstado)
	}

	v.Cuotas[numeroCuota-1].Estado = nuevoEstado
	return nil
}

// CalculateSaldoPendiente calcula el saldo pendiente
func (v *PaymentValidator) CalculateSaldoPendiente() float64 {
	var saldo float64
	for _, cuota := range v.Cuotas {
		if cuota.Estado == "PENDIENTE" {
			saldo += cuota.Saldo
		}
	}
	return v.amountValidator.RoundAmount(saldo)
}
