package supabase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/cursor/FMgo/models"
	"github.com/cursor/FMgo/supabase"
)

// DocumentoTributarioRepository define las operaciones de repositorio para documentos tributarios
type DocumentoTributarioRepository interface {
	Create(ctx context.Context, doc *models.DocumentoTributario) error
	GetByID(ctx context.Context, id string) (*models.DocumentoTributario, error)
	GetByFolio(ctx context.Context, tipo string, folio int) (*models.DocumentoTributario, error)
	UpdateEstado(ctx context.Context, id string, estado models.EstadoDTE) error
	UpdateTrackID(ctx context.Context, id string, trackID string) error
	List(ctx context.Context, filter map[string]interface{}, limit, offset int) ([]*models.DocumentoTributario, error)
	Count(ctx context.Context, filter map[string]interface{}) (int, error)
}

// SupabaseDocumentoRepository implementa DocumentoTributarioRepository usando Supabase
type SupabaseDocumentoRepository struct {
	client     *supabase.Client
	baseURL    string
	apiKey     string
	tableName  string
	httpClient *http.Client
}

// NewSupabaseDocumentoRepository crea una nueva instancia de SupabaseDocumentoRepository
func NewSupabaseDocumentoRepository(client *supabase.Client) DocumentoTributarioRepository {
	config := client.GetConfig()
	return &SupabaseDocumentoRepository{
		client:    client,
		baseURL:   config.Supabase.URL,
		apiKey:    config.Supabase.AnonKey,
		tableName: config.Supabase.TablaDocumentos,
		httpClient: &http.Client{
			Timeout: time.Duration(config.Supabase.Timeout) * time.Second,
		},
	}
}

// Create guarda un nuevo documento
func (r *SupabaseDocumentoRepository) Create(ctx context.Context, doc *models.DocumentoTributario) error {
	// Asegurarse de que el documento tenga fechas asignadas
	now := time.Now()

	// Convertir a formato adecuado para Supabase
	documentoSupabase := map[string]interface{}{
		"tipo":         doc.TipoDTE,
		"rut_emisor":   doc.RUTEmisor,
		"rut_receptor": doc.RUTReceptor,
		"folio":        doc.Folio,
		"monto_total":  doc.MontoTotal,
		"monto_neto":   doc.MontoNeto,
		"monto_iva":    doc.MontoIVA,
		"estado":       string(doc.Estado),
		"created_at":   now,
		"updated_at":   now,
	}

	// Convertir XML a string si existe
	if doc.XML != "" {
		documentoSupabase["xml"] = doc.XML
	}

	// Convertir PDF a string si existe
	if len(doc.PDFData) > 0 {
		documentoSupabase["pdf"] = string(doc.PDFData)
	}

	// Agregar trackID si existe
	if doc.TrackID != "" {
		documentoSupabase["track_id"] = doc.TrackID
	}

	// Convertir a JSON para el cuerpo de la solicitud
	jsonData, err := json.Marshal(documentoSupabase)
	if err != nil {
		return fmt.Errorf("error al serializar documento: %w", err)
	}

	// Construir la URL de la API
	apiURL := fmt.Sprintf("%s/rest/v1/%s", r.baseURL, r.tableName)

	// Crear la solicitud HTTP
	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error al crear la solicitud: %w", err)
	}

	// Agregar encabezados requeridos
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", r.apiKey)
	req.Header.Set("Authorization", "Bearer "+r.apiKey)
	req.Header.Set("Prefer", "return=representation")

	// Enviar la solicitud
	resp, err := r.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error al enviar la solicitud: %w", err)
	}
	defer resp.Body.Close()

	// Verificar el código de estado
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("error en la respuesta de Supabase: código %d, respuesta: %s", resp.StatusCode, string(body))
	}

	// Decodificar la respuesta para obtener el ID asignado
	var result []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("error al decodificar la respuesta: %w", err)
	}

	// Obtener el ID asignado al documento
	if len(result) > 0 {
		if id, ok := result[0]["id"].(string); ok {
			doc.ID = id
		}
	}

	return nil
}

// GetByID obtiene un documento por su ID
func (r *SupabaseDocumentoRepository) GetByID(ctx context.Context, id string) (*models.DocumentoTributario, error) {
	// Construir la URL de la API
	apiURL := fmt.Sprintf("%s/rest/v1/%s?id=eq.%s&select=*", r.baseURL, r.tableName, url.QueryEscape(id))

	// Crear la solicitud HTTP
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error al crear la solicitud: %w", err)
	}

	// Agregar encabezados requeridos
	req.Header.Set("apikey", r.apiKey)
	req.Header.Set("Authorization", "Bearer "+r.apiKey)

	// Enviar la solicitud
	resp, err := r.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error al enviar la solicitud: %w", err)
	}
	defer resp.Body.Close()

	// Verificar el código de estado
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("error en la respuesta de Supabase: código %d, respuesta: %s", resp.StatusCode, string(body))
	}

	// Decodificar la respuesta
	var result []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error al decodificar la respuesta: %w", err)
	}

	// Verificar que se encontró un documento
	if len(result) == 0 {
		return nil, fmt.Errorf("no se encontró ningún documento con ID %s", id)
	}

	// Convertir a modelo de documento
	documento, err := mapToDocumentoTributario(result[0])
	if err != nil {
		return nil, fmt.Errorf("error al convertir documento: %w", err)
	}

	return documento, nil
}

