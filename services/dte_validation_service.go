package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"FMgo/models"
	"FMgo/utils/validation"
	"github.com/go-redis/redis/v8"
)

const (
	// TTLValidationConfig tiempo de vida de la configuración de validación en caché
	TTLValidationConfig = 24 * time.Hour
	// TTLValidationResult tiempo de vida del resultado de validación en caché
	TTLValidationResult = 1 * time.Hour
	// PrefijoValidationConfig prefijo para las claves de configuración
	PrefijoValidationConfig = "val_config:"
	// PrefijoValidationResult prefijo para las claves de resultados
	PrefijoValidationResult = "val_result:"
)

// DTEValidationRule representa una regla de validación específica para DTE
type DTEValidationRule struct {
	Field    string
	Required bool
	Validate func(interface{}) error
}

// DTEValidationConfig representa la configuración de validación para un tipo de DTE
type DTEValidationConfig struct {
	RequiredFields map[string]DTEValidationRule
	FormatRules    []DTEValidationRule
	BusinessRules  []DTEValidationRule
}

// DTESuggestionResult representa una sugerencia de corrección para un DTE
type DTESuggestionResult struct {
	Field   string
	Value   interface{}
	Fix     interface{}
	Message string
}

// DTEValidationService maneja la validación de documentos tributarios
type DTEValidationService struct {
	validations map[string]*models.ValidationConfig
	cache       *redis.Client
	suggestions *SuggestionService
}

// NewDTEValidationService crea una nueva instancia del servicio de validación
func NewDTEValidationService(redisClient *redis.Client) *DTEValidationService {
	return &DTEValidationService{
		validations: make(map[string]*models.ValidationConfig),
		cache:       redisClient,
		suggestions: NewSuggestionService(redisClient),
	}
}

// getConfigCacheKey genera una clave de caché para la configuración
func (s *DTEValidationService) getConfigCacheKey(tipo string) string {
	return fmt.Sprintf("%s%s", PrefijoValidationConfig, tipo)
}

// getResultCacheKey genera una clave de caché para el resultado
func (s *DTEValidationService) getResultCacheKey(doc *models.DocumentoTributario) string {
	return fmt.Sprintf("%s%s:%d", PrefijoValidationResult, doc.TipoDTE, doc.Folio)
}

// RegisterValidation registra una configuración de validación
func (s *DTEValidationService) RegisterValidation(ctx context.Context, tipo string, config *models.ValidationConfig) error {
	if err := config.Validate(); err != nil {
		return fmt.Errorf("configuración de validación inválida: %v", err)
	}

	// Guardar en memoria
	s.validations[tipo] = config

	// Guardar en caché
	data, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("error serializando configuración: %w", err)
	}

	key := s.getConfigCacheKey(tipo)
	if err := s.cache.Set(ctx, key, data, TTLValidationConfig).Err(); err != nil {
		return fmt.Errorf("error guardando configuración en caché: %w", err)
	}

	return nil
}

// getValidationConfig obtiene la configuración de validación
func (s *DTEValidationService) getValidationConfig(ctx context.Context, tipo string) (*models.ValidationConfig, error) {
	// Intentar obtener de memoria
	if config, ok := s.validations[tipo]; ok {
		return config, nil
	}

	// Intentar obtener de caché
	key := s.getConfigCacheKey(tipo)
	data, err := s.cache.Get(ctx, key).Bytes()
	if err == nil {
		var config models.ValidationConfig
		if err := json.Unmarshal(data, &config); err == nil {
			// Guardar en memoria para futuras consultas
			s.validations[tipo] = &config
			return &config, nil
		}
	}

	return nil, fmt.Errorf("no hay configuración de validación para el tipo de documento %s", tipo)
}

