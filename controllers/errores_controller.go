package controllers

import (
	"net/http"
	"runtime"
	"time"

	"FMgo/models"
	"FMgo/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// ErroresController maneja las peticiones relacionadas con errores
type ErroresController struct {
	erroresService *services.ErroresService
}

// NewErroresController crea una nueva instancia del controlador de errores
func NewErroresController(erroresService *services.ErroresService) *ErroresController {
	return &ErroresController{
		erroresService: erroresService,
	}
}

// RegistrarError registra un nuevo error en el sistema
func (c *ErroresController) RegistrarError(ctx *gin.Context) {
	var request struct {
		Tipo        models.TipoError       `json:"tipo" binding:"required"`
		Severidad   models.SeveridadError  `json:"severidad" binding:"required"`
		Codigo      string                 `json:"codigo" binding:"required"`
		Mensaje     string                 `json:"mensaje" binding:"required"`
		Descripcion string                 `json:"descripcion"`
		Contexto    map[string]interface{} `json:"contexto"`
		Entidad     string                 `json:"entidad"`
		EntidadID   string                 `json:"entidad_id"`
		UsuarioID   string                 `json:"usuario_id"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Obtener stacktrace
	buf := make([]byte, 1<<16)
	runtime.Stack(buf, false)
	stacktrace := string(buf)

	errorDetalle, err := c.erroresService.RegistrarError(
		ctx.Request.Context(),
		request.Tipo,
		request.Severidad,
		request.Codigo,
		request.Mensaje,
		request.Descripcion,
		stacktrace,
		request.Contexto,
		request.Entidad,
		request.EntidadID,
		request.UsuarioID,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, errorDetalle)
}

// ObtenerError obtiene los detalles de un error
func (c *ErroresController) ObtenerError(ctx *gin.Context) {
	errorID := ctx.Param("id")
	if errorID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de error requerido"})
		return
	}

	errorDetalle, err := c.erroresService.ObtenerError(ctx.Request.Context(), errorID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Error no encontrado"})
		return
	}

	ctx.JSON(http.StatusOK, errorDetalle)
}

// ListarErrores lista los errores según los criterios especificados
func (c *ErroresController) ListarErrores(ctx *gin.Context) {
	var request struct {
		Tipo        models.TipoError      `json:"tipo"`
		Severidad   models.SeveridadError `json:"severidad"`
		Estado      string                `json:"estado"`
		FechaInicio time.Time             `json:"fecha_inicio"`
		FechaFin    time.Time             `json:"fecha_fin"`
		Limit       int                   `json:"limit" binding:"required"`
		Offset      int                   `json:"offset" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Construir filtro
	filter := bson.M{}
	if request.Tipo != "" {
		filter["tipo"] = request.Tipo
	}
	if request.Severidad != "" {
		filter["severidad"] = request.Severidad
	}
	if request.Estado != "" {
		filter["estado"] = request.Estado
	}
	if !request.FechaInicio.IsZero() && !request.FechaFin.IsZero() {
		filter["fecha_error"] = bson.M{
			"$gte": request.FechaInicio,
			"$lte": request.FechaFin,
		}
	}

	errores, total, err := c.erroresService.ListarErrores(
		ctx.Request.Context(),
		filter,
		request.Limit,
		request.Offset,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"total":   total,
		"errores": errores,
	})
}

// GenerarReporteErrores genera un reporte de errores
func (c *ErroresController) GenerarReporteErrores(ctx *gin.Context) {
	var request struct {
		FechaInicio time.Time `json:"fecha_inicio" binding:"required"`
		FechaFin    time.Time `json:"fecha_fin" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reporte, err := c.erroresService.GenerarReporteErrores(
		ctx.Request.Context(),
		request.FechaInicio,
		request.FechaFin,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, reporte)
}

// ObtenerIntentosRecuperacion obtiene los intentos de recuperación de un error
func (c *ErroresController) ObtenerIntentosRecuperacion(ctx *gin.Context) {
	errorID := ctx.Param("id")
	if errorID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de error requerido"})
		return
	}

	intentos, err := c.erroresService.ObtenerIntentosRecuperacion(ctx.Request.Context(), errorID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, intentos)
}

// ObtenerLogsError obtiene los logs de un error
func (c *ErroresController) ObtenerLogsError(ctx *gin.Context) {
	errorID := ctx.Param("id")
	if errorID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de error requerido"})
		return
	}

	logs, err := c.erroresService.ObtenerLogsError(ctx.Request.Context(), errorID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, logs)
}

// RegisterRoutes registra las rutas del controlador
func (c *ErroresController) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/registrar", c.RegistrarError)
	router.GET("/:id", c.ObtenerError)
	router.POST("/listar", c.ListarErrores)
	router.POST("/reportes", c.GenerarReporteErrores)
	router.GET("/:id/intentos", c.ObtenerIntentosRecuperacion)
	router.GET("/:id/logs", c.ObtenerLogsError)
}
