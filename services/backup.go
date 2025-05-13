package services

import (
	"archive/zip"
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
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

// GetBackupFiles retorna la lista de archivos de respaldo
func (s *BackupService) GetBackupFiles() ([]string, error) {
	// Crear directorio de respaldo si no existe
	if err := os.MkdirAll(s.backupDir, 0755); err != nil {
		return nil, fmt.Errorf("error al crear directorio de respaldo: %w", err)
	}

	// Obtener lista de archivos de respaldo
	files, err := filepath.Glob(filepath.Join(s.backupDir, "backup_*.zip"))
	if err != nil {
		return nil, fmt.Errorf("error al listar archivos de respaldo: %w", err)
	}

	return files, nil
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

	// Obtener columnas
	columns, err := rows.Columns()
	if err != nil {
		return fmt.Errorf("error al obtener columnas: %w", err)
	}

	// Escribir encabezados
	header := strings.Join(columns, ",")
	if _, err := file.Write([]byte(header + "\n")); err != nil {
		return fmt.Errorf("error al escribir encabezados: %w", err)
	}

	// Preparar contenedores para datos
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	// Procesar filas
	for rows.Next() {
		// Escanear fila actual
		if err := rows.Scan(scanArgs...); err != nil {
			return fmt.Errorf("error al escanear fila: %w", err)
		}

		// Procesar valores
		var rowValues []string
		for _, col := range values {
			// Verificar si es NULL
			if col == nil {
				rowValues = append(rowValues, "")
			} else {
				// Escapar comillas y comas para CSV
				value := string(col)
				if strings.ContainsAny(value, "\",") {
					value = fmt.Sprintf("\"%s\"", strings.ReplaceAll(value, "\"", "\"\""))
				}
				rowValues = append(rowValues, value)
			}
		}

		// Escribir línea
		line := strings.Join(rowValues, ",")
		if _, err := file.Write([]byte(line + "\n")); err != nil {
			return fmt.Errorf("error al escribir línea: %w", err)
		}
	}

	// Verificar errores durante la iteración
	if err := rows.Err(); err != nil {
		return fmt.Errorf("error durante la iteración: %w", err)
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
	// Extraer nombre de tabla del nombre de archivo
	tableName := filepath.Base(file.Name)
	if len(tableName) < 5 || tableName[len(tableName)-4:] != ".csv" {
		return fmt.Errorf("formato de archivo inválido: %s", tableName)
	}
	tableName = tableName[:len(tableName)-4]

	// Abrir archivo
	reader, err := file.Open()
	if err != nil {
		return fmt.Errorf("error al abrir archivo: %w", err)
	}
	defer reader.Close()

	// Leer archivo CSV
	csvData, err := readCSVData(reader)
	if err != nil {
		return fmt.Errorf("error al leer datos CSV: %w", err)
	}

	// Si no hay datos, no hay nada que restaurar
	if len(csvData) < 2 {
		return nil
	}

	// Obtener nombres de columnas (primera línea)
	columns := csvData[0]

	// Iniciar transacción
	tx := s.db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("error al iniciar transacción: %w", tx.Error)
	}

	// Limpiar tabla existente
	if err := tx.Exec(fmt.Sprintf("DELETE FROM %s", tableName)).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("error al limpiar tabla: %w", err)
	}

	// Insertar filas
	for i := 1; i < len(csvData); i++ {
		row := csvData[i]

		// Construir mapa para inserción
		data := make(map[string]interface{})
		for j, col := range columns {
			if j < len(row) {
				data[col] = row[j]
			} else {
				data[col] = nil
			}
		}

		// Insertar fila
		if err := tx.Table(tableName).Create(data).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("error al insertar fila %d: %w", i, err)
		}
	}

	// Confirmar transacción
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("error al confirmar transacción: %w", err)
	}

	return nil
}

// readCSVData lee datos de un archivo CSV
func readCSVData(reader io.Reader) ([][]string, error) {
	// Leer todo el contenido
	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("error al leer contenido: %w", err)
	}

	// Dividir en líneas
	lines := strings.Split(string(content), "\n")

	// Procesar cada línea
	var result [][]string
	for _, line := range lines {
		if line == "" {
			continue
		}

		// Dividir por comas (esto es simplificado, un parser CSV real maneja comillas, escapes, etc.)
		fields := strings.Split(line, ",")
		result = append(result, fields)
	}

	return result, nil
}
