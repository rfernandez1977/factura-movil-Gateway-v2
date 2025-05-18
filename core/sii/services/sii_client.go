package services

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"FMgo/core/sii/models"
)

const (
	// URLs de certificación
	urlBaseCertificacion = "https://maullin.sii.cl"
	urlBaseProduccion    = "https://palena.sii.cl"

	// Endpoints
	endpointSemilla  = "/DTEWS/CrSeed.jws"
	endpointToken    = "/DTEWS/GetTokenFromSeed.jws"
	endpointEnvio    = "/cgi_dte/UPL/DTEUpload"
	endpointConsulta = "/DTEWS/QueryEstDte.jws"
)

// ClienteSII implementa la interfaz para comunicación con el SII
type ClienteSII struct {
	httpClient  *http.Client
	urlBase     string
	certificado *models.CertificadoDigital
	token       string
	tokenExp    time.Time
}

// Config contiene la configuración para el cliente SII
type Config struct {
	Ambiente    string // "certificacion" o "produccion"
	Certificado *models.CertificadoDigital
	Timeout     time.Duration
}

// NewClienteSII crea una nueva instancia del cliente SII
func NewClienteSII(config *Config) *ClienteSII {
	urlBase := urlBaseCertificacion
	if config.Ambiente == "produccion" {
		urlBase = urlBaseProduccion
	}

	return &ClienteSII{
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
		urlBase:     urlBase,
		certificado: config.Certificado,
	}
}

// ObtenerSemilla obtiene una semilla del SII
func (c *ClienteSII) ObtenerSemilla() (string, error) {
	soapRequest := `<?xml version="1.0" encoding="UTF-8"?>
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/">
   <soapenv:Header/>
   <soapenv:Body>
      <getSeed xmlns="http://ws.sii.cl"/>
   </soapenv:Body>
</soapenv:Envelope>`

	resp, err := c.httpClient.Post(
		c.urlBase+endpointSemilla,
		"text/xml;charset=UTF-8",
		bytes.NewBufferString(soapRequest),
	)
	if err != nil {
		return "", fmt.Errorf("error al obtener semilla: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error leyendo respuesta: %v", err)
	}

	var soapResponse struct {
		XMLName xml.Name `xml:"Envelope"`
		Body    struct {
			GetSeedResponse struct {
				Seed string `xml:"seed"`
			} `xml:"getSeedResponse"`
		} `xml:"Body"`
	}

	if err := xml.Unmarshal(body, &soapResponse); err != nil {
		return "", fmt.Errorf("error decodificando respuesta: %v", err)
	}

	return soapResponse.Body.GetSeedResponse.Seed, nil
}

// ObtenerToken obtiene un token de autenticación
func (c *ClienteSII) ObtenerToken() error {
	// Verificar si el token actual es válido
	if c.token != "" && time.Now().Before(c.tokenExp) {
		return nil
	}

	// Obtener semilla
	semilla, err := c.ObtenerSemilla()
	if err != nil {
		return fmt.Errorf("error obteniendo semilla: %v", err)
	}

	// Crear solicitud token
	soapRequest := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/">
   <soapenv:Header/>
   <soapenv:Body>
      <getToken xmlns="http://ws.sii.cl">
         <seed>%s</seed>
      </getToken>
   </soapenv:Body>
</soapenv:Envelope>`, semilla)

	resp, err := c.httpClient.Post(
		c.urlBase+endpointToken,
		"text/xml;charset=UTF-8",
		bytes.NewBufferString(soapRequest),
	)
	if err != nil {
		return fmt.Errorf("error obteniendo token: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error leyendo respuesta: %v", err)
	}

	var soapResponse struct {
		XMLName xml.Name `xml:"Envelope"`
		Body    struct {
			GetTokenResponse struct {
				Token string `xml:"token"`
			} `xml:"getTokenResponse"`
		} `xml:"Body"`
	}

	if err := xml.Unmarshal(body, &soapResponse); err != nil {
		return fmt.Errorf("error decodificando respuesta: %v", err)
	}

	c.token = soapResponse.Body.GetTokenResponse.Token
	c.tokenExp = time.Now().Add(1 * time.Hour)
	return nil
}

// EnviarDocumento envía un documento al SII
func (c *ClienteSII) EnviarDocumento(xmlData []byte) (*models.RespuestaEnvio, error) {
	if err := c.ObtenerToken(); err != nil {
		return nil, fmt.Errorf("error obteniendo token: %v", err)
	}

	req, err := http.NewRequest("POST", c.urlBase+endpointEnvio, bytes.NewReader(xmlData))
	if err != nil {
		return nil, fmt.Errorf("error creando solicitud: %v", err)
	}

	req.Header.Set("Content-Type", "text/xml;charset=UTF-8")
	req.Header.Set("Cookie", fmt.Sprintf("TOKEN=%s", c.token))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error enviando documento: %v", err)
	}
	defer resp.Body.Close()

	var respuesta models.RespuestaEnvio
	if err := xml.NewDecoder(resp.Body).Decode(&respuesta); err != nil {
		return nil, fmt.Errorf("error decodificando respuesta: %v", err)
	}

	return &respuesta, nil
}

// ConsultarEstado consulta el estado de un envío
func (c *ClienteSII) ConsultarEstado(trackID string) (*models.EstadoEnvio, error) {
	if err := c.ObtenerToken(); err != nil {
		return nil, fmt.Errorf("error obteniendo token: %v", err)
	}

	soapRequest := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/">
   <soapenv:Header/>
   <soapenv:Body>
      <getEstDte xmlns="http://ws.sii.cl">
         <trackId>%s</trackId>
      </getEstDte>
   </soapenv:Body>
</soapenv:Envelope>`, trackID)

	req, err := http.NewRequest("POST", c.urlBase+endpointConsulta, bytes.NewBufferString(soapRequest))
	if err != nil {
		return nil, fmt.Errorf("error creando solicitud: %v", err)
	}

	req.Header.Set("Content-Type", "text/xml;charset=UTF-8")
	req.Header.Set("Cookie", fmt.Sprintf("TOKEN=%s", c.token))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error consultando estado: %v", err)
	}
	defer resp.Body.Close()

	var estado models.EstadoEnvio
	if err := xml.NewDecoder(resp.Body).Decode(&estado); err != nil {
		return nil, fmt.Errorf("error decodificando respuesta: %v", err)
	}

	return &estado, nil
}
