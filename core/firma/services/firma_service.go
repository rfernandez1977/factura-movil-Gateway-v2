package services

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"time"

	"github.com/fmgo/core/firma/models"
)

// FirmaService implementa el servicio de firma digital
type FirmaService struct {
	certRepo     CertificadoRepository
	cacheService CacheService
	logger       Logger
}

// NewFirmaService crea una nueva instancia del servicio de firma
func NewFirmaService(certRepo CertificadoRepository, cache CacheService, logger Logger) *FirmaService {
	return &FirmaService{
		certRepo:     certRepo,
		cacheService: cache,
		logger:       logger,
	}
}

// FirmarXML firma un documento XML usando un certificado específico
func (s *FirmaService) FirmarXML(ctx context.Context, xmlData []byte, certID string) ([]byte, error) {
	// Obtener certificado (primero del caché, luego del repositorio)
	cert, err := s.obtenerCertificado(ctx, certID)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo certificado: %w", err)
	}

	// Validar vigencia del certificado
	if !cert.ValidarVigencia() {
		return nil, fmt.Errorf("certificado no vigente: %s", certID)
	}

	// Calcular digest del documento
	digest := sha256.Sum256(xmlData)
	digestValue := base64.StdEncoding.EncodeToString(digest[:])

	// Obtener certificado X509
	x509Cert, err := cert.ObtenerCertificadoX509()
	if err != nil {
		return nil, fmt.Errorf("error parseando certificado: %w", err)
	}

	// Crear firma XML
	firma := models.FirmaXML{
		SignedInfo: models.SignedInfo{
			CanonicalizationMethod: models.Method{Algorithm: string(models.C14N)},
			SignatureMethod:        models.Method{Algorithm: string(models.RSA_SHA256)},
			Reference: models.Reference{
				URI: "",
				Transforms: []models.Transform{
					{Algorithm: "http://www.w3.org/2000/09/xmldsig#enveloped-signature"},
				},
				DigestMethod: models.Method{Algorithm: string(models.SHA256)},
				DigestValue:  digestValue,
			},
		},
		KeyInfo: models.KeyInfo{
			X509Data: models.X509Data{
				X509Certificate: base64.StdEncoding.EncodeToString(x509Cert.Raw),
			},
		},
		Timestamp: time.Now(),
	}

	// Firmar el digest
	signature, err := s.firmarDigest(digest[:], cert)
	if err != nil {
		return nil, fmt.Errorf("error firmando digest: %w", err)
	}

	firma.SignatureValue = base64.StdEncoding.EncodeToString(signature)

	// Serializar firma a XML
	firmaXML, err := xml.MarshalIndent(firma, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error serializando firma: %w", err)
	}

	// Insertar firma en el documento
	return s.insertarFirma(xmlData, firmaXML)
}

// ValidarFirma valida la firma de un documento XML
func (s *FirmaService) ValidarFirma(ctx context.Context, xmlData []byte) (*models.EstadoFirma, error) {
	// Extraer firma del documento
	firma, err := s.extraerFirma(xmlData)
	if err != nil {
		return nil, fmt.Errorf("error extrayendo firma: %w", err)
	}

	// Obtener certificado
	certData, err := base64.StdEncoding.DecodeString(firma.KeyInfo.X509Data.X509Certificate)
	if err != nil {
		return nil, fmt.Errorf("error decodificando certificado: %w", err)
	}

	cert, err := x509.ParseCertificate(certData)
	if err != nil {
		return nil, fmt.Errorf("error parseando certificado: %w", err)
	}

	// Validar firma
	digest, err := base64.StdEncoding.DecodeString(firma.SignedInfo.Reference.DigestValue)
	if err != nil {
		return nil, fmt.Errorf("error decodificando digest: %w", err)
	}

	signature, err := base64.StdEncoding.DecodeString(firma.SignatureValue)
	if err != nil {
		return nil, fmt.Errorf("error decodificando firma: %w", err)
	}

	err = rsa.VerifyPKCS1v15(cert.PublicKey.(*rsa.PublicKey), crypto.SHA256, digest, signature)

	estado := &models.EstadoFirma{
		FechaValidacion: time.Now(),
		CertificadoID:   cert.SerialNumber.String(),
	}

	if err != nil {
		estado.Valida = false
		estado.Error = "firma inválida"
	} else {
		estado.Valida = true
	}

	return estado, nil
}

// obtenerCertificado obtiene un certificado del caché o del repositorio
func (s *FirmaService) obtenerCertificado(ctx context.Context, certID string) (*models.Certificado, error) {
	// Intentar obtener del caché
	if cert, err := s.cacheService.ObtenerCertificado(ctx, certID); err == nil {
		return cert, nil
	}

	// Si no está en caché, obtener del repositorio
	cert, err := s.certRepo.ObtenerCertificado(ctx, certID)
	if err != nil {
		return nil, err
	}

	// Guardar en caché para futuras consultas
	if err := s.cacheService.GuardarCertificado(ctx, certID, cert); err != nil {
		s.logger.Warn("error guardando certificado en caché", "error", err)
	}

	return cert, nil
}

