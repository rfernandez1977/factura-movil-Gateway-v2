package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Alerta representa una alerta del sistema
type Alerta struct {
	ID        primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	Level     string                 `bson:"level" json:"level"`
	Message   string                 `bson:"message" json:"message"`
	Component string                 `bson:"component" json:"component"`
	Timestamp time.Time              `bson:"timestamp" json:"timestamp"`
	Metadata  map[string]interface{} `bson:"metadata,omitempty" json:"metadata,omitempty"`
	CreatedAt time.Time              `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time              `bson:"updated_at" json:"updated_at"`
}

// AlertServiceConfig representa la configuración de un servicio de alertas
type AlertServiceConfig struct {
	URL string
}

// NewAlertServiceConfig crea una nueva configuración para el servicio de alertas
func NewAlertServiceConfig(url string) *AlertServiceConfig {
	return &AlertServiceConfig{
		URL: url,
	}
}

// Send envía una alerta al servicio
func (s *AlertServiceConfig) Send(message string) error {
	// TODO: Implementar el envío de alertas al servicio
	return nil
}
