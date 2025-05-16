package caf

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha1"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"math/big"
)

var (
	ErrFirmaNoEncontrada    = errors.New("firma no encontrada en CAF")
	ErrClavePublicaInvalida = errors.New("clave pública inválida")
	ErrVerificacionFirma    = errors.New("error al verificar firma")
)

// RSAKey representa la estructura de la clave pública RSA en el CAF
type RSAKey struct {
	XMLName  xml.Name `xml:"RSAPK"`
	Modulus  string   `xml:"M"`
	Exponent string   `xml:"E"`
}

// SignedInfo representa la información firmada del CAF
type SignedInfo struct {
	XMLName xml.Name `xml:"DA"`
	RE      string   `xml:"RE"` // RUT Emisor
	TD      int      `xml:"TD"` // Tipo de DTE
	RNG     struct {
		D int `xml:"D"` // Desde
		H int `xml:"H"` // Hasta
	} `xml:"RNG"`
	RSAPK RSAKey `xml:"RSAPK"`
}

// SignatureVerifier maneja la verificación de firmas
type SignatureVerifier struct {
	publicKey *rsa.PublicKey
}

// NewSignatureVerifier crea un nuevo verificador de firmas
func NewSignatureVerifier(key *rsa.PublicKey) *SignatureVerifier {
	return &SignatureVerifier{
		publicKey: key,
	}
}

// ParsePublicKey extrae la clave pública del CAF
func ParsePublicKey(key RSAKey) (*rsa.PublicKey, error) {
	// Decodificar módulo y exponente de base64
	modulus, err := base64.StdEncoding.DecodeString(key.Modulus)
	if err != nil {
		return nil, fmt.Errorf("%w: error decodificando módulo: %v", ErrClavePublicaInvalida, err)
	}

	exponent, err := base64.StdEncoding.DecodeString(key.Exponent)
	if err != nil {
		return nil, fmt.Errorf("%w: error decodificando exponente: %v", ErrClavePublicaInvalida, err)
	}

	// Crear clave pública RSA
	publicKey := &rsa.PublicKey{
		N: new(big.Int).SetBytes(modulus),
		E: int(new(big.Int).SetBytes(exponent).Int64()),
	}

	return publicKey, nil
}

// VerifySignature verifica la firma del CAF
func (v *SignatureVerifier) VerifySignature(signedData []byte, signature string) error {
	// Decodificar firma de base64
	sig, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return fmt.Errorf("%w: error decodificando firma: %v", ErrVerificacionFirma, err)
	}

	// Calcular hash SHA1 de los datos firmados
	hashed := sha1.Sum(signedData)

	// Verificar firma PKCS#1 v1.5
	err = rsa.VerifyPKCS1v15(v.publicKey, crypto.SHA1, hashed[:], sig)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrVerificacionFirma, err)
	}

	return nil
}

// ExtractSignedInfo extrae la información firmada del CAF
func ExtractSignedInfo(cafData []byte) (*SignedInfo, error) {
	var info SignedInfo
	if err := xml.Unmarshal(cafData, &info); err != nil {
		return nil, fmt.Errorf("error extrayendo información firmada: %v", err)
	}
	return &info, nil
}

// CanonicalizeXML genera la forma canónica del XML para verificación
func CanonicalizeXML(data []byte) ([]byte, error) {
	// TODO: Implementar canonicalización XML según especificación SII
	// Por ahora retornamos los datos sin procesar
	return data, nil
}
