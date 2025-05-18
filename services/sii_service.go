package services

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"FMgo/config"
	"FMgo/models"
	"FMgo/utils"
	"FMgo/utils/sii"
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

// CrearEnvioMasivoBoletas crea un envío masivo de boletas según el esquema EnvioBOLETA_v11.xsd
func (s *SIIService) CrearEnvioMasivoBoletas(xmls []string, rutEmisor string) (string, error) {
	if len(xmls) == 0 {
		return "", fmt.Errorf("no hay boletas para enviar")
	}
	if len(xmls) > 500 {
		return "", fmt.Errorf("el número máximo de boletas por envío es 500, se recibieron %d", len(xmls))
	}

	// Validar RUT emisor
	if !utils.ValidarRUT(rutEmisor) {
		return "", fmt.Errorf("RUT emisor inválido: %s", rutEmisor)
	}

	// Crear carátula del envío
	caratula := fmt.Sprintf(`
		<Caratula version="1.0">
			<RutEmisor>%s</RutEmisor>
			<RutEnvia>%s</RutEnvia>
			<RutReceptor>60803000-K</RutReceptor>
			<FchResol>%s</FchResol>
			<NroResol>0</NroResol>
			<TmstFirmaEnv>%s</TmstFirmaEnv>
			<SubTotDTE>
				<TpoDTE>39</TpoDTE>
				<NroDTE>%d</NroDTE>
			</SubTotDTE>
		</Caratula>`,
		rutEmisor,
		rutEmisor,
		time.Now().Format("2006-01-02"),
		time.Now().Format("2006-01-02T15:04:05"),
		len(xmls),
	)

	// Crear envío masivo
	envioID := fmt.Sprintf("SetBOL_%s_%s", rutEmisor, time.Now().Format("20060102150405"))
	envioMasivo := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<EnvioBOLETA xmlns="http://www.sii.cl/SiiDte" version="1.0">
	<SetDTE ID="%s">
		%s
		%s
	</SetDTE>
</EnvioBOLETA>`,
		envioID,
		caratula,
		strings.Join(xmls, "\n"),
	)

	// Validar estructura del envío
	if err := s.ValidarEnvioMasivo(envioMasivo); err != nil {
		return "", fmt.Errorf("error validando envío masivo: %v", err)
	}

	return envioMasivo, nil
}

// ValidarEnvioMasivo valida la estructura del envío masivo contra el esquema XSD
func (s *SIIService) ValidarEnvioMasivo(xml string) error {
	// Cargar esquema XSD
	xsdPath := "schemas/EnvioBOLETA_v11.xsd"
	xsd, err := ioutil.ReadFile(xsdPath)
	if err != nil {
		return fmt.Errorf("error leyendo esquema XSD: %v", err)
	}

	// Validar XML contra esquema
	if err := utils.ValidarXMLContraXSD(xml, string(xsd)); err != nil {
		return fmt.Errorf("error validando XML contra esquema: %v", err)
	}

	return nil
}

// FirmarEnvioMasivo firma digitalmente el envío masivo de boletas
func (s *SIIService) FirmarEnvioMasivo(xml string) (string, error) {
	// Obtener certificado
	cert, err := s.obtenerCertificado()
	if err != nil {
		return "", fmt.Errorf("error obteniendo certificado: %v", err)
	}

	// Firmar XML
	xmlFirmado, err := s.firmaService.FirmarXML(xml, cert)
	if err != nil {
		return "", fmt.Errorf("error firmando XML: %v", err)
	}

	return xmlFirmado, nil
}

// GenerarXMLBoleta genera el XML para una boleta individual
func (s *SIIService) GenerarXMLBoleta(boleta *models.Boleta) (string, error) {
	// Convertir boleta a BOLETAType
	boletaType := models.ConvertirBoleta(boleta)

	// Generar XML
	xml, err := xml.MarshalIndent(boletaType, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error generando XML: %v", err)
	}

	return string(xml), nil
}

// ValidarXMLBoleta valida el XML de una boleta individual
func (s *SIIService) ValidarXMLBoleta(xml string) error {
	// Cargar esquema XSD
	xsdPath := "schemas/DTE_v10.xsd"
	xsd, err := ioutil.ReadFile(xsdPath)
	if err != nil {
		return fmt.Errorf("error leyendo esquema XSD: %v", err)
	}

	// Validar XML contra esquema
	if err := utils.ValidarXMLContraXSD(xml, string(xsd)); err != nil {
		return fmt.Errorf("error validando XML contra esquema: %v", err)
	}

	return nil
}
