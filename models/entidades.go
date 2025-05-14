package models

import "time"

// Emisor representa a un emisor de documentos tributarios
type Emisor struct {
	ID            string    `json:"id" bson:"_id,omitempty"`
	RUT           string    `json:"rut" bson:"rut"`
	RazonSocial   string    `json:"razon_social" bson:"razon_social"`
	GiroComercial string    `json:"giro_comercial" bson:"giro_comercial"`
	Direccion     string    `json:"direccion" bson:"direccion"`
	Comuna        string    `json:"comuna" bson:"comuna"`
	Ciudad        string    `json:"ciudad" bson:"ciudad"`
	Telefono      string    `json:"telefono,omitempty" bson:"telefono,omitempty"`
	Email         string    `json:"email,omitempty" bson:"email,omitempty"`
	CreatedAt     time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" bson:"updated_at"`
}

// Receptor representa a un receptor de documentos tributarios
type Receptor struct {
	ID            string    `json:"id" bson:"_id,omitempty"`
	RUT           string    `json:"rut" bson:"rut"`
	RazonSocial   string    `json:"razon_social" bson:"razon_social"`
	GiroComercial string    `json:"giro_comercial,omitempty" bson:"giro_comercial,omitempty"`
	Direccion     string    `json:"direccion" bson:"direccion"`
	Comuna        string    `json:"comuna,omitempty" bson:"comuna,omitempty"`
	Ciudad        string    `json:"ciudad,omitempty" bson:"ciudad,omitempty"`
	Telefono      string    `json:"telefono,omitempty" bson:"telefono,omitempty"`
	Email         string    `json:"email,omitempty" bson:"email,omitempty"`
	CreatedAt     time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" bson:"updated_at"`
}

// EstadoSII representa el estado de un documento en el SII
type EstadoSII struct {
	Estado      string            `json:"estado"`
	Glosa       string            `json:"glosa"`
	Codigo      int               `json:"codigo"`
	Descripcion string            `json:"descripcion"`
	Timestamp   time.Time         `json:"timestamp"`
	TrackID     string            `json:"track_id,omitempty"`
	Errores     []ErrorReporteSII `json:"errores,omitempty"`
}
