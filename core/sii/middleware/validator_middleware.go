package middleware

import (
	"context"
	"fmt"
	"sync"

	"github.com/fmgo/core/sii/models"
	"github.com/fmgo/core/sii/services"
)

// ValidatorMiddleware proporciona middleware para validación automática de DTEs
type ValidatorMiddleware struct {
	validator  *services.ValidatorService
	initOnce   sync.Once
	initError  error
	schemasDir string
}

// NewValidatorMiddleware crea una nueva instancia del middleware de validación
func NewValidatorMiddleware(schemasDir string) *ValidatorMiddleware {
	return &ValidatorMiddleware{
		validator:  services.NewValidatorService(),
		schemasDir: schemasDir,
	}
}

// ValidateDTE valida un DTE antes de su procesamiento
func (m *ValidatorMiddleware) ValidateDTE(ctx context.Context, dte *models.DTE) error {
	// Inicializar validador de forma lazy
	m.initOnce.Do(func() {
		m.initError = m.validator.CargarEsquemasBase(m.schemasDir)
	})

	if m.initError != nil {
		return fmt.Errorf("error inicializando validador: %w", m.initError)
	}

	// Validar el DTE
	if err := m.validator.ValidarDTE(dte); err != nil {
		return fmt.Errorf("error de validación DTE: %w", err)
	}

	return nil
}

// ValidateEnvioDTE valida un envío de DTE antes de su procesamiento
func (m *ValidatorMiddleware) ValidateEnvioDTE(ctx context.Context, envioDTE []byte) error {
	// Inicializar validador de forma lazy
	m.initOnce.Do(func() {
		m.initError = m.validator.CargarEsquemasBase(m.schemasDir)
	})

	if m.initError != nil {
		return fmt.Errorf("error inicializando validador: %w", m.initError)
	}

	// Validar el envío
	if err := m.validator.ValidarEnvioDTE(envioDTE); err != nil {
		return fmt.Errorf("error de validación EnvioDTE: %w", err)
	}

	return nil
}

// WithDTEValidation envuelve un handler con validación de DTE
func (m *ValidatorMiddleware) WithDTEValidation(next func(context.Context, *models.DTE) error) func(context.Context, *models.DTE) error {
	return func(ctx context.Context, dte *models.DTE) error {
		if err := m.ValidateDTE(ctx, dte); err != nil {
			return err
		}
		return next(ctx, dte)
	}
}

// WithEnvioDTEValidation envuelve un handler con validación de EnvioDTE
func (m *ValidatorMiddleware) WithEnvioDTEValidation(next func(context.Context, []byte) error) func(context.Context, []byte) error {
	return func(ctx context.Context, envioDTE []byte) error {
		if err := m.ValidateEnvioDTE(ctx, envioDTE); err != nil {
			return err
		}
		return next(ctx, envioDTE)
	}
}
