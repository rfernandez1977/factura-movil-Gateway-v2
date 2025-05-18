package controllers

import (
	"net/http"
	"time"

	"FMgo/models"
	"FMgo/services"
	"FMgo/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ConfiguracionController maneja las peticiones HTTP relacionadas con configuración
type ConfiguracionController struct {
	configuracionService *services.ConfiguracionService
	logger               *zap.Logger
}

// NewConfiguracionController crea una nueva instancia del controlador de configuración
func NewConfiguracionController(configuracionService *services.ConfiguracionService, logger *zap.Logger) *ConfiguracionController {
	return &ConfiguracionController{
		configuracionService: configuracionService,
		logger:               logger,
	}
}

// ObtenerConfiguracion maneja la obtención de la configuración
func (c *ConfiguracionController) ObtenerConfiguracion(ctx *gin.Context) {
	empresaID := ctx.Query("empresa_id")
	if empresaID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de empresa es requerido"})
		return
	}

	configuracion, err := c.configuracionService.ObtenerConfiguracion(empresaID)
	if err != nil {
		c.logger.Error("Error al obtener configuración",
			zap.String("empresaID", empresaID),
			zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, configuracion)
}

// ActualizarConfiguracion maneja la actualización de la configuración
func (c *ConfiguracionController) ActualizarConfiguracion(ctx *gin.Context) {
	start := time.Now()

	var configuracion models.Configuracion
	if err := ctx.ShouldBindJSON(&configuracion); err != nil {
		c.logger.Error("Error al vincular JSON",
			zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.configuracionService.ActualizarConfiguracion(&configuracion); err != nil {
		c.logger.Error("Error al actualizar configuración",
			zap.String("empresaID", configuracion.EmpresaID),
			zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.logger.Info("Configuración actualizada exitosamente",
		zap.String("empresa_id", configuracion.EmpresaID),
	)

	duration := time.Since(start).Seconds()
	utils.RecordHTTPRequest(
		ctx.Request.Method,
		ctx.Request.URL.Path,
		http.StatusOK,
		duration,
		float64(ctx.Request.ContentLength),
		float64(len(configuracion.ID)),
	)

	ctx.JSON(http.StatusOK, configuracion)
}

// ObtenerConfiguracionSII maneja la obtención de la configuración del SII
func (c *ConfiguracionController) ObtenerConfiguracionSII(ctx *gin.Context) {
	empresaID := ctx.Query("empresa_id")
	if empresaID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de empresa es requerido"})
		return
	}

	configuracion, err := c.configuracionService.ObtenerConfiguracionSII(empresaID)
	if err != nil {
		c.logger.Error("Error al obtener configuración SII",
			zap.String("empresaID", empresaID),
			zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, configuracion)
}

// ActualizarConfiguracionSII maneja la actualización de la configuración del SII
func (c *ConfiguracionController) ActualizarConfiguracionSII(ctx *gin.Context) {
	var configuracion models.ConfiguracionSIIEmpresa
	if err := ctx.ShouldBindJSON(&configuracion); err != nil {
		c.logger.Error("Error al vincular JSON",
			zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.configuracionService.ActualizarConfiguracionSII(&configuracion); err != nil {
		c.logger.Error("Error al actualizar configuración SII",
			zap.String("empresaID", configuracion.EmpresaID),
			zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.logger.Info("Configuración SII actualizada exitosamente",
		zap.String("empresa_id", configuracion.EmpresaID),
	)

	ctx.JSON(http.StatusOK, configuracion)
}

// ObtenerConfiguracionEmail maneja la obtención de la configuración de email
func (c *ConfiguracionController) ObtenerConfiguracionEmail(ctx *gin.Context) {
	empresaID := ctx.Query("empresa_id")
	if empresaID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de empresa es requerido"})
		return
	}

	configuracion, err := c.configuracionService.ObtenerConfiguracionEmail(empresaID)
	if err != nil {
		c.logger.Error("Error al obtener configuración de email",
			zap.String("empresaID", empresaID),
			zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, configuracion)
}

// ActualizarConfiguracionEmail maneja la actualización de la configuración de email
func (c *ConfiguracionController) ActualizarConfiguracionEmail(ctx *gin.Context) {
	var configuracion models.ConfiguracionEmail
	if err := ctx.ShouldBindJSON(&configuracion); err != nil {
		c.logger.Error("Error al vincular JSON",
			zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.configuracionService.ActualizarConfiguracionEmail(&configuracion); err != nil {
		c.logger.Error("Error al actualizar configuración de email",
			zap.String("empresaID", configuracion.EmpresaID),
			zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.logger.Info("Configuración de email actualizada exitosamente",
		zap.String("empresa_id", configuracion.EmpresaID),
	)

	ctx.JSON(http.StatusOK, configuracion)
}

// RegisterRoutes registra las rutas del controlador
func (c *ConfiguracionController) RegisterRoutes(router *gin.RouterGroup) {
	configuracion := router.Group("/configuracion")
	{
		configuracion.GET("", c.ObtenerConfiguracion)
		configuracion.PUT("", c.ActualizarConfiguracion)
		configuracion.GET("/sii", c.ObtenerConfiguracionSII)
		configuracion.PUT("/sii", c.ActualizarConfiguracionSII)
		configuracion.GET("/email", c.ObtenerConfiguracionEmail)
		configuracion.PUT("/email", c.ActualizarConfiguracionEmail)
	}
}
