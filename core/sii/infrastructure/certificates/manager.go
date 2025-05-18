package certificates

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"time"

	"software.sslmate.com/src/go-pkcs12"
)

// CertificateInfo contiene información sobre un certificado
type CertificateInfo struct {
	Subject      string    `json:"subject"`
	Issuer       string    `json:"issuer"`
	NotBefore    time.Time `json:"not_before"`
	NotAfter     time.Time `json:"not_after"`
	SerialNumber string    `json:"serial_number"`
}

// CertManager maneja los certificados digitales
type CertManager struct {
	certPath string
	keyPath  string
	cert     *x509.Certificate
}

// NewCertManager crea una nueva instancia de CertManager
func NewCertManager(certPath, keyPath string) (*CertManager, error) {
	manager := &CertManager{
		certPath: certPath,
		keyPath:  keyPath,
	}

	if err := manager.loadCertificate(); err != nil {
		return nil, err
	}

	return manager, nil
}

// loadCertificate carga el certificado desde los archivos
func (m *CertManager) loadCertificate() error {
	cert, err := tls.LoadX509KeyPair(m.certPath, m.keyPath)
	if err != nil {
		return fmt.Errorf("error cargando certificado: %w", err)
	}

	x509Cert, err := x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		return fmt.Errorf("error parseando certificado: %w", err)
	}

	m.cert = x509Cert
	return nil
}

// ValidateCertificate valida que el certificado sea válido
func (m *CertManager) ValidateCertificate() error {
	now := time.Now()

	if now.Before(m.cert.NotBefore) {
		return fmt.Errorf("certificado aún no es válido")
	}

	if now.After(m.cert.NotAfter) {
		return fmt.Errorf("certificado expirado")
	}

	return nil
}

// GetCertificateInfo retorna información sobre el certificado
func (m *CertManager) GetCertificateInfo() *CertificateInfo {
	return &CertificateInfo{
		Subject:      m.cert.Subject.String(),
		Issuer:       m.cert.Issuer.String(),
		NotBefore:    m.cert.NotBefore,
		NotAfter:     m.cert.NotAfter,
		SerialNumber: m.cert.SerialNumber.String(),
	}
}

// IsExpiringSoon verifica si el certificado está por expirar
func (m *CertManager) IsExpiringSoon(daysThreshold int) bool {
	expirationDate := m.cert.NotAfter
	threshold := time.Now().AddDate(0, 0, daysThreshold)
	return threshold.After(expirationDate)
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
