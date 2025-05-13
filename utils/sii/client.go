package sii

import (
	"bytes"
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/cursor/FMgo/models"
)

// SIIClient maneja la comunicación con el SII
type SIIClient struct {
	certPath     string
	certPassword string
	client       *http.Client
}

// NewSIIClient crea una nueva instancia de SIIClient
func NewSIIClient(certPath, certPassword string) (*SIIClient, error) {
	// Cargar certificado
	cert, err := tls.LoadX509KeyPair(certPath, certPath)
	if err != nil {
		return nil, fmt.Errorf("error al cargar certificado: %v", err)
	}

	// Configurar cliente HTTP
	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	transport := &http.Transport{
		TLSClientConfig: config,
	}
	client := &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}

	return &SIIClient{
		certPath:     certPath,
		certPassword: certPassword,
		client:       client,
	}, nil
}

// FirmarDTE firma un documento DTE
func (c *SIIClient) FirmarDTE(xml []byte) ([]byte, error) {
	// TODO: Implementar firma del documento
	return xml, nil
}

// EnviarDTE envía un documento DTE al SII
func (c *SIIClient) EnviarDTE(sobre models.Sobre) (*models.RespuestaSII, error) {
	// Crear request
	req, err := http.NewRequest("POST", "https://palena.sii.cl/DTEWS/CrSeed.jws", bytes.NewReader(sobre.Documento))
	if err != nil {
		return nil, fmt.Errorf("error al crear request: %v", err)
	}

	// Agregar headers
	req.Header.Set("Content-Type", "application/xml")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", sobre.Token))

	// Enviar request
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error al enviar request: %v", err)
	}
	defer resp.Body.Close()

	// Leer respuesta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error al leer respuesta: %v", err)
	}

	// Decodificar respuesta
	var respuesta models.RespuestaSII
	err = xml.Unmarshal(body, &respuesta)
	if err != nil {
		return nil, fmt.Errorf("error al decodificar respuesta: %v", err)
	}

	return &respuesta, nil
}

// ConsultarEstado consulta el estado de un documento DTE en el SII
func (c *SIIClient) ConsultarEstado(trackID string) (*models.RespuestaSII, error) {
	// Crear request
	req, err := http.NewRequest("GET", fmt.Sprintf("https://palena.sii.cl/DTEWS/GetTokenFromSeed.jws?seed=%s", trackID), nil)
	if err != nil {
		return nil, fmt.Errorf("error al crear request: %v", err)
	}

	// Enviar request
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error al enviar request: %v", err)
	}
	defer resp.Body.Close()

	// Leer respuesta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error al leer respuesta: %v", err)
	}

	// Decodificar respuesta
	var respuesta models.RespuestaSII
	err = xml.Unmarshal(body, &respuesta)
	if err != nil {
		return nil, fmt.Errorf("error al decodificar respuesta: %v", err)
	}

	return &respuesta, nil
}
