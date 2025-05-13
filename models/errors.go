package models

import (
	"fmt"
	"time"
)

// ErrorSII representa un error específico del SII
type ErrorSII struct {
	ID          string    `json:"id" bson:"id"`
	Codigo      string    `json:"codigo" bson:"codigo" xml:"CODIGO"`
	Descripcion string    `json:"descripcion" bson:"descripcion" xml:"DESCRIPCION"`
	Mensaje     string    `json:"mensaje" bson:"mensaje"`
	Detalle     string    `json:"detalle,omitempty" bson:"detalle,omitempty" xml:"DETALLE,omitempty"`
	Detalles    string    `json:"detalles" bson:"detalles"`
	Timestamp   time.Time `json:"timestamp" bson:"timestamp"`
}

func (e *ErrorSII) Error() string {
	return e.Descripcion
}

// ErrorLog representa un registro de error en el sistema
type ErrorLog struct {
	ID        string    `json:"id"`
	Tipo      string    `json:"tipo"`
	Mensaje   string    `json:"mensaje"`
	Detalles  string    `json:"detalles"`
	Usuario   string    `json:"usuario"`
	Timestamp time.Time `json:"timestamp"`
}

// ErrorResponse representa una respuesta de error
type ErrorResponse struct {
	Codigo   string `json:"codigo"`
	Mensaje  string `json:"mensaje"`
	Detalles string `json:"detalles,omitempty"`
}

// ErrorValidation representa un error de validación
type ErrorValidation struct {
	Campo   string `json:"campo"`
	Mensaje string `json:"mensaje"`
	Valor   string `json:"valor,omitempty"`
}

// ErrorBusiness representa un error de negocio
type ErrorBusiness struct {
	Codigo   string `json:"codigo"`
	Mensaje  string `json:"mensaje"`
	Detalles string `json:"detalles,omitempty"`
}

// ErrorSecurity representa un error de seguridad
type ErrorSecurity struct {
	Codigo   string `json:"codigo"`
	Mensaje  string `json:"mensaje"`
	Detalles string `json:"detalles,omitempty"`
}

// ErrorCommunication representa un error de comunicación
type ErrorCommunication struct {
	Codigo   string `json:"codigo"`
	Mensaje  string `json:"mensaje"`
	Detalles string `json:"detalles,omitempty"`
}

// Tipos de error
const (
	ErrorTipoSII          = "SII"
	ErrorTipoValidacion   = "VALIDACION"
	ErrorTipoSistema      = "SISTEMA"
	ErrorTipoIntegracion  = "INTEGRACION"
	ErrorTipoNegocio      = "NEGOCIO"
	ErrorTipoSeguridad    = "SEGURIDAD"
	ErrorTipoComunicacion = "COMUNICACION"
)

// GenerateErrorID genera un ID único para los errores
func GenerateErrorID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// Funciones auxiliares para crear errores
func NewErrorSII(codigo, mensaje, detalles string) *ErrorSII {
	return &ErrorSII{
		ID:        GenerateErrorID(),
		Codigo:    codigo,
		Mensaje:   mensaje,
		Detalles:  detalles,
		Timestamp: time.Now(),
	}
}

func NewErrorLog(tipo, mensaje, detalles, usuario string) *ErrorLog {
	return &ErrorLog{
		ID:        GenerateErrorID(),
		Tipo:      tipo,
		Mensaje:   mensaje,
		Detalles:  detalles,
		Usuario:   usuario,
		Timestamp: time.Now(),
	}
}

func NewErrorResponse(codigo, mensaje, detalles string) *ErrorResponse {
	return &ErrorResponse{
		Codigo:   codigo,
		Mensaje:  mensaje,
		Detalles: detalles,
	}
}

func NewErrorValidation(campo, mensaje, valor string) *ErrorValidation {
	return &ErrorValidation{
		Campo:   campo,
		Mensaje: mensaje,
		Valor:   valor,
	}
}

func NewErrorBusiness(codigo, mensaje, detalles string) *ErrorBusiness {
	return &ErrorBusiness{
		Codigo:   codigo,
		Mensaje:  mensaje,
		Detalles: detalles,
	}
}

func NewErrorSecurity(codigo, mensaje, detalles string) *ErrorSecurity {
	return &ErrorSecurity{
		Codigo:   codigo,
		Mensaje:  mensaje,
		Detalles: detalles,
	}
}

func NewErrorCommunication(codigo, mensaje, detalles string) *ErrorCommunication {
	return &ErrorCommunication{
		Codigo:   codigo,
		Mensaje:  mensaje,
		Detalles: detalles,
	}
}

// ErrorValidacion representa un error de validación
type ErrorValidacion struct {
	Campo   string `json:"campo"`
	Mensaje string `json:"mensaje"`
}

