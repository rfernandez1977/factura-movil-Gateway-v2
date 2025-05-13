package utils

import (
	"fmt"
	"time"

	"github.com/cursor/FMgo/domain"
	"github.com/cursor/FMgo/models"
)

// DocumentValidator define la interfaz para validar documentos tributarios
type DocumentValidator interface {
	Validate() error
	CalculateTotals() error
	ValidateBusinessRules() error
}

// BaseDocumentValidator implementa la validación base para todos los documentos
type BaseDocumentValidator struct {
	doc             *models.DocumentoTributario
	amountValidator *AmountValidator
	dateValidator   *DateValidator
}

// NewBaseDocumentValidator crea una nueva instancia de BaseDocumentValidator
func NewBaseDocumentValidator(doc *models.DocumentoTributario) *BaseDocumentValidator {
	return &BaseDocumentValidator{
		doc:             doc,
		amountValidator: NewAmountValidator(),
		dateValidator:   NewDateValidator(),
	}
}

// Validate implementa la validación base
func (v *BaseDocumentValidator) Validate() error {
	if v.doc.Folio <= 0 {
		return fmt.Errorf("folio debe ser mayor que 0")
	}

	if v.doc.RutEmisor == "" {
		return fmt.Errorf("RUT emisor es requerido")
	}

	if err := ValidateRUT(v.doc.RutEmisor); err != nil {
		return fmt.Errorf("RUT emisor inválido: %v", err)
	}

	if v.doc.RutReceptor == "" {
		return fmt.Errorf("RUT receptor es requerido")
	}

	if err := ValidateRUT(v.doc.RutReceptor); err != nil {
		return fmt.Errorf("RUT receptor inválido: %v", err)
	}

	if err := v.amountValidator.ValidateTotalAmount(v.doc.MontoTotal); err != nil {
		return err
	}

	if v.doc.FechaEmision.IsZero() {
		return fmt.Errorf("fecha de emisión es requerida")
	}

	if v.doc.FechaEmision.After(time.Now()) {
		return fmt.Errorf("fecha de emisión no puede ser futura")
	}

	return nil
}

// CalculateTotals calcula los totales del documento
func (v *BaseDocumentValidator) CalculateTotals() error {
	var montoNeto, montoExento, montoIVA, montoImpuestosAdicionales float64

	for i := range v.doc.Items {
		item := &v.doc.Items[i]
		// Validar ítem
		itemValidator := NewDocumentItemValidator(item)
		if err := itemValidator.Validate(); err != nil {
			return err
		}

		// Calcular subtotal
		subtotal := v.amountValidator.CalculateSubtotal(item.Cantidad, item.PrecioUnitario, item.Descuento)
		item.Subtotal = v.amountValidator.RoundAmount(subtotal)

		if item.Exento {
			montoExento += item.Subtotal
		} else {
			montoNeto += item.Subtotal

			// Recalcular el IVA basado en el subtotal y el porcentaje
			ivaItem := v.amountValidator.CalculateIVA(item.Subtotal, item.PorcentajeIVA)
			item.MontoIVA = ivaItem
			montoIVA += ivaItem
		}

		// Calcular impuestos adicionales
		var itemImpuestosAdicionales float64
		for j := range item.ImpuestosAdicionales {
			impuesto := &item.ImpuestosAdicionales[j]
			impuesto.BaseImponible = item.Subtotal
			impuesto.Monto = v.amountValidator.RoundAmount(item.Subtotal * (impuesto.Porcentaje / 100))
			itemImpuestosAdicionales += impuesto.Monto
		}
		montoImpuestosAdicionales += itemImpuestosAdicionales
	}

	// Actualizar totales
	v.doc.MontoNeto = v.amountValidator.RoundAmount(montoNeto)
	v.doc.MontoExento = v.amountValidator.RoundAmount(montoExento)
	v.doc.MontoIVA = v.amountValidator.RoundAmount(montoIVA)
	v.doc.MontoTotal = v.amountValidator.RoundAmount(montoNeto + montoExento + montoIVA + montoImpuestosAdicionales)

	return nil
}

