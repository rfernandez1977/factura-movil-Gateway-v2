package utils

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"FMgo/models"
)

// StorageUtils maneja el almacenamiento de documentos
type StorageUtils struct {
	utils   *DocumentUtils
	sii     *SIIUtils
	pdf     *PDFUtils
	baseDir string
}

// NewStorageUtils crea una nueva instancia de StorageUtils
func NewStorageUtils(baseDir string) *StorageUtils {
	return &StorageUtils{
		utils:   NewDocumentUtils(),
		sii:     NewSIIUtils(),
		pdf:     NewPDFUtils(),
		baseDir: baseDir,
	}
}

// SaveDocument guarda un documento en el almacenamiento
func (s *StorageUtils) SaveDocument(doc *models.DocumentoTributario) error {
	// Crear directorio para el emisor
	emisorDir := filepath.Join(s.baseDir, doc.RUTEmisor)
	if err := os.MkdirAll(emisorDir, 0755); err != nil {
		return fmt.Errorf("error al crear directorio del emisor: %v", err)
	}

	// Crear directorio para el año
	yearDir := filepath.Join(emisorDir, doc.FechaEmision.Format("2006"))
	if err := os.MkdirAll(yearDir, 0755); err != nil {
		return fmt.Errorf("error al crear directorio del año: %v", err)
	}

	// Crear directorio para el mes
	monthDir := filepath.Join(yearDir, doc.FechaEmision.Format("01"))
	if err := os.MkdirAll(monthDir, 0755); err != nil {
		return fmt.Errorf("error al crear directorio del mes: %v", err)
	}

	// Generar nombre del archivo
	filename := fmt.Sprintf("%s-%d-%s", doc.TipoDTE, doc.Folio, doc.FechaEmision.Format("20060102"))

	// Guardar XML
	xmlData, err := s.sii.GenerateSIIXML(doc)
	if err != nil {
		return fmt.Errorf("error al generar XML: %v", err)
	}

	if err := s.saveFile(monthDir, filename+".xml", xmlData); err != nil {
		return err
	}

	// Guardar PDF
	pdfData, err := s.pdf.GeneratePDF(doc)
	if err != nil {
		return fmt.Errorf("error al generar PDF: %v", err)
	}

	if err := s.saveFile(monthDir, filename+".pdf", pdfData); err != nil {
		return err
	}

	return nil
}

// GetDocument obtiene un documento del almacenamiento
func (s *StorageUtils) GetDocument(tipoDocumento string, folio int, rutEmisor string, fechaEmision time.Time) (*models.DocumentoTributario, error) {
	// Construir ruta del archivo XML
	xmlPath := filepath.Join(
		s.baseDir,
		rutEmisor,
		fechaEmision.Format("2006"),
		fechaEmision.Format("01"),
		fmt.Sprintf("%s-%d-%s.xml", tipoDocumento, folio, fechaEmision.Format("20060102")),
	)

	// Leer archivo XML
	xmlData, err := os.ReadFile(xmlPath)
	if err != nil {
		return nil, fmt.Errorf("error al leer archivo XML: %v", err)
	}

	// Decodificar XML
	var doc models.DocumentoTributario
	if err := xml.Unmarshal(xmlData, &doc); err != nil {
		return nil, fmt.Errorf("error al decodificar XML: %v", err)
	}

	return &doc, nil
}

// GetDocumentPDF obtiene el PDF de un documento del almacenamiento
func (s *StorageUtils) GetDocumentPDF(tipoDocumento string, folio int, rutEmisor string, fechaEmision time.Time) ([]byte, error) {
	// Construir ruta del archivo PDF
	pdfPath := filepath.Join(
		s.baseDir,
		rutEmisor,
		fechaEmision.Format("2006"),
		fechaEmision.Format("01"),
		fmt.Sprintf("%s-%d-%s.pdf", tipoDocumento, folio, fechaEmision.Format("20060102")),
	)

	// Leer archivo PDF
	return os.ReadFile(pdfPath)
}

// GetDocumentXML obtiene el XML de un documento del almacenamiento
func (s *StorageUtils) GetDocumentXML(tipoDocumento string, folio int, rutEmisor string, fechaEmision time.Time) ([]byte, error) {
	// Construir ruta del archivo XML
	xmlPath := filepath.Join(
		s.baseDir,
		rutEmisor,
		fechaEmision.Format("2006"),
		fechaEmision.Format("01"),
		fmt.Sprintf("%s-%d-%s.xml", tipoDocumento, folio, fechaEmision.Format("20060102")),
	)

	// Leer archivo XML
	return os.ReadFile(xmlPath)
}

