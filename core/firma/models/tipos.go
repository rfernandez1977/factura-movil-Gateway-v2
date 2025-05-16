package models

// AlgoritmoFirma define los algoritmos de firma soportados
type AlgoritmoFirma string

const (
	// RSA_SHA1 algoritmo RSA con SHA1
	RSA_SHA1 AlgoritmoFirma = "http://www.w3.org/2000/09/xmldsig#rsa-sha1"
	// RSA_SHA256 algoritmo RSA con SHA256
	RSA_SHA256 AlgoritmoFirma = "http://www.w3.org/2001/04/xmldsig-more#rsa-sha256"
)

// AlgoritmoDigest define los algoritmos de digest soportados
type AlgoritmoDigest string

const (
	// SHA1 algoritmo SHA1
	SHA1 AlgoritmoDigest = "http://www.w3.org/2000/09/xmldsig#sha1"
	// SHA256 algoritmo SHA256
	SHA256 AlgoritmoDigest = "http://www.w3.org/2001/04/xmlenc#sha256"
)

// AlgoritmoCanonicalization define los algoritmos de canonicalización
type AlgoritmoCanonicalization string

const (
	// C14N algoritmo de canonicalización XML
	C14N AlgoritmoCanonicalization = "http://www.w3.org/TR/2001/REC-xml-c14n-20010315"
	// C14N_WITH_COMMENTS algoritmo de canonicalización XML con comentarios
	C14N_WITH_COMMENTS AlgoritmoCanonicalization = "http://www.w3.org/TR/2001/REC-xml-c14n-20010315#WithComments"
)

// EstadoCertificado define los estados posibles de un certificado
type EstadoCertificado string

const (
	// ACTIVO certificado activo y válido
	ACTIVO EstadoCertificado = "ACTIVO"
	// REVOCADO certificado revocado
	REVOCADO EstadoCertificado = "REVOCADO"
	// EXPIRADO certificado expirado
	EXPIRADO EstadoCertificado = "EXPIRADO"
	// PENDIENTE_RENOVACION certificado próximo a expirar
	PENDIENTE_RENOVACION EstadoCertificado = "PENDIENTE_RENOVACION"
)

// TipoCertificado define los tipos de certificados soportados
type TipoCertificado string

const (
	// FIRMA certificado para firma digital
	FIRMA TipoCertificado = "FIRMA"
	// SSL certificado SSL/TLS
	SSL TipoCertificado = "SSL"
	// AUTENTICACION certificado para autenticación
	AUTENTICACION TipoCertificado = "AUTENTICACION"
)

// ErrorFirma define los tipos de errores en el proceso de firma
type ErrorFirma string

const (
	// ERROR_CERTIFICADO_INVALIDO certificado no válido
	ERROR_CERTIFICADO_INVALIDO ErrorFirma = "CERTIFICADO_INVALIDO"
	// ERROR_FIRMA_INVALIDA firma no válida
	ERROR_FIRMA_INVALIDA ErrorFirma = "FIRMA_INVALIDA"
	// ERROR_DOCUMENTO_INVALIDO documento no válido
	ERROR_DOCUMENTO_INVALIDO ErrorFirma = "DOCUMENTO_INVALIDO"
	// ERROR_SISTEMA error interno del sistema
	ERROR_SISTEMA ErrorFirma = "ERROR_SISTEMA"
)
