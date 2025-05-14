package models

import "time"

// ErrorSII representa un error del SII en la respuesta
type ErrorSII struct {
	Codigo      string `xml:"Codigo" json:"codigo"`
	Descripcion string `xml:"Descripcion" json:"descripcion"`
	Detalle     string `xml:"Detalle" json:"detalle"`
}

// RespuestaSII representa la respuesta del SII a una consulta o envío
type RespuestaSII struct {
	TrackID      string     `xml:"TRACKID" json:"track_id"`
	Estado       string     `xml:"ESTADO" json:"estado"`
	Glosa        string     `xml:"GLOSA" json:"glosa"`
	NumAtencion  string     `xml:"NUMATENCION,omitempty" json:"num_atencion,omitempty"`
	FechaProceso time.Time  `xml:"FECHA_PROCESO" json:"fecha_proceso"`
	Errores      []ErrorSII `xml:"ERRORES>ERROR,omitempty" json:"errores,omitempty"`
}
