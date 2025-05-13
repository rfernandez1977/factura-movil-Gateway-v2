package models

import (
	"fmt"
)

// TipoDTE representa el tipo de documento tributario electrónico
type TipoDTE int

const (
	// Factura Electrónica
	FacturaElectronica TipoDTE = 33
	// Factura Exenta Electrónica
	FacturaExentaElectronica TipoDTE = 34
	// Boleta Electrónica
	BoletaElectronica TipoDTE = 39
	// Boleta Exenta Electrónica
	BoletaExentaElectronica TipoDTE = 41
	// Liquidación Factura Electrónica
	LiquidacionFacturaElectronica TipoDTE = 43
	// Guía de Despacho Electrónica
	GuiaDespachoElectronica TipoDTE = 52
	// Nota de Débito Electrónica
	NotaDebitoElectronica TipoDTE = 56
	// Nota de Crédito Electrónica
	NotaCreditoElectronica TipoDTE = 61
)

// String retorna la representación en string del tipo de documento
func (t TipoDTE) String() string {
	switch t {
	case FacturaElectronica:
		return "FACTURA ELECTRÓNICA"
	case FacturaExentaElectronica:
		return "FACTURA EXENTA ELECTRÓNICA"
	case BoletaElectronica:
		return "BOLETA ELECTRÓNICA"
	case BoletaExentaElectronica:
		return "BOLETA EXENTA ELECTRÓNICA"
	case LiquidacionFacturaElectronica:
		return "LIQUIDACIÓN FACTURA ELECTRÓNICA"
	case GuiaDespachoElectronica:
		return "GUÍA DE DESPACHO ELECTRÓNICA"
	case NotaDebitoElectronica:
		return "NOTA DE DÉBITO ELECTRÓNICA"
	case NotaCreditoElectronica:
		return "NOTA DE CRÉDITO ELECTRÓNICA"
	default:
		return "DESCONOCIDO"
	}
}

// ToInt retorna el valor numérico del tipo de documento
func (t TipoDTE) ToInt() int {
	return int(t)
}

// FromInt convierte un entero a TipoDTE
func FromInt(i int) TipoDTE {
	return TipoDTE(i)
}

// ToString retorna el valor como string (para compatibilidad con DTEType)
func (t TipoDTE) ToString() string {
	return fmt.Sprintf("%d", t)
}

// FromString convierte un string a TipoDTE
func FromString(s string) TipoDTE {
	switch s {
	case "FACTURA":
		return FacturaElectronica
	case "BOLETA":
		return BoletaElectronica
	case "NOTA_CREDITO":
		return NotaCreditoElectronica
	case "NOTA_DEBITO":
		return NotaDebitoElectronica
	case "GUIA_DESPACHO":
		return GuiaDespachoElectronica
	default:
		return 0
	}
}

// BOLETAType representa el tipo de boleta
type BOLETAType string

const (
	BOLETATypeElectronica BOLETAType = "ELECTRONICA"
	BOLETATypePapel       BOLETAType = "PAPEL"
)
