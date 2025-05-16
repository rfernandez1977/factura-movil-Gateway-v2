# Sistema de Firma Digital

## Descripción General
El sistema de firma digital implementa la funcionalidad necesaria para firmar y validar documentos XML según los requerimientos del SII (Servicio de Impuestos Internos de Chile).

## Componentes Principales

### FirmaService
Servicio principal para la gestión de firmas digitales.

```go
firmaService, err := services.NewFirmaService(
    "ruta/certificado.p12",
    "ruta/clave.key",
    "contraseña",
    "76.555.555-5"
)
```

#### Funcionalidades Principales:
- Firma de documentos XML
- Validación de firmas
- Gestión de certificados digitales
- Caché de certificados

### XMLProcessor
Procesador especializado para documentos XML del SII.

```go
xmlProcessor := services.NewXMLProcessor(logger)
```

#### Características:
- Validación de estructura XML
- Extracción de certificados
- Extracción de firmas
- Limpieza de documentos XML

## Sistema de Logging

### Niveles de Log
- **DEBUG**: Información detallada para desarrollo y debugging
- **INFO**: Información general de operaciones exitosas
- **WARN**: Advertencias y situaciones inesperadas
- **ERROR**: Errores que requieren atención

### Ejemplo de Uso
```go
logger, err := logger.NewLogger("logs/firma_service.log", logger.DEBUG)
if err != nil {
    // Manejar error
}
defer logger.Close()

// Logging de operaciones
logger.Info("Iniciando proceso de firma")
logger.Debug("Detalles del documento: %s", detalles)
logger.Error("Error en proceso: %v", err)
```

## Proceso de Firma

### 1. Preparación
```go
xmlData := []byte(`...documento XML...`)
```

### 2. Validación Previa
```go
if err := xmlProcessor.validarEstructuraXML(xmlData); err != nil {
    return err
}
```

### 3. Firma del Documento
```go
signedXML, err := firmaService.FirmarXML(xmlData)
if err != nil {
    return err
}
```

### 4. Validación de Firma
```go
if err := firmaService.ValidarFirma(signedXML); err != nil {
    return err
}
```

## Manejo de Certificados

### Formato Soportado
- PKCS#12 (.p12)
- Certificados X.509

### Caché de Certificados
```go
certCache := services.NewCertCache(24*time.Hour, 100) // 24 horas, máximo 100 items
```

## Buenas Prácticas

### Seguridad
1. Nunca almacenar contraseñas en texto plano
2. Usar variables de entorno para configuraciones sensibles
3. Mantener los certificados en ubicaciones seguras
4. Rotar los logs periódicamente

### Rendimiento
1. Utilizar el sistema de caché para certificados frecuentes
2. Limpiar los XML antes de procesarlos
3. Configurar niveles de log apropiados en producción

### Manejo de Errores
1. Validar la estructura XML antes de firmar
2. Verificar la validez de los certificados
3. Implementar reintentos para operaciones críticas
4. Mantener logs detallados para debugging

## Ejemplos de Uso

### Firma Básica
```go
func firmarDocumento(xmlData []byte) ([]byte, error) {
    firmaService, err := services.NewFirmaService(
        os.Getenv("CERT_PATH"),
        os.Getenv("KEY_PATH"),
        os.Getenv("CERT_PASSWORD"),
        os.Getenv("RUT_EMPRESA")
    )
    if err != nil {
        return nil, err
    }

    return firmaService.FirmarXML(xmlData)
}
```

### Validación Completa
```go
func validarDocumentoFirmado(xmlData []byte) error {
    firmaService, err := services.NewFirmaService(
        os.Getenv("CERT_PATH"),
        os.Getenv("KEY_PATH"),
        os.Getenv("CERT_PASSWORD"),
        os.Getenv("RUT_EMPRESA")
    )
    if err != nil {
        return err
    }

    // Validar estructura
    if err := firmaService.xmlProc.validarEstructuraXML(xmlData); err != nil {
        return err
    }

    // Validar firma
    return firmaService.ValidarFirma(xmlData)
}
```

## Troubleshooting

### Problemas Comunes

1. **Error de Certificado**
   - Verificar la ruta del certificado
   - Confirmar la contraseña
   - Validar la vigencia del certificado

2. **Error de Firma**
   - Verificar la estructura del XML
   - Confirmar que el certificado tiene permisos de firma
   - Revisar los logs en nivel DEBUG

3. **Problemas de Rendimiento**
   - Ajustar el tamaño del caché
   - Verificar el nivel de logging
   - Monitorear el uso de memoria

### Logs de Debugging
```go
logger.LogXMLOperation("FirmarDocumento", xmlData, err)
logger.LogCertOperation("ValidarCertificado", certInfo, err)
``` 