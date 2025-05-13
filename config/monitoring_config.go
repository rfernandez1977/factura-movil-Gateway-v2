package config

import (
	"time"
)

// MonitoringConfig contiene la configuración del monitoreo
type MonitoringConfig struct {
	// Configuración de métricas
	Metrics struct {
		Enabled    bool               `json:"enabled"`
		Interval   time.Duration      `json:"interval"`
		Retention  time.Duration      `json:"retention"`
		Labels     []string           `json:"labels"`
		Thresholds map[string]float64 `json:"thresholds"`
	} `json:"metrics"`

	// Configuración de alertas
	Alerts struct {
		Enabled    bool          `json:"enabled"`
		Channels   []string      `json:"channels"`
		Severities []string      `json:"severities"`
		Cooldown   time.Duration `json:"cooldown"`
		Recipients []string      `json:"recipients"`
	} `json:"alerts"`

	// Configuración de logs
	Logs struct {
		Enabled   bool          `json:"enabled"`
		Level     string        `json:"level"`
		Format    string        `json:"format"`
		Output    string        `json:"output"`
		Retention time.Duration `json:"retention"`
	} `json:"logs"`

	// Configuración de almacenamiento
	Storage struct {
		Enabled     bool          `json:"enabled"`
		Type        string        `json:"type"`
		Path        string        `json:"path"`
		Compression bool          `json:"compression"`
		Retention   time.Duration `json:"retention"`
	} `json:"storage"`

	// Configuración de caché
	Cache struct {
		Enabled bool          `json:"enabled"`
		Type    string        `json:"type"`
		Size    int64         `json:"size"`
		TTL     time.Duration `json:"ttl"`
	} `json:"cache"`

	// Configuración de reintentos
	Retry struct {
		Enabled       bool          `json:"enabled"`
		MaxAttempts   int           `json:"max_attempts"`
		InitialDelay  time.Duration `json:"initial_delay"`
		MaxDelay      time.Duration `json:"max_delay"`
		BackoffFactor float64       `json:"backoff_factor"`
		JitterFactor  float64       `json:"jitter_factor"`
	} `json:"retry"`
}

// DefaultMonitoringConfig retorna una configuración por defecto
func DefaultMonitoringConfig() *MonitoringConfig {
	config := &MonitoringConfig{}

	// Configuración de métricas
	config.Metrics.Enabled = true
	config.Metrics.Interval = 1 * time.Minute
	config.Metrics.Retention = 7 * 24 * time.Hour
	config.Metrics.Labels = []string{"environment", "service", "instance"}
	config.Metrics.Thresholds = map[string]float64{
		"storage_size": 1073741824, // 1GB
		"cache_size":   536870912,  // 512MB
		"retry_rate":   0.1,        // 10%
		"error_rate":   0.01,       // 1%
		"latency_p95":  1.0,        // 1 segundo
	}

	// Configuración de alertas
	config.Alerts.Enabled = true
	config.Alerts.Channels = []string{"email", "slack"}
	config.Alerts.Severities = []string{"critical", "error", "warning", "info"}
	config.Alerts.Cooldown = 5 * time.Minute
	config.Alerts.Recipients = []string{"admin@example.com"}

	// Configuración de logs
	config.Logs.Enabled = true
	config.Logs.Level = "info"
	config.Logs.Format = "json"
	config.Logs.Output = "stdout"
	config.Logs.Retention = 30 * 24 * time.Hour

	// Configuración de almacenamiento
	config.Storage.Enabled = true
	config.Storage.Type = "local"
	config.Storage.Path = "/var/lib/fmgo/storage"
	config.Storage.Compression = true
	config.Storage.Retention = 365 * 24 * time.Hour

	// Configuración de caché
	config.Cache.Enabled = true
	config.Cache.Type = "redis"
	config.Cache.Size = 536870912 // 512MB
	config.Cache.TTL = 24 * time.Hour

	// Configuración de reintentos
	config.Retry.Enabled = true
	config.Retry.MaxAttempts = 5
	config.Retry.InitialDelay = 1 * time.Second
	config.Retry.MaxDelay = 30 * time.Second
	config.Retry.BackoffFactor = 2.0
	config.Retry.JitterFactor = 0.1

	return config
}
