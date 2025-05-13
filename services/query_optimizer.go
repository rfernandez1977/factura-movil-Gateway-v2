package services

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// QueryOptimizer maneja la optimización de consultas
type QueryOptimizer struct {
	db *mongo.Database
}

// NewQueryOptimizer crea una nueva instancia del optimizador de consultas
func NewQueryOptimizer(db *mongo.Database) *QueryOptimizer {
	return &QueryOptimizer{db: db}
}

// OptimizeFind optimiza una consulta de búsqueda
func (o *QueryOptimizer) OptimizeFind(ctx context.Context, collection string, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	// Aplicar índices compuestos
	opts = append(opts, options.Find().SetHint(bson.D{{Key: "_id", Value: 1}}))

	// Limitar el número de documentos
	opts = append(opts, options.Find().SetLimit(1000))

	// Proyectar solo los campos necesarios
	opts = append(opts, options.Find().SetProjection(bson.D{
		{Key: "_id", Value: 1},
		{Key: "created_at", Value: 1},
		{Key: "updated_at", Value: 1},
	}))

	return o.db.Collection(collection).Find(ctx, filter, opts...)
}

// OptimizeAggregate optimiza una consulta de agregación
func (o *QueryOptimizer) OptimizeAggregate(ctx context.Context, collection string, pipeline []bson.D, opts ...*options.AggregateOptions) (*mongo.Cursor, error) {
	// Aplicar índices compuestos
	opts = append(opts, options.Aggregate().SetHint(bson.D{{Key: "_id", Value: 1}}))

	// Limitar el número de documentos
	opts = append(opts, options.Aggregate().SetMaxTime(10*time.Second))

	// Permitir uso de disco para operaciones grandes
	opts = append(opts, options.Aggregate().SetAllowDiskUse(true))

	return o.db.Collection(collection).Aggregate(ctx, pipeline, opts...)
}

// CreateIndexes crea índices optimizados para una colección
func (o *QueryOptimizer) CreateIndexes(ctx context.Context, collection string, indexes []mongo.IndexModel) error {
	_, err := o.db.Collection(collection).Indexes().CreateMany(ctx, indexes)
	return err
}

// GetIndexes obtiene los índices de una colección
func (o *QueryOptimizer) GetIndexes(ctx context.Context, collection string) ([]bson.M, error) {
	cursor, err := o.db.Collection(collection).Indexes().List(ctx)
	if err != nil {
		return nil, err
	}

	var indexes []bson.M
	if err = cursor.All(ctx, &indexes); err != nil {
		return nil, err
	}

	return indexes, nil
}

// DropIndexes elimina los índices de una colección
func (o *QueryOptimizer) DropIndexes(ctx context.Context, collection string) error {
	_, err := o.db.Collection(collection).Indexes().DropAll(ctx)
	return err
}

// Explain ejecuta una consulta con explain
func (o *QueryOptimizer) Explain(ctx context.Context, collection string, filter interface{}, opts ...*options.FindOptions) (bson.M, error) {
	opts = append(opts, options.Find().SetExplain(true))
	cursor, err := o.db.Collection(collection).Find(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}

	var result bson.M
	if err = cursor.Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

// QueryStats representa estadísticas de una consulta
type QueryStats struct {
	ExecutionTime time.Duration
	DocsExamined  int64
	DocsReturned  int64
	IndexUsed     string
}

// GetQueryStats obtiene estadísticas de una consulta
func (o *QueryOptimizer) GetQueryStats(ctx context.Context, collection string, filter interface{}, opts ...*options.FindOptions) (*QueryStats, error) {
	start := time.Now()

	cursor, err := o.db.Collection(collection).Find(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}

	var docs []bson.M
	if err = cursor.All(ctx, &docs); err != nil {
		return nil, err
	}

	explain, err := o.Explain(ctx, collection, filter, opts...)
	if err != nil {
		return nil, err
	}

	stats := &QueryStats{
		ExecutionTime: time.Since(start),
		DocsReturned:  int64(len(docs)),
	}

	if executionStats, ok := explain["executionStats"].(bson.M); ok {
		if totalDocsExamined, ok := executionStats["totalDocsExamined"].(int64); ok {
			stats.DocsExamined = totalDocsExamined
		}
		if winningPlan, ok := executionStats["winningPlan"].(bson.M); ok {
			if inputStage, ok := winningPlan["inputStage"].(bson.M); ok {
				if indexName, ok := inputStage["indexName"].(string); ok {
					stats.IndexUsed = indexName
				}
			}
		}
	}

	return stats, nil
}
