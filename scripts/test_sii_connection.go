package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"FMgo/core/sii/client"
	"FMgo/core/sii/models"
	"FMgo/utils/logger"
)

func loadConfig() (*models.Config, error) {
	data, err := os.ReadFile("config/sii_config.json")
	if err != nil {
		return nil, fmt.Errorf("error leyendo archivo de configuración: %w", err)
	}

	var config models.Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("error decodificando configuración: %w", err)
	}

	return &config, nil
}

func main() {
	// 1. Crear logger
	logger := logger.NewLogger()

	// 2. Cargar configuración
	config, err := loadConfig()
	if err != nil {
		logger.Fatal("Error cargando configuración: %v", err)
	}

	// 3. Crear cliente HTTP
	httpClient, err := client.NewHTTPClient(config, logger)
	if err != nil {
		logger.Fatal("Error creando cliente HTTP: %v", err)
	}

	// 4. Crear contexto con timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 5. Probar obtención de semilla
	logger.Info("Probando obtención de semilla...")
	semilla, err := httpClient.ObtenerSemilla(ctx)
	if err != nil {
		logger.Fatal("Error obteniendo semilla: %v", err)
	}
	logger.Info("✅ Semilla obtenida exitosamente: %s", semilla)

	// 6. Probar obtención de token
	logger.Info("Probando obtención de token...")
	token, err := httpClient.ObtenerToken(ctx, semilla)
	if err != nil {
		logger.Fatal("Error obteniendo token: %v", err)
	}
	logger.Info("✅ Token obtenido exitosamente: %s", token)

	// 7. Probar verificación de comunicación
	logger.Info("Verificando comunicación general...")
	err = httpClient.VerificarComunicacion(ctx)
	if err != nil {
		logger.Fatal("Error verificando comunicación: %v", err)
	}
	logger.Info("✅ Comunicación verificada exitosamente")

	logger.Info("✅ Todas las pruebas de conexión completadas exitosamente")
}
