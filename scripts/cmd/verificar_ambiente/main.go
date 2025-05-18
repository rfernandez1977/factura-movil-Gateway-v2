package main

import (
	"context"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"FMgo/core/sii/client"
	siilogger "FMgo/core/sii/logger"
	"FMgo/core/sii/models"

	"software.sslmate.com/src/go-pkcs12"
)

type Config struct {
	RutEmpresa       string `json:"rut_empresa"`
	RutCertificado   string `json:"rut_certificado"`
	ClaveCertificado string `json:"clave_certificado"`
	PathCertificado  string `json:"path_certificado"`
	BaseURL          string `json:"base_url"`
	RetryCount       int    `json:"retry_count"`
	Timeout          int    `json:"timeout"`
	Ambiente         string `json:"ambiente"`
	CertPath         string `json:"cert_path"`
	KeyPath          string `json:"key_path"`
	SchemaPath       string `json:"schema_path"`
	RetryDelay       int    `json:"retry_delay"`
	Monitoring       struct {
		Enabled        bool   `json:"enabled"`
		LogLevel       string `json:"log_level"`
		MetricsEnabled bool   `json:"metrics_enabled"`
		AlertThreshold struct {
			ResponseTimeMs   int `json:"response_time_ms"`
			ErrorRatePercent int `json:"error_rate_percent"`
		} `json:"alert_threshold"`
	} `json:"monitoring"`
	Validation struct {
		StrictMode     bool `json:"strict_mode"`
		ValidateSchema bool `json:"validate_schema"`
		ValidateCAF    bool `json:"validate_caf"`
		ValidateFirma  bool `json:"validate_firma"`
	} `json:"validation"`
	Cache struct {
		Enabled    bool `json:"enabled"`
		TTLSeconds int  `json:"ttl_seconds"`
		MaxSizeMB  int  `json:"max_size_mb"`
	} `json:"cache"`
}

type Logger struct {
	debug bool
}

func NewLogger() *Logger {
	return &Logger{debug: true}
}

func (l *Logger) Debug(format string, args ...interface{}) {
	if l.debug {
		fmt.Printf("[DEBUG] "+format+"\n", args...)
	}
}

func (l *Logger) Info(format string, args ...interface{}) {
	fmt.Printf("[INFO] "+format+"\n", args...)
}

func (l *Logger) Warn(format string, args ...interface{}) {
	fmt.Printf("[WARN] "+format+"\n", args...)
}

func (l *Logger) Error(format string, args ...interface{}) {
	fmt.Printf("[ERROR] "+format+"\n", args...)
}

func (l *Logger) Fatal(format string, args ...interface{}) {
	fmt.Printf("[FATAL] "+format+"\n", args...)
	os.Exit(1)
}

func verificarDirectorios(log *Logger) error {
	directorios := map[string]string{
		"Esquemas XSD":       "schema_dte",
		"Certificados":       "firma_test",
		"Logs certificación": "logs",
		"Métricas":           "metrics",
		"Configuración":      "config",
	}

	for desc, dir := range directorios {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			log.Error("❌ %s: directorio no encontrado en %s", desc, dir)
			return fmt.Errorf("directorio no encontrado: %s", dir)
		}
		log.Info("✅ %s: directorio verificado en %s", desc, dir)
	}
	return nil
}

func verificarArchivosConfiguracion(log *Logger) error {
	// Primero leemos la configuración para obtener la ruta del certificado
	configPath := "config/sii_config.json"
	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Error("❌ Error leyendo archivo de configuración: %v", err)
		return fmt.Errorf("error leyendo archivo de configuración: %v", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		log.Error("❌ Error parseando configuración: %v", err)
		return fmt.Errorf("error parseando configuración: %v", err)
	}

	archivos := map[string]string{
		"Configuración SII": configPath,
		"Certificado PFX":   config.PathCertificado,
		"Schema SiiTypes":   "schema_dte/SiiTypes_v10.xsd",
		"Schema DTE":        "schema_dte/DTE_v10.xsd",
		"Schema EnvioDTE":   "schema_dte/EnvioDTE_v10.xsd",
	}

	for desc, archivo := range archivos {
		if _, err := os.Stat(archivo); os.IsNotExist(err) {
			log.Error("❌ %s: archivo no encontrado en %s", desc, archivo)
			return fmt.Errorf("archivo no encontrado: %s", archivo)
		}
		log.Info("✅ %s: archivo verificado en %s", desc, archivo)
	}
	return nil
}

