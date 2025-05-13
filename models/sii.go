package models

import (
	"net/http"
	"time"
)

// SIIResponseHTTP representa una respuesta del SII que incluye headers HTTP
type SIIResponseHTTP struct {
	Codigo           int                    `json:"codigo"`
	Mensaje          string                 `json:"mensaje"`
	Detalle          string                 `json:"detalle,omitempty"`
	Timestamp        time.Time              `json:"timestamp"`
	Header           http.Header            `json:"-"`
	DatosAdicionales map[string]interface{} `json:"datos_adicionales,omitempty"`
}
