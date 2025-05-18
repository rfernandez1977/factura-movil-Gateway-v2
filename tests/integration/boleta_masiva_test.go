package integration

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"FMgo/models"
	"FMgo/services"

	"github.com/stretchr/testify/suite"
)

type BoletaMasivaTestSuite struct {
	suite.Suite
	ctx          context.Context
	cancel       context.CancelFunc
	siiService   *services.SIIService
	boletaRepo   services.BoletaRepository
	boletaService *services.BoletaService
}

func (s *BoletaMasivaTestSuite) SetupSuite() {
	s.ctx, s.cancel = context.WithTimeout(context.Background(), 10*time.Minute)
	s.boletaRepo = &services.BoletaRepositoryMock{}
	s.siiService = services.NewSIIService()
	s.boletaService = services.NewBoletaService(s.siiService, s.boletaRepo)
}

func (s *BoletaMasivaTestSuite) TearDownSuite() {
	s.cancel()
}

func (s *BoletaMasivaTestSuite) TestEnvioMasivoBoletas() {
	// Test de envío masivo de boletas (hasta 500 según esquema XSD)
	s.Run("Envío Masivo de Boletas", func() {
		numBoletas := 100 // Probar con 100 boletas
		boletas := make([]*models.Boleta, numBoletas)

		// 1. Preparar datos de prueba
		for i := 0; i < numBoletas; i++ {
			boletas[i] = &models.Boleta{
				RUTEmisor:           "76555555-5",
				RazonSocialEmisor:   "EMPRESA DE PRUEBA SPA",
				GiroEmisor:          "DESARROLLO DE SOFTWARE",
				DireccionEmisor:     "CALLE PRUEBA 123",
				ComunaEmisor:        "SANTIAGO",
				Folio:               int64(i + 1),
				FechaEmision:        time.Now(),
				MontoNeto:           10000,
				MontoIVA:            1900,
				MontoTotal:          11900,
				Items: []models.Item{
					{
						Descripcion:    fmt.Sprintf("Producto de Prueba %d", i+1),
						Cantidad:       1,
						Precio:         10000,
						Total:          10000,
					},
				},
			}
		}

		// 2. Validar CAF para el rango de folios
		caf, err := s.siiService.ValidarCAF("39", 1, int64(numBoletas))
		s.Require().NoError(err, "Debe obtener un CAF válido")
		s.Require().NotNil(caf, "CAF no debe ser nil")

		// 3. Generar y validar XML para cada boleta
		var wg sync.WaitGroup
		xmlChan := make(chan string, numBoletas)
		errChan := make(chan error, numBoletas)

		for _, boleta := range boletas {
			wg.Add(1)
			go func(b *models.Boleta) {
				defer wg.Done()
				xml, err := s.siiService.GenerarXMLBoleta(b)
				if err != nil {
					errChan <- fmt.Errorf("error generando XML para folio %d: %v", b.Folio, err)
					return
				}
				if err := s.siiService.ValidarXMLBoleta(xml); err != nil {
					errChan <- fmt.Errorf("error validando XML para folio %d: %v", b.Folio, err)
					return
				}
				xmlChan <- xml
			}(boleta)
		}

		// Esperar generación de XMLs
		wg.Wait()
		close(xmlChan)
		close(errChan)

		// Verificar errores
		for err := range errChan {
			s.Require().NoError(err, "No debe haber errores en la generación/validación de XMLs")
		}

		// 4. Crear envío masivo
		xmls := make([]string, 0, numBoletas)
		for xml := range xmlChan {
			xmls = append(xmls, xml)
		}

		envioMasivo, err := s.siiService.CrearEnvioMasivoBoletas(xmls, "76555555-5")
		s.Require().NoError(err, "Debe crear el envío masivo")
		s.Require().NotEmpty(envioMasivo, "El envío masivo no debe estar vacío")

		// 5. Firmar envío masivo
		envioFirmado, err := s.siiService.FirmarEnvioMasivo(envioMasivo)
		s.Require().NoError(err, "Debe firmar el envío masivo")
		s.Require().NotEmpty(envioFirmado, "El envío firmado no debe estar vacío")

		// 6. Enviar al SII
		trackID, err := s.siiService.EnviarAlSII(envioFirmado)
		s.Require().NoError(err, "Debe enviar al SII")
		s.Require().NotEmpty(trackID, "Debe obtener trackID")

		// 7. Verificar estado del envío
		estado, err := s.siiService.ConsultarEstadoEnvio(trackID)
		s.Require().NoError(err, "Debe consultar estado")
		s.Require().Equal("RECIBIDO", estado, "El envío debe ser recibido")

		// 8. Almacenar resultados
		for i, boleta := range boletas {
			err := s.boletaRepo.Crear(boleta)
			s.Require().NoError(err, "Debe almacenar boleta %d", i+1)
		}
	})
}

