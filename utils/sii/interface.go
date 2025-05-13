package sii

import "github.com/cursor/FMgo/models"

// SIIService define la interfaz para el servicio de integración con el SII
type SIIService interface {
	// EnviarDTE envía un DTE al SII
	EnviarDTE(dte *models.DTEXMLModel) (*models.RespuestaSII, error)

	// ConsultarEstado consulta el estado de un DTE en el SII
	ConsultarEstado(trackID string) (*models.EstadoSII, error)

	// ValidarDTE valida un DTE antes de enviarlo al SII
	ValidarDTE(dte *models.DTEXMLModel) (*models.RespuestaSII, error)

	// ConsultarDTE consulta un DTE específico en el SII
	ConsultarDTE(tipoDTE, folio, rutEmisor string) (*models.EstadoSII, error)

	// VerificarComunicacion verifica la comunicación con el SII
	VerificarComunicacion() error
}