func (e *ErrorValidacion) Error() string {
	return fmt.Sprintf("Error de validación en %s: %s", e.Campo, e.Mensaje)
}

// ErrorSistema representa un error interno del sistema
type ErrorSistema struct {
	Codigo    int       `json:"codigo"`
	Mensaje   string    `json:"mensaje"`
	Timestamp time.Time `json:"timestamp"`
}

func (e *ErrorSistema) Error() string {
	return fmt.Sprintf("Error del sistema [%d]: %s", e.Codigo, e.Mensaje)
}

// ErrorIntegracion representa un error de integración con sistemas externos
type ErrorIntegracion struct {
	Sistema   string    `json:"sistema"`
	Codigo    int       `json:"codigo"`
	Mensaje   string    `json:"mensaje"`
	Timestamp time.Time `json:"timestamp"`
}

func (e *ErrorIntegracion) Error() string {
	return fmt.Sprintf("Error de integración con %s [%d]: %s", e.Sistema, e.Codigo, e.Mensaje)
}

// CodigoError representa los códigos de error posibles
type CodigoError string

const (
	// Errores de validación
	ErrorValidacionSchema  CodigoError = "VAL001"
	ErrorValidacionFirma   CodigoError = "VAL002"
	ErrorValidacionRUT     CodigoError = "VAL003"
	ErrorValidacionFormato CodigoError = "VAL004"
	ErrorValidacionNegocio CodigoError = "VAL005"

	// Errores del SII
	ErrorSIIServicio      CodigoError = "SII001"
	ErrorSIIAutenticacion CodigoError = "SII002"
	ErrorSIIAutorizacion  CodigoError = "SII003"
	ErrorSIIComunicacion  CodigoError = "SII004"
	ErrorSIIProcesamiento CodigoError = "SII005"

	// Errores del sistema
	ErrorSistemaDB           CodigoError = "SYS001"
	ErrorSistemaArchivo      CodigoError = "SYS002"
	ErrorSistemaMemoria      CodigoError = "SYS003"
	ErrorSistemaConcurrencia CodigoError = "SYS004"

	// Errores de integración
	ErrorIntegracionTimeout CodigoError = "INT001"
	ErrorIntegracionFormato CodigoError = "INT002"
	ErrorIntegracionEstado  CodigoError = "INT003"
)

// IntentoRecuperacion representa un intento de recuperación de un error
type IntentoRecuperacion struct {
	NumeroIntento int       `json:"numero_intento"`
	Timestamp     time.Time `json:"timestamp"`
	Exitoso       bool      `json:"exitoso"`
	Mensaje       string    `json:"mensaje,omitempty"`
}

// LogError representa un registro de error en el sistema
type LogError struct {
	ID                 int64                 `json:"id"`
	Tipo               string                `json:"tipo"`
	Severidad          string                `json:"severidad"`
	Mensaje            string                `json:"mensaje"`
	Detalles           ErrorSistema          `json:"detalles"`
	IntentoActual      int                   `json:"intento_actual"`
	MaxIntentos        int                   `json:"max_intentos"`
	Intentos           []IntentoRecuperacion `json:"intentos,omitempty"`
	FechaCreacion      time.Time             `json:"fecha_creacion"`
	FechaActualizacion time.Time             `json:"fecha_actualizacion"`
}

// Error implementa la interfaz error
func (e *LogError) Error() string {
	return fmt.Sprintf("[%s] %s: %s", e.Severidad, e.Tipo, e.Mensaje)
}

// Niveles de severidad
const (
	ErrorSeveridadCritico     = "CRITICO"
	ErrorSeveridadError       = "ERROR"
	ErrorSeveridadAdvertencia = "ADVERTENCIA"
	ErrorSeveridadInfo        = "INFO"
	ErrorSeveridadDebug       = "DEBUG"
)

// ReporteErrores contiene estadísticas de errores
type ReporteErrores struct {
	ID                       string         `bson:"_id" json:"id"`
	FechaInicio              time.Time      `bson:"fecha_inicio" json:"fecha_inicio"`
	FechaFin                 time.Time      `bson:"fecha_fin" json:"fecha_fin"`
	TotalErrores             int            `bson:"total_errores" json:"total_errores"`
	ErroresPorTipo           map[string]int `bson:"errores_por_tipo" json:"errores_por_tipo"`
	ErroresPorSeveridad      map[string]int `bson:"errores_por_severidad" json:"errores_por_severidad"`
	ErroresResueltos         int            `bson:"errores_resueltos" json:"errores_resueltos"`
	ErroresPendientes        int            `bson:"errores_pendientes" json:"errores_pendientes"`
	TiempoPromedioResolucion int64          `bson:"tiempo_promedio_resolucion" json:"tiempo_promedio_resolucion"`
	FechaGeneracion          time.Time      `bson:"fecha_generacion" json:"fecha_generacion"`
}
