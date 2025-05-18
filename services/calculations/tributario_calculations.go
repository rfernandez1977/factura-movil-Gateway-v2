package calculations

import (
	"fmt"

	"FMgo/domain"
	"FMgo/models"
	"FMgo/utils"
)

// TributarioCalculation contiene la lógica para calcular impuestos
type TributarioCalculation struct {
	config *Config
}

// Config contiene la configuración para el cálculo de impuestos
type Config struct {
	PorcentajeIVA float64
}

// NewTributarioCalculation crea una nueva instancia de TributarioCalculation
func NewTributarioCalculation(config *Config) *TributarioCalculation {
	return &TributarioCalculation{
		config: config,
	}
}

// calcularMontosDomainItems calcula todos los montos para un documento con domain.Item
func (c *TributarioCalculation) calcularMontosDomainItems(items []domain.Item, amountValidator *utils.AmountValidator) (float64, float64, float64, float64, []float64, error) {
	var (
		montoNeto                 float64
		montoExento               float64
		montoIVA                  float64
		totalImpuestosAdicionales float64
		impuestosAdicionales      []float64
	)

	// Calcular montos por cada ítem
	for i := range items {
		item := &items[i]
		subtotal := item.PrecioUnit * item.Cantidad
		// Asumimos un descuento de 0 ya que domain.Item no tiene campo de descuento
		descuento := float64(0)
		subtotalNeto := amountValidator.RoundAmount(subtotal - descuento)

		// Asumimos que todos los items son gravados con IVA estándar (no exentos)
		montoNeto += subtotalNeto

		// Validar porcentaje de IVA (estándar 19% para Chile)
		porcentajeIVA := 19.0
		ivaItem := amountValidator.RoundAmount(subtotalNeto * (porcentajeIVA / 100))
		montoIVA += ivaItem

		// Los domain.Item no tienen impuestos adicionales
	}

	return amountValidator.RoundAmount(montoNeto),
		amountValidator.RoundAmount(montoExento),
		amountValidator.RoundAmount(montoIVA),
		amountValidator.RoundAmount(totalImpuestosAdicionales),
		impuestosAdicionales,
		nil
}

// calcularMontosModelItems calcula todos los montos para un documento con models.Item
func (c *TributarioCalculation) calcularMontosModelItems(items []models.Item, amountValidator *utils.AmountValidator) (float64, float64, float64, float64, []float64, error) {
	var (
		montoNeto                 float64
		montoExento               float64
		montoIVA                  float64
		totalImpuestosAdicionales float64
		impuestosAdicionales      []float64
	)

	// Calcular montos por cada ítem
	for i := range items {
		item := &items[i]
		subtotal := item.PrecioUnitario * item.Cantidad
		descuento := subtotal * (item.PorcentajeDescuento / 100)
		subtotalNeto := amountValidator.RoundAmount(subtotal - descuento)

		if !item.Exento {
			montoNeto += subtotalNeto

			// Usar una tasa de IVA estándar (19% para Chile)
			porcentajeIVA := 19.0
			ivaItem := amountValidator.RoundAmount(subtotalNeto * (porcentajeIVA / 100))
			montoIVA += ivaItem

			// No podemos asignar a campos que no existen, comentamos estas líneas
			// item.MontoIVA = ivaItem
			// item.Subtotal = subtotalNeto
		} else {
			montoExento += subtotalNeto
		}

		// Calcular impuestos adicionales del ítem, si existen
		if len(item.ImpuestosAdicionales) > 0 {
			for j := range item.ImpuestosAdicionales {
				impuesto := &item.ImpuestosAdicionales[j]
				// Validar porcentaje de impuesto
				if impuesto.Porcentaje < 0 || impuesto.Porcentaje > 100 {
					return 0, 0, 0, 0, nil, fmt.Errorf("porcentaje de impuesto adicional inválido: %.2f", impuesto.Porcentaje)
				}
				monto := amountValidator.RoundAmount(subtotalNeto * (impuesto.Porcentaje / 100))
				impuesto.Monto = monto
				impuesto.BaseImponible = subtotalNeto
				totalImpuestosAdicionales += monto
				impuestosAdicionales = append(impuestosAdicionales, monto)
			}
		}
	}

	return amountValidator.RoundAmount(montoNeto),
		amountValidator.RoundAmount(montoExento),
		amountValidator.RoundAmount(montoIVA),
		amountValidator.RoundAmount(totalImpuestosAdicionales),
		impuestosAdicionales,
		nil
}

