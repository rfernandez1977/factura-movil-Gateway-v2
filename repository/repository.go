package repository

import "context"

// Repository define la interfaz para operaciones con la base de datos
type Repository interface {
	// Operaciones CRUD básicas
	FindOne(collection string, filter map[string]interface{}, result interface{}) error
	Find(collection string, filter map[string]interface{}, results interface{}, options ...QueryOption) error
	InsertOne(collection string, document interface{}) (string, error)
	UpdateOne(collection string, filter map[string]interface{}, update interface{}) error
	DeleteOne(collection string, filter map[string]interface{}) error

	// Operaciones avanzadas
	Count(collection string, filter map[string]interface{}) (int64, error)
	Aggregate(collection string, pipeline interface{}, results interface{}) error

	// Operaciones de transacción
	WithTransaction(ctx context.Context, fn func(txRepo Repository) error) error
}

// QueryOption representa opciones para la consulta
type QueryOption struct {
	Sort  map[string]int
	Skip  int64
	Limit int64
}

// NewQueryOption crea una nueva opción de consulta
func NewQueryOption() *QueryOption {
	return &QueryOption{
		Sort: make(map[string]int),
	}
}

// WithSort agrega ordenamiento a la consulta
func (q *QueryOption) WithSort(field string, direction int) *QueryOption {
	q.Sort[field] = direction
	return q
}

// WithSkip agrega salto a la consulta
func (q *QueryOption) WithSkip(skip int64) *QueryOption {
	q.Skip = skip
	return q
}

// WithLimit agrega límite a la consulta
func (q *QueryOption) WithLimit(limit int64) *QueryOption {
	q.Limit = limit
	return q
}
