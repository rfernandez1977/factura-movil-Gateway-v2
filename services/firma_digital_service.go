package services

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"math/big"
	"strings"
	"time"

	"github.com/beevik/etree"
	"github.com/fmgo/models"
)

// XMLSignatureService proporciona métodos para firmar documentos XML digitalmente
type XMLSignatureService struct {
	privateKey  *rsa.PrivateKey
	certificate *x509.Certificate
	rutFirmante string
}

// NewXMLSignatureService crea una nueva instancia del servicio de firma digital XML
func NewXMLSignatureService(pathCertificado, pathLlave, passwordLlave, rutFirmante string) (*XMLSignatureService, error) {
	// Leer la clave privada del archivo
	keyBytes, err := ioutil.ReadFile(pathLlave)
	if err != nil {
		return nil, fmt.Errorf("error al leer la clave privada: %v", err)
	}

	// Decodificar el PEM
	keyBlock, _ := pem.Decode(keyBytes)
	if keyBlock == nil {
		return nil, fmt.Errorf("error al decodificar el PEM de la clave privada")
	}

	// Parsear la clave privada (con o sin contraseña)
	var privateKey *rsa.PrivateKey
	if passwordLlave != "" {
		decrypted, err := x509.DecryptPEMBlock(keyBlock, []byte(passwordLlave))
		if err != nil {
			return nil, fmt.Errorf("error al desencriptar la llave privada: %v", err)
		}
		privateKey, err = x509.ParsePKCS1PrivateKey(decrypted)
		if err != nil {
			return nil, fmt.Errorf("error al parsear la llave privada PKCS1: %v", err)
		}
	} else {
		privateKey, err = x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
		if err != nil {
			return nil, fmt.Errorf("error al parsear la llave privada: %v", err)
		}
	}

	// Leer certificado
	certBytes, err := ioutil.ReadFile(pathCertificado)
	if err != nil {
		return nil, fmt.Errorf("error al leer el certificado: %v", err)
	}

	// Decodificar el certificado PEM
	certBlock, _ := pem.Decode(certBytes)
	if certBlock == nil {
		return nil, fmt.Errorf("error al decodificar el certificado PEM")
	}

	// Parsear certificado
	cert, err := x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("error al parsear el certificado: %v", err)
	}

	return &XMLSignatureService{
		privateKey:  privateKey,
		certificate: cert,
		rutFirmante: rutFirmante,
	}, nil
}

// FirmarXML firma un documento XML según el estándar XML-DSIG
func (s *XMLSignatureService) FirmarXML(xml string) (string, error) {
	// Crear un documento XML usando etree
	doc := etree.NewDocument()
	if err := doc.ReadFromString(xml); err != nil {
		return "", fmt.Errorf("error al parsear XML: %v", err)
	}

	// Buscar el elemento raíz para identificar el tipo de documento
	root := doc.Root()
	if root == nil {
		return "", fmt.Errorf("XML sin elemento raíz")
	}

	// Determinar ID de referencia (SetDoc para sobres, o el ID específico del documento)
	referenceID := "SetDoc"
	if idAttr := root.SelectAttr("ID"); idAttr != nil {
		referenceID = idAttr.Value
	}

	// Canonicalizar el XML para la firma
	canonicalXML, err := doc.WriteToString()
	if err != nil {
		return "", fmt.Errorf("error al canonicalizar XML: %v", err)
	}

	// Calcular el digest SHA-1 del XML canonicalizado
	digestSHA1 := sha1.Sum([]byte(canonicalXML))
	digestValue := base64.StdEncoding.EncodeToString(digestSHA1[:])

	// Firmar el digest con la llave privada
	signature, err := rsa.SignPKCS1v15(rand.Reader, s.privateKey, crypto.SHA1, digestSHA1[:])
	if err != nil {
		return "", fmt.Errorf("error al firmar: %v", err)
	}
	signatureValue := base64.StdEncoding.EncodeToString(signature)

	// Crear el nodo Signature
	signatureXML := s.crearNodoSignature(referenceID, digestValue, signatureValue)

	// Agregar la firma al XML original
	return s.agregarFirmaAXML(xml, signatureXML)
}

