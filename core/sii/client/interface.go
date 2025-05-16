package client

import (
	"context"

	"github.com/fmgo/core/sii/models"
)

// SIIClient define la interfaz para el cliente del SII
type SIIClient interface {
	// ObtenerSemilla obtiene una semilla del SII
	ObtenerSemilla(ctx context.Context) (string, error)

	// ObtenerToken obtiene un token de autenticación usando una semilla
	ObtenerToken(ctx context.Context, semilla string) (string, error)

	// EnviarDTE envía un DTE al SII
	EnviarDTE(ctx context.Context, sobre []byte, token string) (*models.RespuestaSII, error)

	// ConsultarEstado consulta el estado de un DTE
	ConsultarEstado(ctx context.Context, trackID string) (*models.EstadoSII, error)

	// ConsultarDTE consulta un DTE específico
	ConsultarDTE(ctx context.Context, tipoDTE string, folio int64, rutEmisor string) (*models.EstadoSII, error)

	// VerificarComunicacion verifica la comunicación con el SII
	VerificarComunicacion(ctx context.Context) error
}
