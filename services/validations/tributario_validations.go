package validations

import (
	"errors"
	"fmt"
	"math"

	"github.com/cursor/FMgo/domain"
	"github.com/cursor/FMgo/models"
)

// TributarioValidation maneja las validaciones de negocio para documentos tributarios
type TributarioValidation struct {
	config *ConfigValidacionTributario
}

// ConfigValidacionTributario contiene la configuración para las validaciones tributarias
type ConfigValidacionTributario struct {
	MaxMontoTotal        float64
	MaxItems             int
	MaxDiasAntiguedad    int
	PorcentajeIVA        float64
	PorcentajeRetencion  float64
	MontoMinimoRetencion float64
}

// NewTributarioValidation crea una nueva instancia del servicio de validaciones tributarias
func NewTributarioValidation() *TributarioValidation {
	return &TributarioValidation{
		config: &ConfigValidacionTributario{
			MaxMontoTotal:        1000000000, // 1 billón
			MaxItems:             1000,
			MaxDiasAntiguedad:    365,
			PorcentajeIVA:        0.19,
			PorcentajeRetencion:  0.10,
			MontoMinimoRetencion: 1000000, // 1 millón
		},
	}
}

// ValidarDocumento valida un documento tributario según las reglas de negocio
func (v *TributarioValidation) ValidarDocumento(doc interface{}) error {
	switch d := doc.(type) {
	case *models.Factura:
		return v.validarFactura(d)
	case *models.Boleta:
		return v.validarBoleta(d)
	case *models.NotaCredito:
		return v.validarNotaCredito(d)
	case *models.NotaDebito:
		return v.validarNotaDebito(d)
	case *models.GuiaDespacho:
		return v.validarGuiaDespacho(d)
	default:
		return errors.New("tipo de documento no soportado")
	}
}

// validarFactura valida una factura según las reglas de negocio
func (v *TributarioValidation) validarFactura(factura *models.Factura) error {
	// Validar montos
	if err := v.validarMontos(factura); err != nil {
		return err
	}

	// Validar referencias cruzadas
	if err := v.validarReferenciasFactura(factura); err != nil {
		return err
	}

	// Validar estado tributario
	if err := v.validarEstadoTributario(factura); err != nil {
		return err
	}

	// Validar cálculos de impuestos
	if err := v.validarCalculosImpuestos(factura); err != nil {
		return err
	}

	return nil
}

// validarMontos valida los montos de un documento
func (v *TributarioValidation) validarMontos(doc interface{}) error {
	var montoTotal float64
	var items int

	switch d := doc.(type) {
	case *models.Factura:
		montoTotal = d.MontoTotal
		items = len(d.Items)
	case *models.Boleta:
		montoTotal = d.MontoTotal
		items = len(d.Items)
	case *models.NotaCredito:
		montoTotal = d.MontoTotal
		items = len(d.Items)
	case *models.NotaDebito:
		montoTotal = d.MontoTotal
		items = len(d.Items)
	case *models.GuiaDespacho:
		montoTotal = d.MontoTotal
		items = len(d.Items)
	}

	if montoTotal > v.config.MaxMontoTotal {
		return fmt.Errorf("monto total excede el máximo permitido: %v", v.config.MaxMontoTotal)
	}

	if items > v.config.MaxItems {
		return fmt.Errorf("número de ítems excede el máximo permitido: %v", v.config.MaxItems)
	}

	return nil
}

// validarReferenciasFactura valida las referencias cruzadas de una factura
func (v *TributarioValidation) validarReferenciasFactura(factura *models.Factura) error {
	// Validar referencias a guías de despacho
	for _, ref := range factura.Referencias {
		if ref.TipoDocumento == models.TipoGuiaDespacho {
			// Verificar que la guía de despacho existe y está en estado correcto
			if err := v.validarExistenciaGuiaDespacho(ref.FolioReferencia); err != nil {
				return err
			}
		}
	}

	// Validar referencias a notas de crédito/débito
	for _, ref := range factura.Referencias {
		if ref.TipoDocumento == models.TipoNotaCredito || ref.TipoDocumento == models.TipoNotaDebito {
			// Verificar que la nota existe y está en estado correcto
			if err := v.validarExistenciaNota(ref.FolioReferencia); err != nil {
				return err
			}
		}
	}

	return nil
}

