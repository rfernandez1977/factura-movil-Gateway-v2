package models

import (
	"time"
)

// IntentoRecuperacion representa un intento de recuperación tras un error
type IntentoRecuperacion struct {
	ID                      string    `json:"id" bson:"_id,omitempty"`
	LogErrorID              string    `json:"log_error_id" bson:"log_error_id"`
	UsuarioID               string    `json:"usuario_id" bson:"usuario_id"`
	FechaInicio             time.Time `json:"fecha_inicio" bson:"fecha_inicio"`
	FechaFinalizacion       time.Time `json:"fecha_finalizacion" bson:"fecha_finalizacion,omitempty"`
	Estado                  string    `json:"estado" bson:"estado"` // PENDIENTE, EN_PROCESO, EXITOSO, FALLIDO
	Descripcion             string    `json:"descripcion" bson:"descripcion"`
	AccionesRealizadas      []string  `json:"acciones_realizadas" bson:"acciones_realizadas"`
	ResultadoAcciones       string    `json:"resultado_acciones" bson:"resultado_acciones"`
	ErrorNuevo              string    `json:"error_nuevo" bson:"error_nuevo,omitempty"`
	Notas                   string    `json:"notas" bson:"notas"`
	TiempoInvertidoSegundos int       `json:"tiempo_invertido_segundos" bson:"tiempo_invertido_segundos"`
	EmpresaID               string    `json:"empresa_id" bson:"empresa_id"`
	DocumentoID             string    `json:"documento_id" bson:"documento_id,omitempty"`
	CreatedAt               time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt               time.Time `json:"updated_at" bson:"updated_at"`
}

// NewIntentoRecuperacion crea una nueva instancia de IntentoRecuperacion
func NewIntentoRecuperacion(logErrorID, usuarioID, descripcion string, empresaID string) *IntentoRecuperacion {
	now := time.Now()
	return &IntentoRecuperacion{
		LogErrorID:         logErrorID,
		UsuarioID:          usuarioID,
		FechaInicio:        now,
		Estado:             "PENDIENTE",
		Descripcion:        descripcion,
		AccionesRealizadas: []string{},
		EmpresaID:          empresaID,
		CreatedAt:          now,
		UpdatedAt:          now,
	}
}

// AgregarAccion agrega una acción realizada al intento de recuperación
func (i *IntentoRecuperacion) AgregarAccion(accion string) {
	i.AccionesRealizadas = append(i.AccionesRealizadas, accion)
	i.UpdatedAt = time.Now()
}

// FinalizarExitoso marca el intento de recuperación como exitoso
func (i *IntentoRecuperacion) FinalizarExitoso(resultado string) {
	i.Estado = "EXITOSO"
	i.ResultadoAcciones = resultado
	i.FechaFinalizacion = time.Now()
	i.TiempoInvertidoSegundos = int(i.FechaFinalizacion.Sub(i.FechaInicio).Seconds())
	i.UpdatedAt = time.Now()
}

// FinalizarFallido marca el intento de recuperación como fallido
func (i *IntentoRecuperacion) FinalizarFallido(resultado string, errorNuevo string) {
	i.Estado = "FALLIDO"
	i.ResultadoAcciones = resultado
	i.ErrorNuevo = errorNuevo
	i.FechaFinalizacion = time.Now()
	i.TiempoInvertidoSegundos = int(i.FechaFinalizacion.Sub(i.FechaInicio).Seconds())
	i.UpdatedAt = time.Now()
}
