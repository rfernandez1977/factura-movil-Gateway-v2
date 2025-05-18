package sii

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

// SIIClient interface define los métodos para interactuar con el SII
type SIIClient interface {
	GetSemilla(ctx context.Context) (string, error)
	GetToken(ctx context.Context, semillaFirmada string) (string, error)
}

// Config estructura para la configuración del cliente SII
type Config struct {
	BaseURL     string
	Environment string // certificacion o produccion
	Timeout     int    // timeout en segundos
	MaxRetries  int    // número máximo de reintentos
	TestMode    bool   // modo de prueba
}

// DefaultConfig retorna la configuración por defecto
func DefaultConfig() *Config {
	return &Config{
		BaseURL:     "https://palena.sii.cl",
		Environment: "certificacion",
		Timeout:     30,
		MaxRetries:  3,
		TestMode:    false,
	}
}

// Client implementación del cliente SII
type Client struct {
	config     *Config
	httpClient *http.Client
	logger     *log.Logger
}

// NewClient crea una nueva instancia del cliente SII
func NewClient(config *Config) *Client {
	if config == nil {
		config = DefaultConfig()
	}
	return &Client{
		config: config,
		httpClient: &http.Client{
			Timeout: time.Duration(config.Timeout) * time.Second,
		},
		logger: log.New(log.Writer(), "[SII Client] ", log.LstdFlags),
	}
}

// doRequest ejecuta una petición HTTP con reintentos
func (c *Client) doRequest(req *http.Request) (*http.Response, error) {
	if c.config.TestMode {
		c.logger.Printf("Modo de prueba: simulando respuesta para %s %s", req.Method, req.URL.Path)
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(strings.NewReader(c.getMockResponse(req))),
		}, nil
	}

	var lastErr error
	for i := 0; i <= c.config.MaxRetries; i++ {
		c.logger.Printf("Intento %d de %d para %s %s", i+1, c.config.MaxRetries+1, req.Method, req.URL.Path)
		resp, err := c.httpClient.Do(req)
		if err == nil {
			c.logger.Printf("Respuesta exitosa: status %d", resp.StatusCode)
			return resp, nil
		}
		lastErr = err
		c.logger.Printf("Error en intento %d: %v", i+1, err)
		if i < c.config.MaxRetries {
			delay := time.Duration(i+1) * time.Second
			c.logger.Printf("Esperando %v antes del siguiente intento", delay)
			time.Sleep(delay)
		}
	}
	return nil, fmt.Errorf("después de %d intentos: %w", c.config.MaxRetries, lastErr)
}

// getMockResponse retorna una respuesta simulada para pruebas
func (c *Client) getMockResponse(req *http.Request) string {
	if strings.Contains(req.URL.Path, "CrSeed") {
		return `<?xml version="1.0" encoding="UTF-8"?>
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
   <SOAP-ENV:Body>
      <ns1:getSeedResponse xmlns:ns1="http://DefaultNamespace">
         <return>SEMILLA-DE-PRUEBA-1234567890</return>
      </ns1:getSeedResponse>
   </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`
	}
	if strings.Contains(req.URL.Path, "GetTokenFromSeed") {
		bodyBytes, err := ioutil.ReadAll(req.Body)
		if err != nil {
			c.logger.Printf("Error leyendo body del request: %v", err)
			return ""
		}
		// Restaurar el body para futuras lecturas
		req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		if strings.Contains(string(bodyBytes), "SEMILLA-INVALIDA") {
			return `<?xml version="1.0" encoding="UTF-8"?>
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
   <SOAP-ENV:Body>
      <SOAP-ENV:Fault>
         <faultcode>SOAP-ENV:Client</faultcode>
         <faultstring>Semilla inválida</faultstring>
      </SOAP-ENV:Fault>
   </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`
		}
		return `<?xml version="1.0" encoding="UTF-8"?>
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
   <SOAP-ENV:Body>
      <ns1:getTokenResponse xmlns:ns1="http://DefaultNamespace">
         <return>TOKEN-DE-PRUEBA-1234567890</return>
      </ns1:getTokenResponse>
   </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`
	}
	return ""
}

