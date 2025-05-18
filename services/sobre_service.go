package services

import (
	"encoding/xml"
	"errors"
	"fmt"
	"time"

	"FMgo/models"
)

// SobreService representa el servicio para gestionar sobres de envío
type SobreService struct {
	firmaService *XMLSignatureService
}

// NewSobreService crea una nueva instancia del servicio de sobres
func NewSobreService(firmaService *XMLSignatureService) *SobreService {
	return &SobreService{
		firmaService: firmaService,
	}
}

func (s *SobreService) validarDatosSobre(sobre *models.EnvioDTE) error {
	if sobre == nil {
		return errors.New("el sobre no puede ser nulo")
	}

	if sobre.SetDTE.DTEs == nil || len(sobre.SetDTE.DTEs) == 0 {
		return errors.New("el sobre debe contener al menos un DTE")
	}

	if sobre.SetDTE.Caratula.RutEmisor == "" {
		return errors.New("RutEmisor es requerido en la carátula")
	}

	return nil
}

// GenerarSobreFactura genera un sobre de envío para facturas electrónicas
// según el esquema XML del SII (EnvioDTE_v10.xsd)
func (s *SobreService) GenerarSobreFactura(factura *models.Factura, resolucionSII time.Time, numeroResolucion int) (*models.EnvioDTE, error) {
	if factura == nil {
		return nil, fmt.Errorf("la factura no puede ser nula")
	}

	// Crear el sobre DTE
	sobre := &models.EnvioDTE{
		Version: "1.0",
		SetDTE: models.SetDTE{
			ID: fmt.Sprintf("SetDoc_%s", GenerateID()),
			Caratula: models.Caratula{
				Version:      "1.0",
				RutEmisor:    factura.RutEmisor,
				RutEnvia:     factura.RutEmisor, // Podría ser diferente en algunos casos
				RutReceptor:  factura.RutReceptor,
				FchResol:     resolucionSII.Format("2006-01-02"),
				NroResol:     numeroResolucion,
				TmstFirmaEnv: time.Now().Format("2006-01-02T15:04:05"),
				SubTotDTE: []models.SubTotalDTE{
					{
						TpoDTE: tipoDocumentoACodigo(factura.Tipo),
						NroDTE: 1,
					},
				},
			},
			DTEs: []models.DTEType{
				generarDTEDesdeFactura(factura),
			},
		},
		// La firma se aplicará posteriormente
	}

	// Aplicar firma digital si el servicio está disponible
	if s.firmaService != nil {
		if err := s.firmaService.FirmarEnvioDTE(sobre); err != nil {
			return nil, fmt.Errorf("error al firmar sobre: %v", err)
		}
	}

	return sobre, nil
}

func generarDetalles(items []models.Item) []models.DetalleXML {
	detalles := make([]models.DetalleXML, len(items))
	for i, item := range items {
		detalles[i] = models.DetalleXML{
			NroLinDet: i + 1,
			NmbItem:   item.Descripcion,
			QtyItem:   item.Cantidad,
			PrcItem:   item.PrecioUnitario,
			MontoItem: int(item.MontoItem),
		}

		if item.CodigoProducto != "" {
			detalles[i].CdgItem = models.CdgItem{
				TpoCodigo: "INT1", // Código interno
				VlrCodigo: item.CodigoProducto,
			}
		}
	}
	return detalles
}

func tipoDocumentoACodigo(tipo models.TipoDocumento) string {
	switch tipo {
	case models.TipoFactura:
		return "33"
	case models.TipoBoleta:
		return "39"
	case models.TipoNotaCredito:
		return "61"
	case models.TipoNotaDebito:
		return "56"
	case models.TipoGuiaDespacho:
		return "52"
	default:
		return "33" // Factura por defecto
	}
}

func (s *SobreService) firmarSobre(sobre *models.EnvioDTE) error {
	// Implementar firma del sobre usando el certificado digital
	return nil
}

