package logger

import (
	"FMgo/utils/logger"
)

// Logger define la interfaz que debe implementar cualquier logger usado por el cliente SII
type Logger interface {
	Debug(format string, v ...interface{})
	Info(format string, v ...interface{})
	Warn(format string, v ...interface{})
	Error(format string, v ...interface{})
	Fatal(format string, v ...interface{})
}

// NewLogger crea una nueva instancia del logger por defecto
func NewLogger() Logger {
	return logger.NewLogger()
}
