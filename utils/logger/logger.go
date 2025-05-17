package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	log *zap.Logger
)

// Config representa la configuración del logger
type Config struct {
	Level      string `json:"level"`
	OutputPath string `json:"output_path"`
	MaxSize    int    `json:"max_size"`    // megabytes
	MaxBackups int    `json:"max_backups"` // número de archivos
	MaxAge     int    `json:"max_age"`     // días
	Compress   bool   `json:"compress"`
}

// InitLogger inicializa el sistema de logging
func InitLogger(config *Config) error {
	if config == nil {
		config = &Config{
			Level:      "info",
			OutputPath: "logs/fmgo.log",
			MaxSize:    100,
			MaxBackups: 3,
			MaxAge:     28,
			Compress:   true,
		}
	}

	// Crear directorio de logs si no existe
	if err := os.MkdirAll(filepath.Dir(config.OutputPath), 0755); err != nil {
		return fmt.Errorf("error creando directorio de logs: %w", err)
	}

	// Configurar rotación de logs
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   config.OutputPath,
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
		Compress:   config.Compress,
	})

	// Configurar codificador
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	// Configurar nivel de logging
	var level zapcore.Level
	if err := level.UnmarshalText([]byte(config.Level)); err != nil {
		return fmt.Errorf("nivel de log inválido: %w", err)
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), w),
		level,
	)

	// Crear logger
	log = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	return nil
}

// Debug registra un mensaje de debug
func Debug(msg string, fields ...zap.Field) {
	if log != nil {
		log.Debug(msg, fields...)
	}
}

// Info registra un mensaje informativo
func Info(msg string, fields ...zap.Field) {
	if log != nil {
		log.Info(msg, fields...)
	}
}

// Warn registra un mensaje de advertencia
func Warn(msg string, fields ...zap.Field) {
	if log != nil {
		log.Warn(msg, fields...)
	}
}

// Error registra un mensaje de error
func Error(msg string, err error, fields ...zap.Field) {
	if log != nil {
		fields = append(fields, zap.Error(err))
		log.Error(msg, fields...)
	}
}

// Fatal registra un mensaje fatal y termina la aplicación
func Fatal(msg string, err error, fields ...zap.Field) {
	if log != nil {
		fields = append(fields, zap.Error(err))
		log.Fatal(msg, fields...)
	}
}

// WithFields agrega campos al contexto del logger
func WithFields(fields ...zap.Field) *zap.Logger {
	if log != nil {
		return log.With(fields...)
	}
	return nil
}

// Sync sincroniza los buffers del logger
func Sync() error {
	if log != nil {
		return log.Sync()
	}
	return nil
}

// GetLogger retorna la instancia del logger
func GetLogger() *zap.Logger {
	return log
}

// Field crea un campo de log
func Field(key string, value interface{}) zap.Field {
	return zap.Any(key, value)
}

// TraceID crea un campo de ID de traza
func TraceID(id string) zap.Field {
	return zap.String("trace_id", id)
}

// UserID crea un campo de ID de usuario
func UserID(id string) zap.Field {
	return zap.String("user_id", id)
}

// RequestID crea un campo de ID de solicitud
func RequestID(id string) zap.Field {
	return zap.String("request_id", id)
}

// Duration crea un campo de duración
func Duration(d time.Duration) zap.Field {
	return zap.Duration("duration", d)
}

// Caller crea un campo con información del llamador
func Caller() zap.Field {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		return zap.Skip()
	}
	return zap.String("caller", fmt.Sprintf("%s:%d", filepath.Base(file), line))
}
