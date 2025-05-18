package test_data

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"time"
	
	"software.sslmate.com/src/go-pkcs12"
)

// GenerarCertificadoPrueba genera un certificado X.509 y una llave privada RSA para pruebas
func GenerarCertificadoPrueba(rutPath, keyPath string) error {
	// Generar llave privada
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	// Crear plantilla del certificado
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Empresa de Prueba SpA"},
			Country:     []string{"CL"},
			CommonName:  "Test Certificate",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(24 * time.Hour),
		KeyUsage:             x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:          []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	// Crear certificado
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return err
	}

	// Guardar certificado
	certOut, err := os.Create(rutPath)
	if err != nil {
		return err
	}
	defer certOut.Close()

	err = pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	if err != nil {
		return err
	}

	// Guardar llave privada
	keyOut, err := os.Create(keyPath)
	if err != nil {
		return err
	}
	defer keyOut.Close()

	privBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	err = pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: privBytes})
	if err != nil {
		return err
	}

	return nil
}

// GenerarPFX12 genera un archivo PKCS12 (.pfx) para pruebas
func GenerarPFX12(pfxPath string, password string) error {
	// Generar llave privada
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	// Crear plantilla del certificado
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Empresa de Prueba SpA"},
			Country:     []string{"CL"},
			CommonName:  "Test Certificate",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(24 * time.Hour),
		KeyUsage:             x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:          []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	// Crear certificado
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return err
	}

	// Convertir a PKCS12
	pfxData, err := pkcs12.Encode(rand.Reader, privateKey, &template, nil, password)
	if err != nil {
		return err
	}

	// Guardar archivo PFX
	return os.WriteFile(pfxPath, pfxData, 0644)
} 