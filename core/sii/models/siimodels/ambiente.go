package siimodels

// Ambiente representa el ambiente de ejecución (Producción o Certificación)
type Ambiente string

const (
	Produccion    Ambiente = "PRODUCCION"
	Certificacion Ambiente = "CERTIFICACION"
)

// String implementa la interfaz Stringer
func (a Ambiente) String() string {
	return string(a)
}

// IsValid verifica si el ambiente es válido
func (a Ambiente) IsValid() bool {
	switch a {
	case Produccion, Certificacion:
		return true
	default:
		return false
	}
}
