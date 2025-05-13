package models

import (
	"time"
)

// ConfiguracionArchivoPlano representa la configuración para procesar archivos planos
type ConfiguracionArchivoPlano struct {
	ID                      string                `json:"id" bson:"_id,omitempty"`
	EmpresaID               string                `json:"empresa_id" bson:"empresa_id"`
	Nombre                  string                `json:"nombre" bson:"nombre"`
	Descripcion             string                `json:"descripcion" bson:"descripcion"`
	TipoArchivo             string                `json:"tipo_archivo" bson:"tipo_archivo"`
	Separador               string                `json:"separador" bson:"separador"`
	FormatoFecha            string                `json:"formato_fecha" bson:"formato_fecha"`
	OmitirCabecera          bool                  `json:"omitir_cabecera" bson:"omitir_cabecera"`
	Mappings                []MappingArchivoCampo `json:"mappings" bson:"mappings"`
	Activo                  bool                  `json:"activo" bson:"activo"`
	ValidacionesPreProceso  []ValidacionArchivo   `json:"validaciones_pre_proceso,omitempty" bson:"validaciones_pre_proceso,omitempty"`
	ValidacionesPostProceso []ValidacionArchivo   `json:"validaciones_post_proceso,omitempty" bson:"validaciones_post_proceso,omitempty"`
	CreatedAt               time.Time             `json:"created_at" bson:"created_at"`
	UpdatedAt               time.Time             `json:"updated_at" bson:"updated_at"`
}

// MappingArchivoCampo representa el mapeo entre un campo del archivo y un campo del modelo
type MappingArchivoCampo struct {
	IndiceColumna        int         `json:"indice_columna" bson:"indice_columna"`
	NombreCampo          string      `json:"nombre_campo" bson:"nombre_campo"`
	TipoDato             string      `json:"tipo_dato" bson:"tipo_dato"`
	ObligatorioEnArchivo bool        `json:"obligatorio_en_archivo" bson:"obligatorio_en_archivo"`
	ObligatorioEnModelo  bool        `json:"obligatorio_en_modelo" bson:"obligatorio_en_modelo"`
	ValorPorDefecto      interface{} `json:"valor_por_defecto,omitempty" bson:"valor_por_defecto,omitempty"`
}

// ValidacionArchivo representa una validación para un archivo plano
type ValidacionArchivo struct {
	Tipo      string `json:"tipo" bson:"tipo"`
	Expresion string `json:"expresion" bson:"expresion"`
	Mensaje   string `json:"mensaje" bson:"mensaje"`
}

// Validate valida que todos los campos obligatorios estén presentes
func (c *ConfiguracionArchivoPlano) Validate() error {
	if c.EmpresaID == "" {
		return &ValidationFieldError{Field: "empresa_id", Message: "El ID de la empresa es obligatorio"}
	}
	if c.Nombre == "" {
		return &ValidationFieldError{Field: "nombre", Message: "El nombre de la configuración es obligatorio"}
	}
	if c.TipoArchivo == "" {
		return &ValidationFieldError{Field: "tipo_archivo", Message: "El tipo de archivo es obligatorio"}
	}
	if c.Separador == "" {
		return &ValidationFieldError{Field: "separador", Message: "El separador es obligatorio"}
	}
	if len(c.Mappings) == 0 {
		return &ValidationFieldError{Field: "mappings", Message: "Debe existir al menos un mapeo de campos"}
	}
	return nil
}
