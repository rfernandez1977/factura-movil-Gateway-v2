package utils

import (
	"encoding/xml"
	"fmt"
	"time"

	"FMgo/models"
)

// GeneradorXML representa un generador de documentos XML
type GeneradorXML struct {
	version string
}

// NuevoGeneradorXML crea un nuevo generador de documentos XML
func NuevoGeneradorXML(version string) *GeneradorXML {
	return &GeneradorXML{
		version: version,
	}
}

// GenerarDTE genera un documento DTE en formato XML
func (g *GeneradorXML) GenerarDTE(dte *models.DTEXMLModel) ([]byte, error) {
	// Generar el XML
	xmlData, err := xml.MarshalIndent(dte, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error generando XML: %v", err)
	}

	// Agregar la declaración XML
	xmlStr := fmt.Sprintf(`<?xml version="1.0" encoding="ISO-8859-1"?>
%s`, string(xmlData))

	return []byte(xmlStr), nil
}

// GenerarSobreDTE genera un sobre DTE en formato XML
func (g *GeneradorXML) GenerarSobreDTE(dtes []*models.DTEXMLModel, emisor *models.Emisor) ([]byte, error) {
	// Crear caratula
	caratula := &models.Caratula{
		Version:          g.version,
		RutEmisor:        emisor.RUT,
		RutEnvia:         emisor.RUT,
		RutReceptor:      "60803000-K", // SII
		FechaResolucion:  time.Now().Format("2006-01-02"),
		NumeroResolucion: 0, // Este valor debe ser configurado correctamente
		TmstFirmaEnv:     time.Now().Format("2006-01-02T15:04:05"),
	}

	// Crear el SetDTE
	setDTE := &models.SetDTE{
		ID:       fmt.Sprintf("SetDoc_%s", emisor.RUT),
		Caratula: caratula,
	}

	// Convertir dtes a DTEs de la estructura SetDTE
	dtesList := make([]models.DTEXMLModel, len(dtes))
	for i, dte := range dtes {
		dtesList[i] = *dte
	}
	setDTE.DTEs = dtesList

	// Crear el sobre principal
	sobre := &models.SobreDTEModel{
		SetDTE: setDTE,
	}

	// Generar el XML
	xmlData, err := xml.MarshalIndent(sobre, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error generando XML: %v", err)
	}

	// Agregar la declaración XML
	xmlStr := fmt.Sprintf(`<?xml version="1.0" encoding="ISO-8859-1"?>
%s`, string(xmlData))

	return []byte(xmlStr), nil
}
