package services

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/cursor/FMgo/config"
	"github.com/cursor/FMgo/models"
	"github.com/cursor/FMgo/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// SesionService maneja las sesiones electrónicas con el SII
type SesionService struct {
	config *config.SupabaseConfig
	db     *mongo.Database
	client *SesionSIIClient
}

// NewSesionService crea una nueva instancia del servicio de sesiones
func NewSesionService(config *config.SupabaseConfig, db *mongo.Database) *SesionService {
	return &SesionService{
		config: config,
		db:     db,
		client: NewSesionSIIClient(config),
	}
}

// IniciarSesion inicia una nueva sesión electrónica
func (s *SesionService) IniciarSesion(ctx context.Context, empresa *models.Empresa) (*models.SesionElectronica, error) {
	// Verificar si ya existe una sesión activa
	sesion, err := s.ObtenerSesionActiva(ctx, empresa.ID)
	if err == nil && sesion != nil && sesion.IsValid() {
		return sesion, nil
	}

	// Iniciar nueva sesión con el SII
	resp, err := s.client.IniciarSesion(ctx, empresa)
	if err != nil {
		return nil, fmt.Errorf("error iniciando sesión con SII: %v", err)
	}

	// Crear nueva sesión
	sesion = &models.SesionElectronica{
		EmpresaID:       empresa.ID,
		Token:           resp.Token,
		Estado:          resp.Estado,
		FechaInicio:     time.Now(),
		FechaExpiracion: resp.FechaExpiracion,
		Intentos:        0,
		UltimoAcceso:    time.Now(),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Guardar sesión en la base de datos
	_, err = s.db.Collection("sesiones").InsertOne(ctx, sesion)
	if err != nil {
		return nil, fmt.Errorf("error guardando sesión: %v", err)
	}

	return sesion, nil
}

// CerrarSesion cierra una sesión electrónica
func (s *SesionService) CerrarSesion(ctx context.Context, sesionID string) error {
	// Obtener sesión
	var sesion models.SesionElectronica
	err := s.db.Collection("sesiones").FindOne(ctx, bson.M{"_id": sesionID}).Decode(&sesion)
	if err != nil {
		return fmt.Errorf("error obteniendo sesión: %v", err)
	}

	// Cerrar sesión en el SII
	err = s.client.CerrarSesion(ctx, sesion.Token)
	if err != nil {
		return fmt.Errorf("error cerrando sesión en SII: %v", err)
	}

	// Actualizar estado en la base de datos
	_, err = s.db.Collection("sesiones").UpdateOne(
		ctx,
		bson.M{"_id": sesionID},
		bson.M{
			"$set": bson.M{
				"estado":     "CERRADA",
				"updated_at": time.Now(),
			},
		},
	)
	if err != nil {
		return fmt.Errorf("error actualizando sesión: %v", err)
	}

	return nil
}

// ObtenerSesionActiva obtiene la sesión activa de una empresa
func (s *SesionService) ObtenerSesionActiva(ctx context.Context, empresaID string) (*models.SesionElectronica, error) {
	var sesion models.SesionElectronica
	err := s.db.Collection("sesiones").FindOne(
		ctx,
		bson.M{
			"empresa_id": empresaID,
			"estado":     "ACTIVA",
			"fecha_expiracion": bson.M{
				"$gt": time.Now(),
			},
		},
	).Decode(&sesion)

	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error obteniendo sesión activa: %v", err)
	}

	return &sesion, nil
}

// VerificarEstadoSesion verifica el estado de una sesión
func (s *SesionService) VerificarEstadoSesion(ctx context.Context, sesionID string) (*models.EstadoSesionInfo, error) {
	// Obtener sesión
	var sesion models.SesionElectronica
	err := s.db.Collection("sesiones").FindOne(ctx, bson.M{"_id": sesionID}).Decode(&sesion)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo sesión: %v", err)
	}

	// Verificar estado en el SII
	estado, err := s.client.VerificarEstadoSesion(ctx, sesion.Token)
	if err != nil {
		return nil, fmt.Errorf("error verificando estado de sesión: %v", err)
	}

	// Actualizar estado en la base de datos si es necesario
	if estado.Estado != sesion.Estado {
		_, err = s.db.Collection("sesiones").UpdateOne(
			ctx,
			bson.M{"_id": sesionID},
			bson.M{
				"$set": bson.M{
					"estado":     estado.Estado,
					"updated_at": time.Now(),
				},
			},
		)
		if err != nil {
			return nil, fmt.Errorf("error actualizando estado de sesión: %v", err)
		}
	}

	return estado, nil
}

// RenovarSesion renueva una sesión expirada
func (s *SesionService) RenovarSesion(ctx context.Context, sesionID string) (*models.SesionElectronica, error) {
	// Obtener sesión
	var sesion models.SesionElectronica
	err := s.db.Collection("sesiones").FindOne(ctx, bson.M{"_id": sesionID}).Decode(&sesion)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo sesión: %v", err)
	}

	// Obtener empresa
	var empresa models.Empresa
	err = s.db.Collection("empresas").FindOne(ctx, bson.M{"_id": sesion.EmpresaID}).Decode(&empresa)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo empresa: %v", err)
	}

	// Iniciar nueva sesión
	return s.IniciarSesion(ctx, &empresa)
}

// LimpiarSesionesExpiradas limpia las sesiones expiradas
func (s *SesionService) LimpiarSesionesExpiradas(ctx context.Context) error {
	_, err := s.db.Collection("sesiones").UpdateMany(
		ctx,
		bson.M{
			"fecha_expiracion": bson.M{
				"$lt": time.Now(),
			},
			"estado": "ACTIVA",
		},
		bson.M{
			"$set": bson.M{
				"estado":     "EXPIRADA",
				"updated_at": time.Now(),
			},
		},
	)
	if err != nil {
		return fmt.Errorf("error limpiando sesiones expiradas: %v", err)
	}

	return nil
}

// VerificarToken verifica si un token de sesión es válido
func (s *SesionService) VerificarToken(token string) (bool, error) {
	// Crear request HTTP
	req, err := http.NewRequest("GET", utils.GetSiiEndpoint(s.config)+"/sesion/verificar", nil)
	if err != nil {
		return false, fmt.Errorf("error al crear request: %v", err)
	}

	// Configurar headers
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("User-Agent", "FMgo/1.0")

	// Enviar request usando el cliente HTTP
	httpClient := &http.Client{
		Timeout: time.Second * 30,
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("error al enviar request: %v", err)
	}
	defer resp.Body.Close()

	// Verificar código de respuesta
	return resp.StatusCode == http.StatusOK, nil
}

// VerificarSesion verifica si una sesión es válida
func (s *SesionService) VerificarSesion(ctx context.Context, token string) (*models.Sesion, error) {
	// TODO: Verificar sesión en Supabase
	return nil, fmt.Errorf("no implementado")
}

// generateToken genera un token único para la sesión
func generateToken() string {
	// TODO: Implementar generación de token seguro
	return "token-temporal"
}
