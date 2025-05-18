package services

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// TestService maneja las pruebas de integración con el SII
type TestService struct {
	firmaService *FirmaService
	xmlService   *XMLService
	logService   *LogService
	ambiente     string // "certificacion" o "produccion"
	urlBase      string
}

// NewTestService crea una nueva instancia del servicio de pruebas
func NewTestService(firmaService *FirmaService, ambiente string) (*TestService, error) {
	// Crear servicios necesarios
	xmlService := NewXMLService()
	logService, err := NewLogService("logs/test")
	if err != nil {
		return nil, fmt.Errorf("error inicializando servicio de logs: %w", err)
	}

	// Determinar URL base según ambiente
	var urlBase string
	switch ambiente {
	case "certificacion":
		urlBase = "https://maullin.sii.cl/DTEWS"
	case "produccion":
		urlBase = "https://palena.sii.cl/DTEWS"
	default:
		return nil, fmt.Errorf("ambiente no válido: %s", ambiente)
	}

	return &TestService{
		firmaService: firmaService,
		xmlService:   xmlService,
		logService:   logService,
		ambiente:     ambiente,
		urlBase:      urlBase,
	}, nil
}

// ProbarSemilla prueba la obtención de semilla desde el SII
func (s *TestService) ProbarSemilla() error {
	// URL del servicio de semilla
	url := fmt.Sprintf("%s/CrSeed.jws?WSDL", s.urlBase)

	// Crear request
	xmlRequest := `<?xml version="1.0" encoding="UTF-8"?>
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
	<SOAP-ENV:Body>
		<getSeed xmlns="http://DefaultNamespace"/>
	</SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

	// Log del request
	if err := s.logService.LogXML("semilla_request", []byte(xmlRequest)); err != nil {
		return fmt.Errorf("error logging request semilla: %w", err)
	}

	// Enviar request
	resp, err := http.Post(url, "text/xml", strings.NewReader(xmlRequest))
	if err != nil {
		s.logService.LogError("semilla_http", err)
		return fmt.Errorf("error en request HTTP: %w", err)
	}
	defer resp.Body.Close()

	// Leer respuesta
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		s.logService.LogError("semilla_response", err)
		return fmt.Errorf("error leyendo respuesta: %w", err)
	}

	// Log de la respuesta
	if err := s.logService.LogXML("semilla_response", body); err != nil {
		return fmt.Errorf("error logging response semilla: %w", err)
	}

	return nil
}

// ProbarToken prueba la obtención de token usando una semilla
func (s *TestService) ProbarToken(semilla string) error {
	// URL del servicio de token
	url := fmt.Sprintf("%s/GetTokenFromSeed.jws?WSDL", s.urlBase)

	// Crear XML con semilla
	xmlSemilla := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<getToken>
    <item>
        <Semilla>%s</Semilla>
    </item>
</getToken>`, semilla)

	// Firmar XML
	resultado, err := s.firmaService.FirmarXML([]byte(xmlSemilla), "")
	if err != nil {
		s.logService.LogError("token_firma", err)
		return fmt.Errorf("error firmando semilla: %w", err)
	}

	// Crear request SOAP
	xmlRequest := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
	<SOAP-ENV:Body>
		<getToken xmlns="http://DefaultNamespace">
			%s
		</getToken>
	</SOAP-ENV:Body>
</SOAP-ENV:Envelope>`, resultado.XMLFirmado)

	// Log del request
	if err := s.logService.LogXML("token_request", []byte(xmlRequest)); err != nil {
		return fmt.Errorf("error logging request token: %w", err)
	}

	// Enviar request
	resp, err := http.Post(url, "text/xml", strings.NewReader(xmlRequest))
	if err != nil {
		s.logService.LogError("token_http", err)
		return fmt.Errorf("error en request HTTP: %w", err)
	}
	defer resp.Body.Close()

	// Leer respuesta
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		s.logService.LogError("token_response", err)
		return fmt.Errorf("error leyendo respuesta: %w", err)
	}

	// Log de la respuesta
	if err := s.logService.LogXML("token_response", body); err != nil {
		return fmt.Errorf("error logging response token: %w", err)
	}

	return nil
}

// ValidarFirma realiza una prueba completa de firma y validación
func (s *TestService) ValidarFirma(xmlData []byte) error {
	// Firmar el XML
	resultado, err := s.firmaService.FirmarXML(xmlData, "")
	if err != nil {
		s.logService.LogError("validacion_firma", err)
		return fmt.Errorf("error en firma: %w", err)
	}

	// Validar estructura contra esquema XSD
	if err := s.xmlService.ValidarEstructura([]byte(resultado.XMLFirmado), "schemas/xmldsignature_v10.xsd"); err != nil {
		s.logService.LogError("validacion_xsd", err)
		return fmt.Errorf("error en validación XSD: %w", err)
	}

	// Log del resultado
	if err := s.logService.LogValidacion([]byte(resultado.XMLFirmado), "Validación completa exitosa"); err != nil {
		return fmt.Errorf("error logging validación: %w", err)
	}

	return nil
}

// ProbarFlujoCompleto realiza una prueba del flujo completo de autenticación
func (s *TestService) ProbarFlujoCompleto() error {
	// 1. Obtener semilla
	if err := s.ProbarSemilla(); err != nil {
		return fmt.Errorf("error en prueba de semilla: %w", err)
	}

	// Esperar un momento entre requests
	time.Sleep(2 * time.Second)

	// 2. Obtener token (usando una semilla de ejemplo)
	semilla := "1234567890" // En un caso real, se obtendría de la respuesta anterior
	if err := s.ProbarToken(semilla); err != nil {
		return fmt.Errorf("error en prueba de token: %w", err)
	}

	return nil
}

// GenerarXMLPrueba genera un XML de prueba para validación
func (s *TestService) GenerarXMLPrueba() []byte {
	xmlPrueba := `<?xml version="1.0" encoding="UTF-8"?>
<documento>
    <contenido>Documento de prueba para firma digital</contenido>
    <fecha>%s</fecha>
</documento>`

	return []byte(fmt.Sprintf(xmlPrueba, time.Now().Format("2006-01-02T15:04:05")))
}
