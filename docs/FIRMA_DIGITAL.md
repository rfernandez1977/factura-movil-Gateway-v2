# Servicio de Firma Digital

## Descripción General
El servicio de firma digital proporciona una implementación robusta para la firma de documentos XML, específicamente diseñada para la integración con el Servicio de Impuestos Internos (SII) de Chile. La arquitectura está diseñada para ser modular, extensible y mantener un alto rendimiento bajo carga concurrente.

## Componentes Principales

### 1. Servicio Base de Firma (`BaseFirmaService`)
Proporciona la funcionalidad core para firmar documentos XML:
- Inicialización y manejo de certificados digitales
- Firma de documentos XML
- Validación de firmas
- Gestión de certificados

### 2. Servicio SII (`SIIFirmaService`)
Extiende el servicio base para proporcionar funcionalidad específica del SII:
- Firma de documentos de Semilla
- Firma de documentos de Token
- Firma de DTE (Documentos Tributarios Electrónicos)
- Validaciones específicas del SII

### 3. Caché de Certificados (`CertCache`)
Optimiza el rendimiento mediante el almacenamiento en caché de certificados:
- Gestión de TTL (Time-To-Live)
- Límite de elementos en caché
- Thread-safe para operaciones concurrentes
- Limpieza automática de certificados expirados

## Pruebas Unitarias

### 1. Pruebas del Servicio Base (`base_firma_service_test.go`)
- Inicialización del servicio con diferentes configuraciones
- Firma de documentos XML
- Validación de firmas
- Manejo de certificados
- Pruebas de concurrencia
- Casos de error

### 2. Pruebas del Servicio SII (`sii_firma_service_test.go`)
- Firma de documentos de Semilla
- Firma de documentos de Token
- Firma de DTE
- Validaciones de RUT emisor
- Pruebas de concurrencia con múltiples tipos de documentos
- Casos de error específicos del SII

### 3. Pruebas del Caché (`cert_cache_test.go`)
- Operaciones básicas de caché
- Expiración de elementos
- Límite de elementos
- Pruebas de concurrencia
- Casos límite y edge cases

## Datos de Prueba

### 1. Certificados de Prueba (`test_data/test_cert.go`)
- Generación de certificados X.509
- Generación de llaves RSA
- Soporte para formatos PEM y PKCS12
- Configuración específica para pruebas del SII

### 2. Documentos XML de Prueba
- Semilla
- Token
- DTE
- Documentos inválidos para pruebas de error

## Consideraciones de Seguridad
1. Manejo seguro de certificados y llaves privadas
2. Validación de RUT emisor
3. Verificación de firmas
4. Protección contra ataques de concurrencia
5. Limpieza segura de datos sensibles

## Rendimiento y Escalabilidad
1. Caché de certificados para optimizar el rendimiento
2. Soporte para operaciones concurrentes
3. Gestión eficiente de recursos
4. Limpieza automática de recursos no utilizados

## Integración con el SII
1. Cumplimiento de estándares del SII
2. Validación de esquemas XML
3. Manejo de errores específicos
4. Formato de firma compatible

## Ejemplos de Uso

### 1. Firma de un DTE
```go
config := &models.ConfiguracionFirma{
    RutaCertificado: "ruta/al/certificado.pfx",
    Password:        "password",
    RutEmpresa:     "76555555-5",
}

service, err := NewSIIFirmaService(config)
if err != nil {
    log.Fatal(err)
}

resultado, err := service.FirmarDTE(xmlDTE)
if err != nil {
    log.Fatal(err)
}

// El DTE firmado está en resultado.XMLFirmado
```

### 2. Obtención de Token
```go
// Firmar Semilla
semillaFirmada, err := service.FirmarSemilla(xmlSemilla)
if err != nil {
    log.Fatal(err)
}

// Firmar Token
tokenFirmado, err := service.FirmarToken(xmlToken)
if err != nil {
    log.Fatal(err)
}
```

## Mantenimiento y Extensión
1. Estructura modular para facilitar extensiones
2. Pruebas exhaustivas para garantizar estabilidad
3. Documentación detallada
4. Manejo consistente de errores

## Próximos Pasos
1. Implementación de más validaciones específicas del SII
2. Mejoras en el rendimiento del caché
3. Soporte para más tipos de documentos
4. Integración con sistemas de monitoreo 