// ValidateBusinessRules implementa las reglas de negocio base
func (v *BaseDocumentValidator) ValidateBusinessRules() error {
	if len(v.doc.Referencias) > 0 {
		// Crear validadores de referencia
		var validators []ReferenceValidator
		for _, ref := range v.doc.Referencias {
			validator := NewReferenceValidator(ref)
			validator.TipoDocumentoDestino = getDTEDocumentType(v.doc.TipoDTE)
			validator.FolioDestino = fmt.Sprintf("%d", v.doc.Folio)
			validator.FechaDestino = v.doc.FechaEmision
			validator.RUTEmisorDestino = v.doc.RutEmisor
			validators = append(validators, *validator)
		}

		// Validar cadena de referencias
		if err := ValidateReferenceChain(validators); err != nil {
			return fmt.Errorf("error en la cadena de referencias: %v", err)
		}
	}

	return nil
}

// FacturaValidator implementa la validación específica para facturas
type FacturaValidator struct {
	*BaseDocumentValidator
	factura *models.Factura
}

// NewFacturaValidator crea una nueva instancia de FacturaValidator
func NewFacturaValidator(factura *models.Factura) *FacturaValidator {
	// Convertir de domain.DocumentoTributario a models.DocumentoTributario
	modelDoc := &models.DocumentoTributario{
		ID:           factura.ID,
		TipoDTE:      factura.TipoDocumento.ToString(),
		Folio:        int(factura.Folio),
		FechaEmision: factura.FechaEmision,
		MontoTotal:   factura.MontoTotal,
		MontoNeto:    factura.MontoNeto,
		MontoExento:  factura.MontoExento,
		MontoIVA:     factura.MontoIVA,
		RutEmisor:    factura.RutEmisor,
		RutReceptor:  factura.RutReceptor,
		Estado:       models.EstadoDTE(factura.Estado),
		// Convertir items si es necesario
		// Items: convertItems(factura.Items),
	}

	return &FacturaValidator{
		BaseDocumentValidator: NewBaseDocumentValidator(modelDoc),
		factura:               factura,
	}
}

// Validate implementa la validación específica para facturas
func (v *FacturaValidator) Validate() error {
	if err := v.BaseDocumentValidator.Validate(); err != nil {
		return err
	}

	if v.factura.FormaPago == "" {
		return fmt.Errorf("forma de pago es requerida")
	}

	return nil
}

// ValidateBusinessRules implementa las reglas de negocio específicas para facturas
func (v *FacturaValidator) ValidateBusinessRules() error {
	if err := v.BaseDocumentValidator.ValidateBusinessRules(); err != nil {
		return err
	}

	// Verificar si es una factura exenta o afecta basada en el monto exento
	if v.factura.MontoExento > 0 && v.factura.MontoIVA > 0 {
		return fmt.Errorf("una factura no puede tener montos exentos y afectos simultáneamente")
	}

	// Si es una factura exenta (determinado por el monto exento)
	if v.factura.MontoExento > 0 && v.factura.MontoIVA <= 0 {
		// Verificar que no haya IVA
		if v.factura.MontoIVA > 0 {
			return fmt.Errorf("una factura exenta no puede tener monto de IVA")
		}
	}

	// Si es una factura afecta (normal)
	if v.factura.MontoExento <= 0 && v.factura.MontoNeto > 0 {
		// Verificar que tenga IVA
		if v.factura.MontoIVA <= 0 {
			return fmt.Errorf("una factura afecta debe tener monto de IVA")
		}
	}

	return nil
}

// BoletaValidator implementa la validación específica para boletas
type BoletaValidator struct {
	*BaseDocumentValidator
	boleta *models.Boleta
}

