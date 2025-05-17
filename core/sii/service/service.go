package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/fmgo/core/sii/models"
	"github.com/fmgo/utils/logger"
)

// SIIClient define la interfaz para el cliente del SII.
// Esta interfaz proporciona los métodos necesarios para interactuar con el SII.
type SIIClient interface {
	ObtenerSemilla(ctx context.Context) (string, error)
	ObtenerToken(ctx context.Context, semilla string) (string, error)
	EnviarDTE(ctx context.Context, sobre []byte, token string) (*models.RespuestaSII, error)
	ConsultarEstado(ctx context.Context, trackID string) (*models.EstadoConsulta, error)
	ConsultarDTE(ctx context.Context, tipoDTE models.TipoDocumentoSII, folio int64, rutEmisor string) (*models.EstadoConsulta, error)
	VerificarComunicacion(ctx context.Context) error
}

// SIIServiceImpl implementa la interfaz SIIService
type SIIServiceImpl struct {
	client    SIIClient
	token     string
	tokenExp  time.Time
	tokenLock sync.RWMutex
}

// NewSIIService crea una nueva instancia del servicio SII
func NewSIIService(client SIIClient) SIIService {
	return &SIIServiceImpl{
		client: client,
	}
}

// ObtenerSemilla obtiene una semilla del SII
func (s *SIIServiceImpl) ObtenerSemilla(ctx context.Context) (string, error) {
	logger.Info("Obteniendo semilla del SII",
		logger.RequestID(ctx.Value("request_id").(string)),
	)

	semilla, err := s.client.ObtenerSemilla(ctx)
	if err != nil {
		logger.Error("Error obteniendo semilla del SII", err,
			logger.RequestID(ctx.Value("request_id").(string)),
		)
		return "", fmt.Errorf("error obteniendo semilla: %w", err)
	}

	logger.Info("Semilla obtenida exitosamente",
		logger.Field("semilla", semilla),
		logger.RequestID(ctx.Value("request_id").(string)),
	)

	return semilla, nil
}

// ObtenerToken obtiene un token de autenticación usando una semilla
func (s *SIIServiceImpl) ObtenerToken(ctx context.Context, semilla string) (string, error) {
	logger.Info("Obteniendo token del SII",
		logger.Field("semilla", semilla),
		logger.RequestID(ctx.Value("request_id").(string)),
	)

	token, err := s.client.ObtenerToken(ctx, semilla)
	if err != nil {
		logger.Error("Error obteniendo token del SII", err,
			logger.Field("semilla", semilla),
			logger.RequestID(ctx.Value("request_id").(string)),
		)
		return "", fmt.Errorf("error obteniendo token: %w", err)
	}

	logger.Info("Token obtenido exitosamente",
		logger.RequestID(ctx.Value("request_id").(string)),
	)

	return token, nil
}

// obtenerToken obtiene un token válido para comunicarse con el SII
func (s *SIIServiceImpl) obtenerTokenValido(ctx context.Context) (string, error) {
	s.tokenLock.RLock()
	if s.token != "" && time.Now().Before(s.tokenExp) {
		token := s.token
		s.tokenLock.RUnlock()
		return token, nil
	}
	s.tokenLock.RUnlock()

	s.tokenLock.Lock()
	defer s.tokenLock.Unlock()

	// Verificar nuevamente en caso de que otro goroutine haya actualizado el token
	if s.token != "" && time.Now().Before(s.tokenExp) {
		return s.token, nil
	}

	// Obtener nueva semilla
	semilla, err := s.ObtenerSemilla(ctx)
	if err != nil {
		return "", err
	}

	// Obtener nuevo token
	token, err := s.ObtenerToken(ctx, semilla)
	if err != nil {
		return "", err
	}

	// Guardar token con expiración de 1 hora
	s.token = token
	s.tokenExp = time.Now().Add(1 * time.Hour)

	return token, nil
}

