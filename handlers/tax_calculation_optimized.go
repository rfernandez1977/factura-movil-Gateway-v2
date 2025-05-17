package handlers

import (
	"math"
	"sync"

	"github.com/gin-gonic/gin"
)

// Cache para resultados de cálculos frecuentes
var (
	calculationCache sync.Map
	roundingPool     sync.Pool
)

func init() {
	roundingPool = sync.Pool{
		New: func() interface{} {
			return make([]float64, 0, 8) // Capacidad inicial para los valores más comunes
		},
	}
}

// Estructura optimizada para cálculos tributarios
type OptimizedTaxCalculation struct {
	// Datos base (usando tipos más eficientes)
	MontoNeto   float64     `json:"montoNeto"`
	MontoExento float64     `json:"montoExento"`
	TasaIVA     float64     `json:"tasaIVA"`
	Descuentos  []Descuento `json:"descuentos,omitempty"` // omitempty para reducir JSON
	Recargos    []Recargo   `json:"recargos,omitempty"`   // omitempty para reducir JSON

	// Retenciones específicas
	RetencionHonorarios float64 `json:"retencionHonorarios,omitempty"`
	RetencionILA        float64 `json:"retencionILA,omitempty"`

	// Resultados calculados
	MontoIVA         float64 `json:"montoIVA"`
	MontoTotal       float64 `json:"montoTotal"`
	TotalRetenciones float64 `json:"totalRetenciones"`

	// Cache de cálculos intermedios
	baseImponible float64 `json:"-"`
	totalDesc     float64 `json:"-"`
	totalRec      float64 `json:"-"`
}

// Validador optimizado que usa un pool de errores
type ValidationError struct {
	Code    string
	Message string
	Detail  string
}

var validationErrorPool = sync.Pool{
	New: func() interface{} {
		return &ValidationError{}
	},
}

// Constantes para validación
const (
	TasaIVANormal    = 19.0
	TasaHonorarios   = 10.75
	TasaILAMin       = 10.0
	TasaILAMed       = 11.0
	TasaILAMax       = 13.0
	PorcentajeMaximo = 100.0
	PorcentajeMinimo = 0.0
)

// Función optimizada para redondeo según normas SII
func roundSII(value float64) float64 {
	return math.Round(value*100) / 100
}

// ValidateBase realiza las validaciones básicas comunes
func (calc *OptimizedTaxCalculation) ValidateBase() *ValidationError {
	if calc.MontoNeto < 0 || calc.MontoExento < 0 {
		err := validationErrorPool.Get().(*ValidationError)
		err.Code = "CALC_001"
		err.Message = "Los montos no pueden ser negativos"
		err.Detail = "Monto neto y monto exento deben ser positivos o cero"
		return err
	}

	if calc.MontoNeto == 0 && calc.MontoExento == 0 {
		err := validationErrorPool.Get().(*ValidationError)
		err.Code = "CALC_001"
		err.Message = "Montos inválidos"
		err.Detail = "El documento debe tener al menos monto neto o monto exento mayor a cero"
		return err
	}

	if calc.TasaIVA != TasaIVANormal && calc.TasaIVA != 0 {
		err := validationErrorPool.Get().(*ValidationError)
		err.Code = "CALC_002"
		err.Message = "Tasa de IVA inválida"
		err.Detail = "La tasa de IVA debe ser 19% o 0%"
		return err
	}

	return nil
}

// CalculateIVA calcula el IVA de forma optimizada
func (calc *OptimizedTaxCalculation) CalculateIVA() *ValidationError {
	if calc.MontoNeto > 0 && calc.TasaIVA > 0 {
		calc.MontoIVA = roundSII(calc.MontoNeto * (calc.TasaIVA / 100))
	} else {
		calc.MontoIVA = 0
		if calc.TasaIVA > 0 && calc.MontoNeto <= 0 {
			err := validationErrorPool.Get().(*ValidationError)
			err.Code = "CALC_009"
			err.Message = "Inconsistencia en cálculo de IVA"
			err.Detail = "No se puede calcular IVA con monto neto cero o negativo"
			return err
		}
	}
	return nil
}

// CalculateDiscounts calcula los descuentos de forma optimizada
func (calc *OptimizedTaxCalculation) CalculateDiscounts() *ValidationError {
	calc.baseImponible = calc.MontoNeto + calc.MontoExento
	calc.totalDesc = 0

	for i := range calc.Descuentos {
		desc := &calc.Descuentos[i]
		if desc.Porcentaje < PorcentajeMinimo || desc.Porcentaje > PorcentajeMaximo {
			err := validationErrorPool.Get().(*ValidationError)
			err.Code = "CALC_004"
			err.Message = "Porcentaje de descuento inválido"
			err.Detail = "El porcentaje debe estar entre 0 y 100"
			return err
		}

		if desc.TipoDescuento == "GLOBAL" {
			desc.MontoDescuento = roundSII(calc.baseImponible * (desc.Porcentaje / 100))
			calc.totalDesc += desc.MontoDescuento
		}
	}
	return nil
}

