package models

import "time"

// DetalleTributario representa un detalle de documento tributario
type DetalleTributario struct {
	Descripcion    string  `json:"descripcion" bson:"descripcion"`
	Cantidad       int     `json:"cantidad" bson:"cantidad"`
	PrecioUnitario float64 `json:"precio_unitario" bson:"precio_unitario"`
	MontoItem      float64 `json:"monto_item" bson:"monto_item"`
	Exento         bool    `json:"exento" bson:"exento"`
}

// DocumentoTributario representa la estructura común para todos los documentos tributarios
type DocumentoTributario struct {
	ID                  string       `json:"id" bson:"_id,omitempty"`
	Folio               int          `json:"folio" bson:"folio"`
	FechaEmision        time.Time    `json:"fecha_emision" bson:"fecha_emision"`
	TipoDocumento       TipoDTE      `json:"tipo_documento" bson:"tipo_documento"`
	TipoDTE             string       `json:"tipo_dte" bson:"tipo_dte"` // Representa el DTE como string para interfaz con SII
	RUTEmisor           string       `json:"rut_emisor" bson:"rut_emisor"`
	RazonSocialEmisor   string       `json:"razon_social_emisor" bson:"razon_social_emisor"`
	GiroEmisor          string       `json:"giro_emisor" bson:"giro_emisor"`
	DireccionEmisor     string       `json:"direccion_emisor" bson:"direccion_emisor"`
	ComunaEmisor        string       `json:"comuna_emisor" bson:"comuna_emisor"`
	RUTReceptor         string       `json:"rut_receptor" bson:"rut_receptor"`
	RazonSocialReceptor string       `json:"razon_social_receptor" bson:"razon_social_receptor"`
	GiroReceptor        string       `json:"giro_receptor,omitempty" bson:"giro_receptor,omitempty"`
	DireccionReceptor   string       `json:"direccion_receptor" bson:"direccion_receptor"`
	ComunaReceptor      string       `json:"comuna_receptor,omitempty" bson:"comuna_receptor,omitempty"`
	MontoNeto           float64      `json:"monto_neto" bson:"monto_neto"`
	MontoExento         float64      `json:"monto_exento" bson:"monto_exento"`
	MontoIVA            float64      `json:"monto_iva" bson:"monto_iva"`
	TasaIVA             float64      `json:"tasa_iva" bson:"tasa_iva"`
	MontoTotal          float64      `json:"monto_total" bson:"monto_total"`
	Referencias         []Referencia `json:"referencias,omitempty" bson:"referencias,omitempty"`
	Estado              EstadoDTE    `json:"estado" bson:"estado"`
	TrackID             string       `json:"track_id,omitempty" bson:"track_id,omitempty"`
	PDF                 string       `json:"pdf,omitempty" bson:"pdf,omitempty"`
	PDFData             []byte       `json:"pdf_data,omitempty" bson:"-"`
	XML                 string       `json:"xml,omitempty" bson:"xml,omitempty"`
	CreatedAt           time.Time    `json:"created_at" bson:"created_at"`
	UpdatedAt           time.Time    `json:"updated_at" bson:"updated_at"`
	Timestamps          Timestamps   `json:"timestamps,omitempty" bson:"timestamps,omitempty"`

	// Campos adicionales para la emisión de documentos
	Emisor   *Emisor             `json:"emisor,omitempty" bson:"emisor,omitempty"`
	Receptor *Receptor           `json:"receptor,omitempty" bson:"receptor,omitempty"`
	Detalles []DetalleTributario `json:"detalles,omitempty" bson:"detalles,omitempty"`
}

// GetField obtiene el valor de un campo
func (d *DocumentoTributario) GetField(campo string) interface{} {
	switch campo {
	case "id":
		return d.ID
	case "tipoDTE":
		return d.TipoDTE
	case "folio":
		return d.Folio
	case "fecha_emision":
		return d.FechaEmision
	case "rut_emisor":
		return d.RUTEmisor
	case "rut_receptor":
		return d.RUTReceptor
	case "razon_social_emisor":
		return d.RazonSocialEmisor
	case "razon_social_receptor":
		return d.RazonSocialReceptor
	case "monto_total":
		return d.MontoTotal
	case "monto_neto":
		return d.MontoNeto
	case "monto_exento":
		return d.MontoExento
	case "monto_iva":
		return d.MontoIVA
	case "estado":
		return d.Estado
	}
	return nil
}

// SetField establece el valor de un campo
func (d *DocumentoTributario) SetField(campo string, valor interface{}) error {
	switch campo {
	case "id":
		if id, ok := valor.(string); ok {
			d.ID = id
			return nil
		}
	case "tipoDTE":
		if tipo, ok := valor.(string); ok {
			d.TipoDTE = tipo
			return nil
		}
	case "folio":
		if folio, ok := valor.(int); ok {
			d.Folio = folio
			return nil
		}
	case "fecha_emision":
		if fecha, ok := valor.(time.Time); ok {
			d.FechaEmision = fecha
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
	case "razon_social_emisor":
		if razon, ok := valor.(string); ok {
			d.RazonSocialEmisor = razon
			return nil
		}
	case "razon_social_receptor":
		if razon, ok := valor.(string); ok {
			d.RazonSocialReceptor = razon
			return nil
		}
	case "monto_total":
		if monto, ok := valor.(float64); ok {
			d.MontoTotal = monto
			return nil
		}
	case "monto_neto":
		if monto, ok := valor.(float64); ok {
			d.MontoNeto = monto
			return nil
		}
	case "monto_exento":
		if monto, ok := valor.(float64); ok {
			d.MontoExento = monto
			return nil
		}
	case "monto_iva":
		if monto, ok := valor.(float64); ok {
			d.MontoIVA = monto
			return nil
		}
	case "estado":
		if estado, ok := valor.(string); ok {
			d.Estado = EstadoDTE(estado)
			return nil
		}
	}
	return NewValidationFieldError(campo, "Tipo de dato inválido para el campo", "INVALID_TYPE", valor)
}
