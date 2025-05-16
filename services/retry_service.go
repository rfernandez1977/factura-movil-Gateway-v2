package services

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"

	"github.com/fmgo/models"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	// TTLReintento tiempo de vida de los reintentos en caché
	TTLReintento = 24 * time.Hour
	// TTLFlujo tiempo de vida de los flujos en caché
	TTLFlujo = 12 * time.Hour
	// TTLPaso tiempo de vida de los pasos en caché
	TTLPaso = 12 * time.Hour
	// PrefijoReintento prefijo para las claves de reintento
	PrefijoReintento = "reintento:"
	// PrefijoFlujo prefijo para las claves de flujo
	PrefijoFlujo = "flujo:"
	// PrefijoPaso prefijo para las claves de paso
	PrefijoPaso = "paso:"
	// PrefijoIntentos prefijo para las claves de intentos
	PrefijoIntentos = "intentos:"
)

// RetryService maneja los reintentos de operaciones
type RetryService struct {
	cache *redis.Client
	db    *mongo.Database
}

// RetryConfig contiene la configuración de reintentos
type RetryConfig struct {
	MaxAttempts     int           // Número máximo de intentos
	InitialDelay    time.Duration // Retraso inicial
	MaxDelay        time.Duration // Retraso máximo
	BackoffFactor   float64       // Factor de crecimiento exponencial
	JitterFactor    float64       // Factor de aleatoriedad
	RetryableErrors []string      // Errores que permiten reintento
}

// DefaultRetryConfig retorna una configuración por defecto
func DefaultRetryConfig() *RetryConfig {
	return &RetryConfig{
		MaxAttempts:   5,
		InitialDelay:  1 * time.Second,
		MaxDelay:      30 * time.Second,
		BackoffFactor: 2.0,
		JitterFactor:  0.1,
		RetryableErrors: []string{
			"timeout",
			"connection refused",
			"connection reset",
			"network error",
		},
	}
}

// NewRetryService crea una nueva instancia del servicio de reintentos
func NewRetryService(redisClient *redis.Client, db *mongo.Database) *RetryService {
	return &RetryService{
		cache: redisClient,
		db:    db,
	}
}

// getReintentoCacheKey genera una clave de caché para reintento
func (s *RetryService) getReintentoCacheKey(id primitive.ObjectID) string {
	return fmt.Sprintf("%s%s", PrefijoReintento, id.Hex())
}

// getFlujoCacheKey genera una clave de caché para flujo
func (s *RetryService) getFlujoCacheKey(id primitive.ObjectID) string {
	return fmt.Sprintf("%s%s", PrefijoFlujo, id.Hex())
}

// getPasoCacheKey genera una clave de caché para paso
func (s *RetryService) getPasoCacheKey(id primitive.ObjectID) string {
	return fmt.Sprintf("%s%s", PrefijoPaso, id.Hex())
}

// getIntentosCacheKey genera una clave de caché para intentos
func (s *RetryService) getIntentosCacheKey(operationID string) string {
	return fmt.Sprintf("%s%s", PrefijoIntentos, operationID)
}

// AgregarReintento agrega un elemento a la cola de reintentos
func (s *RetryService) AgregarReintento(ctx context.Context, reintento *models.ColaReintentos) error {
	collection := s.db.Collection("cola_reintentos")

	if reintento.ID.IsZero() {
		reintento.ID = primitive.NewObjectID()
	}

	// Guardar en base de datos
	_, err := collection.InsertOne(ctx, reintento)
	if err != nil {
		return err
	}

	// Guardar en caché
	if err := s.guardarReintentoEnCache(ctx, reintento); err != nil {
		// Solo logear error de caché, no afecta operación principal
		fmt.Printf("error guardando reintento en caché: %v\n", err)
	}

	return nil
}

// guardarReintentoEnCache guarda un reintento en caché
func (s *RetryService) guardarReintentoEnCache(ctx context.Context, reintento *models.ColaReintentos) error {
	data, err := json.Marshal(reintento)
	if err != nil {
		return fmt.Errorf("error serializando reintento: %w", err)
	}

	key := s.getReintentoCacheKey(reintento.ID)
	return s.cache.Set(ctx, key, data, TTLReintento).Err()
}

// obtenerReintentoDeCache obtiene un reintento del caché
func (s *RetryService) obtenerReintentoDeCache(ctx context.Context, id primitive.ObjectID) (*models.ColaReintentos, error) {
	key := s.getReintentoCacheKey(id)
	data, err := s.cache.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}

	var reintento models.ColaReintentos
	if err := json.Unmarshal(data, &reintento); err != nil {
		return nil, fmt.Errorf("error deserializando reintento: %w", err)
	}

	return &reintento, nil
}

