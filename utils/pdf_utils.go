package utils

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"image"
	"io"

	"github.com/cursor/FMgo/models"
	"github.com/jung-kurt/gofpdf"
	"github.com/skip2/go-qrcode"
)

// PDFUtils contiene utilidades para la generación de PDFs y firma digital
type PDFUtils struct {
	utils *DocumentUtils
	sii   *SIIUtils
}

// NewPDFUtils crea una nueva instancia de PDFUtils
func NewPDFUtils() *PDFUtils {
	return &PDFUtils{
		utils: NewDocumentUtils(),
		sii:   NewSIIUtils(),
	}
}

// GeneratePDF genera un PDF moderno para un documento tributario
func (p *PDFUtils) GeneratePDF(doc *models.DocumentoTributario) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// COLORES Y ESTILO
	pdf.SetFillColor(255, 251, 240) // fondo crema
	pdf.Rect(0, 0, 210, 297, "F")
	pdf.SetFillColor(0, 255, 204) // fluor verde-azul

	// LOGO (ejemplo: círculo fluor)
	pdf.SetFillColor(0, 255, 204)
	pdf.Ellipse(20, 20, 10, 10, 0, "F")
	pdf.SetFont("Arial", "B", 20)
	pdf.SetTextColor(0, 102, 204)
	pdf.SetXY(35, 15)
	pdf.Cell(60, 12, "Business Name")

	// DATOS DE LA EMPRESA
	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(60, 60, 60)
	pdf.SetXY(20, 30)
	pdf.Cell(0, 6, "RUT: "+doc.RUTEmisor)
	pdf.Ln(5)
	pdf.Cell(0, 6, "Dirección: [Dirección de la empresa]")
	pdf.Ln(5)
	pdf.Cell(0, 6, "Contacto: [Teléfono/Email]")

	// DATOS DEL DOCUMENTO
	pdf.SetFont("Arial", "B", 16)
	pdf.SetTextColor(0, 102, 204)
	pdf.SetXY(140, 15)
	pdf.Cell(60, 10, doc.TipoDTE+" N° "+fmt.Sprint(doc.Folio))
	pdf.SetFont("Arial", "", 11)
	pdf.SetTextColor(60, 60, 60)
	pdf.SetXY(140, 25)
	pdf.Cell(60, 8, "Emisión: "+doc.FechaEmision.Format("02/01/2006"))

	// DATOS DEL CLIENTE
	pdf.SetXY(20, 45)
	pdf.SetFont("Arial", "B", 12)
	pdf.SetTextColor(0, 102, 204)
	pdf.Cell(0, 8, "Cliente")
	pdf.Ln(7)
	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(60, 60, 60)
	pdf.Cell(0, 6, "RUT: "+doc.RUTReceptor)
	pdf.Ln(5)
	pdf.Cell(0, 6, "Razón Social: "+doc.RUTReceptor)

	// TABLA DE ÍTEMS
	pdf.SetXY(20, 70)
	pdf.SetFont("Arial", "B", 11)
	pdf.SetFillColor(0, 255, 204)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(20, 8, "Código", "1", 0, "C", true, 0, "")
	pdf.CellFormat(60, 8, "Descripción", "1", 0, "C", true, 0, "")
	pdf.CellFormat(20, 8, "Cantidad", "1", 0, "C", true, 0, "")
	pdf.CellFormat(30, 8, "Precio", "1", 0, "C", true, 0, "")
	pdf.CellFormat(30, 8, "Total", "1", 1, "C", true, 0, "")
	pdf.SetFont("Arial", "", 10)
	for _, detalle := range doc.Detalles {
		pdf.CellFormat(20, 8, "CÓDIGO", "1", 0, "C", false, 0, "")
		pdf.CellFormat(60, 8, detalle.Descripcion, "1", 0, "L", false, 0, "")
		pdf.CellFormat(20, 8, fmt.Sprintf("%d", detalle.Cantidad), "1", 0, "C", false, 0, "")
		pdf.CellFormat(30, 8, fmt.Sprintf("$%.2f", detalle.PrecioUnitario), "1", 0, "R", false, 0, "")
		pdf.CellFormat(30, 8, fmt.Sprintf("$%.2f", detalle.MontoItem), "1", 1, "R", false, 0, "")
	}

	// TOTALES
	pdf.Ln(5)
	pdf.SetX(100)
	pdf.SetFont("Arial", "B", 11)
	pdf.SetFillColor(255, 251, 240)
	pdf.CellFormat(50, 8, "Neto:", "0", 0, "R", false, 0, "")
	pdf.CellFormat(30, 8, fmt.Sprintf("$%.2f", doc.MontoNeto), "0", 1, "R", false, 0, "")
	pdf.SetX(100)
	pdf.CellFormat(50, 8, "IVA:", "0", 0, "R", false, 0, "")
	pdf.CellFormat(30, 8, fmt.Sprintf("$%.2f", doc.MontoIVA), "0", 1, "R", false, 0, "")
	pdf.SetX(100)
	pdf.SetFillColor(0, 255, 204)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(50, 10, "TOTAL:", "1", 0, "R", true, 0, "")
	pdf.CellFormat(30, 10, fmt.Sprintf("$%.2f", doc.MontoTotal), "1", 1, "R", true, 0, "")

	// QR CODE (ejemplo: folio y total)
	qrText := fmt.Sprintf("Folio: %d\nTotal: $%.2f", doc.Folio, doc.MontoTotal)
	qr, _ := qrcode.New(qrText, qrcode.Medium)
	qrImg := qr.Image(80)

	// Convertir image.Image a Reader
	qrImageReader := newQRImageReader(qrImg)

	imgOpt := gofpdf.ImageOptions{ImageType: "png", ReadDpi: false}
	pdf.RegisterImageOptionsReader("qr.png", imgOpt, qrImageReader)
	pdf.ImageOptions("qr.png", 160, 230, 30, 30, false, imgOpt, 0, "")

	// PIE DE PÁGINA
	pdf.SetY(270)
	pdf.SetFont("Arial", "I", 8)
	pdf.SetTextColor(120, 120, 120)
	pdf.Cell(0, 6, "Documento generado electrónicamente. Verifique en www.sii.cl")

	// GENERAR PDF EN MEMORIA
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, fmt.Errorf("error al generar PDF: %v", err)
	}

	return buf.Bytes(), nil
}

