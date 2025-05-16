package utils

import (
	"fmt"
	"time"

	"github.com/fmgo/domain"
	"github.com/fmgo/models"
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

	if v.doc.RUTEmisor == "" {
		return fmt.Errorf("RUT emisor es requerido")
	}

	if err := ValidateRUT(v.doc.RUTEmisor); err != nil {
		return fmt.Errorf("RUT emisor inválido: %v", err)
	}

	if v.doc.RUTReceptor == "" {
		return fmt.Errorf("RUT receptor es requerido")
	}

	if err := ValidateRUT(v.doc.RUTReceptor); err != nil {
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
	// En el caso de DocumentoTributario, los items están en el campo Detalles
	if len(v.doc.Detalles) == 0 {
		return nil
	}

	var montoNeto, montoExento, montoIVA, montoImpuestosAdicionales float64

	for _, detalle := range v.doc.Detalles {
		subtotal := v.amountValidator.CalculateSubtotal(float64(detalle.Cantidad), detalle.PrecioUnitario, 0)

		if detalle.Exento {
			montoExento += subtotal
		} else {
			montoNeto += subtotal
			montoIVA += v.amountValidator.CalculateIVA(subtotal, 19.0) // 19% IVA estándar
		}
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
			validator.RUTEmisorDestino = v.doc.RUTEmisor
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
	tipoDTEStr := fmt.Sprintf("%d", factura.TipoDocumento)

	modelDoc := &models.DocumentoTributario{
		ID:           factura.ID,
		TipoDTE:      tipoDTEStr,
		Folio:        int(factura.Folio),
		FechaEmision: factura.FechaEmision,
		MontoTotal:   factura.MontoTotal,
		MontoNeto:    factura.MontoNeto,
		MontoExento:  factura.MontoExento,
		MontoIVA:     factura.MontoIVA,
		RUTEmisor:    factura.RutEmisor,
		RUTReceptor:  factura.RutReceptor,
		Estado:       models.EstadoDTE(factura.Estado),
		// No hay campo Items sino Detalles en DocumentoTributario
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
	boleta *models.BoletaElectronica
}

// NewBoletaValidator crea una nueva instancia de BoletaValidator
func NewBoletaValidator(boleta *models.BoletaElectronica) *BoletaValidator {
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
		RUTEmisor:    boleta.RUTEmisor,
		RUTReceptor:  boleta.RUTReceptor,
		Estado:       boleta.Estado,
		// No hay campo Items sino Detalles en DocumentoTributario
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
		RUTEmisor:    notaCredito.RUTEmisor,
		RUTReceptor:  notaCredito.RUTReceptor,
		Estado:       notaCredito.Estado,
		// No hay campo Items sino Detalles en DocumentoTributario
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
		RUTEmisor:    notaDebito.RUTEmisor,
		RUTReceptor:  notaDebito.RUTReceptor,
		Estado:       notaDebito.Estado,
		// No hay campo Items sino Detalles en DocumentoTributario
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
		RUTEmisor:    guiaDespacho.RUTEmisor,
		RUTReceptor:  guiaDespacho.RUTReceptor,
		Estado:       guiaDespacho.Estado,
		// No hay campo Items sino Detalles en DocumentoTributario
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

	// El ítem actual no tiene estos campos, así que omitimos estas validaciones
	// Subtotal y MontoIVA no aparecen en el modelo Item

	return nil
}

// Funciones auxiliares para la conversión entre domain y models

// convertItemsToDetalles convierte items a detalles
func convertItemsToDetalles(items []models.Item) []models.DetalleTributario {
	detalles := make([]models.DetalleTributario, len(items))
	for i, item := range items {
		detalles[i] = models.DetalleTributario{
			Descripcion:    item.Descripcion,
			Cantidad:       int(item.Cantidad),
			PrecioUnitario: item.PrecioUnitario,
			MontoItem:      item.MontoItem,
			Exento:         item.Exento,
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
			MontoItem:      domainItem.MontoTotal,
			Exento:         false,
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
