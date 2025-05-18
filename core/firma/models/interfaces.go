package models

import (
	"crypto/x509"
)

// FirmaDigital define la interfaz base para servicios de firma digital
type FirmaDigital interface {
	// FirmarDocumento firma un documento XML y retorna el documento firmado
	FirmarDocumento(documento string) (*ResultadoFirma, error)

	// ValidarFirma valida la firma de un documento XML
	ValidarFirma(xmlFirmado string) (*EstadoFirma, error)

	// ObtenerCertificado retorna el certificado actual
	ObtenerCertificado() (*x509.Certificate, error)
}

// FirmaSII extiende FirmaDigital con funcionalidades específicas del SII
type FirmaSII interface {
	FirmaDigital

	// FirmarSemilla firma una semilla del SII
	FirmarSemilla(semilla string) (*ResultadoFirma, error)

	// FirmarToken firma un token de autenticación
	FirmarToken(token string) (*ResultadoFirma, error)

	// ValidarCAF valida un archivo CAF
	ValidarCAF(caf []byte) error
}

// ResultadoFirma contiene el resultado de una operación de firma
type ResultadoFirma struct {
	XMLFirmado     string
	DigestValue    string
	SignatureValue string
	Timestamp      string
}

// EstadoFirma contiene el resultado de una validación de firma
type EstadoFirma struct {
	Valida          bool
	FechaValidacion string
	CertificadoID   string
	Error           string
}

// ConfiguracionFirma contiene la configuración necesaria para el servicio de firma
type ConfiguracionFirma struct {
	RutaCertificado       string
	RutaLlave             string
	Password              string
	RutEmpresa            string
	RutaConfiguracion     string
	AmbienteCertificacion bool
}
