package models

import "time"

// TipoReferencia representa el tipo de referencia entre documentos
type TipoReferencia string

// Tipos de referencias
const (
	TipoAnula             TipoReferencia = "1" // Anula documento de referencia
	TipoCorrige           TipoReferencia = "2" // Corrige texto documento de referencia
	TipoPreciosCantidad   TipoReferencia = "3" // Corrige montos y/o cantidades
	TipoReferenciaInterna TipoReferencia = "4" // Referencia interna
	TipoOtraReferencia    TipoReferencia = "6" // Otros
	TipoSetPruebas        TipoReferencia = "7" // Set de pruebas
	TipoOrdenCompra       TipoReferencia = "8" // Orden de compra
)

// Referencia representa una referencia a otro documento
type Referencia struct {
	ID              string         `json:"id" bson:"_id,omitempty"`
	TipoDocumento   string         `json:"tipo_documento" bson:"tipo_documento"`
	TipoReferencia  TipoReferencia `json:"tipo_referencia" bson:"tipo_referencia"`
	Folio           int            `json:"folio" bson:"folio"`
	FechaReferencia time.Time      `json:"fecha_referencia" bson:"fecha_referencia"`
	RazonReferencia string         `json:"razon_referencia" bson:"razon_referencia"`
	DocumentoID     string         `json:"documento_id" bson:"documento_id"`
}

// ReferenciaDocumento representa una referencia entre documentos tributarios
type ReferenciaDocumento struct {
	ID                  string                `json:"id" bson:"_id"`
	TipoDocumentoOrigen string                `json:"tipoDocumentoOrigen" bson:"tipo_documento_origen"`
	FolioOrigen         int                   `json:"folioOrigen" bson:"folio_origen"`
	TipoDocumentoRef    string                `json:"tipoDocumentoRef" bson:"tipo_documento_ref"`
	FolioRef            int                   `json:"folioRef" bson:"folio_ref"`
	FechaRef            time.Time             `json:"fechaRef" bson:"fecha_ref"`
	CodigoRef           string                `json:"codigoRef" bson:"codigo_ref"`
	RazonRef            string                `json:"razonRef" bson:"razon_ref"`
	ValidacionesSII     []SIIValidationResult `json:"validacionesSII" bson:"validaciones_sii"`
}

// SIIValidationResult representa el resultado de una validación del SII
type SIIValidationResult struct {
	CodigoValidacion string    `json:"codigoValidacion" bson:"codigo_validacion"`
	Resultado        bool      `json:"resultado" bson:"resultado"`
	MensajeError     string    `json:"mensajeError,omitempty" bson:"mensaje_error,omitempty"`
	FechaValidacion  time.Time `json:"fechaValidacion" bson:"fecha_validacion"`
}

// GetReferenciaTipoGuiaDespacho obtiene el tipo de referencia para guías de despacho
func GetReferenciaTipoGuiaDespacho() TipoReferencia {
	return TipoReferencia("5") // Evitamos usar la constante directamente
}

// GetReferenciaTipoNotaCredito obtiene el tipo de referencia para notas de crédito
func GetReferenciaTipoNotaCredito() TipoReferencia {
	return TipoReferencia("9") // Evitamos usar la constante directamente
}

// GetReferenciaTipoNotaDebito obtiene el tipo de referencia para notas de débito
func GetReferenciaTipoNotaDebito() TipoReferencia {
	return TipoReferencia("10") // Evitamos usar la constante directamente
}
