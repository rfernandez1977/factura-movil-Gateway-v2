package services

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"time"

	"FMgo/config"
	"FMgo/models"
	"github.com/supabase-community/postgrest-go"
	"software.sslmate.com/src/go-pkcs12"
)

// CertificadoService maneja la lógica de negocio de certificados
type CertificadoService struct {
	config   *config.SupabaseConfig
	certPath string
	keyPath  string
}

// NewCertificadoService crea una nueva instancia del servicio de certificado
func NewCertificadoService(config *config.SupabaseConfig, certPath, keyPath string) *CertificadoService {
	return &CertificadoService{
		config:   config,
		certPath: certPath,
		keyPath:  keyPath,
	}
}

// GetCertificadoByEmpresaID obtiene el certificado de una empresa
func (s *CertificadoService) GetCertificadoByEmpresaID(empresaID string) (*models.CertificadoDigital, error) {
	var certificado models.CertificadoDigital
	client := s.config.Client.(*postgrest.Client)

	resp, _, err := client.From("certificados_digitales").
		Select("*", "", false).
		Eq("empresa_id", empresaID).
		Single().
		Execute()

	if err != nil {
		return nil, fmt.Errorf("error al obtener certificado: %v", err)
	}

	if err := json.Unmarshal(resp, &certificado); err != nil {
		return nil, fmt.Errorf("error al decodificar certificado: %v", err)
	}

	return &certificado, nil
}

// CrearCertificado crea un nuevo certificado
func (s *CertificadoService) CrearCertificado(certificado *models.CertificadoDigital) (*models.CertificadoDigital, error) {
	// Validar certificado
	if err := s.validarCertificado(certificado); err != nil {
		return nil, err
	}

	certificado.CreatedAt = time.Now()
	certificado.UpdatedAt = time.Now()

	client := s.config.Client.(*postgrest.Client)

	resp, _, err := client.From("certificados_digitales").
		Insert(certificado, false, "", "", "").
		Execute()

	if err != nil {
		return nil, fmt.Errorf("error al guardar certificado: %v", err)
	}

	var nuevosCertificados []models.CertificadoDigital
	if err := json.Unmarshal(resp, &nuevosCertificados); err != nil {
		return nil, fmt.Errorf("error al decodificar respuesta: %v", err)
	}

	if len(nuevosCertificados) == 0 {
		return nil, fmt.Errorf("no se creó el certificado")
	}

	return &nuevosCertificados[0], nil
}

// ActualizarCertificado actualiza un certificado existente
func (s *CertificadoService) ActualizarCertificado(certificado *models.CertificadoDigital) error {
	// Validar certificado
	if err := s.validarCertificado(certificado); err != nil {
		return err
	}

	certificado.UpdatedAt = time.Now()

	client := s.config.Client.(*postgrest.Client)

	_, _, err := client.From("certificados_digitales").
		Update(certificado, "", "").
		Eq("id", certificado.ID).
		Execute()

	if err != nil {
		return fmt.Errorf("error al actualizar certificado: %v", err)
	}

	return nil
}

// EliminarCertificado elimina un certificado
func (s *CertificadoService) EliminarCertificado(id string) error {
	client := s.config.Client.(*postgrest.Client)

	_, _, err := client.From("certificados_digitales").
		Delete("", "").
		Eq("id", id).
		Execute()

	if err != nil {
		return fmt.Errorf("error al eliminar certificado: %v", err)
	}

	return nil
}

// validarCertificado valida un certificado antes de crearlo o actualizarlo
func (s *CertificadoService) validarCertificado(certificado *models.CertificadoDigital) error {
	if certificado.EmpresaID == "" {
		return fmt.Errorf("ID de empresa requerido")
	}
	if certificado.Certificate == "" {
		return fmt.Errorf("certificado requerido")
	}
	if certificado.PrivateKey == "" {
		return fmt.Errorf("llave privada requerida")
	}
	if certificado.ValidTo.Before(time.Now()) {
		return fmt.Errorf("certificado vencido")
	}
	return nil
}

