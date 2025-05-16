package services

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/fmgo/models"
)

// ReportesService maneja la generación de reportes
type ReportesService struct {
	db *mongo.Database
}

// NewReportesService crea una nueva instancia del servicio de reportes
func NewReportesService(db *mongo.Database) *ReportesService {
	return &ReportesService{
		db: db,
	}
}

// GenerarReporteDocumentosEstado genera un reporte de documentos por estado
func (s *ReportesService) GenerarReporteDocumentosEstado(ctx context.Context, fechaInicio, fechaFin time.Time, rutEmisor, rutReceptor string) (*models.ReporteDocumentosEstado, error) {
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

	// Calcular totales
	totalesPorEstado := make(map[models.EstadoDocumento]int)
	totalesPorTipo := make(map[models.TipoDocumento]int)

	for _, doc := range documentos {
		totalesPorEstado[doc.Estado]++
		totalesPorTipo[doc.Tipo]++
	}

	// Crear reporte
	reporte := &models.ReporteDocumentosEstado{
		ID:               GenerateID(),
		FechaInicio:      fechaInicio,
		FechaFin:         fechaFin,
		RutEmisor:        rutEmisor,
		RutReceptor:      rutReceptor,
		TotalesPorEstado: totalesPorEstado,
		TotalesPorTipo:   totalesPorTipo,
		Documentos:       documentos,
		FechaGeneracion:  time.Now(),
	}

	// Guardar reporte
	_, err = s.db.Collection("reportes_estado").InsertOne(ctx, reporte)
	if err != nil {
		return nil, err
	}

	return reporte, nil
}

