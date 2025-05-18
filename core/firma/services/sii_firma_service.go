package services

import (
	"crypto"
	"crypto/sha1"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"time"

	"github.com/beevik/etree"
	"FMgo/core/firma/models"
)

// SIIFirmaService extiende BaseFirmaService para el SII
type SIIFirmaService struct {
	*BaseFirmaService
	certCache *CertCache
}

// NewSIIFirmaService crea una nueva instancia del servicio de firma para el SII
func NewSIIFirmaService(config *models.ConfiguracionFirma) (*SIIFirmaService, error) {
	baseService, err := NewBaseFirmaService(config)
	if err != nil {
		return nil, fmt.Errorf("error creando servicio base: %w", err)
	}

	return &SIIFirmaService{
		BaseFirmaService: baseService,
		certCache:        NewCertCache(24*time.Hour, 100), // Caché por 24 horas, máximo 100 items
	}, nil
}

// FirmarSemilla implementa la firma de semillas del SII
func (s *SIIFirmaService) FirmarSemilla(semilla string) (*models.ResultadoFirma, error) {
	// Construir XML de semilla
	xmlSemilla := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<SemillaXML>
    <Semilla>%s</Semilla>
</SemillaXML>`, semilla)

	// Usar el servicio base para firmar
	return s.FirmarDocumento(xmlSemilla)
}

// FirmarToken implementa la firma de tokens de autenticación
func (s *SIIFirmaService) FirmarToken(token string) (*models.ResultadoFirma, error) {
	// Construir XML de token
	xmlToken := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<TokenXML>
    <Token>%s</Token>
</TokenXML>`, token)

	// Usar el servicio base para firmar
	return s.FirmarDocumento(xmlToken)
}

// ValidarCAF implementa la validación de archivos CAF
func (s *SIIFirmaService) ValidarCAF(caf []byte) error {
	// Parsear CAF
	var cafData struct {
		RutEmpresa        string `xml:"RE"`
		RutFirmante       string `xml:"RS"`
		TipoDocumento     string `xml:"TD"`
		RangoInicial      int    `xml:"RNG>D"`
		RangoFinal        int    `xml:"RNG>H"`
		FechaAutorizacion string `xml:"FA"`
		Firma             string `xml:"FRMA>SignatureValue"`
	}

	if err := xml.Unmarshal(caf, &cafData); err != nil {
		return fmt.Errorf("error al parsear CAF: %w", err)
	}

	// Validar firma del CAF
	doc := etree.NewDocument()
	if err := doc.ReadFromBytes(caf); err != nil {
		return fmt.Errorf("error al leer CAF: %w", err)
	}

	// Extraer y validar firma
	firma, err := s.extraerFirma(doc)
	if err != nil {
		return fmt.Errorf("error al extraer firma del CAF: %w", err)
	}

	// Calcular hash del contenido
	hasher := sha1.New()
	hasher.Write([]byte(fmt.Sprintf("%s%s%s%d%d%s",
		cafData.RutEmpresa,
		cafData.RutFirmante,
		cafData.TipoDocumento,
		cafData.RangoInicial,
		cafData.RangoFinal,
		cafData.FechaAutorizacion)))

	// Verificar firma
	if err := s.privateKey.Verify(hasher.Sum(nil), firma, crypto.SHA1); err != nil {
		return fmt.Errorf("firma del CAF inválida: %w", err)
	}

	return nil
}

// construirXMLFirmadoSII construye un XML firmado según especificaciones del SII
func (s *SIIFirmaService) construirXMLFirmadoSII(doc *etree.Document, digest, firma []byte) (string, error) {
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
	transforms.CreateElement("ds:Transform").
		CreateAttr("Algorithm", "http://www.w3.org/TR/2001/REC-xml-c14n-20010315")
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
	x509IssuerSerial := x509Data.CreateElement("ds:X509IssuerSerial")
	x509IssuerSerial.CreateElement("ds:X509IssuerName").
		SetText(s.cert.Issuer.String())
	x509IssuerSerial.CreateElement("ds:X509SerialNumber").
		SetText(s.cert.SerialNumber.String())

	// Agregar firma al documento
	doc.Root().AddChild(signatureElement)

	// Convertir a string
	xmlString, err := doc.WriteToString()
	if err != nil {
		return "", fmt.Errorf("error al convertir XML a string: %v", err)
	}

	return xmlString, nil
}
