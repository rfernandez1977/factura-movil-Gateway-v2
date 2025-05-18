package pdf

import (
	"bytes"
	"fmt"

	"FMgo/models"
	"github.com/jung-kurt/gofpdf"
)

// TributarioPDF maneja la generación de PDFs para documentos tributarios
type TributarioPDF struct {
	pdf      *gofpdf.Fpdf
	metadata map[string]string
}

// generarTablaImpuestos genera la tabla de impuestos del documento
func (p *TributarioPDF) generarTablaImpuestos(doc interface{}) {
	var (
		montoNeto            float64
		montoIVA             float64
		montoTotal           float64
		montoExento          float64
		impuestosAdicionales []models.ImpuestoAdicional
	)

	switch d := doc.(type) {
	case *models.Factura:
		montoNeto = d.MontoNeto
		montoIVA = d.MontoIVA
		montoTotal = d.MontoTotal
		montoExento = d.MontoExento
		// Obtener impuestos adicionales de los items
		for _, item := range d.Items {
			for _, impAdicional := range item.ImpuestosAdicionales {
				impuestosAdicionales = append(impuestosAdicionales, models.ImpuestoAdicional{
					Codigo:        impAdicional.Codigo,
					Nombre:        impAdicional.Nombre,
					Porcentaje:    impAdicional.Porcentaje,
					Monto:         impAdicional.Monto,
					BaseImponible: impAdicional.BaseImponible,
				})
			}
		}
	case *models.Boleta:
		montoNeto = d.MontoNeto
		montoIVA = d.MontoIVA
		montoTotal = d.MontoTotal
		montoExento = d.MontoExento
		// Obtener impuestos adicionales de los items
		for _, item := range d.Items {
			for _, impAdicional := range item.ImpuestosAdicionales {
				impuestosAdicionales = append(impuestosAdicionales, models.ImpuestoAdicional{
					Codigo:        impAdicional.Codigo,
					Nombre:        impAdicional.Nombre,
					Porcentaje:    impAdicional.Porcentaje,
					Monto:         impAdicional.Monto,
					BaseImponible: impAdicional.BaseImponible,
				})
			}
		}
	case *models.NotaCredito:
		montoNeto = d.MontoNeto
		montoIVA = d.MontoIVA
		montoTotal = d.MontoTotal
		montoExento = d.MontoExento
		// Obtener impuestos adicionales de los items
		for _, item := range d.Items {
			for _, impAdicional := range item.ImpuestosAdicionales {
				impuestosAdicionales = append(impuestosAdicionales, models.ImpuestoAdicional{
					Codigo:        impAdicional.Codigo,
					Nombre:        impAdicional.Nombre,
					Porcentaje:    impAdicional.Porcentaje,
					Monto:         impAdicional.Monto,
					BaseImponible: impAdicional.BaseImponible,
				})
			}
		}
	case *models.NotaDebito:
		montoNeto = d.MontoNeto
		montoIVA = d.MontoIVA
		montoTotal = d.MontoTotal
		montoExento = d.MontoExento
		// Obtener impuestos adicionales de los items
		for _, item := range d.Items {
			for _, impAdicional := range item.ImpuestosAdicionales {
				impuestosAdicionales = append(impuestosAdicionales, models.ImpuestoAdicional{
					Codigo:        impAdicional.Codigo,
					Nombre:        impAdicional.Nombre,
					Porcentaje:    impAdicional.Porcentaje,
					Monto:         impAdicional.Monto,
					BaseImponible: impAdicional.BaseImponible,
				})
			}
		}
	case *models.GuiaDespacho:
		montoNeto = d.MontoNeto
		montoIVA = d.MontoIVA
		montoTotal = d.MontoTotal
		montoExento = d.MontoExento
		// Obtener impuestos adicionales de los items
		for _, item := range d.Items {
			for _, impAdicional := range item.ImpuestosAdicionales {
				impuestosAdicionales = append(impuestosAdicionales, models.ImpuestoAdicional{
					Codigo:        impAdicional.Codigo,
					Nombre:        impAdicional.Nombre,
					Porcentaje:    impAdicional.Porcentaje,
					Monto:         impAdicional.Monto,
					BaseImponible: impAdicional.BaseImponible,
				})
			}
		}
	}

	// Crear tabla de impuestos
	p.pdf.SetY(p.pdf.GetY() + 10)
	p.pdf.SetFont("Arial", "B", 10)
	p.pdf.Cell(40, 10, "Resumen de Impuestos")
	p.pdf.Ln(10)

	p.pdf.SetFont("Arial", "", 10)
	p.pdf.Cell(40, 10, "Monto Neto:")
	p.pdf.Cell(40, 10, fmt.Sprintf("$%.2f", montoNeto))
	p.pdf.Ln(10)

	if montoExento > 0 {
		p.pdf.Cell(40, 10, "Monto Exento:")
		p.pdf.Cell(40, 10, fmt.Sprintf("$%.2f", montoExento))
		p.pdf.Ln(10)
	}

	p.pdf.Cell(40, 10, "IVA (19%):")
	p.pdf.Cell(40, 10, fmt.Sprintf("$%.2f", montoIVA))
	p.pdf.Ln(10)

	// Impuestos adicionales
	if len(impuestosAdicionales) > 0 {
		p.pdf.SetFont("Arial", "B", 10)
		p.pdf.Cell(40, 10, "Impuestos Adicionales:")
		p.pdf.Ln(10)

		p.pdf.SetFont("Arial", "", 10)
		for _, impuesto := range impuestosAdicionales {
			p.pdf.Cell(40, 10, fmt.Sprintf("%s (%.2f%%):", impuesto.Nombre, impuesto.Porcentaje))
			p.pdf.Cell(40, 10, fmt.Sprintf("$%.2f", impuesto.Monto))
			p.pdf.Ln(10)
		}
	}

	p.pdf.SetFont("Arial", "B", 10)
	p.pdf.Cell(40, 10, "Total:")
	p.pdf.Cell(40, 10, fmt.Sprintf("$%.2f", montoTotal))
	p.pdf.Ln(10)
}

