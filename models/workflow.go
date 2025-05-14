package models

import (
	"time"
)

// Workflow representa un flujo de trabajo para procesar documentos
type Workflow struct {
	ID            string         `json:"id" bson:"_id,omitempty"`
	Nombre        string         `json:"nombre" bson:"nombre"`
	Descripcion   string         `json:"descripcion" bson:"descripcion"`
	Version       string         `json:"version" bson:"version"`
	EmpresaID     string         `json:"empresa_id" bson:"empresa_id"`
	TipoDocumento string         `json:"tipo_documento" bson:"tipo_documento"`
	Pasos         []PasoWorkflow `json:"pasos" bson:"pasos"`
	Condiciones   []Condicion    `json:"condiciones" bson:"condiciones"`
	Estado        string         `json:"estado" bson:"estado"` // ACTIVO, INACTIVO, BORRADOR
	Creador       string         `json:"creador" bson:"creador"`
	CreatedAt     time.Time      `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at" bson:"updated_at"`
}

// PasoWorkflow representa un paso específico dentro de un workflow
type PasoWorkflow struct {
	ID                        string                 `json:"id" bson:"_id,omitempty"`
	WorkflowID                string                 `json:"workflow_id" bson:"workflow_id"`
	Nombre                    string                 `json:"nombre" bson:"nombre"`
	Descripcion               string                 `json:"descripcion" bson:"descripcion"`
	Orden                     int                    `json:"orden" bson:"orden"`
	TipoAccion                string                 `json:"tipo_accion" bson:"tipo_accion"`
	Configuracion             map[string]interface{} `json:"configuracion" bson:"configuracion"`
	RequiereAprobacion        bool                   `json:"requiere_aprobacion" bson:"requiere_aprobacion"`
	RequiereNotificacion      bool                   `json:"requiere_notificacion" bson:"requiere_notificacion"`
	DestinatariosNotificacion []string               `json:"destinatarios_notificacion" bson:"destinatarios_notificacion"`
	TiempoEsperaMaximo        int                    `json:"tiempo_espera_maximo" bson:"tiempo_espera_maximo"` // En minutos
	AccionTimeoutID           string                 `json:"accion_timeout_id" bson:"accion_timeout_id"`
	PasosSiguientes           []string               `json:"pasos_siguientes" bson:"pasos_siguientes"`
	CreatedAt                 time.Time              `json:"created_at" bson:"created_at"`
	UpdatedAt                 time.Time              `json:"updated_at" bson:"updated_at"`
}

// Condicion representa una condición para determinar el flujo del workflow
type Condicion struct {
	ID                 string    `json:"id" bson:"_id,omitempty"`
	WorkflowID         string    `json:"workflow_id" bson:"workflow_id"`
	Nombre             string    `json:"nombre" bson:"nombre"`
	Descripcion        string    `json:"descripcion" bson:"descripcion"`
	TipoCondicion      string    `json:"tipo_condicion" bson:"tipo_condicion"` // CAMPO, ESTADO, PERSONALIZADA
	CampoEvaluacion    string    `json:"campo_evaluacion" bson:"campo_evaluacion"`
	Operador           string    `json:"operador" bson:"operador"` // IGUAL, DIFERENTE, MAYOR, MENOR, CONTIENE, etc.
	ValorComparacion   string    `json:"valor_comparacion" bson:"valor_comparacion"`
	PasoDestinoSiTrue  string    `json:"paso_destino_si_true" bson:"paso_destino_si_true"`
	PasoDestinoSiFalse string    `json:"paso_destino_si_false" bson:"paso_destino_si_false"`
	CreatedAt          time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" bson:"updated_at"`
}
