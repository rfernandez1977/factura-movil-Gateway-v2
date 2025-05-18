package validations

import (
	"context"
	"errors"
	"fmt"
	"time"

	"FMgo/core/sii/models"
)

// ValidadorSII maneja las validaciones específicas del SII
type ValidadorSII struct {
	config *ConfiguracionValidacion
}

// ConfiguracionValidacion contiene la configuración para las validaciones
type ConfiguracionValidacion struct {
	MaxDiasAntiguedad int     `json:"max_dias_antiguedad"`
	MaxMontoTotal     float64 `json:"max_monto_total"`
	MaxItems          int     `json:"max_items"`
	TasaIVA           float64 `json:"tasa_iva"`
	ValidarCAF        bool    `json:"validar_caf"`
	ValidarFirma      bool    `json:"validar_firma"`
	ValidarSchema     bool    `json:"validar_schema"`
}

// NewValidadorSII crea una nueva instancia del validador
func NewValidadorSII(config *ConfiguracionValidacion) *ValidadorSII {
	if config == nil {
		config = &ConfiguracionValidacion{
			MaxDiasAntiguedad: 30,
			MaxMontoTotal:     1000000000,
			MaxItems:          1000,
			TasaIVA:           19.0,
			ValidarCAF:        true,
			ValidarFirma:      true,
			ValidarSchema:     true,
		}
	}
	return &ValidadorSII{config: config}
}

// ValidarDocumento realiza todas las validaciones necesarias para un documento
func (v *ValidadorSII) ValidarDocumento(ctx context.Context, doc interface{}) (*models.ValidacionSII, error) {
	result := &models.ValidacionSII{
		FechaValidacion: time.Now(),
		Resultado:       true,
	}

	var err error
	switch d := doc.(type) {
	case *models.DTEXMLModel:
		err = v.validarDTE(ctx, d, result)
	default:
		return nil, errors.New("tipo de documento no soportado")
	}

	if err != nil {
		result.Resultado = false
		result.MensajeError = err.Error()
	}

	return result, nil
}

// validarDTE realiza las validaciones específicas para un DTE
func (v *ValidadorSII) validarDTE(ctx context.Context, dte *models.DTEXMLModel, result *models.ValidacionSII) error {
	if err := v.validarEncabezado(dte.Documento.Encabezado); err != nil {
		return fmt.Errorf("error en encabezado: %w", err)
	}

	if err := v.validarDetalle(dte.Documento.Detalle); err != nil {
		return fmt.Errorf("error en detalle: %w", err)
	}

	if v.config.ValidarFirma {
		if err := v.validarFirma(dte.Signature); err != nil {
			return fmt.Errorf("error en firma: %w", err)
		}
	}

	return nil
}

// validarEncabezado valida el encabezado del documento
func (v *ValidadorSII) validarEncabezado(enc models.EncabezadoXMLModel) error {
	if enc.IdDoc.FechaEmision == "" {
		return errors.New("fecha de emisión es requerida")
	}

	fechaEmision, err := time.Parse("2006-01-02", enc.IdDoc.FechaEmision)
	if err != nil {
		return fmt.Errorf("formato de fecha inválido: %w", err)
	}

	diasAntiguedad := int(time.Since(fechaEmision).Hours() / 24)
	if diasAntiguedad > v.config.MaxDiasAntiguedad {
		return fmt.Errorf("documento excede el máximo de días permitidos: %d", diasAntiguedad)
	}

	return nil
}

// validarDetalle valida el detalle del documento
func (v *ValidadorSII) validarDetalle(detalle []models.DetalleDTEXML) error {
	if len(detalle) == 0 {
		return errors.New("el documento debe tener al menos un detalle")
	}

	if len(detalle) > v.config.MaxItems {
		return fmt.Errorf("el documento excede el máximo de items permitidos: %d", v.config.MaxItems)
	}

	var montoTotal float64
	for _, d := range detalle {
		montoTotal += d.MontoItem
	}

	if montoTotal > v.config.MaxMontoTotal {
		return fmt.Errorf("el monto total excede el máximo permitido: %.2f", montoTotal)
	}

	return nil
}

// validarFirma valida la firma del documento
func (v *ValidadorSII) validarFirma(firma *models.FirmaXMLModel) error {
	if firma == nil {
		return errors.New("firma es requerida")
	}

	if firma.SignatureValue == "" {
		return errors.New("valor de firma es requerido")
	}

	// Aquí se pueden agregar más validaciones específicas de la firma

	return nil
}
