package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

// Level representa el nivel de logging
type Level int

const (
	// Niveles de log
	DEBUG Level = iota
	INFO
	WARN
	ERROR
	FATAL
)

// Logger representa un logger personalizado
type Logger struct {
	debug   *log.Logger
	info    *log.Logger
	warn    *log.Logger
	error   *log.Logger
	level   Level
	logFile *os.File
}

// Config representa la configuración del logger
type Config struct {
	Level      Level
	OutputPath string
	ToConsole  bool
	ToFile     bool
}

// DefaultConfig retorna la configuración por defecto
func DefaultConfig() *Config {
	return &Config{
		Level:      INFO,
		OutputPath: "logs/app.log",
		ToConsole:  true,
		ToFile:     true,
	}
}

// NewLogger crea una nueva instancia del logger
func NewLogger() *Logger {
	return NewLoggerWithConfig(DefaultConfig())
}

// NewLoggerWithConfig crea una nueva instancia del logger con configuración personalizada
func NewLoggerWithConfig(config *Config) *Logger {
	var outputs []string
	writers := []interface{}{os.Stdout}

	if config.ToFile {
		// Crear directorio de logs si no existe
		if err := os.MkdirAll(filepath.Dir(config.OutputPath), 0755); err != nil {
			fmt.Printf("Error creando directorio de logs: %v\n", err)
		}

		// Abrir archivo de log
		f, err := os.OpenFile(config.OutputPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			fmt.Printf("Error abriendo archivo de log: %v\n", err)
		} else {
			writers = append(writers, f)
			outputs = append(outputs, config.OutputPath)
		}
	}

	if config.ToConsole {
		outputs = append(outputs, "console")
	}

	flags := log.Ldate | log.Ltime | log.Lshortfile

	return &Logger{
		debug: log.New(os.Stdout, "[DEBUG] ", flags),
		info:  log.New(os.Stdout, "[INFO] ", flags),
		warn:  log.New(os.Stdout, "[WARN] ", flags),
		error: log.New(os.Stderr, "[ERROR] ", flags),
		level: config.Level,
	}
}

// Debug registra un mensaje de depuración
func (l *Logger) Debug(format string, v ...interface{}) {
	if l.level <= DEBUG {
		l.debug.Output(2, fmt.Sprintf(format, v...))
	}
}

// Info registra un mensaje informativo
func (l *Logger) Info(format string, v ...interface{}) {
	if l.level <= INFO {
		l.info.Output(2, fmt.Sprintf(format, v...))
	}
}

// Warn registra un mensaje de advertencia
func (l *Logger) Warn(format string, v ...interface{}) {
	if l.level <= WARN {
		l.warn.Output(2, fmt.Sprintf(format, v...))
	}
}

// Error registra un mensaje de error
func (l *Logger) Error(format string, v ...interface{}) {
	if l.level <= ERROR {
		l.error.Output(2, fmt.Sprintf(format, v...))
	}
}

// Fatal registra un mensaje de error y termina el programa
func (l *Logger) Fatal(format string, v ...interface{}) {
	l.error.Output(2, fmt.Sprintf(format, v...))
	os.Exit(1)
}

// WithFields agrega campos al mensaje de log
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	return l
}

// Close cierra el archivo de log si existe
func (l *Logger) Close() error {
	if l.logFile != nil {
		return l.logFile.Close()
	}
	return nil
}

// SetLevel establece el nivel de logging
func (l *Logger) SetLevel(level Level) {
	l.level = level
}

// GetLevel obtiene el nivel actual de logging
func (l *Logger) GetLevel() Level {
	return l.level
}

// FormatTime formatea un tiempo para el log
func (l *Logger) FormatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05.000")
}