func (s *BoletaMasivaTestSuite) TestEnvioMasivoBoletasErrores() {
	// Test de manejo de errores en envío masivo
	s.Run("Manejo de Errores en Envío Masivo", func() {
		// 1. Probar límite máximo de boletas (más de 500)
		numBoletas := 501
		boletas := make([]*models.Boleta, numBoletas)
		for i := 0; i < numBoletas; i++ {
			boletas[i] = &models.Boleta{
				Folio: int64(i + 1),
				// ... otros campos necesarios
			}
		}

		xmls := make([]string, numBoletas)
		_, err := s.siiService.CrearEnvioMasivoBoletas(xmls, "76555555-5")
		s.Require().Error(err, "Debe fallar al exceder el límite de 500 boletas")

		// 2. Probar envío con XML inválido
		xmlsInvalidos := []string{"<xml>inválido</xml>"}
		_, err = s.siiService.CrearEnvioMasivoBoletas(xmlsInvalidos, "76555555-5")
		s.Require().Error(err, "Debe fallar con XML inválido")

		// 3. Probar envío sin boletas
		_, err = s.siiService.CrearEnvioMasivoBoletas([]string{}, "76555555-5")
		s.Require().Error(err, "Debe fallar sin boletas")

		// 4. Probar RUT emisor inválido
		_, err = s.siiService.CrearEnvioMasivoBoletas([]string{"<xml>válido</xml>"}, "RUT-INVALIDO")
		s.Require().Error(err, "Debe fallar con RUT inválido")
	})
}

func (s *BoletaMasivaTestSuite) TestEnvioMasivoBoletasConcurrencia() {
	// Test de envío masivo concurrente
	s.Run("Envío Masivo Concurrente", func() {
		numEnvios := 5
		numBoletasPorEnvio := 100
		var wg sync.WaitGroup
		errChan := make(chan error, numEnvios)

		for i := 0; i < numEnvios; i++ {
			wg.Add(1)
			go func(envioID int) {
				defer wg.Done()

				// Crear boletas para este envío
				boletas := make([]*models.Boleta, numBoletasPorEnvio)
				for j := 0; j < numBoletasPorEnvio; j++ {
					boletas[j] = &models.Boleta{
						Folio: int64(envioID*numBoletasPorEnvio + j + 1),
						// ... otros campos necesarios
					}
				}

				// Procesar envío
				xmls := make([]string, numBoletasPorEnvio)
				envioMasivo, err := s.siiService.CrearEnvioMasivoBoletas(xmls, "76555555-5")
				if err != nil {
					errChan <- fmt.Errorf("error en envío %d: %v", envioID, err)
					return
				}

				// Enviar al SII
				_, err = s.siiService.EnviarAlSII(envioMasivo)
				if err != nil {
					errChan <- fmt.Errorf("error enviando al SII envío %d: %v", envioID, err)
					return
				}
			}(i)
		}

		// Esperar todos los envíos
		wg.Wait()
		close(errChan)

		// Verificar errores
		for err := range errChan {
			s.Require().NoError(err, "No debe haber errores en los envíos concurrentes")
		}
	})
}

func TestBoletaMasivaSuite(t *testing.T) {
	suite.Run(t, new(BoletaMasivaTestSuite))
} 