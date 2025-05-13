package services

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"net/smtp"
	"os"
	"path/filepath"
	"time"

	"github.com/cursor/FMgo/models"
	"github.com/go-redis/redis/v8"
	"github.com/jung-kurt/gofpdf"
	"github.com/wcharczuk/go-chart"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// FolioService maneja la integración con el sistema de folios
type FolioService struct {
	db           *mongo.Database
	cafService   *CAFService
	redisClient  *redis.Client
	umbralFolios int // Cantidad mínima de folios antes de solicitar nuevo CAF
}

// Folio representa un folio individual
type Folio struct {
	ID          string    `bson:"_id"`
	RUTEmisor   string    `bson:"rut_emisor"`
	TipoDTE     string    `bson:"tipo_dte"`
	Numero      int       `bson:"numero"`
	Estado      string    `bson:"estado"` // DISPONIBLE, UTILIZADO, ANULADO
	CAFID       string    `bson:"caf_id"`
	FechaUso    time.Time `bson:"fecha_uso,omitempty"`
	DocumentoID string    `bson:"documento_id,omitempty"`
}

// RangoFolios representa un rango de folios de un CAF
type RangoFolios struct {
	RUTEmisor         string
	TipoDTE           string
	FolioInicial      int
	FolioFinal        int
	CAFID             string
	FechaAutorizacion time.Time
	FechaVencimiento  time.Time
}

// EstadoCAF representa el estado de un CAF
type EstadoCAF struct {
	CAFID             string
	RUTEmisor         string
	TipoDTE           string
	Estado            string // VIGENTE, PROXIMO_EXPIRAR, EXPIRADO
	DiasRestantes     int
	FoliosDisponibles int
	FoliosUtilizados  int
	FoliosAnulados    int
}

// ReporteFolios representa un reporte de uso de folios
type ReporteFolios struct {
	RUTEmisor         string
	TipoDTE           string
	PeriodoInicio     time.Time
	PeriodoFin        time.Time
	TotalFolios       int
	FoliosDisponibles int
	FoliosUtilizados  int
	FoliosAnulados    int
	DetalleUso        []DetalleUsoFolio
	CAFActivos        []EstadoCAF
}

// DetalleUsoFolio representa el detalle de uso de un folio
type DetalleUsoFolio struct {
	Numero      int
	Estado      string
	FechaUso    time.Time
	DocumentoID string
	CAFID       string
}

// AlertConfig configura las alertas por email
type AlertConfig struct {
	SMTPHost     string
	SMTPPort     int
	SMTPUser     string
	SMTPPassword string
	FromEmail    string
	ToEmails     []string
}

// EmailAlertService maneja el envío de alertas por email
type EmailAlertService struct {
	config AlertConfig
}

// ChartConfig configura la generación de gráficos
type ChartConfig struct {
	Width     int
	Height    int
	Theme     string
	FontSize  int
	OutputDir string
}

// NotificationConfig configura las notificaciones
type NotificationConfig struct {
	SlackWebhookURL string
	TeamsWebhookURL string
	TelegramToken   string
	TelegramChatID  string
}

// NotificationService maneja las notificaciones en diferentes plataformas
type NotificationService struct {
	config NotificationConfig
}

// NewFolioService crea una nueva instancia del servicio de folios
func NewFolioService(db *mongo.Database, cafService *CAFService, redisClient *redis.Client, umbralFolios int) *FolioService {
	return &FolioService{
		db:           db,
		cafService:   cafService,
		redisClient:  redisClient,
		umbralFolios: umbralFolios,
	}
}

// NewEmailAlertService crea una nueva instancia del servicio de alertas
func NewEmailAlertService(config AlertConfig) *EmailAlertService {
	return &EmailAlertService{config: config}
}

// NewNotificationService crea una nueva instancia del servicio de notificaciones
func NewNotificationService(config NotificationConfig) *NotificationService {
	return &NotificationService{config: config}
}

