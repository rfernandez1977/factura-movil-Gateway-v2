package models

import "time"

// Sobre representa un sobre de envío de documentos electrónicos al SII
type Sobre struct {
	ID               string    `json:"id" bson:"_id,omitempty"`
	RUTEnviador      string    `json:"rut_enviador" bson:"rut_enviador"`
	RUTCompania      string    `json:"rut_compania" bson:"rut_compania"`
	FechaResolucion  time.Time `json:"fecha_resolucion" bson:"fecha_resolucion"`
	NumeroResolucion int       `json:"numero_resolucion" bson:"numero_resolucion"`
	FechaEnvio       time.Time `json:"fecha_envio" bson:"fecha_envio"`
	Documento        []byte    `json:"documento" bson:"documento"`
	Token            string    `json:"token" bson:"token,omitempty"`
	Firma            []byte    `json:"firma" bson:"firma,omitempty"`
	TrackID          string    `json:"track_id" bson:"track_id,omitempty"`
	Estado           string    `json:"estado" bson:"estado"`
	CreatedAt        time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" bson:"updated_at"`
}

// NewSobre crea un nuevo sobre para envío al SII
func NewSobre(rutEnviador, rutCompania string, fechaResolucion time.Time, numeroResolucion int, documento []byte) *Sobre {
	return &Sobre{
		ID:               GenerateID(),
		RUTEnviador:      rutEnviador,
		RUTCompania:      rutCompania,
		FechaResolucion:  fechaResolucion,
		NumeroResolucion: numeroResolucion,
		FechaEnvio:       time.Now(),
		Documento:        documento,
		Estado:           "PENDIENTE",
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
}

// Validate valida que el sobre tenga todos los campos requeridos
func (s *Sobre) Validate() error {
	if s.RUTEnviador == "" {
		return &ValidationFieldError{Field: "rut_enviador", Message: "El RUT del enviador es obligatorio"}
	}
	if s.RUTCompania == "" {
		return &ValidationFieldError{Field: "rut_compania", Message: "El RUT de la compañía es obligatorio"}
	}
	if s.FechaResolucion.IsZero() {
		return &ValidationFieldError{Field: "fecha_resolucion", Message: "La fecha de resolución es obligatoria"}
	}
	if s.NumeroResolucion <= 0 {
		return &ValidationFieldError{Field: "numero_resolucion", Message: "El número de resolución debe ser mayor a cero"}
	}
	if len(s.Documento) == 0 {
		return &ValidationFieldError{Field: "documento", Message: "El documento es obligatorio"}
	}
	return nil
}
