package models

import "time"

// RespuestaSII representa la respuesta del SII al enviar un documento
type RespuestaSII struct {
	TrackID   string    `json:"track_id" xml:"trackid"`
	Estado    string    `json:"estado" xml:"estado"`
	Glosa     string    `json:"glosa" xml:"glosa"`
	Timestamp time.Time `json:"timestamp" xml:"timestamp"`
}

// EstadoSII representa el estado de un documento en el SII
type EstadoSII struct {
	TrackID   string    `json:"track_id" xml:"trackid"`
	Estado    string    `json:"estado" xml:"estado"`
	Glosa     string    `json:"glosa" xml:"glosa"`
	Timestamp time.Time `json:"timestamp" xml:"timestamp"`
}

// Estados comunes del SII
const (
	EstadoOK           = "OK"
	EstadoError        = "ERROR"
	EstadoProcesado    = "EPR" // Envío Procesado
	EstadoRecibido     = "REC" // Recibido
	EstadoRechazado    = "REP" // Rechazado
	EstadoDTERecibido  = "DTE_RECIBIDO"
	EstadoDTERechazado = "DTE_RECHAZADO"
)

// DetalleSII representa el detalle de una respuesta del SII
type DetalleSII struct {
	Tipo     string `json:"tipo" xml:"TIPO"`
	Folio    int64  `json:"folio" xml:"FOLIO"`
	Estado   string `json:"estado" xml:"ESTADO"`
	Glosa    string `json:"glosa" xml:"GLOSA"`
	NumError int    `json:"numError,omitempty" xml:"NUM_ERROR,omitempty"`
}
