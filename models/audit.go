package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AuditLog registra las acciones realizadas en el sistema
type AuditLog struct {
	ID         string    `gorm:"primaryKey" json:"id"`
	UserID     string    `gorm:"index" json:"user_id"`
	Rut        string    `gorm:"index" json:"rut"`
	Action     string    `gorm:"index" json:"action"`
	EntityType string    `gorm:"index" json:"entity_type"`
	EntityID   string    `gorm:"index" json:"entity_id"`
	OldValue   string    `json:"old_value"`
	NewValue   string    `json:"new_value"`
	IPAddress  string    `json:"ip_address"`
	UserAgent  string    `json:"user_agent"`
	CreatedAt  time.Time `json:"created_at"`
}

// AuditService define la interfaz para el servicio de auditoría
type AuditService interface {
	LogAction(userID, rut, action, entityType, entityID string, oldValue, newValue interface{}, ip, userAgent string) error
	GetLogsByUser(userID string, startDate, endDate time.Time) ([]*AuditLog, error)
	GetLogsByEntity(entityType, entityID string, startDate, endDate time.Time) ([]*AuditLog, error)
	GetLogsByRut(rut string, startDate, endDate time.Time) ([]*AuditLog, error)
}

// AuditServiceImpl implementa AuditService
type AuditServiceImpl struct {
	db *gorm.DB
}

// NewAuditService crea una nueva instancia del servicio de auditoría
func NewAuditService(db *gorm.DB) AuditService {
	return &AuditServiceImpl{db: db}
}

// LogAction registra una acción en el log de auditoría
func (s *AuditServiceImpl) LogAction(userID, rut, action, entityType, entityID string, oldValue, newValue interface{}, ip, userAgent string) error {
	log := &AuditLog{
		ID:         uuid.New().String(),
		UserID:     userID,
		Rut:        rut,
		Action:     action,
		EntityType: entityType,
		EntityID:   entityID,
		IPAddress:  ip,
		UserAgent:  userAgent,
		CreatedAt:  time.Now(),
	}

	// Convertir valores a JSON
	if oldValue != nil {
		if json, err := json.Marshal(oldValue); err == nil {
			log.OldValue = string(json)
		}
	}
	if newValue != nil {
		if json, err := json.Marshal(newValue); err == nil {
			log.NewValue = string(json)
		}
	}

	return s.db.Create(log).Error
}

// GetLogsByUser obtiene los logs de un usuario
func (s *AuditServiceImpl) GetLogsByUser(userID string, startDate, endDate time.Time) ([]*AuditLog, error) {
	var logs []*AuditLog
	err := s.db.Where("user_id = ? AND created_at BETWEEN ? AND ?", userID, startDate, endDate).
		Order("created_at DESC").
		Find(&logs).Error
	return logs, err
}

// GetLogsByEntity obtiene los logs de una entidad
func (s *AuditServiceImpl) GetLogsByEntity(entityType, entityID string, startDate, endDate time.Time) ([]*AuditLog, error) {
	var logs []*AuditLog
	err := s.db.Where("entity_type = ? AND entity_id = ? AND created_at BETWEEN ? AND ?", entityType, entityID, startDate, endDate).
		Order("created_at DESC").
		Find(&logs).Error
	return logs, err
}

// GetLogsByRut obtiene los logs de un RUT
func (s *AuditServiceImpl) GetLogsByRut(rut string, startDate, endDate time.Time) ([]*AuditLog, error) {
	var logs []*AuditLog
	err := s.db.Where("rut = ? AND created_at BETWEEN ? AND ?", rut, startDate, endDate).
		Order("created_at DESC").
		Find(&logs).Error
	return logs, err
}
