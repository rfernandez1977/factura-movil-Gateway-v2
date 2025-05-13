package models

import (
	"time"
)

// Empresa representa una empresa en el sistema
type Empresa struct {
	ID          string    `json:"id" db:"id"`
	Nombre      string    `json:"nombre" db:"nombre"`
	RUT         string    `json:"rut" db:"rut"`
	Direccion   string    `json:"direccion" db:"direccion"`
	Telefono    string    `json:"telefono" db:"telefono"`
	Email       string    `json:"email" db:"email"`
	RUTFirma    string    `json:"rut_firma" db:"rut_firma"`
	NombreFirma string    `json:"nombre_firma" db:"nombre_firma"`
	ClaveFirma  string    `json:"clave_firma" db:"clave_firma"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// NewEmpresa crea una nueva instancia de Empresa
func NewEmpresa(nombre, rut, direccion, telefono, email, rutFirma, nombreFirma, claveFirma string) *Empresa {
	return &Empresa{
		Nombre:      nombre,
		RUT:         rut,
		Direccion:   direccion,
		Telefono:    telefono,
		Email:       email,
		RUTFirma:    rutFirma,
		NombreFirma: nombreFirma,
		ClaveFirma:  claveFirma,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// Validate valida que todos los campos obligatorios estén presentes
func (e *Empresa) Validate() error {
	if e.Nombre == "" {
		return ValidationError{Field: "nombre", Message: "El nombre de la empresa es obligatorio"}
	}
	if e.RUT == "" {
		return ValidationError{Field: "rut", Message: "El RUT de la empresa es obligatorio"}
	}
	if e.RUTFirma == "" {
		return ValidationError{Field: "rut_firma", Message: "El RUT de la firma es obligatorio"}
	}
	if e.NombreFirma == "" {
		return ValidationError{Field: "nombre_firma", Message: "El nombre de la firma es obligatorio"}
	}
	if e.ClaveFirma == "" {
		return ValidationError{Field: "clave_firma", Message: "La clave de la firma es obligatoria"}
	}
	return nil
}
