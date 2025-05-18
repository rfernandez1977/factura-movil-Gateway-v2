package controllers

import (
	"net/http"
	"time"

	"FMgo/models"
	"FMgo/services"

	"github.com/gin-gonic/gin"
)

// ERPController maneja las peticiones relacionadas con la integración ERP
type ERPController struct {
	erpService *services.ERPService
}

// NewERPController crea una nueva instancia del controlador ERP
func NewERPController(erpService *services.ERPService) *ERPController {
	return &ERPController{
		erpService: erpService,
	}
}

// RegistrarConfiguracionERP registra una nueva configuración de ERP
func (c *ERPController) RegistrarConfiguracionERP(ctx *gin.Context) {
	var config models.ConfiguracionERP
	if err := ctx.ShouldBindJSON(&config); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.erpService.RegistrarConfiguracionERP(ctx.Request.Context(), &config); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"id": config.ID})
}

// ObtenerConfiguracionERP obtiene la configuración de un ERP
func (c *ERPController) ObtenerConfiguracionERP(ctx *gin.Context) {
	erpID := ctx.Param("id")
	if erpID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de ERP no proporcionado"})
		return
	}

	config, err := c.erpService.ObtenerConfiguracionERP(ctx.Request.Context(), erpID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, config)
}

// RegistrarMapeoCampos registra un mapeo de campos para una entidad
func (c *ERPController) RegistrarMapeoCampos(ctx *gin.Context) {
	var mapeo models.MapeoCamposERP
	if err := ctx.ShouldBindJSON(&mapeo); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.erpService.RegistrarMapeoCampos(ctx.Request.Context(), &mapeo); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"id": mapeo.ID})
}

// ObtenerMapeosCampos obtiene los mapeos de campos para una entidad
func (c *ERPController) ObtenerMapeosCampos(ctx *gin.Context) {
	erpID := ctx.Param("erp_id")
	if erpID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de ERP no proporcionado"})
		return
	}

	entidad := ctx.Query("entidad")

	mapeos, err := c.erpService.ObtenerMapeosCampos(ctx.Request.Context(), erpID, entidad)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, mapeos)
}

// RegistrarEventoERP registra un nuevo evento de sincronización
func (c *ERPController) RegistrarEventoERP(ctx *gin.Context) {
	var evento models.EventoERP
	if err := ctx.ShouldBindJSON(&evento); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.erpService.RegistrarEventoERP(ctx.Request.Context(), &evento); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"id": evento.ID})
}

// ProcesarEventoERP procesa un evento de sincronización
func (c *ERPController) ProcesarEventoERP(ctx *gin.Context) {
	eventoID := ctx.Param("id")
	if eventoID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de evento no proporcionado"})
		return
	}

	if err := c.erpService.ProcesarEventoERP(ctx.Request.Context(), eventoID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Evento procesado exitosamente"})
}

// GenerarReporteIntegracion genera un reporte de integración
func (c *ERPController) GenerarReporteIntegracion(ctx *gin.Context) {
	erpID := ctx.Param("erp_id")
	if erpID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de ERP no proporcionado"})
		return
	}

	fechaInicioStr := ctx.Query("fecha_inicio")
	fechaFinStr := ctx.Query("fecha_fin")

	fechaInicio, err := time.Parse(time.RFC3339, fechaInicioStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Formato de fecha inicio inválido"})
		return
	}

	fechaFin, err := time.Parse(time.RFC3339, fechaFinStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Formato de fecha fin inválido"})
		return
	}

	reporte, err := c.erpService.GenerarReporteIntegracion(ctx.Request.Context(), erpID, fechaInicio, fechaFin)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, reporte)
}

// RegisterRoutes registra las rutas del controlador
func (c *ERPController) RegisterRoutes(router *gin.RouterGroup) {
	erp := router.Group("/erp")
	{
		erp.POST("/configuraciones", c.RegistrarConfiguracionERP)
		erp.GET("/configuraciones/:id", c.ObtenerConfiguracionERP)
		erp.POST("/mapeos", c.RegistrarMapeoCampos)
		erp.GET("/mapeos/:erp_id", c.ObtenerMapeosCampos)
		erp.POST("/eventos", c.RegistrarEventoERP)
		erp.POST("/eventos/:id/procesar", c.ProcesarEventoERP)
		erp.GET("/reportes/:erp_id", c.GenerarReporteIntegracion)
	}
}
