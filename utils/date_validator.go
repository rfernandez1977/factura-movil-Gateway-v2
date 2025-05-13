package utils

import (
	"fmt"
	"time"
)

// DateValidator define la validación de fechas
type DateValidator struct {
	minDate time.Time
	maxDate time.Time
}

// NewDateValidator crea una nueva instancia de DateValidator
func NewDateValidator() *DateValidator {
	return &DateValidator{
		minDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local),
		maxDate: time.Now().AddDate(1, 0, 0), // Un año en el futuro como máximo
	}
}

// ValidateDate valida que una fecha esté dentro del rango permitido
func (v *DateValidator) ValidateDate(date time.Time, field string) error {
	if date.IsZero() {
		return fmt.Errorf("el campo %s no puede estar vacío", field)
	}

	if date.Before(v.minDate) {
		return fmt.Errorf("el campo %s no puede ser anterior a %s", field, v.minDate.Format("02/01/2006"))
	}

	if date.After(v.maxDate) {
		return fmt.Errorf("el campo %s no puede ser posterior a %s", field, v.maxDate.Format("02/01/2006"))
	}

	return nil
}

// ValidateDateRange valida que una fecha de inicio sea anterior a una fecha de fin
func (v *DateValidator) ValidateDateRange(start, end time.Time, startField, endField string) error {
	if start.IsZero() {
		return fmt.Errorf("el campo %s no puede estar vacío", startField)
	}

	if end.IsZero() {
		return fmt.Errorf("el campo %s no puede estar vacío", endField)
	}

	if start.After(end) {
		return fmt.Errorf("el campo %s debe ser anterior al campo %s", startField, endField)
	}

	return nil
}

// ValidateDueDate valida la fecha de vencimiento
func (v *DateValidator) ValidateDueDate(issueDate, dueDate time.Time) error {
	if dueDate.IsZero() {
		return nil // Fecha de vencimiento es opcional
	}

	if issueDate.IsZero() {
		return fmt.Errorf("la fecha de emisión es requerida para validar la fecha de vencimiento")
	}

	if dueDate.Before(issueDate) {
		return fmt.Errorf("la fecha de vencimiento no puede ser anterior a la fecha de emisión")
	}

	maxDueDate := issueDate.AddDate(1, 0, 0) // Máximo un año después de la emisión
	if dueDate.After(maxDueDate) {
		return fmt.Errorf("la fecha de vencimiento no puede ser más de un año posterior a la fecha de emisión")
	}

	return nil
}

// FormatDate formatea una fecha en el formato dd/mm/yyyy
func (v *DateValidator) FormatDate(date time.Time) string {
	return date.Format("02/01/2006")
}

// ValidateDateOrder valida el orden de dos fechas
func (v *DateValidator) ValidateDateOrder(startDate time.Time, endDate time.Time, startFieldName string, endFieldName string) error {
	if err := v.ValidateDate(startDate, startFieldName); err != nil {
		return err
	}

	if err := v.ValidateDate(endDate, endFieldName); err != nil {
		return err
	}

	if endDate.Before(startDate) {
		return fmt.Errorf("%s no puede ser anterior a %s", endFieldName, startFieldName)
	}

	return nil
}

// ValidateBusinessDays valida que una fecha sea día hábil
// isEmission indica si la fecha es de emisión (true) o vencimiento (false)
func (v *DateValidator) ValidateBusinessDays(date time.Time, fieldName string, isEmission bool) error {
	if err := v.ValidateDate(date, fieldName); err != nil {
		return err
	}

	// Si es fecha de emisión, permitimos cualquier día
	if isEmission {
		return nil
	}

	// Solo validamos días hábiles para fechas de vencimiento
	if date.Weekday() == time.Saturday || date.Weekday() == time.Sunday {
		return fmt.Errorf("%s no puede ser fin de semana", fieldName)
	}

	return nil
}

// ValidateHoliday valida que una fecha no sea feriado
// isEmission indica si la fecha es de emisión (true) o vencimiento (false)
func (v *DateValidator) ValidateHoliday(date time.Time, fieldName string, isEmission bool) error {
	if err := v.ValidateDate(date, fieldName); err != nil {
		return err
	}

	// Si es fecha de emisión, permitimos cualquier día
	if isEmission {
		return nil
	}

	// Lista de feriados en Chile (ejemplo)
	holidays := map[string]bool{
		"01-01": true, // Año Nuevo
		"05-01": true, // Día del Trabajo
		"09-18": true, // Fiestas Patrias
		"09-19": true, // Día de las Glorias del Ejército
		"12-25": true, // Navidad
	}

	dateStr := date.Format("01-02")
	if holidays[dateStr] {
		return fmt.Errorf("%s no puede ser feriado", fieldName)
	}

	return nil
}

// ValidateDocumentDates valida las fechas de un documento
func (v *DateValidator) ValidateDocumentDates(fechaEmision time.Time, fechaVencimiento time.Time) error {
	// Validar fecha de emisión
	if err := v.ValidateDate(fechaEmision, "fecha de emisión"); err != nil {
		return err
	}

	// Validar fecha de vencimiento
	if err := v.ValidateDate(fechaVencimiento, "fecha de vencimiento"); err != nil {
		return err
	}

	// Validar orden de fechas
	if fechaVencimiento.Before(fechaEmision) {
		return fmt.Errorf("la fecha de vencimiento no puede ser anterior a la fecha de emisión")
	}

	// La fecha de vencimiento no puede ser más de 365 días posterior a la de emisión
	if fechaVencimiento.Sub(fechaEmision) > 365*24*time.Hour {
		return fmt.Errorf("la fecha de vencimiento no puede ser más de 365 días posterior a la de emisión")
	}

	return nil
}

// ParseDate parsea un string a fecha
func (v *DateValidator) ParseDate(dateStr string) (time.Time, error) {
	return time.Parse("2006-01-02", dateStr)
}
