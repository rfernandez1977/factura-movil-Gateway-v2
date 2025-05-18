package services

import (
	"FMgo/core/sii/logger"
	"context"
	"fmt"
	"time"

	"FMgo/core/firma/services"
	"FMgo/core/sii/client"
	"FMgo/core/sii/models"
)

// IntegrationService maneja la integración entre DTE, Firma y SII
type IntegrationService struct {
	siiClient    client.SIIClient
	firmaService *services.FirmaService
	logger       *logger.Logger
}

// NewIntegrationService crea una nueva instancia del servicio de integración
func NewIntegrationService(siiClient client.SIIClient, firmaService *services.FirmaService, logger *logger.Logger) *IntegrationService {
	return &IntegrationService{
		siiClient:    siiClient,
		firmaService: firmaService,
		logger:       logger,
	}
}

// EnviarDocumento procesa y envía un documento al SII
func (s *IntegrationService) EnviarDocumento(ctx context.Context, doc *models.Documento) (*models.RespuestaEnvio, error) {
	// 1. Firmar documento
	docFirmado, err := s.firmaService.FirmarDocumento(ctx, doc)
	if err != nil {
		return nil, fmt.Errorf("error al firmar documento: %w", err)
	}

	// 2. Obtener token de autenticación
	semilla, err := s.siiClient.ObtenerSemilla(ctx)
	if err != nil {
		return nil, fmt.Errorf("error al obtener semilla: %w", err)
	}

	token, err := s.siiClient.ObtenerToken(ctx, semilla)
	if err != nil {
		return nil, fmt.Errorf("error al obtener token: %w", err)
	}

	// 3. Enviar al SII
	respuesta, err := s.siiClient.EnviarDTE(ctx, docFirmado.XML, token)
	if err != nil {
		return nil, fmt.Errorf("error al enviar DTE: %w", err)
	}

	// 4. Procesar respuesta
	resultado := &models.RespuestaEnvio{
		TrackID:   respuesta.TrackID,
		Estado:    respuesta.Estado,
		Glosa:     respuesta.Glosa,
		Timestamp: time.Now(),
	}

	// 5. Registrar resultado
	s.logger.Info("Documento enviado exitosamente",
		"trackID", resultado.TrackID,
		"estado", resultado.Estado,
		"folio", doc.Folio)

	return resultado, nil
}

// ConsultarEstadoEnvio consulta el estado de un envío
func (s *IntegrationService) ConsultarEstadoEnvio(ctx context.Context, trackID string) (*models.EstadoEnvio, error) {
	estado, err := s.siiClient.ConsultarEstado(ctx, trackID)
	if err != nil {
		return nil, fmt.Errorf("error al consultar estado: %w", err)
	}

	resultado := &models.EstadoEnvio{
		TrackID:   trackID,
		Estado:    estado.Estado,
		Glosa:     estado.Glosa,
		Timestamp: time.Now(),
	}

	s.logger.Info("Estado consultado",
		"trackID", trackID,
		"estado", estado.Estado)

	return resultado, nil
}

// ValidarDocumento valida un documento contra el SII
func (s *IntegrationService) ValidarDocumento(ctx context.Context, doc *models.Documento) (*models.ResultadoValidacion, error) {
	// 1. Validar estructura
	if err := doc.Validar(); err != nil {
		return nil, fmt.Errorf("error de validación: %w", err)
	}

	// 2. Consultar estado en SII
	estado, err := s.siiClient.ConsultarDTE(ctx, doc.TipoDTE, doc.Folio, doc.RutEmisor)
	if err != nil {
		return nil, fmt.Errorf("error al consultar DTE: %w", err)
	}

	resultado := &models.ResultadoValidacion{
		Folio:     doc.Folio,
		TipoDTE:   doc.TipoDTE,
		Estado:    estado.Estado,
		Glosa:     estado.Glosa,
		Timestamp: time.Now(),
	}

	s.logger.Info("Documento validado",
		"folio", doc.Folio,
		"tipo", doc.TipoDTE,
		"estado", estado.Estado)

	return resultado, nil
}
