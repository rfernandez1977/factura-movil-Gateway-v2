package models

import (
	"time"
)

// CAFRequest representa una solicitud de nuevos folios al SII
type CAFRequest struct {
	ID             string    `json:"id" bson:"_id,omitempty"`
	TipoDocumento  string    `json:"tipo_documento" bson:"tipo_documento"`
	RutEmisor      string    `json:"rut_emisor" bson:"rut_emisor"`
	Cantidad       int       `json:"cantidad" bson:"cantidad"`
	FechaSolicitud time.Time `json:"fecha_solicitud" bson:"fecha_solicitud"`
	Estado         string    `json:"estado" bson:"estado"` // PENDIENTE, PROCESANDO, COMPLETADO, ERROR
	TrackID        string    `json:"track_id" bson:"track_id,omitempty"`
	UsuarioID      string    `json:"usuario_id" bson:"usuario_id"`
	EmpresaID      string    `json:"empresa_id" bson:"empresa_id"`
	CAFID          string    `json:"caf_id" bson:"caf_id,omitempty"`
	MensajeError   string    `json:"mensaje_error" bson:"mensaje_error,omitempty"`
	CreatedAt      time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" bson:"updated_at"`
}

// NewCAFRequest crea una nueva solicitud de CAF
func NewCAFRequest(tipoDocumento, rutEmisor string, cantidad int, empresaID, usuarioID string) *CAFRequest {
	now := time.Now()
	return &CAFRequest{
		TipoDocumento:  tipoDocumento,
		RutEmisor:      rutEmisor,
		Cantidad:       cantidad,
		FechaSolicitud: now,
		Estado:         "PENDIENTE",
		UsuarioID:      usuarioID,
		EmpresaID:      empresaID,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

// MarcarCompletado marca la solicitud como completada
func (r *CAFRequest) MarcarCompletado(trackID, cafID string) {
	r.Estado = "COMPLETADO"
	r.TrackID = trackID
	r.CAFID = cafID
	r.UpdatedAt = time.Now()
}

// MarcarError marca la solicitud como erronea
func (r *CAFRequest) MarcarError(mensaje string) {
	r.Estado = "ERROR"
	r.MensajeError = mensaje
	r.UpdatedAt = time.Now()
}
