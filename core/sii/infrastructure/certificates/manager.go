package certificates

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"time"

	"software.sslmate.com/src/go-pkcs12"
)

// CertificateInfo contiene la información extraída del certificado
type CertificateInfo struct {
	Subject      string
	Issuer       string
	SerialNumber string
	ValidFrom    time.Time
	ValidUntil   time.Time
	RutTitular   string
}

// CertificateManager maneja las operaciones con certificados digitales
type CertificateManager struct {
	certPath     string
	certPassword string
	certInfo     *CertificateInfo
	tlsConfig    *tls.Config
}

// NewCertificateManager crea una nueva instancia del gestor de certificados
func NewCertificateManager(certPath, certPassword string) (*CertificateManager, error) {
	manager := &CertificateManager{
		certPath:     certPath,
		certPassword: certPassword,
	}

	if err := manager.loadCertificate(); err != nil {
		return nil, fmt.Errorf("error al cargar certificado: %w", err)
	}

	return manager, nil
}

// loadCertificate carga y valida el certificado
func (m *CertificateManager) loadCertificate() error {
	certData, err := ioutil.ReadFile(m.certPath)
	if err != nil {
		return fmt.Errorf("error al leer archivo de certificado: %w", err)
	}

	// Decodificar el certificado
	block, _ := pem.Decode(certData)
	if block == nil {
		return fmt.Errorf("error al decodificar certificado PEM")
	}

	// Parsear el certificado
	cert, err := x509.ParsePKCS12(block.Bytes, []byte(m.certPassword))
	if err != nil {
		return fmt.Errorf("error al parsear certificado PKCS12: %w", err)
	}

	// Extraer información del certificado
	m.certInfo = &CertificateInfo{
		Subject:      cert.Subject.String(),
		Issuer:       cert.Issuer.String(),
		SerialNumber: cert.SerialNumber.String(),
		ValidFrom:    cert.NotBefore,
		ValidUntil:   cert.NotAfter,
		RutTitular:   extractRutFromSubject(cert.Subject.String()),
	}

	// Configurar TLS
	tlsCert, err := tls.LoadX509KeyPair(m.certPath, m.certPath)
	if err != nil {
		return fmt.Errorf("error al cargar par de claves: %w", err)
	}

	m.tlsConfig = &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		MinVersion:   tls.VersionTLS12,
	}

	return nil
}

// GetTLSConfig retorna la configuración TLS
func (m *CertificateManager) GetTLSConfig() *tls.Config {
	return m.tlsConfig
}

// GetCertificateInfo retorna la información del certificado
func (m *CertificateManager) GetCertificateInfo() *CertificateInfo {
	return m.certInfo
}

// ValidateCertificate valida el certificado actual
func (m *CertificateManager) ValidateCertificate() error {
	if m.certInfo == nil {
		return fmt.Errorf("certificado no cargado")
	}

	now := time.Now()
	if now.Before(m.certInfo.ValidFrom) {
		return fmt.Errorf("certificado aún no es válido")
	}
	if now.After(m.certInfo.ValidUntil) {
		return fmt.Errorf("certificado expirado")
	}

	return nil
}

// IsExpiringSoon verifica si el certificado está por expirar
func (m *CertificateManager) IsExpiringSoon(daysThreshold int) bool {
	if m.certInfo == nil {
		return false
	}

	expirationWarningDate := time.Now().AddDate(0, 0, daysThreshold)
	return m.certInfo.ValidUntil.Before(expirationWarningDate)
}

// extractRutFromSubject extrae el RUT del subject del certificado
func extractRutFromSubject(subject string) string {
	// TODO: Implementar extracción de RUT según formato del certificado
	return ""
}

// Manager maneja las operaciones con certificados
type Manager struct {
	certPath string
	keyPath  string
}

// NewManager crea una nueva instancia del administrador de certificados
func NewManager(certPath, keyPath string) *Manager {
	return &Manager{
		certPath: certPath,
		keyPath:  keyPath,
	}
}

// LoadCertificate carga un certificado desde un archivo PKCS12
func (m *Manager) LoadCertificate(password string) (*x509.Certificate, interface{}, error) {
	data, err := ioutil.ReadFile(m.certPath)
	if err != nil {
		return nil, nil, fmt.Errorf("error leyendo archivo de certificado: %v", err)
	}

	// Usar el paquete go-pkcs12 para parsear el certificado
	privateKey, cert, err := pkcs12.Decode(data, password)
	if err != nil {
		return nil, nil, fmt.Errorf("error decodificando PKCS12: %v", err)
	}

	return cert, privateKey, nil
}
