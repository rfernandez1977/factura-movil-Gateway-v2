package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"strings"
)

// EncryptService maneja el cifrado y descifrado de datos
type EncryptService struct {
	key []byte
}

// NewEncryptService crea una nueva instancia del servicio de cifrado
func NewEncryptService(key string) (*EncryptService, error) {
	if len(key) != 32 {
		return nil, errors.New("la clave debe tener 32 bytes")
	}
	return &EncryptService{key: []byte(key)}, nil
}

// Encrypt cifra un texto plano
func (s *EncryptService) Encrypt(plaintext string) (string, error) {
	// Crear cipher block
	block, err := aes.NewCipher(s.key)
	if err != nil {
		return "", err
	}

	// Crear GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Crear nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// Cifrar texto
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt descifra un texto cifrado
func (s *EncryptService) Decrypt(encrypted string) (string, error) {
	// Decodificar base64
	ciphertext, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", err
	}

	// Crear cipher block
	block, err := aes.NewCipher(s.key)
	if err != nil {
		return "", err
	}

	// Crear GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Extraer nonce
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", errors.New("texto cifrado demasiado corto")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// Descifrar texto
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// HashRUT genera un hash seguro de un RUT
func (s *EncryptService) HashRUT(rut string) (string, error) {
	// Limpiar y validar RUT
	rut = CleanRUT(rut)
	if err := ValidateRUT(rut); err != nil {
		return "", err
	}

	// Cifrar RUT
	return s.Encrypt(rut)
}

// MaskRUT enmascara un RUT para mostrar
func MaskRUT(rut string) string {
	if len(rut) < 4 {
		return rut
	}
	return "***" + rut[len(rut)-4:]
}

// CleanRUT elimina puntos y guiones de un RUT
func CleanRUT(rut string) string {
	rut = strings.ReplaceAll(rut, ".", "")
	rut = strings.ReplaceAll(rut, "-", "")
	return rut
}