// VerificarVigenciaCAF verifica el estado de un CAF
func (s *FolioService) VerificarVigenciaCAF(ctx context.Context, cafID string) (*EstadoCAF, error) {
	// Obtener información del CAF
	var rango RangoFolios
	err := s.db.Collection("cafs").FindOne(ctx, bson.M{"_id": cafID}).Decode(&rango)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo información del CAF: %v", err)
	}

	// Calcular días restantes
	diasRestantes := int(rango.FechaVencimiento.Sub(time.Now()).Hours() / 24)

	// Obtener estadísticas de folios
	estadisticas, err := s.obtenerEstadisticasFolios(ctx, cafID)
	if err != nil {
		return nil, err
	}

	// Determinar estado
	estado := "VIGENTE"
	if diasRestantes <= 0 {
		estado = "EXPIRADO"
	} else if diasRestantes <= 15 {
		estado = "PROXIMO_EXPIRAR"
	}

	return &EstadoCAF{
		CAFID:             cafID,
		RUTEmisor:         rango.RUTEmisor,
		TipoDTE:           rango.TipoDTE,
		Estado:            estado,
		DiasRestantes:     diasRestantes,
		FoliosDisponibles: estadisticas.Disponibles,
		FoliosUtilizados:  estadisticas.Utilizados,
		FoliosAnulados:    estadisticas.Anulados,
	}, nil
}

// MonitorearVigenciaCAF monitorea la vigencia de los CAF y envía alertas
func (s *FolioService) MonitorearVigenciaCAF(ctx context.Context) error {
	// Obtener todos los CAF activos
	cursor, err := s.db.Collection("cafs").Find(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("error obteniendo CAF: %v", err)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var rango RangoFolios
		if err := cursor.Decode(&rango); err != nil {
			continue
		}

		// Verificar vigencia
		estado, err := s.VerificarVigenciaCAF(ctx, rango.CAFID)
		if err != nil {
			continue
		}

		// Enviar alertas según el estado
		switch estado.Estado {
		case "PROXIMO_EXPIRAR":
			s.enviarAlertaProximaExpiracion(ctx, estado)
		case "EXPIRADO":
			s.enviarAlertaExpirado(ctx, estado)
		}
	}

	return nil
}

// enviarAlertaProximaExpiracion envía una alerta cuando el CAF está próximo a expirar
func (s *FolioService) enviarAlertaProximaExpiracion(ctx context.Context, estado *EstadoCAF) error {
	// Crear el mensaje de alerta
	subject := fmt.Sprintf("ALERTA: CAF %s está próximo a expirar", estado.CAFID)
	body := fmt.Sprintf(`
		El CAF %s está próximo a expirar.
		Detalles:
		- RUT Emisor: %s
		- Tipo DTE: %s
		- Días restantes: %d
		- Folios disponibles: %d
		- Folios utilizados: %d
		- Folios anulados: %d
	`, estado.CAFID, estado.RUTEmisor, estado.TipoDTE,
		estado.DiasRestantes, estado.FoliosDisponibles,
		estado.FoliosUtilizados, estado.FoliosAnulados)

	// Enviar email
	emailService := NewEmailAlertService(AlertConfig{
		SMTPHost:     "smtp.gmail.com",
		SMTPPort:     587,
		SMTPUser:     "alertas@empresa.com",
		SMTPPassword: "password",
		FromEmail:    "alertas@empresa.com",
		ToEmails:     []string{"admin@empresa.com"},
	})

	return emailService.SendAlert(subject, body)
}

// enviarAlertaExpirado envía una alerta cuando el CAF ha expirado
func (s *FolioService) enviarAlertaExpirado(ctx context.Context, estado *EstadoCAF) {
	// TODO: Implementar envío de alerta (email, notificación, etc.)
	fmt.Printf("ALERTA: CAF %s ha expirado. Solicitar nuevo CAF.\n", estado.CAFID)
}

// EstadisticasFolios representa las estadísticas de uso de folios
type EstadisticasFolios struct {
	Disponibles int
	Utilizados  int
	Anulados    int
}

// obtenerEstadisticasFolios obtiene las estadísticas de uso de folios
func (s *FolioService) obtenerEstadisticasFolios(ctx context.Context, cafID string) (*EstadisticasFolios, error) {
	// Obtener conteos por estado
	pipeline := []bson.M{
		{"$match": bson.M{"caf_id": cafID}},
		{"$group": bson.M{
			"_id":   "$estado",
			"count": bson.M{"$sum": 1},
		}},
	}

	cursor, err := s.db.Collection("folios").Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo estadísticas: %v", err)
	}
	defer cursor.Close(ctx)

	estadisticas := &EstadisticasFolios{}
	for cursor.Next(ctx) {
		var result struct {
			ID    string `bson:"_id"`
			Count int    `bson:"count"`
		}
		if err := cursor.Decode(&result); err != nil {
			continue
		}

		switch result.ID {
		case "DISPONIBLE":
			estadisticas.Disponibles = result.Count
		case "UTILIZADO":
			estadisticas.Utilizados = result.Count
		case "ANULADO":
			estadisticas.Anulados = result.Count
		}
	}

	return estadisticas, nil
}

