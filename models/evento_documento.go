package models

type EventoDocumento struct {
	ID      string `json:"id"`
	Evento  string `json:"evento"`
	Detalle string `json:"detalle"`
}
