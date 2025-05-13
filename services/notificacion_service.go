package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/cursor/FMgo/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NotificacionService maneja el sistema de notificaciones
type NotificacionService struct {
	db                 *mongo.Database
	config             *models.SistemaNotificaciones
	emailService       *EmailService
	smsService         *SMSService
	pushService        *PushService
	colaNotificaciones chan *models.Notificacion
}

// NewNotificacionService crea una nueva instancia del servicio de notificaciones
func NewNotificacionService(db *mongo.Database, emailService *EmailService, smsService *SMSService, pushService *PushService) *NotificacionService {
	service := &NotificacionService{
		db:                 db,
		emailService:       emailService,
		smsService:         smsService,
		pushService:        pushService,
		colaNotificaciones: make(chan *models.Notificacion, 1000),
		config: &models.SistemaNotificaciones{
			MaxIntentosEnvio:   3,
			IntervaloReintento: 5 * time.Minute,
			NotificarEmail:     true,
			NotificarSMS:       true,
			NotificarPush:      true,
			PlantillasEmail: map[string]string{
				"ALERTA_DOCUMENTO":  "Plantilla de alerta de documento",
				"ESTADO_DOCUMENTO":  "Plantilla de estado de documento",
				"VENCIMIENTO":       "Plantilla de vencimiento",
				"ERROR_SISTEMA":     "Plantilla de error de sistema",
				"ACTUALIZACION_SII": "Plantilla de actualización SII",
			},
			PlantillasSMS: map[string]string{
				"ALERTA_DOCUMENTO":  "Plantilla SMS de alerta",
				"ESTADO_DOCUMENTO":  "Plantilla SMS de estado",
				"VENCIMIENTO":       "Plantilla SMS de vencimiento",
				"ERROR_SISTEMA":     "Plantilla SMS de error",
				"ACTUALIZACION_SII": "Plantilla SMS de actualización",
			},
			PlantillasPush: map[string]string{
				"ALERTA_DOCUMENTO":  "Plantilla Push de alerta",
				"ESTADO_DOCUMENTO":  "Plantilla Push de estado",
				"VENCIMIENTO":       "Plantilla Push de vencimiento",
				"ERROR_SISTEMA":     "Plantilla Push de error",
				"ACTUALIZACION_SII": "Plantilla Push de actualización",
			},
		},
	}

	// Iniciar procesador de notificaciones
	go service.procesarNotificaciones()

	return service
}

// CrearNotificacion crea una nueva notificación
func (s *NotificacionService) CrearNotificacion(ctx context.Context, notificacion *models.Notificacion) error {
	// Validar notificación
	if err := s.validarNotificacion(notificacion); err != nil {
		return err
	}

	// Asignar ID y fechas
	notificacion.ID = primitive.NewObjectID().Hex()
	notificacion.FechaCreacion = time.Now()
	notificacion.Estado = models.EstadoPendiente
	notificacion.IntentosEnvio = 0

	// Insertar en la base de datos
	_, err := s.db.Collection("notificaciones").InsertOne(ctx, notificacion)
	if err != nil {
		return fmt.Errorf("error creando notificación: %v", err)
	}

	// Enviar a la cola de procesamiento
	s.colaNotificaciones <- notificacion

	return nil
}

// procesarNotificaciones procesa las notificaciones pendientes
func (s *NotificacionService) procesarNotificaciones() {
	for notificacion := range s.colaNotificaciones {
		// Obtener preferencias del usuario
		preferencias, err := s.obtenerPreferenciasUsuario(notificacion.UsuarioID)
		if err != nil {
			s.actualizarEstadoNotificacion(notificacion.ID, models.EstadoError, err.Error())
			continue
		}

		// Verificar si el usuario quiere recibir este tipo de notificación
		if !s.usuarioRecibeTipoNotificacion(preferencias, notificacion.Tipo) {
			s.actualizarEstadoNotificacion(notificacion.ID, models.EstadoEnviada, "Usuario no recibe este tipo de notificación")
			continue
		}

		// Verificar horario de notificación
		if !s.enHorarioNotificacion(preferencias) {
			// Reprogramar notificación para el siguiente horario permitido
			s.reprogramarNotificacion(notificacion, preferencias)
			continue
		}

		// Enviar notificaciones según preferencias
		var errores []error

		if s.config.NotificarEmail && preferencias.Email != "" {
			if err := s.enviarEmail(notificacion, preferencias); err != nil {
				errores = append(errores, fmt.Errorf("error enviando email: %v", err))
			}
		}

		if s.config.NotificarSMS && preferencias.Telefono != "" {
			if err := s.enviarSMS(notificacion, preferencias); err != nil {
				errores = append(errores, fmt.Errorf("error enviando SMS: %v", err))
			}
		}

		if s.config.NotificarPush && len(preferencias.TokensPush) > 0 {
			if err := s.enviarPush(notificacion, preferencias); err != nil {
				errores = append(errores, fmt.Errorf("error enviando push: %v", err))
			}
		}

		// Actualizar estado de la notificación
		if len(errores) > 0 {
			s.actualizarEstadoNotificacion(notificacion.ID, models.EstadoError, fmt.Sprintf("%v", errores))
		} else {
			s.actualizarEstadoNotificacion(notificacion.ID, models.EstadoEnviada, "")
		}
	}
}

// validarNotificacion valida una notificación antes de crearla
func (s *NotificacionService) validarNotificacion(notificacion *models.Notificacion) error {
	if notificacion.Tipo == "" {
		return errors.New("tipo de notificación requerido")
	}
	if notificacion.Titulo == "" {
		return errors.New("título requerido")
	}
	if notificacion.Mensaje == "" {
		return errors.New("mensaje requerido")
	}
	if notificacion.UsuarioID == "" {
		return errors.New("usuario ID requerido")
	}
	return nil
}

