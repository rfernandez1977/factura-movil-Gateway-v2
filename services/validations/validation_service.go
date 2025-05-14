package validations

import (
	"github.com/cursor/FMgo/models"
)

// ValidationService proporciona métodos para validar documentos tributarios
type ValidationService struct {
	tributarioValidation *TributarioValidation
	siiValidator         *SIIValidator
}

// NewValidationService crea una nueva instancia del servicio de validación
func NewValidationService(tributarioValidation *TributarioValidation, siiValidator *SIIValidator) *ValidationService {
	return &ValidationService{
		tributarioValidation: tributarioValidation,
		siiValidator:         siiValidator,
	}
}

// ValidateFactura valida una factura
func (s *ValidationService) ValidateFactura(factura *models.Factura) []models.ValidationFieldError {
	var errors []models.ValidationFieldError

	// Validar datos básicos
	errors = append(errors, s.ValidarRUTEmisor(factura.RutEmisor)...)
	errors = append(errors, s.ValidarRUTReceptor(factura.RutReceptor)...)
	errors = append(errors, s.ValidarMontos(factura.MontoNeto, factura.MontoIVA, factura.MontoTotal)...)

	// Validar referencias si existen
	if factura.Referencias != nil && len(factura.Referencias) > 0 {
		errors = append(errors, s.ValidarReferencias(factura)...)
	}

	// Validar items
	if factura.Items != nil && len(factura.Items) > 0 {
		errors = append(errors, s.ValidarItems(factura.Items)...)
	}

	return errors
}

// ValidateBoleta valida una boleta
func (s *ValidationService) ValidateBoleta(boleta *models.Boleta) []models.ValidationFieldError {
	var errors []models.ValidationFieldError

	// Validar datos básicos
	errors = append(errors, s.ValidarRUTEmisor(boleta.RUTEmisor)...)
	errors = append(errors, s.ValidarMontos(boleta.MontoNeto, boleta.MontoIVA, boleta.MontoTotal)...)

	// Validar items si existen
	if boleta.Items != nil && len(boleta.Items) > 0 {
		errors = append(errors, s.ValidarItemsBoleta(boleta.Items)...)
	}

	return errors
}

// ValidarRUTEmisor valida un RUT de emisor
func (s *ValidationService) ValidarRUTEmisor(rut string) []models.ValidationFieldError {
	var errors []models.ValidationFieldError

	// Implementación básica para el ejemplo
	if rut == "" {
		errors = append(errors, models.ValidationFieldError{
			Field:   "rut_emisor",
			Message: "RUT emisor es requerido",
		})
	}

	return errors
}

// ValidarRUTReceptor valida un RUT de receptor
func (s *ValidationService) ValidarRUTReceptor(rut string) []models.ValidationFieldError {
	var errors []models.ValidationFieldError

	// Implementación básica para el ejemplo
	if rut == "" {
		errors = append(errors, models.ValidationFieldError{
			Field:   "rut_receptor",
			Message: "RUT receptor es requerido",
		})
	}

	return errors
}

// ValidarMontos valida los montos
func (s *ValidationService) ValidarMontos(montoNeto, montoIVA, montoTotal float64) []models.ValidationFieldError {
	var errors []models.ValidationFieldError

	// Implementación básica para el ejemplo
	if montoTotal <= 0 {
		errors = append(errors, models.ValidationFieldError{
			Field:   "monto_total",
			Message: "El monto total debe ser mayor a cero",
		})
	}

	return errors
}

// ValidarReferencias valida las referencias
func (s *ValidationService) ValidarReferencias(factura *models.Factura) []models.ValidationFieldError {
	var errors []models.ValidationFieldError

	// Implementación delegada al tributarioValidation
	if s.tributarioValidation != nil {
		return s.tributarioValidation.ValidarReferencias(factura)
	}

	return errors
}

// ValidarItems valida los ítems de una factura
func (s *ValidationService) ValidarItems(items []models.Item) []models.ValidationFieldError {
	var errors []models.ValidationFieldError

	// Implementación básica para el ejemplo
	if len(items) == 0 {
		errors = append(errors, models.ValidationFieldError{
			Field:   "items",
			Message: "La factura debe tener al menos un ítem",
		})
	}

	return errors
}

// ValidarItemsBoleta valida los ítems de una boleta
func (s *ValidationService) ValidarItemsBoleta(items []models.Item) []models.ValidationFieldError {
	var errors []models.ValidationFieldError

	// Implementación básica para el ejemplo
	if len(items) == 0 {
		errors = append(errors, models.ValidationFieldError{
			Field:   "items",
			Message: "La boleta debe tener al menos un ítem",
		})
	}

	return errors
}
