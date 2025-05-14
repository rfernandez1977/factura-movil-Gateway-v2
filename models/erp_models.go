package models

import "time"

// TipoERP representa el tipo de ERP
type TipoERP string

// Tipos de ERP soportados
const (
	ERP_SAP      TipoERP = "SAP"
	ERP_ORACLE   TipoERP = "ORACLE"
	ERP_DYNAMICS TipoERP = "DYNAMICS"
	ERP_NETSUITE TipoERP = "NETSUITE"
	ERP_LEGACY   TipoERP = "LEGACY"
)

// ConfiguracionERP representa la configuración de conexión a un ERP
type ConfiguracionERP struct {
	ID              string    `json:"id" bson:"_id,omitempty"`
	EmpresaID       string    `json:"empresa_id" bson:"empresa_id"`
	TipoERP         TipoERP   `json:"tipo_erp" bson:"tipo_erp"`
	Nombre          string    `json:"nombre" bson:"nombre"`
	URLConexion     string    `json:"url_conexion" bson:"url_conexion"`
	Usuario         string    `json:"usuario" bson:"usuario"`
	Contrasena      string    `json:"contrasena" bson:"contrasena,omitempty"` // Notar que se omite al serializar a BSON
	Token           string    `json:"token" bson:"token,omitempty"`
	ApiKey          string    `json:"api_key" bson:"api_key,omitempty"`
	ConfiguracionDB string    `json:"configuracion_db" bson:"configuracion_db,omitempty"`
	Parametros      []string  `json:"parametros" bson:"parametros,omitempty"`
	Activo          bool      `json:"activo" bson:"activo"`
	FechaCreacion   time.Time `json:"fecha_creacion" bson:"fecha_creacion"`
	FechaModificado time.Time `json:"fecha_modificado" bson:"fecha_modificado"`
}

// ERPIntegracion representa la configuración de integración con un ERP
type ERPIntegracion struct {
	ID              string            `json:"id" bson:"_id,omitempty"`
	EmpresaID       string            `json:"empresa_id" bson:"empresa_id"`
	ConfiguracionID string            `json:"configuracion_id" bson:"configuracion_id"`
	Tipo            string            `json:"tipo" bson:"tipo"` // Productos, Clientes, Facturas, etc.
	Estado          string            `json:"estado" bson:"estado"`
	Frecuencia      string            `json:"frecuencia" bson:"frecuencia"` // Diaria, Semanal, etc.
	UltimaSync      time.Time         `json:"ultima_sync" bson:"ultima_sync"`
	ProximaSync     time.Time         `json:"proxima_sync" bson:"proxima_sync"`
	Parametros      map[string]string `json:"parametros" bson:"parametros"`
	Activo          bool              `json:"activo" bson:"activo"`
}

// ERPSincronizacion representa un registro de sincronización con un ERP
type ERPSincronizacion struct {
	ID            string    `json:"id" bson:"_id,omitempty"`
	IntegracionID string    `json:"integracion_id" bson:"integracion_id"`
	FechaInicio   time.Time `json:"fecha_inicio" bson:"fecha_inicio"`
	FechaFin      time.Time `json:"fecha_fin" bson:"fecha_fin"`
	Estado        string    `json:"estado" bson:"estado"`
	Resultado     string    `json:"resultado" bson:"resultado"`
	ElementosSync int       `json:"elementos_sync" bson:"elementos_sync"`
	Errores       int       `json:"errores" bson:"errores"`
	Advertencias  int       `json:"advertencias" bson:"advertencias"`
	LogID         string    `json:"log_id" bson:"log_id"`
}
