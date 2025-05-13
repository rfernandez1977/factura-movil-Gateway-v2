package api

import (
	"bytes"
	"fmt"
	"net/http"
)

// FacturaMovilClient representa el cliente de Factura Móvil
type FacturaMovilClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewFacturaMovilClient crea una nueva instancia del cliente
func NewFacturaMovilClient(baseURL string) *FacturaMovilClient {
	return &FacturaMovilClient{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{},
	}
}

// SearchClients busca clientes por término
func (c *FacturaMovilClient) SearchClients(searchTerm string) ([]byte, error) {
	url := fmt.Sprintf("%s/api/v1/clients/search?q=%s", c.BaseURL, searchTerm)

	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error al realizar la petición: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("respuesta no exitosa: %d", resp.StatusCode)
	}

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(resp.Body); err != nil {
		return nil, fmt.Errorf("error al leer la respuesta: %v", err)
	}

	return buf.Bytes(), nil
}
