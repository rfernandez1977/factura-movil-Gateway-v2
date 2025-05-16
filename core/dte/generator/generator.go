package generator

import (
	"fmt"
	"time"

	"github.com/fmgo/core/dte/types"
)

// Generator es la interfaz para generar DTEs
type Generator interface {
	GenerateDTE(input *GenerateInput) (*types.DTE, error)
}

// GenerateInput contiene los datos necesarios para generar un DTE
type GenerateInput struct {
	TipoDTE      string
	Folio        int
	FechaEmision time.Time
	Emisor       types.Emisor
	Receptor     types.Receptor
	Detalles     []types.Detalle
	Referencias  []types.Referencia
	Descuentos   []types.Descuento
	Recargos     []types.Recargo
}

// DefaultGenerator es la implementación por defecto del generador de DTEs
type DefaultGenerator struct{}

// NewDefaultGenerator crea una nueva instancia del generador por defecto
func NewDefaultGenerator() *DefaultGenerator {
	return &DefaultGenerator{}
}

// GenerateDTE genera un nuevo DTE a partir de los datos de entrada
func (g *DefaultGenerator) GenerateDTE(input *GenerateInput) (*types.DTE, error) {
	// Calcular totales
	montoNeto, iva, montoTotal := calcularTotales(input.Detalles)

	// Crear el documento
	doc := &types.Documento{
		Encabezado: types.Encabezado{
			IDDocumento: types.IDDocumento{
				TipoDTE:      input.TipoDTE,
				Folio:        input.Folio,
				FechaEmision: input.FechaEmision,
			},
			Emisor:   input.Emisor,
			Receptor: input.Receptor,
			Totales: types.Totales{
				MontoNeto:  montoNeto,
				TasaIVA:    19,
				IVA:        iva,
				MontoTotal: montoTotal,
			},
		},
		Detalles:    input.Detalles,
		Referencias: input.Referencias,
		Descuentos:  input.Descuentos,
		Recargos:    input.Recargos,
	}

	// Crear el DTE
	dte := &types.DTE{
		ID:            generateID(input.TipoDTE, input.Folio),
		Documento:     *doc,
		FechaCreacion: time.Now(),
		Estado:        "CREADO",
		Firmado:       false,
	}

	return dte, nil
}

// calcularTotales calcula los montos totales del documento
func calcularTotales(detalles []types.Detalle) (montoNeto, iva, montoTotal float64) {
	for _, det := range detalles {
		if !det.Exento {
			montoNeto += det.MontoItem
		}
	}

	iva = montoNeto * 0.19
	montoTotal = montoNeto + iva

	return
}

// generateID genera un ID único para el DTE
func generateID(tipoDTE string, folio int) string {
	return fmt.Sprintf("DTE-%s-%d-%d", tipoDTE, folio, time.Now().Unix())
}
