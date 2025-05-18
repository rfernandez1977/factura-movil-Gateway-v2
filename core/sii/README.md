# Módulo de Integración SII

Este módulo proporciona la integración con el Servicio de Impuestos Internos (SII) de Chile para el envío y consulta de Documentos Tributarios Electrónicos (DTE).

## Estructura del Módulo

```
core/sii/
├── client/           # Cliente HTTP para comunicación con el SII
├── infrastructure/   # Implementaciones de infraestructura
│   └── cache/       # Sistema de caché para tokens
├── models/          # Modelos y tipos de datos
└── service/         # Servicios de alto nivel
```

## Características Principales

- Gestión automática de tokens de autenticación
- Sistema de caché con Redis para tokens
- Manejo robusto de errores y reintentos
- Soporte para ambientes de certificación y producción
- Validación de certificados digitales
- Pruebas unitarias completas

## Uso Básico

```go
// Crear cliente HTTP
client, err := client.NewHTTPClient(
    "path/to/cert.pfx",
    "password",
    models.Certificacion,
    client.RetryConfig{
        MaxRetries: 3,
        RetryDelay: time.Second,
    },
)
if err != nil {
    log.Fatal(err)
}

// Crear servicio SII
service := service.NewDefaultSIIService(client)

// Enviar DTE
resp, err := service.EnviarDTE(context.Background(), dteByte)
if err != nil {
    log.Fatal(err)
}

// Consultar estado
estado, err := service.ConsultarEstado(context.Background(), resp.TrackID)
if err != nil {
    log.Fatal(err)
}
```

## Configuración del Caché Redis

```go
// Crear cliente Redis
redisClient := redis.NewClient(&redis.Options{
    Addr:     "localhost:6379",
    Password: "",
    DB:       0,
})

// Crear caché de tokens
tokenCache := redis.NewRedisTokenCache(redisClient, "sii")

// Configurar en el servicio
service.SetTokenCache(tokenCache)
```

## Manejo de Errores

El módulo proporciona tipos de error específicos para cada situación:

```go
if err != nil {
    var siiErr *models.SIIError
    if errors.As(err, &siiErr) {
        switch siiErr.Code {
        case models.ErrAuthInvalid:
            // Manejar error de autenticación
        case models.ErrDTEInvalido:
            // Manejar error de DTE
        case models.ErrTimeout:
            // Manejar timeout
        }
    }
}
```

## Pruebas

Para ejecutar las pruebas:

```bash
go test ./core/sii/... -v
```

## Estado del Proyecto

- [x] Cliente HTTP y Certificados (Fase 2)
- [x] Sistema de Caché de Tokens
- [x] Manejo de Errores
- [x] Pruebas Unitarias (85% cobertura)
- [ ] Procesamiento XML (en progreso)
- [ ] Documentación API (70% completado)

## Contribución

1. Fork el repositorio
2. Crear una rama feature (`git checkout -b feature/nueva-funcionalidad`)
3. Commit los cambios (`git commit -am 'Agregar nueva funcionalidad'`)
4. Push a la rama (`git push origin feature/nueva-funcionalidad`)
5. Crear un Pull Request

## Licencia

Este proyecto está licenciado bajo la Licencia MIT - ver el archivo [LICENSE](LICENSE) para más detalles. 