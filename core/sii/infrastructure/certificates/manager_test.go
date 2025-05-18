package certificates

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewCertificateManager(t *testing.T) {
	t.Run("error al crear con ruta inválida", func(t *testing.T) {
		manager, err := NewCertificateManager("ruta/invalida.p12", "password")
		assert.Error(t, err)
		assert.Nil(t, manager)
	})
}

func TestValidateCertificate(t *testing.T) {
	t.Run("certificado no cargado", func(t *testing.T) {
		manager := &CertificateManager{}
		err := manager.ValidateCertificate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "certificado no cargado")
	})

	t.Run("certificado expirado", func(t *testing.T) {
		manager := &CertificateManager{
			certInfo: &CertificateInfo{
				ValidFrom:  time.Now().AddDate(-1, 0, 0),
				ValidUntil: time.Now().AddDate(0, 0, -1),
			},
		}
		err := manager.ValidateCertificate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "certificado expirado")
	})

	t.Run("certificado aún no válido", func(t *testing.T) {
		manager := &CertificateManager{
			certInfo: &CertificateInfo{
				ValidFrom:  time.Now().AddDate(0, 0, 1),
				ValidUntil: time.Now().AddDate(1, 0, 0),
			},
		}
		err := manager.ValidateCertificate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "certificado aún no es válido")
	})

	t.Run("certificado válido", func(t *testing.T) {
		manager := &CertificateManager{
			certInfo: &CertificateInfo{
				ValidFrom:  time.Now().AddDate(0, 0, -1),
				ValidUntil: time.Now().AddDate(0, 0, 1),
			},
		}
		err := manager.ValidateCertificate()
		assert.NoError(t, err)
	})
}

func TestIsExpiringSoon(t *testing.T) {
	t.Run("certificado no cargado", func(t *testing.T) {
		manager := &CertificateManager{}
		assert.False(t, manager.IsExpiringSoon(30))
	})

	t.Run("certificado expirando pronto", func(t *testing.T) {
		manager := &CertificateManager{
			certInfo: &CertificateInfo{
				ValidUntil: time.Now().AddDate(0, 0, 15),
			},
		}
		assert.True(t, manager.IsExpiringSoon(30))
	})

	t.Run("certificado no expira pronto", func(t *testing.T) {
		manager := &CertificateManager{
			certInfo: &CertificateInfo{
				ValidUntil: time.Now().AddDate(0, 0, 45),
			},
		}
		assert.False(t, manager.IsExpiringSoon(30))
	})
}

func TestGetCertificateInfo(t *testing.T) {
	t.Run("obtener info de certificado", func(t *testing.T) {
		info := &CertificateInfo{
			Subject:      "CN=Test",
			Issuer:       "CN=CA Test",
			SerialNumber: "123456",
			RutTitular:   "12345678-9",
		}
		manager := &CertificateManager{
			certInfo: info,
		}

		result := manager.GetCertificateInfo()
		assert.Equal(t, info, result)
	})
}
