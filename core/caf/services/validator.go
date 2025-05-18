package services

import (
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/sha1"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"math/big"
	"time"

	"FMgo/core/caf/models"
)

// ValidadorCAF define la interfaz para validación de CAF
type ValidadorCAF interface {
	ValidarCAF(ctx context.Context, xmlCAF []byte) (*models.ResultadoValidacion, error)
	ValidarFolio(ctx context.Context, folio int, tipoDTE int) (bool, error)
	ObtenerRangoFolios(ctx context.Context, tipoDTE int) (int, int, error)
}

// ServicioValidadorCAF implementa la interfaz ValidadorCAF
type ServicioValidadorCAF struct {
	cacheService CacheService
	logger       Logger
}

// NewValidadorCAF crea una nueva instancia del validador CAF
func NewValidadorCAF(cache CacheService, logger Logger) ValidadorCAF {
	return &ServicioValidadorCAF{
		cacheService: cache,
		logger:       logger,
	}
}

// ValidarCAF valida un archivo CAF
func (s *ServicioValidadorCAF) ValidarCAF(ctx context.Context, xmlCAF []byte) (*models.ResultadoValidacion, error) {
	resultado := &models.ResultadoValidacion{
		Timestamp: time.Now(),
	}

	var autorizacion models.AutorizacionCAF
	if err := xml.Unmarshal(xmlCAF, &autorizacion); err != nil {
		resultado.Error = fmt.Errorf("error al decodificar XML: %v", err)
		return resultado, nil
	}

	// Validar firma del CAF
	if err := s.validarFirmaCAF(&autorizacion); err != nil {
		resultado.Error = fmt.Errorf("error en firma del CAF: %v", err)
		return resultado, nil
	}

	// Validar fechas
	fechaEmision, err := time.Parse("2006-01-02", autorizacion.CAF.DA.FA)
	if err != nil {
		resultado.Error = fmt.Errorf("fecha de emisión inválida: %v", err)
		return resultado, nil
	}

	if fechaEmision.After(time.Now()) {
		resultado.Error = fmt.Errorf("fecha de emisión posterior a la actual")
		return resultado, nil
	}

	// Guardar en caché para futuras validaciones
	caf := &models.CAF{
		RUT:         autorizacion.CAF.DA.RE,
		RazonSocial: autorizacion.CAF.DA.RS,
		TipoDTE:     autorizacion.CAF.DA.TD,
		FolioDesde:  autorizacion.CAF.DA.RNG.D,
		FolioHasta:  autorizacion.CAF.DA.RNG.H,
		XMLOriginal: xmlCAF,
	}

	cacheKey := fmt.Sprintf("CAF_%s_%d", caf.RUT, caf.TipoDTE)
	if err := s.cacheService.Set(ctx, cacheKey, caf, 24*time.Hour); err != nil {
		s.logger.Warn("error guardando CAF en caché", "error", err)
	}

	resultado.Valido = true
	resultado.Detalles = fmt.Sprintf("CAF válido para %s, tipo DTE %d",
		autorizacion.CAF.DA.RS,
		autorizacion.CAF.DA.TD)

	return resultado, nil
}

// ValidarFolio verifica si un folio está dentro del rango autorizado
func (s *ServicioValidadorCAF) ValidarFolio(ctx context.Context, folio int, tipoDTE int) (bool, error) {
	desde, hasta, err := s.ObtenerRangoFolios(ctx, tipoDTE)
	if err != nil {
		return false, err
	}

	return folio >= desde && folio <= hasta, nil
}

// ObtenerRangoFolios obtiene el rango de folios disponible para un tipo de DTE
func (s *ServicioValidadorCAF) ObtenerRangoFolios(ctx context.Context, tipoDTE int) (int, int, error) {
	// Buscar en caché primero
	cacheKey := fmt.Sprintf("CAF_*_%d", tipoDTE)
	cafInterface, err := s.cacheService.Get(ctx, cacheKey)
	if err != nil {
		return 0, 0, fmt.Errorf("no se encontró CAF para tipo DTE %d", tipoDTE)
	}

	caf, ok := cafInterface.(*models.CAF)
	if !ok {
		return 0, 0, fmt.Errorf("error al convertir CAF desde caché")
	}

	return caf.FolioDesde, caf.FolioHasta, nil
}

// validarFirmaCAF valida la firma del CAF
func (s *ServicioValidadorCAF) validarFirmaCAF(autorizacion *models.AutorizacionCAF) error {
	// Decodificar llave pública
	modulus, err := base64.StdEncoding.DecodeString(autorizacion.CAF.DA.RSAPK.M)
	if err != nil {
		return fmt.Errorf("error decodificando módulo RSA: %v", err)
	}

	exponent, err := base64.StdEncoding.DecodeString(autorizacion.CAF.DA.RSAPK.E)
	if err != nil {
		return fmt.Errorf("error decodificando exponente RSA: %v", err)
	}

	publicKey := &rsa.PublicKey{
		N: new(big.Int).SetBytes(modulus),
		E: int(new(big.Int).SetBytes(exponent).Int64()),
	}

	// Obtener datos a verificar
	dataToVerify := fmt.Sprintf("%s%s%d%d%d%s%s%s",
		autorizacion.CAF.DA.RE,
		autorizacion.CAF.DA.RS,
		autorizacion.CAF.DA.TD,
		autorizacion.CAF.DA.RNG.D,
		autorizacion.CAF.DA.RNG.H,
		autorizacion.CAF.DA.FA,
		autorizacion.CAF.DA.RSAPK.M,
		autorizacion.CAF.DA.RSAPK.E)

	// Calcular hash SHA1
	hashed := sha1.Sum([]byte(dataToVerify))

	// Decodificar firma
	signature, err := base64.StdEncoding.DecodeString(autorizacion.CAF.FRMA.Valor)
	if err != nil {
		return fmt.Errorf("error decodificando firma: %v", err)
	}

	// Verificar firma
	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA1, hashed[:], signature)
	if err != nil {
		return fmt.Errorf("firma inválida: %v", err)
	}

	return nil
}