// GenerarCertificado genera un nuevo certificado autofirmado
func (c *CertificadoService) GenerarCertificado(organizacion, unidad, pais, provincia, localidad string, diasValidez int) error {
	// Generar clave privada
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("error generando clave privada: %v", err)
	}

	// Crear plantilla de certificado
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization:       []string{organizacion},
			OrganizationalUnit: []string{unidad},
			Country:            []string{pais},
			Province:           []string{provincia},
			Locality:           []string{localidad},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(0, 0, diasValidez),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  false,
	}

	// Crear certificado
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return fmt.Errorf("error creando certificado: %v", err)
	}

	// Guardar certificado
	certOut, err := os.Create(c.certPath)
	if err != nil {
		return fmt.Errorf("error creando archivo de certificado: %v", err)
	}
	defer certOut.Close()

	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		return fmt.Errorf("error guardando certificado: %v", err)
	}

	// Guardar clave privada
	keyOut, err := os.OpenFile(c.keyPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("error creando archivo de clave privada: %v", err)
	}
	defer keyOut.Close()

	privBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	if err := pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: privBytes}); err != nil {
		return fmt.Errorf("error guardando clave privada: %v", err)
	}

	return nil
}

// RenovarCertificado renueva un certificado existente
func (c *CertificadoService) RenovarCertificado(diasValidez int) error {
	// Cargar certificado existente
	certData, err := ioutil.ReadFile(c.certPath)
	if err != nil {
		return fmt.Errorf("error leyendo certificado: %v", err)
	}

	block, _ := pem.Decode(certData)
	if block == nil {
		return fmt.Errorf("error decodificando certificado PEM")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return fmt.Errorf("error parseando certificado: %v", err)
	}

	// Cargar clave privada
	keyData, err := ioutil.ReadFile(c.keyPath)
	if err != nil {
		return fmt.Errorf("error leyendo clave privada: %v", err)
	}

	block, _ = pem.Decode(keyData)
	if block == nil {
		return fmt.Errorf("error decodificando clave privada PEM")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("error parseando clave privada: %v", err)
	}

	// Crear nuevo certificado con los mismos datos
	template := x509.Certificate{
		SerialNumber:          big.NewInt(cert.SerialNumber.Int64() + 1),
		Subject:               cert.Subject,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(0, 0, diasValidez),
		KeyUsage:              cert.KeyUsage,
		ExtKeyUsage:           cert.ExtKeyUsage,
		BasicConstraintsValid: cert.BasicConstraintsValid,
		IsCA:                  cert.IsCA,
	}

	// Crear certificado
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, cert, &privateKey.PublicKey, privateKey)
	if err != nil {
		return fmt.Errorf("error creando certificado: %v", err)
	}

	// Guardar certificado
	certOut, err := os.Create(c.certPath)
	if err != nil {
		return fmt.Errorf("error creando archivo de certificado: %v", err)
	}
	defer certOut.Close()

	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		return fmt.Errorf("error guardando certificado: %v", err)
	}

	return nil
}

// VerificarCertificado verifica la validez de un certificado
func (c *CertificadoService) VerificarCertificado() error {
	// Cargar certificado
	certData, err := ioutil.ReadFile(c.certPath)
	if err != nil {
		return fmt.Errorf("error leyendo certificado: %v", err)
	}

	block, _ := pem.Decode(certData)
	if block == nil {
		return fmt.Errorf("error decodificando certificado PEM")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return fmt.Errorf("error parseando certificado: %v", err)
	}

	// Verificar fechas
	if time.Now().After(cert.NotAfter) {
		return fmt.Errorf("certificado expirado")
	}

	if time.Now().Before(cert.NotBefore) {
		return fmt.Errorf("certificado no válido aún")
	}

	// Verificar uso
	if cert.KeyUsage&x509.KeyUsageDigitalSignature == 0 {
		return fmt.Errorf("certificado no tiene permiso para firma digital")
	}

	// Verificar clave pública
	if _, ok := cert.PublicKey.(*rsa.PublicKey); !ok {
		return fmt.Errorf("certificado no contiene una clave pública RSA")
	}

	return nil
}

// ExportarCertificado exporta el certificado en formato PEM
func (c *CertificadoService) ExportarCertificado() (string, error) {
	certData, err := ioutil.ReadFile(c.certPath)
	if err != nil {
		return "", fmt.Errorf("error leyendo certificado: %v", err)
	}

	return string(certData), nil
}

// ExportarClavePublica exporta la clave pública en formato PEM
func (c *CertificadoService) ExportarClavePublica() (string, error) {
	// Cargar certificado
	certData, err := ioutil.ReadFile(c.certPath)
	if err != nil {
		return "", fmt.Errorf("error leyendo certificado: %v", err)
	}

	block, _ := pem.Decode(certData)
	if block == nil {
		return "", fmt.Errorf("error decodificando certificado PEM")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("error parseando certificado: %v", err)
	}

	// Convertir clave pública a PEM
	pubKey, ok := cert.PublicKey.(*rsa.PublicKey)
	if !ok {
		return "", fmt.Errorf("certificado no contiene una clave pública RSA")
	}

	pubKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(pubKey),
	})

	return string(pubKeyPEM), nil
}

