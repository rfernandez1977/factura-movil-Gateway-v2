package services

import (
	"context"
	"fmt"
	"time"

	"FMgo/models"
)

// FacturaService maneja la lógica de negocio de facturas
// Ahora usa SupabaseService para persistencia
// y mapea de models.Factura a models.Documento

type FacturaService struct {
	supabase     *SupabaseService
	xmlService   *XMLService
	firmaService *FirmaService
	siiService   *SIIService
	cafService   *CAFService
}

func NewFacturaService(
	supabase *SupabaseService,
	xmlService *XMLService,
	firmaService *FirmaService,
	siiService *SIIService,
	cafService *CAFService,
) *FacturaService {
	return &FacturaService{
		supabase:     supabase,
		xmlService:   xmlService,
		firmaService: firmaService,
		siiService:   siiService,
		cafService:   cafService,
	}
}

// CrearFactura crea una nueva factura y la guarda como Documento
func (s *FacturaService) CrearFactura(ctx context.Context, empresa *models.Empresa, factura *models.Factura) (*models.Documento, error) {
	// Validar factura
	if err := s.validarFactura(factura); err != nil {
		return nil, err
	}

	// Obtener CAF disponible
	caf, err := s.cafService.GetCAFDisponible(empresa.ID, "33")
	if err != nil {
		return nil, fmt.Errorf("error al obtener CAF: %v", err)
	}

	// Asignar folio (usamos FolioActual para la demo, idealmente deberías llevar el folio actual en la base de datos)
	factura.Folio = int(caf.FolioActual)
	factura.FechaEmision = time.Now()

	// Mapear a Documento
	doc := &models.Documento{
		EmpresaID:     empresa.ID,
		TipoDocumento: "FACTURA",
		Folio:         int64(factura.Folio),
		RutEmisor:     factura.RutEmisor,
		RutReceptor:   factura.RutReceptor,
		MontoTotal:    int64(factura.MontoTotal),
		Estado:        "PENDIENTE",
		FechaEmision:  factura.FechaEmision,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Guardar documento en Supabase
	err = s.supabase.GuardarDocumento(ctx, doc)
	if err != nil {
		return nil, fmt.Errorf("error al guardar documento: %v", err)
	}

	// Actualizar folio del CAF (solo si caf.ID existe)
	// if caf.ID != "" {
	// 	_ = s.cafService.ActualizarFolioCAF(ctx, caf.ID)
	// }

	return doc, nil
}

// ProcesarFactura: genera XML, firma, guarda y envía al SII
func (s *FacturaService) ProcesarFactura(ctx context.Context, docID string, empresa *models.Empresa) error {
	// Obtener documento
	doc, err := s.supabase.ObtenerDocumento(ctx, docID)
	if err != nil {
		return fmt.Errorf("error al obtener documento: %v", err)
	}

	// Mapear a Factura (solo los campos necesarios)
	factura := &models.Factura{
		RutEmisor:   doc.RutEmisor,
		RutReceptor: doc.RutReceptor,
		Folio:       int(doc.Folio),
		// FechaEmision:        doc.FechaEmision, // Omitido si no existe
		MontoTotal: float64(doc.MontoTotal),
		Items:      nil, // Si necesitas los items, deberás guardarlos en el documento o en otra tabla
	}

	// Generar XML
	xmlData, err := s.xmlService.GenerarXML(factura)
	if err != nil {
		return fmt.Errorf("error al generar XML: %v", err)
	}

	// Firmar XML (usamos RutEmisor como empresaID si no hay campo EmpresaID)
	xmlFirmado, err := s.firmaService.FirmarXML(xmlData, doc.RutEmisor)
	if err != nil {
		return fmt.Errorf("error al firmar XML: %v", err)
	}

	// Guardar XML en el documento
	doc.XML = string(xmlFirmado)
	doc.Estado = "FIRMADO"
	doc.UpdatedAt = time.Now()
	err = s.supabase.GuardarDocumento(ctx, doc)
	if err != nil {
		return fmt.Errorf("error al guardar XML: %v", err)
	}

	// Enviar al SII (ajustar argumentos según la firma)
	_, err = s.siiService.EnviarDTE(xmlFirmado, empresa)
	if err != nil {
		return fmt.Errorf("error al enviar al SII: %v", err)
	}

	// Actualizar estado y trackID
	doc.Estado = "ENVIADO"
	doc.UpdatedAt = time.Now()
	// doc.TrackID = respuesta.TrackID // Si tienes este campo en Documento
	err = s.supabase.GuardarDocumento(ctx, doc)
	if err != nil {
		return fmt.Errorf("error al actualizar estado: %v", err)
	}

	return nil
}

// ConsultarEstado consulta el estado de una factura en el SII y actualiza el documento
func (s *FacturaService) ConsultarEstado(ctx context.Context, docID string) (string, error) {
	doc, err := s.supabase.ObtenerDocumento(ctx, docID)
	if err != nil {
		return "", fmt.Errorf("error al obtener documento: %v", err)
	}

	estado, err := s.siiService.ConsultarEstado(docID)
	if err != nil {
		return "", fmt.Errorf("error al consultar estado: %v", err)
	}

	doc.Estado = estado
	doc.UpdatedAt = time.Now()
	err = s.supabase.GuardarDocumento(ctx, doc)
	if err != nil {
		return "", fmt.Errorf("error al actualizar estado: %v", err)
	}

	return estado, nil
}

// validarFactura valida una factura antes de crearla
func (s *FacturaService) validarFactura(factura *models.Factura) error {
	if factura.RutEmisor == "" {
		return fmt.Errorf("RUT emisor requerido")
	}
	if factura.RutReceptor == "" {
		return fmt.Errorf("RUT del receptor requerido")
	}
	if len(factura.Items) == 0 {
		return fmt.Errorf("la factura debe tener al menos un item")
	}
	return nil
}
