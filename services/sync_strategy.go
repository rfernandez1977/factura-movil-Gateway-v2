package services

import (
	"context"
	"time"
)

// SyncStrategy define la interfaz para las estrategias de sincronización
type SyncStrategy interface {
	Execute(ctx context.Context, data interface{}) error
	GetName() string
	GetDescription() string
}

// RealTimeStrategy implementa sincronización en tiempo real
type RealTimeStrategy struct {
	connector ERPConnector
}

func NewRealTimeStrategy(connector ERPConnector) *RealTimeStrategy {
	return &RealTimeStrategy{connector: connector}
}

func (s *RealTimeStrategy) Execute(ctx context.Context, data interface{}) error {
	// Implementación de sincronización en tiempo real
	return nil
}

func (s *RealTimeStrategy) GetName() string {
	return "RealTime"
}

func (s *RealTimeStrategy) GetDescription() string {
	return "Sincronización inmediata de datos"
}

// BatchStrategy implementa sincronización por lotes
type BatchStrategy struct {
	connector ERPConnector
	batchSize int
	interval  time.Duration
}

func NewBatchStrategy(connector ERPConnector, batchSize int, interval time.Duration) *BatchStrategy {
	return &BatchStrategy{
		connector: connector,
		batchSize: batchSize,
		interval:  interval,
	}
}

func (s *BatchStrategy) Execute(ctx context.Context, data interface{}) error {
	// Implementación de sincronización por lotes
	return nil
}

func (s *BatchStrategy) GetName() string {
	return "Batch"
}

func (s *BatchStrategy) GetDescription() string {
	return "Sincronización por lotes programada"
}

// EventDrivenStrategy implementa sincronización basada en eventos
type EventDrivenStrategy struct {
	connector ERPConnector
	queue     Queue
}

func NewEventDrivenStrategy(connector ERPConnector, queue Queue) *EventDrivenStrategy {
	return &EventDrivenStrategy{
		connector: connector,
		queue:     queue,
	}
}

func (s *EventDrivenStrategy) Execute(ctx context.Context, data interface{}) error {
	// Implementación de sincronización basada en eventos
	return nil
}

func (s *EventDrivenStrategy) GetName() string {
	return "EventDriven"
}

func (s *EventDrivenStrategy) GetDescription() string {
	return "Sincronización basada en eventos"
}

// SyncContext maneja la ejecución de estrategias
type SyncContext struct {
	strategy SyncStrategy
}

func NewSyncContext(strategy SyncStrategy) *SyncContext {
	return &SyncContext{strategy: strategy}
}

func (c *SyncContext) SetStrategy(strategy SyncStrategy) {
	c.strategy = strategy
}

func (c *SyncContext) Execute(ctx context.Context, data interface{}) error {
	return c.strategy.Execute(ctx, data)
}
