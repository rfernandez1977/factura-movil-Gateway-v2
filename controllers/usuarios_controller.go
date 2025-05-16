package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/fmgo/models"
	"github.com/fmgo/services"
	"github.com/fmgo/utils"
	"go.uber.org/zap"
)

// UsuariosController maneja las peticiones HTTP relacionadas con usuarios
type UsuariosController struct {
	usuariosService *services.UsuariosService
}

// NewUsuariosController crea una nueva instancia del controlador de usuarios
func NewUsuariosController(usuariosService *services.UsuariosService) *UsuariosController {
	return &UsuariosController{
		usuariosService: usuariosService,
	}
}

// CrearUsuario maneja la creación de un nuevo usuario
func (c *UsuariosController) CrearUsuario(ctx *gin.Context) {
	start := time.Now()

	var usuario models.Usuario
	if err := ctx.ShouldBindJSON(&usuario); err != nil {
		utils.LogError(err, zap.String("endpoint", "CrearUsuario"))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.usuariosService.CrearUsuario(&usuario); err != nil {
		utils.LogError(err, zap.String("endpoint", "CrearUsuario"))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.LogInfo("usuario creado exitosamente",
		zap.String("email", usuario.Email),
		zap.String("nombre", usuario.Nombre),
	)

	duration := time.Since(start).Seconds()
	utils.RecordHTTPRequest(
		ctx.Request.Method,
		ctx.Request.URL.Path,
		http.StatusCreated,
		duration,
		float64(ctx.Request.ContentLength),
		float64(len(usuario.ID)),
	)

	ctx.JSON(http.StatusCreated, usuario)
}

// ObtenerUsuario maneja la obtención de un usuario por ID
func (c *UsuariosController) ObtenerUsuario(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de usuario es requerido"})
		return
	}

	usuario, err := c.usuariosService.ObtenerUsuario(id)
	if err != nil {
		utils.LogError(err, zap.String("endpoint", "ObtenerUsuario"))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, usuario)
}

// ActualizarUsuario maneja la actualización de un usuario
func (c *UsuariosController) ActualizarUsuario(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de usuario es requerido"})
		return
	}

	var usuario models.Usuario
	if err := ctx.ShouldBindJSON(&usuario); err != nil {
		utils.LogError(err, zap.String("endpoint", "ActualizarUsuario"))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	usuario.ID = id
	if err := c.usuariosService.ActualizarUsuario(&usuario); err != nil {
		utils.LogError(err, zap.String("endpoint", "ActualizarUsuario"))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.LogInfo("usuario actualizado exitosamente",
		zap.String("id", id),
		zap.String("email", usuario.Email),
	)

	ctx.JSON(http.StatusOK, usuario)
}

// EliminarUsuario maneja la eliminación de un usuario
func (c *UsuariosController) EliminarUsuario(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de usuario es requerido"})
		return
	}

	if err := c.usuariosService.EliminarUsuario(id); err != nil {
		utils.LogError(err, zap.String("endpoint", "EliminarUsuario"))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.LogInfo("usuario eliminado exitosamente", zap.String("id", id))
	ctx.JSON(http.StatusOK, gin.H{"message": "Usuario eliminado exitosamente"})
}

// ListarUsuarios maneja la obtención de una lista de usuarios
func (c *UsuariosController) ListarUsuarios(ctx *gin.Context) {
	empresaID := ctx.Query("empresa_id")
	if empresaID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de empresa es requerido"})
		return
	}

	usuarios, err := c.usuariosService.ListarUsuarios(empresaID)
	if err != nil {
		utils.LogError(err, zap.String("endpoint", "ListarUsuarios"))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, usuarios)
}

// CambiarContrasena maneja el cambio de contraseña de un usuario
func (c *UsuariosController) CambiarContrasena(ctx *gin.Context) {
	var request struct {
		ContrasenaActual string `json:"contrasena_actual" binding:"required"`
		ContrasenaNueva  string `json:"contrasena_nueva" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jwtUtils := utils.NewJWTUtils()
	userID, _ := utils.GetUserID(ctx.GetHeader("Authorization"), jwtUtils)
	if err := c.usuariosService.CambiarContrasena(userID, request.ContrasenaActual, request.ContrasenaNueva); err != nil {
		utils.LogError(err, zap.String("endpoint", "CambiarContrasena"))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.LogInfo("contraseña cambiada exitosamente", zap.String("user_id", userID))
	ctx.JSON(http.StatusOK, gin.H{"message": "Contraseña cambiada exitosamente"})
}

// RegisterRoutes registra las rutas del controlador
func (c *UsuariosController) RegisterRoutes(router *gin.RouterGroup) {
	usuarios := router.Group("/usuarios")
	{
		usuarios.POST("", c.CrearUsuario)
		usuarios.GET("/:id", c.ObtenerUsuario)
		usuarios.PUT("/:id", c.ActualizarUsuario)
		usuarios.DELETE("/:id", c.EliminarUsuario)
		usuarios.GET("", c.ListarUsuarios)
		usuarios.POST("/cambiar-contrasena", c.CambiarContrasena)
	}
}
