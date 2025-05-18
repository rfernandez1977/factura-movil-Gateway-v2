package utils

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"strings"

	"FMgo/models"
	"github.com/lestrrat-go/libxml2"
	"github.com/lestrrat-go/libxml2/xsd"
)

// ValidadorXML representa un validador de esquema XML
type ValidadorXML struct {
	schema *xsd.Schema
}

// NuevoValidadorXML crea un nuevo validador de esquema XML
func NuevoValidadorXML(schemaPath string) (*ValidadorXML, error) {
	// Leer el archivo de esquema
	schemaBytes, err := ioutil.ReadFile(schemaPath)
	if err != nil {
		return nil, fmt.Errorf("error leyendo esquema: %v", err)
	}

	// Parsear el esquema
	schema, err := xsd.Parse(schemaBytes)
	if err != nil {
		return nil, fmt.Errorf("error parseando esquema: %v", err)
	}

	return &ValidadorXML{schema: schema}, nil
}

// ValidarXML valida un documento XML contra el esquema
func (v *ValidadorXML) ValidarXML(xmlData []byte) error {
	// Parsear el documento XML
	doc, err := libxml2.ParseString(string(xmlData))
	if err != nil {
		return fmt.Errorf("error parseando XML: %v", err)
	}
	defer doc.Free()

	// Validar contra el esquema
	if err := v.schema.Validate(doc); err != nil {
		// Procesar errores de validación
		errList := err.(xsd.SchemaValidationError).Errors()
		var errMsgs []string
		for _, validationErr := range errList {
			errMsgs = append(errMsgs, fmt.Sprintf("Error de validación: %s",
				validationErr.Error()))
		}
		return fmt.Errorf("errores de validación:\n%s", strings.Join(errMsgs, "\n"))
	}

	return nil
}

// ValidarDTE valida un DTE completo
func (v *ValidadorXML) ValidarDTE(xmlData []byte) error {
	// Primero validamos que el XML sea válido
	var dte models.DTEXMLModel
	if err := xml.Unmarshal(xmlData, &dte); err != nil {
		return fmt.Errorf("error validando DTE: %v", err)
	}

	// Luego validamos contra el esquema
	return v.ValidarXML(xmlData)
}

// ValidarRespuestaSII valida la respuesta del SII
func ValidarRespuestaSII(xmlData []byte) (*models.RespuestaSII, error) {
	var resp models.RespuestaSII
	if err := xml.Unmarshal(xmlData, &resp); err != nil {
		return nil, fmt.Errorf("error parseando respuesta SII: %v", err)
	}

	if resp.Estado != "ACEPTADO" {
		return &resp, fmt.Errorf("error SII: %s", resp.Glosa)
	}

	return &resp, nil
}
