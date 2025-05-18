package test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"FMgo/core/caf/services"
)

// mockRedisCache simula un servicio de caché Redis
type mockRedisCache struct {
	data map[string]interface{}
}

func newMockRedisCache() *mockRedisCache {
	return &mockRedisCache{
		data: make(map[string]interface{}),
	}
}

func (m *mockRedisCache) Get(_ context.Context, key string) (interface{}, error) {
	if val, ok := m.data[key]; ok {
		return val, nil
	}
	return nil, fmt.Errorf("clave no encontrada: %s", key)
}

func (m *mockRedisCache) Set(_ context.Context, key string, value interface{}, _ time.Duration) error {
	m.data[key] = value
	return nil
}

// mockProductionLogger implementa un logger real para pruebas
type mockProductionLogger struct{}

func (m *mockProductionLogger) Info(msg string, args ...interface{}) {
	fmt.Printf("[INFO] "+msg+"\n", args...)
}
func (m *mockProductionLogger) Error(msg string, args ...interface{}) {
	fmt.Printf("[ERROR] "+msg+"\n", args...)
}
func (m *mockProductionLogger) Warn(msg string, args ...interface{}) {
	fmt.Printf("[WARN] "+msg+"\n", args...)
}
func (m *mockProductionLogger) Debug(msg string, args ...interface{}) {
	fmt.Printf("[DEBUG] "+msg+"\n", args...)
}

func TestIntegracionValidadorCAF(t *testing.T) {
	// Configurar servicios
	cache := newMockRedisCache()
	logger := &mockProductionLogger{}
	validador := services.NewValidadorCAF(cache, logger)

	// Cargar CAF de prueba
	cafPath := filepath.Join("testdata", "caf_test.xml")
	xmlCAF, err := os.ReadFile(cafPath)
	if err != nil {
		t.Fatalf("Error leyendo archivo CAF: %v", err)
	}

	// Probar flujo completo
	t.Run("Flujo completo de validación", func(t *testing.T) {
		ctx := context.Background()

		// 1. Validar CAF
		resultado, err := validador.ValidarCAF(ctx, xmlCAF)
		if err != nil {
			t.Fatalf("Error validando CAF: %v", err)
		}
		if !resultado.Valido {
			t.Errorf("CAF debería ser válido: %v", resultado.Error)
		}

		// 2. Verificar almacenamiento en caché
		tipoDTE := 33 // Asumiendo que es el tipo en el XML de prueba
		desde, hasta, err := validador.ObtenerRangoFolios(ctx, tipoDTE)
		if err != nil {
			t.Errorf("Error obteniendo rango de folios: %v", err)
		}

		// 3. Validar folios
		foliosAPruebaValidos := []int{desde, desde + 1, hasta - 1, hasta}
		foliosAPruebaInvalidos := []int{desde - 1, hasta + 1}

		for _, folio := range foliosAPruebaValidos {
			valido, err := validador.ValidarFolio(ctx, folio, tipoDTE)
			if err != nil {
				t.Errorf("Error validando folio %d: %v", folio, err)
			}
			if !valido {
				t.Errorf("Folio %d debería ser válido", folio)
			}
		}

		for _, folio := range foliosAPruebaInvalidos {
			valido, err := validador.ValidarFolio(ctx, folio, tipoDTE)
			if err != nil {
				continue // Esperamos error para folios inválidos
			}
			if valido {
				t.Errorf("Folio %d no debería ser válido", folio)
			}
		}
	})
}

func TestIntegracionConcurrencia(t *testing.T) {
	cache := newMockRedisCache()
	logger := &mockProductionLogger{}
	validador := services.NewValidadorCAF(cache, logger)

	// Cargar CAF de prueba
	cafPath := filepath.Join("testdata", "caf_test.xml")
	xmlCAF, err := os.ReadFile(cafPath)
	if err != nil {
		t.Fatalf("Error leyendo archivo CAF: %v", err)
	}

	// Probar validación concurrente
	t.Run("Validación concurrente", func(t *testing.T) {
		const numGoroutines = 10
		errChan := make(chan error, numGoroutines)
		ctx := context.Background()

		for i := 0; i < numGoroutines; i++ {
			go func() {
				_, err := validador.ValidarCAF(ctx, xmlCAF)
				errChan <- err
			}()
		}

		// Esperar resultados
		for i := 0; i < numGoroutines; i++ {
			if err := <-errChan; err != nil {
				t.Errorf("Error en goroutine %d: %v", i, err)
			}
		}
	})
}

func TestIntegracionRecuperacion(t *testing.T) {
	cache := newMockRedisCache()
	logger := &mockProductionLogger{}
	validador := services.NewValidadorCAF(cache, logger)

	ctx := context.Background()

	t.Run("Recuperación tras fallo de caché", func(t *testing.T) {
		// 1. Validar CAF inicial
		cafPath := filepath.Join("testdata", "caf_test.xml")
		xmlCAF, err := os.ReadFile(cafPath)
		if err != nil {
			t.Fatalf("Error leyendo archivo CAF: %v", err)
		}

		resultado, err := validador.ValidarCAF(ctx, xmlCAF)
		if err != nil {
			t.Fatalf("Error en validación inicial: %v", err)
		}

		// 2. Simular fallo de caché limpiando datos
		cache.data = make(map[string]interface{})

		// 3. Intentar validar folio
		tipoDTE := 33
		_, err = validador.ValidarFolio(ctx, 1, tipoDTE)
		if err == nil {
			t.Error("Se esperaba error por caché vacía")
		}

		// 4. Recargar CAF y verificar recuperación
		resultado, err = validador.ValidarCAF(ctx, xmlCAF)
		if err != nil {
			t.Fatalf("Error en recarga de CAF: %v", err)
		}
		if !resultado.Valido {
			t.Error("CAF debería ser válido después de recarga")
		}

		// 5. Verificar que ahora sí funciona la validación de folio
		valido, err := validador.ValidarFolio(ctx, 1, tipoDTE)
		if err != nil {
			t.Errorf("Error validando folio después de recarga: %v", err)
		}
		if !valido {
			t.Error("Folio debería ser válido después de recarga")
		}
	})
}
