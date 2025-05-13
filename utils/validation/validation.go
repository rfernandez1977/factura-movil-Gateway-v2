package validation

import (
	"time"

	"github.com/cursor/FMgo/models"
)

// ValidateFecha valida una fecha en formato string
func ValidateFecha(fecha string) error {
	if fecha == "" {
		return models.NewValidationError("fecha", "REQUIRED_FIELD", "no puede estar vacía", nil)
	}

	// Intentar parsear la fecha
	_, err := time.Parse("2006-01-02", fecha)
	if err != nil {
		return models.NewValidationError("fecha", "INVALID_FORMAT", "formato inválido (YYYY-MM-DD)", fecha)
	}

	return nil
}

// ValidateMonto valida un monto
func ValidateMonto(monto int) error {
	if monto < 0 {
		return models.NewValidationError("monto", "INVALID_VALUE", "debe ser mayor o igual a 0", monto)
	}
	return nil
}

// ValidatePorcentaje valida un porcentaje
func ValidatePorcentaje(porcentaje int) error {
	if porcentaje < 0 || porcentaje > 100 {
		return models.NewValidationError("porcentaje", "INVALID_VALUE", "debe estar entre 0 y 100", porcentaje)
	}
	return nil
}

// ValidateCantidad valida una cantidad
func ValidateCantidad(cantidad float64) error {
	if cantidad <= 0 {
		return models.NewValidationError("cantidad", "INVALID_VALUE", "debe ser mayor a 0", cantidad)
	}
	return nil
}

// ValidatePrecio valida un precio
func ValidatePrecio(precio int) error {
	if precio < 0 {
		return models.NewValidationError("precio", "INVALID_VALUE", "debe ser mayor o igual a 0", precio)
	}
	return nil
}

// ValidateTexto valida un texto con longitud específica
func ValidateTexto(texto string, minLength, maxLength int) error {
	return ValidateText(texto, minLength, maxLength, "texto")
}

// ValidateNumero valida un número dentro de un rango
func ValidateNumero(numero int, min, max int) error {
	return ValidateNumber(numero, min, max, "numero")
}

// ValidateLista valida una lista con longitud específica
func ValidateLista(lista []string, minLength, maxLength int) error {
	return ValidateList(lista, minLength, maxLength, "lista")
}

// esBisiesto determina si un año es bisiesto
func esBisiesto(anio int) bool {
	return anio%4 == 0 && (anio%100 != 0 || anio%400 == 0)
}
