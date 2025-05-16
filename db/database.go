package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/fmgo/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB es la instancia global de la base de datos
var DB *gorm.DB

// InitDB inicializa la conexión a la base de datos
func InitDB(cfg *config.Config) error {
	// Configurar logger de GORM
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	// Conectar a la base de datos usando la configuración centralizada
	var err error
	DB, err = gorm.Open(postgres.Open(config.GetDSN(cfg)), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return fmt.Errorf("error al conectar a la base de datos: %w", err)
	}

	// Configurar conexión
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("error al obtener conexión SQL: %w", err)
	}

	// Configurar pool de conexiones
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return nil
}

// getEnv obtiene una variable de entorno o devuelve un valor por defecto
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
