package models

import (
	"fmt"
	"time"
)

// DocumentoTributario representa un documento tributario electrónico
type DocumentoTributario struct {
	ID                 string                `json:"id" bson:"_id,omitempty"`
	TipoDTE            string                `json:"tipo_dte" bson:"tipo_dte"`
	Tipo               string                `json:"tipo" bson:"tipo"`
	Folio              int                   `json:"folio" bson:"folio"`
	FechaEmision       time.Time             `json:"fecha_emision" bson:"fecha_emision"`
	FechaVencimiento   time.Time             `json:"fecha_vencimiento" bson:"fecha_vencimiento"`
	MontoTotal         float64               `json:"monto_total" bson:"monto_total"`
	Estado             EstadoDTE             `json:"estado" bson:"estado"`
	Emisor             *Emisor               `json:"emisor" bson:"emisor"`
	Receptor           *Receptor             `json:"receptor" bson:"receptor"`
	Detalles           []Detalle             `json:"detalles,omitempty" bson:"detalles,omitempty"`
	EstadoRecepcion    EstadoRecepcion       `json:"estado_recepcion" bson:"estado_recepcion"`
	CertificadoDigital string                `json:"certificado_digital,omitempty" bson:"certificado_digital,omitempty"`
	TrackID            string                `json:"track_id,omitempty" bson:"track_id,omitempty"`
	RutEmisor          string                `json:"rut_emisor" bson:"rut_emisor"`
	RutReceptor        string                `json:"rut_receptor" bson:"rut_receptor"`
	MontoNeto          float64               `json:"monto_neto" bson:"monto_neto"`
	MontoExento        float64               `json:"monto_exento" bson:"monto_exento"`
	MontoIVA           float64               `json:"monto_iva" bson:"monto_iva"`
	Items              []Item                `json:"items" bson:"items"`
	Referencias        []Referencia          `json:"referencias,omitempty" bson:"referencias,omitempty"`
	XML                []byte                `json:"xml,omitempty" bson:"xml,omitempty"`
	PDFData            []byte                `json:"pdf_data,omitempty" bson:"pdf_data,omitempty"`
	EstadoComercial    string                `json:"estado_comercial,omitempty" bson:"estado_comercial,omitempty"`
	EstadoSII          *EstadoSII            `json:"estado_sii,omitempty" bson:"estado_sii,omitempty"`
	Comentarios        []string              `json:"comentarios,omitempty" bson:"comentarios,omitempty"`
	HistorialCambios   []CambioDocumento     `json:"historial_cambios,omitempty" bson:"historial_cambios,omitempty"`
	Validaciones       []ValidacionDocumento `json:"validaciones,omitempty" bson:"validaciones,omitempty"`
	Impuestos          []Impuesto            `json:"impuestos,omitempty" bson:"impuestos,omitempty"`
	MetadatosSII       MetadatosSII          `json:"metadatos_sii,omitempty" bson:"metadatos_sii,omitempty"`
	Timestamps
}

// DocumentoGenerico representa un documento genérico en el sistema
type DocumentoGenerico struct {
	ID          string     `json:"id" bson:"_id,omitempty"`
	Tipo        string     `json:"tipo" bson:"tipo"`
	Nombre      string     `json:"nombre" bson:"nombre"`
	Descripcion string     `json:"descripcion" bson:"descripcion"`
	Contenido   []byte     `json:"contenido" bson:"contenido"`
	MimeType    string     `json:"mime_type" bson:"mime_type"`
	Size        int64      `json:"size" bson:"size"`
	Hash        string     `json:"hash" bson:"hash"`
	Metadata    Metadata   `json:"metadata" bson:"metadata"`
	Timestamps  Timestamps `json:"timestamps" bson:"timestamps"`
}

// DocumentoAlmacenado representa un documento almacenado en el sistema
type DocumentoAlmacenado struct {
	ID          string     `json:"id" bson:"_id,omitempty"`
	Tipo        string     `json:"tipo" bson:"tipo"`
	Nombre      string     `json:"nombre" bson:"nombre"`
	Descripcion string     `json:"descripcion" bson:"descripcion"`
	Ruta        string     `json:"ruta" bson:"ruta"`
	MimeType    string     `json:"mime_type" bson:"mime_type"`
	Size        int64      `json:"size" bson:"size"`
	Hash        string     `json:"hash" bson:"hash"`
	Metadata    Metadata   `json:"metadata" bson:"metadata"`
	Timestamps  Timestamps `json:"timestamps" bson:"timestamps"`
	Contenido   []byte     `json:"contenido,omitempty" bson:"contenido,omitempty"`
	CacheInfo   CacheInfo  `json:"-" bson:"-"`
}

