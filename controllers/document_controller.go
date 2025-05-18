package controllers

import (
	"net/http"
	"strconv"

	"FMgo/domain"
	"FMgo/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DocumentController maneja las peticiones relacionadas con documentos
type DocumentController struct {
	docService domain.DocumentService
}

// NewDocumentController crea una nueva instancia del controlador de documentos
func NewDocumentController(docService domain.DocumentService) *DocumentController {
	return &DocumentController{
		docService: docService,
	}
}

// CrearDocumento maneja la creación de un documento
func (c *DocumentController) CrearDocumento(ctx *gin.Context) {
	var req models.FacturaRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convertir request a DocumentoTributario
	doc := &domain.DocumentoTributario{
		TipoDocumento:      "FACTURA",
		FechaEmision:       req.FechaEmision,
		RutEmisor:          req.RutEmisor,
		RutReceptor:        req.RutReceptor,
		MontoTotal:         0, // Se calculará
		MontoNeto:          0, // Se calculará
		MontoExento:        0, // Se calculará
		MontoIVA:           0, // Se calculará
		Estado:             "PENDIENTE",
		FechaCreacion:      req.FechaEmision,
		FechaActualizacion: req.FechaEmision,
	}

	// Calcular montos
	for _, item := range req.Items {
		doc.MontoNeto += item.MontoNeto
		doc.MontoIVA += item.MontoIVA
		doc.MontoTotal += item.MontoTotal
	}

	// Crear documento
	if err := c.docService.CrearDocumento(ctx.Request.Context(), doc); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"id":     doc.ID,
		"folio":  doc.Folio,
		"estado": doc.Estado,
	})
}

// ObtenerDocumento maneja la obtención de un documento
func (c *DocumentController) ObtenerDocumento(ctx *gin.Context) {
	tipo := ctx.Param("tipo")
	folioStr := ctx.Param("folio")

	folio, err := strconv.ParseInt(folioStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "folio inválido"})
		return
	}

	doc, err := c.docService.ObtenerDocumento(ctx.Request.Context(), tipo, folio)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if doc == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "documento no encontrado"})
		return
	}

	ctx.JSON(http.StatusOK, doc)
}

// ActualizarDocumento maneja la actualización de un documento
func (c *DocumentController) ActualizarDocumento(ctx *gin.Context) {
	var req models.FacturaRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Obtener el folio del parámetro de la URL
	folioStr := ctx.Param("folio")
	folio, err := strconv.ParseInt(folioStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "folio inválido"})
		return
	}

	// Obtener documento actual
	doc, err := c.docService.ObtenerDocumento(ctx.Request.Context(), "FACTURA", folio)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if doc == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "documento no encontrado"})
		return
	}

	// Actualizar campos
	doc.FechaEmision = req.FechaEmision
	doc.RutEmisor = req.RutEmisor
	doc.RutReceptor = req.RutReceptor
	doc.FechaActualizacion = req.FechaEmision

	// Recalcular montos
	doc.MontoTotal = 0
	doc.MontoNeto = 0
	doc.MontoExento = 0
	doc.MontoIVA = 0

	for _, item := range req.Items {
		doc.MontoNeto += item.MontoNeto
		doc.MontoIVA += item.MontoIVA
		doc.MontoTotal += item.MontoTotal
	}

	// Actualizar documento
	if err := c.docService.ActualizarDocumento(ctx.Request.Context(), doc); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":     doc.ID,
		"folio":  doc.Folio,
		"estado": doc.Estado,
	})
}

// CambiarEstadoDocumento maneja el cambio de estado de un documento
func (c *DocumentController) CambiarEstadoDocumento(ctx *gin.Context) {
	id := ctx.Param("id")
	estado := ctx.Param("estado")
	usuario := ctx.GetString("usuario") // Asumiendo que se obtiene del middleware de autenticación

	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	if err := c.docService.CambiarEstadoDocumento(ctx.Request.Context(), docID, estado, usuario); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Estado actualizado correctamente"})
}

// AgregarReferencia maneja la adición de una referencia a un documento
func (c *DocumentController) AgregarReferencia(ctx *gin.Context) {
	var ref domain.ReferenciaDocumento
	if err := ctx.ShouldBindJSON(&ref); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.docService.AgregarReferencia(ctx.Request.Context(), &ref); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"id": ref.ID,
	})
}

// ObtenerReferencias maneja la obtención de referencias de un documento
func (c *DocumentController) ObtenerReferencias(ctx *gin.Context) {
	tipoOrigen := ctx.Param("tipo")
	folioOrigenStr := ctx.Param("folio")

	folioOrigen, err := strconv.ParseInt(folioOrigenStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "folio inválido"})
		return
	}

	refs, err := c.docService.ObtenerReferencias(ctx.Request.Context(), tipoOrigen, folioOrigen)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, refs)
}
