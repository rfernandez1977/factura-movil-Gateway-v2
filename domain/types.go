package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DocumentoTributario representa un documento tributario base
type DocumentoTributario struct {
	ID                  primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	TipoDocumento       string             `json:"tipo_documento" bson:"tipo_documento"`
	Folio               int64              `json:"folio" bson:"folio"`
	FechaEmision        time.Time          `json:"fecha_emision" bson:"fecha_emision"`
	RutEmisor           string             `json:"rut_emisor" bson:"rut_emisor"`
	RazonSocialEmisor   string             `json:"razon_social_emisor" bson:"razon_social_emisor"`
	RutReceptor         string             `json:"rut_receptor" bson:"rut_receptor"`
	RazonSocialReceptor string             `json:"razon_social_receptor" bson:"razon_social_receptor"`
	MontoTotal          float64            `json:"monto_total" bson:"monto_total"`
	MontoNeto           float64            `json:"monto_neto" bson:"monto_neto"`
	MontoExento         float64            `json:"monto_exento" bson:"monto_exento"`
	MontoIVA            float64            `json:"monto_iva" bson:"monto_iva"`
	Estado              string             `json:"estado" bson:"estado"`
	FechaCreacion       time.Time          `json:"fecha_creacion" bson:"fecha_creacion"`
	FechaActualizacion  time.Time          `json:"fecha_actualizacion" bson:"fecha_actualizacion"`
}

// Item representa un ítem de un documento tributario
type Item struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Descripcion string             `json:"descripcion" bson:"descripcion"`
	Cantidad    float64            `json:"cantidad" bson:"cantidad"`
	PrecioUnit  float64            `json:"precio_unit" bson:"precio_unit"`
	MontoNeto   float64            `json:"monto_neto" bson:"monto_neto"`
	MontoIVA    float64            `json:"monto_iva" bson:"monto_iva"`
	MontoTotal  float64            `json:"monto_total" bson:"monto_total"`
}

// EstadoDocumento representa el estado de un documento
type EstadoDocumento struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	DocumentoID primitive.ObjectID `json:"documento_id" bson:"documento_id"`
	Estado      string             `json:"estado" bson:"estado"`
	Fecha       time.Time          `json:"fecha" bson:"fecha"`
	Usuario     string             `json:"usuario" bson:"usuario"`
	Comentario  string             `json:"comentario" bson:"comentario"`
}

// ReferenciaDocumento representa una referencia entre documentos
type ReferenciaDocumento struct {
	ID              primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	TipoOrigen      string             `json:"tipo_origen" bson:"tipo_origen"`
	FolioOrigen     int64              `json:"folio_origen" bson:"folio_origen"`
	TipoReferencia  string             `json:"tipo_referencia" bson:"tipo_referencia"`
	FolioReferencia int64              `json:"folio_referencia" bson:"folio_referencia"`
	FechaCreacion   time.Time          `json:"fecha_creacion" bson:"fecha_creacion"`
}

// CAF representa un Código de Autorización de Folios
type CAF struct {
	ID               primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	TipoDocumento    string             `json:"tipo_documento" bson:"tipo_documento"`
	RangoInicial     int64              `json:"rango_inicial" bson:"rango_inicial"`
	RangoFinal       int64              `json:"rango_final" bson:"rango_final"`
	FolioActual      int64              `json:"folio_actual" bson:"folio_actual"`
	Estado           string             `json:"estado" bson:"estado"`
	FechaCreacion    time.Time          `json:"fecha_creacion" bson:"fecha_creacion"`
	FechaVencimiento time.Time          `json:"fecha_vencimiento" bson:"fecha_vencimiento"`
}
