package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// ValidateRUT valida un RUT chileno
func ValidateRUT(rut string) error {
	if rut == "" {
		return fmt.Errorf("RUT no puede estar vacío")
	}

	// Eliminar puntos y guión
	cleanRut := strings.ReplaceAll(rut, ".", "")
	cleanRut = strings.ReplaceAll(cleanRut, "-", "")

	// Validar formato
	re := regexp.MustCompile(`^\d{1,8}[0-9kK]$`)
	if !re.MatchString(cleanRut) {
		return fmt.Errorf("formato inválido de RUT: %s", rut)
	}

	// Separar número y dígito verificador
	lastChar := cleanRut[len(cleanRut)-1:]
	numberStr := cleanRut[:len(cleanRut)-1]

	// Convertir a número
	number, err := strconv.Atoi(numberStr)
	if err != nil {
		return fmt.Errorf("número inválido en RUT: %s", rut)
	}

	// Calcular dígito verificador
	var sum int
	var factor int = 2

	// Aplicar algoritmo del módulo 11
	tempNum := number
	for tempNum > 0 {
		digit := tempNum % 10
		sum += digit * factor
		factor++
		if factor > 7 {
			factor = 2
		}
		tempNum = tempNum / 10
	}

	mod := 11 - (sum % 11)
	var expectedDV string

	if mod == 11 {
		expectedDV = "0"
	} else if mod == 10 {
		expectedDV = "K"
	} else {
		expectedDV = strconv.Itoa(mod)
	}

	// Comparar con el dígito verificador proporcionado
	lastCharUpper := strings.ToUpper(lastChar)
	if lastCharUpper != expectedDV {
		return fmt.Errorf("dígito verificador incorrecto para RUT %s. Esperado: %s, Recibido: %s",
			rut, expectedDV, lastChar)
	}

	return nil
}
