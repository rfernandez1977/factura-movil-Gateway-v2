package services

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"FMgo/core/sii/logger"

	"software.sslmate.com/src/go-pkcs12"
)

// FirmaService proporciona funcionalidades para firmar documentos XML
type FirmaService struct {
	certPath   string
	keyPath    string
	password   string
	rutEmpresa string
	xmlProc    *XMLProcessor
	certCache  *CertCache
	log        *logger.Logger
}

// NewFirmaService crea una nueva instancia de FirmaService
func NewFirmaService(certPath, keyPath, password, rutEmpresa string) (*FirmaService, error) {
	// Crear logger
	log, err := logger.NewLogger("logs/firma_service.log", logger.DEBUG)
	if err != nil {
		return nil, fmt.Errorf("error inicializando logger: %w", err)
	}

	return &FirmaService{
		certPath:   certPath,
		keyPath:    keyPath,
		password:   password,
		rutEmpresa: rutEmpresa,
		xmlProc:    NewXMLProcessor(log),
		certCache:  NewCertCache(24*time.Hour, 100), // Caché por 24 horas, máximo 100 items
		log:        log,
	}, nil
}

// FirmarXML firma un documento XML usando el certificado digital
func (s *FirmaService) FirmarXML(xmlData []byte) ([]byte, error) {
	s.log.Debug("Iniciando proceso de firma XML")

	// Cargar el certificado P12
	p12Data, err := ioutil.ReadFile(s.certPath)
	if err != nil {
		s.log.Error("Error al leer certificado P12: %v", err)
		return nil, fmt.Errorf("error al leer certificado P12: %w", err)
	}

	// Extraer la clave privada y el certificado
	privateKey, cert, err := pkcs12.Decode(p12Data, s.password)
	if err != nil {
		s.log.Error("Error al decodificar P12: %v", err)
		return nil, fmt.Errorf("error al decodificar P12: %w", err)
	}

	s.log.Info("Certificado cargado exitosamente: %s", cert.Subject.CommonName)

	// Verificar que el certificado sea válido
	if err := cert.CheckSignature(x509.SHA256WithRSA, xmlData, nil); err != nil {
		s.log.Error("Certificado inválido: %v", err)
		return nil, fmt.Errorf("certificado inválido: %w", err)
	}

	// Calcular el hash SHA256 del documento
	hashed := sha256.Sum256(xmlData)

	// Firmar el hash con la clave privada
	rsaPrivateKey, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		s.log.Error("Tipo de clave privada no soportado")
		return nil, fmt.Errorf("tipo de clave privada no soportado")
	}

	signature, err := rsa.SignPKCS1v15(rand.Reader, rsaPrivateKey, crypto.SHA256, hashed[:])
	if err != nil {
		s.log.Error("Error al firmar documento: %v", err)
		return nil, fmt.Errorf("error al firmar documento: %w", err)
	}

	// Codificar la firma en base64
	signatureBase64 := base64.StdEncoding.EncodeToString(signature)

	// Insertar la firma en el documento XML
	signedXML, err := s.insertarFirmaXML(xmlData, signatureBase64, cert)
	if err != nil {
		s.log.Error("Error al insertar firma en XML: %v", err)
		return nil, fmt.Errorf("error al insertar firma en XML: %w", err)
	}

	s.log.Info("Documento XML firmado exitosamente")
	return signedXML, nil
}

// insertarFirmaXML inserta la firma digital en el documento XML
func (s *FirmaService) insertarFirmaXML(xmlData []byte, firma string, cert *x509.Certificate) ([]byte, error) {
	// Crear estructura de firma
	firmaNode := fmt.Sprintf(`
		<Signature xmlns="http://www.w3.org/2000/09/xmldsig#">
			<SignedInfo>
				<CanonicalizationMethod Algorithm="http://www.w3.org/TR/2001/REC-xml-c14n-20010315"/>
				<SignatureMethod Algorithm="http://www.w3.org/2000/09/xmldsig#rsa-sha1"/>
				<Reference URI="#DOC001">
					<Transforms>
						<Transform Algorithm="http://www.w3.org/2000/09/xmldsig#enveloped-signature"/>
					</Transforms>
					<DigestMethod Algorithm="http://www.w3.org/2000/09/xmldsig#sha1"/>
					<DigestValue>%s</DigestValue>
				</Reference>
			</SignedInfo>
			<SignatureValue>%s</SignatureValue>
			<KeyInfo>
				<X509Data>
					<X509Certificate>%s</X509Certificate>
				</X509Data>
			</KeyInfo>
		</Signature>
	`, s.calcularDigestValue(xmlData), firma, base64.StdEncoding.EncodeToString(cert.Raw))

	// Insertar firma antes del cierre del documento
	docStr := string(xmlData)
	closeTag := "</Documento>"
	pos := strings.LastIndex(docStr, closeTag)
	if pos == -1 {
		return nil, fmt.Errorf("no se encontró la etiqueta de cierre del documento")
	}

	signedXML := docStr[:pos] + firmaNode + docStr[pos:]
	return []byte(signedXML), nil
}