// GenerarSobreBoleta genera un sobre de envío para boletas electrónicas
// según el esquema XML del SII (EnvioBOLETA_v11.xsd)
func (s *SobreService) GenerarSobreBoleta(boleta *models.Boleta, resolucionSII time.Time, numeroResolucion int) (*models.EnvioBOLETA, error) {
	if boleta == nil {
		return nil, fmt.Errorf("la boleta no puede ser nula")
	}

	// Crear el sobre BOLETA
	sobre := &models.EnvioBOLETA{
		Version: "1.0",
		SetDTE: models.SetDTEBol{
			ID: fmt.Sprintf("SetDoc_%s", GenerateID()),
			Caratula: models.Caratula{
				Version:      "1.0",
				RutEmisor:    boleta.RutEmisor,
				RutEnvia:     boleta.RutEmisor,
				RutReceptor:  "60803000-K", // RUT SII - Para boletas siempre es el SII
				FchResol:     resolucionSII.Format("2006-01-02"),
				NroResol:     numeroResolucion,
				TmstFirmaEnv: time.Now().Format("2006-01-02T15:04:05"),
				SubTotDTE: []models.SubTotalDTE{
					{
						TpoDTE: tipoDocumentoACodigo(boleta.Tipo),
						NroDTE: 1,
					},
				},
			},
			DTEs: []models.BOLETAType{
				generarBoletaDesdeModelo(boleta),
			},
		},
		// La firma se aplicará posteriormente
	}

	// Aplicar firma digital si el servicio está disponible
	if s.firmaService != nil {
		if err := s.firmaService.FirmarEnvioBOLETA(sobre); err != nil {
			return nil, fmt.Errorf("error al firmar sobre de boleta: %v", err)
		}
	}

	return sobre, nil
}

// GenerarSobreMultiple genera un sobre con múltiples documentos del mismo tipo
func (s *SobreService) GenerarSobreMultiple(documentos []models.DocumentoTributario, resolucionSII time.Time, numeroResolucion int) (*models.EnvioDTE, error) {
	if len(documentos) == 0 {
		return nil, fmt.Errorf("debe proporcionar al menos un documento")
	}

	// Verificar que todos los documentos son del mismo tipo y mismo emisor
	tipoPrimero := documentos[0].Tipo
	rutEmisor := documentos[0].RutEmisor
	for _, doc := range documentos {
		if doc.Tipo != tipoPrimero {
			return nil, fmt.Errorf("todos los documentos deben ser del mismo tipo")
		}
		if doc.RutEmisor != rutEmisor {
			return nil, fmt.Errorf("todos los documentos deben tener el mismo emisor")
		}
	}

	// Determinar receptor del sobre (SII o cliente específico)
	var rutReceptor string
	if esBoleta(tipoPrimero) {
		rutReceptor = "60803000-K" // RUT SII para boletas
	} else {
		rutReceptor = documentos[0].RutReceptor
	}

	// Crear el sobre DTE con múltiples documentos
	sobre := &models.EnvioDTE{
		Version: "1.0",
		SetDTE: models.SetDTE{
			ID: fmt.Sprintf("SetDoc_%s", GenerateID()),
			Caratula: models.Caratula{
				Version:      "1.0",
				RutEmisor:    rutEmisor,
				RutEnvia:     rutEmisor,
				RutReceptor:  rutReceptor,
				FchResol:     resolucionSII.Format("2006-01-02"),
				NroResol:     numeroResolucion,
				TmstFirmaEnv: time.Now().Format("2006-01-02T15:04:05"),
				SubTotDTE: []models.SubTotalDTE{
					{
						TpoDTE: tipoDocumentoACodigo(tipoPrimero),
						NroDTE: len(documentos),
					},
				},
			},
			DTEs: []models.DTEType{},
		},
	}

	// Agregar cada documento al sobre
	for _, doc := range documentos {
		dte := generarDTEDesdeDocumento(&doc)
		sobre.SetDTE.DTEs = append(sobre.SetDTE.DTEs, dte)
	}

	// Aplicar firma digital si el servicio está disponible
	if s.firmaService != nil {
		if err := s.firmaService.FirmarEnvioDTE(sobre); err != nil {
			return nil, fmt.Errorf("error al firmar sobre múltiple: %v", err)
		}
	}

	return sobre, nil
}