// ValidateDTE valida un documento tributario
func (s *DTEValidationService) ValidateDTE(ctx context.Context, doc *models.DocumentoTributario) ([]*models.ValidationFieldError, error) {
	// Intentar obtener resultado de caché
	key := s.getResultCacheKey(doc)
	data, err := s.cache.Get(ctx, key).Bytes()
	if err == nil {
		var errors []*models.ValidationFieldError
		if err := json.Unmarshal(data, &errors); err == nil {
			return errors, nil
		}
	}

	// Obtener configuración de validación
	config, err := s.getValidationConfig(ctx, doc.TipoDTE)
	if err != nil {
		return nil, err
	}

	validator := &models.BaseValidator{}

	// Validar cada regla
	for _, regla := range config.Reglas {
		// Obtener el valor del campo a validar
		valor := doc.GetField(regla.Expresion)
		if valor == nil {
			validator.AddError(regla.Expresion, "campo no encontrado", "FIELD_NOT_FOUND")
			continue
		}

		// Validar según el tipo de regla
		switch regla.Tipo {
		case "required":
			if valor == nil || valor == "" {
				validator.AddError(regla.Expresion, regla.Mensaje, "FIELD_REQUIRED")
			}
		case "format":
			if err := validateFormat(valor, regla.Expresion); err != nil {
				validator.AddError(regla.Expresion, regla.Mensaje, "FORMAT_ERROR", fmt.Sprintf("%v", valor))
			}
		case "range":
			if err := validateRange(valor, regla.Expresion); err != nil {
				validator.AddError(regla.Expresion, regla.Mensaje, "RANGE_ERROR", fmt.Sprintf("%v", valor))
			}
		case "custom":
			if err := validateCustom(valor, regla.Expresion); err != nil {
				validator.AddError(regla.Expresion, regla.Mensaje, "CUSTOM_ERROR", fmt.Sprintf("%v", valor))
			}
		}

		// Si se activó stopOnError y hay errores, detener la validación
		if config.StopOnError && validator.HasErrors() {
			break
		}

		// Si se alcanzó el máximo de errores, detener la validación
		if len(validator.GetErrors()) >= config.MaxErrores {
			break
		}
	}

	errors := validator.GetErrors()

	// Guardar resultado en caché
	if data, err := json.Marshal(errors); err == nil {
		if err := s.cache.Set(ctx, key, data, TTLValidationResult).Err(); err != nil {
			fmt.Printf("error guardando resultado en caché: %v\n", err)
		}
	}

	return errors, nil
}

// validateFormat valida el formato de un valor
func validateFormat(valor interface{}, expresion string) error {
	switch expresion {
	case "rut":
		if str, ok := valor.(string); ok {
			return validation.ValidateRUT(str)
		}
		return fmt.Errorf("valor no es un string")
	case "email":
		if str, ok := valor.(string); ok {
			return validation.ValidateEmail(str)
		}
		return fmt.Errorf("valor no es un string")
	case "date":
		if t, ok := valor.(time.Time); ok {
			return validation.ValidateDate(t, "fecha")
		}
		return fmt.Errorf("valor no es una fecha")
	default:
		return fmt.Errorf("formato no soportado: %s", expresion)
	}
}

// validateRange valida que un valor esté dentro de un rango
func validateRange(valor interface{}, expresion string) error {
	// Implementar validación de rangos según la expresión
	return nil
}

// validateCustom valida un valor con una regla personalizada
func validateCustom(valor interface{}, expresion string) error {
	// Implementar validación personalizada según la expresión
	return nil
}

// ApplySuggestions aplica sugerencias de corrección a un documento
func (s *DTEValidationService) ApplySuggestions(ctx context.Context, doc *models.DocumentoTributario, suggestions []*models.Suggestion) error {
	for _, suggestion := range suggestions {
		if err := doc.SetField(suggestion.Campo, suggestion.Valor); err != nil {
			return fmt.Errorf("error aplicando sugerencia para %s: %v", suggestion.Campo, err)
		}
	}

	// Invalidar resultado de validación en caché
	key := s.getResultCacheKey(doc)
	if err := s.cache.Del(ctx, key).Err(); err != nil {
		fmt.Printf("error invalidando resultado en caché: %v\n", err)
	}

	return nil
}

// LimpiarCache limpia el caché de validación
func (s *DTEValidationService) LimpiarCache(ctx context.Context) error {
	var cursor uint64
	var keys []string

	// Obtener todas las claves con los prefijos
	for {
		var result []string
		var err error
		result, cursor, err = s.cache.Scan(ctx, cursor, PrefijoValidationConfig+"*", 100).Result()
		if err != nil {
			return fmt.Errorf("error escaneando claves de configuración: %w", err)
		}
		keys = append(keys, result...)

		result, cursor, err = s.cache.Scan(ctx, cursor, PrefijoValidationResult+"*", 100).Result()
		if err != nil {
			return fmt.Errorf("error escaneando claves de resultados: %w", err)
		}
		keys = append(keys, result...)

		if cursor == 0 {
			break
		}
	}

	// Eliminar todas las claves encontradas
	if len(keys) > 0 {
		if err := s.cache.Del(ctx, keys...).Err(); err != nil {
			return fmt.Errorf("error eliminando claves del caché: %w", err)
		}
	}

	// Limpiar caché en memoria
	s.validations = make(map[string]*models.ValidationConfig)

	return nil
}