// GetByFolio obtiene un documento por su tipo y folio
func (r *SupabaseDocumentoRepository) GetByFolio(ctx context.Context, tipo string, folio int) (*models.DocumentoTributario, error) {
	// Construir la URL de la API
	apiURL := fmt.Sprintf("%s/rest/v1/%s?tipo=eq.%s&folio=eq.%d&select=*",
		r.baseURL, r.tableName, url.QueryEscape(tipo), folio)

	// Crear la solicitud HTTP
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error al crear la solicitud: %w", err)
	}

	// Agregar encabezados requeridos
	req.Header.Set("apikey", r.apiKey)
	req.Header.Set("Authorization", "Bearer "+r.apiKey)

	// Enviar la solicitud
	resp, err := r.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error al enviar la solicitud: %w", err)
	}
	defer resp.Body.Close()

	// Verificar el código de estado
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("error en la respuesta de Supabase: código %d, respuesta: %s", resp.StatusCode, string(body))
	}

	// Decodificar la respuesta
	var result []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error al decodificar la respuesta: %w", err)
	}

	// Verificar que se encontró un documento
	if len(result) == 0 {
		return nil, fmt.Errorf("no se encontró ningún documento con tipo %s y folio %d", tipo, folio)
	}

	// Convertir a modelo de documento
	documento, err := mapToDocumentoTributario(result[0])
	if err != nil {
		return nil, fmt.Errorf("error al convertir documento: %w", err)
	}

	return documento, nil
}

// UpdateEstado actualiza el estado de un documento
func (r *SupabaseDocumentoRepository) UpdateEstado(ctx context.Context, id string, estado models.EstadoDTE) error {
	// Preparar los datos de actualización
	update := map[string]interface{}{
		"estado":     string(estado),
		"updated_at": time.Now(),
	}

	// Convertir a JSON para el cuerpo de la solicitud
	jsonData, err := json.Marshal(update)
	if err != nil {
		return fmt.Errorf("error al serializar los datos: %w", err)
	}

	// Construir la URL de la API
	apiURL := fmt.Sprintf("%s/rest/v1/%s?id=eq.%s", r.baseURL, r.tableName, url.QueryEscape(id))

	// Crear la solicitud HTTP
	req, err := http.NewRequestWithContext(ctx, "PATCH", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error al crear la solicitud: %w", err)
	}

	// Agregar encabezados requeridos
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", r.apiKey)
	req.Header.Set("Authorization", "Bearer "+r.apiKey)
	req.Header.Set("Prefer", "return=minimal")

	// Enviar la solicitud
	resp, err := r.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error al enviar la solicitud: %w", err)
	}
	defer resp.Body.Close()

	// Verificar el código de estado
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("error en la respuesta de Supabase: código %d, respuesta: %s", resp.StatusCode, string(body))
	}

	return nil
}

// UpdateTrackID actualiza el track ID de un documento
func (r *SupabaseDocumentoRepository) UpdateTrackID(ctx context.Context, id string, trackID string) error {
	// Preparar los datos de actualización
	update := map[string]interface{}{
		"track_id":   trackID,
		"updated_at": time.Now(),
	}

	// Convertir a JSON para el cuerpo de la solicitud
	jsonData, err := json.Marshal(update)
	if err != nil {
		return fmt.Errorf("error al serializar los datos: %w", err)
	}

	// Construir la URL de la API
	apiURL := fmt.Sprintf("%s/rest/v1/%s?id=eq.%s", r.baseURL, r.tableName, url.QueryEscape(id))

	// Crear la solicitud HTTP
	req, err := http.NewRequestWithContext(ctx, "PATCH", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error al crear la solicitud: %w", err)
	}

	// Agregar encabezados requeridos
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", r.apiKey)
	req.Header.Set("Authorization", "Bearer "+r.apiKey)
	req.Header.Set("Prefer", "return=minimal")

	// Enviar la solicitud
	resp, err := r.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error al enviar la solicitud: %w", err)
	}
	defer resp.Body.Close()

	// Verificar el código de estado
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("error en la respuesta de Supabase: código %d, respuesta: %s", resp.StatusCode, string(body))
	}

	return nil
}

