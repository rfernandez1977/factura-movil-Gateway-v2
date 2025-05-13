package utils

import (
	"encoding/xml"
	"fmt"
	"time"
)

// DocumentUtils contiene utilidades comunes para documentos
type DocumentUtils struct {
	amountValidator *AmountValidator
}

// NewDocumentUtils crea una nueva instancia de DocumentUtils
func NewDocumentUtils() *DocumentUtils {
	return &DocumentUtils{
		amountValidator: NewAmountValidator(),
	}
}

// IsValidDate valida que una fecha esté dentro de un rango válido
func (u *DocumentUtils) IsValidDate(date time.Time, minDate, maxDate time.Time) bool {
	return !date.IsZero() && (date.After(minDate) || date.Equal(minDate)) && (date.Before(maxDate) || date.Equal(maxDate))
}

// FormatAmount formatea un monto como string con 2 decimales
func (u *DocumentUtils) FormatAmount(amount float64) string {
	return fmt.Sprintf("%.2f", amount)
}

// GenerateXML genera el XML de un documento
func (u *DocumentUtils) GenerateXML(doc interface{}) ([]byte, error) {
	return xml.Marshal(doc)
}
