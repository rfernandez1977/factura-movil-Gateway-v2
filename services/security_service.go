package services

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"time"

	"github.com/cursor/FMgo/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// SecurityService maneja las operaciones de seguridad
type SecurityService struct {
	db *mongo.Database
}

// NewSecurityService crea una nueva instancia del servicio de seguridad
func NewSecurityService(db *mongo.Database) *SecurityService {
	return &SecurityService{
		db: db,
	}
}

// GenerarCertificado genera un nuevo certificado digital
func (s *SecurityService) GenerarCertificado(ctx context.Context, config *models.ConfiguracionCertificado) (*models.CertificadoDigital, error) {
	// Generar clave privada
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	// Crear certificado
	certificado := &models.CertificadoDigital{
		ID:            primitive.NewObjectID(),
		Nombre:        config.Nombre,
		Organizacion:  config.Organizacion,
		ValidoDesde:   time.Now(),
		ValidoHasta:   time.Now().AddDate(1, 0, 0), // 1 año de validez
		Estado:        models.EstadoCertificadoActivo,
		FechaCreacion: time.Now(),
	}

	// Guardar certificado
	collection := s.db.Collection("certificados")
	_, err = collection.InsertOne(ctx, certificado)
	if err != nil {
		return nil, err
	}

	return certificado, nil
}

// ValidarCertificado valida un certificado digital
func (s *SecurityService) ValidarCertificado(ctx context.Context, certificadoID primitive.ObjectID) error {
	collection := s.db.Collection("certificados")

	var certificado models.CertificadoDigital
	err := collection.FindOne(ctx, bson.M{"_id": certificadoID}).Decode(&certificado)
	if err != nil {
		return err
	}

	// Verificar validez temporal
	if time.Now().Before(certificado.ValidoDesde) || time.Now().After(certificado.ValidoHasta) {
		return errors.New("certificado no válido en este momento")
	}

	// Verificar estado
	if certificado.Estado != models.EstadoCertificadoActivo {
		return errors.New("certificado no activo")
	}

	return nil
}

// RegistrarAcceso registra un acceso al sistema
func (s *SecurityService) RegistrarAcceso(ctx context.Context, acceso *models.RegistroAcceso) error {
	collection := s.db.Collection("registros_acceso")

	if acceso.ID.IsZero() {
		acceso.ID = primitive.NewObjectID()
	}

	_, err := collection.InsertOne(ctx, acceso)
	return err
}

// ObtenerRegistrosAcceso obtiene los registros de acceso
func (s *SecurityService) ObtenerRegistrosAcceso(ctx context.Context, filtro bson.M) ([]models.RegistroAcceso, error) {
	collection := s.db.Collection("registros_acceso")

	cursor, err := collection.Find(ctx, filtro)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var registros []models.RegistroAcceso
	if err = cursor.All(ctx, &registros); err != nil {
		return nil, err
	}

	return registros, nil
}

// GenerarReporteSeguridad genera un reporte de seguridad
func (s *SecurityService) GenerarReporteSeguridad(ctx context.Context, inicio, fin time.Time) (*models.ReporteSeguridad, error) {
	// Obtener registros de acceso del período
	registros, err := s.ObtenerRegistrosAcceso(ctx, bson.M{
		"fecha": bson.M{
			"$gte": inicio,
			"$lte": fin,
		},
	})
	if err != nil {
		return nil, err
	}

	// Calcular estadísticas
	reporte := &models.ReporteSeguridad{
		ID:              primitive.NewObjectID(),
		PeriodoInicio:   inicio,
		PeriodoFin:      fin,
		TotalAccesos:    len(registros),
		AccesosExitosos: 0,
		AccesosFallidos: 0,
		IntentosPorIP:   make(map[string]int),
		FechaGeneracion: time.Now(),
	}

	for _, registro := range registros {
		if registro.Exitoso {
			reporte.AccesosExitosos++
		} else {
			reporte.AccesosFallidos++
		}
		reporte.IntentosPorIP[registro.IP]++
	}

	// Guardar el reporte
	collection := s.db.Collection("reportes_seguridad")
	_, err = collection.InsertOne(ctx, reporte)
	if err != nil {
		return nil, err
	}

	return reporte, nil
}
