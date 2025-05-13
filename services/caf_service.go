package services

import (
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cursor/FMgo/config"
	"github.com/cursor/FMgo/domain"
	"github.com/cursor/FMgo/models"
	"github.com/cursor/FMgo/sii"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// CAFService maneja la lógica de negocio relacionada con los CAF
type CAFService struct {
	db         *mongo.Database
	cache      *redis.Client
	siiService sii.SIIService
	certFile   string
	keyFile    string
	config     *config.SupabaseConfig
}

// SIICAFRequest representa la solicitud de CAF al SII
type SIICAFRequest struct {
	RUTEmisor      string
	TipoDTE        string
	FolioInicial   int
	FolioFinal     int
	FechaSolicitud time.Time
}

// SIICAFResponse representa la respuesta del SII
type SIICAFResponse struct {
	Estado         string    `xml:"ESTADO"`
	Glosa          string    `xml:"GLOSA"`
	TrackID        string    `xml:"TRACKID"`
	FechaRespuesta time.Time `xml:"FECHARESPUESTA"`
	URLDescarga    string    `xml:"URLDESCARGA"`
}

// CAFMetadata contiene metadatos del archivo CAF
type CAFMetadata struct {
	RUTEmisor        string
	TipoDTE          string
	FolioInicial     int
	FolioFinal       int
	FechaEmision     time.Time
	FechaVencimiento time.Time
	Estado           string
	Hash             string
}

// SIICAFXML representa la estructura del archivo CAF del SII
type SIICAFXML struct {
	XMLName           xml.Name `xml:"AUTORIZACION"`
	Version           string   `xml:"CAF>version,attr"`
	RUTEmisor         string   `xml:"CAF>DA>RE"`
	RazonSocial       string   `xml:"CAF>DA>RS"`
	TipoDTE           string   `xml:"CAF>DA>TD"`
	FolioInicial      int      `xml:"CAF>DA>RNG>D"`
	FolioFinal        int      `xml:"CAF>DA>RNG>H"`
	FechaAutorizacion string   `xml:"CAF>DA>FA"`
	Modulo            string   `xml:"CAF>DA>RSAPK>M"`
	Exponente         string   `xml:"CAF>DA>RSAPK>E"`
	IDK               string   `xml:"CAF>DA>IDK"`
	Firma             string   `xml:"CAF>FRMA"`
	PrivateKey        string   `xml:"RSASK"`
	PublicKey         string   `xml:"RSAPUBK"`
}

// SIIError representa un error específico del SII
type SIIError struct {
	Codigo    string
	Mensaje   string
	Severidad string // ERROR, ADVERTENCIA, INFORMACION
	Solucion  string
}

// CAFError representa un error específico del servicio CAF
type CAFError struct {
	Codigo    string
	Mensaje   string
	Severidad string
	Detalles  map[string]interface{}
}

func (e *CAFError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Codigo, e.Mensaje)
}

// NewCAFService crea una nueva instancia del servicio CAF
func NewCAFService(
	db *mongo.Database,
	cache *redis.Client,
	siiService sii.SIIService,
	certFile string,
	keyFile string,
	config *config.SupabaseConfig,
) domain.CAFService {
	return &CAFService{
		db:         db,
		cache:      cache,
		siiService: siiService,
		certFile:   certFile,
		keyFile:    keyFile,
		config:     config,
	}
}

// ObtenerCAF obtiene un CAF por tipo de documento
func (s *CAFService) ObtenerCAF(ctx context.Context, tipoDocumento string) (*domain.CAF, error) {
	collection := s.db.Collection("cafs")
	var caf domain.CAF
	err := collection.FindOne(ctx, bson.M{
		"tipo_documento": tipoDocumento,
		"estado":         "ACTIVO",
	}).Decode(&caf)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("no hay CAF disponible para el tipo de documento")
		}
		return nil, err
	}
	return &caf, nil
}

