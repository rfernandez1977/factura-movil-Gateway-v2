package services

import (
	"testing"
	"time"

	"FMgo/models"
)

func TestValidatorService(t *testing.T) {
	schemaPath := "../schema_dte"
	validator := NewValidatorService(schemaPath)

	t.Run("Cargar Esquema", func(t *testing.T) {
		err := validator.CargarEsquema("DTE_v10.xsd")
		if err != nil {
			t.Errorf("Error al cargar esquema: %v", err)
		}
	})

	t.Run("Cargar Esquema Inexistente", func(t *testing.T) {
		err := validator.CargarEsquema("NoExiste.xsd")
		if err == nil {
			t.Error("Se esperaba un error al cargar esquema inexistente")
		}
	})

	t.Run("Validar Documento", func(t *testing.T) {
		doc := &models.DTEDocument{
			Version: "1.0",
			Documento: models.DTEDocumento{
				ID: "T33F1",
				Encabezado: models.DTEEncabezado{
					IdDoc: models.DTEIdDoc{
						TipoDTE: "33",
						Folio:   1,
						FchEmis: time.Now(),
					},
					Emisor: models.DTEEmisor{
						RUTEmisor: "76.123.456-7",
						RznSoc:    "Empresa de Prueba",
					},
					Receptor: models.DTEReceptor{
						RUTRecep:    "77.888.999-0",
						RznSocRecep: "Cliente de Prueba",
					},
					Totales: models.DTETotales{
						MntNeto:  100000,
						TasaIVA:  19,
						IVA:      19000,
						MntTotal: 119000,
					},
				},
			},
		}

		err := validator.ValidarDocumento(doc)
		if err != nil {
			t.Errorf("Error al validar documento válido: %v", err)
		}
	})

	t.Run("Thread Safety", func(t *testing.T) {
		done := make(chan bool)
		for i := 0; i < 10; i++ {
			go func() {
				err := validator.CargarEsquema("DTE_v10.xsd")
				if err != nil {
					t.Errorf("Error en concurrencia: %v", err)
				}
				done <- true
			}()
		}

		for i := 0; i < 10; i++ {
			<-done
		}
	})

	t.Run("Limpiar Esquemas", func(t *testing.T) {
		validator.LimpiarEsquemas()
		// Intentar usar un esquema después de limpiar
		err := validator.ValidarXML([]byte("<DTE/>"), "DTE_v10.xsd")
		if err == nil {
			t.Error("Se esperaba un error al usar esquema después de limpiar")
		}
	})
}
