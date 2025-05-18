package client

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"crypto/tls"

	"FMgo/core/sii/infrastructure/certificates"
	"FMgo/core/sii/logger"
	"FMgo/core/sii/models"
	"FMgo/core/sii/retry"
	"FMgo/utils/xmlutils"
)

// HTTPClient implementa la interfaz para comunicación con el SII
type HTTPClient struct {
	baseURL     string
	retryCount  int
	timeout     time.Duration
	logger      logger.Logger
	xmlParser   *xmlutils.XMLParser
	client      *http.Client
	certManager *certificates.CertManager
	ambiente    models.Ambiente
}

// Estructuras para parsear las respuestas XML del SII
type respuestaSII struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    respuestaBody
}

type respuestaBody struct {
	XMLName xml.Name    `xml:"Body"`
	Content interface{} `xml:",any"`
}

type getSeedResponse struct {
	XMLName xml.Name `xml:"getSeedResponse"`
	Semilla string   `xml:"seed"`
}

type getTokenResponse struct {
	XMLName xml.Name `xml:"getTokenResponse"`
	Token   string   `xml:"token"`
}

type sendDTEResponse struct {
	XMLName xml.Name `xml:"sendDTEResponse"`
	Estado  string   `xml:"estado"`
	Glosa   string   `xml:"glosa"`
	TrackID string   `xml:"trackid"`
}

type getStatusResponse struct {
	XMLName xml.Name `xml:"getStatusResponse"`
	Estado  string   `xml:"estado"`
	Glosa   string   `xml:"glosa"`
	TrackID string   `xml:"trackid,omitempty"`
}

type getEstDteResponse struct {
	XMLName xml.Name `xml:"getEstDteResponse"`
	Estado  string   `xml:"estado"`
	Glosa   string   `xml:"glosa"`
	TrackID string   `xml:"trackid,omitempty"`
}

// Errores específicos del cliente HTTP
var (
	ErrSemillaNoEncontrada = fmt.Errorf("respuesta del SII no contiene semilla")
	ErrTokenNoEncontrado   = fmt.Errorf("respuesta del SII no contiene token")
	ErrEstadoNoEncontrado  = fmt.Errorf("respuesta del SII no contiene estado")
	ErrDTENoEncontrado     = fmt.Errorf("respuesta del SII no contiene estado del DTE")
	ErrEnvioNoEncontrado   = fmt.Errorf("respuesta del SII no contiene estado del envío")
)

// NewHTTPClient crea una nueva instancia del cliente HTTP
func NewHTTPClient(config *models.Config, log logger.Logger) (*HTTPClient, error) {
	if config == nil {
		return nil, fmt.Errorf("configuración no puede ser nil")
	}
	if log == nil {
		return nil, fmt.Errorf("logger no puede ser nil")
	}

	xmlParser := xmlutils.NewXMLParser(false)

	// Si hay un esquema XSD configurado, establecerlo
	if config.SchemaPath != "" {
		if err := xmlParser.SetValidator(config.SchemaPath); err != nil {
			return nil, fmt.Errorf("error configurando validador XML: %w", err)
		}
	}

	// Configurar certificado
	certManager, err := certificates.NewCertManager(config.CertPath, config.KeyPath)
	if err != nil {
		return nil, fmt.Errorf("error configurando certificados: %w", err)
	}

	// Configurar cliente HTTP con certificado TLS
	cert, err := tls.LoadX509KeyPair(config.CertPath, config.KeyPath)
	if err != nil {
		return nil, fmt.Errorf("error cargando certificado: %w", err)
	}

	client := &http.Client{
		Timeout: time.Duration(config.Timeout) * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates: []tls.Certificate{cert},
			},
		},
	}

	return &HTTPClient{
		baseURL:     config.BaseURL,
		retryCount:  config.RetryCount,
		timeout:     time.Duration(config.Timeout) * time.Second,
		logger:      log,
		xmlParser:   xmlParser,
		client:      client,
		certManager: certManager,
		ambiente:    config.Ambiente,
	}, nil
}

// parseSOAPResponse parsea una respuesta SOAP en la estructura correspondiente
func (c *HTTPClient) parseSOAPResponse(body []byte, content interface{}) error {
	if body == nil {
		return fmt.Errorf("cuerpo de respuesta no puede ser nil")
	}
	if content == nil {
		return fmt.Errorf("contenedor de respuesta no puede ser nil")
	}
	return c.xmlParser.ParseSOAP(body, content)
}

