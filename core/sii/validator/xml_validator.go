package validator

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/lestrrat-go/libxml2/parser"
	"github.com/lestrrat-go/libxml2/xsd"
)

// XMLValidator implementa la validación de documentos XML contra esquemas XSD
type XMLValidator struct {
	schemas map[string]*xsd.Schema
}

// NewXMLValidator crea una nueva instancia del validador XML
func NewXMLValidator(schemaPath string) (*XMLValidator, error) {
	validator := &XMLValidator{
		schemas: make(map[string]*xsd.Schema),
	}

	// Cargar los esquemas XSD principales
	schemaFiles := []string{
		"DTE_v10.xsd",
		"EnvioDTE_v10.xsd",
		"SiiTypes_v10.xsd",
		"xmldsignature_v10.xsd",
	}

	for _, filename := range schemaFiles {
		schema, err := loadSchema(filepath.Join(schemaPath, filename))
		if err != nil {
			return nil, fmt.Errorf("error al cargar esquema %s: %w", filename, err)
		}
		validator.schemas[filename] = schema
	}

	return validator, nil
}

// ValidateDTE valida un DTE contra su esquema XSD
func (v *XMLValidator) ValidateDTE(doc interface{}) error {
	// Convertir el documento a XML
	xmlData, err := xml.MarshalIndent(doc, "", "  ")
	if err != nil {
		return fmt.Errorf("error al serializar documento: %w", err)
	}

	// Parsear el XML
	p := parser.New()
	xmlDoc, err := p.ParseString(string(xmlData))
	if err != nil {
		return fmt.Errorf("error al parsear XML: %w", err)
	}
	defer xmlDoc.Free()

	// Validar contra el esquema DTE
	schema := v.schemas["DTE_v10.xsd"]
	if schema == nil {
		return fmt.Errorf("esquema DTE_v10.xsd no encontrado")
	}

	if err := schema.Validate(xmlDoc); err != nil {
		return fmt.Errorf("error de validación: %w", err)
	}

	return nil
}

// ValidateEnvioDTE valida un sobre de envío contra su esquema XSD
func (v *XMLValidator) ValidateEnvioDTE(doc interface{}) error {
	// Convertir el documento a XML
	xmlData, err := xml.MarshalIndent(doc, "", "  ")
	if err != nil {
		return fmt.Errorf("error al serializar documento: %w", err)
	}

	// Parsear el XML
	p := parser.New()
	xmlDoc, err := p.ParseString(string(xmlData))
	if err != nil {
		return fmt.Errorf("error al parsear XML: %w", err)
	}
	defer xmlDoc.Free()

	// Validar contra el esquema EnvioDTE
	schema := v.schemas["EnvioDTE_v10.xsd"]
	if schema == nil {
		return fmt.Errorf("esquema EnvioDTE_v10.xsd no encontrado")
	}

	if err := schema.Validate(xmlDoc); err != nil {
		return fmt.Errorf("error de validación: %w", err)
	}

	return nil
}

// loadSchema carga un esquema XSD desde un archivo
func loadSchema(path string) (*xsd.Schema, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error al leer archivo XSD: %w", err)
	}

	s, err := xsd.Parse(data)
	if err != nil {
		return nil, fmt.Errorf("error al parsear XSD: %w", err)
	}

	return s, nil
}

// Close libera los recursos utilizados por el validador
func (v *XMLValidator) Close() {
	for _, schema := range v.schemas {
		schema.Free()
	}
}