// NewBoletaValidator crea una nueva instancia de BoletaValidator
func NewBoletaValidator(boleta *models.Boleta) *BoletaValidator {
	// Convertir de DocumentoTributario a models.DocumentoTributario
	modelDoc := &models.DocumentoTributario{
		ID:           boleta.ID,
		TipoDTE:      boleta.TipoDTE,
		Folio:        boleta.Folio,
		FechaEmision: boleta.FechaEmision,
		MontoTotal:   boleta.MontoTotal,
		MontoNeto:    boleta.MontoNeto,
		MontoExento:  boleta.MontoExento,
		MontoIVA:     boleta.MontoIVA,
		RutEmisor:    boleta.RutEmisor,
		RutReceptor:  boleta.RutReceptor,
		Estado:       boleta.Estado,
		Items:        boleta.Items,
		Detalles:     convertItemsToDetalles(boleta.Items),
	}

	return &BoletaValidator{
		BaseDocumentValidator: NewBaseDocumentValidator(modelDoc),
		boleta:                boleta,
	}
}

// Validate implementa la validación específica para boletas
func (v *BoletaValidator) Validate() error {
	if err := v.BaseDocumentValidator.Validate(); err != nil {
		return err
	}

	if v.boleta.Vendedor == "" {
		return fmt.Errorf("vendedor es requerido")
	}

	return nil
}

// ValidateBusinessRules implementa las reglas de negocio específicas para boletas
func (v *BoletaValidator) ValidateBusinessRules() error {
	if v.boleta.MontoExento > 0 {
		return fmt.Errorf("las boletas no pueden tener monto exento")
	}

	return nil
}

// NotaCreditoValidator implementa la validación específica para notas de crédito
type NotaCreditoValidator struct {
	*BaseDocumentValidator
	notaCredito *models.NotaCredito
}

// NewNotaCreditoValidator crea una nueva instancia de NotaCreditoValidator
func NewNotaCreditoValidator(notaCredito *models.NotaCredito) *NotaCreditoValidator {
	// Convertir de DocumentoTributario a models.DocumentoTributario
	modelDoc := &models.DocumentoTributario{
		ID:           notaCredito.ID,
		TipoDTE:      notaCredito.TipoDTE,
		Folio:        notaCredito.Folio,
		FechaEmision: notaCredito.FechaEmision,
		MontoTotal:   notaCredito.MontoTotal,
		MontoNeto:    notaCredito.MontoNeto,
		MontoExento:  notaCredito.MontoExento,
		MontoIVA:     notaCredito.MontoIVA,
		RutEmisor:    notaCredito.RutEmisor,
		RutReceptor:  notaCredito.RutReceptor,
		Estado:       notaCredito.Estado,
		Items:        notaCredito.Items,
		Detalles:     convertItemsToDetalles(notaCredito.Items),
	}

	return &NotaCreditoValidator{
		BaseDocumentValidator: NewBaseDocumentValidator(modelDoc),
		notaCredito:           notaCredito,
	}
}

// Validate implementa la validación específica para notas de crédito
func (v *NotaCreditoValidator) Validate() error {
	if err := v.BaseDocumentValidator.Validate(); err != nil {
		return err
	}

	if v.notaCredito.FolioReferencia <= 0 {
		return fmt.Errorf("folio de referencia es requerido")
	}

	if v.notaCredito.TipoReferencia == "" {
		return fmt.Errorf("tipo de referencia es requerido")
	}

	if v.notaCredito.Motivo == "" {
		return fmt.Errorf("motivo es requerido")
	}

	return nil
}

// ValidateBusinessRules implementa las reglas de negocio específicas para notas de crédito
func (v *NotaCreditoValidator) ValidateBusinessRules() error {
	if v.notaCredito.DocumentoReferencia == nil {
		return fmt.Errorf("documento de referencia es requerido")
	}

	if v.notaCredito.MontoTotal > v.notaCredito.DocumentoReferencia.MontoTotal {
		return fmt.Errorf("el monto total de la nota de crédito no puede ser mayor que el documento original")
	}

	return nil
}

// NotaDebitoValidator implementa la validación específica para notas de débito
type NotaDebitoValidator struct {
	*BaseDocumentValidator
	notaDebito *models.NotaDebito
}