// EnviarDTE envía un DTE al SII
func (s *SIIServiceImpl) EnviarDTE(ctx context.Context, dte []byte) (*models.RespuestaSII, error) {
	token, err := s.obtenerTokenValido(ctx)
	if err != nil {
		return nil, err
	}

	logger.Info("Enviando DTE al SII",
		logger.Field("size", len(dte)),
		logger.RequestID(ctx.Value("request_id").(string)),
	)

	respuesta, err := s.client.EnviarDTE(ctx, dte, token)
	if err != nil {
		logger.Error("Error enviando DTE al SII", err,
			logger.Field("size", len(dte)),
			logger.RequestID(ctx.Value("request_id").(string)),
		)
		return nil, fmt.Errorf("error enviando DTE: %w", err)
	}

	logger.Info("DTE enviado exitosamente",
		logger.Field("track_id", respuesta.TrackID),
		logger.Field("estado", respuesta.Estado),
		logger.RequestID(ctx.Value("request_id").(string)),
	)

	return respuesta, nil
}

// ConsultarEstado consulta el estado de un DTE
func (s *SIIServiceImpl) ConsultarEstado(ctx context.Context, trackID string) (*models.EstadoConsulta, error) {
	logger.Info("Consultando estado de DTE",
		logger.Field("track_id", trackID),
		logger.RequestID(ctx.Value("request_id").(string)),
	)

	estado, err := s.client.ConsultarEstado(ctx, trackID)
	if err != nil {
		logger.Error("Error consultando estado de DTE", err,
			logger.Field("track_id", trackID),
			logger.RequestID(ctx.Value("request_id").(string)),
		)
		return nil, fmt.Errorf("error consultando estado: %w", err)
	}

	logger.Info("Estado de DTE consultado exitosamente",
		logger.Field("track_id", trackID),
		logger.Field("estado", estado.Estado),
		logger.RequestID(ctx.Value("request_id").(string)),
	)

	return estado, nil
}

// ConsultarDTE consulta un DTE específico
func (s *SIIServiceImpl) ConsultarDTE(ctx context.Context, tipoDTE models.TipoDocumentoSII, folio int64, rutEmisor string) (*models.EstadoConsulta, error) {
	logger.Info("Consultando DTE específico",
		logger.Field("tipo_dte", tipoDTE),
		logger.Field("folio", folio),
		logger.Field("rut_emisor", rutEmisor),
		logger.RequestID(ctx.Value("request_id").(string)),
	)

	estado, err := s.client.ConsultarDTE(ctx, tipoDTE, folio, rutEmisor)
	if err != nil {
		logger.Error("Error consultando DTE específico", err,
			logger.Field("tipo_dte", tipoDTE),
			logger.Field("folio", folio),
			logger.Field("rut_emisor", rutEmisor),
			logger.RequestID(ctx.Value("request_id").(string)),
		)
		return nil, fmt.Errorf("error consultando DTE: %w", err)
	}

	logger.Info("DTE consultado exitosamente",
		logger.Field("tipo_dte", tipoDTE),
		logger.Field("folio", folio),
		logger.Field("rut_emisor", rutEmisor),
		logger.Field("estado", estado.Estado),
		logger.RequestID(ctx.Value("request_id").(string)),
	)

	return estado, nil
}

// ValidarDTE valida un DTE antes de enviarlo al SII
func (s *SIIServiceImpl) ValidarDTE(ctx context.Context, dte []byte) (*models.ValidacionSII, error) {
	logger.Info("Validando DTE",
		logger.Field("size", len(dte)),
		logger.RequestID(ctx.Value("request_id").(string)),
	)

	// TODO: Implementar validación real del DTE
	validacion := &models.ValidacionSII{
		CodigoValidacion: "VAL001",
		Resultado:        true,
		FechaValidacion:  time.Now(),
	}

	logger.Info("DTE validado exitosamente",
		logger.Field("codigo", validacion.CodigoValidacion),
		logger.Field("resultado", validacion.Resultado),
		logger.RequestID(ctx.Value("request_id").(string)),
	)

	return validacion, nil
}

// VerificarComunicacion verifica la comunicación con el SII
func (s *SIIServiceImpl) VerificarComunicacion(ctx context.Context) error {
	logger.Info("Verificando comunicación con SII",
		logger.RequestID(ctx.Value("request_id").(string)),
	)

	if err := s.client.VerificarComunicacion(ctx); err != nil {
		logger.Error("Error verificando comunicación con SII", err,
			logger.RequestID(ctx.Value("request_id").(string)),
		)
		return fmt.Errorf("error verificando comunicación: %w", err)
	}

	logger.Info("Comunicación con SII verificada exitosamente",
		logger.RequestID(ctx.Value("request_id").(string)),
	)

	return nil
}
