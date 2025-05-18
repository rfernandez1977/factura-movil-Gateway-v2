package integration

import (
	"testing"
	"time"

	"FMgo/core/dte/generator"
	"FMgo/core/dte/parser"
	"FMgo/core/dte/types"
)

func TestDTEIntegration(t *testing.T) {
	// Crear instancias de los componentes
	gen := generator.NewDefaultGenerator()
	parse := parser.NewXMLParser()

	// Datos de prueba
	input := &generator.GenerateInput{
		TipoDTE:      "33",
		Folio:        1,
		FechaEmision: time.Now(),
		Emisor: types.Emisor{
			RUT:         "76.543.210-9",
			RazonSocial: "Empresa Test",
			Giro:        "Servicios",
			Direccion:   "Calle Test 123",
			Comuna:      "Santiago",
			Ciudad:      "Santiago",
		},
		Receptor: types.Receptor{
			RUT:         "77.654.321-0",
			RazonSocial: "Cliente Test",
			Giro:        "Comercio",
			Direccion:   "Av Test 456",
			Comuna:      "Santiago",
			Ciudad:      "Santiago",
		},
		Detalles: []types.Detalle{
			{
				NumeroLinea: 1,
				Nombre:      "Producto Test",
				Cantidad:    2,
				Precio:      1000,
				MontoItem:   2000,
			},
		},
	}

	// ... existing code ...
}
