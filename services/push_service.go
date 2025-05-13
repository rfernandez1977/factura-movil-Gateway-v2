package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// PushService maneja el envío de notificaciones push
type PushService struct {
	apiKey    string
	apiSecret string
	baseURL   string
}

// NewPushService crea una nueva instancia del servicio de push
func NewPushService(apiKey, apiSecret, baseURL string) *PushService {
	return &PushService{
		apiKey:    apiKey,
		apiSecret: apiSecret,
		baseURL:   baseURL,
	}
}

// EnviarPush envía una notificación push
func (s *PushService) EnviarPush(tokens []string, title, template string, data map[string]interface{}) error {
	// Procesar plantilla
	message := s.procesarPlantilla(template, data)

	// Construir payload
	payload := map[string]interface{}{
		"registration_ids": tokens,
		"notification": map[string]interface{}{
			"title": title,
			"body":  message,
		},
		"data": data,
	}

	// Convertir a JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error codificando payload: %v", err)
	}

	// Crear petición
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/send", s.baseURL), bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creando petición: %v", err)
	}

	// Añadir headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("key=%s", s.apiKey))

	// Enviar petición
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error enviando push: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error del servicio push: %s", resp.Status)
	}

	return nil
}

// procesarPlantilla procesa una plantilla de push con los datos proporcionados
func (s *PushService) procesarPlantilla(template string, data map[string]interface{}) string {
	// En push notifications, el procesamiento de plantillas se hace en el cliente
	// por lo que solo devolvemos la plantilla original
	return template
}
