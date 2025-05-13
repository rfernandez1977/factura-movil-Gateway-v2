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

// SucursalesController maneja las peticiones HTTP relacionadas con sucursales
type SucursalesController struct {
	sucursalesService *services.SucursalesService
}

// NewSucursalesController crea una nueva instancia del controlador de sucursales
func NewSucursalesController(sucursalesService *services.SucursalesService) *SucursalesController {
	return &SucursalesController{
		sucursalesService: sucursalesService,
	}
}

// CrearSucursal maneja la creación de una nueva sucursal
func (c *SucursalesController) CrearSucursal(ctx *gin.Context) {
	start := time.Now()

	var sucursal models.Sucursal
	if err := ctx.ShouldBindJSON(&sucursal); err != nil {
		utils.LogError(err, zap.String("endpoint", "CrearSucursal"))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.sucursalesService.CrearSucursal(&sucursal); err != nil {
		utils.LogError(err, zap.String("endpoint", "CrearSucursal"))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.LogInfo("sucursal creada exitosamente",
		zap.String("codigo", sucursal.Codigo),
		zap.String("nombre", sucursal.Nombre),
	)

	duration := time.Since(start).Seconds()
	utils.RecordHTTPRequest(
		ctx.Request.Method,
		ctx.Request.URL.Path,
		http.StatusCreated,
		duration,
		float64(ctx.Request.ContentLength),
		float64(len(sucursal.ID)),
	)

	ctx.JSON(http.StatusCreated, sucursal)
}

// ObtenerSucursal maneja la obtención de una sucursal por ID
func (c *SucursalesController) ObtenerSucursal(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de sucursal es requerido"})
		return
	}

	sucursal, err := c.sucursalesService.ObtenerSucursal(id)
	if err != nil {
		utils.LogError(err, zap.String("endpoint", "ObtenerSucursal"))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, sucursal)
}

// ActualizarSucursal maneja la actualización de una sucursal
func (c *SucursalesController) ActualizarSucursal(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de sucursal es requerido"})
		return
	}

	var sucursal models.Sucursal
	if err := ctx.ShouldBindJSON(&sucursal); err != nil {
		utils.LogError(err, zap.String("endpoint", "ActualizarSucursal"))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sucursal.ID = id
	if err := c.sucursalesService.ActualizarSucursal(&sucursal); err != nil {
		utils.LogError(err, zap.String("endpoint", "ActualizarSucursal"))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.LogInfo("sucursal actualizada exitosamente",
		zap.String("id", id),
		zap.String("codigo", sucursal.Codigo),
	)

	ctx.JSON(http.StatusOK, sucursal)
}

// EliminarSucursal maneja la eliminación de una sucursal
func (c *SucursalesController) EliminarSucursal(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de sucursal es requerido"})
		return
	}

	if err := c.sucursalesService.EliminarSucursal(id); err != nil {
		utils.LogError(err, zap.String("endpoint", "EliminarSucursal"))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.LogInfo("sucursal eliminada exitosamente", zap.String("id", id))
	ctx.JSON(http.StatusOK, gin.H{"message": "Sucursal eliminada exitosamente"})
}

// ListarSucursales maneja la obtención de una lista de sucursales
func (c *SucursalesController) ListarSucursales(ctx *gin.Context) {
	empresaID := ctx.Query("empresa_id")
	if empresaID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de empresa es requerido"})
		return
	}

	sucursales, err := c.sucursalesService.ListarSucursales(empresaID)
	if err != nil {
		utils.LogError(err, zap.String("endpoint", "ListarSucursales"))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, sucursales)
}

// RegisterRoutes registra las rutas del controlador
func (c *SucursalesController) RegisterRoutes(router *gin.RouterGroup) {
	sucursales := router.Group("/sucursales")
	{
		sucursales.POST("", c.CrearSucursal)
		sucursales.GET("/:id", c.ObtenerSucursal)
		sucursales.PUT("/:id", c.ActualizarSucursal)
		sucursales.DELETE("/:id", c.EliminarSucursal)
		sucursales.GET("", c.ListarSucursales)
	}
}
