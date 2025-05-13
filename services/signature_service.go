package services

import (
    "crypto"
    "crypto/rand"
    "crypto/rsa"
    "crypto/sha256"
    "crypto/x509"
    "encoding/base64"
    "encoding/pem"
    "errors"
    "fmt"
    "io/ioutil"
)

// SignatureService proporciona métodos para firmar documentos electrónicos
type SignatureService struct {
    privateKey *rsa.PrivateKey
    certificate *x509.Certificate
}

// NewSignatureService crea una nueva instancia del servicio de firma
func NewSignatureService(keyPath, certPath string) (*SignatureService, error) {
    // Cargar clave privada
    keyData, err := ioutil.ReadFile(keyPath)
    if err != nil {
        return nil, fmt.Errorf("error al leer clave privada: %w", err)
    }
    
    keyBlock, _ := pem.Decode(keyData)
    if keyBlock == nil {
        return nil, errors.New("no se pudo decodificar la clave privada")
    }
    
    privateKey, err := x509.ParsePKCS8PrivateKey(keyBlock.Bytes)
    if err != nil {
        return nil, fmt.Errorf("error al parsear clave privada: %w", err)
    }
    
    rsaKey, ok := privateKey.(*rsa.PrivateKey)
    if !ok {
        return nil, errors.New("la clave privada no es de tipo RSA")
    }
    
    // Cargar certificado
    certData, err := ioutil.ReadFile(certPath)
    if err != nil {
        return nil, fmt.Errorf("error al leer certificado: %w", err)
    }
    
    certBlock, _ := pem.Decode(certData)
    if certBlock == nil {
        return nil, errors.New("no se pudo decodificar el certificado")
    }
    
    cert, err := x509.ParseCertificate(certBlock.Bytes)
    if err != nil {
        return nil, fmt.Errorf("error al parsear certificado: %w", err)
    }
    
    return &SignatureService{
        privateKey: rsaKey,
        certificate: cert,
    }, nil
}

// FirmarDocumento firma un documento XML
func (s *SignatureService) FirmarDocumento(documento []byte) ([]byte, error) {
    // Calcular hash SHA-256 del documento
    hash := sha256.Sum256(documento)
    
    // Firmar el hash con la clave privada
    signature, err := rsa.SignPKCS1v15(rand.Reader, s.privateKey, crypto.SHA256, hash[:])
    if err != nil {
        return nil, fmt.Errorf("error al firmar documento: %w", err)
    }
    
    // Codificar firma en base64
    signatureBase64 := base64.StdEncoding.EncodeToString(signature)
    
    // Aquí se debería insertar la firma en el documento XML
    // Esta es una implementación simplificada
    
    return documento, nil
}

// VerificarFirma verifica la firma de un documento
func (s *SignatureService) VerificarFirma(documento []byte, firma []byte) (bool, error) {
    // Calcular hash SHA-256 del documento
    hash := sha256.Sum256(documento)
    
    // Verificar firma
    err := rsa.VerifyPKCS1v15(&s.privateKey.PublicKey, crypto.SHA256, hash[:], firma)
    if err != nil {
        return false, fmt.Errorf("firma inválida: %w", err)
    }
    
    return true, nil
}

// ObtenerCertificadoBase64 devuelve el certificado en formato base64
func (s *SignatureService) ObtenerCertificadoBase64() string {
    return base64.StdEncoding.EncodeToString(s.certificate.Raw)
}

// ValidarCertificado valida el certificado
func (s *SignatureService) ValidarCertificado() error {
    // Verificar fecha de validez
    now := time.Now()
    if now.Before(s.certificate.NotBefore) || now.After(s.certificate.NotAfter) {
        return fmt.Errorf("certificado fuera de fecha de validez (válido desde %s hasta %s)",
            s.certificate.NotBefore.Format("2006-01-02"),
            s.certificate.NotAfter.Format("2006-01-02"))
    }
    
    // Aquí se podrían agregar más validaciones
    
    return nil
}