// RegistrarRangoFolios registra un nuevo rango de folios de un CAF
func (s *FolioService) RegistrarRangoFolios(ctx context.Context, rango RangoFolios) error {
	// Verificar superposición con rangos existentes
	if err := s.verificarSuperposicion(ctx, rango); err != nil {
		return err
	}

	// Crear folios en la base de datos
	folios := make([]interface{}, 0)
	for i := rango.FolioInicial; i <= rango.FolioFinal; i++ {
		folio := Folio{
			ID:        fmt.Sprintf("%s-%s-%d", rango.RUTEmisor, rango.TipoDTE, i),
			RUTEmisor: rango.RUTEmisor,
			TipoDTE:   rango.TipoDTE,
			Numero:    i,
			Estado:    "DISPONIBLE",
			CAFID:     rango.CAFID,
		}
		folios = append(folios, folio)
	}

	// Insertar folios en lote
	_, err := s.db.Collection("folios").InsertMany(ctx, folios)
	if err != nil {
		return fmt.Errorf("error registrando folios: %v", err)
	}

	return nil
}

// ObtenerFolioDisponible obtiene el próximo folio disponible con manejo de concurrencia
func (s *FolioService) ObtenerFolioDisponible(ctx context.Context, rutEmisor string, tipoDTE string) (*Folio, error) {
	// Crear una clave única para el bloqueo
	lockKey := fmt.Sprintf("folio_lock:%s:%s", rutEmisor, tipoDTE)

	// Intentar obtener el bloqueo
	locked, err := s.redisClient.SetNX(ctx, lockKey, "1", 30*time.Second).Result()
	if err != nil {
		return nil, fmt.Errorf("error obteniendo bloqueo: %v", err)
	}
	if !locked {
		return nil, fmt.Errorf("no se pudo obtener el bloqueo para asignar folio")
	}
	defer s.redisClient.Del(ctx, lockKey)

	// Verificar disponibilidad de folios
	disponibles, err := s.ContarFoliosDisponibles(ctx, rutEmisor, tipoDTE)
	if err != nil {
		return nil, err
	}

	// Si quedan pocos folios, solicitar nuevo CAF
	if disponibles <= s.umbralFolios {
		if err := s.solicitarNuevoCAF(ctx, rutEmisor, tipoDTE); err != nil {
			return nil, err
		}
	}

	// Obtener el próximo folio disponible usando transacción
	session, err := s.db.Client().StartSession()
	if err != nil {
		return nil, fmt.Errorf("error iniciando sesión: %v", err)
	}
	defer session.EndSession(ctx)

	var folio *Folio
	err = mongo.WithSession(ctx, session, func(sessCtx mongo.SessionContext) error {
		if err := session.StartTransaction(); err != nil {
			return err
		}

		// Obtener y actualizar el folio en una transacción
		err = s.db.Collection("folios").FindOneAndUpdate(
			sessCtx,
			bson.M{
				"rut_emisor": rutEmisor,
				"tipo_dte":   tipoDTE,
				"estado":     "DISPONIBLE",
			},
			bson.M{
				"$set": bson.M{
					"estado":    "UTILIZADO",
					"fecha_uso": time.Now(),
				},
			},
		).Decode(&folio)

		if err != nil {
			session.AbortTransaction(sessCtx)
			return fmt.Errorf("error obteniendo folio disponible: %v", err)
		}

		return session.CommitTransaction(sessCtx)
	})

	if err != nil {
		return nil, err
	}

	return folio, nil
}

// ContarFoliosDisponibles cuenta los folios disponibles
func (s *FolioService) ContarFoliosDisponibles(ctx context.Context, rutEmisor string, tipoDTE string) (int, error) {
	count, err := s.db.Collection("folios").CountDocuments(
		ctx,
		bson.M{
			"rut_emisor": rutEmisor,
			"tipo_dte":   tipoDTE,
			"estado":     "DISPONIBLE",
		},
	)
	if err != nil {
		return 0, fmt.Errorf("error contando folios disponibles: %v", err)
	}
	return int(count), nil
}

// verificarSuperposicion verifica que no haya superposición con rangos existentes
func (s *FolioService) verificarSuperposicion(ctx context.Context, rango RangoFolios) error {
	// Buscar rangos que se superpongan
	cursor, err := s.db.Collection("folios").Find(
		ctx,
		bson.M{
			"rut_emisor": rango.RUTEmisor,
			"tipo_dte":   rango.TipoDTE,
			"numero": bson.M{
				"$gte": rango.FolioInicial,
				"$lte": rango.FolioFinal,
			},
		},
	)
	if err != nil {
		return fmt.Errorf("error verificando superposición: %v", err)
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		return fmt.Errorf("existe superposición de folios")
	}

	return nil
}

