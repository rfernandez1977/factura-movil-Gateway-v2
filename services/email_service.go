package services

import (
	"bytes"
	"fmt"
	"net/smtp"
	"strings"
	"time"

	"github.com/fmgo/config"
	"github.com/fmgo/models"
	"github.com/jordan-wright/email"
)

// EmailService maneja el envío de emails
type EmailService struct {
	config    *config.SupabaseConfig
	smtpHost  string
	smtpPort  int
	smtpUser  string
	smtpPass  string
	fromEmail string
	fromName  string
}

// NewEmailService crea una nueva instancia del servicio de email
func NewEmailService(
	config *config.SupabaseConfig,
	smtpHost string,
	smtpPort int,
	smtpUser string,
	smtpPass string,
	fromEmail string,
	fromName string,
) *EmailService {
	return &EmailService{
		config:    config,
		smtpHost:  smtpHost,
		smtpPort:  smtpPort,
		smtpUser:  smtpUser,
		smtpPass:  smtpPass,
		fromEmail: fromEmail,
		fromName:  fromName,
	}
}

// EnviarFactura envía una factura por email
func (s *EmailService) EnviarFactura(factura *models.Factura, pdfPath string) error {
	// Construir asunto
	asunto := fmt.Sprintf("Factura Electrónica N° %d", factura.Folio)

	// Construir cuerpo
	cuerpo := fmt.Sprintf(`
		Estimado/a %s,

		Adjunto encontrará la Factura Electrónica N° %d.

		Detalles de la factura:
		- Fecha: %s
		- Monto Neto: $%.2f
		- IVA: $%.2f
		- Total: $%.2f

		Saludos cordiales,
		%s
	`, factura.RazonSocialReceptor, factura.Folio, factura.FechaEmision.Format("02/01/2006"),
		factura.MontoNeto, factura.MontoIVA, factura.MontoTotal, factura.RazonSocialEmisor)

	// Enviar email
	// Nota: Estamos asumiendo que necesitaríamos una dirección de email del receptor
	// Ya que la estructura actual no tiene este campo, usaremos un email genérico
	destinatario := "receptor@example.com" // Deberíamos obtener este dato de algún lugar
	err := s.enviarEmail(destinatario, asunto, cuerpo, pdfPath)
	if err != nil {
		return fmt.Errorf("error al enviar email: %v", err)
	}

	// Guardar registro en Supabase
	// Si no tenemos acceso directo a la base de datos, podríamos omitir esta parte
	// o implementarla de otra manera
	/*
		_, err = s.config.Client.DB.From("emails_enviados").
			Insert(map[string]interface{}{
				"documento_id": factura.ID,
				"destinatario": destinatario,
				"asunto":       asunto,
				"cuerpo":       cuerpo,
				"fecha_envio":  factura.FechaEmision,
			}).
			Execute()

		if err != nil {
			return fmt.Errorf("error al guardar registro de email: %v", err)
		}
	*/

	return nil
}

// enviarEmail envía un email
func (s *EmailService) enviarEmail(to, subject, body, attachmentPath string) error {
	// Configurar autenticación
	auth := smtp.PlainAuth("", s.smtpUser, s.smtpPass, s.smtpHost)

	// Construir headers
	headers := make(map[string]string)
	headers["From"] = fmt.Sprintf("%s <%s>", s.fromName, s.fromEmail)
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/plain; charset=\"utf-8\""

	// Construir mensaje
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	// Enviar email
	err := smtp.SendMail(
		fmt.Sprintf("%s:%d", s.smtpHost, s.smtpPort),
		auth,
		s.fromEmail,
		[]string{to},
		[]byte(message),
	)

	if err != nil {
		return fmt.Errorf("error al enviar email: %v", err)
	}

	return nil
}

// ObtenerEmailsEnviados obtiene los emails enviados para un documento
/*
func (s *EmailService) ObtenerEmailsEnviados(documentoID string) ([]*models.EmailEnviado, error) {
	var emails []*models.EmailEnviado
	err := s.config.Client.DB.From("emails_enviados").
		Select("*").
		Eq("documento_id", documentoID).
		Execute(&emails)

	if err != nil {
		return nil, fmt.Errorf("error al obtener emails: %v", err)
	}

	return emails, nil
}
*/

// Enviar envía un email
func (e *EmailService) Enviar(destinatario, asunto, mensaje string) error {
	// Preparar mensaje
	msg := fmt.Sprintf("From: %s <%s>\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"Content-Type: text/plain; charset=UTF-8\r\n"+
		"\r\n"+
		"%s", e.fromEmail, e.fromName, destinatario, asunto, mensaje)

	// Autenticación
	auth := smtp.PlainAuth("", e.smtpUser, e.smtpPass, e.smtpHost)

	// Enviar email
	err := smtp.SendMail(
		fmt.Sprintf("%s:%d", e.smtpHost, e.smtpPort),
		auth,
		e.fromEmail,
		[]string{destinatario},
		[]byte(msg),
	)
	if err != nil {
		return fmt.Errorf("error enviando email: %w", err)
	}

	return nil
}

