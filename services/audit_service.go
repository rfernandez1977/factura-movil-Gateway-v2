package services

import (
	"context"
	"time"

	"FMgo/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// AuditService implementa la interfaz domain.AuditService
type AuditService struct {
	db *mongo.Database
}

// NewAuditService crea una nueva instancia del servicio de auditoría
func NewAuditService(db *mongo.Database) domain.AuditService {
	return &AuditService{
		db: db,
	}
}

// RegistrarOperacion registra una operación en el sistema
func (s *AuditService) RegistrarOperacion(
	ctx context.Context,
	operacion string,
	entidad string,
	entidadID primitive.ObjectID,
	usuario string,
) error {
	collection := s.db.Collection("auditoria_operaciones")
	registro := bson.M{
		"operacion":       operacion,
		"entidad":         entidad,
		"entidad_id":      entidadID,
		"usuario":         usuario,
		"fecha_operacion": time.Now(),
	}

	_, err := collection.InsertOne(ctx, registro)
	return err
}

// ObtenerRegistros obtiene los registros de auditoría
func (s *AuditService) ObtenerRegistros(ctx context.Context, filtro bson.M) ([]bson.M, error) {
	collection := s.db.Collection("auditoria_operaciones")
	cursor, err := collection.Find(ctx, filtro)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var registros []bson.M
	if err = cursor.All(ctx, &registros); err != nil {
		return nil, err
	}
	return registros, nil
}

// GenerarReporte genera un reporte de auditoría
func (s *AuditService) GenerarReporte(ctx context.Context, fechaInicio, fechaFin time.Time) ([]bson.M, error) {
	collection := s.db.Collection("auditoria_operaciones")
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"fecha_operacion": bson.M{
					"$gte": fechaInicio,
					"$lte": fechaFin,
				},
			},
		},
		{
			"$group": bson.M{
				"_id": bson.M{
					"operacion": "$operacion",
					"entidad":   "$entidad",
				},
				"total": bson.M{"$sum": 1},
				"usuarios": bson.M{
					"$addToSet": "$usuario",
				},
			},
		},
		{
			"$project": bson.M{
				"_id":       0,
				"operacion": "$_id.operacion",
				"entidad":   "$_id.entidad",
				"total":     1,
				"usuarios":  1,
			},
		},
		{
			"$sort": bson.M{
				"total": -1,
			},
		},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var reporte []bson.M
	if err = cursor.All(ctx, &reporte); err != nil {
		return nil, err
	}
	return reporte, nil
}
