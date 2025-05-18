package firma

import (
	"crypto/rsa"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"

	"software.sslmate.com/src/go-pkcs12"
)

// BaseFirmaService implementa la interfaz FirmaService
type BaseFirmaService struct {
	cert       *x509.Certificate
	privateKey *rsa.PrivateKey
	config     *ConfiguracionFirma
}

// NewBaseFirmaService crea una nueva instancia de BaseFirmaService
func NewBaseFirmaService(config *ConfiguracionFirma) (*BaseFirmaService, error) {
	if config == nil {
		return nil, errors.New("la configuración no puede ser nil")
	}

	service := &BaseFirmaService{
		config: config,
	}

	if err := service.cargarCertificado(); err != nil {
		return nil, fmt.Errorf("error cargando certificado: %w", err)
	}

	return service, nil
}

// cargarCertificado carga el certificado y la llave privada
func (s *BaseFirmaService) cargarCertificado() error {
	// Leer archivo de certificado
	certData, err := ioutil.ReadFile(s.config.RutaCertificado)
	if err != nil {
		return fmt.Errorf("error leyendo certificado: %w", err)
	}

	// Si es un archivo PFX
	if s.config.Password != "" {
		privateKey, cert, err := pkcs12.Decode(certData, s.config.Password)
		if err != nil {
			return fmt.Errorf("error decodificando PFX: %w", err)
		}
		s.privateKey = privateKey.(*rsa.PrivateKey)
		s.cert = cert
		return nil
	}

	// Si es un archivo PEM
	if s.config.RutaLlave == "" {
		return errors.New("se requiere la ruta de la llave privada para certificados PEM")
	}

	// TODO: Implementar carga de certificados PEM
	return errors.New("formato PEM no implementado aún")
}

// FirmarDocumento implementa la interfaz FirmaService
func (s *BaseFirmaService) FirmarDocumento(xml string) (*ResultadoFirma, error) {
	if xml == "" {
		return nil, errors.New("el XML no puede estar vacío")
	}

	// TODO: Implementar firma de documento
	return nil, errors.New("método no implementado")
}

// ValidarFirma implementa la interfaz FirmaService
func (s *BaseFirmaService) ValidarFirma(xml string) (*EstadoFirma, error) {
	if xml == "" {
		return nil, errors.New("el XML no puede estar vacío")
	}

	// TODO: Implementar validación de firma
	return nil, errors.New("método no implementado")
}

// ObtenerCertificado implementa la interfaz FirmaService
func (s *BaseFirmaService) ObtenerCertificado() (*x509.Certificate, error) {
	if s.cert == nil {
		return nil, errors.New("certificado no inicializado")
	}
	return s.cert, nil
}