// CalculateSurcharges calcula los recargos de forma optimizada
func (calc *OptimizedTaxCalculation) CalculateSurcharges() *ValidationError {
	calc.totalRec = 0

	for i := range calc.Recargos {
		rec := &calc.Recargos[i]
		if rec.Porcentaje < PorcentajeMinimo || rec.Porcentaje > PorcentajeMaximo {
			err := validationErrorPool.Get().(*ValidationError)
			err.Code = "CALC_005"
			err.Message = "Porcentaje de recargo inválido"
			err.Detail = "El porcentaje debe estar entre 0 y 100"
			return err
		}

		if rec.TipoRecargo == "GLOBAL" {
			rec.MontoRecargo = roundSII(calc.baseImponible * (rec.Porcentaje / 100))
			calc.totalRec += rec.MontoRecargo
		}
	}
	return nil
}

// CalculateRetentions calcula las retenciones de forma optimizada
func (calc *OptimizedTaxCalculation) CalculateRetentions() *ValidationError {
	calc.TotalRetenciones = 0

	if calc.RetencionHonorarios > 0 {
		if calc.RetencionHonorarios != TasaHonorarios {
			err := validationErrorPool.Get().(*ValidationError)
			err.Code = "CALC_003"
			err.Message = "Tasa de retención de honorarios inválida"
			err.Detail = "La retención de honorarios debe ser 10.75%"
			return err
		}
		calc.TotalRetenciones += roundSII(calc.MontoNeto * (calc.RetencionHonorarios / 100))
	}

	if calc.RetencionILA > 0 {
		if calc.RetencionILA != TasaILAMin && calc.RetencionILA != TasaILAMed && calc.RetencionILA != TasaILAMax {
			err := validationErrorPool.Get().(*ValidationError)
			err.Code = "CALC_006"
			err.Message = "Tasa de retención ILA inválida"
			err.Detail = "La retención ILA debe ser 10%, 11% o 13%"
			return err
		}
		calc.TotalRetenciones += roundSII(calc.MontoNeto * (calc.RetencionILA / 100))
	}

	return nil
}

// CalculateTotal calcula el total final de forma optimizada
func (calc *OptimizedTaxCalculation) CalculateTotal() {
	subTotal := calc.MontoNeto + calc.MontoExento + calc.MontoIVA
	calc.MontoTotal = roundSII(subTotal - calc.totalDesc + calc.totalRec - calc.TotalRetenciones)
}

// OptimizedCalculationHandler es el nuevo handler optimizado
func (h *TaxCalculationHandlers) OptimizedCalculationHandler(c *gin.Context) {
	var calc OptimizedTaxCalculation

	// Leer y validar la solicitud
	if err := c.BindJSON(&calc); err != nil {
		c.JSON(400, gin.H{
			"error":   "Error al procesar la solicitud",
			"codigo":  "CALC_000",
			"detalle": err.Error(),
		})
		return
	}

	// Validaciones y cálculos
	if err := calc.ValidateBase(); err != nil {
		c.JSON(400, gin.H{
			"error":   err.Message,
			"codigo":  err.Code,
			"detalle": err.Detail,
		})
		validationErrorPool.Put(err)
		return
	}

	if err := calc.CalculateIVA(); err != nil {
		c.JSON(400, gin.H{
			"error":   err.Message,
			"codigo":  err.Code,
			"detalle": err.Detail,
		})
		validationErrorPool.Put(err)
		return
	}

	if err := calc.CalculateDiscounts(); err != nil {
		c.JSON(400, gin.H{
			"error":   err.Message,
			"codigo":  err.Code,
			"detalle": err.Detail,
		})
		validationErrorPool.Put(err)
		return
	}

	if err := calc.CalculateSurcharges(); err != nil {
		c.JSON(400, gin.H{
			"error":   err.Message,
			"codigo":  err.Code,
			"detalle": err.Detail,
		})
		validationErrorPool.Put(err)
		return
	}

	if err := calc.CalculateRetentions(); err != nil {
		c.JSON(400, gin.H{
			"error":   err.Message,
			"codigo":  err.Code,
			"detalle": err.Detail,
		})
		validationErrorPool.Put(err)
		return
	}

	calc.CalculateTotal()

	// Preparar respuesta usando el pool de slices para el redondeo
	values := roundingPool.Get().([]float64)
	values = append(values,
		calc.MontoNeto,
		calc.MontoExento,
		calc.MontoIVA,
		calc.totalDesc,
		calc.totalRec,
		calc.TotalRetenciones,
		calc.MontoTotal,
		math.Round(calc.MontoTotal),
	)

	c.JSON(200, gin.H{
		"resultado": calc,
		"desglose": map[string]float64{
			"montoNeto":        values[0],
			"montoExento":      values[1],
			"montoIVA":         values[2],
			"descuentos":       values[3],
			"recargos":         values[4],
			"retenciones":      values[5],
			"montoTotal":       values[6],
			"montoTotalVisual": values[7],
		},
	})

	// Devolver el slice al pool
	values = values[:0]
	roundingPool.Put(values)
}