// firmarDigest firma un digest usando la llave privada del certificado
func (s *FirmaService) firmarDigest(digest []byte, cert *models.Certificado) ([]byte, error) {
	privateKey, err := x509.ParsePKCS1PrivateKey(cert.LlavePrivadaPEM)
	if err != nil {
		return nil, fmt.Errorf("error parseando llave privada: %w", err)
	}

	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, digest)
	if err != nil {
		return nil, fmt.Errorf("error firmando digest: %w", err)
	}

	return signature, nil
}

// insertarFirma inserta una firma XML en un documento
func (s *FirmaService) insertarFirma(xmlData, firma []byte) ([]byte, error) {
	// Implementar lógica de inserción de firma en el documento XML
	// Esta es una implementación simplificada
	result := append(xmlData[:len(xmlData)-2], firma...)
	result = append(result, xmlData[len(xmlData)-2:]...)
	return result, nil
}

// extraerFirma extrae la firma de un documento XML
func (s *FirmaService) extraerFirma(xmlData []byte) (*models.FirmaXML, error) {
	var firma models.FirmaXML
	if err := xml.Unmarshal(xmlData, &firma); err != nil {
		return nil, fmt.Errorf("error deserializando firma: %w", err)
	}
	return &firma, nil
}

// RenovarCertificado renueva un certificado próximo a vencer
func (s *FirmaService) RenovarCertificado(ctx context.Context, certID string) error {
	cert, err := s.obtenerCertificado(ctx, certID)
	if err != nil {
		return fmt.Errorf("error obteniendo certificado: %w", err)
	}

	// Verificar si necesita renovación
	if !cert.NecesitaRenovacion() {
		return nil
	}

	// Generar nuevo par de llaves
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("error generando llaves: %w", err)
	}

	// Crear solicitud de certificado
	template := x509.Certificate{
		SerialNumber: cert.SerialNumber,
		Subject:      cert.Subject,
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(1, 0, 0), // 1 año de validez
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	// Crear nuevo certificado
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return fmt.Errorf("error creando certificado: %w", err)
	}

	// Actualizar certificado en el repositorio
	nuevoCert := &models.Certificado{
		ID:              cert.ID,
		Nombre:          cert.Nombre,
		RUT:            cert.RUT,
		SerialNumber:    cert.SerialNumber,
		Subject:         cert.Subject,
		CertificadoPEM: certDER,
		LlavePrivadaPEM: x509.MarshalPKCS1PrivateKey(privateKey),
		FechaEmision:    time.Now(),
		FechaVencimiento: time.Now().AddDate(1, 0, 0),
	}

	if err := s.certRepo.GuardarCertificado(ctx, nuevoCert); err != nil {
		return fmt.Errorf("error guardando certificado renovado: %w", err)
	}

	// Limpiar caché
	if err := s.cacheService.EliminarCertificado(ctx, certID); err != nil {
		s.logger.Warn("error eliminando certificado del caché", "error", err)
	}

	s.logger.Info("certificado renovado exitosamente", 
		"id", certID,
		"fecha_vencimiento", nuevoCert.FechaVencimiento)

	return nil
}

// ValidarCadenaCertificados valida la cadena de certificados
func (s *FirmaService) ValidarCadenaCertificados(ctx context.Context, certID string) error {
	cert, err := s.obtenerCertificado(ctx, certID)
	if err != nil {
		return fmt.Errorf("error obteniendo certificado: %w", err)
	}

	// Crear pool de certificados raíz
	rootPool := x509.NewCertPool()
	
	// Obtener certificados raíz (implementación dependerá de la CA)
	raices, err := s.obtenerCertificadosRaiz()
	if err != nil {
		return fmt.Errorf("error obteniendo certificados raíz: %w", err)
	}

	for _, raiz := range raices {
		rootPool.AddCert(raiz)
	}

	// Parsear certificado a validar
	x509Cert, err := cert.ObtenerCertificadoX509()
	if err != nil {
		return fmt.Errorf("error parseando certificado: %w", err)
	}

	// Crear pool de certificados intermedios
	intermediatePool := x509.NewCertPool()
	
	// Obtener certificados intermedios (implementación dependerá de la CA)
	intermedios, err := s.obtenerCertificadosIntermedios()
	if err != nil {
		return fmt.Errorf("error obteniendo certificados intermedios: %w", err)
	}

	for _, intermedio := range intermedios {
		intermediatePool.AddCert(intermedio)
	}

	// Validar cadena de certificados
	opts := x509.VerifyOptions{
		Roots:         rootPool,
		Intermediates: intermediatePool,
		CurrentTime:   time.Now(),
		KeyUsages:     []x509.ExtKeyUsage{x509.ExtKeyUsageAny},
	}

	if _, err := x509Cert.Verify(opts); err != nil {
		return fmt.Errorf("error validando cadena de certificados: %w", err)
	}

	return nil
}

// obtenerCertificadosRaiz obtiene los certificados raíz de la CA
func (s *FirmaService) obtenerCertificadosRaiz() ([]*x509.Certificate, error) {
	// TODO: Implementar obtención de certificados raíz según la CA utilizada
	return []*x509.Certificate{}, nil
}

// obtenerCertificadosIntermedios obtiene los certificados intermedios de la CA
func (s *FirmaService) obtenerCertificadosIntermedios() ([]*x509.Certificate, error) {
	// TODO: Implementar obtención de certificados intermedios según la CA utilizada