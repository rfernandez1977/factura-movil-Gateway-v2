package caf

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/fmgo/config"
	"github.com/fmgo/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestBackupDir(t *testing.T) string {
	dir, err := os.MkdirTemp("", "caf_backup_test")
	require.NoError(t, err)
	t.Cleanup(func() {
		os.RemoveAll(dir)
	})
	return dir
}

func TestBackupCAF(t *testing.T) {
	backupDir := setupTestBackupDir(t)
	cfg := &config.Config{}
	svc := NewService(cfg, nil, nil, backupDir, 30)

	t.Run("backup exitoso", func(t *testing.T) {
		caf := &models.CAF{
			ID:            "TEST-123",
			EmpresaID:     "EMP-001",
			TipoDocumento: "33",
			FolioInicial:  1,
			FolioFinal:    100,
			Archivo:       []byte("contenido de prueba"),
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		result, err := svc.BackupCAF(context.Background(), caf)
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, caf.ID, result.CAFId)
		assert.FileExists(t, result.Ubicacion)
	})

	t.Run("backup con CAF nulo", func(t *testing.T) {
		result, err := svc.BackupCAF(context.Background(), nil)
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestRestoreCAF(t *testing.T) {
	backupDir := setupTestBackupDir(t)
	cfg := &config.Config{}
	svc := NewService(cfg, nil, nil, backupDir, 30)

	t.Run("restauración con backup inexistente", func(t *testing.T) {
		caf, err := svc.RestoreCAF(context.Background(), "NOEXISTE")
		assert.Error(t, err)
		assert.Nil(t, caf)
	})

	t.Run("restauración con backup corrupto", func(t *testing.T) {
		// Crear archivo de backup corrupto
		backupPath := filepath.Join(backupDir, "corrupto.json")
		err := os.WriteFile(backupPath, []byte("datos corruptos"), 0644)
		require.NoError(t, err)

		caf, err := svc.RestoreCAF(context.Background(), "corrupto")
		assert.Error(t, err)
		assert.Nil(t, caf)
	})
}

func TestCleanOldBackups(t *testing.T) {
	backupDir := setupTestBackupDir(t)
	cfg := &config.Config{}
	svc := NewService(cfg, nil, nil, backupDir, 7) // 7 días de retención

	t.Run("limpieza exitosa", func(t *testing.T) {
		// Crear archivos de prueba
		oldFile := filepath.Join(backupDir, "old_backup.json")
		newFile := filepath.Join(backupDir, "new_backup.json")

		require.NoError(t, os.WriteFile(oldFile, []byte("old"), 0644))
		require.NoError(t, os.WriteFile(newFile, []byte("new"), 0644))

		// Modificar tiempo de oldFile a 8 días atrás
		oldTime := time.Now().AddDate(0, 0, -8)
		require.NoError(t, os.Chtimes(oldFile, oldTime, oldTime))

		// Ejecutar limpieza
		err := svc.CleanOldBackups(context.Background())
		require.NoError(t, err)

		// Verificar resultados
		assert.NoFileExists(t, oldFile)
		assert.FileExists(t, newFile)
	})

	t.Run("limpieza con directorio inexistente", func(t *testing.T) {
		svc := NewService(cfg, nil, nil, "/no/existe", 7)
		err := svc.CleanOldBackups(context.Background())
		assert.Error(t, err)
	})
}
