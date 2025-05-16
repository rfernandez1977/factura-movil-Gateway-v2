# Arquitectura del Sistema FMgo

## Visión General

FMgo es un sistema modular diseñado para la integración con el Servicio de Impuestos Internos (SII) de Chile, enfocado en la gestión de documentos tributarios electrónicos (DTE). La arquitectura está diseñada para ser escalable, mantenible y altamente disponible.

## Estructura del Proyecto

```
core/
├── firma/                    # Módulo de firma digital
│   ├── common/              # Componentes comunes
│   │   └── interfaces.go    # Interfaces compartidas (Logger, etc.)
│   ├── models/              # Modelos de dominio
│   │   ├── certificado.go   # Modelo de certificado digital
│   │   ├── firma.go        # Modelo de firma XML
│   │   ├── caf.go          # Modelo de CAF
│   │   └── tipos.go        # Tipos y constantes comunes
│   ├── services/           # Servicios de negocio
│   │   ├── firma_service.go # Servicio de firma digital
│   │   ├── caf_service.go  # Servicio de gestión de CAF
│   │   └── interfaces.go   # Interfaces de servicios
│   └── infrastructure/     # Implementaciones de infraestructura
│       ├── storage/        # Almacenamiento
│       │   └── filesystem/ # Implementación en sistema de archivos
│       └── cache/         # Caché
│           └── redis/     # Implementación en Redis
```

## Componentes Principales

### 1. Módulo de Firma Digital (`core/firma`)

#### 1.1 Modelos (`models/`)
- **Certificado**: Gestión de certificados digitales
  - Validación de vigencia
  - Gestión de estados (Activo, Revocado, Expirado)
  - Conversión PEM/X509
  
- **FirmaXML**: Estructura de firma XML-DSig
  - Soporte para SignedInfo
  - Gestión de referencias y transformaciones
  - Información de certificados X509

- **CAF**: Código de Autorización de Folios
  - Control de rangos de folios
  - Validación de vigencia
  - Gestión de estados

#### 1.2 Servicios (`services/`)
- **FirmaService**: 
  - Firma de documentos XML
  - Validación de firmas
  - Gestión de certificados

- **CAFService**:
  - Gestión de CAF
  - Control de stock
  - Sistema de alertas

#### 1.3 Infraestructura (`infrastructure/`)

##### Storage (`storage/filesystem/`)
- **CertificadoStorage**:
  - Almacenamiento seguro de certificados
  - Gestión de metadatos
  - Control de acceso

- **CAFStorage**:
  - Almacenamiento de CAF
  - Gestión de XML y firmas SII
  - Organización por tipo de documento

##### Cache (`cache/redis/`)
- **CertificadoCache**:
  - Caché de certificados
  - TTL configurable
  - Gestión de invalidación

- **CAFCache**:
  - Caché de CAF
  - Indexación por tipo
  - Gestión de conjuntos

### 2. Interfaces Comunes (`common/`)

#### 2.1 Logger
- Niveles de logging (Debug, Info, Warn, Error)
- Soporte para contextos y metadatos
- Implementación flexible

## Patrones de Diseño

### 1. Repository Pattern
- Separación de la lógica de acceso a datos
- Interfaces bien definidas
- Implementaciones intercambiables

### 2. Service Layer
- Encapsulación de lógica de negocio
- Gestión de transacciones
- Coordinación entre componentes

### 3. Cache-Aside
- Caché transparente
- Política de expiración
- Consistencia eventual

## Seguridad

### 1. Certificados Digitales
- Almacenamiento seguro de llaves privadas
- Validación de certificados
- Control de acceso granular

### 2. CAF
- Validación de firmas SII
- Control de folios
- Alertas de stock

## Monitoreo y Logging

### 1. Sistema de Logging
- Logging multinivel
- Contexto enriquecido
- Trazabilidad de operaciones

### 2. Alertas
- Alertas de stock de CAF
- Notificación de certificados por expirar
- Monitoreo de estado del sistema

## Consideraciones de Rendimiento

### 1. Caché
- Caché en memoria con Redis
- TTL configurable
- Invalidación selectiva

### 2. Concurrencia
- Locks de lectura/escritura
- Operaciones atómicas
- Control de concurrencia optimista

## Extensibilidad

El sistema está diseñado para ser extensible en varios aspectos:

1. **Almacenamiento**: Nuevas implementaciones de storage
2. **Caché**: Soporte para otros sistemas de caché
3. **Logging**: Integración con diferentes sistemas de logging
4. **Alertas**: Nuevos canales y tipos de alertas 