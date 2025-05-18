package services

import (
	"context"
	"fmt"
	"time"

	"FMgo/core/firma/models"
)

// CAFService implementa el servicio de gestión de CAF
type CAFService struct {
	cafRepo CAFRepository
	logger  Logger
	alerter AlertService
}

// NewCAFService crea una nueva instancia del servicio CAF
func NewCAFService(repo CAFRepository, logger Logger, alerter AlertService) *CAFService {
	return &CAFService{
		cafRepo: repo,
		logger:  logger,
		alerter: alerter,
	}
}

// ValidarCAF valida un CAF
func (s *CAFService) ValidarCAF(ctx context.Context, caf *models.CAF) error {
	// Validar fecha de vencimiento
	if time.Now().After(caf.FechaVencimiento) {
		return fmt.Errorf("CAF vencido: %s", caf.ID)
	}

	// Validar rango de folios
	if caf.FolioInicial > caf.FolioFinal {
		return fmt.Errorf("rango de folios inválido: %d-%d", caf.FolioInicial, caf.FolioFinal)
	}

	// Validar firma del SII
	if err := s.validarFirmaSII(caf); err != nil {
		return fmt.Errorf("firma SII inválida: %w", err)
	}

	return nil
}

// ObtenerCAF obtiene un CAF disponible para un tipo de documento
func (s *CAFService) ObtenerCAF(ctx context.Context, tipo string, folio int64) (*models.CAF, error) {
	caf, err := s.cafRepo.ObtenerCAFPorFolio(ctx, tipo, folio)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo CAF: %w", err)
	}

	// Validar CAF
	if err := s.ValidarCAF(ctx, caf); err != nil {
		return nil, err
	}

	// Verificar stock y alertar si es necesario
	if s.requiereAlerta(caf) {
		s.alerter.EnviarAlerta(ctx, &models.Alerta{
			Tipo:     "CAF_BAJO_STOCK",
			Mensaje:  fmt.Sprintf("Stock bajo de CAF para tipo %s", tipo),
			Detalles: map[string]interface{}{"tipo": tipo, "disponibles": caf.FoliosDisponibles()},
		})
	}

	return caf, nil
}

// RegistrarCAF registra un nuevo CAF
func (s *CAFService) RegistrarCAF(ctx context.Context, caf *models.CAF) error {
	// Validar CAF antes de registrar
	if err := s.ValidarCAF(ctx, caf); err != nil {
		return err
	}

	// Registrar en repositorio
	if err := s.cafRepo.GuardarCAF(ctx, caf); err != nil {
		return fmt.Errorf("error guardando CAF: %w", err)
	}

	s.logger.Info("CAF registrado exitosamente",
		"id", caf.ID,
		"tipo", caf.TipoDocumento,
		"rango", fmt.Sprintf("%d-%d", caf.FolioInicial, caf.FolioFinal))

	return nil
}

// ConsultarDisponibilidad consulta la disponibilidad de CAF
func (s *CAFService) ConsultarDisponibilidad(ctx context.Context, tipo string) (*models.DisponibilidadCAF, error) {
	cafs, err := s.cafRepo.ListarCAFsPorTipo(ctx, tipo)
	if err != nil {
		return nil, fmt.Errorf("error consultando CAFs: %w", err)
	}

	disp := &models.DisponibilidadCAF{
		TipoDocumento: tipo,
		FechaConsulta: time.Now(),
	}

	for _, caf := range cafs {
		if time.Now().Before(caf.FechaVencimiento) {
			disp.FoliosDisponibles += caf.FoliosDisponibles()
			disp.CAFsActivos++
		}
	}

	return disp, nil
}

// validarFirmaSII valida la firma del SII en el CAF
func (s *CAFService) validarFirmaSII(caf *models.CAF) error {
	// Implementar validación de firma SII
	// Esta es una implementación simulada
	return nil
}

// requiereAlerta verifica si se debe enviar una alerta por bajo stock
func (s *CAFService) requiereAlerta(caf *models.CAF) bool {
	disponibles := caf.FoliosDisponibles()
	return disponibles < 1000 // Umbral de alerta
}
