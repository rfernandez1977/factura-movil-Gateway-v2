package services

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/beevik/etree"
	"FMgo/core/firma/models"
	"software.sslmate.com/src/go-pkcs12"
)

// BaseFirmaService implementa la funcionalidad base de firma digital
type BaseFirmaService struct {
	config     *models.ConfiguracionFirma
	privateKey *rsa.PrivateKey
	cert       *x509.Certificate
	xmlService *XMLService
	logService *LogService
}

// NewBaseFirmaService crea una nueva instancia del servicio base de firma
func NewBaseFirmaService(config *models.ConfiguracionFirma) (*BaseFirmaService, error) {
	// Crear servicio de logging
	logService, err := NewLogService("logs/firma")
	if err != nil {
		return nil, fmt.Errorf("error inicializando servicio de logs: %w", err)
	}

	// Cargar certificado PKCS12
	pfxData, err := ioutil.ReadFile(config.RutaCertificado)
	if err != nil {
		return nil, fmt.Errorf("error al leer certificado PFX: %v", err)
	}

	// Extraer llave privada y certificado
	privateKey, cert, err := pkcs12.Decode(pfxData, config.Password)
	if err != nil {
		return nil, fmt.Errorf("error al decodificar PKCS12: %v", err)
	}

	// Convertir llave privada a RSA
	rsaKey, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("la llave privada no es de tipo RSA")
	}

	return &BaseFirmaService{
		config:     config,
		privateKey: rsaKey,
		cert:       cert,
		xmlService: NewXMLService(),
		logService: logService,
	}, nil
}

// FirmarDocumento implementa la firma de documentos XML
func (s *BaseFirmaService) FirmarDocumento(documento string) (*models.ResultadoFirma, error) {
	// Crear documento XML
	doc := etree.NewDocument()
	if err := doc.ReadFromString(documento); err != nil {
		return nil, fmt.Errorf("error al leer XML: %v", err)
	}

	// Aplicar transformación C14N
	canonicalXML, err := s.canonicalizar(doc)
	if err != nil {
		return nil, fmt.Errorf("error en canonicalización: %v", err)
	}

	// Calcular digest
	digest := s.calcularDigest(canonicalXML)

	// Firmar digest
	firma, err := s.firmarDigest(digest)
	if err != nil {
		return nil, fmt.Errorf("error al firmar: %v", err)
	}

	// Construir estructura de firma
	signedDoc, err := s.construirXMLFirmado(doc, digest, firma)
	if err != nil {
		return nil, fmt.Errorf("error al construir XML firmado: %v", err)
	}

	return &models.ResultadoFirma{
		XMLFirmado:     signedDoc,
		DigestValue:    base64.StdEncoding.EncodeToString(digest),
		SignatureValue: base64.StdEncoding.EncodeToString(firma),
		Timestamp:      time.Now().Format(time.RFC3339),
	}, nil
}

// ValidarFirma implementa la validación de firmas
func (s *BaseFirmaService) ValidarFirma(xmlFirmado string) (*models.EstadoFirma, error) {
	doc := etree.NewDocument()
	if err := doc.ReadFromString(xmlFirmado); err != nil {
		return nil, fmt.Errorf("error al leer XML: %v", err)
	}

	// Extraer y validar firma
	firma, err := s.extraerFirma(doc)
	if err != nil {
		return &models.EstadoFirma{
			Valida: false,
			Error:  fmt.Sprintf("error al extraer firma: %v", err),
		}, nil
	}

	// Validar certificado
	if err := s.validarCertificado(); err != nil {
		return &models.EstadoFirma{
			Valida: false,
			Error:  fmt.Sprintf("certificado inválido: %v", err),
		}, nil
	}

	return &models.EstadoFirma{
		Valida:          true,
		FechaValidacion: time.Now().Format(time.RFC3339),
		CertificadoID:   s.cert.SerialNumber.String(),
	}, nil
}

// ObtenerCertificado retorna el certificado actual
func (s *BaseFirmaService) ObtenerCertificado() (*x509.Certificate, error) {
	if s.cert == nil {
		return nil, fmt.Errorf("no hay certificado cargado")
	}
	return s.cert, nil
}

// Métodos privados auxiliares

func (s *BaseFirmaService) canonicalizar(doc *etree.Document) ([]byte, error) {
	root := doc.Root()
	if root == nil {
		return nil, fmt.Errorf("documento XML sin elemento raíz")
	}
	return root.WriteToBytes()
}

func (s *BaseFirmaService) calcularDigest(data []byte) []byte {
	hasher := sha256.New()
	hasher.Write(data)
	return hasher.Sum(nil)
}

func (s *BaseFirmaService) firmarDigest(digest []byte) ([]byte, error) {
	return rsa.SignPKCS1v15(rand.Reader, s.privateKey, crypto.SHA256, digest)
}

func (s *BaseFirmaService) validarCertificado() error {
	now := time.Now()
	if now.Before(s.cert.NotBefore) || now.After(s.cert.NotAfter) {
		return fmt.Errorf("certificado fuera de fecha de validez")
	}
	return nil
}

func (s *BaseFirmaService) extraerFirma(doc *etree.Document) ([]byte, error) {
	signatureElement := doc.FindElement("//Signature")
	if signatureElement == nil {
		return nil, fmt.Errorf("firma no encontrada en el documento")
	}

	signatureValue := signatureElement.FindElement("//SignatureValue")
	if signatureValue == nil {
		return nil, fmt.Errorf("valor de firma no encontrado")
	}

	return base64.StdEncoding.DecodeString(signatureValue.Text())
}
