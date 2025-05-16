package test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/fmgo/core/firma/services"
	"github.com/fmgo/core/sii/client"
	"github.com/fmgo/core/sii/logger"
	"github.com/fmgo/core/sii/models"
	siiservices "github.com/fmgo/core/sii/services"
	"github.com/stretchr/testify/assert"
)

func setupTestEnvironment(t *testing.T) (*siiservices.IntegrationService, *models.Documento) {
	// Inicializar logger
	logger, err := logger.NewLogger("e2e_test", logger.DEBUG)
	assert.NoError(t, err)

	// Configurar ambiente de pruebas
	testConfig := &models.Config{
		SII: models.SIIConfig{
			BaseURL:    os.Getenv("SII_URL"),
			CertPath:   os.Getenv("CERT_PATH"),
			KeyPath:    os.Getenv("KEY_PATH"),
			RetryCount: 3,
			Timeout:    30,
		},
	}

	// Inicializar cliente SII
	siiClient, err := client.NewHTTPClient(testConfig.SII.CertPath, "test", models.Certificacion, client.DefaultRetryConfig())
	assert.NoError(t, err)

	// Inicializar servicio de firma
	firmaService := services.NewFirmaService(
		services.NewCertificadoRepository(),
		services.NewCacheService(),
		logger,
	)

	// Inicializar servicio de integración
	integrationService := siiservices.NewIntegrationService(siiClient, firmaService, logger)

	// Crear documento de prueba
	doc := &models.Documento{
		ID:           "E2E_TEST_001",
		TipoDTE:      "33",
		Folio:        1,
		RutEmisor:    os.Getenv("RUT_EMISOR"),
		RutReceptor:  os.Getenv("RUT_RECEPTOR"),
		MontoTotal:   1000,
		FechaEmision: time.Now(),
	}

	return integrationService, doc
}

func TestE2EFlujoCompleto(t *testing.T) {
	if testing.Short() {
		t.Skip("Saltando prueba E2E en modo corto")
	}

	// Configurar ambiente
	service, doc := setupTestEnvironment(t)
	ctx := context.Background()

	// 1. Enviar documento
	respEnvio, err := service.EnviarDocumento(ctx, doc)
	assert.NoError(t, err)
	assert.NotEmpty(t, respEnvio.TrackID)

	// Esperar un momento para que el SII procese el documento
	time.Sleep(5 * time.Second)

	// 2. Consultar estado
	estado, err := service.ConsultarEstadoEnvio(ctx, respEnvio.TrackID)
	assert.NoError(t, err)
	assert.NotNil(t, estado)

	// 3. Validar documento
	validacion, err := service.ValidarDocumento(ctx, doc)
	assert.NoError(t, err)
	assert.NotNil(t, validacion)
	assert.Equal(t, doc.Folio, validacion.Folio)
}

func TestE2EErrorHandling(t *testing.T) {
	if testing.Short() {
		t.Skip("Saltando prueba E2E en modo corto")
	}

	// Configurar ambiente
	service, doc := setupTestEnvironment(t)
	ctx := context.Background()

	// Intentar enviar documento con datos inválidos
	doc.RutEmisor = "INVALID_RUT"
	_, err := service.EnviarDocumento(ctx, doc)
	assert.Error(t, err)

	// Intentar consultar estado con trackID inválido
	_, err = service.ConsultarEstadoEnvio(ctx, "INVALID_TRACK_ID")
	assert.Error(t, err)

	// Intentar validar documento con folio inválido
	doc.Folio = -1
	_, err = service.ValidarDocumento(ctx, doc)
	assert.Error(t, err)
}
