package models

// CertificadoDigital representa un certificado digital
type CertificadoDigital struct {
	RUT      string
	Nombre   string
	Archivo  string
	Password string
}

// ResultadoFirma contiene el resultado de la operación de firma
type ResultadoFirma struct {
	XMLFirmado     string `json:"xml_firmado"`
	DigestValue    string `json:"digest_value"`
	SignatureValue string `json:"signature_value"`
}
