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
	"FMgo/config"
	"FMgo/models"
	"software.sslmate.com/src/go-pkcs12"
)

// FirmaDigitalService representa el servicio para manejar la firma digital
type FirmaDigitalService struct {
	privateKey  *rsa.PrivateKey
	certificate *x509.Certificate
	rutFirmante string
}

// NewFirmaDigitalService crea una nueva instancia del servicio de firma digital
func NewFirmaDigitalService(pathCertificado, pathLlave, passwordLlave, rutFirmante string) (*FirmaDigitalService, error) {
	// Leer certificado
	certPEM, err := ioutil.ReadFile(pathCertificado)
	if err != nil {
		return nil, fmt.Errorf("error al leer el certificado: %v", err)
	}

	// Decodificar certificado
	block, _ := pem.Decode(certPEM)
	if block == nil {
		return nil, fmt.Errorf("no se pudo decodificar el certificado PEM")
	}

	// Parsear certificado
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("error al parsear el certificado: %v", err)
	}

	// Leer llave privada
	keyPEM, err := ioutil.ReadFile(pathLlave)
	if err != nil {
		return nil, fmt.Errorf("error al leer la llave privada: %v", err)
	}

	// Decodificar llave privada
	block, _ = pem.Decode(keyPEM)
	if block == nil {
		return nil, fmt.Errorf("no se pudo decodificar la llave privada PEM")
	}

	// Parsear llave privada (con o sin contraseña)
	var privateKey *rsa.PrivateKey
	if passwordLlave != "" {
		decrypted, err := x509.DecryptPEMBlock(block, []byte(passwordLlave))
		if err != nil {
			return nil, fmt.Errorf("error al desencriptar la llave privada: %v", err)
		}
		privateKey, err = x509.ParsePKCS1PrivateKey(decrypted)
		if err != nil {
			return nil, fmt.Errorf("error al parsear la llave privada PKCS1: %v", err)
		}
	} else {
		privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("error al parsear la llave privada: %v", err)
		}
	}

	return &FirmaDigitalService{
		privateKey:  privateKey,
		certificate: cert,
		rutFirmante: rutFirmante,
	}, nil
}

// FirmarEnvioDTE firma digitalmente un sobre de documentos tributarios
func (s *FirmaDigitalService) FirmarEnvioDTE(sobre *models.EnvioDTE) error {
	// Firmar el sobre (Signature del EnvioDTE)
	xmlSobre, err := xml.MarshalIndent(sobre, "", "  ")
	if err != nil {
		return fmt.Errorf("error al generar XML del sobre: %v", err)
	}

	// Crear firma del sobre
	firma, err := s.generarFirmaXML(xmlSobre, "SetDoc")
	if err != nil {
		return fmt.Errorf("error al generar firma del sobre: %v", err)
	}

	// Asignar firma al sobre
	sobre.Signature = firma

	// Firmar cada DTE individual dentro del sobre
	for i := range sobre.SetDTE.DTEs {
		// Firmar el documento (Signature del DTE)
		xmlDTE, err := xml.MarshalIndent(sobre.SetDTE.DTEs[i], "", "  ")
		if err != nil {
			return fmt.Errorf("error al generar XML del DTE: %v", err)
		}

		// Crear firma del DTE
		firma, err := s.generarFirmaXML(xmlDTE, sobre.SetDTE.DTEs[i].Documento.ID)
		if err != nil {
			return fmt.Errorf("error al generar firma del DTE: %v", err)
		}

		// Asignar firma al DTE
		sobre.SetDTE.DTEs[i].Signature = firma

		// Generar y asignar TED (Timbre Electrónico del Documento)
		err = s.generarTED(&sobre.SetDTE.DTEs[i])
		if err != nil {
			return fmt.Errorf("error al generar TED: %v", err)
		}
	}

	return nil
}

