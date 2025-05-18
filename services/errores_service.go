package services

import (
	"context"
	"log"
	"math"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"FMgo/models"
)

// ErroresService maneja el sistema de errores y recuperación
type ErroresService struct {
	db *mongo.Database
}

// NewErroresService crea una nueva instancia del servicio de errores
func NewErroresService(db *mongo.Database) *ErroresService {
	return &ErroresService{
		db: db,
	}
}

// RegistrarError registra un nuevo error en el sistema
func (s *ErroresService) RegistrarError(
	ctx context.Context,
	tipo string,
	severidad string,
	codigo, mensaje, descripcion string,
	stacktrace string,
	contexto map[string]interface{},
	entidad, entidadID, usuarioID string,
) (*models.ErrorLog, error) {
	errorLog := &models.ErrorLog{
		ID:        models.GenerateErrorID(),
		Tipo:      tipo,
		Mensaje:   mensaje,
		Detalles:  descripcion,
		Usuario:   usuarioID,
		Timestamp: time.Now(),
	}

	_, err := s.db.Collection("errores").InsertOne(ctx, errorLog)
	if err != nil {
		return nil, err
	}

	// Registrar log del error
	if err := s.registrarLogError(ctx, errorLog.ID, "ERROR", mensaje, contexto); err != nil {
		log.Printf("Error al registrar log: %v", err)
	}

	return errorLog, nil
}

// RegistrarLogError registra un log de error
func (s *ErroresService) registrarLogError(
	ctx context.Context,
	errorID string,
	nivel, mensaje string,
	contexto map[string]interface{},
) error {
	logError := &models.LogError{
		ID:        models.GenerateErrorID(),
		ErrorID:   errorID,
		Nivel:     nivel,
		Mensaje:   mensaje,
		Contexto:  contexto,
		FechaLog:  time.Now(),
		IP:        getIPFromContext(ctx),
		UserAgent: getUserAgentFromContext(ctx),
		UsuarioID: getUsuarioIDFromContext(ctx),
	}

	_, err := s.db.Collection("logs_errores").InsertOne(ctx, logError)
	return err
}

// iniciarProcesoRecuperacion inicia el proceso de recuperación de un error
func (s *ErroresService) iniciarProcesoRecuperacion(errorID string) {
	ctx := context.Background()
	config, err := s.obtenerConfiguracionRecuperacion(ctx)
	if err != nil {
		log.Printf("Error al obtener configuración de recuperación: %v", err)
		return
	}

	for intento := 1; intento <= config.MaxIntentos; intento++ {
		// Calcular intervalo de espera con backoff exponencial
		intervalo := int(math.Min(
			float64(config.IntervaloBase)*math.Pow(config.FactorExponencial, float64(intento-1)),
			float64(config.MaxIntervalo),
		))

		// Esperar antes del siguiente intento
		time.Sleep(time.Duration(intervalo) * time.Second)

		// Realizar intento de recuperación
		exitoso, mensaje, detalles := s.intentarRecuperacion(ctx, errorID, intento)

		// Registrar intento
		intentoRecuperacion := &models.IntentoRecuperacion{
			ID:              models.GenerateErrorID(),
			ErrorID:         errorID,
			NumeroIntento:   intento,
			Exitoso:         exitoso,
			Mensaje:         mensaje,
			Detalles:        detalles,
			FechaIntento:    time.Now(),
			TiempoEjecucion: time.Now().UnixMilli(),
		}

		_, err := s.db.Collection("intentos_recuperacion").InsertOne(ctx, intentoRecuperacion)
		if err != nil {
			log.Printf("Error al registrar intento de recuperación: %v", err)
		}

		if exitoso {
			// Actualizar estado del error
			_, err := s.db.Collection("errores").UpdateOne(
				ctx,
				bson.M{"_id": errorID},
				bson.M{"$set": bson.M{
					"estado":           "RESUELTO",
					"fecha_resolucion": time.Now(),
				}},
			)
			if err != nil {
				log.Printf("Error al actualizar estado del error: %v", err)
			}
			return
		}

		// Notificar al administrador si se alcanza el umbral
		if intento >= config.UmbralNotificacion && config.NotificarAdmin {
			s.notificarAdministrador(ctx, errorID, intento)
		}
	}

	// Si se agotaron los intentos, marcar como no recuperable
	_, err = s.db.Collection("errores").UpdateOne(
		ctx,
		bson.M{"_id": errorID},
		bson.M{"$set": bson.M{
			"estado":           "NO_RECUPERABLE",
			"fecha_resolucion": time.Now(),
		}},
	)
	if err != nil {
		log.Printf("Error al marcar error como no recuperable: %v", err)
	}
}

