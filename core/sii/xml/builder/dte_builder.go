package builder

import (
	"encoding/xml"
	"fmt"
	"time"

	"FMgo/core/sii/models"
)

// DTEBuilder construye documentos XML para el SII
type DTEBuilder struct {
	doc *models.Documento
}

// NewDTEBuilder crea una nueva instancia del builder
func NewDTEBuilder(doc *models.Documento) *DTEBuilder {
	return &DTEBuilder{
		doc: doc,
	}
}

// BuildEnvelope construye el sobre XML para envío al SII
func (b *DTEBuilder) BuildEnvelope() ([]byte, error) {
	envelope := &Envelope{
		SetDTE: SetDTE{
			ID: fmt.Sprintf("SetDoc_%s", b.doc.ID),
			Caratula: Caratula{
				RutEmisor:   b.doc.RutEmisor,
				RutEnvia:    b.doc.RutEmisor,
				RutReceptor: b.doc.RutReceptor,
				FechaEnvio:  time.Now(),
				Version:     "1.0",
				SubTotDTE:   []SubTotDTE{{TipoDTE: b.doc.TipoDTE, NroDTE: 1}},
			},
			DTE: []DTE{{
				Documento: b.doc,
			}},
		},
	}

	return xml.MarshalIndent(envelope, "", "  ")
}

// Envelope representa el sobre XML para envío al SII
type Envelope struct {
	XMLName xml.Name `xml:"EnvioDTE"`
	SetDTE  SetDTE   `xml:"SetDTE"`
}

// SetDTE representa el conjunto de documentos a enviar
type SetDTE struct {
	ID       string   `xml:"ID,attr"`
	Caratula Caratula `xml:"Caratula"`
	DTE      []DTE    `xml:"DTE"`
}

// Caratula representa la carátula del envío
type Caratula struct {
	RutEmisor   string      `xml:"RutEmisor"`
	RutEnvia    string      `xml:"RutEnvia"`
	RutReceptor string      `xml:"RutReceptor"`
	FechaEnvio  time.Time   `xml:"FechaEnvio"`
	Version     string      `xml:"Version"`
	SubTotDTE   []SubTotDTE `xml:"SubTotDTE"`
}

// SubTotDTE representa el subtotal por tipo de DTE
type SubTotDTE struct {
	TipoDTE string `xml:"TipoDTE"`
	NroDTE  int    `xml:"NroDTE"`
}

// DTE representa un documento tributario electrónico
type DTE struct {
	Documento *models.Documento `xml:"Documento"`
}
