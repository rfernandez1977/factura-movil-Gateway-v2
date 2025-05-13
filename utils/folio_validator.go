package utils

import (
	"fmt"
	"regexp"
	"strconv"
)

// FolioValidator valida los folios de un documento
type FolioValidator struct {
	Folio         string
	TipoDocumento string
	Serie         string
}

// Validate valida el folio del documento
func (v *FolioValidator) Validate() error {
	// Validar formato básico
	if !isValidFormat(v.Folio) {
		return fmt.Errorf("formato de folio inválido")
	}

	// Validar longitud según tipo de documento
	if err := v.validateLength(); err != nil {
		return err
	}

	// Validar que sea numérico
	if !isNumeric(v.Folio) {
		return fmt.Errorf("el folio debe contener solo números")
	}

	// Validar rango según tipo de documento
	if err := v.validateRange(); err != nil {
		return err
	}

	return nil
}

// isValidFormat valida el formato del folio
func isValidFormat(folio string) bool {
	// El folio debe contener solo números
	pattern := `^[0-9]+$`
	matched, _ := regexp.MatchString(pattern, folio)
	return matched
}

// validateLength valida la longitud del folio según el tipo de documento
func (v *FolioValidator) validateLength() error {
	minLength := 1
	maxLength := 8

	switch v.TipoDocumento {
	case "FACTURA":
		minLength = 1
		maxLength = 8
	case "BOLETA":
		minLength = 1
		maxLength = 8
	case "NOTA_CREDITO":
		minLength = 1
		maxLength = 8
	case "NOTA_VENTA":
		minLength = 1
		maxLength = 8
	default:
		return fmt.Errorf("tipo de documento no válido")
	}

	if len(v.Folio) < minLength || len(v.Folio) > maxLength {
		return fmt.Errorf("la longitud del folio debe estar entre %d y %d caracteres", minLength, maxLength)
	}

	return nil
}

// validateRange valida el rango del folio según el tipo de documento
func (v *FolioValidator) validateRange() error {
	folioNum, err := strconv.Atoi(v.Folio)
	if err != nil {
		return fmt.Errorf("el folio debe ser un número válido")
	}

	minValue := 1
	maxValue := 99999999

	switch v.TipoDocumento {
	case "FACTURA":
		minValue = 1
		maxValue = 99999999
	case "BOLETA":
		minValue = 1
		maxValue = 99999999
	case "NOTA_CREDITO":
		minValue = 1
		maxValue = 99999999
	case "NOTA_VENTA":
		minValue = 1
		maxValue = 99999999
	}

	if folioNum < minValue || folioNum > maxValue {
		return fmt.Errorf("el folio debe estar entre %d y %d", minValue, maxValue)
	}

	return nil
}

// isNumeric verifica si una cadena contiene solo números
func isNumeric(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

// FormatFolio formatea el folio con ceros a la izquierda
func FormatFolio(folio string, length int) string {
	return fmt.Sprintf("%0"+strconv.Itoa(length)+"s", folio)
}

// GenerateNextFolio genera el siguiente folio en la secuencia
func GenerateNextFolio(currentFolio string) (string, error) {
	if !isNumeric(currentFolio) {
		return "", fmt.Errorf("el folio actual no es un número válido")
	}

	current, err := strconv.Atoi(currentFolio)
	if err != nil {
		return "", err
	}

	next := current + 1
	return strconv.Itoa(next), nil
}
