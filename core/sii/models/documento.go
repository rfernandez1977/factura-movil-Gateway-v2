package models

import (
	"fmt"
	"time"
)

// Documento representa un documento tributario electrónico
type Documento struct {
	ID           string    `json:"id"`
	TipoDTE      string    `json:"tipo_dte"`
	Folio        int64     `json:"folio"`
	RutEmisor    string    `json:"rut_emisor"`
	RutReceptor  string    `json:"rut_receptor"`
	FechaEmision time.Time `json:"fecha_emision"`
	MontoTotal   float64   `json:"monto_total"`
	XML          []byte    `json:"xml"`
	Estado       string    `json:"estado"`
}

// Validar verifica que el documento tenga todos los campos requeridos
func (d *Documento) Validar() error {
	if d.TipoDTE == "" {
		return fmt.Errorf("tipo DTE es requerido")
	}
	if d.Folio <= 0 {
		return fmt.Errorf("folio debe ser mayor a 0")
	}
	if d.RutEmisor == "" {
		return fmt.Errorf("RUT emisor es requerido")
	}
	if d.RutReceptor == "" {
		return fmt.Errorf("RUT receptor es requerido")
	}
	if d.MontoTotal <= 0 {
		return fmt.Errorf("monto total debe ser mayor a 0")
	}
	return nil
}

// RespuestaEnvio representa la respuesta al enviar un documento
type RespuestaEnvio struct {
	TrackID   string    `json:"track_id"`
	Estado    string    `json:"estado"`
	Glosa     string    `json:"glosa"`
	Timestamp time.Time `json:"timestamp"`
}

// EstadoEnvio representa el estado de un envío
type EstadoEnvio struct {
	TrackID   string    `json:"track_id"`
	Estado    string    `json:"estado"`
	Glosa     string    `json:"glosa"`
	Timestamp time.Time `json:"timestamp"`
}

// ResultadoValidacion representa el resultado de validar un documento
type ResultadoValidacion struct {
	Folio     int64     `json:"folio"`
	TipoDTE   string    `json:"tipo_dte"`
	Estado    string    `json:"estado"`
	Glosa     string    `json:"glosa"`
	Timestamp time.Time `json:"timestamp"`
}
