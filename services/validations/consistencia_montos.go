package validations

import (
	"fmt"
	"math"

	"FMgo/models"
)

// ValidadorConsistenciaMontos maneja la validación centralizada de consistencia de montos
type ValidadorConsistenciaMontos struct {
	toleranciaIVA    float64
	toleranciaTotal  float64
	limiteMontoNeto  float64
	limiteMontoTotal float64
	limiteMontoIVA   float64
	tasaIVA          float64
}

// NewValidadorConsistenciaMontos crea una nueva instancia del validador
func NewValidadorConsistenciaMontos() *ValidadorConsistenciaMontos {
	return &ValidadorConsistenciaMontos{
		toleranciaIVA:    0.02,       // 2 centavos
		toleranciaTotal:  0.02,       // 2 centavos
		limiteMontoNeto:  1000000000, // 1 billón
		limiteMontoTotal: 1190000000, // 1.19 billón
		limiteMontoIVA:   190000000,  // 190 millones
		tasaIVA:          19.0,       // 19%
	}
}

// ValidarConsistencia valida la consistencia de todos los montos
func (v *ValidadorConsistenciaMontos) ValidarConsistencia(doc interface{}) error {
	var (
		montoNeto            float64
		montoIVA             float64
		montoTotal           float64
		montoExento          float64
		descuentos           float64
		recargos             float64
		impuestosAdicionales []models.ImpuestoAdicional
	)

	// Extraer montos según el tipo de documento
	switch d := doc.(type) {
	case *models.Factura:
		montoNeto = d.MontoNeto
		montoIVA = d.MontoIVA
		montoTotal = d.MontoTotal
		montoExento = d.MontoExento
		impuestosAdicionales = d.ImpuestosAdicionales
	case *models.NotaCredito:
		montoNeto = d.MontoNeto
		montoIVA = d.MontoIVA
		montoTotal = d.MontoTotal
		montoExento = d.MontoExento
		impuestosAdicionales = d.ImpuestosAdicionales
	case *models.NotaDebito:
		montoNeto = d.MontoNeto
		montoIVA = d.MontoIVA
		montoTotal = d.MontoTotal
		montoExento = d.MontoExento
		impuestosAdicionales = d.ImpuestosAdicionales
	case *models.GuiaDespacho:
		montoNeto = d.MontoNeto
		montoIVA = d.MontoIVA
		montoTotal = d.MontoTotal
		montoExento = d.MontoExento
		impuestosAdicionales = d.ImpuestosAdicionales
	default:
		return fmt.Errorf("tipo de documento no soportado")
	}

	// Validar límites de montos
	if err := v.validarLimites(montoNeto, montoIVA, montoTotal); err != nil {
		return err
	}

	// Validar IVA
	if err := v.validarIVA(montoNeto, montoIVA); err != nil {
		return err
	}

	// Calcular total de impuestos adicionales
	totalImpuestosAdicionales := 0.0
	if len(impuestosAdicionales) > 0 {
		for _, imp := range impuestosAdicionales {
			totalImpuestosAdicionales += imp.Monto
		}
	}

	// Validar total
	totalCalculado := math.Round((montoNeto+montoExento+montoIVA+totalImpuestosAdicionales-descuentos+recargos)*100) / 100
	if math.Abs(totalCalculado-montoTotal) > v.toleranciaTotal {
		return fmt.Errorf("el total calculado (%.2f) no coincide con el monto total proporcionado (%.2f), diferencia: %.2f, tolerancia máxima: %.2f",
			totalCalculado, montoTotal, math.Abs(totalCalculado-montoTotal), v.toleranciaTotal)
	}

	return nil
}

// validarLimites valida los límites de los montos
func (v *ValidadorConsistenciaMontos) validarLimites(montoNeto, montoIVA, montoTotal float64) error {
	if montoNeto > v.limiteMontoNeto {
		return fmt.Errorf("monto neto excede el límite máximo permitido (%.2f)", v.limiteMontoNeto)
	}
	if montoIVA > v.limiteMontoIVA {
		return fmt.Errorf("monto IVA excede el límite máximo permitido (%.2f)", v.limiteMontoIVA)
	}
	if montoTotal > v.limiteMontoTotal {
		return fmt.Errorf("monto total excede el límite máximo permitido (%.2f)", v.limiteMontoTotal)
	}
	return nil
}

// validarIVA valida el cálculo y consistencia del IVA
func (v *ValidadorConsistenciaMontos) validarIVA(montoNeto, montoIVA float64) error {
	if montoNeto > 0 {
		ivaCalculado := math.Round((montoNeto*v.tasaIVA/100)*100) / 100
		if math.Abs(ivaCalculado-montoIVA) > v.toleranciaIVA {
			return fmt.Errorf("el IVA calculado (%.2f) no coincide con el monto IVA proporcionado (%.2f), diferencia: %.2f, tolerancia máxima: %.2f",
				ivaCalculado, montoIVA, math.Abs(ivaCalculado-montoIVA), v.toleranciaIVA)
		}
	} else if montoIVA > 0 {
		return fmt.Errorf("se ha proporcionado un monto de IVA (%.2f) pero el monto neto es cero o negativo", montoIVA)
	}
	return nil
}

// ValidarConsistenciaItems valida la consistencia entre los items y los totales
func (v *ValidadorConsistenciaMontos) ValidarConsistenciaItems(items interface{}, montoNeto, montoExento float64) error {
	var totalNeto, totalExento float64

	switch items := items.(type) {
	case []models.Item:
		for _, item := range items {
			subtotal := item.PrecioUnitario * item.Cantidad
			if item.Exento {
				totalExento += subtotal
			} else {
				totalNeto += subtotal
			}
		}
	case []*models.DetalleBoleta:
		for _, item := range items {
			subtotal := item.Precio * float64(item.Cantidad)
			if item.Exento {
				totalExento += subtotal
			} else {
				totalNeto += subtotal
			}
		}
	default:
		return fmt.Errorf("tipo de items no soportado")
	}

	// Redondear totales calculados
	totalNeto = math.Round(totalNeto*100) / 100
	totalExento = math.Round(totalExento*100) / 100

	// Validar consistencia con tolerancia
	if math.Abs(totalNeto-montoNeto) > v.toleranciaTotal {
		return fmt.Errorf("el total neto calculado de items (%.2f) no coincide con el monto neto del documento (%.2f)",
			totalNeto, montoNeto)
	}
	if math.Abs(totalExento-montoExento) > v.toleranciaTotal {
		return fmt.Errorf("el total exento calculado de items (%.2f) no coincide con el monto exento del documento (%.2f)",
			totalExento, montoExento)
	}

	return nil
}
