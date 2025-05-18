package filesystem

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"FMgo/core/firma/common"
	"FMgo/core/firma/models"
)

// CertificadoStorage implementa el almacenamiento de certificados en sistema de archivos
type CertificadoStorage struct {
	basePath string
	mu       sync.RWMutex
	logger   common.Logger
}

// NewCertificadoStorage crea una nueva instancia de almacenamiento de certificados
func NewCertificadoStorage(basePath string, logger common.Logger) (*CertificadoStorage, error) {
	// Crear directorio base si no existe
	if err := os.MkdirAll(basePath, 0750); err != nil {
		return nil, fmt.Errorf("error creando directorio base: %w", err)
	}

	return &CertificadoStorage{
		basePath: basePath,
		logger:   logger,
	}, nil
}

// GuardarCertificado guarda un certificado en el sistema de archivos
func (s *CertificadoStorage) GuardarCertificado(ctx context.Context, cert *models.Certificado) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Crear directorio para el certificado
	certPath := filepath.Join(s.basePath, cert.ID)
	if err := os.MkdirAll(certPath, 0750); err != nil {
		return fmt.Errorf("error creando directorio del certificado: %w", err)
	}

	// Guardar metadatos
	metadataPath := filepath.Join(certPath, "metadata.json")
	metadata, err := json.MarshalIndent(cert, "", "  ")
	if err != nil {
		return fmt.Errorf("error serializando metadatos: %w", err)
	}

	if err := os.WriteFile(metadataPath, metadata, 0600); err != nil {
		return fmt.Errorf("error guardando metadatos: %w", err)
	}

	// Guardar certificado PEM
	certPEMPath := filepath.Join(certPath, "certificado.pem")
	if err := os.WriteFile(certPEMPath, cert.CertificadoPEM, 0600); err != nil {
		return fmt.Errorf("error guardando certificado PEM: %w", err)
	}

	// Guardar llave privada PEM
	keyPEMPath := filepath.Join(certPath, "llave.pem")
	if err := os.WriteFile(keyPEMPath, cert.LlavePrivadaPEM, 0600); err != nil {
		return fmt.Errorf("error guardando llave privada: %w", err)
	}

	s.logger.Info("Certificado guardado exitosamente",
		"id", cert.ID,
		"rut", cert.RutEmpresa,
		"path", certPath)

	return nil
}

// ObtenerCertificado obtiene un certificado del sistema de archivos
func (s *CertificadoStorage) ObtenerCertificado(ctx context.Context, id string) (*models.Certificado, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	certPath := filepath.Join(s.basePath, id)

	// Leer metadatos
	metadataPath := filepath.Join(certPath, "metadata.json")
	metadata, err := os.ReadFile(metadataPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("certificado no encontrado: %s", id)
		}
		return nil, fmt.Errorf("error leyendo metadatos: %w", err)
	}

	var cert models.Certificado
	if err := json.Unmarshal(metadata, &cert); err != nil {
		return nil, fmt.Errorf("error deserializando metadatos: %w", err)
	}

	// Leer certificado PEM
	certPEMPath := filepath.Join(certPath, "certificado.pem")
	certPEM, err := os.ReadFile(certPEMPath)
	if err != nil {
		return nil, fmt.Errorf("error leyendo certificado PEM: %w", err)
	}
	cert.CertificadoPEM = certPEM

	// Leer llave privada PEM
	keyPEMPath := filepath.Join(certPath, "llave.pem")
	keyPEM, err := os.ReadFile(keyPEMPath)
	if err != nil {
		return nil, fmt.Errorf("error leyendo llave privada: %w", err)
	}
	cert.LlavePrivadaPEM = keyPEM

	return &cert, nil
}

// ListarCertificados lista todos los certificados
func (s *CertificadoStorage) ListarCertificados(ctx context.Context) ([]*models.Certificado, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entries, err := os.ReadDir(s.basePath)
	if err != nil {
		return nil, fmt.Errorf("error listando directorio: %w", err)
	}

	var certificados []*models.Certificado
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		cert, err := s.ObtenerCertificado(ctx, entry.Name())
		if err != nil {
			s.logger.Warn("Error obteniendo certificado",
				"id", entry.Name(),
				"error", err)
			continue
		}

		certificados = append(certificados, cert)
	}

	return certificados, nil
}

// EliminarCertificado elimina un certificado
func (s *CertificadoStorage) EliminarCertificado(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	certPath := filepath.Join(s.basePath, id)
	if err := os.RemoveAll(certPath); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("certificado no encontrado: %s", id)
		}
		return fmt.Errorf("error eliminando certificado: %w", err)
	}

	s.logger.Info("Certificado eliminado exitosamente",
		"id", id,
		"path", certPath)

	return nil
}