// solicitarNuevoCAF solicita un nuevo CAF cuando los folios están por agotarse
func (s *FolioService) solicitarNuevoCAF(ctx context.Context, rutEmisor string, tipoDTE string) error {
	// Crear una clave única para el bloqueo de solicitud
	lockKey := fmt.Sprintf("caf_request_lock:%s:%s", rutEmisor, tipoDTE)

	// Verificar si ya hay una solicitud en proceso
	exists, err := s.redisClient.Exists(ctx, lockKey).Result()
	if err != nil {
		return fmt.Errorf("error verificando solicitud existente: %v", err)
	}
	if exists {
		return fmt.Errorf("ya existe una solicitud de CAF en proceso")
	}

	// Obtener bloqueo para la solicitud
	locked, err := s.redisClient.SetNX(ctx, lockKey, "1", 5*time.Minute).Result()
	if err != nil {
		return fmt.Errorf("error obteniendo bloqueo: %v", err)
	}
	if !locked {
		return fmt.Errorf("no se pudo obtener el bloqueo para solicitar CAF")
	}
	defer s.redisClient.Del(ctx, lockKey)

	// Obtener el último folio utilizado
	var ultimoFolio Folio
	err = s.db.Collection("folios").FindOne(
		ctx,
		bson.M{
			"rut_emisor": rutEmisor,
			"tipo_dte":   tipoDTE,
		},
		bson.M{"sort": bson.M{"numero": -1}},
	).Decode(&ultimoFolio)

	if err != nil && err != mongo.ErrNoDocuments {
		return fmt.Errorf("error obteniendo último folio: %v", err)
	}

	// Calcular rango para nuevo CAF
	folioInicial := 1
	if err == nil {
		folioInicial = ultimoFolio.Numero + 1
	}

	// Solicitar nuevo CAF
	req := &models.CAFRequest{
		TipoDocumento: tipoDTE,
		RutEmisor:     rutEmisor,
		Cantidad:      1000, // Solicitar 1000 folios
	}

	// Iniciar monitoreo en segundo plano
	go s.monitorearSolicitudCAF(context.Background(), req)

	_, err = s.cafService.SolicitarCAF(ctx, req)
	if err != nil {
		return fmt.Errorf("error solicitando nuevo CAF: %v", err)
	}

	return nil
}

// monitorearSolicitudCAF monitorea el estado de una solicitud de CAF
func (s *FolioService) monitorearSolicitudCAF(ctx context.Context, req *models.CAFRequest) {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	// Crear un ID de seguimiento local para esta solicitud
	trackID := fmt.Sprintf("%s-%s-%d", req.RutEmisor, req.TipoDocumento, time.Now().Unix())

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// Verificar estado de la solicitud
			estado, err := s.cafService.ConsultarEstadoCAF(ctx, trackID)
			if err != nil {
				continue
			}

			switch estado {
			case "ACEPTADO":
				// Descargar y registrar el nuevo CAF
				if err := s.procesarNuevoCAF(ctx, req, trackID); err != nil {
					// TODO: Implementar manejo de errores
					continue
				}
				return
			case "RECHAZADO":
				// TODO: Implementar manejo de rechazo
				return
			}
		}
	}
}

// procesarNuevoCAF procesa un nuevo CAF recibido
func (s *FolioService) procesarNuevoCAF(ctx context.Context, req *models.CAFRequest, trackID string) error {
	// Descargar el CAF
	metadata, err := s.cafService.DescargarCAF(ctx, trackID)
	if err != nil {
		return fmt.Errorf("error descargando CAF: %v", err)
	}

	// Registrar el nuevo rango de folios
	rango := RangoFolios{
		RUTEmisor:         req.RutEmisor,
		TipoDTE:           req.TipoDocumento,
		FolioInicial:      metadata.FolioInicial,
		FolioFinal:        metadata.FolioFinal,
		CAFID:             metadata.CAFID,
		FechaAutorizacion: metadata.FechaAutorizacion,
		FechaVencimiento:  metadata.FechaVencimiento,
	}

	if err := s.RegistrarRangoFolios(ctx, rango); err != nil {
		return fmt.Errorf("error registrando rango de folios: %v", err)
	}

	return nil
}

