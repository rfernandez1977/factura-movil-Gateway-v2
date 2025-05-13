package services

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// SMSService maneja el envío de mensajes SMS
type SMSService struct {
	apiKey    string
	apiSecret string
	from      string
	baseURL   string
}

// NewSMSService crea una nueva instancia del servicio de SMS
func NewSMSService(apiKey, apiSecret, from, baseURL string) *SMSService {
	return &SMSService{
		apiKey:    apiKey,
		apiSecret: apiSecret,
		from:      from,
		baseURL:   baseURL,
	}
}

// EnviarSMS envía un mensaje SMS
func (s *SMSService) EnviarSMS(to, template string, data map[string]interface{}) error {
	// Procesar plantilla
	message := s.procesarPlantilla(template, data)

	// Construir URL
	params := url.Values{}
	params.Set("api_key", s.apiKey)
	params.Set("api_secret", s.apiSecret)
	params.Set("from", s.from)
	params.Set("to", to)
	params.Set("text", message)

	// Enviar petición
	resp, err := http.Post(
		fmt.Sprintf("%s/send?%s", s.baseURL, params.Encode()),
		"application/x-www-form-urlencoded",
		nil,
	)
	if err != nil {
		return fmt.Errorf("error enviando SMS: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error del servicio SMS: %s", resp.Status)
	}

	return nil
}

// procesarPlantilla procesa una plantilla de SMS con los datos proporcionados
func (s *SMSService) procesarPlantilla(template string, data map[string]interface{}) string {
	result := template
	for key, value := range data {
		placeholder := fmt.Sprintf("{{%s}}", key)
		result = strings.ReplaceAll(result, placeholder, fmt.Sprintf("%v", value))
	}
	return result
}
