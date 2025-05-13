package repository

import (
	"fmt"
	"time"
)

// ValidationError representa un error de validación en un modelo
type ValidationError struct {
	Field   string
	Message string
}

// Error implementa la interfaz error
func (e *ValidationError) Error() string {
	return fmt.Sprintf("Error de validación en campo '%s': %s", e.Field, e.Message)
}

// NewValidationError crea un nuevo error de validación
func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Message: message,
	}
}

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
		return NewValidationError("nombre", "El nombre de la empresa es obligatorio")
	}
	if e.RUT == "" {
		return NewValidationError("rut", "El RUT de la empresa es obligatorio")
	}
	if e.RUTFirma == "" {
		return NewValidationError("rut_firma", "El RUT de la firma es obligatorio")
	}
	if e.NombreFirma == "" {
		return NewValidationError("nombre_firma", "El nombre de la firma es obligatorio")
	}
	if e.ClaveFirma == "" {
		return NewValidationError("clave_firma", "La clave de la firma es obligatoria")
	}
	return nil
}

// Documento representa un documento tributario
type Documento struct {
	ID              string    `json:"id" db:"id"`
	EmpresaID       string    `json:"empresa_id" db:"empresa_id"`
	TipoDocumento   string    `json:"tipo_documento" db:"tipo_documento"`
	NumeroDocumento string    `json:"numero_documento" db:"numero_documento"`
	FechaEmision    string    `json:"fecha_emision" db:"fecha_emision"`
	Monto           float64   `json:"monto" db:"monto"`
	Estado          string    `json:"estado" db:"estado"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

// NewDocumento crea una nueva instancia de Documento
func NewDocumento(empresaID, tipoDocumento, numeroDocumento, fechaEmision string, monto float64) *Documento {
	return &Documento{
		EmpresaID:       empresaID,
		TipoDocumento:   tipoDocumento,
		NumeroDocumento: numeroDocumento,
		FechaEmision:    fechaEmision,
		Monto:           monto,
		Estado:          "PENDIENTE",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
}

// Validate valida que todos los campos obligatorios estén presentes
func (d *Documento) Validate() error {
	if d.EmpresaID == "" {
		return NewValidationError("empresa_id", "El ID de la empresa es obligatorio")
	}
	if d.TipoDocumento == "" {
		return NewValidationError("tipo_documento", "El tipo de documento es obligatorio")
	}
	if d.NumeroDocumento == "" {
		return NewValidationError("numero_documento", "El número de documento es obligatorio")
	}
	if d.FechaEmision == "" {
		return NewValidationError("fecha_emision", "La fecha de emisión es obligatoria")
	}
	return nil
}

// EstadosDocumento define los posibles estados de un documento
var EstadosDocumento = struct {
	Pendiente  string
	Enviado    string
	Aceptado   string
	Rechazado  string
	Anulado    string
	Procesando string
	Error      string
	Completado string
}{
	Pendiente:  "PENDIENTE",
	Enviado:    "ENVIADO",
	Aceptado:   "ACEPTADO",
	Rechazado:  "RECHAZADO",
	Anulado:    "ANULADO",
	Procesando: "PROCESANDO",
	Error:      "ERROR",
	Completado: "COMPLETADO",
}

// CAF representa un Código de Autorización de Folios
type CAF struct {
	ID               string    `json:"id" db:"id"`
	EmpresaID        string    `json:"empresa_id" db:"empresa_id"`
	TipoDocumento    string    `json:"tipo_documento" db:"tipo_documento"`
	Desde            int       `json:"desde" db:"desde"`
	Hasta            int       `json:"hasta" db:"hasta"`
	Archivo          []byte    `json:"archivo" db:"archivo"`
	FechaVencimiento string    `json:"fecha_vencimiento" db:"fecha_vencimiento"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

// NewCAF crea una nueva instancia de CAF
func NewCAF(empresaID, tipoDocumento string, desde, hasta int, archivo []byte, fechaVencimiento string) *CAF {
	return &CAF{
		EmpresaID:        empresaID,
		TipoDocumento:    tipoDocumento,
		Desde:            desde,
		Hasta:            hasta,
		Archivo:          archivo,
		FechaVencimiento: fechaVencimiento,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
}

// Validate valida que todos los campos obligatorios estén presentes
func (c *CAF) Validate() error {
	if c.EmpresaID == "" {
		return NewValidationError("empresa_id", "El ID de la empresa es obligatorio")
	}
	if c.TipoDocumento == "" {
		return NewValidationError("tipo_documento", "El tipo de documento es obligatorio")
	}
	if c.Desde <= 0 {
		return NewValidationError("desde", "El rango inicial debe ser mayor a cero")
	}
	if c.Hasta <= 0 {
		return NewValidationError("hasta", "El rango final debe ser mayor a cero")
	}
	if c.Hasta < c.Desde {
		return NewValidationError("hasta", "El rango final debe ser mayor o igual al rango inicial")
	}
	if c.Archivo == nil || len(c.Archivo) == 0 {
		return NewValidationError("archivo", "El archivo del CAF es obligatorio")
	}
	return nil
}