// AnularFolio marca un folio como anulado
func (s *FolioService) AnularFolio(ctx context.Context, folioID string) error {
	_, err := s.db.Collection("folios").UpdateOne(
		ctx,
		bson.M{"_id": folioID},
		bson.M{
			"$set": bson.M{
				"estado": "ANULADO",
			},
		},
	)
	if err != nil {
		return fmt.Errorf("error anulando folio: %v", err)
	}
	return nil
}

// ObtenerHistorialFolio obtiene el historial de un folio
func (s *FolioService) ObtenerHistorialFolio(ctx context.Context, folioID string) (*Folio, error) {
	var folio Folio
	err := s.db.Collection("folios").FindOne(ctx, bson.M{"_id": folioID}).Decode(&folio)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo historial de folio: %v", err)
	}
	return &folio, nil
}

// GenerarReporteFolios genera un reporte de uso de folios
func (s *FolioService) GenerarReporteFolios(ctx context.Context, rutEmisor string, tipoDTE string, periodoInicio, periodoFin time.Time) (*ReporteFolios, error) {
	// Obtener estadísticas generales
	pipeline := []bson.M{
		{"$match": bson.M{
			"rut_emisor": rutEmisor,
			"tipo_dte":   tipoDTE,
			"fecha_uso": bson.M{
				"$gte": periodoInicio,
				"$lte": periodoFin,
			},
		}},
		{"$group": bson.M{
			"_id":   "$estado",
			"count": bson.M{"$sum": 1},
		}},
	}

	cursor, err := s.db.Collection("folios").Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo estadísticas: %v", err)
	}
	defer cursor.Close(ctx)

	reporte := &ReporteFolios{
		RUTEmisor:     rutEmisor,
		TipoDTE:       tipoDTE,
		PeriodoInicio: periodoInicio,
		PeriodoFin:    periodoFin,
	}

	for cursor.Next(ctx) {
		var result struct {
			ID    string `bson:"_id"`
			Count int    `bson:"count"`
		}
		if err := cursor.Decode(&result); err != nil {
			continue
		}

		switch result.ID {
		case "DISPONIBLE":
			reporte.FoliosDisponibles = result.Count
		case "UTILIZADO":
			reporte.FoliosUtilizados = result.Count
		case "ANULADO":
			reporte.FoliosAnulados = result.Count
		}
	}

	reporte.TotalFolios = reporte.FoliosDisponibles + reporte.FoliosUtilizados + reporte.FoliosAnulados

	// Obtener detalle de uso
	detalleCursor, err := s.db.Collection("folios").Find(
		ctx,
		bson.M{
			"rut_emisor": rutEmisor,
			"tipo_dte":   tipoDTE,
			"fecha_uso": bson.M{
				"$gte": periodoInicio,
				"$lte": periodoFin,
			},
		},
		bson.M{"sort": bson.M{"numero": 1}},
	)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo detalle de uso: %v", err)
	}
	defer detalleCursor.Close(ctx)

	for detalleCursor.Next(ctx) {
		var folio Folio
		if err := detalleCursor.Decode(&folio); err != nil {
			continue
		}

		reporte.DetalleUso = append(reporte.DetalleUso, DetalleUsoFolio{
			Numero:      folio.Numero,
			Estado:      folio.Estado,
			FechaUso:    folio.FechaUso,
			DocumentoID: folio.DocumentoID,
			CAFID:       folio.CAFID,
		})
	}

	// Obtener estado de CAF activos
	cafCursor, err := s.db.Collection("cafs").Find(
		ctx,
		bson.M{
			"rut_emisor": rutEmisor,
			"tipo_dte":   tipoDTE,
			"fecha_vencimiento": bson.M{
				"$gt": time.Now(),
			},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo CAF activos: %v", err)
	}
	defer cafCursor.Close(ctx)

	for cafCursor.Next(ctx) {
		var rango RangoFolios
		if err := cafCursor.Decode(&rango); err != nil {
			continue
		}

		estado, err := s.VerificarVigenciaCAF(ctx, rango.CAFID)
		if err != nil {
			continue
		}

		reporte.CAFActivos = append(reporte.CAFActivos, *estado)
	}

	return reporte, nil
}

// ExportarReporteFolios exporta el reporte en diferentes formatos
func (s *FolioService) ExportarReporteFolios(ctx context.Context, reporte *ReporteFolios, formato string) ([]byte, error) {
	switch formato {
	case "PDF":
		return s.exportarPDF(reporte)
	case "CSV":
		return s.exportarCSV(reporte)
	case "JSON":
		return s.exportarJSON(reporte)
	default:
		return nil, fmt.Errorf("formato no soportado: %s", formato)
	}
}

