package models

type Transformacion struct {
	ID          int64  `json:"id"`
	TipoOrigen  string `json:"tipo_origen"`
	TipoDestino string `json:"tipo_destino"`
	Reglas      string `json:"reglas"`
	Estado      string `json:"estado"`
}
