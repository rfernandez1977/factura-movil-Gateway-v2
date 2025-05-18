package controllers

import (
	"net/http"
	"os"
	"path/filepath"
	"time"

	"FMgo/services"
	"FMgo/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// BackupController maneja las peticiones HTTP relacionadas con respaldos
type BackupController struct {
	backupService *services.BackupService
}

// NewBackupController crea una nueva instancia del controlador de respaldos
func NewBackupController(backupService *services.BackupService) *BackupController {
	return &BackupController{
		backupService: backupService,
	}
}

// CreateBackup crea un nuevo respaldo
func (c *BackupController) CreateBackup(ctx *gin.Context) {
	start := time.Now()

	// Verificar permisos
	jwtUtils := utils.NewJWTUtils()
	userID, _ := utils.GetUserID(ctx.GetHeader("Authorization"), jwtUtils)
	ok, _ := utils.HasRole(ctx.GetHeader("Authorization"), "admin", jwtUtils)
	if !ok {
		utils.LogWarning("intento de crear respaldo sin permisos",
			zap.String("user_id", userID),
		)
		ctx.JSON(http.StatusForbidden, gin.H{"error": "acceso denegado"})
		return
	}

	// Crear respaldo
	if err := c.backupService.CreateBackup(); err != nil {
		utils.LogError(err, zap.String("endpoint", "CreateBackup"))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Registrar éxito
	utils.LogInfo("respaldo creado exitosamente",
		zap.String("user_id", userID),
	)

	// Registrar métricas HTTP
	duration := time.Since(start).Seconds()
	utils.RecordHTTPRequest(
		ctx.Request.Method,
		ctx.Request.URL.Path,
		http.StatusOK,
		duration,
		float64(ctx.Request.ContentLength),
		0,
	)

	ctx.JSON(http.StatusOK, gin.H{"message": "respaldo creado exitosamente"})
}

// RestoreBackup restaura un respaldo
func (c *BackupController) RestoreBackup(ctx *gin.Context) {
	start := time.Now()

	// Verificar permisos
	jwtUtils := utils.NewJWTUtils()
	userID, _ := utils.GetUserID(ctx.GetHeader("Authorization"), jwtUtils)
	ok, _ := utils.HasRole(ctx.GetHeader("Authorization"), "admin", jwtUtils)
	if !ok {
		utils.LogWarning("intento de restaurar respaldo sin permisos",
			zap.String("user_id", userID),
		)
		ctx.JSON(http.StatusForbidden, gin.H{"error": "acceso denegado"})
		return
	}

	// Obtener archivo de respaldo
	file, err := ctx.FormFile("backup")
	if err != nil {
		utils.LogError(err, zap.String("endpoint", "RestoreBackup"))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "archivo de respaldo no proporcionado"})
		return
	}

	// Guardar archivo temporalmente
	tempFile := filepath.Join(os.TempDir(), file.Filename)
	if err := ctx.SaveUploadedFile(file, tempFile); err != nil {
		utils.LogError(err, zap.String("endpoint", "RestoreBackup"))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error al guardar archivo"})
		return
	}
	defer os.Remove(tempFile)

	// Restaurar respaldo
	if err := c.backupService.RestoreBackup(tempFile); err != nil {
		utils.LogError(err, zap.String("endpoint", "RestoreBackup"))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Registrar éxito
	utils.LogInfo("respaldo restaurado exitosamente",
		zap.String("user_id", userID),
		zap.String("file", file.Filename),
	)

	// Registrar métricas HTTP
	duration := time.Since(start).Seconds()
	utils.RecordHTTPRequest(
		ctx.Request.Method,
		ctx.Request.URL.Path,
		http.StatusOK,
		duration,
		float64(ctx.Request.ContentLength),
		0,
	)

	ctx.JSON(http.StatusOK, gin.H{"message": "respaldo restaurado exitosamente"})
}

// ListBackups lista los respaldos disponibles
func (c *BackupController) ListBackups(ctx *gin.Context) {
	start := time.Now()

	// Verificar permisos
	jwtUtils := utils.NewJWTUtils()
	userID, _ := utils.GetUserID(ctx.GetHeader("Authorization"), jwtUtils)
	ok, _ := utils.HasRole(ctx.GetHeader("Authorization"), "admin", jwtUtils)
	if !ok {
		utils.LogWarning("intento de listar respaldos sin permisos",
			zap.String("user_id", userID),
		)
		ctx.JSON(http.StatusForbidden, gin.H{"error": "acceso denegado"})
		return
	}

	// Obtenemos los archivos de backup a través del método GetBackupFiles que debemos implementar en el servicio
	backups, err := c.GetBackupFiles()
	if err != nil {
		utils.LogError(err, zap.String("endpoint", "ListBackups"))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Registrar métricas HTTP
	duration := time.Since(start).Seconds()
	utils.RecordHTTPRequest(
		ctx.Request.Method,
		ctx.Request.URL.Path,
		http.StatusOK,
		duration,
		float64(ctx.Request.ContentLength),
		float64(len(backups)),
	)

	ctx.JSON(http.StatusOK, backups)
}

// GetBackupFiles obtiene la lista de archivos de respaldo
func (c *BackupController) GetBackupFiles() ([]gin.H, error) {
	// Obtener lista de archivos del servicio
	files, err := c.backupService.GetBackupFiles()
	if err != nil {
		return nil, err
	}

	// Obtener información de cada archivo
	backups := make([]gin.H, len(files))
	for i, file := range files {
		info, err := os.Stat(file)
		if err != nil {
			continue
		}

		backups[i] = gin.H{
			"filename": filepath.Base(file),
			"size":     info.Size(),
			"date":     info.ModTime(),
		}
	}

	return backups, nil
}
