package services

import (
	"sync"
	"time"
)

// AutoScaler maneja el escalado automático de recursos
type AutoScaler struct {
	mu            sync.RWMutex
	metrics       *MetricsCollector
	config        *ScalingConfig
	currentScale  int
	minScale      int
	maxScale      int
	scaleUpChan   chan struct{}
	scaleDownChan chan struct{}
}

// MetricsCollector recopila métricas para el escalado
type MetricsCollector struct {
	mu           sync.RWMutex
	cpuUsage     float64
	memoryUsage  float64
	requestRate  float64
	responseTime time.Duration
	errorRate    float64
	lastUpdated  time.Time
}

// ScalingConfig contiene la configuración de escalado
type ScalingConfig struct {
	ScaleUpThreshold   float64
	ScaleDownThreshold float64
	CooldownPeriod     time.Duration
	ScaleUpStep        int
	ScaleDownStep      int
	MaxScale           int
	MinScale           int
}

// NewAutoScaler crea una nueva instancia del auto-scaler
func NewAutoScaler(config *ScalingConfig) *AutoScaler {
	scaler := &AutoScaler{
		metrics:       &MetricsCollector{},
		config:        config,
		currentScale:  config.MinScale,
		minScale:      config.MinScale,
		maxScale:      config.MaxScale,
		scaleUpChan:   make(chan struct{}, 1),
		scaleDownChan: make(chan struct{}, 1),
	}
	go scaler.startScalingLoop()
	return scaler
}

// UpdateMetrics actualiza las métricas del sistema
func (s *AutoScaler) UpdateMetrics(metrics *MetricsCollector) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.metrics = metrics
}

// GetCurrentScale obtiene la escala actual
func (s *AutoScaler) GetCurrentScale() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.currentScale
}

// startScalingLoop inicia el bucle de escalado
func (s *AutoScaler) startScalingLoop() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		s.evaluateScaling()
	}
}

// evaluateScaling evalúa si es necesario escalar
func (s *AutoScaler) evaluateScaling() {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Verificar si estamos en período de cooldown
	if time.Since(s.metrics.lastUpdated) < s.config.CooldownPeriod {
		return
	}

	// Evaluar métricas para escalado
	if s.shouldScaleUp() {
		s.scaleUp()
	} else if s.shouldScaleDown() {
		s.scaleDown()
	}
}

// shouldScaleUp determina si se debe escalar hacia arriba
func (s *AutoScaler) shouldScaleUp() bool {
	return s.metrics.cpuUsage > s.config.ScaleUpThreshold ||
		s.metrics.memoryUsage > s.config.ScaleUpThreshold ||
		s.metrics.responseTime > 500*time.Millisecond ||
		s.metrics.errorRate > 0.05
}

// shouldScaleDown determina si se debe escalar hacia abajo
func (s *AutoScaler) shouldScaleDown() bool {
	return s.metrics.cpuUsage < s.config.ScaleDownThreshold &&
		s.metrics.memoryUsage < s.config.ScaleDownThreshold &&
		s.metrics.responseTime < 100*time.Millisecond &&
		s.metrics.errorRate < 0.01
}

// scaleUp aumenta la escala
func (s *AutoScaler) scaleUp() {
	if s.currentScale >= s.maxScale {
		return
	}

	newScale := s.currentScale + s.config.ScaleUpStep
	if newScale > s.maxScale {
		newScale = s.maxScale
	}

	s.currentScale = newScale
	s.scaleUpChan <- struct{}{}
}

// scaleDown disminuye la escala
func (s *AutoScaler) scaleDown() {
	if s.currentScale <= s.minScale {
		return
	}

	newScale := s.currentScale - s.config.ScaleDownStep
	if newScale < s.minScale {
		newScale = s.minScale
	}

	s.currentScale = newScale
	s.scaleDownChan <- struct{}{}
}

// UpdateCPUUsage actualiza el uso de CPU
func (m *MetricsCollector) UpdateCPUUsage(usage float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.cpuUsage = usage
	m.lastUpdated = time.Now()
}

// UpdateMemoryUsage actualiza el uso de memoria
func (m *MetricsCollector) UpdateMemoryUsage(usage float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.memoryUsage = usage
	m.lastUpdated = time.Now()
}

// UpdateRequestRate actualiza la tasa de solicitudes
func (m *MetricsCollector) UpdateRequestRate(rate float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.requestRate = rate
	m.lastUpdated = time.Now()
}

// UpdateResponseTime actualiza el tiempo de respuesta
func (m *MetricsCollector) UpdateResponseTime(duration time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.responseTime = duration
	m.lastUpdated = time.Now()
}

// UpdateErrorRate actualiza la tasa de errores
func (m *MetricsCollector) UpdateErrorRate(rate float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.errorRate = rate
	m.lastUpdated = time.Now()
}

// GetMetrics obtiene las métricas actuales
func (m *MetricsCollector) GetMetrics() *MetricsCollector {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m
}

// NewScalingConfig crea una nueva configuración de escalado
func NewScalingConfig() *ScalingConfig {
	return &ScalingConfig{
		ScaleUpThreshold:   0.8, // 80% de uso
		ScaleDownThreshold: 0.2, // 20% de uso
		CooldownPeriod:     5 * time.Minute,
		ScaleUpStep:        1,
		ScaleDownStep:      1,
		MaxScale:           10,
		MinScale:           1,
	}
}
