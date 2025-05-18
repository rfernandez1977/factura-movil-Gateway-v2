package services

import (
	"crypto/x509"
	"time"
)

// ResultadoFirma contiene el resultado de una operación de firma
type ResultadoFirma struct {
	XMLFirmado     string
	DigestValue    string
	SignatureValue string
	Timestamp      time.Time
}

// EstadoFirma contiene el estado de validación de una firma
type EstadoFirma struct {
	Valida          bool
	FechaValidacion time.Time
	CertificadoID   string
	Error           string
}

// FirmaService define la interfaz base para servicios de firma digital
type FirmaService interface {
	// FirmarDocumento firma un documento XML
	FirmarDocumento(xml string) (*ResultadoFirma, error)

	// ValidarFirma valida la firma de un documento XML
	ValidarFirma(xml string) (*EstadoFirma, error)

	// ObtenerCertificado obtiene el certificado actual del servicio
	ObtenerCertificado() (*x509.Certificate, error)
}

// SIIFirmaService define la interfaz para el servicio de firma del SII
type SIIFirmaService interface {
	FirmaService

	// FirmarSemilla firma un documento de semilla del SII
	FirmarSemilla(xml string) (*ResultadoFirma, error)

	// FirmarToken firma un documento de token del SII
	FirmarToken(xml string) (*ResultadoFirma, error)

	// FirmarDTE firma un Documento Tributario Electrónico
	FirmarDTE(xml string) (*ResultadoFirma, error)
}

// CertCache define la interfaz para el caché de certificados
type CertCache interface {
	// Get obtiene un certificado del caché
	Get(key string) *x509.Certificate

	// Set almacena un certificado en el caché
	Set(key string, cert *x509.Certificate)
} 