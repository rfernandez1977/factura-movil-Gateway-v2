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
	errors = append(errors, s.tributarioValidation.ValidarRUTEmisor(factura.RUTEmisor)...)
	errors = append(errors, s.tributarioValidation.ValidarRUTReceptor(factura.RUTReceptor)...)
	errors = append(errors, s.tributarioValidation.ValidarMontos(factura.MontoNeto, factura.MontoIVA, factura.MontoTotal)...)

	// Validar referencias si existen
	if factura.Referencias != nil && len(factura.Referencias) > 0 {
		errors = append(errors, s.tributarioValidation.ValidarReferencias(factura)...)
	}

	// Validar items
	if factura.Items != nil && len(factura.Items) > 0 {
		errors = append(errors, s.tributarioValidation.ValidarItems(factura.Items)...)
	}

	return errors
}

// ValidateBoleta valida una boleta
func (s *ValidationService) ValidateBoleta(boleta *models.Boleta) []models.ValidationFieldError {
	var errors []models.ValidationFieldError

	// Validar datos básicos
	errors = append(errors, s.tributarioValidation.ValidarRUTEmisor(boleta.RUTEmisor)...)
	errors = append(errors, s.tributarioValidation.ValidarMontos(boleta.MontoNeto, boleta.MontoIVA, boleta.MontoTotal)...)

	// Validar items si existen
	if boleta.Items != nil && len(boleta.Items) > 0 {
		errors = append(errors, s.tributarioValidation.ValidarItemsBoleta(boleta.Items)...)
	}

	return errors
}
