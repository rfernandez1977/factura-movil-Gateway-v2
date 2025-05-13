package integration

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/cursor/FMgo/config"
	"github.com/cursor/FMgo/models"
	"github.com/cursor/FMgo/services"
	"github.com/cursor/FMgo/services/firma"
	"github.com/cursor/FMgo/services/sii"
	"github.com/cursor/FMgo/tests/mocks"
	"github.com/cursor/FMgo/utils"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestIntegracionFlujoCompletoDTE prueba el flujo completo de un DTE
func TestIntegracionFlujoCompletoDTE(t *testing.T) {
	// Configuración inicial
	config := config.LoadConfig()

	// Crear mock del repositorio
	mockRepo := mocks.NewMockRepository()

	// Crear mock del cliente Redis
	mockRedis := mocks.NewMockRedisClient()

	// Crear mock del servicio SII
	mockSII := mocks.NewMockSIIService()

	// Crear servicio de firma digital
	firmaService, err := firma.NewService(config)
	assert.NoError(t, err)

	// Crear servicio de DTE
	dteService, err := services.NewDTEService(&config.Supabase, nil)
	assert.NoError(t, err)
	assert.NotNil(t, dteService)

	// Crear generador de DTE
	dteGenerator := utils.NewDTEGenerator()

	// Datos del emisor
	emisor := utils.Emisor{
		RUT:         "76.123.456-7",
		RazonSocial: "EMPRESA DE PRUEBA",
		Giro:        "DESARROLLO DE SOFTWARE",
		Acteco:      "620100",
		Direccion:   "CALLE PRINCIPAL 123",
		Comuna:      "SANTIAGO",
		Ciudad:      "SANTIAGO",
	}

	// Datos del receptor
	receptor := utils.Receptor{
		RUT:         "56.789.012-3",
		RazonSocial: "CLIENTE DE PRUEBA",
		Giro:        "COMERCIO",
		Direccion:   "AVENIDA SECUNDARIA 456",
		Comuna:      "PROVIDENCIA",
	}

	// Detalles del documento
	detalles := []utils.DetalleDTE{
		{
			NroLinDet:  1,
			NombreItem: "Producto de prueba",
			Cantidad:   1,
			PrecioUnit: 1000,
			MontoItem:  1000,
		},
	}

	// Configurar mocks del repositorio
	mockRepo.On("GetControlFolio", "33").Return(&models.ControlFolio{
		TipoDocumento: "33",
		FolioActual:   1,
		RangoFinal:    100,
	}, nil)

	mockRepo.On("SaveDocumentoTributario", mock.Anything).Return(nil)
	mockRepo.On("UpdateDocumentoTributario", mock.Anything).Return(nil)
	mockRepo.On("SaveEstadoDocumento", mock.Anything).Return(nil)

	// Configurar mocks del cliente Redis
	mockRedis.On("Get", mock.Anything, mock.Anything).Return(redis.NewStringResult("", nil))
	mockRedis.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(redis.NewStatusResult("OK", nil))

	// Generar documento
	documento := dteGenerator.GenerarDTE("33", emisor, receptor, detalles)

	// Crear sobre para envío
	sobre := &models.SobreDTEModel{
		Caratula: models.CaratulaXMLModel{
			Version:      "1.0",
			RutEmisor:    emisor.RUT,
			RutEnvia:     emisor.RUT,
			RutReceptor:  receptor.RUT,
			FchResol:     "2024-01-01",
			NroResol:     "0",
			TmstFirmaEnv: time.Now().Format("2006-01-02T15:04:05"),
		},
		Documentos: []models.DTEXMLModel{
			{
				DocumentoXML: *documento,
			},
		},
	}

	// Configurar mocks del servicio SII
	mockSII.On("EnviarDTE", mock.Anything).Return(&models.RespuestaSII{
		Estado:  "OK",
		Glosa:   "Documento Recibido",
		TrackID: "123456",
	}, nil)

	mockSII.On("ConsultarEstado", "123456").Return(&sii.EstadoSII{
		Estado: "OK",
		Glosa:  "Documento Aceptado",
	}, nil)

	// Firmar sobre
	err = firmaService.FirmarSobre(sobre)
	assert.NoError(t, err)
	assert.NotEmpty(t, sobre.Signature)

	// Firmar cada DTE y generar TED
	for i := range sobre.Documentos {
		err = firmaService.FirmarDTE(&sobre.Documentos[i])
		assert.NoError(t, err)
		assert.NotEmpty(t, sobre.Documentos[i].Signature)

		ted, err := firmaService.GenerarTED(&sobre.Documentos[i])
		assert.NoError(t, err)
		assert.NotEmpty(t, ted)
		// Asignar TED al documento
		if doc, ok := sobre.Documentos[i].DocumentoXML.(map[string]interface{}); ok {
			doc["TED"] = ted
		}
	}

	// Enviar a SII
	respuesta, err := mockSII.EnviarDTE(&sobre.Documentos[0])
	assert.NoError(t, err)
	assert.Equal(t, "OK", respuesta.Estado)
	assert.Equal(t, "123456", respuesta.TrackID)

	// Consultar estado
	estado, err := mockSII.ConsultarEstado(respuesta.TrackID)
	assert.NoError(t, err)
	assert.Equal(t, "OK", estado.Estado)

	// Verificar que se llamaron los métodos del repositorio
	mockRepo.AssertExpectations(t)
	mockRedis.AssertExpectations(t)
	mockSII.AssertExpectations(t)
}

