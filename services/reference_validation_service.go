package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
)

// ReferenceValidationService maneja las validaciones de documentos referenciados
type ReferenceValidationService struct {
	db          *mongo.Database
	cache       *redis.Client
	suggestions *SuggestionService
	siiService  *SIIService
}

// ReferenceValidationResult contiene el resultado de la validación de referencias
type ReferenceValidationResult struct {
	IsValid       bool
	Errors        []ReferenceError
	Suggestions   []*SuggestionResult
	ReferencedDoc interface{}
}

// ReferenceError representa un error en la validación de referencias
type ReferenceError struct {
	Type     string // DOCUMENT_NOT_FOUND, INVALID_REFERENCE, AMOUNT_MISMATCH, etc.
	Message  string
	Field    string
	Value    interface{}
	Severity string // ERROR, WARNING
}

// NewReferenceValidationService crea una nueva instancia del servicio
func NewReferenceValidationService(db *mongo.Database, redisClient *redis.Client, siiService *SIIService) *ReferenceValidationService {
	return &ReferenceValidationService{
		db:          db,
		cache:       redisClient,
		suggestions: NewSuggestionService(redisClient),
		siiService:  siiService,
	}
}

// ValidateReference valida un documento referenciado
func (s *ReferenceValidationService) ValidateReference(ctx context.Context, docType string, reference string) (*ReferenceValidationResult, error) {
	result := &ReferenceValidationResult{
		IsValid: true,
	}

	// Intentar obtener del cache
	cacheKey := fmt.Sprintf("reference:%s:%s", docType, reference)
	if cached, err := s.cache.Get(ctx, cacheKey).Result(); err == nil {
		if err := json.Unmarshal([]byte(cached), &result); err == nil {
			return result, nil
		}
	}

	// Obtener documento referenciado
	referencedDoc, err := s.getReferencedDocument(ctx, docType, reference)
	if err != nil {
		result.IsValid = false
		result.Errors = append(result.Errors, ReferenceError{
			Type:     "DOCUMENT_NOT_FOUND",
			Message:  "Documento referenciado no encontrado",
			Field:    "Referencia",
			Value:    reference,
			Severity: "ERROR",
		})
		return result, nil
	}

	// Validar estado del documento en SII
	siiStatus, err := s.siiService.ConsultarEstado(referencedDoc.TrackID)
	if err != nil {
		result.Errors = append(result.Errors, ReferenceError{
			Type:     "SII_ERROR",
			Message:  "Error al consultar estado en SII",
			Field:    "TrackID",
			Value:    referencedDoc.TrackID,
			Severity: "ERROR",
		})
	}

	// Validar que el documento esté aceptado
	if siiStatus.Estado != "ACEPTADO" {
		result.Errors = append(result.Errors, ReferenceError{
			Type:     "INVALID_STATUS",
			Message:  "El documento referenciado no está aceptado en SII",
			Field:    "Estado",
			Value:    siiStatus.Estado,
			Severity: "ERROR",
		})
	}

	// Obtener sugerencias para los errores encontrados
	if len(result.Errors) > 0 {
		suggestions, err := s.getSIIErrorSuggestions(ctx, result.Errors)
		if err == nil {
			result.Suggestions = suggestions
		}
	}

	// Guardar en cache
	if data, err := json.Marshal(result); err == nil {
		s.cache.Set(ctx, cacheKey, data, 1*time.Hour)
	}

	return result, nil
}

// getReferencedDocument obtiene un documento referenciado
func (s *ReferenceValidationService) getReferencedDocument(ctx context.Context, docType string, reference string) (interface{}, error) {
	// Implementar lógica para obtener el documento según el tipo
	// Por ejemplo, para Nota de Crédito (61) que referencia una Factura (33)
	switch docType {
	case "61": // Nota de Crédito
		return s.getFacturaReferenciada(ctx, reference)
	case "56": // Nota de Débito
		return s.getFacturaReferenciada(ctx, reference)
	default:
		return nil, fmt.Errorf("tipo de documento no soportado: %s", docType)
	}
}

// getFacturaReferenciada obtiene una factura referenciada
func (s *ReferenceValidationService) getFacturaReferenciada(ctx context.Context, reference string) (interface{}, error) {
	// Implementar lógica para obtener la factura
	return nil, nil
}

// getSIIErrorSuggestions obtiene sugerencias para errores del SII
func (s *ReferenceValidationService) getSIIErrorSuggestions(ctx context.Context, errors []ReferenceError) ([]*SuggestionResult, error) {
	var suggestions []*SuggestionResult

	for _, err := range errors {
		// Sugerencias específicas para errores del SII
		switch err.Type {
		case "DOCUMENT_NOT_FOUND":
			suggestions = append(suggestions, &SuggestionResult{
				Field: err.Field,
				Error: err.Message,
				Suggestions: []string{
					"Verifique que el número de referencia sea correcto",
					"Confirme que el documento existe en el sistema",
					"Si el documento es reciente, espere unos minutos y reintente",
				},
				Severity: err.Severity,
				Priority: 1,
			})

		case "INVALID_STATUS":
			suggestions = append(suggestions, &SuggestionResult{
				Field: err.Field,
				Error: err.Message,
				Suggestions: []string{
					"Espere a que el documento sea aceptado por el SII",
					"Verifique el estado del documento en el portal del SII",
					"Si el documento fue rechazado, corrija los errores y reenvíelo",
				},
				Severity: err.Severity,
				Priority: 1,
			})

		case "SII_ERROR":
			suggestions = append(suggestions, &SuggestionResult{
				Field: err.Field,
				Error: err.Message,
				Suggestions: []string{
					"Verifique la conexión con el SII",
					"Intente nuevamente en unos minutos",
					"Si el problema persiste, contacte al soporte",
				},
				Severity: err.Severity,
				Priority: 2,
			})
		}
	}

	return suggestions, nil
}

// ValidateAmount valida el monto de un documento referenciado
func (s *ReferenceValidationService) ValidateAmount(ctx context.Context, docType string, reference string, amount float64) error {
	// Obtener documento referenciado
	referencedDoc, err := s.getReferencedDocument(ctx, docType, reference)
	if err != nil {
		return err
	}

	// Validar montos según tipo de documento
	switch docType {
	case "61": // Nota de Crédito
		if amount > referencedDoc.MontoTotal {
			return fmt.Errorf("el monto de la nota de crédito (%v) excede el monto total del documento referenciado (%v)",
				amount, referencedDoc.MontoTotal)
		}
	case "56": // Nota de Débito
		// Validaciones específicas para nota de débito
		if amount <= 0 {
			return fmt.Errorf("el monto de la nota de débito debe ser positivo")
		}
	}

	return nil
}

// ValidateTiming valida el timing de un documento referenciado
func (s *ReferenceValidationService) ValidateTiming(ctx context.Context, docType string, reference string, date time.Time) error {
	// Obtener documento referenciado
	referencedDoc, err := s.getReferencedDocument(ctx, docType, reference)
	if err != nil {
		return err
	}

	// Validar fechas según tipo de documento
	switch docType {
	case "61": // Nota de Crédito
		if date.Before(referencedDoc.FechaEmision) {
			return fmt.Errorf("la fecha de la nota de crédito no puede ser anterior a la fecha del documento referenciado")
		}
	case "56": // Nota de Débito
		// Validaciones específicas para nota de débito
		if date.Before(referencedDoc.FechaEmision) {
			return fmt.Errorf("la fecha de la nota de débito no puede ser anterior a la fecha del documento referenciado")
		}
	}

	return nil
}
