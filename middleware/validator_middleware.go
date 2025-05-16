package middleware

import (
	"context"
	"fmt"
	"sync"

	"github.com/fmgo/models"
	"github.com/fmgo/services"
)

// ValidatorMiddleware proporciona middleware para validación automática de documentos
type ValidatorMiddleware struct {
	validator *services.ValidatorService
	once      sync.Once
}

// NewValidatorMiddleware crea una nueva instancia de ValidatorMiddleware
func NewValidatorMiddleware(schemaPath string) *ValidatorMiddleware {
	return &ValidatorMiddleware{
		validator: services.NewValidatorService(schemaPath),
	}
}

// initValidator inicializa el validador de manera lazy
func (m *ValidatorMiddleware) initValidator() error {
	var initErr error
	m.once.Do(func() {
		// Cargar esquemas necesarios
		if err := m.validator.CargarEsquema("DTE_v10.xsd"); err != nil {
			initErr = fmt.Errorf("error al cargar esquema DTE: %w", err)
			return
		}
		if err := m.validator.CargarEsquema("EnvioDTE_v10.xsd"); err != nil {
			initErr = fmt.Errorf("error al cargar esquema EnvioDTE: %w", err)
			return
		}
		if err := m.validator.CargarEsquema("EnvioBOLETA_v11.xsd"); err != nil {
			initErr = fmt.Errorf("error al cargar esquema Boleta: %w", err)
			return
		}
	})
	return initErr
}

// WithDocumentoValidation envuelve un handler con validación de documento
func (m *ValidatorMiddleware) WithDocumentoValidation(next func(ctx context.Context, doc *models.DTEDocument) error) func(ctx context.Context, doc *models.DTEDocument) error {
	return func(ctx context.Context, doc *models.DTEDocument) error {
		if err := m.initValidator(); err != nil {
			return fmt.Errorf("error al inicializar validador: %w", err)
		}

		if err := m.validator.ValidarDocumento(doc); err != nil {
			return fmt.Errorf("error de validación de documento: %w", err)
		}

		return next(ctx, doc)
	}
}

// WithEnvioValidation envuelve un handler con validación de envío
func (m *ValidatorMiddleware) WithEnvioValidation(next func(ctx context.Context, envio *models.EnvioDTEDocument) error) func(ctx context.Context, envio *models.EnvioDTEDocument) error {
	return func(ctx context.Context, envio *models.EnvioDTEDocument) error {
		if err := m.initValidator(); err != nil {
			return fmt.Errorf("error al inicializar validador: %w", err)
		}

		if err := m.validator.ValidarEnvio(envio); err != nil {
			return fmt.Errorf("error de validación de envío: %w", err)
		}

		return next(ctx, envio)
	}
}