// exportarPDF exporta el reporte en formato PDF
func (s *FolioService) exportarPDF(reporte *ReporteFolios) ([]byte, error) {
	// Crear nuevo documento PDF
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Configurar fuente
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Reporte de Uso de Folios")
	pdf.Ln(20)

	// Información general
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, fmt.Sprintf("RUT Emisor: %s", reporte.RUTEmisor))
	pdf.Ln(10)
	pdf.Cell(40, 10, fmt.Sprintf("Tipo DTE: %s", reporte.TipoDTE))
	pdf.Ln(10)
	pdf.Cell(40, 10, fmt.Sprintf("Periodo: %s - %s",
		reporte.PeriodoInicio.Format("02/01/2006"),
		reporte.PeriodoFin.Format("02/01/2006")))
	pdf.Ln(20)

	// Estadísticas
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 10, "Estadísticas")
	pdf.Ln(10)
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, fmt.Sprintf("Total Folios: %d", reporte.TotalFolios))
	pdf.Ln(10)
	pdf.Cell(40, 10, fmt.Sprintf("Folios Disponibles: %d", reporte.FoliosDisponibles))
	pdf.Ln(10)
	pdf.Cell(40, 10, fmt.Sprintf("Folios Utilizados: %d", reporte.FoliosUtilizados))
	pdf.Ln(10)
	pdf.Cell(40, 10, fmt.Sprintf("Folios Anulados: %d", reporte.FoliosAnulados))
	pdf.Ln(20)

	// Detalle de uso
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 10, "Detalle de Uso")
	pdf.Ln(10)
	pdf.SetFont("Arial", "", 10)

	// Encabezados de tabla
	headers := []string{"Número", "Estado", "Fecha Uso", "Documento ID", "CAF ID"}
	widths := []float64{30, 30, 40, 50, 50}

	for i, header := range headers {
		pdf.Cell(widths[i], 10, header)
	}
	pdf.Ln(10)

	// Datos de la tabla
	for _, detalle := range reporte.DetalleUso {
		pdf.Cell(widths[0], 10, fmt.Sprintf("%d", detalle.Numero))
		pdf.Cell(widths[1], 10, detalle.Estado)
		pdf.Cell(widths[2], 10, detalle.FechaUso.Format("02/01/2006"))
		pdf.Cell(widths[3], 10, detalle.DocumentoID)
		pdf.Cell(widths[4], 10, detalle.CAFID)
		pdf.Ln(10)
	}

	// CAF Activos
	pdf.Ln(20)
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 10, "CAF Activos")
	pdf.Ln(10)
	pdf.SetFont("Arial", "", 10)

	for _, caf := range reporte.CAFActivos {
		pdf.Cell(40, 10, fmt.Sprintf("CAF ID: %s", caf.CAFID))
		pdf.Ln(10)
		pdf.Cell(40, 10, fmt.Sprintf("Estado: %s", caf.Estado))
		pdf.Ln(10)
		pdf.Cell(40, 10, fmt.Sprintf("Días Restantes: %d", caf.DiasRestantes))
		pdf.Ln(10)
	}

	// Guardar PDF en buffer
	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, fmt.Errorf("error generando PDF: %v", err)
	}

	return buf.Bytes(), nil
}

