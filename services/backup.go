package services

import (
	"archive/zip"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gorm.io/gorm"
)

// BackupService maneja los respaldos de la base de datos
type BackupService struct {
	db        *gorm.DB
	backupDir string
	retention int // días de retención
}

// NewBackupService crea una nueva instancia del servicio de respaldo
func NewBackupService(db *gorm.DB, backupDir string, retention int) *BackupService {
	return &BackupService{
		db:        db,
		backupDir: backupDir,
		retention: retention,
	}
}

// CreateBackup crea un nuevo respaldo
func (s *BackupService) CreateBackup() error {
	// Crear directorio de respaldo si no existe
	if err := os.MkdirAll(s.backupDir, 0755); err != nil {
		return fmt.Errorf("error al crear directorio de respaldo: %w", err)
	}

	// Generar nombre de archivo
	timestamp := time.Now().Format("20060102150405")
	backupFile := filepath.Join(s.backupDir, fmt.Sprintf("backup_%s.zip", timestamp))

	// Crear archivo ZIP
	zipFile, err := os.Create(backupFile)
	if err != nil {
		return fmt.Errorf("error al crear archivo ZIP: %w", err)
	}
	defer zipFile.Close()

	// Crear writer ZIP
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Exportar datos de la base de datos
	if err := s.exportData(zipWriter); err != nil {
		return fmt.Errorf("error al exportar datos: %w", err)
	}

	// Limpiar respaldos antiguos
	if err := s.cleanOldBackups(); err != nil {
		return fmt.Errorf("error al limpiar respaldos antiguos: %w", err)
	}

	return nil
}

// exportData exporta los datos de la base de datos
func (s *BackupService) exportData(zipWriter *zip.Writer) error {
	// Exportar boletas
	if err := s.exportTable(zipWriter, "boletas", "SELECT * FROM boletas"); err != nil {
		return err
	}

	// Exportar detalles de boletas
	if err := s.exportTable(zipWriter, "detalles_boletas", "SELECT * FROM detalles_boletas"); err != nil {
		return err
	}

	// Exportar logs de auditoría
	if err := s.exportTable(zipWriter, "audit_logs", "SELECT * FROM audit_logs"); err != nil {
		return err
	}

	return nil
}

// exportTable exporta una tabla específica
func (s *BackupService) exportTable(zipWriter *zip.Writer, tableName, query string) error {
	// Crear archivo en el ZIP
	file, err := zipWriter.Create(fmt.Sprintf("%s.csv", tableName))
	if err != nil {
		return fmt.Errorf("error al crear archivo en ZIP: %w", err)
	}

	// Ejecutar consulta y exportar a CSV
	rows, err := s.db.Raw(query).Rows()
	if err != nil {
		return fmt.Errorf("error al ejecutar consulta: %w", err)
	}
	defer rows.Close()

	// Escribir datos en CSV
	columns, err := rows.Columns()
	if err != nil {
		return fmt.Errorf("error al obtener columnas: %w", err)
	}

	// Escribir encabezados
	for i, col := range columns {
		if i > 0 {
			file.Write([]byte(","))
		}
		file.Write([]byte(col))
	}
	file.Write([]byte("\n"))

	// Escribir datos
	values := make([]interface{}, len(columns))
	for i := range values {
		values[i] = new(interface{})
	}

	for rows.Next() {
		if err := rows.Scan(values...); err != nil {
			return fmt.Errorf("error al escanear fila: %w", err)
		}

		for i, val := range values {
			if i > 0 {
				file.Write([]byte(","))
			}
			if val != nil {
				file.Write([]byte(fmt.Sprintf("%v", *val.(*interface{}))))
			}
		}
		file.Write([]byte("\n"))
	}

	return nil
}

// cleanOldBackups elimina respaldos antiguos
func (s *BackupService) cleanOldBackups() error {
	// Obtener lista de archivos de respaldo
	files, err := filepath.Glob(filepath.Join(s.backupDir, "backup_*.zip"))
	if err != nil {
		return fmt.Errorf("error al listar archivos de respaldo: %w", err)
	}

	// Calcular fecha límite
	limit := time.Now().AddDate(0, 0, -s.retention)

	// Eliminar archivos antiguos
	for _, file := range files {
		info, err := os.Stat(file)
		if err != nil {
			continue
		}

		if info.ModTime().Before(limit) {
			if err := os.Remove(file); err != nil {
				return fmt.Errorf("error al eliminar archivo antiguo: %w", err)
			}
		}
	}

	return nil
}

// RestoreBackup restaura un respaldo
func (s *BackupService) RestoreBackup(backupFile string) error {
	// Abrir archivo ZIP
	reader, err := zip.OpenReader(backupFile)
	if err != nil {
		return fmt.Errorf("error al abrir archivo ZIP: %w", err)
	}
	defer reader.Close()

	// Restaurar cada tabla
	for _, file := range reader.File {
		if err := s.restoreTable(file); err != nil {
			return fmt.Errorf("error al restaurar tabla: %w", err)
		}
	}

	return nil
}

// restoreTable restaura una tabla desde un archivo CSV
func (s *BackupService) restoreTable(file *zip.File) error {
	// Abrir archivo
	reader, err := file.Open()
	if err != nil {
		return fmt.Errorf("error al abrir archivo: %w", err)
	}
	defer reader.Close()

	// TODO: Implementar restauración de datos
	// Esto requerirá un parser CSV y la inserción de datos en la base de datos

	return nil
}