// validarEstadoTributario valida el estado tributario de un documento
func (v *TributarioValidation) validarEstadoTributario(doc interface{}) error {
	// Verificar que el emisor esté al día en sus obligaciones tributarias
	if err := v.verificarEstadoTributarioEmisor(doc); err != nil {
		return err
	}

	// Verificar que el receptor esté al día en sus obligaciones tributarias
	if err := v.verificarEstadoTributarioReceptor(doc); err != nil {
		return err
	}

	return nil
}

// validarCalculosImpuestos valida los cálculos de impuestos de un documento
func (v *TributarioValidation) validarCalculosImpuestos(doc interface{}) error {
	var (
		montoNeto                 float64
		montoIVA                  float64
		montoTotal                float64
		montoExento               float64
		totalImpuestosAdicionales float64
	)

	// Tolerancia para comparaciones de montos
	const toleranciaMontos = 0.02 // 2 centavos de tolerancia

	switch d := doc.(type) {
	case *models.Factura:
		montoNeto = d.MontoNeto
		montoIVA = d.MontoIVA
		montoTotal = d.MontoTotal
		montoExento = d.MontoExento
		// Los impuestos adicionales se calcularán de los items
		for _, item := range d.Items {
			for _, impuesto := range item.ImpuestosAdicionales {
				totalImpuestosAdicionales += impuesto.Monto

				// Validar código de impuesto
				if !v.validarCodigoImpuesto(impuesto.Codigo) {
					return fmt.Errorf("código de impuesto no válido: %s", impuesto.Codigo)
				}

				// Validar porcentaje
				if impuesto.Porcentaje < 0 || impuesto.Porcentaje > 100 {
					return fmt.Errorf("porcentaje de impuesto no válido: %.2f", impuesto.Porcentaje)
				}

				// Validar base imponible
				if impuesto.BaseImponible < 0 {
					return fmt.Errorf("base imponible no puede ser negativa: %.2f", impuesto.BaseImponible)
				}

				// Validar monto calculado con tolerancia
				montoCalculado := math.Round((impuesto.BaseImponible*(impuesto.Porcentaje/100))*100) / 100
				if math.Abs(montoCalculado-impuesto.Monto) > toleranciaMontos {
					return fmt.Errorf("el monto del impuesto calculado (%.2f) no coincide con el monto proporcionado (%.2f), diferencia: %.2f, tolerancia máxima: %.2f",
						montoCalculado, impuesto.Monto, math.Abs(montoCalculado-impuesto.Monto), toleranciaMontos)
				}
			}
		}
	case *models.Boleta:
		montoNeto = d.MontoNeto
		montoIVA = d.MontoIVA
		montoTotal = d.MontoTotal
		montoExento = d.MontoExento
		// Los impuestos adicionales se calcularán de los items
		for _, item := range d.Items {
			for _, impuesto := range item.ImpuestosAdicionales {
				totalImpuestosAdicionales += impuesto.Monto
			}
		}
	case *models.NotaCredito:
		montoNeto = d.MontoNeto
		montoIVA = d.MontoIVA
		montoTotal = d.MontoTotal
		montoExento = d.MontoExento
		// Los impuestos adicionales se calcularán de los items
		for _, item := range d.Items {
			for _, impuesto := range item.ImpuestosAdicionales {
				totalImpuestosAdicionales += impuesto.Monto
			}
		}
	case *models.NotaDebito:
		montoNeto = d.MontoNeto
		montoIVA = d.MontoIVA
		montoTotal = d.MontoTotal
		montoExento = d.MontoExento
		// Los impuestos adicionales se calcularán de los items
		for _, item := range d.Items {
			for _, impuesto := range item.ImpuestosAdicionales {
				totalImpuestosAdicionales += impuesto.Monto
			}
		}
	case *models.GuiaDespacho:
		montoNeto = d.MontoNeto
		montoIVA = d.MontoIVA
		montoTotal = d.MontoTotal
		montoExento = d.MontoExento
		// Los impuestos adicionales se calcularán de los items
		for _, item := range d.Items {
			for _, impuesto := range item.ImpuestosAdicionales {
				totalImpuestosAdicionales += impuesto.Monto
			}
		}
	}

	// Validar IVA solo si hay monto neto afecto a IVA
	if montoNeto > 0 {
		ivaCalculado := math.Round((montoNeto*v.config.PorcentajeIVA)*100) / 100
		if math.Abs(ivaCalculado-montoIVA) > toleranciaMontos {
			return fmt.Errorf("el IVA calculado (%.2f) no coincide con el monto IVA proporcionado (%.2f), diferencia: %.2f, tolerancia máxima: %.2f",
				ivaCalculado, montoIVA, math.Abs(ivaCalculado-montoIVA), toleranciaMontos)
		}
	} else if montoIVA > 0 {
		// Si no hay monto neto pero hay IVA, es un error
		return fmt.Errorf("se ha proporcionado un monto de IVA (%.2f) pero el monto neto es cero o negativo", montoIVA)
	}

	// Validar total con tolerancia
	totalCalculado := math.Round((montoNeto+montoExento+montoIVA+totalImpuestosAdicionales)*100) / 100
	if math.Abs(totalCalculado-montoTotal) > toleranciaMontos {
		return fmt.Errorf("el total calculado (%.2f) no coincide con el monto total proporcionado (%.2f), diferencia: %.2f, tolerancia máxima: %.2f",
			totalCalculado, montoTotal, math.Abs(totalCalculado-montoTotal), toleranciaMontos)
	}

	return nil
}

