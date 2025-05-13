package validation

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/cursor/FMgo/models"
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

	// Calcular dígito verificador
	var suma int
	var multiplicador = 2
	for i := len(numero) - 1; i >= 0; i-- {
		digito, _ := strconv.Atoi(string(numero[i]))
		suma += digito * multiplicador
		multiplicador++
		if multiplicador > 7 {
			multiplicador = 2
		}
	}

	// Calcular dígito verificador esperado
	dvEsperado := 11 - (suma % 11)
	var dvCalculado string
	if dvEsperado == 11 {
		dvCalculado = "0"
	} else if dvEsperado == 10 {
		dvCalculado = "K"
	} else {
		dvCalculado = strconv.Itoa(dvEsperado)
	}

	// Comparar dígito verificador
	if dv != dvCalculado {
		return models.NewValidationFieldError("RUT", "INVALID_CHECK_DIGIT", "dígito verificador inválido", rut)
	}

	return nil
}

// ValidateEmail valida un correo electrónico
func ValidateEmail(email string) error {
	if email == "" {
		return models.NewValidationFieldError("Email", "REQUIRED_FIELD", "no puede estar vacío", nil)
	}

	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !re.MatchString(email) {
		return models.NewValidationFieldError("Email", "INVALID_FORMAT", "formato inválido", email)
	}

	return nil
}

// ValidateDate valida una fecha
func ValidateDate(date time.Time, fieldName string) error {
	if date.IsZero() {
		return models.NewValidationFieldError(fieldName, "REQUIRED_FIELD", "no puede estar vacía", nil)
	}
	return nil
}

// ValidateDateRange valida que una fecha esté dentro de un rango
func ValidateDateRange(date, minDate, maxDate time.Time, fieldName string) error {
	if err := ValidateDate(date, fieldName); err != nil {
		return err
	}

	if date.Before(minDate) {
		return models.NewValidationFieldError(fieldName,
			"DATE_BEFORE_MIN",
			fmt.Sprintf("no puede ser anterior a %s", minDate.Format("2006-01-02")),
			date.Format("2006-01-02"))
	}

	if date.After(maxDate) {
		return models.NewValidationFieldError(fieldName,
			"DATE_AFTER_MAX",
			fmt.Sprintf("no puede ser posterior a %s", maxDate.Format("2006-01-02")),
			date.Format("2006-01-02"))
	}

	return nil
}

// ValidateText valida un texto
func ValidateText(text string, minLength, maxLength int, fieldName string) error {
	if text == "" {
		return models.NewValidationFieldError(fieldName, "REQUIRED_FIELD", "no puede estar vacío", nil)
	}

	length := len(strings.TrimSpace(text))
	if length < minLength {
		return models.NewValidationFieldError(fieldName,
			"TEXT_TOO_SHORT",
			fmt.Sprintf("debe tener al menos %d caracteres", minLength),
			text)
	}
	if length > maxLength {
		return models.NewValidationFieldError(fieldName,
			"TEXT_TOO_LONG",
			fmt.Sprintf("no debe exceder %d caracteres", maxLength),
			text)
	}

	return nil
}

// ValidateNumber valida un número
func ValidateNumber(number int, min, max int, fieldName string) error {
	if number < min {
		return models.NewValidationFieldError(fieldName,
			"NUMBER_BELOW_MIN",
			fmt.Sprintf("debe ser mayor o igual a %d", min),
			number)
	}
	if number > max {
		return models.NewValidationFieldError(fieldName,
			"NUMBER_ABOVE_MAX",
			fmt.Sprintf("debe ser menor o igual a %d", max),
			number)
	}
	return nil
}

// ValidateList valida una lista
func ValidateList(list []string, minLength, maxLength int, fieldName string) error {
	if len(list) < minLength {
		return models.NewValidationFieldError(fieldName,
			"LIST_TOO_SHORT",
			fmt.Sprintf("debe tener al menos %d elementos", minLength),
			list)
	}
	if len(list) > maxLength {
		return models.NewValidationFieldError(fieldName,
			"LIST_TOO_LONG",
			fmt.Sprintf("no debe exceder %d elementos", maxLength),
			list)
	}
	return nil
}
