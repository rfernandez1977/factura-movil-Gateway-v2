package models

import (
	"time"
)

// ColaReintentos representa una cola para reintentar operaciones fallidas
type ColaReintentos struct {
	ID               string                 `json:"id" bson:"_id,omitempty"`
	TipoOperacion    string                 `json:"tipo_operacion" bson:"tipo_operacion"`   // ENVIO_DTE, TIMBRE, FIRMA, CONSULTA, ANULACION, etc.
	ReferenciaID     string                 `json:"referencia_id" bson:"referencia_id"`     // ID del documento u objeto referenciado
	TipoReferencia   string                 `json:"tipo_referencia" bson:"tipo_referencia"` // DOCUMENTO, SOBRE, BOLETA, etc.
	FechaCreacion    time.Time              `json:"fecha_creacion" bson:"fecha_creacion"`
	UltimoReintento  time.Time              `json:"ultimo_reintento" bson:"ultimo_reintento"`
	ProximoReintento time.Time              `json:"proximo_reintento" bson:"proximo_reintento"`
	NumeroIntentos   int                    `json:"numero_intentos" bson:"numero_intentos"`
	MaximoIntentos   int                    `json:"maximo_intentos" bson:"maximo_intentos"`
	Estado           string                 `json:"estado" bson:"estado"` // PENDIENTE, EN_PROCESO, COMPLETADO, FALLIDO, CANCELADO
	UltimoError      string                 `json:"ultimo_error" bson:"ultimo_error"`
	Prioridad        int                    `json:"prioridad" bson:"prioridad"` // 1-alta, 2-normal, 3-baja
	Contexto         map[string]interface{} `json:"contexto" bson:"contexto,omitempty"`
	EmpresaID        string                 `json:"empresa_id" bson:"empresa_id"`
	Activo           bool                   `json:"activo" bson:"activo"`
	CreatedAt        time.Time              `json:"created_at" bson:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at" bson:"updated_at"`
}

// NewColaReintentos crea una nueva entrada en la cola de reintentos
func NewColaReintentos(tipoOperacion, referenciaID, tipoReferencia, empresaID string, maximoIntentos int) *ColaReintentos {
	now := time.Now()
	return &ColaReintentos{
		TipoOperacion:    tipoOperacion,
		ReferenciaID:     referenciaID,
		TipoReferencia:   tipoReferencia,
		FechaCreacion:    now,
		ProximoReintento: now,
		NumeroIntentos:   0,
		MaximoIntentos:   maximoIntentos,
		Estado:           "PENDIENTE",
		Prioridad:        2, // Normal por defecto
		EmpresaID:        empresaID,
		Activo:           true,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
}

// IncrementarIntentos incrementa el contador de intentos y actualiza las fechas
func (c *ColaReintentos) IncrementarIntentos(error string, minutosEspera int) {
	now := time.Now()
	c.NumeroIntentos++
	c.UltimoReintento = now
	c.ProximoReintento = now.Add(time.Duration(minutosEspera) * time.Minute)
	c.UltimoError = error
	c.UpdatedAt = now

	// Si alcanzamos el máximo de intentos, marcar como fallido
	if c.NumeroIntentos >= c.MaximoIntentos {
		c.Estado = "FALLIDO"
		c.Activo = false
	}
}

// MarcarCompletado marca la entrada como completada
func (c *ColaReintentos) MarcarCompletado() {
	now := time.Now()
	c.Estado = "COMPLETADO"
	c.Activo = false
	c.UpdatedAt = now
}

// MarcarCancelado marca la entrada como cancelada
func (c *ColaReintentos) MarcarCancelado(motivo string) {
	now := time.Now()
	c.Estado = "CANCELADO"
	c.UltimoError = motivo
	c.Activo = false
	c.UpdatedAt = now
}