// GenerarReporteRechazos genera un reporte de análisis de rechazos
func (s *ReportesService) GenerarReporteRechazos(ctx context.Context, fechaInicio, fechaFin time.Time, rutEmisor, rutReceptor string) (*models.ReporteRechazos, error) {
	// Construir filtro
	filtro := bson.M{
		"fecha_rechazo": bson.M{
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

	// Obtener documentos rechazados
	cursor, err := s.db.Collection("documentos_rechazados").Find(ctx, filtro)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var documentos []models.DocumentoRechazado
	if err = cursor.All(ctx, &documentos); err != nil {
		return nil, err
	}

	// Calcular errores comunes
	errores := make(map[string]*models.ErrorFrecuente)
	totalErrores := 0

	for _, doc := range documentos {
		for _, error := range doc.Errores {
			if err, exists := errores[error.Codigo]; exists {
				err.Frecuencia++
			} else {
				errores[error.Codigo] = &models.ErrorFrecuente{
					Codigo:      error.Codigo,
					Descripcion: error.Descripcion,
					Frecuencia:  1,
				}
			}
			totalErrores++
		}
	}

	// Calcular porcentajes
	var erroresComunes []models.ErrorFrecuente
	for _, err := range errores {
		err.Porcentaje = float64(err.Frecuencia) / float64(totalErrores) * 100
		erroresComunes = append(erroresComunes, *err)
	}

	// Crear reporte
	reporte := &models.ReporteRechazos{
		ID:                   GenerateID(),
		FechaInicio:          fechaInicio,
		FechaFin:             fechaFin,
		RutEmisor:            rutEmisor,
		RutReceptor:          rutReceptor,
		TotalRechazos:        len(documentos),
		ErroresComunes:       erroresComunes,
		DocumentosRechazados: documentos,
		FechaGeneracion:      time.Now(),
	}

	// Guardar reporte
	_, err = s.db.Collection("reportes_rechazos").InsertOne(ctx, reporte)
	if err != nil {
		return nil, err
	}

	return reporte, nil
}

// GenerarReporteMetricasRendimiento genera un reporte de métricas de rendimiento
func (s *ReportesService) GenerarReporteMetricasRendimiento(ctx context.Context, fechaInicio, fechaFin time.Time, rutEmisor string) (*models.ReporteMetricasRendimiento, error) {
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

	// Calcular métricas
	metricas := models.MetricasRendimiento{}
	var totalTiempoRespuesta time.Duration
	totalErrores := 0

	for _, doc := range documentos {
		metricas.TotalDocumentos++
		switch doc.Estado {
		case models.EstadoAceptado:
			metricas.DocumentosAceptados++
		case models.EstadoRechazado:
			metricas.DocumentosRechazados++
		case models.EstadoPendiente:
			metricas.DocumentosPendientes++
		}

		// Calcular tiempo de respuesta
		if !doc.FechaCreacion.IsZero() && !doc.FechaModificacion.IsZero() {
			totalTiempoRespuesta += doc.FechaModificacion.Sub(doc.FechaCreacion)
		}

		// Contar errores
		if doc.Estado == models.EstadoRechazado {
			totalErrores++
		}
	}

	// Calcular métricas derivadas
	if metricas.TotalDocumentos > 0 {
		metricas.TasaAceptacion = float64(metricas.DocumentosAceptados) / float64(metricas.TotalDocumentos) * 100
		metricas.TiempoPromedioRespuesta = totalTiempoRespuesta / time.Duration(metricas.TotalDocumentos)
		horas := fechaFin.Sub(fechaInicio).Hours()
		if horas > 0 {
			metricas.DocumentosPorHora = float64(metricas.TotalDocumentos) / horas
		}
		metricas.ErroresPorDocumento = float64(totalErrores) / float64(metricas.TotalDocumentos)
	}

	// Crear reporte
	reporte := &models.ReporteMetricasRendimiento{
		ID:              GenerateID(),
		FechaInicio:     fechaInicio,
		FechaFin:        fechaFin,
		RutEmisor:       rutEmisor,
		Metricas:        metricas,
		FechaGeneracion: time.Now(),
	}

	// Guardar reporte
	_, err = s.db.Collection("reportes_metricas").InsertOne(ctx, reporte)
	if err != nil {
		return nil, err
	}

	return reporte, nil
}

// GenerarReporteTributario genera un reporte tributario
func (s *ReportesService) GenerarReporteTributario(ctx context.Context, fechaInicio, fechaFin time.Time, rutEmisor, rutReceptor string) (*models.ReporteTributario, error) {
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

	// Calcular totales
	totales := models.TotalesTributarios{
		TotalesPorTipo: make(map[models.TipoDocumento]models.TotalesTipo),
	}

	for _, doc := range documentos {
		// Actualizar totales generales
		totales.MontoNetoTotal += doc.MontoNeto
		totales.MontoIVATotal += doc.MontoIVA
		totales.MontoRetencionTotal += doc.MontoRetencion
		totales.MontoTotal += doc.MontoTotal

		// Actualizar totales por tipo
		tipo := doc.Tipo
		if _, exists := totales.TotalesPorTipo[tipo]; !exists {
			totales.TotalesPorTipo[tipo] = models.TotalesTipo{}
		}
		totalesTipo := totales.TotalesPorTipo[tipo]
		totalesTipo.Cantidad++
		totalesTipo.MontoNeto += doc.MontoNeto
		totalesTipo.MontoIVA += doc.MontoIVA
		totalesTipo.MontoRetencion += doc.MontoRetencion
		totalesTipo.MontoTotal += doc.MontoTotal
		totales.TotalesPorTipo[tipo] = totalesTipo
	}

	// Crear reporte
	reporte := &models.ReporteTributario{
		ID:                 GenerateID(),
		FechaInicio:        fechaInicio,
		FechaFin:           fechaFin,
		RutEmisor:          rutEmisor,
		RutReceptor:        rutReceptor,
		TotalesTributarios: totales,
		Documentos:         documentos,
		FechaGeneracion:    time.Now(),
	}

	// Guardar reporte
	_, err = s.db.Collection("reportes_tributarios").InsertOne(ctx, reporte)
	if err != nil {
		return nil, err
	}

	return reporte, nil
}

// ObtenerReporte obtiene un reporte por su ID
func (s *ReportesService) ObtenerReporte(ctx context.Context, id string, tipo string) (interface{}, error) {
	var collection string
	switch tipo {
	case "estado":
		collection = "reportes_estado"
	case "rechazos":
		collection = "reportes_rechazos"
	case "metricas":
		collection = "reportes_metricas"
	case "tributario":
		collection = "reportes_tributarios"
	default:
		return nil, errors.New("tipo de reporte inválido")
	}

	var result interface{}
	err := s.db.Collection(collection).FindOne(ctx, bson.M{"_id": id}).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// ListarReportes obtiene una lista de reportes
func (s *ReportesService) ListarReportes(ctx context.Context, tipo string, fechaInicio, fechaFin time.Time, rutEmisor string) ([]interface{}, error) {
	var collection string
	switch tipo {
	case "estado":
		collection = "reportes_estado"
	case "rechazos":
		collection = "reportes_rechazos"
	case "metricas":
		collection = "reportes_metricas"
	case "tributario":
		collection = "reportes_tributarios"
	default:
		return nil, errors.New("tipo de reporte inválido")
	}

	// Construir filtro
	filtro := bson.M{
		"fecha_generacion": bson.M{
			"$gte": fechaInicio,
			"$lte": fechaFin,
		},
	}
	if rutEmisor != "" {
		filtro["rut_emisor"] = rutEmisor
	}

	// Obtener reportes
	cursor, err := s.db.Collection(collection).Find(ctx, filtro)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []interface{}
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}
