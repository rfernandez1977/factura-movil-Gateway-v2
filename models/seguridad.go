package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Usuario representa un usuario del sistema
type Usuario struct {
	ID                string    `json:"id" bson:"_id"`
	Rut               string    `json:"rut" bson:"rut"`
	Nombre            string    `json:"nombre" bson:"nombre"`
	Email             string    `json:"email" bson:"email"`
	HashContrasena    string    `json:"-" bson:"hash_contrasena"`
	Salt              string    `json:"-" bson:"salt"`
	Roles             []string  `json:"roles" bson:"roles"`
	Permisos          []string  `json:"permisos" bson:"permisos"`
	Estado            string    `json:"estado" bson:"estado"` // ACTIVO, INACTIVO, BLOQUEADO
	UltimoAcceso      time.Time `json:"ultimo_acceso" bson:"ultimo_acceso"`
	IntentosFallidos  int       `json:"intentos_fallidos" bson:"intentos_fallidos"`
	FechaCreacion     time.Time `json:"fecha_creacion" bson:"fecha_creacion"`
	FechaModificacion time.Time `json:"fecha_modificacion" bson:"fecha_modificacion"`
}

// RegistroAuditoriaAcceso representa un registro de acceso al sistema
type RegistroAuditoriaAcceso struct {
	ID          string    `json:"id" bson:"_id"`
	UsuarioID   string    `json:"usuario_id" bson:"usuario_id"`
	Rut         string    `json:"rut" bson:"rut"`
	Accion      string    `json:"accion" bson:"accion"` // LOGIN, LOGOUT, CAMBIO_CONTRASENA
	IP          string    `json:"ip" bson:"ip"`
	UserAgent   string    `json:"user_agent" bson:"user_agent"`
	Exitoso     bool      `json:"exitoso" bson:"exitoso"`
	Detalles    string    `json:"detalles" bson:"detalles"`
	FechaAcceso time.Time `json:"fecha_acceso" bson:"fecha_acceso"`
}

// RegistroAuditoriaOperacion representa un registro de operaciones en el sistema
type RegistroAuditoriaOperacion struct {
	ID             string                 `json:"id" bson:"_id"`
	UsuarioID      string                 `json:"usuario_id" bson:"usuario_id"`
	Rut            string                 `json:"rut" bson:"rut"`
	Operacion      string                 `json:"operacion" bson:"operacion"`
	Entidad        string                 `json:"entidad" bson:"entidad"`
	EntidadID      string                 `json:"entidad_id" bson:"entidad_id"`
	Cambios        map[string]interface{} `json:"cambios" bson:"cambios"`
	EstadoAnterior map[string]interface{} `json:"estado_anterior" bson:"estado_anterior"`
	EstadoNuevo    map[string]interface{} `json:"estado_nuevo" bson:"estado_nuevo"`
	IP             string                 `json:"ip" bson:"ip"`
	UserAgent      string                 `json:"user_agent" bson:"user_agent"`
	FechaOperacion time.Time              `json:"fecha_operacion" bson:"fecha_operacion"`
}

// FirmaDigital representa una firma digital
type FirmaDigital struct {
	ID                string    `json:"id" bson:"_id"`
	UsuarioID         string    `json:"usuario_id" bson:"usuario_id"`
	Rut               string    `json:"rut" bson:"rut"`
	Certificado       []byte    `json:"certificado" bson:"certificado"`
	ClavePrivada      []byte    `json:"-" bson:"clave_privada"`
	ClavePublica      []byte    `json:"clave_publica" bson:"clave_publica"`
	VigenciaDesde     time.Time `json:"vigencia_desde" bson:"vigencia_desde"`
	VigenciaHasta     time.Time `json:"vigencia_hasta" bson:"vigencia_hasta"`
	Estado            string    `json:"estado" bson:"estado"` // ACTIVA, REVOCADA, VENCIDA
	FechaCreacion     time.Time `json:"fecha_creacion" bson:"fecha_creacion"`
	FechaModificacion time.Time `json:"fecha_modificacion" bson:"fecha_modificacion"`
}

// DatosEncriptados representa datos sensibles encriptados
type DatosEncriptados struct {
	ID                string    `json:"id" bson:"_id"`
	Entidad           string    `json:"entidad" bson:"entidad"`
	EntidadID         string    `json:"entidad_id" bson:"entidad_id"`
	Campo             string    `json:"campo" bson:"campo"`
	ValorEncriptado   []byte    `json:"valor_encriptado" bson:"valor_encriptado"`
	IV                []byte    `json:"iv" bson:"iv"`
	Algoritmo         string    `json:"algoritmo" bson:"algoritmo"`
	Version           int       `json:"version" bson:"version"`
	FechaCreacion     time.Time `json:"fecha_creacion" bson:"fecha_creacion"`
	FechaModificacion time.Time `json:"fecha_modificacion" bson:"fecha_modificacion"`
}

// RegistroAcceso representa un registro de acceso al sistema
type RegistroAcceso struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UsuarioID primitive.ObjectID `json:"usuario_id" bson:"usuario_id"`
	IP        string             `json:"ip" bson:"ip"`
	Fecha     time.Time          `json:"fecha" bson:"fecha"`
	Exitoso   bool               `json:"exitoso" bson:"exitoso"`
	Metodo    string             `json:"metodo" bson:"metodo"`
}

// ReporteSeguridad representa un reporte de seguridad
type ReporteSeguridad struct {
	ID              primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	PeriodoInicio   time.Time          `json:"periodo_inicio" bson:"periodo_inicio"`
	PeriodoFin      time.Time          `json:"periodo_fin" bson:"periodo_fin"`
	TotalAccesos    int                `json:"total_accesos" bson:"total_accesos"`
	AccesosExitosos int                `json:"accesos_exitosos" bson:"accesos_exitosos"`
	AccesosFallidos int                `json:"accesos_fallidos" bson:"accesos_fallidos"`
	IntentosPorIP   map[string]int     `json:"intentos_por_ip" bson:"intentos_por_ip"`
	FechaGeneracion time.Time          `json:"fecha_generacion" bson:"fecha_generacion"`
}

// AlertaSeguridad representa una alerta de seguridad
type AlertaSeguridad struct {
	ID          string    `json:"id" bson:"_id"`
	Tipo        string    `json:"tipo" bson:"tipo"`
	Severidad   string    `json:"severidad" bson:"severidad"` // BAJA, MEDIA, ALTA, CRITICA
	Descripcion string    `json:"descripcion" bson:"descripcion"`
	Detalles    string    `json:"detalles" bson:"detalles"`
	Estado      string    `json:"estado" bson:"estado"` // PENDIENTE, RESUELTA, DESCARTADA
	FechaAlerta time.Time `json:"fecha_alerta" bson:"fecha_alerta"`
}
