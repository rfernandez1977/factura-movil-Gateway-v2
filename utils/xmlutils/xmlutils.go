package xmlutils

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"os"
	"strings"

	"github.com/lestrrat-go/libxml2"
	"github.com/lestrrat-go/libxml2/xsd"
)

// XMLParser proporciona métodos para manejar XML
type XMLParser struct {
	strict    bool
	validator *XMLValidator
}

// NewXMLParser crea una nueva instancia de XMLParser
func NewXMLParser(strict bool) *XMLParser {
	return &XMLParser{
		strict: strict,
	}
}

// SetValidator establece un validador XSD para el parser
func (p *XMLParser) SetValidator(schemaPath string) error {
	validator, err := NewXMLValidator(schemaPath)
	if err != nil {
		return fmt.Errorf("error configurando validador: %w", err)
	}
	p.validator = validator
	return nil
}

// ParseSOAP parsea una respuesta SOAP en la estructura correspondiente
func (p *XMLParser) ParseSOAP(data []byte, v interface{}) error {
	// Limpiar namespaces para simplificar el parsing
	cleanXML := p.cleanSOAPNamespaces(data)

	decoder := xml.NewDecoder(bytes.NewReader(cleanXML))
	decoder.Strict = p.strict

	if err := decoder.Decode(v); err != nil {
		return fmt.Errorf("error decodificando SOAP: %w", err)
	}

	return nil
}

// ParseXML parsea un XML en la estructura correspondiente
func (p *XMLParser) ParseXML(data []byte, v interface{}) error {
	// Validar si hay un esquema configurado
	if p.validator != nil {
		if err := p.validator.ValidateXML(data); err != nil {
			return fmt.Errorf("error de validación: %w", err)
		}
	}

	decoder := xml.NewDecoder(bytes.NewReader(data))
	decoder.Strict = p.strict

	if err := decoder.Decode(v); err != nil {
		return fmt.Errorf("error decodificando XML: %w", err)
	}

	return nil
}

// GenerateXML genera un XML a partir de una estructura
func (p *XMLParser) GenerateXML(v interface{}) ([]byte, error) {
	buf := &bytes.Buffer{}

	// Escribir cabecera XML
	buf.WriteString(xml.Header)

	encoder := xml.NewEncoder(buf)
	encoder.Indent("", "  ")

	if err := encoder.Encode(v); err != nil {
		return nil, fmt.Errorf("error codificando XML: %w", err)
	}

	result := buf.Bytes()

	// Validar si hay un esquema configurado
	if p.validator != nil {
		if err := p.validator.ValidateXML(result); err != nil {
			return nil, fmt.Errorf("error de validación: %w", err)
		}
	}

	return result, nil
}

// SaveToFile guarda un XML en un archivo
func (p *XMLParser) SaveToFile(v interface{}, filename string) error {
	xmlData, err := p.GenerateXML(v)
	if err != nil {
		return err
	}

	return os.WriteFile(filename, xmlData, 0644)
}

// LoadFromFile carga un XML desde un archivo
func (p *XMLParser) LoadFromFile(filename string, v interface{}) error {
	xmlData, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("error leyendo archivo: %w", err)
	}

	return p.ParseXML(xmlData, v)
}

// cleanSOAPNamespaces limpia los namespaces de un XML SOAP
func (p *XMLParser) cleanSOAPNamespaces(data []byte) []byte {
	xmlStr := string(data)
	replacements := map[string]string{
		"soap:":                           "",
		"ns1:":                            "",
		"ns2:":                            "",
		`xmlns="http://DefaultNamespace"`: "",
		`xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/"`: "",
	}

	for old, new := range replacements {
		xmlStr = strings.ReplaceAll(xmlStr, old, new)
	}

	return []byte(xmlStr)
}

// XMLValidator proporciona validación contra esquemas XSD
type XMLValidator struct {
	schema *xsd.Schema
}

// NewXMLValidator crea un nuevo validador XML
func NewXMLValidator(schemaPath string) (*XMLValidator, error) {
	schemaBytes, err := os.ReadFile(schemaPath)
	if err != nil {
		return nil, fmt.Errorf("error leyendo esquema: %w", err)
	}

	schema, err := xsd.Parse(schemaBytes)
	if err != nil {
		return nil, fmt.Errorf("error parseando esquema: %w", err)
	}

	return &XMLValidator{
		schema: schema,
	}, nil
}

// ValidateXML valida un documento XML contra el esquema
func (v *XMLValidator) ValidateXML(data []byte) error {
	doc, err := libxml2.ParseString(string(data))
	if err != nil {
		return fmt.Errorf("error parseando XML: %w", err)
	}
	defer doc.Free()

	if err := v.schema.Validate(doc); err != nil {
		validationErr := err.(xsd.SchemaValidationError)
		var errMsgs []string
		for _, e := range validationErr.Errors() {
			errMsgs = append(errMsgs, e.Error())
		}
		return fmt.Errorf("errores de validación:\n%s", strings.Join(errMsgs, "\n"))
	}

	return nil
}
