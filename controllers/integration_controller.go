package controllers

import (
	"net/http"

	"FMgo/models"
	"FMgo/services"

	"github.com/gin-gonic/gin"
)

// IntegrationController maneja las peticiones relacionadas con la integración
type IntegrationController struct {
	integrationService *services.IntegrationService
}

// NewIntegrationController crea una nueva instancia del controlador de integración
func NewIntegrationController(integrationService *services.IntegrationService) *IntegrationController {
	return &IntegrationController{
		integrationService: integrationService,
	}
}

// IniciarSincronizacion inicia un proceso de sincronización
func (c *IntegrationController) IniciarSincronizacion(ctx *gin.Context) {
	var request struct {
		ERPID     string                 `json:"erp_id" binding:"required"`
		Entidad   string                 `json:"entidad" binding:"required"`
		Direccion string                 `json:"direccion" binding:"required"`
		Datos     map[string]interface{} `json:"datos" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	registro, err := c.integrationService.IniciarSincronizacion(ctx.Request.Context(), request.ERPID, request.Entidad, request.Direccion, request.Datos)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, registro)
}

// ProcesarSincronizacion procesa un registro de sincronización
func (c *IntegrationController) ProcesarSincronizacion(ctx *gin.Context) {
	registroID := ctx.Param("id")
	if registroID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de registro no proporcionado"})
		return
	}

	if err := c.integrationService.ProcesarSincronizacion(ctx.Request.Context(), registroID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Sincronización procesada exitosamente"})
}

// RegistrarMetrica registra una métrica de integración
func (c *IntegrationController) RegistrarMetrica(ctx *gin.Context) {
	var metrica models.MetricaIntegracion
	if err := ctx.ShouldBindJSON(&metrica); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.integrationService.RegistrarMetrica(ctx.Request.Context(), &metrica); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"id": metrica.ID})
}

// RegistrarAlerta registra una alerta
func (c *IntegrationController) RegistrarAlerta(ctx *gin.Context) {
	var alerta models.Alerta
	if err := ctx.ShouldBindJSON(&alerta); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.integrationService.RegistrarAlerta(ctx.Request.Context(), &alerta); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"id": alerta.ID})
}

// AgregarReintento agrega un elemento a la cola de reintentos
func (c *IntegrationController) AgregarReintento(ctx *gin.Context) {
	var reintento models.ColaReintentos
	if err := ctx.ShouldBindJSON(&reintento); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.integrationService.AgregarReintento(ctx.Request.Context(), &reintento); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"id": reintento.ID})
}

// RegisterRoutes registra las rutas del controlador
func (c *IntegrationController) RegisterRoutes(router *gin.RouterGroup) {
	integration := router.Group("/integration")
	{
		integration.POST("/sincronizacion", c.IniciarSincronizacion)
		integration.POST("/sincronizacion/:id/procesar", c.ProcesarSincronizacion)
		integration.POST("/metricas", c.RegistrarMetrica)
		integration.POST("/alertas", c.RegistrarAlerta)
		integration.POST("/reintentos", c.AgregarReintento)
	}
}
