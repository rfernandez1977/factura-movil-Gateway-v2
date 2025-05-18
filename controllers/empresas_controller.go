package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"FMgo/models"
	"FMgo/services"
	"FMgo/utils"
	"go.uber.org/zap"
)

// EmpresasController maneja las peticiones HTTP relacionadas con empresas
type EmpresasController struct {
	empresasService *services.EmpresasService
}

// NewEmpresasController crea una nueva instancia del controlador de empresas
func NewEmpresasController(empresasService *services.EmpresasService) *EmpresasController {
	return &EmpresasController{
		empresasService: empresasService,
	}
}

// CrearEmpresa maneja la creación de una nueva empresa
func (c *EmpresasController) CrearEmpresa(ctx *gin.Context) {
	start := time.Now()

	var empresa models.Empresa
	if err := ctx.ShouldBindJSON(&empresa); err != nil {
		utils.LogError(err, zap.String("endpoint", "CrearEmpresa"))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.empresasService.CrearEmpresa(&empresa); err != nil {
		utils.LogError(err, zap.String("endpoint", "CrearEmpresa"))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.LogInfo("empresa creada exitosamente",
		zap.String("rut", empresa.Rut),
		zap.String("razon_social", empresa.RazonSocial),
	)

	duration := time.Since(start).Seconds()
	utils.RecordHTTPRequest(
		ctx.Request.Method,
		ctx.Request.URL.Path,
		http.StatusCreated,
		duration,
		float64(ctx.Request.ContentLength),
		float64(len(empresa.ID)),
	)

	ctx.JSON(http.StatusCreated, empresa)
}

// ObtenerEmpresa maneja la obtención de una empresa por ID
func (c *EmpresasController) ObtenerEmpresa(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de empresa es requerido"})
		return
	}

	empresa, err := c.empresasService.ObtenerEmpresa(id)
	if err != nil {
		utils.LogError(err, zap.String("endpoint", "ObtenerEmpresa"))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, empresa)
}

// ActualizarEmpresa maneja la actualización de una empresa
func (c *EmpresasController) ActualizarEmpresa(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de empresa es requerido"})
		return
	}

	var empresa models.Empresa
	if err := ctx.ShouldBindJSON(&empresa); err != nil {
		utils.LogError(err, zap.String("endpoint", "ActualizarEmpresa"))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	empresa.ID = id
	if err := c.empresasService.ActualizarEmpresa(&empresa); err != nil {
		utils.LogError(err, zap.String("endpoint", "ActualizarEmpresa"))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.LogInfo("empresa actualizada exitosamente",
		zap.String("id", id),
		zap.String("rut", empresa.Rut),
	)

	ctx.JSON(http.StatusOK, empresa)
}

// EliminarEmpresa maneja la eliminación de una empresa
func (c *EmpresasController) EliminarEmpresa(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de empresa es requerido"})
		return
	}

	if err := c.empresasService.EliminarEmpresa(id); err != nil {
		utils.LogError(err, zap.String("endpoint", "EliminarEmpresa"))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.LogInfo("empresa eliminada exitosamente", zap.String("id", id))
	ctx.JSON(http.StatusOK, gin.H{"message": "Empresa eliminada exitosamente"})
}

// ListarEmpresas maneja la obtención de una lista de empresas
func (c *EmpresasController) ListarEmpresas(ctx *gin.Context) {
	empresas, err := c.empresasService.ListarEmpresas()
	if err != nil {
		utils.LogError(err, zap.String("endpoint", "ListarEmpresas"))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, empresas)
}

// RegisterRoutes registra las rutas del controlador
func (c *EmpresasController) RegisterRoutes(router *gin.RouterGroup) {
	empresas := router.Group("/empresas")
	{
		empresas.POST("", c.CrearEmpresa)
		empresas.GET("/:id", c.ObtenerEmpresa)
		empresas.PUT("/:id", c.ActualizarEmpresa)
		empresas.DELETE("/:id", c.EliminarEmpresa)
		empresas.GET("", c.ListarEmpresas)
	}
}
