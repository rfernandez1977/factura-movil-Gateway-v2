package controllers

import (
	"net/http"

	"github.com/fmgo/models"
	"github.com/fmgo/services"

	"github.com/gin-gonic/gin"
)

// RetryController maneja las peticiones relacionadas con los reintentos
type RetryController struct {
	retryService *services.RetryService
}

// NewRetryController crea una nueva instancia del controlador de reintentos
func NewRetryController(retryService *services.RetryService) *RetryController {
	return &RetryController{
		retryService: retryService,
	}
}

// AgregarReintento agrega un elemento a la cola de reintentos
func (c *RetryController) AgregarReintento(ctx *gin.Context) {
	var reintento models.ColaReintentos
	if err := ctx.ShouldBindJSON(&reintento); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.retryService.AgregarReintento(ctx.Request.Context(), &reintento); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"id": reintento.ID})
}

// ProcesarReintentos procesa los reintentos pendientes
func (c *RetryController) ProcesarReintentos(ctx *gin.Context) {
	if err := c.retryService.ProcesarReintentos(ctx.Request.Context()); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Reintentos procesados exitosamente"})
}

// RegisterRoutes registra las rutas del controlador
func (c *RetryController) RegisterRoutes(router *gin.RouterGroup) {
	retry := router.Group("/retry")
	{
		retry.POST("/reintentos", c.AgregarReintento)
		retry.POST("/procesar", c.ProcesarReintentos)
	}
}
