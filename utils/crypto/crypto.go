package crypto

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"time"
)

// GenerateKeyPair genera un par de claves RSA
func GenerateKeyPair(bits int) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, fmt.Errorf("error al generar clave privada: %w", err)
	}

	return privateKey, &privateKey.PublicKey, nil
}

// SavePrivateKey guarda una clave privada en formato PEM
func SavePrivateKey(key *rsa.PrivateKey, filename string) error {
	keyBytes := x509.MarshalPKCS1PrivateKey(key)
	keyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: keyBytes,
	})

	err := os.WriteFile(filename, keyPEM, 0600)
	if err != nil {
		return fmt.Errorf("error al guardar clave privada: %w", err)
	}

	return nil
}

// SavePublicKey guarda una clave pública en formato PEM
func SavePublicKey(key *rsa.PublicKey, filename string) error {
	keyBytes := x509.MarshalPKCS1PublicKey(key)
	keyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: keyBytes,
	})

	err := os.WriteFile(filename, keyPEM, 0644)
	if err != nil {
		return fmt.Errorf("error al guardar clave pública: %w", err)
	}

	return nil
}

// LoadPrivateKey carga una clave privada desde un archivo PEM
func LoadPrivateKey(filename string) (*rsa.PrivateKey, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error al leer archivo de clave privada: %w", err)
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("error al decodificar PEM")
	}

	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("error al parsear clave privada: %w", err)
	}

	return key, nil
}

// LoadPublicKey carga una clave pública desde un archivo PEM
func LoadPublicKey(filename string) (*rsa.PublicKey, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error al leer archivo de clave pública: %w", err)
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("error al decodificar PEM")
	}

	key, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("error al parsear clave pública: %w", err)
	}

	return key, nil
}

// Sign firma un mensaje usando una clave privada
func Sign(message []byte, privateKey *rsa.PrivateKey) (string, error) {
	hash := sha256.Sum256(message)
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hash[:])
	if err != nil {
		return "", fmt.Errorf("error al firmar mensaje: %w", err)
	}

	return base64.StdEncoding.EncodeToString(signature), nil
}

// Verify verifica una firma usando una clave pública
func Verify(message []byte, signature string, publicKey *rsa.PublicKey) error {
	sig, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return fmt.Errorf("error al decodificar firma: %w", err)
	}

	hash := sha256.Sum256(message)
	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hash[:], sig)
	if err != nil {
		return fmt.Errorf("error al verificar firma: %w", err)
	}

	return nil
}

// Encrypt encripta un mensaje usando una clave pública
func Encrypt(message []byte, publicKey *rsa.PublicKey) (string, error) {
	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, message)
	if err != nil {
		return "", fmt.Errorf("error al encriptar mensaje: %w", err)
	}

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt desencripta un mensaje usando una clave privada
func Decrypt(ciphertext string, privateKey *rsa.PrivateKey) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return nil, fmt.Errorf("error al decodificar mensaje encriptado: %w", err)
	}

	plaintext, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, data)
	if err != nil {
		return nil, fmt.Errorf("error al desencriptar mensaje: %w", err)
	}

	return plaintext, nil
}

// GenerateCertificate genera un certificado X.509
func GenerateCertificate(privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey, subject string) ([]byte, error) {
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, fmt.Errorf("error al generar número de serie: %w", err)
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName: subject,
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(365 * 24 * time.Hour), // 1 año
		KeyUsage:  x509.KeyUsageDigitalSignature,
	}

	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, publicKey, privateKey)
	if err != nil {
		return nil, fmt.Errorf("error al generar certificado: %w", err)
	}

	return certDER, nil
}

// SaveCertificate guarda un certificado en formato PEM
func SaveCertificate(cert []byte, filename string) error {
	certPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert,
	})

	err := os.WriteFile(filename, certPEM, 0644)
	if err != nil {
		return fmt.Errorf("error al guardar certificado: %w", err)
	}

	return nil
}

// LoadCertificate carga un certificado desde un archivo PEM
func LoadCertificate(filename string) (*x509.Certificate, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error al leer archivo de certificado: %w", err)
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("error al decodificar PEM")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("error al parsear certificado: %w", err)
	}

	return cert, nil
}
