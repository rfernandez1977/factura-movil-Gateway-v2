package models

// EstadoDTE representa el estado de un DTE en el SII
type EstadoDTE struct {
	Estado       string `xml:"ESTADO"`
	Glosa        string `xml:"GLOSA"`
	NumTrack     string `xml:"NUMTRACK,omitempty"`
	RutEmisor    string `xml:"RUTEMISOR,omitempty"`
	RutReceptor  string `xml:"RUTRECEPTOR,omitempty"`
	TipoDTE      string `xml:"TIPODTE,omitempty"`
	FolioDTE     string `xml:"FOLIODTE,omitempty"`
	FechaEmision string `xml:"FECHAEMISION,omitempty"`
	MontoTotal   string `xml:"MONTOTOTAL,omitempty"`
	TrackID      string `xml:"TRACKID,omitempty"`
}

// EstadoEnvioSII representa el estado de un envío al SII
type EstadoEnvioSII struct {
	TrackID     string `xml:"trackId"`
	Estado      string `xml:"estado"`
	Glosa       string `xml:"glosa"`
	NumAtencion string `xml:"numAtencion,omitempty"`
	FechaRecep  string `xml:"fechaRecep,omitempty"`
}

// EstadoSII representa el estado de un documento en el SII
type EstadoSII string

const (
	// Estados generales
	EstadoOK           EstadoSII = "00" // OK
	EstadoErrorGeneral EstadoSII = "99" // Error general
	EstadoNoAutorizado EstadoSII = "01" // No autorizado

	// Estados de documentos
	EstadoRecibido    EstadoSII = "REC"  // Recibido
	EstadoProcesando  EstadoSII = "EPR"  // En proceso
	EstadoRechazado   EstadoSII = "RCH"  // Rechazado
	EstadoAceptado    EstadoSII = "SOK"  // Aceptado
	EstadoReprocesado EstadoSII = "RPR"  // Reprocesado
	EstadoError       EstadoSII = "ERR"  // Error
	EstadoNoRecibido  EstadoSII = "NRE"  // No recibido
	EstadoDesconocido EstadoSII = "DESC" // Desconocido
)

// String devuelve la representación en string del estado
func (e EstadoSII) String() string {
	switch e {
	case EstadoOK:
		return "OK"
	case EstadoErrorGeneral:
		return "Error general"
	case EstadoNoAutorizado:
		return "No autorizado"
	case EstadoRecibido:
		return "Recibido"
	case EstadoProcesando:
		return "En proceso"
	case EstadoRechazado:
		return "Rechazado"
	case EstadoAceptado:
		return "Aceptado"
	case EstadoReprocesado:
		return "Reprocesado"
	case EstadoError:
		return "Error"
	case EstadoNoRecibido:
		return "No recibido"
	default:
		return "Desconocido"
	}
}

// IsTerminal indica si el estado es terminal (no cambiará)
func (e EstadoSII) IsTerminal() bool {
	switch e {
	case EstadoAceptado, EstadoRechazado, EstadoError, EstadoErrorGeneral:
		return true
	default:
		return false
	}
}

// IsSuccess indica si el estado representa un éxito
func (e EstadoSII) IsSuccess() bool {
	switch e {
	case EstadoOK, EstadoAceptado, EstadoReprocesado:
		return true
	default:
		return false
	}
}

// ParseEstadoSII convierte un string a EstadoSII
func ParseEstadoSII(estado string) EstadoSII {
	switch estado {
	case string(EstadoOK):
		return EstadoOK
	case string(EstadoErrorGeneral):
		return EstadoErrorGeneral
	case string(EstadoNoAutorizado):
		return EstadoNoAutorizado
	case string(EstadoRecibido):
		return EstadoRecibido
	case string(EstadoProcesando):
		return EstadoProcesando
	case string(EstadoRechazado):
		return EstadoRechazado
	case string(EstadoAceptado):
		return EstadoAceptado
	case string(EstadoReprocesado):
		return EstadoReprocesado
	case string(EstadoError):
		return EstadoError
	case string(EstadoNoRecibido):
		return EstadoNoRecibido
	default:
		return EstadoDesconocido
	}
}
