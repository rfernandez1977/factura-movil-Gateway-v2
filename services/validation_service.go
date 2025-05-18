package services

import (
	"errors"
	"regexp"
	"time"

	"FMgo/domain"
)

// ValidationService implementa la interfaz domain.ValidationService
type ValidationService struct{}

// NewValidationService crea una nueva instancia del servicio de validación
func NewValidationService() domain.ValidationService {
	return &ValidationService{}
}

// ValidarDocumento valida un documento tributario
func (s *ValidationService) ValidarDocumento(doc *domain.DocumentoTributario) error {
	if doc == nil {
		return errors.New("documento no puede ser nulo")
	}

	// Validar RUT emisor
	if err := s.ValidarRUT(doc.RutEmisor); err != nil {
		return err
	}

	// Validar RUT receptor
	if err := s.ValidarRUT(doc.RutReceptor); err != nil {
		return err
	}

	// Validar fecha de emisión
	if err := s.ValidarFecha(doc.FechaEmision); err != nil {
		return err
	}

	// Validar montos
	if err := s.ValidarMonto(doc.MontoTotal); err != nil {
		return err
	}
	if err := s.ValidarMonto(doc.MontoNeto); err != nil {
		return err
	}
	if err := s.ValidarMonto(doc.MontoIVA); err != nil {
		return err
	}
	if err := s.ValidarMonto(doc.MontoExento); err != nil {
		return err
	}

	// Validar que el monto total sea la suma de los otros montos
	montoCalculado := doc.MontoNeto + doc.MontoIVA + doc.MontoExento
	if montoCalculado != doc.MontoTotal {
		return errors.New("el monto total no coincide con la suma de los montos parciales")
	}

	return nil
}

// ValidarRUT valida un RUT chileno
func (s *ValidationService) ValidarRUT(rut string) error {
	// Eliminar puntos y guión
	rut = regexp.MustCompile(`[^0-9kK]`).ReplaceAllString(rut, "")

	// Validar longitud
	if len(rut) < 2 {
		return errors.New("RUT inválido: longitud incorrecta")
	}

	// Separar número y dígito verificador
	numero := rut[:len(rut)-1]
	dv := rut[len(rut)-1:]

	// Calcular dígito verificador
	var suma int
	var multiplicador = 2
	for i := len(numero) - 1; i >= 0; i-- {
		suma += int(numero[i]-'0') * multiplicador
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
		dvCalculado = string(dvEsperado + '0')
	}

	// Comparar dígito verificador
	if dvCalculado != dv {
		return errors.New("RUT inválido: dígito verificador incorrecto")
	}

	return nil
}

// ValidarMonto valida un monto
func (s *ValidationService) ValidarMonto(monto float64) error {
	if monto < 0 {
		return errors.New("el monto no puede ser negativo")
	}
	return nil
}

// ValidarFecha valida una fecha
func (s *ValidationService) ValidarFecha(fecha time.Time) error {
	if fecha.IsZero() {
		return errors.New("la fecha no puede ser cero")
	}
	if fecha.After(time.Now()) {
		return errors.New("la fecha no puede ser futura")
	}
	return nil
}
