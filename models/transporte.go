package models

// Transporte representa la información de transporte en guías de despacho
type Transporte struct {
	Patente          string `json:"patente" bson:"patente"`
	RUTTransportista string `json:"rut_transportista" bson:"rut_transportista"`
	Conductor        string `json:"conductor" bson:"conductor,omitempty"`
	DireccionOrigen  string `json:"direccion_origen" bson:"direccion_origen,omitempty"`
	DireccionDestino string `json:"direccion_destino" bson:"direccion_destino,omitempty"`
}
