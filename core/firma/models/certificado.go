package models

import (
	"crypto/x509"
	"time"
)

// Certificado representa un certificado digital
type Certificado struct {
	ID              string    `json:"id"`
	RutEmpresa      string    `json:"rut_empresa"`
	NombreEmpresa   string    `json:"nombre_empresa"`
	CertificadoPEM  []byte    `json:"certificado_pem"`
	LlavePrivadaPEM []byte    `json:"llave_privada_pem"`
	FechaEmision    time.Time `json:"fecha_emision"`
	FechaExpiracion time.Time `json:"fecha_expiracion"`
	Estado          string    `json:"estado"` // Activo, Revocado, Expirado
	Tipo            string    `json:"tipo"`   // Firma, SSL, etc.
}

// ValidarVigencia verifica si el certificado está vigente
func (c *Certificado) ValidarVigencia() bool {
	now := time.Now()
	return now.After(c.FechaEmision) && now.Before(c.FechaExpiracion)
}

// ObtenerCertificadoX509 parsea el certificado PEM a x509.Certificate
func (c *Certificado) ObtenerCertificadoX509() (*x509.Certificate, error) {
	return x509.ParseCertificate(c.CertificadoPEM)
}

// DiasParaExpiracion calcula los días restantes hasta la expiración
func (c *Certificado) DiasParaExpiracion() int {
	return int(time.Until(c.FechaExpiracion).Hours() / 24)
}

// RequiereRenovacion verifica si el certificado necesita ser renovado (menos de 30 días)
func (c *Certificado) RequiereRenovacion() bool {
	return c.DiasParaExpiracion() < 30
}