// ListDocuments lista los documentos de un contribuyente
func (s *StorageUtils) ListDocuments(rutEmisor string, fechaInicio time.Time, fechaFin time.Time) ([]*models.DocumentoTributario, error) {
	var docs []*models.DocumentoTributario

	// Recorrer directorios de años
	for year := fechaInicio.Year(); year <= fechaFin.Year(); year++ {
		yearDir := filepath.Join(s.baseDir, rutEmisor, fmt.Sprintf("%d", year))
		if _, err := os.Stat(yearDir); os.IsNotExist(err) {
			continue
		}

		// Recorrer directorios de meses
		for month := 1; month <= 12; month++ {
			if year == fechaInicio.Year() && month < int(fechaInicio.Month()) {
				continue
			}
			if year == fechaFin.Year() && month > int(fechaFin.Month()) {
				continue
			}

			monthDir := filepath.Join(yearDir, fmt.Sprintf("%02d", month))
			if _, err := os.Stat(monthDir); os.IsNotExist(err) {
				continue
			}

			// Recorrer archivos XML
			err := filepath.Walk(monthDir, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				if !info.IsDir() && filepath.Ext(path) == ".xml" {
					// Leer archivo XML
					xmlData, err := os.ReadFile(path)
					if err != nil {
						return err
					}

					// Decodificar XML
					var doc models.DocumentoTributario
					if err := xml.Unmarshal(xmlData, &doc); err != nil {
						return err
					}

					// Verificar fecha
					if doc.FechaEmision.After(fechaInicio) && doc.FechaEmision.Before(fechaFin) {
						docs = append(docs, &doc)
					}
				}

				return nil
			})

			if err != nil {
				return nil, fmt.Errorf("error al listar documentos: %v", err)
			}
		}
	}

	return docs, nil
}

// DeleteDocument elimina un documento del almacenamiento
func (s *StorageUtils) DeleteDocument(tipoDocumento string, folio int, rutEmisor string, fechaEmision time.Time) error {
	// Construir rutas de archivos
	baseDir := filepath.Join(
		s.baseDir,
		rutEmisor,
		fechaEmision.Format("2006"),
		fechaEmision.Format("01"),
	)
	baseFilename := fmt.Sprintf("%s-%d-%s", tipoDocumento, folio, fechaEmision.Format("20060102"))

	// Eliminar XML
	xmlPath := filepath.Join(baseDir, baseFilename+".xml")
	if err := os.Remove(xmlPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("error al eliminar XML: %v", err)
	}

	// Eliminar PDF
	pdfPath := filepath.Join(baseDir, baseFilename+".pdf")
	if err := os.Remove(pdfPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("error al eliminar PDF: %v", err)
	}

	return nil
}

// saveFile guarda un archivo en el sistema de archivos
func (s *StorageUtils) saveFile(dir string, filename string, data []byte) error {
	// Crear archivo
	file, err := os.Create(filepath.Join(dir, filename))
	if err != nil {
		return fmt.Errorf("error al crear archivo: %v", err)
	}
	defer file.Close()

	// Escribir datos
	if _, err := io.Copy(file, bytes.NewReader(data)); err != nil {
		return fmt.Errorf("error al escribir archivo: %v", err)
	}

	return nil
}

// GetDocumentsByPeriod obtiene los documentos de un período
func (s *StorageUtils) GetDocumentsByPeriod(rutEmisor string, fechaInicio time.Time, fechaFin time.Time) ([]*models.DocumentoTributario, error) {
	return s.ListDocuments(rutEmisor, fechaInicio, fechaFin)
}

// GetDocumentStatus obtiene el estado de un documento
func (s *StorageUtils) GetDocumentStatus(tipoDocumento string, folio int, rutEmisor string, fechaEmision time.Time) (*models.EstadoDocumento, error) {
	// TODO: Implementar obtención de estado
	return nil, nil
}

// UpdateDocumentStatus actualiza el estado de un documento
func (s *StorageUtils) UpdateDocumentStatus(tipoDocumento string, folio int, rutEmisor string, fechaEmision time.Time, estado *models.EstadoDocumento) error {
	// TODO: Implementar actualización de estado
	return nil
}
