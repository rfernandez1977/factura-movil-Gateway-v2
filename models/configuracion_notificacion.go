package models

import "time"

// SistemaNotificaciones representa la configuración del sistema de notificaciones
type SistemaNotificaciones struct {
	ID                 string            `json:"id" bson:"_id,omitempty"`
	MaxIntentosEnvio   int               `json:"max_intentos_envio" bson:"max_intentos_envio"`
	IntervaloReintento time.Duration     `json:"intervalo_reintento" bson:"intervalo_reintento"`
	NotificarEmail     bool              `json:"notificar_email" bson:"notificar_email"`
	NotificarSMS       bool              `json:"notificar_sms" bson:"notificar_sms"`
	NotificarPush      bool              `json:"notificar_push" bson:"notificar_push"`
	PlantillasEmail    map[string]string `json:"plantillas_email" bson:"plantillas_email"`
	PlantillasSMS      map[string]string `json:"plantillas_sms" bson:"plantillas_sms"`
	PlantillasPush     map[string]string `json:"plantillas_push" bson:"plantillas_push"`
	HorarioEnvio       []HorarioEnvio    `json:"horario_envio,omitempty" bson:"horario_envio,omitempty"`
	FechaCreacion      time.Time         `json:"fecha_creacion" bson:"fecha_creacion"`
	FechaActualizacion time.Time         `json:"fecha_actualizacion" bson:"fecha_actualizacion"`
}

// HorarioEnvio representa un horario para el envío de notificaciones
type HorarioEnvio struct {
	DiaSemana  int       `json:"dia_semana" bson:"dia_semana"`
	HoraInicio time.Time `json:"hora_inicio" bson:"hora_inicio"`
	HoraFin    time.Time `json:"hora_fin" bson:"hora_fin"`
}

// Validate valida que todos los campos obligatorios estén presentes
func (c *SistemaNotificaciones) Validate() error {
	if c.MaxIntentosEnvio <= 0 {
		return &ValidationFieldError{Field: "max_intentos_envio", Message: "El número máximo de intentos de envío debe ser mayor que cero"}
	}
	if c.IntervaloReintento <= 0 {
		return &ValidationFieldError{Field: "intervalo_reintento", Message: "El intervalo de reintento debe ser mayor que cero"}
	}
	return nil
}
