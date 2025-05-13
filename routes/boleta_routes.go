package routes

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/cursor/FMgo/controllers"
	"github.com/cursor/FMgo/middleware"
	"github.com/cursor/FMgo/repository"
	"github.com/cursor/FMgo/services"
)

// SetupBoletaRoutes configura las rutas para las boletas electrónicas
func SetupBoletaRoutes(router *gin.Engine, siiService *services.SIIService) {
	// Crear repositorio y servicios
	boletaRepo := repository.NewBoletaRepository()
	boletaService := services.NewBoletaService(siiService, boletaRepo)
	boletaController := controllers.NewBoletaController(boletaService)

	// Grupo de rutas para boletas
	boletas := router.Group("/api/boletas")
	{
		// Aplicar middlewares
		boletas.Use(middleware.AuthMiddleware("user", "admin"))
		boletas.Use(middleware.RateLimitMiddleware(100, time.Minute))

		// Rutas básicas
		boletas.POST("/", boletaController.CrearBoleta)
		boletas.GET("/:id", boletaController.GetBoleta)
		boletas.GET("/", boletaController.ListarBoletas)
		boletas.GET("/estado/:trackID/:rutEmisor", boletaController.ConsultarEstadoBoleta)

		// Rutas avanzadas (requieren rol admin)
		admin := boletas.Group("")
		admin.Use(middleware.AuthMiddleware("admin"))
		{
			admin.POST("/:id/anular", boletaController.AnularBoleta)
			admin.POST("/:id/reenviar", boletaController.ReenviarBoleta)
		}

		// Rutas de PDF y email
		boletas.GET("/:id/pdf", boletaController.DescargarPDF)
		boletas.POST("/:id/email", boletaController.EnviarPorEmail)
	}
}
