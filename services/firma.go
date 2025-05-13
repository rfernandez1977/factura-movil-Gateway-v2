package services

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/pkcs12"
)

// FirmaConfig contiene la configuración de la firma electrónica
type FirmaConfig struct {
	RutFirmante    string
	Nombre         string
	CertificadoPEM []byte
	LlavePEM       []byte
	Token          string // Token de autenticación SII
}

// FirmaManager maneja las operaciones relacionadas con la firma electrónica
type FirmaManager struct {
	config      *FirmaConfig
	certificado *x509.Certificate
	llave       *rsa.PrivateKey
}

// NewFirmaManager crea un nuevo manejador de firma
func NewFirmaManager(rutaCertificado, contraseña string) (*FirmaManager, error) {
	// Leer archivo PFX
	datos, err := ioutil.ReadFile(rutaCertificado)
	if err != nil {
		return nil, fmt.Errorf("error leyendo certificado: %v", err)
	}

	// Decodificar PFX
	privateKey, cert, err := pkcs12.Decode(datos, contraseña)
	if err != nil {
		return nil, fmt.Errorf("error decodificando PFX: %v", err)
	}

	// Convertir la llave privada a RSA
	rsaKey, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("la llave privada no es RSA")
	}

	return &FirmaManager{
		certificado: cert,
		llave:       rsaKey,
	}, nil
}

// FirmarXML firma un documento XML usando el certificado
func (fm *FirmaManager) FirmarXML(xmlData []byte) ([]byte, error) {
	// Calcular hash SHA-1 del XML
	hash := sha1.Sum(xmlData)

	// Firmar hash
	firma, err := rsa.SignPKCS1v15(rand.Reader, fm.llave, crypto.SHA1, hash[:])
	if err != nil {
		return nil, fmt.Errorf("error firmando XML: %v", err)
	}

	// Codificar firma en base64
	firmaBase64 := base64.StdEncoding.EncodeToString(firma)

	// Crear nodo de firma
	nodoFirma := fmt.Sprintf(`<Signature xmlns="http://www.w3.org/2000/09/xmldsig#">
  <SignedInfo>
    <CanonicalizationMethod Algorithm="http://www.w3.org/TR/2001/REC-xml-c14n-20010315"/>
    <SignatureMethod Algorithm="http://www.w3.org/2000/09/xmldsig#rsa-sha1"/>
    <Reference URI="">
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
</Signature>`, base64.StdEncoding.EncodeToString(hash[:]), firmaBase64, base64.StdEncoding.EncodeToString(fm.certificado.Raw))

	// Insertar nodo de firma en el XML
	xmlStr := string(xmlData)
	xmlStr = strings.Replace(xmlStr, "</Documento>", nodoFirma+"</Documento>", 1)

	return []byte(xmlStr), nil
}

