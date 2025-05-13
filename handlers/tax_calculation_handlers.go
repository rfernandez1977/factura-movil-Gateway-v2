package handlers

import (
	"fmt"
	"math"

	"github.com/cursor/FMgo/api"

	"github.com/gin-gonic/gin"
)

type TaxCalculationHandlers struct {
	client *api.FacturaMovilClient
}

type CalculoTributario struct {
	// Datos base
	MontoNeto   float64     `json:"montoNeto"`
	MontoExento float64     `json:"montoExento"`
	TasaIVA     float64     `json:"tasaIVA"`
	Descuentos  []Descuento `json:"descuentos"`
	Recargos    []Recargo   `json:"recargos"`

	// Retenciones específicas
	RetencionHonorarios float64 `json:"retencionHonorarios,omitempty"`
	RetencionILA        float64 `json:"retencionILA,omitempty"`

	// Resultados calculados
	MontoIVA         float64 `json:"montoIVA"`
	MontoTotal       float64 `json:"montoTotal"`
	TotalRetenciones float64 `json:"totalRetenciones"`
}

type Descuento struct {
	TipoDescuento  string  `json:"tipoDescuento"` // GLOBAL, POR_ITEM
	Porcentaje     float64 `json:"porcentaje"`
	MontoDescuento float64 `json:"montoDescuento"`
	Glosa          string  `json:"glosa"`
}

type Recargo struct {
	TipoRecargo  string  `json:"tipoRecargo"`
	Porcentaje   float64 `json:"porcentaje"`
	MontoRecargo float64 `json:"montoRecargo"`
	Glosa        string  `json:"glosa"`
}

func (h *TaxCalculationHandlers) CalculateHandler(c *gin.Context) {
	var calculo CalculoTributario

	// Leer el cuerpo de la solicitud
	if err := c.BindJSON(&calculo); err != nil {
		c.JSON(400, gin.H{
			"error":   "Error al procesar la solicitud",
			"codigo":  "CALC_000",
			"detalle": err.Error(),
		})
		return
	}

	// Validación inicial de montos
	if calculo.MontoNeto < 0 || calculo.MontoExento < 0 {
		c.JSON(400, gin.H{
			"error":   "Los montos no pueden ser negativos",
			"codigo":  "CALC_001",
			"detalle": "Monto neto y monto exento deben ser positivos o cero",
		})
		return
	}

	// Validar tasa de IVA
	if calculo.TasaIVA != 19 && calculo.TasaIVA != 0 {
		c.JSON(400, gin.H{
			"error":   "Tasa de IVA inválida",
			"codigo":  "CALC_002",
			"detalle": "La tasa de IVA debe ser 19% o 0%",
		})
		return
	}

	// Cálculo de IVA con redondeo
	if calculo.MontoNeto > 0 && calculo.TasaIVA > 0 {
		calculo.MontoIVA = math.Round((calculo.MontoNeto*(calculo.TasaIVA/100))*100) / 100
	} else {
		calculo.MontoIVA = 0
		// Si no hay monto neto pero hay tasa IVA positiva, es mejor no generar un monto de IVA
		if calculo.TasaIVA > 0 && calculo.MontoNeto <= 0 {
			c.JSON(400, gin.H{
				"error":   "Inconsistencia en cálculo de IVA",
				"codigo":  "CALC_009",
				"detalle": "No se puede calcular IVA con monto neto cero o negativo",
			})
			return
		}
	}

	// Aplicar descuentos
	totalDescuentos := 0.0
	for i, desc := range calculo.Descuentos {
		// Validar porcentaje de descuento
		if desc.Porcentaje < 0 || desc.Porcentaje > 100 {
			c.JSON(400, gin.H{
				"error":   "Porcentaje de descuento inválido",
				"codigo":  "CALC_004",
				"detalle": "El porcentaje debe estar entre 0 y 100",
			})
			return
		}

		if desc.TipoDescuento == "GLOBAL" {
			montoDesc := math.Round(((calculo.MontoNeto+calculo.MontoExento)*(desc.Porcentaje/100))*100) / 100
			totalDescuentos += montoDesc
			calculo.Descuentos[i].MontoDescuento = montoDesc
		}
	}

	// Aplicar recargos
	totalRecargos := 0.0
	for i, rec := range calculo.Recargos {
		// Validar porcentaje de recargo
		if rec.Porcentaje < 0 || rec.Porcentaje > 100 {
			c.JSON(400, gin.H{
				"error":   "Porcentaje de recargo inválido",
				"codigo":  "CALC_005",
				"detalle": "El porcentaje debe estar entre 0 y 100",
			})
			return
		}

		if rec.TipoRecargo == "GLOBAL" {
			montoRec := math.Round(((calculo.MontoNeto+calculo.MontoExento)*(rec.Porcentaje/100))*100) / 100
			totalRecargos += montoRec
			calculo.Recargos[i].MontoRecargo = montoRec
		}
	}

	// Cálculo de retenciones con validación
	calculo.TotalRetenciones = 0

	if calculo.RetencionHonorarios > 0 {
		if calculo.RetencionHonorarios != 10.75 {
			c.JSON(400, gin.H{
				"error":   "Tasa de retención de honorarios inválida",
				"codigo":  "CALC_003",
				"detalle": "La retención de honorarios debe ser 10.75%",
			})
			return
		}
		retencion := math.Round((calculo.MontoNeto*(calculo.RetencionHonorarios/100))*100) / 100
		calculo.TotalRetenciones += retencion
	}

	if calculo.RetencionILA > 0 {
		if calculo.RetencionILA != 10 && calculo.RetencionILA != 11 && calculo.RetencionILA != 13 {
			c.JSON(400, gin.H{
				"error":   "Tasa de retención ILA inválida",
				"codigo":  "CALC_006",
				"detalle": "La retención ILA debe ser 10%, 11% o 13%",
			})
			return
		}
		retencion := math.Round((calculo.MontoNeto*(calculo.RetencionILA/100))*100) / 100
		calculo.TotalRetenciones += retencion
	}

	// Cálculo del monto total con redondeo
	subTotal := calculo.MontoNeto + calculo.MontoExento + calculo.MontoIVA
	calculo.MontoTotal = math.Round((subTotal-totalDescuentos+totalRecargos-calculo.TotalRetenciones)*100) / 100

	// Redondeo final según normativa SII:
	// - Para documentos electrónicos, se mantienen 2 decimales
	// - Solo para efectos de visualización en el documento impreso se redondea a entero
	// Por lo tanto, para los cálculos seguimos manteniendo la precisión de 2 decimales
	montoTotalVisual := math.Round(calculo.MontoTotal)

	c.JSON(200, gin.H{
		"resultado": calculo,
		"desglose": map[string]float64{
			"montoNeto":        calculo.MontoNeto,
			"montoExento":      calculo.MontoExento,
			"montoIVA":         calculo.MontoIVA,
			"descuentos":       totalDescuentos,
			"recargos":         totalRecargos,
			"retenciones":      calculo.TotalRetenciones,
			"montoTotal":       calculo.MontoTotal,
			"montoTotalVisual": montoTotalVisual, // Monto redondeado a entero para mostrar en documentos
		},
	})
}

