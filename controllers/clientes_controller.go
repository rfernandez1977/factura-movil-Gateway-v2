package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/fmgo/models"
	"github.com/fmgo/services"
	"github.com/fmgo/utils"
	"go.uber.org/zap"
)

// ClientesController maneja las peticiones HTTP relacionadas con clientes
type ClientesController struct {
	clientesService *services.ClienteService
}

// NewClientesController crea una nueva instancia del controlador de clientes
func NewClientesController(clientesService *services.ClienteService) *ClientesController {
	return &ClientesController{
		clientesService: clientesService,
	}
}

// CrearCliente maneja la creación de un nuevo cliente
func (c *ClientesController) CrearCliente(ctx *gin.Context) {
	start := time.Now()

	var cliente models.Client
	if err := ctx.ShouldBindJSON(&cliente); err != nil {
		utils.LogError(err, zap.String("endpoint", "CrearCliente"))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.clientesService.CrearCliente(&cliente); err != nil {
		utils.LogError(err, zap.String("endpoint", "CrearCliente"))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.LogInfo("cliente creado exitosamente",
		zap.String("code", cliente.Code),
		zap.String("name", cliente.Name),
	)

	duration := time.Since(start).Seconds()
	utils.RecordHTTPRequest(
		ctx.Request.Method,
		ctx.Request.URL.Path,
		http.StatusCreated,
		duration,
		float64(ctx.Request.ContentLength),
		float64(cliente.ID),
	)

	ctx.JSON(http.StatusCreated, cliente)
}

// ObtenerCliente maneja la obtención de un cliente por ID
func (c *ClientesController) ObtenerCliente(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de cliente es requerido"})
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de cliente inválido"})
		return
	}

	cliente, err := c.clientesService.GetClienteByID(idInt)
	if err != nil {
		utils.LogError(err, zap.String("endpoint", "ObtenerCliente"))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, cliente)
}

// ActualizarCliente maneja la actualización de un cliente
func (c *ClientesController) ActualizarCliente(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de cliente es requerido"})
		return
	}

	var cliente models.Client
	if err := ctx.ShouldBindJSON(&cliente); err != nil {
		utils.LogError(err, zap.String("endpoint", "ActualizarCliente"))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cliente.ID, _ = strconv.Atoi(id)
	if err := c.clientesService.ActualizarCliente(&cliente); err != nil {
		utils.LogError(err, zap.String("endpoint", "ActualizarCliente"))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.LogInfo("cliente actualizado exitosamente",
		zap.Int("id", cliente.ID),
		zap.String("code", cliente.Code),
	)

	ctx.JSON(http.StatusOK, cliente)
}

// EliminarCliente maneja la eliminación de un cliente
func (c *ClientesController) EliminarCliente(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de cliente es requerido"})
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de cliente inválido"})
		return
	}

	if err := c.clientesService.EliminarCliente(idInt); err != nil {
		utils.LogError(err, zap.String("endpoint", "EliminarCliente"))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.LogInfo("cliente eliminado exitosamente", zap.String("id", id))
	ctx.JSON(http.StatusOK, gin.H{"message": "Cliente eliminado exitosamente"})
}

// ListarClientes maneja la obtención de una lista de clientes
/*
func (c *ClientesController) ListarClientes(ctx *gin.Context) {
	empresaID := ctx.Query("empresa_id")
	if empresaID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de empresa es requerido"})
		return
	}

	clientes, err := c.clientesService.ListarClientes(empresaID)
	if err != nil {
		utils.LogError(err, zap.String("endpoint", "ListarClientes"))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, clientes)
}
*/
// TODO: Implementar ListarClientes en el servicio

// BuscarClientes maneja la búsqueda de clientes por término
/*
func (c *ClientesController) BuscarClientes(ctx *gin.Context) {
	termino := ctx.Query("q")
	if termino == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Término de búsqueda es requerido"})
		return
	}

	empresaID := ctx.Query("empresa_id")
	if empresaID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de empresa es requerido"})
		return
	}

	clientes, err := c.clientesService.BuscarClientes(empresaID, termino)
	if err != nil {
		utils.LogError(err, zap.String("endpoint", "BuscarClientes"))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, clientes)
}
*/
// TODO: Implementar BuscarClientes en el servicio

// RegisterRoutes registra las rutas del controlador
func (c *ClientesController) RegisterRoutes(router *gin.RouterGroup) {
	clientes := router.Group("/clientes")
	{
		clientes.POST("", c.CrearCliente)
		clientes.GET("/:id", c.ObtenerCliente)
		clientes.PUT("/:id", c.ActualizarCliente)
		clientes.DELETE("/:id", c.EliminarCliente)
		// clientes.GET("", c.ListarClientes) // TODO: Implementar
		// clientes.GET("/buscar", c.BuscarClientes) // TODO: Implementar
	}
}
