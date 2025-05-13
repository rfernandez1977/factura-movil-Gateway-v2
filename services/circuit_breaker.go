package services

type CircuitBreaker struct {
	Name         string
	MaxFailures  int
	ResetTimeout int64
}

func NewCircuitBreaker(name string, maxFailures int, resetTimeout int64) *CircuitBreaker {
	return &CircuitBreaker{
		Name:         name,
		MaxFailures:  maxFailures,
		ResetTimeout: resetTimeout,
	}
}
