package dte

import (
	"context"
	"fmt"

	"FMgo/core/caf/services"
)

// ValidadorDTEImpl implementa la interfaz ValidadorDTE
type ValidadorDTEImpl struct {
	validadorCAF services.ValidadorCAF
}

// NewValidadorDTE crea una nueva instancia del validador DTE
func NewValidadorDTE(validadorCAF services.ValidadorCAF) ValidadorDTE {
	return &ValidadorDTEImpl{
		validadorCAF: validadorCAF,
	}
}

// ValidarDTE valida un documento tributario electrónico
func (v *ValidadorDTEImpl) ValidarDTE(dte *DTE) error {
	// Validar campos obligatorios
	if dte.TipoDTE == "" {
		return fmt.Errorf("tipo DTE es obligatorio")
	}
	if dte.Folio <= 0 {
		return fmt.Errorf("folio debe ser mayor a 0")
	}
	if dte.RUTEmisor == "" {
		return fmt.Errorf("RUT emisor es obligatorio")
	}
	if dte.RUTReceptor == "" {
		return fmt.Errorf("RUT receptor es obligatorio")
	}
	if dte.MontoTotal <= 0 {
		return fmt.Errorf("monto total debe ser mayor a 0")
	}
	if dte.FechaEmision == "" {
		return fmt.Errorf("fecha de emisión es obligatoria")
	}

	return nil
}

// ValidarCAF valida el CAF del documento
func (v *ValidadorDTEImpl) ValidarCAF(dte *DTE) error {
	tipoDTE := 0
	switch dte.TipoDTE {
	case "33":
		tipoDTE = 33
	case "34":
		tipoDTE = 34
	case "39":
		tipoDTE = 39
	case "41":
		tipoDTE = 41
	default:
		return fmt.Errorf("tipo DTE no soportado: %s", dte.TipoDTE)
	}

	// Validar folio usando el validador CAF
	valido, err := v.validadorCAF.ValidarFolio(context.Background(), int(dte.Folio), tipoDTE)
	if err != nil {
		return fmt.Errorf("error validando folio: %v", err)
	}
	if !valido {
		return fmt.Errorf("folio %d no está autorizado para tipo DTE %d", dte.Folio, tipoDTE)
	}

	return nil
}
