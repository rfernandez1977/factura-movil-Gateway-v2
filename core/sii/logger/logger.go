package logger

import (
	"fmt"
	"log"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// LogLevel representa el nivel de logging
type LogLevel string

const (
	DEBUG LogLevel = "DEBUG"
	INFO  LogLevel = "INFO"
	WARN  LogLevel = "WARN"
	ERROR LogLevel = "ERROR"
)

// Logger es un wrapper para zap.Logger con funcionalidad adicional
type Logger struct {
	*zap.Logger
	level  LogLevel
	logger *log.Logger
}

// New crea una nueva instancia del logger
func New(development bool) (*Logger, error) {
	var config zap.Config
	if development {
		config = zap.NewDevelopmentConfig()
	} else {
		config = zap.NewProductionConfig()
	}

	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	zapLogger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return &Logger{Logger: zapLogger}, nil
}

// NewLogger crea una nueva instancia del logger
func NewLogger(name string, level LogLevel) (*Logger, error) {
	logger := log.New(os.Stdout, fmt.Sprintf("[%s] ", name), log.LstdFlags)
	return &Logger{
		level:  level,
		logger: logger,
	}, nil
}

// Debug registra mensajes de nivel debug
func (l *Logger) Debug(format string, v ...interface{}) {
	if l.level == DEBUG {
		l.logger.Printf("[DEBUG] "+format, v...)
	}
}

// Info registra mensajes de nivel info
func (l *Logger) Info(format string, v ...interface{}) {
	if l.level <= INFO {
		l.logger.Printf("[INFO] "+format, v...)
	}
}

// Warn registra mensajes de nivel warning
func (l *Logger) Warn(format string, v ...interface{}) {
	if l.level <= WARN {
		l.logger.Printf("[WARN] "+format, v...)
	}
}

// Error registra mensajes de nivel error
func (l *Logger) Error(format string, v ...interface{}) {
	l.logger.Printf("[ERROR] "+format, v...)
}

// LogXMLOperation registra operaciones relacionadas con XML
func (l *Logger) LogXMLOperation(operation string, xmlData []byte, err error) {
	if l.level > DEBUG {
		return
	}

	// Truncar XML largo para el log
	xmlStr := string(xmlData)
	if len(xmlStr) > 1000 {
		xmlStr = xmlStr[:1000] + "... [truncado]"
	}

	if err != nil {
		l.Error("Operación XML '%s' falló: %v\nXML:\n%s", operation, err, xmlStr)
	} else {
		l.Debug("Operación XML '%s' exitosa\nXML:\n%s", operation, xmlStr)
	}
}

// LogCertOperation registra operaciones relacionadas con certificados
func (l *Logger) LogCertOperation(operation string, certInfo string, err error) {
	if err != nil {
		l.Error("Operación de certificado '%s' falló: %v\nInfo: %s", operation, err, certInfo)
	} else {
		l.Info("Operación de certificado '%s' exitosa\nInfo: %s", operation, certInfo)
	}
}

// Close cierra el logger y libera recursos
func (l *Logger) Close() error {
	return l.Sync()
}
