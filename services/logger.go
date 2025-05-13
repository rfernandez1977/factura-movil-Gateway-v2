package services

// Logger define una interfaz para logging
type Logger interface {
	Log(format string, args ...interface{})
	Info(format string, args ...interface{})
	Warn(format string, args ...interface{})
	Error(format string, args ...interface{})
}

// DefaultLogger is the default implementation of Logger
type DefaultLogger struct {
	LogPath string
}

func NewLogger(logPath string) *DefaultLogger {
	return &DefaultLogger{
		LogPath: logPath,
	}
}

// Log logs a message
func (l *DefaultLogger) Log(format string, args ...interface{}) {
	// Implementation would go here
}

// Info logs an informational message
func (l *DefaultLogger) Info(format string, args ...interface{}) {
	// Implementation would go here
}

// Warn logs a warning message
func (l *DefaultLogger) Warn(format string, args ...interface{}) {
	// Implementation would go here
}

// Error logs an error message
func (l *DefaultLogger) Error(format string, args ...interface{}) {
	// Implementation would go here
}

// NoOpLogger implementa Logger sin hacer nada
type NoOpLogger struct{}

func NewNoOpLogger() *NoOpLogger {
	return &NoOpLogger{}
}

func (l *NoOpLogger) Log(format string, args ...interface{})   {}
func (l *NoOpLogger) Info(format string, args ...interface{})  {}
func (l *NoOpLogger) Warn(format string, args ...interface{})  {}
func (l *NoOpLogger) Error(format string, args ...interface{}) {}
