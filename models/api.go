package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// API representa una configuración de API
type API struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	Version     string             `json:"version" bson:"version"`
	Enabled     bool               `json:"enabled" bson:"enabled"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

// VersionAPI representa una versión de la API
type VersionAPI struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	APIID       primitive.ObjectID `json:"api_id" bson:"api_id"`
	Version     string             `json:"version" bson:"version"`
	Description string             `json:"description" bson:"description"`
	Enabled     bool               `json:"enabled" bson:"enabled"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

// RegistroAPI representa un registro de uso de la API
type RegistroAPI struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	APIID        primitive.ObjectID `json:"api_id" bson:"api_id"`
	VersionID    primitive.ObjectID `json:"version_id" bson:"version_id"`
	Endpoint     string             `json:"endpoint" bson:"endpoint"`
	Method       string             `json:"method" bson:"method"`
	Status       int                `json:"status" bson:"status"`
	ResponseTime int64              `json:"response_time" bson:"response_time"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
}

// ReporteAPI representa un reporte de uso de la API
type ReporteAPI struct {
	ID                      primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	PeriodoInicio           time.Time          `json:"periodo_inicio" bson:"periodo_inicio"`
	PeriodoFin              time.Time          `json:"periodo_fin" bson:"periodo_fin"`
	TotalLlamadas           int                `json:"total_llamadas" bson:"total_llamadas"`
	LlamadasExitosas        int                `json:"llamadas_exitosas" bson:"llamadas_exitosas"`
	LlamadasFallidas        int                `json:"llamadas_fallidas" bson:"llamadas_fallidas"`
	TiempoPromedioRespuesta time.Duration      `json:"tiempo_promedio_respuesta" bson:"tiempo_promedio_respuesta"`
	ErroresPorTipo          map[string]int     `json:"errores_por_tipo" bson:"errores_por_tipo"`
	FechaGeneracion         time.Time          `json:"fecha_generacion" bson:"fecha_generacion"`
}
