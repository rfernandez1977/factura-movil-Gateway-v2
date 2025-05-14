package models

// EstadoDTE representa el estado de un documento tributario electrónico
type EstadoDTE string

// Estados de los documentos DTE
const (
	EstadoDTEEmitido   EstadoDTE = "EMITIDO"   // Documento emitido pero no enviado al SII
	EstadoDTEEnviado   EstadoDTE = "ENVIADO"   // Documento enviado al SII
	EstadoDTEAceptado  EstadoDTE = "ACEPTADO"  // Documento aceptado por el SII
	EstadoDTERechazado EstadoDTE = "RECHAZADO" // Documento rechazado por el SII
	EstadoDTEPendiente EstadoDTE = "PENDIENTE" // Documento pendiente de resolución por el SII
	EstadoDTEAnulado   EstadoDTE = "ANULADO"   // Documento anulado
	EstadoDTEErroneo   EstadoDTE = "ERRONEO"   // Documento con errores
	EstadoDTEBorrador  EstadoDTE = "BORRADOR"  // Documento en estado borrador
)

// Timestamps representa las fechas importantes de un documento
type Timestamps struct {
	FechaEmision    string `json:"fecha_emision,omitempty" bson:"fecha_emision,omitempty"`
	FechaEnvio      string `json:"fecha_envio,omitempty" bson:"fecha_envio,omitempty"`
	FechaRecepcion  string `json:"fecha_recepcion,omitempty" bson:"fecha_recepcion,omitempty"`
	FechaAceptacion string `json:"fecha_aceptacion,omitempty" bson:"fecha_aceptacion,omitempty"`
	FechaRechazo    string `json:"fecha_rechazo,omitempty" bson:"fecha_rechazo,omitempty"`
	Creado          string `json:"creado,omitempty" bson:"creado,omitempty"`
	Modificado      string `json:"modificado,omitempty" bson:"modificado,omitempty"`
}