// validarCodigoImpuesto valida si un código de impuesto es válido según el SII
func (v *TributarioValidation) validarCodigoImpuesto(codigo string) bool {
	// Mapa de códigos de impuestos adicionales según la normativa SII
	codigosValidos := map[string]bool{
		"14": true, // IVA anticipado faenamiento carne
		"15": true, // IVA anticipado carne
		"17": true, // Impuesto a las bebidas analcohólicas
		"18": true, // Impuesto a las bebidas alcohólicas
		"19": true, // IVA
		"23": true, // Impuesto específico a los combustibles
		"24": true, // Impuesto específico derecho de extracción
		"25": true, // Impuesto específico joyas y piedras preciosas
		"26": true, // Impuesto específico alfombras, tapices
		"27": true, // Impuesto específico vehículos fuera de catálogo
		"28": true, // Impuesto específico vehículos casa rodante, aeronaves, embarcaciones recreo
		"30": true, // Impuesto puros
		"31": true, // Impuesto cigarrillos
		"32": true, // Impuesto tabacos elaborados
		"33": true, // Impuesto pirotecnia
		"45": true, // Impuesto adicional bienes suntuarios
		"46": true, // Impuesto bebidas alcohólicas
		"47": true, // Impuesto particular diesel
		"48": true, // Impuesto único a los combustibles
		"49": true, // IVA compra activo fijo
		"50": true, // IVA legítimo
		"51": true, // IVA turista
	}

	return codigosValidos[codigo]
}

// validarExistenciaGuiaDespacho verifica la existencia y estado de una guía de despacho
func (v *TributarioValidation) validarExistenciaGuiaDespacho(folio int) error {
	// TODO: Implementar verificación real contra base de datos
	return nil
}

// validarExistenciaNota verifica la existencia y estado de una nota de crédito/débito
func (v *TributarioValidation) validarExistenciaNota(folio int) error {
	// TODO: Implementar verificación real contra base de datos
	return nil
}

