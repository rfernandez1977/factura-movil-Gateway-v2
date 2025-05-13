package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/cursor/FMgo/models"
	"github.com/cursor/FMgo/supabase"
	"go.uber.org/zap"
)

// ConfiguracionService maneja la lógica de negocio relacionada con la configuración
type ConfiguracionService struct {
	logger         *zap.Logger
	supabaseClient *supabase.Client
	cacheService   *CacheService
	empresaService *EmpresaService
}

// NewConfiguracionService crea una nueva instancia del servicio de configuración
func NewConfiguracionService(
	logger *zap.Logger,
	supabaseClient *supabase.Client,
	cacheService *CacheService,
	empresaService *EmpresaService,
) *ConfiguracionService {
	return &ConfiguracionService{
		logger:         logger,
		supabaseClient: supabaseClient,
		cacheService:   cacheService,
		empresaService: empresaService,
	}
}

// ObtenerConfiguracion obtiene la configuración de una empresa
func (s *ConfiguracionService) ObtenerConfiguracion(empresaID string) (*models.Configuracion, error) {
	// Intentar obtener de caché primero
	if s.cacheService != nil {
		cacheKey := fmt.Sprintf("config:%s", empresaID)
		cachedConfig, found := s.cacheService.Get(cacheKey)
		if found {
			s.logger.Info("Configuración obtenida de caché", zap.String("empresaID", empresaID))
			return cachedConfig.(*models.Configuracion), nil
		}
	}

	// Consultar en Supabase
	var config models.Configuracion
	err := s.supabaseClient.GetClient().DB.From("configuraciones").
		Select("*").
		Eq("empresa_id", empresaID).
		Single().
		Execute(&config)

	if err != nil {
		s.logger.Error("Error al obtener configuración", zap.String("empresaID", empresaID), zap.Error(err))
		return nil, err
	}

	// Guardar en caché
	if s.cacheService != nil {
		cacheKey := fmt.Sprintf("config:%s", empresaID)
		s.cacheService.Set(cacheKey, &config, 30*time.Minute)
	}

	return &config, nil
}

// ActualizarConfiguracion actualiza la configuración de una empresa
func (s *ConfiguracionService) ActualizarConfiguracion(config *models.Configuracion) error {
	if config == nil {
		return errors.New("configuración no puede ser nil")
	}

	// Verificar que la empresa existe
	empresa, err := s.empresaService.ObtenerEmpresa(config.EmpresaID)
	if err != nil {
		s.logger.Error("Error al obtener empresa para actualizar configuración",
			zap.String("empresaID", config.EmpresaID), zap.Error(err))
		return err
	}

	// Asegurarse de que los datos sean consistentes
	config.RUT = empresa.RUT
	config.UpdatedAt = time.Now()

	// Actualizar en Supabase
	_, err = s.supabaseClient.GetClient().DB.From("configuraciones").
		Update(config).
		Eq("id", config.ID).
		Execute()

	if err != nil {
		s.logger.Error("Error al actualizar configuración",
			zap.String("id", config.ID), zap.Error(err))
		return err
	}

	// Invalidar caché
	if s.cacheService != nil {
		cacheKey := fmt.Sprintf("config:%s", config.EmpresaID)
		s.cacheService.Delete(cacheKey)
	}

	return nil
}

// ObtenerConfiguracionSII obtiene la configuración del SII para una empresa
func (s *ConfiguracionService) ObtenerConfiguracionSII(empresaID string) (*models.ConfiguracionSIIEmpresa, error) {
	// Obtener la configuración completa
	config, err := s.ObtenerConfiguracion(empresaID)
	if err != nil {
		return nil, err
	}

	return &config.ConfigSII, nil
}

// ActualizarConfiguracionSII actualiza la configuración del SII para una empresa
func (s *ConfiguracionService) ActualizarConfiguracionSII(configSII *models.ConfiguracionSIIEmpresa) error {
	if configSII == nil {
		return errors.New("configuración SII no puede ser nil")
	}

	// Obtener la configuración completa
	config, err := s.ObtenerConfiguracion(configSII.EmpresaID)
	if err != nil {
		return err
	}

	// Actualizar la sección de configuración SII
	config.ConfigSII = *configSII
	config.UpdatedAt = time.Now()

	// Actualizar en la base de datos
	return s.ActualizarConfiguracion(config)
}

// ObtenerConfiguracionEmail obtiene la configuración de email para una empresa
func (s *ConfiguracionService) ObtenerConfiguracionEmail(empresaID string) (*models.ConfiguracionEmail, error) {
	// Obtener la configuración completa
	config, err := s.ObtenerConfiguracion(empresaID)
	if err != nil {
		return nil, err
	}

	return &config.ConfigEmail, nil
}

// ActualizarConfiguracionEmail actualiza la configuración de email para una empresa
func (s *ConfiguracionService) ActualizarConfiguracionEmail(configEmail *models.ConfiguracionEmail) error {
	if configEmail == nil {
		return errors.New("configuración Email no puede ser nil")
	}

	// Obtener la configuración completa
	config, err := s.ObtenerConfiguracion(configEmail.EmpresaID)
	if err != nil {
		return err
	}

	// Actualizar la sección de configuración Email
	config.ConfigEmail = *configEmail
	config.UpdatedAt = time.Now()

	// Actualizar en la base de datos
	return s.ActualizarConfiguracion(config)
}
