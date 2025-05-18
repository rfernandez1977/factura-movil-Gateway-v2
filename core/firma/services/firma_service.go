package services

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"math/big"
	"time"

	"github.com/beevik/etree"
	"FMgo/core/firma/models"
	"software.sslmate.com/src/go-pkcs12"
)

// FirmaService maneja las operaciones de firma digital
type FirmaService struct {
	config     *models.CertConfig
	privateKey interface{}
	cert       *x509.Certificate
	xmlService *XMLService
	logService *LogService
}

// NewFirmaService crea una nueva instancia del servicio de firma
func NewFirmaService(configPath string) (*FirmaService, error) {
	// Crear servicio de logging
	logService, err := NewLogService("logs/firma")
	if err != nil {
		return nil, fmt.Errorf("error inicializando servicio de logs: %w", err)
	}

	// Leer configuración
	configData, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("error al leer configuración: %v", err)
	}

	var config models.CertConfig
	if err := xml.Unmarshal(configData, &config); err != nil {
		return nil, fmt.Errorf("error al parsear configuración: %v", err)
	}

	// Cargar certificado PKCS12
	pfxData, err := ioutil.ReadFile(config.Certificado.RutaPfx)
	if err != nil {
		return nil, fmt.Errorf("error al leer certificado PFX: %v", err)
	}

	// Extraer llave privada y certificado
	privateKey, cert, err := pkcs12.Decode(pfxData, config.Certificado.Password)
	if err != nil {
		return nil, fmt.Errorf("error al decodificar PKCS12: %v", err)
	}

	return &FirmaService{
		config:     &config,
		privateKey: privateKey,
		cert:       cert,
		xmlService: NewXMLService(),
		logService: logService,
	}, nil
}

// FirmarDocumento firma un documento XML
func (s *FirmaService) FirmarDocumento(documento string) (*models.ResultadoFirma, error) {
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
	}, nil
}

