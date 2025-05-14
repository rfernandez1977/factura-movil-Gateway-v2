package models

import "time"

// DocumentoAlmacenado representa un documento almacenado en la base de datos
type DocumentoAlmacenado struct {
	ID             string                 `json:"id" bson:"_id,omitempty"`
	TipoDocumento  string                 `json:"tipo_documento" bson:"tipo_documento"`
	Folio          int                    `json:"folio" bson:"folio"`
	RUTEmisor      string                 `json:"rut_emisor" bson:"rut_emisor"`
	RUTReceptor    string                 `json:"rut_receptor" bson:"rut_receptor"`
	FechaEmision   time.Time              `json:"fecha_emision" bson:"fecha_emision"`
	FechaRecepcion time.Time              `json:"fecha_recepcion" bson:"fecha_recepcion"`
	MontoTotal     float64                `json:"monto_total" bson:"monto_total"`
	XML            string                 `json:"xml,omitempty" bson:"xml,omitempty"`
	PDF            []byte                 `json:"pdf,omitempty" bson:"pdf,omitempty"`
	Estado         string                 `json:"estado" bson:"estado"`
	Origen         string                 `json:"origen" bson:"origen"` // API, Email, SII, etc.
	Metadata       map[string]interface{} `json:"metadata,omitempty" bson:"metadata,omitempty"`
	Validaciones   []ValidationResult     `json:"validaciones,omitempty" bson:"validaciones,omitempty"`
	CreatedAt      time.Time              `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at" bson:"updated_at"`
	CacheInfo      CacheInfo              `json:"cache_info,omitempty" bson:"cache_info,omitempty"`
}

// CacheInfo contiene información sobre el almacenamiento en caché
type CacheInfo struct {
	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	ExpiresAt time.Time `json:"expires_at,omitempty" bson:"expires_at,omitempty"`
}

// ValidationResult contiene el resultado de la validación de un documento
type ValidationResult struct {
	Tipo        string                  `json:"tipo" bson:"tipo"`
	Estado      string                  `json:"estado" bson:"estado"`
	Mensaje     string                  `json:"mensaje" bson:"mensaje"`
	Errores     []*ValidationFieldError `json:"errores,omitempty" bson:"errores,omitempty"`
	Timestamp   time.Time               `json:"timestamp" bson:"timestamp"`
	ValidadorID string                  `json:"validador_id,omitempty" bson:"validador_id,omitempty"`
	Metadata    map[string]interface{}  `json:"metadata,omitempty" bson:"metadata,omitempty"`
}

// GetField obtiene el valor de un campo de la estructura
func (d *DocumentoAlmacenado) GetField(campo string) interface{} {
	switch campo {
	case "id":
		return d.ID
	case "tipo_documento":
		return d.TipoDocumento
	case "folio":
		return d.Folio
	case "rut_emisor":
		return d.RUTEmisor
	case "rut_receptor":
		return d.RUTReceptor
	case "fecha_emision":
		return d.FechaEmision
	case "fecha_recepcion":
		return d.FechaRecepcion
	case "monto_total":
		return d.MontoTotal
	case "estado":
		return d.Estado
	case "origen":
		return d.Origen
	case "created_at":
		return d.CreatedAt
	case "updated_at":
		return d.UpdatedAt
	}

	// Si no se encuentra el campo específico, buscar en metadata
	if d.Metadata != nil {
		if valor, ok := d.Metadata[campo]; ok {
			return valor
		}
	}

	return nil
}

// SetField establece el valor de un campo
func (d *DocumentoAlmacenado) SetField(campo string, valor interface{}) error {
	switch campo {
	case "id":
		if id, ok := valor.(string); ok {
			d.ID = id
			return nil
		}
	case "tipo_documento":
		if tipo, ok := valor.(string); ok {
			d.TipoDocumento = tipo
			return nil
		}
	case "folio":
		if folio, ok := valor.(int); ok {
			d.Folio = folio
			return nil
		}
	case "rut_emisor":
		if rut, ok := valor.(string); ok {
			d.RUTEmisor = rut
			return nil
		}
	case "rut_receptor":
		if rut, ok := valor.(string); ok {
			d.RUTReceptor = rut
			return nil
		}
	case "fecha_emision":
		if fecha, ok := valor.(time.Time); ok {
			d.FechaEmision = fecha
			return nil
		}
	case "fecha_recepcion":
		if fecha, ok := valor.(time.Time); ok {
			d.FechaRecepcion = fecha
			return nil
		}
	case "monto_total":
		if monto, ok := valor.(float64); ok {
			d.MontoTotal = monto
			return nil
		}
	case "estado":
		if estado, ok := valor.(string); ok {
			d.Estado = estado
			return nil
		}
	case "origen":
		if origen, ok := valor.(string); ok {
			d.Origen = origen
			return nil
		}
	default:
		// Si no es un campo específico, almacenar en metadata
		if d.Metadata == nil {
			d.Metadata = make(map[string]interface{})
		}
		d.Metadata[campo] = valor
		return nil
	}

	return NewValidationFieldError(campo, "Tipo de dato inválido para el campo", "INVALID_TYPE", valor)
}
