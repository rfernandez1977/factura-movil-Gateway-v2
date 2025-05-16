# Plan de Modularización - Firma Digital y CAF

## 1. Estructura del Módulo

### Nueva Estructura de Directorios
```
firma/
├── core/
│   ├── models/
│   │   ├── certificado.go
│   │   ├── caf.go
│   │   ├── firma.go
│   │   └── tipos.go
│   ├── services/
│   │   ├── firma_service.go
│   │   ├── caf_service.go
│   │   └── validacion.go
│   └── interfaces/
│       ├── repository.go
│       └── service.go
├── infrastructure/
│   ├── storage/
│   │   ├── filesystem/
│   │   └── database/
│   ├── crypto/
│   │   ├── rsa/
│   │   └── x509/
│   └── cache/
│       └── redis/
└── api/
    ├── handlers/
    ├── routes/
    └── middleware/
```

### Plan de Migración

1. **Fase 1: Modelos Core (2-3 días)**
   - Definir estructura de certificados
   - Implementar modelo CAF
   - Crear tipos de firma
   - Establecer interfaces base

2. **Fase 2: Servicios de Firma (3-4 días)**
   - Implementar firma XML
   - Validación de certificados
   - Gestión de CAF
   - Manejo de errores

3. **Fase 3: Almacenamiento (2-3 días)**
   - Sistema de archivos seguro
   - Persistencia en base de datos
   - Caché de certificados
   - Backup automático

4. **Fase 4: API y Utilidades (2-3 días)**
   - Endpoints REST
   - CLI para gestión
   - Herramientas de diagnóstico
   - Documentación

### Componentes Principales

1. **Servicio de Firma**
```go
type FirmaService interface {
    // Operaciones de firma
    FirmarXML(ctx context.Context, documento []byte, cert *Certificado) ([]byte, error)
    ValidarFirma(ctx context.Context, documentoFirmado []byte) error
    ObtenerDatosFirma(ctx context.Context, documentoFirmado []byte) (*DatosFirma, error)

    // Gestión de certificados
    CargarCertificado(ctx context.Context, path string, password string) (*Certificado, error)
    ValidarCertificado(ctx context.Context, cert *Certificado) error
    RenovarCertificado(ctx context.Context, cert *Certificado) (*Certificado, error)
}
```

2. **Servicio CAF**
```go
type CAFService interface {
    // Gestión de CAF
    ObtenerCAF(ctx context.Context, tipo string, folio int64) (*CAF, error)
    RegistrarCAF(ctx context.Context, caf *CAF) error
    ValidarCAF(ctx context.Context, caf *CAF) error
    ConsultarDisponibilidad(ctx context.Context, tipo string) (*DisponibilidadCAF, error)

    // Monitoreo
    ObtenerEstadisticas(ctx context.Context) (*EstadisticasCAF, error)
    NotificarBajoStock(ctx context.Context, tipo string) error
}
```

3. **Repositorio**
```go
type FirmaRepository interface {
    // Certificados
    GuardarCertificado(ctx context.Context, cert *Certificado) error
    ObtenerCertificado(ctx context.Context, id string) (*Certificado, error)
    ListarCertificados(ctx context.Context) ([]*Certificado, error)
    EliminarCertificado(ctx context.Context, id string) error

    // CAF
    GuardarCAF(ctx context.Context, caf *CAF) error
    ObtenerCAF(ctx context.Context, tipo string, folio int64) (*CAF, error)
    ListarCAFs(ctx context.Context, tipo string) ([]*CAF, error)
    ActualizarEstadoCAF(ctx context.Context, id string, estado EstadoCAF) error
}
```

### Aspectos de Seguridad

1. **Almacenamiento Seguro**
   - Encriptación en reposo
   - Control de acceso
   - Auditoría de uso
   - Respaldos seguros

2. **Manejo de Claves**
   - Rotación periódica
   - Almacenamiento seguro
   - Control de acceso
   - Logs de uso

3. **Validaciones**
   - Integridad de certificados
   - Vigencia de CAF
   - Firmas válidas
   - Permisos adecuados

### Monitoreo y Alertas

1. **Métricas**
   - Uso de certificados
   - Consumo de CAF
   - Errores de firma
   - Performance

2. **Alertas**
   - CAF bajo stock
   - Certificados por vencer
   - Errores críticos
   - Intentos no autorizados

### Plan de Pruebas

1. **Unitarias**
   - Firma de documentos
   - Validación de certificados
   - Gestión de CAF
   - Manejo de errores

2. **Integración**
   - Flujo completo de firma
   - Almacenamiento seguro
   - Caché y performance
   - Recuperación de errores

3. **Seguridad**
   - Penetration testing
   - Validación de encriptación
   - Auditoría de accesos
   - Pruebas de recuperación

### Siguientes Pasos

1. **Inmediatos (1-2 días)**
   - Crear estructura base
   - Migrar código existente
   - Configurar ambiente

2. **Corto Plazo (1 semana)**
   - Implementar servicios core
   - Configurar almacenamiento
   - Establecer pruebas base

3. **Mediano Plazo (2-3 semanas)**
   - Completar funcionalidades
   - Implementar monitoreo
   - Documentar API
   - Realizar pruebas de seguridad

### Consideraciones Adicionales

1. **Performance**
   - Optimizar operaciones criptográficas
   - Implementar caché efectivo
   - Minimizar operaciones I/O
   - Manejar concurrencia

2. **Mantenibilidad**
   - Documentación detallada
   - Logs comprensivos
   - Herramientas de diagnóstico
   - Scripts de mantenimiento

3. **Recuperación**
   - Backup automático
   - Procedimientos de recuperación
   - Planes de contingencia
   - Documentación de procesos 