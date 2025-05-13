package utils

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"

	"github.com/cursor/FMgo/models"
)

// FirmadorXML representa un firmador de documentos XML
type FirmadorXML struct {
	privateKey *rsa.PrivateKey
}

// NuevoFirmadorXML crea un nuevo firmador de documentos XML
func NuevoFirmadorXML(privateKeyPEM []byte) (*FirmadorXML, error) {
	block, _ := pem.Decode(privateKeyPEM)
	if block == nil {
		return nil, fmt.Errorf("error decodificando PEM")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("error parseando clave privada: %v", err)
	}

	return &FirmadorXML{
		privateKey: privateKey,
	}, nil
}

// FirmarDTE firma un documento DTE
func (f *FirmadorXML) FirmarDTE(dte *models.DTEXMLModel) error {
	// Generar el hash del documento
	hash := sha1.New()
	hash.Write([]byte(dte.Documento.ID))
	hash.Write([]byte(dte.Documento.Encabezado.IdDoc.FechaEmision))
	hash.Write([]byte(dte.Documento.Encabezado.IdDoc.TipoDTE))
	hash.Write([]byte(dte.Documento.Encabezado.Emisor.RUT))
	hash.Write([]byte(dte.Documento.Encabezado.Receptor.RUT))
	montoTotal := fmt.Sprintf("%d", dte.Documento.Encabezado.Totales.MntTotal)
	hash.Write([]byte(montoTotal))

	// Firmar el hash
	signature, err := rsa.SignPKCS1v15(rand.Reader, f.privateKey, crypto.SHA1, hash.Sum(nil))
	if err != nil {
		return fmt.Errorf("error firmando documento: %v", err)
	}

	// Agregar la firma al documento
	dte.Signature = &models.FirmaXMLModel{
		SignedInfo: models.SignedInfoXML{
			CanonicalizationMethod: models.CanonicalizationMethodXML{
				Algorithm: "http://www.w3.org/TR/2001/REC-xml-c14n-20010315",
			},
			SignatureMethod: models.SignatureMethodXML{
				Algorithm: "http://www.w3.org/2000/09/xmldsig#rsa-sha1",
			},
			Reference: models.ReferenceSignatureXML{
				URI: "#" + dte.Documento.ID,
				DigestMethod: models.DigestMethodXML{
					Algorithm: "http://www.w3.org/2000/09/xmldsig#sha1",
				},
				DigestValue: base64.StdEncoding.EncodeToString(hash.Sum(nil)),
			},
		},
		SignatureValue: base64.StdEncoding.EncodeToString(signature),
		KeyInfo: models.KeyInfoXML{
			KeyValue: models.KeyValueXML{
				RSAKeyValue: models.RSAKeyValueXML{
					Modulus:  base64.StdEncoding.EncodeToString(f.privateKey.PublicKey.N.Bytes()),
					Exponent: "AQAB", // RSA estándar
				},
			},
		},
	}

	return nil
}

// VerificarFirmaDTE verifica la firma de un documento DTE
func (f *FirmadorXML) VerificarFirmaDTE(dte *models.DTEXMLModel) error {
	if dte.Signature == nil {
		return fmt.Errorf("documento no firmado")
	}

	// Generar el hash del documento
	hash := sha1.New()
	hash.Write([]byte(dte.Documento.ID))
	hash.Write([]byte(dte.Documento.Encabezado.IdDoc.FechaEmision))
	hash.Write([]byte(dte.Documento.Encabezado.IdDoc.TipoDTE))
	hash.Write([]byte(dte.Documento.Encabezado.Emisor.RUT))
	hash.Write([]byte(dte.Documento.Encabezado.Receptor.RUT))
	montoTotal := fmt.Sprintf("%d", dte.Documento.Encabezado.Totales.MntTotal)
	hash.Write([]byte(montoTotal))

	// Decodificar la firma
	signature, err := base64.StdEncoding.DecodeString(dte.Signature.SignatureValue)
	if err != nil {
		return fmt.Errorf("error decodificando firma: %v", err)
	}

	// Verificar la firma
	err = rsa.VerifyPKCS1v15(&f.privateKey.PublicKey, crypto.SHA1, hash.Sum(nil), signature)
	if err != nil {
		return fmt.Errorf("error verificando firma: %v", err)
	}

	return nil
}
