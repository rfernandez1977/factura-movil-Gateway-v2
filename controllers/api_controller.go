package controllers

import (
	"net/http"
	"time"

	"FMgo/models"
	"FMgo/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// APIController maneja las peticiones relacionadas con las APIs
type APIController struct {
	apiService *services.APIService
}

// NewAPIController crea una nueva instancia del controlador de API
func NewAPIController(apiService *services.APIService) *APIController {
	return &APIController{
		apiService: apiService,
	}
}

// RegistrarAPI registra una API
func (c *APIController) RegistrarAPI(ctx *gin.Context) {
	var api models.API
	if err := ctx.ShouldBindJSON(&api); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.apiService.RegistrarAPI(ctx.Request.Context(), &api); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"id": api.ID})
}

// ObtenerAPIs obtiene las APIs
func (c *APIController) ObtenerAPIs(ctx *gin.Context) {
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

	apis, err := c.apiService.ObtenerAPIs(ctx.Request.Context(), query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, apis)
}

// RegistrarVersionAPI registra una versión de API
func (c *APIController) RegistrarVersionAPI(ctx *gin.Context) {
	var version models.VersionAPI
	if err := ctx.ShouldBindJSON(&version); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.apiService.RegistrarVersionAPI(ctx.Request.Context(), &version); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"id": version.ID})
}

// ObtenerVersionesAPI obtiene las versiones de una API
func (c *APIController) ObtenerVersionesAPI(ctx *gin.Context) {
	apiID, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de API inválido"})
		return
	}

	versiones, err := c.apiService.ObtenerVersionesAPI(ctx.Request.Context(), apiID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, versiones)
}

// RegistrarRegistroAPI registra un registro de API
func (c *APIController) RegistrarRegistroAPI(ctx *gin.Context) {
	var registro models.RegistroAPI
	if err := ctx.ShouldBindJSON(&registro); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.apiService.RegistrarRegistroAPI(ctx.Request.Context(), &registro); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"id": registro.ID})
}

// ObtenerRegistrosAPI obtiene los registros de API
func (c *APIController) ObtenerRegistrosAPI(ctx *gin.Context) {
	var filtro struct {
		Inicio time.Time `form:"inicio" binding:"required"`
		Fin    time.Time `form:"fin" binding:"required"`
	}

	if err := ctx.ShouldBindQuery(&filtro); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	registros, err := c.apiService.ObtenerRegistrosAPI(ctx.Request.Context(), map[string]interface{}{
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

// GenerarReporteAPI genera un reporte de API
func (c *APIController) GenerarReporteAPI(ctx *gin.Context) {
	var request struct {
		Inicio time.Time `json:"inicio" binding:"required"`
		Fin    time.Time `json:"fin" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reporte, err := c.apiService.GenerarReporteAPI(ctx.Request.Context(), request.Inicio, request.Fin)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, reporte)
}

// RegisterRoutes registra las rutas del controlador
func (c *APIController) RegisterRoutes(router *gin.RouterGroup) {
	api := router.Group("/api")
	{
		api.POST("/apis", c.RegistrarAPI)
		api.GET("/apis", c.ObtenerAPIs)
		api.POST("/versiones", c.RegistrarVersionAPI)
		api.GET("/apis/:id/versiones", c.ObtenerVersionesAPI)
		api.POST("/registros", c.RegistrarRegistroAPI)
		api.GET("/registros", c.ObtenerRegistrosAPI)
		api.POST("/reportes", c.GenerarReporteAPI)
	}
}