// intentarRecuperacion realiza un intento de recuperación
func (s *ErroresService) intentarRecuperacion(
	ctx context.Context,
	errorID string,
	intento int,
) (bool, string, string) {
	// Obtener detalles del error
	var errorDetalle models.ErrorDetalle
	err := s.db.Collection("errores").FindOne(ctx, bson.M{"_id": errorID}).Decode(&errorDetalle)
	if err != nil {
		return false, "Error al obtener detalles del error", err.Error()
	}

	// Implementar lógica específica de recuperación según el tipo de error
	switch errorDetalle.Tipo {
	case models.ErrorBaseDatos:
		return s.recuperarErrorBaseDatos(ctx, errorDetalle)
	case models.ErrorIntegracion:
		return s.recuperarErrorIntegracion(ctx, errorDetalle)
	case models.ErrorRed:
		return s.recuperarErrorRed(ctx, errorDetalle)
	default:
		return false, "Tipo de error no soportado para recuperación", ""
	}
}

// recuperarErrorBaseDatos implementa la recuperación de errores de base de datos
func (s *ErroresService) recuperarErrorBaseDatos(
	ctx context.Context,
	errorDetalle models.ErrorDetalle,
) (bool, string, string) {
	// Implementar lógica específica para recuperación de errores de base de datos
	// Por ejemplo: reconexión, reintento de transacción, etc.
	return false, "Recuperación de base de datos no implementada", ""
}

// recuperarErrorIntegracion implementa la recuperación de errores de integración
func (s *ErroresService) recuperarErrorIntegracion(
	ctx context.Context,
	errorDetalle models.ErrorDetalle,
) (bool, string, string) {
	// Implementar lógica específica para recuperación de errores de integración
	// Por ejemplo: reintento de llamada a API, uso de caché, etc.
	return false, "Recuperación de integración no implementada", ""
}

// recuperarErrorRed implementa la recuperación de errores de red
func (s *ErroresService) recuperarErrorRed(
	ctx context.Context,
	errorDetalle models.ErrorDetalle,
) (bool, string, string) {
	// Implementar lógica específica para recuperación de errores de red
	// Por ejemplo: reintento de conexión, cambio de servidor, etc.
	return false, "Recuperación de red no implementada", ""
}

// obtenerConfiguracionRecuperacion obtiene la configuración de recuperación
func (s *ErroresService) obtenerConfiguracionRecuperacion(
	ctx context.Context,
) (*models.ConfiguracionRecuperacion, error) {
	var config models.ConfiguracionRecuperacion
	err := s.db.Collection("configuracion_recuperacion").FindOne(ctx, bson.M{}).Decode(&config)
	if err != nil {
		// Si no existe configuración, usar valores por defecto
		config = models.ConfiguracionRecuperacion{
			ID:                 models.GenerateErrorID(),
			MaxIntentos:        3,
			IntervaloBase:      5,
			FactorExponencial:  2,
			MaxIntervalo:       60,
			NotificarAdmin:     true,
			UmbralNotificacion: 2,
			ActivarLogging:     true,
			NivelLogging:       "ERROR",
		}
		_, err = s.db.Collection("configuracion_recuperacion").InsertOne(ctx, config)
		if err != nil {
			return nil, err
		}
	}
	return &config, nil
}

// notificarAdministrador envía una notificación al administrador
func (s *ErroresService) notificarAdministrador(
	ctx context.Context,
	errorID string,
	intento int,
) {
	// Implementar lógica de notificación
	// Por ejemplo: email, SMS, notificación push, etc.
	log.Printf("Notificando al administrador sobre el error %s (intento %d)", errorID, intento)
}

