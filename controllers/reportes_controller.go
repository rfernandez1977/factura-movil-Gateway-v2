package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/fmgo/services"
)

// ReportesController maneja las peticiones relacionadas con reportes
type ReportesController struct {
	reportesService *services.ReportesService
}

// NewReportesController crea una nueva instancia del controlador de reportes
func NewReportesController(reportesService *services.ReportesService) *ReportesController {
	return &ReportesController{
		reportesService: reportesService,
	}
}

// GenerarReporteDocumentosEstado genera un reporte de documentos por estado
func (c *ReportesController) GenerarReporteDocumentosEstado(ctx *gin.Context) {
	var request struct {
		FechaInicio time.Time `json:"fecha_inicio" binding:"required"`
		FechaFin    time.Time `json:"fecha_fin" binding:"required"`
		RutEmisor   string    `json:"rut_emisor"`
		RutReceptor string    `json:"rut_receptor"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reporte, err := c.reportesService.GenerarReporteDocumentosEstado(
		context.Background(),
		request.FechaInicio,
		request.FechaFin,
		request.RutEmisor,
		request.RutReceptor,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, reporte)
}

// GenerarReporteRechazos genera un reporte de análisis de rechazos
func (c *ReportesController) GenerarReporteRechazos(ctx *gin.Context) {
	var request struct {
		FechaInicio time.Time `json:"fecha_inicio" binding:"required"`
		FechaFin    time.Time `json:"fecha_fin" binding:"required"`
		RutEmisor   string    `json:"rut_emisor"`
		RutReceptor string    `json:"rut_receptor"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reporte, err := c.reportesService.GenerarReporteRechazos(
		context.Background(),
		request.FechaInicio,
		request.FechaFin,
		request.RutEmisor,
		request.RutReceptor,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, reporte)
}

// GenerarReporteMetricasRendimiento genera un reporte de métricas de rendimiento
func (c *ReportesController) GenerarReporteMetricasRendimiento(ctx *gin.Context) {
	var request struct {
		FechaInicio time.Time `json:"fecha_inicio" binding:"required"`
		FechaFin    time.Time `json:"fecha_fin" binding:"required"`
		RutEmisor   string    `json:"rut_emisor"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reporte, err := c.reportesService.GenerarReporteMetricasRendimiento(
		context.Background(),
		request.FechaInicio,
		request.FechaFin,
		request.RutEmisor,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, reporte)
}

// GenerarReporteTributario genera un reporte tributario
func (c *ReportesController) GenerarReporteTributario(ctx *gin.Context) {
	var request struct {
		FechaInicio time.Time `json:"fecha_inicio" binding:"required"`
		FechaFin    time.Time `json:"fecha_fin" binding:"required"`
		RutEmisor   string    `json:"rut_emisor"`
		RutReceptor string    `json:"rut_receptor"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reporte, err := c.reportesService.GenerarReporteTributario(
		context.Background(),
		request.FechaInicio,
		request.FechaFin,
		request.RutEmisor,
		request.RutReceptor,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, reporte)
}

// ObtenerReporte obtiene un reporte por su ID
func (c *ReportesController) ObtenerReporte(ctx *gin.Context) {
	id := ctx.Param("id")
	tipo := ctx.Param("tipo")

	reporte, err := c.reportesService.ObtenerReporte(context.Background(), id, tipo)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, reporte)
}

// ListarReportes obtiene una lista de reportes
func (c *ReportesController) ListarReportes(ctx *gin.Context) {
	tipo := ctx.Param("tipo")
	fechaInicio := ctx.Query("fecha_inicio")
	fechaFin := ctx.Query("fecha_fin")
	rutEmisor := ctx.Query("rut_emisor")

	// Parsear fechas
	inicio, err := time.Parse(time.RFC3339, fechaInicio)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "fecha_inicio inválida"})
		return
	}

	fin, err := time.Parse(time.RFC3339, fechaFin)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "fecha_fin inválida"})
		return
	}

	reportes, err := c.reportesService.ListarReportes(
		context.Background(),
		tipo,
		inicio,
		fin,
		rutEmisor,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, reportes)
}

// RegisterRoutes registra las rutas del controlador
func (c *ReportesController) RegisterRoutes(router *gin.RouterGroup) {
	reportes := router.Group("/reportes")
	{
		reportes.POST("/estado", c.GenerarReporteDocumentosEstado)
		reportes.POST("/rechazos", c.GenerarReporteRechazos)
		reportes.POST("/metricas", c.GenerarReporteMetricasRendimiento)
		reportes.POST("/tributario", c.GenerarReporteTributario)
		reportes.GET("/:tipo/:id", c.ObtenerReporte)
		reportes.GET("/:tipo", c.ListarReportes)
	}
}
