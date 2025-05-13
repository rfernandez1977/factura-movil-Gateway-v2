package services

import (
	"crypto/rand"
	"fmt"
	"time"

	"github.com/cursor/FMgo/utils"
)

// UtilsService proporciona métodos de utilidad
type UtilsService struct {
	emailService *EmailService
}

// NewUtilsService crea una nueva instancia del servicio de utilidades
func NewUtilsService(emailService *EmailService) *UtilsService {
	return &UtilsService{
		emailService: emailService,
	}
}

// IsValidRUT valida si un RUT es válido
func (u *UtilsService) IsValidRUT(rut string) bool {
	return utils.ValidateRUT(rut) == nil
}

// EnviarEmail envía un email
func (u *UtilsService) EnviarEmail(destinatario, asunto, mensaje string) error {
	return u.emailService.Enviar(destinatario, asunto, mensaje)
}

// FormatDate formatea una fecha en el formato requerido por el SII
func FormatDate(date time.Time) string {
	return date.Format("2006-01-02")
}

// FormatDateTime formatea una fecha y hora en el formato requerido por el SII
func FormatDateTime(date time.Time) string {
	return date.Format("2006-01-02T15:04:05")
}

// GenerateID genera un ID único para identificadores en el sistema
// Esta función reemplaza las implementaciones duplicadas en:
// - seguridad_service.go
// - ecommerce_service.go
// - reportes_service.go
// y otros servicios que generan IDs
func GenerateID() string {
	// Generar bytes aleatorios
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		// En caso de error, usar timestamp como fallback
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}

	// Formatear como string hexadecimal
	return fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
