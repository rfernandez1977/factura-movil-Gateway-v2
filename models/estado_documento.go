package models

import (
	"time"
)

// EstadoDocumentoResponse representa la respuesta del estado de un documento
// Esta estructura es utilizada por el servicio de boletas para devolver el estado consultado
type EstadoDocumentoResponse struct {
	TrackID        string    `json:"track_id"`
	Estado         string    `json:"estado"`
	Glosa          string    `json:"glosa"`
	FechaRecepcion time.Time `json:"fecha_recepcion"`
	Errores        []string  `json:"errores,omitempty"`
}
