package caf

import (
	"fmt"
	"sync"
	"time"
)

// Service maneja la gestión de CAFs
type Service struct {
	mu         sync.RWMutex
	validators map[string]*Validator // key: RUT-TipoDTE
}

// NewService crea un nuevo servicio de CAF
func NewService() *Service {
	return &Service{
		validators: make(map[string]*Validator),
	}
}

// RegistrarCAF registra un nuevo CAF en el servicio
func (s *Service) RegistrarCAF(cafXML []byte) error {
	validator, err := NewValidator(cafXML)
	if err != nil {
		return fmt.Errorf("error creando validador: %w", err)
	}

	key := fmt.Sprintf("%s-%d", validator.caf.RutEmisor, validator.caf.TipoDTE)

	s.mu.Lock()
	s.validators[key] = validator
	s.mu.Unlock()

	return nil
}

// ValidarFolio valida un folio para un RUT y tipo de DTE específicos
func (s *Service) ValidarFolio(rut string, tipoDTE int, folio int) error {
	key := fmt.Sprintf("%s-%d", rut, tipoDTE)

	s.mu.RLock()
	validator, exists := s.validators[key]
	s.mu.RUnlock()

	if !exists {
		return fmt.Errorf("no hay CAF registrado para RUT %s y tipo DTE %d", rut, tipoDTE)
	}

	if err := validator.ValidarCompleto(rut, tipoDTE, folio); err != nil {
		return err
	}

	return validator.MarcarFolioUsado(folio)
}

// ObtenerEstadoCAF obtiene el estado actual de un CAF
func (s *Service) ObtenerEstadoCAF(rut string, tipoDTE int) (*EstadoCAF, error) {
	key := fmt.Sprintf("%s-%d", rut, tipoDTE)

	s.mu.RLock()
	validator, exists := s.validators[key]
	s.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("no hay CAF registrado para RUT %s y tipo DTE %d", rut, tipoDTE)
	}

	return &EstadoCAF{
		RutEmisor:  validator.caf.RutEmisor,
		TipoDTE:    validator.caf.TipoDTE,
		RangoDesde: validator.caf.RangoDesde,
		RangoHasta: validator.caf.RangoHasta,
		FechaDesde: validator.caf.FechaDesde,
		FechaHasta: validator.caf.FechaHasta,
	}, nil
}

// EstadoCAF representa el estado actual de un CAF
type EstadoCAF struct {
	RutEmisor  string
	TipoDTE    int
	RangoDesde int
	RangoHasta int
	FechaDesde time.Time
	FechaHasta time.Time
}
