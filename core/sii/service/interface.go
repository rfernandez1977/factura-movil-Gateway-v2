package service

import (
	"context"

	"github.com/fmgo/core/sii/models"
)

// SIIService define la interfaz para el servicio de integración con el SII
type SIIService interface {
	// ObtenerSemilla obtiene una semilla del SII
	ObtenerSemilla(ctx context.Context) (string, error)

	// ObtenerToken obtiene un token de autenticación usando una semilla
	ObtenerToken(ctx context.Context, semilla string) (string, error)

	// EnviarDTE envía un DTE al SII
	EnviarDTE(ctx context.Context, dte []byte) (*models.RespuestaSII, error)

	// ConsultarEstado consulta el estado de un DTE
	ConsultarEstado(ctx context.Context, trackID string) (*models.EstadoConsulta, error)

	// ConsultarDTE consulta un DTE específico
	ConsultarDTE(ctx context.Context, tipoDTE models.TipoDocumentoSII, folio int64, rutEmisor string) (*models.EstadoConsulta, error)

	// ValidarDTE valida un DTE antes de enviarlo al SII
	ValidarDTE(ctx context.Context, dte []byte) (*models.ValidacionSII, error)

	// VerificarComunicacion verifica la comunicación con el SII
	VerificarComunicacion(ctx context.Context) error
}