// GetSemilla obtiene una semilla desde el SII
func (c *Client) GetSemilla(ctx context.Context) (string, error) {
	url := fmt.Sprintf("%s/DTEWS/CrSeed.jws?WSDL", c.config.BaseURL)

	c.logger.Printf("Solicitando semilla a %s", url)

	// Crear request SOAP
	soapReq := &Envelope{
		XMLNS: "http://schemas.xmlsoap.org/soap/envelope/",
		Body: Body{
			Content: &SemillaRequest{},
		},
	}

	reqBody, err := xml.MarshalIndent(soapReq, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error creando XML request: %w", err)
	}

	c.logger.Printf("Request XML: %s", string(reqBody))

	req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(string(reqBody)))
	if err != nil {
		return "", fmt.Errorf("error creando request: %w", err)
	}

	req.Header.Set("Content-Type", "text/xml; charset=utf-8")
	req.Header.Set("SOAPAction", "getSeed")

	resp, err := c.doRequest(req)
	if err != nil {
		return "", fmt.Errorf("error en la petición: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error del servidor SII: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error leyendo respuesta: %w", err)
	}

	c.logger.Printf("Respuesta XML: %s", string(body))

	var soapResp Envelope
	decoder := xml.NewDecoder(bytes.NewReader(body))
	decoder.Strict = false
	if err := decoder.Decode(&soapResp); err != nil {
		return "", fmt.Errorf("error parseando respuesta SOAP: %w", err)
	}

	// Verificar si hay error SOAP
	if soapResp.Body.Fault != nil {
		return "", fmt.Errorf("error SOAP: %s - %s", soapResp.Body.Fault.FaultCode, soapResp.Body.Fault.FaultString)
	}

	// Convertir el contenido a SemillaResponse
	respContent, err := xml.Marshal(soapResp.Body.Content)
	if err != nil {
		return "", fmt.Errorf("error procesando respuesta: %w", err)
	}

	var semillaResp SemillaResponse
	decoder = xml.NewDecoder(bytes.NewReader(respContent))
	decoder.Strict = false
	if err := decoder.Decode(&semillaResp); err != nil {
		return "", fmt.Errorf("error parseando respuesta semilla: %w", err)
	}

	if semillaResp.Return == "" {
		return "", fmt.Errorf("semilla vacía en respuesta")
	}

	c.logger.Printf("Semilla obtenida exitosamente")
	return semillaResp.Return, nil
}

// GetToken obtiene un token usando una semilla firmada
func (c *Client) GetToken(ctx context.Context, semillaFirmada string) (string, error) {
	url := fmt.Sprintf("%s/DTEWS/GetTokenFromSeed.jws?WSDL", c.config.BaseURL)

	c.logger.Printf("Solicitando token a %s", url)

	// Crear request SOAP
	soapReq := &Envelope{
		XMLNS: "http://schemas.xmlsoap.org/soap/envelope/",
		Body: Body{
			Content: &TokenRequest{
				Token: semillaFirmada,
			},
		},
	}

	reqBody, err := xml.MarshalIndent(soapReq, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error creando XML request: %w", err)
	}

	c.logger.Printf("Request XML: %s", string(reqBody))

	req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(string(reqBody)))
	if err != nil {
		return "", fmt.Errorf("error creando request: %w", err)
	}

	req.Header.Set("Content-Type", "text/xml; charset=utf-8")
	req.Header.Set("SOAPAction", "getToken")

	resp, err := c.doRequest(req)
	if err != nil {
		return "", fmt.Errorf("error en la petición: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error del servidor SII: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error leyendo respuesta: %w", err)
	}

	c.logger.Printf("Respuesta XML: %s", string(body))

	var soapResp Envelope
	decoder := xml.NewDecoder(bytes.NewReader(body))
	decoder.Strict = false
	if err := decoder.Decode(&soapResp); err != nil {
		return "", fmt.Errorf("error parseando respuesta SOAP: %w", err)
	}

	// Verificar si hay error SOAP
	if soapResp.Body.Fault != nil {
		return "", fmt.Errorf("error SOAP: %s - %s", soapResp.Body.Fault.FaultCode, soapResp.Body.Fault.FaultString)
	}

	// Convertir el contenido a TokenResponse
	respContent, err := xml.Marshal(soapResp.Body.Content)
	if err != nil {
		return "", fmt.Errorf("error procesando respuesta: %w", err)
	}

	var tokenResp TokenResponse
	decoder = xml.NewDecoder(bytes.NewReader(respContent))
	decoder.Strict = false
	if err := decoder.Decode(&tokenResp); err != nil {
		return "", fmt.Errorf("error parseando respuesta token: %w", err)
	}

	if tokenResp.Return == "" {
		return "", fmt.Errorf("token vacío en respuesta")
	}

	c.logger.Printf("Token obtenido exitosamente")
	return tokenResp.Return, nil
}
