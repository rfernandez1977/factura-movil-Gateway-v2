package models

import "time"

// CertificadoDigital representa un certificado digital
type CertificadoDigital struct {
	ID           string    `json:"id"`
	EmpresaID    string    `json:"empresa_id"`
	SerialNumber string    `json:"serial_number"`
	Issuer       string    `json:"issuer"`
	Subject      string    `json:"subject"`
	ValidFrom    time.Time `json:"valid_from"`
	ValidTo      time.Time `json:"valid_to"`
	Certificate  string    `json:"certificate"`
	PrivateKey   string    `json:"private_key"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
