package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"FMgo/models"
	"github.com/go-redis/redis/v8"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	// TTLWorkflow tiempo de vida del workflow en caché
	TTLWorkflow = 24 * time.Hour
	// TTLRegistro tiempo de vida del registro en caché
	TTLRegistro = 1 * time.Hour
	// TTLMetricas tiempo de vida de las métricas en caché
	TTLMetricas = 12 * time.Hour
	// PrefijoWorkflow prefijo para las claves de workflow
	PrefijoWorkflow = "workflow:"
	// PrefijoRegistro prefijo para las claves de registro
	PrefijoRegistro = "registro:"
	// PrefijoMetricas prefijo para las claves de métricas
	PrefijoMetricas = "metricas:"
)

// IntegrationService maneja la lógica de integración con ERPs
type IntegrationService struct {
	db             *mongo.Database
	cache          *redis.Client
	async          *AsyncService
	queryOptimizer *QueryOptimizer
	circuitBreaker *CircuitBreaker
	eventSubject   *EventSubject
}

// NewIntegrationService crea una nueva instancia del servicio de integración
func NewIntegrationService(db *mongo.Database, redisClient *redis.Client, queue *amqp.Channel) *IntegrationService {
	return &IntegrationService{
		db:             db,
		cache:          redisClient,
		async:          NewAsyncService(queue, 5),
		queryOptimizer: NewQueryOptimizer(db),
		circuitBreaker: NewCircuitBreaker(5, 30*time.Second),
		eventSubject:   NewEventSubject(),
	}
}

// getWorkflowCacheKey genera una clave de caché para workflow
func (s *IntegrationService) getWorkflowCacheKey(erpID, entidad string) string {
	return fmt.Sprintf("%s%s:%s", PrefijoWorkflow, erpID, entidad)
}

// getRegistroCacheKey genera una clave de caché para registro
func (s *IntegrationService) getRegistroCacheKey(registroID string) string {
	return fmt.Sprintf("%s%s", PrefijoRegistro, registroID)
}

// getMetricasCacheKey genera una clave de caché para métricas
func (s *IntegrationService) getMetricasCacheKey(erpID string) string {
	return fmt.Sprintf("%s%s", PrefijoMetricas, erpID)
}

// IniciarSincronizacion inicia un proceso de sincronización
func (s *IntegrationService) IniciarSincronizacion(ctx context.Context, erpID string, entidad string, direccion string, datos map[string]interface{}) (*models.RegistroSincronizacion, error) {
	registro := &models.RegistroSincronizacion{
		ID:                 generateID(),
		ERPID:              erpID,
		Entidad:            entidad,
		Direccion:          direccion,
		Estado:             models.EstadoPendiente,
		DatosOriginales:    datos,
		FechaCreacion:      time.Now(),
		FechaActualizacion: time.Now(),
	}

	// Guardar en base de datos
	_, err := s.db.Collection("registros_sincronizacion").InsertOne(ctx, registro)
	if err != nil {
		return nil, fmt.Errorf("error al crear registro de sincronización: %v", err)
	}

	// Guardar en caché
	if err := s.guardarRegistroEnCache(ctx, registro); err != nil {
		log.Printf("error guardando registro en caché: %v", err)
	}

	return registro, nil
}

// guardarRegistroEnCache guarda un registro en caché
func (s *IntegrationService) guardarRegistroEnCache(ctx context.Context, registro *models.RegistroSincronizacion) error {
	data, err := json.Marshal(registro)
	if err != nil {
		return fmt.Errorf("error serializando registro: %w", err)
	}

	key := s.getRegistroCacheKey(registro.ID)
	return s.cache.Set(ctx, key, data, TTLRegistro).Err()
}

// obtenerRegistroDeCache obtiene un registro del caché
func (s *IntegrationService) obtenerRegistroDeCache(ctx context.Context, registroID string) (*models.RegistroSincronizacion, error) {
	key := s.getRegistroCacheKey(registroID)
	data, err := s.cache.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}

	var registro models.RegistroSincronizacion
	if err := json.Unmarshal(data, &registro); err != nil {
		return nil, fmt.Errorf("error deserializando registro: %w", err)
	}

	return &registro, nil
}