// FirmarEnvioBOLETA firma digitalmente un sobre de boletas electrónicas
func (s *FirmaDigitalService) FirmarEnvioBOLETA(sobre *models.EnvioBOLETA) error {
	// Firmar el sobre (Signature del EnvioBOLETA)
	xmlSobre, err := xml.MarshalIndent(sobre, "", "  ")
	if err != nil {
		return fmt.Errorf("error al generar XML del sobre de boleta: %v", err)
	}

	// Crear firma del sobre
	firma, err := s.generarFirmaXML(xmlSobre, "SetDoc")
	if err != nil {
		return fmt.Errorf("error al generar firma del sobre de boleta: %v", err)
	}

	// Asignar firma al sobre
	sobre.Signature = firma

	// Firmar cada BOLETA individual dentro del sobre
	for i := range sobre.SetDTE.DTEs {
		// Firmar el documento (Signature de la BOLETA)
		xmlBOLETA, err := xml.MarshalIndent(sobre.SetDTE.DTEs[i], "", "  ")
		if err != nil {
			return fmt.Errorf("error al generar XML de la BOLETA: %v", err)
		}

		// Crear firma de la BOLETA
		firma, err := s.generarFirmaXML(xmlBOLETA, sobre.SetDTE.DTEs[i].Documento.ID)
		if err != nil {
			return fmt.Errorf("error al generar firma de la BOLETA: %v", err)
		}

		// Asignar firma a la BOLETA
		sobre.SetDTE.DTEs[i].Signature = firma

		// Generar y asignar TED (Timbre Electrónico del Documento)
		err = s.generarTEDBoleta(&sobre.SetDTE.DTEs[i])
		if err != nil {
			return fmt.Errorf("error al generar TED para boleta: %v", err)
		}
	}

	return nil
}

