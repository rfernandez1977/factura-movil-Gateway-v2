package services

// AlertService define una interfaz para enviar alertas
type AlertService interface {
	SendAlert(subject string, payload interface{}) error
}

// DefaultAlertService is the default implementation of AlertService
type DefaultAlertService struct {
	NotificationEndpoint string
}

func NewAlertService(endpoint string) *DefaultAlertService {
	return &DefaultAlertService{
		NotificationEndpoint: endpoint,
	}
}

// SendAlert sends an alert to the notification endpoint
func (s *DefaultAlertService) SendAlert(subject string, payload interface{}) error {
	// Implementation would go here
	return nil
}

// NoOpAlertService implementa AlertService sin hacer nada
type NoOpAlertService struct{}

func NewNoOpAlertService() *NoOpAlertService {
	return &NoOpAlertService{}
}

func (s *NoOpAlertService) SendAlert(subject string, payload interface{}) error {
	// No hace nada, solo implementa la interfaz
	return nil
}
