package integration

import "time"

// Documento representa un documento tributario electrónico
type Documento struct {
	ID                string
	TipoDTE           string
	Folio             int64
	FechaEmision      time.Time
	RutEmisor         string
	RazonEmisor       string
	GiroEmisor        string
	DireccionEmisor   string
	ComunaEmisor      string
	RutReceptor       string
	RazonReceptor     string
	GiroReceptor      string
	DireccionReceptor string
	ComunaReceptor    string
	MontoNeto         float64
	MontoExento       float64
	TasaIVA           float64
	MontoIVA          float64
	MontoTotal        float64
	Estado            string
	TrackID           string
	EstadoSII         string
}

// crearDocumentoPrueba crea un documento de prueba
func crearDocumentoPrueba(tipoDTE string) *Documento {
	return &Documento{
		ID:           "TEST-" + tipoDTE,
		TipoDTE:      tipoDTE,
		Folio:        1,
		FechaEmision: time.Now(),
		RutEmisor:    "76.123.456-7",
		RazonEmisor:  "Empresa de Prueba",
		GiroEmisor:   "Servicios Informáticos",
		Estado:       "BORRADOR",
	}
}