// crearNodoSignature genera el nodo XML de firma según el estándar XML-DSIG
func (s *XMLSignatureService) crearNodoSignature(referenceID, digestValue, signatureValue string) string {
	// Certificado en Base64
	certDer := base64.StdEncoding.EncodeToString(s.certificate.Raw)

	// Timestamp actual (usado para algunos elementos)
	timestamp := time.Now().Format("2006-01-02T15:04:05")

	// Construir la firma XML según estándar XML-DSIG
	return fmt.Sprintf(`<Signature xmlns="http://www.w3.org/2000/09/xmldsig#">
  <SignedInfo>
    <CanonicalizationMethod Algorithm="http://www.w3.org/TR/2001/REC-xml-c14n-20010315"/>
    <SignatureMethod Algorithm="http://www.w3.org/2000/09/xmldsig#rsa-sha1"/>
    <Reference URI="#%s">
      <Transforms>
        <Transform Algorithm="http://www.w3.org/2000/09/xmldsig#enveloped-signature"/>
      </Transforms>
      <DigestMethod Algorithm="http://www.w3.org/2000/09/xmldsig#sha1"/>
      <DigestValue>%s</DigestValue>
    </Reference>
  </SignedInfo>
  <SignatureValue>%s</SignatureValue>
  <KeyInfo>
    <KeyValue>
      <RSAKeyValue>
        <Modulus>%s</Modulus>
        <Exponent>%s</Exponent>
      </RSAKeyValue>
    </KeyValue>
    <X509Data>
      <X509Certificate>%s</X509Certificate>
    </X509Data>
  </KeyInfo>
</Signature>`, referenceID, digestValue, signatureValue,
		base64.StdEncoding.EncodeToString(s.privateKey.PublicKey.N.Bytes()),
		base64.StdEncoding.EncodeToString(big.NewInt(int64(s.privateKey.PublicKey.E)).Bytes()),
		certDer)
}

// agregarFirmaAXML agrega el nodo de firma al XML original
func (s *XMLSignatureService) agregarFirmaAXML(xmlOriginal, firma string) (string, error) {
	// Para agregar la firma correctamente, necesitamos analizar el XML
	doc := etree.NewDocument()
	if err := doc.ReadFromString(xmlOriginal); err != nil {
		return "", fmt.Errorf("error al leer XML original: %v", err)
	}

	// Parsear la firma
	firmaDoc := etree.NewDocument()
	if err := firmaDoc.ReadFromString(firma); err != nil {
		return "", fmt.Errorf("error al parsear firma: %v", err)
	}

	// Agregar la firma como último elemento del documento
	if firmaDoc.Root() != nil {
		doc.Root().AddChild(firmaDoc.Root())
	}

	// Convertir a string
	resultado, err := doc.WriteToString()
	if err != nil {
		return "", fmt.Errorf("error al escribir XML firmado: %v", err)
	}

	return resultado, nil
}

// FirmarEnvioDTE firma un sobre de envío de DTE completo
func (s *XMLSignatureService) FirmarEnvioDTE(sobre *models.EnvioDTE) error {
	// Para firmar correctamente, debemos:
	// 1. Generar primero el XML del sobre sin la firma
	// 2. Firmar ese XML
	// 3. Asignar la firma al campo Signature del sobre

	// Limpiamos cualquier firma previa
	sobre.Signature = ""

	// Convertir el sobre a XML sin firma
	xmlSinFirma, err := xml.MarshalIndent(sobre, "", "  ")
	if err != nil {
		return fmt.Errorf("error al generar XML del sobre: %v", err)
	}

	// Configurar atributos necesarios en el sobre
	sobreXML := etree.NewDocument()
	if err := sobreXML.ReadFromBytes(xmlSinFirma); err != nil {
		return fmt.Errorf("error al leer XML del sobre: %v", err)
	}

	// Asegurar que el elemento SetDTE tenga un ID
	if setDTE := sobreXML.FindElement("//SetDTE"); setDTE != nil {
		if setDTE.SelectAttr("ID") == nil {
			setDTE.CreateAttr("ID", "SetDoc")
		}
	}

	// Volver a serializar con los atributos correctos
	xmlConAtributos, err := sobreXML.WriteToBytes()
	if err != nil {
		return fmt.Errorf("error al generar XML con atributos: %v", err)
	}

	// Firmar el XML
	xmlFirmado, err := s.FirmarXML(string(xmlConAtributos))
	if err != nil {
		return fmt.Errorf("error al firmar sobre: %v", err)
	}

	// Extraer solo el nodo Signature
	signatureDoc := etree.NewDocument()
	if err := signatureDoc.ReadFromString(xmlFirmado); err != nil {
		return fmt.Errorf("error al leer XML firmado: %v", err)
	}

	signatureNode := signatureDoc.FindElement("//Signature")
	if signatureNode == nil {
		return fmt.Errorf("no se encontró el nodo Signature en el XML firmado")
	}

	signatureXML, err := signatureNode.WriteToString()
	if err != nil {
		return fmt.Errorf("error al extraer nodo Signature: %v", err)
	}

	// Asignar la firma al sobre
	sobre.Signature = signatureXML

	// Agregar TmstFirma a cada documento
	timestamp := time.Now().Format("2006-01-02T15:04:05")
	for i := range sobre.SetDTE.DTEs {
		sobre.SetDTE.DTEs[i].Documento.TmstFirma = timestamp
	}

	return nil
}

