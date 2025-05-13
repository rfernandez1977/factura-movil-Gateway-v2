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

// ConfiguracionController maneja las peticiones HTTP relacionadas con configuración
type ConfiguracionController struct {
	configuracionService *services.ConfiguracionService
}

// NewConfiguracionController crea una nueva instancia del controlador de configuración
func NewConfiguracionController(configuracionService *services.ConfiguracionService) *ConfiguracionController {
	return &ConfiguracionController{
		configuracionService: configuracionService,
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
		utils.LogError(err, zap.String("endpoint", "ObtenerConfiguracion"))
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
		utils.LogError(err, zap.String("endpoint", "ActualizarConfiguracion"))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.configuracionService.ActualizarConfiguracion(&configuracion); err != nil {
		utils.LogError(err, zap.String("endpoint", "ActualizarConfiguracion"))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.LogInfo("configuración actualizada exitosamente",
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
		utils.LogError(err, zap.String("endpoint", "ObtenerConfiguracionSII"))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, configuracion)
}

// ActualizarConfiguracionSII maneja la actualización de la configuración del SII
func (c *ConfiguracionController) ActualizarConfiguracionSII(ctx *gin.Context) {
	var configuracion models.ConfiguracionSII
	if err := ctx.ShouldBindJSON(&configuracion); err != nil {
		utils.LogError(err, zap.String("endpoint", "ActualizarConfiguracionSII"))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.configuracionService.ActualizarConfiguracionSII(&configuracion); err != nil {
		utils.LogError(err, zap.String("endpoint", "ActualizarConfiguracionSII"))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.LogInfo("configuración SII actualizada exitosamente",
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
		utils.LogError(err, zap.String("endpoint", "ObtenerConfiguracionEmail"))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, configuracion)
}

// ActualizarConfiguracionEmail maneja la actualización de la configuración de email
func (c *ConfiguracionController) ActualizarConfiguracionEmail(ctx *gin.Context) {
	var configuracion models.ConfiguracionEmail
	if err := ctx.ShouldBindJSON(&configuracion); err != nil {
		utils.LogError(err, zap.String("endpoint", "ActualizarConfiguracionEmail"))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.configuracionService.ActualizarConfiguracionEmail(&configuracion); err != nil {
		utils.LogError(err, zap.String("endpoint", "ActualizarConfiguracionEmail"))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.LogInfo("configuración email actualizada exitosamente",
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
