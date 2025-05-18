package client

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"sync"
	"time"

	"FMgo/core/sii/models/siimodels"
)

// AuthClient maneja la autenticación con el SII
type AuthClient struct {
	soapClient *SOAPClient
	config     *siimodels.ConfigSII
	token      string
	tokenMutex sync.RWMutex
	tokenExp   time.Time
}

// NewAuthClient crea una nueva instancia del cliente de autenticación
func NewAuthClient(soapClient *SOAPClient, config *siimodels.ConfigSII) *AuthClient {
	return &AuthClient{
		soapClient: soapClient,
		config:     config,
	}
}

// GetToken obtiene un token válido, renovándolo si es necesario
func (c *AuthClient) GetToken(ctx context.Context) (string, error) {
	c.tokenMutex.RLock()
	if c.token != "" && time.Now().Before(c.tokenExp) {
		token := c.token
		c.tokenMutex.RUnlock()
		return token, nil
	}
	c.tokenMutex.RUnlock()

	// Necesitamos un nuevo token
	c.tokenMutex.Lock()
	defer c.tokenMutex.Unlock()

	// Verificar nuevamente en caso de que otro goroutine haya actualizado el token
	if c.token != "" && time.Now().Before(c.tokenExp) {
		return c.token, nil
	}

	// Obtener semilla
	seed, err := c.getSeed(ctx)
	if err != nil {
		return "", fmt.Errorf("error al obtener semilla: %w", err)
	}

	// Obtener token
	token, err := c.getToken(ctx, seed)
	if err != nil {
		return "", fmt.Errorf("error al obtener token: %w", err)
	}

	// Almacenar el nuevo token
	c.token = token
	c.tokenExp = time.Now().Add(1 * time.Hour) // Los tokens del SII expiran después de 1 hora

	return token, nil
}

// getSeed obtiene una semilla del SII
func (c *AuthClient) getSeed(ctx context.Context) (string, error) {
	request := struct {
		XMLName struct{} `xml:"getSeed"`
	}{}

	response := &siimodels.RespuestaSII{}

	err := c.soapClient.Call(ctx, siimodels.EndpointSemillaCert, request, response)
	if err != nil {
		return "", err
	}

	if response.Header.Estado != "00" {
		return "", fmt.Errorf("error al obtener semilla: %s", response.Header.Glosa)
	}

	return response.Body.Semilla, nil
}

// getToken obtiene un token usando la semilla
func (c *AuthClient) getToken(ctx context.Context, seed string) (string, error) {
	request := struct {
		XMLName struct{} `xml:"getToken"`
		Seed    string   `xml:"seed"`
	}{
		Seed: seed,
	}

	response := &siimodels.RespuestaSII{}

	err := c.soapClient.Call(ctx, siimodels.EndpointTokenCert, request, response)
	if err != nil {
		return "", err
	}

	if response.Header.Estado != "00" {
		return "", fmt.Errorf("error al obtener token: %s", response.Header.Glosa)
	}

	return response.Body.Token, nil
}

// IniciarSesion inicia una sesión con el SII
func (c *AuthClient) IniciarSesion(ctx context.Context) (*siimodels.RespuestaSesion, error) {
	// Leer el certificado
	certData, err := ioutil.ReadFile(c.config.CertPath)
	if err != nil {
		return nil, fmt.Errorf("error al leer certificado: %w", err)
	}

	// Decodificar el certificado
	block, _ := pem.Decode(certData)
	if block == nil {
		return nil, fmt.Errorf("error al decodificar certificado PEM")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("error al parsear certificado: %w", err)
	}

	// Crear la solicitud de sesión
	certificadoInfo := &siimodels.CertificadoDigital{
		SerialNumber: cert.SerialNumber.String(),
		Issuer:       cert.Issuer.String(),
		ValidFrom:    cert.NotBefore,
		ValidTo:      cert.NotAfter,
	}

	solicitud := siimodels.NewSolicitudSesion(
		c.config.RutEmpresa,
		certificadoInfo,
		c.config.Ambiente,
	)

	// Preparar la respuesta
	respuesta := &siimodels.RespuestaSesion{}

	// Enviar la solicitud
	err = c.soapClient.Call(ctx, siimodels.EndpointTokenCert, solicitud, respuesta)
	if err != nil {
		return nil, fmt.Errorf("error al enviar solicitud de sesión: %w", err)
	}

	return respuesta, nil
}

// RenovarSesion renueva una sesión existente
func (c *AuthClient) RenovarSesion(ctx context.Context, token string) (*siimodels.RespuestaSesion, error) {
	solicitud := &siimodels.RenovacionSesion{
		Version:   "1.0",
		Token:     token,
		Timestamp: time.Now(),
	}

	respuesta := &siimodels.RespuestaSesion{}

	err := c.soapClient.Call(ctx, siimodels.EndpointTokenCert, solicitud, respuesta)
	if err != nil {
		return nil, fmt.Errorf("error al renovar sesión: %w", err)
	}

	return respuesta, nil
}

// CerrarSesion cierra una sesión existente
func (c *AuthClient) CerrarSesion(ctx context.Context, token string) error {
	solicitud := &siimodels.CierreSesion{
		Version:   "1.0",
		Token:     token,
		Timestamp: time.Now(),
	}

	respuesta := &siimodels.RespuestaSesion{}

	err := c.soapClient.Call(ctx, siimodels.EndpointTokenCert, solicitud, respuesta)
	if err != nil {
		return fmt.Errorf("error al cerrar sesión: %w", err)
	}

	if respuesta.Estado != siimodels.EstadoOK {
		return fmt.Errorf("error al cerrar sesión: %s", respuesta.Mensaje)
	}

	return nil
}
