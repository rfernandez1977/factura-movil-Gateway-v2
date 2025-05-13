package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cursor/FMgo/models"
	"github.com/go-redis/redis/v8"
)

// SuggestionService proporciona sugerencias para diferentes tipos de errores
type SuggestionService struct {
	cache *redis.Client
}

// SuggestionResult representa el resultado de una sugerencia
type SuggestionResult struct {
	Field       string   `json:"field"`
	Error       string   `json:"error"`
	Suggestions []string `json:"suggestions"`
	Severity    string   `json:"severity"`
	Priority    int      `json:"priority"`
}

// NewSuggestionService crea una nueva instancia del servicio de sugerencias
func NewSuggestionService(redisClient *redis.Client) *SuggestionService {
	return &SuggestionService{
		cache: redisClient,
	}
}

// GetSuggestions obtiene sugerencias para un error específico
func (s *SuggestionService) GetSuggestions(ctx context.Context, errorType string, field string) ([]*SuggestionResult, error) {
	cacheKey := fmt.Sprintf("suggestion:%s:%s", errorType, field)

	// Intentar obtener del cache
	if cached, err := s.cache.Get(ctx, cacheKey).Result(); err == nil {
		var suggestions []*SuggestionResult
		if err := json.Unmarshal([]byte(cached), &suggestions); err == nil {
			return suggestions, nil
		}
	}

	// Si no está en cache, generar sugerencias
	suggestions := s.generateSuggestions(errorType, field)

	// Guardar en cache
	if data, err := json.Marshal(suggestions); err == nil {
		s.cache.Set(ctx, cacheKey, data, 24*time.Hour)
	}

	return suggestions, nil
}

// generateSuggestions genera sugerencias para un error específico
func (s *SuggestionService) generateSuggestions(errorType string, field string) []*SuggestionResult {
	var suggestions []*SuggestionResult

	switch errorType {
	case "INVALID_FORMAT":
		suggestions = append(suggestions, &SuggestionResult{
			Field: field,
			Error: "Formato inválido",
			Suggestions: []string{
				"Verifique que el formato sea correcto",
				"Ejemplo de formato correcto: ...",
			},
			Severity: "ERROR",
			Priority: 1,
		})
	case "REQUIRED_FIELD":
		suggestions = append(suggestions, &SuggestionResult{
			Field: field,
			Error: "Campo requerido",
			Suggestions: []string{
				"Este campo es obligatorio",
				"Debe proporcionar un valor para este campo",
			},
			Severity: "ERROR",
			Priority: 1,
		})
	case "INVALID_VALUE":
		suggestions = append(suggestions, &SuggestionResult{
			Field: field,
			Error: "Valor inválido",
			Suggestions: []string{
				"El valor proporcionado no es válido",
				"Intente con un valor diferente",
			},
			Severity: "ERROR",
			Priority: 1,
		})
	}

	return suggestions
}

// GenerateSuggestions genera sugerencias de corrección para un documento
func (s *SuggestionService) GenerateSuggestions(ctx context.Context, doc *models.DocumentoTributario, errores []*models.ValidationError) ([]*models.Suggestion, error) {
	var sugerencias []*models.Suggestion

	// Generar sugerencias para cada error
	for _, err := range errores {
		// Generar sugerencia según el tipo de error
		var tipo, mensaje string
		switch err.Code {
		case "FIELD_REQUIRED":
			tipo = "required"
			mensaje = fmt.Sprintf("El campo %s es requerido", err.Field)
		case "FORMAT_ERROR":
			tipo = "format"
			mensaje = fmt.Sprintf("El formato del campo %s es inválido", err.Field)
		case "RANGE_ERROR":
			tipo = "range"
			mensaje = fmt.Sprintf("El valor del campo %s está fuera de rango", err.Field)
		case "CUSTOM_ERROR":
			tipo = "custom"
			mensaje = fmt.Sprintf("El campo %s no cumple con la validación personalizada", err.Field)
		default:
			tipo = "unknown"
			mensaje = fmt.Sprintf("Error desconocido en el campo %s", err.Field)
		}

		sugerencia := models.NewSuggestion(
			doc.ID,
			err.ID,
			err.Field,
			tipo,
			mensaje,
			err.Value,
		)

		sugerencias = append(sugerencias, sugerencia)
	}

	// Guardar sugerencias en caché
	if err := s.cacheSuggestions(ctx, doc.ID, sugerencias); err != nil {
		return nil, fmt.Errorf("error al guardar sugerencias en caché: %v", err)
	}

	return sugerencias, nil
}

// cacheSuggestions guarda las sugerencias en caché
func (s *SuggestionService) cacheSuggestions(ctx context.Context, docID string, sugerencias []*models.Suggestion) error {
	// Convertir sugerencias a JSON
	data, err := json.Marshal(sugerencias)
	if err != nil {
		return fmt.Errorf("error al serializar sugerencias: %v", err)
	}

	// Guardar en Redis con expiración de 1 hora
	key := fmt.Sprintf("suggestions:%s", docID)
	if err := s.cache.Set(ctx, key, data, time.Hour).Err(); err != nil {
		return fmt.Errorf("error al guardar en Redis: %v", err)
	}

	return nil
}

// getCachedSuggestions obtiene las sugerencias de caché
func (s *SuggestionService) getCachedSuggestions(ctx context.Context, docID string) ([]*models.Suggestion, error) {
	// Obtener de Redis
	key := fmt.Sprintf("suggestions:%s", docID)
	data, err := s.cache.Get(ctx, key).Bytes()
	if err != nil {
		return nil, fmt.Errorf("error al obtener de Redis: %v", err)
	}

	// Deserializar sugerencias
	var sugerencias []*models.Suggestion
	if err := json.Unmarshal(data, &sugerencias); err != nil {
		return nil, fmt.Errorf("error al deserializar sugerencias: %v", err)
	}

	return sugerencias, nil
}