// QRImageReader implementa io.Reader para image.Image
type QRImageReader struct {
	img     image.Image
	encoded []byte
	pos     int
}

func newQRImageReader(img image.Image) *QRImageReader {
	// Codificar la imagen (simulación simple)
	return &QRImageReader{
		img:     img,
		encoded: []byte("simulated_qr_image_data"),
		pos:     0,
	}
}

func (r *QRImageReader) Read(p []byte) (n int, err error) {
	if r.pos >= len(r.encoded) {
		return 0, io.EOF
	}
	n = copy(p, r.encoded[r.pos:])
	r.pos += n
	return n, nil
}

// SignPDF firma digitalmente un PDF
func (p *PDFUtils) SignPDF(pdfData []byte, cert *x509.Certificate, key *rsa.PrivateKey) ([]byte, error) {
	// TODO: Implementar firma digital de PDF
	return pdfData, nil
}

// ValidateSignature valida una firma digital
func (p *PDFUtils) ValidateSignature(pdfData []byte, cert *x509.Certificate) error {
	// TODO: Implementar validación de firma digital
	return nil
}

// LoadCertificate carga un certificado desde un archivo PEM
func (p *PDFUtils) LoadCertificate(certPEM []byte) (*x509.Certificate, error) {
	block, _ := pem.Decode(certPEM)
	if block == nil {
		return nil, fmt.Errorf("no se pudo decodificar el certificado PEM")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("error al parsear certificado: %v", err)
	}

	return cert, nil
}

// LoadPrivateKey carga una llave privada desde un archivo PEM
func (p *PDFUtils) LoadPrivateKey(keyPEM []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(keyPEM)
	if block == nil {
		return nil, fmt.Errorf("no se pudo decodificar la llave privada PEM")
	}

	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("error al parsear llave privada: %v", err)
	}

	return key, nil
}
