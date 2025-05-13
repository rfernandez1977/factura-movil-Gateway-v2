package controllers

import (
	"net/http"
	"time"

	"github.com/cursor/FMgo/models"
	"github.com/cursor/FMgo/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// TransformationController maneja las peticiones relacionadas con las transformaciones
type TransformationController struct {
	transformationService *services.TransformationService
}

// NewTransformationController crea una nueva instancia del controlador de transformación
func NewTransformationController(transformationService *services.TransformationService) *TransformationController {
	return &TransformationController{
		transformationService: transformationService,
	}
}

// RegistrarTransformacion registra una transformación
func (c *TransformationController) RegistrarTransformacion(ctx *gin.Context) {
	var transformacion models.Transformacion
	if err := ctx.ShouldBindJSON(&transformacion); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.transformationService.RegistrarTransformacion(ctx.Request.Context(), &transformacion); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"id": transformacion.ID})
}

// ObtenerTransformaciones obtiene las transformaciones
func (c *TransformationController) ObtenerTransformaciones(ctx *gin.Context) {
	var filtro struct {
		Tipo string `form:"tipo"`
	}

	if err := ctx.ShouldBindQuery(&filtro); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := bson.M{}
	if filtro.Tipo != "" {
		query["tipo"] = filtro.Tipo
	}

	transformaciones, err := c.transformationService.ObtenerTransformaciones(ctx.Request.Context(), query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, transformaciones)
}

// AplicarTransformacion aplica una transformación a los datos
func (c *TransformationController) AplicarTransformacion(ctx *gin.Context) {
	transformacionID, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de transformación inválido"})
		return
	}

	var datos map[string]interface{}
	if err := ctx.ShouldBindJSON(&datos); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resultado, err := c.transformationService.AplicarTransformacion(ctx.Request.Context(), transformacionID, datos)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resultado)
}

// RegistrarRegistroTransformacion registra un registro de transformación
func (c *TransformationController) RegistrarRegistroTransformacion(ctx *gin.Context) {
	var registro models.RegistroTransformacion
	if err := ctx.ShouldBindJSON(&registro); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.transformationService.RegistrarRegistroTransformacion(ctx.Request.Context(), &registro); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"id": registro.ID})
}

// ObtenerRegistrosTransformacion obtiene los registros de transformación
func (c *TransformationController) ObtenerRegistrosTransformacion(ctx *gin.Context) {
	var filtro struct {
		Inicio time.Time `form:"inicio" binding:"required"`
		Fin    time.Time `form:"fin" binding:"required"`
	}

	if err := ctx.ShouldBindQuery(&filtro); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	registros, err := c.transformationService.ObtenerRegistrosTransformacion(ctx.Request.Context(), map[string]interface{}{
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

// GenerarReporteTransformacion genera un reporte de transformación
func (c *TransformationController) GenerarReporteTransformacion(ctx *gin.Context) {
	var request struct {
		Inicio time.Time `json:"inicio" binding:"required"`
		Fin    time.Time `json:"fin" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reporte, err := c.transformationService.GenerarReporteTransformacion(ctx.Request.Context(), request.Inicio, request.Fin)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, reporte)
}

// RegisterRoutes registra las rutas del controlador
func (c *TransformationController) RegisterRoutes(router *gin.RouterGroup) {
	transformation := router.Group("/transformation")
	{
		transformation.POST("/transformaciones", c.RegistrarTransformacion)
		transformation.GET("/transformaciones", c.ObtenerTransformaciones)
		transformation.POST("/transformaciones/:id/aplicar", c.AplicarTransformacion)
		transformation.POST("/registros", c.RegistrarRegistroTransformacion)
		transformation.GET("/registros", c.ObtenerRegistrosTransformacion)
		transformation.POST("/reportes", c.GenerarReporteTransformacion)
	}
}