// FirmarSemilla firma una semilla del SII
func (s *FirmaService) FirmarSemilla(semilla string) (*models.ResultadoFirma, error) {
	// Construir XML de semilla
	xmlSemilla := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<SemillaXML>
    <Semilla>%s</Semilla>
</SemillaXML>`, semilla)

	return s.FirmarDocumento(xmlSemilla)
}

// ValidarFirma valida una firma digital
func (s *FirmaService) ValidarFirma(xmlFirmado string) (*models.EstadoFirma, error) {
	doc := etree.NewDocument()
	if err := doc.ReadFromString(xmlFirmado); err != nil {
		return nil, fmt.Errorf("error al leer XML: %v", err)
	}

	// Implementar validación de firma
	// TODO: Implementar lógica de validación

	return &models.EstadoFirma{
		Valida:          true,
		FechaValidacion: time.Now(),
		CertificadoID:   s.cert.SerialNumber.String(),
	}, nil
}

// Métodos privados auxiliares
func (s *FirmaService) canonicalizar(doc *etree.Document) ([]byte, error) {
	// Obtener el elemento raíz
	root := doc.Root()
	if root == nil {
		return nil, fmt.Errorf("documento XML sin elemento raíz")
	}

	// Aplicar canonicalización C14N
	c14nBytes, err := root.WriteToBytes()
	if err != nil {
		return nil, err
	}

	return c14nBytes, nil
}

func (s *FirmaService) calcularDigest(data []byte) []byte {
	// Calcular hash SHA-256
	hasher := sha256.New()
	hasher.Write(data)
	return hasher.Sum(nil)
}

func (s *FirmaService) firmarDigest(digest []byte) ([]byte, error) {
	// Convertir la llave privada genérica a RSA
	rsaKey, ok := s.privateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("la llave privada no es de tipo RSA")
	}

	// Firmar el digest usando PKCS1v15 y SHA-256
	signature, err := rsa.SignPKCS1v15(rand.Reader, rsaKey, crypto.SHA256, digest)
	if err != nil {
		return nil, fmt.Errorf("error al firmar digest: %v", err)
	}

	return signature, nil
}

func (s *FirmaService) construirXMLFirmado(doc *etree.Document, digest, firma []byte) (string, error) {
	// Crear elemento Signature
	signatureElement := etree.NewElement("ds:Signature")
	signatureElement.CreateAttr("xmlns:ds", "http://www.w3.org/2000/09/xmldsig#")

	// SignedInfo
	signedInfo := signatureElement.CreateElement("ds:SignedInfo")
	signedInfo.CreateElement("ds:CanonicalizationMethod").
		CreateAttr("Algorithm", "http://www.w3.org/TR/2001/REC-xml-c14n-20010315")
	signedInfo.CreateElement("ds:SignatureMethod").
		CreateAttr("Algorithm", "http://www.w3.org/2000/09/xmldsig#rsa-sha1")

	// Reference
	reference := signedInfo.CreateElement("ds:Reference")
	reference.CreateAttr("URI", "")
	transforms := reference.CreateElement("ds:Transforms")
	transforms.CreateElement("ds:Transform").
		CreateAttr("Algorithm", "http://www.w3.org/2000/09/xmldsig#enveloped-signature")
	reference.CreateElement("ds:DigestMethod").
		CreateAttr("Algorithm", "http://www.w3.org/2000/09/xmldsig#sha1")
	reference.CreateElement("ds:DigestValue").
		SetText(base64.StdEncoding.EncodeToString(digest))

	// SignatureValue
	signatureElement.CreateElement("ds:SignatureValue").
		SetText(base64.StdEncoding.EncodeToString(firma))

	// KeyInfo
	keyInfo := signatureElement.CreateElement("ds:KeyInfo")
	x509Data := keyInfo.CreateElement("ds:X509Data")
	x509Data.CreateElement("ds:X509Certificate").
		SetText(base64.StdEncoding.EncodeToString(s.cert.Raw))

	// Agregar firma al documento original
	doc.Root().AddChild(signatureElement)

	// Convertir a string
	xmlString, err := doc.WriteToString()
	if err != nil {
		return "", fmt.Errorf("error al convertir XML a string: %v", err)
	}

	return xmlString, nil
}

// FirmarXML firma un documento XML
func (s *FirmaService) FirmarXML(xmlData []byte, referenceID string) (*models.ResultadoFirma, error) {
	// Log del XML original
	if err := s.logService.LogXML("original", xmlData); err != nil {
		return nil, fmt.Errorf("error logging XML original: %w", err)
	}

	// Validar namespaces antes de procesar
	if err := s.xmlService.ValidarNamespaces(xmlData); err != nil {
		s.logService.LogError("namespaces", err)
		return nil, fmt.Errorf("error en namespaces: %w", err)
	}

	// Agregar namespaces necesarios
	xmlConNamespaces, err := s.xmlService.AgregarNamespaces(xmlData)
	if err != nil {
		s.logService.LogError("agregar_namespaces", err)
		return nil, fmt.Errorf("error agregando namespaces: %w", err)
	}

	// Log del XML con namespaces
	if err := s.logService.LogXML("con_namespaces", xmlConNamespaces); err != nil {
		return nil, fmt.Errorf("error logging XML con namespaces: %w", err)
	}

	// Aplicar canonicalización C14N
	canonicalXML, err := s.xmlService.Canonicalizar(xmlConNamespaces)
	if err != nil {
		s.logService.LogError("canonicalizacion", err)
		return nil, fmt.Errorf("error en canonicalización: %w", err)
	}

	// Log del XML canonicalizado
	if err := s.logService.LogXML("canonicalizado", canonicalXML); err != nil {
		return nil, fmt.Errorf("error logging XML canonicalizado: %w", err)
	}

	// Calcular el digest del contenido canonicalizado
	hasher := sha1.New()
	hasher.Write(canonicalXML)
	digestValue := base64.StdEncoding.EncodeToString(hasher.Sum(nil))

	// Firmar el digest con la llave privada
	signature, err := rsa.SignPKCS1v15(rand.Reader, s.privateKey.(rsa.PrivateKey), crypto.SHA1, hasher.Sum(nil))
	if err != nil {
		return nil, fmt.Errorf("error al firmar: %w", err)
	}
	signatureValue := base64.StdEncoding.EncodeToString(signature)

	// Certificado en Base64
	certDer := base64.StdEncoding.EncodeToString(s.cert.Raw)

	// Obtener información del emisor del certificado
	issuerName := s.cert.Issuer.String()
	serialNumber := s.cert.SerialNumber.String()

	// Construir la firma XML según el esquema del SII
	firmaXML := fmt.Sprintf(`<Signature xmlns="http://www.w3.org/2000/09/xmldsig#">
		<SignedInfo>
			<CanonicalizationMethod Algorithm="http://www.w3.org/TR/2001/REC-xml-c14n-20010315"/>
			<SignatureMethod Algorithm="http://www.w3.org/2000/09/xmldsig#rsa-sha1"/>
			<Reference URI="">
				<Transforms>
					<Transform Algorithm="http://www.w3.org/2000/09/xmldsig#enveloped-signature"/>
					<Transform Algorithm="http://www.w3.org/TR/2001/REC-xml-c14n-20010315"/>
				</Transforms>
				<DigestMethod Algorithm="http://www.w3.org/2000/09/xmldsig#sha1"/>
				<DigestValue>%s</DigestValue>
			</Reference>
		</SignedInfo>
		<SignatureValue>%s</SignatureValue>
		<KeyInfo>
			<X509Data>
				<X509Certificate>%s</X509Certificate>
				<X509IssuerSerial>
					<X509IssuerName>%s</X509IssuerName>
					<X509SerialNumber>%s</X509SerialNumber>
				</X509IssuerSerial>
			</X509Data>
			<KeyValue>
				<RSAKeyValue>
					<Modulus>%s</Modulus>
					<Exponent>%s</Exponent>
				</RSAKeyValue>
			</KeyValue>
		</KeyInfo>
	</Signature>`, digestValue, signatureValue,
		certDer, issuerName, serialNumber,
		base64.StdEncoding.EncodeToString(s.privateKey.(rsa.PrivateKey).N.Bytes()),
		base64.StdEncoding.EncodeToString(big.NewInt(int64(s.privateKey.(rsa.PrivateKey).E)).Bytes()))

	// Log de la firma XML final
	if err := s.logService.LogXML("firma_final", []byte(firmaXML)); err != nil {
		return nil, fmt.Errorf("error logging firma final: %w", err)
	}

	// Formatear la firma XML para mejor legibilidad
	firmaFormateada, err := s.xmlService.FormatearXML([]byte(firmaXML))
	if err != nil {
		s.logService.LogError("formateo", err)
		return nil, fmt.Errorf("error formateando firma: %w", err)
	}

	// Canonicalizar la firma XML
	firmaCanonicalizada, err := s.xmlService.Canonicalizar(firmaFormateada)
	if err != nil {
		s.logService.LogError("canonicalizacion_firma", err)
		return nil, fmt.Errorf("error canonicalizando firma: %w", err)
	}

	// Validar estructura contra esquema XSD
	if xsdPath := "schemas/xmldsignature_v10.xsd"; xsdPath != "" {
		if err := s.xmlService.ValidarEstructura(firmaCanonicalizada, xsdPath); err != nil {
			s.logService.LogError("validacion_xsd", err)
			return nil, fmt.Errorf("error de validación contra esquema: %w", err)
		}
		// Log de validación exitosa
		if err := s.logService.LogValidacion(firmaCanonicalizada, "Validación XSD exitosa"); err != nil {
			return nil, fmt.Errorf("error logging validación: %w", err)
		}
	}

	// Crear el resultado
	resultado := &models.ResultadoFirma{
		XMLFirmado:     string(firmaCanonicalizada),
		DigestValue:    digestValue,
		SignatureValue: signatureValue,
	}

	return resultado, nil
}
