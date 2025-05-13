package services

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AsyncService maneja el procesamiento asíncrono de tareas
type AsyncService struct {
	queue    *amqp.Channel
	workers  int
	handlers map[string]func(context.Context, []byte) error
	mu       sync.RWMutex
}

// NewAsyncService crea una nueva instancia del servicio asíncrono
func NewAsyncService(queue *amqp.Channel, workers int) *AsyncService {
	return &AsyncService{
		queue:    queue,
		workers:  workers,
		handlers: make(map[string]func(context.Context, []byte) error),
	}
}

// RegisterHandler registra un manejador para un tipo de tarea
func (s *AsyncService) RegisterHandler(taskType string, handler func(context.Context, []byte) error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.handlers[taskType] = handler
}

// Enqueue agrega una tarea a la cola
func (s *AsyncService) Enqueue(ctx context.Context, taskType string, data interface{}) error {
	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	msg := amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
		Timestamp:   time.Now(),
	}

	return s.queue.Publish(
		"",       // exchange
		taskType, // routing key
		false,    // mandatory
		false,    // immediate
		msg,
	)
}

// StartWorkers inicia los workers para procesar tareas
func (s *AsyncService) StartWorkers(ctx context.Context) error {
	for i := 0; i < s.workers; i++ {
		go s.worker(ctx, i)
	}
	return nil
}

// worker procesa tareas de la cola
func (s *AsyncService) worker(ctx context.Context, id int) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			s.processNextTask(ctx)
		}
	}
}

// processNextTask procesa la siguiente tarea de la cola
func (s *AsyncService) processNextTask(ctx context.Context) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for taskType, handler := range s.handlers {
		msgs, err := s.queue.Consume(
			taskType, // queue
			"",       // consumer
			false,    // auto-ack
			false,    // exclusive
			false,    // no-local
			false,    // no-wait
			nil,      // args
		)
		if err != nil {
			continue
		}

		for msg := range msgs {
			if err := handler(ctx, msg.Body); err != nil {
				// Reintentar la tarea
				msg.Nack(false, true)
				continue
			}
			msg.Ack(false)
		}
	}
}

// Task representa una tarea asíncrona
type Task struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Type      string             `json:"type" bson:"type"`
	Data      interface{}        `json:"data" bson:"data"`
	Status    string             `json:"status" bson:"status"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

// TaskStatus representa el estado de una tarea
type TaskStatus string

const (
	TaskStatusPending   TaskStatus = "pending"
	TaskStatusRunning   TaskStatus = "running"
	TaskStatusCompleted TaskStatus = "completed"
	TaskStatusFailed    TaskStatus = "failed"
)

// NewTask crea una nueva tarea
func NewTask(taskType string, data interface{}) *Task {
	return &Task{
		ID:        primitive.NewObjectID(),
		Type:      taskType,
		Data:      data,
		Status:    string(TaskStatusPending),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