// validarRespuestaHTTP valida la respuesta HTTP del SII y retorna un error tipado si hay problemas
func validarRespuestaHTTP(resp *http.Response) error {
	if resp == nil {
		return models.NewSIIError(models.ErrProcesamiento, "Respuesta HTTP es nil", nil)
	}

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			body = []byte("no se pudo leer el cuerpo de la respuesta")
		}

		switch resp.StatusCode {
		case http.StatusUnauthorized:
			return models.NewSIIError(models.ErrAuthInvalid, "Credenciales inválidas", nil)
		case http.StatusForbidden:
			return models.NewSIIError(models.ErrCertInvalid, "Certificado inválido o no autorizado", nil)
		case http.StatusRequestTimeout:
			return models.NewSIIError(models.ErrTimeout, "Tiempo de espera agotado", nil)
		case http.StatusBadGateway, http.StatusServiceUnavailable:
			return models.NewSIIError(models.ErrServidor, "Servicio no disponible", nil)
		default:
			return models.NewSIIError(models.ErrConexion, fmt.Sprintf("Error del SII (código %d): %s", resp.StatusCode, string(body)), nil)
		}
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(strings.ToLower(contentType), "xml") {
		return models.NewSIIError(models.ErrProcesamiento, fmt.Sprintf("Tipo de contenido no válido: %s", contentType), nil)
	}

	return nil
}

// ObtenerSemilla obtiene una semilla del SII
func (c *HTTPClient) ObtenerSemilla(ctx context.Context) (string, error) {
	// Verificar certificado antes de cada operación crítica
	if err := c.certManager.ValidateCertificate(); err != nil {
		return "", fmt.Errorf("error de certificado: %w", err)
	}

	if ctx == nil {
		ctx = context.Background()
	}

	url := models.URLSemillaCert
	if c.ambiente == models.Produccion {
		url = models.URLSemillaProd
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("error al crear request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error al obtener semilla: %w", err)
	}
	defer resp.Body.Close()

	if err := validarRespuestaHTTP(resp); err != nil {
		return "", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error al leer respuesta: %w", err)
	}

	var respuesta respuestaSII
	respuesta.Body.Content = &getSeedResponse{}
	if err := c.parseSOAPResponse(body, &respuesta); err != nil {
		return "", fmt.Errorf("error al decodificar respuesta: %w", err)
	}

	if seedResp, ok := respuesta.Body.Content.(*getSeedResponse); ok {
		if seedResp.Semilla == "" {
			return "", models.NewSIIError(models.ErrProcesamiento, "Respuesta del SII no contiene semilla", nil)
		}
		return seedResp.Semilla, nil
	}

	return "", models.NewSIIError(models.ErrProcesamiento, "Respuesta del SII no contiene semilla", nil)
}

// ObtenerToken obtiene un token usando una semilla
func (c *HTTPClient) ObtenerToken(ctx context.Context, semilla string) (string, error) {
	// Verificar certificado
	if err := c.certManager.ValidateCertificate(); err != nil {
		return "", fmt.Errorf("error de certificado: %w", err)
	}

	if ctx == nil {
		ctx = context.Background()
	}

	if semilla == "" {
		return "", models.NewSIIError(models.ErrProcesamiento, "Semilla no puede estar vacía", nil)
	}

	url := models.URLTokenCert
	if c.ambiente == models.Produccion {
		url = models.URLTokenProd
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBufferString(semilla))
	if err != nil {
		return "", fmt.Errorf("error al crear request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error al obtener token: %w", err)
	}
	defer resp.Body.Close()

	if err := validarRespuestaHTTP(resp); err != nil {
		return "", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error al leer respuesta: %w", err)
	}

	var respuesta respuestaSII
	respuesta.Body.Content = &getTokenResponse{}
	if err := c.parseSOAPResponse(body, &respuesta); err != nil {
		return "", fmt.Errorf("error al decodificar respuesta: %w", err)
	}

	if tokenResp, ok := respuesta.Body.Content.(*getTokenResponse); ok {
		if tokenResp.Token == "" {
			return "", models.NewSIIError(models.ErrAuthInvalid, "Respuesta del SII no contiene token", nil)
		}
		return tokenResp.Token, nil
	}

	return "", models.NewSIIError(models.ErrAuthInvalid, "Respuesta del SII no contiene token", nil)
}

// EnviarDTE envía un DTE al SII con reintentos y manejo de errores
func (c *HTTPClient) EnviarDTE(ctx context.Context, sobre []byte, token string) (*models.RespuestaSII, error) {
	endpoint := fmt.Sprintf("%s/recepciondte", c.baseURL)

	// Preparar request con certificado
	req, err := c.prepareRequest(ctx, "POST", endpoint, sobre)
	if err != nil {
		return nil, fmt.Errorf("error preparando request: %w", err)
	}

	// Agregar token de autenticación
	req.Header.Set("Cookie", fmt.Sprintf("TOKEN=%s", token))

	// Configurar reintentos
	retryConfig := &retry.Config{
		MaxRetries:  c.retryCount,
		WaitTime:    time.Second,
		MaxWaitTime: c.timeout,
	}

	// Ejecutar con reintentos
	var resp *http.Response
	err = retry.Do(ctx, retryConfig, func() error {
		var reqErr error
		resp, reqErr = c.doRequest(req)
		return reqErr
	})

	if err != nil {
		return nil, fmt.Errorf("error enviando DTE: %w", err)
	}

	// Procesar respuesta
	respuesta := &models.RespuestaSII{}
	if err := xml.NewDecoder(resp.Body).Decode(respuesta); err != nil {
		return nil, fmt.Errorf("error decodificando respuesta: %w", err)
	}

	return respuesta, nil
}

