package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/fmgo/core/sii/client"
	"github.com/fmgo/core/sii/models"
)

// DefaultSIIService implementa la interfaz SIIService
type DefaultSIIService struct {
	client   client.SIIClient
	token    string
	tokenMu  sync.RWMutex
	tokenExp time.Time
}

// NewDefaultSIIService crea una nueva instancia del servicio SII
func NewDefaultSIIService(client client.SIIClient) *DefaultSIIService {
	return &DefaultSIIService{
		client: client,
	}
}

// actualizarToken actualiza el token si es necesario
func (s *DefaultSIIService) actualizarToken(ctx context.Context) error {
	s.tokenMu.RLock()
	tokenValido := s.token != "" && time.Now().Before(s.tokenExp)
	s.tokenMu.RUnlock()

	if tokenValido {
		return nil
	}

	s.tokenMu.Lock()
	defer s.tokenMu.Unlock()

	// Verificar nuevamente en caso de que otro goroutine haya actualizado el token
	if s.token != "" && time.Now().Before(s.tokenExp) {
		return nil
	}

	// Obtener semilla
	semilla, err := s.client.ObtenerSemilla(ctx)
	if err != nil {
		return fmt.Errorf("error al obtener semilla: %v", err)
	}

	// Obtener token
	token, err := s.client.ObtenerToken(ctx, semilla)
	if err != nil {
		return fmt.Errorf("error al obtener token: %v", err)
	}

	s.token = token
	s.tokenExp = time.Now().Add(1 * time.Hour) // Los tokens del SII expiran después de 1 hora
	return nil
}

// EnviarDTE implementa SIIService.EnviarDTE
func (s *DefaultSIIService) EnviarDTE(ctx context.Context, dte []byte) (*models.RespuestaSII, error) {
	if err := s.actualizarToken(ctx); err != nil {
		return nil, fmt.Errorf("error al actualizar token: %v", err)
	}

	s.tokenMu.RLock()
	token := s.token
	s.tokenMu.RUnlock()

	return s.client.EnviarDTE(ctx, dte, token)
}

// ConsultarEstado implementa SIIService.ConsultarEstado
func (s *DefaultSIIService) ConsultarEstado(ctx context.Context, trackID string) (*models.EstadoSII, error) {
	if err := s.actualizarToken(ctx); err != nil {
		return nil, fmt.Errorf("error al actualizar token: %v", err)
	}

	return s.client.ConsultarEstado(ctx, trackID)
}

// ConsultarDTE implementa SIIService.ConsultarDTE
func (s *DefaultSIIService) ConsultarDTE(ctx context.Context, tipoDTE string, folio int64, rutEmisor string) (*models.EstadoSII, error) {
	if err := s.actualizarToken(ctx); err != nil {
		return nil, fmt.Errorf("error al actualizar token: %v", err)
	}

	return s.client.ConsultarDTE(ctx, tipoDTE, folio, rutEmisor)
}

// VerificarComunicacion implementa SIIService.VerificarComunicacion
func (s *DefaultSIIService) VerificarComunicacion(ctx context.Context) error {
	return s.client.VerificarComunicacion(ctx)
}
