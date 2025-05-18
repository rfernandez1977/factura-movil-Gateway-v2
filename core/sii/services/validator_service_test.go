package services

import (
	"path/filepath"
	"testing"

	"FMgo/core/sii/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidatorService(t *testing.T) {
	// Configurar servicio de validación
	validator := NewValidatorService()
	schemaDir := filepath.Join("..", "..", "..", "schema_dte")
	err := validator.CargarEsquemasBase(schemaDir)
	require.NoError(t, err, "Error cargando esquemas base")

	t.Run("ValidarDTE_Valido", func(t *testing.T) {
		dte := &models.DTE{
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

		err := validator.ValidarDTE(dte)
		assert.NoError(t, err, "DTE válido debería pasar la validación")
	})

	t.Run("ValidarDTE_Invalido", func(t *testing.T) {
		dte := &models.DTE{
			Documento: models.Documento{
				// DTE inválido sin campos requeridos
			},
		}

		err := validator.ValidarDTE(dte)
		assert.Error(t, err, "DTE inválido debería fallar la validación")
	})

	t.Run("ValidarEnvioDTE", func(t *testing.T) {
		envioDTEXML := []byte(`<?xml version="1.0" encoding="UTF-8"?>
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

		err := validator.ValidarEnvioDTE(envioDTEXML)
		assert.NoError(t, err, "EnvioDTE válido debería pasar la validación")
	})

	t.Run("ValidarXMLContraEsquema", func(t *testing.T) {
		xmlData := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<DTE version="1.0">
    <Documento>
        <Encabezado>
            <IdDoc>
                <TipoDTE>33</TipoDTE>
                <Folio>1</Folio>
            </IdDoc>
        </Encabezado>
    </Documento>
</DTE>`)

		err := validator.ValidarXMLContraEsquema(xmlData, "DTE")
		assert.NoError(t, err, "XML válido debería pasar la validación")
	})
}

func TestCargarEsquema(t *testing.T) {
	validator := NewValidatorService()

	t.Run("CargarEsquema_Existente", func(t *testing.T) {
		err := validator.CargarEsquema("DTE", filepath.Join("..", "..", "..", "schema_dte", "DTE_v10.xsd"))
		assert.NoError(t, err, "Debería cargar un esquema existente")
	})

	t.Run("CargarEsquema_NoExistente", func(t *testing.T) {
		err := validator.CargarEsquema("NoExiste", "archivo_no_existente.xsd")
		assert.Error(t, err, "Debería fallar al cargar un esquema no existente")
	})
}
