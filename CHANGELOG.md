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

### Integración SII - Mejoras y Consolidación
- Consolidación de implementaciones del cliente SII
  - Eliminación de archivos duplicados (cert_manager.go, documento.go)
  - Unificación del cliente SII en core/sii/client
  - Consolidación de modelos en core/sii/models

- Mejoras en el Sistema de Logging
  - Implementación de logger unificado en utils/logger
  - Mejora en el manejo de niveles de log
  - Mejor trazabilidad de operaciones

- Mejoras en la Configuración
  - Implementación de validaciones robustas
  - Mejor manejo de ambientes (certificación/producción)
  - Verificación de archivos y rutas

- Mejoras en el Cliente HTTP
  - Eliminación de campos redundantes
  - Mejor manejo de errores y validaciones
  - Implementación de sistema de reintentos
  - Mejora en el manejo de certificados

- Mejoras en los Modelos
  - Implementación de validaciones de tipos
  - Mejora en la documentación
  - Implementación de métodos de validación
  - Optimización de tags XML y JSON

- Documentación
  - Creación de documentación detallada de la integración
  - Actualización de ejemplos y guías
  - Documentación de configuración y ambientes

## [0.1.0] - 2024-03-XX

### Agregado
- Implementación básica del validador CAF (MVP)
  - Validación de RUT emisor
  - Validación de tipo DTE
  - Control de rango de folios
  - Validación de fechas de vigencia
  - Control de folios usados en memoria
- Servicio de gestión de CAFs
  - Registro de CAFs
  - Validación de folios
  - Consulta de estado
- Pruebas unitarias completas
  - Cobertura > 80%
  - Casos de prueba para validaciones básicas
  - Pruebas de concurrencia básicas

### Pendiente para próximas versiones
- Verificación de firmas XML
- Persistencia de folios usados
- Métricas y monitoreo
- Pruebas de carga
- Manejo avanzado de concurrencia

[1.2.0]: https://github.com/tu-usuario/FMgo/compare/v1.1.0...v1.2.0
[1.1.0]: https://github.com/tu-usuario/FMgo/compare/v1.0.0...v1.1.0
[1.0.0]: https://github.com/tu-usuario/FMgo/releases/tag/v1.0.0 