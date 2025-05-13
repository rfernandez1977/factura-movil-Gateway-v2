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

// RolesController maneja las peticiones HTTP relacionadas con roles
type RolesController struct {
	rolesService *services.RolesService
}

// NewRolesController crea una nueva instancia del controlador de roles
func NewRolesController(rolesService *services.RolesService) *RolesController {
	return &RolesController{
		rolesService: rolesService,
	}
}

// CrearRol maneja la creación de un nuevo rol
func (c *RolesController) CrearRol(ctx *gin.Context) {
	start := time.Now()

	var rol models.Rol
	if err := ctx.ShouldBindJSON(&rol); err != nil {
		utils.LogError(err, zap.String("endpoint", "CrearRol"))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.rolesService.CrearRol(&rol); err != nil {
		utils.LogError(err, zap.String("endpoint", "CrearRol"))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.LogInfo("rol creado exitosamente",
		zap.String("nombre", rol.Nombre),
		zap.String("descripcion", rol.Descripcion),
	)

	duration := time.Since(start).Seconds()
	utils.RecordHTTPRequest(
		ctx.Request.Method,
		ctx.Request.URL.Path,
		http.StatusCreated,
		duration,
		float64(ctx.Request.ContentLength),
		float64(len(rol.ID)),
	)

	ctx.JSON(http.StatusCreated, rol)
}

// ObtenerRol maneja la obtención de un rol por ID
func (c *RolesController) ObtenerRol(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de rol es requerido"})
		return
	}

	rol, err := c.rolesService.ObtenerRol(id)
	if err != nil {
		utils.LogError(err, zap.String("endpoint", "ObtenerRol"))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, rol)
}

// ActualizarRol maneja la actualización de un rol
func (c *RolesController) ActualizarRol(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de rol es requerido"})
		return
	}

	var rol models.Rol
	if err := ctx.ShouldBindJSON(&rol); err != nil {
		utils.LogError(err, zap.String("endpoint", "ActualizarRol"))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rol.ID = id
	if err := c.rolesService.ActualizarRol(&rol); err != nil {
		utils.LogError(err, zap.String("endpoint", "ActualizarRol"))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.LogInfo("rol actualizado exitosamente",
		zap.String("id", id),
		zap.String("nombre", rol.Nombre),
	)

	ctx.JSON(http.StatusOK, rol)
}

// EliminarRol maneja la eliminación de un rol
func (c *RolesController) EliminarRol(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de rol es requerido"})
		return
	}

	if err := c.rolesService.EliminarRol(id); err != nil {
		utils.LogError(err, zap.String("endpoint", "EliminarRol"))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.LogInfo("rol eliminado exitosamente", zap.String("id", id))
	ctx.JSON(http.StatusOK, gin.H{"message": "Rol eliminado exitosamente"})
}

// ListarRoles maneja la obtención de una lista de roles
func (c *RolesController) ListarRoles(ctx *gin.Context) {
	empresaID := ctx.Query("empresa_id")
	if empresaID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de empresa es requerido"})
		return
	}

	roles, err := c.rolesService.ListarRoles(empresaID)
	if err != nil {
		utils.LogError(err, zap.String("endpoint", "ListarRoles"))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, roles)
}

// AsignarPermisos maneja la asignación de permisos a un rol
func (c *RolesController) AsignarPermisos(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de rol es requerido"})
		return
	}

	var request struct {
		Permisos []string `json:"permisos" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.rolesService.AsignarPermisos(id, request.Permisos); err != nil {
		utils.LogError(err, zap.String("endpoint", "AsignarPermisos"))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.LogInfo("permisos asignados exitosamente",
		zap.String("rol_id", id),
		zap.Strings("permisos", request.Permisos),
	)

	ctx.JSON(http.StatusOK, gin.H{"message": "Permisos asignados exitosamente"})
}

// RegisterRoutes registra las rutas del controlador
func (c *RolesController) RegisterRoutes(router *gin.RouterGroup) {
	roles := router.Group("/roles")
	{
		roles.POST("", c.CrearRol)
		roles.GET("/:id", c.ObtenerRol)
		roles.PUT("/:id", c.ActualizarRol)
		roles.DELETE("/:id", c.EliminarRol)
		roles.GET("", c.ListarRoles)
		roles.POST("/:id/permisos", c.AsignarPermisos)
	}
}
