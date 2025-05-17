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

// TestDTERedisFlow prueba el flujo completo de DTE con integración Redis
func TestDTERedisFlow(t *testing.T) {
	// Configuración de Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	ctx := context.Background()

	// Limpiar datos de prueba
	t.Cleanup(func() {
		rdb.FlushDB(ctx)
	})

	// Test del flujo completo
	t.Run("TestFlujoCompleto", func(t *testing.T) {
		// 1. Crear documento
		doc := crearDocumentoPrueba("33")
		cacheKey := "dte:" + doc.ID

		// 2. Guardar en caché (estado inicial)
		docJSON, err := json.Marshal(doc)
		require.NoError(t, err, "Debe poder convertir documento a JSON")
		err = rdb.Set(ctx, cacheKey, docJSON, 24*time.Hour).Err()
		require.NoError(t, err, "Debe poder guardar documento en caché")

		// 3. Simular validación
		doc.Estado = "VALIDADO"
		docJSON, _ = json.Marshal(doc)
		err = rdb.Set(ctx, cacheKey, docJSON, 24*time.Hour).Err()
		require.NoError(t, err, "Debe poder actualizar estado a VALIDADO")

		// 4. Simular generación XML
		doc.Estado = "XML_GENERADO"
		docJSON, _ = json.Marshal(doc)
		err = rdb.Set(ctx, cacheKey, docJSON, 24*time.Hour).Err()
		require.NoError(t, err, "Debe poder actualizar estado a XML_GENERADO")

		// 5. Simular envío al SII
		doc.Estado = "ENVIADO_SII"
		doc.TrackID = "123456789"
		docJSON, _ = json.Marshal(doc)
		err = rdb.Set(ctx, cacheKey, docJSON, 24*time.Hour).Err()
		require.NoError(t, err, "Debe poder actualizar estado a ENVIADO_SII")

		// 6. Verificar estado final
		val, err := rdb.Get(ctx, cacheKey).Result()
		require.NoError(t, err, "Debe poder recuperar documento final")

		var docFinal Documento
		err = json.Unmarshal([]byte(val), &docFinal)
		require.NoError(t, err, "Debe poder convertir documento final")

		assert.Equal(t, "ENVIADO_SII", docFinal.Estado, "Estado final debe ser ENVIADO_SII")
		assert.Equal(t, "123456789", docFinal.TrackID, "Debe tener TrackID")
	})

	// Test de recuperación ante fallos
	t.Run("TestRecuperacionFallos", func(t *testing.T) {
		doc := crearDocumentoPrueba("33")
		cacheKey := "dte:" + doc.ID

		// 1. Simular fallo de Redis
		rdb.Close()

		// 2. Reconectar
		rdb = redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		})

		// 3. Verificar recuperación
		docJSON, err := json.Marshal(doc)
		require.NoError(t, err, "Debe poder convertir documento después de reconexión")

		err = rdb.Set(ctx, cacheKey, docJSON, 24*time.Hour).Err()
		require.NoError(t, err, "Debe poder guardar documento después de reconexión")
	})

	// Test de concurrencia en flujo DTE
	t.Run("TestConcurrenciaDTE", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			go func(i int) {
				doc := crearDocumentoPrueba("33")
				doc.ID = doc.ID + "-" + string(i)
				cacheKey := "dte:" + doc.ID

				docJSON, err := json.Marshal(doc)
				require.NoError(t, err, "Debe poder convertir documento concurrente")

				err = rdb.Set(ctx, cacheKey, docJSON, 24*time.Hour).Err()
				require.NoError(t, err, "Debe poder guardar documento concurrente")
			}(i)
		}

		// Esperar a que terminen las operaciones concurrentes
		time.Sleep(100 * time.Millisecond)
	})

	// Test de expiración de caché
	t.Run("TestExpiracionCache", func(t *testing.T) {
		doc := crearDocumentoPrueba("33")
		cacheKey := "dte:" + doc.ID

		// 1. Guardar con expiración corta
		docJSON, _ := json.Marshal(doc)
		err := rdb.Set(ctx, cacheKey, docJSON, 1*time.Second).Err()
		require.NoError(t, err, "Debe poder guardar con expiración")

		// 2. Verificar antes de expirar
		_, err = rdb.Get(ctx, cacheKey).Result()
		require.NoError(t, err, "Debe existir antes de expirar")

		// 3. Esperar expiración
		time.Sleep(2 * time.Second)

		// 4. Verificar después de expirar
		_, err = rdb.Get(ctx, cacheKey).Result()
		assert.Equal(t, redis.Nil, err, "Debe haber expirado")
	})
}
