package sii

import (
	"fmt"

	"github.com/cursor/FMgo/models"
)

// CalcularIVA calcula el IVA de un monto
func CalcularIVA(monto float64) float64 {
	return monto * 0.19
}

// CalcularMontoTotal calcula el monto total incluyendo IVA
func CalcularMontoTotal(montoNeto float64) float64 {
	return montoNeto + CalcularIVA(montoNeto)
}

// ValidarRUT valida un RUT chileno
func ValidarRUT(rut string) bool {
	// Implementar validación de RUT
	return true
}

// ValidarFolio valida un folio de DTE
func ValidarFolio(folio int, tipoDTE models.TipoDTE) error {
	if folio <= 0 {
		return fmt.Errorf("el folio debe ser mayor a 0")
	}

	// Validar rango según tipo de DTE
	if folio > 999999999 {
		return fmt.Errorf("el folio excede el máximo permitido")
	}

	// Validar tipo de DTE
	switch tipoDTE {
	case 33, 34, 39, 52, 56, 61:
		return nil
	default:
		return fmt.Errorf("tipo de DTE no válido")
	}
}

// ValidarMontos valida los montos de un DTE
func ValidarMontos(montoNeto, montoIVA, montoTotal float64) error {
	if montoNeto < 0 {
		return fmt.Errorf("el monto neto no puede ser negativo")
	}

	if montoIVA < 0 {
		return fmt.Errorf("el monto de IVA no puede ser negativo")
	}

	if montoTotal < 0 {
		return fmt.Errorf("el monto total no puede ser negativo")
	}

	ivaCalculado := CalcularIVA(montoNeto)
	if montoIVA != ivaCalculado {
		return fmt.Errorf("el monto de IVA no corresponde al monto neto")
	}

	totalCalculado := CalcularMontoTotal(montoNeto)
	if montoTotal != totalCalculado {
		return fmt.Errorf("el monto total no corresponde a la suma del monto neto e IVA")
	}

	return nil
}

// TributarioSII representa el servicio para interactuar con el SII
type TributarioSII struct {
	ambiente        string // PRODUCCION, CERTIFICACION
	credenciales    map[string]string
	certificadoPath string
	clavePrivada    string
}

// NewTributarioSII crea una nueva instancia del servicio TributarioSII
func NewTributarioSII(ambiente, certificadoPath, clavePrivada string) *TributarioSII {
	return &TributarioSII{
		ambiente:        ambiente,
		credenciales:    make(map[string]string),
		certificadoPath: certificadoPath,
		clavePrivada:    clavePrivada,
	}
}

// generarXMLImpuestos genera el XML de impuestos para el SII
func (s *TributarioSII) generarXMLImpuestos(doc interface{}) (string, error) {
	var (
		montoNeto        float64
		montoIVA         float64
		montoExento      float64
		impuestosParaXML []models.ImpuestoAdicionalItem // Nombre más claro sobre su propósito
	)

	switch d := doc.(type) {
	case *models.Factura:
		montoNeto = d.MontoNeto
		montoIVA = d.MontoIVA
		montoExento = d.MontoExento
		impuestosParaXML = s.obtenerImpuestosParaXML(d.Items)
	case *models.Boleta:
		montoNeto = d.MontoNeto
		montoIVA = d.MontoIVA
		montoExento = d.MontoExento
		impuestosParaXML = s.obtenerImpuestosParaXML(d.Items)
	case *models.NotaCredito:
		montoNeto = d.MontoNeto
		montoIVA = d.MontoIVA
		montoExento = d.MontoExento
		impuestosParaXML = s.obtenerImpuestosParaXML(d.Items)
	case *models.NotaDebito:
		montoNeto = d.MontoNeto
		montoIVA = d.MontoIVA
		montoExento = d.MontoExento
		impuestosParaXML = s.obtenerImpuestosParaXML(d.Items)
	case *models.GuiaDespacho:
		montoNeto = d.MontoNeto
		montoIVA = d.MontoIVA
		montoExento = d.MontoExento
		impuestosParaXML = s.obtenerImpuestosParaXML(d.Items)
	}

	// Crear XML de impuestos
	xml := `<?xml version="1.0" encoding="ISO-8859-1"?>
<DTE version="1.0">
	<Documento>
		<Encabezado>
			<Totales>
				<MntNeto>%.2f</MntNeto>
				<MntExe>%.2f</MntExe>
				<IVA>%.2f</IVA>`

	// Agregar impuestos adicionales
	if len(impuestosParaXML) > 0 {
		xml += `
				<ImptoReten>
					<Indicador>1</Indicador>
					<Codigo>15</Codigo>
					<Tasa>%.2f</Tasa>
					<Valor>%.2f</Valor>
				</ImptoReten>`
	}

	// Cerrar XML
	xml += `
			</Totales>
		</Encabezado>
	</Documento>
</DTE>`

	// Formatear XML con los valores
	if len(impuestosParaXML) > 0 {
		return fmt.Sprintf(xml, montoNeto, montoExento, montoIVA, impuestosParaXML[0].Porcentaje, impuestosParaXML[0].Monto), nil
	}
	return fmt.Sprintf(xml, montoNeto, montoExento, montoIVA), nil
}

