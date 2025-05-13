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

const (
	// URLs del SII
	urlSemilla = "https://apicert.sii.cl/recursos/v1/boleta.electronica.semilla"
	urlToken   = "https://apicert.sii.cl/recursos/v1/boleta.electronica.token"
	urlEnvio   = "https://apicert.sii.cl/recursos/v1/boleta.electronica.envio"

	// Headers comunes
	contentType     = "Content-Type"
	userAgent       = "User-Agent"
	accept          = "Accept"
	authorization   = "Authorization"
	rutEmisor       = "Rut-Emisor"
	rutEnvia        = "Rut-Envia"
	applicationJSON = "application/json"
	applicationXML  = "application/xml"
	mozillaAgent    = "Mozilla/5.0"

	// Algoritmos de firma
	algoritmoCanonicalizacion = "http://www.w3.org/TR/2001/REC-xml-c14n-20010315"
	algoritmoFirma            = "http://www.w3.org/2000/09/xmldsig#rsa-sha1"
	algoritmoDigest           = "http://www.w3.org/2000/09/xmldsig#sha1"

	// Estados válidos del SII
	estadoOK    = "OK"
	estadoERROR = "ERROR"
)

// ObtenerSemilla obtiene una semilla del SII
func ObtenerSemilla() (string, error) {
	resp, err := http.Get(urlSemilla)
	if err != nil {
		return "", fmt.Errorf("error obteniendo semilla: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error en respuesta del SII: código %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error leyendo respuesta: %v", err)
	}

	return string(body), nil
}

// ObtenerToken obtiene un token de autenticación del SII
func ObtenerToken(semilla, certPath, keyPath string) (string, error) {
	client, err := crearClienteHTTP(certPath, keyPath)
	if err != nil {
		return "", fmt.Errorf("error creando cliente HTTP: %v", err)
	}

	req, err := http.NewRequest(http.MethodGet, urlToken, nil)
	if err != nil {
		return "", fmt.Errorf("error creando request: %v", err)
	}

	req.Header.Set(contentType, applicationJSON)
	req.Header.Set(userAgent, mozillaAgent)
	req.Header.Set(accept, applicationJSON)

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error realizando request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("error en respuesta del SII: código %d, mensaje: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error leyendo respuesta: %v", err)
	}

	return string(body), nil
}

// GenerarXMLDTE genera un documento DTE en formato XML
func GenerarXMLDTE(doc *models.DocumentoTributario, empresa *models.Empresa) ([]byte, error) {
	// Calcular el monto neto como 84% del monto total (aproximación para Chile con IVA del 19%)
	montoNeto := doc.MontoTotal / 1.19
	montoIVA := doc.MontoTotal - montoNeto

	// Convertir valores a tipos apropiados
	montoNetoInt64 := int64(montoNeto)
	montoIVAInt64 := int64(montoIVA)
	montoTotalInt64 := int64(doc.MontoTotal)
	tasaIVA := float64(19) // Tasa estándar de IVA en Chile

	// Crear estructura para el XML
	dte := models.DTEXMLModel{
		Version: "1.0",
		Documento: models.DocumentoXMLModel{
			ID: fmt.Sprintf("DTE_%s_%d", doc.RutEmisor, doc.Folio),
			Encabezado: models.EncabezadoXMLModel{
				IdDoc: models.IDDocumentoXMLModel{
					TipoDTE:      doc.TipoDTE,
					Folio:        int(doc.Folio),
					FechaEmision: doc.FechaEmision.Format("2006-01-02"),
				},
				Emisor: models.EmisorXMLModel{
					RUT:         empresa.RUT,
					RazonSocial: empresa.RazonSocial,
					Giro:        empresa.Giro,
					Direccion:   empresa.Direccion,
					Comuna:      empresa.Comuna,
					Ciudad:      empresa.Ciudad,
				},
				Receptor: models.ReceptorXMLModel{
					RUT:         doc.RutReceptor,
					RazonSocial: "Receptor", // Valor por defecto
					Direccion:   "",
					Comuna:      "",
					Ciudad:      "",
				},
				Totales: models.TotalesXMLModel{
					MntNeto:  &montoNetoInt64,
					TasaIVA:  &tasaIVA,
					IVA:      &montoIVAInt64,
					MntTotal: montoTotalInt64,
				},
			},
			Detalle: make([]models.DetalleDTEXML, len(doc.Detalles)),
		},
	}

	// Agregar items
	for i, item := range doc.Detalles {
		descripcion := item.Descripcion
		cantidad := float64(item.Cantidad)
		precio := float64(item.PrecioUnitario)
		montoItem := int64(item.MontoItem)

		dte.Documento.Detalle[i] = models.DetalleDTEXML{
			NroLinDet:   i + 1,
			Nombre:      descripcion,
			Cantidad:    &cantidad,
			Precio:      &precio,
			MontoItem:   montoItem,
			Descripcion: &descripcion,
		}
	}

	// Generar XML
	xmlData, err := xml.MarshalIndent(dte, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error generando XML: %v", err)
	}

	// Agregar declaración XML
	xmlDeclaration := []byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
	xmlData = append(xmlDeclaration, xmlData...)

	return xmlData, nil
}

