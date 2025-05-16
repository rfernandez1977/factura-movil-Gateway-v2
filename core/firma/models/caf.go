package models

import (
	"time"
)

// CAF representa un Código de Autorización de Folios
type CAF struct {
	ID               string    `json:"id"`
	TipoDocumento    string    `json:"tipo_documento"`
	RutEmisor        string    `json:"rut_emisor"`
	RazonSocial      string    `json:"razon_social"`
	FolioInicial     int64     `json:"folio_inicial"`
	FolioFinal       int64     `json:"folio_final"`
	FolioUltimo      int64     `json:"folio_ultimo"`
	FechaEmision     time.Time `json:"fecha_emision"`
	FechaVencimiento time.Time `json:"fecha_vencimiento"`
	Estado           string    `json:"estado"` // ACTIVO, VENCIDO, AGOTADO
	XML              []byte    `json:"xml"`    // XML original del CAF
	FirmaSII         []byte    `json:"firma_sii"`
}

// CAFBackup representa una copia de respaldo de un CAF
type CAFBackup struct {
	CAF           CAF       `json:"caf"`
	FechaBackup   time.Time `json:"fecha_backup"`
	HashContenido string    `json:"hash_contenido"`
	Ubicacion     string    `json:"ubicacion"`
}

// CAFMetadata contiene información adicional sobre el uso del CAF
type CAFMetadata struct {
	ID                    string    `json:"id"`
	CAFID                 string    `json:"caf_id"`
	FoliosUsados          int       `json:"folios_usados"`
	UltimoUso             time.Time `json:"ultimo_uso"`
	PromedioUso           float64   `json:"promedio_uso"`
	EstimacionAgotamiento time.Time `json:"estimacion_agotamiento"`
}

// FoliosDisponibles calcula la cantidad de folios disponibles
func (c *CAF) FoliosDisponibles() int64 {
	if c.FolioUltimo == 0 {
		return c.FolioFinal - c.FolioInicial + 1
	}
	return c.FolioFinal - c.FolioUltimo
}

// Contienefolio verifica si el CAF contiene un folio específico
func (c *CAF) ContieneFolio(folio int64) bool {
	return folio >= c.FolioInicial && folio <= c.FolioFinal
}

// EstaVigente verifica si el CAF está vigente
func (c *CAF) EstaVigente() bool {
	return time.Now().Before(c.FechaVencimiento)
}

// TieneFoliosDisponibles verifica si quedan folios disponibles
func (c *CAF) TieneFoliosDisponibles() bool {
	return c.FoliosDisponibles() > 0
}

// DisponibilidadCAF representa la disponibilidad de CAF para un tipo de documento
type DisponibilidadCAF struct {
	TipoDocumento     string    `json:"tipo_documento"`
	FoliosDisponibles int64     `json:"folios_disponibles"`
	CAFsActivos       int       `json:"cafs_activos"`
	FechaConsulta     time.Time `json:"fecha_consulta"`
}

// Alerta representa una alerta del sistema
type Alerta struct {
	ID        string                 `json:"id"`
	Tipo      string                 `json:"tipo"`
	Mensaje   string                 `json:"mensaje"`
	Detalles  map[string]interface{} `json:"detalles,omitempty"`
	Fecha     time.Time              `json:"fecha"`
	Estado    string                 `json:"estado"`    // NUEVA, LEIDA, RESUELTA
	Prioridad string                 `json:"prioridad"` // BAJA, MEDIA, ALTA, CRITICA
}
