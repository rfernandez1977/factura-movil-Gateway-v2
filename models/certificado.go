package models

import "time"

// CertificadoDigital representa un certificado digital para firma electrónica
type CertificadoDigital struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	EmpresaID string    `json:"empresa_id" bson:"empresa_id"`
	RutFirma  string    `json:"rut_firma" bson:"rut_firma"`
	Nombre    string    `json:"nombre" bson:"nombre"`
	Contenido []byte    `json:"contenido" bson:"contenido"`
	Password  string    `json:"password" bson:"password,omitempty"`
	Vigencia  time.Time `json:"vigencia" bson:"vigencia"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

// CAF representa un Código de Autorización de Folios
type CAF struct {
	ID            string    `json:"id" bson:"_id,omitempty"`
	EmpresaID     string    `json:"empresa_id" bson:"empresa_id"`
	TipoDocumento string    `json:"tipo_documento" bson:"tipo_documento"`
	FolioInicial  int       `json:"folio_inicial" bson:"folio_inicial"`
	FolioFinal    int       `json:"folio_final" bson:"folio_final"`
	FolioActual   int       `json:"folio_actual" bson:"folio_actual"`
	Estado        string    `json:"estado" bson:"estado"`
	Contenido     []byte    `json:"contenido" bson:"contenido"`
	CreatedAt     time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" bson:"updated_at"`
}

// ConfiguracionCertificado contiene la configuración para generar un certificado
type ConfiguracionCertificado struct {
	Nombre       string `json:"nombre"`
	Organizacion string `json:"organizacion"`
}
