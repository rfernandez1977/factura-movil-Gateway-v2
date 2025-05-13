package utils

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"mime/quotedprintable"
	"net/smtp"
	"net/textproto"
	"time"

	"github.com/cursor/FMgo/models"
)

// EmailUtils maneja el envío de correos electrónicos
type EmailUtils struct {
	utils    *DocumentUtils
	sii      *SIIUtils
	pdf      *PDFUtils
	storage  *StorageUtils
	smtpHost string
	smtpPort int
	username string
	password string
	from     string
}

// NewEmailUtils crea una nueva instancia de EmailUtils
func NewEmailUtils(smtpHost string, smtpPort int, username string, password string, from string) *EmailUtils {
	return &EmailUtils{
		utils:    NewDocumentUtils(),
		sii:      NewSIIUtils(),
		pdf:      NewPDFUtils(),
		storage:  NewStorageUtils("/var/documents"),
		smtpHost: smtpHost,
		smtpPort: smtpPort,
		username: username,
		password: password,
		from:     from,
	}
}

// SendDocumentEmail envía un documento por correo electrónico
func (e *EmailUtils) SendDocumentEmail(doc *models.DocumentoTributario, to []string, subject string, body string) error {
	// Obtener PDF del documento
	tipoDoc := string(doc.Tipo)

	pdfData, err := e.storage.GetDocumentPDF(tipoDoc, doc.Folio, doc.RutEmisor, doc.FechaEmision)
	if err != nil {
		return fmt.Errorf("error al obtener PDF: %v", err)
	}

	// Obtener XML del documento
	xmlData, err := e.storage.GetDocumentXML(tipoDoc, doc.Folio, doc.RutEmisor, doc.FechaEmision)
	if err != nil {
		return fmt.Errorf("error al obtener XML: %v", err)
	}

	// Crear mensaje multipart
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Agregar encabezados
	headers := make(textproto.MIMEHeader)
	headers.Set("From", e.from)
	headers.Set("To", to[0])
	headers.Set("Subject", subject)
	headers.Set("MIME-Version", "1.0")
	headers.Set("Content-Type", fmt.Sprintf("multipart/mixed; boundary=%s", writer.Boundary()))

	for key, values := range headers {
		for _, value := range values {
			fmt.Fprintf(&buf, "%s: %s\r\n", key, value)
		}
	}
	fmt.Fprintf(&buf, "\r\n")

	// Agregar cuerpo del mensaje
	partHeaders := make(textproto.MIMEHeader)
	partHeaders.Set("Content-Type", "text/plain; charset=utf-8")
	partHeaders.Set("Content-Transfer-Encoding", "quoted-printable")

	part, err := writer.CreatePart(partHeaders)
	if err != nil {
		return fmt.Errorf("error al crear parte del mensaje: %v", err)
	}

	qp := quotedprintable.NewWriter(part)
	if _, err := qp.Write([]byte(body)); err != nil {
		return fmt.Errorf("error al escribir cuerpo del mensaje: %v", err)
	}
	qp.Close()

	// Agregar PDF
	partHeaders = make(textproto.MIMEHeader)
	partHeaders.Set("Content-Type", "application/pdf")
	partHeaders.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s-%d.pdf", tipoDoc, doc.Folio))

	part, err = writer.CreatePart(partHeaders)
	if err != nil {
		return fmt.Errorf("error al crear parte del PDF: %v", err)
	}

	if _, err := part.Write(pdfData); err != nil {
		return fmt.Errorf("error al escribir PDF: %v", err)
	}

	// Agregar XML
	partHeaders = make(textproto.MIMEHeader)
	partHeaders.Set("Content-Type", "application/xml")
	partHeaders.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s-%d.xml", tipoDoc, doc.Folio))

	part, err = writer.CreatePart(partHeaders)
	if err != nil {
		return fmt.Errorf("error al crear parte del XML: %v", err)
	}

	if _, err := part.Write(xmlData); err != nil {
		return fmt.Errorf("error al escribir XML: %v", err)
	}

	writer.Close()

	// Enviar correo
	auth := smtp.PlainAuth("", e.username, e.password, e.smtpHost)
	if err := smtp.SendMail(fmt.Sprintf("%s:%d", e.smtpHost, e.smtpPort), auth, e.from, to, buf.Bytes()); err != nil {
		return fmt.Errorf("error al enviar correo: %v", err)
	}

	return nil
}

// SendDocumentNotification envía una notificación de documento
func (e *EmailUtils) SendDocumentNotification(doc *models.DocumentoTributario, to []string, notificationType string) error {
	var subject, body string
	tipoDoc := string(doc.Tipo)

	switch notificationType {
	case "emision":
		subject = fmt.Sprintf("Nuevo documento emitido: %s %d", tipoDoc, doc.Folio)
		body = fmt.Sprintf("Se ha emitido un nuevo documento:\n\n"+
			"Tipo: %s\n"+
			"Folio: %d\n"+
			"Emisor: %s\n"+
			"Receptor: %s\n"+
			"Fecha: %s\n"+
			"Monto Total: %.2f\n",
			tipoDoc, doc.Folio, doc.RutEmisor, doc.RutReceptor, doc.FechaEmision.Format("02/01/2006"), doc.MontoTotal)
	case "recepcion":
		subject = fmt.Sprintf("Nuevo documento recibido: %s %d", tipoDoc, doc.Folio)
		body = fmt.Sprintf("Se ha recibido un nuevo documento:\n\n"+
			"Tipo: %s\n"+
			"Folio: %d\n"+
			"Emisor: %s\n"+
			"Receptor: %s\n"+
			"Fecha: %s\n"+
			"Monto Total: %.2f\n",
			tipoDoc, doc.Folio, doc.RutEmisor, doc.RutReceptor, doc.FechaEmision.Format("02/01/2006"), doc.MontoTotal)
	default:
		return fmt.Errorf("tipo de notificación no válido: %s", notificationType)
	}

	return e.SendDocumentEmail(doc, to, subject, body)
}

// SendDocumentSummary envía un resumen de documentos
func (e *EmailUtils) SendDocumentSummary(rutEmisor string, fechaInicio time.Time, fechaFin time.Time, to []string) error {
	// Obtener documentos del período
	docs, err := e.storage.GetDocumentsByPeriod(rutEmisor, fechaInicio, fechaFin)
	if err != nil {
		return fmt.Errorf("error al obtener documentos: %v", err)
	}

	// Calcular totales
	var totalEmitidos, totalRecibidos, totalNeto, totalIVA float64
	for _, doc := range docs {
		if doc.RutEmisor == rutEmisor {
			totalEmitidos++
		} else {
			totalRecibidos++
		}
		totalNeto += doc.MontoNeto
		totalIVA += doc.MontoIVA
	}

	// Crear cuerpo del mensaje
	body := fmt.Sprintf("Resumen de documentos del período %s al %s:\n\n"+
		"Documentos emitidos: %.0f\n"+
		"Documentos recibidos: %.0f\n"+
		"Total neto: %.2f\n"+
		"Total IVA: %.2f\n"+
		"Total general: %.2f\n",
		fechaInicio.Format("02/01/2006"), fechaFin.Format("02/01/2006"),
		totalEmitidos, totalRecibidos, totalNeto, totalIVA, totalNeto+totalIVA)

	// Enviar correo
	subject := fmt.Sprintf("Resumen de documentos %s - %s", fechaInicio.Format("02/01/2006"), fechaFin.Format("02/01/2006"))
	return e.SendDocumentEmail(nil, to, subject, body)
}