// EnviarDocumento envía un documento por correo electrónico
func (s *EmailService) EnviarDocumento(toEmail, toName string, doc interface{}, pdfData []byte, xmlData []byte) error {
	// Crear correo
	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", s.fromName, s.fromEmail)
	e.To = []string{fmt.Sprintf("%s <%s>", toName, toEmail)}
	e.Subject = s.generarAsunto(doc)

	// Generar cuerpo del correo
	body, err := s.generarCuerpo(doc)
	if err != nil {
		return fmt.Errorf("error generando cuerpo del correo: %w", err)
	}
	e.HTML = []byte(body)

	// Adjuntar PDF
	if pdfData != nil {
		nombrePDF := s.generarNombreArchivo(doc, "pdf")
		e.Attach(bytes.NewReader(pdfData), nombrePDF, "application/pdf")
	}

	// Adjuntar XML
	if xmlData != nil {
		nombreXML := s.generarNombreArchivo(doc, "xml")
		e.Attach(bytes.NewReader(xmlData), nombreXML, "application/xml")
	}

	// Enviar correo
	auth := smtp.PlainAuth("", s.smtpUser, s.smtpPass, s.smtpHost)
	if err := e.Send(fmt.Sprintf("%s:%d", s.smtpHost, s.smtpPort), auth); err != nil {
		return fmt.Errorf("error enviando correo: %w", err)
	}

	return nil
}

// generarAsunto genera el asunto del correo según el tipo de documento
func (s *EmailService) generarAsunto(doc interface{}) string {
	switch d := doc.(type) {
	case *models.Factura:
		return fmt.Sprintf("Factura Electrónica N° %d", d.Folio)
	case *models.Boleta:
		return fmt.Sprintf("Boleta Electrónica N° %d", d.Folio)
	default:
		return "Documento Tributario Electrónico"
	}
}

// generarCuerpo genera el cuerpo del correo según el tipo de documento
func (s *EmailService) generarCuerpo(doc interface{}) (string, error) {
	var tipo, folio string
	var fecha string
	var emisor, receptor string
	var montoTotal float64

	switch d := doc.(type) {
	case *models.Factura:
		tipo = "Factura Electrónica"
		folio = fmt.Sprintf("%d", d.Folio)
		fecha = d.FechaEmision.Format("02/01/2006")
		emisor = d.RazonSocialEmisor
		receptor = d.RazonSocialReceptor
		montoTotal = d.MontoTotal
	case *models.Boleta:
		tipo = "Boleta Electrónica"
		folio = fmt.Sprintf("%d", d.Folio)
		fecha = d.FechaEmision.Format("02/01/2006")
		emisor = d.RazonSocialEmisor
		receptor = d.RazonSocialReceptor
		montoTotal = d.MontoTotal
	default:
		// Para otros tipos de documentos, usar valores genéricos o extraer campos si es posible
		tipo = "Documento Tributario Electrónico"
		folio = "0"
		fecha = time.Now().Format("02/01/2006")
		emisor = "Emisor"
		receptor = "Receptor"
		montoTotal = 0.0
	}

	// Generar HTML del correo
	html := fmt.Sprintf(`
		<html>
			<body>
				<h1>%s</h1>
				<p>Estimado(a) %s,</p>
				<p>Adjunto encontrará su %s N° %s emitida el %s.</p>
				<p>Detalles del documento:</p>
				<ul>
					<li>Emisor: %s</li>
					<li>Receptor: %s</li>
					<li>Monto Total: $%.2f</li>
				</ul>
				<p>Este es un correo automático, por favor no responda.</p>
			</body>
		</html>
	`, tipo, receptor, tipo, folio, fecha, emisor, receptor, montoTotal)

	return html, nil
}

// generarNombreArchivo genera el nombre del archivo según el tipo de documento
func (s *EmailService) generarNombreArchivo(doc interface{}, extension string) string {
	var tipo, folio string

	switch d := doc.(type) {
	case *models.Factura:
		tipo = "FE"
		folio = fmt.Sprintf("%d", d.Folio)
	case *models.Boleta:
		tipo = "BE"
		folio = fmt.Sprintf("%d", d.Folio)
	default:
		tipo = "DTE"
		folio = "00000000"
	}

	return fmt.Sprintf("%s_%s.%s", tipo, folio, extension)
}

// EnviarNotificacion envía una notificación por correo electrónico
func (s *EmailService) EnviarNotificacion(toEmail, toName, asunto, mensaje string) error {
	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", s.fromName, s.fromEmail)
	e.To = []string{fmt.Sprintf("%s <%s>", toName, toEmail)}
	e.Subject = asunto
	e.HTML = []byte(mensaje)

	auth := smtp.PlainAuth("", s.smtpUser, s.smtpPass, s.smtpHost)
	if err := e.Send(fmt.Sprintf("%s:%d", s.smtpHost, s.smtpPort), auth); err != nil {
		return fmt.Errorf("error enviando notificación: %w", err)
	}

	return nil
}

// EnviarEmail envía un correo electrónico
func (s *EmailService) EnviarEmail(to, subject, template string, data map[string]interface{}) error {
	// Construir mensaje
	body := s.procesarPlantilla(template, data)
	msg := fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\n\n%s", s.fromEmail, to, subject, body)

	// Configurar autenticación
	auth := smtp.PlainAuth("", s.smtpUser, s.smtpPass, s.smtpHost)

	// Enviar correo
	err := smtp.SendMail(
		fmt.Sprintf("%s:%d", s.smtpHost, s.smtpPort),
		auth,
		s.fromEmail,
		[]string{to},
		[]byte(msg),
	)
	if err != nil {
		return fmt.Errorf("error enviando email: %v", err)
	}

	return nil
}

// procesarPlantilla procesa una plantilla de email con los datos proporcionados
func (s *EmailService) procesarPlantilla(template string, data map[string]interface{}) string {
	result := template
	for key, value := range data {
		placeholder := fmt.Sprintf("{{%s}}", key)
		result = strings.ReplaceAll(result, placeholder, fmt.Sprintf("%v", value))
	}
	return result
}
