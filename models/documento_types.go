package models

// TipoDTE representa el tipo de DTE
type TipoDTE int

// Tipos de DTE
const (
	TipoFactura                TipoDTE = 33  // Factura Electrónica
	TipoFacturaExenta          TipoDTE = 34  // Factura Electrónica Exenta
	TipoBoleta                 TipoDTE = 39  // Boleta Electrónica
	TipoBoletaExenta           TipoDTE = 41  // Boleta Electrónica Exenta
	TipoFacturaCompra          TipoDTE = 46  // Factura de Compra Electrónica
	TipoGuiaDespacho           TipoDTE = 52  // Guía de Despacho Electrónica
	TipoNotaDebito             TipoDTE = 56  // Nota de Débito Electrónica
	TipoNotaCredito            TipoDTE = 61  // Nota de Crédito Electrónica
	TipoFacturaExportacion     TipoDTE = 110 // Factura de Exportación Electrónica
	TipoNotaDebitoExportacion  TipoDTE = 111 // Nota de Débito de Exportación Electrónica
	TipoNotaCreditoExportacion TipoDTE = 112 // Nota de Crédito de Exportación Electrónica
)

// TipoDocumento representa los tipos de documento en el sistema
type TipoDocumento string

// Tipos de documento
const (
	DocumentoFactura                TipoDocumento = "FACTURA"
	DocumentoFacturaExenta          TipoDocumento = "FACTURA_EXENTA"
	DocumentoBoleta                 TipoDocumento = "BOLETA"
	DocumentoBoletaExenta           TipoDocumento = "BOLETA_EXENTA"
	DocumentoFacturaCompra          TipoDocumento = "FACTURA_COMPRA"
	DocumentoGuiaDespacho           TipoDocumento = "GUIA_DESPACHO"
	DocumentoNotaDebito             TipoDocumento = "NOTA_DEBITO"
	DocumentoNotaCredito            TipoDocumento = "NOTA_CREDITO"
	DocumentoFacturaExportacion     TipoDocumento = "FACTURA_EXPORTACION"
	DocumentoNotaDebitoExportacion  TipoDocumento = "NOTA_DEBITO_EXPORTACION"
	DocumentoNotaCreditoExportacion TipoDocumento = "NOTA_CREDITO_EXPORTACION"
	DocumentoOrdenCompra            TipoDocumento = "ORDEN_COMPRA"
	DocumentoLiquidacionFactura     TipoDocumento = "LIQUIDACION_FACTURA"
	DocumentoCobranza               TipoDocumento = "COBRANZA"
	DocumentoInformeComercial       TipoDocumento = "INFORME_COMERCIAL"
	DocumentoRecibo                 TipoDocumento = "RECIBO"
	DocumentoPresupuesto            TipoDocumento = "PRESUPUESTO"
	DocumentoProforma               TipoDocumento = "PROFORMA"
	DocumentoMemo                   TipoDocumento = "MEMO"
)

// EstadoSIIType representa el tipo de estado del SII
type EstadoSIIType string

// Estados SII
const (
	EstadoSIIAceptado  EstadoSIIType = "ACEPTADO"
	EstadoSIIRechazado EstadoSIIType = "RECHAZADO"
	EstadoSIIPendiente EstadoSIIType = "PENDIENTE"
	EstadoSIIErroneo   EstadoSIIType = "ERRONEO"
)

// ErrorValidacion representa un error de validación
type ErrorValidacion struct {
	Codigo    string `json:"codigo" bson:"codigo"`
	Mensaje   string `json:"mensaje" bson:"mensaje"`
	Campo     string `json:"campo,omitempty" bson:"campo,omitempty"`
	Valor     string `json:"valor,omitempty" bson:"valor,omitempty"`
	Timestamp string `json:"timestamp" bson:"timestamp"`
}

