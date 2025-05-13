package models

import "time"

// CSF representa un Código Seguimiento de Folios
type CSF struct {
	ID            string    `json:"id" bson:"_id,omitempty"`
	EmpresaID     string    `json:"empresa_id" bson:"empresa_id"`
	TipoDocumento string    `json:"tipo_documento" bson:"tipo_documento"`
	FolioInicial  int       `json:"folio_inicial" bson:"folio_inicial"`
	FolioFinal    int       `json:"folio_final" bson:"folio_final"`
	FolioActual   int       `json:"folio_actual" bson:"folio_actual"`
	Activo        bool      `json:"activo" bson:"activo"`
	CreatedAt     time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" bson:"updated_at"`
}
