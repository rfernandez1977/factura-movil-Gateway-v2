package services

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"FMgo/models"
)

// ReportesAuditoriaService maneja la generación de reportes de auditoría
type ReportesAuditoriaService struct {
	db *mongo.Database
}

// NewReportesAuditoriaService crea una nueva instancia del servicio de reportes de auditoría
func NewReportesAuditoriaService(db *mongo.Database) *ReportesAuditoriaService {
	return &ReportesAuditoriaService{
		db: db,
	}
}

// GenerarReporteAuditoria genera un reporte de auditoría
func (s *ReportesAuditoriaService) GenerarReporteAuditoria(
	ctx context.Context,
	fechaInicio, fechaFin time.Time,
	rutEmisor, rutReceptor string,
) (*models.ReporteAuditoria, error) {
	// Construir filtro
	filtro := bson.M{
		"fecha_cambio": bson.M{
			"$gte": fechaInicio,
			"$lte": fechaFin,
		},
	}
	if rutEmisor != "" {
		filtro["rut_emisor"] = rutEmisor
	}
	if rutReceptor != "" {
		filtro["rut_receptor"] = rutReceptor
	}

	// Obtener registros de auditoría
	cursor, err := s.db.Collection("auditoria").Find(ctx, filtro)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var registros []models.RegistroAuditoria
	if err = cursor.All(ctx, &registros); err != nil {
		return nil, err
	}

	// Calcular estadísticas
	cambiosPorTipo := make(map[models.TipoDocumento]int)
	cambiosPorEstado := make(map[string]int)
	cambiosPorUsuario := make(map[string]int)

	for _, reg := range registros {
		cambiosPorTipo[reg.TipoDocumento]++
		cambiosPorEstado[reg.EstadoNuevo]++
		cambiosPorUsuario[reg.Usuario]++
	}

	// Crear reporte
	reporte := &models.ReporteAuditoria{
		ID:                 generateID(),
		FechaInicio:        fechaInicio,
		FechaFin:           fechaFin,
		RutEmisor:          rutEmisor,
		RutReceptor:        rutReceptor,
		CambiosPorTipo:     cambiosPorTipo,
		CambiosPorEstado:   cambiosPorEstado,
		CambiosPorUsuario:  cambiosPorUsuario,
		RegistrosAuditoria: registros,
		FechaGeneracion:    time.Now(),
	}

	// Guardar reporte
	_, err = s.db.Collection("reportes_auditoria").InsertOne(ctx, reporte)
	if err != nil {
		return nil, err
	}

	return reporte, nil
}

// GenerarReporteCumplimiento genera un reporte de cumplimiento
func (s *ReportesAuditoriaService) GenerarReporteCumplimiento(
	ctx context.Context,
	fechaInicio, fechaFin time.Time,
	rutEmisor, rutReceptor string,
) (*models.ReporteCumplimiento, error) {
	// Construir filtro
	filtro := bson.M{
		"fecha_emision": bson.M{
			"$gte": fechaInicio,
			"$lte": fechaFin,
		},
	}
	if rutEmisor != "" {
		filtro["rut_emisor"] = rutEmisor
	}
	if rutReceptor != "" {
		filtro["rut_receptor"] = rutReceptor
	}

	// Obtener documentos
	cursor, err := s.db.Collection("documentos").Find(ctx, filtro)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var documentos []models.DocumentoTributario
	if err = cursor.All(ctx, &documentos); err != nil {
		return nil, err
	}

	// Identificar documentos vencidos y generar alertas
	var documentosVencidos []models.DocumentoVencido
	var alertas []models.AlertaCumplimiento

	for _, doc := range documentos {
		// Verificar vencimientos
		if !doc.FechaVencimiento.IsZero() && time.Now().After(doc.FechaVencimiento) {
			diasVencido := int(time.Now().Sub(doc.FechaVencimiento).Hours() / 24)
			documentosVencidos = append(documentosVencidos, models.DocumentoVencido{
				DocumentoID:      doc.ID,
				TipoDocumento:    doc.Tipo,
				FechaEmision:     doc.FechaEmision,
				FechaVencimiento: doc.FechaVencimiento,
				DiasVencido:      diasVencido,
				Estado:           doc.Estado,
			})

			// Generar alerta
			severidad := "BAJA"
			if diasVencido > 30 {
				severidad = "ALTA"
			} else if diasVencido > 15 {
				severidad = "MEDIA"
			}

			alertas = append(alertas, models.AlertaCumplimiento{
				Tipo:        "VENCIMIENTO",
				Descripcion: "Documento vencido",
				Severidad:   severidad,
				FechaAlerta: time.Now(),
				DocumentoID: doc.ID,
				Estado:      "PENDIENTE",
			})
		}

		// Verificar otros criterios de cumplimiento
		if doc.Estado == models.EstadoPendiente && time.Now().Sub(doc.FechaEmision).Hours() > 24 {
			alertas = append(alertas, models.AlertaCumplimiento{
				Tipo:        "PENDIENTE",
				Descripcion: "Documento pendiente por más de 24 horas",
				Severidad:   "MEDIA",
				FechaAlerta: time.Now(),
				DocumentoID: doc.ID,
				Estado:      "PENDIENTE",
			})
		}
	}

	// Crear reporte
	reporte := &models.ReporteCumplimiento{
		ID:                 generateID(),
		FechaInicio:        fechaInicio,
		FechaFin:           fechaFin,
		RutEmisor:          rutEmisor,
		RutReceptor:        rutReceptor,
		DocumentosVencidos: documentosVencidos,
		Alertas:            alertas,
		FechaGeneracion:    time.Now(),
	}

	// Guardar reporte
	_, err = s.db.Collection("reportes_cumplimiento").InsertOne(ctx, reporte)
	if err != nil {
		return nil, err
	}

	return reporte, nil
}

