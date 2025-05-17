package integration

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Documento representa un documento tributario electrónico
type Documento struct {
	ID                string
	TipoDTE           string
	Folio             int64
	FechaEmision      time.Time
	RutEmisor         string
	RazonEmisor       string
	GiroEmisor        string
	DireccionEmisor   string
	ComunaEmisor      string
	RutReceptor       string
	RazonReceptor     string
	GiroReceptor      string
	DireccionReceptor string
	ComunaReceptor    string
	MontoNeto         float64
	MontoExento       float64
	TasaIVA           float64
	MontoIVA          float64
	MontoTotal        float64
	Estado            string
	TrackID           string
	EstadoSII         string
}

// crearDocumentoPrueba crea un documento de prueba
func crearDocumentoPrueba(tipoDTE string) *Documento {
	return &Documento{
		ID:           "TEST-" + tipoDTE,
		TipoDTE:      tipoDTE,
		Folio:        1,
		FechaEmision: time.Now(),
		RutEmisor:    "76.123.456-7",
		RazonEmisor:  "Empresa de Prueba",
		Estado:       "BORRADOR",
	}
}

// TestRedisIntegration prueba la integración completa con Redis
func TestRedisIntegration(t *testing.T) {
	// Configuración de Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // sin contraseña por defecto
		DB:       0,  // base de datos por defecto
	})

	ctx := context.Background()

	// Limpiar datos de prueba anteriores
	t.Cleanup(func() {
		rdb.FlushDB(ctx)
	})

	// Test 1: Conexión básica
	t.Run("TestConexionRedis", func(t *testing.T) {
		_, err := rdb.Ping(ctx).Result()
		require.NoError(t, err, "Debe poder conectarse a Redis")
	})

	// Test 2: Almacenamiento y recuperación de documento
	t.Run("TestAlmacenamientoDocumento", func(t *testing.T) {
		doc := &Documento{
			ID:           "TEST-001",
			TipoDTE:      "33",
			Folio:        1,
			FechaEmision: time.Now(),
			RutEmisor:    "76.123.456-7",
			Estado:       "BORRADOR",
		}

		// Convertir a JSON
		docJSON, err := json.Marshal(doc)
		require.NoError(t, err, "Debe poder convertir documento a JSON")

		// Guardar en Redis
		err = rdb.Set(ctx, "doc:"+doc.ID, docJSON, 24*time.Hour).Err()
		require.NoError(t, err, "Debe poder guardar documento en Redis")

		// Recuperar de Redis
		val, err := rdb.Get(ctx, "doc:"+doc.ID).Result()
		require.NoError(t, err, "Debe poder recuperar documento de Redis")

		var docRecuperado Documento
		err = json.Unmarshal([]byte(val), &docRecuperado)
		require.NoError(t, err, "Debe poder convertir JSON a documento")

		assert.Equal(t, doc.ID, docRecuperado.ID, "Los IDs deben coincidir")
		assert.Equal(t, doc.TipoDTE, docRecuperado.TipoDTE, "Los tipos DTE deben coincidir")
	})

	// Test 3: Expiración de datos
	t.Run("TestExpiracionDatos", func(t *testing.T) {
		key := "test:expiracion"
		value := "datos_temporales"

		// Guardar con expiración de 1 segundo
		err := rdb.Set(ctx, key, value, 1*time.Second).Err()
		require.NoError(t, err, "Debe poder guardar datos con expiración")

		// Verificar que existe
		val, err := rdb.Get(ctx, key).Result()
		require.NoError(t, err, "Debe poder recuperar datos antes de expirar")
		assert.Equal(t, value, val, "Los valores deben coincidir")

		// Esperar a que expire
		time.Sleep(2 * time.Second)

		// Verificar que ya no existe
		_, err = rdb.Get(ctx, key).Result()
		assert.Equal(t, redis.Nil, err, "Los datos deben haber expirado")
	})

	// Test 4: Manejo de concurrencia
	t.Run("TestConcurrencia", func(t *testing.T) {
		key := "test:concurrencia"

		// Simular operaciones concurrentes
		for i := 0; i < 100; i++ {
			go func(i int) {
				err := rdb.Set(ctx, key, i, 0).Err()
				require.NoError(t, err, "Debe manejar operaciones concurrentes")
			}(i)
		}

		// Esperar un momento para que terminen las operaciones
		time.Sleep(100 * time.Millisecond)

		// Verificar que el valor existe
		_, err := rdb.Get(ctx, key).Result()
		require.NoError(t, err, "Debe poder recuperar valor después de operaciones concurrentes")
	})

	// Test 5: Integración con flujo DTE
	t.Run("TestIntegracionDTE", func(t *testing.T) {
		// Crear documento de prueba
		doc := crearDocumentoPrueba("33")

		// Simular proceso de DTE usando caché
		// 1. Guardar documento en caché
		docJSON, err := json.Marshal(doc)
		require.NoError(t, err, "Debe poder convertir documento a JSON")

		err = rdb.Set(ctx, "dte:"+doc.ID, docJSON, 24*time.Hour).Err()
		require.NoError(t, err, "Debe poder guardar DTE en caché")

		// 2. Simular procesamiento
		doc.Estado = "PROCESADO"
		docJSON, err = json.Marshal(doc)
		require.NoError(t, err, "Debe poder actualizar documento en JSON")

		err = rdb.Set(ctx, "dte:"+doc.ID, docJSON, 24*time.Hour).Err()
		require.NoError(t, err, "Debe poder actualizar DTE en caché")

		// 3. Verificar estado final
		val, err := rdb.Get(ctx, "dte:"+doc.ID).Result()
		require.NoError(t, err, "Debe poder recuperar DTE de caché")

		var docFinal Documento
		err = json.Unmarshal([]byte(val), &docFinal)
		require.NoError(t, err, "Debe poder convertir caché a documento")

		assert.Equal(t, "PROCESADO", docFinal.Estado, "El estado final debe ser PROCESADO")
	})
}