// ProcesarSincronizacion procesa un registro de sincronización
func (s *IntegrationService) ProcesarSincronizacion(ctx context.Context, registroID string) error {
	// Intentar obtener registro del caché
	registro, err := s.obtenerRegistroDeCache(ctx, registroID)
	if err != nil {
		// Si no está en caché, obtener de base de datos
		var registroTemp models.RegistroSincronizacion
		err = s.db.Collection("registros_sincronizacion").FindOne(ctx, bson.M{"_id": registroID}).Decode(&registroTemp)
		if err != nil {
			return fmt.Errorf("error al obtener registro: %v", err)
		}
		registro = &registroTemp

		// Guardar en caché para futuras consultas
		if err := s.guardarRegistroEnCache(ctx, registro); err != nil {
			log.Printf("error guardando registro en caché: %v", err)
		}
	}

	// Actualizar estado
	registro.Estado = models.EstadoEnProceso
	registro.FechaActualizacion = time.Now()

	// Actualizar en base de datos
	_, err = s.db.Collection("registros_sincronizacion").UpdateOne(ctx,
		bson.M{"_id": registroID},
		bson.M{"$set": bson.M{
			"estado":              registro.Estado,
			"fecha_actualizacion": registro.FechaActualizacion,
		}},
	)
	if err != nil {
		return fmt.Errorf("error al actualizar estado: %v", err)
	}

	// Actualizar en caché
	if err := s.guardarRegistroEnCache(ctx, registro); err != nil {
		log.Printf("error actualizando registro en caché: %v", err)
	}

	// Obtener workflow del caché
	key := s.getWorkflowCacheKey(registro.ERPID, registro.Entidad)
	var workflow models.Workflow
	data, err := s.cache.Get(ctx, key).Bytes()
	if err == nil {
		if err := json.Unmarshal(data, &workflow); err == nil {
			// Ejecutar workflow
			if err := s.ejecutarWorkflow(ctx, &workflow, registro); err != nil {
				s.manejarError(ctx, registro, err)
			} else {
				registro.Estado = models.EstadoCompletado
			}
		}
	}

	// Si no está en caché, obtener de base de datos
	if err != nil {
		err = s.db.Collection("workflows").FindOne(ctx, bson.M{
			"erp_id":  registro.ERPID,
			"entidad": registro.Entidad,
		}).Decode(&workflow)
		if err != nil {
			return fmt.Errorf("error al obtener workflow: %v", err)
		}

		// Guardar workflow en caché
		if data, err := json.Marshal(workflow); err == nil {
			if err := s.cache.Set(ctx, key, data, TTLWorkflow).Err(); err != nil {
				log.Printf("error guardando workflow en caché: %v", err)
			}
		}

		// Ejecutar workflow
		if err := s.ejecutarWorkflow(ctx, &workflow, registro); err != nil {
			s.manejarError(ctx, registro, err)
		} else {
			registro.Estado = models.EstadoCompletado
		}
	}

	// Actualizar registro final
	registro.FechaActualizacion = time.Now()
	_, err = s.db.Collection("registros_sincronizacion").UpdateOne(ctx,
		bson.M{"_id": registroID},
		bson.M{"$set": registro},
	)
	if err != nil {
		return fmt.Errorf("error al actualizar registro: %v", err)
	}

	// Actualizar caché final
	if err := s.guardarRegistroEnCache(ctx, registro); err != nil {
		log.Printf("error actualizando registro en caché: %v", err)
	}

	return nil
}

// manejarError maneja un error durante la ejecución del workflow
func (s *IntegrationService) manejarError(ctx context.Context, registro *models.RegistroSincronizacion, err error) {
	registro.Errores = append(registro.Errores, models.ErrorSincronizacion{
		Codigo:    "WORKFLOW_ERROR",
		Mensaje:   err.Error(),
		Timestamp: time.Now(),
	})
	registro.Estado = models.EstadoError
}

// ejecutarWorkflow ejecuta un workflow paso a paso
func (s *IntegrationService) ejecutarWorkflow(ctx context.Context, workflow *models.Workflow, registro *models.RegistroSincronizacion) error {
	for _, paso := range workflow.Pasos {
		// Ejecutar paso
		err := s.ejecutarPaso(ctx, paso, registro)
		if err != nil {
			return fmt.Errorf("error en paso %s: %v", paso.Nombre, err)
		}

		// Verificar condiciones de salida
		if !s.verificarCondiciones(paso.CondicionesSalida, registro.DatosTransformados) {
			return fmt.Errorf("no se cumplieron las condiciones de salida en el paso %s", paso.Nombre)
		}

		// Actualizar registro en caché después de cada paso
		if err := s.guardarRegistroEnCache(ctx, registro); err != nil {
			log.Printf("error actualizando registro en caché: %v", err)
		}
	}

	return nil
}

