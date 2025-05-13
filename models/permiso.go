package models

type Permiso struct {
	ID   string `json:"id" bson:"_id"`
	Name string `json:"name" bson:"name"`
}
