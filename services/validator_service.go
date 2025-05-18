package services

import (
	"encoding/xml"
	"fmt"
	"path/filepath"
	"sync"

	"FMgo/models"
	"github.com/lestrrat-go/libxml2"
	"github.com/lestrrat-go/libxml2/xsd"
)

// ValidatorService proporciona funcionalidades para validar documentos XML contra esquemas XSD
type ValidatorService struct {
	mu       sync.RWMutex
	schemas  map[string]*xsd.Schema
	basePath string
}

// NewValidatorService crea una nueva instancia de ValidatorService
func NewValidatorService(basePath string) *ValidatorService {
	return &ValidatorService{
		schemas:  make(map[string]*xsd.Schema),
		basePath: basePath,
	}
}

// CargarEsquema carga un esquema XSD específico
func (v *ValidatorService) CargarEsquema(nombre string) error {
	v.mu.Lock()
	defer v.mu.Unlock()

	if _, exists := v.schemas[nombre]; exists {
		return nil // Ya está cargado
	}

	rutaEsquema := filepath.Join(v.basePath, nombre)
	schema, err := xsd.ParseFromFile(rutaEsquema)
	if err != nil {
		return fmt.Errorf("error al cargar esquema %s: %w", nombre, err)
	}

	v.schemas[nombre] = schema
	return nil
}

// ValidarXML valida un documento XML contra un esquema específico
func (v *ValidatorService) ValidarXML(xml []byte, nombreEsquema string) error {
	v.mu.RLock()
	schema, exists := v.schemas[nombreEsquema]
	v.mu.RUnlock()

	if !exists {
		return fmt.Errorf("esquema %s no encontrado", nombreEsquema)
	}

	doc, err := libxml2.ParseString(string(xml))
	if err != nil {
		return fmt.Errorf("error al parsear XML: %w", err)
	}
	defer doc.Free()

	if err := schema.Validate(doc); err != nil {
		return fmt.Errorf("error de validación: %w", err)
	}

	return nil
}

// ValidarDocumento valida un documento tributario contra el esquema correspondiente
func (v *ValidatorService) ValidarDocumento(doc *models.DTEDocument) error {
	xmlBytes, err := xml.Marshal(doc)
	if err != nil {
		return fmt.Errorf("error al convertir documento a XML: %w", err)
	}

	// Seleccionar el esquema según el tipo de documento
	var esquema string
	switch doc.Documento.Encabezado.IdDoc.TipoDTE {
	case "33", "34", "56", "61":
		esquema = "DTE_v10.xsd"
	case "39", "41":
		esquema = "EnvioBOLETA_v11.xsd"
	default:
		return fmt.Errorf("tipo de documento no soportado: %s", doc.Documento.Encabezado.IdDoc.TipoDTE)
	}

	return v.ValidarXML(xmlBytes, esquema)
}

// ValidarEnvio valida un envío de documentos contra el esquema correspondiente
func (v *ValidatorService) ValidarEnvio(envio *models.EnvioDTEDocument) error {
	xmlBytes, err := xml.Marshal(envio)
	if err != nil {
		return fmt.Errorf("error al convertir envío a XML: %w", err)
	}

	return v.ValidarXML(xmlBytes, "EnvioDTE_v10.xsd")
}

// LimpiarEsquemas libera los esquemas cargados
func (v *ValidatorService) LimpiarEsquemas() {
	v.mu.Lock()
	defer v.mu.Unlock()

	for _, schema := range v.schemas {
		schema.Free()
	}
	v.schemas = make(map[string]*xsd.Schema)
}