// ejecutarPaso ejecuta un paso específico del workflow
func (s *IntegrationService) ejecutarPaso(ctx context.Context, paso models.PasoWorkflow, registro *models.RegistroSincronizacion) error {
	switch paso.Tipo {
	case "VALIDACION":
		return s.ejecutarValidacion(ctx, paso, registro)
	case "TRANSFORMACION":
		return s.ejecutarTransformacion(ctx, paso, registro)
	case "SINCRONIZACION":
		return s.ejecutarSincronizacion(ctx, paso, registro)
	default:
		return fmt.Errorf("tipo de paso no soportado: %s", paso.Tipo)
	}
}

// ejecutarValidacion ejecuta una validación
func (s *IntegrationService) ejecutarValidacion(ctx context.Context, paso models.PasoWorkflow, registro *models.RegistroSincronizacion) error {
	// Implementar lógica de validación
	return nil
}

// ejecutarTransformacion ejecuta una transformación
func (s *IntegrationService) ejecutarTransformacion(ctx context.Context, paso models.PasoWorkflow, registro *models.RegistroSincronizacion) error {
	// Implementar lógica de transformación
	return nil
}

// ejecutarSincronizacion ejecuta una sincronización
func (s *IntegrationService) ejecutarSincronizacion(ctx context.Context, paso models.PasoWorkflow, registro *models.RegistroSincronizacion) error {
	// Implementar lógica de sincronización
	return nil
}

// verificarCondiciones verifica si se cumplen las condiciones de salida
func (s *IntegrationService) verificarCondiciones(condiciones []models.Condicion, datos map[string]interface{}) bool {
	for _, condicion := range condiciones {
		valor, ok := datos[condicion.Campo]
		if !ok {
			return false
		}

		switch condicion.Operador {
		case "==":
			if valor != condicion.Valor {
				return false
			}
		case "!=":
			if valor == condicion.Valor {
				return false
			}
		case ">":
			if !s.compararMayor(valor, condicion.Valor) {
				return false
			}
		case "<":
			if !s.compararMenor(valor, condicion.Valor) {
				return false
			}
		default:
			log.Printf("Operador no soportado: %s", condicion.Operador)
			return false
		}
	}

	return true
}

// compararMayor compara si un valor es mayor que otro
func (s *IntegrationService) compararMayor(a, b interface{}) bool {
	switch a.(type) {
	case int:
		return a.(int) > b.(int)
	case float64:
		return a.(float64) > b.(float64)
	case time.Time:
		return a.(time.Time).After(b.(time.Time))
	default:
		return false
	}
}

// compararMenor compara si un valor es menor que otro
func (s *IntegrationService) compararMenor(a, b interface{}) bool {
	switch a.(type) {
	case int:
		return a.(int) < b.(int)
	case float64:
		return a.(float64) < b.(float64)
	case time.Time:
		return a.(time.Time).Before(b.(time.Time))
	default:
		return false
	}
}

// RegistrarMetrica registra una métrica de integración
func (s *IntegrationService) RegistrarMetrica(ctx context.Context, metrica *models.MetricaIntegracion) error {
	metrica.ID = generateID()
	metrica.Timestamp = time.Now()

	// Guardar en base de datos
	_, err := s.db.Collection("metricas_integracion").InsertOne(ctx, metrica)
	if err != nil {
		return fmt.Errorf("error al registrar métrica: %v", err)
	}

	// Actualizar métricas en caché
	key := s.getMetricasCacheKey(metrica.ERPID)
	var metricas []models.MetricaIntegracion
	data, err := s.cache.Get(ctx, key).Bytes()
	if err == nil {
		if err := json.Unmarshal(data, &metricas); err == nil {
			metricas = append(metricas, *metrica)
			if data, err := json.Marshal(metricas); err == nil {
				if err := s.cache.Set(ctx, key, data, TTLMetricas).Err(); err != nil {
					log.Printf("error actualizando métricas en caché: %v", err)
				}
			}
		}
	} else {
		// Si no hay métricas en caché, crear nueva lista
		metricas = []models.MetricaIntegracion{*metrica}
		if data, err := json.Marshal(metricas); err == nil {
			if err := s.cache.Set(ctx, key, data, TTLMetricas).Err(); err != nil {
				log.Printf("error guardando métricas en caché: %v", err)
			}
		}
	}

	return nil
}

