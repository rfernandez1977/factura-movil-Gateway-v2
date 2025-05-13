package caf

import (
	"context"
	"testing"

	"FMgo/config"
	"FMgo/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRedisClient es un mock del cliente Redis
type MockRedisClient struct {
	mock.Mock
}

// MockSIIService es un mock del servicio SII
type MockSIIService struct {
	mock.Mock
}

func TestCAFService(t *testing.T) {
	// Configurar servicio
	config := &config.Config{
		CertPath: "testdata/cert.pem",
		KeyPath:  "testdata/key.pem",
	}
	redisClient := new(MockRedisClient)
	siiService := new(MockSIIService)
	service := NewService(config, redisClient, siiService)

	t.Run("GetCAFDisponible", func(t *testing.T) {
		caf, err := service.GetCAFDisponible(context.Background(), models.TipoFacturaElectronica, "76.123.456-7")
		assert.NoError(t, err)
		assert.NotNil(t, caf)
		assert.NotEmpty(t, caf.Version)
		assert.NotEmpty(t, caf.DA.RE.Rut)
		assert.NotEmpty(t, caf.DA.RSAPK.M)
		assert.NotEmpty(t, caf.DA.RSAPK.E)
		assert.NotEmpty(t, caf.FRMA.Value)
	})

	t.Run("ValidarCAF", func(t *testing.T) {
		caf := &models.CAFDTEXML{
			Version: "1.0",
			DA: models.DAXMLModel{
				RE: models.RutXMLModel{
					Rut:         "76.123.456-7",
					Dv:          "7",
					RazonSocial: "Empresa Test",
					Giro:        "Servicios",
					Acteco:      "123456",
					Direccion:   "Calle Test 123",
					Comuna:      "Santiago",
					Ciudad:      "Santiago",
				},
				RSAPK: models.RSAPKXMLModel{
					M: "test_modulus",
					E: "test_exponent",
				},
				IDK: 1,
			},
			FRMA: models.FRMAXMLModel{
				Algoritmo: "SHA256withRSA",
				Value:     "test_signature",
			},
		}

		err := service.ValidarCAF(caf)
		assert.NoError(t, err)
	})

	t.Run("ValidarCAF_InvalidData", func(t *testing.T) {
		caf := &models.CAFDTEXML{
			Version: "",
			DA: models.DAXMLModel{
				RE: models.RutXMLModel{
					Rut: "",
				},
			},
		}

		err := service.ValidarCAF(caf)
		assert.Error(t, err)
	})
}
