package services

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"

	"github.com/beevik/etree"
	"github.com/lestrrat-go/libxml2"
	"github.com/lestrrat-go/libxml2/xsd"
)

// XMLService maneja las operaciones relacionadas con XML
type XMLService struct{}

// NewXMLService crea una nueva instancia del servicio XML
func NewXMLService() *XMLService {
	return &XMLService{}
}

// Canonicalizar aplica la transformación de canonicalización C14N al XML
func (s *XMLService) Canonicalizar(xmlData []byte) ([]byte, error) {
	// Crear un nuevo documento
	doc := etree.NewDocument()
	if err := doc.ReadFromBytes(xmlData); err != nil {
		return nil, fmt.Errorf("error al leer XML: %w", err)
	}

	// Aplicar canonicalización C14N
	var buf bytes.Buffer
	if _, err := doc.WriteTo(&buf); err != nil {
		return nil, fmt.Errorf("error en canonicalización: %w", err)
	}

	return buf.Bytes(), nil
}

// AgregarNamespaces agrega los namespaces necesarios al documento XML
func (s *XMLService) AgregarNamespaces(xmlData []byte) ([]byte, error) {
	doc := etree.NewDocument()
	if err := doc.ReadFromBytes(xmlData); err != nil {
		return nil, fmt.Errorf("error al leer XML: %w", err)
	}

	// Obtener el elemento raíz
	root := doc.Root()
	if root == nil {
		return nil, fmt.Errorf("documento XML sin elemento raíz")
	}

	// Agregar namespaces requeridos
	root.CreateAttr("xmlns", "http://www.sii.cl/SiiDte")
	root.CreateAttr("xmlns:xsd", "http://www.w3.org/2001/XMLSchema")
	root.CreateAttr("xmlns:xsi", "http://www.w3.org/2001/XMLSchema-instance")
	root.CreateAttr("xmlns:ds", "http://www.w3.org/2000/09/xmldsig#")
	root.CreateAttr("xsi:schemaLocation", "http://www.sii.cl/SiiDte xmldsignature_v10.xsd")

	var buf bytes.Buffer
	if _, err := doc.WriteTo(&buf); err != nil {
		return nil, fmt.Errorf("error al escribir XML: %w", err)
	}

	return buf.Bytes(), nil
}

// ValidarEstructura valida la estructura del XML contra un esquema XSD
func (s *XMLService) ValidarEstructura(xmlData []byte, xsdPath string) error {
	// Leer el esquema XSD
	xsdBytes, err := ioutil.ReadFile(xsdPath)
	if err != nil {
		return fmt.Errorf("error al leer esquema XSD: %w", err)
	}

	// Crear el esquema
	schema, err := xsd.Parse(xsdBytes)
	if err != nil {
		return fmt.Errorf("error al parsear esquema XSD: %w", err)
	}
	defer schema.Free()

	// Crear el documento XML
	doc, err := libxml2.Parse(xmlData)
	if err != nil {
		return fmt.Errorf("error al parsear XML: %w", err)
	}
	defer doc.Free()

	// Validar el documento contra el esquema
	if err := schema.Validate(doc); err != nil {
		return fmt.Errorf("error de validación XML: %w", err)
	}

	return nil
}

// FormatearXML formatea el XML para mejor legibilidad
func (s *XMLService) FormatearXML(xmlData []byte) ([]byte, error) {
	// Crear buffer para XML formateado
	var prettyXML bytes.Buffer

	// Crear decoder y encoder
	decoder := xml.NewDecoder(bytes.NewReader(xmlData))
	encoder := xml.NewEncoder(&prettyXML)
	encoder.Indent("", "  ")

	// Copiar y formatear
	for {
		token, err := decoder.Token()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, fmt.Errorf("error al leer token XML: %w", err)
		}
		err = encoder.EncodeToken(token)
		if err != nil {
			return nil, fmt.Errorf("error al codificar token XML: %w", err)
		}
	}

	err := encoder.Flush()
	if err != nil {
		return nil, fmt.Errorf("error al finalizar codificación XML: %w", err)
	}

	return prettyXML.Bytes(), nil
}

// ValidarNamespaces verifica que todos los namespaces requeridos estén presentes
func (s *XMLService) ValidarNamespaces(xmlData []byte) error {
	doc := etree.NewDocument()
	if err := doc.ReadFromBytes(xmlData); err != nil {
		return fmt.Errorf("error al leer XML: %w", err)
	}

	root := doc.Root()
	if root == nil {
		return fmt.Errorf("documento XML sin elemento raíz")
	}

	// Lista de namespaces requeridos
	requiredNamespaces := map[string]string{
		"xmlns":     "http://www.sii.cl/SiiDte",
		"xmlns:ds":  "http://www.w3.org/2000/09/xmldsig#",
		"xmlns:xsd": "http://www.w3.org/2001/XMLSchema",
		"xmlns:xsi": "http://www.w3.org/2001/XMLSchema-instance",
	}

	// Verificar cada namespace requerido
	for attrName, expectedValue := range requiredNamespaces {
		attr := root.SelectAttr(attrName)
		if attr == nil {
			return fmt.Errorf("falta namespace requerido: %s", attrName)
		}
		if attr.Value != expectedValue {
			return fmt.Errorf("valor incorrecto para namespace %s: esperado %s, encontrado %s",
				attrName, expectedValue, attr.Value)
		}
	}

	return nil
}
