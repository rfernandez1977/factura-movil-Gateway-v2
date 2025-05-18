package services

import (
	"context"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"FMgo/models"
)

// SeguridadService maneja la seguridad del sistema
type SeguridadService struct {
	db *mongo.Database
}

// NewSeguridadService crea una nueva instancia del servicio de seguridad
func NewSeguridadService(db *mongo.Database) *SeguridadService {
	return &SeguridadService{
		db: db,
	}
}

// RegistrarAcceso registra un intento de acceso al sistema
func (s *SeguridadService) RegistrarAcceso(
	ctx context.Context,
	usuarioID, rut, accion, ip, userAgent string,
	exitoso bool,
	detalles string,
) error {
	registro := models.RegistroAuditoriaAcceso{
		ID:          GenerateID(),
		UsuarioID:   usuarioID,
		Rut:         rut,
		Accion:      accion,
		IP:          ip,
		UserAgent:   userAgent,
		Exitoso:     exitoso,
		Detalles:    detalles,
		FechaAcceso: time.Now(),
	}

	_, err := s.db.Collection("auditoria_accesos").InsertOne(ctx, registro)
	return err
}

// RegistrarOperacion registra una operación en el sistema
func (s *SeguridadService) RegistrarOperacion(
	ctx context.Context,
	usuarioID, rut, operacion, entidad, entidadID string,
	cambios, estadoAnterior, estadoNuevo map[string]interface{},
	ip, userAgent string,
) error {
	registro := models.RegistroAuditoriaOperacion{
		ID:             GenerateID(),
		UsuarioID:      usuarioID,
		Rut:            rut,
		Operacion:      operacion,
		Entidad:        entidad,
		EntidadID:      entidadID,
		Cambios:        cambios,
		EstadoAnterior: estadoAnterior,
		EstadoNuevo:    estadoNuevo,
		IP:             ip,
		UserAgent:      userAgent,
		FechaOperacion: time.Now(),
	}

	_, err := s.db.Collection("auditoria_operaciones").InsertOne(ctx, registro)
	return err
}

// ValidarFirmaDigital valida una firma digital
func (s *SeguridadService) ValidarFirmaDigital(
	ctx context.Context,
	usuarioID string,
	documento []byte,
	firma []byte,
) (bool, error) {
	// Obtener firma digital del usuario
	var firmaDigital models.FirmaDigital
	err := s.db.Collection("firmas_digitales").FindOne(ctx, bson.M{
		"usuario_id":     usuarioID,
		"estado":         "ACTIVA",
		"vigencia_desde": bson.M{"$lte": time.Now()},
		"vigencia_hasta": bson.M{"$gte": time.Now()},
	}).Decode(&firmaDigital)
	if err != nil {
		return false, err
	}

	// Decodificar clave pública
	block, _ := pem.Decode(firmaDigital.ClavePublica)
	if block == nil {
		return false, errors.New("clave pública inválida")
	}

	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return false, err
	}

	rsaPubKey, ok := pubKey.(*rsa.PublicKey)
	if !ok {
		return false, errors.New("clave pública no es RSA")
	}

	// Calcular hash del documento
	hash := sha256.Sum256(documento)

	// Verificar firma
	err = rsa.VerifyPKCS1v15(rsaPubKey, crypto.SHA256, hash[:], firma)
	if err != nil {
		return false, nil
	}

	return true, nil
}

// EncriptarDatos encripta datos sensibles
func (s *SeguridadService) EncriptarDatos(
	ctx context.Context,
	entidad, entidadID, campo string,
	valor []byte,
) (*models.DatosEncriptados, error) {
	// Generar clave AES
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return nil, err
	}

	// Crear bloque AES
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Generar IV
	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return nil, err
	}

	// Encriptar datos
	stream := cipher.NewCFBEncrypter(block, iv)
	ciphertext := make([]byte, len(valor))
	stream.XORKeyStream(ciphertext, valor)

	// Crear registro de datos encriptados
	datos := &models.DatosEncriptados{
		ID:                GenerateID(),
		Entidad:           entidad,
		EntidadID:         entidadID,
		Campo:             campo,
		ValorEncriptado:   ciphertext,
		IV:                iv,
		Algoritmo:         "AES-256-CFB",
		Version:           1,
		FechaCreacion:     time.Now(),
		FechaModificacion: time.Now(),
	}

	// Guardar en base de datos
	_, err = s.db.Collection("datos_encriptados").InsertOne(ctx, datos)
	if err != nil {
		return nil, err
	}

	return datos, nil
}

// DesencriptarDatos desencripta datos sensibles
func (s *SeguridadService) DesencriptarDatos(
	ctx context.Context,
	entidad, entidadID, campo string,
) ([]byte, error) {
	// Obtener datos encriptados
	var datos models.DatosEncriptados
	err := s.db.Collection("datos_encriptados").FindOne(ctx, bson.M{
		"entidad":    entidad,
		"entidad_id": entidadID,
		"campo":      campo,
	}).Decode(&datos)
	if err != nil {
		return nil, err
	}

	// Obtener clave de encriptación (en un sistema real, esto debería ser más seguro)
	key := make([]byte, 32)
	// TODO: Implementar obtención segura de la clave

	// Crear bloque AES
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Desencriptar datos
	stream := cipher.NewCFBDecrypter(block, datos.IV)
	plaintext := make([]byte, len(datos.ValorEncriptado))
	stream.XORKeyStream(plaintext, datos.ValorEncriptado)

	return plaintext, nil
}

// GenerarReporteSeguridad genera un reporte de seguridad
func (s *SeguridadService) GenerarReporteSeguridad(
	ctx context.Context,
	fechaInicio, fechaFin time.Time,
) (*models.ReporteSeguridad, error) {
	// Contar accesos fallidos
	accesosFallidos, err := s.db.Collection("auditoria_accesos").CountDocuments(ctx, bson.M{
		"fecha_acceso": bson.M{
			"$gte": fechaInicio,
			"$lte": fechaFin,
		},
		"exitoso": false,
	})
	if err != nil {
		return nil, err
	}

	// Contar usuarios bloqueados
	usuariosBloqueados, err := s.db.Collection("usuarios").CountDocuments(ctx, bson.M{
		"estado": "BLOQUEADO",
	})
	if err != nil {
		return nil, err
	}

	// Contar firmas revocadas
	firmasRevocadas, err := s.db.Collection("firmas_digitales").CountDocuments(ctx, bson.M{
		"estado": "REVOCADA",
	})
	if err != nil {
		return nil, err
	}

	// Obtener alertas de seguridad
	var alertas []models.AlertaSeguridad
	cursor, err := s.db.Collection("alertas_seguridad").Find(ctx, bson.M{
		"fecha_alerta": bson.M{
			"$gte": fechaInicio,
			"$lte": fechaFin,
		},
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &alertas); err != nil {
		return nil, err
	}

	// Crear reporte
	reporte := &models.ReporteSeguridad{
		ID:                 GenerateID(),
		FechaInicio:        fechaInicio,
		FechaFin:           fechaFin,
		AccesosFallidos:    int(accesosFallidos),
		UsuariosBloqueados: int(usuariosBloqueados),
		FirmasRevocadas:    int(firmasRevocadas),
		AlertasSeguridad:   alertas,
		FechaGeneracion:    time.Now(),
	}

	// Guardar reporte
	_, err = s.db.Collection("reportes_seguridad").InsertOne(ctx, reporte)
	if err != nil {
		return nil, err
	}

	return reporte, nil
}