// ConsultarEstado consulta el estado de un envío
func (c *HTTPClient) ConsultarEstado(ctx context.Context, trackID string) (*models.EstadoConsulta, error) {
	endpoint := fmt.Sprintf("%s/consultaestadodte", c.baseURL)

	// Preparar datos de consulta
	data := url.Values{}
	data.Set("trackID", trackID)

	// Preparar request
	req, err := c.prepareRequest(ctx, "POST", endpoint, []byte(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("error preparando request: %w", err)
	}

	// Configurar reintentos
	retryConfig := &retry.Config{
		MaxRetries:  c.retryCount,
		WaitTime:    time.Second,
		MaxWaitTime: c.timeout,
	}

	// Ejecutar con reintentos
	var resp *http.Response
	err = retry.Do(ctx, retryConfig, func() error {
		var reqErr error
		resp, reqErr = c.doRequest(req)
		return reqErr
	})

	if err != nil {
		return nil, fmt.Errorf("error consultando estado: %w", err)
	}

	// Procesar respuesta
	estado := &models.EstadoConsulta{
		Estado: models.EstadoDesconocido,
		Glosa:  "Estado desconocido",
	}
	if err := xml.NewDecoder(resp.Body).Decode(estado); err != nil {
		return nil, fmt.Errorf("error decodificando estado: %w", err)
	}

	return estado, nil
}

// doRequest ejecuta un request HTTP con el certificado configurado
func (c *HTTPClient) doRequest(req *http.Request) (*http.Response, error) {
	// Ejecutar request
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error en request HTTP: %w", err)
	}

	// Validar código de respuesta
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("error del SII: código %d, respuesta: %s", resp.StatusCode, string(body))
	}

	return resp, nil
}

// ConsultarDTE consulta el estado de un DTE específico
func (c *HTTPClient) ConsultarDTE(ctx context.Context, tipoDTE models.TipoDocumentoSII, folio int64, rutEmisor string) (*models.EstadoConsulta, error) {
	// Verificar certificado
	if err := c.certManager.ValidateCertificate(); err != nil {
		return nil, fmt.Errorf("error de certificado: %w", err)
	}

	if ctx == nil {
		ctx = context.Background()
	}

	if tipoDTE == "" {
		return nil, fmt.Errorf("tipoDTE no puede estar vacío")
	}

	if folio <= 0 {
		return nil, fmt.Errorf("folio debe ser mayor a 0")
	}

	if rutEmisor == "" {
		return nil, fmt.Errorf("rutEmisor no puede estar vacío")
	}

	url := models.URLEstadoDTECert
	if c.ambiente == models.Produccion {
		url = models.URLEstadoDTEProd
	}

	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s?tipoDTE=%s&folio=%d&rutEmisor=%s", url, tipoDTE, folio, rutEmisor), nil)
	if err != nil {
		return nil, fmt.Errorf("error al crear request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error al consultar DTE: %w", err)
	}
	defer resp.Body.Close()

	if err := validarRespuestaHTTP(resp); err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error al leer respuesta: %w", err)
	}

	var respuesta respuestaSII
	respuesta.Body.Content = &getEstDteResponse{}
	if err := c.parseSOAPResponse(body, &respuesta); err != nil {
		return nil, fmt.Errorf("error al decodificar respuesta: %w", err)
	}

	if dteResp, ok := respuesta.Body.Content.(*getEstDteResponse); ok {
		if dteResp.Estado == "" {
			return nil, ErrDTENoEncontrado
		}
		return &models.EstadoConsulta{
			Estado:  models.ParseEstadoSII(dteResp.Estado),
			Glosa:   dteResp.Glosa,
			TrackID: dteResp.TrackID,
		}, nil
	}

	return nil, ErrDTENoEncontrado
}

// VerificarComunicacion verifica la comunicación con el SII
func (c *HTTPClient) VerificarComunicacion(ctx context.Context) error {
	// Verificar certificado
	if err := c.certManager.ValidateCertificate(); err != nil {
		return fmt.Errorf("error de certificado: %w", err)
	}

	// Intentar obtener una semilla como prueba de comunicación
	_, err := c.ObtenerSemilla(ctx)
	if err != nil {
		return fmt.Errorf("error al verificar comunicación: %v", err)
	}
	return nil
}

// GetCertificateInfo retorna la información del certificado actual
func (c *HTTPClient) GetCertificateInfo() *certificates.CertificateInfo {
	return c.certManager.GetCertificateInfo()
}

// IsExpiringSoon verifica si el certificado está por expirar
func (c *HTTPClient) IsExpiringSoon(daysThreshold int) bool {
	return c.certManager.IsExpiringSoon(daysThreshold)
}

// prepareRequest prepara un request HTTP con los headers necesarios
func (c *HTTPClient) prepareRequest(ctx context.Context, method, url string, body []byte) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("error creando request: %w", err)
	}

	req.Header.Set("Content-Type", "text/xml; charset=utf-8")
	return req, nil
}