func (s *TributarioSII) GenerateXML(doc interface{}) (string, error) {
	var (
		montoNeto        float64
		montoIVA         float64
		montoTotal       float64
		montoExento      float64
		impuestosParaXML []models.ImpuestoAdicionalItem
	)

	switch d := doc.(type) {
	case *models.Factura:
		montoNeto = d.MontoNeto
		montoIVA = d.MontoIVA
		montoTotal = d.MontoTotal
		montoExento = d.MontoExento
		impuestosParaXML = s.obtenerImpuestosParaXML(d.Items)
	case *models.Boleta:
		montoNeto = d.MontoNeto
		montoIVA = d.MontoIVA
		montoTotal = d.MontoTotal
		montoExento = d.MontoExento
		impuestosParaXML = s.obtenerImpuestosParaXML(d.Items)
	case *models.NotaCredito:
		montoNeto = d.MontoNeto
		montoIVA = d.MontoIVA
		montoTotal = d.MontoTotal
		montoExento = d.MontoExento
		impuestosParaXML = s.obtenerImpuestosParaXML(d.Items)
	case *models.NotaDebito:
		montoNeto = d.MontoNeto
		montoIVA = d.MontoIVA
		montoTotal = d.MontoTotal
		montoExento = d.MontoExento
		impuestosParaXML = s.obtenerImpuestosParaXML(d.Items)
	case *models.GuiaDespacho:
		montoNeto = d.MontoNeto
		montoIVA = d.MontoIVA
		montoTotal = d.MontoTotal
		montoExento = d.MontoExento
		impuestosParaXML = s.obtenerImpuestosParaXML(d.Items)
	default:
		return "", fmt.Errorf("tipo de documento no soportado")
	}

	// Generar XML según tipo de documento
	xml := `
	<DTE>
		<Encabezado>
			<IdDoc>
				<TipoDTE>33</TipoDTE>
			</IdDoc>
			<Totales>
				<MntNeto>%.2f</MntNeto>
				<MntExe>%.2f</MntExe>
				<IVA>%.2f</IVA>
				<MntTotal>%.2f</MntTotal>
			</Totales>
		</Encabezado>
	</DTE>`

	// Si hay impuestos adicionales, incluirlos en el XML
	if len(impuestosParaXML) > 0 {
		xml = `
		<DTE>
			<Encabezado>
				<IdDoc>
					<TipoDTE>33</TipoDTE>
				</IdDoc>
				<Totales>
					<MntNeto>%.2f</MntNeto>
					<MntExe>%.2f</MntExe>
					<IVA>%.2f</IVA>
					<ImptoReten>
						<TipoImp>%s</TipoImp>
						<TasaImp>%.2f</TasaImp>
						<MontoImp>%.2f</MontoImp>
					</ImptoReten>
					<MntTotal>%.2f</MntTotal>
				</Totales>
			</Encabezado>
		</DTE>`
		return fmt.Sprintf(xml, montoNeto, montoExento, montoIVA, impuestosParaXML[0].Codigo, impuestosParaXML[0].Porcentaje, impuestosParaXML[0].Monto, montoTotal), nil
	}

	return fmt.Sprintf(xml, montoNeto, montoExento, montoIVA, montoTotal), nil
}

// obtenerImpuestosParaXML extrae los impuestos adicionales para usar en el XML
// Solo se incluye el primer impuesto encontrado, según requerimiento del SII
func (s *TributarioSII) obtenerImpuestosParaXML(items []models.Item) []models.ImpuestoAdicionalItem {
	var impuestosParaXML []models.ImpuestoAdicionalItem

	// Buscar el primer impuesto adicional encontrado en cualquier ítem
	for _, item := range items {
		if len(item.ImpuestosAdicionales) > 0 {
			// Solo agregamos el primer impuesto para el XML
			impuestosParaXML = append(impuestosParaXML, item.ImpuestosAdicionales[0])
			return impuestosParaXML // Retornamos inmediatamente después de encontrar el primero
		}
	}

	return impuestosParaXML
}
