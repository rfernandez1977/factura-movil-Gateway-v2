package monitoring

import (
	"runtime"
	"time"

	"FMgo/models"
)

// MonitoringService implementa el servicio de monitoreo
type MonitoringService struct {
	startTime time.Time
	config    *models.Config
}

// ServiceStatus representa el estado actual del servicio
type ServiceStatus struct {
	Status     string           `json:"status"`
	Uptime     string           `json:"uptime"`
	StartTime  time.Time        `json:"start_time"`
	Memory     MemoryStats      `json:"memory"`
	GoRoutines int              `json:"go_routines"`
	Counters   map[string]int64 `json:"counters"`
	Errors     map[string]int64 `json:"errors"`
}

// MemoryStats contiene las estadísticas de memoria
type MemoryStats struct {
	Alloc        uint64 `json:"alloc"`
	TotalAlloc   uint64 `json:"total_alloc"`
	Sys          uint64 `json:"sys"`
	NumGC        uint32 `json:"num_gc"`
	PauseTotalNs uint64 `json:"pause_total_ns"`
}

// NewMonitoringService crea una nueva instancia de MonitoringService
func NewMonitoringService(config *models.Config) *MonitoringService {
	return &MonitoringService{
		startTime: time.Now(),
		config:    config,
	}
}

// GetStatus obtiene el estado actual del servicio
func (s *MonitoringService) GetStatus() *ServiceStatus {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	return &ServiceStatus{
		Status:     "OK",
		Uptime:     time.Since(s.startTime).String(),
		StartTime:  s.startTime,
		GoRoutines: runtime.NumGoroutine(),
		Memory: MemoryStats{
			Alloc:        memStats.Alloc,
			TotalAlloc:   memStats.TotalAlloc,
			Sys:          memStats.Sys,
			NumGC:        memStats.NumGC,
			PauseTotalNs: memStats.PauseTotalNs,
		},
		Counters: make(map[string]int64),
		Errors:   make(map[string]int64),
	}
}

// GetMetrics obtiene las métricas del sistema
func (s *MonitoringService) GetMetrics() map[string]interface{} {
	metrics := make(map[string]interface{})

	// Métricas básicas del sistema
	metrics["system"] = map[string]interface{}{
		"num_cpu":       runtime.NumCPU(),
		"num_goroutine": runtime.NumGoroutine(),
		"go_version":    runtime.Version(),
	}

	// Métricas de memoria
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	metrics["memory"] = map[string]uint64{
		"alloc":       memStats.Alloc,
		"total_alloc": memStats.TotalAlloc,
		"sys":         memStats.Sys,
		"heap_alloc":  memStats.HeapAlloc,
		"heap_sys":    memStats.HeapSys,
	}

	// Otras métricas personalizadas
	metrics["uptime"] = time.Since(s.startTime).Seconds()

	return metrics
}
