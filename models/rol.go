package models

type Rol struct {
	ID   string `json:"id" bson:"_id"`
	Name string `json:"name" bson:"name"`
}
