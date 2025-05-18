package middleware

import (
	"context"
	"testing"

	"FMgo/core/sii/models"
	"github.com/stretchr/testify/assert"
)

func TestValidatorMiddleware(t *testing.T) {
	middleware := NewValidatorMiddleware("../../../schema_dte")
	ctx := context.Background()

	t.Run("ValidateDTE", func(t *testing.T) {
		// Caso válido
		dteCorrecto := &models.DTE{
			Documento: models.Documento{
				Encabezado: models.Encabezado{
					IdDoc: models.IdDoc{
						TipoDTE: "33",
						Folio:   1,
					},
					Emisor: models.Emisor{
						RUTEmisor:  "76.123.456-7",
						RznSoc:     "Empresa de Prueba",
						GiroEmis:   "Servicios Informáticos",
						Acteco:     "722000",
						DirOrigen:  "Calle Principal 123",
						CmnaOrigen: "Santiago",
					},
					Receptor: models.Receptor{
						RUTRecep:    "77.888.999-0",
						RznSocRecep: "Cliente de Prueba",
						GiroRecep:   "Comercio",
						DirRecep:    "Av. Cliente 456",
						CmnaRecep:   "Providencia",
					},
					Totales: models.Totales{
						MntNeto:  1000000,
						TasaIVA:  19.0,
						IVA:      190000,
						MntTotal: 1190000,
					},
				},
				Detalle: []models.Detalle{
					{
						NroLinDet: 1,
						NmbItem:   "Servicio de Desarrollo",
						QtyItem:   1.0,
						PrcItem:   1000000.0,
						MontoItem: 1000000,
					},
				},
			},
		}

		err := middleware.ValidateDTE(ctx, dteCorrecto)
		assert.NoError(t, err, "DTE válido debería pasar la validación")

		// Caso inválido
		dteIncorrecto := &models.DTE{}
		err = middleware.ValidateDTE(ctx, dteIncorrecto)
		assert.Error(t, err, "DTE inválido debería fallar la validación")
	})

	t.Run("ValidateEnvioDTE", func(t *testing.T) {
		// Caso válido
		envioDTECorrecto := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<EnvioDTE version="1.0">
    <SetDTE ID="SetDoc">
        <Caratula version="1.0">
            <RutEmisor>76123456-7</RutEmisor>
            <RutEnvia>12345678-9</RutEnvia>
            <RutReceptor>77888999-0</RutReceptor>
            <FchResol>2023-01-01</FchResol>
            <NroResol>0</NroResol>
            <TmstFirmaEnv>2024-03-20T12:00:00</TmstFirmaEnv>
        </Caratula>
    </SetDTE>
</EnvioDTE>`)

		err := middleware.ValidateEnvioDTE(ctx, envioDTECorrecto)
		assert.NoError(t, err, "EnvioDTE válido debería pasar la validación")

		// Caso inválido
		envioDTEIncorrecto := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<EnvioDTE version="1.0">
    <SetDTE>
        <!-- Falta ID y otros campos requeridos -->
    </SetDTE>
</EnvioDTE>`)

		err = middleware.ValidateEnvioDTE(ctx, envioDTEIncorrecto)
		assert.Error(t, err, "EnvioDTE inválido debería fallar la validación")
	})

	t.Run("WithDTEValidation", func(t *testing.T) {
		handlerLlamado := false
		handler := func(ctx context.Context, dte *models.DTE) error {
			handlerLlamado = true
			return nil
		}

		wrappedHandler := middleware.WithDTEValidation(handler)

		// Probar con DTE válido
		dteCorrecto := &models.DTE{
			Documento: models.Documento{
				Encabezado: models.Encabezado{
					IdDoc: models.IdDoc{
						TipoDTE: "33",
						Folio:   1,
					},
				},
			},
		}

		err := wrappedHandler(ctx, dteCorrecto)
		assert.NoError(t, err, "Handler con DTE válido debería ejecutarse sin error")
		assert.True(t, handlerLlamado, "Handler debería haber sido llamado")

		// Probar con DTE inválido
		handlerLlamado = false
		dteIncorrecto := &models.DTE{}

		err = wrappedHandler(ctx, dteIncorrecto)
		assert.Error(t, err, "Handler con DTE inválido debería retornar error")
		assert.False(t, handlerLlamado, "Handler no debería haber sido llamado")
	})

	t.Run("WithEnvioDTEValidation", func(t *testing.T) {
		handlerLlamado := false
		handler := func(ctx context.Context, envioDTE []byte) error {
			handlerLlamado = true
			return nil
		}

		wrappedHandler := middleware.WithEnvioDTEValidation(handler)

		// Probar con EnvioDTE válido
		envioDTECorrecto := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<EnvioDTE version="1.0">
    <SetDTE ID="SetDoc">
        <Caratula version="1.0">
            <RutEmisor>76123456-7</RutEmisor>
            <RutEnvia>12345678-9</RutEnvia>
            <RutReceptor>77888999-0</RutReceptor>
            <FchResol>2023-01-01</FchResol>
            <NroResol>0</NroResol>
            <TmstFirmaEnv>2024-03-20T12:00:00</TmstFirmaEnv>
        </Caratula>
    </SetDTE>
</EnvioDTE>`)

		err := wrappedHandler(ctx, envioDTECorrecto)
		assert.NoError(t, err, "Handler con EnvioDTE válido debería ejecutarse sin error")
		assert.True(t, handlerLlamado, "Handler debería haber sido llamado")

		// Probar con EnvioDTE inválido
		handlerLlamado = false
		envioDTEIncorrecto := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<EnvioDTE version="1.0">
    <SetDTE>
        <!-- XML inválido -->
    </SetDTE>
</EnvioDTE>`)

		err = wrappedHandler(ctx, envioDTEIncorrecto)
		assert.Error(t, err, "Handler con EnvioDTE inválido debería retornar error")
		assert.False(t, handlerLlamado, "Handler no debería haber sido llamado")
	})
}
