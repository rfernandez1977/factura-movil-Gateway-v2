package caf

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/fmgo/config"
	"github.com/fmgo/models"
)

// Service representa el servicio de CAF
type Service struct {
	config        *config.Config
	redis         interface{}
	sii           interface{}
	backupDir     string
	retentionDays int
}

// NewService crea una nueva instancia del servicio de CAF
func NewService(config *config.Config, redis interface{}, sii interface{}, backupDir string, retentionDays int) *Service {
	return &Service{
		config:        config,
		redis:         redis,
		sii:           sii,
		backupDir:     backupDir,
		retentionDays: retentionDays,
	}
}

// BackupResult representa el resultado de un respaldo
type BackupResult struct {
	CAFId         string    `json:"caf_id"`
	FechaBackup   time.Time `json:"fecha_backup"`
	HashContenido string    `json:"hash_contenido"`
	Ubicacion     string    `json:"ubicacion"`
}

// BackupCAF realiza una copia de respaldo de un CAF
func (s *Service) BackupCAF(ctx context.Context, caf *models.CAF) (*BackupResult, error) {
	// Generar hash del contenido
	hash := sha256.Sum256(caf.Archivo)
	hashStr := hex.EncodeToString(hash[:])

	// Crear nombre de archivo
	backupFile := fmt.Sprintf("%s_%s.json", caf.ID, hashStr[:8])
	backupPath := filepath.Join(s.backupDir, backupFile)

	// Crear estructura de backup
	backup := &BackupResult{
		CAFId:         caf.ID,
		FechaBackup:   time.Now(),
		HashContenido: hashStr,
		Ubicacion:     backupPath,
	}

	// Serializar y guardar backup
	data, err := json.MarshalIndent(backup, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error serializando backup: %w", err)
	}

	if err := os.WriteFile(backupPath, data, 0644); err != nil {
		return nil, fmt.Errorf("error guardando backup: %w", err)
	}

	return backup, nil
}

// RestoreCAF restaura un CAF desde su backup
func (s *Service) RestoreCAF(ctx context.Context, backupID string) (*models.CAF, error) {
	backupPath := filepath.Join(s.backupDir, backupID+".json")

	data, err := os.ReadFile(backupPath)
	if err != nil {
		return nil, fmt.Errorf("error leyendo backup: %w", err)
	}

	var backup BackupResult
	if err := json.Unmarshal(data, &backup); err != nil {
		return nil, fmt.Errorf("error deserializando backup: %w", err)
	}

	// Obtener CAF original
	caf, err := s.GetCAFByID(ctx, backup.CAFId)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo CAF: %w", err)
	}

	// Verificar hash del contenido
	hash := sha256.Sum256(caf.Archivo)
	hashStr := hex.EncodeToString(hash[:])
	if hashStr != backup.HashContenido {
		return nil, fmt.Errorf("error de integridad: hash no coincide")
	}

	return caf, nil
}

// CleanOldBackups elimina backups antiguos según la política de retención
func (s *Service) CleanOldBackups(ctx context.Context) error {
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

// GetCAFByID obtiene un CAF por su ID
func (s *Service) GetCAFByID(ctx context.Context, id string) (*models.CAF, error) {
	// TODO: Implementar obtención de CAF desde la base de datos
	return nil, fmt.Errorf("método no implementado")
}

// GetCAFDisponible obtiene un CAF disponible
func (s *Service) GetCAFDisponible(ctx context.Context, tipoDTE models.TipoDTE, rutEmisor string) (*models.CAFDTEXML, error) {
	// Implementación mock
	return &models.CAFDTEXML{
		Version: "1.0",
		DA: models.DAXMLModel{
			RUT: models.RutXMLModel{
				Numero: rutEmisor,
			},
			RazonSocial: "EMPRESA DE PRUEBA",
			TipoDTE:     string(tipoDTE),
			RangoDesde:  1,
			RangoHasta:  100,
			FechaAut:    "2023-01-01",
			RSAPK: models.RSAPKXMLModel{
				Modulo:    "test-modulus",
				Exponente: "test-exponent",
			},
			IDK: 1,
		},
		FRMA: models.FRMAXMLModel{
			Algoritmo: "SHA1withRSA",
			Valor:     "test-signature",
		},
	}, nil
}

// ValidarCAF valida un CAF
func (s *Service) ValidarCAF(caf *models.CAFDTEXML) error {
	// Implementación mock
	return nil
}
