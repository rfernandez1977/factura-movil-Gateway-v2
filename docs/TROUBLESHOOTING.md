# Guía de Troubleshooting - FMgo

## Problemas Comunes y Soluciones

### 1. Errores de Firma Digital

#### 1.1 Error: "Certificado no válido"
**Síntomas:**
- Error al firmar documentos
- Mensaje de certificado inválido
- Rechazo de documentos por el SII

**Soluciones:**
1. Verificar fecha de vencimiento del certificado
2. Comprobar que el certificado sea de firma electrónica avanzada
3. Validar que el RUT del certificado coincida con el configurado
4. Revisar permisos de acceso al archivo del certificado

#### 1.2 Error: "Firma inválida"
**Síntomas:**
- Documento se firma pero es rechazado
- Error de validación de firma
- Inconsistencia en el XML firmado

**Soluciones:**
1. Verificar formato del XML antes de firmar
2. Comprobar que el certificado tenga los permisos correctos
3. Validar que la llave privada corresponda al certificado
4. Revisar el algoritmo de firma utilizado

### 2. Problemas de Conexión con SII

#### 2.1 Error: "Timeout en conexión"
**Síntomas:**
- Timeout al enviar documentos
- Conexión interrumpida
- Sin respuesta del servidor

**Soluciones:**
1. Verificar conectividad a internet
2. Comprobar estado de servicios SII
3. Ajustar timeouts en la configuración
4. Implementar reintentos automáticos

#### 2.2 Error: "Token inválido"
**Síntomas:**
- Rechazo de autenticación
- Error 401 en llamadas API
- Token expirado

**Soluciones:**
1. Renovar token de autenticación
2. Verificar credenciales configuradas
3. Sincronizar reloj del servidor
4. Limpiar caché de tokens

### 3. Errores de Validación XML

#### 3.1 Error: "Schema inválido"
**Síntomas:**
- Rechazo de documentos por formato
- Errores de validación XSD
- Inconsistencia en estructura

**Soluciones:**
1. Verificar versión del schema
2. Validar estructura del documento
3. Comprobar namespace correcto
4. Actualizar schemas desde SII

#### 3.2 Error: "Datos obligatorios faltantes"
**Síntomas:**
- Campos requeridos no presentes
- Validación fallida
- Rechazo de documento

**Soluciones:**
1. Revisar documentación de campos obligatorios
2. Implementar validaciones previas
3. Verificar mapeo de datos
4. Actualizar plantillas de documentos

### 4. Problemas de Performance

#### 4.1 Error: "Lentitud en procesamiento"
**Síntomas:**
- Tiempos de respuesta altos
- Consumo excesivo de recursos
- Timeout en operaciones

**Soluciones:**
1. Optimizar consultas a base de datos
2. Implementar/ajustar caché
3. Revisar uso de memoria
4. Monitorear métricas de performance

#### 4.2 Error: "Memoria insuficiente"
**Síntomas:**
- OutOfMemoryError
- Crash de aplicación
- Degradación de performance

**Soluciones:**
1. Ajustar configuración de memoria
2. Implementar limpieza de recursos
3. Optimizar procesamiento por lotes
4. Monitorear uso de memoria

### 5. Errores de Logging

#### 5.1 Error: "Logs no generados"
**Síntomas:**
- Ausencia de logs
- Información incompleta
- Errores de permisos

**Soluciones:**
1. Verificar configuración de logging
2. Comprobar permisos de escritura
3. Validar rotación de logs
4. Ajustar niveles de log

#### 5.2 Error: "Logs corruptos"
**Síntomas:**
- Archivos de log dañados
- Información ilegible
- Errores de formato

**Soluciones:**
1. Implementar backup de logs
2. Verificar rotación correcta
3. Ajustar formato de logs
4. Monitorear espacio en disco

### 6. Problemas de Certificados

#### 6.1 Error: "Certificado expirado"
**Síntomas:**
- Rechazo de operaciones
- Error de validación
- Certificado vencido

**Soluciones:**
1. Renovar certificado digital
2. Actualizar configuración
3. Implementar alertas de vencimiento
4. Mantener backup de certificados

#### 6.2 Error: "Certificado no encontrado"
**Síntomas:**
- Error al cargar certificado
- Archivo no encontrado
- Permisos incorrectos

**Soluciones:**
1. Verificar ruta del certificado
2. Comprobar permisos de archivo
3. Restaurar desde backup
4. Actualizar configuración

### 7. Errores de Base de Datos

#### 7.1 Error: "Conexión perdida"
**Síntomas:**
- Timeout en consultas
- Conexión rechazada
- Errores de comunicación

**Soluciones:**
1. Verificar configuración de conexión
2. Comprobar estado del servidor
3. Implementar reconexión automática
4. Monitorear pool de conexiones

#### 7.2 Error: "Deadlock detectado"
**Síntomas:**
- Bloqueo en transacciones
- Timeout en operaciones
- Errores de concurrencia

**Soluciones:**
1. Optimizar queries
2. Ajustar timeouts
3. Implementar reintentos
4. Monitorear locks

## Herramientas de Diagnóstico

### 1. Logs
```bash
# Ver últimos errores
tail -f logs/error.log

# Buscar errores específicos
grep "ERROR" logs/application.log

# Analizar logs de firma
grep "FirmaService" logs/debug.log
```

### 2. Monitoreo
```bash
# Verificar estado de servicios
systemctl status fmgo

# Monitorear recursos
top -p $(pgrep fmgo)

# Ver conexiones activas
netstat -an | grep 8080
```

### 3. Validación
```bash
# Validar XML
xmllint --schema schema.xsd documento.xml

# Verificar firma
openssl verify -CAfile ca.crt certificado.crt

# Probar conexión SII
curl -v https://palena.sii.cl/DTEWS/
```

## Contacto de Soporte

Para problemas no resueltos:
- Email: soporte@fmgo.cl
- Teléfono: +56 2 2123 4567
- Portal: https://soporte.fmgo.cl

## Referencias

- [Documentación SII](https://www.sii.cl/factura_electronica/)
- [Estándares XML-DSIG](https://www.w3.org/TR/xmldsig-core/)
- [Documentación Go](https://golang.org/doc/)
- [Documentación FMgo](https://docs.fmgo.cl) 