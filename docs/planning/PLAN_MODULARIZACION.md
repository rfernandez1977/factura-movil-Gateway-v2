# Plan de Modularización - FMgo

## 1. Módulo DTE (Documentos Tributarios Electrónicos)

### Nueva Estructura de Directorios
```
dte/
├── core/
│   ├── models/
│   │   ├── documento.go
│   │   ├── emisor.go
│   │   ├── receptor.go
│   │   └── tipos.go
│   ├── services/
│   │   ├── generacion.go
│   │   ├── validacion.go
│   │   └── firma.go
│   └── interfaces/
│       ├── repository.go
│       └── service.go
├── infrastructure/
│   ├── persistence/
│   │   ├── supabase/
│   │   └── mongodb/
│   └── sii/
│       ├── client.go
│       ├── xml/
│       └── validation/
└── api/
    ├── handlers/
    ├── routes/
    └── middleware/
```

### Plan de Migración

1. **Fase 1: Modelos Core (1 semana)**
   - Crear estructura base del módulo DTE
   - Migrar modelos desde `models/` a `dte/core/models/`
   - Implementar interfaces base
   - Actualizar importaciones

2. **Fase 2: Servicios Core (1 semana)**
   - Migrar servicios de DTE
   - Implementar nuevas interfaces de servicio
   - Consolidar lógica de validación
   - Actualizar dependencias

3. **Fase 3: Infraestructura (1 semana)**
   - Implementar capa de persistencia
   - Migrar cliente SII
   - Consolidar manejo de XML
   - Implementar nuevos repositorios

4. **Fase 4: API y Handlers (1 semana)**
   - Migrar handlers existentes
   - Implementar nuevas rutas
   - Actualizar middleware
   - Documentar API

### Dependencias a Consolidar

1. **Internas**
   - Servicios de validación
   - Generación de XML
   - Firma digital
   - Manejo de CAF

2. **Externas**
   - Cliente HTTP para SII
   - Almacenamiento (Supabase/MongoDB)
   - Caché (Redis)

### Interfaces Principales

1. **DocumentoService**
```go
type DocumentoService interface {
    Generar(ctx context.Context, doc *models.DocumentoRequest) (*models.Documento, error)
    Validar(ctx context.Context, doc *models.Documento) error
    Firmar(ctx context.Context, doc *models.Documento) error
    Enviar(ctx context.Context, doc *models.Documento) (*models.RespuestaSII, error)
    ConsultarEstado(ctx context.Context, trackID string) (*models.EstadoSII, error)
}
```

2. **DocumentoRepository**
```go
type DocumentoRepository interface {
    Guardar(ctx context.Context, doc *models.Documento) error
    BuscarPorID(ctx context.Context, id string) (*models.Documento, error)
    BuscarPorFolio(ctx context.Context, tipo, folio string) (*models.Documento, error)
    Actualizar(ctx context.Context, doc *models.Documento) error
    Eliminar(ctx context.Context, id string) error
}
```

### Pruebas

1. **Unitarias**
   - Modelos
   - Servicios
   - Validaciones
   - Generación XML

2. **Integración**
   - Persistencia
   - Cliente SII
   - Firma digital
   - Flujo completo

3. **End-to-End**
   - Flujo de emisión
   - Consulta de estado
   - Validación de documentos

### Métricas de Éxito

1. **Calidad de Código**
   - Cobertura de pruebas > 80%
   - Complejidad ciclomática < 10
   - Duplicación de código < 5%

2. **Rendimiento**
   - Tiempo de respuesta < 500ms
   - Uso de memoria estable
   - Conexiones concurrentes > 100

### Siguientes Pasos

1. **Inmediatos**
   - Crear estructura de directorios
   - Migrar primer modelo (Documento)
   - Implementar interfaces base

2. **Corto Plazo**
   - Completar migración de modelos
   - Implementar servicios core
   - Configurar pruebas unitarias

3. **Mediano Plazo**
   - Implementar persistencia
   - Migrar cliente SII
   - Documentar API 