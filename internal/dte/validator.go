package dte

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

// DTE representa la estructura de un Documento Tributario Electrónico
type DTE struct {
	TipoDTE  string        `json:"tipo_dte"`
	Emisor   Contribuyente `json:"emisor"`
	Receptor Contribuyente `json:"receptor"`
	Detalles []Detalle     `json:"detalles"`
	Totales  Totales       `json:"totales"`
}

// Contribuyente representa a un emisor o receptor
type Contribuyente struct {
	RUT         string `json:"rut"`
	RazonSocial string `json:"razon_social"`
}

// Detalle representa un ítem del DTE
type Detalle struct {
	Cantidad       int    `json:"cantidad"`
	Descripcion    string `json:"descripcion"`
	PrecioUnitario int    `json:"precio_unitario"`
	MontoTotal     int    `json:"monto_total"`
}

// Totales representa los montos totales del DTE
type Totales struct {
	MontoNeto int `json:"monto_neto"`
	TasaIVA   int `json:"tasa_iva"`
	IVA       int `json:"iva"`
	Total     int `json:"total"`
}

// ValidateRUT valida el formato y dígito verificador de un RUT
func ValidateRUT(rut string) error {
	// Expresión regular para validar formato XXXXXXXX-Y
	rutRegex := regexp.MustCompile(`^[0-9]{1,8}-[0-9K]$`)
	if !rutRegex.MatchString(rut) {
		return errors.New("formato de RUT inválido")
	}

	// Separar número y dígito verificador
	parts := regexp.MustCompile(`-`).Split(rut, -1)
	numero := parts[0]
	dv := parts[1]

	// Calcular dígito verificador
	suma := 0
	multiplicador := 2

	// Recorrer el número de derecha a izquierda
	for i := len(numero) - 1; i >= 0; i-- {
		d, _ := strconv.Atoi(string(numero[i]))
		suma += d * multiplicador
		multiplicador++
		if multiplicador > 7 {
			multiplicador = 2
		}
	}

	// Calcular dígito verificador
	dvCalculado := 11 - (suma % 11)
	dvEsperado := ""

	switch dvCalculado {
	case 11:
		dvEsperado = "0"
	case 10:
		dvEsperado = "K"
	default:
		dvEsperado = strconv.Itoa(dvCalculado)
	}

	// Comparar con el dígito verificador proporcionado
	if dv != dvEsperado {
		return errors.New("dígito verificador inválido")
	}

	return nil
}

// Validate valida la estructura y datos del DTE
func (d *DTE) Validate() error {
	// Validar tipo DTE
	tiposValidos := map[string]bool{"33": true, "34": true, "39": true, "41": true}
	if !tiposValidos[d.TipoDTE] {
		return fmt.Errorf("tipo de DTE inválido: %s", d.TipoDTE)
	}

	// Validar RUTs
	if err := ValidateRUT(d.Emisor.RUT); err != nil {
		return fmt.Errorf("RUT emisor inválido: %v", err)
	}
	if err := ValidateRUT(d.Receptor.RUT); err != nil {
		return fmt.Errorf("RUT receptor inválido: %v", err)
	}

	// Validar detalles
	if len(d.Detalles) == 0 {
		return errors.New("DTE debe tener al menos un detalle")
	}

	montoNeto := 0
	for _, detalle := range d.Detalles {
		if detalle.Cantidad <= 0 {
			return errors.New("cantidad debe ser mayor a 0")
		}
		if detalle.PrecioUnitario <= 0 {
			return errors.New("precio unitario debe ser mayor a 0")
		}
		if detalle.MontoTotal != detalle.Cantidad*detalle.PrecioUnitario {
			return errors.New("monto total del detalle no coincide")
		}
		montoNeto += detalle.MontoTotal
	}

	// Validar totales
	if d.Totales.MontoNeto != montoNeto {
		return errors.New("monto neto total no coincide con la suma de detalles")
	}

	iva := (montoNeto * d.Totales.TasaIVA) / 100
	if d.Totales.IVA != iva {
		return errors.New("IVA calculado no coincide")
	}

	if d.Totales.Total != montoNeto+iva {
		return errors.New("monto total no coincide")
	}

	return nil
}
