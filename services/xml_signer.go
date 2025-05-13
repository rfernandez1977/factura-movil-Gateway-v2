package services

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/beevik/etree"
)

// XMLSigner maneja la firma digital de documentos XML
type XMLSigner struct {
	certificado *x509.Certificate
	privateKey  *rsa.PrivateKey
}

// NewXMLSigner crea una nueva instancia del firmador XML
func NewXMLSigner(certPath, keyPath string) (*XMLSigner, error) {
	// Cargar certificado
	certData, err := ioutil.ReadFile(certPath)
	if err != nil {
		return nil, fmt.Errorf("error leyendo certificado: %w", err)
	}

	block, _ := pem.Decode(certData)
	if block == nil {
		return nil, fmt.Errorf("error decodificando certificado PEM")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("error parseando certificado: %w", err)
	}

	// Cargar clave privada
	keyData, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("error leyendo clave privada: %w", err)
	}

	block, _ = pem.Decode(keyData)
	if block == nil {
		return nil, fmt.Errorf("error decodificando clave privada PEM")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("error parseando clave privada: %w", err)
	}

	return &XMLSigner{
		certificado: cert,
		privateKey:  privateKey,
	}, nil
}

// Firmar firma digitalmente un documento XML
func (s *XMLSigner) Firmar(xmlData []byte) ([]byte, error) {
	// Crear documento XML
	doc := etree.NewDocument()
	if err := doc.ReadFromBytes(xmlData); err != nil {
		return nil, fmt.Errorf("error parseando XML: %w", err)
	}

	// Crear elemento Signature
	signature := doc.CreateElement("ds:Signature")
	signature.CreateAttr("xmlns:ds", "http://www.w3.org/2000/09/xmldsig#")

	// Crear elemento SignedInfo
	signedInfo := signature.CreateElement("ds:SignedInfo")
	signedInfo.CreateAttr("xmlns:ds", "http://www.w3.org/2000/09/xmldsig#")

	// Crear elemento CanonicalizationMethod
	canonicalizationMethod := signedInfo.CreateElement("ds:CanonicalizationMethod")
	canonicalizationMethod.CreateAttr("Algorithm", "http://www.w3.org/TR/2001/REC-xml-c14n-20010315")

	// Crear elemento SignatureMethod
	signatureMethod := signedInfo.CreateElement("ds:SignatureMethod")
	signatureMethod.CreateAttr("Algorithm", "http://www.w3.org/2000/09/xmldsig#rsa-sha1")

	// Crear elemento Reference
	reference := signedInfo.CreateElement("ds:Reference")
	reference.CreateAttr("URI", "")

	// Crear elemento Transforms
	transforms := reference.CreateElement("ds:Transforms")
	transform := transforms.CreateElement("ds:Transform")
	transform.CreateAttr("Algorithm", "http://www.w3.org/2000/09/xmldsig#enveloped-signature")

	// Crear elemento DigestMethod
	digestMethod := reference.CreateElement("ds:DigestMethod")
	digestMethod.CreateAttr("Algorithm", "http://www.w3.org/2000/09/xmldsig#sha1")

	// Calcular digest
	digestValue := s.calcularDigest(xmlData)
	digestValueElement := reference.CreateElement("ds:DigestValue")
	digestValueElement.SetText(base64.StdEncoding.EncodeToString(digestValue))

	// Calcular firma
	signatureValue := s.calcularFirma(signedInfo)
	signatureValueElement := signature.CreateElement("ds:SignatureValue")
	signatureValueElement.SetText(base64.StdEncoding.EncodeToString(signatureValue))

	// Agregar certificado
	keyInfo := signature.CreateElement("ds:KeyInfo")
	x509Data := keyInfo.CreateElement("ds:X509Data")
	x509Certificate := x509Data.CreateElement("ds:X509Certificate")
	x509Certificate.SetText(base64.StdEncoding.EncodeToString(s.certificado.Raw))

	// Generar XML final
	doc.Indent(2)
	return doc.WriteToBytes()
}

// calcularDigest calcula el digest SHA-1 del documento
func (s *XMLSigner) calcularDigest(xmlData []byte) []byte {
	hash := sha1.New()
	hash.Write(xmlData)
	return hash.Sum(nil)
}

// calcularFirma calcula la firma RSA del documento
func (s *XMLSigner) calcularFirma(signedInfo *etree.Element) []byte {
	// Canonicalizar el elemento SignedInfo
	canonicalized, err := signedInfo.Canonicalize()
	if err != nil {
		return nil
	}

	// Calcular hash SHA-1
	hash := sha1.New()
	hash.Write(canonicalized)
	hashed := hash.Sum(nil)

	// Firmar con RSA
	signature, err := rsa.SignPKCS1v15(rand.Reader, s.privateKey, crypto.SHA1, hashed)
	if err != nil {
		return nil
	}

	return signature
}

// VerificarFirma verifica la firma digital de un documento XML
func (s *XMLSigner) VerificarFirma(xmlData []byte) error {
	// TODO: Implementar verificación de firma
	return nil
}

// GuardarCertificado guarda el certificado en formato PEM
func (s *XMLSigner) GuardarCertificado(path string) error {
	// Crear bloque PEM
	block := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: s.certificado.Raw,
	}

	// Guardar en archivo
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("error creando directorio: %w", err)
	}

	if err := ioutil.WriteFile(path, pem.EncodeToMemory(block), 0644); err != nil {
		return fmt.Errorf("error guardando certificado: %w", err)
	}

	return nil
}

// GuardarClavePrivada guarda la clave privada en formato PEM
func (s *XMLSigner) GuardarClavePrivada(path string) error {
	// Crear bloque PEM
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(s.privateKey),
	}

	// Guardar en archivo
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("error creando directorio: %w", err)
	}

	if err := ioutil.WriteFile(path, pem.EncodeToMemory(block), 0600); err != nil {
		return fmt.Errorf("error guardando clave privada: %w", err)
	}

	return nil
}
