package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// FacturaMovilClient es el cliente para interactuar con la API de Factura Móvil
type FacturaMovilClient struct {
	baseURL    string
	apiToken   string
	httpClient *http.Client
}

// NewFacturaMovilClient crea una nueva instancia del cliente de Factura Móvil
func NewFacturaMovilClient(baseURL, apiToken string) *FacturaMovilClient {
	return &FacturaMovilClient{
		baseURL:  baseURL,
		apiToken: apiToken,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// CreateProduct crea un nuevo producto en Factura Móvil
func (c *FacturaMovilClient) CreateProduct(product interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(product)
	if err != nil {
		return nil, fmt.Errorf("error al serializar producto: %v", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/productos", c.baseURL), bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error al crear request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiToken))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error al enviar request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		errorBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("error de API (código %d): %s", resp.StatusCode, string(errorBytes))
	}

	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error al leer respuesta: %v", err)
	}

	return responseData, nil
}

// ListProducts obtiene una lista de productos de Factura Móvil
func (c *FacturaMovilClient) ListProducts(params map[string]string) ([]byte, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/productos", c.baseURL), nil)
	if err != nil {
		return nil, fmt.Errorf("error al crear request: %v", err)
	}

	// Agregar parámetros a la URL
	q := req.URL.Query()
	for key, value := range params {
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode()

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiToken))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error al enviar request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		errorBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("error de API (código %d): %s", resp.StatusCode, string(errorBytes))
	}

	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error al leer respuesta: %v", err)
	}

	return responseData, nil
}
