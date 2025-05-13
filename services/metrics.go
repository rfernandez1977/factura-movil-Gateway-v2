package services

// Metrics define una interfaz para registrar métricas
type Metrics interface {
	RecordEvent(eventType string)
	IncrementCounter(name string)
	RecordValue(name string, value float64)
	StartTimer(name string) Timer
}

// Timer define una interfaz para medir tiempos de ejecución
type Timer interface {
	Stop() float64
}

// NoOpMetrics implementación de métricas que no hace nada
type NoOpMetrics struct{}

func NewNoOpMetrics() *NoOpMetrics {
	return &NoOpMetrics{}
}

func (m *NoOpMetrics) RecordEvent(eventType string) {}

func (m *NoOpMetrics) IncrementCounter(name string) {}

func (m *NoOpMetrics) RecordValue(name string, value float64) {}

func (m *NoOpMetrics) StartTimer(name string) Timer {
	return &NoOpTimer{}
}

// NoOpTimer implementación de timer que no hace nada
type NoOpTimer struct{}

func (t *NoOpTimer) Stop() float64 {
	return 0
}