// ValidarDTE valida un DTE contra el esquema XSD
func ValidarDTE(xmlData []byte) error {
	// Por ahora, solo validamos que el XML sea válido
	var doc models.DTEXMLModel
	if err := xml.Unmarshal(xmlData, &doc); err != nil {
		return fmt.Errorf("error validando DTE: %v", err)
	}
	return nil
}

// crearFirmaXML crea la firma XML según la especificación del SII
func crearFirmaXML(doc *models.DocumentoTributario) (*models.FirmaXMLModel, error) {
	// TODO: Implementar firma XML
	return nil, nil
}

// insertarFirmaEnXML inserta la firma en el documento XML
func insertarFirmaEnXML(xmlData []byte, firma *models.FirmaXMLModel) ([]byte, error) {
	// Parsear XML original
	var doc models.DTEXMLModel
	if err := xml.Unmarshal(xmlData, &doc); err != nil {
		return nil, fmt.Errorf("error parseando XML: %v", err)
	}

	// Agregar firma
	doc.Signature = firma

	// Convertir a XML
	xmlFirmado, err := xml.MarshalIndent(doc, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error convirtiendo documento firmado a XML: %v", err)
	}

	return xmlFirmado, nil
}

// FirmarDTE firma un DTE con el certificado digital
func FirmarDTE(xmlData []byte, certPath, keyPath string) ([]byte, error) {
	// Parsear XML original
	var doc models.DTEXMLModel
	if err := xml.Unmarshal(xmlData, &doc); err != nil {
		return nil, fmt.Errorf("error parseando XML: %v", err)
	}

	// Extraer información relevante para crear un DocumentoTributario temporal
	tempDoc := &models.DocumentoTributario{
		TipoDTE:      doc.Documento.Encabezado.IdDoc.TipoDTE,
		Folio:        doc.Documento.Encabezado.IdDoc.Folio,
		FechaEmision: time.Now(), // TODO: Parsear fecha desde XML
		MontoTotal:   float64(doc.Documento.Encabezado.Totales.MntTotal),
		Emisor: &models.Emisor{
			RUT: doc.Documento.Encabezado.Emisor.RUT,
		},
		Receptor: &models.Receptor{
			RUT: doc.Documento.Encabezado.Receptor.RUT,
		},
	}

	// Crear firma
	firma, err := crearFirmaXML(tempDoc)
	if err != nil {
		return nil, fmt.Errorf("error creando firma: %v", err)
	}

	// Insertar firma en XML
	xmlFirmado, err := insertarFirmaEnXML(xmlData, firma)
	if err != nil {
		return nil, fmt.Errorf("error insertando firma: %v", err)
	}

	return xmlFirmado, nil
}

// EnviarDTE envía un DTE al SII
func EnviarDTE(xmlData []byte, token, rutEmisor, rutEnvia string) ([]byte, error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequest(http.MethodPost, urlEnvio, bytes.NewBuffer(xmlData))
	if err != nil {
		return nil, fmt.Errorf("error creando request: %v", err)
	}

	req.Header.Set(contentType, applicationXML)
	req.Header.Set(userAgent, mozillaAgent)
	req.Header.Set(accept, applicationXML)
	req.Header.Set(authorization, "Bearer "+token)
	req.Header.Set(rutEmisor, rutEmisor)
	req.Header.Set(rutEnvia, rutEnvia)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error realizando request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error leyendo respuesta: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error en respuesta del SII: código %d, mensaje: %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// Funciones auxiliares
func crearClienteHTTP(certPath, keyPath string) (*http.Client, error) {
	// Cargar certificado y llave privada
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, err
	}

	// Configurar TLS
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
	}

	// Crear cliente HTTP con configuración TLS
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	return client, nil
}

