package services

import (
	"context"
	"sync"
	"time"
)

// ParallelService maneja el procesamiento en paralelo
type ParallelService struct {
	maxWorkers int
	timeout    time.Duration
}

// NewParallelService crea una nueva instancia del servicio de procesamiento paralelo
func NewParallelService(maxWorkers int, timeout time.Duration) *ParallelService {
	return &ParallelService{
		maxWorkers: maxWorkers,
		timeout:    timeout,
	}
}

// ProcessItems procesa elementos en paralelo
func (s *ParallelService) ProcessItems(ctx context.Context, items []interface{}, processor func(context.Context, interface{}) error) error {
	// Crear un canal para los resultados
	results := make(chan error, len(items))
	// Crear un canal para controlar el número de workers
	semaphore := make(chan struct{}, s.maxWorkers)

	// Crear un contexto con timeout
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	var wg sync.WaitGroup

	// Procesar cada elemento en un goroutine
	for _, item := range items {
		wg.Add(1)
		go func(item interface{}) {
			defer wg.Done()
			// Adquirir un slot en el semáforo
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			// Procesar el elemento
			err := processor(ctx, item)
			results <- err
		}(item)
	}

	// Esperar a que todos los workers terminen
	go func() {
		wg.Wait()
		close(results)
	}()

	// Recolectar resultados
	var errors []error
	for err := range results {
		if err != nil {
			errors = append(errors, err)
		}
	}

	// Retornar el primer error encontrado
	if len(errors) > 0 {
		return errors[0]
	}

	return nil
}

// BatchProcess procesa elementos en lotes
func (s *ParallelService) BatchProcess(ctx context.Context, items []interface{}, batchSize int, processor func(context.Context, []interface{}) error) error {
	// Dividir los elementos en lotes
	batches := make([][]interface{}, 0)
	for i := 0; i < len(items); i += batchSize {
		end := i + batchSize
		if end > len(items) {
			end = len(items)
		}
		batches = append(batches, items[i:end])
	}

	// Procesar los lotes en paralelo
	return s.ProcessItems(ctx, batches, func(ctx context.Context, batch interface{}) error {
		return processor(ctx, batch.([]interface{}))
	})
}

// MapReduce implementa el patrón MapReduce
func (s *ParallelService) MapReduce(ctx context.Context, items []interface{},
	mapper func(context.Context, interface{}) (interface{}, error),
	reducer func(context.Context, []interface{}) (interface{}, error)) (interface{}, error) {

	// Crear un canal para los resultados del mapper
	mappedResults := make(chan interface{}, len(items))
	// Crear un canal para controlar el número de workers
	semaphore := make(chan struct{}, s.maxWorkers)

	var wg sync.WaitGroup

	// Aplicar el mapper en paralelo
	for _, item := range items {
		wg.Add(1)
		go func(item interface{}) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			result, err := mapper(ctx, item)
			if err != nil {
				mappedResults <- err
				return
			}
			mappedResults <- result
		}(item)
	}

	// Esperar a que todos los mappers terminen
	go func() {
		wg.Wait()
		close(mappedResults)
	}()

	// Recolectar resultados del mapper
	var mappedItems []interface{}
	for result := range mappedResults {
		if err, ok := result.(error); ok {
			return nil, err
		}
		mappedItems = append(mappedItems, result)
	}

	// Aplicar el reducer
	return reducer(ctx, mappedItems)
}

// WorkerPool implementa un pool de workers
type WorkerPool struct {
	workers    int
	jobQueue   chan interface{}
	resultChan chan interface{}
	processor  func(context.Context, interface{}) (interface{}, error)
}

// NewWorkerPool crea un nuevo pool de workers
func NewWorkerPool(workers int, processor func(context.Context, interface{}) (interface{}, error)) *WorkerPool {
	return &WorkerPool{
		workers:    workers,
		jobQueue:   make(chan interface{}, workers*2),
		resultChan: make(chan interface{}, workers*2),
		processor:  processor,
	}
}

// Start inicia el pool de workers
func (p *WorkerPool) Start(ctx context.Context) {
	for i := 0; i < p.workers; i++ {
		go p.worker(ctx)
	}
}

// worker procesa trabajos del pool
func (p *WorkerPool) worker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case job := <-p.jobQueue:
			result, err := p.processor(ctx, job)
			if err != nil {
				p.resultChan <- err
			} else {
				p.resultChan <- result
			}
		}
	}
}

// Submit agrega un trabajo al pool
func (p *WorkerPool) Submit(job interface{}) {
	p.jobQueue <- job
}

// GetResult obtiene un resultado del pool
func (p *WorkerPool) GetResult() (interface{}, error) {
	result := <-p.resultChan
	if err, ok := result.(error); ok {
		return nil, err
	}
	return result, nil
}
