package tests

import (
	"time"

	"FMgo/models"
)

// ConfigEmpresa contiene la configuración de la empresa para las pruebas
type ConfigEmpresa struct {
	Empresa struct {
		RutEmisor    string   `json:"rutEmisor"`
		RutEnvia     string   `json:"rutEnvia"`
		RutReceptor  string   `json:"rutReceptor"`
		FchResol     string   `json:"fchResol"`
		NroResol     string   `json:"nroResol"`
		RazonSocial  string   `json:"razonSocial"`
		GiroEmis     string   `json:"giroEmis"`
		CorreoEmisor string   `json:"correoEmisor"`
		Acteco       []string `json:"acteco"`
		DirOrigen    string   `json:"dirOrigen"`
		CmnaOrigen   string   `json:"cmnaOrigen"`
		CiudadOrigen string   `json:"ciudadOrigen"`
	} `json:"empresa"`
	Folios struct {
		Factura int `json:"factura"`
	} `json:"folios"`
	Firma struct {
		Path string `json:"path"`
		CAF  string `json:"caf"`
	} `json:"firma"`
}

// Factura representa una factura electrónica
type Factura struct {
	RutEmisor           string
	RutReceptor         string
	RazonSocialEmisor   string
	RazonSocialReceptor string
	Folio               int
	FechaEmision        time.Time
	MontoNeto           float64
	MontoIVA            float64
	MontoTotal          float64
	Items               []models.Item
}
