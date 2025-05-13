package controllers

import (
	"encoding/base64"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/cursor/FMgo/services"
)

// SeguridadController maneja las peticiones relacionadas con la seguridad
type SeguridadController struct {
	seguridadService *services.SeguridadService
}

// NewSeguridadController crea una nueva instancia del controlador de seguridad
func NewSeguridadController(seguridadService *services.SeguridadService) *SeguridadController {
	return &SeguridadController{
		seguridadService: seguridadService,
	}
}

// RegistrarAcceso registra un intento de acceso al sistema
func (c *SeguridadController) RegistrarAcceso(ctx *gin.Context) {
	var request struct {
		UsuarioID string `json:"usuario_id" binding:"required"`
		Rut       string `json:"rut" binding:"required"`
		Accion    string `json:"accion" binding:"required"`
		IP        string `json:"ip" binding:"required"`
		UserAgent string `json:"user_agent" binding:"required"`
		Exitoso   bool   `json:"exitoso" binding:"required"`
		Detalles  string `json:"detalles"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := c.seguridadService.RegistrarAcceso(
		ctx.Request.Context(),
		request.UsuarioID,
		request.Rut,
		request.Accion,
		request.IP,
		request.UserAgent,
		request.Exitoso,
		request.Detalles,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Acceso registrado correctamente"})
}

// RegistrarOperacion registra una operación en el sistema
func (c *SeguridadController) RegistrarOperacion(ctx *gin.Context) {
	var request struct {
		UsuarioID      string                 `json:"usuario_id" binding:"required"`
		Rut            string                 `json:"rut" binding:"required"`
		Operacion      string                 `json:"operacion" binding:"required"`
		Entidad        string                 `json:"entidad" binding:"required"`
		EntidadID      string                 `json:"entidad_id" binding:"required"`
		Cambios        map[string]interface{} `json:"cambios" binding:"required"`
		EstadoAnterior map[string]interface{} `json:"estado_anterior"`
		EstadoNuevo    map[string]interface{} `json:"estado_nuevo"`
		IP             string                 `json:"ip" binding:"required"`
		UserAgent      string                 `json:"user_agent" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := c.seguridadService.RegistrarOperacion(
		ctx.Request.Context(),
		request.UsuarioID,
		request.Rut,
		request.Operacion,
		request.Entidad,
		request.EntidadID,
		request.Cambios,
		request.EstadoAnterior,
		request.EstadoNuevo,
		request.IP,
		request.UserAgent,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Operación registrada correctamente"})
}

// ValidarFirmaDigital valida una firma digital
func (c *SeguridadController) ValidarFirmaDigital(ctx *gin.Context) {
	var request struct {
		UsuarioID string `json:"usuario_id" binding:"required"`
		Documento string `json:"documento" binding:"required"`
		Firma     string `json:"firma" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	documento, err := base64.StdEncoding.DecodeString(request.Documento)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "documento inválido"})
		return
	}

	firma, err := base64.StdEncoding.DecodeString(request.Firma)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "firma inválida"})
		return
	}

	valido, err := c.seguridadService.ValidarFirmaDigital(
		ctx.Request.Context(),
		request.UsuarioID,
		documento,
		firma,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"valido": valido})
}

// EncriptarDatos encripta datos sensibles
func (c *SeguridadController) EncriptarDatos(ctx *gin.Context) {
	var request struct {
		Entidad   string `json:"entidad" binding:"required"`
		EntidadID string `json:"entidad_id" binding:"required"`
		Campo     string `json:"campo" binding:"required"`
		Valor     string `json:"valor" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	datos, err := c.seguridadService.EncriptarDatos(
		ctx.Request.Context(),
		request.Entidad,
		request.EntidadID,
		request.Campo,
		[]byte(request.Valor),
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, datos)
}

// DesencriptarDatos desencripta datos sensibles
func (c *SeguridadController) DesencriptarDatos(ctx *gin.Context) {
	var request struct {
		Entidad   string `json:"entidad" binding:"required"`
		EntidadID string `json:"entidad_id" binding:"required"`
		Campo     string `json:"campo" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	datos, err := c.seguridadService.DesencriptarDatos(
		ctx.Request.Context(),
		request.Entidad,
		request.EntidadID,
		request.Campo,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"valor": string(datos)})
}

// GenerarReporteSeguridad genera un reporte de seguridad
func (c *SeguridadController) GenerarReporteSeguridad(ctx *gin.Context) {
	var request struct {
		FechaInicio time.Time `json:"fecha_inicio" binding:"required"`
		FechaFin    time.Time `json:"fecha_fin" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reporte, err := c.seguridadService.GenerarReporteSeguridad(
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

// RegisterRoutes registra las rutas del controlador
func (c *SeguridadController) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/accesos", c.RegistrarAcceso)
	router.POST("/operaciones", c.RegistrarOperacion)
	router.POST("/firmas/validar", c.ValidarFirmaDigital)
	router.POST("/datos/encriptar", c.EncriptarDatos)
	router.POST("/datos/desencriptar", c.DesencriptarDatos)
	router.POST("/reportes", c.GenerarReporteSeguridad)
}
