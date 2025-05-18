package utils

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"

	"FMgo/models"
)

// SIIClient representa un cliente para interactuar con el SII
type SIIClient struct {
	BaseURL    string
	HTTPClient *http.Client
	Token      string
}

// NewSIIClient crea una nueva instancia de SIIClient
func NewSIIClient(baseURL string) *SIIClient {
	return &SIIClient{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: time.Second * 30,
		},
	}
}

// Login realiza el inicio de sesión en el SII
func (c *SIIClient) Login(username, password string) error {
	// Implementación del login
	return nil
}

// SendDocument envía un documento al SII
func (c *SIIClient) SendDocument(envio *models.EnvioDTE) (*models.RespuestaSII, error) {
	xmlData, err := xml.Marshal(envio)
	if err != nil {
		return nil, fmt.Errorf("error al serializar el documento: %v", err)
	}

	req, err := http.NewRequest("POST", c.BaseURL+"/envio", bytes.NewBuffer(xmlData))
	if err != nil {
		return nil, fmt.Errorf("error al crear la petición: %v", err)
	}

	req.Header.Set("Content-Type", "application/xml")
	req.Header.Set("Authorization", "Bearer "+c.Token)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error al enviar la petición: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error al leer la respuesta: %v", err)
	}

	var respuesta models.RespuestaSII
	if err := xml.Unmarshal(body, &respuesta); err != nil {
		return nil, fmt.Errorf("error al deserializar la respuesta: %v", err)
	}

	return &respuesta, nil
}

// GetDocumentStatus obtiene el estado de un documento
func (c *SIIClient) GetDocumentStatus(trackID string) (*models.EstadoDocumento, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/estado/%s", c.BaseURL, trackID), nil)
	if err != nil {
		return nil, fmt.Errorf("error al crear la petición: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.Token)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error al enviar la petición: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error al leer la respuesta: %v", err)
	}

	var estado models.EstadoDocumento
	if err := xml.Unmarshal(body, &estado); err != nil {
		return nil, fmt.Errorf("error al deserializar la respuesta: %v", err)
	}

	return &estado, nil
}

// GetContributorInfo obtiene la información de un contribuyente
func (c *SIIClient) GetContributorInfo(rut string) (*models.InformacionContribuyente, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/contribuyente/%s", c.BaseURL, rut), nil)
	if err != nil {
		return nil, fmt.Errorf("error al crear la petición: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.Token)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error al enviar la petición: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error al leer la respuesta: %v", err)
	}

	var info models.InformacionContribuyente
	if err := xml.Unmarshal(body, &info); err != nil {
		return nil, fmt.Errorf("error al deserializar la respuesta: %v", err)
	}

	return &info, nil
}

// GetContributorStatus obtiene el estado de un contribuyente
func (c *SIIClient) GetContributorStatus(rut string) (*models.EstadoContribuyente, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/estado/%s", c.BaseURL, rut), nil)
	if err != nil {
		return nil, fmt.Errorf("error al crear la petición: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.Token)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error al enviar la petición: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error al leer la respuesta: %v", err)
	}

	var estado models.EstadoContribuyente
	if err := xml.Unmarshal(body, &estado); err != nil {
		return nil, fmt.Errorf("error al deserializar la respuesta: %v", err)
	}

	return &estado, nil
}

// GetContributorSummary obtiene el resumen de un contribuyente
func (c *SIIClient) GetContributorSummary(rut, periodo string) (*models.ResumenContribuyente, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/resumen/%s/%s", c.BaseURL, rut, periodo), nil)
	if err != nil {
		return nil, fmt.Errorf("error al crear la petición: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.Token)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error al enviar la petición: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error al leer la respuesta: %v", err)
	}

	var resumen models.ResumenContribuyente
	if err := xml.Unmarshal(body, &resumen); err != nil {
		return nil, fmt.Errorf("error al deserializar la respuesta: %v", err)
	}

	return &resumen, nil
}
