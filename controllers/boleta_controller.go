package controllers

import (
	"fmt"
	"net/http"
	"time"

	"FMgo/config"
	"FMgo/models"
	"FMgo/services"
	"FMgo/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// BoletaController maneja las peticiones HTTP relacionadas con boletas
type BoletaController struct {
	boletaService *services.BoletaService
	pdfService    *services.PDFService
	emailService  *services.EmailService
}

// NewBoletaController crea una nueva instancia del controlador de boletas
func NewBoletaController(boletaService *services.BoletaService, config *config.SupabaseConfig) *BoletaController {
	// Inicialización de servicios auxiliares utilizando configuración
	pdfConfig := config.GetPDFConfig()
	emailConfig := config.GetEmailConfig()

	pdfService := services.NewPDFService(
		config,
		pdfConfig.TemplatePath,
		pdfConfig.TempPath,
	)

	emailService := services.NewEmailService(
		config,
		emailConfig.SMTPServer,
		emailConfig.SMTPPort,
		emailConfig.SMTPUser,
		emailConfig.SMTPPassword,
		emailConfig.FromEmail,
		emailConfig.FromName,
	)

	return &BoletaController{
		boletaService: boletaService,
		pdfService:    pdfService,
		emailService:  emailService,
	}
}

// verificarPermisoBoleta verifica si el usuario tiene permisos para acceder a una boleta
// Retorna la boleta si tiene permisos, o nil y un error si no tiene permiso o hay algún problema
func (c *BoletaController) verificarPermisoBoleta(ctx *gin.Context, id string, endpoint string) (*models.Boleta, int, error) {
	// Verificar permisos
	jwtUtils := utils.NewJWTUtils()
	rutEmisor, err := utils.GetRut(ctx.GetHeader("Authorization"), jwtUtils)
	if err != nil {
		utils.LogError(err, zap.String("endpoint", endpoint))
		return nil, http.StatusUnauthorized, fmt.Errorf("error al verificar autenticación: %w", err)
	}

	// Obtener boleta para verificar que pertenezca al emisor
	boleta, err := c.boletaService.GetBoleta(id)
	if err != nil {
		utils.LogError(err, zap.String("endpoint", endpoint), zap.String("id", id))
		return nil, http.StatusInternalServerError, err
	}

	if boleta == nil {
		utils.LogWarning("boleta no encontrada", zap.String("id", id))
		return nil, http.StatusNotFound, fmt.Errorf("boleta no encontrada")
	}

	// Verificar que la boleta pertenezca al emisor
	if boleta.RUTEmisor != rutEmisor {
		utils.LogWarning("intento de acceder a boleta de otro emisor",
			zap.String("rut_token", rutEmisor),
			zap.String("rut_boleta", boleta.RUTEmisor),
		)
		return nil, http.StatusForbidden, fmt.Errorf("no tiene permisos para acceder a esta boleta")
	}

	return boleta, http.StatusOK, nil
}

// CrearBoleta maneja la creación de una nueva boleta
func (c *BoletaController) CrearBoleta(ctx *gin.Context) {
	start := time.Now()

	var request models.BoletaRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		utils.LogError(err, zap.String("endpoint", "CrearBoleta"))
		utils.RecordBoletaError()
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validar RUT del emisor
	jwtUtils := utils.NewJWTUtils()
	rutEmisor, _ := utils.GetRut(ctx.GetHeader("Authorization"), jwtUtils)
	if rutEmisor != request.RutEmisor {
		utils.LogWarning("intento de crear boleta con RUT emisor diferente",
			zap.String("rut_token", rutEmisor),
			zap.String("rut_request", request.RutEmisor),
		)
		ctx.JSON(http.StatusForbidden, gin.H{"error": "RUT emisor no coincide"})
		return
	}

	// Crear boleta
	boleta, err := c.boletaService.CrearBoleta(&request)
	if err != nil {
		utils.LogError(err, zap.String("endpoint", "CrearBoleta"))
		utils.RecordBoletaError()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Registrar éxito
	utils.LogInfo("boleta creada exitosamente",
		zap.String("id", boleta.ID),
		zap.String("track_id", boleta.TrackID),
	)
	utils.RecordBoletaCreada()

	// Registrar métricas HTTP
	duration := time.Since(start).Seconds()
	utils.RecordHTTPRequest(
		ctx.Request.Method,
		ctx.Request.URL.Path,
		http.StatusOK,
		duration,
		float64(ctx.Request.ContentLength),
		float64(len(boleta.ID)),
	)

	ctx.JSON(http.StatusOK, boleta)
}

// ConsultarEstadoBoleta maneja la consulta del estado de una boleta
func (c *BoletaController) ConsultarEstadoBoleta(ctx *gin.Context) {
	trackID := ctx.Param("trackID")
	rutEmisor := ctx.Param("rutEmisor")

	if trackID == "" || rutEmisor == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "TrackID y RUT emisor son requeridos"})
		return
	}

	estado, err := c.boletaService.ConsultarEstadoBoleta(trackID, rutEmisor)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, estado)
}

