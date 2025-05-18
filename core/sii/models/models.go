package models

// RespuestaEnvio representa la respuesta del SII al enviar un documento
type RespuestaEnvio struct {
	TrackID string
	Estado  string
	Glosa   string
}

// EstadoEnvio representa el estado de un envío al SII
type EstadoEnvio struct {
	TrackID string
	Estado  string
	Glosa   string
	Fecha   string
}
