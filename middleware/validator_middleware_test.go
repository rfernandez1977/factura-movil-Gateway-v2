package middleware

import (
	"context"
	"testing"
	"time"

	"github.com/fmgo/models"
)

func TestValidatorMiddleware(t *testing.T) {
	schemaPath := "../schema_dte"
	middleware := NewValidatorMiddleware(schemaPath)

	t.Run("Inicialización Lazy", func(t *testing.T) {
		err := middleware.initValidator()
		if err != nil {
			t.Errorf("Error en inicialización lazy: %v", err)
		}
	})

	t.Run("Validar Documento", func(t *testing.T) {
		doc := &models.DocumentoTributarioBasico{
			ID:           "T33F1",
			TipoDTE:      "33",
			Folio:        1,
			FechaEmision: time.Now(),
			RutEmisor:    "76.123.456-7",
			RutReceptor:  "77.888.999-0",
			MontoTotal:   119000,
			MontoNeto:    100000,
			MontoIVA:     19000,
			Estado:       string(models.EstadoDocumentoPendiente),
		}

		handler := middleware.WithDocumentoValidation(func(ctx context.Context, doc *models.DocumentoTributarioBasico) error {
			return nil
		})

		err := handler(context.Background(), doc)
		if err != nil {
			t.Errorf("Error al validar documento en middleware: %v", err)
		}
	})

	t.Run("Validar Documento Inválido", func(t *testing.T) {
		doc := &models.DocumentoTributarioBasico{
			ID:      "T33F1",
			TipoDTE: "99", // Tipo inválido
		}

		handler := middleware.WithDocumentoValidation(func(ctx context.Context, doc *models.DocumentoTributarioBasico) error {
			return nil
		})

		err := handler(context.Background(), doc)
		if err == nil {
			t.Error("Se esperaba un error al validar documento inválido")
		}
	})

	t.Run("Validar Envío", func(t *testing.T) {
		envio := struct {
			Version string `xml:"version,attr"`
			ID      string `xml:"ID"`
		}{
			Version: "1.0",
			ID:      "SetDTE001",
		}

		handler := middleware.WithEnvioValidation(func(ctx context.Context, envio interface{}) error {
			return nil
		})

		err := handler(context.Background(), envio)
		if err != nil {
			t.Errorf("Error al validar envío en middleware: %v", err)
		}
	})
}