// GetBoleta maneja la obtención de una boleta por ID
func (c *BoletaController) GetBoleta(ctx *gin.Context) {
	start := time.Now()
	id := ctx.Param("id")
	if id == "" {
		utils.LogWarning("intento de obtener boleta sin proporcionar ID",
			zap.String("endpoint", "GetBoleta"),
		)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de boleta es requerido"})
		return
	}

	// Obtener boleta desde el servicio
	boleta, err := c.boletaService.GetBoleta(id)
	if err != nil {
		utils.LogError(err, zap.String("endpoint", "GetBoleta"), zap.String("id", id))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Si la boleta no existe, el servicio debería retornar nil y un error específico
	if boleta == nil {
		utils.LogWarning("boleta no encontrada", zap.String("id", id))
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Boleta no encontrada"})
		return
	}

	// Obtener detalles de la boleta
	detalles, err := c.boletaService.GetDetalles(id)
	if err != nil {
		utils.LogError(err, zap.String("endpoint", "GetBoleta"), zap.String("id", id))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener detalles de la boleta"})
		return
	}

	// Crear respuesta con boleta y detalles
	response := gin.H{
		"boleta":   boleta,
		"detalles": detalles,
	}

	// Registrar métricas HTTP
	duration := time.Since(start).Seconds()
	utils.RecordHTTPRequest(
		ctx.Request.Method,
		ctx.Request.URL.Path,
		http.StatusOK,
		duration,
		0,
		float64(len(detalles)),
	)

	ctx.JSON(http.StatusOK, response)
}

// ListarBoletas maneja la obtención de una lista de boletas
func (c *BoletaController) ListarBoletas(ctx *gin.Context) {
	start := time.Now()

	// Obtener parámetros de consulta
	rutEmisor := ctx.Query("rut_emisor")
	startDateStr := ctx.Query("start_date")
	endDateStr := ctx.Query("end_date")
	limit := 100 // Valor por defecto

	if rutEmisor == "" {
		utils.LogWarning("intento de listar boletas sin RUT emisor",
			zap.String("endpoint", "ListarBoletas"),
		)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "RUT emisor es requerido"})
		return
	}

	// Parsear fechas
	var startDate, endDate time.Time
	var err error

	if startDateStr != "" {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			utils.LogWarning("formato de fecha inicial inválido",
				zap.String("endpoint", "ListarBoletas"),
				zap.String("start_date", startDateStr),
			)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Formato de fecha inicial inválido"})
			return
		}
	} else {
		startDate = time.Now().AddDate(0, -1, 0) // Último mes por defecto
	}

	if endDateStr != "" {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			utils.LogWarning("formato de fecha final inválido",
				zap.String("endpoint", "ListarBoletas"),
				zap.String("end_date", endDateStr),
			)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Formato de fecha final inválido"})
			return
		}
	} else {
		endDate = time.Now() // Fecha actual por defecto
	}

	// Obtener boletas desde el servicio
	boletas, err := c.boletaService.ListarBoletas(rutEmisor, startDate, endDate, limit)
	if err != nil {
		utils.LogError(err,
			zap.String("endpoint", "ListarBoletas"),
			zap.String("rut_emisor", rutEmisor),
		)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Registrar métricas HTTP
	duration := time.Since(start).Seconds()
	utils.RecordHTTPRequest(
		ctx.Request.Method,
		ctx.Request.URL.Path,
		http.StatusOK,
		duration,
		0,
		float64(len(boletas)),
	)

	ctx.JSON(http.StatusOK, gin.H{
		"boletas": boletas,
		"total":   len(boletas),
		"filtros": gin.H{
			"rut_emisor":   rutEmisor,
			"fecha_inicio": startDate.Format("2006-01-02"),
			"fecha_fin":    endDate.Format("2006-01-02"),
		},
	})
}