// DocumentoSeguro representa un documento con características de seguridad
type DocumentoSeguro struct {
	ID          string     `json:"id" bson:"_id,omitempty"`
	Tipo        string     `json:"tipo" bson:"tipo"`
	Nombre      string     `json:"nombre" bson:"nombre"`
	Descripcion string     `json:"descripcion" bson:"descripcion"`
	Contenido   []byte     `json:"contenido" bson:"contenido"`
	MimeType    string     `json:"mime_type" bson:"mime_type"`
	Size        int64      `json:"size" bson:"size"`
	Hash        string     `json:"hash" bson:"hash"`
	Metadata    Metadata   `json:"metadata" bson:"metadata"`
	Encriptado  bool       `json:"encriptado" bson:"encriptado"`
	Firmado     bool       `json:"firmado" bson:"firmado"`
	Certificado bool       `json:"certificado" bson:"certificado"`
	Timestamps  Timestamps `json:"timestamps" bson:"timestamps"`
}

// Tipos auxiliares
type Impuesto struct {
	Tipo  string  `json:"tipo" bson:"tipo"`
	Base  float64 `json:"base" bson:"base"`
	Tasa  float64 `json:"tasa" bson:"tasa"`
	Monto float64 `json:"monto" bson:"monto"`
}

type Metadata struct {
	Version   string            `json:"version" bson:"version"`
	Tags      []string          `json:"tags" bson:"tags"`
	Atributos map[string]string `json:"atributos" bson:"atributos"`
}

type Timestamps struct {
	Creado     time.Time `json:"creado" bson:"creado"`
	Modificado time.Time `json:"modificado" bson:"modificado"`
	Eliminado  time.Time `json:"eliminado,omitempty" bson:"eliminado,omitempty"`
}

type CambioDocumento struct {
	Timestamp     time.Time   `json:"timestamp" bson:"timestamp"`
	TipoCambio    string      `json:"tipo_cambio" bson:"tipo_cambio"`
	CampoAfectado string      `json:"campo_afectado" bson:"campo_afectado"`
	ValorAnterior interface{} `json:"valor_anterior" bson:"valor_anterior"`
	ValorNuevo    interface{} `json:"valor_nuevo" bson:"valor_nuevo"`
	Usuario       string      `json:"usuario" bson:"usuario"`
	Motivo        string      `json:"motivo" bson:"motivo"`
}

type ValidacionDocumento struct {
	Resultado       bool      `json:"resultado" bson:"resultado"`
	Timestamp       time.Time `json:"timestamp" bson:"timestamp"`
	DetalleError    string    `json:"detalle_error,omitempty" bson:"detalle_error,omitempty"`
	NivelCriticidad string    `json:"nivel_criticidad" bson:"nivel_criticidad"`
}

type MetadatosSII struct {
	TrackID          string    `json:"track_id" bson:"track_id"`
	EstadoSII        string    `json:"estado_sii" bson:"estado_sii"`
	FechaRecepcion   time.Time `json:"fecha_recepcion" bson:"fecha_recepcion"`
	NumeroAtencion   string    `json:"numero_atencion" bson:"numero_atencion"`
	ObservacionesSII []string  `json:"observaciones_sii" bson:"observaciones_sii"`
}

// CacheInfo contiene información sobre el caché de un documento
type CacheInfo struct {
	CreatedAt  time.Time `json:"created_at"`
	ExpiresAt  time.Time `json:"expires_at"`
	LastAccess time.Time `json:"last_access"`
}

// GetField obtiene el valor de un campo del documento por su nombre
func (d *DocumentoTributario) GetField(field string) interface{} {
	switch field {
	case "id":
		return d.ID
	case "tipo_dte":
		return d.TipoDTE
	case "folio":
		return d.Folio
	case "fecha_emision":
		return d.FechaEmision
	case "fecha_vencimiento":
		return d.FechaVencimiento
	case "monto_total":
		return d.MontoTotal
	case "estado":
		return d.Estado
	case "emisor":
		return d.Emisor
	case "receptor":
		return d.Receptor
	case "detalles":
		return d.Detalles
	case "impuestos":
		return d.Impuestos
	case "referencias":
		return d.Referencias
	case "xml":
		return d.XML
	case "pdf":
		return d.PDFData
	case "timestamps":
		return d.Timestamps
	case "metadatos_sii":
		return d.MetadatosSII
	case "historial_cambios":
		return d.HistorialCambios
	case "validaciones":
		return d.Validaciones
	case "track_id":
		return d.TrackID
	case "rut_emisor":
		return d.RutEmisor
	case "rut_receptor":
		return d.RutReceptor
	case "monto_neto":
		return d.MontoNeto
	case "monto_exento":
		return d.MontoExento
	case "monto_iva":
		return d.MontoIVA
	case "items":
		return d.Items
	case "estado_comercial":
		return d.EstadoComercial
	case "estado_sii":
		return d.EstadoSII
	case "comentarios":
		return d.Comentarios
	default:
		return nil
	}
}

