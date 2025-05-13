package firma

import (
	"testing"

	"github.com/cursor/FMgo/config"
	"github.com/cursor/FMgo/models"

	"github.com/stretchr/testify/assert"
)

func TestFirmaService(t *testing.T) {
	// Configurar servicio
	config := &config.Config{
		CertPath: "testdata/cert.pem",
		KeyPath:  "testdata/key.pem",
	}
	service, err := NewService(config)
	assert.NoError(t, err)
	assert.NotNil(t, service)

	t.Run("FirmarDTE", func(t *testing.T) {
		dte := &models.DTEXMLModel{
			DocumentoXML: models.DocumentoXMLModel{
				Encabezado: models.EncabezadoXMLModel{
					IdDoc: models.IdDocXMLModel{
						TipoDTE:     models.TipoFacturaElectronica,
						Folio:       1,
						FchEmis:     "2024-03-20",
						IndServicio: 1,
						IndMntNeto:  1,
					},
					Emisor: models.EmisorXMLModel{
						RUTEmisor: "76.123.456-7",
					},
					Receptor: models.ReceptorXMLModel{
						RUTRecep: "56.789.012-3",
					},
					Totales: models.TotalesXMLModel{
						MntNeto:  84034,
						TasaIVA:  19,
						IVA:      15966,
						MntTotal: 100000,
					},
				},
			},
		}

		err := service.FirmarDTE(dte)
		assert.NoError(t, err)
		assert.NotEmpty(t, dte.Signature)
	})

	t.Run("GenerarTED", func(t *testing.T) {
		dte := &models.DTEXMLModel{
			DocumentoXML: models.DocumentoXMLModel{
				Encabezado: models.EncabezadoXMLModel{
					IdDoc: models.IdDocXMLModel{
						TipoDTE:     models.TipoFacturaElectronica,
						Folio:       1,
						FchEmis:     "2024-03-20",
						IndServicio: 1,
						IndMntNeto:  1,
					},
					Emisor: models.EmisorXMLModel{
						RUTEmisor: "76.123.456-7",
					},
					Receptor: models.ReceptorXMLModel{
						RUTRecep: "56.789.012-3",
					},
					Totales: models.TotalesXMLModel{
						MntNeto:  84034,
						TasaIVA:  19,
						IVA:      15966,
						MntTotal: 100000,
					},
				},
			},
		}

		ted, err := service.GenerarTED(dte)
		assert.NoError(t, err)
		assert.NotEmpty(t, ted)
	})

	t.Run("FirmarSobre", func(t *testing.T) {
		dte := &models.DTEXMLModel{
			DocumentoXML: models.DocumentoXMLModel{
				Encabezado: models.EncabezadoXMLModel{
					IdDoc: models.IdDocXMLModel{
						TipoDTE:     models.TipoFacturaElectronica,
						Folio:       1,
						FchEmis:     "2024-03-20",
						IndServicio: 1,
						IndMntNeto:  1,
					},
					Emisor: models.EmisorXMLModel{
						RUTEmisor: "76.123.456-7",
					},
					Receptor: models.ReceptorXMLModel{
						RUTRecep: "56.789.012-3",
					},
					Totales: models.TotalesXMLModel{
						MntNeto:  84034,
						TasaIVA:  19,
						IVA:      15966,
						MntTotal: 100000,
					},
				},
			},
		}

		sobre := &models.SobreDTEModel{
			Caratula: models.CaratulaXMLModel{
				Version:     "1.0",
				RutEmisor:   "76.123.456-7",
				RutEnvia:    "76.123.456-7",
				RutReceptor: "56.789.012-3",
			},
			Documentos: []models.DTEXMLModel{*dte},
		}

		err := service.FirmarSobre(sobre)
		assert.NoError(t, err)
		assert.NotEmpty(t, sobre.Signature)
	})

	t.Run("InvalidData", func(t *testing.T) {
		// Test con DTE nulo
		err := service.FirmarDTE(nil)
		assert.Error(t, err)

		// Test con DTE sin datos requeridos
		dte := &models.DTEXMLModel{}
		err = service.FirmarDTE(dte)
		assert.Error(t, err)

		// Test con sobre nulo
		err = service.FirmarSobre(nil)
		assert.Error(t, err)

		// Test con sobre sin documentos
		sobre := &models.SobreDTEModel{}
		err = service.FirmarSobre(sobre)
		assert.Error(t, err)
	})
}
