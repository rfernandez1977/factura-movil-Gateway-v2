package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/cursor/FMgo/models"
	"github.com/cursor/FMgo/services"
	"github.com/cursor/FMgo/utils"
	"go.uber.org/zap"
)

// BoletaController maneja las peticiones HTTP relacionadas con boletas
type BoletaController struct {
	boletaService *services.BoletaService
	pdfService    *services.PDFService
	emailService  *services.EmailService
}

// NewBoletaController crea una nueva instancia del controlador de boletas
func NewBoletaController(boletaService *services.BoletaService) *BoletaController {
	return &BoletaController{
		boletaService: boletaService,
		pdfService:    services.NewPDFService(),
		emailService:  services.NewEmailService(),
	}
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
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de boleta es requerido"})
		return
	}

	// TODO: Implementar obtención de boleta desde base de datos
	ctx.JSON(http.StatusNotImplemented, gin.H{"error": "Función no implementada"})
}

// ListarBoletas maneja la obtención de una lista de boletas
func (c *BoletaController) ListarBoletas(ctx *gin.Context) {
	// Obtener parámetros de consulta
	rutEmisor := ctx.Query("rut_emisor")
	startDateStr := ctx.Query("start_date")
	endDateStr := ctx.Query("end_date")

	if rutEmisor == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "RUT emisor es requerido"})
		return
	}

	// Parsear fechas
	// var startDate, endDate time.Time
	var err error

	if startDateStr != "" {
		_, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Formato de fecha inicial inválido"})
			return
		}
	} else {
		// startDate = time.Now().AddDate(0, -1, 0) // Último mes por defecto
	}

	if endDateStr != "" {
		_, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Formato de fecha final inválido"})
			return
		}
	} else {
		// endDate = time.Now()
	}

	// TODO: Implementar listado de boletas desde base de datos
	ctx.JSON(http.StatusNotImplemented, gin.H{"error": "Función no implementada"})
}

// AnularBoleta maneja la anulación de una boleta
func (c *BoletaController) AnularBoleta(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de boleta es requerido"})
		return
	}

	// TODO: Implementar anulación de boleta
	ctx.JSON(http.StatusNotImplemented, gin.H{"error": "Función no implementada"})
}

// ReenviarBoleta maneja el reenvío de una boleta
func (c *BoletaController) ReenviarBoleta(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de boleta es requerido"})
		return
	}

	// TODO: Implementar reenvío de boleta
	ctx.JSON(http.StatusNotImplemented, gin.H{"error": "Función no implementada"})
}

// DescargarPDF maneja la descarga de una boleta en formato PDF
func (c *BoletaController) DescargarPDF(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de boleta es requerido"})
		return
	}

	// Obtener boleta y detalles
	boleta, err := c.boletaService.GetBoleta(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	detalles, err := c.boletaService.GetDetalles(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Generar PDF
	pdfData, err := c.pdfService.GenerarBoletaPDF(boleta, detalles)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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

	// Enviar PDF
	ctx.Data(http.StatusOK, "application/pdf", pdfData)
}

// EnviarPorEmail maneja el envío de una boleta por email
func (c *BoletaController) EnviarPorEmail(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de boleta es requerido"})
		return
	}

	var request struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Obtener boleta y detalles
	boleta, err := c.boletaService.GetBoleta(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	detalles, err := c.boletaService.GetDetalles(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Generar PDF
	pdfData, err := c.pdfService.GenerarBoletaPDF(boleta, detalles)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Preparar mensaje
	asunto := fmt.Sprintf("Boleta Electrónica N° %d", boleta.Folio)
	mensaje := fmt.Sprintf("Adjunto encontrará la boleta electrónica N° %d.\n\nFecha de emisión: %s\nMonto total: $%d",
		boleta.Folio,
		boleta.FechaEmision.Format("02/01/2006"),
		boleta.MontoTotal)

	// Enviar email
	if err := c.emailService.EnviarBoletaPDF(request.Email, asunto, mensaje, pdfData); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Boleta enviada por email exitosamente"})
}
