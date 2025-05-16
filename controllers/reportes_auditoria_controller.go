package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/fmgo/services"
)

// ReportesAuditoriaController maneja las peticiones relacionadas con reportes de auditoría
type ReportesAuditoriaController struct {
	reportesAuditoriaService *services.ReportesAuditoriaService
}

// NewReportesAuditoriaController crea una nueva instancia del controlador de reportes de auditoría
func NewReportesAuditoriaController(reportesAuditoriaService *services.ReportesAuditoriaService) *ReportesAuditoriaController {
	return &ReportesAuditoriaController{
		reportesAuditoriaService: reportesAuditoriaService,
	}
}

// GenerarReporteAuditoria genera un reporte de auditoría
func (c *ReportesAuditoriaController) GenerarReporteAuditoria(ctx *gin.Context) {
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

	reporte, err := c.reportesAuditoriaService.GenerarReporteAuditoria(
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

// GenerarReporteCumplimiento genera un reporte de cumplimiento
func (c *ReportesAuditoriaController) GenerarReporteCumplimiento(ctx *gin.Context) {
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

	reporte, err := c.reportesAuditoriaService.GenerarReporteCumplimiento(
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

// ObtenerReporteAuditoria obtiene un reporte de auditoría por su ID
func (c *ReportesAuditoriaController) ObtenerReporteAuditoria(ctx *gin.Context) {
	id := ctx.Param("id")

	reporte, err := c.reportesAuditoriaService.ObtenerReporteAuditoria(context.Background(), id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, reporte)
}

// ObtenerReporteCumplimiento obtiene un reporte de cumplimiento por su ID
func (c *ReportesAuditoriaController) ObtenerReporteCumplimiento(ctx *gin.Context) {
	id := ctx.Param("id")

	reporte, err := c.reportesAuditoriaService.ObtenerReporteCumplimiento(context.Background(), id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, reporte)
}

// ListarReportesAuditoria obtiene una lista de reportes de auditoría
func (c *ReportesAuditoriaController) ListarReportesAuditoria(ctx *gin.Context) {
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

	reportes, err := c.reportesAuditoriaService.ListarReportesAuditoria(
		context.Background(),
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

// ListarReportesCumplimiento obtiene una lista de reportes de cumplimiento
func (c *ReportesAuditoriaController) ListarReportesCumplimiento(ctx *gin.Context) {
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

	reportes, err := c.reportesAuditoriaService.ListarReportesCumplimiento(
		context.Background(),
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
func (c *ReportesAuditoriaController) RegisterRoutes(router *gin.RouterGroup) {
	auditoria := router.Group("/auditoria")
	{
		auditoria.POST("/reporte", c.GenerarReporteAuditoria)
		auditoria.POST("/cumplimiento", c.GenerarReporteCumplimiento)
		auditoria.GET("/reporte/:id", c.ObtenerReporteAuditoria)
		auditoria.GET("/cumplimiento/:id", c.ObtenerReporteCumplimiento)
		auditoria.GET("/reportes", c.ListarReportesAuditoria)
		auditoria.GET("/cumplimiento", c.ListarReportesCumplimiento)
	}
}
