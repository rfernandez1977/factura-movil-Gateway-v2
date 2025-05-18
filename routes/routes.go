package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"FMgo/controllers"
	"FMgo/services"
)

// SetupRouter configura las rutas de la aplicación
func SetupRouter(siiService *services.SIIService, signatureService *services.SignatureService) *gin.Engine {
	router := gin.Default()

	// Controladores
	docController := controllers.NewDocumentController(siiService, signatureService)

	// Grupo de rutas para documentos
	docGroup := router.Group("/api/documentos")
	{
		docGroup.POST("/enviar", docController.EnviarDocumentoHandler)
		docGroup.POST("/consultar", docController.ConsultarEstadoHandler)
		docGroup.POST("/sobre/enviar", docController.EnviarSobreDTEHandler)
		docGroup.GET("/sobre/estado/:trackId", docController.ConsultarEstadoSobreHandler)
	}

	// Métricas Prometheus
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	return router
}