// NewNotaDebitoValidator crea una nueva instancia de NotaDebitoValidator
func NewNotaDebitoValidator(notaDebito *models.NotaDebito) *NotaDebitoValidator {
	// Convertir de DocumentoTributario a models.DocumentoTributario
	modelDoc := &models.DocumentoTributario{
		ID:           notaDebito.ID,
		TipoDTE:      notaDebito.TipoDTE,
		Folio:        notaDebito.Folio,
		FechaEmision: notaDebito.FechaEmision,
		MontoTotal:   notaDebito.MontoTotal,
		MontoNeto:    notaDebito.MontoNeto,
		MontoExento:  notaDebito.MontoExento,
		MontoIVA:     notaDebito.MontoIVA,
		RutEmisor:    notaDebito.RutEmisor,
		RutReceptor:  notaDebito.RutReceptor,
		Estado:       notaDebito.Estado,
		Items:        notaDebito.Items,
		Detalles:     convertItemsToDetalles(notaDebito.Items),
	}

	return &NotaDebitoValidator{
		BaseDocumentValidator: NewBaseDocumentValidator(modelDoc),
		notaDebito:            notaDebito,
	}
}

// Validate implementa la validación específica para notas de débito
func (v *NotaDebitoValidator) Validate() error {
	if err := v.BaseDocumentValidator.Validate(); err != nil {
		return err
	}

	if v.notaDebito.FolioReferencia <= 0 {
		return fmt.Errorf("folio de referencia es requerido")
	}

	if v.notaDebito.TipoReferencia == "" {
		return fmt.Errorf("tipo de referencia es requerido")
	}

	if v.notaDebito.Motivo == "" {
		return fmt.Errorf("motivo es requerido")
	}

	return nil
}

// ValidateBusinessRules implementa las reglas de negocio específicas para notas de débito
func (v *NotaDebitoValidator) ValidateBusinessRules() error {
	// Las notas de débito no tienen restricción de monto total
	return nil
}

// GuiaDespachoValidator implementa la validación específica para guías de despacho
type GuiaDespachoValidator struct {
	*BaseDocumentValidator
	guiaDespacho *models.GuiaDespacho
}

// NewGuiaDespachoValidator crea una nueva instancia de GuiaDespachoValidator
func NewGuiaDespachoValidator(guiaDespacho *models.GuiaDespacho) *GuiaDespachoValidator {
	// Convertir de DocumentoTributario a models.DocumentoTributario
	modelDoc := &models.DocumentoTributario{
		ID:           guiaDespacho.ID,
		TipoDTE:      guiaDespacho.TipoDTE,
		Folio:        guiaDespacho.Folio,
		FechaEmision: guiaDespacho.FechaEmision,
		MontoTotal:   guiaDespacho.MontoTotal,
		MontoNeto:    guiaDespacho.MontoNeto,
		MontoExento:  guiaDespacho.MontoExento,
		MontoIVA:     guiaDespacho.MontoIVA,
		RutEmisor:    guiaDespacho.RutEmisor,
		RutReceptor:  guiaDespacho.RutReceptor,
		Estado:       guiaDespacho.Estado,
		Items:        guiaDespacho.Items,
		Detalles:     convertItemsToDetalles(guiaDespacho.Items),
	}

	return &GuiaDespachoValidator{
		BaseDocumentValidator: NewBaseDocumentValidator(modelDoc),
		guiaDespacho:          guiaDespacho,
	}
}

// Validate implementa la validación específica para guías de despacho
func (v *GuiaDespachoValidator) Validate() error {
	if err := v.BaseDocumentValidator.Validate(); err != nil {
		return err
	}

	if v.guiaDespacho.DireccionDestino == "" {
		return fmt.Errorf("dirección de destino es requerida")
	}

	if v.guiaDespacho.Transportista == "" {
		return fmt.Errorf("transportista es requerido")
	}

	if v.guiaDespacho.Patente == "" {
		return fmt.Errorf("patente es requerida")
	}

	return nil
}

// ValidateBusinessRules implementa las reglas de negocio específicas para guías de despacho
func (v *GuiaDespachoValidator) ValidateBusinessRules() error {
	return nil
}

// DocumentItemValidator implementa la validación de ítems de documento
type DocumentItemValidator struct {
	Item            *models.Item
	amountValidator *AmountValidator
}

// NewDocumentItemValidator crea una nueva instancia de DocumentItemValidator
func NewDocumentItemValidator(item *models.Item) *DocumentItemValidator {
	return &DocumentItemValidator{
		Item:            item,
		amountValidator: NewAmountValidator(),
	}
}