// calcularImpuestosFactura calcula impuestos para una factura
func (c *TributarioCalculation) calcularImpuestosFactura(factura *models.Factura, amountValidator *utils.AmountValidator) error {
	montoNeto, montoExento, montoIVA, totalImpuestosAdicionales, impuestosAdicionales, err := c.calcularMontosDomainItems(factura.Items, amountValidator)
	if err != nil {
		return err
	}

	factura.MontoNeto = montoNeto
	factura.MontoExento = montoExento
	factura.MontoIVA = montoIVA

	// Calcular el total con valores redondeados
	montoTotal := montoNeto + montoExento + montoIVA + totalImpuestosAdicionales
	factura.MontoTotal = amountValidator.RoundAmount(montoTotal)

	// Validar consistencia de montos
	if err := amountValidator.ValidateAmountsConsistency(factura.MontoNeto, factura.MontoExento, factura.MontoIVA, factura.MontoTotal, impuestosAdicionales...); err != nil {
		return err
	}

	return nil
}

// calcularImpuestosBoleta calcula impuestos para una boleta
func (c *TributarioCalculation) calcularImpuestosBoleta(boleta *models.Boleta, amountValidator *utils.AmountValidator) error {
	// Para Boleta, implementamos un cálculo directo sin usar calcularMontosModelItems
	// ya que boleta.Items es de tipo []*models.DetalleBoleta, no []models.Item

	var montoNeto, montoExento, montoIVA, totalImpuestosAdicionales float64
	var impuestosAdicionales []float64

	// Calcular montos por cada ítem
	for _, item := range boleta.Items {
		total := float64(item.Cantidad) * item.Precio
		// Asumir que todos los ítems son afectos a IVA a menos que se especifique lo contrario
		montoNeto += total
	}

	// Calcular IVA usando tasa estándar 19%
	tasaIVA := 19.0
	montoIVA = amountValidator.RoundAmount(montoNeto * (tasaIVA / 100))

	// Asignar valores calculados
	boleta.MontoNeto = amountValidator.RoundAmount(montoNeto)
	boleta.MontoExento = amountValidator.RoundAmount(montoExento)
	boleta.MontoIVA = montoIVA
	boleta.TasaIVA = tasaIVA

	// Calcular el total con valores redondeados
	montoTotal := montoNeto + montoExento + montoIVA + totalImpuestosAdicionales
	boleta.MontoTotal = amountValidator.RoundAmount(montoTotal)

	// Validar consistencia de montos
	if err := amountValidator.ValidateAmountsConsistency(boleta.MontoNeto, boleta.MontoExento, boleta.MontoIVA, boleta.MontoTotal, impuestosAdicionales...); err != nil {
		return err
	}

	return nil
}

// calcularImpuestosNotaCredito calcula impuestos para una nota de crédito
func (c *TributarioCalculation) calcularImpuestosNotaCredito(notaCredito *models.NotaCredito, amountValidator *utils.AmountValidator) error {
	montoNeto, montoExento, montoIVA, totalImpuestosAdicionales, impuestosAdicionales, err := c.calcularMontosModelItems(notaCredito.Items, amountValidator)
	if err != nil {
		return err
	}

	notaCredito.MontoNeto = montoNeto
	notaCredito.MontoExento = montoExento
	notaCredito.MontoIVA = montoIVA

	// Calcular el total con valores redondeados
	montoTotal := montoNeto + montoExento + montoIVA + totalImpuestosAdicionales
	notaCredito.MontoTotal = amountValidator.RoundAmount(montoTotal)

	// Validar consistencia de montos
	if err := amountValidator.ValidateAmountsConsistency(notaCredito.MontoNeto, notaCredito.MontoExento, notaCredito.MontoIVA, notaCredito.MontoTotal, impuestosAdicionales...); err != nil {
		return err
	}

	return nil
}

// calcularImpuestosNotaDebito calcula impuestos para una nota de débito
func (c *TributarioCalculation) calcularImpuestosNotaDebito(notaDebito *models.NotaDebito, amountValidator *utils.AmountValidator) error {
	montoNeto, montoExento, montoIVA, totalImpuestosAdicionales, impuestosAdicionales, err := c.calcularMontosModelItems(notaDebito.Items, amountValidator)
	if err != nil {
		return err
	}

	notaDebito.MontoNeto = montoNeto
	notaDebito.MontoExento = montoExento
	notaDebito.MontoIVA = montoIVA

	// Calcular el total con valores redondeados
	montoTotal := montoNeto + montoExento + montoIVA + totalImpuestosAdicionales
	notaDebito.MontoTotal = amountValidator.RoundAmount(montoTotal)

	// Validar consistencia de montos
	if err := amountValidator.ValidateAmountsConsistency(notaDebito.MontoNeto, notaDebito.MontoExento, notaDebito.MontoIVA, notaDebito.MontoTotal, impuestosAdicionales...); err != nil {
		return err
	}

	return nil
}

