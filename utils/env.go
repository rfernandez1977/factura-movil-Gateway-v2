package utils

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// LoadEnv carga las variables de entorno desde un archivo .env
func LoadEnv() error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("error loading .env file: %v", err)
	}
	return nil
}

// GetEnv obtiene una variable de entorno
func GetEnv(key string) string {
	return os.Getenv(key)
}

// GetEnvAsInt obtiene una variable de entorno como entero
func GetEnvAsInt(key string) (int, error) {
	value := GetEnv(key)
	if value == "" {
		return 0, fmt.Errorf("environment variable %s not found", key)
	}
	return strconv.Atoi(value)
}

// GetEnvAsInt64 obtiene una variable de entorno como entero de 64 bits
func GetEnvAsInt64(key string) (int64, error) {
	value := GetEnv(key)
	if value == "" {
		return 0, fmt.Errorf("environment variable %s not found", key)
	}
	return strconv.ParseInt(value, 10, 64)
}

// GetEnvAsBool obtiene una variable de entorno como booleano
func GetEnvAsBool(key string) (bool, error) {
	value := GetEnv(key)
	if value == "" {
		return false, fmt.Errorf("environment variable %s not found", key)
	}
	return strconv.ParseBool(value)
}

// GetEnvAsFloat obtiene una variable de entorno como float64
func GetEnvAsFloat(key string) (float64, error) {
	value := GetEnv(key)
	if value == "" {
		return 0, fmt.Errorf("environment variable %s not found", key)
	}
	return strconv.ParseFloat(value, 64)
}

// GetEnvAsStringSlice obtiene una variable de entorno como slice de strings
func GetEnvAsStringSlice(key string) ([]string, error) {
	value := GetEnv(key)
	if value == "" {
		return nil, fmt.Errorf("environment variable %s not found", key)
	}
	return strings.Split(value, ","), nil
}
