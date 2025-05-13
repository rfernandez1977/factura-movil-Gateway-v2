package firma

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

	"github.com/cursor/FMgo/config"
	"github.com/cursor/FMgo/models"
)

// Service representa el servicio de firma digital
type Service struct {
	privateKey  *rsa.PrivateKey
	certificate *x509.Certificate
}

// NewService crea una nueva instancia del servicio de firma digital
func NewService(config *config.Config) (*Service, error) {
	// Cargar clave privada
	keyData, err := ioutil.ReadFile(config.KeyPath)
	if err != nil {
		return nil, fmt.Errorf("error al leer clave privada: %w", err)
	}

	keyBlock, _ := pem.Decode(keyData)
	if keyBlock == nil {
		return nil, errors.New("no se pudo decodificar la clave privada")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("error al parsear clave privada: %w", err)
	}

	// Cargar certificado
	certData, err := ioutil.ReadFile(config.CertPath)
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

	return &Service{
		privateKey:  privateKey,
		certificate: cert,
	}, nil
}

// FirmarDTE firma un DTE
func (s *Service) FirmarDTE(dte *models.DTEXMLModel) error {
	if dte == nil {
		return errors.New("el DTE no puede ser nulo")
	}

	// Validar datos requeridos
	if dte.DocumentoXML.Encabezado.IdDoc.TipoDTE == 0 {
		return errors.New("el tipo de DTE es requerido")
	}

	// Calcular hash del documento
	hash := sha256.Sum256([]byte(fmt.Sprintf("%v", dte)))

	// Firmar hash
	signature, err := rsa.SignPKCS1v15(rand.Reader, s.privateKey, crypto.SHA256, hash[:])
	if err != nil {
		return fmt.Errorf("error al firmar DTE: %w", err)
	}

	// Asignar firma
	dte.Signature = base64.StdEncoding.EncodeToString(signature)

	return nil
}

// GenerarTED genera el Timbre Electrónico del Documento
func (s *Service) GenerarTED(dte *models.DTEXMLModel) (string, error) {
	if dte == nil {
		return "", errors.New("el DTE no puede ser nulo")
	}

	// Validar datos requeridos
	if dte.DocumentoXML.Encabezado.IdDoc.TipoDTE == 0 {
		return "", errors.New("el tipo de DTE es requerido")
	}

	// Calcular hash del documento
	hash := sha256.Sum256([]byte(fmt.Sprintf("%v", dte)))

	// Firmar hash
	signature, err := rsa.SignPKCS1v15(rand.Reader, s.privateKey, crypto.SHA256, hash[:])
	if err != nil {
		return "", fmt.Errorf("error al generar TED: %w", err)
	}

	// Generar TED
	ted := base64.StdEncoding.EncodeToString(signature)

	return ted, nil
}

// FirmarSobre firma un sobre de DTE
func (s *Service) FirmarSobre(sobre *models.SobreDTEModel) error {
	if sobre == nil {
		return errors.New("el sobre no puede ser nulo")
	}

	if len(sobre.Documentos) == 0 {
		return errors.New("el sobre debe contener al menos un documento")
	}

	// Validar datos requeridos
	if sobre.Caratula.RutEmisor == "" {
		return errors.New("el RUT del emisor es requerido")
	}

	// Calcular hash del sobre
	hash := sha256.Sum256([]byte(fmt.Sprintf("%v", sobre)))

	// Firmar hash
	signature, err := rsa.SignPKCS1v15(rand.Reader, s.privateKey, crypto.SHA256, hash[:])
	if err != nil {
		return fmt.Errorf("error al firmar sobre: %w", err)
	}

	// Asignar firma
	sobre.Signature = base64.StdEncoding.EncodeToString(signature)

	return nil
}
