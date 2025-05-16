package filesystem

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/fmgo/core/firma/common"
	"github.com/fmgo/core/firma/models"
)

// CAFStorage implementa el almacenamiento de CAF en sistema de archivos
type CAFStorage struct {
	basePath string
	mu       sync.RWMutex
	logger   common.Logger
}

// NewCAFStorage crea una nueva instancia de almacenamiento de CAF
func NewCAFStorage(basePath string, logger common.Logger) (*CAFStorage, error) {
	// Crear directorio base si no existe
	if err := os.MkdirAll(basePath, 0750); err != nil {
		return nil, fmt.Errorf("error creando directorio base: %w", err)
	}

	return &CAFStorage{
		basePath: basePath,
		logger:   logger,
	}, nil
}

// GuardarCAF guarda un CAF en el sistema de archivos
func (s *CAFStorage) GuardarCAF(ctx context.Context, caf *models.CAF) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Crear directorio para el CAF
	cafPath := filepath.Join(s.basePath, caf.ID)
	if err := os.MkdirAll(cafPath, 0750); err != nil {
		return fmt.Errorf("error creando directorio del CAF: %w", err)
	}

	// Guardar metadatos
	metadataPath := filepath.Join(cafPath, "metadata.json")
	metadata, err := json.MarshalIndent(caf, "", "  ")
	if err != nil {
		return fmt.Errorf("error serializando metadatos: %w", err)
	}

	if err := os.WriteFile(metadataPath, metadata, 0600); err != nil {
		return fmt.Errorf("error guardando metadatos: %w", err)
	}

	// Guardar XML original
	xmlPath := filepath.Join(cafPath, "caf.xml")
	if err := os.WriteFile(xmlPath, caf.XML, 0600); err != nil {
		return fmt.Errorf("error guardando XML: %w", err)
	}

	// Guardar firma SII
	firmaPath := filepath.Join(cafPath, "firma_sii.bin")
	if err := os.WriteFile(firmaPath, caf.FirmaSII, 0600); err != nil {
		return fmt.Errorf("error guardando firma SII: %w", err)
	}

	s.logger.Info("CAF guardado exitosamente",
		"id", caf.ID,
		"tipo", caf.TipoDocumento,
		"rango", fmt.Sprintf("%d-%d", caf.FolioInicial, caf.FolioFinal))

	return nil
}

// ObtenerCAF obtiene un CAF del sistema de archivos
func (s *CAFStorage) ObtenerCAF(ctx context.Context, id string) (*models.CAF, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	cafPath := filepath.Join(s.basePath, id)

	// Leer metadatos
	metadataPath := filepath.Join(cafPath, "metadata.json")
	metadata, err := os.ReadFile(metadataPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("CAF no encontrado: %s", id)
		}
		return nil, fmt.Errorf("error leyendo metadatos: %w", err)
	}

	var caf models.CAF
	if err := json.Unmarshal(metadata, &caf); err != nil {
		return nil, fmt.Errorf("error deserializando metadatos: %w", err)
	}

	// Leer XML original
	xmlPath := filepath.Join(cafPath, "caf.xml")
	xml, err := os.ReadFile(xmlPath)
	if err != nil {
		return nil, fmt.Errorf("error leyendo XML: %w", err)
	}
	caf.XML = xml

	// Leer firma SII
	firmaPath := filepath.Join(cafPath, "firma_sii.bin")
	firma, err := os.ReadFile(firmaPath)
	if err != nil {
		return nil, fmt.Errorf("error leyendo firma SII: %w", err)
	}
	caf.FirmaSII = firma

	return &caf, nil
}

// ObtenerCAFPorFolio obtiene un CAF que contiene el folio especificado
func (s *CAFStorage) ObtenerCAFPorFolio(ctx context.Context, tipo string, folio int64) (*models.CAF, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Listar todos los CAFs del tipo especificado
	cafs, err := s.ListarCAFsPorTipo(ctx, tipo)
	if err != nil {
		return nil, err
	}

	// Buscar el CAF que contiene el folio
	for _, caf := range cafs {
		if caf.ContieneFolio(folio) && caf.EstaVigente() {
			return caf, nil
		}
	}

	return nil, fmt.Errorf("no se encontró CAF vigente para folio %d tipo %s", folio, tipo)
}

// ListarCAFsPorTipo lista todos los CAFs de un tipo de documento
func (s *CAFStorage) ListarCAFsPorTipo(ctx context.Context, tipo string) ([]*models.CAF, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entries, err := os.ReadDir(s.basePath)
	if err != nil {
		return nil, fmt.Errorf("error listando directorio: %w", err)
	}

	var cafs []*models.CAF
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		caf, err := s.ObtenerCAF(ctx, entry.Name())
		if err != nil {
			s.logger.Warn("Error obteniendo CAF",
				"id", entry.Name(),
				"error", err)
			continue
		}

		if caf.TipoDocumento == tipo {
			cafs = append(cafs, caf)
		}
	}

	return cafs, nil
}

// ActualizarEstadoCAF actualiza el estado de un CAF
func (s *CAFStorage) ActualizarEstadoCAF(ctx context.Context, id string, estado string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Obtener CAF actual
	caf, err := s.ObtenerCAF(ctx, id)
	if err != nil {
		return err
	}

	// Actualizar estado
	caf.Estado = estado

	// Guardar cambios
	return s.GuardarCAF(ctx, caf)
}
