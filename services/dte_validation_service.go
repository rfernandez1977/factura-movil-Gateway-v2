package services

import (
	"context"
	"fmt"
	"time"

	"github.com/cursor/FMgo/models"
	"github.com/cursor/FMgo/utils/validation"
	"github.com/go-redis/redis/v8"
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

// RegisterValidation registra una configuración de validación
func (s *DTEValidationService) RegisterValidation(tipo string, config *models.ValidationConfig) error {
	if err := config.Validate(); err != nil {
		return fmt.Errorf("configuración de validación inválida: %v", err)
	}
	s.validations[tipo] = config
	return nil
}

// ValidateDTE valida un documento tributario
func (s *DTEValidationService) ValidateDTE(ctx context.Context, doc *models.DocumentoTributario) ([]*models.ValidationFieldError, error) {
	// Obtener configuración de validación para el tipo de documento
	config, ok := s.validations[doc.TipoDTE]
	if !ok {
		return nil, fmt.Errorf("no hay configuración de validación para el tipo de documento %s", doc.TipoDTE)
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

	return validator.GetErrors(), nil
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
	return nil
}
