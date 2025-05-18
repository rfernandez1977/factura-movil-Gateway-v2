package services

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

// LogLevel define los niveles de log
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

// LogService proporciona funcionalidades de logging
type LogService struct {
	logger *log.Logger
	level  LogLevel
	file   *os.File
}

// NewLogService crea una nueva instancia del servicio de logging
func NewLogService(logPath string) (*LogService, error) {
	// Crear directorio si no existe
	dir := filepath.Dir(logPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("error creando directorio de logs: %w", err)
	}

	// Abrir archivo de log
	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("error abriendo archivo de log: %w", err)
	}

	return &LogService{
		logger: log.New(file, "", log.LstdFlags),
		level:  DEBUG,
		file:   file,
	}, nil
}

// SetLevel establece el nivel de logging
func (s *LogService) SetLevel(level LogLevel) {
	s.level = level
}

// Debug registra un mensaje de debug
func (s *LogService) Debug(format string, v ...interface{}) {
	if s.level <= DEBUG {
		s.log("DEBUG", format, v...)
	}
}

// Info registra un mensaje informativo
func (s *LogService) Info(format string, v ...interface{}) {
	if s.level <= INFO {
		s.log("INFO", format, v...)
	}
}

// Warn registra una advertencia
func (s *LogService) Warn(format string, v ...interface{}) {
	if s.level <= WARN {
		s.log("WARN", format, v...)
	}
}

// Error registra un error
func (s *LogService) Error(format string, v ...interface{}) {
	if s.level <= ERROR {
		s.log("ERROR", format, v...)
	}
}

// Close cierra el archivo de log
func (s *LogService) Close() error {
	if s.file != nil {
		return s.file.Close()
	}
	return nil
}

// log registra un mensaje con el nivel especificado
func (s *LogService) log(level string, format string, v ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05.000")
	message := fmt.Sprintf(format, v...)
	s.logger.Printf("[%s] %s - %s", level, timestamp, message)
}

// LogXML registra un documento XML con formato
func (s *LogService) LogXML(operation string, xmlData []byte) error {
	s.Info("Operación: %s", operation)
	s.Debug("XML:\n%s", string(xmlData))
	return nil
}

// LogError registra un error con detalles
func (s *LogService) LogError(operation string, err error) {
	s.Error("Error en %s: %v", operation, err)
}

// LogValidacion registra el resultado de una validación
func (s *LogService) LogValidacion(xmlData []byte, resultado string) error {
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("validacion_%s.log", timestamp)
	filepath := filepath.Join(s.logDir, filename)

	logContent := fmt.Sprintf("Timestamp: %s\nResultado: %s\nXML:\n%s\n",
		timestamp, resultado, string(xmlData))

	if err := os.WriteFile(filepath, []byte(logContent), 0644); err != nil {
		return fmt.Errorf("error guardando log de validación: %w", err)
	}

	return nil
}

// LimpiarLogs elimina logs antiguos
func (s *LogService) LimpiarLogs(diasRetencion int) error {
	entries, err := os.ReadDir(s.logDir)
	if err != nil {
		return fmt.Errorf("error leyendo directorio de logs: %w", err)
	}

	limite := time.Now().AddDate(0, 0, -diasRetencion)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		if info.ModTime().Before(limite) {
			if err := os.Remove(filepath.Join(s.logDir, entry.Name())); err != nil {
				return fmt.Errorf("error eliminando log antiguo %s: %w", entry.Name(), err)
			}
		}
	}

	return nil
}