// GenerarSobreReceptor genera un sobre para el receptor (clientes)
func (s *SobreService) GenerarSobreReceptor(documentos []models.DocumentoTributario, resolucionSII time.Time, numeroResolucion int) (*models.EnvioDTE, error) {
	if len(documentos) == 0 {
		return nil, fmt.Errorf("debe proporcionar al menos un documento")
	}

	// En este caso, el sobre va dirigido al receptor/cliente
	rutEmisor := documentos[0].RutEmisor
	rutReceptor := documentos[0].RutReceptor

	// Agrupar documentos por tipo
	docsPorTipo := make(map[models.TipoDocumento]int)
	for _, doc := range documentos {
		docsPorTipo[doc.Tipo]++
	}

	// Crear subtotales para cada tipo de documento
	subtotales := []models.SubTotalDTE{}
	for tipo, cantidad := range docsPorTipo {
		subtotales = append(subtotales, models.SubTotalDTE{
			TpoDTE: tipoDocumentoACodigo(tipo),
			NroDTE: cantidad,
		})
	}

	// Crear el sobre DTE para el receptor
	sobre := &models.EnvioDTE{
		Version: "1.0",
		SetDTE: models.SetDTE{
			ID: fmt.Sprintf("SetDoc_%s", GenerateID()),
			Caratula: models.Caratula{
				Version:      "1.0",
				RutEmisor:    rutEmisor,
				RutEnvia:     rutEmisor,
				RutReceptor:  rutReceptor,
				FchResol:     resolucionSII.Format("2006-01-02"),
				NroResol:     numeroResolucion,
				TmstFirmaEnv: time.Now().Format("2006-01-02T15:04:05"),
				SubTotDTE:    subtotales,
			},
			DTEs: []models.DTEType{},
		},
	}

	// Agregar cada documento al sobre
	for _, doc := range documentos {
		dte := generarDTEDesdeDocumento(&doc)
		sobre.SetDTE.DTEs = append(sobre.SetDTE.DTEs, dte)
	}

	return sobre, nil
}

// GenerarSobreSII genera un sobre para el SII
func (s *SobreService) GenerarSobreSII(documentos []models.DocumentoTributario, resolucionSII time.Time, numeroResolucion int) (*models.EnvioDTE, error) {
	// El SII necesita sobres con todos los documentos
	// Dirigido especialmente al SII con RUT 60803000-K
	if len(documentos) == 0 {
		return nil, fmt.Errorf("debe proporcionar al menos un documento")
	}

	rutEmisor := documentos[0].RutEmisor

	// Agrupar documentos por tipo
	docsPorTipo := make(map[models.TipoDocumento]int)
	for _, doc := range documentos {
		if doc.RutEmisor != rutEmisor {
			return nil, fmt.Errorf("todos los documentos deben tener el mismo emisor")
		}
		docsPorTipo[doc.Tipo]++
	}

	// Crear subtotales para cada tipo de documento
	subtotales := []models.SubTotalDTE{}
	for tipo, cantidad := range docsPorTipo {
		subtotales = append(subtotales, models.SubTotalDTE{
			TpoDTE: tipoDocumentoACodigo(tipo),
			NroDTE: cantidad,
		})
	}

	// Crear el sobre DTE para el SII
	sobre := &models.EnvioDTE{
		Version: "1.0",
		SetDTE: models.SetDTE{
			ID: fmt.Sprintf("SetDoc_%s", GenerateID()),
			Caratula: models.Caratula{
				Version:      "1.0",
				RutEmisor:    rutEmisor,
				RutEnvia:     rutEmisor,
				RutReceptor:  "60803000-K", // RUT SII
				FchResol:     resolucionSII.Format("2006-01-02"),
				NroResol:     numeroResolucion,
				TmstFirmaEnv: time.Now().Format("2006-01-02T15:04:05"),
				SubTotDTE:    subtotales,
			},
			DTEs: []models.DTEType{},
		},
	}

	// Agregar cada documento al sobre
	for _, doc := range documentos {
		dte := generarDTEDesdeDocumento(&doc)
		sobre.SetDTE.DTEs = append(sobre.SetDTE.DTEs, dte)
	}

	return sobre, nil
}

// ConvertirAXML convierte un sobre EnvioDTE a formato XML
func (s *SobreService) ConvertirAXML(sobre *models.EnvioDTE) (string, error) {
	// Si el sobre no tiene firma y el servicio de firma está disponible, aplicarla
	if sobre.Signature == "" && s.firmaService != nil {
		if err := s.firmaService.FirmarEnvioDTE(sobre); err != nil {
			return "", fmt.Errorf("error al firmar sobre durante conversión a XML: %v", err)
		}
	}

	// Generar XML con indentación
	xmlBytes, err := xml.MarshalIndent(sobre, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error al generar XML: %v", err)
	}

	// Agregar declaración XML
	xmlString := "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n" + string(xmlBytes)

	return xmlString, nil
}

// ConvertirBoletaAXML convierte un sobre EnvioBOLETA a formato XML
func (s *SobreService) ConvertirBoletaAXML(sobre *models.EnvioBOLETA) (string, error) {
	// Si el sobre no tiene firma y el servicio de firma está disponible, aplicarla
	if sobre.Signature == "" && s.firmaService != nil {
		if err := s.firmaService.FirmarEnvioBOLETA(sobre); err != nil {
			return "", fmt.Errorf("error al firmar sobre de boleta durante conversión a XML: %v", err)
		}
	}

	// Generar XML con indentación
	xmlBytes, err := xml.MarshalIndent(sobre, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error al generar XML de boleta: %v", err)
	}

	// Agregar declaración XML
	xmlString := "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n" + string(xmlBytes)

	return xmlString, nil
}

