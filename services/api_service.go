package services

import (
	"context"
	"fmt"
	"time"

	"github.com/fmgo/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// APIService maneja las operaciones relacionadas con la API
type APIService struct {
	db *mongo.Database
}

// NewAPIService crea una nueva instancia de APIService
func NewAPIService(db *mongo.Database) *APIService {
	return &APIService{
		db: db,
	}
}

// API representa una configuración de API
type API struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	Version     string             `bson:"version" json:"version"`
	Enabled     bool               `bson:"enabled" json:"enabled"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

// Version representa una versión de la API
type Version struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	APIID       primitive.ObjectID `bson:"api_id" json:"api_id"`
	Version     string             `bson:"version" json:"version"`
	Description string             `bson:"description" json:"description"`
	Enabled     bool               `bson:"enabled" json:"enabled"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

// Registro representa un registro de uso de la API
type Registro struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	APIID        primitive.ObjectID `bson:"api_id" json:"api_id"`
	VersionID    primitive.ObjectID `bson:"version_id" json:"version_id"`
	Endpoint     string             `bson:"endpoint" json:"endpoint"`
	Method       string             `bson:"method" json:"method"`
	Status       int                `bson:"status" json:"status"`
	ResponseTime int64              `bson:"response_time" json:"response_time"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
}

// CreateAPI crea una nueva API
func (s *APIService) CreateAPI(ctx context.Context, api *API) error {
	if api.ID.IsZero() {
		api.ID = primitive.NewObjectID()
	}
	api.CreatedAt = time.Now()
	api.UpdatedAt = time.Now()

	_, err := s.db.Collection("apis").InsertOne(ctx, api)
	return err
}

// CreateVersion crea una nueva versión de API
func (s *APIService) CreateVersion(ctx context.Context, version *Version) error {
	if version.ID.IsZero() {
		version.ID = primitive.NewObjectID()
	}
	version.CreatedAt = time.Now()
	version.UpdatedAt = time.Now()

	_, err := s.db.Collection("api_versions").InsertOne(ctx, version)
	return err
}

// CreateRegistro crea un nuevo registro de uso de API
func (s *APIService) CreateRegistro(ctx context.Context, registro *Registro) error {
	if registro.ID.IsZero() {
		registro.ID = primitive.NewObjectID()
	}
	registro.CreatedAt = time.Now()

	_, err := s.db.Collection("api_registros").InsertOne(ctx, registro)
	return err
}

// GetAPI obtiene una API por su ID
func (s *APIService) GetAPI(ctx context.Context, id primitive.ObjectID) (*API, error) {
	var api API
	err := s.db.Collection("apis").FindOne(ctx, bson.M{"_id": id}).Decode(&api)
	if err != nil {
		return nil, err
	}
	return &api, nil
}

// GetVersion obtiene una versión de API por su ID
func (s *APIService) GetVersion(ctx context.Context, id primitive.ObjectID) (*Version, error) {
	var version Version
	err := s.db.Collection("api_versions").FindOne(ctx, bson.M{"_id": id}).Decode(&version)
	if err != nil {
		return nil, err
	}
	return &version, nil
}

// GetRegistros obtiene los registros de uso de una API
func (s *APIService) GetRegistros(ctx context.Context, apiID primitive.ObjectID, limit int64) ([]Registro, error) {
	cursor, err := s.db.Collection("api_registros").
		Find(ctx, bson.M{"api_id": apiID}, options.Find().SetLimit(limit).SetSort(bson.M{"created_at": -1}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var registros []Registro
	if err = cursor.All(ctx, &registros); err != nil {
		return nil, err
	}
	return registros, nil
}

// RegistrarAPI registra una API
func (s *APIService) RegistrarAPI(ctx context.Context, api *models.API) error {
	collection := s.db.Collection("apis")

	if api.ID.IsZero() {
		api.ID = primitive.NewObjectID()
	}

	_, err := collection.InsertOne(ctx, api)
	return err
}

// ObtenerAPIs obtiene las APIs
func (s *APIService) ObtenerAPIs(ctx context.Context, filtro bson.M) ([]models.API, error) {
	collection := s.db.Collection("apis")

	cursor, err := collection.Find(ctx, filtro)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var apis []models.API
	if err = cursor.All(ctx, &apis); err != nil {
		return nil, err
	}

	return apis, nil
}

// RegistrarVersionAPI registra una versión de API
func (s *APIService) RegistrarVersionAPI(ctx context.Context, version *models.VersionAPI) error {
	collection := s.db.Collection("versiones_api")

	if version.ID.IsZero() {
		version.ID = primitive.NewObjectID()
	}

	_, err := collection.InsertOne(ctx, version)
	return err
}

// ObtenerVersionesAPI obtiene las versiones de una API
func (s *APIService) ObtenerVersionesAPI(ctx context.Context, apiID primitive.ObjectID) ([]models.VersionAPI, error) {
	collection := s.db.Collection("versiones_api")

	cursor, err := collection.Find(ctx, bson.M{"api_id": apiID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var versiones []models.VersionAPI
	if err = cursor.All(ctx, &versiones); err != nil {
		return nil, err
	}

	return versiones, nil
}

// RegistrarRegistroAPI registra un registro de API
func (s *APIService) RegistrarRegistroAPI(ctx context.Context, registro *models.RegistroAPI) error {
	collection := s.db.Collection("registros_api")

	if registro.ID.IsZero() {
		registro.ID = primitive.NewObjectID()
	}

	_, err := collection.InsertOne(ctx, registro)
	return err
}

// ObtenerRegistrosAPI obtiene los registros de API
func (s *APIService) ObtenerRegistrosAPI(ctx context.Context, filtro bson.M) ([]models.RegistroAPI, error) {
	collection := s.db.Collection("registros_api")

	cursor, err := collection.Find(ctx, filtro)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var registros []models.RegistroAPI
	if err = cursor.All(ctx, &registros); err != nil {
		return nil, err
	}

	return registros, nil
}

// ReporteAPI representa un reporte de uso de la API
type ReporteAPI struct {
	ID                      primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	PeriodoInicio           time.Time          `bson:"periodo_inicio" json:"periodo_inicio"`
	PeriodoFin              time.Time          `bson:"periodo_fin" json:"periodo_fin"`
	TotalLlamadas           int                `bson:"total_llamadas" json:"total_llamadas"`
	LlamadasExitosas        int                `bson:"llamadas_exitosas" json:"llamadas_exitosas"`
	LlamadasFallidas        int                `bson:"llamadas_fallidas" json:"llamadas_fallidas"`
	TiempoPromedioRespuesta time.Duration      `bson:"tiempo_promedio_respuesta" json:"tiempo_promedio_respuesta"`
	ErroresPorTipo          map[string]int     `bson:"errores_por_tipo" json:"errores_por_tipo"`
	FechaGeneracion         time.Time          `bson:"fecha_generacion" json:"fecha_generacion"`
}

// GenerarReporteAPI genera un reporte de uso de la API
func (s *APIService) GenerarReporteAPI(ctx context.Context, inicio, fin time.Time) (*ReporteAPI, error) {
	reporte := &ReporteAPI{
		ID:              primitive.NewObjectID(),
		PeriodoInicio:   inicio,
		PeriodoFin:      fin,
		FechaGeneracion: time.Now(),
	}

	// Obtener total de llamadas
	totalLlamadas, err := s.db.Collection("registros_api").CountDocuments(ctx, bson.M{
		"created_at": bson.M{
			"$gte": inicio,
			"$lte": fin,
		},
	})
	if err != nil {
		return nil, err
	}
	reporte.TotalLlamadas = int(totalLlamadas)

	// Obtener llamadas exitosas
	llamadasExitosas, err := s.db.Collection("registros_api").CountDocuments(ctx, bson.M{
		"created_at": bson.M{
			"$gte": inicio,
			"$lte": fin,
		},
		"status": bson.M{"$lt": 400},
	})
	if err != nil {
		return nil, err
	}
	reporte.LlamadasExitosas = int(llamadasExitosas)
	reporte.LlamadasFallidas = reporte.TotalLlamadas - reporte.LlamadasExitosas

	// Calcular tiempo promedio de respuesta
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"created_at": bson.M{
					"$gte": inicio,
					"$lte": fin,
				},
			},
		},
		{
			"$group": bson.M{
				"_id": nil,
				"avg_response_time": bson.M{
					"$avg": "$response_time",
				},
			},
		},
	}

	cursor, err := s.db.Collection("registros_api").Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var result []bson.M
	if err = cursor.All(ctx, &result); err != nil {
		return nil, err
	}

	if len(result) > 0 {
		reporte.TiempoPromedioRespuesta = time.Duration(result[0]["avg_response_time"].(float64))
	}

	// Obtener errores por tipo
	pipeline = []bson.M{
		{
			"$match": bson.M{
				"created_at": bson.M{
					"$gte": inicio,
					"$lte": fin,
				},
				"status": bson.M{"$gte": 400},
			},
		},
		{
			"$group": bson.M{
				"_id": "$status",
				"count": bson.M{
					"$sum": 1,
				},
			},
		},
	}

	cursor, err = s.db.Collection("registros_api").Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &result); err != nil {
		return nil, err
	}

	reporte.ErroresPorTipo = make(map[string]int)
	for _, r := range result {
		status := fmt.Sprintf("%v", r["_id"])
		count := r["count"].(int32)
		reporte.ErroresPorTipo[status] = int(count)
	}

	return reporte, nil
}
