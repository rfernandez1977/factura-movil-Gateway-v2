package models

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

// GenerateID genera un ID único
func GenerateID() string {
	// Generar 16 bytes aleatorios
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		// En caso de error, usar timestamp
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}

	// Convertir a string hexadecimal
	return hex.EncodeToString(b)
}

// ParseFecha convierte un string a fecha
func ParseFecha(fecha string) (time.Time, error) {
	layouts := []string{
		"2006-01-02",
		"02/01/2006",
		"2006/01/02",
		time.RFC3339,
	}

	var t time.Time
	var err error

	for _, layout := range layouts {
		t, err = time.Parse(layout, fecha)
		if err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("formato de fecha no reconocido: %s", fecha)
}