// Funciones auxiliares

// generarDTEDesdeFactura convierte una factura en un DTE para XML
func generarDTEDesdeFactura(factura *models.Factura) models.DTEType {
	// Crear ID único para el documento
	docID := fmt.Sprintf("F%dT%s", factura.Folio, tipoDocumentoACodigo(factura.Tipo))

	// Configurar encabezado
	idDoc := models.IdDoc{
		TipoDTE: tipoDocumentoACodigo(factura.Tipo),
		Folio:   factura.Folio,
		FchEmis: factura.FechaEmision.Format("2006-01-02"),
	}

	// Si hay fecha de vencimiento, agregarla
	if !factura.FechaVencimiento.IsZero() {
		idDoc.FchVenc = factura.FechaVencimiento.Format("2006-01-02")
	}

	// Construir detalles
	detalles := []models.DetalleXML{}
	for i, item := range factura.Items {
		detalle := models.DetalleXML{
			NroLinDet: i + 1,
			NmbItem:   item.Descripcion,
			QtyItem:   item.Cantidad,
			PrcItem:   item.PrecioUnitario,
			MontoItem: int(item.MontoItem),
		}

		// Agregar código de producto si existe
		if item.CodigoProducto != "" {
			detalle.CdgItem = models.CdgItem{
				TpoCodigo: "INT1", // Código interno
				VlrCodigo: item.CodigoProducto,
			}
		}

		// Agregar descuento si existe
		if item.MontoDescuento > 0 {
			detalle.DescuentoMonto = models.MontoDescuento{
				TipoDesc:  "$",
				ValorDesc: int(item.MontoDescuento),
			}
		}

		detalles = append(detalles, detalle)
	}

	// Crear DTE
	dte := models.DTEType{
		Version: "1.0",
		Documento: models.DocumentoDTE{
			ID: docID,
			Encabezado: models.Encabezado{
				IdDoc: idDoc,
				Emisor: models.Emisor{
					RUTEmisor:    factura.RutEmisor,
					RznSoc:       factura.RazonSocialEmisor,
					GiroEmis:     "Actividad Comercial",  // Debe personalizarse
					Acteco:       "123456",               // Debe personalizarse
					DirOrigen:    "Dirección del Emisor", // Debe personalizarse
					CmnaOrigen:   "Comuna",               // Debe personalizarse
					CiudadOrigen: "Ciudad",               // Debe personalizarse
				},
				Receptor: models.Receptor{
					RUTRecep:    factura.RutReceptor,
					RznSocRecep: factura.RazonSocialReceptor,
					GiroRecep:   "Actividad Receptor",     // Debe personalizarse
					DirRecep:    "Dirección del Receptor", // Debe personalizarse
					CmnaRecep:   "Comuna",                 // Debe personalizarse
					CiudadRecep: "Ciudad",                 // Debe personalizarse
				},
				Totales: models.Totales{
					MntNeto:  int(factura.MontoNeto),
					MntExe:   int(factura.MontoExento),
					TasaIVA:  19, // IVA Chileno (debe personalizarse)
					IVA:      int(factura.MontoIVA),
					MntTotal: int(factura.MontoTotal),
				},
			},
			Detalle:   detalles,
			TmstFirma: time.Now().Format("2006-01-02T15:04:05"),
		},
	}

	return dte
}