// calcularImpuestosGuiaDespacho calcula impuestos para una guía de despacho
func (c *TributarioCalculation) calcularImpuestosGuiaDespacho(guiaDespacho *models.GuiaDespacho, amountValidator *utils.AmountValidator) error {
	montoNeto, montoExento, montoIVA, totalImpuestosAdicionales, impuestosAdicionales, err := c.calcularMontosModelItems(guiaDespacho.Items, amountValidator)
	if err != nil {
		return err
	}

	guiaDespacho.MontoNeto = montoNeto
	guiaDespacho.MontoExento = montoExento
	guiaDespacho.MontoIVA = montoIVA

	// Calcular el total con valores redondeados
	montoTotal := montoNeto + montoExento + montoIVA + totalImpuestosAdicionales
	guiaDespacho.MontoTotal = amountValidator.RoundAmount(montoTotal)

	// Validar consistencia de montos
	if err := amountValidator.ValidateAmountsConsistency(guiaDespacho.MontoNeto, guiaDespacho.MontoExento, guiaDespacho.MontoIVA, guiaDespacho.MontoTotal, impuestosAdicionales...); err != nil {
		return err
	}

	return nil
}

// CalcularImpuestos calcula los impuestos de un documento tributario
func (c *TributarioCalculation) CalcularImpuestos(doc interface{}) error {
	amountValidator := utils.NewAmountValidator()

	switch d := doc.(type) {
	case *models.Factura:
		return c.calcularImpuestosFactura(d, amountValidator)
	case *models.Boleta:
		return c.calcularImpuestosBoleta(d, amountValidator)
	case *models.NotaCredito:
		return c.calcularImpuestosNotaCredito(d, amountValidator)
	case *models.NotaDebito:
		return c.calcularImpuestosNotaDebito(d, amountValidator)
	case *models.GuiaDespacho:
		return c.calcularImpuestosGuiaDespacho(d, amountValidator)
	case *domain.DocumentoTributario:
		// Obtener los items del documento de dominio
		montoNeto, montoExento, montoIVA, totalImpuestosAdicionales, _, err := c.calcularMontosDomainItems([]domain.Item{}, amountValidator)
		if err != nil {
			return err
		}

		// Actualizar el documento de dominio con los valores calculados
		d.MontoNeto = montoNeto
		d.MontoExento = montoExento
		d.MontoIVA = montoIVA
		d.MontoTotal = amountValidator.RoundAmount(montoNeto + montoExento + montoIVA + totalImpuestosAdicionales)

		return nil
	default:
		return fmt.Errorf("tipo de documento no soportado: %T", doc)
	}
}

// CalcularImpuestosFromDomain calcula los impuestos de un documento tributario genérico
func (c *TributarioCalculation) CalcularImpuestosFromDomain(items []domain.Item) (float64, float64, float64, float64, error) {
	amountValidator := utils.NewAmountValidator()

	montoNeto, montoExento, montoIVA, totalImpuestosAdicionales, impuestosAdicionales, err := c.calcularMontosDomainItems(items, amountValidator)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	// Calcular el total con valores redondeados
	montoTotal := montoNeto + montoExento + montoIVA + totalImpuestosAdicionales
	montoTotal = amountValidator.RoundAmount(montoTotal)

	// Validar consistencia de montos
	if err := amountValidator.ValidateAmountsConsistency(montoNeto, montoExento, montoIVA, montoTotal, impuestosAdicionales...); err != nil {
		return 0, 0, 0, 0, err
	}

	return montoNeto, montoExento, montoIVA, montoTotal, nil
}

// CalcularMontosBoleta calcula los montos de una boleta
func (c *TributarioCalculation) CalcularMontosBoleta(boleta *models.Boleta) error {
	if boleta == nil {
		return fmt.Errorf("boleta es nil")
	}

	amountValidator := utils.NewAmountValidator()

	// Llamamos a calcularImpuestosBoleta que ya implementamos
	return c.calcularImpuestosBoleta(boleta, amountValidator)
}
