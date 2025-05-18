package repository

import (
	"context"
	"errors"

	"FMgo/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// DocumentRepository define las operaciones para interactuar con la base de datos de documentos.
type DocumentRepository interface {
	// Documentos
	SaveDocumentoTributario(doc domain.DocumentoTributario) error
	GetDocumentoTributario(tipo string, folio int64) (*domain.DocumentoTributario, error)
	GetDocumentoTributarioByID(id primitive.ObjectID) (*domain.DocumentoTributario, error)
	UpdateDocumentoTributario(doc domain.DocumentoTributario) error
	GetDocumentosPorEstado(estado string) ([]domain.DocumentoTributario, error)

	// Estados
	SaveEstadoDocumento(estado domain.EstadoDocumento) error
	GetEstadoDocumento(docID primitive.ObjectID) (*domain.EstadoDocumento, error)
	UpdateEstadoDocumento(estado domain.EstadoDocumento) error

	// Referencias
	SaveReferenciaDocumento(ref domain.ReferenciaDocumento) error
	GetReferenciasPorDocumento(tipoOrigen string, folioOrigen int64) ([]domain.ReferenciaDocumento, error)
}

// DocumentRepositoryImpl implementa la interfaz DocumentRepository
type DocumentRepositoryImpl struct {
	db *mongo.Database
}

// NewDocumentRepository crea una nueva instancia de DocumentRepository
func NewDocumentRepository(db *mongo.Database) domain.DocumentRepository {
	return &DocumentRepositoryImpl{
		db: db,
	}
}

// SaveDocumentoTributario guarda un documento tributario
func (r *DocumentRepositoryImpl) SaveDocumentoTributario(doc domain.DocumentoTributario) error {
	collection := r.db.Collection("documentos_tributarios")
	_, err := collection.InsertOne(context.Background(), doc)
	return err
}

// GetDocumentoTributario obtiene un documento tributario
func (r *DocumentRepositoryImpl) GetDocumentoTributario(tipo string, folio int64) (*domain.DocumentoTributario, error) {
	collection := r.db.Collection("documentos_tributarios")
	var doc domain.DocumentoTributario
	err := collection.FindOne(context.Background(), bson.M{
		"tipo_documento": tipo,
		"folio":          folio,
	}).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &doc, nil
}

// GetDocumentoTributarioByID obtiene un documento tributario por su ID
func (r *DocumentRepositoryImpl) GetDocumentoTributarioByID(id primitive.ObjectID) (*domain.DocumentoTributario, error) {
	collection := r.db.Collection("documentos_tributarios")
	var doc domain.DocumentoTributario
	err := collection.FindOne(context.Background(), bson.M{
		"_id": id,
	}).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &doc, nil
}

// UpdateDocumentoTributario actualiza un documento tributario
func (r *DocumentRepositoryImpl) UpdateDocumentoTributario(doc domain.DocumentoTributario) error {
	collection := r.db.Collection("documentos_tributarios")
	_, err := collection.UpdateOne(
		context.Background(),
		bson.M{"_id": doc.ID},
		bson.M{"$set": doc},
	)
	return err
}

// GetDocumentosPorEstado obtiene documentos por estado
func (r *DocumentRepositoryImpl) GetDocumentosPorEstado(estado string) ([]domain.DocumentoTributario, error) {
	collection := r.db.Collection("documentos_tributarios")
	cursor, err := collection.Find(context.Background(), bson.M{"estado": estado})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var docs []domain.DocumentoTributario
	if err = cursor.All(context.Background(), &docs); err != nil {
		return nil, err
	}
	return docs, nil
}

// SaveEstadoDocumento guarda un estado de documento
func (r *DocumentRepositoryImpl) SaveEstadoDocumento(estado domain.EstadoDocumento) error {
	collection := r.db.Collection("estados_documentos")
	_, err := collection.InsertOne(context.Background(), estado)
	return err
}

// GetEstadoDocumento obtiene un estado de documento
func (r *DocumentRepositoryImpl) GetEstadoDocumento(docID primitive.ObjectID) (*domain.EstadoDocumento, error) {
	collection := r.db.Collection("estados_documentos")
	var estado domain.EstadoDocumento
	err := collection.FindOne(context.Background(), bson.M{
		"documento_id": docID,
	}).Decode(&estado)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &estado, nil
}

// UpdateEstadoDocumento actualiza un estado de documento
func (r *DocumentRepositoryImpl) UpdateEstadoDocumento(estado domain.EstadoDocumento) error {
	collection := r.db.Collection("estados_documentos")
	_, err := collection.UpdateOne(
		context.Background(),
		bson.M{"_id": estado.ID},
		bson.M{"$set": estado},
	)
	return err
}

// SaveReferenciaDocumento guarda una referencia de documento
func (r *DocumentRepositoryImpl) SaveReferenciaDocumento(ref domain.ReferenciaDocumento) error {
	collection := r.db.Collection("referencias_documentos")
	_, err := collection.InsertOne(context.Background(), ref)
	return err
}

// GetReferenciasPorDocumento obtiene las referencias de un documento
func (r *DocumentRepositoryImpl) GetReferenciasPorDocumento(tipoOrigen string, folioOrigen int64) ([]domain.ReferenciaDocumento, error) {
	collection := r.db.Collection("referencias_documentos")
	cursor, err := collection.Find(context.Background(), bson.M{
		"tipo_origen":  tipoOrigen,
		"folio_origen": folioOrigen,
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var refs []domain.ReferenciaDocumento
	if err = cursor.All(context.Background(), &refs); err != nil {
		return nil, err
	}
	return refs, nil
}