// ValidarCAF valida un CAF
func (s *CAFService) ValidarCAF(ctx context.Context, caf *domain.CAF) error {
	if caf == nil {
		return errors.New("CAF no puede ser nulo")
	}

	// Validar rango de folios
	if caf.FolioActual < caf.RangoInicial || caf.FolioActual > caf.RangoFinal {
		return errors.New("folio actual fuera de rango")
	}

	// Validar vencimiento
	if err := s.VerificarVencimientoCAF(ctx, caf); err != nil {
		return err
	}

	return nil
}

// ActualizarFolioActual actualiza el folio actual de un CAF
func (s *CAFService) ActualizarFolioActual(ctx context.Context, caf *domain.CAF) error {
	collection := s.db.Collection("cafs")
	_, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": caf.ID},
		bson.M{
			"$set": bson.M{
				"folio_actual": caf.FolioActual + 1,
				"updated_at":   time.Now(),
			},
		},
	)
	return err
}

// VerificarVencimientoCAF verifica si un CAF está vencido
func (s *CAFService) VerificarVencimientoCAF(ctx context.Context, caf *domain.CAF) error {
	if caf.FechaVencimiento.Before(time.Now()) {
		// Actualizar estado del CAF
		collection := s.db.Collection("cafs")
		_, err := collection.UpdateOne(
			ctx,
			bson.M{"_id": caf.ID},
			bson.M{
				"$set": bson.M{
					"estado":     "VENCIDO",
					"updated_at": time.Now(),
				},
			},
		)
		if err != nil {
			return err
		}
		return errors.New("CAF vencido")
	}
	return nil
}

// SolicitarCAF solicita un nuevo CAF al SII
func (s *CAFService) SolicitarCAF(ctx context.Context, req *SIICAFRequest) (*SIICAFResponse, error) {
	// Verificar si ya existe una solicitud en proceso
	cacheKey := fmt.Sprintf("caf_request:%s:%s", req.RUTEmisor, req.TipoDTE)
	if cached, err := s.cache.Get(ctx, cacheKey).Result(); err == nil {
		var response SIICAFResponse
		if err := json.Unmarshal([]byte(cached), &response); err == nil {
			return &response, nil
		}
	}

	// Validar parámetros de la solicitud
	if err := s.validarSolicitud(req); err != nil {
		return nil, fmt.Errorf("error validando solicitud: %v", err)
	}

	// Construir XML de solicitud
	xmlRequest := s.buildCAFRequestXML(req)

	// Enviar solicitud al SII
	resp, err := s.sendRequest(ctx, xmlRequest)
	if err != nil {
		return nil, fmt.Errorf("error enviando solicitud CAF: %v", err)
	}

	// Guardar en cache
	if data, err := json.Marshal(resp); err == nil {
		s.cache.Set(ctx, cacheKey, data, 1*time.Hour)
	}

	return resp, nil
}

// ConsultarEstadoCAF consulta el estado de una solicitud de CAF
func (s *CAFService) ConsultarEstadoCAF(ctx context.Context, trackID string) (*SIICAFResponse, error) {
	// Verificar cache
	cacheKey := fmt.Sprintf("caf_status:%s", trackID)
	if cached, err := s.cache.Get(ctx, cacheKey).Result(); err == nil {
		var response SIICAFResponse
		if err := json.Unmarshal([]byte(cached), &response); err == nil {
			return &response, nil
		}
	}

	// Consultar estado en SII
	resp, err := s.siiService.ConsultarEstado(trackID)
	if err != nil {
		return nil, fmt.Errorf("error consultando estado CAF: %v", err)
	}

	// Convertir respuesta
	response := &SIICAFResponse{
		Estado:         resp.Estado,
		Glosa:          resp.Glosa,
		TrackID:        resp.TrackID,
		FechaRespuesta: time.Now(),
	}

	// Guardar en cache
	if data, err := json.Marshal(response); err == nil {
		s.cache.Set(ctx, cacheKey, data, 5*time.Minute)
	}

	return response, nil
}

