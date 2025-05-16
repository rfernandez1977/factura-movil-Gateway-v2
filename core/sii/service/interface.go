package service

import (
	"context"

	"github.com/fmgo/core/sii/models"
)

// SIIService define la interfaz para el servicio de integración con el SII
type SIIService interface {
	// EnviarDTE envía un DTE al SII
	EnviarDTE(ctx context.Context, dte []byte) (*models.RespuestaSII, error)

	// ConsultarEstado consulta el estado de un DTE
	ConsultarEstado(ctx context.Context, trackID string) (*models.EstadoSII, error)

	// ConsultarDTE consulta un DTE específico
	ConsultarDTE(ctx context.Context, tipoDTE string, folio int64, rutEmisor string) (*models.EstadoSII, error)

	// VerificarComunicacion verifica la comunicación con el SII
	VerificarComunicacion(ctx context.Context) error
}