// TestIntegracionRepositorioDocumento prueba la integración entre el repositorio y el modelo de documento
func TestIntegracionRepositorioDocumento(t *testing.T) {
	// Crear el mock del repositorio
	mockRepo := new(MockRepository)

	// Crear un documento de prueba
	tiempo := time.Now()
	doc := models.DocumentoTributario{
		ID:            "5f50cf13c56e0a1d9b4fbe5a",
		TipoDocumento: models.TipoFacturaElectronica,
		Folio:         1,
		FechaEmision:  tiempo,
		MontoTotal:    10000,
		Estado:        models.EstadoEnviado,
	}

	// Configurar expectativas del mock
	mockRepo.On("SaveDocumentoTributario", mock.Anything, &doc).Return(nil)
	mockRepo.On("GetDocumentoTributario", mock.Anything, "5f50cf13c56e0a1d9b4fbe5a").Return(&doc, nil)

	// Guardar el documento
	err := mockRepo.SaveDocumentoTributario(context.Background(), &doc)
	// Verificar que no hay error
	assert.NoError(t, err)

	// Recuperar el documento
	docRecuperado, err := mockRepo.GetDocumentoTributario(context.Background(), "5f50cf13c56e0a1d9b4fbe5a")

	// Verificar resultados
	assert.NoError(t, err)
	assert.NotNil(t, docRecuperado)
	assert.Equal(t, doc.ID, docRecuperado.ID)
	assert.Equal(t, doc.TipoDocumento, docRecuperado.TipoDocumento)
	assert.Equal(t, doc.Folio, docRecuperado.Folio)
	assert.Equal(t, doc.MontoTotal, docRecuperado.MontoTotal)

	// Verificar que se llamaron los métodos esperados
	mockRepo.AssertExpectations(t)
}

