package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"FMgo/core/sii/models/siimodels"
)

// SOAPClient maneja las peticiones SOAP al SII
type SOAPClient struct {
	httpClient *http.Client
	config     *siimodels.ConfigSII
}

// NewSOAPClient crea una nueva instancia del cliente SOAP
func NewSOAPClient(config *siimodels.ConfigSII) (*SOAPClient, error) {
	// Cargar certificados TLS
	cert, err := tls.LoadX509KeyPair(config.CertPath, config.KeyPath)
	if err != nil {
		return nil, fmt.Errorf("error al cargar certificados: %w", err)
	}

	// Configurar cliente HTTP con TLS
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
		},
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   config.Timeout,
	}

	return &SOAPClient{
		httpClient: client,
		config:     config,
	}, nil
}

// Call realiza una llamada SOAP al SII
func (c *SOAPClient) Call(ctx context.Context, endpoint string, request interface{}, response interface{}) error {
	// Crear sobre SOAP
	envelope := &siimodels.SoapEnvelope{
		XMLNS: "http://schemas.xmlsoap.org/soap/envelope/",
		Body: siimodels.SoapBody{
			Content: request,
		},
	}

	// Serializar a XML
	xmlData, err := xml.MarshalIndent(envelope, "", "  ")
	if err != nil {
		return fmt.Errorf("error al serializar request: %w", err)
	}

	// Preparar URL completa
	url := fmt.Sprintf("%s%s", c.config.BaseURL, endpoint)

	// Crear request HTTP
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(xmlData))
	if err != nil {
		return fmt.Errorf("error al crear request: %w", err)
	}

	// Configurar headers
	req.Header.Set("Content-Type", "text/xml; charset=utf-8")
	req.Header.Set("SOAPAction", "")

	// Agregar token si está presente en el contexto
	if token, ok := ctx.Value("token").(string); ok && token != "" {
		req.Header.Set("Cookie", fmt.Sprintf("TOKEN=%s", token))
	}

	// Realizar request con reintentos
	var resp *http.Response
	for i := 0; i <= c.config.RetryCount; i++ {
		resp, err = c.httpClient.Do(req)
		if err == nil {
			break
		}
		if i < c.config.RetryCount {
			time.Sleep(c.config.RetryDelay)
			continue
		}
		return fmt.Errorf("error después de %d intentos: %w", c.config.RetryCount, err)
	}
	defer resp.Body.Close()

	// Leer respuesta
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error al leer respuesta: %w", err)
	}

	// Verificar código de estado
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error del servidor: %s - %s", resp.Status, string(body))
	}

	// Deserializar respuesta
	respEnvelope := &siimodels.SoapEnvelope{
		Body: siimodels.SoapBody{
			Content: response,
		},
	}

	if err := xml.Unmarshal(body, respEnvelope); err != nil {
		return fmt.Errorf("error al deserializar respuesta: %w", err)
	}

	// Verificar si hay error SOAP
	if respEnvelope.Body.Fault != nil {
		return fmt.Errorf("error SOAP: %s - %s",
			respEnvelope.Body.Fault.FaultCode,
			respEnvelope.Body.Fault.FaultString)
	}

	return nil
}
