package utils

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"strings"

	"FMgo/models"
)

// LeerXML lee un archivo XML y lo convierte a una estructura
func LeerXML(filePath string, v interface{}) error {
	xmlData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error leyendo archivo: %v", err)
	}

	if err := xml.Unmarshal(xmlData, v); err != nil {
		return fmt.Errorf("error parseando XML: %v", err)
	}

	return nil
}

// EscribirXML escribe una estructura a un archivo XML
func EscribirXML(filePath string, v interface{}) error {
	xmlData, err := xml.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("error generando XML: %v", err)
	}

	xmlStr := fmt.Sprintf(`<?xml version="1.0" encoding="ISO-8859-1"?>
%s`, string(xmlData))

	if err := ioutil.WriteFile(filePath, []byte(xmlStr), 0644); err != nil {
		return fmt.Errorf("error escribiendo archivo: %v", err)
	}

	return nil
}

// ValidarDTE valida un DTE
func ValidarDTE(dte *models.DTEXMLModel) error {
	if dte == nil {
		return fmt.Errorf("DTE no puede ser nil")
	}

	if strings.TrimSpace(dte.Documento.ID) == "" {
		return fmt.Errorf("ID no puede estar vacío")
	}

	if strings.TrimSpace(dte.Documento.Encabezado.IdDoc.TipoDTE) == "" {
		return fmt.Errorf("tipo DTE no puede estar vacío")
	}

	if strings.TrimSpace(dte.Documento.Encabezado.Emisor.RUT) == "" {
		return fmt.Errorf("RUT emisor no puede estar vacío")
	}

	if strings.TrimSpace(dte.Documento.Encabezado.Receptor.RUT) == "" {
		return fmt.Errorf("RUT receptor no puede estar vacío")
	}

	return nil
}
