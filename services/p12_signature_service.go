package services

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"time"
)

// P12SignatureService maneja la firma digital de documentos usando certificados P12
type P12SignatureService struct {
	certificate *x509.Certificate
	privateKey  *rsa.PrivateKey
}

// NewP12SignatureService crea una nueva instancia del servicio de firma P12
func NewP12SignatureService(certPath, keyPath string) (*P12SignatureService, error) {
	// Leer el certificado
	certPEM, err := ioutil.ReadFile(certPath)
	if err != nil {
		return nil, fmt.Errorf("error leyendo certificado: %v", err)
	}

	// Decodificar el certificado
	certBlock, _ := pem.Decode(certPEM)
	if certBlock == nil {
		return nil, errors.New("error decodificando certificado PEM")
	}

	cert, err := x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("error parseando certificado: %v", err)
	}

	// Leer la clave privada
	keyPEM, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("error leyendo clave privada: %v", err)
	}

	// Decodificar la clave privada
	keyBlock, _ := pem.Decode(keyPEM)
	if keyBlock == nil {
		return nil, errors.New("error decodificando clave privada PEM")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("error parseando clave privada: %v", err)
	}

	return &P12SignatureService{
		certificate: cert,
		privateKey:  privateKey,
	}, nil
}

// Sign firma un documento usando el certificado P12
func (s *P12SignatureService) Sign(data []byte) ([]byte, error) {
	// Calcular el hash del documento
	hash := sha256.Sum256(data)

	// Firmar el hash
	signature, err := rsa.SignPKCS1v15(rand.Reader, s.privateKey, crypto.SHA256, hash[:])
	if err != nil {
		return nil, fmt.Errorf("error firmando documento: %v", err)
	}

	return signature, nil
}

// Verify verifica la firma de un documento
func (s *P12SignatureService) Verify(data, signature []byte) error {
	// Calcular el hash del documento
	hash := sha256.Sum256(data)

	// Verificar la firma
	err := rsa.VerifyPKCS1v15(&s.certificate.PublicKey, crypto.SHA256, hash[:], signature)
	if err != nil {
		return fmt.Errorf("error verificando firma: %v", err)
	}

	return nil
}

// GetCertificateInfo retorna información del certificado
func (s *P12SignatureService) GetCertificateInfo() map[string]interface{} {
	return map[string]interface{}{
		"subject":     s.certificate.Subject,
		"issuer":      s.certificate.Issuer,
		"valid_from":  s.certificate.NotBefore,
		"valid_until": s.certificate.NotAfter,
		"serial":      s.certificate.SerialNumber.String(),
		"fingerprint": fmt.Sprintf("%x", s.certificate.Signature),
	}
}

// IsValid verifica si el certificado es válido
func (s *P12SignatureService) IsValid() bool {
	now := time.Now()
	return now.After(s.certificate.NotBefore) && now.Before(s.certificate.NotAfter)
}
