package metrics

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// RedisConnections es un gauge que mide el número de conexiones activas a Redis
	RedisConnections = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "fmgo_redis_connections",
		Help: "Número de conexiones activas a Redis",
	})

	// RedisOperationLatency mide la latencia de las operaciones de Redis
	RedisOperationLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "fmgo_redis_operation_latency_seconds",
			Help:    "Latencia de operaciones Redis en segundos",
			Buckets: []float64{.001, .002, .005, .01, .025, .05, .1, .25, .5},
		},
		[]string{"operation"},
	)

	// RedisMemoryUsage mide el uso de memoria de Redis
	RedisMemoryUsage = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "fmgo_redis_memory_bytes",
		Help: "Uso de memoria de Redis en bytes",
	})

	// RedisKeyspaceHits mide los aciertos en el keyspace
	RedisKeyspaceHits = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "fmgo_redis_keyspace_hits_total",
		Help: "Número total de aciertos en el keyspace de Redis",
	})

	// RedisKeyspaceMisses mide los fallos en el keyspace
	RedisKeyspaceMisses = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "fmgo_redis_keyspace_misses_total",
		Help: "Número total de fallos en el keyspace de Redis",
	})

	// RedisExpiredKeys mide el número de claves expiradas
	RedisExpiredKeys = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "fmgo_redis_expired_keys_total",
		Help: "Número total de claves expiradas en Redis",
	})

	// RedisEvictedKeys mide el número de claves desalojadas
	RedisEvictedKeys = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "fmgo_redis_evicted_keys_total",
		Help: "Número total de claves desalojadas en Redis",
	})
)

func init() {
	// Registro de métricas de Redis
	prometheus.MustRegister(RedisConnections)
	prometheus.MustRegister(RedisOperationLatency)
	prometheus.MustRegister(RedisMemoryUsage)
	prometheus.MustRegister(RedisKeyspaceHits)
	prometheus.MustRegister(RedisKeyspaceMisses)
	prometheus.MustRegister(RedisExpiredKeys)
	prometheus.MustRegister(RedisEvictedKeys)
}

// RedisMetricsCollector es un colector que actualiza las métricas de Redis periódicamente
type RedisMetricsCollector struct {
	client *redis.Client
	done   chan struct{}
}

// NewRedisMetricsCollector crea un nuevo colector de métricas de Redis
func NewRedisMetricsCollector(client *redis.Client) *RedisMetricsCollector {
	return &RedisMetricsCollector{
		client: client,
		done:   make(chan struct{}),
	}
}

// Start inicia la recolección de métricas
func (c *RedisMetricsCollector) Start(ctx context.Context) {
	ticker := time.NewTicker(15 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				c.collect(ctx)
			case <-c.done:
				ticker.Stop()
				return
			case <-ctx.Done():
				ticker.Stop()
				return
			}
		}
	}()
}

// Stop detiene la recolección de métricas
func (c *RedisMetricsCollector) Stop() {
	close(c.done)
}

// collect recolecta las métricas de Redis
func (c *RedisMetricsCollector) collect(ctx context.Context) {
	// Obtener estadísticas de Redis
	info := c.client.Info(ctx, "stats", "memory").Val()
	stats := parseRedisInfo(info)

	// Actualizar métricas
	if v, ok := stats["connected_clients"]; ok {
		RedisConnections.Set(v)
	}
	if v, ok := stats["used_memory"]; ok {
		RedisMemoryUsage.Set(v)
	}
	if v, ok := stats["keyspace_hits"]; ok {
		RedisKeyspaceHits.Add(v)
	}
	if v, ok := stats["keyspace_misses"]; ok {
		RedisKeyspaceMisses.Add(v)
	}
	if v, ok := stats["expired_keys"]; ok {
		RedisExpiredKeys.Add(v)
	}
	if v, ok := stats["evicted_keys"]; ok {
		RedisEvictedKeys.Add(v)
	}
}

// RedisOperationWrapper envuelve una operación de Redis para medir su latencia
func RedisOperationWrapper(operation string, fn func() error) error {
	start := time.Now()
	err := fn()
	duration := time.Since(start).Seconds()

	RedisOperationLatency.WithLabelValues(operation).Observe(duration)

	return err
}

// parseRedisInfo parsea la salida del comando INFO de Redis
func parseRedisInfo(info string) map[string]float64 {
	result := make(map[string]float64)
	lines := strings.Split(info, "\n")

	for _, line := range lines {
		// Ignorar líneas vacías, comentarios o secciones
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "[") {
			continue
		}

		// Dividir la línea en clave y valor
		parts := strings.Split(line, ":")
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Convertir el valor a float64 si es posible
		if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
			result[key] = floatVal
		}
	}

	return result
}
