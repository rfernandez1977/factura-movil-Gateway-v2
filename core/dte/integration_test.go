package dte

import (
	"testing"
	"time"

	"github.com/fmgo/core/dte/generator"
	"github.com/fmgo/core/dte/parser"
	"github.com/fmgo/core/dte/types"
	"github.com/fmgo/core/dte/validator"
	"github.com/stretchr/testify/assert"
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
			RUT:         "76212889-6",
			RazonSocial: "Empresa Test",
			Giro:        "Servicios",
			Direccion:   "Calle Test 123",
			Comuna:      "Santiago",
			Ciudad:      "Santiago",
		},
		Receptor: types.Receptor{
			RUT:         "13195458-1",
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

	// Test 1: Generar DTE
	t.Run("Generar DTE", func(t *testing.T) {
		dte, err := gen.GenerateDTE(input)
		assert.NoError(t, err)
		assert.NotNil(t, dte)
		assert.Equal(t, "33", dte.Documento.Encabezado.IDDocumento.TipoDTE)
		assert.Equal(t, 1, dte.Documento.Encabezado.IDDocumento.Folio)
		assert.Equal(t, "CREADO", dte.Estado)
		assert.False(t, dte.Firmado)
	})

	// Test 2: Validar DTE
	t.Run("Validar DTE", func(t *testing.T) {
		dte, err := gen.GenerateDTE(input)
		assert.NoError(t, err)

		err = validator.ValidateDTE(dte)
		assert.NoError(t, err)
	})

	// Test 3: Generar y Parsear XML
	t.Run("Generar y Parsear XML", func(t *testing.T) {
		// Generar DTE
		dte, err := gen.GenerateDTE(input)
		assert.NoError(t, err)

		// Generar XML
		xmlData, err := parse.GenerateXML(dte)
		assert.NoError(t, err)
		assert.NotEmpty(t, xmlData)

		// Parsear XML de vuelta a DTE
		parsedDTE, err := parse.ParseXML(xmlData)
		assert.NoError(t, err)
		assert.NotNil(t, parsedDTE)

		// Validar el DTE parseado
		err = validator.ValidateDTE(parsedDTE)
		assert.NoError(t, err)
	})

	// Test 4: Validar cálculos de totales
	t.Run("Validar cálculos", func(t *testing.T) {
		dte, err := gen.GenerateDTE(input)
		assert.NoError(t, err)

		totales := dte.Documento.Encabezado.Totales
		assert.Equal(t, 2000.0, totales.MontoNeto)
		assert.Equal(t, 19.0, totales.TasaIVA)
		assert.Equal(t, 380.0, totales.IVA)
		assert.Equal(t, 2380.0, totales.MontoTotal)
	})
}