// TestIntegracionControlFolio prueba la integración entre el repositorio y el control de folios
func TestIntegracionControlFolio(t *testing.T) {
	// Crear el mock del repositorio
	mockRepo := new(MockRepository)

	// Crear un control de folio de prueba
	tiempo := time.Now()
	control := models.ControlFolio{
		TipoDocumento:     "33",
		RangoInicial:      1,
		RangoFinal:        100,
		FolioActual:       5,
		FoliosDisponibles: 95,
		UltimoUso:         tiempo,
		EstadoCAF:         "ACTIVO",
		AlertaGenerada:    false,
	}

	// Configurar expectativas
	mockRepo.On("GetControlFolio", mock.Anything, 33).Return(&control, nil)

	// Actualizar el control con un nuevo folio
	controlActualizado := control
	controlActualizado.FolioActual = 6
	controlActualizado.FoliosDisponibles = 94
	controlActualizado.UltimoUso = time.Now()

	mockRepo.On("UpdateControlFolio", mock.Anything, &controlActualizado).Return(nil)

	// Obtener el control de folio
	controlObtenido, err := mockRepo.GetControlFolio(context.Background(), 33)

	// Verificar resultados de obtención
	assert.NoError(t, err)
	assert.NotNil(t, controlObtenido)
	assert.Equal(t, control.TipoDocumento, controlObtenido.TipoDocumento)
	assert.Equal(t, control.FolioActual, controlObtenido.FolioActual)

	// Actualizar el control de folio
	err = mockRepo.UpdateControlFolio(context.Background(), &controlActualizado)

	// Verificar resultados de actualización
	assert.NoError(t, err)

	// Verificar que se llamaron los métodos esperados
	mockRepo.AssertExpectations(t)
}

// TestIntegracionTodosTiposDocumentos prueba la integración con todos los tipos de documentos
func TestIntegracionTodosTiposDocumentos(t *testing.T) {
	// Configuración inicial
	config := config.LoadConfig()

	// Crear mock del repositorio
	mockRepo := mocks.NewMockRepository()

	// Crear mock del cliente Redis
	mockRedis := mocks.NewMockRedisClient()

	// Crear mock del servicio SII
	mockSII := mocks.NewMockSIIService()

	// Crear servicio de firma digital
	firmaService, err := firma.NewService(config)
	assert.NoError(t, err)

	// Crear servicio de DTE
	dteService, err := services.NewDTEService(&config.Supabase, nil)
	assert.NoError(t, err)
	assert.NotNil(t, dteService)

	// Crear generador de DTE
	dteGenerator := utils.NewDTEGenerator()

	// Datos del emisor
	emisor := utils.Emisor{
		RUT:         "76.123.456-7",
		RazonSocial: "EMPRESA DE PRUEBA",
		Giro:        "DESARROLLO DE SOFTWARE",
		Acteco:      "620100",
		Direccion:   "CALLE PRINCIPAL 123",
		Comuna:      "SANTIAGO",
		Ciudad:      "SANTIAGO",
	}

	// Datos del receptor
	receptor := utils.Receptor{
		RUT:         "56.789.012-3",
		RazonSocial: "CLIENTE DE PRUEBA",
		Giro:        "COMERCIO",
		Direccion:   "AVENIDA SECUNDARIA 456",
		Comuna:      "PROVIDENCIA",
	}

	// Detalles del documento
	detalles := []utils.DetalleDTE{
		{
			NroLinDet:  1,
			NombreItem: "Producto de prueba",
			Cantidad:   1,
			PrecioUnit: 1000,
			MontoItem:  1000,
		},
	}

	// Tipos de documentos a probar
	tiposDocumento := []string{"33", "52", "61"}

	for _, tipo := range tiposDocumento {
		t.Run(fmt.Sprintf("Tipo %s", tipo), func(t *testing.T) {
			// Configurar mocks del repositorio
			mockRepo.On("GetControlFolio", tipo).Return(&models.ControlFolio{
				TipoDocumento: tipo,
				FolioActual:   1,
				RangoFinal:    100,
			}, nil)

			mockRepo.On("SaveDocumentoTributario", mock.Anything).Return(nil)
			mockRepo.On("UpdateDocumentoTributario", mock.Anything).Return(nil)
			mockRepo.On("SaveEstadoDocumento", mock.Anything).Return(nil)

			// Configurar mocks del cliente Redis
			mockRedis.On("Get", mock.Anything, mock.Anything).Return(redis.NewStringResult("", nil))
			mockRedis.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(redis.NewStatusResult("OK", nil))

			// Generar documento
			documento := dteGenerator.GenerarDTE(tipo, emisor, receptor, detalles)

			// Crear sobre para envío
			sobre := &models.SobreDTEModel{
				Caratula: models.CaratulaXMLModel{
					Version:      "1.0",
					RutEmisor:    emisor.RUT,
					RutEnvia:     emisor.RUT,
					RutReceptor:  receptor.RUT,
					FchResol:     "2024-01-01",
					NroResol:     "0",
					TmstFirmaEnv: time.Now().Format("2006-01-02T15:04:05"),
				},
				Documentos: []models.DTEXMLModel{
					{
						DocumentoXML: *documento,
					},
				},
			}

			// Configurar mocks del servicio SII
			mockSII.On("EnviarDTE", mock.Anything).Return(&models.RespuestaSII{
				Estado:  "OK",
				Glosa:   "Documento Recibido",
				TrackID: "123456",
			}, nil)

			mockSII.On("ConsultarEstado", "123456").Return(&sii.EstadoSII{
				Estado: "OK",
				Glosa:  "Documento Aceptado",
			}, nil)

			// Firmar sobre
			err = firmaService.FirmarSobre(sobre)
			assert.NoError(t, err)
			assert.NotEmpty(t, sobre.Signature)

			// Firmar cada DTE y generar TED
			for i := range sobre.Documentos {
				err = firmaService.FirmarDTE(&sobre.Documentos[i])
				assert.NoError(t, err)
				assert.NotEmpty(t, sobre.Documentos[i].Signature)

				ted, err := firmaService.GenerarTED(&sobre.Documentos[i])
				assert.NoError(t, err)
				assert.NotEmpty(t, ted)
				// Asignar TED al documento
				if doc, ok := sobre.Documentos[i].DocumentoXML.(map[string]interface{}); ok {
					doc["TED"] = ted
				}
			}

			// Enviar a SII
			respuesta, err := mockSII.EnviarDTE(&sobre.Documentos[0])
			assert.NoError(t, err)
			assert.Equal(t, "OK", respuesta.Estado)
			assert.Equal(t, "123456", respuesta.TrackID)

			// Consultar estado
			estado, err := mockSII.ConsultarEstado(respuesta.TrackID)
			assert.NoError(t, err)
			assert.Equal(t, "OK", estado.Estado)

			// Verificar que se llamaron los métodos del repositorio
			mockRepo.AssertExpectations(t)
			mockRedis.AssertExpectations(t)
			mockSII.AssertExpectations(t)
		})
	}
}

