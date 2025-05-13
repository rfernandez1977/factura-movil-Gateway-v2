package caf

import (
	"crypto/rsa"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// CAFXml representa un archivo de autorización de folios en formato XML
type CAFXml struct {
	XMLName     xml.Name `xml:"AUTORIZACION"`
	Version     string   `xml:"version,attr"`
	RUTEmisor   string   `xml:"RUTEmisor"`
	TipoDTE     int      `xml:"TipoDTE"`
	FolioInicio int      `xml:"FolioInicio"`
	FolioFinal  int      `xml:"FolioFinal"`
	FechaResol  string   `xml:"FechaResol"`
	NumResol    string   `xml:"NumResol"`
	Signature   string   `xml:"Signature"`
}

// Manager maneja los CAF disponibles
type Manager struct {
	mu      sync.RWMutex
	cafs    map[int]*CAFXml // key: tipoDTE
	keyPath string
}

// NewManager crea un nuevo administrador de CAF
func NewManager(keyPath string) *Manager {
	return &Manager{
		cafs:    make(map[int]*CAFXml),
		keyPath: keyPath,
	}
}

// LoadCAF carga un CAF desde un archivo XML
func (m *Manager) LoadCAF(filename string) error {
	// Leer archivo
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("error al leer archivo CAF: %w", err)
	}

	// Decodificar XML
	var caf CAFXml
	err = xml.Unmarshal(data, &caf)
	if err != nil {
		return fmt.Errorf("error al decodificar CAF: %w", err)
	}

	// Validar CAF
	if err := m.validateCAF(&caf); err != nil {
		return fmt.Errorf("error al validar CAF: %w", err)
	}

	// Guardar CAF
	m.mu.Lock()
	m.cafs[caf.TipoDTE] = &caf
	m.mu.Unlock()

	return nil
}

// LoadCAFs carga todos los CAF de un directorio
func (m *Manager) LoadCAFs(dir string) error {
	// Leer directorio
	files, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("error al leer directorio de CAF: %w", err)
	}

	// Cargar cada archivo
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".xml" {
			err := m.LoadCAF(filepath.Join(dir, file.Name()))
			if err != nil {
				return fmt.Errorf("error al cargar CAF %s: %w", file.Name(), err)
			}
		}
	}

	return nil
}

// GetCAF obtiene un CAF para un tipo de DTE
func (m *Manager) GetCAF(tipoDTE int) (*CAFXml, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	caf, ok := m.cafs[tipoDTE]
	if !ok {
		return nil, fmt.Errorf("no hay CAF disponible para el tipo de DTE %d", tipoDTE)
	}

	return caf, nil
}

// GetNextFolio obtiene el siguiente folio disponible para un tipo de DTE
func (m *Manager) GetNextFolio(tipoDTE int) (int, error) {
	caf, err := m.GetCAF(tipoDTE)
	if err != nil {
		return 0, err
	}

	// TODO: Implementar lógica para obtener el siguiente folio disponible
	// Esto requerirá mantener un registro de los folios usados

	return caf.FolioInicio, nil
}

// validateCAF valida un CAF
func (m *Manager) validateCAF(caf *CAFXml) error {
	// Validar versión
	if caf.Version == "" {
		return fmt.Errorf("versión no puede estar vacía")
	}

	// Validar RUT emisor
	if caf.RUTEmisor == "" {
		return fmt.Errorf("RUT emisor no puede estar vacío")
	}

	// Validar tipo DTE
	if caf.TipoDTE <= 0 {
		return fmt.Errorf("tipo DTE inválido")
	}

	// Validar folios
	if caf.FolioInicio <= 0 {
		return fmt.Errorf("folio inicial inválido")
	}
	if caf.FolioFinal <= 0 {
		return fmt.Errorf("folio final inválido")
	}
	if caf.FolioInicio > caf.FolioFinal {
		return fmt.Errorf("folio inicial no puede ser mayor que el final")
	}

	// Validar fecha resolución
	if caf.FechaResol == "" {
		return fmt.Errorf("fecha resolución no puede estar vacía")
	}

	// Validar número resolución
	if caf.NumResol == "" {
		return fmt.Errorf("número resolución no puede estar vacío")
	}

	// Validar firma
	if caf.Signature == "" {
		return fmt.Errorf("firma no puede estar vacía")
	}

	// TODO: Implementar validación de firma usando la clave pública del SII

	return nil
}

// VerifySignature verifica la firma de un CAF
func (m *Manager) VerifySignature(caf *CAFXml, publicKey *rsa.PublicKey) error {
	// TODO: Implementar verificación de firma
	return nil
}

// SaveCAF guarda un CAF en un archivo XML
func (m *Manager) SaveCAF(caf *CAFXml, filename string) error {
	// Convertir a XML
	data, err := xml.MarshalIndent(caf, "", "  ")
	if err != nil {
		return fmt.Errorf("error al convertir CAF a XML: %w", err)
	}

	// Escribir archivo
	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return fmt.Errorf("error al escribir archivo CAF: %w", err)
	}

	return nil
}
