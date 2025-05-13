package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cursor/FMgo/config"
	"github.com/supabase-community/supabase-go"
)

// SupabaseService maneja las operaciones con Supabase
type SupabaseService struct {
	client *supabase.Client
	config *config.SupabaseConfig
}

// NewSupabaseService crea una nueva instancia del servicio Supabase
func NewSupabaseService(cfg *config.SupabaseConfig) *SupabaseService {
	return &SupabaseService{
		client: cfg.GetClient(),
		config: cfg,
	}
}

// SupabaseDocumento representa un documento tributario en Supabase
type SupabaseDocumento struct {
	ID          string    `json:"id"`
	Tipo        string    `json:"tipo"` // DTE, BOLETA, etc.
	RutEmisor   string    `json:"rut_emisor"`
	RutReceptor string    `json:"rut_receptor"`
	Folio       int       `json:"folio"`
	MontoTotal  float64   `json:"monto_total"`
	Estado      string    `json:"estado"`
	XML         string    `json:"xml"`
	PDF         string    `json:"pdf"`
	Firma       string    `json:"firma"`
	TED         string    `json:"ted"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// GuardarDocumento guarda un documento en Supabase
func (s *SupabaseService) GuardarDocumento(ctx context.Context, doc *SupabaseDocumento) error {
	// Convertir documento a JSON
	data, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("error al serializar documento: %v", err)
	}

	// Insertar en la tabla documentos
	var result map[string]interface{}
	err = s.client.DB.From("documentos").Insert(data).Execute(&result)
	if err != nil {
		return fmt.Errorf("error al guardar documento: %v", err)
	}

	return nil
}

// ObtenerDocumento obtiene un documento por su ID
func (s *SupabaseService) ObtenerDocumento(ctx context.Context, id string) (*SupabaseDocumento, error) {
	var doc SupabaseDocumento
	err := s.client.DB.From("documentos").Select("*").Eq("id", id).Single().Execute(&doc)
	if err != nil {
		return nil, fmt.Errorf("error al obtener documento: %v", err)
	}

	return &doc, nil
}

// ActualizarEstadoDocumento actualiza el estado de un documento
func (s *SupabaseService) ActualizarEstadoDocumento(ctx context.Context, id string, estado string) error {
	update := map[string]interface{}{
		"estado":     estado,
		"updated_at": time.Now(),
	}

	err := s.client.DB.From("documentos").Update(update).Eq("id", id).Execute(nil)
	if err != nil {
		return fmt.Errorf("error al actualizar estado: %v", err)
	}

	return nil
}

// ListarDocumentos obtiene una lista de documentos con filtros opcionales
func (s *SupabaseService) ListarDocumentos(ctx context.Context, filtros map[string]interface{}) ([]SupabaseDocumento, error) {
	query := s.client.DB.From("documentos").Select("*")

	// Aplicar filtros
	for key, value := range filtros {
		query = query.Eq(key, value)
	}

	var docs []SupabaseDocumento
	err := query.Execute(&docs)
	if err != nil {
		return nil, fmt.Errorf("error al listar documentos: %v", err)
	}

	return docs, nil
}

// GuardarCertificado guarda un certificado digital en Supabase
func (s *SupabaseService) GuardarCertificado(ctx context.Context, rut string, certificado []byte, llavePrivada []byte) error {
	data := map[string]interface{}{
		"rut":           rut,
		"certificado":   string(certificado),
		"llave_privada": string(llavePrivada),
		"created_at":    time.Now(),
		"updated_at":    time.Now(),
	}

	err := s.client.DB.From("certificados").Insert(data).Execute(nil)
	if err != nil {
		return fmt.Errorf("error al guardar certificado: %v", err)
	}

	return nil
}

// ObtenerCertificado obtiene el certificado digital de un contribuyente
func (s *SupabaseService) ObtenerCertificado(ctx context.Context, rut string) ([]byte, []byte, error) {
	var result struct {
		Certificado  string `json:"certificado"`
		LlavePrivada string `json:"llave_privada"`
	}

	err := s.client.DB.From("certificados").Select("*").Eq("rut", rut).Single().Execute(&result)
	if err != nil {
		return nil, nil, fmt.Errorf("error al obtener certificado: %v", err)
	}

	return []byte(result.Certificado), []byte(result.LlavePrivada), nil
}

// GuardarSesion guarda una sesión electrónica
func (s *SupabaseService) GuardarSesion(ctx context.Context, rut string, token string, expiracion time.Time) error {
	data := map[string]interface{}{
		"rut":        rut,
		"token":      token,
		"expiracion": expiracion,
		"created_at": time.Now(),
		"updated_at": time.Now(),
	}

	err := s.client.DB.From("sesiones").Insert(data).Execute(nil)
	if err != nil {
		return fmt.Errorf("error al guardar sesión: %v", err)
	}

	return nil
}

// VerificarSesion verifica si una sesión es válida
func (s *SupabaseService) VerificarSesion(ctx context.Context, token string) (bool, error) {
	var result struct {
		Expiracion time.Time `json:"expiracion"`
	}

	err := s.client.DB.From("sesiones").Select("expiracion").Eq("token", token).Single().Execute(&result)
	if err != nil {
		return false, fmt.Errorf("error al verificar sesión: %v", err)
	}

	return time.Now().Before(result.Expiracion), nil
}
