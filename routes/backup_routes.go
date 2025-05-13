package routes

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/cursor/FMgo/controllers"
	"github.com/cursor/FMgo/middleware"
	"github.com/cursor/FMgo/services"
	"gorm.io/gorm"
)

// SetupBackupRoutes configura las rutas para el manejo de respaldos
func SetupBackupRoutes(router *gin.Engine, db *gorm.DB) {
	// Crear servicios
	backupService := services.NewBackupService(db, "backups", 30) // 30 días de retención
	backupController := controllers.NewBackupController(backupService)

	// Grupo de rutas para respaldos
	backupGroup := router.Group("/api/backups")
	backupGroup.Use(middleware.AuthMiddleware("admin"))
	backupGroup.Use(middleware.RateLimitMiddleware(10, time.Minute)) // 10 peticiones por minuto

	// Rutas básicas
	backupGroup.POST("/", backupController.CreateBackup)
	backupGroup.POST("/restore", backupController.RestoreBackup)
	backupGroup.GET("/", backupController.ListBackups)
}
