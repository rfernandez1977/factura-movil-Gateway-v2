package services

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/lestrrat-go/libxml2"
	"github.com/lestrrat-go/libxml2/xsd"
)

// XMLValidator maneja la validación de XML contra esquemas XSD
type XMLValidator struct {
	schemas map[string]*xsd.Schema
}

// NewXMLValidator crea una nueva instancia del validador XML
func NewXMLValidator(schemasPaths map[string]string) (*XMLValidator, error) {
	schemas := make(map[string]*xsd.Schema)

	for tipo, path := range schemasPaths {
		schema, err := xsd.ParseFromFile(path)
		if err != nil {
			return nil, fmt.Errorf("error cargando esquema %s: %w", tipo, err)
		}
		schemas[tipo] = schema
	}

	return &XMLValidator{schemas: schemas}, nil
}

// cargarSchema carga un esquema XSD desde una URL o archivo local
func (v *XMLValidator) cargarSchema(url string) (*xsd.Schema, error) {
	// Primero intentar cargar desde archivo local
	localPath := filepath.Join("schemas", filepath.Base(url))
	if _, err := os.Stat(localPath); err == nil {
		schema, err := xsd.ParseFromFile(localPath)
		if err == nil {
			return schema, nil
		}
	}

	// Si no existe localmente, descargar de la URL
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error descargando schema: %w", err)
	}
	defer resp.Body.Close()

	// Leer el contenido del schema
	schemaData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error leyendo schema: %w", err)
	}

	// Parsear el schema
	schema, err := xsd.Parse(schemaData)
	if err != nil {
		return nil, fmt.Errorf("error parseando schema: %w", err)
	}

	// Guardar el schema localmente para uso futuro
	if err := os.MkdirAll("schemas", 0755); err != nil {
		return nil, fmt.Errorf("error creando directorio schemas: %w", err)
	}

	if err := os.WriteFile(localPath, schemaData, 0644); err != nil {
		return nil, fmt.Errorf("error guardando schema: %w", err)
	}

	return schema, nil
}

// ValidateXML valida un XML contra el esquema correspondiente
func (v *XMLValidator) ValidateXML(xmlData []byte, schemaType string) error {
	schema, ok := v.schemas[schemaType]
	if !ok {
		return fmt.Errorf("esquema no encontrado: %s", schemaType)
	}

	// Convertir []byte a Document
	doc, err := libxml2.Parse(xmlData)
	if err != nil {
		return fmt.Errorf("error parseando XML: %w", err)
	}
	defer doc.Free()

	return schema.Validate(doc)
}

// ValidarEstructura valida la estructura básica del XML
func (v *XMLValidator) ValidarEstructura(xmlData []byte) error {
	// Verificar que sea un XML válido
	decoder := xml.NewDecoder(bytes.NewReader(xmlData))
	for {
		_, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("error decodificando XML: %w", err)
		}
	}

	return nil
}

// ValidarCamposObligatorios valida que los campos obligatorios estén presentes
func (v *XMLValidator) ValidarCamposObligatorios(xmlData []byte, tipoDocumento string) error {
	// TODO: Implementar validación de campos obligatorios según tipo de documento
	return nil
}

// ValidarValores valida los valores de los campos según las reglas del SII
func (v *XMLValidator) ValidarValores(xmlData []byte, tipoDocumento string) error {
	// TODO: Implementar validación de valores según reglas del SII
	return nil
}