// GenerarReporteErrores genera un reporte de errores
func (s *ErroresService) GenerarReporteErrores(
	ctx context.Context,
	fechaInicio, fechaFin time.Time,
) (*models.ReporteErrores, error) {
	// Contar total de errores
	totalErrores, err := s.db.Collection("errores").CountDocuments(ctx, bson.M{
		"fecha_error": bson.M{
			"$gte": fechaInicio,
			"$lte": fechaFin,
		},
	})
	if err != nil {
		return nil, err
	}

	// Contar errores por tipo
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"fecha_error": bson.M{
					"$gte": fechaInicio,
					"$lte": fechaFin,
				},
			},
		},
		{
			"$group": bson.M{
				"_id":   "$tipo",
				"count": bson.M{"$sum": 1},
			},
		},
	}

	cursor, err := s.db.Collection("errores").Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	erroresPorTipo := make(map[models.TipoError]int)
	for cursor.Next(ctx) {
		var result struct {
			ID    models.TipoError `bson:"_id"`
			Count int              `bson:"count"`
		}
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
		erroresPorTipo[result.ID] = result.Count
	}

	// Contar errores por severidad
	pipeline = []bson.M{
		{
			"$match": bson.M{
				"fecha_error": bson.M{
					"$gte": fechaInicio,
					"$lte": fechaFin,
				},
			},
		},
		{
			"$group": bson.M{
				"_id":   "$severidad",
				"count": bson.M{"$sum": 1},
			},
		},
	}

	cursor, err = s.db.Collection("errores").Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	erroresPorSeveridad := make(map[models.SeveridadError]int)
	for cursor.Next(ctx) {
		var result struct {
			ID    models.SeveridadError `bson:"_id"`
			Count int                   `bson:"count"`
		}
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
		erroresPorSeveridad[result.ID] = result.Count
	}

	// Contar errores resueltos y pendientes
	erroresResueltos, err := s.db.Collection("errores").CountDocuments(ctx, bson.M{
		"fecha_error": bson.M{
			"$gte": fechaInicio,
			"$lte": fechaFin,
		},
		"estado": "RESUELTO",
	})
	if err != nil {
		return nil, err
	}

	erroresPendientes, err := s.db.Collection("errores").CountDocuments(ctx, bson.M{
		"fecha_error": bson.M{
			"$gte": fechaInicio,
			"$lte": fechaFin,
		},
		"estado": "PENDIENTE",
	})
	if err != nil {
		return nil, err
	}

	// Calcular tiempo promedio de resolución
	pipeline = []bson.M{
		{
			"$match": bson.M{
				"fecha_error": bson.M{
					"$gte": fechaInicio,
					"$lte": fechaFin,
				},
				"estado": "RESUELTO",
			},
		},
		{
			"$group": bson.M{
				"_id": nil,
				"avg_time": bson.M{
					"$avg": bson.M{
						"$subtract": []string{"$fecha_resolucion", "$fecha_error"},
					},
				},
			},
		},
	}

	cursor, err = s.db.Collection("errores").Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tiempoPromedio int64
	if cursor.Next(ctx) {
		var result struct {
			AvgTime int64 `bson:"avg_time"`
		}
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
		tiempoPromedio = result.AvgTime
	}

	// Crear reporte
	reporte := &models.ReporteErrores{
		ID:                       models.GenerateErrorID(),
		FechaInicio:              fechaInicio,
		FechaFin:                 fechaFin,
		TotalErrores:             int(totalErrores),
		ErroresPorTipo:           erroresPorTipo,
		ErroresPorSeveridad:      erroresPorSeveridad,
		ErroresResueltos:         int(erroresResueltos),
		ErroresPendientes:        int(erroresPendientes),
		TiempoPromedioResolucion: tiempoPromedio,
		FechaGeneracion:          time.Now(),
	}

	// Guardar reporte
	_, err = s.db.Collection("reportes_errores").InsertOne(ctx, reporte)
	if err != nil {
		return nil, err
	}

	return reporte, nil
}

// Funciones auxiliares para obtener información del contexto
func getIPFromContext(ctx context.Context) string {
	// Implementar lógica para obtener IP del contexto
	return ""
}

func getUserAgentFromContext(ctx context.Context) string {
	// Implementar lógica para obtener User-Agent del contexto
	return ""
}

func getUsuarioIDFromContext(ctx context.Context) string {
	// Implementar lógica para obtener ID de usuario del contexto
	return ""
}

// ObtenerError obtiene un error por su ID
func (s *ErroresService) ObtenerError(ctx context.Context, errorID string) (*models.ErrorDetalle, error) {
	var errorDetalle models.ErrorDetalle
	err := s.db.Collection("errores").FindOne(ctx, bson.M{"_id": errorID}).Decode(&errorDetalle)
	if err != nil {
		return nil, err
	}
	return &errorDetalle, nil
}

// ListarErrores obtiene una lista paginada de errores según los filtros
func (s *ErroresService) ListarErrores(
	ctx context.Context,
	filtro map[string]interface{},
	limit, offset int,
) ([]models.ErrorDetalle, int64, error) {
	collection := s.db.Collection("errores")

	// Construir filtro
	filter := bson.M{}
	for k, v := range filtro {
		filter[k] = v
	}

	// Contar total
	total, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// Obtener resultados paginados
	options := options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(offset)).
		SetSort(bson.M{"fecha_error": -1})

	cursor, err := collection.Find(ctx, filter, options)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var errores []models.ErrorDetalle
	if err := cursor.All(ctx, &errores); err != nil {
		return nil, 0, err
	}

	return errores, total, nil
}

// ObtenerIntentosRecuperacion obtiene los intentos de recuperación de un error
func (s *ErroresService) ObtenerIntentosRecuperacion(ctx context.Context, errorID string) ([]models.IntentoRecuperacion, error) {
	collection := s.db.Collection("intentos_recuperacion")

	cursor, err := collection.Find(ctx, bson.M{"error_id": errorID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var intentos []models.IntentoRecuperacion
	if err := cursor.All(ctx, &intentos); err != nil {
		return nil, err
	}

	return intentos, nil
}

// ObtenerLogsError obtiene los logs de un error
func (s *ErroresService) ObtenerLogsError(ctx context.Context, errorID string) ([]models.LogError, error) {
	collection := s.db.Collection("logs_errores")

	cursor, err := collection.Find(ctx, bson.M{"error_id": errorID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var logs []models.LogError
	if err := cursor.All(ctx, &logs); err != nil {
		return nil, err
	}

	return logs, nil
}

// Helper function to generate IDs
func generateID() string {
	return models.GenerateErrorID()
}
