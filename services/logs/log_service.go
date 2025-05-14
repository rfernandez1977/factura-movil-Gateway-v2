package logs

import (
	"fmt"
	"time"

	"github.com/cursor/FMgo/models"
)

// LogLevel define el nivel de log
type LogLevel string

const (
	LogLevelDebug   LogLevel = "DEBUG"
	LogLevelInfo    LogLevel = "INFO"
	LogLevelWarning LogLevel = "WARNING"
	LogLevelError   LogLevel = "ERROR"
	LogLevelFatal   LogLevel = "FATAL"
)

// LogEntry representa una entrada de log
type LogEntry struct {
	ID        string                 `json:"id"`
	Timestamp time.Time              `json:"timestamp"`
	Level     LogLevel               `json:"level"`
	Source    string                 `json:"source"`
	Message   string                 `json:"message"`
	Data      map[string]interface{} `json:"data,omitempty"`
}

// LogService implementa el servicio de logs
type LogService struct {
	config *models.Config
	logs   []LogEntry
}

// NewLogService crea una nueva instancia de LogService
func NewLogService(config *models.Config) *LogService {
	return &LogService{
		config: config,
		logs:   make([]LogEntry, 0),
	}
}

// GetLogs obtiene todos los logs
func (s *LogService) GetLogs() ([]LogEntry, error) {
	return s.logs, nil
}

// GetLog obtiene un log por su ID
func (s *LogService) GetLog(id string) (*LogEntry, error) {
	for _, log := range s.logs {
		if log.ID == id {
			return &log, nil
		}
	}
	return nil, fmt.Errorf("log no encontrado: %s", id)
}

// AddLog agrega un log
func (s *LogService) AddLog(level LogLevel, source, message string, data map[string]interface{}) string {
	log := LogEntry{
		ID:        models.GenerateID(),
		Timestamp: time.Now(),
		Level:     level,
		Source:    source,
		Message:   message,
		Data:      data,
	}
	s.logs = append(s.logs, log)
	return log.ID
}

// Debug agrega un log de nivel debug
func (s *LogService) Debug(source, message string, data map[string]interface{}) string {
	if s.config.LogLevel == "debug" {
		return s.AddLog(LogLevelDebug, source, message, data)
	}
	return ""
}

// Info agrega un log de nivel info
func (s *LogService) Info(source, message string, data map[string]interface{}) string {
	if s.config.LogLevel == "debug" || s.config.LogLevel == "info" {
		return s.AddLog(LogLevelInfo, source, message, data)
	}
	return ""
}

// Warning agrega un log de nivel warning
func (s *LogService) Warning(source, message string, data map[string]interface{}) string {
	if s.config.LogLevel != "error" && s.config.LogLevel != "fatal" {
		return s.AddLog(LogLevelWarning, source, message, data)
	}
	return ""
}

// Error agrega un log de nivel error
func (s *LogService) Error(source, message string, data map[string]interface{}) string {
	if s.config.LogLevel != "fatal" {
		return s.AddLog(LogLevelError, source, message, data)
	}
	return ""
}

// Fatal agrega un log de nivel fatal
func (s *LogService) Fatal(source, message string, data map[string]interface{}) string {
	return s.AddLog(LogLevelFatal, source, message, data)
}
