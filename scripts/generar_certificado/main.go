package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"time"

	"software.sslmate.com/src/go-pkcs12"
)

func main() {
	// Generar clave privada
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		fmt.Printf("Error generando clave privada: %v\n", err)
		return
	}

	// Crear plantilla de certificado
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName:   "Certificado de Prueba",
			Organization: []string{"Empresa de Prueba"},
			Country:      []string{"CL"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
	}

	// Firmar certificado
	cert, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		fmt.Printf("Error firmando certificado: %v\n", err)
		return
	}

	// Crear directorio si no existe
	dir := "firma_test/1"
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Printf("Error creando directorio: %v\n", err)
		return
	}

	// Convertir el certificado DER a x509.Certificate
	parsedCert, err := x509.ParseCertificate(cert)
	if err != nil {
		fmt.Printf("Error parseando certificado: %v\n", err)
		return
	}

	// Generar archivo PFX
	pfxBytes, err := pkcs12.Encode(rand.Reader, privateKey, parsedCert, []*x509.Certificate{parsedCert}, "123456")
	if err != nil {
		fmt.Printf("Error generando archivo PFX: %v\n", err)
		return
	}

	// Guardar archivo PFX
	pfxPath := filepath.Join(dir, "firma.pfx")
	if err := os.WriteFile(pfxPath, pfxBytes, 0600); err != nil {
		fmt.Printf("Error guardando archivo PFX: %v\n", err)
		return
	}

	fmt.Printf("Certificado de prueba generado exitosamente en: %s\n", pfxPath)
	fmt.Println("Contraseña: 123456")
}