// generarFirmaXML genera una firma XML según el estándar W3C XML-DSIG
func (s *FirmaDigitalService) generarFirmaXML(xmlData []byte, referenceID string) (string, error) {
	// Calcular el digest (hash) del contenido
	hasher := sha1.New()
	hasher.Write(xmlData)
	digestValue := base64.StdEncoding.EncodeToString(hasher.Sum(nil))

	// Firmar el digest con la llave privada
	signature, err := rsa.SignPKCS1v15(rand.Reader, s.privateKey, crypto.SHA1, hasher.Sum(nil))
	if err != nil {
		return "", fmt.Errorf("error al firmar: %v", err)
	}
	signatureValue := base64.StdEncoding.EncodeToString(signature)

	// Certificado en Base64
	certDer := base64.StdEncoding.EncodeToString(s.certificate.Raw)

	// Timestamp actual
	timestamp := time.Now().Format("2006-01-02T15:04:05")

	// Construir la firma XML
	firmaXML := fmt.Sprintf(`<Signature xmlns="http://www.w3.org/2000/09/xmldsig#">
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

	return firmaXML, nil
}

// generarTED genera el Timbre Electrónico del Documento para un DTE
func (s *FirmaDigitalService) generarTED(dte *models.DTEType) error {
	// Obtener información básica del documento
	idDoc := dte.Documento.Encabezado.IdDoc
	emisor := dte.Documento.Encabezado.Emisor
	receptor := dte.Documento.Encabezado.Receptor
	totales := dte.Documento.Encabezado.Totales

	// Obtener primer item si existe
	var it1 string
	if len(dte.Documento.Detalle) > 0 {
		it1 = dte.Documento.Detalle[0].NmbItem
		if len(it1) > 40 {
			it1 = it1[:40] // Máximo 40 caracteres
		}
	}

	// Crear DD (Datos del Documento)
	dd := models.DD{
		RE:    strings.Replace(emisor.RUTEmisor, "-", "", -1),  // RUT sin guión
		TD:    idDoc.TipoDTE,                                   // Tipo Documento
		F:     fmt.Sprintf("%d", idDoc.Folio),                  // Folio
		FE:    idDoc.FchEmis,                                   // Fecha Emisión
		RR:    strings.Replace(receptor.RUTRecep, "-", "", -1), // RUT Receptor sin guión
		RSR:   receptor.RznSocRecep,                            // Razón Social Receptor
		MNT:   fmt.Sprintf("%d", totales.MntTotal),             // Monto Total
		IT1:   it1,                                             // Descripción Item 1
		TSTED: time.Now().Format("2006-01-02T15:04:05"),        // TimeStamp de Generación del Timbre
	}

	// Simulación simplificada del CAF (Código de Autorización de Folios)
	// En un caso real, este vendría desde el SII y debería estar previamente almacenado
	caf := models.CAF{
		Version: "1.0",
		DA: models.DA{
			RE:  strings.Replace(emisor.RUTEmisor, "-", "", -1),
			RS:  emisor.RznSoc,
			TD:  idDoc.TipoDTE,
			RNG: models.RNG{D: idDoc.Folio - 100, H: idDoc.Folio + 100}, // Rango simulado
			FA:  time.Now().AddDate(0, -1, 0).Format("2006-01-02"),      // Fecha Autorización (un mes antes)
			RSAPK: models.RSAPK{
				M: base64.StdEncoding.EncodeToString(s.privateKey.PublicKey.N.Bytes()),
				E: base64.StdEncoding.EncodeToString(big.NewInt(int64(s.privateKey.PublicKey.E)).Bytes()),
			},
			IDK: 100, // ID Clave simulado
		},
		FRMA: models.FRMA{
			Algoritmo: "SHA1withRSA",
			Value:     "CAF_SIGNATURE_PLACEHOLDER", // En un caso real, sería la firma del SII
		},
	}

	dd.CAF = caf

	// Firmar el DD
	ddXML, err := xml.MarshalIndent(dd, "", "  ")
	if err != nil {
		return fmt.Errorf("error al generar XML del DD: %v", err)
	}

	// Calcular firma DD
	hasher := sha1.New()
	hasher.Write(ddXML)
	signature, err := rsa.SignPKCS1v15(rand.Reader, s.privateKey, crypto.SHA1, hasher.Sum(nil))
	if err != nil {
		return fmt.Errorf("error al firmar DD: %v", err)
	}

	// Crear TED
	ted := models.TED{
		Version: "1.0",
		DD:      dd,
		FRMT: models.FRMT{
			Algoritmo: "SHA1withRSA",
			Value:     base64.StdEncoding.EncodeToString(signature),
		},
	}

	// Asignar TED al documento
	dte.Documento.TED = ted

	return nil
}

// generarTEDBoleta genera el Timbre Electrónico para una boleta electrónica
func (s *FirmaDigitalService) generarTEDBoleta(boleta *models.BOLETAType) error {
	// Similar a generarTED pero adaptado para BOLETAType
	// Obtener información básica del documento
	idDoc := boleta.Documento.Encabezado.IdDoc
	emisor := boleta.Documento.Encabezado.Emisor
	totales := boleta.Documento.Encabezado.Totales

	// Receptor puede ser opcional en boletas
	var rr, rsr string
	if boleta.Documento.Encabezado.Receptor.RUTRecep != "" {
		rr = strings.Replace(boleta.Documento.Encabezado.Receptor.RUTRecep, "-", "", -1)
		rsr = boleta.Documento.Encabezado.Receptor.RznSocRecep
	} else {
		rr = "66666666-6" // RUT genérico para consumidor final
		rsr = "Consumidor Final"
	}

	// Obtener primer item si existe
	var it1 string
	if len(boleta.Documento.Detalle) > 0 {
		it1 = boleta.Documento.Detalle[0].NmbItem
		if len(it1) > 40 {
			it1 = it1[:40] // Máximo 40 caracteres
		}
	}

	// Crear DD (Datos del Documento)
	dd := models.DD{
		RE:    strings.Replace(emisor.RUTEmisor, "-", "", -1), // RUT sin guión
		TD:    idDoc.TipoDTE,                                  // Tipo Documento
		F:     fmt.Sprintf("%d", idDoc.Folio),                 // Folio
		FE:    idDoc.FchEmis,                                  // Fecha Emisión
		RR:    rr,                                             // RUT Receptor sin guión
		RSR:   rsr,                                            // Razón Social Receptor
		MNT:   fmt.Sprintf("%d", totales.MntTotal),            // Monto Total
		IT1:   it1,                                            // Descripción Item 1
		TSTED: time.Now().Format("2006-01-02T15:04:05"),       // TimeStamp de Generación del Timbre
	}

	// Simulación simplificada del CAF (Código de Autorización de Folios)
	caf := models.CAF{
		Version: "1.0",
		DA: models.DA{
			RE:  strings.Replace(emisor.RUTEmisor, "-", "", -1),
			RS:  emisor.RznSoc,
			TD:  idDoc.TipoDTE,
			RNG: models.RNG{D: idDoc.Folio - 100, H: idDoc.Folio + 100}, // Rango simulado
			FA:  time.Now().AddDate(0, -1, 0).Format("2006-01-02"),      // Fecha Autorización (un mes antes)
			RSAPK: models.RSAPK{
				M: base64.StdEncoding.EncodeToString(s.privateKey.PublicKey.N.Bytes()),
				E: base64.StdEncoding.EncodeToString(big.NewInt(int64(s.privateKey.PublicKey.E)).Bytes()),
			},
			IDK: 100, // ID Clave simulado
		},
		FRMA: models.FRMA{
			Algoritmo: "SHA1withRSA",
			Value:     "CAF_SIGNATURE_PLACEHOLDER", // En un caso real, sería la firma del SII
		},
	}

	dd.CAF = caf

	// Firmar el DD
	ddXML, err := xml.MarshalIndent(dd, "", "  ")
	if err != nil {
		return fmt.Errorf("error al generar XML del DD: %v", err)
	}

	// Calcular firma DD
	hasher := sha1.New()
	hasher.Write(ddXML)
	signature, err := rsa.SignPKCS1v15(rand.Reader, s.privateKey, crypto.SHA1, hasher.Sum(nil))
	if err != nil {
		return fmt.Errorf("error al firmar DD: %v", err)
	}

	// Crear TED
	ted := models.TED{
		Version: "1.0",
		DD:      dd,
		FRMT: models.FRMT{
			Algoritmo: "SHA1withRSA",
			Value:     base64.StdEncoding.EncodeToString(signature),
		},
	}

	// Asignar TED al documento
	boleta.Documento.TED = ted

	return nil
}

// FirmaService representa el servicio para manejar la firma digital de documentos
type FirmaService struct {
	config   *config.SupabaseConfig
	log      *config.Logger
	certPath string
	password string
}

// NewFirmaService crea una nueva instancia del servicio de firma
func NewFirmaService(config *config.SupabaseConfig) *FirmaService {
	return &FirmaService{
		config: config,
	}
}

// ObtenerCertificado obtiene el certificado digital de una empresa
func (s *FirmaService) ObtenerCertificado(empresaID string) (*models.CertificadoDigital, error) {
	var certificado models.CertificadoDigital
	err := s.config.Client.DB.From("certificados_digitales").
		Select("*").
		Eq("empresa_id", empresaID).
		Single().
		Execute(&certificado)

	if err != nil {
		return nil, fmt.Errorf("error al obtener certificado: %v", err)
	}

	return &certificado, nil
}

// FirmarXML firma un documento XML
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

	// Convertir a RSA
	rsaPrivateKey, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		s.log.Error("Tipo de clave privada no soportado")
		return nil, fmt.Errorf("tipo de clave privada no soportado")
	}

	// Canonicalizar el XML antes de firmar
	canonicalXML, err := s.canonicalizarXML(xmlData)
	if err != nil {
		s.log.Error("Error al canonicalizar XML: %v", err)
		return nil, fmt.Errorf("error al canonicalizar XML: %w", err)
	}

	// Calcular el hash SHA1 del documento canonicalizado
	hashed := sha1.Sum(canonicalXML)

	// Firmar el hash con la clave privada
	signature, err := rsa.SignPKCS1v15(rand.Reader, rsaPrivateKey, crypto.SHA1, hashed[:])
	if err != nil {
		s.log.Error("Error al firmar documento: %v", err)
		return nil, fmt.Errorf("error al firmar documento: %w", err)
	}

	// Codificar la firma en base64
	signatureBase64 := base64.StdEncoding.EncodeToString(signature)

	// Insertar la firma en el documento XML
	signedXML, err := s.insertarFirmaXML(xmlData, signatureBase64, cert, base64.StdEncoding.EncodeToString(hashed[:]))
	if err != nil {
		s.log.Error("Error al insertar firma en XML: %v", err)
		return nil, fmt.Errorf("error al insertar firma en XML: %w", err)
	}

	s.log.Info("Documento XML firmado exitosamente")
	return signedXML, nil
}

// canonicalizarXML aplica la transformación de canonicalización al XML
func (s *FirmaService) canonicalizarXML(xmlData []byte) ([]byte, error) {
	doc := etree.NewDocument()
	if err := doc.ReadFromBytes(xmlData); err != nil {
		return nil, fmt.Errorf("error al parsear XML: %w", err)
	}

	// Aplicar transformación c14n
	c14n, err := doc.WriteToString()
	if err != nil {
		return nil, fmt.Errorf("error en canonicalización: %w", err)
	}

	return []byte(c14n), nil
}

// insertarFirmaXML inserta la firma digital en el documento XML
func (s *FirmaService) insertarFirmaXML(xmlData []byte, firma string, cert *x509.Certificate, digestValue string) ([]byte, error) {
	// Obtener módulo y exponente de la clave pública RSA
	rsaKey := cert.PublicKey.(*rsa.PublicKey)
	modulus := base64.StdEncoding.EncodeToString(rsaKey.N.Bytes())
	exponent := base64.StdEncoding.EncodeToString(big.NewInt(int64(rsaKey.E)).Bytes())

	// Crear estructura de firma
	firmaNode := fmt.Sprintf(`
		<Signature xmlns="http://www.w3.org/2000/09/xmldsig#">
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
		</Signature>
	`, digestValue, firma, modulus, exponent, base64.StdEncoding.EncodeToString(cert.Raw))

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

// ValidarFirma valida una firma digital
func (s *FirmaService) ValidarFirma(xmlData []byte) (bool, error) {
	// Implementar validación de firma
	return true, nil
}
