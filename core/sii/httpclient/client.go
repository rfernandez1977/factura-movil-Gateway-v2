package httpclient

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"

	"FMgo/core/sii/infrastructure/certificates"
	"FMgo/core/sii/models/siimodels"
)

// Client implementa el cliente HTTP para el SII
type Client struct {
	httpClient  *http.Client
	certManager *certificates.Manager
	config      *siimodels.Config
}

// New crea una nueva instancia del cliente HTTP
func New(config *siimodels.Config) (*Client, error) {
	if err := validateConfig(config); err != nil {
		return nil, fmt.Errorf("configuración inválida: %w", err)
	}

	// Configurar el manejador de certificados
	certManager, err := certificates.New(config.CertPath, config.KeyPath)
	if err != nil {
		return nil, fmt.Errorf("error configurando certificados: %w", err)
	}

	// Configurar cliente HTTP
	client := &http.Client{
		Timeout: config.Timeout,
	}

	return &Client{
		httpClient:  client,
		certManager: certManager,
		config:      config,
	}, nil
}

// validateConfig valida la configuración
func validateConfig(config *siimodels.Config) error {
	if config == nil {
		return fmt.Errorf("la configuración no puede ser nula")
	}
	if config.CertPath == "" {
		return fmt.Errorf("la ruta del certificado es requerida")
	}
	if config.KeyPath == "" {
		return fmt.Errorf("la ruta de la llave privada es requerida")
	}
	if config.Timeout < 1*time.Second {
		return fmt.Errorf("el timeout debe ser al menos 1 segundo")
	}
	if config.RetryCount < 0 {
		return fmt.Errorf("el número de reintentos no puede ser negativo")
	}
	return nil
}

// doRequest ejecuta una petición HTTP con reintentos
func (c *Client) doRequest(ctx context.Context, req *http.Request) (*http.Response, error) {
	var lastErr error
	for attempt := 0; attempt <= c.config.RetryCount; attempt++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			// Si no es el primer intento, esperar antes de reintentar
			if attempt > 0 {
				time.Sleep(c.config.RetryDelay)
			}

			resp, err := c.httpClient.Do(req)
			if err != nil {
				lastErr = err
				continue
			}

			// Verificar si el error es recuperable
			if isRetryableError(resp.StatusCode) && attempt < c.config.RetryCount {
				resp.Body.Close()
				continue
			}

			return resp, nil
		}
	}

	return nil, fmt.Errorf("máximo número de reintentos alcanzado: %w", lastErr)
}

// Post realiza una petición POST
func (c *Client) Post(ctx context.Context, url string, contentType string, body []byte) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("error creando request: %w", err)
	}

	req.Header.Set("Content-Type", contentType)
	return c.doRequest(ctx, req)
}

// Get realiza una petición GET
func (c *Client) Get(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creando request: %w", err)
	}

	return c.doRequest(ctx, req)
}

// SendSOAP envía una petición SOAP
func (c *Client) SendSOAP(ctx context.Context, url string, envelope *siimodels.SoapEnvelope) (*siimodels.SoapEnvelope, error) {
	// Serializar el envelope
	xmlData, err := xml.MarshalIndent(envelope, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error serializando envelope: %w", err)
	}

	// Realizar la petición
	resp, err := c.Post(ctx, url, "text/xml; charset=utf-8", xmlData)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Leer la respuesta
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error leyendo respuesta: %w", err)
	}

	// Decodificar la respuesta
	var respEnvelope siimodels.SoapEnvelope
	if err := xml.Unmarshal(respBody, &respEnvelope); err != nil {
		return nil, fmt.Errorf("error decodificando respuesta: %w", err)
	}

	// Verificar si hay error SOAP
	if respEnvelope.Body.Fault != nil {
		return nil, fmt.Errorf("error SOAP: %s - %s",
			respEnvelope.Body.Fault.FaultCode,
			respEnvelope.Body.Fault.FaultString)
	}

	return &respEnvelope, nil
}

// isRetryableError determina si un código de estado HTTP es recuperable
func isRetryableError(statusCode int) bool {
	return statusCode >= 500 || statusCode == 429
}
