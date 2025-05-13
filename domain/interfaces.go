package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DocumentRepository define las operaciones para interactuar con la base de datos de documentos
type DocumentRepository interface {
	SaveDocumentoTributario(doc DocumentoTributario) error
	GetDocumentoTributario(tipo string, folio int64) (*DocumentoTributario, error)
	GetDocumentoTributarioByID(id primitive.ObjectID) (*DocumentoTributario, error)
	UpdateDocumentoTributario(doc DocumentoTributario) error
	GetDocumentosPorEstado(estado string) ([]DocumentoTributario, error)
	SaveEstadoDocumento(estado EstadoDocumento) error
	GetEstadoDocumento(docID primitive.ObjectID) (*EstadoDocumento, error)
	UpdateEstadoDocumento(estado EstadoDocumento) error
	SaveReferenciaDocumento(ref ReferenciaDocumento) error
	GetReferenciasPorDocumento(tipoOrigen string, folioOrigen int64) ([]ReferenciaDocumento, error)
}

// DocumentService define las operaciones de negocio para documentos
type DocumentService interface {
	CrearDocumento(ctx context.Context, doc *DocumentoTributario) error
	ObtenerDocumento(ctx context.Context, tipo string, folio int64) (*DocumentoTributario, error)
	ActualizarDocumento(ctx context.Context, doc *DocumentoTributario) error
	CambiarEstadoDocumento(ctx context.Context, docID primitive.ObjectID, nuevoEstado string, usuario string) error
	AgregarReferencia(ctx context.Context, ref *ReferenciaDocumento) error
	ObtenerReferencias(ctx context.Context, tipoOrigen string, folioOrigen int64) ([]ReferenciaDocumento, error)
}

// CAFService define las operaciones relacionadas con los CAF
type CAFService interface {
	ObtenerCAF(ctx context.Context, tipoDocumento string) (*CAF, error)
	ValidarCAF(ctx context.Context, caf *CAF) error
	ActualizarFolioActual(ctx context.Context, caf *CAF) error
	VerificarVencimientoCAF(ctx context.Context, caf *CAF) error
}

// ValidationService define las operaciones de validación
type ValidationService interface {
	ValidarDocumento(doc *DocumentoTributario) error
	ValidarRUT(rut string) error
	ValidarMonto(monto float64) error
	ValidarFecha(fecha time.Time) error
}

// AuditService define las operaciones de auditoría
type AuditService interface {
	RegistrarOperacion(ctx context.Context, operacion string, entidad string, entidadID primitive.ObjectID, usuario string) error
	ObtenerRegistros(ctx context.Context, filtro bson.M) ([]bson.M, error)
	GenerarReporte(ctx context.Context, fechaInicio, fechaFin time.Time) ([]bson.M, error)
}
