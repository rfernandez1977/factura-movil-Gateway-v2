package tests

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	// URLs de prueba del SII
	siiTestUploadURL = "https://maullin.sii.cl/DTEWS/CrSeed.jws"
	siiTestStatusURL = "https://maullin.sii.cl/DTEWS/QueryEstUp.jws"
	siiTestTokenURL  = "https://maullin.sii.cl/DTEWS/GetTokenFromSeed.jws"
)

// SIIClient maneja la comunicación con el SII
type SIIClient struct {
	httpClient *http.Client
	token      string
}

// NewSIIClient crea un nuevo cliente para el SII
func NewSIIClient() *SIIClient {
	return &SIIClient{
		httpClient: &http.Client{},
	}
}

// SendDTE envía un DTE al SII
func (c *SIIClient) SendDTE(xmlData string) (string, error) {
	// Obtener semilla
	seed, err := c.getSeed()
	if err != nil {
		return "", fmt.Errorf("error getting seed: %v", err)
	}

	// Obtener token
	token, err := c.getToken(seed)
	if err != nil {
		return "", fmt.Errorf("error getting token: %v", err)
	}
	c.token = token

	// Preparar el envío
	req, err := http.NewRequest("POST", siiTestUploadURL, bytes.NewBufferString(xmlData))
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	// Agregar headers necesarios
	req.Header.Set("Content-Type", "text/xml; charset=ISO-8859-1")
	req.Header.Set("Cookie", fmt.Sprintf("TOKEN=%s", token))

	// Realizar el envío
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending DTE: %v", err)
	}
	defer resp.Body.Close()

	// Leer la respuesta
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %v", err)
	}

	// Extraer el ID de seguimiento
	trackID, err := extractTrackID(string(body))
	if err != nil {
		return "", fmt.Errorf("error extracting track ID: %v", err)
	}

	return trackID, nil
}

// CheckStatus consulta el estado de un envío
func (c *SIIClient) CheckStatus(trackID string) (string, error) {
	// Preparar la consulta
	req, err := http.NewRequest("POST", siiTestStatusURL, strings.NewReader(fmt.Sprintf(`
		<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/">
			<soapenv:Body>
				<ws:getEstUp xmlns:ws="http://ws.status.dte.sii.cl">
					<rutEmpresa>%s</rutEmpresa>
					<trackId>%s</trackId>
				</ws:getEstUp>
			</soapenv:Body>
		</soapenv:Envelope>
	`, "76212889-6", trackID)))
	if err != nil {
		return "", fmt.Errorf("error creating status request: %v", err)
	}

	// Agregar headers necesarios
	req.Header.Set("Content-Type", "text/xml; charset=ISO-8859-1")
	req.Header.Set("Cookie", fmt.Sprintf("TOKEN=%s", c.token))

	// Realizar la consulta
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error checking status: %v", err)
	}
	defer resp.Body.Close()

	// Leer la respuesta
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading status response: %v", err)
	}

	// Extraer el estado
	status, err := extractStatus(string(body))
	if err != nil {
		return "", fmt.Errorf("error extracting status: %v", err)
	}

	return status, nil
}

// getSeed obtiene una semilla del SII
func (c *SIIClient) getSeed() (string, error) {
	resp, err := c.httpClient.Get(siiTestUploadURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Extraer la semilla de la respuesta
	// Esta es una implementación simulada
	return "SEED123", nil
}

// getToken obtiene un token usando la semilla
func (c *SIIClient) getToken(seed string) (string, error) {
	// Esta es una implementación simulada
	return "TOKEN123", nil
}

// extractTrackID extrae el ID de seguimiento de la respuesta del SII
func extractTrackID(response string) (string, error) {
	// Esta es una implementación simulada
	return "123456789", nil
}

// extractStatus extrae el estado de la respuesta del SII
func extractStatus(response string) (string, error) {
	// Esta es una implementación simulada
	return "RECIBIDO", nil
}
