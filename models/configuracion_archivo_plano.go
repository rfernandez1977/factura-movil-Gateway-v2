package models

import (
	"time"
)

// ConfiguracionArchivoPlano representa la configuración para procesar archivos planos en sistemas legacy
type ConfiguracionArchivoPlano struct {
	ID                  string                 `json:"id" bson:"_id,omitempty"`
	EmpresaID           string                 `json:"empresa_id" bson:"empresa_id"`
	Nombre              string                 `json:"nombre" bson:"nombre"`
	Descripcion         string                 `json:"descripcion" bson:"descripcion"`
	TipoArchivo         string                 `json:"tipo_archivo" bson:"tipo_archivo"` // CSV, TXT, FIXED, etc.
	DelimitadorCampo    string                 `json:"delimitador_campo" bson:"delimitador_campo"`
	DelimitadorLinea    string                 `json:"delimitador_linea" bson:"delimitador_linea"`
	CaracterEscape      string                 `json:"caracter_escape" bson:"caracter_escape"`
	IncluirCabecera     bool                   `json:"incluir_cabecera" bson:"incluir_cabecera"`
	CamposCabecera      []string               `json:"campos_cabecera" bson:"campos_cabecera"`
	DefinicionCampos    []DefinicionCampo      `json:"definicion_campos" bson:"definicion_campos"`
	MapeoCamposDTE      map[string]string      `json:"mapeo_campos_dte" bson:"mapeo_campos_dte"`
	ReglasProcesamiento []ReglaProcesamiento   `json:"reglas_procesamiento" bson:"reglas_procesamiento"`
	ProcesarAutomatico  bool                   `json:"procesar_automatico" bson:"procesar_automatico"`
	DirectorioEntrada   string                 `json:"directorio_entrada" bson:"directorio_entrada"`
	DirectorioSalida    string                 `json:"directorio_salida" bson:"directorio_salida"`
	DirectorioErrores   string                 `json:"directorio_errores" bson:"directorio_errores"`
	PrefijosArchivo     []string               `json:"prefijos_archivo" bson:"prefijos_archivo"`
	ExtensionesArchivo  []string               `json:"extensiones_archivo" bson:"extensiones_archivo"`
	CodigoEmpresa       string                 `json:"codigo_empresa" bson:"codigo_empresa"`
	FormatosPlantilla   map[string]string      `json:"formatos_plantilla" bson:"formatos_plantilla"`
	ParametrosExtras    map[string]interface{} `json:"parametros_extras" bson:"parametros_extras,omitempty"`
	CreatedAt           time.Time              `json:"created_at" bson:"created_at"`
	UpdatedAt           time.Time              `json:"updated_at" bson:"updated_at"`
}

// DefinicionCampo representa la configuración de un campo en un archivo plano
type DefinicionCampo struct {
	Nombre           string `json:"nombre" bson:"nombre"`
	TipoDato         string `json:"tipo_dato" bson:"tipo_dato"` // STRING, INT, FLOAT, DATE, etc.
	Posicion         int    `json:"posicion" bson:"posicion"`
	Longitud         int    `json:"longitud" bson:"longitud"` // Para archivos de ancho fijo
	FormatoFecha     string `json:"formato_fecha" bson:"formato_fecha,omitempty"`
	ValorDefecto     string `json:"valor_defecto" bson:"valor_defecto,omitempty"`
	Requerido        bool   `json:"requerido" bson:"requerido"`
	ExpresionRegular string `json:"expresion_regular" bson:"expresion_regular,omitempty"`
	Descripcion      string `json:"descripcion" bson:"descripcion"`
}

// ReglaProcesamiento representa una regla para el procesamiento de archivos planos
type ReglaProcesamiento struct {
	Nombre       string `json:"nombre" bson:"nombre"`
	Descripcion  string `json:"descripcion" bson:"descripcion"`
	TipoRegla    string `json:"tipo_regla" bson:"tipo_regla"` // VALIDACION, TRANSFORMACION, CALCULO
	Campo        string `json:"campo" bson:"campo"`
	Expresion    string `json:"expresion" bson:"expresion"` // Expresión a evaluar
	MensajeError string `json:"mensaje_error" bson:"mensaje_error,omitempty"`
}

// NewConfiguracionArchivoPlano crea una nueva configuración de archivo plano
func NewConfiguracionArchivoPlano(nombre, descripcion, tipoArchivo, empresaID string) *ConfiguracionArchivoPlano {
	now := time.Now()
	return &ConfiguracionArchivoPlano{
		Nombre:           nombre,
		Descripcion:      descripcion,
		TipoArchivo:      tipoArchivo,
		EmpresaID:        empresaID,
		DefinicionCampos: []DefinicionCampo{},
		MapeoCamposDTE:   make(map[string]string),
		ParametrosExtras: make(map[string]interface{}),
		CreatedAt:        now,
		UpdatedAt:        now,
	}
}
