package models

type Sucursal struct {
	ID        string `json:"id" bson:"_id"`
	Name      string `json:"name" bson:"name"`
	Direccion string `json:"direccion" bson:"direccion"`
}
