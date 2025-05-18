package services

import (
	"context"

	"FMgo/core/firma/common"
	"FMgo/core/firma/models"
)

// CertificadoRepository define las operaciones del repositorio de certificados
type CertificadoRepository interface {
	// ObtenerCertificado obtiene un certificado por su ID
	ObtenerCertificado(ctx context.Context, id string) (*models.Certificado, error)
	// GuardarCertificado guarda un certificado
	GuardarCertificado(ctx context.Context, cert *models.Certificado) error
	// ListarCertificados lista todos los certificados
	ListarCertificados(ctx context.Context) ([]*models.Certificado, error)
	// EliminarCertificado elimina un certificado
	EliminarCertificado(ctx context.Context, id string) error
}

// CacheService define las operaciones del servicio de caché
type CacheService interface {
	// ObtenerCertificado obtiene un certificado del caché
	ObtenerCertificado(ctx context.Context, id string) (*models.Certificado, error)
	// GuardarCertificado guarda un certificado en el caché
	GuardarCertificado(ctx context.Context, id string, cert *models.Certificado) error
	// EliminarCertificado elimina un certificado del caché
	EliminarCertificado(ctx context.Context, id string) error
	// LimpiarCache limpia todo el caché
	LimpiarCache(ctx context.Context) error
}

// Logger es un alias para common.Logger
type Logger = common.Logger

// CAFRepository define las operaciones del repositorio de CAF
type CAFRepository interface {
	// GuardarCAF guarda un nuevo CAF
	GuardarCAF(ctx context.Context, caf *models.CAF) error
	// ObtenerCAF obtiene un CAF por su ID
	ObtenerCAF(ctx context.Context, id string) (*models.CAF, error)
	// ObtenerCAFPorFolio obtiene un CAF que contiene el folio especificado
	ObtenerCAFPorFolio(ctx context.Context, tipo string, folio int64) (*models.CAF, error)
	// ListarCAFsPorTipo lista todos los CAFs de un tipo de documento
	ListarCAFsPorTipo(ctx context.Context, tipo string) ([]*models.CAF, error)
	// ActualizarEstadoCAF actualiza el estado de un CAF
	ActualizarEstadoCAF(ctx context.Context, id string, estado string) error
}

// AlertService define las operaciones del servicio de alertas
type AlertService interface {
	// EnviarAlerta envía una alerta
	EnviarAlerta(ctx context.Context, alerta *models.Alerta) error
	// ObtenerAlertas obtiene las alertas activas
	ObtenerAlertas(ctx context.Context) ([]*models.Alerta, error)
	// MarcarAlertaComoLeida marca una alerta como leída
	MarcarAlertaComoLeida(ctx context.Context, id string) error
}