// Item representa un ítem de un documento
type Item struct {
	ID                     string                  `json:"id" bson:"_id,omitempty"`
	NumeroLinea            int                     `json:"numero_linea" bson:"numero_linea"`
	TipoItem               string                  `json:"tipo_item" bson:"tipo_item,omitempty"`
	Codigo                 string                  `json:"codigo" bson:"codigo,omitempty"`
	Nombre                 string                  `json:"nombre" bson:"nombre"`
	Descripcion            string                  `json:"descripcion,omitempty" bson:"descripcion,omitempty"`
	Cantidad               float64                 `json:"cantidad" bson:"cantidad"`
	UnidadMedida           string                  `json:"unidad_medida,omitempty" bson:"unidad_medida,omitempty"`
	PrecioUnitario         float64                 `json:"precio_unitario" bson:"precio_unitario"`
	MontoItem              float64                 `json:"monto_item" bson:"monto_item"`
	Descuento              float64                 `json:"descuento,omitempty" bson:"descuento,omitempty"`
	PorcentajeDescuento    float64                 `json:"porcentaje_descuento,omitempty" bson:"porcentaje_descuento,omitempty"`
	Recargo                float64                 `json:"recargo,omitempty" bson:"recargo,omitempty"`
	PorcentajeRecargo      float64                 `json:"porcentaje_recargo,omitempty" bson:"porcentaje_recargo,omitempty"`
	Exento                 bool                    `json:"exento" bson:"exento"`
	ImpuestosAdicionales   []ImpuestoAdicionalItem `json:"impuestos_adicionales,omitempty" bson:"impuestos_adicionales,omitempty"`
	MontoImpuestoAdicional float64                 `json:"monto_impuesto_adicional,omitempty" bson:"monto_impuesto_adicional,omitempty"`
	Metadata               map[string]interface{}  `json:"metadata,omitempty" bson:"metadata,omitempty"`
}

// ImpuestoAdicionalItem representa un impuesto adicional aplicado a un ítem
type ImpuestoAdicionalItem struct {
	Tipo        string  `json:"tipo" bson:"tipo"`
	Codigo      string  `json:"codigo" bson:"codigo"`
	Tasa        float64 `json:"tasa" bson:"tasa"`
	Monto       float64 `json:"monto" bson:"monto"`
	Descripcion string  `json:"descripcion,omitempty" bson:"descripcion,omitempty"`
}

// Convertir TipoDTE a TipoDocumento
func (t TipoDTE) ToTipoDocumento() TipoDocumento {
	switch t {
	case TipoFactura:
		return DocumentoFactura
	case TipoFacturaExenta:
		return DocumentoFacturaExenta
	case TipoBoleta:
		return DocumentoBoleta
	case TipoBoletaExenta:
		return DocumentoBoletaExenta
	case TipoFacturaCompra:
		return DocumentoFacturaCompra
	case TipoGuiaDespacho:
		return DocumentoGuiaDespacho
	case TipoNotaDebito:
		return DocumentoNotaDebito
	case TipoNotaCredito:
		return DocumentoNotaCredito
	case TipoFacturaExportacion:
		return DocumentoFacturaExportacion
	case TipoNotaDebitoExportacion:
		return DocumentoNotaDebitoExportacion
	case TipoNotaCreditoExportacion:
		return DocumentoNotaCreditoExportacion
	default:
		return ""
	}
}

// Convertir TipoDocumento a TipoDTE
func (t TipoDocumento) ToTipoDTE() TipoDTE {
	switch t {
	case DocumentoFactura:
		return TipoFactura
	case DocumentoFacturaExenta:
		return TipoFacturaExenta
	case DocumentoBoleta:
		return TipoBoleta
	case DocumentoBoletaExenta:
		return TipoBoletaExenta
	case DocumentoFacturaCompra:
		return TipoFacturaCompra
	case DocumentoGuiaDespacho:
		return TipoGuiaDespacho
	case DocumentoNotaDebito:
		return TipoNotaDebito
	case DocumentoNotaCredito:
		return TipoNotaCredito
	case DocumentoFacturaExportacion:
		return TipoFacturaExportacion
	case DocumentoNotaDebitoExportacion:
		return TipoNotaDebitoExportacion
	case DocumentoNotaCreditoExportacion:
		return TipoNotaCreditoExportacion
	default:
		return 0
	}
}
