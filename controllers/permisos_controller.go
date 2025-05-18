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

// PermisosController maneja las peticiones HTTP relacionadas con permisos
type PermisosController struct {
	permisosService *services.PermisosService
}

// NewPermisosController crea una nueva instancia del controlador de permisos
func NewPermisosController(permisosService *services.PermisosService) *PermisosController {
	return &PermisosController{
		permisosService: permisosService,
	}
}

// CrearPermiso maneja la creación de un nuevo permiso
func (c *PermisosController) CrearPermiso(ctx *gin.Context) {
	start := time.Now()

	var permiso models.Permiso
	if err := ctx.ShouldBindJSON(&permiso); err != nil {
		utils.LogError(err, zap.String("endpoint", "CrearPermiso"))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.permisosService.CrearPermiso(&permiso); err != nil {
		utils.LogError(err, zap.String("endpoint", "CrearPermiso"))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.LogInfo("permiso creado exitosamente",
		zap.String("nombre", permiso.Nombre),
		zap.String("descripcion", permiso.Descripcion),
	)

	duration := time.Since(start).Seconds()
	utils.RecordHTTPRequest(
		ctx.Request.Method,
		ctx.Request.URL.Path,
		http.StatusCreated,
		duration,
		float64(ctx.Request.ContentLength),
		float64(len(permiso.ID)),
	)

	ctx.JSON(http.StatusCreated, permiso)
}

// ObtenerPermiso maneja la obtención de un permiso por ID
func (c *PermisosController) ObtenerPermiso(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de permiso es requerido"})
		return
	}

	permiso, err := c.permisosService.ObtenerPermiso(id)
	if err != nil {
		utils.LogError(err, zap.String("endpoint", "ObtenerPermiso"))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, permiso)
}

// ActualizarPermiso maneja la actualización de un permiso
func (c *PermisosController) ActualizarPermiso(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de permiso es requerido"})
		return
	}

	var permiso models.Permiso
	if err := ctx.ShouldBindJSON(&permiso); err != nil {
		utils.LogError(err, zap.String("endpoint", "ActualizarPermiso"))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	permiso.ID = id
	if err := c.permisosService.ActualizarPermiso(&permiso); err != nil {
		utils.LogError(err, zap.String("endpoint", "ActualizarPermiso"))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.LogInfo("permiso actualizado exitosamente",
		zap.String("id", id),
		zap.String("nombre", permiso.Nombre),
	)

	ctx.JSON(http.StatusOK, permiso)
}

// EliminarPermiso maneja la eliminación de un permiso
func (c *PermisosController) EliminarPermiso(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de permiso es requerido"})
		return
	}

	if err := c.permisosService.EliminarPermiso(id); err != nil {
		utils.LogError(err, zap.String("endpoint", "EliminarPermiso"))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.LogInfo("permiso eliminado exitosamente", zap.String("id", id))
	ctx.JSON(http.StatusOK, gin.H{"message": "Permiso eliminado exitosamente"})
}

// ListarPermisos maneja la obtención de una lista de permisos
func (c *PermisosController) ListarPermisos(ctx *gin.Context) {
	empresaID := ctx.Query("empresa_id")
	if empresaID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de empresa es requerido"})
		return
	}

	permisos, err := c.permisosService.ListarPermisos(empresaID)
	if err != nil {
		utils.LogError(err, zap.String("endpoint", "ListarPermisos"))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, permisos)
}

// VerificarPermiso maneja la verificación de un permiso para un usuario
func (c *PermisosController) VerificarPermiso(ctx *gin.Context) {
	var request struct {
		UsuarioID string `json:"usuario_id" binding:"required"`
		Permiso   string `json:"permiso" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tienePermiso, err := c.permisosService.VerificarPermiso(request.UsuarioID, request.Permiso)
	if err != nil {
		utils.LogError(err, zap.String("endpoint", "VerificarPermiso"))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"tiene_permiso": tienePermiso})
}

// RegisterRoutes registra las rutas del controlador
func (c *PermisosController) RegisterRoutes(router *gin.RouterGroup) {
	permisos := router.Group("/permisos")
	{
		permisos.POST("", c.CrearPermiso)
		permisos.GET("/:id", c.ObtenerPermiso)
		permisos.PUT("/:id", c.ActualizarPermiso)
		permisos.DELETE("/:id", c.EliminarPermiso)
		permisos.GET("", c.ListarPermisos)
		permisos.POST("/verificar", c.VerificarPermiso)
	}
}