// Validate implementa la validación de ítems
func (v *DocumentItemValidator) Validate() error {
	if v.Item == nil {
		return fmt.Errorf("el ítem no puede ser nulo")
	}

	if v.Item.Codigo == "" {
		return fmt.Errorf("código del ítem es requerido")
	}

	if v.Item.Descripcion == "" {
		return fmt.Errorf("descripción del ítem es requerida")
	}

	if err := v.amountValidator.ValidateQuantity(v.Item.Cantidad); err != nil {
		return fmt.Errorf("cantidad inválida: %v", err)
	}

	if err := v.amountValidator.ValidateUnitPrice(v.Item.PrecioUnitario); err != nil {
		return fmt.Errorf("precio unitario inválido: %v", err)
	}

	// Validar que el descuento sea un porcentaje válido
	if err := v.amountValidator.ValidatePercentage(v.Item.Descuento, "descuento"); err != nil {
		return err
	}

	if err := v.amountValidator.ValidatePercentage(v.Item.PorcentajeIVA, "porcentaje de IVA"); err != nil {
		return err
	}

	// Calcular y validar subtotal
	subtotal := v.amountValidator.CalculateSubtotal(v.Item.Cantidad, v.Item.PrecioUnitario, v.Item.Descuento)
	if v.amountValidator.RoundAmount(subtotal) != v.amountValidator.RoundAmount(v.Item.Subtotal) {
		return fmt.Errorf("el subtotal calculado (%.2f) no coincide con el subtotal del ítem (%.2f)",
			v.amountValidator.RoundAmount(subtotal), v.amountValidator.RoundAmount(v.Item.Subtotal))
	}

	// Calcular y validar IVA
	iva := v.amountValidator.CalculateIVA(subtotal, v.Item.PorcentajeIVA)
	if v.amountValidator.RoundAmount(iva) != v.amountValidator.RoundAmount(v.Item.MontoIVA) {
		return fmt.Errorf("el monto de IVA calculado (%.2f) no coincide con el monto de IVA del ítem (%.2f)",
			v.amountValidator.RoundAmount(iva), v.amountValidator.RoundAmount(v.Item.MontoIVA))
	}

	return nil
}

// Funciones auxiliares para la conversión entre domain y models

// convertItemsToDetalles convierte items a detalles
func convertItemsToDetalles(items []models.Item) []models.Detalle {
	detalles := make([]models.Detalle, len(items))
	for i, item := range items {
		detalles[i] = models.Detalle{
			Descripcion:    item.Descripcion,
			Cantidad:       int(item.Cantidad),
			PrecioUnitario: item.PrecioUnitario,
			MontoItem:      item.MontoItem,
		}
	}
	return detalles
}

// convertDomainItemsToModels convierte items de domain a models
func convertDomainItemsToModels(domainItems []domain.Item) []models.Item {
	modelItems := make([]models.Item, len(domainItems))
	for i, domainItem := range domainItems {
		modelItems[i] = models.Item{
			Codigo:         domainItem.ID.String(),
			Descripcion:    domainItem.Descripcion,
			Cantidad:       domainItem.Cantidad,
			PrecioUnitario: domainItem.PrecioUnit,
			Subtotal:       domainItem.MontoNeto,
			MontoIVA:       domainItem.MontoIVA,
			MontoItem:      domainItem.MontoTotal,
			Exento:         false,
			PorcentajeIVA:  19.0, // Valor predeterminado para Chile
			Descuento:      0,
		}
	}
	return modelItems
}

// getDTEDocumentType obtiene el tipo de documento
func getDTEDocumentType(tipoDTE string) string {
	switch tipoDTE {
	case "33":
		return "FACTURA_ELECTRONICA"
	case "39":
		return "BOLETA_ELECTRONICA"
	case "56":
		return "NOTA_DEBITO_ELECTRONICA"
	case "61":
		return "NOTA_CREDITO_ELECTRONICA"
	case "52":
		return "GUIA_DESPACHO_ELECTRONICA"
	default:
		return "DOCUMENTO_DESCONOCIDO"
	}
}