func (h *TaxCalculationHandlers) ValidateCalculationHandler(c *gin.Context) {
	var calculo CalculoTributario

	// Leer el cuerpo de la solicitud
	if err := c.BindJSON(&calculo); err != nil {
		c.JSON(400, gin.H{
			"error":   "Error al procesar la solicitud",
			"codigo":  "CALC_000",
			"detalle": err.Error(),
		})
		return
	}

	// Validación de montos
	if calculo.MontoNeto < 0 {
		c.JSON(400, gin.H{
			"error":   "Monto neto inválido",
			"codigo":  "CALC_001",
			"detalle": "El monto neto no puede ser negativo",
		})
		return
	}

	if calculo.MontoExento < 0 {
		c.JSON(400, gin.H{
			"error":   "Monto exento inválido",
			"codigo":  "CALC_001",
			"detalle": "El monto exento no puede ser negativo",
		})
		return
	}

	if calculo.MontoNeto == 0 && calculo.MontoExento == 0 {
		c.JSON(400, gin.H{
			"error":   "Montos inválidos",
			"codigo":  "CALC_001",
			"detalle": "El documento debe tener al menos monto neto o monto exento mayor a cero",
		})
		return
	}

	// Validación de tasas de IVA
	if calculo.TasaIVA != 19 && calculo.TasaIVA != 0 {
		c.JSON(400, gin.H{
			"error":   "Tasa de IVA inválida",
			"codigo":  "CALC_002",
			"detalle": "La tasa de IVA debe ser 19% o 0%",
		})
		return
	}

	// Validación de retenciones
	if calculo.RetencionHonorarios > 0 && calculo.RetencionHonorarios != 10.75 {
		c.JSON(400, gin.H{
			"error":   "Tasa de retención de honorarios inválida",
			"codigo":  "CALC_003",
			"detalle": "La retención de honorarios debe ser 10.75%",
		})
		return
	}

	if calculo.RetencionILA > 0 && (calculo.RetencionILA != 10 && calculo.RetencionILA != 11 && calculo.RetencionILA != 13) {
		c.JSON(400, gin.H{
			"error":   "Tasa de retención ILA inválida",
			"codigo":  "CALC_006",
			"detalle": "La retención ILA debe ser 10%, 11% o 13%",
		})
		return
	}

	// Validación de descuentos
	totalDescuentos := 0.0
	for _, desc := range calculo.Descuentos {
		if desc.Porcentaje < 0 || desc.Porcentaje > 100 {
			c.JSON(400, gin.H{
				"error":   "Porcentaje de descuento inválido",
				"codigo":  "CALC_004",
				"detalle": "El porcentaje debe estar entre 0 y 100",
			})
			return
		}

		if desc.MontoDescuento < 0 {
			c.JSON(400, gin.H{
				"error":   "Monto de descuento inválido",
				"codigo":  "CALC_004",
				"detalle": "El monto de descuento no puede ser negativo",
			})
			return
		}

		totalDescuentos += desc.MontoDescuento
	}

	// Validación de recargos
	totalRecargos := 0.0
	for _, rec := range calculo.Recargos {
		if rec.Porcentaje < 0 || rec.Porcentaje > 100 {
			c.JSON(400, gin.H{
				"error":   "Porcentaje de recargo inválido",
				"codigo":  "CALC_005",
				"detalle": "El porcentaje debe estar entre 0 y 100",
			})
			return
		}

		if rec.MontoRecargo < 0 {
			c.JSON(400, gin.H{
				"error":   "Monto de recargo inválido",
				"codigo":  "CALC_005",
				"detalle": "El monto de recargo no puede ser negativo",
			})
			return
		}

		totalRecargos += rec.MontoRecargo
	}

	// Verificar cálculo de IVA
	ivaEsperado := 0.0
	if calculo.MontoNeto > 0 && calculo.TasaIVA > 0 {
		ivaEsperado = math.Round((calculo.MontoNeto*(calculo.TasaIVA/100))*100) / 100

		// Tolerancia para comparación de IVA
		const toleranciaIVA = 0.02
		if math.Abs(ivaEsperado-calculo.MontoIVA) > toleranciaIVA {
			c.JSON(400, gin.H{
				"error":  "Cálculo de IVA incorrecto",
				"codigo": "CALC_007",
				"detalle": fmt.Sprintf("El IVA calculado debería ser %.2f, pero se proporcionó %.2f (diferencia: %.2f, tolerancia máxima: %.2f)",
					ivaEsperado, calculo.MontoIVA, math.Abs(ivaEsperado-calculo.MontoIVA), toleranciaIVA),
			})
			return
		}
	} else if calculo.MontoIVA > 0 {
		// Si no hay monto neto o tasa de IVA pero hay IVA, es un error
		c.JSON(400, gin.H{
			"error":   "Inconsistencia en IVA",
			"codigo":  "CALC_009",
			"detalle": fmt.Sprintf("Se proporcionó un monto de IVA (%.2f) pero el monto neto es cero o la tasa de IVA es cero", calculo.MontoIVA),
		})
		return
	}

	// Verificar cálculo del monto total
	const toleranciaTotal = 0.02
	totalEsperado := math.Round((calculo.MontoNeto+calculo.MontoExento+calculo.MontoIVA-totalDescuentos+totalRecargos-calculo.TotalRetenciones)*100) / 100

	if math.Abs(totalEsperado-calculo.MontoTotal) > toleranciaTotal {
		c.JSON(400, gin.H{
			"error":  "Cálculo del monto total incorrecto",
			"codigo": "CALC_008",
			"detalle": fmt.Sprintf("El monto total debería ser %.2f, pero se proporcionó %.2f (diferencia: %.2f, tolerancia máxima: %.2f)",
				totalEsperado, calculo.MontoTotal, math.Abs(totalEsperado-calculo.MontoTotal), toleranciaTotal),
		})
		return
	}

	// Redondeo visual para documentos impresos
	montoTotalVisual := math.Round(totalEsperado)

	c.JSON(200, gin.H{
		"mensaje": "Cálculos validados correctamente",
		"calculo": calculo,
		"montos_validados": map[string]float64{
			"montoNeto":        calculo.MontoNeto,
			"montoExento":      calculo.MontoExento,
			"montoIVA":         ivaEsperado,
			"descuentos":       totalDescuentos,
			"recargos":         totalRecargos,
			"retenciones":      calculo.TotalRetenciones,
			"montoTotal":       totalEsperado,
			"montoTotalVisual": montoTotalVisual, // Monto redondeado a entero para mostrar en documentos
		},
	})
}
