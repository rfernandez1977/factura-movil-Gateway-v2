# Sistema de Logging

## Descripción General
El sistema de logging proporciona una infraestructura robusta para el registro y monitoreo de operaciones en el sistema FMgo, con énfasis especial en las operaciones relacionadas con documentos tributarios electrónicos.

## Características Principales

### Niveles de Log
El sistema soporta cuatro niveles de logging, ordenados por prioridad:

1. **DEBUG** (Nivel 0)
   - Información detallada para desarrollo y debugging
   - Incluye datos de operaciones XML completas
   - Útil durante el desarrollo y pruebas

2. **INFO** (Nivel 1)
   - Información general de operaciones exitosas
   - Estado de procesos importantes
   - Operaciones de negocio relevantes

3. **WARN** (Nivel 2)
   - Advertencias sobre situaciones inesperadas
   - Problemas no críticos
   - Situaciones que requieren atención

4. **ERROR** (Nivel 3)
   - Errores críticos que afectan la operación
   - Fallos en operaciones importantes
   - Situaciones que requieren intervención inmediata

### Formato de Log
Cada entrada de log incluye:
- Timestamp preciso (hasta milisegundos)
- Nivel de log
- Información del llamador (archivo:línea)
- Mensaje detallado

Ejemplo:
```
[2024-03-20 15:04:05.123] [firma_service.go:45] INFO: Certificado cargado exitosamente
```

## Uso del Sistema

### Inicialización
```go
logger, err := logger.NewLogger("logs/app.log", logger.DEBUG)
if err != nil {
    panic(err)
}
defer logger.Close()
```

### Logging Básico
```go
// Mensaje de debug
logger.Debug("Procesando documento %s", docID)

// Información importante
logger.Info("Documento %s firmado exitosamente", docID)

// Advertencia
logger.Warn("Certificado próximo a expirar: %s", certInfo)

// Error
logger.Error("Fallo en firma de documento: %v", err)
```

### Logging Especializado

#### Operaciones XML
```go
// Logging de operación XML exitosa
logger.LogXMLOperation("ValidarDocumento", xmlData, nil)

// Logging de error en operación XML
logger.LogXMLOperation("FirmarDocumento", xmlData, err)
```

#### Operaciones con Certificados
```go
// Logging de operación con certificado exitosa
logger.LogCertOperation("CargarCertificado", "CN=MiCertificado", nil)

// Logging de error en operación con certificado
logger.LogCertOperation("ValidarCertificado", certInfo, err)
```

## Configuración

### Estructura de Directorios
```
/logs
  ├── firma_service.log    # Logs del servicio de firma
  ├── xml_processor.log    # Logs del procesador XML
  └── app.log             # Logs generales
```

### Rotación de Logs
Se recomienda implementar una política de rotación de logs:
- Rotar logs diariamente
- Mantener logs por 30 días
- Comprimir logs antiguos

Ejemplo usando logrotate:
```conf
/path/to/logs/*.log {
    daily
    rotate 30
    compress
    delaycompress
    missingok
    notifempty
    create 0644 user group
}
```

## Buenas Prácticas

### 1. Nivel de Log Apropiado
- Desarrollo: DEBUG
- Testing: INFO/DEBUG
- Producción: INFO/WARN
- Monitoreo: ERROR

### 2. Mensajes Efectivos
- Incluir identificadores únicos
- Mantener mensajes concisos pero informativos
- Incluir datos relevantes para debugging
- Evitar información sensible

### 3. Gestión de Recursos
- Cerrar los loggers apropiadamente
- Monitorear el espacio en disco
- Implementar rotación de logs
- Limpiar logs antiguos

### 4. Seguridad
- No registrar datos sensibles (contraseñas, tokens)
- Limitar acceso a archivos de log
- Sanitizar entrada de usuario antes de registrar
- Mantener logs en ubicación segura

## Troubleshooting

### Problemas Comunes

1. **Archivos de Log Grandes**
   - Implementar rotación de logs
   - Ajustar nivel de logging
   - Revisar frecuencia de mensajes DEBUG

2. **Rendimiento**
   - Reducir logging en producción
   - Usar buffering apropiado
   - Implementar logging asíncrono

3. **Pérdida de Logs**
   - Verificar permisos de escritura
   - Monitorear espacio en disco
   - Implementar logging redundante

### Monitoreo

#### Métricas a Observar
- Tamaño de archivos de log
- Frecuencia de mensajes ERROR
- Latencia de operaciones logging
- Uso de disco

#### Alertas Recomendadas
- Errores críticos
- Espacio en disco bajo
- Fallos en rotación de logs
- Picos en frecuencia de errores

## Integración con Otras Herramientas

### ELK Stack
```go
// Configurar formato compatible con ELK
logger.SetFormat(logger.JSONFormat)
```

### Prometheus/Grafana
```go
// Métricas de logging
metrics.LogErrors.Inc()
metrics.LogSize.Set(float64(logSize))
```

### Alertmanager
```go
// Enviar alerta crítica
if errorCount > threshold {
    alertmanager.SendAlert("ErrorRateHigh", errorCount)
}
``` 