// ObtenerReporteAuditoria obtiene un reporte de auditoría por su ID
func (s *ReportesAuditoriaService) ObtenerReporteAuditoria(ctx context.Context, id string) (*models.ReporteAuditoria, error) {
	var reporte models.ReporteAuditoria
	err := s.db.Collection("reportes_auditoria").FindOne(ctx, bson.M{"_id": id}).Decode(&reporte)
	if err != nil {
		return nil, err
	}
	return &reporte, nil
}

// ObtenerReporteCumplimiento obtiene un reporte de cumplimiento por su ID
func (s *ReportesAuditoriaService) ObtenerReporteCumplimiento(ctx context.Context, id string) (*models.ReporteCumplimiento, error) {
	var reporte models.ReporteCumplimiento
	err := s.db.Collection("reportes_cumplimiento").FindOne(ctx, bson.M{"_id": id}).Decode(&reporte)
	if err != nil {
		return nil, err
	}
	return &reporte, nil
}

// ListarReportesAuditoria obtiene una lista de reportes de auditoría
func (s *ReportesAuditoriaService) ListarReportesAuditoria(
	ctx context.Context,
	fechaInicio, fechaFin time.Time,
	rutEmisor string,
) ([]models.ReporteAuditoria, error) {
	filtro := bson.M{
		"fecha_generacion": bson.M{
			"$gte": fechaInicio,
			"$lte": fechaFin,
		},
	}
	if rutEmisor != "" {
		filtro["rut_emisor"] = rutEmisor
	}

	cursor, err := s.db.Collection("reportes_auditoria").Find(ctx, filtro)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var reportes []models.ReporteAuditoria
	if err = cursor.All(ctx, &reportes); err != nil {
		return nil, err
	}

	return reportes, nil
}

// ListarReportesCumplimiento obtiene una lista de reportes de cumplimiento
func (s *ReportesAuditoriaService) ListarReportesCumplimiento(
	ctx context.Context,
	fechaInicio, fechaFin time.Time,
	rutEmisor string,
) ([]models.ReporteCumplimiento, error) {
	filtro := bson.M{
		"fecha_generacion": bson.M{
			"$gte": fechaInicio,
			"$lte": fechaFin,
		},
	}
	if rutEmisor != "" {
		filtro["rut_emisor"] = rutEmisor
	}

	cursor, err := s.db.Collection("reportes_cumplimiento").Find(ctx, filtro)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var reportes []models.ReporteCumplimiento
	if err = cursor.All(ctx, &reportes); err != nil {
		return nil, err
	}

	return reportes, nil
}