// obtenerPreferenciasUsuario obtiene las preferencias de notificación de un usuario
func (s *NotificacionService) obtenerPreferenciasUsuario(usuarioID string) (*models.PreferenciasNotificacion, error) {
	var preferencias models.PreferenciasNotificacion
	err := s.db.Collection("preferencias_notificaciones").FindOne(context.Background(), bson.M{"_id": usuarioID}).Decode(&preferencias)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo preferencias: %v", err)
	}
	return &preferencias, nil
}

// usuarioRecibeTipoNotificacion verifica si un usuario recibe un tipo específico de notificación
func (s *NotificacionService) usuarioRecibeTipoNotificacion(preferencias *models.PreferenciasNotificacion, tipo models.TipoNotificacion) bool {
	for _, t := range preferencias.TiposNotificacion {
		if t == tipo {
			return true
		}
	}
	return false
}

// enHorarioNotificacion verifica si es un horario permitido para notificar
func (s *NotificacionService) enHorarioNotificacion(preferencias *models.PreferenciasNotificacion) bool {
	ahora := time.Now()
	return ahora.After(preferencias.HorarioNotificacion.Inicio) && ahora.Before(preferencias.HorarioNotificacion.Fin)
}

// reprogramarNotificacion reprograma una notificación para el siguiente horario permitido
func (s *NotificacionService) reprogramarNotificacion(notificacion *models.Notificacion, preferencias *models.PreferenciasNotificacion) {
	// TODO: Implementar lógica de reprogramación
}

// actualizarEstadoNotificacion actualiza el estado de una notificación
func (s *NotificacionService) actualizarEstadoNotificacion(id string, estado models.EstadoNotificacion, error string) {
	update := bson.M{
		"$set": bson.M{
			"estado": estado,
			"error":  error,
		},
	}
	if estado == models.EstadoEnviada {
		update["$set"].(bson.M)["fecha_envio"] = time.Now()
	}
	_, err := s.db.Collection("notificaciones").UpdateOne(
		context.Background(),
		bson.M{"_id": id},
		update,
	)
	if err != nil {
		fmt.Printf("Error actualizando estado de notificación: %v\n", err)
	}
}

// enviarEmail envía una notificación por email
func (s *NotificacionService) enviarEmail(notificacion *models.Notificacion, preferencias *models.PreferenciasNotificacion) error {
	plantilla := s.config.PlantillasEmail[string(notificacion.Tipo)]
	if plantilla == "" {
		return errors.New("plantilla de email no encontrada")
	}
	return s.emailService.EnviarEmail(preferencias.Email, notificacion.Titulo, plantilla, notificacion.Data)
}

// enviarSMS envía una notificación por SMS
func (s *NotificacionService) enviarSMS(notificacion *models.Notificacion, preferencias *models.PreferenciasNotificacion) error {
	plantilla := s.config.PlantillasSMS[string(notificacion.Tipo)]
	if plantilla == "" {
		return errors.New("plantilla de SMS no encontrada")
	}
	return s.smsService.EnviarSMS(preferencias.Telefono, plantilla, notificacion.Data)
}

// enviarPush envía una notificación push
func (s *NotificacionService) enviarPush(notificacion *models.Notificacion, preferencias *models.PreferenciasNotificacion) error {
	plantilla := s.config.PlantillasPush[string(notificacion.Tipo)]
	if plantilla == "" {
		return errors.New("plantilla de push no encontrada")
	}
	return s.pushService.EnviarPush(preferencias.TokensPush, notificacion.Titulo, plantilla, notificacion.Data)
}

// MarcarComoLeida marca una notificación como leída
func (s *NotificacionService) MarcarComoLeida(id string) error {
	_, err := s.db.Collection("notificaciones").UpdateOne(
		context.Background(),
		bson.M{"_id": id},
		bson.M{
			"$set": bson.M{
				"estado":      models.EstadoLeida,
				"fecha_leida": time.Now(),
			},
		},
	)
	return err
}

// ObtenerNotificacionesUsuario obtiene las notificaciones de un usuario
func (s *NotificacionService) ObtenerNotificacionesUsuario(usuarioID string, limit int) ([]*models.Notificacion, error) {
	options := options.Find().SetLimit(int64(limit)).SetSort(bson.M{"fecha_creacion": -1})
	cursor, err := s.db.Collection("notificaciones").Find(
		context.Background(),
		bson.M{"usuario_id": usuarioID},
		options,
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var notificaciones []*models.Notificacion
	if err := cursor.All(context.Background(), &notificaciones); err != nil {
		return nil, err
	}

	return notificaciones, nil
}

// ObtenerNotificacionesPendientes obtiene las notificaciones pendientes de un usuario
func (s *NotificacionService) ObtenerNotificacionesPendientes(usuarioID string) ([]*models.Notificacion, error) {
	cursor, err := s.db.Collection("notificaciones").Find(
		context.Background(),
		bson.M{
			"usuario_id": usuarioID,
			"estado":     models.EstadoPendiente,
		},
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var notificaciones []*models.Notificacion
	if err := cursor.All(context.Background(), &notificaciones); err != nil {
		return nil, err
	}

	return notificaciones, nil
}

// ActualizarPreferenciasUsuario actualiza las preferencias de notificación de un usuario
func (s *NotificacionService) ActualizarPreferenciasUsuario(preferencias *models.PreferenciasNotificacion) error {
	options := options.Update().SetUpsert(true)
	_, err := s.db.Collection("preferencias_notificaciones").UpdateOne(
		context.Background(),
		bson.M{"_id": preferencias.UsuarioID},
		bson.M{"$set": preferencias},
		options,
	)
	return err
}
