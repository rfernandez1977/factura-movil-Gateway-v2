package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DocumentoTributarioBasico representa un documento tributario electrónico básico
type DocumentoTributarioBasico struct {
	ID           string    `json:"id" bson:"_id"`
	TipoDTE      string    `json:"tipo_dte" bson:"tipo_dte"`
	Folio        int64     `json:"folio" bson:"folio"`
	FechaEmision time.Time `json:"fecha_emision" bson:"fecha_emision"`
	RutEmisor    string    `json:"rut_emisor" bson:"rut_emisor"`
	RutReceptor  string    `json:"rut_receptor" bson:"rut_receptor"`
	MontoTotal   float64   `json:"monto_total" bson:"monto_total"`
	MontoNeto    float64   `json:"monto_neto" bson:"monto_neto"`
	MontoIVA     float64   `json:"monto_iva" bson:"monto_iva"`
	Estado       string    `json:"estado" bson:"estado"`
	TrackID      string    `json:"track_id,omitempty" bson:"track_id,omitempty"`
	XML          string    `json:"xml,omitempty" bson:"xml,omitempty"`
	XMLFirmado   string    `json:"xml_firmado,omitempty" bson:"xml_firmado,omitempty"`
}

// DocumentoRequest representa una solicitud para crear un documento
type DocumentoRequest struct {
	TipoDTE     string  `json:"tipo_dte" binding:"required"`
	RutEmisor   string  `json:"rut_emisor" binding:"required"`
	RutReceptor string  `json:"rut_receptor" binding:"required"`
	MontoNeto   float64 `json:"monto_neto" binding:"required,gte=0"`
	MontoIVA    float64 `json:"monto_iva" binding:"required,gte=0"`
	MontoTotal  float64 `json:"monto_total" binding:"required,gt=0"`
	Items       []Item  `json:"items" binding:"required,min=1"`
}

// DocumentoResponse representa la respuesta de un documento
type DocumentoResponse struct {
	ID           string    `json:"id"`
	TipoDTE      string    `json:"tipo_dte"`
	Folio        int64     `json:"folio"`
	FechaEmision time.Time `json:"fecha_emision"`
	RutEmisor    string    `json:"rut_emisor"`
	RutReceptor  string    `json:"rut_receptor"`
	MontoTotal   float64   `json:"monto_total"`
	MontoNeto    float64   `json:"monto_neto"`
	MontoIVA     float64   `json:"monto_iva"`
	Estado       string    `json:"estado"`
	TrackID      string    `json:"track_id,omitempty"`
}

// Documento representa un documento almacenado
type Documento struct {
	ID            primitive.ObjectID     `json:"id" bson:"_id,omitempty"`
	TipoDocumento string                 `json:"tipo_documento" bson:"tipo_documento"`
	Folio         int                    `json:"folio" bson:"folio"`
	RutEmisor     string                 `json:"rut_emisor" bson:"rut_emisor"`
	RutReceptor   string                 `json:"rut_receptor" bson:"rut_receptor"`
	FechaEmision  time.Time              `json:"fecha_emision" bson:"fecha_emision"`
	MontoTotal    float64                `json:"monto_total" bson:"monto_total"`
	Estado        string                 `json:"estado" bson:"estado"`
	XML           string                 `json:"xml,omitempty" bson:"xml,omitempty"`
	PDF           string                 `json:"pdf,omitempty" bson:"pdf,omitempty"`
	Timbre        string                 `json:"timbre,omitempty" bson:"timbre,omitempty"`
	TrackID       string                 `json:"track_id,omitempty" bson:"track_id,omitempty"`
	CreatedAt     time.Time              `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at" bson:"updated_at"`
	Metadata      map[string]interface{} `json:"metadata,omitempty" bson:"metadata,omitempty"`
	CacheInfo     *DocCacheInfo          `json:"cache_info,omitempty" bson:"cache_info,omitempty"`
	Contenido     []byte                 `json:"contenido,omitempty" bson:"contenido,omitempty"`
}

// NewDocumento crea una nueva instancia de Documento
func NewDocumento(empresaID, tipoDocumento, numeroDocumento, fechaEmision string, monto float64) *Documento {
	return &Documento{
		EmpresaID:       empresaID,
		TipoDocumento:   tipoDocumento,
		NumeroDocumento: numeroDocumento,
		FechaEmision:    fechaEmision,
		Monto:           monto,
		Estado:          "PENDIENTE",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
}

// Validate valida que todos los campos obligatorios estén presentes
func (d *Documento) Validate() error {
	if d.EmpresaID == "" {
		return &ValidationFieldError{Field: "empresa_id", Message: "El ID de la empresa es obligatorio"}
	}
	if d.TipoDocumento == "" {
		return &ValidationFieldError{Field: "tipo_documento", Message: "El tipo de documento es obligatorio"}
	}
	if d.NumeroDocumento == "" {
		return &ValidationFieldError{Field: "numero_documento", Message: "El número de documento es obligatorio"}
	}
	if d.FechaEmision == "" {
		return &ValidationFieldError{Field: "fecha_emision", Message: "La fecha de emisión es obligatoria"}
	}
	return nil
}

// EstadosDocumento define los posibles estados de un documento
var EstadosDocumento = struct {
	Pendiente  string
	Enviado    string
	Aceptado   string
	Rechazado  string
	Anulado    string
	Procesando string
	Error      string
	Completado string
}{
	Pendiente:  "PENDIENTE",
	Enviado:    "ENVIADO",
	Aceptado:   "ACEPTADO",
	Rechazado:  "RECHAZADO",
	Anulado:    "ANULADO",
	Procesando: "PROCESANDO",
	Error:      "ERROR",
	Completado: "COMPLETADO",
}

// EstadoDocumento representa el estado de un documento
type EstadoDocumento string

// Estados de documentos
const (
	EstadoDocumentoPendiente  EstadoDocumento = "PENDIENTE"
	EstadoDocumentoEnviado    EstadoDocumento = "ENVIADO"
	EstadoDocumentoAceptado   EstadoDocumento = "ACEPTADO"
	EstadoDocumentoRechazado  EstadoDocumento = "RECHAZADO"
	EstadoDocumentoCancelado  EstadoDocumento = "CANCELADO"
	EstadoDocumentoCompletado EstadoDocumento = "COMPLETADO"
)

// DocCacheInfo contiene información de caché para un documento
type DocCacheInfo struct {
	LastAccessed time.Time `json:"last_accessed" bson:"last_accessed"`
	AccessCount  int       `json:"access_count" bson:"access_count"`
	CacheExpiry  time.Time `json:"cache_expiry" bson:"cache_expiry"`
}