// DescargarCAF descarga el archivo CAF del SII
func (s *CAFService) DescargarCAF(ctx context.Context, url string, savePath string, rutEmisor string) (*CAFMetadata, error) {
	// Crear directorio si no existe
	if err := os.MkdirAll(filepath.Dir(savePath), 0755); err != nil {
		return nil, fmt.Errorf("error creando directorio: %v", err)
	}

	// Crear archivo temporal
	tempFile := savePath + ".tmp"
	file, err := os.Create(tempFile)
	if err != nil {
		return nil, fmt.Errorf("error creando archivo temporal: %v", err)
	}
	defer file.Close()

	// Configurar request con certificado
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creando request: %v", err)
	}

	// Cargar certificado
	cert, err := tls.LoadX509KeyPair(s.certFile, s.keyFile)
	if err != nil {
		return nil, fmt.Errorf("error cargando certificado: %v", err)
	}

	// Configurar cliente con certificado
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates:       []tls.Certificate{cert},
				InsecureSkipVerify: s.config.Ambiente == "CERTIFICACION",
			},
		},
	}

	// Descargar archivo
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error descargando archivo: %v", err)
	}
	defer resp.Body.Close()

	// Verificar respuesta
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error en la descarga: %s", resp.Status)
	}

	// Copiar contenido
	if _, err := io.Copy(file, resp.Body); err != nil {
		return nil, fmt.Errorf("error guardando archivo: %v", err)
	}

	// Cerrar archivo temporal
	file.Close()

	// Validar archivo CAF
	metadata, err := s.validarArchivoCAF(tempFile, rutEmisor)
	if err != nil {
		os.Remove(tempFile)
		return nil, fmt.Errorf("error validando archivo CAF: %v", err)
	}

	// Renombrar archivo temporal a definitivo
	if err := os.Rename(tempFile, savePath); err != nil {
		os.Remove(tempFile)
		return nil, fmt.Errorf("error renombrando archivo: %v", err)
	}

	return metadata, nil
}

// validarSolicitud valida los parámetros de la solicitud
func (s *CAFService) validarSolicitud(req *SIICAFRequest) error {
	if req.RUTEmisor == "" {
		return fmt.Errorf("RUT emisor es requerido")
	}
	if req.TipoDTE == "" {
		return fmt.Errorf("tipo DTE es requerido")
	}
	if req.FolioInicial <= 0 {
		return fmt.Errorf("folio inicial debe ser mayor a 0")
	}
	if req.FolioFinal <= req.FolioInicial {
		return fmt.Errorf("folio final debe ser mayor al inicial")
	}
	if req.FolioFinal-req.FolioInicial > 10000 {
		return fmt.Errorf("rango de folios excede el máximo permitido")
	}
	return nil
}

// validarArchivoCAF valida el archivo CAF descargado
func (s *CAFService) validarArchivoCAF(filePath string, rutEmisor string) (*CAFMetadata, error) {
	// Abrir archivo
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error abriendo archivo: %v", err)
	}
	defer file.Close()

	// Leer contenido
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("error leyendo archivo: %v", err)
	}

	// Parsear XML
	var caf SIICAFXML
	if err := xml.Unmarshal(data, &caf); err != nil {
		return nil, fmt.Errorf("error parseando XML: %v", err)
	}

	// Validar RUT emisor
	if caf.RUTEmisor != rutEmisor {
		return nil, s.manejarErrorSII("002", "RUT emisor no coincide")
	}

	// Validar firma digital
	if err := s.validarFirmaCAF(&caf); err != nil {
		return nil, s.manejarErrorSII("003", "Firma digital inválida")
	}

	// Validar fechas
	fechaAutorizacion, err := time.Parse("2006-01-02", caf.FechaAutorizacion)
	if err != nil {
		return nil, fmt.Errorf("error parseando fecha de autorización: %v", err)
	}

	// Validar que el CAF no esté expirado (6 meses de vigencia)
	if time.Now().After(fechaAutorizacion.AddDate(0, 6, 0)) {
		return nil, s.manejarErrorSII("005", "CAF expirado")
	}

	// Validar rango de folios
	if caf.FolioFinal <= caf.FolioInicial {
		return nil, s.manejarErrorSII("004", "Rango de folios inválido")
	}

	// Validar que el rango de folios no se superponga con otros CAF
	if err := s.validarSuperposicionFolios(caf.RUTEmisor, caf.TipoDTE, caf.FolioInicial, caf.FolioFinal); err != nil {
		return nil, err
	}

	// Retornar metadatos
	return &CAFMetadata{
		RUTEmisor:        caf.RUTEmisor,
		TipoDTE:          caf.TipoDTE,
		FolioInicial:     caf.FolioInicial,
		FolioFinal:       caf.FolioFinal,
		FechaEmision:     fechaAutorizacion,
		FechaVencimiento: fechaAutorizacion.AddDate(0, 6, 0), // 6 meses de validez
		Estado:           "VALIDO",
		Hash:             s.calcularHashCAF(data),
	}, nil
}

