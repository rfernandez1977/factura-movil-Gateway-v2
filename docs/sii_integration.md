# Integración SII - FMgo

## Descripción General
Este documento detalla la implementación de la integración con el Servicio de Impuestos Internos (SII) en el proyecto FMgo.

## Estructura del Código
La integración está organizada en los siguientes paquetes principales:

```
core/sii/
├── client/          # Cliente HTTP para comunicación con el SII
├── models/          # Modelos y tipos de datos
├── infrastructure/  # Infraestructura (certificados, etc.)
└── retry/          # Lógica de reintentos
```

## Componentes Principales

### Cliente HTTP (core/sii/client/http_client.go)
- Implementa la comunicación con los servicios web del SII
- Manejo de certificados digitales
- Gestión de semillas y tokens
- Envío y consulta de DTEs
- Validación de respuestas
- Sistema de reintentos configurable

### Modelos (core/sii/models/)
- Definición de tipos y estructuras
- Configuración del cliente
- Estados y respuestas del SII
- Validaciones de datos

### Infraestructura
- Gestión de certificados digitales
- Validación de certificados
- Información de certificados

## Funcionalidades Implementadas

### 1. Autenticación
- Obtención de semilla
- Generación de token
- Validación de certificados

### 2. Operaciones DTE
- Envío de documentos
- Consulta de estado
- Verificación de comunicación

### 3. Manejo de Errores
- Errores tipados
- Reintentos configurables
- Validación de respuestas

## Configuración

### Ejemplo de Configuración
```json
{
  "ambiente": "certificacion",
  "cert_path": "/ruta/al/certificado.crt",
  "key_path": "/ruta/a/llave.key",
  "schema_path": "/ruta/al/schema.xsd",
  "timeout": 30,
  "retry_count": 3,
  "retry_delay": 5
}
```

### Ambientes Disponibles
- Certificación: https://maullin.sii.cl
- Producción: https://palena.sii.cl

## Pruebas
Se incluye un script de prueba de conexión en `scripts/test_sii_connection.go` que verifica:
- Obtención de semilla
- Generación de token
- Verificación de comunicación

## Mejoras Realizadas

### 1. Consolidación de Código
- Eliminación de implementaciones duplicadas
- Unificación de cliente SII
- Consolidación de modelos

### 2. Sistema de Logging
- Logger unificado
- Niveles de log configurables
- Mejor trazabilidad

### 3. Manejo de Configuración
- Validación robusta
- Manejo de ambientes
- Verificación de archivos

### 4. Seguridad
- Validación de certificados
- Manejo seguro de credenciales
- Verificación de expiración

## Mantenimiento

### Certificados
- Verificar regularmente la validez
- Monitorear fechas de expiración
- Mantener respaldos seguros

### Monitoreo
- Logging de operaciones
- Registro de errores
- Métricas de uso

## Próximos Pasos
1. Implementar pruebas automatizadas adicionales
2. Mejorar el sistema de métricas
3. Agregar documentación de API
4. Implementar validaciones adicionales de documentos 