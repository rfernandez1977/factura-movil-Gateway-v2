package services

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/cursor/FMgo/config"
	"github.com/cursor/FMgo/models"
	"github.com/cursor/FMgo/utils/sii"
)

// SIIService es un alias para la interfaz sii.SIIService
type SIIService = sii.SIIService

// NewSIIService crea una nueva instancia del servicio de SII
func NewSIIService(config *config.SupabaseConfig) *SIIService {
	return &SIIService{
		config: config,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// EnviarDTE envía un documento al SII
func (s *SIIService) EnviarDTE(xmlData []byte, empresa *models.Empresa) (*models.RespuestaSII, error) {
	// Crear estructura de envío
	dte := &models.DTEXMLModel{
		Version: "1.0",
		Documento: models.DocumentoXMLModel{
			ID: fmt.Sprintf("DTE_%d", time.Now().Unix()),
			Encabezado: models.EncabezadoXMLModel{
				ID: models.IDDocumentoXMLModel{
					TipoDTE:      models.FacturaElectronica,
					Folio:        1,
					FechaEmision: time.Now().Format("2006-01-02"),
				},
				Emisor: models.EmisorXMLModel{
					RUT:         empresa.RUT,
					RazonSocial: empresa.RazonSocial,
					Giro:        empresa.Giro,
					Direccion:   empresa.Direccion,
					Comuna:      empresa.Comuna,
					Ciudad:      empresa.Ciudad,
				},
			},
		},
	}

	// Generar XML de envío
	envioXML, err := xml.MarshalIndent(dte, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error al generar XML de envío: %v", err)
	}

	// Enviar al SII
	resp, err := s.client.Post(
		"https://palena.sii.cl/DTEWS/CrSeed.jws",
		"application/xml",
		bytes.NewReader(envioXML),
	)
	if err != nil {
		return nil, fmt.Errorf("error al enviar al SII: %v", err)
	}
	defer resp.Body.Close()

	// Leer respuesta
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error al leer respuesta del SII: %v", err)
	}

	// Parsear respuesta
	var respuesta models.RespuestaSII
	if err := xml.Unmarshal(respBody, &respuesta); err != nil {
		return nil, fmt.Errorf("error al parsear respuesta del SII: %v", err)
	}

	return &respuesta, nil
}

// ConsultarEstado consulta el estado de un documento en el SII
func (s *SIIService) ConsultarEstado(trackID string) (*models.EstadoSII, error) {
	if trackID == "" {
		return nil, fmt.Errorf("trackID es requerido")
	}

	// Consultar estado al SII
	resp, err := s.client.Get(fmt.Sprintf("https://palena.sii.cl/DTEWS/QueryEstDte.jws?trackID=%s", trackID))
	if err != nil {
		return nil, fmt.Errorf("error al consultar estado al SII: %v", err)
	}
	defer resp.Body.Close()

	// Leer respuesta
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error al leer respuesta del SII: %v", err)
	}

	// Parsear respuesta
	var estado models.EstadoSII
	if err := xml.Unmarshal(respBody, &estado); err != nil {
		return nil, fmt.Errorf("error al parsear estado del SII: %v", err)
	}

	return &estado, nil
}

// ConsultarDTE consulta un DTE específico en el SII
func (s *SIIService) ConsultarDTE(tipoDTE, folio, rutEmisor string) (*models.EstadoSII, error) {
	if tipoDTE == "" || folio == "" || rutEmisor == "" {
		return nil, fmt.Errorf("tipoDTE, folio y rutEmisor son requeridos")
	}

	// Consultar DTE al SII
	resp, err := s.client.Get(fmt.Sprintf("https://palena.sii.cl/DTEWS/QueryDte.jws?tipoDTE=%s&folio=%s&rutEmisor=%s", tipoDTE, folio, rutEmisor))
	if err != nil {
		return nil, fmt.Errorf("error al consultar DTE al SII: %v", err)
	}
	defer resp.Body.Close()

	// Leer respuesta
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error al leer respuesta del SII: %v", err)
	}

	// Parsear respuesta
	var estado models.EstadoSII
	if err := xml.Unmarshal(respBody, &estado); err != nil {
		return nil, fmt.Errorf("error al parsear estado del DTE: %v", err)
	}

	return &estado, nil
}

// VerificarComunicacion verifica la comunicación con el SII
func (s *SIIService) VerificarComunicacion() error {
	// Verificar comunicación con el SII
	resp, err := s.client.Get("https://palena.sii.cl/DTEWS/CrSeed.jws")
	if err != nil {
		return fmt.Errorf("error al verificar comunicación con el SII: %v", err)
	}
	defer resp.Body.Close()

	// Leer respuesta
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error al leer respuesta del SII: %v", err)
	}

	// Parsear respuesta
	var estado models.EstadoSII
	if err := xml.Unmarshal(respBody, &estado); err != nil {
		return fmt.Errorf("error al parsear estado del SII: %v", err)
	}

	if estado.Estado != models.EstadoSIIAceptado {
		return fmt.Errorf("error de comunicación con el SII: %s", estado.Glosa)
	}

	return nil
}
