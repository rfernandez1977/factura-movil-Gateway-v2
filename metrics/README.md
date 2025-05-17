# Sistema de Métricas FMgo

Este directorio contiene la implementación del sistema de métricas para FMgo, utilizando Prometheus como backend de almacenamiento y visualización.

## Componentes

1. **metrics.go**: Definición de métricas generales del sistema
2. **middleware.go**: Middlewares para capturar métricas automáticamente
3. **redis.go**: Métricas específicas de Redis y colector
4. **prometheus.yml**: Configuración de Prometheus
5. **alerts.yml**: Reglas de alertas

## Métricas Implementadas

### Métricas de API
- `fmgo_api_request_duration_seconds`: Histograma de duración de peticiones
- `fmgo_api_requests_total`: Contador total de peticiones

### Métricas de Caché
- `fmgo_cache_hits_total`: Contador de aciertos en caché
- `fmgo_cache_misses_total`: Contador de fallos en caché
- `fmgo_cache_operation_seconds`: Histograma de latencia de operaciones

### Métricas de Sistema
- `fmgo_goroutines`: Gauge de goroutines activas
- `fmgo_memory_bytes`: Gauge de uso de memoria

### Métricas de SII
- `fmgo_sii_request_duration_seconds`: Histograma de duración de peticiones al SII
- `fmgo_sii_requests_total`: Contador de peticiones al SII

### Métricas de Redis
- `fmgo_redis_connections`: Gauge de conexiones activas
- `fmgo_redis_operation_latency_seconds`: Histograma de latencia de operaciones
- `fmgo_redis_memory_bytes`: Gauge de uso de memoria
- `fmgo_redis_keyspace_hits_total`: Contador de aciertos en keyspace
- `fmgo_redis_keyspace_misses_total`: Contador de fallos en keyspace
- `fmgo_redis_expired_keys_total`: Contador de claves expiradas
- `fmgo_redis_evicted_keys_total`: Contador de claves desalojadas

## Alertas Configuradas

1. **HighLatency**: P95 de latencia > 500ms
2. **HighErrorRate**: >5% de errores 500
3. **HighMemoryUsage**: Uso de memoria > 1GB
4. **HighCacheMissRate**: >30% de fallos en caché
5. **HighSIIErrorRate**: >10% de errores en SII
6. **HighRedisLatency**: P95 de latencia Redis > 50ms
7. **HighRedisMemoryUsage**: Uso de memoria Redis > 500MB
8. **HighRedisEvictionRate**: >10 claves desalojadas/s
9. **HighRedisKeyspaceMissRate**: >40% de fallos en keyspace

## Uso

### Integración en Handlers

```go
// Ejemplo de uso del middleware de métricas
router.Use(metrics.MetricsMiddleware)

// Ejemplo de uso del middleware de caché
err := metrics.CacheMetricsMiddleware(func() error {
    return cache.Get(key)
})

// Ejemplo de uso del middleware de SII
err := metrics.SIIMetricsMiddleware("enviarDTE", func() error {
    return siiClient.EnviarDTE(dte)
})

// Ejemplo de uso del colector de Redis
redisCollector := metrics.NewRedisMetricsCollector(redisClient)
redisCollector.Start(context.Background())
defer redisCollector.Stop()

// Ejemplo de uso del wrapper de operaciones Redis
err := metrics.RedisOperationWrapper("get", func() error {
    return redisClient.Get(ctx, key).Err()
})
```

### Endpoint de Métricas

El endpoint `/metrics` expone todas las métricas en formato Prometheus.

### Configuración de Prometheus

1. Copiar `prometheus.yml` y `alerts.yml` al directorio de configuración de Prometheus
2. Reiniciar Prometheus para aplicar la configuración
3. Acceder al dashboard de Prometheus (por defecto en `localhost:9090`)

## Dashboards Recomendados

Se recomienda crear los siguientes dashboards en Grafana:

1. **Overview**: Vista general del sistema
   - Tasa de peticiones
   - Latencia P95
   - Errores por minuto
   - Uso de memoria

2. **Caché**: Rendimiento del caché
   - Hit rate
   - Miss rate
   - Latencia de operaciones

3. **SII**: Monitoreo de interacciones con SII
   - Tasa de éxito/error
   - Latencia de operaciones
   - Distribución de operaciones

4. **Redis**: Monitoreo de Redis
   - Conexiones activas
   - Uso de memoria
   - Hit/Miss rate
   - Latencia de operaciones
   - Tasa de desalojo de claves
   - Tasa de expiración de claves

## Mantenimiento

- Revisar y ajustar umbrales de alertas según necesidad
- Monitorear uso de recursos y ajustar límites
- Mantener actualizadas las etiquetas de métricas
- Revisar y actualizar dashboards periódicamente
- Monitorear el rendimiento de Redis y ajustar la configuración según sea necesario
- Revisar periódicamente las tasas de aciertos/fallos en caché para optimizar estrategias de cacheo 