// generarBoletaDesdeModelo convierte una boleta en un BOLETAType para XML
func generarBoletaDesdeModelo(boleta *models.Boleta) models.BOLETAType {
	// Crear ID único para el documento
	docID := fmt.Sprintf("F%dT%s", boleta.Folio, tipoDocumentoACodigo(boleta.Tipo))

	// Configurar encabezado
	idDoc := models.IdDocBOL{
		TipoDTE:     tipoDocumentoACodigo(boleta.Tipo),
		Folio:       boleta.Folio,
		FchEmis:     boleta.FechaEmision.Format("2006-01-02"),
		IndServicio: 3, // Boleta de Ventas y Servicio
	}

	// Construir detalles
	detalles := []models.DetalleXML{}
	for i, item := range boleta.Items {
		detalle := models.DetalleXML{
			NroLinDet: i + 1,
			NmbItem:   item.Descripcion,
			QtyItem:   item.Cantidad,
			PrcItem:   item.PrecioUnitario,
			MontoItem: int(item.MontoItem),
		}

		detalles = append(detalles, detalle)
	}

	// Crear BOLETAType
	boletaType := models.BOLETAType{
		Version: "1.0",
		Documento: models.DocumentoBoletaDTE{
			ID: docID,
			Encabezado: models.EncabezadoBOL{
				IdDoc: idDoc,
				Emisor: models.Emisor{
					RUTEmisor:    boleta.RutEmisor,
					RznSoc:       boleta.RazonSocialEmisor,
					GiroEmis:     "Actividad Comercial",  // Debe personalizarse
					Acteco:       "123456",               // Debe personalizarse
					DirOrigen:    "Dirección del Emisor", // Debe personalizarse
					CmnaOrigen:   "Comuna",               // Debe personalizarse
					CiudadOrigen: "Ciudad",               // Debe personalizarse
				},
				Totales: models.Totales{
					MntNeto:  int(boleta.MontoNeto),
					TasaIVA:  19, // IVA Chileno
					IVA:      int(boleta.MontoIVA),
					MntTotal: int(boleta.MontoTotal),
				},
			},
			Detalle:   detalles,
			TmstFirma: time.Now().Format("2006-01-02T15:04:05"),
		},
	}

	// Si hay datos del receptor, agregarlos (opcional en boletas)
	if boleta.RutReceptor != "" {
		boletaType.Documento.Encabezado.Receptor = models.ReceptorBOL{
			RUTRecep:    boleta.RutReceptor,
			RznSocRecep: boleta.RazonSocialReceptor,
		}
	}

	return boletaType
}

// generarDTEDesdeDocumento convierte un documento genérico en DTE para XML
func generarDTEDesdeDocumento(doc *models.DocumentoTributario) models.DTEType {
	// Crear ID único para el documento
	docID := fmt.Sprintf("F%dT%s", doc.Folio, tipoDocumentoACodigo(doc.Tipo))

	// Configurar encabezado
	idDoc := models.IdDoc{
		TipoDTE: tipoDocumentoACodigo(doc.Tipo),
		Folio:   doc.Folio,
		FchEmis: doc.FechaEmision.Format("2006-01-02"),
	}

	// Si hay fecha de vencimiento, agregarla
	if !doc.FechaVencimiento.IsZero() {
		idDoc.FchVenc = doc.FechaVencimiento.Format("2006-01-02")
	}

	// Construir detalles
	detalles := []models.DetalleXML{}
	for i, item := range doc.Items {
		detalle := models.DetalleXML{
			NroLinDet: i + 1,
			NmbItem:   item.Descripcion,
			QtyItem:   item.Cantidad,
			PrcItem:   item.PrecioUnitario,
			MontoItem: int(item.MontoItem),
		}

		// Agregar código de producto si existe
		if item.CodigoProducto != "" {
			detalle.CdgItem = models.CdgItem{
				TpoCodigo: "INT1", // Código interno
				VlrCodigo: item.CodigoProducto,
			}
		}

		detalles = append(detalles, detalle)
	}

	// Crear DTE
	dte := models.DTEType{
		Version: "1.0",
		Documento: models.DocumentoDTE{
			ID: docID,
			Encabezado: models.Encabezado{
				IdDoc: idDoc,
				Emisor: models.Emisor{
					RUTEmisor:    doc.RutEmisor,
					RznSoc:       doc.RazonSocialEmisor,
					GiroEmis:     "Actividad Comercial",  // Debe personalizarse
					Acteco:       "123456",               // Debe personalizarse
					DirOrigen:    "Dirección del Emisor", // Debe personalizarse
					CmnaOrigen:   "Comuna",               // Debe personalizarse
					CiudadOrigen: "Ciudad",               // Debe personalizarse
				},
				Receptor: models.Receptor{
					RUTRecep:    doc.RutReceptor,
					RznSocRecep: doc.RazonSocialReceptor,
					GiroRecep:   "Actividad Receptor",     // Debe personalizarse
					DirRecep:    "Dirección del Receptor", // Debe personalizarse
					CmnaRecep:   "Comuna",                 // Debe personalizarse
					CiudadRecep: "Ciudad",                 // Debe personalizarse
				},
				Totales: models.Totales{
					MntNeto:  int(doc.MontoNeto),
					MntExe:   int(doc.MontoExento),
					TasaIVA:  19, // IVA Chileno
					IVA:      int(doc.MontoIVA),
					MntTotal: int(doc.MontoTotal),
				},
			},
			Detalle:   detalles,
			TmstFirma: time.Now().Format("2006-01-02T15:04:05"),
		},
	}

	return dte
}

// esBoleta determina si un tipo de documento es boleta
func esBoleta(tipo models.TipoDocumento) bool {
	return tipo == models.TipoBoleta
}