// AnularBoleta maneja la anulación de una boleta
func (c *BoletaController) AnularBoleta(ctx *gin.Context) {
	start := time.Now()

	id := ctx.Param("id")
	if id == "" {
		utils.LogWarning("intento de anular boleta sin proporcionar ID",
			zap.String("endpoint", "AnularBoleta"),
		)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de boleta es requerido"})
		return
	}

	var request struct {
		Motivo string `json:"motivo" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		utils.LogWarning("datos inválidos para anulación de boleta",
			zap.String("endpoint", "AnularBoleta"),
			zap.Error(err),
		)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Motivo de anulación es requerido"})
		return
	}

	// Verificar permisos usando el método centralizado
	// No necesitamos la boleta, solo verificar permisos
	if _, status, err := c.verificarPermisoBoleta(ctx, id, "AnularBoleta"); err != nil {
		ctx.JSON(status, gin.H{"error": err.Error()})
		return
	}

	// Anular boleta
	if err := c.boletaService.AnularBoleta(id, request.Motivo); err != nil {
		utils.LogError(err, zap.String("endpoint", "AnularBoleta"), zap.String("id", id))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Registrar éxito
	utils.LogInfo("boleta anulada exitosamente",
		zap.String("id", id),
		zap.String("motivo", request.Motivo),
	)

	// Registrar métricas HTTP
	duration := time.Since(start).Seconds()
	utils.RecordHTTPRequest(
		ctx.Request.Method,
		ctx.Request.URL.Path,
		http.StatusOK,
		duration,
		float64(ctx.Request.ContentLength),
		0,
	)

	ctx.JSON(http.StatusOK, gin.H{"message": "Boleta anulada exitosamente"})
}

// ReenviarBoleta maneja el reenvío de una boleta
func (c *BoletaController) ReenviarBoleta(ctx *gin.Context) {
	start := time.Now()

	id := ctx.Param("id")
	if id == "" {
		utils.LogWarning("intento de reenviar boleta sin proporcionar ID",
			zap.String("endpoint", "ReenviarBoleta"),
		)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de boleta es requerido"})
		return
	}

	// Verificar permisos usando el método centralizado
	// No necesitamos la boleta, solo verificar permisos
	if _, status, err := c.verificarPermisoBoleta(ctx, id, "ReenviarBoleta"); err != nil {
		ctx.JSON(status, gin.H{"error": err.Error()})
		return
	}

	// Reenviar boleta
	if err := c.boletaService.ReenviarBoleta(id); err != nil {
		utils.LogError(err, zap.String("endpoint", "ReenviarBoleta"), zap.String("id", id))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Registrar éxito
	utils.LogInfo("boleta reenviada exitosamente",
		zap.String("id", id),
	)

	// Registrar métricas HTTP
	duration := time.Since(start).Seconds()
	utils.RecordHTTPRequest(
		ctx.Request.Method,
		ctx.Request.URL.Path,
		http.StatusOK,
		duration,
		0,
		0,
	)

	ctx.JSON(http.StatusOK, gin.H{"message": "Boleta reenviada exitosamente"})
}

// DescargarPDF maneja la descarga de una boleta en formato PDF
func (c *BoletaController) DescargarPDF(ctx *gin.Context) {
	start := time.Now()

	id := ctx.Param("id")
	if id == "" {
		utils.LogWarning("intento de descargar PDF sin proporcionar ID",
			zap.String("endpoint", "DescargarPDF"),
		)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de boleta es requerido"})
		return
	}

	// Obtener boleta para verificar que exista
	boleta, err := c.boletaService.GetBoleta(id)
	if err != nil {
		utils.LogError(err, zap.String("endpoint", "DescargarPDF"), zap.String("id", id))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if boleta == nil {
		utils.LogWarning("boleta no encontrada", zap.String("id", id))
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Boleta no encontrada"})
		return
	}

	// Obtener PDF
	pdfData, err := c.pdfService.ObtenerPDF(id)
	if err != nil {
		// Si no existe el PDF ya generado, podríamos implementar la generación aquí
		utils.LogError(err, zap.String("endpoint", "DescargarPDF"), zap.String("id", id))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo obtener el PDF de la boleta"})
		return
	}

	// Configurar headers para descarga
	filename := fmt.Sprintf("boleta_%d.pdf", boleta.Folio)
	ctx.Header("Content-Description", "File Transfer")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	ctx.Header("Content-Type", "application/pdf")
	ctx.Header("Content-Transfer-Encoding", "binary")
	ctx.Header("Expires", "0")
	ctx.Header("Cache-Control", "must-revalidate")
	ctx.Header("Pragma", "public")
	ctx.Header("Content-Length", fmt.Sprintf("%d", len(pdfData)))

	// Registrar métricas HTTP
	duration := time.Since(start).Seconds()
	utils.RecordHTTPRequest(
		ctx.Request.Method,
		ctx.Request.URL.Path,
		http.StatusOK,
		duration,
		0,
		float64(len(pdfData)),
	)

	// Enviar PDF
	ctx.Data(http.StatusOK, "application/pdf", pdfData)
}

// EnviarPorEmail maneja el envío de una boleta por email
func (c *BoletaController) EnviarPorEmail(ctx *gin.Context) {
	start := time.Now()

	id := ctx.Param("id")
	if id == "" {
		utils.LogWarning("intento de enviar boleta por email sin proporcionar ID",
			zap.String("endpoint", "EnviarPorEmail"),
		)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de boleta es requerido"})
		return
	}

	var request struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		utils.LogWarning("email inválido para enviar boleta",
			zap.String("endpoint", "EnviarPorEmail"),
			zap.Error(err),
		)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verificar permisos usando el método centralizado
	boleta, status, err := c.verificarPermisoBoleta(ctx, id, "EnviarPorEmail")
	if err != nil {
		ctx.JSON(status, gin.H{"error": err.Error()})
		return
	}

	// Obtener PDF de la boleta
	pdfData, err := c.pdfService.ObtenerPDF(id)
	if err != nil {
		utils.LogError(err, zap.String("endpoint", "EnviarPorEmail"), zap.String("id", id))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo obtener el PDF de la boleta"})
		return
	}

	// Usar el método EnviarDocumento que sí existe en el servicio de email
	err = c.emailService.EnviarDocumento(request.Email, boleta.RazonSocialReceptor, boleta, pdfData, nil)
	if err != nil {
		utils.LogError(err, zap.String("endpoint", "EnviarPorEmail"), zap.String("id", id))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Registrar éxito
	utils.LogInfo("boleta enviada por email exitosamente",
		zap.String("id", id),
		zap.String("email", request.Email),
	)

	// Registrar métricas HTTP
	duration := time.Since(start).Seconds()
	utils.RecordHTTPRequest(
		ctx.Request.Method,
		ctx.Request.URL.Path,
		http.StatusOK,
		duration,
		float64(ctx.Request.ContentLength),
		0,
	)

	ctx.JSON(http.StatusOK, gin.H{"message": "Boleta enviada por email exitosamente"})
}
