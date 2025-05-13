package services

import (
	"context"
	"errors"
	"time"

	"github.com/cursor/FMgo/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DocumentService implementa la interfaz domain.DocumentService
type DocumentService struct {
	repo          domain.DocumentRepository
	validationSvc domain.ValidationService
	cafSvc        domain.CAFService
	auditSvc      domain.AuditService
}

// NewDocumentService crea una nueva instancia del servicio de documentos
func NewDocumentService(
	repo domain.DocumentRepository,
	validationSvc domain.ValidationService,
	cafSvc domain.CAFService,
	auditSvc domain.AuditService,
) domain.DocumentService {
	return &DocumentService{
		repo:          repo,
		validationSvc: validationSvc,
		cafSvc:        cafSvc,
		auditSvc:      auditSvc,
	}
}

// CrearDocumento crea un nuevo documento
func (s *DocumentService) CrearDocumento(ctx context.Context, doc *domain.DocumentoTributario) error {
	// Validar documento
	if err := s.validationSvc.ValidarDocumento(doc); err != nil {
		return err
	}

	// Obtener y validar CAF
	caf, err := s.cafSvc.ObtenerCAF(ctx, doc.TipoDocumento)
	if err != nil {
		return err
	}
	if err := s.cafSvc.ValidarCAF(ctx, caf); err != nil {
		return err
	}

	// Asignar folio
	doc.Folio = caf.FolioActual

	// Guardar documento
	if err := s.repo.SaveDocumentoTributario(*doc); err != nil {
		return err
	}

	// Actualizar folio actual en CAF
	if err := s.cafSvc.ActualizarFolioActual(ctx, caf); err != nil {
		return err
	}

	// Registrar operación
	if err := s.auditSvc.RegistrarOperacion(ctx, "CREAR_DOCUMENTO", "DocumentoTributario", doc.ID, "sistema"); err != nil {
		return err
	}

	return nil
}

// ObtenerDocumento obtiene un documento por tipo y folio
func (s *DocumentService) ObtenerDocumento(ctx context.Context, tipo string, folio int64) (*domain.DocumentoTributario, error) {
	doc, err := s.repo.GetDocumentoTributario(tipo, folio)
	if err != nil {
		return nil, err
	}
	if doc == nil {
		return nil, errors.New("documento no encontrado")
	}
	return doc, nil
}

// ActualizarDocumento actualiza un documento existente
func (s *DocumentService) ActualizarDocumento(ctx context.Context, doc *domain.DocumentoTributario) error {
	// Validar documento
	if err := s.validationSvc.ValidarDocumento(doc); err != nil {
		return err
	}

	// Actualizar documento
	if err := s.repo.UpdateDocumentoTributario(*doc); err != nil {
		return err
	}

	// Registrar operación
	if err := s.auditSvc.RegistrarOperacion(ctx, "ACTUALIZAR_DOCUMENTO", "DocumentoTributario", doc.ID, "sistema"); err != nil {
		return err
	}

	return nil
}

// CambiarEstadoDocumento cambia el estado de un documento
func (s *DocumentService) CambiarEstadoDocumento(ctx context.Context, docID primitive.ObjectID, nuevoEstado string, usuario string) error {
	// Obtener documento actual
	doc, err := s.repo.GetDocumentoTributarioByID(docID)
	if err != nil {
		return err
	}
	if doc == nil {
		return errors.New("documento no encontrado")
	}

	// Crear nuevo estado
	estado := domain.EstadoDocumento{
		DocumentoID: docID,
		Estado:      nuevoEstado,
		Fecha:       time.Now(),
		Usuario:     usuario,
	}

	// Guardar estado
	if err := s.repo.SaveEstadoDocumento(estado); err != nil {
		return err
	}

	// Actualizar estado en documento
	doc.Estado = nuevoEstado
	if err := s.repo.UpdateDocumentoTributario(*doc); err != nil {
		return err
	}

	// Registrar operación
	if err := s.auditSvc.RegistrarOperacion(ctx, "CAMBIAR_ESTADO", "DocumentoTributario", docID, usuario); err != nil {
		return err
	}

	return nil
}

// AgregarReferencia agrega una referencia a un documento
func (s *DocumentService) AgregarReferencia(ctx context.Context, ref *domain.ReferenciaDocumento) error {
	// Validar documento origen
	docOrigen, err := s.repo.GetDocumentoTributario(ref.TipoOrigen, ref.FolioOrigen)
	if err != nil {
		return err
	}
	if docOrigen == nil {
		return errors.New("documento origen no encontrado")
	}

	// Validar documento referencia
	docRef, err := s.repo.GetDocumentoTributario(ref.TipoReferencia, ref.FolioReferencia)
	if err != nil {
		return err
	}
	if docRef == nil {
		return errors.New("documento referencia no encontrado")
	}

	// Guardar referencia
	if err := s.repo.SaveReferenciaDocumento(*ref); err != nil {
		return err
	}

	// Registrar operación
	if err := s.auditSvc.RegistrarOperacion(ctx, "AGREGAR_REFERENCIA", "ReferenciaDocumento", ref.ID, "sistema"); err != nil {
		return err
	}

	return nil
}

// ObtenerReferencias obtiene las referencias de un documento
func (s *DocumentService) ObtenerReferencias(ctx context.Context, tipoOrigen string, folioOrigen int64) ([]domain.ReferenciaDocumento, error) {
	return s.repo.GetReferenciasPorDocumento(tipoOrigen, folioOrigen)
}
