package sii

// ClienteSII define la interfaz para interactuar con el SII
type ClienteSII interface {
	Autenticar() error
	EnviarDocumento(xmlData string) (string, error)
	ConsultarEstado(trackID string) (string, error)
}

// Config representa la configuración para el cliente SII
type Config struct {
	URL         string
	Certificado string
	Password    string
	Ambiente    string
}
