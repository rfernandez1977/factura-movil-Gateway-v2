package models

import "time"

// GuiaDespacho representa una guía de despacho electrónica
type GuiaDespacho struct {
	DocumentoTributario
	IndicadorTraslado          string    `json:"indicador_traslado" bson:"indicador_traslado"`
	IndicadorServicio          string    `json:"indicador_servicio,omitempty" bson:"indicador_servicio,omitempty"`
	IndicadorVentas            string    `json:"indicador_ventas,omitempty" bson:"indicador_ventas,omitempty"`
	IndicadorTransporte        string    `json:"indicador_transporte,omitempty" bson:"indicador_transporte,omitempty"`
	IndicadorExportacion       string    `json:"indicador_exportacion,omitempty" bson:"indicador_exportacion,omitempty"`
	RutTransportista           string    `json:"rut_transportista" bson:"rut_transportista"`
	RazonSocialTransportista   string    `json:"razon_social_transportista" bson:"razon_social_transportista"`
	Patente                    string    `json:"patente,omitempty" bson:"patente,omitempty"`
	FechaInicioTransporte      time.Time `json:"fecha_inicio_transporte,omitempty" bson:"fecha_inicio_transporte,omitempty"`
	FechaFinTransporte         time.Time `json:"fecha_fin_transporte,omitempty" bson:"fecha_fin_transporte,omitempty"`
	FechaInicioServicio        time.Time `json:"fecha_inicio_servicio,omitempty" bson:"fecha_inicio_servicio,omitempty"`
	FechaFinServicio           time.Time `json:"fecha_fin_servicio,omitempty" bson:"fecha_fin_servicio,omitempty"`
	Periodicidad               string    `json:"periodicidad,omitempty" bson:"periodicidad,omitempty"`
	CodigoAduana               string    `json:"codigo_aduana,omitempty" bson:"codigo_aduana,omitempty"`
	NumeroDocumentoExportacion string    `json:"numero_documento_exportacion,omitempty" bson:"numero_documento_exportacion,omitempty"`
	FechaDocumentoExportacion  time.Time `json:"fecha_documento_exportacion,omitempty" bson:"fecha_documento_exportacion,omitempty"`
	DireccionOrigen            string    `json:"direccion_origen" bson:"direccion_origen"`
	ComunaOrigen               string    `json:"comuna_origen" bson:"comuna_origen"`
	CiudadOrigen               string    `json:"ciudad_origen" bson:"ciudad_origen"`
	DireccionDestino           string    `json:"direccion_destino" bson:"direccion_destino"`
	ComunaDestino              string    `json:"comuna_destino" bson:"comuna_destino"`
	CiudadDestino              string    `json:"ciudad_destino" bson:"ciudad_destino"`
	CAF                        *CAF      `json:"caf,omitempty" bson:"caf,omitempty"`
	TimbreElectronico          string    `json:"timbre_electronico,omitempty" bson:"timbre_electronico,omitempty"`
	FirmaElectronica           string    `json:"firma_electronica,omitempty" bson:"firma_electronica,omitempty"`
	Transportista              string    `json:"transportista" bson:"transportista"`
	TipoTraslado               string    `json:"tipo_traslado" bson:"tipo_traslado"`
	MontoNeto                  float64   `json:"monto_neto" bson:"monto_neto"`
	MontoExento                float64   `json:"monto_exento" bson:"monto_exento"`
	MontoIVA                   float64   `json:"monto_iva" bson:"monto_iva"`
	Items                      []Item    `json:"items" bson:"items"`
}

// GuiaDespachoRequest representa la solicitud para crear una guía de despacho
type GuiaDespachoRequest struct {
	TipoDTE                    TipoDTE   `json:"tipo_dte"`
	Folio                      int       `json:"folio"`
	FechaEmision               time.Time `json:"fecha_emision"`
	RutEmisor                  string    `json:"rut_emisor"`
	RazonSocialEmisor          string    `json:"razon_social_emisor"`
	RutReceptor                string    `json:"rut_receptor"`
	RazonSocialReceptor        string    `json:"razon_social_receptor"`
	IndicadorTraslado          string    `json:"indicador_traslado"`
	IndicadorServicio          string    `json:"indicador_servicio,omitempty"`
	IndicadorVentas            string    `json:"indicador_ventas,omitempty"`
	IndicadorTransporte        string    `json:"indicador_transporte,omitempty"`
	IndicadorExportacion       string    `json:"indicador_exportacion,omitempty"`
	RutTransportista           string    `json:"rut_transportista"`
	RazonSocialTransportista   string    `json:"razon_social_transportista"`
	Patente                    string    `json:"patente,omitempty"`
	FechaInicioTransporte      time.Time `json:"fecha_inicio_transporte,omitempty"`
	FechaFinTransporte         time.Time `json:"fecha_fin_transporte,omitempty"`
	FechaInicioServicio        time.Time `json:"fecha_inicio_servicio,omitempty"`
	FechaFinServicio           time.Time `json:"fecha_fin_servicio,omitempty"`
	Periodicidad               string    `json:"periodicidad,omitempty"`
	CodigoAduana               string    `json:"codigo_aduana,omitempty"`
	NumeroDocumentoExportacion string    `json:"numero_documento_exportacion,omitempty"`
	FechaDocumentoExportacion  time.Time `json:"fecha_documento_exportacion,omitempty"`
	DireccionOrigen            string    `json:"direccion_origen"`
	ComunaOrigen               string    `json:"comuna_origen"`
	CiudadOrigen               string    `json:"ciudad_origen"`
	DireccionDestino           string    `json:"direccion_destino"`
	ComunaDestino              string    `json:"comuna_destino"`
	CiudadDestino              string    `json:"ciudad_destino"`
	Items                      []Item    `json:"items"`
}

// GuiaDespachoResponse representa la respuesta de una guía de despacho
type GuiaDespachoResponse struct {
	TrackID string `json:"track_id"`
	Estado  string `json:"estado"`
	Glosa   string `json:"glosa"`
}
