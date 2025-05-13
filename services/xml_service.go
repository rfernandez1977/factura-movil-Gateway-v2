package services

import (
	"context"
	"encoding/xml"
	"fmt"
	"strconv"
	"time"

	"github.com/cursor/FMgo/config"
	"github.com/cursor/FMgo/models"
	"github.com/cursor/FMgo/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Constantes para tipos de documento
const (
	TipoDTEFactura      = models.FacturaElectronica
	TipoDTEBoleta       = models.BoletaElectronica
	TipoDTEGuiaDespacho = models.GuiaDespachoElectronica
	TipoDTENotaDebito   = models.NotaDebitoElectronica
	TipoDTENotaCredito  = models.NotaCreditoElectronica
)

// Constantes para mensajes de error
const (
	ErrTipoDocumentoNoSoportado = "tipo de documento no soportado"
	ErrTipoDocumentoIncorrecto  = "tipo de documento incorrecto para %s"
	ErrFolioInvalido            = "folio inválido"
	ErrFechaEmisionRequerida    = "fecha de emisión requerida"
	ErrRUTEmisorRequerido       = "RUT del emisor requerido"
	ErrRUTReceptorRequerido     = "RUT del receptor requerido"
	ErrItemsRequeridos          = "el documento debe tener al menos un item"
	ErrMontoTotalInvalido       = "monto total debe ser mayor que cero"
	ErrIVANegativo              = "IVA no puede ser negativo"
	ErrPatenteRequerida         = "patente del vehículo requerida"
	ErrRUTChoferRequerido       = "RUT del chofer requerido"
	ErrReferenciaRequerida      = "la %s requiere una referencia al documento original"
)

// Constantes para nombres de campos XML
const (
	XMLVersion = "1.0"

	// Nombres de elementos
	XMLDocumento  = "Documento"
	XMLEncabezado = "Encabezado"
	XMLIdDoc      = "IdDoc"
	XMLEmisor     = "Emisor"
	XMLReceptor   = "Receptor"
	XMLTotales    = "Totales"
	XMLDetalle    = "Detalle"
	XMLItem       = "Item"
	XMLReferencia = "Referencia"
	XMLTransporte = "Transporte"

	// Atributos
	XMLAttrVersion = "version"

	// Campos de IdDoc
	XMLTipoDTE      = "TipoDTE"
	XMLFolio        = "Folio"
	XMLFchEmis      = "FchEmis"
	XMLIndServicio  = "IndServicio"
	XMLFmaPago      = "FmaPago"
	XMLFchVenc      = "FchVenc"
	XMLIndMntNeto   = "IndMntNeto"
	XMLPeriodoDesde = "PeriodoDesde"
	XMLPeriodoHasta = "PeriodoHasta"
	XMLMntPago      = "MntPago"
	XMLIndTraslado  = "IndTraslado"
	XMLTpoDespacho  = "TpoDespacho"
	XMLIndDocRef    = "IndDocRef"
	XMLFolioDocRef  = "FolioDocRef"
	XMLFchDocRef    = "FchDocRef"

	// Campos de Emisor/Receptor
	XMLRUT       = "RUT"
	XMLRznSoc    = "RznSoc"
	XMLGiro      = "Giro"
	XMLDireccion = "Direccion"
	XMLComuna    = "Comuna"
	XMLCiudad    = "Ciudad"

	// Campos de Totales
	XMLMntNeto  = "MntNeto"
	XMLMntExe   = "MntExe"
	XMLIVA      = "IVA"
	XMLMntTotal = "MntTotal"
	XMLMontoNF  = "MontoNF"
	XMLVlrPagar = "VlrPagar"

	// Campos de Item
	XMLNroLinDet  = "NroLinDet"
	XMLNmItem     = "NmItem"
	XMLQtyItem    = "QtyItem"
	XMLUnmdItem   = "UnmdItem"
	XMLPrcItem    = "PrcItem"
	XMLMontoItem  = "MontoItem"
	XMLDescuentoP = "DescuentoP"
	XMLDescuentoM = "DescuentoM"

	// Campos de Referencia
	XMLTpoDocRef = "TpoDocRef"
	XMLFolioRef  = "FolioRef"
	XMLFchRef    = "FchRef"
	XMLCodRef    = "CodRef"
	XMLRazonRef  = "RazonRef"

	// Campos de Transporte
	XMLPatente   = "Patente"
	XMLRUTChofer = "RUTChofer"
	XMLNombre    = "Nombre"
)

// XMLService maneja la generación y almacenamiento de XML para documentos tributarios
type XMLService struct {
	config *config.SupabaseConfig
	db     *mongo.Database
}

// NewXMLService crea una nueva instancia del servicio XML
func NewXMLService(config *config.SupabaseConfig, db *mongo.Database) *XMLService {
	return &XMLService{
		config: config,
		db:     db,
	}
}

// Helper para convertir TipoDocumento a int
func tipoDocumentoToInt(tipo models.TipoDTE) int {
	return tipo.ToInt()
}

// GenerarXMLFactura genera el XML para una factura electrónica
func (s *XMLService) GenerarXMLFactura(ctx context.Context, documento *models.DocumentoTributario, empresa *models.Empresa) (string, error) {
	// Estructura para el XML de factura
	type FacturaXML struct {
		XMLName    xml.Name `xml:"Documento"`
		Version    string   `xml:"version,attr"`
		Encabezado struct {
			ID           string    `xml:"IdDoc>Folio"`
			FechaEmision time.Time `xml:"IdDoc>FchEmis"`
			TipoDTE      string    `xml:"IdDoc>TipoDTE"`
			RutEmisor    string    `xml:"Emisor>RUTEmisor"`
			RazonSocial  string    `xml:"Emisor>RznSoc"`
			Giro         string    `xml:"Emisor>GiroEmis"`
			Direccion    string    `xml:"Emisor>DirOrigen"`
			Comuna       string    `xml:"Emisor>CmnaOrigen"`
			RutReceptor  string    `xml:"Receptor>RUTRecep"`
			RznReceptor  string    `xml:"Receptor>RznSocRecep"`
			MontoTotal   float64   `xml:"Totales>MntTotal"`
			MontoNeto    float64   `xml:"Totales>MntNeto"`
			IVA          float64   `xml:"Totales>IVA"`
		}
		Detalle []struct {
			NmbItem   string  `xml:"NmbItem"`
			QtyItem   float64 `xml:"QtyItem"`
			PrcItem   float64 `xml:"PrcItem"`
			MontoItem float64 `xml:"MontoItem"`
		}
	}

	// Crear instancia de factura
	factura := FacturaXML{
		Version: "1.0",
	}

	// Llenar datos del encabezado
	factura.Encabezado.ID = fmt.Sprintf("%d", documento.Folio)
	factura.Encabezado.FechaEmision = documento.FechaEmision
	factura.Encabezado.TipoDTE = documento.TipoDocumento.String()
	factura.Encabezado.RutEmisor = documento.RutEmisor
	factura.Encabezado.RazonSocial = documento.RazonSocialEmisor
	factura.Encabezado.Giro = documento.GiroEmisor
	factura.Encabezado.Direccion = documento.DireccionEmisor
	factura.Encabezado.Comuna = documento.ComunaEmisor

	// Convertir a XML
	xmlData, err := xml.MarshalIndent(factura, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error al generar XML: %v", err)
	}

	// Agregar encabezado XML
	xmlString := fmt.Sprintf(`<?xml version="1.0" encoding="ISO-8859-1"?>
%s`, string(xmlData))

	return xmlString, nil
}

// GuardarXML guarda el XML en la base de datos
func (s *XMLService) GuardarXML(ctx context.Context, documentoID string, xmlContent string) error {
	// Actualizar el documento con el XML generado
	_, err := s.db.Collection("documentos").UpdateOne(
		ctx,
		bson.M{"_id": documentoID},
		bson.M{
			"$set": bson.M{
				"xml":        xmlContent,
				"updated_at": time.Now(),
			},
		},
	)

	if err != nil {
		return fmt.Errorf("error al guardar XML: %v", err)
	}

	return nil
}

// ObtenerXML obtiene el XML de un documento
func (s *XMLService) ObtenerXML(ctx context.Context, documentoID string) (string, error) {
	var documento models.DocumentoTributario
	err := s.db.Collection("documentos").FindOne(
		ctx,
		bson.M{"_id": documentoID},
	).Decode(&documento)

	if err != nil {
		return "", fmt.Errorf("error al obtener XML: %v", err)
	}

	// Aquí deberías reconstruir el XML a partir del documento si es necesario
	return "<xml>no implementado</xml>", nil
}

// ValidarXML valida un documento XML según su tipo
func (s *XMLService) ValidarXML(xmlData []byte, tipoDTE int) error {
	// Validar que sea un XML válido
	var temp interface{}
	if err := xml.Unmarshal(xmlData, &temp); err != nil {
		return fmt.Errorf("XML inválido: %v", err)
	}

	// Validar contra el esquema XSD correspondiente
	validator, err := utils.NuevoValidadorXML("schemas/DTE_v10.xsd")
	if err != nil {
		return fmt.Errorf("error creando validador: %v", err)
	}
	switch tipoDTE {
	case TipoDTEFactura.ToInt():
		return validator.ValidarXML(xmlData)
	case TipoDTEBoleta.ToInt():
		return validator.ValidarXML(xmlData)
	case TipoDTENotaCredito.ToInt():
		return validator.ValidarXML(xmlData)
	case TipoDTENotaDebito.ToInt():
		return validator.ValidarXML(xmlData)
	case TipoDTEGuiaDespacho.ToInt():
		return validator.ValidarXML(xmlData)
	default:
		return fmt.Errorf("tipo de documento no soportado: %d", tipoDTE)
	}
}

// GenerarXML genera el XML para un documento según su tipo
func (s *XMLService) GenerarXML(doc *models.DocumentoTributario) ([]byte, error) {
	xmlDoc, err := s.generarXMLDocumentoTributario(doc)
	if err != nil {
		return nil, fmt.Errorf("error generando XML: %v", err)
	}

	return xml.MarshalIndent(xmlDoc, "", "  ")
}

// generarXMLDocumentoTributario genera el XML para un documento tributario
func (s *XMLService) generarXMLDocumentoTributario(doc *models.DocumentoTributario) (*models.DTEXMLModel, error) {
	xmlDoc := &models.DTEXMLModel{
		Version: "1.0",
		Documento: models.DocumentoXMLModel{
			ID: fmt.Sprintf("T%dF%d", doc.TipoDocumento, doc.Folio),
			Encabezado: models.EncabezadoXMLModel{
				ID: models.IDDocumentoXMLModel{
					TipoDTE:      doc.TipoDocumento,
					Folio:        doc.Folio,
					FechaEmision: doc.FechaEmision.Format("2006-01-02"),
				},
				Emisor: models.EmisorXMLModel{
					RUT:         doc.RutEmisor,
					RazonSocial: doc.RazonSocialEmisor,
					Giro:        doc.GiroEmisor,
					Direccion:   doc.DireccionEmisor,
					Comuna:      doc.ComunaEmisor,
					Ciudad:      doc.CiudadEmisor,
				},
				Receptor: models.ReceptorXMLModel{
					RUT:         doc.RutReceptor,
					RazonSocial: doc.RazonSocialReceptor,
					Giro:        doc.GiroReceptor,
					Direccion:   doc.DireccionReceptor,
					Comuna:      doc.ComunaReceptor,
					Ciudad:      doc.CiudadReceptor,
				},
				Totales: models.TotalesXMLModel{
					MontoNeto:  int(doc.MontoNeto),
					IVA:        doc.TasaIVA,
					MontoIVA:   int(doc.MontoIVA),
					MontoTotal: int(doc.MontoTotal),
				},
			},
		},
	}

	// Agregar detalles
	for i, item := range doc.Items {
		detalle := models.DetalleDTEXML{
			NroLinDet: i + 1,
			NmbItem:   item.Nombre,
			QtyItem:   item.Cantidad,
			PrcItem:   item.PrecioUnitario,
			MontoItem: item.PrecioUnitario * int64(item.Cantidad),
			UnmdItem:  item.UnidadMedida,
		}
		xmlDoc.Documento.Detalle = append(xmlDoc.Documento.Detalle, detalle)
	}

	// Agregar referencias
	for i, ref := range doc.Referencias {
		referencia := models.ReferenciaXMLModel{
			NroLinRef: i + 1,
			TpoDocRef: strconv.Itoa(ref.TipoDocumento.ToInt()),
			FolioRef:  ref.Folio,
			FchRef:    ref.Fecha.Format("2006-01-02"),
			RazonRef:  ref.Razon,
			CodRef:    ref.Codigo,
		}
		xmlDoc.Documento.Referencia = append(xmlDoc.Documento.Referencia, referencia)
	}

	return xmlDoc, nil
}

// ParsearXML parsea un documento XML según su tipo
func (s *XMLService) ParsearXML(xmlData []byte, tipoDTE int) (interface{}, error) {
	var doc models.DocumentoTributario
	if err := xml.Unmarshal(xmlData, &doc); err != nil {
		return nil, fmt.Errorf("error al parsear documento: %v", err)
	}
	return &doc, nil
}
