package routes

import (
	"time"

	"github.com/gin-gonic/gin"
	"FMgo/controllers"
	"FMgo/middleware"
	"FMgo/services"
	"gorm.io/gorm"
)

// SetupFacturaRoutes configura las rutas para las facturas electrónicas
func SetupFacturaRoutes(router *gin.Engine, db *gorm.DB, siiService *services.SIIService) {
	// Crear servicios
	facturaService := services.NewFacturaService(db, siiService)
	facturaController := controllers.NewFacturaController(facturaService)

	// Grupo de rutas para facturas
	facturas := router.Group("/api/facturas")
	{
		// Aplicar middlewares
		facturas.Use(middleware.AuthMiddleware("user", "admin"))
		facturas.Use(middleware.RateLimitMiddleware(100, time.Minute))

		// Rutas básicas
		facturas.POST("/", facturaController.CrearFactura)
		facturas.GET("/:id", facturaController.GetFactura)
		facturas.GET("/", facturaController.ListarFacturas)
		facturas.GET("/estado/:trackID/:rutEmisor", facturaController.ConsultarEstadoFactura)

		// Rutas avanzadas (requieren rol admin)
		admin := facturas.Group("")
		admin.Use(middleware.AuthMiddleware("admin"))
		{
			admin.POST("/:id/anular", facturaController.AnularFactura)
			admin.POST("/:id/reenviar", facturaController.ReenviarFactura)
		}

		// Rutas de PDF y email
		facturas.GET("/:id/pdf", facturaController.DescargarPDF)
		facturas.POST("/:id/email", facturaController.EnviarPorEmail)
	}
}
