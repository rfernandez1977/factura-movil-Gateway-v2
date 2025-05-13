package models

import (
	"time"
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

// Documento representa un documento tributario
type Documento struct {
	ID              string    `json:"id" db:"id"`
	EmpresaID       string    `json:"empresa_id" db:"empresa_id"`
	TipoDocumento   string    `json:"tipo_documento" db:"tipo_documento"`
	NumeroDocumento string    `json:"numero_documento" db:"numero_documento"`
	FechaEmision    string    `json:"fecha_emision" db:"fecha_emision"`
	Monto           float64   `json:"monto" db:"monto"`
	Estado          string    `json:"estado" db:"estado"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
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
		return ValidationError{Field: "empresa_id", Message: "El ID de la empresa es obligatorio"}
	}
	if d.TipoDocumento == "" {
		return ValidationError{Field: "tipo_documento", Message: "El tipo de documento es obligatorio"}
	}
	if d.NumeroDocumento == "" {
		return ValidationError{Field: "numero_documento", Message: "El número de documento es obligatorio"}
	}
	if d.FechaEmision == "" {
		return ValidationError{Field: "fecha_emision", Message: "La fecha de emisión es obligatoria"}
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
