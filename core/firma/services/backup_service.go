package services

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"FMgo/core/firma/models"
)

// BackupService maneja el respaldo de CAFs
type BackupService struct {
	backupDir     string
	metadataDir   string
	retentionDays int
}

// NewBackupService crea una nueva instancia del servicio de respaldo
func NewBackupService(backupDir, metadataDir string, retentionDays int) *BackupService {
	return &BackupService{
		backupDir:     backupDir,
		metadataDir:   metadataDir,
		retentionDays: retentionDays,
	}
}

// BackupCAF realiza una copia de respaldo de un CAF
func (s *BackupService) BackupCAF(caf *models.CAF) (*models.CAFBackup, error) {
	// Generar hash del contenido
	hash := sha256.Sum256([]byte(caf.XMLContenido))
	hashStr := hex.EncodeToString(hash[:])

	// Crear estructura de backup
	backup := &models.CAFBackup{
		CAF:           *caf,
		FechaBackup:   time.Now(),
		HashContenido: hashStr,
		Ubicacion:     filepath.Join(s.backupDir, fmt.Sprintf("%s_%s.json", caf.ID, hashStr[:8])),
	}

	// Serializar y guardar backup
	data, err := json.MarshalIndent(backup, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error serializando backup: %w", err)
	}

	if err := os.WriteFile(backup.Ubicacion, data, 0644); err != nil {
		return nil, fmt.Errorf("error guardando backup: %w", err)
	}

	return backup, nil
}

// RestoreCAF restaura un CAF desde su backup
func (s *BackupService) RestoreCAF(backupID string) (*models.CAF, error) {
	backupPath := filepath.Join(s.backupDir, backupID+".json")

	data, err := os.ReadFile(backupPath)
	if err != nil {
		return nil, fmt.Errorf("error leyendo backup: %w", err)
	}

	var backup models.CAFBackup
	if err := json.Unmarshal(data, &backup); err != nil {
		return nil, fmt.Errorf("error deserializando backup: %w", err)
	}

	// Verificar hash del contenido
	hash := sha256.Sum256([]byte(backup.CAF.XMLContenido))
	hashStr := hex.EncodeToString(hash[:])
	if hashStr != backup.HashContenido {
		return nil, fmt.Errorf("error de integridad: hash no coincide")
	}

	return &backup.CAF, nil
}

// CleanOldBackups elimina backups antiguos según la política de retención
func (s *BackupService) CleanOldBackups() error {
	cutoffDate := time.Now().AddDate(0, 0, -s.retentionDays)

	entries, err := os.ReadDir(s.backupDir)
	if err != nil {
		return fmt.Errorf("error leyendo directorio de backups: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		if info.ModTime().Before(cutoffDate) {
			if err := os.Remove(filepath.Join(s.backupDir, entry.Name())); err != nil {
				return fmt.Errorf("error eliminando backup antiguo %s: %w", entry.Name(), err)
			}
		}
	}

	return nil
}