// validarSuperposicionFolios verifica que el rango de folios no se superponga con otros CAF
func (s *CAFService) validarSuperposicionFolios(rutEmisor string, tipoDTE string, folioInicial int, folioFinal int) error {
	// TODO: Implementar verificación en base de datos
	// Por ahora retornamos nil
	return nil
}

// validarFirmaCAF valida la firma digital del CAF
func (s *CAFService) validarFirmaCAF(caf *SIICAFXML) error {
	// Decodificar clave pública
	block, _ := pem.Decode([]byte(caf.PublicKey))
	if block == nil {
		return fmt.Errorf("error decodificando clave pública")
	}

	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("error parseando clave pública: %v", err)
	}

	rsaPubKey, ok := pubKey.(*rsa.PublicKey)
	if !ok {
		return fmt.Errorf("clave pública no es RSA")
	}

	// Decodificar firma
	signature, err := base64.StdEncoding.DecodeString(caf.Firma)
	if err != nil {
		return fmt.Errorf("error decodificando firma: %v", err)
	}

	// Calcular hash del documento
	hash := sha1.New()
	hash.Write([]byte(fmt.Sprintf("%s%s%s%d%d%s%s%s",
		caf.RUTEmisor,
		caf.RazonSocial,
		caf.TipoDTE,
		caf.FolioInicial,
		caf.FolioFinal,
		caf.FechaAutorizacion,
		caf.Modulo,
		caf.Exponente)))

	// Verificar firma
	if err := rsa.VerifyPKCS1v15(rsaPubKey, crypto.SHA1, hash.Sum(nil), signature); err != nil {
		return fmt.Errorf("error verificando firma: %v", err)
	}

	return nil
}

// calcularHashCAF calcula el hash del archivo CAF
func (s *CAFService) calcularHashCAF(data []byte) string {
	hash := sha1.New()
	hash.Write(data)
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}

// manejarErrorSII maneja errores específicos del SII
func (s *CAFService) manejarErrorSII(codigo string, mensaje string) error {
	return &models.ErrorSII{
		Codigo:      codigo,
		Descripcion: mensaje,
	}
}

// buildCAFRequestXML construye el XML de solicitud de CAF
func (s *CAFService) buildCAFRequestXML(req *SIICAFRequest) string {
	return fmt.Sprintf(`
		<CAFRequest>
			<RUTEmisor>%s</RUTEmisor>
			<TipoDTE>%s</TipoDTE>
			<FolioInicial>%d</FolioInicial>
			<FolioFinal>%d</FolioFinal>
			<FechaSolicitud>%s</FechaSolicitud>
			<Ambiente>%s</Ambiente>
		</CAFRequest>
	`, req.RUTEmisor, req.TipoDTE, req.FolioInicial, req.FolioFinal,
		req.FechaSolicitud.Format(time.RFC3339), s.config.Ambiente)
}

