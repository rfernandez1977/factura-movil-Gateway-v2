package services

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/fmgo/models"
)

const (
	// Endpoints de certificación
	SIISemillaEndpoint = "https://maullin.sii.cl/DTEWS/CrSeed.jws?wsdl"
	SIITokenEndpoint   = "https://maullin.sii.cl/DTEWS/GetTokenFromSeed.jws?wsdl"
	SIIEnvioEndpoint   = "https://maullin.sii.cl/cgi_dte/UPL/DTEUpload"
	SIIEstadoEndpoint  = "https://maullin.sii.cl/DTEWS/QueryEstDte.jws?wsdl"

	// Endpoints de producción (comentados por ahora)
	// SIISemillaEndpoint = "https://palena.sii.cl/DTEWS/CrSeed.jws?wsdl"
	// SIITokenEndpoint   = "https://palena.sii.cl/DTEWS/GetTokenFromSeed.jws?wsdl"
	// SIIEnvioEndpoint   = "https://palena.sii.cl/cgi_dte/UPL/DTEUpload"
	// SIIEstadoEndpoint  = "https://palena.sii.cl/DTEWS/QueryEstDte.jws?wsdl"
)

// SIIClientInterface define la interfaz común para clientes SII
type SIIClientInterface interface {
	ObtenerSemilla() (string, error)
	ObtenerToken(semilla string) (string, error)
	EnviarDTE(sobre *models.Sobre, token string) error
}

// SIIClient maneja la comunicación con el SII
type SIIClient struct {
	client     *http.Client
	firma      *FirmaManager
	simulacion bool
}

// MockSIIClient es un cliente mock para pruebas
type MockSIIClient struct{}

func (m *MockSIIClient) ObtenerSemilla() (string, error) {
	return "SEMILLA_TEST", nil
}

func (m *MockSIIClient) ObtenerToken(semilla string) (string, error) {
	return "TOKEN_TEST", nil
}

func (m *MockSIIClient) EnviarDTE(sobre *models.Sobre, token string) error {
	return nil
}

// NewMockSIIClient crea un nuevo cliente mock para el SII
func NewMockSIIClient() SIIClientInterface {
	return &MockSIIClient{}
}

// NewSIIClient crea un nuevo cliente para el SII
func NewSIIClient(rutaCertificado, contraseña string) (SIIClientInterface, error) {
	// Crear manejador de firma
	firma, err := NewFirmaManager(rutaCertificado, contraseña)
	if err != nil {
		return nil, fmt.Errorf("error creando manejador de firma: %v", err)
	}

	client := &http.Client{
		Timeout: time.Second * 30,
	}

	return &SIIClient{
		client:     client,
		firma:      firma,
		simulacion: true, // Por defecto en modo simulación
	}, nil
}

// SetSimulacion establece el modo de simulación
func (c *SIIClient) SetSimulacion(simular bool) {
	c.simulacion = simular
}

// ObtenerSemilla obtiene una semilla del SII
func (c *SIIClient) ObtenerSemilla() (string, error) {
	// Crear solicitud SOAP
	soapRequest := `<?xml version="1.0" encoding="UTF-8"?>
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:ws="http://ws.sii.dte">
   <soapenv:Header/>
   <soapenv:Body>
      <ws:getSeed/>
   </soapenv:Body>
</soapenv:Envelope>`

	// Enviar solicitud
	resp, err := c.client.Post(
		"https://maullin.sii.cl/DTEWS/CrSeed.jws",
		"text/xml;charset=UTF-8",
		bytes.NewBufferString(soapRequest),
	)
	if err != nil {
		return "", fmt.Errorf("error enviando solicitud: %v", err)
	}
	defer resp.Body.Close()

	// Leer respuesta
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error leyendo respuesta: %v", err)
	}

	// Decodificar respuesta SOAP
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