// ObtenerToken obtiene un token de autenticación del SII
func (fm *FirmaManager) ObtenerToken() (string, error) {
	// Obtener semilla
	semilla, err := fm.GenerarSemilla()
	if err != nil {
		return "", fmt.Errorf("error obteniendo semilla: %w", err)
	}

	// Firmar semilla
	semillaFirmada, err := fm.FirmarSemilla(semilla)
	if err != nil {
		return "", fmt.Errorf("error firmando semilla: %w", err)
	}

	// Crear request SOAP para obtener token
	soapRequest := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/"
xmlns:ns1="http://DefaultNamespace">
<SOAP-ENV:Body>
<ns1:getToken>
<ns1:xml>%s</ns1:xml>
</ns1:getToken>
</SOAP-ENV:Body>
</SOAP-ENV:Envelope>`, semillaFirmada)

	// Preparar request
	req, err := http.NewRequest("POST", SIITokenEndpoint, strings.NewReader(soapRequest))
	if err != nil {
		return "", fmt.Errorf("error creando request: %w", err)
	}

	// Agregar headers SOAP
	req.Header.Set("Content-Type", "text/xml; charset=utf-8")
	req.Header.Set("SOAPAction", "getToken")

	// Enviar al SII
	client := &http.Client{Timeout: time.Second * 30}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error enviando al SII: %w", err)
	}
	defer resp.Body.Close()

	// Procesar respuesta
	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error leyendo respuesta: %w", err)
	}

	// Parsear respuesta SOAP
	var soapResp struct {
		XMLName xml.Name `xml:"Envelope"`
		Body    struct {
			GetTokenResponse struct {
				GetTokenResult string `xml:"getTokenResult"`
			} `xml:"getTokenResponse"`
		} `xml:"Body"`
	}

	if err := xml.Unmarshal(respData, &soapResp); err != nil {
		return "", fmt.Errorf("error decodificando respuesta SOAP: %w", err)
	}

	// Parsear resultado
	var result struct {
		XMLName xml.Name `xml:"RESPUESTA"`
		Header  struct {
			Estado string `xml:"ESTADO"`
			Glosa  string `xml:"GLOSA"`
		} `xml:"RESP_HDR"`
		Body struct {
			Token string `xml:"TOKEN"`
		} `xml:"RESP_BODY"`
	}

	if err := xml.Unmarshal([]byte(soapResp.Body.GetTokenResponse.GetTokenResult), &result); err != nil {
		return "", fmt.Errorf("error decodificando resultado: %w", err)
	}

	if result.Header.Estado != "00" {
		return "", fmt.Errorf("error obteniendo token: %s", result.Header.Glosa)
	}

	return result.Body.Token, nil
}

// GetConfig retorna la configuración actual de la firma
func (fm *FirmaManager) GetConfig() *FirmaConfig {
	return fm.config
}

// extractRutFromCert extrae el RUT del certificado
func extractRutFromCert(cert *x509.Certificate) string {
	// El RUT está en el campo Subject.CommonName
	return cert.Subject.CommonName
}

// GenerarSemilla genera una semilla para autenticación con el SII
func (fm *FirmaManager) GenerarSemilla() (string, error) {
	// Crear request SOAP para obtener semilla
	soapRequest := `<?xml version="1.0" encoding="UTF-8"?>
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/"
xmlns:ns1="http://DefaultNamespace">
<SOAP-ENV:Body>
<ns1:getSeed/>
</SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

	// Preparar request
	req, err := http.NewRequest("POST", SIISemillaEndpoint, strings.NewReader(soapRequest))
	if err != nil {
		return "", fmt.Errorf("error creando request: %w", err)
	}

	// Agregar headers SOAP
	req.Header.Set("Content-Type", "text/xml; charset=utf-8")
	req.Header.Set("SOAPAction", "getSeed")

	// Enviar al SII
	client := &http.Client{Timeout: time.Second * 30}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error enviando al SII: %w", err)
	}
	defer resp.Body.Close()

	// Procesar respuesta
	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error leyendo respuesta: %w", err)
	}

	// Parsear respuesta SOAP
	var soapResp struct {
		XMLName xml.Name `xml:"Envelope"`
		Body    struct {
			GetSeedResponse struct {
				GetSeedResult string `xml:"getSeedResult"`
			} `xml:"getSeedResponse"`
		} `xml:"Body"`
	}

	if err := xml.Unmarshal(respData, &soapResp); err != nil {
		return "", fmt.Errorf("error decodificando respuesta SOAP: %w", err)
	}

	// Parsear resultado
	var result struct {
		XMLName xml.Name `xml:"RESPUESTA"`
		Header  struct {
			Estado string `xml:"ESTADO"`
			Glosa  string `xml:"GLOSA"`
		} `xml:"RESP_HDR"`
		Body struct {
			Semilla string `xml:"SEMILLA"`
		} `xml:"RESP_BODY"`
	}

	if err := xml.Unmarshal([]byte(soapResp.Body.GetSeedResponse.GetSeedResult), &result); err != nil {
		return "", fmt.Errorf("error decodificando resultado: %w", err)
	}

	if result.Header.Estado != "00" {
		return "", fmt.Errorf("error obteniendo semilla: %s", result.Header.Glosa)
	}

	return result.Body.Semilla, nil
}

// FirmarSemilla firma una semilla para obtener el token
func (fm *FirmaManager) FirmarSemilla(semilla string) (string, error) {
	// Calcular hash SHA-1 de la semilla
	hash := sha1.Sum([]byte(semilla))

	// Firmar hash
	firma, err := rsa.SignPKCS1v15(rand.Reader, fm.llave, crypto.SHA1, hash[:])
	if err != nil {
		return "", fmt.Errorf("error firmando semilla: %v", err)
	}

	// Codificar firma en base64
	return base64.StdEncoding.EncodeToString(firma), nil
}
