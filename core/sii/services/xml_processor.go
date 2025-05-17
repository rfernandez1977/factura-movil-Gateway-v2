package services

import (
	"encoding/xml"
	"fmt"
	"regexp"
	"strings"

	"github.com/fmgo/core/sii/logger"
)

// XMLProcessor proporciona funcionalidades para procesar documentos XML
type XMLProcessor struct {
	// Namespaces comunes en documentos SII
	nsMap map[string]string
	log   *logger.Logger
}

// NewXMLProcessor crea una nueva instancia de XMLProcessor
func NewXMLProcessor(log *logger.Logger) *XMLProcessor {
	return &XMLProcessor{
		nsMap: map[string]string{
			"ds":  "http://www.w3.org/2000/09/xmldsig#",
			"sii": "http://www.sii.cl/SiiDte",
		},
		log: log,
	}
}

// extraerCertificado extrae el certificado X509 del documento XML
func (p *XMLProcessor) extraerCertificado(xmlData []byte) (string, error) {
	p.log.Debug("Iniciando extracción de certificado")

	// Usar expresión regular para extraer el certificado
	certRegex := regexp.MustCompile(`<X509Certificate>([^<]+)</X509Certificate>`)
	matches := certRegex.FindSubmatch(xmlData)
	if len(matches) < 2 {
		err := fmt.Errorf("certificado no encontrado en el documento")
		p.log.LogXMLOperation("extraerCertificado", xmlData, err)
		return "", err
	}

	p.log.Debug("Certificado extraído exitosamente")
	return string(matches[1]), nil
}

// extraerFirma extrae el valor de la firma del documento XML
func (p *XMLProcessor) extraerFirma(xmlData []byte) (string, error) {
	p.log.Debug("Iniciando extracción de firma")

	// Usar expresión regular para extraer la firma
	signatureRegex := regexp.MustCompile(`<SignatureValue>([^<]+)</SignatureValue>`)
	matches := signatureRegex.FindSubmatch(xmlData)
	if len(matches) < 2 {
		err := fmt.Errorf("firma no encontrada en el documento")
		p.log.LogXMLOperation("extraerFirma", xmlData, err)
		return "", err
	}

	p.log.Debug("Firma extraída exitosamente")
	return string(matches[1]), nil
}

// validarEstructuraXML valida la estructura básica del documento XML
func (p *XMLProcessor) validarEstructuraXML(xmlData []byte) error {
	p.log.Debug("Iniciando validación de estructura XML")

	var doc struct {
		XMLName xml.Name `xml:"DTE"`
	}

	if err := xml.Unmarshal(xmlData, &doc); err != nil {
		p.log.LogXMLOperation("validarEstructura", xmlData, err)
		return fmt.Errorf("documento XML inválido: %w", err)
	}

	// Validar presencia de elementos requeridos
	requiredElements := []string{
		"<Documento",
		"<Encabezado",
		"<IdDoc",
		"<Emisor",
		"<Receptor",
	}

	xmlStr := string(xmlData)
	for _, elem := range requiredElements {
		if !strings.Contains(xmlStr, elem) {
			err := fmt.Errorf("elemento requerido no encontrado: %s", elem)
			p.log.LogXMLOperation("validarEstructura", xmlData, err)
			return err
		}
	}

	p.log.Debug("Estructura XML validada exitosamente")
	return nil
}

// limpiarXML limpia el documento XML de caracteres no válidos
func (p *XMLProcessor) limpiarXML(xmlData []byte) []byte {
	p.log.Debug("Iniciando limpieza de XML")

	// Eliminar caracteres de control excepto espacios en blanco
	re := regexp.MustCompile(`[\x00-\x09\x0B\x0C\x0E-\x1F\x7F]`)
	cleanXML := re.ReplaceAll(xmlData, []byte{})

	p.log.Debug("XML limpiado exitosamente")
	return cleanXML
}
