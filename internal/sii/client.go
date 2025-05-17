package sii

import (
	"crypto/tls"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"software.sslmate.com/src/go-pkcs12"
)

const (
	// Clave del certificado PFX para ambiente de pruebas
	certPassword = "83559705FM"

	// Configuración de reintentos
	maxRetries    = 3
	retryInterval = 2 * time.Second
)

// Errores específicos
var (
	ErrNotAuthenticated    = errors.New("cliente no autenticado")
	ErrInvalidXML          = errors.New("XML inválido")
	ErrCertificateNotFound = errors.New("certificado no encontrado")
	ErrConnectionFailed    = errors.New("error de conexión con SII")
	ErrInvalidResponse     = errors.New("respuesta inválida del SII")
)

// RetryConfig contiene la configuración para reintentos
type RetryConfig struct {
	MaxRetries    int
	RetryInterval time.Duration
}

// Client representa un cliente para interactuar con el SII
type Client struct {
	certPath    string
	httpClient  *http.Client
	baseURL     string
	token       string
	retryConfig RetryConfig
}

// NewClient crea una nueva instancia del cliente SII
func NewClient(certPath string) (*Client, error) {
	// Verificar que el archivo exista
	if _, err := os.Stat(certPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("%w: %s", ErrCertificateNotFound, certPath)
	}

	// Leer archivo PFX
	pfxData, err := ioutil.ReadFile(certPath)
	if err != nil {
		return nil, fmt.Errorf("error leyendo certificado PFX: %v", err)
	}

	// Extraer certificado y clave privada del PFX usando la clave proporcionada
	privateKey, cert, err := pkcs12.Decode(pfxData, strings.TrimSpace(certPassword))
	if err != nil {
		return nil, fmt.Errorf("error decodificando PFX: %v", err)
	}

	// Crear certificado TLS
	tlsCert := tls.Certificate{
		Certificate: [][]byte{cert.Raw},
		PrivateKey:  privateKey,
		Leaf:        cert,
	}

	// Configurar cliente HTTP con TLS y timeouts
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
		Timeout: 30 * time.Second,
	}

	return &Client{
		certPath:   certPath,
		httpClient: client,
		baseURL:    "https://maullin.sii.cl/DTEWS/",
		retryConfig: RetryConfig{
			MaxRetries:    maxRetries,
			RetryInterval: retryInterval,
		},
	}, nil
}

// withRetry ejecuta una función con reintentos
func (c *Client) withRetry(operation func() error) error {
	var lastErr error
	for attempt := 0; attempt < c.retryConfig.MaxRetries; attempt++ {
		if err := operation(); err != nil {
			lastErr = err
			// Si es un error que no debemos reintentar, retornamos inmediatamente
			if errors.Is(err, ErrInvalidXML) || errors.Is(err, ErrNotAuthenticated) {
				return err
			}
			// Esperar antes del siguiente intento
			time.Sleep(c.retryConfig.RetryInterval)
			continue
		}
		return nil
	}
	return fmt.Errorf("después de %d intentos: %v", c.retryConfig.MaxRetries, lastErr)
}

// Authenticate obtiene un token de sesión del SII
func (c *Client) Authenticate() error {
	if c.httpClient == nil {
		return fmt.Errorf("%w: cliente no inicializado", ErrConnectionFailed)
	}

	return c.withRetry(func() error {
		// Simular autenticación en modo prueba
		c.token = "TEST_TOKEN"
		return nil
	})
}

// SendDTE envía un DTE al SII
func (c *Client) SendDTE(dteXML string) (string, error) {
	if c.token == "" {
		return "", ErrNotAuthenticated
	}

	// Validar XML
	if err := xml.Unmarshal([]byte(dteXML), new(interface{})); err != nil {
		return "", fmt.Errorf("%w: %v", ErrInvalidXML, err)
	}

	var trackID string
	err := c.withRetry(func() error {
		// Simular envío en modo prueba
		trackID = "123456789"
		return nil
	})

	return trackID, err
}

// ValidateResponse valida la respuesta del SII
func (c *Client) ValidateResponse(response []byte) error {
	if len(response) == 0 {
		return ErrInvalidResponse
	}

	// Validar estructura de la respuesta
	var result struct {
		Status  string `xml:"status"`
		Message string `xml:"message"`
	}

	if err := xml.Unmarshal(response, &result); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidResponse, err)
	}

	if result.Status != "OK" {
		return fmt.Errorf("%w: %s", ErrInvalidResponse, result.Message)
	}

	return nil
}

// SetRetryConfig configura los parámetros de reintento
func (c *Client) SetRetryConfig(config RetryConfig) {
	if config.MaxRetries > 0 {
		c.retryConfig.MaxRetries = config.MaxRetries
	}
	if config.RetryInterval > 0 {
		c.retryConfig.RetryInterval = config.RetryInterval
	}
}
