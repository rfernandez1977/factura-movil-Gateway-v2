package models

import (
	"time"
)

// SesionElectronica representa una sesión electrónica con el SII
type SesionElectronica struct {
	ID              string    `json:"id" bson:"_id,omitempty"`
	EmpresaID       string    `json:"empresa_id" bson:"empresa_id"`
	Token           string    `json:"token" bson:"token"`
	Estado          string    `json:"estado" bson:"estado"`
	FechaInicio     time.Time `json:"fecha_inicio" bson:"fecha_inicio"`
	FechaExpiracion time.Time `json:"fecha_expiracion" bson:"fecha_expiracion"`
	Intentos        int       `json:"intentos" bson:"intentos"`
	UltimoAcceso    time.Time `json:"ultimo_acceso" bson:"ultimo_acceso"`
	CreatedAt       time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" bson:"updated_at"`
}

// IsValid verifica si la sesión electrónica es válida
func (s *SesionElectronica) IsValid() bool {
	return s.Estado == "ACTIVA" && time.Now().Before(s.FechaExpiracion)
}

// EstadoSesionInfo representa la información del estado de una sesión
type EstadoSesionInfo struct {
	Estado    string    `json:"estado" bson:"estado"`
	Mensaje   string    `json:"mensaje" bson:"mensaje"`
	Timestamp time.Time `json:"timestamp" bson:"timestamp"`
}

// Sesion representa una sesión de usuario en el sistema
type Sesion struct {
	ID              string    `json:"id" bson:"_id,omitempty"`
	UsuarioID       string    `json:"usuario_id" bson:"usuario_id"`
	Token           string    `json:"token" bson:"token"`
	IP              string    `json:"ip" bson:"ip"`
	UserAgent       string    `json:"user_agent" bson:"user_agent"`
	FechaInicio     time.Time `json:"fecha_inicio" bson:"fecha_inicio"`
	FechaExpiracion time.Time `json:"fecha_expiracion" bson:"fecha_expiracion"`
	UltimoAcceso    time.Time `json:"ultimo_acceso" bson:"ultimo_acceso"`
	Activa          bool      `json:"activa" bson:"activa"`
	CreatedAt       time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" bson:"updated_at"`
}

// IsValid verifica si la sesión es válida
func (s *Sesion) IsValid() bool {
	return s.Activa && time.Now().Before(s.FechaExpiracion)
}

// SesionResponse representa la respuesta de inicio de sesión del SII
type SesionResponse struct {
	Token           string    `json:"token" bson:"token"`
	Estado          string    `json:"estado" bson:"estado"`
	FechaExpiracion time.Time `json:"fecha_expiracion" bson:"fecha_expiracion"`
}
