package utils

import (
	"math"
	"testing"
)

func TestAmountValidator(t *testing.T) {
	validator := NewAmountValidator()

	// Test ValidateAmount
	t.Run("ValidateAmount", func(t *testing.T) {
		tests := []struct {
			name      string
			amount    float64
			fieldName string
			wantErr   bool
		}{
			{"valid amount", 100.50, "monto", false},
			{"negative amount", -100.50, "monto", true},
			{"zero amount", 0, "monto", false},
			{"max amount", 999999999.99, "monto", false},
			{"exceed max amount", 1000000000.00, "monto", true},
			{"infinity", math.Inf(1), "monto", true},
			{"NaN", math.NaN(), "monto", true},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := validator.ValidateAmount(tt.amount, tt.fieldName)
				if (err != nil) != tt.wantErr {
					t.Errorf("ValidateAmount() error = %v, wantErr %v", err, tt.wantErr)
				}
			})
		}
	})

	// Test ValidatePercentage
	t.Run("ValidatePercentage", func(t *testing.T) {
		tests := []struct {
			name       string
			percentage float64
			fieldName  string
			wantErr    bool
		}{
			{"valid percentage", 50.0, "porcentaje", false},
			{"negative percentage", -50.0, "porcentaje", true},
			{"zero percentage", 0, "porcentaje", false},
			{"max percentage", 100.0, "porcentaje", false},
			{"exceed max percentage", 101.0, "porcentaje", true},
			{"infinity", math.Inf(1), "porcentaje", true},
			{"NaN", math.NaN(), "porcentaje", true},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := validator.ValidatePercentage(tt.percentage, tt.fieldName)
				if (err != nil) != tt.wantErr {
					t.Errorf("ValidatePercentage() error = %v, wantErr %v", err, tt.wantErr)
				}
			})
		}
	})

	// Test ValidateDiscount
	t.Run("ValidateDiscount", func(t *testing.T) {
		tests := []struct {
			name      string
			amount    float64
			discount  float64
			fieldName string
			wantErr   bool
		}{
			{"valid discount", 100.0, 10.0, "descuento", false},
			{"discount equals amount", 100.0, 100.0, "descuento", false},
			{"discount exceeds amount", 100.0, 101.0, "descuento", true},
			{"negative discount", 100.0, -10.0, "descuento", true},
			{"zero discount", 100.0, 0, "descuento", false},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := validator.ValidateDiscount(tt.amount, tt.discount, tt.fieldName)
				if (err != nil) != tt.wantErr {
					t.Errorf("ValidateDiscount() error = %v, wantErr %v", err, tt.wantErr)
				}
			})
		}
	})

	// Test CalculateIVA
	t.Run("CalculateIVA", func(t *testing.T) {
		tests := []struct {
			name          string
			montoNeto     float64
			porcentajeIVA float64
			want          float64
		}{
			{"valid IVA", 100.0, 19.0, 19.0},
			{"zero amount", 0, 19.0, 0},
			{"zero percentage", 100.0, 0, 0},
			{"negative amount", -100.0, 19.0, 0},
			{"negative percentage", 100.0, -19.0, 0},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got := validator.CalculateIVA(tt.montoNeto, tt.porcentajeIVA)
				if got != tt.want {
					t.Errorf("CalculateIVA() = %v, want %v", got, tt.want)
				}
			})
		}
	})

	// Test CalculateDiscount
	t.Run("CalculateDiscount", func(t *testing.T) {
		tests := []struct {
			name       string
			amount     float64
			percentage float64
			want       float64
		}{
			{"valid discount", 100.0, 10.0, 10.0},
			{"zero amount", 0, 10.0, 0},
			{"zero percentage", 100.0, 0, 0},
			{"negative amount", -100.0, 10.0, 0},
			{"negative percentage", 100.0, -10.0, 0},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got := validator.CalculateDiscount(tt.amount, tt.percentage)
				if got != tt.want {
					t.Errorf("CalculateDiscount() = %v, want %v", got, tt.want)
				}
			})
		}
	})

	// Test FormatAmount
	t.Run("FormatAmount", func(t *testing.T) {
		tests := []struct {
			name   string
			amount float64
			want   string
		}{
			{"valid amount", 100.50, "100.50"},
			{"zero amount", 0, "0.00"},
			{"negative amount", -100.50, "0.00"},
			{"infinity", math.Inf(1), "0.00"},
			{"NaN", math.NaN(), "0.00"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got := validator.FormatAmount(tt.amount)
				if got != tt.want {
					t.Errorf("FormatAmount() = %v, want %v", got, tt.want)
				}
			})
		}
	})

	// Test ValidateAmountsConsistency
	t.Run("ValidateAmountsConsistency", func(t *testing.T) {
		tests := []struct {
			name                 string
			montoNeto            float64
			montoExento          float64
			montoIVA             float64
			montoTotal           float64
			impuestosAdicionales []float64
			wantErr              bool
		}{
			{
				"valid amounts",
				100.0,
				50.0,
				19.0,
				169.0,
				[]float64{},
				false,
			},
			{
				"with additional taxes",
				100.0,
				50.0,
				19.0,
				179.0,
				[]float64{10.0},
				false,
			},
			{
				"inconsistent amounts",
				100.0,
				50.0,
				19.0,
				170.0,
				[]float64{},
				true,
			},
			{
				"negative amount",
				-100.0,
				50.0,
				19.0,
				169.0,
				[]float64{},
				true,
			},
			{
				"infinity",
				math.Inf(1),
				50.0,
				19.0,
				169.0,
				[]float64{},
				true,
			},
			{
				"NaN",
				math.NaN(),
				50.0,
				19.0,
				169.0,
				[]float64{},
				true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := validator.ValidateAmountsConsistency(tt.montoNeto, tt.montoExento, tt.montoIVA, tt.montoTotal, tt.impuestosAdicionales...)
				if (err != nil) != tt.wantErr {
					t.Errorf("ValidateAmountsConsistency() error = %v, wantErr %v", err, tt.wantErr)
				}
			})
		}
	})

	// Test isValidDecimal
	t.Run("isValidDecimal", func(t *testing.T) {
		tests := []struct {
			name  string
			value float64
			want  bool
		}{
			{"valid decimal", 100.50, true},
			{"no decimal", 100.0, true},
			{"too many decimals", 100.555, false},
			{"negative", -100.50, true},
			{"zero", 0.0, true},
			{"infinity", math.Inf(1), false},
			{"NaN", math.NaN(), false},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got := validator.isValidDecimal(tt.value)
				if got != tt.want {
					t.Errorf("isValidDecimal() = %v, want %v", got, tt.want)
				}
			})
		}
	})

	// Test ValidateDecimal
	t.Run("ValidateDecimal", func(t *testing.T) {
		tests := []struct {
			name      string
			number    float64
			decimals  int
			fieldName string
			wantErr   bool
		}{
			{"valid decimal", 100.50, 2, "monto", false},
			{"too many decimals", 100.555, 2, "monto", true},
			{"invalid decimals", 100.50, -1, "monto", true},
			{"exceed max decimals", 100.50, 3, "monto", true},
			{"infinity", math.Inf(1), 2, "monto", true},
			{"NaN", math.NaN(), 2, "monto", true},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := validator.ValidateDecimal(tt.number, tt.decimals, tt.fieldName)
				if (err != nil) != tt.wantErr {
					t.Errorf("ValidateDecimal() error = %v, wantErr %v", err, tt.wantErr)
				}
			})
		}
	})

	// Test ValidateQuantity
	t.Run("ValidateQuantity", func(t *testing.T) {
		tests := []struct {
			name     string
			quantity float64
			wantErr  bool
		}{
			{"valid quantity", 100.50, false},
			{"zero quantity", 0.0, true},
			{"negative quantity", -100.50, true},
			{"max quantity", 999999.999, false},
			{"exceed max quantity", 1000000.0, true},
			{"too many decimals", 100.5555, true},
			{"infinity", math.Inf(1), true},
			{"NaN", math.NaN(), true},
			{"valid decimal places", 100.123, false},
			{"valid decimal places with trailing zeros", 100.100, false},
			{"valid decimal places with leading zeros", 100.001, false},
			{"valid integer", 100.0, false},
			{"valid large number with decimals", 999999.999, false},
			{"valid small number with decimals", 0.001, false},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := validator.ValidateQuantity(tt.quantity)
				if (err != nil) != tt.wantErr {
					t.Errorf("ValidateQuantity() error = %v, wantErr %v", err, tt.wantErr)
				}
			})
		}
	})

	// Test ValidateQuantityWithUnit
	t.Run("ValidateQuantityWithUnit", func(t *testing.T) {
		tests := []struct {
			name     string
			quantity float64
			unit     string
			wantErr  bool
		}{
			{"valid kg", 100.50, "kg", false},
			{"exceed max kg", 1001.0, "kg", true},
			{"valid units", 100.0, "unidades", false},
			{"invalid units decimal", 100.5, "unidades", true},
			{"valid liters", 500.0, "litros", false},
			{"exceed max liters", 1001.0, "litros", true},
			{"invalid unit", 100.0, "invalid", false},
			{"valid kg with decimals", 999.999, "kg", false},
			{"valid liters with decimals", 999.999, "litros", false},
			{"zero quantity", 0.0, "kg", true},
			{"negative quantity", -100.0, "kg", true},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := validator.ValidateQuantityWithUnit(tt.quantity, tt.unit)
				if (err != nil) != tt.wantErr {
					t.Errorf("ValidateQuantityWithUnit() error = %v, wantErr %v", err, tt.wantErr)
				}
			})
		}
	})

	// Test ValidateUnitPrice
	t.Run("ValidateUnitPrice", func(t *testing.T) {
		tests := []struct {
			name    string
			price   float64
			wantErr bool
		}{
			{"valid price", 100.50, false},
			{"zero price", 0.0, true},
			{"negative price", -100.50, true},
			{"max price", 999999999.99, false},
			{"exceed max price", 1000000000.0, true},
			{"too many decimals", 100.555, true},
			{"infinity", math.Inf(1), true},
			{"NaN", math.NaN(), true},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := validator.ValidateUnitPrice(tt.price)
				if (err != nil) != tt.wantErr {
					t.Errorf("ValidateUnitPrice() error = %v, wantErr %v", err, tt.wantErr)
				}
			})
		}
	})
}
