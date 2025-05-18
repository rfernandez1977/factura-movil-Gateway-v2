package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"FMgo/domain"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DocumentService implementa la interfaz domain.DocumentService
type DocumentService struct {
	repo          domain.DocumentRepository
	validationSvc domain.ValidationService
	cafSvc        domain.CAFService
	auditSvc      domain.AuditService
	cache         *redis.Client
}

const (
	// TTLDocumento tiempo de vida del documento en caché
	TTLDocumento = 24 * time.Hour
	// PrefijoDocumento prefijo para las claves de documentos
	PrefijoDocumento = "doc:"
	// PrefijoReferencias prefijo para las claves de referencias
	PrefijoReferencias = "ref:"
)

// NewDocumentService crea una nueva instancia del servicio de documentos
func NewDocumentService(
	repo domain.DocumentRepository,
	validationSvc domain.ValidationService,
	cafSvc domain.CAFService,
	auditSvc domain.AuditService,
	cache *redis.Client,
) domain.DocumentService {
	return &DocumentService{
		repo:          repo,
		validationSvc: validationSvc,
		cafSvc:        cafSvc,
		auditSvc:      auditSvc,
		cache:         cache,
	}
}

// getCacheKey genera una clave de caché para un documento
func (s *DocumentService) getCacheKey(tipo string, folio int64) string {
	return fmt.Sprintf("%s%s:%d", PrefijoDocumento, tipo, folio)
}

// getRefCacheKey genera una clave de caché para referencias
func (s *DocumentService) getRefCacheKey(tipo string, folio int64) string {
	return fmt.Sprintf("%s%s:%d", PrefijoReferencias, tipo, folio)
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

	// Guardar en caché
	if err := s.guardarEnCache(ctx, doc); err != nil {
		// Solo logear error de caché, no afecta operación principal
		fmt.Printf("error guardando en caché: %v\n", err)
	}

	// Registrar operación
	if err := s.auditSvc.RegistrarOperacion(ctx, "CREAR_DOCUMENTO", "DocumentoTributario", doc.ID, "sistema"); err != nil {
		return err
	}

	return nil
}

// guardarEnCache guarda un documento en el caché
func (s *DocumentService) guardarEnCache(ctx context.Context, doc *domain.DocumentoTributario) error {
	data, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("error serializando documento: %w", err)
	}

	key := s.getCacheKey(doc.TipoDocumento, doc.Folio)
	return s.cache.Set(ctx, key, data, TTLDocumento).Err()
}

// ObtenerDocumento obtiene un documento por tipo y folio
func (s *DocumentService) ObtenerDocumento(ctx context.Context, tipo string, folio int64) (*domain.DocumentoTributario, error) {
	// Intentar obtener del caché
	key := s.getCacheKey(tipo, folio)
	data, err := s.cache.Get(ctx, key).Bytes()
	if err == nil {
		var doc domain.DocumentoTributario
		if err := json.Unmarshal(data, &doc); err == nil {
			return &doc, nil
		}
	}

	// Si no está en caché, obtener de la base de datos
	doc, err := s.repo.GetDocumentoTributario(tipo, folio)
	if err != nil {
		return nil, err
	}
	if doc == nil {
		return nil, errors.New("documento no encontrado")
	}

	// Guardar en caché para futuras consultas
	if err := s.guardarEnCache(ctx, doc); err != nil {
		fmt.Printf("error guardando en caché: %v\n", err)
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

	// Actualizar caché
	if err := s.guardarEnCache(ctx, doc); err != nil {
		fmt.Printf("error actualizando caché: %v\n", err)
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

	// Actualizar caché
	if err := s.guardarEnCache(ctx, doc); err != nil {
		fmt.Printf("error actualizando caché: %v\n", err)
	}

	// Registrar operación
	if err := s.auditSvc.RegistrarOperacion(ctx, "CAMBIAR_ESTADO", "DocumentoTributario", docID, usuario); err != nil {
		return err
	}

	return nil
}

// guardarReferenciasEnCache guarda las referencias en el caché
func (s *DocumentService) guardarReferenciasEnCache(ctx context.Context, tipo string, folio int64, refs []domain.ReferenciaDocumento) error {
	data, err := json.Marshal(refs)
	if err != nil {
		return fmt.Errorf("error serializando referencias: %w", err)
	}

	key := s.getRefCacheKey(tipo, folio)
	return s.cache.Set(ctx, key, data, TTLDocumento).Err()
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

	// Invalidar caché de referencias
	key := s.getRefCacheKey(ref.TipoOrigen, ref.FolioOrigen)
	if err := s.cache.Del(ctx, key).Err(); err != nil {
		fmt.Printf("error invalidando caché de referencias: %v\n", err)
	}

	// Registrar operación
	if err := s.auditSvc.RegistrarOperacion(ctx, "AGREGAR_REFERENCIA", "ReferenciaDocumento", ref.ID, "sistema"); err != nil {
		return err
	}

	return nil
}

// ObtenerReferencias obtiene las referencias de un documento
func (s *DocumentService) ObtenerReferencias(ctx context.Context, tipoOrigen string, folioOrigen int64) ([]domain.ReferenciaDocumento, error) {
	// Intentar obtener del caché
	key := s.getRefCacheKey(tipoOrigen, folioOrigen)
	data, err := s.cache.Get(ctx, key).Bytes()
	if err == nil {
		var refs []domain.ReferenciaDocumento
		if err := json.Unmarshal(data, &refs); err == nil {
			return refs, nil
		}
	}

	// Si no está en caché, obtener de la base de datos
	refs, err := s.repo.GetReferenciasPorDocumento(tipoOrigen, folioOrigen)
	if err != nil {
		return nil, err
	}

	// Guardar en caché para futuras consultas
	if err := s.guardarReferenciasEnCache(ctx, tipoOrigen, folioOrigen, refs); err != nil {
		fmt.Printf("error guardando referencias en caché: %v\n", err)
	}

	return refs, nil
}

// LimpiarCache limpia el caché de documentos y referencias
func (s *DocumentService) LimpiarCache(ctx context.Context) error {
	// Obtener todas las claves con los prefijos
	var cursor uint64
	var keys []string

	for {
		var result []string
		var err error
		result, cursor, err = s.cache.Scan(ctx, cursor, PrefijoDocumento+"*", 100).Result()
		if err != nil {
			return fmt.Errorf("error escaneando claves de documentos: %w", err)
		}
		keys = append(keys, result...)

		result, cursor, err = s.cache.Scan(ctx, cursor, PrefijoReferencias+"*", 100).Result()
		if err != nil {
			return fmt.Errorf("error escaneando claves de referencias: %w", err)
		}
		keys = append(keys, result...)

		if cursor == 0 {
			break
		}
	}

	// Eliminar todas las claves encontradas
	if len(keys) > 0 {
		if err := s.cache.Del(ctx, keys...).Err(); err != nil {
			return fmt.Errorf("error eliminando claves del caché: %w", err)
		}
	}

	return nil
}
