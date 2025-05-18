package services

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sync"

	"FMgo/core/sii/models"
	"github.com/lestrrat-go/libxml2"
	"github.com/lestrrat-go/libxml2/xsd"
)

// ValidatorService proporciona funcionalidades de validación de documentos XML contra esquemas XSD
type ValidatorService struct {
	schemas     map[string]*xsd.Schema
	schemasLock sync.RWMutex
}

// NewValidatorService crea una nueva instancia del servicio de validación
func NewValidatorService() *ValidatorService {
	return &ValidatorService{
		schemas: make(map[string]*xsd.Schema),
	}
}

// CargarEsquema carga un esquema XSD desde un archivo
func (s *ValidatorService) CargarEsquema(nombre, rutaArchivo string) error {
	contenido, err := ioutil.ReadFile(rutaArchivo)
	if err != nil {
		return fmt.Errorf("error leyendo archivo XSD %s: %w", rutaArchivo, err)
	}

	schema, err := xsd.Parse(contenido)
	if err != nil {
		return fmt.Errorf("error parseando esquema XSD %s: %w", rutaArchivo, err)
	}

	s.schemasLock.Lock()
	s.schemas[nombre] = schema
	s.schemasLock.Unlock()

	return nil
}

// CargarEsquemasBase carga los esquemas base necesarios para validación
func (s *ValidatorService) CargarEsquemasBase(dirBase string) error {
	esquemas := map[string]string{
		"DTE":      "DTE_v10.xsd",
		"EnvioDTE": "EnvioDTE_v10.xsd",
		"SiiTypes": "SiiTypes_v10.xsd",
		"XMLDSig":  "xmldsignature_v10.xsd",
	}

	for nombre, archivo := range esquemas {
		ruta := filepath.Join(dirBase, archivo)
		if err := s.CargarEsquema(nombre, ruta); err != nil {
			return fmt.Errorf("error cargando esquema %s: %w", nombre, err)
		}
	}

	return nil
}

// ValidarDTE valida un DTE contra el esquema XSD correspondiente
func (s *ValidatorService) ValidarDTE(dte *models.DTE) error {
	s.schemasLock.RLock()
	schema, exists := s.schemas["DTE"]
	s.schemasLock.RUnlock()

	if !exists {
		return fmt.Errorf("esquema DTE no cargado")
	}

	// Convertir DTE a XML
	xmlBytes, err := xml.MarshalIndent(dte, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling DTE a XML: %w", err)
	}

	// Parsear XML para validación
	doc, err := libxml2.ParseString(string(xmlBytes))
	if err != nil {
		return fmt.Errorf("error parseando XML del DTE: %w", err)
	}
	defer doc.Free()

	// Validar contra esquema
	if err := schema.Validate(doc); err != nil {
		return fmt.Errorf("error de validación XSD: %w", err)
	}

	return nil
}

// ValidarEnvioDTE valida un envío de DTE contra el esquema correspondiente
func (s *ValidatorService) ValidarEnvioDTE(envioDTE []byte) error {
	s.schemasLock.RLock()
	schema, exists := s.schemas["EnvioDTE"]
	s.schemasLock.RUnlock()

	if !exists {
		return fmt.Errorf("esquema EnvioDTE no cargado")
	}

	// Parsear XML para validación
	doc, err := libxml2.ParseString(string(envioDTE))
	if err != nil {
		return fmt.Errorf("error parseando XML del EnvioDTE: %w", err)
	}
	defer doc.Free()

	// Validar contra esquema
	if err := schema.Validate(doc); err != nil {
		return fmt.Errorf("error de validación XSD: %w", err)
	}

	return nil
}

// ValidarXMLContraEsquema valida un XML contra un esquema específico
func (s *ValidatorService) ValidarXMLContraEsquema(xml []byte, nombreEsquema string) error {
	s.schemasLock.RLock()
	schema, exists := s.schemas[nombreEsquema]
	s.schemasLock.RUnlock()

	if !exists {
		return fmt.Errorf("esquema %s no cargado", nombreEsquema)
	}

	// Parsear XML para validación
	doc, err := libxml2.ParseString(string(xml))
	if err != nil {
		return fmt.Errorf("error parseando XML: %w", err)
	}
	defer doc.Free()

	// Validar contra esquema
	if err := schema.Validate(doc); err != nil {
		return fmt.Errorf("error de validación XSD: %w", err)
	}

	return nil
}
