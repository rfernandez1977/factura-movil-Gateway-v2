package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// Datos de prueba
var testData = []struct {
	name string
	calc CalculoTributario
}{
	{
		name: "Cálculo simple",
		calc: CalculoTributario{
			MontoNeto:   1000,
			MontoExento: 0,
			TasaIVA:     19,
		},
	},
	{
		name: "Cálculo con descuentos",
		calc: CalculoTributario{
			MontoNeto:   1000,
			MontoExento: 0,
			TasaIVA:     19,
			Descuentos: []Descuento{
				{TipoDescuento: "GLOBAL", Porcentaje: 10},
			},
		},
	},
	{
		name: "Cálculo complejo",
		calc: CalculoTributario{
			MontoNeto:           1000,
			MontoExento:         500,
			TasaIVA:             19,
			RetencionHonorarios: 10.75,
			RetencionILA:        10,
			Descuentos: []Descuento{
				{TipoDescuento: "GLOBAL", Porcentaje: 5},
				{TipoDescuento: "GLOBAL", Porcentaje: 3},
			},
			Recargos: []Recargo{
				{TipoRecargo: "GLOBAL", Porcentaje: 2},
			},
		},
	},
}

func setupRouter() (*gin.Engine, *TaxCalculationHandlers) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	handlers := &TaxCalculationHandlers{}

	router.POST("/calculate", handlers.CalculateHandler)
	router.POST("/calculate-optimized", handlers.OptimizedCalculationHandler)

	return router, handlers
}

func BenchmarkCalculationHandlers(b *testing.B) {
	router, _ := setupRouter()

	for _, tt := range testData {
		b.Run("Original/"+tt.name, func(b *testing.B) {
			jsonData, _ := json.Marshal(tt.calc)
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				w := httptest.NewRecorder()
				req, _ := http.NewRequest("POST", "/calculate", bytes.NewBuffer(jsonData))
				router.ServeHTTP(w, req)
			}
		})

		b.Run("Optimized/"+tt.name, func(b *testing.B) {
			jsonData, _ := json.Marshal(tt.calc)
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				w := httptest.NewRecorder()
				req, _ := http.NewRequest("POST", "/calculate-optimized", bytes.NewBuffer(jsonData))
				router.ServeHTTP(w, req)
			}
		})
	}
}

func TestCalculationResults(t *testing.T) {
	router, _ := setupRouter()

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			// Obtener resultado original
			jsonData, _ := json.Marshal(tt.calc)
			w1 := httptest.NewRecorder()
			req1, _ := http.NewRequest("POST", "/calculate", bytes.NewBuffer(jsonData))
			router.ServeHTTP(w1, req1)

			// Obtener resultado optimizado
			w2 := httptest.NewRecorder()
			req2, _ := http.NewRequest("POST", "/calculate-optimized", bytes.NewBuffer(jsonData))
			router.ServeHTTP(w2, req2)

			// Comparar resultados
			var result1, result2 map[string]interface{}
			json.Unmarshal(w1.Body.Bytes(), &result1)
			json.Unmarshal(w2.Body.Bytes(), &result2)

			// Verificar que los montos totales son iguales
			total1 := result1["desglose"].(map[string]interface{})["montoTotal"]
			total2 := result2["desglose"].(map[string]interface{})["montoTotal"]

			if total1 != total2 {
				t.Errorf("Los resultados no coinciden: original=%v, optimizado=%v", total1, total2)
			}
		})
	}
}

func TestValidationErrors(t *testing.T) {
	router, _ := setupRouter()
	tests := []struct {
		name    string
		calc    CalculoTributario
		wantErr string
	}{
		{
			name: "Monto negativo",
			calc: CalculoTributario{
				MontoNeto: -1000,
				TasaIVA:   19,
			},
			wantErr: "CALC_001",
		},
		{
			name: "Tasa IVA inválida",
			calc: CalculoTributario{
				MontoNeto: 1000,
				TasaIVA:   20,
			},
			wantErr: "CALC_002",
		},
		{
			name: "Retención honorarios inválida",
			calc: CalculoTributario{
				MontoNeto:           1000,
				TasaIVA:             19,
				RetencionHonorarios: 11,
			},
			wantErr: "CALC_003",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData, _ := json.Marshal(tt.calc)

			// Probar handler original
			w1 := httptest.NewRecorder()
			req1, _ := http.NewRequest("POST", "/calculate", bytes.NewBuffer(jsonData))
			router.ServeHTTP(w1, req1)

			// Probar handler optimizado
			w2 := httptest.NewRecorder()
			req2, _ := http.NewRequest("POST", "/calculate-optimized", bytes.NewBuffer(jsonData))
			router.ServeHTTP(w2, req2)

			// Verificar que ambos handlers retornan el mismo error
			var resp1, resp2 map[string]interface{}
			json.Unmarshal(w1.Body.Bytes(), &resp1)
			json.Unmarshal(w2.Body.Bytes(), &resp2)

			if resp1["codigo"] != tt.wantErr || resp2["codigo"] != tt.wantErr {
				t.Errorf("Error codes don't match: original=%v, optimized=%v, want=%v",
					resp1["codigo"], resp2["codigo"], tt.wantErr)
			}
		})
	}
}
