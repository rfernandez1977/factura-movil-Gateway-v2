package supabase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/fmgo/config"
)

// GetClientFromConfig crea un cliente de Supabase a partir del archivo de configuración
func GetClientFromConfig(configPath string) (*Client, error) {
	cfg, err := config.Load(configPath)
	if err != nil {
		return nil, fmt.Errorf("error al cargar la configuración: %w", err)
	}

	client, err := NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("error al crear cliente de Supabase: %w", err)
	}

	return client, nil
}

// VerifyConnection verifica la conexión con Supabase
func VerifyConnection(ctx context.Context, client *Client) error {
	return client.Ping(ctx)
}

// ListAllTables lista todas las tablas disponibles en Supabase
func ListAllTables(ctx context.Context, client *Client) ([]string, error) {
	// Lista predefinida de tablas que sabemos que existen
	tables := []string{
		"empresas",
		"documentos",
		"certificados",
		"sesiones",
		"cafs",
		"xml_files",
		"sobres_xml",
	}
	return tables, nil
}

// GetTableCount obtiene el número de registros en una tabla
func GetTableCount(ctx context.Context, client *Client, table string) (int, error) {
	url := fmt.Sprintf("%s/rest/v1/%s?select=count", client.config.Supabase.URL, table)

	req, err := http.NewRequestWithContext(ctx, "HEAD", url, nil)
	if err != nil {
		return 0, fmt.Errorf("error creando solicitud: %w", err)
	}

	req.Header.Set("apikey", client.config.Supabase.AnonKey)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", client.config.Supabase.AnonKey))
	req.Header.Set("Prefer", "count=exact")

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("error verificando tabla: %w", err)
	}
	defer resp.Body.Close()

	contentRange := resp.Header.Get("Content-Range")
	if contentRange == "" {
		return 0, fmt.Errorf("no se pudo obtener el conteo")
	}

	var count int
	_, err = fmt.Sscanf(contentRange, "%*d-%*d/%d", &count)
	if err != nil {
		return 0, fmt.Errorf("error parseando conteo: %w", err)
	}

	return count, nil
}

// ListTableData obtiene los datos de una tabla (con límite)
func ListTableData(ctx context.Context, client *Client, table string, limit int) ([]map[string]interface{}, error) {
	url := fmt.Sprintf("%s/rest/v1/%s?select=*&limit=%d", client.config.Supabase.URL, table, limit)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creando solicitud: %w", err)
	}

	req.Header.Set("apikey", client.config.Supabase.AnonKey)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", client.config.Supabase.AnonKey))

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo datos: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error leyendo respuesta: %w", err)
	}

	var data []map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("error decodificando JSON: %w", err)
	}

	return data, nil
}

// InsertRecord inserta un nuevo registro en una tabla
func InsertRecord(ctx context.Context, client *Client, table string, data map[string]interface{}) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/rest/v1/%s", client.config.Supabase.URL, table)

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("error serializando datos: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creando solicitud: %w", err)
	}

	req.Header.Set("apikey", client.config.Supabase.AnonKey)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", client.config.Supabase.AnonKey))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Prefer", "return=representation")

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error insertando datos: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("error insertando datos (código %d): %s", resp.StatusCode, body)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error leyendo respuesta: %w", err)
	}

	var result []map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("error decodificando respuesta: %w", err)
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("no se recibió respuesta después de la inserción")
	}

	return result[0], nil
}

// UpdateRecord actualiza un registro existente
func UpdateRecord(ctx context.Context, client *Client, table, id string, data map[string]interface{}) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/rest/v1/%s?id=eq.%s", client.config.Supabase.URL, table, id)

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("error serializando datos: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "PATCH", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creando solicitud: %w", err)
	}

	req.Header.Set("apikey", client.config.Supabase.AnonKey)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", client.config.Supabase.AnonKey))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Prefer", "return=representation")

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error actualizando datos: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("error actualizando datos (código %d): %s", resp.StatusCode, body)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error leyendo respuesta: %w", err)
	}

	var result []map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("error decodificando respuesta: %w", err)
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("no se recibió respuesta después de la actualización")
	}

	return result[0], nil
}

// GetRecordByID obtiene un registro por su ID
func GetRecordByID(ctx context.Context, client *Client, table, id string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/rest/v1/%s?id=eq.%s", client.config.Supabase.URL, table, id)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creando solicitud: %w", err)
	}

	req.Header.Set("apikey", client.config.Supabase.AnonKey)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", client.config.Supabase.AnonKey))

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo datos: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error leyendo respuesta: %w", err)
	}

	var data []map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("error decodificando JSON: %w", err)
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("registro no encontrado")
	}

	return data[0], nil
}

// DeleteRecord elimina un registro por su ID
func DeleteRecord(ctx context.Context, client *Client, table, id string) error {
	url := fmt.Sprintf("%s/rest/v1/%s?id=eq.%s", client.config.Supabase.URL, table, id)

	req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("error creando solicitud: %w", err)
	}

	req.Header.Set("apikey", client.config.Supabase.AnonKey)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", client.config.Supabase.AnonKey))

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error eliminando datos: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("error eliminando datos (código %d): %s", resp.StatusCode, body)
	}

	return nil
}

// QueryRecords busca registros por filtros personalizados
func QueryRecords(ctx context.Context, client *Client, table string, filters map[string]string, limit int) ([]map[string]interface{}, error) {
	url := fmt.Sprintf("%s/rest/v1/%s?select=*&limit=%d", client.config.Supabase.URL, table, limit)

	// Añadir filtros a la URL
	for k, v := range filters {
		url = fmt.Sprintf("%s&%s=eq.%s", url, k, v)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creando solicitud: %w", err)
	}

	req.Header.Set("apikey", client.config.Supabase.AnonKey)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", client.config.Supabase.AnonKey))

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo datos: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error leyendo respuesta: %w", err)
	}

	var data []map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("error decodificando JSON: %w", err)
	}

	return data, nil
}

// InitClientWithConfig inicializa un cliente con la configuración y verifica la conexión
func InitClientWithConfig(configPath string) (*Client, error) {
	log.Printf("Inicializando cliente Supabase con configuración desde: %s", configPath)

	// Cargar la configuración
	client, err := GetClientFromConfig(configPath)
	if err != nil {
		return nil, err
	}

	// Verificar la conexión
	if err := VerifyConnection(context.Background(), client); err != nil {
		return nil, fmt.Errorf("error verificando conexión a Supabase: %w", err)
	}

	log.Println("Conexión a Supabase establecida correctamente")
	return client, nil
}
