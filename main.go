package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/usuario/gateway/db"
	"github.com/usuario/gateway/handlers"
	"github.com/usuario/gateway/middleware"
	"github.com/usuario/gateway/metrics"
)

func main() {
	// Initialize database connection
	dbConn, err := db.InitDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer dbConn.Close()

	// Initialize Prometheus metrics
	metrics.InitMetrics()

	// Set up router
	r := gin.Default()

	// Apply API key middleware for authentication
	r.Use(middleware.APIKeyMiddleware())

	// Prometheus metrics endpoint
	r.GET("/metrics", metrics.MetricsHandler())

	// Endpoints for document creation
	r.POST("/facturas", handlers.CreateDocumentHandler(dbConn, "factura"))
	r.POST("/boletas", handlers.CreateDocumentHandler(dbConn, "boleta"))
	r.POST("/notas", handlers.CreateDocumentHandler(dbConn, "nota"))
	r.POST("/guias", handlers.CreateDocumentHandler(dbConn, "guia"))

	// Endpoints for entity creation
	r.POST("/clientes", handlers.CreateEntityHandler("cliente"))
	r.POST("/productos", handlers.CreateEntityHandler("producto"))

	// Endpoints for document queries
	r.GET("/documents/:id", handlers.GetDocumentStatusHandler)
	r.GET("/documents/:id/pdf", handlers.GetDocumentPDFHandler)

	// Start the server
	if err := r.Run(":3000"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}