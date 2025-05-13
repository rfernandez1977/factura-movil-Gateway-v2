package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Empresa representa una empresa en el sistema
type Empresa struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	RUT          string             `json:"rut" bson:"rut"`
	RazonSocial  string             `json:"razon_social" bson:"razon_social"`
	Giro         string             `json:"giro" bson:"giro"`
	Direccion    string             `json:"direccion" bson:"direccion"`
	Comuna       string             `json:"comuna" bson:"comuna"`
	Ciudad       string             `json:"ciudad" bson:"ciudad"`
	Correo       string             `json:"correo" bson:"correo"`
	Actecos      []string           `json:"actecos" bson:"actecos"`
	FechaInicio  time.Time          `json:"fecha_inicio" bson:"fecha_inicio"`
	FechaTermino *time.Time         `json:"fecha_termino,omitempty" bson:"fecha_termino,omitempty"`
	Estado       string             `json:"estado" bson:"estado"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at"`
}