// ProcesarReintentos procesa los reintentos pendientes
func (s *RetryService) ProcesarReintentos(ctx context.Context) error {
	collection := s.db.Collection("cola_reintentos")

	// Obtener reintentos pendientes
	cursor, err := collection.Find(ctx, bson.M{
		"estado": models.EstadoReintentoPendiente,
		"fecha_proximo_intento": bson.M{
			"$lte": time.Now(),
		},
	})
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	var reintentos []models.ColaReintentos
	if err = cursor.All(ctx, &reintentos); err != nil {
		return err
	}

	// Procesar cada reintento
	for _, reintento := range reintentos {
		if err := s.procesarReintento(ctx, &reintento); err != nil {
			// Registrar el error pero continuar con el siguiente reintento
			s.registrarErrorReintento(ctx, &reintento, err)
		}
	}

	return nil
}

// procesarReintento procesa un reintento individual
func (s *RetryService) procesarReintento(ctx context.Context, reintento *models.ColaReintentos) error {
	// Obtener el flujo y paso correspondientes
	flujo, err := s.obtenerFlujo(ctx, reintento.FlujoID)
	if err != nil {
		return err
	}

	paso, err := s.obtenerPaso(ctx, reintento.PasoID)
	if err != nil {
		return err
	}

	// Ejecutar el paso
	if err := s.ejecutarPaso(ctx, flujo, paso); err != nil {
		// Si el error persiste, actualizar el reintento
		if reintento.Intento < paso.MaxReintentos {
			reintento.Intento++
			reintento.FechaProximoIntento = time.Now().Add(s.calcularIntervaloReintento(reintento.Intento))
			reintento.Estado = models.EstadoReintentoPendiente
		} else {
			reintento.Estado = models.EstadoReintentoFallido
		}
		reintento.Error = err.Error()
	} else {
		// Éxito
		reintento.Estado = models.EstadoReintentoCompletado
	}

	// Actualizar en base de datos
	collection := s.db.Collection("cola_reintentos")
	_, err = collection.ReplaceOne(ctx, bson.M{"_id": reintento.ID}, reintento)
	if err != nil {
		return err
	}

	// Actualizar en caché
	if err := s.guardarReintentoEnCache(ctx, reintento); err != nil {
		fmt.Printf("error actualizando reintento en caché: %v\n", err)
	}

	return nil
}

// obtenerFlujo obtiene un flujo por su ID
func (s *RetryService) obtenerFlujo(ctx context.Context, id primitive.ObjectID) (*models.FlujoIntegracion, error) {
	// Intentar obtener del caché
	key := s.getFlujoCacheKey(id)
	data, err := s.cache.Get(ctx, key).Bytes()
	if err == nil {
		var flujo models.FlujoIntegracion
		if err := json.Unmarshal(data, &flujo); err == nil {
			return &flujo, nil
		}
	}

	// Si no está en caché, obtener de base de datos
	collection := s.db.Collection("flujos_integracion")
	var flujo models.FlujoIntegracion
	err = collection.FindOne(ctx, bson.M{"_id": id}).Decode(&flujo)
	if err != nil {
		return nil, err
	}

	// Guardar en caché
	if data, err := json.Marshal(flujo); err == nil {
		if err := s.cache.Set(ctx, key, data, TTLFlujo).Err(); err != nil {
			fmt.Printf("error guardando flujo en caché: %v\n", err)
		}
	}

	return &flujo, nil
}

// obtenerPaso obtiene un paso por su ID
func (s *RetryService) obtenerPaso(ctx context.Context, id primitive.ObjectID) (*models.PasoFlujo, error) {
	// Intentar obtener del caché
	key := s.getPasoCacheKey(id)
	data, err := s.cache.Get(ctx, key).Bytes()
	if err == nil {
		var paso models.PasoFlujo
		if err := json.Unmarshal(data, &paso); err == nil {
			return &paso, nil
		}
	}

	// Si no está en caché, obtener de base de datos
	collection := s.db.Collection("pasos_flujo")
	var paso models.PasoFlujo
	err = collection.FindOne(ctx, bson.M{"_id": id}).Decode(&paso)
	if err != nil {
		return nil, err
	}

	// Guardar en caché
	if data, err := json.Marshal(paso); err == nil {
		if err := s.cache.Set(ctx, key, data, TTLPaso).Err(); err != nil {
			fmt.Printf("error guardando paso en caché: %v\n", err)
		}
	}

	return &paso, nil
}

// ejecutarPaso ejecuta un paso del flujo
func (s *RetryService) ejecutarPaso(ctx context.Context, flujo *models.FlujoIntegracion, paso *models.PasoFlujo) error {
	// TODO: Implementar la lógica de ejecución del paso
	return nil
}

