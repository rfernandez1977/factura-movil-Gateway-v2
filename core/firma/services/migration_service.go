package services

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"fmgo/core/firma/models"
)

// MigrationService maneja la migración de datos de CAF
type MigrationService struct {
	sourceDir     string
	targetDir     string
	backupService *BackupService
}

// NewMigrationService crea una nueva instancia del servicio de migración
func NewMigrationService(sourceDir, targetDir string, backupService *BackupService) *MigrationService {
	return &MigrationService{
		sourceDir:     sourceDir,
		targetDir:     targetDir,
		backupService: backupService,
	}
}

// MigrationResult contiene el resultado de una migración
type MigrationResult struct {
	TotalCAFs   int       `json:"total_cafs"`
	MigradosOK  int       `json:"migrados_ok"`
	Errores     []string  `json:"errores"`
	FechaInicio time.Time `json:"fecha_inicio"`
	FechaFin    time.Time `json:"fecha_fin"`
}

// MigrateCAFs realiza la migración de CAFs desde el directorio fuente
func (s *MigrationService) MigrateCAFs() (*MigrationResult, error) {
	result := &MigrationResult{
		FechaInicio: time.Now(),
		Errores:     make([]string, 0),
	}

	// Leer archivos del directorio fuente
	entries, err := os.ReadDir(s.sourceDir)
	if err != nil {
		return nil, fmt.Errorf("error leyendo directorio fuente: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".xml" {
			continue
		}

		result.TotalCAFs++

		// Leer y procesar archivo CAF
		cafPath := filepath.Join(s.sourceDir, entry.Name())
		xmlData, err := os.ReadFile(cafPath)
		if err != nil {
			result.Errores = append(result.Errores,
				fmt.Sprintf("Error leyendo %s: %v", entry.Name(), err))
			continue
		}

		// Crear estructura CAF
		caf := &models.CAF{
			ID:            filepath.Base(entry.Name()),
			XMLContenido:  string(xmlData),
			FechaCreacion: time.Now(),
			Estado:        "ACTIVO",
		}

		// Realizar backup del CAF
		if _, err := s.backupService.BackupCAF(caf); err != nil {
			result.Errores = append(result.Errores,
				fmt.Sprintf("Error respaldando %s: %v", entry.Name(), err))
			continue
		}

		// Guardar CAF en nuevo formato
		targetPath := filepath.Join(s.targetDir, entry.Name())
		if err := s.saveCAF(caf, targetPath); err != nil {
			result.Errores = append(result.Errores,
				fmt.Sprintf("Error guardando %s: %v", entry.Name(), err))
			continue
		}

		result.MigradosOK++
	}

	result.FechaFin = time.Now()
	return result, nil
}

// saveCAF guarda un CAF en el nuevo formato
func (s *MigrationService) saveCAF(caf *models.CAF, targetPath string) error {
	data, err := json.MarshalIndent(caf, "", "  ")
	if err != nil {
		return fmt.Errorf("error serializando CAF: %w", err)
	}

	if err := os.WriteFile(targetPath, data, 0644); err != nil {
		return fmt.Errorf("error guardando CAF: %w", err)
	}

	return nil
}

// ValidateMigration valida que la migración se haya realizado correctamente
func (s *MigrationService) ValidateMigration() error {
	sourceEntries, err := os.ReadDir(s.sourceDir)
	if err != nil {
		return fmt.Errorf("error leyendo directorio fuente: %w", err)
	}

	for _, entry := range sourceEntries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".xml" {
			continue
		}

		targetPath := filepath.Join(s.targetDir, entry.Name())
		if _, err := os.Stat(targetPath); os.IsNotExist(err) {
			return fmt.Errorf("archivo migrado no encontrado: %s", entry.Name())
		}
	}

	return nil
}
