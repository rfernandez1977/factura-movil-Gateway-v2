package validations

import (
	"fmt"
	"strings"

	"FMgo/models"
)

// ValidarRespuestaSII valida una respuesta del SII
func ValidarRespuestaSII(respuesta *models.RespuestaSII) error {
	if respuesta == nil {
		return fmt.Errorf("respuesta del SII nula")
	}

	if respuesta.TrackID == "" {
		return fmt.Errorf("track ID no proporcionado")
	}

	if respuesta.Estado == "" {
		return fmt.Errorf("estado no proporcionado")
	}

	// Verificar si hay errores (en lugar de usar TieneErrores)
	if len(respuesta.Errores) > 0 {
		var errores []string
		for _, err := range respuesta.Errores {
			errores = append(errores, fmt.Sprintf("%s: %s", err.Codigo, err.Descripcion))
		}
		return fmt.Errorf("respuesta con errores: %s", strings.Join(errores, "; "))
	}

	return nil
}
