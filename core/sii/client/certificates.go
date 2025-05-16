package client

import (
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"time"

	"github.com/fmgo/core/sii/models"
)

// DatosFirmante contiene la información del titular del certificado
type DatosFirmante struct {
	RUT         string
	RazonSocial string
	Email       string
}

// CertificateManager gestiona los certificados digitales para el SII
type CertificateManager interface {
	CargarCertificado(path string, password string) (*x509.Certificate, error)
	ValidarCertificado(cert *x509.Certificate) error
	ObtenerDatosFirmante(cert *x509.Certificate) (*DatosFirmante, error)
	RenovarCertificado(cert *x509.Certificate) error
}

// DefaultCertManager implementa CertificateManager
type DefaultCertManager struct{}

// NewCertificateManager crea una nueva instancia del gestor de certificados
func NewCertificateManager() CertificateManager {
	return &DefaultCertManager{}
}

// CargarCertificado carga un certificado desde un archivo
func (m *DefaultCertManager) CargarCertificado(path string, password string) (*x509.Certificate, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, models.NewSIIError(models.ErrCertInvalid, "Error al leer certificado", err)
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, models.NewSIIError(models.ErrCertInvalid, "Formato de certificado inválido", nil)
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, models.NewSIIError(models.ErrCertInvalid, "Error al parsear certificado", err)
	}

	return cert, nil
}

// ValidarCertificado verifica la validez del certificado
func (m *DefaultCertManager) ValidarCertificado(cert *x509.Certificate) error {
	now := time.Now()

	if now.Before(cert.NotBefore) {
		return models.NewSIIError(models.ErrCertInvalid, "Certificado aún no es válido", nil)
	}

	if now.After(cert.NotAfter) {
		return models.NewSIIError(models.ErrCertInvalid, "Certificado expirado", nil)
	}

	// Aquí se podrían agregar más validaciones específicas del SII
	return nil
}

// ObtenerDatosFirmante extrae la información del titular del certificado
func (m *DefaultCertManager) ObtenerDatosFirmante(cert *x509.Certificate) (*DatosFirmante, error) {
	if cert == nil {
		return nil, models.NewSIIError(models.ErrCertInvalid, "Certificado no proporcionado", nil)
	}

	// Extraer información del Subject del certificado
	// Nota: Esto es un ejemplo, la estructura exacta dependerá del formato de los certificados del SII
	datos := &DatosFirmante{
		RazonSocial: cert.Subject.CommonName,
		Email:       cert.EmailAddresses[0],
	}

	// Extraer RUT del certificado (implementación específica según formato SII)
	// Esta es una implementación de ejemplo
	for _, v := range cert.Subject.Organization {
		if len(v) > 4 && v[:4] == "RUT:" {
			datos.RUT = v[4:]
			break
		}
	}

	if datos.RUT == "" {
		return nil, models.NewSIIError(models.ErrCertInvalid, "No se pudo extraer RUT del certificado", nil)
	}

	return datos, nil
}

// RenovarCertificado renueva un certificado existente
func (m *DefaultCertManager) RenovarCertificado(cert *x509.Certificate) error {
	// Esta función debería implementar la lógica de renovación específica del SII
	// Por ahora retornamos un error indicando que no está implementado
	return models.NewSIIError(models.ErrCertInvalid, "Renovación de certificados no implementada", nil)
}