// ImportarCertificado importa un certificado desde un archivo PEM
func (c *CertificadoService) ImportarCertificado(certPEM string) error {
	// Decodificar certificado
	block, _ := pem.Decode([]byte(certPEM))
	if block == nil {
		return fmt.Errorf("error decodificando certificado PEM")
	}

	// Parsear y verificar certificado
	if _, err := x509.ParseCertificate(block.Bytes); err != nil {
		return fmt.Errorf("error parseando certificado: %v", err)
	}

	// Verificar certificado
	if err := c.VerificarCertificado(); err != nil {
		return fmt.Errorf("error verificando certificado: %v", err)
	}

	// Guardar certificado
	err := ioutil.WriteFile(c.certPath, []byte(certPEM), 0644)
	if err != nil {
		return fmt.Errorf("error guardando certificado: %v", err)
	}

	return nil
}

// ObtenerInformacionCertificado retorna información detallada del certificado
func (c *CertificadoService) ObtenerInformacionCertificado() (*CertificadoInfo, error) {
	// Cargar certificado
	certData, err := ioutil.ReadFile(c.certPath)
	if err != nil {
		return nil, fmt.Errorf("error leyendo certificado: %v", err)
	}

	block, _ := pem.Decode(certData)
	if block == nil {
		return nil, fmt.Errorf("error decodificando certificado PEM")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("error parseando certificado: %v", err)
	}

	info := &CertificadoInfo{
		SerialNumber:     cert.SerialNumber.String(),
		Subject:          cert.Subject.String(),
		Issuer:           cert.Issuer.String(),
		NotBefore:        cert.NotBefore,
		NotAfter:         cert.NotAfter,
		KeyUsage:         cert.KeyUsage,
		ExtKeyUsage:      cert.ExtKeyUsage,
		IsCA:             cert.IsCA,
		MaxPathLen:       cert.MaxPathLen,
		MaxPathLenZero:   cert.MaxPathLenZero,
		BasicConstraints: cert.BasicConstraintsValid,
	}

	return info, nil
}

// CertificadoInfo contiene información detallada de un certificado
type CertificadoInfo struct {
	SerialNumber     string
	Subject          string
	Issuer           string
	NotBefore        time.Time
	NotAfter         time.Time
	KeyUsage         x509.KeyUsage
	ExtKeyUsage      []x509.ExtKeyUsage
	IsCA             bool
	MaxPathLen       int
	MaxPathLenZero   bool
	BasicConstraints bool
}

// ExtraerInfoCertificado extrae la información de un certificado PFX
func (s *CertificadoService) ExtraerInfoCertificado(pfxData []byte, password string) (*models.CertificadoDigital, error) {
	// Decodificar el archivo PFX
	privateKey, certificate, err := pkcs12.Decode(pfxData, password)
	if err != nil {
		return nil, fmt.Errorf("error al decodificar PFX: %v", err)
	}

	// Extraer información del certificado
	cert := &models.CertificadoDigital{
		SerialNumber: certificate.SerialNumber.String(),
		Issuer:       certificate.Issuer.CommonName,
		Subject:      certificate.Subject.CommonName,
		ValidFrom:    certificate.NotBefore,
		ValidTo:      certificate.NotAfter,
	}

	// Convertir la llave privada a PEM
	privateKeyPEM := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey.(*rsa.PrivateKey)),
	}
	cert.PrivateKey = string(pem.EncodeToMemory(privateKeyPEM))

	// Convertir el certificado a PEM
	certificatePEM := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certificate.Raw,
	}
	cert.Certificate = string(pem.EncodeToMemory(certificatePEM))

	return cert, nil
}

// ValidarCertificado valida que un certificado sea válido
func (s *CertificadoService) ValidarCertificado(cert *models.CertificadoDigital) error {
	// Verificar fechas de validez
	now := time.Now()
	if now.Before(cert.ValidFrom) {
		return fmt.Errorf("el certificado aún no es válido (válido desde %s)", cert.ValidFrom)
	}
	if now.After(cert.ValidTo) {
		return fmt.Errorf("el certificado ha expirado (expiró el %s)", cert.ValidTo)
	}

	// Verificar que el certificado y la llave privada sean válidos
	_, err := tls.X509KeyPair([]byte(cert.Certificate), []byte(cert.PrivateKey))
	if err != nil {
		return fmt.Errorf("error al validar certificado y llave privada: %v", err)
	}

	return nil
}
