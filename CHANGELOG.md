# Changelog

Todos los cambios notables en este proyecto serán documentados en este archivo.

El formato está basado en [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
y este proyecto adhiere a [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.2.0] - 2024-01-20

### Added
- Implementación completa del sistema de firma digital XML-DSIG
- Nuevo servicio XMLSignatureService para manejo de firmas
- Sistema de caché para certificados digitales
- Validación de firmas y certificados
- Soporte para múltiples formatos de certificados (P12/PEM)
- Sistema de logging multinivel con rotación
- Documentación técnica completa
- Nuevas guías de troubleshooting y desarrollo

### Changed
- Refactorización completa del sistema de firma digital
- Mejora en el manejo de errores y validaciones
- Optimización del sistema de caché
- Actualización de dependencias core
- Mejora en la estructura del proyecto

### Fixed
- Corrección en la validación de certificados expirados
- Mejora en el manejo de memoria en procesamiento XML
- Corrección de race conditions en caché de certificados
- Optimización de queries y conexiones a base de datos

## [1.1.0] - 2024-01-10

### Added
- Implementación inicial del sistema de firma digital
- Validación básica de documentos XML
- Integración con servicios SII
- Sistema básico de logging

### Changed
- Actualización de estructura de proyecto
- Mejora en manejo de configuraciones
- Optimización de validaciones XSD

### Fixed
- Corrección de errores en validación de schemas
- Mejora en manejo de errores de conexión
- Corrección de memory leaks en procesamiento XML

## [1.0.0] - 2024-01-01

### Added
- Primera versión estable del sistema
- Funcionalidades básicas de validación
- Integración inicial con SII
- Documentación básica

## [Unreleased]

### Dependencias Actualizadas
- `github.com/stretchr/testify` actualizado a v1.10.0
- `github.com/go-redis/redis/v8` en v8.11.5
- `go.mongodb.org/mongo-driver` en v1.17.3

### Dependencias Principales
- gin-gonic/gin v1.9.1 - Framework web
- go-redis/redis/v8 v8.11.5 - Cliente Redis
- mongodb/mongo-driver v1.17.3 - Driver MongoDB
- stretchr/testify v1.10.0 - Framework de testing
- streadway/amqp v1.1.0 - Cliente RabbitMQ
- prometheus/client_golang v1.11.1 - Métricas
- uber/zap v1.24.0 - Logging

### Notas de Compatibilidad
- Todas las dependencias son compatibles con Go 1.23
- Se utiliza el toolchain go1.24.2

[1.2.0]: https://github.com/tu-usuario/FMgo/compare/v1.1.0...v1.2.0
[1.1.0]: https://github.com/tu-usuario/FMgo/compare/v1.0.0...v1.1.0
[1.0.0]: https://github.com/tu-usuario/FMgo/releases/tag/v1.0.0 