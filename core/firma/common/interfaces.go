package common

// Logger define las operaciones del logger
type Logger interface {
	// Debug registra un mensaje de debug
	Debug(msg string, args ...interface{})
	// Info registra un mensaje informativo
	Info(msg string, args ...interface{})
	// Warn registra una advertencia
	Warn(msg string, args ...interface{})
	// Error registra un error
	Error(msg string, args ...interface{})
}
