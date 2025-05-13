package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/cursor/FMgo/models"
	"github.com/cursor/FMgo/services"
	"github.com/cursor/FMgo/utils"
	"go.uber.org/zap"
)

// FacturaController maneja las peticiones HTTP relacionadas con facturas
type FacturaController struct {
	facturaService *services.FacturaService
	pdfService     *services.PDFService
	emailService   *services.EmailService
}

// NewFacturaController crea una nueva instancia del controlador de facturas
func NewFacturaController(facturaService *services.FacturaService) *FacturaController {
	return &FacturaController{
		facturaService: facturaService,
		pdfService:     services.NewPDFService(),
		emailService:   services.NewEmailService(),
	}
}

// CrearFactura maneja la creación de una nueva factura
func (c *FacturaController) CrearFactura(ctx *gin.Context) {
	start := time.Now()

	var request models.FacturaRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		utils.LogError(err, zap.String("endpoint", "CrearFactura"))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validar RUT del emisor
	jwtUtils := utils.NewJWTUtils()
	rutEmisor, _ := utils.GetRut(ctx.GetHeader("Authorization"), jwtUtils)
	if rutEmisor != request.RutEmisor {
		utils.LogWarning("intento de crear factura con RUT emisor diferente",
			zap.String("rut_token", rutEmisor),
			zap.String("rut_request", request.RutEmisor),
		)
		ctx.JSON(http.StatusForbidden, gin.H{"error": "RUT emisor no coincide"})
		return
	}

	// Crear factura
	response, err := c.facturaService.CrearFactura(&request)
	if err != nil {
		utils.LogError(err, zap.String("endpoint", "CrearFactura"))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Registrar éxito
	utils.LogInfo("factura creada exitosamente",
		zap.String("id", response.Factura.ID),
		zap.String("track_id", response.Factura.TrackID),
	)

	// Registrar métricas HTTP
	duration := time.Since(start).Seconds()
	utils.RecordHTTPRequest(
		ctx.Request.Method,
		ctx.Request.URL.Path,
		http.StatusOK,
		duration,
		float64(ctx.Request.ContentLength),
		float64(len(response.Factura.ID)),
	)

	ctx.JSON(http.StatusOK, response)
}

// ConsultarEstadoFactura maneja la consulta del estado de una factura
func (c *FacturaController) ConsultarEstadoFactura(ctx *gin.Context) {
	trackID := ctx.Param("trackID")
	rutEmisor := ctx.Param("rutEmisor")

	if trackID == "" || rutEmisor == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "TrackID y RUT emisor son requeridos"})
		return
	}

	estado, err := c.facturaService.ConsultarEstadoFactura(trackID, rutEmisor)
	if err != nil {
		utils.LogError(err, zap.String("endpoint", "ConsultarEstadoFactura"))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, estado)
}

// GetFactura maneja la obtención de una factura por ID
func (c *FacturaController) GetFactura(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de factura es requerido"})
		return
	}

	response, err := c.facturaService.GetFactura(id)
	if err != nil {
		utils.LogError(err, zap.String("endpoint", "GetFactura"))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// ListarFacturas maneja la obtención de una lista de facturas
func (c *FacturaController) ListarFacturas(ctx *gin.Context) {
	rutEmisor := ctx.Query("rut_emisor")
	startDateStr := ctx.Query("start_date")
	endDateStr := ctx.Query("end_date")

	if rutEmisor == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "RUT emisor es requerido"})
		return
	}

	// Parsear fechas
	var startDate, endDate time.Time
	var err error

	if startDateStr != "" {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Formato de fecha inicial inválido"})
			return
		}
	}

	if endDateStr != "" {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Formato de fecha final inválido"})
			return
		}
	}

	facturas, err := c.facturaService.ListarFacturas(rutEmisor, startDate, endDate)
	if err != nil {
		utils.LogError(err, zap.String("endpoint", "ListarFacturas"))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, facturas)
}

// AnularFactura maneja la anulación de una factura
func (c *FacturaController) AnularFactura(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de factura es requerido"})
		return
	}

	if err := c.facturaService.AnularFactura(id); err != nil {
		utils.LogError(err, zap.String("endpoint", "AnularFactura"))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.LogInfo("factura anulada exitosamente", zap.String("id", id))
	ctx.JSON(http.StatusOK, gin.H{"message": "Factura anulada exitosamente"})
}

// ReenviarFactura maneja el reenvío de una factura
func (c *FacturaController) ReenviarFactura(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de factura es requerido"})
		return
	}

	if err := c.facturaService.ReenviarFactura(id); err != nil {
		utils.LogError(err, zap.String("endpoint", "ReenviarFactura"))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.LogInfo("factura reenviada exitosamente", zap.String("id", id))
	ctx.JSON(http.StatusOK, gin.H{"message": "Factura reenviada exitosamente"})
}

// DescargarPDF maneja la descarga de una factura en formato PDF
func (c *FacturaController) DescargarPDF(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de factura es requerido"})
		return
	}

	response, err := c.facturaService.GetFactura(id)
	if err != nil {
		utils.LogError(err, zap.String("endpoint", "DescargarPDF"))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Configurar headers para descarga
	filename := "factura_" + response.Factura.ID + ".pdf"
	ctx.Header("Content-Description", "File Transfer")
	ctx.Header("Content-Disposition", "attachment; filename="+filename)
	ctx.Header("Content-Type", "application/pdf")
	ctx.Header("Content-Transfer-Encoding", "binary")
	ctx.Header("Expires", "0")
	ctx.Header("Cache-Control", "must-revalidate")
	ctx.Header("Pragma", "public")
	ctx.Header("Content-Length", string(len(response.Factura.PDF)))

	ctx.Data(http.StatusOK, "application/pdf", response.Factura.PDF)
}

// EnviarPorEmail maneja el envío de una factura por email
func (c *FacturaController) EnviarPorEmail(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de factura es requerido"})
		return
	}

	var request struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := c.facturaService.GetFactura(id)
	if err != nil {
		utils.LogError(err, zap.String("endpoint", "EnviarPorEmail"))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Preparar mensaje
	asunto := "Factura Electrónica N° " + response.Factura.ID
	mensaje := "Adjunto encontrará la factura electrónica N° " + response.Factura.ID

	// Enviar email
	if err := c.emailService.EnviarFacturaPDF(request.Email, asunto, mensaje, response.Factura.PDF); err != nil {
		utils.LogError(err, zap.String("endpoint", "EnviarPorEmail"))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.LogInfo("factura enviada por email exitosamente",
		zap.String("id", id),
		zap.String("email", request.Email),
	)
	ctx.JSON(http.StatusOK, gin.H{"message": "Factura enviada por email exitosamente"})
}