// verificarEstadoTributarioEmisor verifica el estado tributario del emisor
func (v *TributarioValidation) verificarEstadoTributarioEmisor(doc interface{}) error {
	// TODO: Implementar verificación real contra SII
	return nil
}

// verificarEstadoTributarioReceptor verifica el estado tributario del receptor
func (v *TributarioValidation) verificarEstadoTributarioReceptor(doc interface{}) error {
	// TODO: Implementar verificación real contra SII
	return nil
}

// validarBoleta valida una boleta según las reglas de negocio
func (v *TributarioValidation) validarBoleta(boleta *models.Boleta) error {
	// Validar montos
	if err := v.validarMontos(boleta); err != nil {
		return err
	}

	// Validar estado tributario
	if err := v.validarEstadoTributario(boleta); err != nil {
		return err
	}

	// Validar cálculos de impuestos
	if err := v.validarCalculosImpuestos(boleta); err != nil {
		return err
	}

	return nil
}

// validarNotaCredito valida una nota de crédito según las reglas de negocio
func (v *TributarioValidation) validarNotaCredito(nota *models.NotaCredito) error {
	// Validar montos
	if err := v.validarMontos(nota); err != nil {
		return err
	}

	// Validar referencias
	if err := v.validarReferenciasNota(nota); err != nil {
		return err
	}

	// Validar estado tributario
	if err := v.validarEstadoTributario(nota); err != nil {
		return err
	}

	// Validar cálculos de impuestos
	if err := v.validarCalculosImpuestos(nota); err != nil {
		return err
	}

	return nil
}

// validarNotaDebito valida una nota de débito según las reglas de negocio
func (v *TributarioValidation) validarNotaDebito(nota *models.NotaDebito) error {
	// Validar montos
	if err := v.validarMontos(nota); err != nil {
		return err
	}

	// Validar referencias
	if err := v.validarReferenciasNota(nota); err != nil {
		return err
	}

	// Validar estado tributario
	if err := v.validarEstadoTributario(nota); err != nil {
		return err
	}

	// Validar cálculos de impuestos
	if err := v.validarCalculosImpuestos(nota); err != nil {
		return err
	}

	return nil
}

// validarGuiaDespacho valida una guía de despacho según las reglas de negocio
func (v *TributarioValidation) validarGuiaDespacho(guia *models.GuiaDespacho) error {
	// Validar montos
	if err := v.validarMontos(guia); err != nil {
		return err
	}

	// Validar estado tributario
	if err := v.validarEstadoTributario(guia); err != nil {
		return err
	}

	// Validar cálculos de impuestos
	if err := v.validarCalculosImpuestos(guia); err != nil {
		return err
	}

	return nil
}

// validarReferenciasNota valida las referencias de una nota de crédito/débito
func (v *TributarioValidation) validarReferenciasNota(nota interface{}) error {
	var folioReferencia int
	var tipoReferencia models.TipoReferencia

	switch n := nota.(type) {
	case *models.NotaCredito:
		folioReferencia = int(n.FolioReferencia)
		tipoReferencia = models.TipoReferencia(n.TipoReferencia)
	case *models.NotaDebito:
		folioReferencia = int(n.FolioReferencia)
		tipoReferencia = models.TipoReferencia(n.TipoReferencia)
	}

	// Verificar que el tipo de referencia sea válido
	validTipos := map[models.TipoReferencia]bool{
		models.TipoReferenciaAnulacion:    true,
		models.TipoReferenciaCorreccion:   true,
		models.TipoReferenciaDevolucion:   true,
		models.TipoReferenciaFactura:      true,
		models.TipoReferenciaGuiaDespacho: true,
	}

	if !validTipos[tipoReferencia] {
		return fmt.Errorf("tipo de referencia no válido: %s", tipoReferencia)
	}

	// Verificar que el documento referenciado existe
	if err := v.validarExistenciaDocumento(folioReferencia, string(tipoReferencia)); err != nil {
		return err
	}

	// Verificar que el documento referenciado no está anulado
	if err := v.verificarDocumentoNoAnulado(folioReferencia, string(tipoReferencia)); err != nil {
		return err
	}

	return nil
}

