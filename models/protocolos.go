package models

import (
	"time"
)

// ConfiguracionProtocolo representa la configuración para transferir archivos
type ConfiguracionProtocolo struct {
	ID                  string            `json:"id" bson:"_id,omitempty"`
	EmpresaID           string            `json:"empresa_id" bson:"empresa_id"`
	ERPID               string            `json:"erp_id" bson:"erp_id"`
	Nombre              string            `json:"nombre" bson:"nombre"`
	Descripcion         string            `json:"descripcion" bson:"descripcion"`
	Protocolo           string            `json:"protocolo" bson:"protocolo"`
	Host                string            `json:"host" bson:"host"`
	Puerto              int               `json:"puerto" bson:"puerto"`
	Usuario             string            `json:"usuario" bson:"usuario"`
	Contrasena          string            `json:"contrasena" bson:"contrasena"`
	Ruta                string            `json:"ruta" bson:"ruta"`
	Timeout             int               `json:"timeout" bson:"timeout"`
	UsarTLS             bool              `json:"usar_tls" bson:"usar_tls"`
	TipoAutenticacion   string            `json:"tipo_autenticacion" bson:"tipo_autenticacion"`
	ParametrosExtra     map[string]string `json:"parametros_extra,omitempty" bson:"parametros_extra,omitempty"`
	FrecuenciaEjecucion string            `json:"frecuencia_ejecucion" bson:"frecuencia_ejecucion"`
	UltimaEjecucion     time.Time         `json:"ultima_ejecucion,omitempty" bson:"ultima_ejecucion,omitempty"`
	ProximaEjecucion    time.Time         `json:"proxima_ejecucion,omitempty" bson:"proxima_ejecucion,omitempty"`
	FechaCreacion       time.Time         `json:"fecha_creacion" bson:"fecha_creacion"`
	FechaActualizacion  time.Time         `json:"fecha_actualizacion" bson:"fecha_actualizacion"`
}

// TransformacionLegacy representa una transformación de datos legacy
type TransformacionLegacy struct {
	ID                 string            `json:"id" bson:"_id,omitempty"`
	EmpresaID          string            `json:"empresa_id" bson:"empresa_id"`
	ERPID              string            `json:"erp_id" bson:"erp_id"`
	Nombre             string            `json:"nombre" bson:"nombre"`
	Descripcion        string            `json:"descripcion" bson:"descripcion"`
	CampoOrigen        string            `json:"campo_origen" bson:"campo_origen"`
	CampoDestino       string            `json:"campo_destino" bson:"campo_destino"`
	TipoTransformacion string            `json:"tipo_transformacion" bson:"tipo_transformacion"`
	Parametros         map[string]string `json:"parametros,omitempty" bson:"parametros,omitempty"`
	ExpresionRegular   string            `json:"expresion_regular,omitempty" bson:"expresion_regular,omitempty"`
	ValorPorDefecto    interface{}       `json:"valor_por_defecto,omitempty" bson:"valor_por_defecto,omitempty"`
	Obligatorio        bool              `json:"obligatorio" bson:"obligatorio"`
	Activo             bool              `json:"activo" bson:"activo"`
	FechaCreacion      time.Time         `json:"fecha_creacion" bson:"fecha_creacion"`
	FechaActualizacion time.Time         `json:"fecha_actualizacion" bson:"fecha_actualizacion"`
}

// Constantes para tipos de protocolos
const (
	ProtocoloFTP  = "FTP"
	ProtocoloSFTP = "SFTP"
	ProtocoloFTPS = "FTPS"
)

// Validate valida que todos los campos obligatorios estén presentes
func (c *ConfiguracionProtocolo) Validate() error {
	if c.EmpresaID == "" {
		return &ValidationFieldError{Field: "empresa_id", Message: "El ID de la empresa es obligatorio"}
	}
	if c.ERPID == "" {
		return &ValidationFieldError{Field: "erp_id", Message: "El ID del ERP es obligatorio"}
	}
	if c.Nombre == "" {
		return &ValidationFieldError{Field: "nombre", Message: "El nombre de la configuración es obligatorio"}
	}
	if c.Protocolo == "" {
		return &ValidationFieldError{Field: "protocolo", Message: "El protocolo es obligatorio"}
	}
	if c.Host == "" {
		return &ValidationFieldError{Field: "host", Message: "El host es obligatorio"}
	}
	if c.Puerto <= 0 {
		return &ValidationFieldError{Field: "puerto", Message: "El puerto debe ser mayor que cero"}
	}
	return nil
}

// Validate valida que todos los campos obligatorios estén presentes
func (t *TransformacionLegacy) Validate() error {
	if t.ERPID == "" {
		return &ValidationFieldError{Field: "erp_id", Message: "El ID del ERP es obligatorio"}
	}
	if t.Nombre == "" {
		return &ValidationFieldError{Field: "nombre", Message: "El nombre de la transformación es obligatorio"}
	}
	if t.CampoOrigen == "" {
		return &ValidationFieldError{Field: "campo_origen", Message: "El campo origen es obligatorio"}
	}
	if t.CampoDestino == "" {
		return &ValidationFieldError{Field: "campo_destino", Message: "El campo destino es obligatorio"}
	}
	if t.TipoTransformacion == "" {
		return &ValidationFieldError{Field: "tipo_transformacion", Message: "El tipo de transformación es obligatorio"}
	}
	return nil
}
