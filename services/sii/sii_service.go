package sii

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/cursor/FMgo/config"
	"github.com/cursor/FMgo/models"
)

// SIIService define la interfaz para el servicio SII
type SIIService interface {
	// ConsultarEstado consulta el estado de una solicitud al SII
	ConsultarEstado(trackID string) (*models.EstadoSII, error)
	// EnviarDTE envía un DTE al SII
	EnviarDTE(dte []byte) (*models.EstadoSII, error)
	// ConsultarDTE consulta un DTE específico en el SII
	ConsultarDTE(tipoDTE, folio, rutEmisor string) (*models.EstadoSII, error)
	// VerificarComunicacion verifica la comunicación con el SII
	VerificarComunicacion() error
}

// SIIServiceImpl implementa la interfaz sii.SIIService
type SIIServiceImpl struct {
	baseURL    string
	token      string
	ambiente   string
	certFile   string
	keyFile    string
	httpClient *http.Client
}

// NewSIIServiceImpl crea una nueva instancia del servicio SII
func NewSIIServiceImpl(baseURL, token, ambiente, certFile, keyFile string) (*SIIServiceImpl, error) {
	// Validar parámetros
	if baseURL == "" {
		return nil, fmt.Errorf("baseURL es requerido")
	}
	if token == "" {
		return nil, fmt.Errorf("token es requerido")
	}
	if certFile == "" {
		return nil, fmt.Errorf("certFile es requerido")
	}
	if keyFile == "" {
		return nil, fmt.Errorf("keyFile es requerido")
	}

	// Cargar certificado
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, fmt.Errorf("error cargando certificado: %v", err)
	}

	// Configurar cliente HTTP
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates:       []tls.Certificate{cert},
				InsecureSkipVerify: ambiente == "CERTIFICACION",
			},
		},
		Timeout: 30 * time.Second,
	}

	return &SIIServiceImpl{
		baseURL:    baseURL,
		token:      token,
		ambiente:   ambiente,
		certFile:   certFile,
		keyFile:    keyFile,
		httpClient: client,
	}, nil
}

// ConsultarEstado consulta el estado de un DTE
func (s *SIIServiceImpl) ConsultarEstado(trackID string) (*models.EstadoSII, error) {
	// Crear request
	url := fmt.Sprintf("%s/ConsultaEstado?trackID=%s", s.baseURL, trackID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error al crear request: %w", err)
	}

	// Configurar headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.token))

	// Enviar request
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error al enviar request: %w", err)
	}
	defer resp.Body.Close()

	// Leer respuesta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error al leer respuesta: %w", err)
	}

	// Verificar status code
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusBadRequest {
		return nil, fmt.Errorf("error en respuesta del SII: %s", string(body))
	}

	// Decodificar respuesta
	var respuesta models.RespuestaSII
	if err := json.Unmarshal(body, &respuesta); err != nil {
		return nil, fmt.Errorf("error al decodificar respuesta: %w", err)
	}

	// Convertir a EstadoSII
	estado := &models.EstadoSII{
		Estado:  respuesta.Estado,
		Glosa:   respuesta.Glosa,
		TrackID: respuesta.TrackID,
		Fecha:   respuesta.FechaProceso,
	}

	// Si hay errores, agregarlos al estado
	if len(respuesta.Errores) > 0 {
		estado.Errores = make([]models.ErrorSII, len(respuesta.Errores))
		copy(estado.Errores, respuesta.Errores)
	}

	return estado, nil
}

// EnviarDTE envía un DTE al SII
func (s *SIIServiceImpl) EnviarDTE(dte []byte) (*models.EstadoSII, error) {
	url := fmt.Sprintf("%s/dte/enviar", s.baseURL)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(dte))
	if err != nil {
		return nil, fmt.Errorf("error al crear request: %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.token))
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error al enviar DTE: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error al leer respuesta: %v", err)
	}

	var respuesta models.EstadoSII
	if err := json.Unmarshal(body, &respuesta); err != nil {
		return nil, fmt.Errorf("error al decodificar respuesta: %v", err)
	}

	return &respuesta, nil
}

// ConsultarDTE consulta un DTE específico en el SII
func (s *SIIServiceImpl) ConsultarDTE(tipoDTE, folio, rutEmisor string) (*models.EstadoSII, error) {
	url := fmt.Sprintf("%s/dte/consulta/%s/%s/%s", s.baseURL, tipoDTE, folio, rutEmisor)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error al crear request: %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.token))
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error al consultar DTE: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error al leer respuesta: %v", err)
	}

	var respuesta models.EstadoSII
	if err := json.Unmarshal(body, &respuesta); err != nil {
		return nil, fmt.Errorf("error al decodificar respuesta: %v", err)
	}

	return &respuesta, nil
}

// VerificarComunicacion verifica la comunicación con el SII
func (s *SIIServiceImpl) VerificarComunicacion() error {
	url := fmt.Sprintf("%s/verificar", s.baseURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("error al crear request: %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.token))

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error al verificar comunicación: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error en respuesta del SII: %d", resp.StatusCode)
	}

	return nil
}

// Service representa el servicio SII legacy
type Service struct {
	config *config.Config
}

// NewService crea una nueva instancia del servicio SII legacy
func NewService(config *config.Config) *Service {
	return &Service{
		config: config,
	}
}

// EnviarDTE envía un DTE al SII
func (s *Service) EnviarDTE(sobre *models.SobreDTE) (*models.RespuestaSII, error) {
	// Implementación mock
	return &models.RespuestaSII{
		Estado:  "OK",
		Glosa:   "Documento Recibido",
		TrackID: "123456",
	}, nil
}

// ConsultarEstado consulta el estado de un DTE
func (s *Service) ConsultarEstado(trackID string) (*models.EstadoSII, error) {
	// Implementación mock
	return &models.EstadoSII{
		Estado:  "OK",
		Glosa:   "Documento Aceptado",
		TrackID: trackID,
		Fecha:   time.Now(),
	}, nil
}
