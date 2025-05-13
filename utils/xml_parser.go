package utils

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"

	"github.com/cursor/FMgo/models"
)

// ParserXML representa un parser de documentos XML
type ParserXML struct {
	version string
}

// NuevoParserXML crea un nuevo parser de documentos XML
func NuevoParserXML(version string) *ParserXML {
	return &ParserXML{
		version: version,
	}
}

// ParsearDTE parsea un documento DTE desde XML
func (p *ParserXML) ParsearDTE(xmlData []byte) (*models.DTEXMLModel, error) {
	var dte models.DTEXMLModel
	if err := xml.Unmarshal(xmlData, &dte); err != nil {
		return nil, fmt.Errorf("error parseando DTE: %v", err)
	}
	return &dte, nil
}

// ParsearSobreDTE parsea un sobre DTE desde XML
func (p *ParserXML) ParsearSobreDTE(xmlData []byte) (*models.SobreDTEModel, error) {
	var sobre models.SobreDTEModel
	if err := xml.Unmarshal(xmlData, &sobre); err != nil {
		return nil, fmt.Errorf("error parseando sobre DTE: %v", err)
	}
	return &sobre, nil
}

// ParsearDTEFromFile parsea un documento DTE desde un archivo XML
func (p *ParserXML) ParsearDTEFromFile(filePath string) (*models.DTEXMLModel, error) {
	xmlData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error leyendo archivo: %v", err)
	}
	return p.ParsearDTE(xmlData)
}

// ParsearSobreDTEFromFile parsea un sobre DTE desde un archivo XML
func (p *ParserXML) ParsearSobreDTEFromFile(filePath string) (*models.SobreDTEModel, error) {
	xmlData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error leyendo archivo: %v", err)
	}
	return p.ParsearSobreDTE(xmlData)
}
