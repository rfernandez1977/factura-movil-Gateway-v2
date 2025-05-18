package dte

// DTE representa un Documento Tributario Electrónico
type DTE struct {
	TipoDTE      string
	Folio        int64
	RUTEmisor    string
	RUTReceptor  string
	MontoTotal   int64
	MontoNeto    int64
	MontoIVA     int64
	FechaEmision string
}

// ValidadorDTE define la interfaz para validación de documentos
type ValidadorDTE interface {
	ValidarDTE(dte *DTE) error
	ValidarCAF(dte *DTE) error
}
