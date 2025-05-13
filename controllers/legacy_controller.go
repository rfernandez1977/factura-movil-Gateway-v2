package controllers

import (
	"net/http"

	"github.com/cursor/FMgo/models"
	"github.com/cursor/FMgo/services"

	"github.com/gin-gonic/gin"
)

// LegacyController maneja las peticiones relacionadas con la integración legacy
type LegacyController struct {
	legacyService *services.LegacyService
}

// NewLegacyController crea una nueva instancia del controlador legacy
func NewLegacyController(legacyService *services.LegacyService) *LegacyController {
	return &LegacyController{
		legacyService: legacyService,
	}
}

// RegistrarConfiguracionArchivoPlano registra una nueva configuración de archivo plano
func (c *LegacyController) RegistrarConfiguracionArchivoPlano(ctx *gin.Context) {
	var config models.ConfiguracionArchivoPlano
	if err := ctx.ShouldBindJSON(&config); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.legacyService.RegistrarConfiguracionArchivoPlano(ctx.Request.Context(), &config); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"id": config.ID})
}

// RegistrarConfiguracionProtocolo registra una nueva configuración de protocolo
func (c *LegacyController) RegistrarConfiguracionProtocolo(ctx *gin.Context) {
	var config models.ConfiguracionProtocolo
	if err := ctx.ShouldBindJSON(&config); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.legacyService.RegistrarConfiguracionProtocolo(ctx.Request.Context(), &config); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"id": config.ID})
}

// RegistrarTransformacionLegacy registra una nueva transformación legacy
func (c *LegacyController) RegistrarTransformacionLegacy(ctx *gin.Context) {
	var transformacion models.TransformacionLegacy
	if err := ctx.ShouldBindJSON(&transformacion); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.legacyService.RegistrarTransformacionLegacy(ctx.Request.Context(), &transformacion); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"id": transformacion.ID})
}

// ProcesarArchivoPlano procesa un archivo plano
func (c *LegacyController) ProcesarArchivoPlano(ctx *gin.Context) {
	erpID := ctx.Param("erp_id")
	if erpID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de ERP no proporcionado"})
		return
	}

	archivo, err := ctx.FormFile("archivo")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Archivo no proporcionado"})
		return
	}

	// Guardar archivo temporalmente
	rutaTemporal := "/tmp/" + archivo.Filename
	if err := ctx.SaveUploadedFile(archivo, rutaTemporal); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al guardar archivo"})
		return
	}

	// Procesar archivo
	if err := c.legacyService.ProcesarArchivoPlano(ctx.Request.Context(), erpID, rutaTemporal); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Archivo procesado exitosamente"})
}

// TransferirArchivo transfiere un archivo usando el protocolo configurado
func (c *LegacyController) TransferirArchivo(ctx *gin.Context) {
	erpID := ctx.Param("erp_id")
	if erpID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de ERP no proporcionado"})
		return
	}

	archivo, err := ctx.FormFile("archivo")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Archivo no proporcionado"})
		return
	}

	// Guardar archivo temporalmente
	rutaTemporal := "/tmp/" + archivo.Filename
	if err := ctx.SaveUploadedFile(archivo, rutaTemporal); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al guardar archivo"})
		return
	}

	// Transferir archivo
	if err := c.legacyService.TransferirArchivo(ctx.Request.Context(), erpID, rutaTemporal); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Archivo transferido exitosamente"})
}

// RegisterRoutes registra las rutas del controlador
func (c *LegacyController) RegisterRoutes(router *gin.RouterGroup) {
	legacy := router.Group("/legacy")
	{
		legacy.POST("/configuraciones/archivos", c.RegistrarConfiguracionArchivoPlano)
		legacy.POST("/configuraciones/protocolos", c.RegistrarConfiguracionProtocolo)
		legacy.POST("/transformaciones", c.RegistrarTransformacionLegacy)
		legacy.POST("/:erp_id/procesar", c.ProcesarArchivoPlano)
		legacy.POST("/:erp_id/transferir", c.TransferirArchivo)
	}
}
