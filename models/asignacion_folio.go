package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AsignacionFolio representa la asignación de un folio a un documento
type AsignacionFolio struct {
	ID              primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	TipoDocumento   string             `json:"tipoDocumento" bson:"tipo_documento"`
	Folio           int                `json:"folio" bson:"folio"`
	FechaAsignacion time.Time          `json:"fechaAsignacion" bson:"fecha_asignacion"`
	Usuario         string             `json:"usuario" bson:"usuario"`
	EstadoUso       string             `json:"estadoUso" bson:"estado_uso"`
	RutEmisor       string             `json:"rutEmisor" bson:"rut_emisor,omitempty"`
	CreatedAt       time.Time          `json:"createdAt" bson:"created_at"`
	UpdatedAt       time.Time          `json:"updatedAt" bson:"updated_at"`
}

// ControlFolio representa el control de folios para un tipo de documento
type ControlFolio struct {
	ID                primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	TipoDocumento     string             `json:"tipoDocumento" bson:"tipo_documento"`
	RangoInicial      int                `json:"rangoInicial" bson:"rango_inicial"`
	RangoFinal        int                `json:"rangoFinal" bson:"rango_final"`
	FolioActual       int                `json:"folioActual" bson:"folio_actual"`
	FoliosDisponibles int                `json:"foliosDisponibles" bson:"folios_disponibles"`
	UltimoUso         time.Time          `json:"ultimoUso" bson:"ultimo_uso"`
	EstadoCAF         string             `json:"estadoCAF" bson:"estado_caf"`
	AlertaGenerada    bool               `json:"alertaGenerada" bson:"alerta_generada"`
	RutEmisor         string             `json:"rutEmisor" bson:"rut_emisor"`
	CreatedAt         time.Time          `json:"createdAt" bson:"created_at"`
	UpdatedAt         time.Time          `json:"updatedAt" bson:"updated_at"`
}
