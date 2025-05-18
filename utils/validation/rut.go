package validation

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"FMgo/models"
)

// ValidateRUT valida un RUT chileno
func ValidateRUT(rut string) error {
	if rut == "" {
		return models.NewValidationFieldError("RUT", "REQUIRED_FIELD", "no puede estar vacío", nil)
	}

	// Eliminar puntos y guión
	rut = strings.ReplaceAll(rut, ".", "")
	rut = strings.ReplaceAll(rut, "-", "")

	// Validar formato
	re := regexp.MustCompile(`^\d{1,8}[0-9kK]$`)
	if !re.MatchString(rut) {
		return models.NewValidationFieldError("RUT", "INVALID_FORMAT", "formato inválido", rut)
	}

	// Separar número y dígito verificador
	numero := rut[:len(rut)-1]
	dv := strings.ToUpper(rut[len(rut)-1:])

	// Convertir a número
	num, err := strconv.Atoi(numero)
	if err != nil {
		return models.NewValidationFieldError("RUT", "INVALID_FORMAT", "número inválido", rut)
	}

	// Calcular dígito verificador
	dvCalculado := calcularDV(num)
	if dv != dvCalculado {
		return models.NewValidationFieldError("RUT", "INVALID_CHECK_DIGIT", "dígito verificador inválido", rut)
	}

	return nil
}

// calcularDV calcula el dígito verificador de un RUT según el algoritmo del SII
func calcularDV(rut int) string {
	var suma int
	var multiplicador = 2

	// Convertir a string para procesar dígito por dígito
	rutStr := fmt.Sprintf("%d", rut)
	for i := len(rutStr) - 1; i >= 0; i-- {
		digito, _ := strconv.Atoi(string(rutStr[i]))
		suma += digito * multiplicador
		multiplicador++
		if multiplicador > 7 {
			multiplicador = 2
		}
	}

	// Calcular dígito verificador
	resultado := 11 - (suma % 11)

	// Convertir a string según las reglas del SII
	if resultado == 11 {
		return "0"
	}
	if resultado == 10 {
		return "K"
	}
	return fmt.Sprintf("%d", resultado)
}