// TestIntegracionCAF prueba la integración con el servicio CAF
func TestIntegracionCAF(t *testing.T) {
	// Configuración inicial
	config := config.LoadConfig()

	// Crear mock del cliente Redis
	mockRedis := &MockRedisClient{}

	// Crear mock del servicio SII
	mockSII := &MockSIIService{}

	// Crear servicio CAF
	cafService := services.NewCAFService(&config.Supabase, mockRedis, mockSII, "cert.pem", "key.pem", "CERTIFICACION", "https://api.sii.cl")

	// Configurar mocks
	mockRedis.On("Get", mock.Anything, mock.Anything).Return(redis.NewStringResult("", nil))
	mockRedis.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(redis.NewStatusResult("OK", nil))

	mockSII.On("EnviarDTE", mock.Anything).Return(&sii.RespuestaSII{
		Estado:  "OK",
		Glosa:   "Documento Recibido",
		TrackID: "123456",
	}, nil)

	mockSII.On("ConsultarEstado", "123456").Return(&sii.EstadoSII{
		Estado: "OK",
		Glosa:  "Documento Aceptado",
	}, nil)

	// Solicitar CAF
	req := &services.CAFRequest{
		RUTEmisor:      "76.123.456-7",
		TipoDTE:        "33",
		FolioInicial:   1,
		FolioFinal:     100,
		FechaSolicitud: time.Now(),
	}

	resp, err := cafService.SolicitarCAF(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, "OK", resp.Estado)
	assert.Equal(t, "123456", resp.TrackID)

	// Consultar estado
	estado, err := cafService.ConsultarEstadoCAF(context.Background(), resp.TrackID)
	assert.NoError(t, err)
	assert.Equal(t, "OK", estado.Estado)

	// Verificar que se llamaron los métodos
	mockRedis.AssertExpectations(t)
	mockSII.AssertExpectations(t)
}
