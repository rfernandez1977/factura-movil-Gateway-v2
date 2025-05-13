package services

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/jung-kurt/gofpdf"
	"github.com/cursor/FMgo/config"
	"github.com/cursor/FMgo/models"
)

// PDFService maneja la generación de PDFs
type PDFService struct {
	config       *config.SupabaseConfig
	templatePath string
	outputPath   string
}

// NewPDFService crea una nueva instancia del servicio de PDF
func NewPDFService(config *config.SupabaseConfig, templatePath, outputPath string) *PDFService {
	return &PDFService{
		config:       config,
		templatePath: templatePath,
		outputPath:   outputPath,
	}
}

// GenerarPDF genera un PDF para una factura
func (s *PDFService) GenerarPDF(factura *models.Factura) (string, error) {
	// Crear PDF
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Configurar fuente
	pdf.SetFont("Arial", "B", 16)

	// Título
	pdf.Cell(190, 10, "FACTURA ELECTRÓNICA")
	pdf.Ln(20)

	// Información de la empresa
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(190, 10, fmt.Sprintf("RUT: %s", factura.Empresa.Rut))
	pdf.Ln(10)
	pdf.Cell(190, 10, fmt.Sprintf("Razón Social: %s", factura.Empresa.RazonSocial))
	pdf.Ln(10)
	pdf.Cell(190, 10, fmt.Sprintf("Dirección: %s", factura.Empresa.Direccion))
	pdf.Ln(20)

	// Información del cliente
	pdf.Cell(190, 10, fmt.Sprintf("Cliente: %s", factura.Cliente.RazonSocial))
	pdf.Ln(10)
	pdf.Cell(190, 10, fmt.Sprintf("RUT: %s", factura.Cliente.Rut))
	pdf.Ln(10)
	pdf.Cell(190, 10, fmt.Sprintf("Dirección: %s", factura.Cliente.Direccion))
	pdf.Ln(20)

	// Detalles de la factura
	pdf.Cell(190, 10, fmt.Sprintf("Número: %d", factura.Numero))
	pdf.Ln(10)
	pdf.Cell(190, 10, fmt.Sprintf("Fecha: %s", factura.FechaEmision.Format("02/01/2006")))
	pdf.Ln(20)

	// Tabla de items
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 10, "Código")
	pdf.Cell(80, 10, "Descripción")
	pdf.Cell(35, 10, "Cantidad")
	pdf.Cell(35, 10, "Precio")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 12)
	for _, item := range factura.Items {
		pdf.Cell(40, 10, item.Codigo)
		pdf.Cell(80, 10, item.Descripcion)
		pdf.Cell(35, 10, fmt.Sprintf("%d", item.Cantidad))
		pdf.Cell(35, 10, fmt.Sprintf("$%d", item.PrecioUnitario))
		pdf.Ln(10)
	}

	// Totales
	pdf.Ln(10)
	pdf.Cell(120, 10, "")
	pdf.Cell(35, 10, "Neto:")
	pdf.Cell(35, 10, fmt.Sprintf("$%d", factura.TotalNeto))
	pdf.Ln(10)
	pdf.Cell(120, 10, "")
	pdf.Cell(35, 10, "IVA:")
	pdf.Cell(35, 10, fmt.Sprintf("$%d", factura.IVA))
	pdf.Ln(10)
	pdf.Cell(120, 10, "")
	pdf.Cell(35, 10, "Total:")
	pdf.Cell(35, 10, fmt.Sprintf("$%d", factura.Total))

	// Guardar PDF
	outputFile := filepath.Join(s.outputPath, fmt.Sprintf("factura_%d.pdf", factura.Numero))
	err := pdf.OutputFileAndClose(outputFile)
	if err != nil {
		return "", fmt.Errorf("error al guardar PDF: %v", err)
	}

	// Guardar PDF en Supabase
	file, err := os.Open(outputFile)
	if err != nil {
		return "", fmt.Errorf("error al abrir PDF: %v", err)
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("error al leer PDF: %v", err)
	}

	_, err = s.config.Client.DB.From("documentos_pdf").
		Insert(map[string]interface{}{
			"documento_id":     factura.ID,
			"pdf_data":         fileBytes,
			"fecha_generacion": factura.FechaEmision,
		}).
		Execute()

	if err != nil {
		return "", fmt.Errorf("error al guardar PDF en Supabase: %v", err)
	}

	return outputFile, nil
}

// ObtenerPDF obtiene el PDF de un documento
func (s *PDFService) ObtenerPDF(documentoID string) ([]byte, error) {
	var pdfDoc struct {
		PDFData []byte `json:"pdf_data"`
	}

	err := s.config.Client.DB.From("documentos_pdf").
		Select("pdf_data").
		Eq("documento_id", documentoID).
		Single().
		Execute(&pdfDoc)

	if err != nil {
		return nil, fmt.Errorf("error al obtener PDF: %v", err)
	}

	return pdfDoc.PDFData, nil
}
