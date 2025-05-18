package controllers

import (
	"net/http"
	"time"

	"FMgo/models"
	"FMgo/services"

	"github.com/gin-gonic/gin"
)

// MonitoringController maneja las peticiones relacionadas con el monitoreo
type MonitoringController struct {
	monitoringService *services.MonitoringService
}

// NewMonitoringController crea una nueva instancia del controlador de monitoreo
func NewMonitoringController(monitoringService *services.MonitoringService) *MonitoringController {
	return &MonitoringController{
		monitoringService: monitoringService,
	}
}

// RegistrarMetrica registra una métrica de integración
func (c *MonitoringController) RegistrarMetrica(ctx *gin.Context) {
	var metrica models.MetricaIntegracion
	if err := ctx.ShouldBindJSON(&metrica); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.monitoringService.RegistrarMetrica(ctx.Request.Context(), &metrica); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"id": metrica.ID})
}

// RegistrarAlerta registra una alerta
func (c *MonitoringController) RegistrarAlerta(ctx *gin.Context) {
	var alerta models.Alerta
	if err := ctx.ShouldBindJSON(&alerta); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.monitoringService.RegistrarAlerta(ctx.Request.Context(), &alerta); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"id": alerta.ID})
}

// ObtenerMetricas obtiene las métricas de integración
func (c *MonitoringController) ObtenerMetricas(ctx *gin.Context) {
	var filtro struct {
		Inicio time.Time `form:"inicio" binding:"required"`
		Fin    time.Time `form:"fin" binding:"required"`
	}

	if err := ctx.ShouldBindQuery(&filtro); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	metricas, err := c.monitoringService.ObtenerMetricas(ctx.Request.Context(), map[string]interface{}{
		"fecha": map[string]interface{}{
			"$gte": filtro.Inicio,
			"$lte": filtro.Fin,
		},
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, metricas)
}

// ObtenerAlertas obtiene las alertas
func (c *MonitoringController) ObtenerAlertas(ctx *gin.Context) {
	var filtro struct {
		Inicio time.Time `form:"inicio" binding:"required"`
		Fin    time.Time `form:"fin" binding:"required"`
	}

	if err := ctx.ShouldBindQuery(&filtro); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	alertas, err := c.monitoringService.ObtenerAlertas(ctx.Request.Context(), map[string]interface{}{
		"fecha_creacion": map[string]interface{}{
			"$gte": filtro.Inicio,
			"$lte": filtro.Fin,
		},
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, alertas)
}

// GenerarReporte genera un reporte de integración
func (c *MonitoringController) GenerarReporte(ctx *gin.Context) {
	var request struct {
		Inicio time.Time `json:"inicio" binding:"required"`
		Fin    time.Time `json:"fin" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reporte, err := c.monitoringService.GenerarReporte(ctx.Request.Context(), request.Inicio, request.Fin)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, reporte)
}

// RegisterRoutes registra las rutas del controlador
func (c *MonitoringController) RegisterRoutes(router *gin.RouterGroup) {
	monitoring := router.Group("/monitoring")
	{
		monitoring.POST("/metricas", c.RegistrarMetrica)
		monitoring.POST("/alertas", c.RegistrarAlerta)
		monitoring.GET("/metricas", c.ObtenerMetricas)
		monitoring.GET("/alertas", c.ObtenerAlertas)
		monitoring.POST("/reportes", c.GenerarReporte)
	}
}
