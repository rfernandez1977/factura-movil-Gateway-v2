package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// LoadConfig representa la configuración específica para el script de carga
type LoadConfig struct {
	Supabase struct {
		URL         string `json:"url"`
		Key         string `json:"key"`
		JWTSecret   string `json:"jwt_secret"`
		Environment string `json:"environment"`
	} `json:"supabase"`
	Database struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		Name     string `json:"name"`
		User     string `json:"user"`
		Password string `json:"password"`
		SSLMode  string `json:"ssl_mode"`
	} `json:"database"`
	SII struct {
		URL         string `json:"url"`
		Environment string `json:"environment"`
		Timeout     int    `json:"timeout"`
	} `json:"sii"`
	Storage struct {
		Path               string `json:"path"`
		CompressionEnabled bool   `json:"compression_enabled"`
	} `json:"storage"`
	Redis struct {
		Host         string `json:"host"`
		Port         int    `json:"port"`
		Password     string `json:"password"`
		DB           int    `json:"db"`
		MaxRetries   int    `json:"max_retries"`
		PoolSize     int    `json:"pool_size"`
		MinIdleConns int    `json:"min_idle_conns"`
	} `json:"redis"`
}

func main() {
	// Leer archivo de configuración
	configPath := filepath.Join("config", "supabase-config.json")
	data, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Printf("Error leyendo archivo de configuración: %v\n", err)
		os.Exit(1)
	}

	// Decodificar JSON
	var config LoadConfig
	if err := json.Unmarshal(data, &config); err != nil {
		fmt.Printf("Error decodificando configuración: %v\n", err)
		os.Exit(1)
	}

	// Establecer variables de entorno
	os.Setenv("SUPABASE_URL", config.Supabase.URL)
	os.Setenv("SUPABASE_KEY", config.Supabase.Key)
	os.Setenv("SUPABASE_JWT_SECRET", config.Supabase.JWTSecret)
	os.Setenv("SUPABASE_ENVIRONMENT", config.Supabase.Environment)

	os.Setenv("DB_HOST", config.Database.Host)
	os.Setenv("DB_PORT", fmt.Sprintf("%d", config.Database.Port))
	os.Setenv("DB_NAME", config.Database.Name)
	os.Setenv("DB_USER", config.Database.User)
	os.Setenv("DB_PASSWORD", config.Database.Password)
	os.Setenv("DB_SSL_MODE", config.Database.SSLMode)

	os.Setenv("SII_URL", config.SII.URL)
	os.Setenv("SII_ENVIRONMENT", config.SII.Environment)
	os.Setenv("SII_TIMEOUT", fmt.Sprintf("%d", config.SII.Timeout))

	os.Setenv("STORAGE_PATH", config.Storage.Path)
	os.Setenv("COMPRESSION_ENABLED", fmt.Sprintf("%v", config.Storage.CompressionEnabled))

	os.Setenv("REDIS_HOST", config.Redis.Host)
	os.Setenv("REDIS_PORT", fmt.Sprintf("%d", config.Redis.Port))
	os.Setenv("REDIS_PASSWORD", config.Redis.Password)
	os.Setenv("REDIS_DB", fmt.Sprintf("%d", config.Redis.DB))
	os.Setenv("REDIS_MAX_RETRIES", fmt.Sprintf("%d", config.Redis.MaxRetries))
	os.Setenv("REDIS_POOL_SIZE", fmt.Sprintf("%d", config.Redis.PoolSize))
	os.Setenv("REDIS_MIN_IDLE_CONNS", fmt.Sprintf("%d", config.Redis.MinIdleConns))

	fmt.Println("Configuración cargada exitosamente")
}
