package models

type RegistroTransformacion struct {
	ID   string `json:"id" bson:"_id"`
	Tipo string `json:"tipo" bson:"tipo"`
}
