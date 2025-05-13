package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Transformacion struct {
	ID          int64  `json:"id"`
	TipoOrigen  string `json:"tipo_origen"`
	TipoDestino string `json:"tipo_destino"`
	Reglas      string `json:"reglas"`
	Estado      string `json:"estado"`
}

// ReporteTransformacion representa un reporte de transformaciones
type ReporteTransformacion struct {
	ID                       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	PeriodoInicio            time.Time          `json:"periodo_inicio" bson:"periodo_inicio"`
	PeriodoFin               time.Time          `json:"periodo_fin" bson:"periodo_fin"`
	TotalTransformaciones    int                `json:"total_transformaciones" bson:"total_transformaciones"`
	TransformacionesExitosas int                `json:"transformaciones_exitosas" bson:"transformaciones_exitosas"`
	TransformacionesFallidas int                `json:"transformaciones_fallidas" bson:"transformaciones_fallidas"`
	ErroresPorTipo           map[string]int     `json:"errores_por_tipo" bson:"errores_por_tipo"`
	FechaGeneracion          time.Time          `json:"fecha_generacion" bson:"fecha_generacion"`
}