// calcularIntervaloReintento calcula el intervalo para el próximo reintento
func (s *RetryService) calcularIntervaloReintento(intento int) time.Duration {
	// Implementar una estrategia de backoff exponencial
	base := time.Second * 5
	return base * time.Duration(1<<uint(intento-1))
}

// registrarErrorReintento registra un error en el reintento
func (s *RetryService) registrarErrorReintento(ctx context.Context, reintento *models.ColaReintentos, err error) {
	// TODO: Implementar registro de errores
}

// Retry ejecuta una operación con reintentos
func (s *RetryService) Retry(ctx context.Context, operationID string, config *RetryConfig, operation func() error) error {
	// Obtener intentos previos
	attempts, err := s.getAttempts(ctx, operationID)
	if err != nil {
		return fmt.Errorf("error al obtener intentos previos: %v", err)
	}

	// Verificar límite de intentos
	if attempts >= config.MaxAttempts {
		return fmt.Errorf("se excedió el número máximo de intentos (%d)", config.MaxAttempts)
	}

	// Ejecutar operación
	err = operation()
	if err == nil {
		// Éxito, limpiar intentos
		s.clearAttempts(ctx, operationID)
		return nil
	}

	// Verificar si el error permite reintento
	if !s.isRetryableError(err, config.RetryableErrors) {
		return err
	}

	// Calcular retraso con backoff exponencial y jitter
	delay := s.calculateDelay(attempts, config)

	// Registrar intento
	if err := s.recordAttempt(ctx, operationID); err != nil {
		return fmt.Errorf("error al registrar intento: %v", err)
	}

	// Esperar antes del siguiente intento
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(delay):
		return s.Retry(ctx, operationID, config, operation)
	}
}

// getAttempts obtiene el número de intentos previos
func (s *RetryService) getAttempts(ctx context.Context, operationID string) (int, error) {
	key := s.getIntentosCacheKey(operationID)
	attempts, err := s.cache.Get(ctx, key).Int()
	if err == redis.Nil {
		return 0, nil
	}
	return attempts, err
}

// recordAttempt registra un nuevo intento
func (s *RetryService) recordAttempt(ctx context.Context, operationID string) error {
	key := s.getIntentosCacheKey(operationID)
	pipe := s.cache.Pipeline()
	pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, TTLReintento)
	_, err := pipe.Exec(ctx)
	return err
}

// clearAttempts limpia los intentos registrados
func (s *RetryService) clearAttempts(ctx context.Context, operationID string) error {
	key := s.getIntentosCacheKey(operationID)
	return s.cache.Del(ctx, key).Err()
}

// isRetryableError verifica si un error permite reintento
func (s *RetryService) isRetryableError(err error, retryableErrors []string) bool {
	errStr := err.Error()
	for _, retryableErr := range retryableErrors {
		if strings.Contains(strings.ToLower(errStr), strings.ToLower(retryableErr)) {
			return true
		}
	}
	return false
}

// calculateDelay calcula el retraso para el siguiente intento
func (s *RetryService) calculateDelay(attempts int, config *RetryConfig) time.Duration {
	// Calcular retraso base con backoff exponencial
	delay := float64(config.InitialDelay) * math.Pow(config.BackoffFactor, float64(attempts))

	// Aplicar límite máximo
	if delay > float64(config.MaxDelay) {
		delay = float64(config.MaxDelay)
	}

	// Aplicar jitter
	jitter := (rand.Float64()*2 - 1) * config.JitterFactor * delay
	delay = delay + jitter

	return time.Duration(delay)
}

// LimpiarCache limpia el caché del servicio
func (s *RetryService) LimpiarCache(ctx context.Context) error {
	var cursor uint64
	var keys []string

	// Obtener todas las claves con los prefijos
	for {
		var result []string
		var err error
		result, cursor, err = s.cache.Scan(ctx, cursor, PrefijoReintento+"*", 100).Result()
		if err != nil {
			return fmt.Errorf("error escaneando claves de reintento: %w", err)
		}
		keys = append(keys, result...)

		result, cursor, err = s.cache.Scan(ctx, cursor, PrefijoFlujo+"*", 100).Result()
		if err != nil {
			return fmt.Errorf("error escaneando claves de flujo: %w", err)
		}
		keys = append(keys, result...)

		result, cursor, err = s.cache.Scan(ctx, cursor, PrefijoPaso+"*", 100).Result()
		if err != nil {
			return fmt.Errorf("error escaneando claves de paso: %w", err)
		}
		keys = append(keys, result...)

		result, cursor, err = s.cache.Scan(ctx, cursor, PrefijoIntentos+"*", 100).Result()
		if err != nil {
			return fmt.Errorf("error escaneando claves de intentos: %w", err)
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
