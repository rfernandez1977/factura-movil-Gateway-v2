package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/cursor/FMgo/models"
	"github.com/cursor/FMgo/supabase"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

// ConfiguracionService maneja la lógica de negocio relacionada con la configuración
type ConfiguracionService struct {
	logger         *zap.Logger
	supabaseClient *supabase.Client
	cacheService   *CacheService
	empresaService *EmpresaService
	db             *mongo.Database
}

// NewConfiguracionService crea una nueva instancia del servicio de configuración
func NewConfiguracionService(
	logger *zap.Logger,
	supabaseClient *supabase.Client,
	cacheService *CacheService,
	empresaService *EmpresaService,
	db *mongo.Database,
) *ConfiguracionService {
	return &ConfiguracionService{
		logger:         logger,
		supabaseClient: supabaseClient,
		cacheService:   cacheService,
		empresaService: empresaService,
		db:             db,
	}
}

// ObtenerConfiguracion obtiene la configuración del sistema
func (s *ConfiguracionService) ObtenerConfiguracion(clave string) (*models.Configuracion, error) {
	cacheKey := fmt.Sprintf("config:%s", clave)

	// Verificar si está en caché
	if s.cacheService != nil {
		ctx := context.Background()
		var cachedConfig models.Configuracion
		err := s.cacheService.Get(ctx, cacheKey, &cachedConfig)
		if err == nil {
			return &cachedConfig, nil
		}
	}

	// No está en caché, obtener de la base de datos
	var config models.Configuracion
	var err error

	if s.supabaseClient != nil {
		// Obtener de Supabase
		req := s.supabaseClient.GetClient().From("configuraciones").Select("*", "", false).Eq("clave", clave)
		data, count, err := req.Execute()
		if err != nil {
			return nil, fmt.Errorf("error al obtener configuración de Supabase: %v", err)
		}
		if count == 0 {
			return nil, fmt.Errorf("configuración no encontrada: %s", clave)
		}

		// Decodificar respuesta
		err = json.Unmarshal(data, &config)
		if err != nil {
			return nil, fmt.Errorf("error al decodificar configuración: %v", err)
		}
	} else if s.db != nil {
		// Obtener de MongoDB
		err = s.db.Collection("configuraciones").FindOne(
			context.Background(),
			bson.M{"clave": clave},
		).Decode(&config)
		if err != nil {
			return nil, fmt.Errorf("error al obtener configuración: %v", err)
		}
	} else {
		return nil, errors.New("no se ha configurado una fuente de datos")
	}

	// Guardar en caché para futuras consultas
	if s.cacheService != nil {
		ctx := context.Background()
		err = s.cacheService.SetWithExpiration(ctx, cacheKey, &config, 30*time.Minute)
		if err != nil {
			// Solo loguear el error, continuar con la respuesta
			s.logger.Warn("Error al guardar en caché", zap.Error(err))
		}
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

// EliminarConfiguracion elimina una configuración
func (s *ConfiguracionService) EliminarConfiguracion(clave string) error {
	// Eliminar de la base de datos
	var err error

	if s.supabaseClient != nil {
		// Eliminar de Supabase
		req := s.supabaseClient.GetClient().From("configuraciones").Delete("*", "").Eq("clave", clave)
		_, count, err := req.Execute()
		if err != nil {
			return fmt.Errorf("error al eliminar configuración de Supabase: %v", err)
		}
		if count == 0 {
			return fmt.Errorf("configuración no encontrada: %s", clave)
		}
	} else if s.db != nil {
		// Eliminar de MongoDB
		result, err := s.db.Collection("configuraciones").DeleteOne(
			context.Background(),
			bson.M{"clave": clave},
		)
		if err != nil {
			return fmt.Errorf("error al eliminar configuración: %v", err)
		}
		if result.DeletedCount == 0 {
			return fmt.Errorf("configuración no encontrada: %s", clave)
		}
	} else {
		return errors.New("no se ha configurado una fuente de datos")
	}

	// Eliminar de la caché
	cacheKey := fmt.Sprintf("config:%s", clave)
	if s.cacheService != nil {
		ctx := context.Background()
		err = s.cacheService.Delete(ctx, cacheKey)
		if err != nil {
			// Solo loguear el error, continuar con la respuesta
			s.logger.Warn("Error al eliminar de caché", zap.Error(err))
		}
	}

	return nil
}