// NewTributarioPDF crea una nueva instancia de TributarioPDF
func NewTributarioPDF() *TributarioPDF {
	pdf := gofpdf.New("P", "mm", "Letter", "")
	pdf.AddPage()

	return &TributarioPDF{
		pdf:      pdf,
		metadata: make(map[string]string),
	}
}

// GeneratePDF genera un PDF para un documento tributario
func (p *TributarioPDF) GeneratePDF(doc interface{}) error {
	p.agregarEncabezado(doc)
	p.agregarDatosEmisorReceptor(doc)
	p.agregarItems(doc)
	p.generarTablaImpuestos(doc)
	p.agregarObservaciones(doc)
	p.agregarTimbreElectronico(doc)
	return nil
}

// SaveToFile guarda el PDF en un archivo
func (p *TributarioPDF) SaveToFile(path string) error {
	return p.pdf.OutputFileAndClose(path)
}

// GetPDFBytes devuelve el PDF como bytes
func (p *TributarioPDF) GetPDFBytes() ([]byte, error) {
	var buf bytes.Buffer
	err := p.pdf.Output(&buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// agregarEncabezado agrega el encabezado al PDF
func (p *TributarioPDF) agregarEncabezado(doc interface{}) {
	// Implementar según el tipo de documento
	p.pdf.SetFont("Arial", "B", 16)
	switch doc.(type) {
	case *models.Factura:
		p.pdf.Cell(0, 10, "FACTURA ELECTRÓNICA")
	case *models.Boleta:
		p.pdf.Cell(0, 10, "BOLETA ELECTRÓNICA")
	case *models.NotaCredito:
		p.pdf.Cell(0, 10, "NOTA DE CRÉDITO ELECTRÓNICA")
	case *models.NotaDebito:
		p.pdf.Cell(0, 10, "NOTA DE DÉBITO ELECTRÓNICA")
	case *models.GuiaDespacho:
		p.pdf.Cell(0, 10, "GUÍA DE DESPACHO ELECTRÓNICA")
	}
	p.pdf.Ln(10)
}

// agregarDatosEmisorReceptor agrega los datos del emisor y receptor al PDF
func (p *TributarioPDF) agregarDatosEmisorReceptor(doc interface{}) {
	var (
		rutEmisor           string
		razonSocialEmisor   string
		rutReceptor         string
		razonSocialReceptor string
	)

	switch d := doc.(type) {
	case *models.Factura:
		rutEmisor = d.RutEmisor
		razonSocialEmisor = d.RazonSocialEmisor
		rutReceptor = d.RutReceptor
		razonSocialReceptor = d.RazonSocialReceptor
	case *models.Boleta:
		rutEmisor = d.RutEmisor
		razonSocialEmisor = d.RazonSocialEmisor
		rutReceptor = d.RutReceptor
		razonSocialReceptor = d.RazonSocialReceptor
	case *models.NotaCredito:
		rutEmisor = d.RutEmisor
		razonSocialEmisor = d.RazonSocialEmisor
		rutReceptor = d.RutReceptor
		razonSocialReceptor = d.RazonSocialReceptor
	case *models.NotaDebito:
		rutEmisor = d.RutEmisor
		razonSocialEmisor = d.RazonSocialEmisor
		rutReceptor = d.RutReceptor
		razonSocialReceptor = d.RazonSocialReceptor
	case *models.GuiaDespacho:
		rutEmisor = d.RutEmisor
		razonSocialEmisor = d.RazonSocialEmisor
		rutReceptor = d.RutReceptor
		razonSocialReceptor = d.RazonSocialReceptor
	}

	p.pdf.SetFont("Arial", "B", 12)
	p.pdf.Cell(40, 10, "Emisor:")
	p.pdf.SetFont("Arial", "", 12)
	p.pdf.Cell(0, 10, fmt.Sprintf("%s - %s", rutEmisor, razonSocialEmisor))
	p.pdf.Ln(10)

	p.pdf.SetFont("Arial", "B", 12)
	p.pdf.Cell(40, 10, "Receptor:")
	p.pdf.SetFont("Arial", "", 12)
	p.pdf.Cell(0, 10, fmt.Sprintf("%s - %s", rutReceptor, razonSocialReceptor))
	p.pdf.Ln(15)
}

// agregarItems agrega los items al PDF
func (p *TributarioPDF) agregarItems(doc interface{}) {
	var items []models.Item

	switch d := doc.(type) {
	case *models.Factura:
		items = d.Items
	case *models.Boleta:
		items = d.Items
	case *models.NotaCredito:
		items = d.Items
	case *models.NotaDebito:
		items = d.Items
	case *models.GuiaDespacho:
		items = d.Items
	}

	// Cabecera de tabla
	p.pdf.SetFont("Arial", "B", 10)
	p.pdf.Cell(20, 10, "Código")
	p.pdf.Cell(70, 10, "Descripción")
	p.pdf.Cell(20, 10, "Cantidad")
	p.pdf.Cell(30, 10, "Precio Unit.")
	p.pdf.Cell(20, 10, "Descuento")
	p.pdf.Cell(30, 10, "Total")
	p.pdf.Ln(10)

	// Datos
	p.pdf.SetFont("Arial", "", 10)
	for _, item := range items {
		p.pdf.Cell(20, 10, item.Codigo)
		p.pdf.Cell(70, 10, item.Descripcion)
		p.pdf.Cell(20, 10, fmt.Sprintf("%.2f", item.Cantidad))
		p.pdf.Cell(30, 10, fmt.Sprintf("$%.2f", item.PrecioUnitario))
		p.pdf.Cell(20, 10, fmt.Sprintf("$%.2f", item.Descuento))
		p.pdf.Cell(30, 10, fmt.Sprintf("$%.2f", item.MontoItem))
		p.pdf.Ln(10)
	}

	p.pdf.Ln(10)
}

// agregarObservaciones agrega observaciones al PDF
func (p *TributarioPDF) agregarObservaciones(doc interface{}) {
	// Implementar según el tipo de documento (si es necesario)
}

// agregarTimbreElectronico agrega el timbre electrónico al PDF
func (p *TributarioPDF) agregarTimbreElectronico(doc interface{}) {
	// Implementar la lógica para agregar el timbre electrónico (si es necesario)
}
