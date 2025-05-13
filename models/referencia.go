package models

import "time"

// Referencia representa una referencia a otro documento tributario
type Referencia struct {
	ID               string    `json:"id" bson:"id"`
	NumeroLinea      int       `json:"numero_linea" bson:"numero_linea"`
	TipoDocumento    string    `json:"tipo_documento" bson:"tipo_documento"`
	Folio            string    `json:"folio" bson:"folio"`
	FechaReferencia  time.Time `json:"fecha_referencia" bson:"fecha_referencia"`
	CodigoReferencia string    `json:"codigo_referencia,omitempty" bson:"codigo_referencia,omitempty"`
	RazonReferencia  string    `json:"razon_referencia" bson:"razon_referencia"`
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