// validarExistenciaDocumento verifica la existencia de un documento
func (v *TributarioValidation) validarExistenciaDocumento(folio int, tipo string) error {
	// TODO: Implementar verificación real contra base de datos
	return nil
}

// verificarDocumentoNoAnulado verifica que un documento no esté anulado
func (v *TributarioValidation) verificarDocumentoNoAnulado(folio int, tipo string) error {
	// TODO: Implementar verificación real contra base de datos
	return nil
}

// ValidarReferencias valida las referencias de una factura
func (s *TributarioValidation) ValidarReferencias(factura *models.Factura) []models.ValidationFieldError {
	var errors []models.ValidationFieldError

	if factura.Referencias == nil {
		return errors
	}

	for i, ref := range factura.Referencias {
		// Validar que el tipo de referencia sea válido
		if ref.TipoReferencia == "" {
			errors = append(errors, models.ValidationFieldError{
				Field:   fmt.Sprintf("referencias[%d].tipo_referencia", i),
				Message: "Tipo de referencia es requerido",
			})
		}

		// Validar referencias a guías de despacho
		if ref.TipoReferencia == models.TipoGuiaDespacho && ref.TipoDocumento != fmt.Sprintf("%d", models.TipoGuiaDespacho) {
			errors = append(errors, models.ValidationFieldError{
				Field:   fmt.Sprintf("referencias[%d].tipo_documento", i),
				Message: "La referencia debe ser a una guía de despacho",
			})
		}

		// Validar referencias a notas de crédito/débito
		if (ref.TipoReferencia == models.TipoNotaCredito || ref.TipoReferencia == models.TipoNotaDebito) &&
			(ref.TipoDocumento != fmt.Sprintf("%d", models.TipoNotaCredito) && ref.TipoDocumento != fmt.Sprintf("%d", models.TipoNotaDebito)) {
			errors = append(errors, models.ValidationFieldError{
				Field:   fmt.Sprintf("referencias[%d].tipo_documento", i),
				Message: "La referencia debe ser a una nota de crédito o débito",
			})
		}
	}

	return errors
}

// ValidarImpuestosAdicionalesItems valida los impuestos adicionales de los items de domain.Item
func (s *TributarioValidation) ValidarImpuestosAdicionalesItems(items []domain.Item) []models.ValidationFieldError {
	var errors []models.ValidationFieldError

	for i, item := range items {
		// No se validan impuestos adicionales en items del dominio
		// ya que este tipo no tiene el campo ImpuestosAdicionales en la estructura original
	}

	return errors
}

// ValidarImpuestosAdicionalesBoleta valida los impuestos adicionales de los items de una boleta
func (s *TributarioValidation) ValidarImpuestosAdicionalesBoleta(items []*models.DetalleBoleta) []models.ValidationFieldError {
	var errors []models.ValidationFieldError

	for i, item := range items {
		// No se validan impuestos adicionales en los detalles de boleta
		// ya que este tipo no tiene el campo ImpuestosAdicionales en la estructura original
	}

	return errors
}

// ValidarTipoReferencia valida un tipo de referencia
func (s *TributarioValidation) ValidarTipoReferencia(tipoReferencia models.TipoReferencia) error {
	switch tipoReferencia {
	case models.TipoAnula, models.TipoCorrige, models.TipoPreciosCantidad,
		models.TipoReferenciaInterna, models.TipoGuiaDespacho, models.TipoOtraReferencia,
		models.TipoSetPruebas, models.TipoOrdenCompra, models.TipoNotaCredito, models.TipoNotaDebito:
		return nil
	default:
		return fmt.Errorf("tipo de referencia no válido: %s", tipoReferencia)
	}
}