// ObtenerToken obtiene un token de autenticación del SII
func (c *SIIClient) ObtenerToken(semilla string) (string, error) {
	// Firmar semilla
	semillaFirmada, err := c.firma.FirmarSemilla(semilla)
	if err != nil {
		return "", fmt.Errorf("error firmando semilla: %v", err)
	}

	// Crear solicitud SOAP
	soapRequest := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:ws="http://ws.sii.dte">
   <soapenv:Header/>
   <soapenv:Body>
      <ws:getToken>
         <ws:seed>%s</ws:seed>
         <ws:signature>%s</ws:signature>
      </ws:getToken>
   </soapenv:Body>
</soapenv:Envelope>`, semilla, semillaFirmada)

	// Enviar solicitud
	resp, err := c.client.Post(
		"https://maullin.sii.cl/DTEWS/GetTokenFromSeed.jws",
		"text/xml;charset=UTF-8",
		bytes.NewBufferString(soapRequest),
	)
	if err != nil {
		return "", fmt.Errorf("error enviando solicitud: %v", err)
	}
	defer resp.Body.Close()

	// Leer respuesta
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error leyendo respuesta: %v", err)
	}

	// Decodificar respuesta SOAP
	var soapResponse struct {
		XMLName xml.Name `xml:"Envelope"`
		Body    struct {
			GetTokenResponse struct {
				Token string `xml:"token"`
			} `xml:"getTokenResponse"`
		} `xml:"Body"`
	}

	if err := xml.Unmarshal(body, &soapResponse); err != nil {
		return "", fmt.Errorf("error decodificando respuesta: %v", err)
	}

	return soapResponse.Body.GetTokenResponse.Token, nil
}

// EnviarDTE envía un documento al SII
func (c *SIIClient) EnviarDTE(sobre *models.Sobre, token string) error {
	// Convertir sobre a XML
	sobreXML, err := xml.MarshalIndent(sobre, "", "  ")
	if err != nil {
		return fmt.Errorf("error generando XML del sobre: %v", err)
	}

	// Firmar XML
	sobreXMLFirmado, err := c.firma.FirmarXML(sobreXML)
	if err != nil {
		return fmt.Errorf("error firmando XML: %v", err)
	}

	// Crear solicitud SOAP
	soapRequest := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:ws="http://ws.sii.dte">
   <soapenv:Header/>
   <soapenv:Body>
      <ws:sendDTE>
         <ws:token>%s</ws:token>
         <ws:xml>%s</ws:xml>
      </ws:sendDTE>
   </soapenv:Body>
</soapenv:Envelope>`, token, string(sobreXMLFirmado))

	// Enviar solicitud
	resp, err := c.client.Post(
		"https://maullin.sii.cl/DTEWS/RecepcionDTE.jws",
		"text/xml;charset=UTF-8",
		bytes.NewBufferString(soapRequest),
	)
	if err != nil {
		return fmt.Errorf("error enviando solicitud: %v", err)
	}
	defer resp.Body.Close()

	// Leer respuesta
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error leyendo respuesta: %v", err)
	}

	// Decodificar respuesta SOAP
	var soapResponse struct {
		XMLName xml.Name `xml:"Envelope"`
		Body    struct {
			SendDTEResponse struct {
				TrackID string `xml:"trackid"`
				Estado  string `xml:"estado"`
				Glosa   string `xml:"glosa"`
			} `xml:"sendDTEResponse"`
		} `xml:"Body"`
	}

	if err := xml.Unmarshal(body, &soapResponse); err != nil {
		return fmt.Errorf("error decodificando respuesta: %v", err)
	}

	// Verificar estado
	if soapResponse.Body.SendDTEResponse.Estado != "00" {
		return fmt.Errorf("error enviando DTE: %s", soapResponse.Body.SendDTEResponse.Glosa)
	}

	return nil
}

// ConsultarEstado consulta el estado de un DTE
func (c *SIIClient) ConsultarEstado(trackID string) (*models.EstadoResult, error) {
	// Crear solicitud SOAP
	soapRequest := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:ws="http://ws.sii.dte">
   <soapenv:Header/>
   <soapenv:Body>
      <ws:getEstDte>
         <ws:trackid>%s</ws:trackid>
      </ws:getEstDte>
   </soapenv:Body>
</soapenv:Envelope>`, trackID)

	// Enviar solicitud
	resp, err := c.client.Post(
		"https://maullin.sii.cl/DTEWS/QueryEstDte.jws",
		"text/xml;charset=UTF-8",
		bytes.NewBufferString(soapRequest),
	)
	if err != nil {
		return nil, fmt.Errorf("error enviando solicitud: %v", err)
	}
	defer resp.Body.Close()

	// Leer respuesta
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error leyendo respuesta: %v", err)
	}

	// Decodificar respuesta SOAP
	var soapResponse struct {
		XMLName xml.Name `xml:"Envelope"`
		Body    struct {
			GetEstDteResponse struct {
				Estado *models.EstadoResult `xml:"estado"`
			} `xml:"getEstDteResponse"`
		} `xml:"Body"`
	}

	if err := xml.Unmarshal(body, &soapResponse); err != nil {
		return nil, fmt.Errorf("error decodificando respuesta: %v", err)
	}

	return soapResponse.Body.GetEstDteResponse.Estado, nil
}