// exportarCSV exporta el reporte en formato CSV
func (s *FolioService) exportarCSV(reporte *ReporteFolios) ([]byte, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Escribir encabezados
	headers := []string{
		"RUT Emisor", "Tipo DTE", "Periodo Inicio", "Periodo Fin",
		"Total Folios", "Folios Disponibles", "Folios Utilizados", "Folios Anulados",
	}
	if err := writer.Write(headers); err != nil {
		return nil, fmt.Errorf("error escribiendo encabezados CSV: %v", err)
	}

	// Escribir datos generales
	row := []string{
		reporte.RUTEmisor,
		reporte.TipoDTE,
		reporte.PeriodoInicio.Format("02/01/2006"),
		reporte.PeriodoFin.Format("02/01/2006"),
		fmt.Sprintf("%d", reporte.TotalFolios),
		fmt.Sprintf("%d", reporte.FoliosDisponibles),
		fmt.Sprintf("%d", reporte.FoliosUtilizados),
		fmt.Sprintf("%d", reporte.FoliosAnulados),
	}
	if err := writer.Write(row); err != nil {
		return nil, fmt.Errorf("error escribiendo datos CSV: %v", err)
	}

	// Escribir detalle de uso
	if err := writer.Write([]string{"", "", "", "", "", "", "", ""}); err != nil {
		return nil, fmt.Errorf("error escribiendo separador CSV: %v", err)
	}

	detailHeaders := []string{"Número", "Estado", "Fecha Uso", "Documento ID", "CAF ID"}
	if err := writer.Write(detailHeaders); err != nil {
		return nil, fmt.Errorf("error escribiendo encabezados de detalle CSV: %v", err)
	}

	for _, detalle := range reporte.DetalleUso {
		row := []string{
			fmt.Sprintf("%d", detalle.Numero),
			detalle.Estado,
			detalle.FechaUso.Format("02/01/2006"),
			detalle.DocumentoID,
			detalle.CAFID,
		}
		if err := writer.Write(row); err != nil {
			return nil, fmt.Errorf("error escribiendo detalle CSV: %v", err)
		}
	}

	// Escribir CAF activos
	if err := writer.Write([]string{"", "", "", "", "", "", "", ""}); err != nil {
		return nil, fmt.Errorf("error escribiendo separador CSV: %v", err)
	}

	cafHeaders := []string{"CAF ID", "Estado", "Días Restantes", "Folios Disponibles", "Folios Utilizados", "Folios Anulados"}
	if err := writer.Write(cafHeaders); err != nil {
		return nil, fmt.Errorf("error escribiendo encabezados de CAF CSV: %v", err)
	}

	for _, caf := range reporte.CAFActivos {
		row := []string{
			caf.CAFID,
			caf.Estado,
			fmt.Sprintf("%d", caf.DiasRestantes),
			fmt.Sprintf("%d", caf.FoliosDisponibles),
			fmt.Sprintf("%d", caf.FoliosUtilizados),
			fmt.Sprintf("%d", caf.FoliosAnulados),
		}
		if err := writer.Write(row); err != nil {
			return nil, fmt.Errorf("error escribiendo CAF CSV: %v", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, fmt.Errorf("error finalizando CSV: %v", err)
	}

	return buf.Bytes(), nil
}

// exportarJSON exporta el reporte en formato JSON
func (s *FolioService) exportarJSON(reporte *ReporteFolios) ([]byte, error) {
	return json.MarshalIndent(reporte, "", "  ")
}

// SendAlert envía una alerta por email
func (s *EmailAlertService) SendAlert(subject, body string) error {
	// Configurar autenticación
	auth := smtp.PlainAuth("", s.config.SMTPUser, s.config.SMTPPassword, s.config.SMTPHost)

	// Construir mensaje
	msg := []byte(fmt.Sprintf("Subject: %s\r\n\r\n%s", subject, body))

	// Enviar email
	return smtp.SendMail(
		fmt.Sprintf("%s:%d", s.config.SMTPHost, s.config.SMTPPort),
		auth,
		s.config.FromEmail,
		s.config.ToEmails,
		msg,
	)
}

// SendNotification envía una notificación a las plataformas configuradas
func (s *NotificationService) SendNotification(title, message string) error {
	var errs []error

	// Enviar a Slack
	if s.config.SlackWebhookURL != "" {
		if err := s.sendToSlack(title, message); err != nil {
			errs = append(errs, fmt.Errorf("error enviando a Slack: %v", err))
		}
	}

	// Enviar a Teams
	if s.config.TeamsWebhookURL != "" {
		if err := s.sendToTeams(title, message); err != nil {
			errs = append(errs, fmt.Errorf("error enviando a Teams: %v", err))
		}
	}

	// Enviar a Telegram
	if s.config.TelegramToken != "" && s.config.TelegramChatID != "" {
		if err := s.sendToTelegram(title, message); err != nil {
			errs = append(errs, fmt.Errorf("error enviando a Telegram: %v", err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("errores enviando notificaciones: %v", errs)
	}

	return nil
}

// sendToSlack envía una notificación a Slack
func (s *NotificationService) sendToSlack(title, message string) error {
	payload := map[string]interface{}{
		"blocks": []map[string]interface{}{
			{
				"type": "header",
				"text": map[string]interface{}{
					"type": "plain_text",
					"text": title,
				},
			},
			{
				"type": "section",
				"text": map[string]interface{}{
					"type": "mrkdwn",
					"text": message,
				},
			},
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(s.config.SlackWebhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error enviando a Slack: %s", resp.Status)
	}

	return nil
}

// sendToTeams envía una notificación a Microsoft Teams
func (s *NotificationService) sendToTeams(title, message string) error {
	payload := map[string]interface{}{
		"@type":      "MessageCard",
		"@context":   "http://schema.org/extensions",
		"themeColor": "0076D7",
		"summary":    title,
		"sections": []map[string]interface{}{
			{
				"activityTitle": title,
				"text":          message,
			},
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(s.config.TeamsWebhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error enviando a Teams: %s", resp.Status)
	}

	return nil
}

// sendToTelegram envía una notificación a Telegram
func (s *NotificationService) sendToTelegram(title, message string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", s.config.TelegramToken)

	payload := map[string]interface{}{
		"chat_id":    s.config.TelegramChatID,
		"text":       fmt.Sprintf("*%s*\n\n%s", title, message),
		"parse_mode": "Markdown",
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error enviando a Telegram: %s", resp.Status)
	}

	return nil
}

// generarGraficos genera gráficos para el reporte
func (s *FolioService) generarGraficos(reporte *ReporteFolios, config ChartConfig) error {
	// Crear directorio de salida si no existe
	if err := os.MkdirAll(config.OutputDir, 0755); err != nil {
		return fmt.Errorf("error creando directorio de salida: %v", err)
	}

	// Gráfico de distribución de folios
	if err := s.generarGraficoDistribucion(reporte, config); err != nil {
		return err
	}

	// Gráfico de uso por día
	if err := s.generarGraficoUsoDiario(reporte, config); err != nil {
		return err
	}

	// Gráfico de estado de CAF
	if err := s.generarGraficoEstadoCAF(reporte, config); err != nil {
		return err
	}

	return nil
}

// generarGraficoDistribucion genera un gráfico de distribución de folios
func (s *FolioService) generarGraficoDistribucion(reporte *ReporteFolios, config ChartConfig) error {
	// Crear gráfico de torta
	pie := chart.PieChart{
		Width:  config.Width,
		Height: config.Height,
		Values: []chart.Value{
			{Value: float64(reporte.FoliosDisponibles), Label: "Disponibles"},
			{Value: float64(reporte.FoliosUtilizados), Label: "Utilizados"},
			{Value: float64(reporte.FoliosAnulados), Label: "Anulados"},
		},
	}

	// Guardar gráfico
	file, err := os.Create(filepath.Join(config.OutputDir, "distribucion.png"))
	if err != nil {
		return err
	}
	defer file.Close()

	return pie.Render(chart.PNG, file)
}

// generarGraficoUsoDiario genera un gráfico de uso diario de folios
func (s *FolioService) generarGraficoUsoDiario(reporte *ReporteFolios, config ChartConfig) error {
	// Agrupar uso por día
	usoDiario := make(map[string]int)
	for _, detalle := range reporte.DetalleUso {
		if detalle.Estado == "UTILIZADO" {
			fecha := detalle.FechaUso.Format("02/01")
			usoDiario[fecha]++
		}
	}

	// Crear gráfico de barras
	bar := chart.BarChart{
		Width:  config.Width,
		Height: config.Height,
		Bars:   make([]chart.Value, 0, len(usoDiario)),
	}

	for fecha, count := range usoDiario {
		bar.Bars = append(bar.Bars, chart.Value{
			Value: float64(count),
			Label: fecha,
		})
	}

	// Guardar gráfico
	file, err := os.Create(filepath.Join(config.OutputDir, "uso_diario.png"))
	if err != nil {
		return err
	}
	defer file.Close()

	return bar.Render(chart.PNG, file)
}

// generarGraficoEstadoCAF genera un gráfico del estado de los CAF
func (s *FolioService) generarGraficoEstadoCAF(reporte *ReporteFolios, config ChartConfig) error {
	// Crear gráfico de barras horizontales
	bar := chart.BarChart{
		Width:  config.Width,
		Height: config.Height,
		Bars:   make([]chart.Value, 0, len(reporte.CAFActivos)),
	}

	for _, caf := range reporte.CAFActivos {
		bar.Bars = append(bar.Bars, chart.Value{
			Value: float64(caf.DiasRestantes),
			Label: fmt.Sprintf("%s (%d días)", caf.CAFID, caf.DiasRestantes),
		})
	}

	// Guardar gráfico
	file, err := os.Create(filepath.Join(config.OutputDir, "estado_caf.png"))
	if err != nil {
		return err
	}
	defer file.Close()

	return bar.Render(chart.PNG, file)
}
