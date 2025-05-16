package validator

import (
	"testing"
	"time"

	"github.com/fmgo/core/dte/types"
	"github.com/stretchr/testify/assert"
)

func TestValidateDTE(t *testing.T) {
	t.Run("DTE válido", func(t *testing.T) {
		dte := &types.DTE{
			ID: "TEST-001",
			Documento: types.Documento{
				Encabezado: types.Encabezado{
					IDDocumento: types.IDDocumento{
						TipoDTE:      "33",
						Folio:        1,
						FechaEmision: time.Now(),
					},
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
					Totales: types.Totales{
						MontoNeto:  1000,
						TasaIVA:    19,
						IVA:        190,
						MontoTotal: 1190,
					},
				},
				Detalles: []types.Detalle{
					{
						NumeroLinea: 1,
						Nombre:      "Producto Test",
						Cantidad:    1,
						Precio:      1000,
						MontoItem:   1000,
					},
				},
			},
			Estado:  "CREADO",
			Firmado: false,
		}

		err := ValidateDTE(dte)
		assert.NoError(t, err)
	})

	t.Run("DTE sin ID", func(t *testing.T) {
		dte := &types.DTE{
			Documento: types.Documento{
				Encabezado: types.Encabezado{
					IDDocumento: types.IDDocumento{
						TipoDTE:      "33",
						Folio:        1,
						FechaEmision: time.Now(),
					},
				},
			},
		}

		err := ValidateDTE(dte)
		assert.Error(t, err)
	})

	t.Run("DTE con XML firmado pero no marcado como firmado", func(t *testing.T) {
		dte := &types.DTE{
			ID:         "TEST-001",
			XMLFirmado: "<xml>test</xml>",
			Firmado:    false,
		}

		err := ValidateDTE(dte)
		assert.Error(t, err)
	})
}

func TestValidateDocumento(t *testing.T) {
	t.Run("Documento válido", func(t *testing.T) {
		doc := types.Documento{
			Encabezado: types.Encabezado{
				IDDocumento: types.IDDocumento{
					TipoDTE:      "33",
					Folio:        1,
					FechaEmision: time.Now(),
				},
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
				Totales: types.Totales{
					MontoNeto:  1000,
					TasaIVA:    19,
					IVA:        190,
					MontoTotal: 1190,
				},
			},
			Detalles: []types.Detalle{
				{
					NumeroLinea: 1,
					Nombre:      "Producto Test",
					Cantidad:    1,
					Precio:      1000,
					MontoItem:   1000,
				},
			},
		}

		err := ValidateDocumento(&doc)
		assert.NoError(t, err)
	})

	t.Run("Documento sin detalles", func(t *testing.T) {
		doc := types.Documento{
			Encabezado: types.Encabezado{
				IDDocumento: types.IDDocumento{
					TipoDTE:      "33",
					Folio:        1,
					FechaEmision: time.Now(),
				},
			},
		}

		err := ValidateDocumento(&doc)
		assert.Error(t, err)
	})
}