// RegistrarAlerta registra una alerta
func (s *IntegrationService) RegistrarAlerta(ctx context.Context, alerta *models.Alerta) error {
	alerta.ID = generateID()
	alerta.FechaCreacion = time.Now()
	alerta.FechaActualizacion = time.Now()

	// Guardar en base de datos
	_, err := s.db.Collection("alertas").InsertOne(ctx, alerta)
	if err != nil {
		return fmt.Errorf("error al registrar alerta: %v", err)
	}

	// Notificar a los observadores
	s.eventSubject.Notify(ctx, "ALERTA", alerta)

	return nil
}

// LimpiarCache limpia el caché del servicio
func (s *IntegrationService) LimpiarCache(ctx context.Context) error {
	var cursor uint64
	var keys []string

	// Obtener todas las claves con los prefijos
	for {
		var result []string
		var err error
		result, cursor, err = s.cache.Scan(ctx, cursor, PrefijoWorkflow+"*", 100).Result()
		if err != nil {
			return fmt.Errorf("error escaneando claves de workflow: %w", err)
		}
		keys = append(keys, result...)

		result, cursor, err = s.cache.Scan(ctx, cursor, PrefijoRegistro+"*", 100).Result()
		if err != nil {
			return fmt.Errorf("error escaneando claves de registro: %w", err)
		}
		keys = append(keys, result...)

		result, cursor, err = s.cache.Scan(ctx, cursor, PrefijoMetricas+"*", 100).Result()
		if err != nil {
			return fmt.Errorf("error escaneando claves de métricas: %w", err)
		}
		keys = append(keys, result...)

		if cursor == 0 {
			break
		}
	}

	// Eliminar todas las claves encontradas
	if len(keys) > 0 {
		if err := s.cache.Del(ctx, keys...).Err(); err != nil {
			return fmt.Errorf("error eliminando claves del caché: %w", err)
		}
	}

	return nil
}

func (s *IntegrationService) StartSync(erpID primitive.ObjectID) error {
	return s.circuitBreaker.Execute(func() error {
		// Obtener configuración del ERP
		config, err := s.getERPConfig(context.Background(), erpID)
		if err != nil {
			s.eventSubject.Notify(Event{
				Type:    "ERROR",
				Payload: err,
			})
			return err
		}

		// Crear registro de sincronización
		record := models.SyncRecord{
			ID:        primitive.NewObjectID(),
			ERPID:     erpID,
			Status:    models.SyncStatusInProgress,
			StartTime: time.Now(),
		}

		// Publicar evento de inicio
		s.eventSubject.Notify(Event{
			Type:    "SYNC_STARTED",
			Payload: record,
		})

		// Ejecutar sincronización
		err = s.executeSync(context.Background(), config, &record)
		if err != nil {
			s.eventSubject.Notify(Event{
				Type:    "ERROR",
				Payload: err,
			})
			return err
		}

		// Publicar evento de finalización
		s.eventSubject.Notify(Event{
			Type:    "SYNC_COMPLETED",
			Payload: record,
		})

		return nil
	})
}

func (s *IntegrationService) getERPConfig(ctx context.Context, erpID primitive.ObjectID) (*models.ConfiguracionERP, error) {
	// Intentar obtener de la caché
	var config models.ConfiguracionERP
	key := s.cache.CacheKey("erp_configs", erpID)
	err := s.cache.GetOrSet(ctx, key, &config, func() (interface{}, error) {
		// Si no está en caché, obtener de la base de datos
		var result models.ConfiguracionERP
		err := s.db.Collection("erp_configs").FindOne(ctx, bson.M{"_id": erpID}).Decode(&result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func (s *IntegrationService) executeSync(ctx context.Context, config *models.ConfiguracionERP, record *models.SyncRecord) error {
	// Optimizar la consulta de documentos
	cursor, err := s.queryOptimizer.OptimizeFind(ctx, "documents", bson.M{
		"erp_id": config.ID,
		"status": "pending",
	})
	if err != nil {
		return err
	}

	var documents []models.Document
	if err = cursor.All(ctx, &documents); err != nil {
		return err
	}

	// Procesar documentos de forma asíncrona
	for _, doc := range documents {
		task := NewTask("process_document", doc)
		if err := s.async.Enqueue(ctx, task.Type, task); err != nil {
			return err
		}
	}

	return nil
}
