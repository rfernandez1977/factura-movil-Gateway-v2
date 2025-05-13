package services

import (
	"sync"
)

// Event representa un evento del sistema
type Event struct {
	Type    string
	Payload interface{}
}

// EventObserver define la interfaz para los observadores de eventos
type EventObserver interface {
	OnEvent(event Event)
}

// EventSubject maneja la suscripción y notificación de eventos
type EventSubject struct {
	observers []EventObserver
	mu        sync.RWMutex
}

func NewEventSubject() *EventSubject {
	return &EventSubject{
		observers: make([]EventObserver, 0),
	}
}

// Subscribe agrega un nuevo observador
func (s *EventSubject) Subscribe(observer EventObserver) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.observers = append(s.observers, observer)
}

// Unsubscribe elimina un observador
func (s *EventSubject) Unsubscribe(observer EventObserver) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, obs := range s.observers {
		if obs == observer {
			s.observers = append(s.observers[:i], s.observers[i+1:]...)
			break
		}
	}
}

// Notify notifica a todos los observadores sobre un evento
func (s *EventSubject) Notify(event Event) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, observer := range s.observers {
		observer.OnEvent(event)
	}
}

// LoggingObserver implementa un observador para logging
type LoggingObserver struct {
	logger Logger
}

func NewLoggingObserver(logger Logger) *LoggingObserver {
	return &LoggingObserver{logger: logger}
}

func (o *LoggingObserver) OnEvent(event Event) {
	o.logger.Log("Evento recibido: %s, Payload: %v", event.Type, event.Payload)
}

// MetricsObserver implementa un observador para métricas
type MetricsObserver struct {
	metrics Metrics
}

func NewMetricsObserver(metrics Metrics) *MetricsObserver {
	return &MetricsObserver{metrics: metrics}
}

func (o *MetricsObserver) OnEvent(event Event) {
	o.metrics.RecordEvent(event.Type)
}

// AlertObserver implementa un observador para alertas
type AlertObserver struct {
	alertService AlertService
}

func NewAlertObserver(alertService AlertService) *AlertObserver {
	return &AlertObserver{alertService: alertService}
}

func (o *AlertObserver) OnEvent(event Event) {
	if event.Type == "ERROR" {
		o.alertService.SendAlert("Error detectado", event.Payload)
	}
}