// FirmarEnvioBOLETA firma un sobre de envío de Boletas
func (s *XMLSignatureService) FirmarEnvioBOLETA(sobre *models.EnvioBOLETA) error {
	// Seguimos el mismo proceso que para FirmarEnvioDTE
	sobre.Signature = ""

	xmlSinFirma, err := xml.MarshalIndent(sobre, "", "  ")
	if err != nil {
		return fmt.Errorf("error al generar XML del sobre de boleta: %v", err)
	}

	// Configurar atributos necesarios
	sobreXML := etree.NewDocument()
	if err := sobreXML.ReadFromBytes(xmlSinFirma); err != nil {
		return fmt.Errorf("error al leer XML del sobre de boleta: %v", err)
	}

	if setDTE := sobreXML.FindElement("//SetDTE"); setDTE != nil {
		if setDTE.SelectAttr("ID") == nil {
			setDTE.CreateAttr("ID", "SetDoc")
		}
	}

	xmlConAtributos, err := sobreXML.WriteToBytes()
	if err != nil {
		return fmt.Errorf("error al generar XML con atributos: %v", err)
	}

	xmlFirmado, err := s.FirmarXML(string(xmlConAtributos))
	if err != nil {
		return fmt.Errorf("error al firmar sobre de boleta: %v", err)
	}

	signatureDoc := etree.NewDocument()
	if err := signatureDoc.ReadFromString(xmlFirmado); err != nil {
		return fmt.Errorf("error al leer XML firmado: %v", err)
	}

	signatureNode := signatureDoc.FindElement("//Signature")
	if signatureNode == nil {
		return fmt.Errorf("no se encontró el nodo Signature en el XML firmado")
	}

	signatureXML, err := signatureNode.WriteToString()
	if err != nil {
		return fmt.Errorf("error al extraer nodo Signature: %v", err)
	}

	sobre.Signature = signatureXML

	// Agregar TmstFirma a cada boleta
	timestamp := time.Now().Format("2006-01-02T15:04:05")
	for i := range sobre.SetDTE.DTEs {
		sobre.SetDTE.DTEs[i].Documento.TmstFirma = timestamp
	}

	return nil
}

// VerificarFirma verifica la firma de un documento XML
func (s *XMLSignatureService) VerificarFirma(xmlFirmado string) (bool, error) {
	doc := etree.NewDocument()
	if err := doc.ReadFromString(xmlFirmado); err != nil {
		return false, fmt.Errorf("error al leer XML: %v", err)
	}

	// Encontrar el nodo Signature
	signatureNode := doc.FindElement("//Signature")
	if signatureNode == nil {
		return false, fmt.Errorf("no se encontró firma en el documento")
	}

	// Extraer el nodo SignedInfo
	signedInfoNode := signatureNode.FindElement("SignedInfo")
	if signedInfoNode == nil {
		return false, fmt.Errorf("no se encontró nodo SignedInfo")
	}

	// Extraer el valor del digest
	digestValueNode := signedInfoNode.FindElement(".//DigestValue")
	if digestValueNode == nil {
		return false, fmt.Errorf("no se encontró nodo DigestValue")
	}
	digestValue := digestValueNode.Text()

	// Extraer el valor de la firma
	signatureValueNode := signatureNode.FindElement("SignatureValue")
	if signatureValueNode == nil {
		return false, fmt.Errorf("no se encontró nodo SignatureValue")
	}
	signatureValue := signatureValueNode.Text()

	// Extraer el ID de referencia
	referenceNode := signedInfoNode.FindElement(".//Reference")
	if referenceNode == nil {
		return false, fmt.Errorf("no se encontró nodo Reference")
	}

	uriAttr := referenceNode.SelectAttr("URI")
	if uriAttr == nil {
		return false, fmt.Errorf("falta atributo URI en Reference")
	}

	referenceID := strings.TrimPrefix(uriAttr.Value, "#")

	// Recrear el hash del documento original (sin la firma)
	// Esto requiere clonar el documento, quitar la firma y calcular el hash
	docClone := doc.Copy()
	signatureNodeClone := docClone.FindElement("//Signature")
	if signatureNodeClone != nil {
		signatureNodeClone.Parent().RemoveChild(signatureNodeClone)
	}

	// Canonicalizar y calcular hash
	xmlSinFirma, err := docClone.WriteToString()
	if err != nil {
		return false, fmt.Errorf("error al canonicalizar XML sin firma: %v", err)
	}

	hasher := sha1.New()
	hasher.Write([]byte(xmlSinFirma))
	calculatedDigest := base64.StdEncoding.EncodeToString(hasher.Sum(nil))

	// Verificar que el digest coincida
	if calculatedDigest != digestValue {
		return false, fmt.Errorf("el digest calculado no coincide con el almacenado")
	}

	// Decodificar la firma
	signatureBytes, err := base64.StdEncoding.DecodeString(signatureValue)
	if err != nil {
		return false, fmt.Errorf("error al decodificar firma: %v", err)
	}

	// Verificar la firma
	err = rsa.VerifyPKCS1v15(&s.certificate.PublicKey.(*rsa.PublicKey), crypto.SHA1, hasher.Sum(nil), signatureBytes)
	if err != nil {
		return false, fmt.Errorf("la firma no es válida: %v", err)
	}

	return true, nil
}
