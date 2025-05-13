package controllers

import (
	"net/http"

	"github.com/cursor/FMgo/models"
	"github.com/cursor/FMgo/services"

	"github.com/gin-gonic/gin"
)

// OrchestrationController maneja las peticiones relacionadas con la orquestación
type OrchestrationController struct {
	orchestrationService *services.OrchestrationService
}

// NewOrchestrationController crea una nueva instancia del controlador de orquestación
func NewOrchestrationController(orchestrationService *services.OrchestrationService) *OrchestrationController {
	return &OrchestrationController{
		orchestrationService: orchestrationService,
	}
}

// EjecutarFlujo ejecuta un flujo de trabajo
func (c *OrchestrationController) EjecutarFlujo(ctx *gin.Context) {
	var flujo models.FlujoIntegracion
	if err := ctx.ShouldBindJSON(&flujo); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.orchestrationService.EjecutarFlujo(ctx.Request.Context(), &flujo); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"id": flujo.ID})
}

// RegisterRoutes registra las rutas del controlador
func (c *OrchestrationController) RegisterRoutes(router *gin.RouterGroup) {
	orchestration := router.Group("/orchestration")
	{
		orchestration.POST("/flujos", c.EjecutarFlujo)
	}
}
