package sii

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"

	"FMgo/models"
)

// HTTPClient define la interfaz para el cliente HTTP
type HTTPClient interface {
	Get(ctx context.Context, url string) ([]byte, error)
	Post(ctx context.Context, url string, body []byte) ([]byte, error)
}

// SIIClient implementa la interfaz para el cliente SII
type SIIClient struct {
	httpClient HTTPClient
	baseURL    string
}

// NewSIIClient crea una nueva instancia del cliente SII
func NewSIIClient(httpClient HTTPClient, baseURL string) *SIIClient {
	return &SIIClient{
		httpClient: httpClient,
		baseURL:    baseURL,
	}
}

// EnviarDTE envía un DTE al SII
func (c *SIIClient) EnviarDTE(dte *models.DTEXMLModel) (*models.RespuestaSII, error) {
	if dte == nil {
		return nil, fmt.Errorf("el DTE no puede ser nulo")
	}

	// Convertir el DTE a XML
	xmlData, err := xml.Marshal(dte)
	if err != nil {
		return nil, fmt.Errorf("error al convertir DTE a XML: %v", err)
	}

	// Enviar el DTE al SII
	resp, err := c.httpClient.Post(context.Background(), fmt.Sprintf("%s/api/dte", c.baseURL), xmlData)
	if err != nil {
		return nil, fmt.Errorf("error al enviar DTE al SII: %v", err)
	}

	// Procesar la respuesta
	var respuesta models.RespuestaSII
	if err := xml.Unmarshal(resp, &respuesta); err != nil {
		return nil, fmt.Errorf("error al procesar respuesta del SII: %v", err)
	}

	return &respuesta, nil
}

// ConsultarEstado consulta el estado de un DTE en el SII
func (c *SIIClient) ConsultarEstado(trackID string) (*models.EstadoSII, error) {
	if trackID == "" {
		return nil, fmt.Errorf("trackID es requerido")
	}

	// Consultar estado al SII
	resp, err := c.httpClient.Get(context.Background(), fmt.Sprintf("%s/api/dte/estado/%s", c.baseURL, trackID))
	if err != nil {
		return nil, fmt.Errorf("error al consultar estado al SII: %v", err)
	}

	// Procesar la respuesta
	var estado models.EstadoSII
	if err := xml.Unmarshal(resp, &estado); err != nil {
		return nil, fmt.Errorf("error al procesar estado del SII: %v", err)
	}

	return &estado, nil
}

// ConsultarDTE consulta un DTE específico en el SII
func (c *SIIClient) ConsultarDTE(tipoDTE, folio, rutEmisor string) (*models.EstadoSII, error) {
	if tipoDTE == "" || folio == "" || rutEmisor == "" {
		return nil, fmt.Errorf("tipoDTE, folio y rutEmisor son requeridos")
	}

	// Consultar DTE al SII
	resp, err := c.httpClient.Get(context.Background(), fmt.Sprintf("%s/api/dte/consulta/%s/%s/%s", c.baseURL, tipoDTE, folio, rutEmisor))
	if err != nil {
		return nil, fmt.Errorf("error al consultar DTE al SII: %v", err)
	}

	// Procesar la respuesta
	var estado models.EstadoSII
	if err := xml.Unmarshal(resp, &estado); err != nil {
		return nil, fmt.Errorf("error al procesar estado del DTE: %v", err)
	}

	return &estado, nil
}

// VerificarComunicacion verifica la comunicación con el SII
func (c *SIIClient) VerificarComunicacion() error {
	// Verificar comunicación con el SII
	resp, err := c.httpClient.Get(context.Background(), fmt.Sprintf("%s/api/estado", c.baseURL))
	if err != nil {
		return fmt.Errorf("error al verificar comunicación con el SII: %v", err)
	}

	// Procesar la respuesta
	var estado models.EstadoSII
	if err := xml.Unmarshal(resp, &estado); err != nil {
		return fmt.Errorf("error al procesar estado del SII: %v", err)
	}

	if models.EstadoSIIType(estado.Estado) != models.EstadoSIIAceptado {
		return fmt.Errorf("error de comunicación con el SII: %s", estado.Glosa)
	}

	return nil
}

// Client representa un cliente para interactuar con el SII
type Client struct {
	baseURL string
	client  *http.Client
}

// NewClient crea una nueva instancia del cliente SII
func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// ConsultarEstado consulta el estado de un DTE
func (c *Client) ConsultarEstado(tipoDTE int, folio int, rutEmisor string) (*models.EstadoSII, error) {
	// Crear request
	url := fmt.Sprintf("%s/ConsultaEstado?tipoDTE=%d&folio=%d&rutEmisor=%s",
		c.baseURL, tipoDTE, folio, rutEmisor)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error al crear request: %w", err)
	}

	// Configurar headers
	req.Header.Set("Accept", "application/xml")

	// Enviar request
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error al enviar request: %w", err)
	}
	defer resp.Body.Close()

	// Leer respuesta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error al leer respuesta: %w", err)
	}

	// Verificar status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error en respuesta del SII: %s", string(body))
	}

	// Decodificar respuesta
	var estado models.EstadoSII
	err = xml.Unmarshal(body, &estado)
	if err != nil {
		return nil, fmt.Errorf("error al decodificar respuesta: %w", err)
	}

	return &estado, nil
}

// EstadoDTE representa el estado de un DTE
type EstadoDTE struct {
	XMLName      xml.Name `xml:"Estado"`
	Version      string   `xml:"version,attr"`
	Estado       string   `xml:"Estado"`
	Glosa        string   `xml:"Glosa"`
	TrackID      string   `xml:"TrackID,omitempty"`
	Errores      []Error  `xml:"Errores>Error,omitempty"`
	Advertencias []Error  `xml:"Advertencias>Error,omitempty"`
}

// Error representa un error o advertencia del SII
type Error struct {
	XMLName xml.Name `xml:"Error"`
	Codigo  string   `xml:"Codigo"`
	Glosa   string   `xml:"Glosa"`
}
