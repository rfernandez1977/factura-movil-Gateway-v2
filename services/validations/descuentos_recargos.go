package validations

import (
	"fmt"
	"math"
)

// TipoMovimiento representa el tipo de movimiento (Descuento o Recargo)
type TipoMovimiento string

const (
	TipoDescuento TipoMovimiento = "D"
	TipoRecargo   TipoMovimiento = "R"
)

// DescuentoRecargo representa un descuento o recargo global
type DescuentoRecargo struct {
	NumeroLinea int            `json:"numero_linea"`
	TipoMov     TipoMovimiento `json:"tipo_mov"`
	Glosa       string         `json:"glosa,omitempty"`
	TipoValor   string         `json:"tipo_valor"` // % o $
	ValorMonto  float64        `json:"valor_monto"`
	MontoMovim  float64        `json:"monto_movim,omitempty"`
}

// ValidadorDescuentosRecargos maneja la validación de descuentos y recargos
type ValidadorDescuentosRecargos struct {
	montoNeto   float64
	montoExento float64
}

// NewValidadorDescuentosRecargos crea una nueva instancia del validador
func NewValidadorDescuentosRecargos(montoNeto, montoExento float64) *ValidadorDescuentosRecargos {
	return &ValidadorDescuentosRecargos{
		montoNeto:   montoNeto,
		montoExento: montoExento,
	}
}

// ValidarSecuencia valida la secuencia de descuentos y recargos según el esquema XSD
func (v *ValidadorDescuentosRecargos) ValidarSecuencia(movimientos []DescuentoRecargo) error {
	if len(movimientos) > 20 {
		return fmt.Errorf("número máximo de descuentos/recargos excedido (máximo 20, actual: %d)", len(movimientos))
	}

	numerosLinea := make(map[int]bool)
	baseImponible := v.montoNeto + v.montoExento

	for i, mov := range movimientos {
		// Validar número de línea
		if mov.NumeroLinea <= 0 {
			return fmt.Errorf("línea %d: número de línea debe ser mayor a 0", i+1)
		}
		if numerosLinea[mov.NumeroLinea] {
			return fmt.Errorf("línea %d: número de línea duplicado", mov.NumeroLinea)
		}
		numerosLinea[mov.NumeroLinea] = true

		// Validar tipo de movimiento
		if mov.TipoMov != TipoDescuento && mov.TipoMov != TipoRecargo {
			return fmt.Errorf("línea %d: tipo de movimiento inválido", mov.NumeroLinea)
		}

		// Validar glosa
		if len(mov.Glosa) > 45 {
			return fmt.Errorf("línea %d: glosa excede longitud máxima (45 caracteres)", mov.NumeroLinea)
		}

		// Validar tipo de valor y monto
		switch mov.TipoValor {
		case "%":
			if mov.ValorMonto < 0 || mov.ValorMonto > 100 {
				return fmt.Errorf("línea %d: porcentaje debe estar entre 0 y 100", mov.NumeroLinea)
			}
			montoCalculado := math.Round((baseImponible*(mov.ValorMonto/100))*100) / 100
			if math.Abs(montoCalculado-mov.MontoMovim) > 0.01 {
				return fmt.Errorf("línea %d: monto calculado (%.2f) no coincide con el proporcionado (%.2f)",
					mov.NumeroLinea, montoCalculado, mov.MontoMovim)
			}
		case "$":
			if mov.ValorMonto < 0 {
				return fmt.Errorf("línea %d: monto debe ser mayor o igual a 0", mov.NumeroLinea)
			}
			if math.Abs(mov.ValorMonto-mov.MontoMovim) > 0.01 {
				return fmt.Errorf("línea %d: monto no coincide con el valor especificado", mov.NumeroLinea)
			}
		default:
			return fmt.Errorf("línea %d: tipo de valor inválido (debe ser '%' o '$')", mov.NumeroLinea)
		}

		// Actualizar base imponible para el siguiente movimiento
		if mov.TipoMov == TipoDescuento {
			baseImponible -= mov.MontoMovim
		} else {
			baseImponible += mov.MontoMovim
		}

		// Validar que la base imponible no quede negativa
		if baseImponible < 0 {
			return fmt.Errorf("línea %d: los descuentos no pueden resultar en un monto negativo", mov.NumeroLinea)
		}
	}

	return nil
}

// CalcularTotales calcula los totales después de aplicar descuentos y recargos
func (v *ValidadorDescuentosRecargos) CalcularTotales(movimientos []DescuentoRecargo) (float64, float64, error) {
	var totalDescuentos, totalRecargos float64
	baseImponible := v.montoNeto + v.montoExento

	for _, mov := range movimientos {
		if mov.TipoMov == TipoDescuento {
			totalDescuentos += mov.MontoMovim
		} else {
			totalRecargos += mov.MontoMovim
		}
	}

	montoFinal := baseImponible - totalDescuentos + totalRecargos
	if montoFinal < 0 {
		return 0, 0, fmt.Errorf("el total después de descuentos y recargos no puede ser negativo")
	}

	return totalDescuentos, totalRecargos, nil
}
