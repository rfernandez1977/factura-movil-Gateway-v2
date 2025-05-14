package models

import (
	"time"
)

// LogError representa un registro detallado de error en el sistema
type LogError struct {
	ID                    string                 `json:"id" bson:"_id,omitempty"`
	Codigo                string                 `json:"codigo" bson:"codigo"`
	Nivel                 string                 `json:"nivel" bson:"nivel"` // ERROR, WARNING, INFO, DEBUG
	Mensaje               string                 `json:"mensaje" bson:"mensaje"`
	DetallesTecnicos      string                 `json:"detalles_tecnicos" bson:"detalles_tecnicos"`
	StackTrace            string                 `json:"stack_trace" bson:"stack_trace"`
	FechaHora             time.Time              `json:"fecha_hora" bson:"fecha_hora"`
	ModuloAfectado        string                 `json:"modulo_afectado" bson:"modulo_afectado"`
	ClaseMetodo           string                 `json:"clase_metodo" bson:"clase_metodo"`
	LineaError            int                    `json:"linea_error" bson:"linea_error"`
	UsuarioID             string                 `json:"usuario_id" bson:"usuario_id"`
	DireccionIP           string                 `json:"direccion_ip" bson:"direccion_ip"`
	Contexto              map[string]interface{} `json:"contexto" bson:"contexto,omitempty"`
	DocumentoID           string                 `json:"documento_id" bson:"documento_id,omitempty"`
	TipoDocumento         string                 `json:"tipo_documento" bson:"tipo_documento,omitempty"`
	EmpresaID             string                 `json:"empresa_id" bson:"empresa_id,omitempty"`
	EstadoResolucion      string                 `json:"estado_resolucion" bson:"estado_resolucion"` // PENDIENTE, EN_PROCESO, RESUELTO, IGNORADO
	IntentoRecuperacionID string                 `json:"intento_recuperacion_id" bson:"intento_recuperacion_id,omitempty"`
	CreatedAt             time.Time              `json:"created_at" bson:"created_at"`
	UpdatedAt             time.Time              `json:"updated_at" bson:"updated_at"`
}

// NewLogError crea una nueva instancia de LogError
func NewLogError(codigo, nivel, mensaje, detallesTecnicos, moduloAfectado, claseMetodo string, lineaError int, usuarioID string) *LogError {
	return &LogError{
		Codigo:           codigo,
		Nivel:            nivel,
		Mensaje:          mensaje,
		DetallesTecnicos: detallesTecnicos,
		ModuloAfectado:   moduloAfectado,
		ClaseMetodo:      claseMetodo,
		LineaError:       lineaError,
		UsuarioID:        usuarioID,
		FechaHora:        time.Now(),
		EstadoResolucion: "PENDIENTE",
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
}
