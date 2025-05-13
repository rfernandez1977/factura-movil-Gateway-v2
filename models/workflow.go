package models

import (
	"time"
)

// Workflow representa un flujo de trabajo
type Workflow struct {
	ID          string            `json:"id" bson:"_id,omitempty"`
	Nombre      string            `json:"nombre" bson:"nombre"`
	Descripcion string            `json:"descripcion" bson:"descripcion"`
	Estado      EstadoFlujo       `json:"estado" bson:"estado"`
	Pasos       []PasoWorkflow    `json:"pasos" bson:"pasos"`
	Condiciones []Condicion       `json:"condiciones,omitempty" bson:"condiciones,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty" bson:"metadata,omitempty"`
	CreatedAt   time.Time         `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at" bson:"updated_at"`
}

// PasoWorkflow representa un paso dentro de un flujo de trabajo
type PasoWorkflow struct {
	ID              string            `json:"id" bson:"_id,omitempty"`
	WorkflowID      string            `json:"workflow_id" bson:"workflow_id"`
	Nombre          string            `json:"nombre" bson:"nombre"`
	Descripcion     string            `json:"descripcion" bson:"descripcion"`
	Orden           int               `json:"orden" bson:"orden"`
	TipoPaso        string            `json:"tipo_paso" bson:"tipo_paso"`
	Estado          EstadoPaso        `json:"estado" bson:"estado"`
	PasoAnteriorID  string            `json:"paso_anterior_id,omitempty" bson:"paso_anterior_id,omitempty"`
	PasoSiguienteID string            `json:"paso_siguiente_id,omitempty" bson:"paso_siguiente_id,omitempty"`
	Parametros      map[string]string `json:"parametros,omitempty" bson:"parametros,omitempty"`
	Timeout         int               `json:"timeout" bson:"timeout"`
	Reintentos      int               `json:"reintentos" bson:"reintentos"`
	CreatedAt       time.Time         `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at" bson:"updated_at"`
}

// Condicion representa una condición para la ejecución de pasos en un workflow
type Condicion struct {
	ID          string            `json:"id" bson:"_id,omitempty"`
	WorkflowID  string            `json:"workflow_id" bson:"workflow_id"`
	PasoID      string            `json:"paso_id" bson:"paso_id"`
	Nombre      string            `json:"nombre" bson:"nombre"`
	Descripcion string            `json:"descripcion" bson:"descripcion"`
	Tipo        string            `json:"tipo" bson:"tipo"`
	Expresion   string            `json:"expresion" bson:"expresion"`
	Parametros  map[string]string `json:"parametros,omitempty" bson:"parametros,omitempty"`
	CreatedAt   time.Time         `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at" bson:"updated_at"`
}

// ColaReintentos representa una cola de reintentos para pasos fallidos
type ColaReintentos struct {
	ID             string    `json:"id" bson:"_id,omitempty"`
	WorkflowID     string    `json:"workflow_id" bson:"workflow_id"`
	PasoID         string    `json:"paso_id" bson:"paso_id"`
	Estado         string    `json:"estado" bson:"estado"`
	NumeroIntento  int       `json:"numero_intento" bson:"numero_intento"`
	TiempoRetardo  int       `json:"tiempo_retardo" bson:"tiempo_retardo"`
	UltimoError    string    `json:"ultimo_error,omitempty" bson:"ultimo_error,omitempty"`
	ProximoIntento time.Time `json:"proximo_intento" bson:"proximo_intento"`
	CreatedAt      time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" bson:"updated_at"`
}

// Validate valida que todos los campos obligatorios estén presentes en un Workflow
func (w *Workflow) Validate() error {
	if w.Nombre == "" {
		return &ValidationFieldError{Field: "nombre", Message: "El nombre del workflow es obligatorio"}
	}
	if len(w.Pasos) == 0 {
		return &ValidationFieldError{Field: "pasos", Message: "El workflow debe tener al menos un paso"}
	}
	return nil
}

// Validate valida que todos los campos obligatorios estén presentes en un PasoWorkflow
func (p *PasoWorkflow) Validate() error {
	if p.WorkflowID == "" {
		return &ValidationFieldError{Field: "workflow_id", Message: "El ID del workflow es obligatorio"}
	}
	if p.Nombre == "" {
		return &ValidationFieldError{Field: "nombre", Message: "El nombre del paso es obligatorio"}
	}
	if p.TipoPaso == "" {
		return &ValidationFieldError{Field: "tipo_paso", Message: "El tipo de paso es obligatorio"}
	}
	return nil
}

// Validate valida que todos los campos obligatorios estén presentes en una Condicion
func (c *Condicion) Validate() error {
	if c.WorkflowID == "" {
		return &ValidationFieldError{Field: "workflow_id", Message: "El ID del workflow es obligatorio"}
	}
	if c.PasoID == "" {
		return &ValidationFieldError{Field: "paso_id", Message: "El ID del paso es obligatorio"}
	}
	if c.Expresion == "" {
		return &ValidationFieldError{Field: "expresion", Message: "La expresión de la condición es obligatoria"}
	}
	return nil
}

// Validate valida que todos los campos obligatorios estén presentes en una ColaReintentos
func (c *ColaReintentos) Validate() error {
	if c.WorkflowID == "" {
		return &ValidationFieldError{Field: "workflow_id", Message: "El ID del workflow es obligatorio"}
	}
	if c.PasoID == "" {
		return &ValidationFieldError{Field: "paso_id", Message: "El ID del paso es obligatorio"}
	}
	if c.NumeroIntento < 0 {
		return &ValidationFieldError{Field: "numero_intento", Message: "El número de intento debe ser mayor o igual a cero"}
	}
	return nil
}
