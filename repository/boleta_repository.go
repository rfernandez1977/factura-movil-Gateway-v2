package repository

import (
	"context"
	"errors"
	"time"

	"FMgo/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// BoletaRepository define las operaciones de base de datos para boletas
type BoletaRepository interface {
	Create(ctx context.Context, doc *models.DocumentoTributario) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*models.DocumentoTributario, error)
	GetByTrackID(ctx context.Context, trackID string) (*models.DocumentoTributario, error)
	GetByFolio(ctx context.Context, folio int) (*models.DocumentoTributario, error)
	UpdateEstado(ctx context.Context, id primitive.ObjectID, estado models.EstadoDocumento) error
	UpdateTrackID(ctx context.Context, id primitive.ObjectID, trackID string) error
}

// BoletaRepositoryImpl implementa BoletaRepository
type BoletaRepositoryImpl struct {
	collection *mongo.Collection
}

// NewBoletaRepository crea una nueva instancia de BoletaRepository
func NewBoletaRepository(collection *mongo.Collection) BoletaRepository {
	return &BoletaRepositoryImpl{
		collection: collection,
	}
}

// Create guarda una nueva boleta
func (r *BoletaRepositoryImpl) Create(ctx context.Context, doc *models.DocumentoTributario) error {
	_, err := r.collection.InsertOne(ctx, doc)
	return err
}

// GetByID obtiene una boleta por su ID
func (r *BoletaRepositoryImpl) GetByID(ctx context.Context, id primitive.ObjectID) (*models.DocumentoTributario, error) {
	var doc models.DocumentoTributario
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("boleta no encontrada")
		}
		return nil, err
	}
	return &doc, nil
}

// GetByTrackID obtiene una boleta por su track ID
func (r *BoletaRepositoryImpl) GetByTrackID(ctx context.Context, trackID string) (*models.DocumentoTributario, error) {
	var doc models.DocumentoTributario
	err := r.collection.FindOne(ctx, bson.M{"track_id": trackID}).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("boleta no encontrada")
		}
		return nil, err
	}
	return &doc, nil
}

// GetByFolio obtiene una boleta por su folio
func (r *BoletaRepositoryImpl) GetByFolio(ctx context.Context, folio int) (*models.DocumentoTributario, error) {
	var doc models.DocumentoTributario
	err := r.collection.FindOne(ctx, bson.M{"folio": folio}).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("boleta no encontrada")
		}
		return nil, err
	}
	return &doc, nil
}

// UpdateEstado actualiza el estado de una boleta
func (r *BoletaRepositoryImpl) UpdateEstado(ctx context.Context, id primitive.ObjectID, estado models.EstadoDocumento) error {
	update := bson.M{
		"$set": bson.M{
			"estado":     estado,
			"updated_at": time.Now(),
		},
	}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

// UpdateTrackID actualiza el track ID de una boleta
func (r *BoletaRepositoryImpl) UpdateTrackID(ctx context.Context, id primitive.ObjectID, trackID string) error {
	update := bson.M{
		"$set": bson.M{
			"track_id":   trackID,
			"updated_at": time.Now(),
		},
	}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

// Errores del repositorio
var (
	ErrIDRequired     = errors.New("ID es requerido")
	ErrBoletaExists   = errors.New("la boleta ya existe")
	ErrBoletaNotFound = errors.New("boleta no encontrada")
)
