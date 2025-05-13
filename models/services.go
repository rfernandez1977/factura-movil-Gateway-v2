package models

import (
	"time"
)

// DTEService representa el servicio de gestión de DTE
type DTEService struct {
	ID          string    `json:"id"`
	Nombre      string    `json:"nombre"`
	Descripcion string    `json:"descripcion"`
	Version     string    `json:"version"`
	Estado      string    `json:"estado"`
	Config      Config    `json:"config"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// DTEGenerator representa el servicio de generación de DTE
type DTEGenerator struct {
	ID          string    `json:"id"`
	Nombre      string    `json:"nombre"`
	Descripcion string    `json:"descripcion"`
	Version     string    `json:"version"`
	Estado      string    `json:"estado"`
	Config      Config    `json:"config"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// ValidationService representa el servicio de validación
type ValidationService struct {
	ID          string    `json:"id"`
	Nombre      string    `json:"nombre"`
	Descripcion string    `json:"descripcion"`
	Version     string    `json:"version"`
	Estado      string    `json:"estado"`
	Config      Config    `json:"config"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// NotificationService representa el servicio de notificaciones
type NotificationService struct {
	ID          string    `json:"id"`
	Nombre      string    `json:"nombre"`
	Descripcion string    `json:"descripcion"`
	Version     string    `json:"version"`
	Estado      string    `json:"estado"`
	Config      Config    `json:"config"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// StorageService representa el servicio de almacenamiento
type StorageService struct {
	ID          string    `json:"id"`
	Nombre      string    `json:"nombre"`
	Descripcion string    `json:"descripcion"`
	Version     string    `json:"version"`
	Estado      string    `json:"estado"`
	Config      Config    `json:"config"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// ERPService representa el servicio de integración con ERP
type ERPService struct {
	ID          string    `json:"id"`
	Nombre      string    `json:"nombre"`
	Descripcion string    `json:"descripcion"`
	Version     string    `json:"version"`
	Estado      string    `json:"estado"`
	Config      Config    `json:"config"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// MonitoringService representa el servicio de monitoreo
type MonitoringService struct {
	ID          string    `json:"id"`
	Nombre      string    `json:"nombre"`
	Descripcion string    `json:"descripcion"`
	Version     string    `json:"version"`
	Estado      string    `json:"estado"`
	Config      Config    `json:"config"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// Funciones auxiliares para crear nuevas instancias de servicios
func NewDTEService(nombre, descripcion, version string, config Config) *DTEService {
	return &DTEService{
		ID:          GenerateID(),
		Nombre:      nombre,
		Descripcion: descripcion,
		Version:     version,
		Estado:      "ACTIVO",
		Config:      config,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func NewDTEGenerator(nombre, descripcion, version string, config Config) *DTEGenerator {
	return &DTEGenerator{
		ID:          GenerateID(),
		Nombre:      nombre,
		Descripcion: descripcion,
		Version:     version,
		Estado:      "ACTIVO",
		Config:      config,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func NewValidationService(nombre, descripcion, version string, config Config) *ValidationService {
	return &ValidationService{
		ID:          GenerateID(),
		Nombre:      nombre,
		Descripcion: descripcion,
		Version:     version,
		Estado:      "ACTIVO",
		Config:      config,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func NewNotificationService(nombre, descripcion, version string, config Config) *NotificationService {
	return &NotificationService{
		ID:          GenerateID(),
		Nombre:      nombre,
		Descripcion: descripcion,
		Version:     version,
		Estado:      "ACTIVO",
		Config:      config,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func NewStorageService(nombre, descripcion, version string, config Config) *StorageService {
	return &StorageService{
		ID:          GenerateID(),
		Nombre:      nombre,
		Descripcion: descripcion,
		Version:     version,
		Estado:      "ACTIVO",
		Config:      config,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func NewERPService(nombre, descripcion, version string, config Config) *ERPService {
	return &ERPService{
		ID:          GenerateID(),
		Nombre:      nombre,
		Descripcion: descripcion,
		Version:     version,
		Estado:      "ACTIVO",
		Config:      config,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func NewMonitoringService(nombre, descripcion, version string, config Config) *MonitoringService {
	return &MonitoringService{
		ID:          GenerateID(),
		Nombre:      nombre,
		Descripcion: descripcion,
		Version:     version,
		Estado:      "ACTIVO",
		Config:      config,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}
