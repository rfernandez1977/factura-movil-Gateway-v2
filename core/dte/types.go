package dte

import (
	"time"
)

// DTE representa un documento tributario electrónico
type DTE struct {
	ID            string    `json:"id"`
	Documento     Documento `json:"documento"`
	FechaCreacion time.Time `json:"fecha_creacion"`
	Estado        string    `json:"estado"`
	Firmado       bool      `json:"firmado"`
	XML           string    `json:"xml,omitempty"`
	XMLFirmado    string    `json:"xml_firmado,omitempty"`
}

// Documento representa el contenido de un documento tributario
type Documento struct {
	Encabezado  Encabezado   `json:"encabezado"`
	Detalles    []Detalle    `json:"detalles"`
	Referencias []Referencia `json:"referencias,omitempty"`
	Descuentos  []Descuento  `json:"descuentos,omitempty"`
	Recargos    []Recargo    `json:"recargos,omitempty"`
}

// Encabezado contiene la información principal del documento
type Encabezado struct {
	IDDocumento IDDocumento `json:"id_documento"`
	Emisor      Emisor      `json:"emisor"`
	Receptor    Receptor    `json:"receptor"`
	Totales     Totales     `json:"totales"`
}

// IDDocumento contiene la identificación del documento
type IDDocumento struct {
	TipoDTE          string    `json:"tipo_dte"`
	Folio            int       `json:"folio"`
	FechaEmision     time.Time `json:"fecha_emision"`
	FechaVencimiento time.Time `json:"fecha_vencimiento,omitempty"`
	FormaPago        string    `json:"forma_pago,omitempty"`
}

// Emisor contiene la información del emisor del documento
type Emisor struct {
	RUT         string   `json:"rut"`
	RazonSocial string   `json:"razon_social"`
	Giro        string   `json:"giro"`
	Direccion   string   `json:"direccion"`
	Comuna      string   `json:"comuna"`
	Ciudad      string   `json:"ciudad"`
	Telefono    string   `json:"telefono,omitempty"`
	Email       string   `json:"email,omitempty"`
	SucursalSII string   `json:"sucursal_sii,omitempty"`
	Acteco      []string `json:"acteco,omitempty"`
}

// Receptor contiene la información del receptor del documento
type Receptor struct {
	RUT         string `json:"rut"`
	RazonSocial string `json:"razon_social"`
	Giro        string `json:"giro"`
	Direccion   string `json:"direccion"`
	Comuna      string `json:"comuna"`
	Ciudad      string `json:"ciudad"`
	Telefono    string `json:"telefono,omitempty"`
	Email       string `json:"email,omitempty"`
	Contacto    string `json:"contacto,omitempty"`
}

// Totales contiene los montos totales del documento
type Totales struct {
	MontoNeto      float64    `json:"monto_neto"`
	MontoExento    float64    `json:"monto_exento"`
	TasaIVA        float64    `json:"tasa_iva"`
	IVA            float64    `json:"iva"`
	MontoTotal     float64    `json:"monto_total"`
	OtrosImpuestos []Impuesto `json:"otros_impuestos,omitempty"`
}

// Detalle representa un ítem del documento
type Detalle struct {
	NumeroLinea int     `json:"numero_linea"`
	Nombre      string  `json:"nombre"`
	Descripcion string  `json:"descripcion,omitempty"`
	Cantidad    float64 `json:"cantidad"`
	Unidad      string  `json:"unidad,omitempty"`
	Precio      float64 `json:"precio"`
	Descuento   float64 `json:"descuento,omitempty"`
	Recargo     float64 `json:"recargo,omitempty"`
	MontoItem   float64 `json:"monto_item"`
	Exento      bool    `json:"exento"`
}

// Referencia representa una referencia a otro documento
type Referencia struct {
	NumeroLinea   int       `json:"numero_linea"`
	TipoDocumento string    `json:"tipo_documento"`
	Folio         string    `json:"folio"`
	Fecha         time.Time `json:"fecha"`
	Codigo        string    `json:"codigo,omitempty"`
	Razon         string    `json:"razon"`
}

// Descuento representa un descuento global aplicado al documento
type Descuento struct {
	NumeroLinea int     `json:"numero_linea"`
	Tipo        string  `json:"tipo"` // Porcentaje o Monto
	Valor       float64 `json:"valor"`
	Glosa       string  `json:"glosa,omitempty"`
}

// Recargo representa un recargo global aplicado al documento
type Recargo struct {
	NumeroLinea int     `json:"numero_linea"`
	Tipo        string  `json:"tipo"` // Porcentaje o Monto
	Valor       float64 `json:"valor"`
	Glosa       string  `json:"glosa,omitempty"`
}

// Impuesto representa un impuesto aplicado al documento
type Impuesto struct {
	Tipo  string  `json:"tipo"`
	Tasa  float64 `json:"tasa"`
	Monto float64 `json:"monto"`
}