// SetField establece el valor de un campo del documento por su nombre
func (d *DocumentoTributario) SetField(field string, value interface{}) error {
	switch field {
	case "id":
		if str, ok := value.(string); ok {
			d.ID = str
		} else {
			return fmt.Errorf("tipo inválido para el campo id")
		}
	case "tipo_dte":
		if str, ok := value.(string); ok {
			d.TipoDTE = str
		} else {
			return fmt.Errorf("tipo inválido para el campo tipo_dte")
		}
	case "folio":
		if num, ok := value.(int); ok {
			d.Folio = num
		} else {
			return fmt.Errorf("tipo inválido para el campo folio")
		}
	case "fecha_emision":
		if t, ok := value.(time.Time); ok {
			d.FechaEmision = t
		} else {
			return fmt.Errorf("tipo inválido para el campo fecha_emision")
		}
	case "fecha_vencimiento":
		if t, ok := value.(time.Time); ok {
			d.FechaVencimiento = t
		} else {
			return fmt.Errorf("tipo inválido para el campo fecha_vencimiento")
		}
	case "monto_total":
		if num, ok := value.(float64); ok {
			d.MontoTotal = num
		} else {
			return fmt.Errorf("tipo inválido para el campo monto_total")
		}
	case "estado":
		if str, ok := value.(EstadoDTE); ok {
			d.Estado = str
		} else {
			return fmt.Errorf("tipo inválido para el campo estado")
		}
	case "emisor":
		if e, ok := value.(Emisor); ok {
			d.Emisor = &e
		} else {
			return fmt.Errorf("tipo inválido para el campo emisor")
		}
	case "receptor":
		if r, ok := value.(Receptor); ok {
			d.Receptor = &r
		} else {
			return fmt.Errorf("tipo inválido para el campo receptor")
		}
	case "detalles":
		if detalles, ok := value.([]Detalle); ok {
			d.Detalles = detalles
		} else {
			return fmt.Errorf("tipo inválido para el campo detalles")
		}
	case "impuestos":
		if i, ok := value.([]Impuesto); ok {
			d.Impuestos = i
		} else {
			return fmt.Errorf("tipo inválido para el campo impuestos")
		}
	case "referencias":
		if r, ok := value.([]Referencia); ok {
			d.Referencias = r
		} else {
			return fmt.Errorf("tipo inválido para el campo referencias")
		}
	case "xml":
		if b, ok := value.([]byte); ok {
			d.XML = b
		} else {
			return fmt.Errorf("tipo inválido para el campo xml")
		}
	case "pdf":
		if b, ok := value.([]byte); ok {
			d.PDFData = b
		} else {
			return fmt.Errorf("tipo inválido para el campo pdf")
		}
	case "timestamps":
		if t, ok := value.(Timestamps); ok {
			d.Timestamps = t
		} else {
			return fmt.Errorf("tipo inválido para el campo timestamps")
		}
	case "metadatos_sii":
		if m, ok := value.(MetadatosSII); ok {
			d.MetadatosSII = m
		} else {
			return fmt.Errorf("tipo inválido para el campo metadatos_sii")
		}
	case "historial_cambios":
		if h, ok := value.([]CambioDocumento); ok {
			d.HistorialCambios = h
		} else {
			return fmt.Errorf("tipo inválido para el campo historial_cambios")
		}
	case "validaciones":
		if v, ok := value.([]ValidacionDocumento); ok {
			d.Validaciones = v
		} else {
			return fmt.Errorf("tipo inválido para el campo validaciones")
		}
	case "track_id":
		if str, ok := value.(string); ok {
			d.TrackID = str
		} else {
			return fmt.Errorf("tipo inválido para el campo track_id")
		}
	case "rut_emisor":
		if str, ok := value.(string); ok {
			d.RutEmisor = str
		} else {
			return fmt.Errorf("tipo inválido para el campo rut_emisor")
		}
	case "rut_receptor":
		if str, ok := value.(string); ok {
			d.RutReceptor = str
		} else {
			return fmt.Errorf("tipo inválido para el campo rut_receptor")
		}
	case "monto_neto":
		if num, ok := value.(float64); ok {
			d.MontoNeto = num
		} else {
			return fmt.Errorf("tipo inválido para el campo monto_neto")
		}
	case "monto_exento":
		if num, ok := value.(float64); ok {
			d.MontoExento = num
		} else {
			return fmt.Errorf("tipo inválido para el campo monto_exento")
		}
	case "monto_iva":
		if num, ok := value.(float64); ok {
			d.MontoIVA = num
		} else {
			return fmt.Errorf("tipo inválido para el campo monto_iva")
		}
	case "items":
		if items, ok := value.([]Item); ok {
			d.Items = items
		} else {
			return fmt.Errorf("tipo inválido para el campo items")
		}
	case "estado_comercial":
		if str, ok := value.(string); ok {
			d.EstadoComercial = str
		} else {
			return fmt.Errorf("tipo inválido para el campo estado_comercial")
		}
	case "estado_sii":
		if es, ok := value.(EstadoSII); ok {
			d.EstadoSII = &es
		} else {
			return fmt.Errorf("tipo inválido para el campo estado_sii")
		}
	case "comentarios":
		if comments, ok := value.([]string); ok {
			d.Comentarios = comments
		} else {
			return fmt.Errorf("tipo inválido para el campo comentarios")
		}
	default:
		return fmt.Errorf("campo desconocido: %s", field)
	}
	return nil
}
