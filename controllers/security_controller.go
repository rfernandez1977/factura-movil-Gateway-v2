package controllers

import (
	"net/http"
	"time"

	"github.com/cursor/FMgo/models"
	"github.com/cursor/FMgo/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SecurityController maneja las peticiones relacionadas con la seguridad
type SecurityController struct {
	securityService *services.SecurityService
}

// NewSecurityController crea una nueva instancia del controlador de seguridad
func NewSecurityController(securityService *services.SecurityService) *SecurityController {
	return &SecurityController{
		securityService: securityService,
	}
}

// GenerarCertificado genera un nuevo certificado digital
func (c *SecurityController) GenerarCertificado(ctx *gin.Context) {
	var config models.ConfiguracionCertificado
	if err := ctx.ShouldBindJSON(&config); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	certificado, err := c.securityService.GenerarCertificado(ctx.Request.Context(), &config)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, certificado)
}

// ValidarCertificado valida un certificado digital
func (c *SecurityController) ValidarCertificado(ctx *gin.Context) {
	certificadoID, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de certificado inválido"})
		return
	}

	if err := c.securityService.ValidarCertificado(ctx.Request.Context(), certificadoID); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Certificado válido"})
}

// RegistrarAcceso registra un acceso al sistema
func (c *SecurityController) RegistrarAcceso(ctx *gin.Context) {
	var acceso models.RegistroAcceso
	if err := ctx.ShouldBindJSON(&acceso); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.securityService.RegistrarAcceso(ctx.Request.Context(), &acceso); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"id": acceso.ID})
}

// ObtenerRegistrosAcceso obtiene los registros de acceso
func (c *SecurityController) ObtenerRegistrosAcceso(ctx *gin.Context) {
	var filtro struct {
		Inicio time.Time `form:"inicio" binding:"required"`
		Fin    time.Time `form:"fin" binding:"required"`
	}

	if err := ctx.ShouldBindQuery(&filtro); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	registros, err := c.securityService.ObtenerRegistrosAcceso(ctx.Request.Context(), map[string]interface{}{
		"fecha": map[string]interface{}{
			"$gte": filtro.Inicio,
			"$lte": filtro.Fin,
		},
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, registros)
}

// GenerarReporteSeguridad genera un reporte de seguridad
func (c *SecurityController) GenerarReporteSeguridad(ctx *gin.Context) {
	var request struct {
		Inicio time.Time `json:"inicio" binding:"required"`
		Fin    time.Time `json:"fin" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reporte, err := c.securityService.GenerarReporteSeguridad(ctx.Request.Context(), request.Inicio, request.Fin)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, reporte)
}

// RegisterRoutes registra las rutas del controlador
func (c *SecurityController) RegisterRoutes(router *gin.RouterGroup) {
	security := router.Group("/security")
	{
		security.POST("/certificados", c.GenerarCertificado)
		security.GET("/certificados/:id/validar", c.ValidarCertificado)
		security.POST("/accesos", c.RegistrarAcceso)
		security.GET("/accesos", c.ObtenerRegistrosAcceso)
		security.POST("/reportes", c.GenerarReporteSeguridad)
	}
}