// ProcesarRespuestaSII procesa la respuesta del SII y devuelve el estado
func ProcesarRespuestaSII(resp *models.RespuestaSII) (*models.EstadoSII, error) {
	if resp == nil {
		return nil, fmt.Errorf("respuesta nula")
	}

	estado := &models.EstadoSII{
		Codigo:      0,
		Descripcion: resp.Glosa,
		Timestamp:   resp.FechaProceso,
	}

	// Determinar el estado basado en la respuesta
	if resp.Estado == estadoOK {
		estado.Estado = "ACEPTADO"
	} else {
		estado.Estado = "RECHAZADO"
		estado.Codigo = 1 // Código genérico de error
	}

	// Agregar información adicional
	estado.Glosa = resp.Glosa
	estado.TrackID = resp.TrackID
	estado.Errores = resp.Errores

	return estado, nil
}

// ObtenerTrackID obtiene el TrackID de un documento XML
func ObtenerTrackID(xmlData []byte) (string, error) {
	var resp models.RespuestaSII
	if err := xml.Unmarshal(xmlData, &resp); err != nil {
		return "", fmt.Errorf("error al decodificar XML: %v", err)
	}

	if err := ValidarRespuestaSII(&resp); err != nil {
		return "", fmt.Errorf("error validando respuesta: %v", err)
	}

	if resp.TrackID == "" {
		return "", fmt.Errorf("trackID no encontrado en la respuesta")
	}

	return resp.TrackID, nil
}

// ValidarRespuestaSII valida la respuesta del SII
func ValidarRespuestaSII(resp *models.RespuestaSII) error {
	if resp == nil {
		return fmt.Errorf("respuesta nula")
	}

	if err := ValidarEstadoSII(resp); err != nil {
		return fmt.Errorf("error en estado: %v", err)
	}

	if err := ValidarFechasSII(resp); err != nil {
		return fmt.Errorf("error en fechas: %v", err)
	}

	// Verificar si hay errores en la respuesta SII
	if resp.Estado == estadoERROR || len(resp.Errores) > 0 {
		return fmt.Errorf("la respuesta contiene errores: %s", resp.Glosa)
	}

	if err := ValidarDetallesSII(resp); err != nil {
		return fmt.Errorf("error en detalles: %v", err)
	}

	return nil
}

// ValidarEstadoSII valida el estado de la respuesta del SII
func ValidarEstadoSII(resp *models.RespuestaSII) error {
	if resp.Estado == "" {
		return fmt.Errorf("estado no encontrado en la respuesta")
	}

	if resp.Estado != estadoOK && resp.Estado != estadoERROR {
		return fmt.Errorf("estado inválido: %s", resp.Estado)
	}

	return nil
}

// ValidarFechasSII valida las fechas de la respuesta del SII
func ValidarFechasSII(resp *models.RespuestaSII) error {
	if resp.FechaProceso.IsZero() {
		return fmt.Errorf("fecha no encontrada en la respuesta")
	}

	// Validar que la fecha no sea futura
	if resp.FechaProceso.After(time.Now()) {
		return fmt.Errorf("fecha futura no válida: %v", resp.FechaProceso)
	}

	// Validar que la fecha no sea muy antigua (más de 24 horas)
	if time.Since(resp.FechaProceso) > 24*time.Hour {
		return fmt.Errorf("fecha demasiado antigua: %v", resp.FechaProceso)
	}

	return nil
}

// ValidarErroresSII valida los errores de la respuesta del SII
func ValidarErroresSII(resp *models.RespuestaSII) error {
	if resp.Estado == estadoERROR && len(resp.Errores) == 0 {
		return fmt.Errorf("estado ERROR sin errores especificados")
	}

	for i, err := range resp.Errores {
		if err.Codigo == "" {
			return fmt.Errorf("error %d: código no encontrado", i+1)
		}
		if err.Descripcion == "" {
			return fmt.Errorf("error %d: descripción no encontrada", i+1)
		}
		if err.Detalle == "" {
			return fmt.Errorf("error %d: detalle no encontrado", i+1)
		}
	}

	return nil
}

// ValidarDetallesSII valida los detalles de la respuesta del SII
func ValidarDetallesSII(resp *models.RespuestaSII) error {
	// No hay campo detalles en la respuesta SII actual, por lo que esta función no hace nada
	return nil
}
