package firma

import (
	"encoding/xml"
	"errors"
	"fmt"
	"time"
)

// SIIFirmaService implementa la interfaz SIIFirmaService
type SIIFirmaService struct {
	*BaseFirmaService
	certCache CertCache
}

// NewSIIFirmaService crea una nueva instancia de SIIFirmaService
func NewSIIFirmaService(config *ConfiguracionFirma) (*SIIFirmaService, error) {
	if config.RutEmpresa == "" {
		return nil, errors.New("el RUT de la empresa es requerido")
	}

	base, err := NewBaseFirmaService(config)
	if err != nil {
		return nil, fmt.Errorf("error creando servicio base: %w", err)
	}

	return &SIIFirmaService{
		BaseFirmaService: base,
		certCache:        NewCertCache(24*time.Hour, 100),
	}, nil
}

// FirmarSemilla firma un documento de semilla del SII
func (s *SIIFirmaService) FirmarSemilla(semilla string) (*ResultadoFirma, error) {
	if semilla == "" {
		return nil, errors.New("la semilla no puede estar vacía")
	}

	// Crear XML de semilla
	xmlSemilla := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<SemillaXML>
    <Semilla>%s</Semilla>
</SemillaXML>`, semilla)

	return s.FirmarDocumento(xmlSemilla)
}

// FirmarToken firma un documento de token del SII
func (s *SIIFirmaService) FirmarToken(token string) (*ResultadoFirma, error) {
	if token == "" {
		return nil, errors.New("el token no puede estar vacío")
	}

	// Crear XML de token
	xmlToken := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<TokenXML>
    <Token>%s</Token>
</TokenXML>`, token)

	return s.FirmarDocumento(xmlToken)
}

// FirmarDTE firma un Documento Tributario Electrónico
func (s *SIIFirmaService) FirmarDTE(xml string) (*ResultadoFirma, error) {
	if xml == "" {
		return nil, errors.New("el XML no puede estar vacío")
	}

	// Validar que el RUT del emisor coincida con el configurado
	var dte struct {
		Documento struct {
			Encabezado struct {
				Emisor struct {
					RUTEmisor string
				}
			}
		}
	}

	if err := xml.Unmarshal([]byte(xml), &dte); err != nil {
		return nil, fmt.Errorf("error parseando DTE: %w", err)
	}

	if dte.Documento.Encabezado.Emisor.RUTEmisor != s.config.RutEmpresa {
		return nil, errors.New("el RUT del emisor no coincide con el configurado")
	}

	return s.FirmarDocumento(xml)
}

// ValidarCAF valida un archivo CAF
func (s *SIIFirmaService) ValidarCAF(caf []byte) error {
	if len(caf) == 0 {
		return errors.New("el CAF no puede estar vacío")
	}

	var cafXML struct {
		Version string `xml:"version,attr"`
		DA      struct {
			RE string // RUT Empresa
			FA string // Fecha Autorización
		}
	}

	if err := xml.Unmarshal(caf, &cafXML); err != nil {
		return fmt.Errorf("error parseando CAF: %w", err)
	}

	// Validar versión
	if cafXML.Version != "1.0" {
		return fmt.Errorf("versión de CAF no soportada: %s", cafXML.Version)
	}

	// Validar RUT empresa
	if cafXML.DA.RE != s.config.RutEmpresa {
		return errors.New("el RUT del CAF no coincide con el configurado")
	}

	// Validar fecha de autorización
	fechaAuth, err := time.Parse("2006-01-02", cafXML.DA.FA)
	if err != nil {
		return fmt.Errorf("formato de fecha inválido en CAF: %w", err)
	}

	if fechaAuth.After(time.Now()) {
		return errors.New("el CAF tiene una fecha de autorización futura")
	}

	return nil
}