func cargarConfiguracion(log *Logger) (*models.Config, error) {
	configPath := "config/sii_config.json"
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("error leyendo archivo de configuración: %v", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("error parseando configuración: %v", err)
	}

	// Leer el archivo PFX
	pfxData, err := os.ReadFile(config.PathCertificado)
	if err != nil {
		return nil, fmt.Errorf("error leyendo certificado PFX: %v", err)
	}

	// Decodificar el PFX
	privateKey, cert, err := pkcs12.Decode(pfxData, config.ClaveCertificado)
	if err != nil {
		return nil, fmt.Errorf("error decodificando PFX: %v", err)
	}

	// Crear un archivo temporal para el certificado PEM
	certFile, err := os.CreateTemp("", "cert*.pem")
	if err != nil {
		return nil, fmt.Errorf("error creando archivo temporal para certificado: %v", err)
	}
	defer os.Remove(certFile.Name())

	// Crear un archivo temporal para la llave privada PEM
	keyFile, err := os.CreateTemp("", "key*.pem")
	if err != nil {
		return nil, fmt.Errorf("error creando archivo temporal para llave privada: %v", err)
	}
	defer os.Remove(keyFile.Name())

	// Escribir el certificado y la llave privada en formato PEM
	if err := os.WriteFile(certFile.Name(), cert.Raw, 0600); err != nil {
		return nil, fmt.Errorf("error escribiendo certificado temporal: %v", err)
	}

	keyPEM, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return nil, fmt.Errorf("error codificando llave privada: %v", err)
	}

	if err := os.WriteFile(keyFile.Name(), keyPEM, 0600); err != nil {
		return nil, fmt.Errorf("error escribiendo llave privada temporal: %v", err)
	}

	siiConfig := models.NewConfig()
	siiConfig.Ambiente = models.Certificacion
	siiConfig.CertPath = certFile.Name()
	siiConfig.KeyPath = keyFile.Name()
	siiConfig.SchemaPath = "schema_dte/SiiTypes_v10.xsd"
	siiConfig.Timeout = time.Duration(config.Timeout) * time.Second
	siiConfig.RetryCount = config.RetryCount
	siiConfig.RetryDelay = time.Duration(config.RetryDelay) * time.Second

	return siiConfig, nil
}

func verificarConectividadSII(ctx context.Context, config *models.Config, log siilogger.Logger) error {
	httpClient, err := client.NewHTTPClient(config, log)
	if err != nil {
		return fmt.Errorf("error creando cliente HTTP: %v", err)
	}

	log.Info("Verificando comunicación con SII...")
	if err := httpClient.VerificarComunicacion(ctx); err != nil {
		log.Error("❌ Error en comunicación con SII")
		return fmt.Errorf("error verificando comunicación: %v", err)
	}
	log.Info("✅ Comunicación con SII verificada")

	log.Info("Probando obtención de semilla...")
	semilla, err := httpClient.ObtenerSemilla(ctx)
	if err != nil {
		log.Error("❌ Error obteniendo semilla")
		return fmt.Errorf("error obteniendo semilla: %v", err)
	}
	log.Info("✅ Semilla obtenida exitosamente: %s", semilla)

	log.Info("Probando obtención de token...")
	token, err := httpClient.ObtenerToken(ctx, semilla)
	if err != nil {
		log.Error("❌ Error obteniendo token")
		return fmt.Errorf("error obteniendo token: %v", err)
	}
	log.Info("✅ Token obtenido exitosamente: %s", token)

	return nil
}

func verificarAmbiente() error {
	log := NewLogger()
	log.Info("🔍 Iniciando verificación del ambiente de certificación...")

	if err := verificarDirectorios(log); err != nil {
		return fmt.Errorf("error en estructura de directorios: %v", err)
	}

	if err := verificarArchivosConfiguracion(log); err != nil {
		return fmt.Errorf("error en archivos de configuración: %v", err)
	}

	config, err := cargarConfiguracion(log)
	if err != nil {
		return fmt.Errorf("error cargando configuración: %v", err)
	}
	log.Info("✅ Configuración cargada correctamente")

	ctx := context.Background()
	if err := verificarConectividadSII(ctx, config, log); err != nil {
		return fmt.Errorf("error verificando conectividad con SII: %v", err)
	}

	log.Info("✅ Verificación del ambiente completada exitosamente")
	return nil
}

func main() {
	if err := verificarAmbiente(); err != nil {
		fmt.Printf("❌ Error en la verificación del ambiente: %v\n", err)
		os.Exit(1)
	}
}