// calcularDigestValue calcula el valor digest del documento
func (s *FirmaService) calcularDigestValue(xmlData []byte) string {
	hash := sha1.Sum(xmlData)
	return base64.StdEncoding.EncodeToString(hash[:])
}

// ValidarFirma valida la firma de un documento XML
func (s *FirmaService) ValidarFirma(xmlData []byte) error {
	s.log.Debug("Iniciando validación de firma")

	// Extraer el certificado del documento
	certB64, err := s.extraerCertificadoDesdeXML(xmlData)
	if err != nil {
		s.log.Error("Error al extraer certificado: %v", err)
		return fmt.Errorf("error al extraer certificado: %w", err)
	}

	// Decodificar el certificado
	certDer, err := base64.StdEncoding.DecodeString(certB64)
	if err != nil {
		s.log.Error("Error al decodificar certificado: %v", err)
		return fmt.Errorf("error al decodificar certificado: %w", err)
	}

	// Parsear el certificado
	cert, err := x509.ParseCertificate(certDer)
	if err != nil {
		s.log.Error("Error al parsear certificado: %v", err)
		return fmt.Errorf("error al parsear certificado: %w", err)
	}

	s.log.Info("Certificado extraído y parseado exitosamente: %s", cert.Subject.CommonName)

	// Extraer la firma
	firma, err := s.extraerFirmaDesdeXML(xmlData)
	if err != nil {
		s.log.Error("Error al extraer firma: %v", err)
		return fmt.Errorf("error al extraer firma: %w", err)
	}

	// Decodificar la firma
	firmaBytes, err := base64.StdEncoding.DecodeString(firma)
	if err != nil {
		s.log.Error("Error al decodificar firma: %v", err)
		return fmt.Errorf("error al decodificar firma: %w", err)
	}

	// Calcular el hash del documento original
	docHash := s.calcularDigestValue(xmlData)

	// Verificar la firma
	err = rsa.VerifyPKCS1v15(cert.PublicKey.(*rsa.PublicKey), crypto.SHA1, []byte(docHash), firmaBytes)
	if err != nil {
		s.log.Error("Firma inválida: %v", err)
		return fmt.Errorf("firma inválida: %w", err)
	}

	s.log.Info("Firma validada exitosamente")
	return nil
}

// extraerCertificadoDesdeXML extrae el certificado X509 del documento XML
func (s *FirmaService) extraerCertificadoDesdeXML(xmlData []byte) (string, error) {
	// Limpiar XML antes de procesar
	cleanXML := s.xmlProc.limpiarXML(xmlData)

	// Validar estructura del documento
	if err := s.xmlProc.validarEstructuraXML(cleanXML); err != nil {
		return "", fmt.Errorf("documento XML inválido: %w", err)
	}

	return s.xmlProc.extraerCertificado(cleanXML)
}

// extraerFirmaDesdeXML extrae la firma del documento XML
func (s *FirmaService) extraerFirmaDesdeXML(xmlData []byte) (string, error) {
	// Limpiar XML antes de procesar
	cleanXML := s.xmlProc.limpiarXML(xmlData)

	// Validar estructura del documento
	if err := s.xmlProc.validarEstructuraXML(cleanXML); err != nil {
		return "", fmt.Errorf("documento XML inválido: %w", err)
	}

	return s.xmlProc.extraerFirma(cleanXML)
}

// ObtenerCertificado obtiene información del certificado digital
func (s *FirmaService) ObtenerCertificado() (*x509.Certificate, error) {
	p12Data, err := ioutil.ReadFile(s.certPath)
	if err != nil {
		return nil, fmt.Errorf("error al leer certificado P12: %w", err)
	}

	_, cert, err := pkcs12.Decode(p12Data, s.password)
	if err != nil {
		return nil, fmt.Errorf("error al decodificar P12: %w", err)
	}

	return cert, nil
}
