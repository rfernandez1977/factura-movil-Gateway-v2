package caf

import (
	"encoding/xml"
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	// Errores específicos de validación CAF
	ErrCAFInvalido     = errors.New("CAF inválido o mal formado")
	ErrCAFExpirado     = errors.New("CAF expirado")
	ErrFolioNoValido   = errors.New("folio fuera de rango")
	ErrRUTNoCoincide   = errors.New("RUT no coincide con CAF")
	ErrTipoDTEInvalido = errors.New("tipo de DTE no corresponde al CAF")
	ErrFolioUsado      = errors.New("folio ya ha sido utilizado")
)

// CAF representa la estructura básica del CAF
type CAF struct {
	XMLName    xml.Name  `xml:"AUTORIZACION"`
	RutEmisor  string    `xml:"CAF>DA>RE"`
	TipoDTE    int       `xml:"CAF>DA>TD"`
	RangoDesde int       `xml:"CAF>DA>RNG>D"`
	RangoHasta int       `xml:"CAF>DA>RNG>H"`
	FechaDesde time.Time `xml:"CAF>DA>RSAPK>M"`
	FechaHasta time.Time `xml:"CAF>DA>RSAPK>E"`
}

// MemoryFolioStore implementa almacenamiento simple en memoria
type MemoryFolioStore struct {
	mu     sync.RWMutex
	folios map[string]map[int]bool
}

func NewMemoryFolioStore() *MemoryFolioStore {
	return &MemoryFolioStore{
		folios: make(map[string]map[int]bool),
	}
}

func (m *MemoryFolioStore) MarcarUsado(rut string, folio int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.folios[rut]; !exists {
		m.folios[rut] = make(map[int]bool)
	}
	m.folios[rut][folio] = true
	return nil
}

func (m *MemoryFolioStore) EstaUsado(rut string, folio int) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if folios, exists := m.folios[rut]; exists {
		return folios[folio]
	}
	return false
}

// Validator maneja la validación básica de CAF
type Validator struct {
	caf   *CAF
	store *MemoryFolioStore
}

// NewValidator crea un nuevo validador de CAF
func NewValidator(cafXML []byte) (*Validator, error) {
	caf := &CAF{}
	if err := xml.Unmarshal(cafXML, caf); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrCAFInvalido, err)
	}

	return &Validator{
		caf:   caf,
		store: NewMemoryFolioStore(),
	}, nil
}

// ValidarFolio verifica si un folio es válido
func (v *Validator) ValidarFolio(folio int) error {
	// Verificar rango
	if folio < v.caf.RangoDesde || folio > v.caf.RangoHasta {
		return fmt.Errorf("%w: folio %d fuera de rango [%d-%d]",
			ErrFolioNoValido, folio, v.caf.RangoDesde, v.caf.RangoHasta)
	}

	// Verificar si ya fue usado
	if v.store.EstaUsado(v.caf.RutEmisor, folio) {
		return fmt.Errorf("%w: folio %d", ErrFolioUsado, folio)
	}

	// Verificar vigencia
	now := time.Now()
	if now.Before(v.caf.FechaDesde) || now.After(v.caf.FechaHasta) {
		return fmt.Errorf("%w: CAF válido entre %v y %v",
			ErrCAFExpirado, v.caf.FechaDesde, v.caf.FechaHasta)
	}

	return nil
}

// MarcarFolioUsado registra un folio como utilizado
func (v *Validator) MarcarFolioUsado(folio int) error {
	return v.store.MarcarUsado(v.caf.RutEmisor, folio)
}

// ValidarRUTEmisor verifica si el RUT del emisor coincide
func (v *Validator) ValidarRUTEmisor(rut string) error {
	if rut != v.caf.RutEmisor {
		return fmt.Errorf("%w: esperado %s, recibido %s",
			ErrRUTNoCoincide, v.caf.RutEmisor, rut)
	}
	return nil
}

// ValidarTipoDTE verifica si el tipo de DTE coincide
func (v *Validator) ValidarTipoDTE(tipo int) error {
	if tipo != v.caf.TipoDTE {
		return fmt.Errorf("%w: esperado %d, recibido %d",
			ErrTipoDTEInvalido, v.caf.TipoDTE, tipo)
	}
	return nil
}

// ValidarCompleto realiza todas las validaciones básicas
func (v *Validator) ValidarCompleto(rut string, tipo int, folio int) error {
	if err := v.ValidarRUTEmisor(rut); err != nil {
		return err
	}

	if err := v.ValidarTipoDTE(tipo); err != nil {
		return err
	}

	if err := v.ValidarFolio(folio); err != nil {
		return err
	}

	return nil
}
