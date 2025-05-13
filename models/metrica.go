package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MetricaIntegracion representa una métrica de integración
type MetricaIntegracion struct {
	ID         primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	Tipo       string                 `bson:"tipo" json:"tipo"`
	Valor      float64                `bson:"valor" json:"valor"`
	Timestamp  time.Time              `bson:"timestamp" json:"timestamp"`
	Componente string                 `bson:"componente" json:"componente"`
	Estado     string                 `bson:"estado" json:"estado"`
	Metadata   map[string]interface{} `bson:"metadata,omitempty" json:"metadata,omitempty"`
	CreatedAt  time.Time              `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time              `bson:"updated_at" json:"updated_at"`
}

// Metric representa una métrica genérica
type Metric struct {
	Name      string                 `json:"name"`
	Value     float64                `json:"value"`
	Labels    map[string]string      `json:"labels,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}
