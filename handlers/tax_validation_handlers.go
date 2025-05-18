package handlers

import (
	"fmt"
	"math"
	"strings"
	"time"

	"FMgo/api"

	"github.com/gin-gonic/gin"
)

type TaxValidationHandlers struct {
	client *api.FacturaMovilClient
}

type TaxDocument struct {
	TipoDocumento string    `json:"tipoDocumento"` // FACTURA, BOLETA, NOTA_CREDITO, NOTA_VENTA
	Folio         int       `json:"folio"`
	FechaEmision  time.Time `json:"fechaEmision"`
	MontoNeto     float64   `json:"montoNeto"`
	MontoIVA      float64   `json:"montoIVA"`
	MontoTotal    float64   `json:"montoTotal"`
	RUTEmisor     string    `json:"rutEmisor"`
	RUTReceptor   string    `json:"rutReceptor"`
	EstadoSII     string    `json:"estadoSII"`
}

func (h *TaxValidationHandlers) ValidateTaxDocumentHandler(c *gin.Context) {
	var doc TaxDocument

	// Obtener los datos del body de la solicitud
	if err := c.BindJSON(&doc); err != nil {
		c.JSON(400, gin.H{
			"error":   "Error al procesar la solicitud",
			"detalle": err.Error(),
		})
		return
	}

	// Validar que se proporcionó un tipo de documento
	if doc.TipoDocumento == "" {
		c.JSON(400, gin.H{
			"error":   "Tipo de documento requerido",
			"detalle": "Debe especificar el tipo de documento (FACTURA, NOTA_CREDITO, NOTA_VENTA)",
		})
		return
	}

	// Validaciones específicas por tipo de documento
	switch doc.TipoDocumento {
	case "FACTURA":
		if err := validateFactura(doc); err != nil {
			c.JSON(400, gin.H{
				"error":   "Error en validación de factura",
				"detalle": err.Error(),
			})
			return
		}
	case "NOTA_CREDITO":
		if err := validateNotaCredito(doc); err != nil {
			c.JSON(400, gin.H{
				"error":   "Error en validación de nota de crédito",
				"detalle": err.Error(),
			})
			return
		}
	case "NOTA_VENTA":
		if err := validateNotaVenta(doc); err != nil {
			c.JSON(400, gin.H{
				"error":   "Error en validación de nota de venta",
				"detalle": err.Error(),
			})
			return
		}
	default:
		c.JSON(400, gin.H{
			"error":   "Tipo de documento no soportado",
			"detalle": fmt.Sprintf("El tipo de documento '%s' no es válido", doc.TipoDocumento),
		})
		return
	}

	// Validación de montos
	if err := validateMontos(doc); err != nil {
		c.JSON(400, gin.H{
			"error":   "Error en validación de montos",
			"detalle": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"mensaje":   "Documento validado correctamente",
		"documento": doc,
	})
}

func validateFactura(doc TaxDocument) error {
	// Validación de RUTs
	if doc.RUTEmisor == "" {
		return fmt.Errorf("RUT emisor es obligatorio")
	}
	if doc.RUTReceptor == "" {
		return fmt.Errorf("RUT receptor es obligatorio")
	}

	// Validación de montos
	if doc.MontoNeto < 0 {
		return fmt.Errorf("monto neto no puede ser negativo")
	}
	if doc.MontoIVA < 0 {
		return fmt.Errorf("monto IVA no puede ser negativo")
	}
	if doc.MontoTotal <= 0 {
		return fmt.Errorf("monto total debe ser mayor a 0")
	}

	// Validación de cálculo de IVA con redondeo adecuado
	ivaCalculado := 0.0
	if doc.MontoNeto > 0 {
		ivaCalculado = math.Round((doc.MontoNeto*0.19)*100) / 100
	}

	const toleranciaIVA = 0.02 // Ajustado a 2 centavos para mayor compatibilidad
	if math.Abs(ivaCalculado-doc.MontoIVA) > toleranciaIVA {
		return fmt.Errorf("el IVA calculado (%.2f) no coincide con el monto IVA proporcionado (%.2f), diferencia: %.2f, tolerancia máxima: %.2f",
			ivaCalculado, doc.MontoIVA, math.Abs(ivaCalculado-doc.MontoIVA), toleranciaIVA)
	}

	// Validación de total con redondeo adecuado (considerando montos exentos)
	var montoExento float64 = doc.MontoTotal - doc.MontoNeto - doc.MontoIVA
	if montoExento < 0 && math.Abs(montoExento) > toleranciaIVA {
		return fmt.Errorf("el monto total es inconsistente con los montos neto e IVA (posible monto exento negativo): %.2f", montoExento)
	}

	totalCalculado := math.Round((doc.MontoNeto+doc.MontoIVA+montoExento)*100) / 100
	montoTotalRedondeado := math.Round(doc.MontoTotal*100) / 100
	const toleranciaTotal = 0.02 // Tolerancia para el total
	if math.Abs(totalCalculado-montoTotalRedondeado) > toleranciaTotal {
		return fmt.Errorf("el total calculado (%.2f) no coincide con el monto total proporcionado (%.2f), diferencia: %.2f, tolerancia máxima: %.2f",
			totalCalculado, montoTotalRedondeado, math.Abs(totalCalculado-montoTotalRedondeado), toleranciaTotal)
	}

	// Validación de fecha
	if doc.FechaEmision.IsZero() {
		return fmt.Errorf("fecha de emisión es obligatoria")
	}
	if doc.FechaEmision.After(time.Now()) {
		return fmt.Errorf("fecha de emisión no puede ser futura")
	}

	// Validación de folio
	if doc.Folio <= 0 {
		return fmt.Errorf("folio debe ser mayor a 0")
	}

	return nil
}

func validateNotaCredito(doc TaxDocument) error {
	// Validación de RUTs
	if doc.RUTEmisor == "" {
		return fmt.Errorf("RUT emisor es obligatorio")
	}
	if doc.RUTReceptor == "" {
		return fmt.Errorf("RUT receptor es obligatorio")
	}

	// Validación de montos
	if (doc.MontoNeto <= 0 && doc.MontoIVA > 0) || (doc.MontoTotal <= 0) {
		return fmt.Errorf("los montos de la nota de crédito son inválidos")
	}

	// Validación de cálculo de IVA con redondeo adecuado
	ivaCalculado := 0.0
	if doc.MontoNeto > 0 {
		ivaCalculado = math.Round((doc.MontoNeto*0.19)*100) / 100
	}

	const toleranciaIVA = 0.02 // Ajustado a 2 centavos para mayor compatibilidad
	if math.Abs(ivaCalculado-doc.MontoIVA) > toleranciaIVA {
		return fmt.Errorf("el IVA calculado (%.2f) no coincide con el monto IVA proporcionado (%.2f), diferencia: %.2f, tolerancia máxima: %.2f",
			ivaCalculado, doc.MontoIVA, math.Abs(ivaCalculado-doc.MontoIVA), toleranciaIVA)
	}

	// Validación de total con redondeo adecuado (considerando montos exentos)
	var montoExento float64 = doc.MontoTotal - doc.MontoNeto - doc.MontoIVA
	if montoExento < 0 && math.Abs(montoExento) > toleranciaIVA {
		return fmt.Errorf("el monto total es inconsistente con los montos neto e IVA (posible monto exento negativo): %.2f", montoExento)
	}

	totalCalculado := math.Round((doc.MontoNeto+doc.MontoIVA+montoExento)*100) / 100
	montoTotalRedondeado := math.Round(doc.MontoTotal*100) / 100
	const toleranciaTotal = 0.02 // Tolerancia para el total
	if math.Abs(totalCalculado-montoTotalRedondeado) > toleranciaTotal {
		return fmt.Errorf("el total calculado (%.2f) no coincide con el monto total proporcionado (%.2f), diferencia: %.2f, tolerancia máxima: %.2f",
			totalCalculado, montoTotalRedondeado, math.Abs(totalCalculado-montoTotalRedondeado), toleranciaTotal)
	}

	// Validación de fecha
	if doc.FechaEmision.IsZero() {
		return fmt.Errorf("fecha de emisión es obligatoria")
	}
	if doc.FechaEmision.After(time.Now()) {
		return fmt.Errorf("fecha de emisión no puede ser futura")
	}

	// Validación de folio
	if doc.Folio <= 0 {
		return fmt.Errorf("folio debe ser mayor a 0")
	}

	return nil
}

func validateNotaVenta(doc TaxDocument) error {
	// Validación de RUTs
	if doc.RUTEmisor == "" {
		return fmt.Errorf("RUT emisor es obligatorio")
	}
	if doc.RUTReceptor == "" {
		return fmt.Errorf("RUT receptor es obligatorio")
	}

	// Validación de montos
	if doc.MontoNeto <= 0 {
		return fmt.Errorf("monto neto debe ser mayor a 0")
	}
	if doc.MontoIVA < 0 {
		return fmt.Errorf("monto IVA no puede ser negativo")
	}
	if doc.MontoTotal <= 0 {
		return fmt.Errorf("monto total debe ser mayor a 0")
	}

	// Validación de cálculo de IVA con redondeo adecuado
	ivaCalculado := math.Round((doc.MontoNeto*0.19)*100) / 100
	const toleranciaIVA = 0.02 // Ajustado a 2 centavos para mayor compatibilidad
	if math.Abs(ivaCalculado-doc.MontoIVA) > toleranciaIVA {
		return fmt.Errorf("el IVA calculado (%.2f) no coincide con el monto IVA proporcionado (%.2f), diferencia: %.2f, tolerancia máxima: %.2f",
			ivaCalculado, doc.MontoIVA, math.Abs(ivaCalculado-doc.MontoIVA), toleranciaIVA)
	}

	// Validación de total con redondeo adecuado (considerando montos exentos)
	var montoExento float64 = doc.MontoTotal - doc.MontoNeto - doc.MontoIVA
	if montoExento < 0 && math.Abs(montoExento) > toleranciaIVA {
		return fmt.Errorf("el monto total es inconsistente con los montos neto e IVA (posible monto exento negativo): %.2f", montoExento)
	}

	totalCalculado := math.Round((doc.MontoNeto+doc.MontoIVA+montoExento)*100) / 100
	montoTotalRedondeado := math.Round(doc.MontoTotal*100) / 100
	const toleranciaTotal = 0.02 // Tolerancia para el total
	if math.Abs(totalCalculado-montoTotalRedondeado) > toleranciaTotal {
		return fmt.Errorf("el total calculado (%.2f) no coincide con el monto total proporcionado (%.2f), diferencia: %.2f, tolerancia máxima: %.2f",
			totalCalculado, montoTotalRedondeado, math.Abs(totalCalculado-montoTotalRedondeado), toleranciaTotal)
	}

	// Validación de fecha
	if doc.FechaEmision.IsZero() {
		return fmt.Errorf("fecha de emisión es obligatoria")
	}
	if doc.FechaEmision.After(time.Now()) {
		return fmt.Errorf("fecha de emisión no puede ser futura")
	}

	// Validación de folio
	if doc.Folio <= 0 {
		return fmt.Errorf("folio debe ser mayor a 0")
	}

	return nil
}

func validateMontos(doc TaxDocument) error {
	// Validación de rangos para montos
	if doc.MontoNeto > 1000000000 {
		return fmt.Errorf("monto neto excede el límite máximo permitido (1.000.000.000)")
	}

	if doc.MontoIVA > 190000000 {
		return fmt.Errorf("monto IVA excede el límite máximo permitido (190.000.000)")
	}

	if doc.MontoTotal > 1190000000 {
		return fmt.Errorf("monto total excede el límite máximo permitido (1.190.000.000)")
	}

	// Validación de proporcionalidad del IVA
	if doc.MontoIVA > 0 && doc.MontoNeto > 0 {
		// Si hay IVA y monto neto, verificar que el IVA esté en el rango esperado (19% +/- 1%)
		porcentajeIVA := (doc.MontoIVA / doc.MontoNeto) * 100
		if porcentajeIVA < 18 || porcentajeIVA > 20 {
			return fmt.Errorf("el porcentaje de IVA calculado (%.2f%%) está fuera del rango esperado (19%% +/- 1%%)",
				porcentajeIVA)
		}
	} else if doc.MontoIVA > 0 && doc.MontoNeto <= 0 {
		// Si hay IVA pero no hay monto neto, es un error
		return fmt.Errorf("existe monto de IVA (%.2f) pero el monto neto es cero o negativo", doc.MontoIVA)
	}

	// Validación de decimales excesivos
	if !validarDecimales(doc.MontoNeto, 2) {
		return fmt.Errorf("monto neto tiene más de 2 decimales")
	}

	if !validarDecimales(doc.MontoIVA, 2) {
		return fmt.Errorf("monto IVA tiene más de 2 decimales")
	}

	if !validarDecimales(doc.MontoTotal, 2) {
		return fmt.Errorf("monto total tiene más de 2 decimales")
	}

	return nil
}

// validarDecimales verifica si un número tiene un número máximo de decimales
func validarDecimales(numero float64, maxDecimales int) bool {
	// Convertir a string con alta precisión para contar decimales
	str := fmt.Sprintf("%f", numero)
	partes := strings.Split(str, ".")
	if len(partes) != 2 {
		return true // No tiene decimales
	}

	// Contar decimales significativos (eliminar ceros a la derecha)
	decimales := strings.TrimRight(partes[1], "0")
	return len(decimales) <= maxDecimales
}
