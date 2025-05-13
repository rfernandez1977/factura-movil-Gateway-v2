package services

import (
	"context"
	"sync"
	"time"
)

// LoadBalancer maneja el balanceo de carga
type LoadBalancer struct {
	mu            sync.RWMutex
	workers       map[string]*WorkerStats
	strategy      LoadBalancingStrategy
	healthChecker *HealthChecker
}

// WorkerStats contiene estadísticas de un worker
type WorkerStats struct {
	ID            string
	ActiveTasks   int
	TotalTasks    int
	ErrorCount    int
	ResponseTime  time.Duration
	LastHeartbeat time.Time
	IsHealthy     bool
}

// LoadBalancingStrategy define la estrategia de balanceo de carga
type LoadBalancingStrategy interface {
	SelectWorker(workers map[string]*WorkerStats) string
}

// RoundRobinStrategy implementa balanceo de carga round-robin
type RoundRobinStrategy struct {
	lastIndex int
}

// LeastConnectionsStrategy implementa balanceo por menor número de conexiones
type LeastConnectionsStrategy struct{}

// WeightedRoundRobinStrategy implementa balanceo round-robin con pesos
type WeightedRoundRobinStrategy struct {
	lastIndex int
	weights   map[string]int
}

// HealthChecker verifica la salud de los workers
type HealthChecker struct {
	interval time.Duration
	timeout  time.Duration
}

// NewLoadBalancer crea una nueva instancia del balanceador de carga
func NewLoadBalancer(strategy LoadBalancingStrategy, healthCheckInterval, healthCheckTimeout time.Duration) *LoadBalancer {
	lb := &LoadBalancer{
		workers:  make(map[string]*WorkerStats),
		strategy: strategy,
		healthChecker: &HealthChecker{
			interval: healthCheckInterval,
			timeout:  healthCheckTimeout,
		},
	}
	go lb.startHealthChecks()
	return lb
}

// RegisterWorker registra un nuevo worker
func (lb *LoadBalancer) RegisterWorker(id string) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	lb.workers[id] = &WorkerStats{
		ID:            id,
		LastHeartbeat: time.Now(),
		IsHealthy:     true,
	}
}

// UnregisterWorker elimina un worker
func (lb *LoadBalancer) UnregisterWorker(id string) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	delete(lb.workers, id)
}

// UpdateWorkerStats actualiza las estadísticas de un worker
func (lb *LoadBalancer) UpdateWorkerStats(id string, activeTasks int, responseTime time.Duration) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if stats, ok := lb.workers[id]; ok {
		stats.ActiveTasks = activeTasks
		stats.TotalTasks++
		stats.ResponseTime = responseTime
		stats.LastHeartbeat = time.Now()
	}
}

// SelectWorker selecciona un worker según la estrategia
func (lb *LoadBalancer) SelectWorker() string {
	lb.mu.RLock()
	defer lb.mu.RUnlock()

	// Filtrar workers saludables
	healthyWorkers := make(map[string]*WorkerStats)
	for id, stats := range lb.workers {
		if stats.IsHealthy {
			healthyWorkers[id] = stats
		}
	}

	if len(healthyWorkers) == 0 {
		return ""
	}

	return lb.strategy.SelectWorker(healthyWorkers)
}

// startHealthChecks inicia las verificaciones de salud
func (lb *LoadBalancer) startHealthChecks() {
	ticker := time.NewTicker(lb.healthChecker.interval)
	defer ticker.Stop()

	for range ticker.C {
		lb.checkWorkersHealth()
	}
}

// checkWorkersHealth verifica la salud de todos los workers
func (lb *LoadBalancer) checkWorkersHealth() {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	now := time.Now()
	for _, stats := range lb.workers {
		if now.Sub(stats.LastHeartbeat) > lb.healthChecker.timeout {
			stats.IsHealthy = false
		}
	}
}

// SelectWorker implementa la estrategia round-robin
func (s *RoundRobinStrategy) SelectWorker(workers map[string]*WorkerStats) string {
	ids := make([]string, 0, len(workers))
	for id := range workers {
		ids = append(ids, id)
	}

	if len(ids) == 0 {
		return ""
	}

	s.lastIndex = (s.lastIndex + 1) % len(ids)
	return ids[s.lastIndex]
}

// SelectWorker implementa la estrategia de menor número de conexiones
func (s *LeastConnectionsStrategy) SelectWorker(workers map[string]*WorkerStats) string {
	var selectedID string
	minTasks := int(^uint(0) >> 1) // Max int

	for id, stats := range workers {
		if stats.ActiveTasks < minTasks {
			minTasks = stats.ActiveTasks
			selectedID = id
		}
	}

	return selectedID
}

// SelectWorker implementa la estrategia round-robin con pesos
func (s *WeightedRoundRobinStrategy) SelectWorker(workers map[string]*WorkerStats) string {
	// Implementar lógica de selección con pesos
	// ...
	return ""
}

// NewHealthChecker crea un nuevo verificador de salud
func NewHealthChecker(interval, timeout time.Duration) *HealthChecker {
	return &HealthChecker{
		interval: interval,
		timeout:  timeout,
	}
}

// Check verifica la salud de un worker
func (hc *HealthChecker) Check(ctx context.Context, workerID string) bool {
	// Implementar lógica de verificación de salud
	// ...
	return true
}
