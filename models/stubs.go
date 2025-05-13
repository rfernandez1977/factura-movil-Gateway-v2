package models

import "time"

// ErrorDetalle representa los detalles de un error
type ErrorDetalle struct {
	ID          string                 `json:"id" bson:"_id,omitempty"`
	Tipo        string                 `json:"tipo" bson:"tipo"`
	Severidad   string                 `json:"severidad" bson:"severidad"`
	Codigo      string                 `json:"codigo" bson:"codigo"`
	Mensaje     string                 `json:"mensaje" bson:"mensaje"`
	Descripcion string                 `json:"descripcion" bson:"descripcion"`
	Stacktrace  string                 `json:"stacktrace" bson:"stacktrace"`
	Contexto    map[string]interface{} `json:"contexto" bson:"contexto"`
	Entidad     string                 `json:"entidad" bson:"entidad"`
	EntidadID   string                 `json:"entidad_id" bson:"entidad_id"`
	UsuarioID   string                 `json:"usuario_id" bson:"usuario_id"`
	FechaError  time.Time              `json:"fecha_error" bson:"fecha_error"`
	Estado      string                 `json:"estado" bson:"estado"`
}

// ConfiguracionRecuperacion representa la configuración para recuperación de errores
type ConfiguracionRecuperacion struct {
	ID                 string  `json:"id" bson:"_id,omitempty"`
	MaxIntentos        int     `json:"max_intentos" bson:"max_intentos"`
	IntervaloBase      int     `json:"intervalo_base" bson:"intervalo_base"`
	FactorExponencial  float64 `json:"factor_exponencial" bson:"factor_exponencial"`
	MaxIntervalo       int     `json:"max_intervalo" bson:"max_intervalo"`
	NotificarAdmin     bool    `json:"notificar_admin" bson:"notificar_admin"`
	UmbralNotificacion int     `json:"umbral_notificacion" bson:"umbral_notificacion"`
	ActivarLogging     bool    `json:"activar_logging" bson:"activar_logging"`
	NivelLogging       string  `json:"nivel_logging" bson:"nivel_logging"`
}

// DTEType representa un documento tributario electrónico para servicios
type DTEType struct {
	ID     string `json:"id" bson:"_id,omitempty"`
	Tipo   string `json:"tipo" bson:"tipo"`
	Nombre string `json:"nombre" bson:"nombre"`
}

// RegistroSincronizacion representa un registro de sincronización con sistemas externos
type RegistroSincronizacion struct {
	ID             string    `json:"id" bson:"_id,omitempty"`
	Sistema        string    `json:"sistema" bson:"sistema"`
	TipoOperacion  string    `json:"tipo_operacion" bson:"tipo_operacion"`
	Timestamp      time.Time `json:"timestamp" bson:"timestamp"`
	DatosEnviados  string    `json:"datos_enviados" bson:"datos_enviados"`
	DatosRecibidos string    `json:"datos_recibidos" bson:"datos_recibidos"`
	Resultado      string    `json:"resultado" bson:"resultado"`
	Error          string    `json:"error,omitempty" bson:"error,omitempty"`
}
