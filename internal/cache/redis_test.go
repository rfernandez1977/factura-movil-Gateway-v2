package cache

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	ID    string `json:"id"`
	Value int    `json:"value"`
}

func getTestConfig() Config {
	return Config{
		Host:     "localhost",
		Port:     6379,
		Password: "", // Sin contraseña para pruebas locales
		DB:       0,
	}
}

func TestRedisCache(t *testing.T) {
	cfg := getTestConfig()
	cache, err := NewRedisCache(cfg)
	if err != nil {
		t.Skipf("Redis no disponible: %v", err)
		return
	}
	defer cache.Close()

	t.Run("set_and_get", func(t *testing.T) {
		testData := TestStruct{
			ID:    "test1",
			Value: 123,
		}

		// Almacenar en caché
		err := cache.Set("test_key", testData, 0)
		assert.NoError(t, err)

		// Recuperar del caché
		var result TestStruct
		err = cache.Get("test_key", &result)
		assert.NoError(t, err)
		assert.Equal(t, testData, result)
	})

	t.Run("expiration", func(t *testing.T) {
		testData := TestStruct{
			ID:    "test2",
			Value: 456,
		}

		// Almacenar con expiración corta
		err := cache.Set("test_exp", testData, 1*time.Second)
		assert.NoError(t, err)

		// Esperar a que expire
		time.Sleep(2 * time.Second)

		// Intentar recuperar
		var result TestStruct
		err = cache.Get("test_exp", &result)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "clave no encontrada")
	})

	t.Run("delete", func(t *testing.T) {
		testData := TestStruct{
			ID:    "test3",
			Value: 789,
		}

		// Almacenar y eliminar
		key := "test_del"
		err := cache.Set(key, testData, 0)
		assert.NoError(t, err)

		err = cache.Delete(key)
		assert.NoError(t, err)

		// Intentar recuperar después de eliminar
		var result TestStruct
		err = cache.Get(key, &result)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "clave no encontrada")
	})

	t.Run("clear", func(t *testing.T) {
		// Almacenar múltiples valores
		for i := 0; i < 3; i++ {
			testData := TestStruct{
				ID:    "test_clear",
				Value: i,
			}
			err := cache.Set(fmt.Sprintf("clear_test_%d", i), testData, 0)
			assert.NoError(t, err)
		}

		// Limpiar caché
		err := cache.Clear()
		assert.NoError(t, err)

		// Verificar que las claves fueron eliminadas
		for i := 0; i < 3; i++ {
			var result TestStruct
			err := cache.Get(fmt.Sprintf("clear_test_%d", i), &result)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "clave no encontrada")
		}
	})
}