// List obtiene una lista de documentos con filtros
func (r *SupabaseDocumentoRepository) List(ctx context.Context, filter map[string]interface{}, limit, offset int) ([]*models.DocumentoTributario, error) {
	// Construir la URL base de la API
	apiURL := fmt.Sprintf("%s/rest/v1/%s?select=*", r.baseURL, r.tableName)

	// Agregar filtros a la URL
	for key, value := range filter {
		// Convertir valor a string para la URL
		var strValue string
		switch v := value.(type) {
		case string:
			strValue = v
		case int:
			strValue = fmt.Sprintf("%d", v)
		case float64:
			strValue = fmt.Sprintf("%.2f", v)
		case bool:
			if v {
				strValue = "true"
			} else {
				strValue = "false"
			}
		default:
			continue // Omitir filtros no compatibles
		}
		apiURL += fmt.Sprintf("&%s=eq.%s", url.QueryEscape(key), url.QueryEscape(strValue))
	}

	// Agregar paginación
	if limit > 0 {
		apiURL += fmt.Sprintf("&limit=%d", limit)
	}
	if offset > 0 {
		apiURL += fmt.Sprintf("&offset=%d", offset)
	}

	// Ordenar por fecha de creación
	apiURL += "&order=created_at.desc"

	// Crear la solicitud HTTP
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error al crear la solicitud: %w", err)
	}

	// Agregar encabezados requeridos
	req.Header.Set("apikey", r.apiKey)
	req.Header.Set("Authorization", "Bearer "+r.apiKey)

	// Enviar la solicitud
	resp, err := r.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error al enviar la solicitud: %w", err)
	}
	defer resp.Body.Close()

	// Verificar el código de estado
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("error en la respuesta de Supabase: código %d, respuesta: %s", resp.StatusCode, string(body))
	}

	// Decodificar la respuesta
	var result []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error al decodificar la respuesta: %w", err)
	}

	// Convertir a modelos de documentos
	documentos := make([]*models.DocumentoTributario, 0, len(result))
	for _, item := range result {
		documento, err := mapToDocumentoTributario(item)
		if err != nil {
			continue // Ignorar documentos que no se puedan convertir
		}
		documentos = append(documentos, documento)
	}

	return documentos, nil
}

// Count cuenta documentos con filtros
func (r *SupabaseDocumentoRepository) Count(ctx context.Context, filter map[string]interface{}) (int, error) {
	// Construir la URL base de la API con la función de conteo
	apiURL := fmt.Sprintf("%s/rest/v1/%s?select=count", r.baseURL, r.tableName)

	// Agregar filtros a la URL
	for key, value := range filter {
		// Convertir valor a string para la URL
		var strValue string
		switch v := value.(type) {
		case string:
			strValue = v
		case int:
			strValue = fmt.Sprintf("%d", v)
		case float64:
			strValue = fmt.Sprintf("%.2f", v)
		case bool:
			if v {
				strValue = "true"
			} else {
				strValue = "false"
			}
		default:
			continue // Omitir filtros no compatibles
		}
		apiURL += fmt.Sprintf("&%s=eq.%s", url.QueryEscape(key), url.QueryEscape(strValue))
	}

	// Crear la solicitud HTTP
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return 0, fmt.Errorf("error al crear la solicitud: %w", err)
	}

	// Agregar encabezados requeridos
	req.Header.Set("apikey", r.apiKey)
	req.Header.Set("Authorization", "Bearer "+r.apiKey)
	req.Header.Set("Prefer", "count=exact")

	// Enviar la solicitud
	resp, err := r.httpClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("error al enviar la solicitud: %w", err)
	}
	defer resp.Body.Close()

	// Verificar el código de estado
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return 0, fmt.Errorf("error en la respuesta de Supabase: código %d, respuesta: %s", resp.StatusCode, string(body))
	}

	// Obtener el total del header de respuesta
	countHeader := resp.Header.Get("Content-Range")
	var count int
	if countHeader != "" {
		_, err := fmt.Sscanf(countHeader, "0-%d/%d", &count, &count)
		if err != nil {
			// Si no se puede parsear, usar una aproximación
			docs, err := r.List(ctx, filter, 1000000, 0)
			if err != nil {
				return 0, err
			}
			return len(docs), nil
		}
	} else {
		// Si no hay header, hacer una aproximación
		docs, err := r.List(ctx, filter, 1000000, 0)
		if err != nil {
			return 0, err
		}
		return len(docs), nil
	}

	return count, nil
}

// mapToDocumentoTributario convierte un mapa de Supabase a DocumentoTributario
func mapToDocumentoTributario(data map[string]interface{}) (*models.DocumentoTributario, error) {
	// Convertir el mapa a JSON para facilitar la conversión
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("error al serializar datos: %w", err)
	}

	// Crear documento base
	doc := &models.DocumentoTributario{}

	// Deserializar JSON a estructura
	if err := json.Unmarshal(jsonData, doc); err != nil {
		return nil, fmt.Errorf("error al deserializar documento: %w", err)
	}

	// Tratamiento especial para campos que pueden ser problemáticos
	if id, ok := data["id"].(string); ok {
		doc.ID = id
	}

	if tipo, ok := data["tipo"].(string); ok {
		doc.TipoDTE = tipo
	}

	if estado, ok := data["estado"].(string); ok {
		doc.Estado = models.EstadoDTE(estado)
	}

	if xml, ok := data["xml"].(string); ok && xml != "" {
		doc.XML = xml
	}

	if pdf, ok := data["pdf"].(string); ok && pdf != "" {
		doc.PDFData = []byte(pdf)
	}

	if trackID, ok := data["track_id"].(string); ok {
		doc.TrackID = trackID
	}

	return doc, nil
}
