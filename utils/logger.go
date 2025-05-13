package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger es la instancia global del logger
var Logger *zap.Logger

// InitLogger inicializa el logger
func InitLogger() error {
	// Configurar encoder
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Configurar salida
	file, err := os.OpenFile("logs/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error al abrir archivo de log: %w", err)
	}

	// Crear core
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(
			zapcore.AddSync(os.Stdout),
			zapcore.AddSync(file),
		),
		zapcore.InfoLevel,
	)

	// Crear logger
	Logger = zap.New(core, zap.AddCaller())
	return nil
}

// LogRequest registra una petición HTTP
func LogRequest(c *gin.Context, status int, duration time.Duration) {
	Logger.Info("request",
		zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
		zap.String("ip", c.ClientIP()),
		zap.Int("status", status),
		zap.Duration("duration", duration),
		zap.String("user_agent", c.Request.UserAgent()),
	)
}

// LogError registra un error
func LogError(err error, fields ...zap.Field) {
	Logger.Error("error",
		append([]zap.Field{
			zap.Error(err),
			zap.String("stack", fmt.Sprintf("%+v", err)),
		}, fields...)...,
	)
}

// LogInfo registra información
func LogInfo(msg string, fields ...zap.Field) {
	Logger.Info(msg, fields...)
}

// LogWarning registra una advertencia
func LogWarning(msg string, fields ...zap.Field) {
	Logger.Warn(msg, fields...)
}

// LogDebug registra información de depuración
func LogDebug(msg string, fields ...zap.Field) {
	Logger.Debug(msg, fields...)
}