// sendRequest envía la solicitud al SII
func (s *CAFService) sendRequest(ctx context.Context, xmlRequest string) (*SIICAFResponse, error) {
	// Crear request
	req, err := http.NewRequestWithContext(ctx, "POST", s.config.BaseURL+"/caf/request", strings.NewReader(xmlRequest))
	if err != nil {
		return nil, fmt.Errorf("error creando request: %v", err)
	}

	// Configurar headers
	req.Header.Set("Content-Type", "application/xml")
	req.Header.Set("Authorization", "Bearer "+s.config.Token)

	// Cargar certificado
	cert, err := tls.LoadX509KeyPair(s.certFile, s.keyFile)
	if err != nil {
		return nil, fmt.Errorf("error cargando certificado: %v", err)
	}

	// Configurar cliente con certificado
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates:       []tls.Certificate{cert},
				InsecureSkipVerify: s.config.Ambiente == "CERTIFICACION",
			},
		},
	}

	// Enviar request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error enviando request: %v", err)
	}
	defer resp.Body.Close()

	// Leer respuesta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error leyendo respuesta: %v", err)
	}

	// Parsear respuesta XML
	var response SIICAFResponse
	if err := xml.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("error parseando respuesta: %v", err)
	}

	return &response, nil
}

// VerificarFoliosDisponibles verifica los folios disponibles
func (s *CAFService) VerificarFoliosDisponibles(ctx context.Context, rutEmisor string, tipoDTE string) (int, error) {
	// Implementar verificación de folios disponibles
	return 0, nil
}

// ProgramarSolicitudCAF programa una solicitud automática de CAF cuando los folios están por agotarse
func (s *CAFService) ProgramarSolicitudCAF(ctx context.Context, rutEmisor string, tipoDTE string, umbralFolios int) error {
	// Verificar si ya existe una solicitud programada
	cacheKey := fmt.Sprintf("caf_programado:%s:%s", rutEmisor, tipoDTE)
	if _, err := s.cache.Get(ctx, cacheKey).Result(); err == nil {
		return &CAFError{
			Codigo:    "CAF008",
			Mensaje:   "Ya existe una solicitud programada para este CAF",
			Severidad: "ADVERTENCIA",
			Detalles: map[string]interface{}{
				"rut_emisor": rutEmisor,
				"tipo_dte":   tipoDTE,
				"umbral":     umbralFolios,
			},
		}
	}

	// Obtener folios disponibles
	foliosDisponibles, err := s.VerificarFoliosDisponibles(ctx, rutEmisor, tipoDTE)
	if err != nil {
		return err
	}

	// Si los folios están por debajo del umbral, programar solicitud
	if foliosDisponibles <= umbralFolios {
		// Crear solicitud
		req := &SIICAFRequest{
			RUTEmisor:      rutEmisor,
			TipoDTE:        tipoDTE,
			FolioInicial:   1,
			FolioFinal:     1000, // Valor por defecto, ajustar según necesidades
			FechaSolicitud: time.Now(),
		}

		// Enviar solicitud
		_, err := s.SolicitarCAF(ctx, req)
		if err != nil {
			return err
		}

		// Marcar como programado en cache
		s.cache.Set(ctx, cacheKey, "true", 24*time.Hour)
	}

	return nil
}

// MonitorearEstadoCAF monitorea el estado de una solicitud de CAF
func (s *CAFService) MonitorearEstadoCAF(ctx context.Context, trackID string, interval time.Duration) (*SIICAFResponse, error) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-ticker.C:
			estado, err := s.ConsultarEstadoCAF(ctx, trackID)
			if err != nil {
				return nil, err
			}

			if estado.Estado == "ACEPTADO" {
				return estado, nil
			}

			if estado.Estado == "RECHAZADO" {
				return nil, fmt.Errorf("solicitud rechazada: %s", estado.Glosa)
			}
		}
	}
}
